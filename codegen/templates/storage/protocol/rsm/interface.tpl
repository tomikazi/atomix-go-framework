// Code generated by atomix-go-framework. DO NOT EDIT.

// SPDX-FileCopyrightText: 2019-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package {{ .Package.Name }}

{{ $serviceInt := printf "%sService" .Generator.Prefix }}
{{ $serviceContextInt := printf "%sServiceContext" .Generator.Prefix }}
{{ $serviceContextImpl := ( printf "%sServiceContext" .Generator.Prefix | toLowerCamel ) }}

{{ define "type" }}
{{- if .Package.Import -}}
{{- printf "%s.%s" .Package.Alias .Name -}}
{{- else -}}
{{- .Name -}}
{{- end -}}
{{- end -}}

import (
	"io"
	"fmt"
	{{ import "github.com/atomix/atomix-go-framework/pkg/atomix/storage/protocol/rsm" }}
	{{ import "github.com/atomix/atomix-go-framework/pkg/atomix/util" }}
	{{ import "github.com/golang/protobuf/proto" }}
	{{- range .Primitive.Methods }}
	{{- if ( or .Response.IsStream .Type.IsAsync ) }}
	{{ import "github.com/atomix/atomix-go-framework/pkg/atomix/errors" }}
	{{- end }}
	{{- end }}
	{{- $package := .Package }}
	{{- range .Imports }}
	{{ .Alias }} {{ .Path | quote }}
	{{- end }}
)

{{- $primitive := .Primitive }}

{{- $serviceWatcherInt := printf "%sWatcher" .Generator.Prefix }}
{{- $serviceWatcherImpl := ( ( printf "%sServiceWatcher" .Generator.Prefix ) | toLowerCamel ) }}
{{- $newServiceWatcher := printf "new%sWatcher" .Generator.Prefix }}

{{- $serviceSnapshotWriterInt := printf "%sSnapshotWriter" .Generator.Prefix }}
{{- $serviceSnapshotWriterImpl := ( ( printf "%sServiceSnapshotWriter" .Generator.Prefix ) | toLowerCamel ) }}
{{- $newServiceSnapshotWriter := printf "new%sSnapshotWriter" .Generator.Prefix }}

{{- $serviceSnapshotReaderInt := printf "%sSnapshotReader" .Generator.Prefix }}
{{- $serviceSnapshotReaderImpl := ( ( printf "%sServiceSnapshotReader" .Generator.Prefix ) | toLowerCamel ) }}
{{- $newServiceSnapshotReader := printf "new%sSnapshotReader" .Generator.Prefix }}

{{- $serviceSessionsInt := printf "%sSessions" .Generator.Prefix }}
{{- $serviceSessionsImpl := ( ( printf "%sServiceSessions" .Generator.Prefix ) | toLowerCamel ) }}
{{- $newServiceSessions := printf "new%sSessions" .Generator.Prefix }}

{{- $serviceSessionID := printf "%sSessionID" .Generator.Prefix }}
{{- $serviceSessionState := printf "%sSessionState" .Generator.Prefix }}
{{- $serviceSessionInt := printf "%sSession" .Generator.Prefix }}
{{- $serviceSessionImpl := ( ( printf "%sServiceSession" .Generator.Prefix ) | toLowerCamel ) }}
{{- $newServiceSession := printf "new%sSession" .Generator.Prefix }}

{{- $serviceProposalsInt := printf "%sProposals" .Generator.Prefix }}
{{- $serviceProposalsImpl := ( ( printf "%sServiceProposals" .Generator.Prefix ) | toLowerCamel ) }}
{{- $newServiceProposals := printf "new%sProposals" .Generator.Prefix }}

{{- $serviceProposalID := printf "%sProposalID" .Generator.Prefix }}
{{- $serviceProposalState := printf "%sProposalState" .Generator.Prefix }}
{{- $serviceProposalInt := printf "%sProposal" .Generator.Prefix }}
{{- $serviceProposalImpl := ( ( printf "%sServiceProposal" .Generator.Prefix ) | toLowerCamel ) }}
{{- $newServiceProposal := printf "new%sProposal" .Generator.Prefix }}

