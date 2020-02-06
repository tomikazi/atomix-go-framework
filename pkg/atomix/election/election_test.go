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

package election

import (
	"context"
	client "github.com/atomix/go-client/pkg/client/election"
	"github.com/atomix/go-client/pkg/client/primitive"
	_ "github.com/atomix/go-framework/pkg/atomix/session"
	"github.com/atomix/go-framework/pkg/atomix/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestElection(t *testing.T) {
	partition, node := test.StartTestNode()
	defer node.Stop()

	session1, err := primitive.NewSession(context.TODO(), partition)
	assert.NoError(t, err)
	defer session1.Close()

	name := primitive.NewName("default", "test", "default", "test")
	election1, err := client.New(context.TODO(), name, []*primitive.Session{session1})
	assert.NoError(t, err)
	assert.NotNil(t, election1)

	session2, err := primitive.NewSession(context.TODO(), partition)
	assert.NoError(t, err)
	defer session2.Close()

	election2, err := client.New(context.TODO(), name, []*primitive.Session{session2})
	assert.NoError(t, err)
	assert.NotNil(t, election2)

	session3, err := primitive.NewSession(context.TODO(), partition)
	assert.NoError(t, err)
	defer session3.Close()

	election3, err := client.New(context.TODO(), name, []*primitive.Session{session3})
	assert.NoError(t, err)
	assert.NotNil(t, election3)

	ch := make(chan *client.Event)
	err = election1.Watch(context.TODO(), ch)
	assert.NoError(t, err)

	term, err := election1.GetTerm(context.TODO())
	assert.NoError(t, err)
	assert.Equal(t, uint64(0), term.ID)
	assert.Equal(t, "", term.Leader)
	assert.Len(t, term.Candidates, 0)

	term, err = election1.Enter(context.TODO())
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), term.ID)
	assert.Equal(t, election1.ID(), term.Leader)
	assert.Len(t, term.Candidates, 1)
	assert.Equal(t, election1.ID(), term.Candidates[0])

	event := <-ch
	assert.Equal(t, client.EventChanged, event.Type)
	assert.Equal(t, uint64(1), event.Term.ID)
	assert.Equal(t, election1.ID(), event.Term.Leader)
	assert.Len(t, event.Term.Candidates, 1)
	assert.Equal(t, election1.ID(), event.Term.Candidates[0])

	term, err = election2.Enter(context.TODO())
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), term.ID)
	assert.Equal(t, election1.ID(), term.Leader)
	assert.Len(t, term.Candidates, 2)
	assert.Equal(t, election1.ID(), term.Candidates[0])
	assert.Equal(t, election2.ID(), term.Candidates[1])

	event = <-ch
	assert.Equal(t, client.EventChanged, event.Type)
	assert.Equal(t, uint64(1), event.Term.ID)
	assert.Equal(t, election1.ID(), event.Term.Leader)
	assert.Len(t, event.Term.Candidates, 2)
	assert.Equal(t, election1.ID(), event.Term.Candidates[0])
	assert.Equal(t, election2.ID(), event.Term.Candidates[1])

	term, err = election3.Enter(context.TODO())
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), term.ID)
	assert.Equal(t, election1.ID(), term.Leader)
	assert.Len(t, term.Candidates, 3)
	assert.Equal(t, election1.ID(), term.Candidates[0])
	assert.Equal(t, election2.ID(), term.Candidates[1])
	assert.Equal(t, election3.ID(), term.Candidates[2])

	event = <-ch
	assert.Equal(t, client.EventChanged, event.Type)
	assert.Equal(t, uint64(1), event.Term.ID)
	assert.Equal(t, election1.ID(), event.Term.Leader)
	assert.Len(t, event.Term.Candidates, 3)
	assert.Equal(t, election1.ID(), event.Term.Candidates[0])
	assert.Equal(t, election2.ID(), event.Term.Candidates[1])
	assert.Equal(t, election3.ID(), event.Term.Candidates[2])

	term, err = election3.Promote(context.TODO(), election3.ID())
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), term.ID)
	assert.Equal(t, election1.ID(), term.Leader)
	assert.Len(t, term.Candidates, 3)
	assert.Equal(t, election1.ID(), term.Candidates[0])
	assert.Equal(t, election3.ID(), term.Candidates[1])
	assert.Equal(t, election2.ID(), term.Candidates[2])

	event = <-ch
	assert.Equal(t, client.EventChanged, event.Type)
	assert.Equal(t, uint64(1), event.Term.ID)
	assert.Equal(t, election1.ID(), event.Term.Leader)
	assert.Len(t, event.Term.Candidates, 3)
	assert.Equal(t, election1.ID(), event.Term.Candidates[0])
	assert.Equal(t, election3.ID(), event.Term.Candidates[1])
	assert.Equal(t, election2.ID(), event.Term.Candidates[2])

	term, err = election3.Promote(context.TODO(), election3.ID())
	assert.NoError(t, err)
	assert.Equal(t, uint64(2), term.ID)
	assert.Equal(t, election3.ID(), term.Leader)
	assert.Len(t, term.Candidates, 3)
	assert.Equal(t, election3.ID(), term.Candidates[0])
	assert.Equal(t, election1.ID(), term.Candidates[1])
	assert.Equal(t, election2.ID(), term.Candidates[2])

	event = <-ch
	assert.Equal(t, client.EventChanged, event.Type)
	assert.Equal(t, uint64(2), event.Term.ID)
	assert.Equal(t, election3.ID(), event.Term.Leader)
	assert.Len(t, event.Term.Candidates, 3)
	assert.Equal(t, election3.ID(), event.Term.Candidates[0])
	assert.Equal(t, election1.ID(), event.Term.Candidates[1])
	assert.Equal(t, election2.ID(), event.Term.Candidates[2])

	term, err = election2.Anoint(context.TODO(), election2.ID())
	assert.NoError(t, err)
	assert.Equal(t, uint64(3), term.ID)
	assert.Equal(t, election2.ID(), term.Leader)
	assert.Len(t, term.Candidates, 3)
	assert.Equal(t, election2.ID(), term.Candidates[0])
	assert.Equal(t, election3.ID(), term.Candidates[1])
	assert.Equal(t, election1.ID(), term.Candidates[2])

	event = <-ch
	assert.Equal(t, client.EventChanged, event.Type)
	assert.Equal(t, uint64(3), event.Term.ID)
	assert.Equal(t, election2.ID(), event.Term.Leader)
	assert.Len(t, event.Term.Candidates, 3)
	assert.Equal(t, election2.ID(), event.Term.Candidates[0])
	assert.Equal(t, election3.ID(), event.Term.Candidates[1])
	assert.Equal(t, election1.ID(), event.Term.Candidates[2])

	term, err = election2.Leave(context.TODO())
	assert.NoError(t, err)
	assert.Equal(t, uint64(4), term.ID)
	assert.Equal(t, election3.ID(), term.Leader)
	assert.Len(t, term.Candidates, 2)
	assert.Equal(t, election3.ID(), term.Candidates[0])
	assert.Equal(t, election1.ID(), term.Candidates[1])

	event = <-ch
	assert.Equal(t, client.EventChanged, event.Type)
	assert.Equal(t, uint64(4), event.Term.ID)
	assert.Equal(t, election3.ID(), event.Term.Leader)
	assert.Len(t, event.Term.Candidates, 2)
	assert.Equal(t, election3.ID(), event.Term.Candidates[0])
	assert.Equal(t, election1.ID(), event.Term.Candidates[1])

	term, err = election3.Evict(context.TODO(), election3.ID())
	assert.NoError(t, err)
	assert.Equal(t, uint64(5), term.ID)
	assert.Equal(t, election1.ID(), term.Leader)
	assert.Len(t, term.Candidates, 1)
	assert.Equal(t, election1.ID(), term.Candidates[0])

	event = <-ch
	assert.Equal(t, client.EventChanged, event.Type)
	assert.Equal(t, uint64(5), event.Term.ID)
	assert.Equal(t, election1.ID(), event.Term.Leader)
	assert.Len(t, event.Term.Candidates, 1)
	assert.Equal(t, election1.ID(), event.Term.Candidates[0])

	term, err = election2.Enter(context.TODO())
	assert.NoError(t, err)
	assert.Equal(t, uint64(5), term.ID)
	assert.Equal(t, election1.ID(), term.Leader)
	assert.Len(t, term.Candidates, 2)
	assert.Equal(t, election1.ID(), term.Candidates[0])

	event = <-ch
	assert.Equal(t, client.EventChanged, event.Type)
	assert.Equal(t, uint64(5), event.Term.ID)
	assert.Equal(t, election1.ID(), event.Term.Leader)
	assert.Len(t, event.Term.Candidates, 2)
	assert.Equal(t, election1.ID(), event.Term.Candidates[0])

	term, err = election1.Anoint(context.TODO(), election2.ID())
	assert.NoError(t, err)
	assert.Equal(t, uint64(6), term.ID)
	assert.Equal(t, election2.ID(), term.Leader)
	assert.Len(t, term.Candidates, 2)
	assert.Equal(t, election2.ID(), term.Candidates[0])

	event = <-ch
	assert.Equal(t, client.EventChanged, event.Type)
	assert.Equal(t, uint64(6), event.Term.ID)
	assert.Equal(t, election2.ID(), event.Term.Leader)
	assert.Len(t, event.Term.Candidates, 2)
	assert.Equal(t, election2.ID(), event.Term.Candidates[0])

	err = election2.Close(context.Background())
	assert.NoError(t, err)

	event = <-ch
	assert.Equal(t, client.EventChanged, event.Type)
	assert.Equal(t, uint64(7), event.Term.ID)
	assert.Equal(t, election1.ID(), event.Term.Leader)
	assert.Len(t, event.Term.Candidates, 1)
	assert.Equal(t, election1.ID(), event.Term.Candidates[0])

	err = election1.Close(context.Background())
	assert.NoError(t, err)
	err = election3.Close(context.Background())
	assert.NoError(t, err)

	election1, err = client.New(context.TODO(), name, []*primitive.Session{session1})
	assert.NoError(t, err)

	election2, err = client.New(context.TODO(), name, []*primitive.Session{session2})
	assert.NoError(t, err)

	term, err = election1.GetTerm(context.TODO())
	assert.NoError(t, err)
	assert.Equal(t, uint64(7), term.ID)
	assert.Equal(t, "", term.Leader)
	assert.Len(t, term.Candidates, 0)

	err = election1.Close(context.Background())
	assert.NoError(t, err)

	err = election1.Delete(context.Background())
	assert.NoError(t, err)

	err = election2.Delete(context.Background())
	assert.NoError(t, err)

	session, err := primitive.NewSession(context.TODO(), partition)
	assert.NoError(t, err)
	defer session.Close()

	election, err := client.New(context.TODO(), name, []*primitive.Session{session})
	assert.NoError(t, err)

	term, err = election.GetTerm(context.TODO())
	assert.NoError(t, err)
	assert.Equal(t, uint64(0), term.ID)
	assert.Equal(t, "", term.Leader)
	assert.Len(t, term.Candidates, 0)

}
