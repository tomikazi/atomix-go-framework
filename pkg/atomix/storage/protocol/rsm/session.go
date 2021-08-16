// Copyright 2019-present Open Networking Foundation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rsm

import (
	"container/list"
	"encoding/binary"
	"github.com/atomix/atomix-go-framework/pkg/atomix/errors"
	streams "github.com/atomix/atomix-go-framework/pkg/atomix/stream"
	"github.com/bits-and-blooms/bloom/v3"
	"github.com/google/uuid"
	"time"
)

type SessionState int

const (
	SessionClosed SessionState = iota
	SessionOpen
)

// SessionID is a session identifier
type SessionID uint64

// Sessions provides access to open sessions
type Sessions interface {
	// Get gets a session by ID
	Get(SessionID) (Session, bool)
	// List lists all open sessions
	List() []Session
}

func newServiceSessions() *primitiveServiceSessions {
	return &primitiveServiceSessions{
		sessions: make(map[SessionID]*primitiveServiceSession),
	}
}

type primitiveServiceSessions struct {
	sessions map[SessionID]*primitiveServiceSession
}

func (s *primitiveServiceSessions) add(session *primitiveServiceSession) {
	s.sessions[session.ID()] = session
}

func (s *primitiveServiceSessions) remove(session *primitiveServiceSession) {
	delete(s.sessions, session.ID())
}

func (s *primitiveServiceSessions) Get(sessionID SessionID) (Session, bool) {
	session, ok := s.sessions[sessionID]
	return session, ok
}

func (s *primitiveServiceSessions) List() []Session {
	sessions := make([]Session, 0, len(s.sessions))
	for _, session := range s.sessions {
		sessions = append(sessions, session)
	}
	return sessions
}

var _ Sessions = (*primitiveServiceSessions)(nil)

// Session is a service session
type Session interface {
	// ID returns the session identifier
	ID() SessionID
	// State returns the current session state
	State() SessionState
	// Watch watches the session state
	Watch(f func(SessionState)) Watcher
	// Commands returns the session commands
	Commands() Commands
}

// Watcher is a context for a Watch call
type Watcher interface {
	// Cancel cancels the watcher
	Cancel()
}

func newSession(manager *primitiveServiceManager) *primitiveSession {
	return &primitiveSession{
		manager:  manager,
		services: make(map[ServiceID]*primitiveServiceSession),
	}
}

// primitiveSession is a Session implementation
type primitiveSession struct {
	manager     *primitiveServiceManager
	sessionID   SessionID
	timeout     time.Duration
	lastUpdated time.Time
	services    map[ServiceID]*primitiveServiceSession
}

func (s *primitiveSession) getService(serviceID ServiceID) (*primitiveServiceSession, bool) {
	session, ok := s.services[serviceID]
	return session, ok
}

func (s *primitiveSession) open(sessionID SessionID, timeout time.Duration) error {
	s.sessionID = sessionID
	s.timeout = timeout
	s.lastUpdated = s.manager.timestamp
	return nil
}

func (s *primitiveSession) snapshot() (*SessionSnapshot, error) {
	return &SessionSnapshot{
		SessionID:   s.sessionID,
		Timeout:     s.timeout,
		LastUpdated: s.lastUpdated,
	}, nil
}

func (s *primitiveSession) restore(snapshot *SessionSnapshot) error {
	s.sessionID = snapshot.SessionID
	s.timeout = snapshot.Timeout
	s.lastUpdated = snapshot.LastUpdated
	return nil
}

func (s *primitiveSession) keepAlive(lastRequestID RequestID, requestFilter *bloom.BloomFilter) error {
	for _, serviceSession := range s.services {
		if err := serviceSession.keepAlive(lastRequestID, requestFilter); err != nil {
			return err
		}
	}
	return nil
}

func (s *primitiveSession) close() error {
	return nil
}

func newServiceSession(service *primitiveService) *primitiveServiceSession {
	return &primitiveServiceSession{
		service:  service,
		watchers: make(map[string]func(SessionState)),
	}
}

