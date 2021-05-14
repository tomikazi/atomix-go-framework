// Code generated by atomix-go-framework. DO NOT EDIT.
package log

import (
	log "github.com/atomix/atomix-api/go/atomix/primitive/log"
	"github.com/atomix/atomix-go-framework/pkg/atomix/errors"
	"github.com/atomix/atomix-go-framework/pkg/atomix/logging"
	"github.com/atomix/atomix-go-framework/pkg/atomix/storage/protocol/rsm"
	"github.com/golang/protobuf/proto"
	"io"
)

const Type = "Log"

const (
	sizeOp       = "Size"
	appendOp     = "Append"
	getOp        = "Get"
	firstEntryOp = "FirstEntry"
	lastEntryOp  = "LastEntry"
	prevEntryOp  = "PrevEntry"
	nextEntryOp  = "NextEntry"
	removeOp     = "Remove"
	clearOp      = "Clear"
	eventsOp     = "Events"
	entriesOp    = "Entries"
)

var newServiceFunc rsm.NewServiceFunc

func registerServiceFunc(rsmf NewServiceFunc) {
	newServiceFunc = func(scheduler rsm.Scheduler, context rsm.ServiceContext) rsm.Service {
		service := &ServiceAdaptor{
			Service: rsm.NewService(scheduler, context),
			rsm:     rsmf(newServiceContext(scheduler)),
			log:     logging.GetLogger("atomix", "log", "service"),
		}
		service.init()
		return service
	}
}

type NewServiceFunc func(ServiceContext) Service

// RegisterService registers the election primitive service on the given node
func RegisterService(node *rsm.Node) {
	node.RegisterService(Type, newServiceFunc)
}

type ServiceAdaptor struct {
	rsm.Service
	rsm Service
	log logging.Logger
}

func (s *ServiceAdaptor) init() {
	s.RegisterUnaryOperation(sizeOp, s.size)
	s.RegisterUnaryOperation(appendOp, s.append)
	s.RegisterUnaryOperation(getOp, s.get)
	s.RegisterUnaryOperation(firstEntryOp, s.firstEntry)
	s.RegisterUnaryOperation(lastEntryOp, s.lastEntry)
	s.RegisterUnaryOperation(prevEntryOp, s.prevEntry)
	s.RegisterUnaryOperation(nextEntryOp, s.nextEntry)
	s.RegisterUnaryOperation(removeOp, s.remove)
	s.RegisterUnaryOperation(clearOp, s.clear)
	s.RegisterStreamOperation(eventsOp, s.events)
	s.RegisterStreamOperation(entriesOp, s.entries)
}
func (s *ServiceAdaptor) SessionOpen(rsmSession rsm.Session) {
	s.rsm.Sessions().open(newSession(rsmSession))
}

func (s *ServiceAdaptor) SessionExpired(session rsm.Session) {
	s.rsm.Sessions().expire(SessionID(session.ID()))
}

func (s *ServiceAdaptor) SessionClosed(session rsm.Session) {
	s.rsm.Sessions().close(SessionID(session.ID()))
}
func (s *ServiceAdaptor) Backup(writer io.Writer) error {
	err := s.rsm.Backup(newSnapshotWriter(writer))
	if err != nil {
		s.log.Error(err)
		return err
	}
	return nil
}

