// Code generated by atomix-go-framework. DO NOT EDIT.
package _map

import (
	_map "github.com/atomix/atomix-api/go/atomix/primitive/map"
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
	// Size returns the size of the map
	Size(SizeProposal) error
	// Put puts an entry into the map
	Put(PutProposal) error
	// Get gets the entry for a key
	Get(GetProposal) error
	// Remove removes an entry from the map
	Remove(RemoveProposal) error
	// Clear removes all entries from the map
	Clear(ClearProposal) error
	// Events listens for change events
	Events(EventsProposal) error
	// Entries lists all entries in the map
	Entries(EntriesProposal) error
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
	WriteState(*MapState) error
}

func newSnapshotWriter(writer io.Writer) SnapshotWriter {
	return &serviceSnapshotWriter{
		writer: writer,
	}
}

type serviceSnapshotWriter struct {
	writer io.Writer
}

func (w *serviceSnapshotWriter) WriteState(state *MapState) error {
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
	ReadState() (*MapState, error)
}

func newSnapshotReader(reader io.Reader) SnapshotReader {
	return &serviceSnapshotReader{
		reader: reader,
	}
}

type serviceSnapshotReader struct {
	reader io.Reader
}

func (r *serviceSnapshotReader) ReadState() (*MapState, error) {
	bytes, err := util.ReadBytes(r.reader)
	if err != nil {
		return nil, err
	}
	state := &MapState{}
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
	Size() SizeProposals
	Put() PutProposals
	Get() GetProposals
	Remove() RemoveProposals
	Clear() ClearProposals
	Events() EventsProposals
	Entries() EntriesProposals
}

func newProposals() Proposals {
	return &serviceProposals{
		sizeProposals:    newSizeProposals(),
		putProposals:     newPutProposals(),
		getProposals:     newGetProposals(),
		removeProposals:  newRemoveProposals(),
		clearProposals:   newClearProposals(),
		eventsProposals:  newEventsProposals(),
		entriesProposals: newEntriesProposals(),
	}
}

type serviceProposals struct {
	sizeProposals    SizeProposals
	putProposals     PutProposals
	getProposals     GetProposals
	removeProposals  RemoveProposals
	clearProposals   ClearProposals
	eventsProposals  EventsProposals
	entriesProposals EntriesProposals
}

func (s *serviceProposals) Size() SizeProposals {
	return s.sizeProposals
}
func (s *serviceProposals) Put() PutProposals {
	return s.putProposals
}
func (s *serviceProposals) Get() GetProposals {
	return s.getProposals
}
func (s *serviceProposals) Remove() RemoveProposals {
	return s.removeProposals
}
func (s *serviceProposals) Clear() ClearProposals {
	return s.clearProposals
}
func (s *serviceProposals) Events() EventsProposals {
	return s.eventsProposals
}
func (s *serviceProposals) Entries() EntriesProposals {
	return s.entriesProposals
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

type SizeProposals interface {
	register(SizeProposal)
	unregister(ProposalID)
	Get(ProposalID) (SizeProposal, bool)
	List() []SizeProposal
}

func newSizeProposals() SizeProposals {
	return &sizeProposals{
		proposals: make(map[ProposalID]SizeProposal),
	}
}

type sizeProposals struct {
	proposals map[ProposalID]SizeProposal
}

func (p *sizeProposals) register(proposal SizeProposal) {
	p.proposals[proposal.ID()] = proposal
}

func (p *sizeProposals) unregister(id ProposalID) {
	delete(p.proposals, id)
}

func (p *sizeProposals) Get(id ProposalID) (SizeProposal, bool) {
	proposal, ok := p.proposals[id]
	return proposal, ok
}

func (p *sizeProposals) List() []SizeProposal {
	proposals := make([]SizeProposal, 0, len(p.proposals))
	for _, proposal := range p.proposals {
		proposals = append(proposals, proposal)
	}
	return proposals
}

var _ SizeProposals = &sizeProposals{}

type SizeProposal interface {
	Proposal
	Request() *_map.SizeRequest
	Reply(*_map.SizeResponse) error
	response() *_map.SizeResponse
}

func newSizeProposal(id ProposalID, session Session, request *_map.SizeRequest) SizeProposal {
	return &sizeProposal{
		Proposal: newProposal(id, session),
		req:      request,
	}
}

type sizeProposal struct {
	Proposal
	req *_map.SizeRequest
	res *_map.SizeResponse
}

func (p *sizeProposal) Request() *_map.SizeRequest {
	return p.req
}

func (p *sizeProposal) Reply(reply *_map.SizeResponse) error {
	if p.res != nil {
		return errors.NewConflict("reply already sent")
	}
	p.res = reply
	return nil
}

func (p *sizeProposal) response() *_map.SizeResponse {
	return p.res
}

var _ SizeProposal = &sizeProposal{}

type PutProposals interface {
	register(PutProposal)
	unregister(ProposalID)
	Get(ProposalID) (PutProposal, bool)
	List() []PutProposal
}

func newPutProposals() PutProposals {
	return &putProposals{
		proposals: make(map[ProposalID]PutProposal),
	}
}

type putProposals struct {
	proposals map[ProposalID]PutProposal
}

func (p *putProposals) register(proposal PutProposal) {
	p.proposals[proposal.ID()] = proposal
}

func (p *putProposals) unregister(id ProposalID) {
	delete(p.proposals, id)
}

func (p *putProposals) Get(id ProposalID) (PutProposal, bool) {
	proposal, ok := p.proposals[id]
	return proposal, ok
}

func (p *putProposals) List() []PutProposal {
	proposals := make([]PutProposal, 0, len(p.proposals))
	for _, proposal := range p.proposals {
		proposals = append(proposals, proposal)
	}
	return proposals
}

var _ PutProposals = &putProposals{}

type PutProposal interface {
	Proposal
	Request() *_map.PutRequest
	Reply(*_map.PutResponse) error
	response() *_map.PutResponse
}

func newPutProposal(id ProposalID, session Session, request *_map.PutRequest) PutProposal {
	return &putProposal{
		Proposal: newProposal(id, session),
		req:      request,
	}
}

type putProposal struct {
	Proposal
	req *_map.PutRequest
	res *_map.PutResponse
}

func (p *putProposal) Request() *_map.PutRequest {
	return p.req
}

func (p *putProposal) Reply(reply *_map.PutResponse) error {
	if p.res != nil {
		return errors.NewConflict("reply already sent")
	}
	p.res = reply
	return nil
}

func (p *putProposal) response() *_map.PutResponse {
	return p.res
}

var _ PutProposal = &putProposal{}

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
	Request() *_map.GetRequest
	Reply(*_map.GetResponse) error
	response() *_map.GetResponse
}