type primitiveServiceSession struct {
	service  *primitiveService
	session  *primitiveSession
	commands *primitiveSessionCommands
	requests map[RequestID]*primitiveServiceSessionCommand
	state    SessionState
	watchers map[string]func(SessionState)
}

func (s *primitiveServiceSession) ID() SessionID {
	return s.session.sessionID
}

func (s *primitiveServiceSession) State() SessionState {
	return s.state
}

func (s *primitiveServiceSession) Watch(f func(SessionState)) Watcher {
	id := uuid.New().String()
	s.watchers[id] = f
	return &primitiveSessionWatcher{func() {
		delete(s.watchers, id)
	}}
}

func (s *primitiveServiceSession) Commands() Commands {
	return s.commands
}

func (s *primitiveServiceSession) command(requestID RequestID) *primitiveServiceSessionCommand {
	command, ok := s.requests[requestID]
	if ok {
		return command
	}
	return newServiceSessionCommand(s)
}

func (s *primitiveServiceSession) query() *primitiveServiceSessionQuery {
	return newServiceSessionQuery(s)
}

func (s *primitiveServiceSession) open(sessionID SessionID) error {
	session, ok := s.service.manager.sessions[sessionID]
	if !ok {
		return errors.NewInvalid("unknown session %d", sessionID)
	}
	s.session = session
	s.requests = make(map[RequestID]*primitiveServiceSessionCommand)
	s.commands = newSessionCommands()
	s.session.services[s.service.serviceID] = s
	s.service.sessions.add(s)
	s.state = SessionOpen
	return nil
}

func (s *primitiveServiceSession) snapshot() (*ServiceSessionSnapshot, error) {
	commands := make([]*SessionCommandSnapshot, 0, len(s.commands.commands))
	for _, command := range s.requests {
		commandSnapshot, err := command.snapshot()
		if err != nil {
			return nil, err
		}
		commands = append(commands, commandSnapshot)
	}
	return &ServiceSessionSnapshot{
		SessionID: s.session.sessionID,
		Commands:  commands,
	}, nil
}

func (s *primitiveServiceSession) restore(snapshot *ServiceSessionSnapshot) error {
	session, ok := s.service.manager.sessions[snapshot.SessionID]
	if !ok {
		return errors.NewInvalid("unknown session %d", snapshot.SessionID)
	}
	s.session = session
	s.requests = make(map[RequestID]*primitiveServiceSessionCommand)
	s.commands = newSessionCommands()
	for _, commandSnapshot := range snapshot.Commands {
		command := newServiceSessionCommand(s)
		if err := command.restore(commandSnapshot); err != nil {
			return err
		}
	}
	s.session.services[s.service.serviceID] = s
	s.state = SessionOpen
	s.service.sessions.add(s)
	return nil
}

func (s *primitiveServiceSession) keepAlive(lastRequestID RequestID, requestFilter *bloom.BloomFilter) error {
	for _, command := range s.commands.commands {
		if err := command.keepAlive(lastRequestID, requestFilter); err != nil {
			return err
		}
	}
	return nil
}

func (s *primitiveServiceSession) close() error {
	s.service.sessions.remove(s)
	s.state = SessionClosed
	for _, watcher := range s.watchers {
		watcher(SessionClosed)
	}
	return nil
}

var _ Session = (*primitiveServiceSession)(nil)

func newSessionCommands() *primitiveSessionCommands {
	return &primitiveSessionCommands{
		commands: make(map[CommandID]*primitiveServiceSessionCommand),
	}
}

type primitiveSessionCommands struct {
	commands map[CommandID]*primitiveServiceSessionCommand
}

func (s *primitiveSessionCommands) add(command *primitiveServiceSessionCommand) {
	s.commands[command.ID()] = command
}

func (s *primitiveSessionCommands) remove(command *primitiveServiceSessionCommand) {
	delete(s.commands, command.ID())
}

func (s *primitiveSessionCommands) Get(commandID CommandID) (Command, bool) {
	command, ok := s.commands[commandID]
	return command, ok
}

