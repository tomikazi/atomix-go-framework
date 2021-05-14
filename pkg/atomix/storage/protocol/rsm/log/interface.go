// Code generated by atomix-go-framework. DO NOT EDIT.
package log

import (
	log "github.com/atomix/atomix-api/go/atomix/primitive/log"
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
	// Size returns the size of the log
	Size(SizeProposal) error
	// Appends appends an entry into the log
	Append(AppendProposal) error
	// Get gets the entry for an index
	Get(GetProposal) error
	// FirstEntry gets the first entry in the log
	FirstEntry(FirstEntryProposal) error
	// LastEntry gets the last entry in the log
	LastEntry(LastEntryProposal) error
	// PrevEntry gets the previous entry in the log
	PrevEntry(PrevEntryProposal) error
	// NextEntry gets the next entry in the log
	NextEntry(NextEntryProposal) error
	// Remove removes an entry from the log
	Remove(RemoveProposal) error
	// Clear removes all entries from the log
	Clear(ClearProposal) error
	// Events listens for change events
	Events(EventsProposal) error
	// Entries lists all entries in the log
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
	WriteState(*LogState) error
}

func newSnapshotWriter(writer io.Writer) SnapshotWriter {
	return &serviceSnapshotWriter{
		writer: writer,
	}
}

type serviceSnapshotWriter struct {
	writer io.Writer
}

func (w *serviceSnapshotWriter) WriteState(state *LogState) error {
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
	ReadState() (*LogState, error)
}

func newSnapshotReader(reader io.Reader) SnapshotReader {
	return &serviceSnapshotReader{
		reader: reader,
	}
}

type serviceSnapshotReader struct {
	reader io.Reader
}

func (r *serviceSnapshotReader) ReadState() (*LogState, error) {
	bytes, err := util.ReadBytes(r.reader)
	if err != nil {
		return nil, err
	}
	state := &LogState{}
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
	Append() AppendProposals
	Get() GetProposals
	FirstEntry() FirstEntryProposals
	LastEntry() LastEntryProposals
	PrevEntry() PrevEntryProposals
	NextEntry() NextEntryProposals
	Remove() RemoveProposals
	Clear() ClearProposals
	Events() EventsProposals
	Entries() EntriesProposals
}

func newProposals() Proposals {
	return &serviceProposals{
		sizeProposals:       newSizeProposals(),
		appendProposals:     newAppendProposals(),
		getProposals:        newGetProposals(),
		firstEntryProposals: newFirstEntryProposals(),
		lastEntryProposals:  newLastEntryProposals(),
		prevEntryProposals:  newPrevEntryProposals(),
		nextEntryProposals:  newNextEntryProposals(),
		removeProposals:     newRemoveProposals(),
		clearProposals:      newClearProposals(),
		eventsProposals:     newEventsProposals(),
		entriesProposals:    newEntriesProposals(),
	}
}

type serviceProposals struct {
	sizeProposals       SizeProposals
	appendProposals     AppendProposals
	getProposals        GetProposals
	firstEntryProposals FirstEntryProposals
	lastEntryProposals  LastEntryProposals
	prevEntryProposals  PrevEntryProposals
	nextEntryProposals  NextEntryProposals
	removeProposals     RemoveProposals
	clearProposals      ClearProposals
	eventsProposals     EventsProposals
	entriesProposals    EntriesProposals
}

func (s *serviceProposals) Size() SizeProposals {
	return s.sizeProposals
}
func (s *serviceProposals) Append() AppendProposals {
	return s.appendProposals
}
func (s *serviceProposals) Get() GetProposals {
	return s.getProposals
}
func (s *serviceProposals) FirstEntry() FirstEntryProposals {
	return s.firstEntryProposals
}
func (s *serviceProposals) LastEntry() LastEntryProposals {
	return s.lastEntryProposals
}
func (s *serviceProposals) PrevEntry() PrevEntryProposals {
	return s.prevEntryProposals
}
func (s *serviceProposals) NextEntry() NextEntryProposals {
	return s.nextEntryProposals
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
	Request() *log.SizeRequest
	Reply(*log.SizeResponse) error
	response() *log.SizeResponse
}

func newSizeProposal(id ProposalID, session Session, request *log.SizeRequest) SizeProposal {
	return &sizeProposal{
		Proposal: newProposal(id, session),
		req:      request,
	}
}

type sizeProposal struct {
	Proposal
	req *log.SizeRequest
	res *log.SizeResponse
}

func (p *sizeProposal) Request() *log.SizeRequest {
	return p.req
}

func (p *sizeProposal) Reply(reply *log.SizeResponse) error {
	if p.res != nil {
		return errors.NewConflict("reply already sent")
	}
	p.res = reply
	return nil
}