{{- $serviceQueryInt := printf "%sQuery" .Generator.Prefix }}
{{- $serviceQueryImpl := ( ( printf "%sServiceQuery" .Generator.Prefix ) | toLowerCamel ) }}
{{- $newServiceQuery := printf "new%sQuery" .Generator.Prefix }}

type {{ $serviceInt }} interface {
    {{ $serviceContextInt }}

    {{- if .Primitive.State }}
    Backup({{ $serviceSnapshotWriterInt }}) error
    Restore({{ $serviceSnapshotReaderInt }}) error
    {{- end }}

    {{- range .Primitive.Methods }}
    {{- $comments := split .Comment "\n" }}
    {{- range $comment := $comments }}
    {{- if $comment }}
    // {{ $comment | trim }}
    {{- end }}
    {{- end }}

    {{- $streamInt := printf "%s%sStream" $serviceInt .Name }}
    {{- $informerInt := printf "%s%sInformer" $serviceInt .Name }}
    {{- $writerInt := printf "%s%sWriter" $serviceInt .Name }}

    {{- if .Type.IsCommand }}
    {{- $proposalInt := printf "%sProposal" .Name }}
    {{- if ( and .Response.IsUnary .Type.IsSync) }}
    {{ .Name }}({{ $proposalInt }}) (*{{ template "type" .Response.Type }}, error)
    {{- else }}
    {{ .Name }}({{ $proposalInt }})
    {{- end }}
    {{- else if .Type.IsQuery }}
    {{- $queryInt := printf "%sQuery" .Name }}
    {{- if .Response.IsUnary }}
    {{ .Name }}({{ $queryInt }}) (*{{ template "type" .Response.Type }}, error)
    {{- else if .Response.IsStream }}
    {{ .Name }}({{ $queryInt }})
    {{- end }}
    {{- end }}
    {{- end }}
}

type {{ $serviceContextInt }} interface {
    Scheduler() rsm.Scheduler
    Sessions() {{ $serviceSessionsInt }}
    Proposals() {{ $serviceProposalsInt }}
}

func new{{ $serviceContextInt }}(service rsm.ServiceContext) {{ $serviceContextInt }} {
    return &{{ $serviceContextImpl }}{
        scheduler: service.Scheduler(),
        sessions:  {{ $newServiceSessions }}(service.Sessions()),
        proposals: {{ $newServiceProposals }}(service.Commands()),
    }
}

type {{ $serviceContextImpl }} struct {
    scheduler rsm.Scheduler
    sessions {{ $serviceSessionsInt }}
    proposals {{ $serviceProposalsInt }}
}

func (s *{{ $serviceContextImpl }}) Scheduler() rsm.Scheduler {
    return s.scheduler
}

func (s *{{ $serviceContextImpl }}) Sessions() {{ $serviceSessionsInt }} {
    return s.sessions
}

func (s *{{ $serviceContextImpl }}) Proposals() {{ $serviceProposalsInt }} {
    return s.proposals
}

var _ {{ $serviceContextInt }} = &{{ $serviceContextImpl }}{}

type {{ $serviceSnapshotWriterInt }} interface {
    WriteState(*{{ template "type" .Primitive.State.Type }}) error
}

func {{ $newServiceSnapshotWriter }}(writer io.Writer) {{ $serviceSnapshotWriterInt }} {
    return &{{ $serviceSnapshotWriterImpl }}{
        writer: writer,
    }
}

type {{ $serviceSnapshotWriterImpl }} struct {
    writer io.Writer
}

