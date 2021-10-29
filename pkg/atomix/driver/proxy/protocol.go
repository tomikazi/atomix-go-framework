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

package proxy

import (
	protocolapi "github.com/atomix/atomix-api/go/atomix/protocol"
	"github.com/atomix/atomix-go-sdk/pkg/atomix/cluster"
	"github.com/atomix/atomix-go-sdk/pkg/atomix/driver/env"
	"github.com/atomix/atomix-go-sdk/pkg/atomix/driver/primitive"
	"github.com/atomix/atomix-go-sdk/pkg/atomix/server"
)

// ProtocolFunc is a protocol factory function
type ProtocolFunc func(cluster cluster.Cluster, env env.DriverEnv) Protocol

// Protocol is a proxy protocol
type Protocol interface {
	server.Node

	// Name returns the protocol name
	Name() string

	// Primitives returns the protocol primitives
	Primitives() *primitive.PrimitiveTypeRegistry

	// Configure configures the protocol
	Configure(config protocolapi.ProtocolConfig) error
}
