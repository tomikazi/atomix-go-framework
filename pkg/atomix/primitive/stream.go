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

package primitive

import (
	"container/list"
	streams "github.com/atomix/go-framework/pkg/atomix/stream"
	"github.com/atomix/go-framework/pkg/atomix/util"
	"github.com/gogo/protobuf/proto"
)

// sessionStream manages a single stream for a session
type sessionStream struct {
	ID         uint64
	Type       string
	session    *Session
	responseID uint64
	completeID uint64
	lastIndex  uint64
	ctx        PartitionContext
	stream     streams.WriteStream
	results    *list.List
}

// sessionStreamResult contains a single stream result
type sessionStreamResult struct {
	id     uint64
	index  uint64
	result streams.Result
}

// open opens the stream
func (s *sessionStream) open() {
	s.updateClock()

	bytes, err := proto.Marshal(&SessionResponse{
		Response: &SessionResponse_Command{
			Command: &SessionCommandResponse{
				Context: &SessionResponseContext{
					StreamID: s.ID,
					Index:    s.lastIndex,
					Sequence: s.responseID,
					Type:     SessionResponseType_OPEN_STREAM,
				},
			},
		},
	})
	result := streams.Result{
		Value: bytes,
		Error: err,
	}

	out := sessionStreamResult{
		id:     s.responseID,
		index:  s.ctx.Index(),
		result: result,
	}
	s.results.PushBack(out)

	util.StreamEntry(s.ctx.NodeID(), s.session.ID, s.ID).
		Tracef("Sending stream open %d %v", s.responseID, out.result)
	s.stream.Send(out.result)
}

func (s *sessionStream) updateClock() {
	// If the client acked a sequence number greater than the current event sequence number since we know the
	// client must have received it from another server.
	s.responseID++
	if s.completeID > s.responseID {
		util.StreamEntry(s.ctx.NodeID(), s.session.ID, s.ID).
			Debugf("Skipped completed result %d", s.responseID)
		return
	}

	// Record the last index sent on the stream
	s.lastIndex = s.ctx.Index()
}

func (s *sessionStream) Send(result streams.Result) {
	s.updateClock()

	// Create the stream result and add it to the results list.
	if result.Succeeded() {
		bytes, err := proto.Marshal(&SessionResponse{
			Response: &SessionResponse_Command{
				Command: &SessionCommandResponse{
					Context: &SessionResponseContext{
						StreamID: s.ID,
						Index:    s.lastIndex,
						Sequence: s.responseID,
					},
					Response: &ServiceCommandResponse{
						Response: &ServiceCommandResponse_Operation{
							Operation: &ServiceOperationResponse{
								result.Value.([]byte),
							},
						},
					},
				},
			},
		})
		result = streams.Result{
			Value: bytes,
			Error: err,
		}
	}

	out := sessionStreamResult{
		id:     s.responseID,
		index:  s.ctx.Index(),
		result: result,
	}
	s.results.PushBack(out)
	util.StreamEntry(s.ctx.NodeID(), s.session.ID, s.ID).
		Tracef("Cached response %d", s.responseID)

	// If the out channel is set, send the result
	util.StreamEntry(s.ctx.NodeID(), s.session.ID, s.ID).
		Tracef("Sending response %d %v", s.responseID, out.result)
	s.stream.Send(out.result)
}

func (s *sessionStream) Result(value interface{}, err error) {
	s.Send(streams.Result{
		Value: value,
		Error: err,
	})
}

func (s *sessionStream) Value(value interface{}) {
	s.Result(value, nil)
}

func (s *sessionStream) Error(err error) {
	s.Result(nil, err)
}

func (s *sessionStream) Close() {
	util.StreamEntry(s.ctx.NodeID(), s.session.ID, s.ID).
		Trace("Stream closed")
	s.updateClock()

	bytes, err := proto.Marshal(&SessionResponse{
		Response: &SessionResponse_Command{
			Command: &SessionCommandResponse{
				Context: &SessionResponseContext{
					StreamID: s.ID,
					Index:    s.lastIndex,
					Sequence: s.responseID,
					Type:     SessionResponseType_CLOSE_STREAM,
				},
			},
		},
	})
	result := streams.Result{
		Value: bytes,
		Error: err,
	}

	out := sessionStreamResult{
		id:     s.responseID,
		index:  s.ctx.Index(),
		result: result,
	}
	s.results.PushBack(out)

	util.StreamEntry(s.ctx.NodeID(), s.session.ID, s.ID).
		Tracef("Sending stream close %d %v", s.responseID, out.result)
	s.stream.Send(out.result)
	s.stream.Close()
}

// LastIndex returns the last index in the stream
func (s *sessionStream) LastIndex() uint64 {
	if s.results.Len() > 0 {
		return s.lastIndex
	}
	return s.ctx.Index()
}

// ack acknowledges results up to the given ID
func (s *sessionStream) ack(id uint64) {
	if id > s.completeID {
		event := s.results.Front()
		for event != nil && event.Value.(sessionStreamResult).id <= id {
			next := event.Next()
			s.results.Remove(event)
			s.completeID = event.Value.(sessionStreamResult).id
			event = next
		}
		util.StreamEntry(s.ctx.NodeID(), s.session.ID, s.ID).
			Tracef("Discarded cached responses up to %d", id)
	}
}

// replay resends results on the given channel
func (s *sessionStream) replay(stream streams.WriteStream) {
	result := s.results.Front()
	for result != nil {
		response := result.Value.(sessionStreamResult)
		util.StreamEntry(s.ctx.NodeID(), s.session.ID, s.ID).
			Tracef("Sending response %d %v", response.id, response.result)
		stream.Send(response.result)
		result = result.Next()
	}
	s.stream = stream
}
