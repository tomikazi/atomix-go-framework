package lock

import (
	"context"
	"github.com/atomix/atomix-go-node/pkg/atomix/server"
	"github.com/atomix/atomix-go-node/pkg/atomix/service"
	"github.com/atomix/atomix-go-node/proto/atomix/headers"
	pb "github.com/atomix/atomix-go-node/proto/atomix/lock"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func RegisterLockServer(server *grpc.Server, client service.Client) {
	pb.RegisterLockServiceServer(server, NewLockServiceServer(client))
}

func NewLockServiceServer(client service.Client) pb.LockServiceServer {
	return &lockServer{
		SessionizedServer: &server.SessionizedServer{
			Type:   "lock",
			Client: client,
		},
	}
}

// lockServer is an implementation of MapServiceServer for the map primitive
type lockServer struct {
	*server.SessionizedServer
}

func (s *lockServer) Create(ctx context.Context, request *pb.CreateRequest) (*pb.CreateResponse, error) {
	log.Tracef("Received CreateRequest %+v", request)
	session, err := s.OpenSession(ctx, request.Header, request.Timeout)
	if err != nil {
		return nil, err
	}
	response := &pb.CreateResponse{
		Header: &headers.ResponseHeader{
			SessionId: session,
		},
	}
	log.Tracef("Sending CreateResponse %+v", response)
	return response, nil
}

func (s *lockServer) KeepAlive(ctx context.Context, request *pb.KeepAliveRequest) (*pb.KeepAliveResponse, error) {
	log.Tracef("Received KeepAliveRequest %+v", request)
	if err := s.KeepAliveSession(ctx, request.Header); err != nil {
		return nil, err
	}
	response := &pb.KeepAliveResponse{
		Header: &headers.ResponseHeader{
			SessionId: request.Header.SessionId,
		},
	}
	log.Tracef("Sending KeepAliveResponse %+v", response)
	return response, nil
}

func (s *lockServer) Close(ctx context.Context, request *pb.CloseRequest) (*pb.CloseResponse, error) {
	log.Tracef("Received CloseRequest %+v", request)
	if err := s.CloseSession(ctx, request.Header); err != nil {
		return nil, err
	}
	response := &pb.CloseResponse{
		Header: &headers.ResponseHeader{
			SessionId: request.Header.SessionId,
		},
	}
	log.Tracef("Sending CloseResponse %+v", response)
	return response, nil
}

func (s *lockServer) Lock(ctx context.Context, request *pb.LockRequest) (*pb.LockResponse, error) {
	log.Tracef("Received LockRequest %+v", request)
	timeout, err := ptypes.Duration(request.Timeout)
	if err != nil {
		return nil, err
	}

	in, err := proto.Marshal(&LockRequest{
		Timeout: int64(timeout),
	})
	if err != nil {
		return nil, err
	}

	out, header, err := s.Command(ctx, "lock", in, request.Header)
	if err != nil {
		return nil, err
	}

	lockResponse := &LockResponse{}
	if err = proto.Unmarshal(out, lockResponse); err != nil {
		return nil, err
	}

	response := &pb.LockResponse{
		Header:  header,
		Version: uint64(lockResponse.Index),
	}
	log.Tracef("Sending LockResponse %+v", response)
	return response, nil
}

func (s *lockServer) Unlock(ctx context.Context, request *pb.UnlockRequest) (*pb.UnlockResponse, error) {
	log.Tracef("Received UnlockRequest %+v", request)
	in, err := proto.Marshal(&UnlockRequest{
		Index: int64(request.Version),
	})
	if err != nil {
		return nil, err
	}

	out, header, err := s.Command(ctx, "unlock", in, request.Header)
	if err != nil {
		return nil, err
	}

	unlockResponse := &UnlockResponse{}
	if err = proto.Unmarshal(out, unlockResponse); err != nil {
		return nil, err
	}

	response := &pb.UnlockResponse{
		Header:   header,
		Unlocked: unlockResponse.Succeeded,
	}
	log.Tracef("Sending UnlockResponse %+v", response)
	return response, nil
}

func (s *lockServer) IsLocked(ctx context.Context, request *pb.IsLockedRequest) (*pb.IsLockedResponse, error) {
	log.Tracef("Received IsLockedRequest %+v", request)
	in, err := proto.Marshal(&IsLockedRequest{
		Index: int64(request.Version),
	})
	if err != nil {
		return nil, err
	}

	out, header, err := s.Query(ctx, "islocked", in, request.Header)
	if err != nil {
		return nil, err
	}

	isLockedResponse := &IsLockedResponse{}
	if err = proto.Unmarshal(out, isLockedResponse); err != nil {
		return nil, err
	}

	response := &pb.IsLockedResponse{
		Header:   header,
		IsLocked: isLockedResponse.Locked,
	}
	log.Tracef("Sending IsLockedResponse %+v", response)
	return response, nil
}