func (s *primitiveSessionCommands) List(operationID OperationID) []Command {
	commands := make([]Command, 0, len(s.commands))
	for _, command := range s.commands {
		if command.OperationID() == operationID {
			commands = append(commands, command)
		}
	}
	return commands
}

var _ Commands = (*primitiveSessionCommands)(nil)

func newServiceSessionCommand(session *primitiveServiceSession) *primitiveServiceSessionCommand {
	return &primitiveServiceSessionCommand{
		primitiveOperation: newOperation(session),
		session:            session,
		watchers:           make(map[string]func(CommandState)),
	}
}

type primitiveServiceSessionCommand struct {
	*primitiveOperation
	commandID  CommandID
	session    *primitiveServiceSession
	state      CommandState
	watchers   map[string]func(CommandState)
	request    *ServiceCommandRequest
	responses  *list.List
	responseID ResponseID
	stream     streams.WriteStream
}

func (c *primitiveServiceSessionCommand) ID() CommandID {
	return c.commandID
}

func (c *primitiveServiceSessionCommand) OperationID() OperationID {
	return c.request.Operation.OperationID
}

func (c *primitiveServiceSessionCommand) State() CommandState {
	return c.state
}

func (c *primitiveServiceSessionCommand) Watch(f func(state CommandState)) Watcher {
	id := uuid.New().String()
	c.watchers[id] = f
	return &primitiveSessionWatcher{func() {
		delete(c.watchers, id)
	}}
}

func (c *primitiveServiceSessionCommand) Input() []byte {
	return c.request.Operation.Value
}

func (c *primitiveServiceSessionCommand) execute(request *ServiceCommandRequest, stream streams.WriteStream) {
	switch c.state {
	case CommandPending:
		c.commandID = CommandID(c.session.service.Index())
		c.request = request
		c.responses = list.New()
		c.stream = stream
		c.session.requests[request.RequestID] = c
		c.session.commands.add(c)
		c.session.service.commands.add(c)
		c.state = CommandRunning
		log.Debugf("Executing command %d: %+v", c.commandID, request)
		c.session.service.service.ExecuteCommand(c)
	case CommandRunning:
		if c.responses.Len() > 0 {
			log.Debugf("Replaying %d responses for command %d: %+v", c.responses.Len(), c.commandID, request)
			elem := c.responses.Front()
			for elem != nil {
				response := elem.Value.(*ServiceCommandResponse)
				stream.Value(response)
				elem = elem.Next()
			}
		}
		c.stream = stream
	case CommandComplete:
		if c.responses.Len() > 0 {
			log.Debugf("Replaying %d responses for command %d: %+v", c.responses.Len(), c.commandID, request)
			elem := c.responses.Front()
			for elem != nil {
				response := elem.Value.(*ServiceCommandResponse)
				stream.Value(response)
				elem = elem.Next()
			}
		}
		stream.Close()
	}
}

func (c *primitiveServiceSessionCommand) snapshot() (*SessionCommandSnapshot, error) {
	responses := make([]ServiceCommandResponse, 0, c.responses.Len())
	elem := c.responses.Front()
	for elem != nil {
		responses = append(responses, *elem.Value.(*ServiceCommandResponse))
		elem = elem.Next()
	}
	var state SessionCommandState
	switch c.state {
	case CommandRunning:
		state = SessionCommandState_COMMAND_OPEN
	case CommandComplete:
		state = SessionCommandState_COMMAND_COMPLETE
	}
	return &SessionCommandSnapshot{
		CommandID:        c.commandID,
		State:            state,
		Request:          c.request,
		PendingResponses: responses,
	}, nil
}

func (c *primitiveServiceSessionCommand) restore(snapshot *SessionCommandSnapshot) error {
	c.commandID = snapshot.CommandID
	switch snapshot.State {
	case SessionCommandState_COMMAND_OPEN:
		c.state = CommandRunning
	case SessionCommandState_COMMAND_COMPLETE:
		c.state = CommandComplete
	}
	c.request = snapshot.Request
	c.responses = list.New()
	for _, response := range snapshot.PendingResponses {
		r := response
		c.responses.PushBack(&r)
	}
	c.stream = streams.NewNilStream()
	c.session.requests[c.request.RequestID] = c
	c.session.commands.add(c)
	c.session.service.commands.add(c)
	return nil
}

