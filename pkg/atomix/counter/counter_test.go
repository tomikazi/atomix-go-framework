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
	"github.com/atomix/atomix-api/proto/atomix/counter"
	"github.com/atomix/atomix-go-node/pkg/atomix/test"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"testing"
)

func TestCounter(t *testing.T) {
	node := test.NewTestNode()
	defer node.Stop()


	conn, err := grpc.Dial(":5678", grpc.WithInsecure())
	assert.NoError(t, err)
	defer conn.Close()

	client := counter.NewCounterServiceClient(conn)
	assert.NotNil(t, client)
}