func (s *ServiceAdaptor) Restore(reader io.Reader) error {
	err := s.rsm.Restore(newSnapshotReader(reader))
	if err != nil {
		s.log.Error(err)
		return err
	}
	return nil
}
func (s *ServiceAdaptor) size(input []byte, rsmSession rsm.Session) ([]byte, error) {
	request := &log.SizeRequest{}
	err := proto.Unmarshal(input, request)
	if err != nil {
		s.log.Error(err)
		return nil, err
	}

	session, ok := s.rsm.Sessions().Get(SessionID(rsmSession.ID()))
	if !ok {
		err := errors.NewConflict("session %d not found", rsmSession.ID())
		s.log.Error(err)
		return nil, err
	}

	proposal := newSizeProposal(ProposalID(s.Index()), session, request)

	s.rsm.Proposals().Size().register(proposal)
	session.Proposals().Size().register(proposal)

	defer func() {
		session.Proposals().Size().unregister(proposal.ID())
		s.rsm.Proposals().Size().unregister(proposal.ID())
	}()

	err = s.rsm.Size(proposal)
	if err != nil {
		s.log.Error(err)
		return nil, err
	}

	output, err := proto.Marshal(proposal.response())
	if err != nil {
		s.log.Error(err)
		return nil, err
	}
	return output, nil
}
func (s *ServiceAdaptor) append(input []byte, rsmSession rsm.Session) ([]byte, error) {
	request := &log.AppendRequest{}
	err := proto.Unmarshal(input, request)
	if err != nil {
		s.log.Error(err)
		return nil, err
	}

	session, ok := s.rsm.Sessions().Get(SessionID(rsmSession.ID()))
	if !ok {
		err := errors.NewConflict("session %d not found", rsmSession.ID())
		s.log.Error(err)
		return nil, err
	}

	proposal := newAppendProposal(ProposalID(s.Index()), session, request)

	s.rsm.Proposals().Append().register(proposal)
	session.Proposals().Append().register(proposal)

	defer func() {
		session.Proposals().Append().unregister(proposal.ID())
		s.rsm.Proposals().Append().unregister(proposal.ID())
	}()

	err = s.rsm.Append(proposal)
	if err != nil {
		s.log.Error(err)
		return nil, err
	}

	output, err := proto.Marshal(proposal.response())
	if err != nil {
		s.log.Error(err)
		return nil, err
	}
	return output, nil
}
func (s *ServiceAdaptor) get(input []byte, rsmSession rsm.Session) ([]byte, error) {
	request := &log.GetRequest{}
	err := proto.Unmarshal(input, request)
	if err != nil {
		s.log.Error(err)
		return nil, err
	}

	session, ok := s.rsm.Sessions().Get(SessionID(rsmSession.ID()))
	if !ok {
		err := errors.NewConflict("session %d not found", rsmSession.ID())
		s.log.Error(err)
		return nil, err
	}

	proposal := newGetProposal(ProposalID(s.Index()), session, request)

	s.rsm.Proposals().Get().register(proposal)
	session.Proposals().Get().register(proposal)

	defer func() {
		session.Proposals().Get().unregister(proposal.ID())
		s.rsm.Proposals().Get().unregister(proposal.ID())
	}()

	err = s.rsm.Get(proposal)
	if err != nil {
		s.log.Error(err)
		return nil, err
	}

	output, err := proto.Marshal(proposal.response())
	if err != nil {
		s.log.Error(err)
		return nil, err
	}
	return output, nil
}
func (s *ServiceAdaptor) firstEntry(input []byte, rsmSession rsm.Session) ([]byte, error) {
	request := &log.FirstEntryRequest{}
	err := proto.Unmarshal(input, request)
	if err != nil {
		s.log.Error(err)
		return nil, err
	}

	session, ok := s.rsm.Sessions().Get(SessionID(rsmSession.ID()))
	if !ok {
		err := errors.NewConflict("session %d not found", rsmSession.ID())
		s.log.Error(err)
		return nil, err
	}

	proposal := newFirstEntryProposal(ProposalID(s.Index()), session, request)

	s.rsm.Proposals().FirstEntry().register(proposal)
	session.Proposals().FirstEntry().register(proposal)

	defer func() {
		session.Proposals().FirstEntry().unregister(proposal.ID())
		s.rsm.Proposals().FirstEntry().unregister(proposal.ID())
	}()

	err = s.rsm.FirstEntry(proposal)
	if err != nil {
		s.log.Error(err)
		return nil, err
	}

	output, err := proto.Marshal(proposal.response())
	if err != nil {
		s.log.Error(err)
		return nil, err
	}
	return output, nil
}
func (s *ServiceAdaptor) lastEntry(input []byte, rsmSession rsm.Session) ([]byte, error) {
	request := &log.LastEntryRequest{}
	err := proto.Unmarshal(input, request)
	if err != nil {
		s.log.Error(err)
		return nil, err
	}

	session, ok := s.rsm.Sessions().Get(SessionID(rsmSession.ID()))
	if !ok {
		err := errors.NewConflict("session %d not found", rsmSession.ID())
		s.log.Error(err)
		return nil, err
	}

	proposal := newLastEntryProposal(ProposalID(s.Index()), session, request)

	s.rsm.Proposals().LastEntry().register(proposal)
	session.Proposals().LastEntry().register(proposal)

	defer func() {
		session.Proposals().LastEntry().unregister(proposal.ID())
		s.rsm.Proposals().LastEntry().unregister(proposal.ID())
	}()

	err = s.rsm.LastEntry(proposal)
	if err != nil {
		s.log.Error(err)
		return nil, err
	}

	output, err := proto.Marshal(proposal.response())
	if err != nil {
		s.log.Error(err)
		return nil, err
	}
	return output, nil
}
func (s *ServiceAdaptor) prevEntry(input []byte, rsmSession rsm.Session) ([]byte, error) {
	request := &log.PrevEntryRequest{}
	err := proto.Unmarshal(input, request)
	if err != nil {
		s.log.Error(err)
		return nil, err
	}

	session, ok := s.rsm.Sessions().Get(SessionID(rsmSession.ID()))
	if !ok {
		err := errors.NewConflict("session %d not found", rsmSession.ID())
		s.log.Error(err)
		return nil, err
	}

	proposal := newPrevEntryProposal(ProposalID(s.Index()), session, request)

	s.rsm.Proposals().PrevEntry().register(proposal)
	session.Proposals().PrevEntry().register(proposal)

	defer func() {
		session.Proposals().PrevEntry().unregister(proposal.ID())
		s.rsm.Proposals().PrevEntry().unregister(proposal.ID())
	}()

	err = s.rsm.PrevEntry(proposal)
	if err != nil {
		s.log.Error(err)
		return nil, err
	}

	output, err := proto.Marshal(proposal.response())
	if err != nil {
		s.log.Error(err)
		return nil, err
	}
	return output, nil
}
func (s *ServiceAdaptor) nextEntry(input []byte, rsmSession rsm.Session) ([]byte, error) {
	request := &log.NextEntryRequest{}
	err := proto.Unmarshal(input, request)
	if err != nil {
		s.log.Error(err)
		return nil, err
	}

	session, ok := s.rsm.Sessions().Get(SessionID(rsmSession.ID()))
	if !ok {
		err := errors.NewConflict("session %d not found", rsmSession.ID())
		s.log.Error(err)
		return nil, err
	}

	proposal := newNextEntryProposal(ProposalID(s.Index()), session, request)

	s.rsm.Proposals().NextEntry().register(proposal)
	session.Proposals().NextEntry().register(proposal)

	defer func() {
		session.Proposals().NextEntry().unregister(proposal.ID())
		s.rsm.Proposals().NextEntry().unregister(proposal.ID())
	}()

	err = s.rsm.NextEntry(proposal)
	if err != nil {
		s.log.Error(err)
		return nil, err
	}

	output, err := proto.Marshal(proposal.response())
	if err != nil {
		s.log.Error(err)
		return nil, err
	}
	return output, nil
}
func (s *ServiceAdaptor) remove(input []byte, rsmSession rsm.Session) ([]byte, error) {
	request := &log.RemoveRequest{}
	err := proto.Unmarshal(input, request)
	if err != nil {
		s.log.Error(err)
		return nil, err
	}

	session, ok := s.rsm.Sessions().Get(SessionID(rsmSession.ID()))
	if !ok {
		err := errors.NewConflict("session %d not found", rsmSession.ID())
		s.log.Error(err)
		return nil, err
	}

	proposal := newRemoveProposal(ProposalID(s.Index()), session, request)

	s.rsm.Proposals().Remove().register(proposal)
	session.Proposals().Remove().register(proposal)

	defer func() {
		session.Proposals().Remove().unregister(proposal.ID())
		s.rsm.Proposals().Remove().unregister(proposal.ID())
	}()

	err = s.rsm.Remove(proposal)
	if err != nil {
		s.log.Error(err)
		return nil, err
	}

	output, err := proto.Marshal(proposal.response())
	if err != nil {
		s.log.Error(err)
		return nil, err
	}
	return output, nil
}
func (s *ServiceAdaptor) clear(input []byte, rsmSession rsm.Session) ([]byte, error) {
	request := &log.ClearRequest{}
	err := proto.Unmarshal(input, request)
	if err != nil {
		s.log.Error(err)
		return nil, err
	}

	session, ok := s.rsm.Sessions().Get(SessionID(rsmSession.ID()))
	if !ok {
		err := errors.NewConflict("session %d not found", rsmSession.ID())
		s.log.Error(err)
		return nil, err
	}

	proposal := newClearProposal(ProposalID(s.Index()), session, request)

	s.rsm.Proposals().Clear().register(proposal)
	session.Proposals().Clear().register(proposal)

	defer func() {
		session.Proposals().Clear().unregister(proposal.ID())
		s.rsm.Proposals().Clear().unregister(proposal.ID())
	}()

	err = s.rsm.Clear(proposal)
	if err != nil {
		s.log.Error(err)
		return nil, err
	}

	output, err := proto.Marshal(proposal.response())
	if err != nil {
		s.log.Error(err)
		return nil, err
	}
	return output, nil
}
func (s *ServiceAdaptor) events(input []byte, rsmSession rsm.Session, stream rsm.Stream) (rsm.StreamCloser, error) {
	request := &log.EventsRequest{}
	err := proto.Unmarshal(input, request)
	if err != nil {
		s.log.Error(err)
		return nil, err
	}

	session, ok := s.rsm.Sessions().Get(SessionID(rsmSession.ID()))
	if !ok {
		err := errors.NewConflict("session %d not found", rsmSession.ID())
		s.log.Error(err)
		return nil, err
	}

	proposal := newEventsProposal(ProposalID(stream.ID()), session, request, stream)

	s.rsm.Proposals().Events().register(proposal)
	session.Proposals().Events().register(proposal)

	err = s.rsm.Events(proposal)
	if err != nil {
		s.log.Error(err)
		return nil, err
	}
	return func() {
		session.Proposals().Events().unregister(proposal.ID())
		s.rsm.Proposals().Events().unregister(proposal.ID())
	}, nil
}

func (s *ServiceAdaptor) entries(input []byte, rsmSession rsm.Session, stream rsm.Stream) (rsm.StreamCloser, error) {
	request := &log.EntriesRequest{}
	err := proto.Unmarshal(input, request)
	if err != nil {
		s.log.Error(err)
		return nil, err
	}

	session, ok := s.rsm.Sessions().Get(SessionID(rsmSession.ID()))
	if !ok {
		err := errors.NewConflict("session %d not found", rsmSession.ID())
		s.log.Error(err)
		return nil, err
	}

	proposal := newEntriesProposal(ProposalID(stream.ID()), session, request, stream)

	s.rsm.Proposals().Entries().register(proposal)
	session.Proposals().Entries().register(proposal)

	err = s.rsm.Entries(proposal)
	if err != nil {
		s.log.Error(err)
		return nil, err
	}
	return func() {
		session.Proposals().Entries().unregister(proposal.ID())
		s.rsm.Proposals().Entries().unregister(proposal.ID())
	}, nil
}

var _ rsm.Service = &ServiceAdaptor{}