func newGetProposal(id ProposalID, session Session, request *_map.GetRequest) GetProposal {
	return &getProposal{
		Proposal: newProposal(id, session),
		req:      request,
	}
}

type getProposal struct {
	Proposal
	req *_map.GetRequest
	res *_map.GetResponse
}

func (p *getProposal) Request() *_map.GetRequest {
	return p.req
}

func (p *getProposal) Reply(reply *_map.GetResponse) error {
	if p.res != nil {
		return errors.NewConflict("reply already sent")
	}
	p.res = reply
	return nil
}

func (p *getProposal) response() *_map.GetResponse {
	return p.res
}

var _ GetProposal = &getProposal{}

type RemoveProposals interface {
	register(RemoveProposal)
	unregister(ProposalID)
	Get(ProposalID) (RemoveProposal, bool)
	List() []RemoveProposal
}

func newRemoveProposals() RemoveProposals {
	return &removeProposals{
		proposals: make(map[ProposalID]RemoveProposal),
	}
}

type removeProposals struct {
	proposals map[ProposalID]RemoveProposal
}

func (p *removeProposals) register(proposal RemoveProposal) {
	p.proposals[proposal.ID()] = proposal
}

func (p *removeProposals) unregister(id ProposalID) {
	delete(p.proposals, id)
}

func (p *removeProposals) Get(id ProposalID) (RemoveProposal, bool) {
	proposal, ok := p.proposals[id]
	return proposal, ok
}

func (p *removeProposals) List() []RemoveProposal {
	proposals := make([]RemoveProposal, 0, len(p.proposals))
	for _, proposal := range p.proposals {
		proposals = append(proposals, proposal)
	}
	return proposals
}

var _ RemoveProposals = &removeProposals{}

type RemoveProposal interface {
	Proposal
	Request() *_map.RemoveRequest
	Reply(*_map.RemoveResponse) error
	response() *_map.RemoveResponse
}

func newRemoveProposal(id ProposalID, session Session, request *_map.RemoveRequest) RemoveProposal {
	return &removeProposal{
		Proposal: newProposal(id, session),
		req:      request,
	}
}

type removeProposal struct {
	Proposal
	req *_map.RemoveRequest
	res *_map.RemoveResponse
}

func (p *removeProposal) Request() *_map.RemoveRequest {
	return p.req
}

func (p *removeProposal) Reply(reply *_map.RemoveResponse) error {
	if p.res != nil {
		return errors.NewConflict("reply already sent")
	}
	p.res = reply
	return nil
}

func (p *removeProposal) response() *_map.RemoveResponse {
	return p.res
}

