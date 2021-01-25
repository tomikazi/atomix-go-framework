{{- define "type" -}}
{{- if .Package.Import -}}
{{- printf "%s.%s" .Package.Alias .Name -}}
{{- else -}}
{{- .Name -}}
{{- end -}}
{{- end -}}

{{- define "field" }}
{{- $path := .Field.Path }}
{{- range $index, $element := $path -}}
{{- if eq $index 0 -}}
{{- if isLast $path $index -}}
{{- if $element.Type.IsPointer -}}
.Get{{ $element.Name }}()
{{- else -}}
.{{ $element.Name }}
{{- end -}}
{{- else -}}
{{- if $element.Type.IsPointer -}}
.Get{{ $element.Name }}().
{{- else -}}
.{{ $element.Name }}.
{{- end -}}
{{- end -}}
{{- else -}}
{{- if isLast $path $index -}}
{{- if $element.Type.IsPointer -}}
    Get{{ $element.Name }}()
{{- else -}}
    {{ $element.Name -}}
{{- end -}}
{{- else -}}
{{- if $element.Type.IsPointer -}}
    Get{{ $element.Name }}().
{{- else -}}
    {{ $element.Name }}.
{{- end -}}
{{- end -}}
{{- end -}}
{{- end -}}
{{- end }}