func (p *sizeProposal) response() *log.SizeResponse {
	return p.res
}

var _ SizeProposal = &sizeProposal{}

type AppendProposals interface {
	register(AppendProposal)
	unregister(ProposalID)
	Get(ProposalID) (AppendProposal, bool)
	List() []AppendProposal
}

func newAppendProposals() AppendProposals {
	return &appendProposals{
		proposals: make(map[ProposalID]AppendProposal),
	}
}

type appendProposals struct {
	proposals map[ProposalID]AppendProposal
}

func (p *appendProposals) register(proposal AppendProposal) {
	p.proposals[proposal.ID()] = proposal
}

func (p *appendProposals) unregister(id ProposalID) {
	delete(p.proposals, id)
}

func (p *appendProposals) Get(id ProposalID) (AppendProposal, bool) {
	proposal, ok := p.proposals[id]
	return proposal, ok
}

func (p *appendProposals) List() []AppendProposal {
	proposals := make([]AppendProposal, 0, len(p.proposals))
	for _, proposal := range p.proposals {
		proposals = append(proposals, proposal)
	}
	return proposals
}

var _ AppendProposals = &appendProposals{}

type AppendProposal interface {
	Proposal
	Request() *log.AppendRequest
	Reply(*log.AppendResponse) error
	response() *log.AppendResponse
}

func newAppendProposal(id ProposalID, session Session, request *log.AppendRequest) AppendProposal {
	return &appendProposal{
		Proposal: newProposal(id, session),
		req:      request,
	}
}

type appendProposal struct {
	Proposal
	req *log.AppendRequest
	res *log.AppendResponse
}

func (p *appendProposal) Request() *log.AppendRequest {
	return p.req
}

func (p *appendProposal) Reply(reply *log.AppendResponse) error {
	if p.res != nil {
		return errors.NewConflict("reply already sent")
	}
	p.res = reply
	return nil
}

func (p *appendProposal) response() *log.AppendResponse {
	return p.res
}

var _ AppendProposal = &appendProposal{}

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
	Request() *log.GetRequest
	Reply(*log.GetResponse) error
	response() *log.GetResponse
}

func newGetProposal(id ProposalID, session Session, request *log.GetRequest) GetProposal {
	return &getProposal{
		Proposal: newProposal(id, session),
		req:      request,
	}
}

type getProposal struct {
	Proposal
	req *log.GetRequest
	res *log.GetResponse
}

func (p *getProposal) Request() *log.GetRequest {
	return p.req
}

func (p *getProposal) Reply(reply *log.GetResponse) error {
	if p.res != nil {
		return errors.NewConflict("reply already sent")
	}
	p.res = reply
	return nil
}

func (p *getProposal) response() *log.GetResponse {
	return p.res
}

var _ GetProposal = &getProposal{}

type FirstEntryProposals interface {
	register(FirstEntryProposal)
	unregister(ProposalID)
	Get(ProposalID) (FirstEntryProposal, bool)
	List() []FirstEntryProposal
}

func newFirstEntryProposals() FirstEntryProposals {
	return &firstEntryProposals{
		proposals: make(map[ProposalID]FirstEntryProposal),
	}
}

type firstEntryProposals struct {
	proposals map[ProposalID]FirstEntryProposal
}

func (p *firstEntryProposals) register(proposal FirstEntryProposal) {
	p.proposals[proposal.ID()] = proposal
}

func (p *firstEntryProposals) unregister(id ProposalID) {
	delete(p.proposals, id)
}

func (p *firstEntryProposals) Get(id ProposalID) (FirstEntryProposal, bool) {
	proposal, ok := p.proposals[id]
	return proposal, ok
}

func (p *firstEntryProposals) List() []FirstEntryProposal {
	proposals := make([]FirstEntryProposal, 0, len(p.proposals))
	for _, proposal := range p.proposals {
		proposals = append(proposals, proposal)
	}
	return proposals
}

var _ FirstEntryProposals = &firstEntryProposals{}

type FirstEntryProposal interface {
	Proposal
	Request() *log.FirstEntryRequest
	Reply(*log.FirstEntryResponse) error
	response() *log.FirstEntryResponse
}

func newFirstEntryProposal(id ProposalID, session Session, request *log.FirstEntryRequest) FirstEntryProposal {
	return &firstEntryProposal{
		Proposal: newProposal(id, session),
		req:      request,
	}
}

type firstEntryProposal struct {
	Proposal
	req *log.FirstEntryRequest
	res *log.FirstEntryResponse
}

func (p *firstEntryProposal) Request() *log.FirstEntryRequest {
	return p.req
}

