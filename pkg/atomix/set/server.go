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

package set

import (
	"context"
	"github.com/atomix/atomix-api/proto/atomix/headers"
	api "github.com/atomix/atomix-api/proto/atomix/set"
	"github.com/atomix/atomix-go-node/pkg/atomix/node"
	"github.com/atomix/atomix-go-node/pkg/atomix/server"
	streams "github.com/atomix/atomix-go-node/pkg/atomix/stream"
	"github.com/gogo/protobuf/proto"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func init() {
	node.RegisterServer(registerServer)
}

// registerServer registers a set server with the given gRPC server
func registerServer(server *grpc.Server, protocol node.Protocol) {
	api.RegisterSetServiceServer(server, newServer(protocol.Client()))
}

func newServer(client node.Client) api.SetServiceServer {
	return &Server{
		SessionizedServer: &server.SessionizedServer{
			Type:   setType,
			Client: client,
		},
	}
}

// Server is an implementation of SetServiceServer for the set primitive
type Server struct {
	*server.SessionizedServer
}

// Size gets the number of elements in the set
func (s *Server) Size(ctx context.Context, request *api.SizeRequest) (*api.SizeResponse, error) {
	log.Tracef("Received SizeRequest %+v", request)
	in, err := proto.Marshal(&SizeRequest{})
	if err != nil {
		return nil, err
	}

	out, header, err := s.Query(ctx, opSize, in, request.Header)
	if err != nil {
		return nil, err
	}

	sizeResponse := &SizeResponse{}
	if err = proto.Unmarshal(out, sizeResponse); err != nil {
		return nil, err
	}

	response := &api.SizeResponse{
		Header: header,
		Size_:  sizeResponse.Size_,
	}
	log.Tracef("Sending SizeResponse %+v", response)
	return response, nil
}

// Contains checks whether the set contains an element
func (s *Server) Contains(ctx context.Context, request *api.ContainsRequest) (*api.ContainsResponse, error) {
	log.Tracef("Received ContainsRequest %+v", request)
	in, err := proto.Marshal(&ContainsRequest{
		Value: request.Value,
	})
	if err != nil {
		return nil, err
	}

	out, header, err := s.Query(ctx, opContains, in, request.Header)
	if err != nil {
		return nil, err
	}

	containsResponse := &ContainsResponse{}
	if err = proto.Unmarshal(out, containsResponse); err != nil {
		return nil, err
	}

	response := &api.ContainsResponse{
		Header:   header,
		Contains: containsResponse.Contains,
	}
	log.Tracef("Sending ContainsResponse %+v", response)
	return response, nil
}

// Add adds an element to the set
func (s *Server) Add(ctx context.Context, request *api.AddRequest) (*api.AddResponse, error) {
	log.Tracef("Received AddRequest %+v", request)
	in, err := proto.Marshal(&AddRequest{
		Value: request.Value,
	})
	if err != nil {
		return nil, err
	}

	out, header, err := s.Command(ctx, opAdd, in, request.Header)
	if err != nil {
		return nil, err
	}

	addResponse := &AddResponse{}
	if err = proto.Unmarshal(out, addResponse); err != nil {
		return nil, err
	}

	response := &api.AddResponse{
		Header: header,
		Added:  addResponse.Added,
	}
	log.Tracef("Sending AddResponse %+v", response)
	return response, nil
}

// Remove removes an element from the set
func (s *Server) Remove(ctx context.Context, request *api.RemoveRequest) (*api.RemoveResponse, error) {
	log.Tracef("Received RemoveRequest %+v", request)
	in, err := proto.Marshal(&RemoveRequest{
		Value: request.Value,
	})
	if err != nil {
		return nil, err
	}

	out, header, err := s.Command(ctx, opRemove, in, request.Header)
	if err != nil {
		return nil, err
	}

	removeResponse := &RemoveResponse{}
	if err = proto.Unmarshal(out, removeResponse); err != nil {
		return nil, err
	}

	response := &api.RemoveResponse{
		Header:  header,
		Removed: removeResponse.Removed,
	}
	log.Tracef("Sending RemoveResponse %+v", response)
	return response, nil
}

// Clear removes all elements from the set
func (s *Server) Clear(ctx context.Context, request *api.ClearRequest) (*api.ClearResponse, error) {
	log.Tracef("Received ClearRequest %+v", request)
	in, err := proto.Marshal(&ClearRequest{})
	if err != nil {
		return nil, err
	}

	out, header, err := s.Command(ctx, opClear, in, request.Header)
	if err != nil {
		return nil, err
	}

	clearResponse := &ClearResponse{}
	if err = proto.Unmarshal(out, clearResponse); err != nil {
		return nil, err
	}

	response := &api.ClearResponse{
		Header: header,
	}
	log.Tracef("Sending ClearResponse %+v", response)
	return response, nil
}

// Events listens for set change events
func (s *Server) Events(request *api.EventRequest, srv api.SetService_EventsServer) error {
	log.Tracef("Received EventRequest %+v", request)
	in, err := proto.Marshal(&ListenRequest{
		Replay: request.Replay,
	})
	if err != nil {
		return err
	}

	stream := streams.NewBufferedStream()
	if err := s.CommandStream(srv.Context(), opEvents, in, request.Header, stream); err != nil {
		return err
	}

	for {
		result, ok := stream.Receive()
		if !ok {
			break
		}

		if result.Failed() {
			return result.Error
		}

		response := &ListenResponse{}
		output := result.Value.(server.SessionOutput)
		if err = proto.Unmarshal(output.Value.([]byte), response); err != nil {
			return err
		}

		var eventResponse *api.EventResponse
		switch output.Header.Type {
		case headers.ResponseType_OPEN_STREAM:
			eventResponse = &api.EventResponse{
				Header: output.Header,
			}
		case headers.ResponseType_CLOSE_STREAM:
			eventResponse = &api.EventResponse{
				Header: output.Header,
			}
		default:
			eventResponse = &api.EventResponse{
				Header: output.Header,
				Type:   getEventType(response.Type),
				Value:  response.Value,
			}
		}

		log.Tracef("Sending EventResponse %+v", response)
		if err = srv.Send(eventResponse); err != nil {
			return err
		}
	}
	log.Tracef("Finished EventRequest %+v", request)
	return nil
}

// Iterate lists all elements currently in the set
func (s *Server) Iterate(request *api.IterateRequest, srv api.SetService_IterateServer) error {
	log.Tracef("Received IterateRequest %+v", request)
	in, err := proto.Marshal(&IterateRequest{})
	if err != nil {
		return err
	}

	stream := streams.NewBufferedStream()
	if err := s.QueryStream(srv.Context(), opIterate, in, request.Header, stream); err != nil {
		return err
	}

	for {
		result, ok := stream.Receive()
		if !ok {
			break
		}

		if result.Failed() {
			return result.Error
		}

		response := &IterateResponse{}
		output := result.Value.(server.SessionOutput)
		if err = proto.Unmarshal(output.Value.([]byte), response); err != nil {
			return err
		}

		var iterateResponse *api.IterateResponse
		switch output.Header.Type {
		case headers.ResponseType_OPEN_STREAM:
			iterateResponse = &api.IterateResponse{
				Header: output.Header,
			}
		case headers.ResponseType_CLOSE_STREAM:
			iterateResponse = &api.IterateResponse{
				Header: output.Header,
			}
		default:
			iterateResponse = &api.IterateResponse{
				Header: output.Header,
				Value:  response.Value,
			}
		}

		log.Tracef("Sending IterateResponse %+v", response)
		if err = srv.Send(iterateResponse); err != nil {
			return err
		}
	}
	log.Tracef("Finished IterateRequest %+v", request)
	return nil
}

// Create opens a new session
func (s *Server) Create(ctx context.Context, request *api.CreateRequest) (*api.CreateResponse, error) {
	log.Tracef("Received CreateRequest %+v", request)
	header, err := s.OpenSession(ctx, request.Header, request.Timeout)
	if err != nil {
		return nil, err
	}
	response := &api.CreateResponse{
		Header: header,
	}
	log.Tracef("Sending CreateResponse %+v", response)
	return response, nil
}

// KeepAlive keeps an existing session alive
func (s *Server) KeepAlive(ctx context.Context, request *api.KeepAliveRequest) (*api.KeepAliveResponse, error) {
	log.Tracef("Received KeepAliveRequest %+v", request)
	header, err := s.KeepAliveSession(ctx, request.Header)
	if err != nil {
		return nil, err
	}
	response := &api.KeepAliveResponse{
		Header: header,
	}
	log.Tracef("Sending KeepAliveResponse %+v", response)
	return response, nil
}

// Close closes a session
func (s *Server) Close(ctx context.Context, request *api.CloseRequest) (*api.CloseResponse, error) {
	log.Tracef("Received CloseRequest %+v", request)
	if request.Delete {
		header, err := s.Delete(ctx, request.Header)
		if err != nil {
			return nil, err
		}
		response := &api.CloseResponse{
			Header: header,
		}
		log.Tracef("Sending CloseResponse %+v", response)
		return response, nil
	}

	header, err := s.CloseSession(ctx, request.Header)
	if err != nil {
		return nil, err
	}
	response := &api.CloseResponse{
		Header: header,
	}
	log.Tracef("Sending CloseResponse %+v", response)
	return response, nil
}

func getEventType(eventType ListenResponse_Type) api.EventResponse_Type {
	switch eventType {
	case ListenResponse_NONE:
		return api.EventResponse_NONE
	case ListenResponse_ADDED:
		return api.EventResponse_ADDED
	case ListenResponse_REMOVED:
		return api.EventResponse_REMOVED
	}
	return api.EventResponse_NONE
}