func (c *primitiveServiceSessionCommand) keepAlive(lastRequestID RequestID, filter *bloom.BloomFilter) error {
	if lastRequestID < c.request.RequestID {
		return nil
	}

	// If the request ID is not in the keep-alive filter, the client canceled the request.
	// Close the canceled request and remove it from the session.
	requestBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(requestBytes, uint64(c.request.RequestID))
	if c.state == CommandRunning && !filter.Test(requestBytes) {
		c.Close()
		delete(c.session.requests, c.request.RequestID)
		return nil
	}

	// The keep-alive filter indicates the next response ID the client is waiting for.
	// Remove pending responses up to the first response ID matching the keep-alive filter.
	elem := c.responses.Front()
	for elem != nil {
		response := elem.Value.(*ServiceCommandResponse)
		bytes := make([]byte, 8)
		binary.BigEndian.PutUint64(bytes, uint64(response.ResponseID))
		responseBytes := append(requestBytes, bytes...)
		if !filter.Test(responseBytes) {
			c.responses.Remove(elem)
		} else {
			break
		}
		elem = elem.Next()
	}

	// If the command is complete and the client has acknowledged receipt of all responses,
	// remove the command from the session.
	if c.state == CommandComplete && c.responses.Len() == 0 {
		delete(c.session.requests, c.request.RequestID)
	}
	return nil
}

func (c *primitiveServiceSessionCommand) Output(bytes []byte, err error) {
	if c.state == CommandComplete {
		return
	}

	c.responseID++
	response := &ServiceCommandResponse{
		ResponseID: c.responseID,
		Operation: &OperationResponse{
			Status: ResponseStatus{
				Code:    getCode(err),
				Message: getMessage(err),
			},
			Value: bytes,
		},
	}
	c.responses.PushBack(response)
	c.stream.Value(response)
}

func (c *primitiveServiceSessionCommand) Close() {
	c.session.service.commands.remove(c)
	c.session.commands.remove(c)
	c.state = CommandComplete
	for _, watcher := range c.watchers {
		watcher(CommandComplete)
	}
	c.stream.Close()
}

var _ Command = (*primitiveServiceSessionCommand)(nil)

func newServiceSessionQuery(session *primitiveServiceSession) *primitiveServiceSessionQuery {
	return &primitiveServiceSessionQuery{
		primitiveOperation: newOperation(session),
		session:            session,
	}
}

type primitiveServiceSessionQuery struct {
	*primitiveOperation
	session    *primitiveServiceSession
	request    *ServiceQueryRequest
	stream     streams.WriteStream
	responseID ResponseID
}

func (q *primitiveServiceSessionQuery) OperationID() OperationID {
	return q.request.Operation.OperationID
}

func (q *primitiveServiceSessionQuery) Input() []byte {
	return q.request.Operation.Value
}

func (q *primitiveServiceSessionQuery) execute(request *ServiceQueryRequest, stream streams.WriteStream) {
	q.request = request
	q.stream = stream
	log.Debugf("Executing query at index %d: %+v", q.session.service.Index(), request)
	q.session.service.service.ExecuteQuery(q)
}

func (q *primitiveServiceSessionQuery) Output(bytes []byte, err error) {
	q.responseID++
	response := &ServiceQueryResponse{
		ResponseID: q.responseID,
		Operation: &OperationResponse{
			Status: ResponseStatus{
				Code:    getCode(err),
				Message: getMessage(err),
			},
			Value: bytes,
		},
	}
	q.stream.Value(response)
}

func (q *primitiveServiceSessionQuery) Close() {
	q.stream.Close()
}

var _ Query = (*primitiveServiceSessionQuery)(nil)

type primitiveSessionWatcher struct {
	f func()
}

func (w *primitiveSessionWatcher) Cancel() {
	w.f()
}

var _ Watcher = (*primitiveSessionWatcher)(nil)