{{- define "stateKey" -}}
{{- if .Primitive.State.IsDiscrete -}}
""
{{- else if .Primitive.State.IsContinuous -}}
{{- if .Primitive.State.Entry.Key -}}
state{{ template "field" .Primitive.State.Entry.Key -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{- define "stateDigest" -}}
state{{ template "field" .Primitive.State.Entry.Digest }}
{{- end -}}

package {{ .Package.Name }}

import (
	"context"
	"math/rand"
	"time"

	"github.com/atomix/go-framework/pkg/atomix/meta"
	"github.com/atomix/go-framework/pkg/atomix/protocol/gossip"
	atime "github.com/atomix/go-framework/pkg/atomix/time"

	"github.com/golang/protobuf/proto"

	{{- if .Primitive.State.Entry.Type.Package.Import }}
	{{ .Primitive.State.Entry.Type.Package.Alias }} {{ .Primitive.State.Entry.Type.Package.Path | quote }}
	{{- end }}
)

const antiEntropyPeriod = time.Second

func newGossipProtocol(serviceID gossip.ServiceID, partition *gossip.Partition, clock atime.Clock) (GossipProtocol, error) {
	peers, err := gossip.NewPeerGroup(partition, ServiceType, serviceID)
	if err != nil {
		return nil, err
	}
	return &gossipProtocol{
		clock:  clock,
		group:  newGossipGroup(peers),
		server: newGossipServer(partition),
	}, nil
}

type GossipProtocol interface {
	Clock() atime.Clock
	Group() GossipGroup
	Server() GossipServer
}

type GossipHandler interface {
    {{- if .Primitive.State.IsDiscrete }}
	Read(ctx context.Context) (*{{ template "type" .Primitive.State.Entry.Type }}, error)
	Update(ctx context.Context, state *{{ template "type" .Primitive.State.Entry.Type }}) error
    {{- else if .Primitive.State.IsContinuous }}
    {{- if .Primitive.State.Entry.Key }}
    {{- if .Primitive.State.Entry.Key.Field.Type.IsScalar }}
	Read(ctx context.Context, key {{ .Primitive.State.Entry.Key.Field.Type.Name }}) (*{{ template "type" .Primitive.State.Entry.Type }}, error)
    {{- else }}
	Read(ctx context.Context, key {{ template "type" .Primitive.State.Entry.Key.Field.Type }}) (*{{ template "type" .Primitive.State.Entry.Type }}, error)
	{{- end }}
    {{- else }}
	Read(ctx context.Context) (*{{ template "type" .Primitive.State.Entry.Type }}, error)
    {{- end }}
    List(ctx context.Context, ch chan<- {{ template "type" .Primitive.State.Entry.Type }}) error
	Update(ctx context.Context, entry *{{ template "type" .Primitive.State.Entry.Type }}) error
    {{- end }}
}

type GossipServer interface {
	Register(GossipHandler) error
	handler() GossipHandler
}

type GossipClient interface {
    {{- if .Primitive.State.IsDiscrete }}
	Bootstrap(ctx context.Context) (*{{ template "type" .Primitive.State.Entry.Type }}, error)
	{{- else if .Primitive.State.IsContinuous }}
	Bootstrap(ctx context.Context, ch chan<- {{ template "type" .Primitive.State.Entry.Type }}) error
	{{- end }}
	Repair(ctx context.Context, state *{{ template "type" .Primitive.State.Entry.Type }}) (*{{ template "type" .Primitive.State.Entry.Type }}, error)
	Advertise(ctx context.Context, state *{{ template "type" .Primitive.State.Entry.Type }}) error
	Update(ctx context.Context, state *{{ template "type" .Primitive.State.Entry.Type }}) error
}

type GossipGroup interface {
	GossipClient
	MemberID() GossipMemberID
	Members() []GossipMember
	Member(GossipMemberID) GossipMember
}

type GossipMemberID gossip.PeerID

func (i GossipMemberID) String() string {
    return string(i)
}

type GossipMember interface {
	GossipClient
	ID() GossipMemberID
	Client() *gossip.Peer
}

type gossipProtocol struct {
	clock  atime.Clock
	group  GossipGroup
	server GossipServer
}

func (p *gossipProtocol) Clock() atime.Clock {
	return p.clock
}

func (p *gossipProtocol) Group() GossipGroup {
	return p.group
}

func (p *gossipProtocol) Server() GossipServer {
	return p.server
}

var _ GossipProtocol = &gossipProtocol{}

func newGossipGroup(group *gossip.PeerGroup) GossipGroup {
	peers := group.Peers()
	members := make([]GossipMember, 0, len(peers))
	memberIDs := make(map[GossipMemberID]GossipMember)
	for _, peer := range peers {
		member := newGossipMember(peer)
		members = append(members, member)
		memberIDs[member.ID()] = member
	}
	return &gossipGroup{
		group: group,
	}
}

type gossipGroup struct {
	group     *gossip.PeerGroup
	members   []GossipMember
	memberIDs map[GossipMemberID]GossipMember
}

func (p *gossipGroup) MemberID() GossipMemberID {
    return GossipMemberID(p.group.MemberID())
}

func (p *gossipGroup) Members() []GossipMember {
	return p.members
}

func (p *gossipGroup) Member(id GossipMemberID) GossipMember {
	return p.memberIDs[id]
}

{{- if .Primitive.State.IsDiscrete }}
func (p *gossipGroup) Bootstrap(ctx context.Context) (*{{ template "type" .Primitive.State.Entry.Type }}, error) {
	objects, err := p.group.Read(ctx, "")
	if err != nil {
		return nil, err
	}

    state := &{{ template "type" .Primitive.State.Entry.Type }}{}
	for _, object := range objects {
        {{- if .Primitive.State.Entry.Digest }}
		if meta.FromProto(object.ObjectMeta).After(meta.FromProto({{ template "stateDigest" . }})) {
		{{- end }}
			err = proto.Unmarshal(object.Value, state)
			if err != nil {
				return nil, err
			}
        {{- if .Primitive.State.Entry.Digest }}
		}
		{{- end }}
	}
	return state, nil
}
{{- else if .Primitive.State.IsContinuous }}
func (p *gossipGroup) Bootstrap(ctx context.Context, ch chan<- {{ template "type" .Primitive.State.Entry.Type }}) error {
	objectCh := make(chan gossip.Object)
	if err := p.group.ReadAll(ctx, objectCh); err != nil {
		return err
	}
	go func() {
		for object := range objectCh {
			var entry {{ template "type" .Primitive.State.Entry.Type }}
			err := proto.Unmarshal(object.Value, &entry)
			if err != nil {
				log.Errorf("Bootstrap failed: %v", err)
			} else {
				ch <- entry
			}
		}
	}()
	return nil
}
{{- end }}

func (p *gossipGroup) Repair(ctx context.Context, state *{{ template "type" .Primitive.State.Entry.Type }}) (*{{ template "type" .Primitive.State.Entry.Type }}, error) {
	objects, err := p.group.Read(ctx, {{ template "stateKey" . }})
	if err != nil {
		return nil, err
	}

	for _, object := range objects {
        {{- if .Primitive.State.Entry.Digest }}
		if meta.FromProto(object.ObjectMeta).After(meta.FromProto({{ template "stateDigest" . }})) {
		{{- end }}
			err = proto.Unmarshal(object.Value, state)
			if err != nil {
				return nil, err
			}
        {{- if .Primitive.State.Entry.Digest }}
		}
		{{- end }}
	}
	return state, nil
}

func (p *gossipGroup) Advertise(ctx context.Context, state *{{ template "type" .Primitive.State.Entry.Type }}) error {
	peers := p.group.Peers()
	peer := peers[rand.Intn(len(peers))]
    {{- if .Primitive.State.Entry.Digest }}
	peer.Advertise(ctx, {{ template "stateKey" . }}, meta.FromProto({{ template "stateDigest" . }}))
	{{- else }}
	peer.Advertise(ctx, {{ template "stateKey" . }}, meta.ObjectMeta{})
    {{- end }}
	return nil
}

func (p *gossipGroup) Update(ctx context.Context, state *{{ template "type" .Primitive.State.Entry.Type }}) error {
	bytes, err := proto.Marshal(state)
	if err != nil {
		return err
	}
	object := &gossip.Object{
        {{- if .Primitive.State.Entry.Digest }}
		ObjectMeta: {{ template "stateDigest" . }},
		{{- end }}
		Value:      bytes,
	}
	p.group.Update(ctx, object)
	return nil
}

var _ GossipGroup = &gossipGroup{}

func newGossipServer(partition *gossip.Partition) GossipServer {
	return &gossipServer{
		partition: partition,
	}
}

type gossipServer struct {
	partition     *gossip.Partition
	gossipHandler GossipHandler
}

func (s *gossipServer) Register(handler GossipHandler) error {
	s.gossipHandler = handler
	return s.partition.RegisterReplica(newReplica(handler))
}

func (s *gossipServer) handler() GossipHandler {
	return s.gossipHandler
}

var _ GossipServer = &gossipServer{}

func newGossipMember(peer *gossip.Peer) GossipMember {
	return &gossipMember{
		id:   GossipMemberID(peer.ID),
		peer: peer,
	}
}

type gossipMember struct {
	id   GossipMemberID
	peer *gossip.Peer
}

func (p *gossipMember) ID() GossipMemberID {
	return p.id
}

func (p *gossipMember) Client() *gossip.Peer {
	return p.peer
}

{{- if .Primitive.State.IsDiscrete }}
func (p *gossipMember) Bootstrap(ctx context.Context) (*{{ template "type" .Primitive.State.Entry.Type }}, error) {
	object, err := p.peer.Read(ctx, "")
	if err != nil {
		return nil, err
	}

    state := &{{ template "type" .Primitive.State.Entry.Type }}{}
    err = proto.Unmarshal(object.Value, state)
    if err != nil {
        return nil, err
    }
	return state, nil
}
{{- else if .Primitive.State.IsContinuous }}
func (p *gossipMember) Bootstrap(ctx context.Context, ch chan<- {{ template "type" .Primitive.State.Entry.Type }}) error {
	objectCh := make(chan gossip.Object)
	if err := p.peer.ReadAll(ctx, objectCh); err != nil {
		return err
	}
	go func() {
		for object := range objectCh {
			var entry {{ template "type" .Primitive.State.Entry.Type }}
			err := proto.Unmarshal(object.Value, &entry)
			if err != nil {
				log.Errorf("Bootstrap failed: %v", err)
			} else {
				ch <- entry
			}
		}
	}()
	return nil
}
{{- end }}

func (p *gossipMember) Repair(ctx context.Context, state *{{ template "type" .Primitive.State.Entry.Type }}) (*{{ template "type" .Primitive.State.Entry.Type }}, error) {
	object, err := p.peer.Read(ctx, {{ template "stateKey" . }})
	if err != nil {
		return nil, err
	}

    {{- if .Primitive.State.Entry.Digest }}
    if meta.FromProto(object.ObjectMeta).After(meta.FromProto({{ template "stateDigest" . }})) {
	{{- end }}
        err = proto.Unmarshal(object.Value, state)
        if err != nil {
            return nil, err
        }
    {{- if .Primitive.State.Entry.Digest }}
	}
	{{- end }}
	return state, nil
}

func (p *gossipMember) Advertise(ctx context.Context, state *{{ template "type" .Primitive.State.Entry.Type }}) error {
    {{- if .Primitive.State.Entry.Digest }}
	p.peer.Advertise(ctx, "", meta.FromProto({{ template "stateDigest" . }}))
	{{- else }}
	p.peer.Advertise(ctx, "", meta.ObjectMeta{})
	{{- end }}
	return nil
}

func (p *gossipMember) Update(ctx context.Context, state *{{ template "type" .Primitive.State.Entry.Type }}) error {
	bytes, err := proto.Marshal(state)
	if err != nil {
		return err
	}
	object := &gossip.Object{
        {{- if .Primitive.State.Entry.Digest }}
		ObjectMeta: {{ template "stateDigest" . }},
		{{- end }}
		Value:      bytes,
	}
	p.peer.Update(ctx, object)
	return nil
}

var _ GossipMember = &gossipMember{}

func newReplica(handler GossipHandler) gossip.Replica {
	return &gossipReplica{
		handler: handler,
	}
}

type gossipReplica struct {
	id      gossip.ServiceID
	handler GossipHandler
}

func (s *gossipReplica) ID() gossip.ServiceID {
	return s.id
}

func (s *gossipReplica) Type() gossip.ServiceType {
	return ServiceType
}

func (s *gossipReplica) Update(ctx context.Context, object *gossip.Object) error {
	state := &{{ template "type" .Primitive.State.Entry.Type }}{}
	err := proto.Unmarshal(object.Value, state)
	if err != nil {
		return err
	}
	return s.handler.Update(ctx, state)
}

func (s *gossipReplica) Read(ctx context.Context, key string) (*gossip.Object, error) {
    {{- if .Primitive.State.IsDiscrete }}
	state, err := s.handler.Read(ctx)
    {{- else if .Primitive.State.IsContinuous }}
    {{- if .Primitive.State.Entry.Key }}
	state, err := s.handler.Read(ctx, key)
	{{- else }}
	state, err := s.handler.Read(ctx)
	{{- end }}
	{{- end }}
	if err != nil {
		return nil, err
	} else if state == nil {
		return nil, nil
	}

	bytes, err := proto.Marshal(state)
	if err != nil {
		return nil, err
	}
	return &gossip.Object{
        {{- if .Primitive.State.Entry.Digest }}
        ObjectMeta: {{ template "stateDigest" . }},
        {{- end }}
        {{- if .Primitive.State.Entry.Key }}
        Key:        {{ template "stateKey" . }},
        {{- end }}
		Value:      bytes,
	}, nil
}

func (s *gossipReplica) ReadAll(ctx context.Context, ch chan<- gossip.Object) error {
    {{- if .Primitive.State.IsDiscrete }}
    errCh := make(chan error)
    go func() {
        defer close(errCh)
        state, err := s.handler.Read(ctx)
        if err != nil {
            errCh <- err
            return
        }
        bytes, err := proto.Marshal(state)
        if err != nil {
            errCh <- err
            return
        }
        object := gossip.Object{
            {{- if .Primitive.State.Entry.Digest }}
            ObjectMeta: {{ template "stateDigest" . }},
            {{- end }}
            Value:      bytes,
        }
        ch <- object
    }()
	return <-errCh
    {{- else if .Primitive.State.IsContinuous }}
	entriesCh := make(chan {{ template "type" .Primitive.State.Entry.Type }})
	errCh := make(chan error)
	go func() {
		err := s.handler.List(ctx, entriesCh)
		if err != nil {
			errCh <- err
		}
	}()
	go func() {
		defer close(errCh)
		for state := range entriesCh {
			bytes, err := proto.Marshal(&state)
			if err != nil {
				errCh <- err
				return
			}
			object := gossip.Object{
			    {{- if .Primitive.State.Entry.Digest }}
				ObjectMeta: {{ template "stateDigest" . }},
				{{- end }}
				{{- if .Primitive.State.Entry.Key }}
				Key: {{ template "stateKey" . }},
				{{- end }}
				Value:      bytes,
			}
			ch <- object
		}
	}()
	return <-errCh
    {{- end }}
}

var _ gossip.Replica = &gossipReplica{}

type GossipEngine interface {
	start()
	stop()
}

func newGossipEngine(protocol GossipProtocol) GossipEngine {
	return &gossipEngine{
		protocol: protocol,
	}
}

type gossipEngine struct {
	protocol GossipProtocol
	ticker   *time.Ticker
	cancel   context.CancelFunc
}

func (m *gossipEngine) start() {
	ctx, cancel := context.WithCancel(context.Background())
	m.cancel = cancel
	if err := m.bootstrap(ctx); err != nil {
		log.Errorf("Failed to bootstrap service: %v", err)
	}
	m.runAntiEntropy(ctx)
}

func (m *gossipEngine) bootstrap(ctx context.Context) error {
    {{- if .Primitive.State.IsDiscrete }}
    state, err := m.protocol.Group().Bootstrap(ctx)
    if err != nil {
        return err
    }
    if err := m.protocol.Server().handler().Update(ctx, state); err != nil {
        return err
    }
    {{- else if .Primitive.State.IsContinuous }}
	stateCh := make(chan {{ template "type" .Primitive.State.Entry.Type }})
	if err := m.protocol.Group().Bootstrap(ctx, stateCh); err != nil {
		return err
	}
	for state := range stateCh {
		if err := m.protocol.Server().handler().Update(ctx, &state); err != nil {
			return err
		}
	}
    {{- end }}
	return nil
}

func (m *gossipEngine) runAntiEntropy(ctx context.Context) {
	m.ticker = time.NewTicker(antiEntropyPeriod)
	for range m.ticker.C {
		if err := m.advertise(ctx); err != nil {
			log.Errorf("Anti-entropy protocol failed: %v", err)
		}
	}
}

func (m *gossipEngine) advertise(ctx context.Context) error {
	{{- if .Primitive.State.IsDiscrete }}
    state, err := m.protocol.Server().handler().Read(context.Background())
    if err != nil {
        return err
    }
    if err := m.protocol.Group().Advertise(ctx, state); err != nil {
        return err
    }
    {{- else if .Primitive.State.IsContinuous }}
	stateCh := make(chan {{ template "type" .Primitive.State.Entry.Type }})
	if err := m.protocol.Server().handler().List(ctx, stateCh); err != nil {
		return err
	}
	for state := range stateCh {
		if err := m.protocol.Group().Advertise(ctx, &state); err != nil {
			return err
		}
	}
    {{- end }}
	return nil
}

func (m *gossipEngine) stop() {
	m.ticker.Stop()
	m.cancel()
}