func (p *firstEntryProposal) Reply(reply *log.FirstEntryResponse) error {
	if p.res != nil {
		return errors.NewConflict("reply already sent")
	}
	p.res = reply
	return nil
}

func (p *firstEntryProposal) response() *log.FirstEntryResponse {
	return p.res
}

var _ FirstEntryProposal = &firstEntryProposal{}

type LastEntryProposals interface {
	register(LastEntryProposal)
	unregister(ProposalID)
	Get(ProposalID) (LastEntryProposal, bool)
	List() []LastEntryProposal
}

func newLastEntryProposals() LastEntryProposals {
	return &lastEntryProposals{
		proposals: make(map[ProposalID]LastEntryProposal),
	}
}

type lastEntryProposals struct {
	proposals map[ProposalID]LastEntryProposal
}

func (p *lastEntryProposals) register(proposal LastEntryProposal) {
	p.proposals[proposal.ID()] = proposal
}

func (p *lastEntryProposals) unregister(id ProposalID) {
	delete(p.proposals, id)
}

func (p *lastEntryProposals) Get(id ProposalID) (LastEntryProposal, bool) {
	proposal, ok := p.proposals[id]
	return proposal, ok
}

func (p *lastEntryProposals) List() []LastEntryProposal {
	proposals := make([]LastEntryProposal, 0, len(p.proposals))
	for _, proposal := range p.proposals {
		proposals = append(proposals, proposal)
	}
	return proposals
}

var _ LastEntryProposals = &lastEntryProposals{}

type LastEntryProposal interface {
	Proposal
	Request() *log.LastEntryRequest
	Reply(*log.LastEntryResponse) error
	response() *log.LastEntryResponse
}

func newLastEntryProposal(id ProposalID, session Session, request *log.LastEntryRequest) LastEntryProposal {
	return &lastEntryProposal{
		Proposal: newProposal(id, session),
		req:      request,
	}
}

type lastEntryProposal struct {
	Proposal
	req *log.LastEntryRequest
	res *log.LastEntryResponse
}

func (p *lastEntryProposal) Request() *log.LastEntryRequest {
	return p.req
}

func (p *lastEntryProposal) Reply(reply *log.LastEntryResponse) error {
	if p.res != nil {
		return errors.NewConflict("reply already sent")
	}
	p.res = reply
	return nil
}

func (p *lastEntryProposal) response() *log.LastEntryResponse {
	return p.res
}

var _ LastEntryProposal = &lastEntryProposal{}

type PrevEntryProposals interface {
	register(PrevEntryProposal)
	unregister(ProposalID)
	Get(ProposalID) (PrevEntryProposal, bool)
	List() []PrevEntryProposal
}

func newPrevEntryProposals() PrevEntryProposals {
	return &prevEntryProposals{
		proposals: make(map[ProposalID]PrevEntryProposal),
	}
}

type prevEntryProposals struct {
	proposals map[ProposalID]PrevEntryProposal
}

func (p *prevEntryProposals) register(proposal PrevEntryProposal) {
	p.proposals[proposal.ID()] = proposal
}

func (p *prevEntryProposals) unregister(id ProposalID) {
	delete(p.proposals, id)
}

func (p *prevEntryProposals) Get(id ProposalID) (PrevEntryProposal, bool) {
	proposal, ok := p.proposals[id]
	return proposal, ok
}

func (p *prevEntryProposals) List() []PrevEntryProposal {
	proposals := make([]PrevEntryProposal, 0, len(p.proposals))
	for _, proposal := range p.proposals {
		proposals = append(proposals, proposal)
	}
	return proposals
}

var _ PrevEntryProposals = &prevEntryProposals{}

type PrevEntryProposal interface {
	Proposal
	Request() *log.PrevEntryRequest
	Reply(*log.PrevEntryResponse) error
	response() *log.PrevEntryResponse
}

func newPrevEntryProposal(id ProposalID, session Session, request *log.PrevEntryRequest) PrevEntryProposal {
	return &prevEntryProposal{
		Proposal: newProposal(id, session),
		req:      request,
	}
}

type prevEntryProposal struct {
	Proposal
	req *log.PrevEntryRequest
	res *log.PrevEntryResponse
}

func (p *prevEntryProposal) Request() *log.PrevEntryRequest {
	return p.req
}

func (p *prevEntryProposal) Reply(reply *log.PrevEntryResponse) error {
	if p.res != nil {
		return errors.NewConflict("reply already sent")
	}
	p.res = reply
	return nil
}

func (p *prevEntryProposal) response() *log.PrevEntryResponse {
	return p.res
}

var _ PrevEntryProposal = &prevEntryProposal{}

type NextEntryProposals interface {
	register(NextEntryProposal)
	unregister(ProposalID)
	Get(ProposalID) (NextEntryProposal, bool)
	List() []NextEntryProposal
}

