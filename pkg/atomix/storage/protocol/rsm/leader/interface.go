// Code generated by atomix-go-framework. DO NOT EDIT.
package leader

import (
	leader "github.com/atomix/atomix-api/go/atomix/primitive/leader"
	errors "github.com/atomix/atomix-go-framework/pkg/atomix/errors"
	rsm "github.com/atomix/atomix-go-framework/pkg/atomix/storage/protocol/rsm"
	util "github.com/atomix/atomix-go-framework/pkg/atomix/util"
	proto "github.com/golang/protobuf/proto"
	uuid "github.com/google/uuid"
	"io"
)

type Service interface {
	ServiceContext
	Backup(SnapshotWriter) error
	Restore(SnapshotReader) error
	// Latch attempts to acquire the leader latch
	Latch(LatchProposal) error
	// Get gets the current leader
	Get(GetProposal) error
	// Events listens for leader change events
	Events(EventsProposal) error
}

type ServiceContext interface {
	Scheduler() rsm.Scheduler
	Sessions() Sessions
	Proposals() Proposals
}

func newServiceContext(scheduler rsm.Scheduler) ServiceContext {
	return &serviceContext{
		scheduler: scheduler,
		sessions:  newSessions(),
		proposals: newProposals(),
	}
}

type serviceContext struct {
	scheduler rsm.Scheduler
	sessions  Sessions
	proposals Proposals
}

func (s *serviceContext) Scheduler() rsm.Scheduler {
	return s.scheduler
}

func (s *serviceContext) Sessions() Sessions {
	return s.sessions
}

func (s *serviceContext) Proposals() Proposals {
	return s.proposals
}

var _ ServiceContext = &serviceContext{}

type SnapshotWriter interface {
	WriteState(*LeaderLatchState) error
}

func newSnapshotWriter(writer io.Writer) SnapshotWriter {
	return &serviceSnapshotWriter{
		writer: writer,
	}
}

type serviceSnapshotWriter struct {
	writer io.Writer
}

func (w *serviceSnapshotWriter) WriteState(state *LeaderLatchState) error {
	bytes, err := proto.Marshal(state)
	if err != nil {
		return err
	}
	err = util.WriteBytes(w.writer, bytes)
	if err != nil {
		return err
	}
	return err
}

var _ SnapshotWriter = &serviceSnapshotWriter{}

type SnapshotReader interface {
	ReadState() (*LeaderLatchState, error)
}

func newSnapshotReader(reader io.Reader) SnapshotReader {
	return &serviceSnapshotReader{
		reader: reader,
	}
}

type serviceSnapshotReader struct {
	reader io.Reader
}

func (r *serviceSnapshotReader) ReadState() (*LeaderLatchState, error) {
	bytes, err := util.ReadBytes(r.reader)
	if err != nil {
		return nil, err
	}
	state := &LeaderLatchState{}
	err = proto.Unmarshal(bytes, state)
	if err != nil {
		return nil, err
	}
	return state, nil
}

var _ SnapshotReader = &serviceSnapshotReader{}

type Sessions interface {
	open(Session)
	expire(SessionID)
	close(SessionID)
	Get(SessionID) (Session, bool)
	List() []Session
}

func newSessions() Sessions {
	return &serviceSessions{
		sessions: make(map[SessionID]Session),
	}
}

type serviceSessions struct {
	sessions map[SessionID]Session
}

func (s *serviceSessions) open(session Session) {
	s.sessions[session.ID()] = session
	session.setState(SessionOpen)
}

func (s *serviceSessions) expire(sessionID SessionID) {
	session, ok := s.sessions[sessionID]
	if ok {
		session.setState(SessionClosed)
		delete(s.sessions, sessionID)
	}
}

func (s *serviceSessions) close(sessionID SessionID) {
	session, ok := s.sessions[sessionID]
	if ok {
		session.setState(SessionClosed)
		delete(s.sessions, sessionID)
	}
}

func (s *serviceSessions) Get(id SessionID) (Session, bool) {
	session, ok := s.sessions[id]
	return session, ok
}

func (s *serviceSessions) List() []Session {
	sessions := make([]Session, 0, len(s.sessions))
	for _, session := range s.sessions {
		sessions = append(sessions, session)
	}
	return sessions
}

