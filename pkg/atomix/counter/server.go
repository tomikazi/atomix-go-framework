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

package counter

import (
	"context"
	api "github.com/atomix/atomix-api/proto/atomix/counter"
	"github.com/atomix/atomix-api/proto/atomix/headers"
	"github.com/atomix/atomix-go-node/pkg/atomix/server"
	"github.com/atomix/atomix-go-node/pkg/atomix/service"
	"github.com/gogo/protobuf/proto"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func init() {
	service.RegisterServer(registerServer)
}

// registerServer registers a counter server with the given gRPC server
func registerServer(server *grpc.Server, protocol service.Protocol) {
	api.RegisterCounterServiceServer(server, newServer(protocol.Client()))
}

func newServer(client service.Client) api.CounterServiceServer {
	return &Server{
		SimpleServer: &server.SimpleServer{
			Type:   "counter",
			Client: client,
		},
	}
}

// Server is an implementation of CounterServiceServer for the counter primitive
type Server struct {
	*server.SimpleServer
}

// Create opens a new session
func (s *Server) Create(ctx context.Context, request *api.CreateRequest) (*api.CreateResponse, error) {
	log.Tracef("Received CreateRequest %+v", request)
	if err := s.Open(ctx, request.Header); err != nil {
		return nil, err
	}

	response := &api.CreateResponse{
		Header: &headers.ResponseHeader{},
	}
	log.Tracef("Sending CreateResponse %+v", response)
	return response, nil
}

// Set sets the current value of the counter
func (s *Server) Set(ctx context.Context, request *api.SetRequest) (*api.SetResponse, error) {
	log.Tracef("Received SetRequest %+v", request)

	in, err := proto.Marshal(&SetRequest{
		Value: request.Value,
	})
	if err != nil {
		return nil, err
	}

	out, header, err := s.Command(ctx, "set", in, request.Header)
	if err != nil {
		return nil, err
	}

	setResponse := &SetResponse{}
	if err = proto.Unmarshal(out, setResponse); err != nil {
		return nil, err
	}

	response := &api.SetResponse{
		Header: header,
	}
	log.Tracef("Sending SetResponse %+v", response)
	return response, nil
}

// Get gets the current value of the counter
func (s *Server) Get(ctx context.Context, request *api.GetRequest) (*api.GetResponse, error) {
	log.Tracef("Received GetRequest %+v", request)

	in, err := proto.Marshal(&GetRequest{})
	if err != nil {
		return nil, err
	}

	out, header, err := s.Query(ctx, "get", in, request.Header)
	if err != nil {
		return nil, err
	}

	getResponse := &GetResponse{}
	if err = proto.Unmarshal(out, getResponse); err != nil {
		return nil, err
	}

	response := &api.GetResponse{
		Header: header,
		Value:  getResponse.Value,
	}
	log.Tracef("Sending GetResponse %+v", response)
	return response, nil
}

// Increment increments the value of the counter by a delta
func (s *Server) Increment(ctx context.Context, request *api.IncrementRequest) (*api.IncrementResponse, error) {
	log.Tracef("Received IncrementRequest %+v", request)

	in, err := proto.Marshal(&IncrementRequest{
		Delta: request.Delta,
	})
	if err != nil {
		return nil, err
	}

	out, header, err := s.Command(ctx, "increment", in, request.Header)
	if err != nil {
		return nil, err
	}

	incrementResponse := &IncrementResponse{}
	if err = proto.Unmarshal(out, incrementResponse); err != nil {
		return nil, err
	}

	response := &api.IncrementResponse{
		Header:        header,
		PreviousValue: incrementResponse.PreviousValue,
		NextValue:     incrementResponse.NextValue,
	}
	log.Tracef("Sending IncrementResponse %+v", response)
	return response, nil
}

// Decrement decrements the value of the counter by a delta
func (s *Server) Decrement(ctx context.Context, request *api.DecrementRequest) (*api.DecrementResponse, error) {
	log.Tracef("Received DecrementRequest %+v", request)

	in, err := proto.Marshal(&DecrementRequest{
		Delta: request.Delta,
	})
	if err != nil {
		return nil, err
	}

	out, header, err := s.Command(ctx, "decrement", in, request.Header)
	if err != nil {
		return nil, err
	}

	decrementResponse := &DecrementResponse{}
	if err = proto.Unmarshal(out, decrementResponse); err != nil {
		return nil, err
	}

	response := &api.DecrementResponse{
		Header:        header,
		PreviousValue: decrementResponse.PreviousValue,
		NextValue:     decrementResponse.NextValue,
	}
	log.Tracef("Sending DecrementResponse %+v", response)
	return response, nil
}

// CheckAndSet updates the value of the counter conditionally
func (s *Server) CheckAndSet(ctx context.Context, request *api.CheckAndSetRequest) (*api.CheckAndSetResponse, error) {
	log.Tracef("Received CheckAndSetRequest %+v", request)

	in, err := proto.Marshal(&CheckAndSetRequest{
		Expect: request.Expect,
		Update: request.Update,
	})
	if err != nil {
		return nil, err
	}

	out, header, err := s.Command(ctx, "cas", in, request.Header)
	if err != nil {
		return nil, err
	}

	casResponse := &CheckAndSetResponse{}
	if err = proto.Unmarshal(out, casResponse); err != nil {
		return nil, err
	}

	response := &api.CheckAndSetResponse{
		Header:    header,
		Succeeded: casResponse.Succeeded,
	}
	log.Tracef("Sending CheckAndSetResponse %+v", response)
	return response, nil
}

// Close closes a session
func (s *Server) Close(ctx context.Context, request *api.CloseRequest) (*api.CloseResponse, error) {
	log.Tracef("Received CloseRequest %+v", request)
	if request.Delete {
		if err := s.Delete(ctx, request.Header); err != nil {
			return nil, err
		}
	}

	response := &api.CloseResponse{
		Header: &headers.ResponseHeader{},
	}
	log.Tracef("Sending CloseResponse %+v", response)
	return response, nil
}