func newNextEntryProposals() NextEntryProposals {
	return &nextEntryProposals{
		proposals: make(map[ProposalID]NextEntryProposal),
	}
}

type nextEntryProposals struct {
	proposals map[ProposalID]NextEntryProposal
}

func (p *nextEntryProposals) register(proposal NextEntryProposal) {
	p.proposals[proposal.ID()] = proposal
}

func (p *nextEntryProposals) unregister(id ProposalID) {
	delete(p.proposals, id)
}

func (p *nextEntryProposals) Get(id ProposalID) (NextEntryProposal, bool) {
	proposal, ok := p.proposals[id]
	return proposal, ok
}

func (p *nextEntryProposals) List() []NextEntryProposal {
	proposals := make([]NextEntryProposal, 0, len(p.proposals))
	for _, proposal := range p.proposals {
		proposals = append(proposals, proposal)
	}
	return proposals
}

var _ NextEntryProposals = &nextEntryProposals{}

type NextEntryProposal interface {
	Proposal
	Request() *log.NextEntryRequest
	Reply(*log.NextEntryResponse) error
	response() *log.NextEntryResponse
}

func newNextEntryProposal(id ProposalID, session Session, request *log.NextEntryRequest) NextEntryProposal {
	return &nextEntryProposal{
		Proposal: newProposal(id, session),
		req:      request,
	}
}

type nextEntryProposal struct {
	Proposal
	req *log.NextEntryRequest
	res *log.NextEntryResponse
}

func (p *nextEntryProposal) Request() *log.NextEntryRequest {
	return p.req
}

func (p *nextEntryProposal) Reply(reply *log.NextEntryResponse) error {
	if p.res != nil {
		return errors.NewConflict("reply already sent")
	}
	p.res = reply
	return nil
}

func (p *nextEntryProposal) response() *log.NextEntryResponse {
	return p.res
}

var _ NextEntryProposal = &nextEntryProposal{}

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
	Request() *log.RemoveRequest
	Reply(*log.RemoveResponse) error
	response() *log.RemoveResponse
}

func newRemoveProposal(id ProposalID, session Session, request *log.RemoveRequest) RemoveProposal {
	return &removeProposal{
		Proposal: newProposal(id, session),
		req:      request,
	}
}

type removeProposal struct {
	Proposal
	req *log.RemoveRequest
	res *log.RemoveResponse
}

func (p *removeProposal) Request() *log.RemoveRequest {
	return p.req
}

func (p *removeProposal) Reply(reply *log.RemoveResponse) error {
	if p.res != nil {
		return errors.NewConflict("reply already sent")
	}
	p.res = reply
	return nil
}

func (p *removeProposal) response() *log.RemoveResponse {
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
	Request() *log.ClearRequest
	Reply(*log.ClearResponse) error
	response() *log.ClearResponse
}

func newClearProposal(id ProposalID, session Session, request *log.ClearRequest) ClearProposal {
	return &clearProposal{
		Proposal: newProposal(id, session),
		req:      request,
	}
}

type clearProposal struct {
	Proposal
	req *log.ClearRequest
	res *log.ClearResponse
}

func (p *clearProposal) Request() *log.ClearRequest {
	return p.req
}

func (p *clearProposal) Reply(reply *log.ClearResponse) error {
	if p.res != nil {
		return errors.NewConflict("reply already sent")
	}
	p.res = reply
	return nil
}

func (p *clearProposal) response() *log.ClearResponse {
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
	Request() *log.EventsRequest
	Notify(*log.EventsResponse) error
	Close() error
}

func newEventsProposal(id ProposalID, session Session, request *log.EventsRequest, stream rsm.Stream) EventsProposal {
	return &eventsProposal{
		Proposal: newProposal(id, session),
		request:  request,
		stream:   stream,
	}
}

type eventsProposal struct {
	Proposal
	request *log.EventsRequest
	stream  rsm.Stream
}

func (p *eventsProposal) Request() *log.EventsRequest {
	return p.request
}

func (p *eventsProposal) Notify(notification *log.EventsResponse) error {
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
	Request() *log.EntriesRequest
	Notify(*log.EntriesResponse) error
	Close() error
}

func newEntriesProposal(id ProposalID, session Session, request *log.EntriesRequest, stream rsm.Stream) EntriesProposal {
	return &entriesProposal{
		Proposal: newProposal(id, session),
		request:  request,
		stream:   stream,
	}
}

type entriesProposal struct {
	Proposal
	request *log.EntriesRequest
	stream  rsm.Stream
}

func (p *entriesProposal) Request() *log.EntriesRequest {
	return p.request
}

func (p *entriesProposal) Notify(notification *log.EntriesResponse) error {
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