var _ RemoveProposal = &removeProposal{}

type ClearProposals interface {
	register(ClearProposal)
	unregister(ProposalID)
	Get(ProposalID) (ClearProposal, bool)
	List() []ClearProposal
}

func newClearProposals() ClearProposals {
	return &clearProposals{
		proposals: make(map[ProposalID]ClearProposal),
	}
}

type clearProposals struct {
	proposals map[ProposalID]ClearProposal
}

func (p *clearProposals) register(proposal ClearProposal) {
	p.proposals[proposal.ID()] = proposal
}

func (p *clearProposals) unregister(id ProposalID) {
	delete(p.proposals, id)
}

func (p *clearProposals) Get(id ProposalID) (ClearProposal, bool) {
	proposal, ok := p.proposals[id]
	return proposal, ok
}

func (p *clearProposals) List() []ClearProposal {
	proposals := make([]ClearProposal, 0, len(p.proposals))
	for _, proposal := range p.proposals {
		proposals = append(proposals, proposal)
	}
	return proposals
}

var _ ClearProposals = &clearProposals{}

type ClearProposal interface {
	Proposal
	Request() *_map.ClearRequest
	Reply(*_map.ClearResponse) error
	response() *_map.ClearResponse
}

func newClearProposal(id ProposalID, session Session, request *_map.ClearRequest) ClearProposal {
	return &clearProposal{
		Proposal: newProposal(id, session),
		req:      request,
	}
}

type clearProposal struct {
	Proposal
	req *_map.ClearRequest
	res *_map.ClearResponse
}

func (p *clearProposal) Request() *_map.ClearRequest {
	return p.req
}

func (p *clearProposal) Reply(reply *_map.ClearResponse) error {
	if p.res != nil {
		return errors.NewConflict("reply already sent")
	}
	p.res = reply
	return nil
}

func (p *clearProposal) response() *_map.ClearResponse {
	return p.res
}

var _ ClearProposal = &clearProposal{}

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
	Request() *_map.EventsRequest
	Notify(*_map.EventsResponse) error
	Close() error
}

func newEventsProposal(id ProposalID, session Session, request *_map.EventsRequest, stream rsm.Stream) EventsProposal {
	return &eventsProposal{
		Proposal: newProposal(id, session),
		request:  request,
		stream:   stream,
	}
}

type eventsProposal struct {
	Proposal
	request *_map.EventsRequest
	stream  rsm.Stream
}

func (p *eventsProposal) Request() *_map.EventsRequest {
	return p.request
}

func (p *eventsProposal) Notify(notification *_map.EventsResponse) error {
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

type EntriesProposals interface {
	register(EntriesProposal)
	unregister(ProposalID)
	Get(ProposalID) (EntriesProposal, bool)
	List() []EntriesProposal
}

func newEntriesProposals() EntriesProposals {
	return &entriesProposals{
		proposals: make(map[ProposalID]EntriesProposal),
	}
}

type entriesProposals struct {
	proposals map[ProposalID]EntriesProposal
}

func (p *entriesProposals) register(proposal EntriesProposal) {
	p.proposals[proposal.ID()] = proposal
}

func (p *entriesProposals) unregister(id ProposalID) {
	delete(p.proposals, id)
}

func (p *entriesProposals) Get(id ProposalID) (EntriesProposal, bool) {
	proposal, ok := p.proposals[id]
	return proposal, ok
}

func (p *entriesProposals) List() []EntriesProposal {
	proposals := make([]EntriesProposal, 0, len(p.proposals))
	for _, proposal := range p.proposals {
		proposals = append(proposals, proposal)
	}
	return proposals
}

var _ EntriesProposals = &entriesProposals{}

type EntriesProposal interface {
	Proposal
	Request() *_map.EntriesRequest
	Notify(*_map.EntriesResponse) error
	Close() error
}

func newEntriesProposal(id ProposalID, session Session, request *_map.EntriesRequest, stream rsm.Stream) EntriesProposal {
	return &entriesProposal{
		Proposal: newProposal(id, session),
		request:  request,
		stream:   stream,
	}
}

type entriesProposal struct {
	Proposal
	request *_map.EntriesRequest
	stream  rsm.Stream
}

func (p *entriesProposal) Request() *_map.EntriesRequest {
	return p.request
}

func (p *entriesProposal) Notify(notification *_map.EntriesResponse) error {
	bytes, err := proto.Marshal(notification)
	if err != nil {
		return err
	}
	p.stream.Value(bytes)
	return nil
}

func (p *entriesProposal) Close() error {
	p.stream.Close()
	return nil
}

var _ EntriesProposal = &entriesProposal{}