var _ Sessions = &serviceSessions{}

type SessionID uint64

type SessionState int

const (
	SessionClosed SessionState = iota
	SessionOpen
)

type Watcher interface {
	Cancel()
}

func newWatcher(f func()) Watcher {
	return &serviceWatcher{
		f: f,
	}
}

type serviceWatcher struct {
	f func()
}

func (s *serviceWatcher) Cancel() {
	s.f()
}

var _ Watcher = &serviceWatcher{}

type Session interface {
	ID() SessionID
	State() SessionState
	setState(SessionState)
	Watch(func(SessionState)) Watcher
	Proposals() Proposals
}

func newSession(session rsm.Session) Session {
	return &serviceSession{
		session:   session,
		proposals: newProposals(),
		watchers:  make(map[string]func(SessionState)),
	}
}

type serviceSession struct {
	session   rsm.Session
	proposals Proposals
	state     SessionState
	watchers  map[string]func(SessionState)
}

func (s *serviceSession) ID() SessionID {
	return SessionID(s.session.ID())
}

func (s *serviceSession) Proposals() Proposals {
	return s.proposals
}

func (s *serviceSession) State() SessionState {
	return s.state
}

func (s *serviceSession) setState(state SessionState) {
	if state != s.state {
		s.state = state
		for _, watcher := range s.watchers {
			watcher(state)
		}
	}
}

func (s *serviceSession) Watch(f func(SessionState)) Watcher {
	id := uuid.New().String()
	s.watchers[id] = f
	return newWatcher(func() {
		delete(s.watchers, id)
	})
}

var _ Session = &serviceSession{}

type Proposals interface {
	Latch() LatchProposals
	Get() GetProposals
	Events() EventsProposals
}

func newProposals() Proposals {
	return &serviceProposals{
		latchProposals:  newLatchProposals(),
		getProposals:    newGetProposals(),
		eventsProposals: newEventsProposals(),
	}
}

type serviceProposals struct {
	latchProposals  LatchProposals
	getProposals    GetProposals
	eventsProposals EventsProposals
}

func (s *serviceProposals) Latch() LatchProposals {
	return s.latchProposals
}
func (s *serviceProposals) Get() GetProposals {
	return s.getProposals
}
func (s *serviceProposals) Events() EventsProposals {
	return s.eventsProposals
}

var _ Proposals = &serviceProposals{}

type ProposalID uint64

type Proposal interface {
	ID() ProposalID
	Session() Session
}

func newProposal(id ProposalID, session Session) Proposal {
	return &serviceProposal{
		id:      id,
		session: session,
	}
}

type serviceProposal struct {
	id      ProposalID
	session Session
}

func (p *serviceProposal) ID() ProposalID {
	return p.id
}

func (p *serviceProposal) Session() Session {
	return p.session
}

var _ Proposal = &serviceProposal{}

type LatchProposals interface {
	register(LatchProposal)
	unregister(ProposalID)
	Get(ProposalID) (LatchProposal, bool)
	List() []LatchProposal
}

func newLatchProposals() LatchProposals {
	return &latchProposals{
		proposals: make(map[ProposalID]LatchProposal),
	}
}

type latchProposals struct {
	proposals map[ProposalID]LatchProposal
}

func (p *latchProposals) register(proposal LatchProposal) {
	p.proposals[proposal.ID()] = proposal
}

func (p *latchProposals) unregister(id ProposalID) {
	delete(p.proposals, id)
}

func (p *latchProposals) Get(id ProposalID) (LatchProposal, bool) {
	proposal, ok := p.proposals[id]
	return proposal, ok
}

func (p *latchProposals) List() []LatchProposal {
	proposals := make([]LatchProposal, 0, len(p.proposals))
	for _, proposal := range p.proposals {
		proposals = append(proposals, proposal)
	}
	return proposals
}

var _ LatchProposals = &latchProposals{}

type LatchProposal interface {
	Proposal
	Request() *leader.LatchRequest
	Reply(*leader.LatchResponse) error
	response() *leader.LatchResponse
}

func newLatchProposal(id ProposalID, session Session, request *leader.LatchRequest) LatchProposal {
	return &latchProposal{
		Proposal: newProposal(id, session),
		req:      request,
	}
}