func (w *{{ $serviceSnapshotWriterImpl }}) WriteState(state *{{ template "type" .Primitive.State.Type }}) error {
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

var _ {{ $serviceSnapshotWriterInt }} = &{{ $serviceSnapshotWriterImpl }}{}

type {{ $serviceSnapshotReaderInt }} interface {
    ReadState() (*{{ template "type" .Primitive.State.Type }}, error)
}

func {{ $newServiceSnapshotReader }}(reader io.Reader) {{ $serviceSnapshotReaderInt }} {
    return &{{ $serviceSnapshotReaderImpl }}{
        reader: reader,
    }
}

type {{ $serviceSnapshotReaderImpl }} struct {
    reader io.Reader
}

func (r *{{ $serviceSnapshotReaderImpl }}) ReadState() (*{{ template "type" .Primitive.State.Type }}, error) {
    bytes, err := util.ReadBytes(r.reader)
	if err != nil {
		return nil, err
	}
	state := &{{ template "type" .Primitive.State.Type }}{}
	err = proto.Unmarshal(bytes, state)
	if err != nil {
		return nil, err
	}
	return state, nil
}

var _ {{ $serviceSnapshotReaderInt }} = &{{ $serviceSnapshotReaderImpl }}{}

type {{ $serviceSessionsInt }} interface {
    Get({{ $serviceSessionID }}) ({{ $serviceSessionInt }}, bool)
    List() []{{ $serviceSessionInt }}
}

func {{ $newServiceSessions }}(sessions rsm.Sessions) {{ $serviceSessionsInt }} {
    return &{{ $serviceSessionsImpl }}{
        sessions: sessions,
    }
}

type {{ $serviceSessionsImpl }} struct {
    sessions rsm.Sessions
}

func (s *{{ $serviceSessionsImpl }}) Get(id {{ $serviceSessionID }}) ({{ $serviceSessionInt }}, bool) {
    session, ok := s.sessions.Get(rsm.SessionID(id))
    if !ok {
        return nil, false
    }
    return {{ $newServiceSession }}(session), true
}

func (s *{{ $serviceSessionsImpl }}) List() []{{ $serviceSessionInt }} {
    serviceSessions := s.sessions.List()
    sessions := make([]{{ $serviceSessionInt }}, len(serviceSessions))
    for i, serviceSession := range serviceSessions {
        sessions[i] = {{ $newServiceSession }}(serviceSession)
    }
    return sessions
}

var _ {{ $serviceSessionsInt }} = &{{ $serviceSessionsImpl }}{}

type {{ $serviceSessionID }} uint64

type {{ $serviceSessionState }} int

const (
	SessionClosed {{ $serviceSessionState }} = iota
	SessionOpen
)

type {{ $serviceWatcherInt }} interface {
	Cancel()
}

func {{ $newServiceWatcher }}(watcher rsm.Watcher) {{ $serviceWatcherInt }} {
	return &{{ $serviceWatcherImpl }}{
		watcher: watcher,
	}
}

type {{ $serviceWatcherImpl }} struct {
	watcher rsm.Watcher
}

func (s *{{ $serviceWatcherImpl }}) Cancel() {
    s.watcher.Cancel()
}

var _ {{ $serviceWatcherInt }} = &{{ $serviceWatcherImpl }}{}

type {{ $serviceSessionInt }} interface {
    ID() {{ $serviceSessionID }}
	State() {{ $serviceSessionState }}
	Watch(func({{ $serviceSessionState }})) {{ $serviceWatcherInt }}
    Proposals() {{ $serviceProposalsInt }}
}

func {{ $newServiceSession }}(session rsm.Session) {{ $serviceSessionInt }} {
    return &{{ $serviceSessionImpl }}{
        session:    session,
        proposals: {{ $newServiceProposals }}(session.Commands()),
    }
}

type {{ $serviceSessionImpl }} struct {
    session   rsm.Session
    proposals {{ $serviceProposalsInt }}
}

func (s *{{ $serviceSessionImpl }}) ID() {{ $serviceSessionID }} {
    return {{ $serviceSessionID }}(s.session.ID())
}

func (s *{{ $serviceSessionImpl }}) Proposals() {{ $serviceProposalsInt }} {
    return s.proposals
}

func (s *{{ $serviceSessionImpl }}) State() {{ $serviceSessionState }} {
	return {{ $serviceSessionState }}(s.session.State())
}

func (s *{{ $serviceSessionImpl }}) Watch(f func({{ $serviceSessionState }})) {{ $serviceWatcherInt }} {
	return {{ $newServiceWatcher }}(s.session.Watch(func(state rsm.SessionState) {
	    f(SessionState(state))
	}))
}

var _ {{ $serviceSessionInt }} = &{{ $serviceSessionImpl }}{}

type {{ $serviceProposalsInt }} interface {
    {{- range .Primitive.Methods }}
    {{- if .Type.IsCommand }}
    {{- $proposalsInt := printf "%sProposals" .Name }}
    {{ .Name }}() {{ $proposalsInt }}
    {{- end }}
    {{- end }}
}

func {{ $newServiceProposals }}(commands rsm.Commands) {{ $serviceProposalsInt }} {
    return &{{ $serviceProposalsImpl }}{
        {{- range .Primitive.Methods }}
        {{- if .Type.IsCommand }}
        {{- $proposalsField := printf "%sProposals" ( .Name | toLowerCamel ) }}
        {{- $newProposals := printf "new%sProposals" .Name }}
        {{ $proposalsField }}: {{ $newProposals }}(commands),
        {{- end }}
        {{- end }}
    }
}

type {{ $serviceProposalsImpl }} struct {
    {{- range .Primitive.Methods }}
    {{- if .Type.IsCommand }}
    {{- $proposalsInt := printf "%sProposals" .Name }}
    {{- $proposalsField := printf "%sProposals" ( .Name | toLowerCamel ) }}
    {{ $proposalsField }} {{ $proposalsInt }}
    {{- end }}
    {{- end }}
}

{{- range .Primitive.Methods }}
{{- if .Type.IsCommand }}
{{- $proposalsInt := printf "%sProposals" .Name }}
{{- $proposalsField := printf "%sProposals" ( .Name | toLowerCamel ) }}
func (s *{{ $serviceProposalsImpl }}) {{ .Name }}() {{ $proposalsInt }} {
    return s.{{ $proposalsField }}
}
{{- end }}
{{- end }}

var _ {{ $serviceProposalsInt }} = &{{ $serviceProposalsImpl }}{}

type {{ $serviceProposalID }} uint64

type {{ $serviceProposalState }} int

const (
	ProposalComplete {{ $serviceProposalState }} = iota
	ProposalOpen
)

type {{ $serviceProposalInt }} interface {
    fmt.Stringer
	ID() {{ $serviceProposalID }}
	Session() {{ $serviceSessionInt }}
	State() {{ $serviceProposalState }}
	Watch(func({{ $serviceProposalState }})) {{ $serviceWatcherInt }}
}

func {{ $newServiceProposal }}(command rsm.Command) {{ $serviceProposalInt }} {
    return &{{ $serviceProposalImpl }}{
        command: command,
    }
}

type {{ $serviceProposalImpl }} struct {
    command rsm.Command
}

func (p *{{ $serviceProposalImpl }}) ID() {{ $serviceProposalID }} {
    return {{ $serviceProposalID }}(p.command.ID())
}

func (p *{{ $serviceProposalImpl }}) Session() {{ $serviceSessionInt }} {
    return {{ $newServiceSession }}(p.command.Session())
}

func (p *{{ $serviceProposalImpl }}) State() {{ $serviceProposalState }} {
    return {{ $serviceProposalState }}(p.command.State())
}

func (p *{{ $serviceProposalImpl }}) Watch(f func({{ $serviceProposalState }})) {{ $serviceWatcherInt }} {
	return {{ $newServiceWatcher }}(p.command.Watch(func(state rsm.CommandState) {
	    f({{ $serviceProposalState }}(state))
	}))
}

func (p *{{ $serviceProposalImpl }}) String() string {
    return fmt.Sprintf("ProposalID: %d, SessionID: %d", p.ID(), p.Session().ID())
}

var _ {{ $serviceProposalInt }} = &{{ $serviceProposalImpl }}{}

type {{ $serviceQueryInt }} interface {
    fmt.Stringer
	Session() {{ $serviceSessionInt }}
}

func {{ $newServiceQuery }}(query rsm.Query) {{ $serviceQueryInt }} {
    return &{{ $serviceQueryImpl }}{
        query: query,
    }
}

type {{ $serviceQueryImpl }} struct {
    query rsm.Query
}

func (p *{{ $serviceQueryImpl }}) Session() {{ $serviceSessionInt }} {
    return {{ $newServiceSession }}(p.query.Session())
}

func (p *{{ $serviceQueryImpl }}) String() string {
    return fmt.Sprintf("SessionID: %d", p.Session().ID())
}

var _ {{ $serviceQueryInt }} = &{{ $serviceQueryImpl }}{}

{{- range .Primitive.Methods }}
{{- if .Type.IsCommand }}
{{- $proposalsInt := printf "%sProposals" .Name }}
{{- $proposalsImpl := printf "%sProposals" ( .Name | toLowerCamel ) }}
{{- $newProposals := printf "new%sProposals" .Name }}
{{- $proposalInt := printf "%sProposal" .Name }}
{{- $proposalImpl := printf "%sProposal" ( .Name | toLowerCamel ) }}
{{- $newProposal := printf "new%sProposal" .Name }}
type {{ $proposalsInt }} interface {
    Get({{ $serviceProposalID }}) ({{ $proposalInt }}, bool)
    List() []{{ $proposalInt }}
}

func {{ $newProposals }}(commands rsm.Commands) {{ $proposalsInt }} {
    return &{{ $proposalsImpl }}{
        commands: commands,
    }
}

type {{ $proposalsImpl }} struct {
    commands rsm.Commands
}

func (p *{{ $proposalsImpl }}) Get(id {{ $serviceProposalID }}) ({{ $proposalInt }}, bool) {
    command, ok := p.commands.Get(rsm.CommandID(id))
    if !ok {
        return nil, false
    }
    proposal, err := {{ $newProposal }}(command)
    if err != nil {
        log.Error(err)
        return nil, false
    }
    return proposal, true
}

func (p *{{ $proposalsImpl }}) List() []{{ $proposalInt }} {
    commands := p.commands.List(rsm.OperationID({{ .ID }}))
    proposals := make([]{{ $proposalInt }}, len(commands))
    for i, command := range commands {
        proposal, err := {{ $newProposal }}(command)
        if err != nil {
            log.Error(err)
        } else {
            proposals[i] = proposal
        }
    }
    return proposals
}

var _ {{ $proposalsInt }} = &{{ $proposalsImpl }}{}

type {{ $proposalInt }} interface {
    {{ $serviceProposalInt }}
    Request() *{{ template "type" .Request.Type }}
    {{- if ( and .Response.IsUnary .Type.IsAsync ) }}
    Reply(*{{ template "type" .Response.Type }})
    Fail(error)
    {{- else if .Response.IsStream }}
    Notify(*{{ template "type" .Response.Type }})
    Close()
    {{- end }}
}

{{- if .Response.IsUnary }}
func {{ $newProposal }}(command rsm.Command) ({{ $proposalInt }}, error) {
    request := &{{ template "type" .Request.Type }}{}
    if err := proto.Unmarshal(command.Input(), request); err != nil {
        return nil, err
    }
    return &{{ $proposalImpl }}{
        {{ $serviceProposalInt }}: {{ $newServiceProposal }}(command),
        command: command,
        request: request,
    }, nil
}

type {{ $proposalImpl }} struct {
    {{ $serviceProposalInt }}
    command  rsm.Command
    request  *{{ template "type" .Request.Type }}
    {{- if .Type.IsAsync }}
    complete bool
    {{- end }}
}

func (p *{{ $proposalImpl }}) Request() *{{ template "type" .Request.Type }} {
    return p.request
}

{{- if .Type.IsAsync }}
func (p *{{ $proposalImpl }}) Reply(response *{{ template "type" .Response.Type }}) {
    if p.complete {
        return
    }
    log.Debugf("Sending {{ $proposalInt }} %s: %s", p, response)
    output, err := proto.Marshal(response)
    if err != nil {
        err = errors.NewInternal(err.Error())
        log.Errorf("Sending {{ $proposalInt }} %s response failed: %v", p, err)
        p.command.Output(nil, err)
    } else {
        p.command.Output(output, nil)
    }
    p.command.Close()
    p.complete = true
}

func (p *{{ $proposalImpl }}) Fail(err error) {
    if p.complete {
        return
    }
    log.Debugf("Failing {{ $proposalInt }} %s: %s", p, err)
    p.command.Output(nil, err)
    p.command.Close()
    p.complete = true
}
{{- end }}

func (p *{{ $proposalImpl }}) String() string {
    return fmt.Sprintf("ProposalID=%d, SessionID=%d", p.ID(), p.Session().ID())
}
{{- else if .Response.IsStream }}
func {{ $newProposal }}(command rsm.Command) ({{ $proposalInt }}, error) {
    request := &{{ template "type" .Request.Type }}{}
    if err := proto.Unmarshal(command.Input(), request); err != nil {
        return nil, err
    }
    return &{{ $proposalImpl }}{
        {{ $serviceProposalInt }}: {{ $newServiceProposal }}(command),
        command: command,
        request: request,
    }, nil
}

type {{ $proposalImpl }} struct {
    {{ $serviceProposalInt }}
    command rsm.Command
    request *{{ template "type" .Request.Type }}
    closed  bool
}

func (p *{{ $proposalImpl }}) Request() *{{ template "type" .Request.Type }} {
    return p.request
}

func (p *{{ $proposalImpl }}) Notify(response *{{ template "type" .Response.Type }}) {
    if p.closed {
        return
    }
    log.Debugf("Notifying {{ $proposalInt }} %s: %s", p, response)
    output, err := proto.Marshal(response)
    if err != nil {
        err = errors.NewInternal(err.Error())
        log.Errorf("Notifying {{ $proposalInt }} %s failed: %v", p, err)
        p.command.Output(nil, err)
        p.command.Close()
        p.closed = true
    } else {
        p.command.Output(output, nil)
    }
}

func (p *{{ $proposalImpl }}) Close() {
    p.command.Close()
    p.closed = true
}

func (p *{{ $proposalImpl }}) String() string {
    return fmt.Sprintf("ProposalID=%d, SessionID=%d", p.ID(), p.Session().ID())
}
{{- end }}

var _ {{ $proposalInt }} = &{{ $proposalImpl }}{}
{{- else if .Type.IsQuery }}
{{- $queryInt := printf "%sQuery" .Name }}
{{- $queryImpl := printf "%sQuery" ( .Name | toLowerCamel ) }}
{{- $newQuery := printf "new%sQuery" .Name }}

type {{ $queryInt }} interface {
    {{ $serviceQueryInt }}
    Request() *{{ template "type" .Request.Type }}
    {{- if .Response.IsStream }}
    Notify(*{{ template "type" .Response.Type }})
    Close()
    {{- end }}
}

{{- if .Response.IsUnary }}
func {{ $newQuery }}(query rsm.Query) ({{ $queryInt }}, error) {
    request := &{{ template "type" .Request.Type }}{}
    if err := proto.Unmarshal(query.Input(), request); err != nil {
        return nil, err
    }
    return &{{ $queryImpl }}{
        {{ $serviceQueryInt }}: {{ $newServiceQuery }}(query),
        query:   query,
        request: request,
    }, nil
}

type {{ $queryImpl }} struct {
    {{ $serviceQueryInt }}
    query   rsm.Query
    request *{{ template "type" .Request.Type }}
}

func (p *{{ $queryImpl }}) Request() *{{ template "type" .Request.Type }} {
    return p.request
}

func (p *{{ $queryImpl }}) String() string {
    return fmt.Sprintf("SessionID=%d", p.Session().ID())
}
{{- else if .Response.IsStream }}
func {{ $newQuery }}(query rsm.Query) ({{ $queryInt }}, error) {
    request := &{{ template "type" .Request.Type }}{}
    if err := proto.Unmarshal(query.Input(), request); err != nil {
        return nil, err
    }
    return &{{ $queryImpl }}{
        {{ $serviceQueryInt }}: {{ $newServiceQuery }}(query),
        query:   query,
        request: request,
    }, nil
}

type {{ $queryImpl }} struct {
    {{ $serviceQueryInt }}
    query   rsm.Query
    request *{{ template "type" .Request.Type }}
    closed  bool
}

func (p *{{ $queryImpl }}) Request() *{{ template "type" .Request.Type }} {
    return p.request
}

func (p *{{ $queryImpl }}) Notify(response *{{ template "type" .Response.Type }}) {
    if p.closed {
        return
    }
    log.Debugf("Notifying {{ $queryInt }} %s: %s", p, response)
    output, err := proto.Marshal(response)
    if err != nil {
        err = errors.NewInternal(err.Error())
        log.Errorf("Notifying {{ $queryInt }} %s failed: %v", p, err)
        p.query.Output(nil, err)
        p.query.Close()
        p.closed = true
    } else {
        p.query.Output(output, nil)
    }
}

func (p *{{ $queryImpl }}) Close() {
    p.query.Close()
    p.closed = true
}

func (p *{{ $queryImpl }}) String() string {
    return fmt.Sprintf("SessionID=%d", p.Session().ID())
}
{{- end }}

var _ {{ $queryInt }} = &{{ $queryImpl }}{}
{{- end }}
{{- end }}