type latchProposal struct {
	Proposal
	req *leader.LatchRequest
	res *leader.LatchResponse
}

func (p *latchProposal) Request() *leader.LatchRequest {
	return p.req
}

func (p *latchProposal) Reply(reply *leader.LatchResponse) error {
	if p.res != nil {
		return errors.NewConflict("reply already sent")
	}
	p.res = reply
	return nil
}

func (p *latchProposal) response() *leader.LatchResponse {
	return p.res
}

var _ LatchProposal = &latchProposal{}

type GetProposals interface {
	register(GetProposal)
	unregister(ProposalID)
	Get(ProposalID) (GetProposal, bool)
	List() []GetProposal
}

func newGetProposals() GetProposals {
	return &getProposals{
		proposals: make(map[ProposalID]GetProposal),
	}
}

type getProposals struct {
	proposals map[ProposalID]GetProposal
}

func (p *getProposals) register(proposal GetProposal) {
	p.proposals[proposal.ID()] = proposal
}

func (p *getProposals) unregister(id ProposalID) {
	delete(p.proposals, id)
}

func (p *getProposals) Get(id ProposalID) (GetProposal, bool) {
	proposal, ok := p.proposals[id]
	return proposal, ok
}

func (p *getProposals) List() []GetProposal {
	proposals := make([]GetProposal, 0, len(p.proposals))
	for _, proposal := range p.proposals {
		proposals = append(proposals, proposal)
	}
	return proposals
}

var _ GetProposals = &getProposals{}

type GetProposal interface {
	Proposal
	Request() *leader.GetRequest
	Reply(*leader.GetResponse) error
	response() *leader.GetResponse
}

func newGetProposal(id ProposalID, session Session, request *leader.GetRequest) GetProposal {
	return &getProposal{
		Proposal: newProposal(id, session),
		req:      request,
	}
}

type getProposal struct {
	Proposal
	req *leader.GetRequest
	res *leader.GetResponse
}

func (p *getProposal) Request() *leader.GetRequest {
	return p.req
}

func (p *getProposal) Reply(reply *leader.GetResponse) error {
	if p.res != nil {
		return errors.NewConflict("reply already sent")
	}
	p.res = reply
	return nil
}

func (p *getProposal) response() *leader.GetResponse {
	return p.res
}

var _ GetProposal = &getProposal{}

type EventsProposals interface {
	register(EventsProposal)
	unregister(ProposalID)
	Get(ProposalID) (EventsProposal, bool)
	List() []EventsProposal
}

func newEventsProposals() EventsProposals {
	return &eventsProposals{
		proposals: make(map[ProposalID]EventsProposal),
	}
}

type eventsProposals struct {
	proposals map[ProposalID]EventsProposal
}

func (p *eventsProposals) register(proposal EventsProposal) {
	p.proposals[proposal.ID()] = proposal
}

func (p *eventsProposals) unregister(id ProposalID) {
	delete(p.proposals, id)
}

func (p *eventsProposals) Get(id ProposalID) (EventsProposal, bool) {
	proposal, ok := p.proposals[id]
	return proposal, ok
}

func (p *eventsProposals) List() []EventsProposal {
	proposals := make([]EventsProposal, 0, len(p.proposals))
	for _, proposal := range p.proposals {
		proposals = append(proposals, proposal)
	}
	return proposals
}

var _ EventsProposals = &eventsProposals{}

type EventsProposal interface {
	Proposal
	Request() *leader.EventsRequest
	Notify(*leader.EventsResponse) error
	Close() error
}

func newEventsProposal(id ProposalID, session Session, request *leader.EventsRequest, stream rsm.Stream) EventsProposal {
	return &eventsProposal{
		Proposal: newProposal(id, session),
		request:  request,
		stream:   stream,
	}
}

type eventsProposal struct {
	Proposal
	request *leader.EventsRequest
	stream  rsm.Stream
}

func (p *eventsProposal) Request() *leader.EventsRequest {
	return p.request
}

func (p *eventsProposal) Notify(notification *leader.EventsResponse) error {
	bytes, err := proto.Marshal(notification)
	if err != nil {
		return err
	}
	p.stream.Value(bytes)
	return nil
}

func (p *eventsProposal) Close() error {
	p.stream.Close()
	return nil
}

var _ EventsProposal = &eventsProposal{}
