{{- $serviceType := printf "%sServiceType" .Generator.Prefix }}

{{- define "type" -}}
{{- if .Package.Import -}}
{{- printf "%s.%s" .Package.Alias .Name -}}
{{- else -}}
{{- .Name -}}
{{- end -}}
{{- end -}}

package {{ .Package.Name }}

import (
	"context"
	{{- $package := .Package }}
	{{- range .Imports }}
	{{ .Alias }} {{ .Path | quote }}
	{{- end }}
	"github.com/atomix/go-framework/pkg/atomix/storage/protocol/gossip"
	"github.com/atomix/go-framework/pkg/atomix/logging"
	"github.com/atomix/go-framework/pkg/atomix/time"
)

var log = logging.GetLogger("atomix", "protocol", "gossip", {{ .Primitive.Name | lower | quote }})

const {{ $serviceType }} gossip.ServiceType = {{ .Primitive.Name | quote }}

// RegisterService registers the service on the given node
func RegisterService(node *gossip.Node) {
	node.RegisterService(ServiceType, func(ctx context.Context, serviceID gossip.ServiceId, partition *gossip.Partition, clock time.Clock) (gossip.Service, error) {
		protocol, err := newGossipProtocol(serviceID, partition, clock)
		if err != nil {
			return nil, err
		}
		service, err := newService(protocol)
		if err != nil {
			return nil, err
		}
		engine := newGossipEngine(protocol)
		go engine.start()
		return service, nil
	})
}

var newService func(protocol GossipProtocol) (Service, error)

func registerService(f func(protocol GossipProtocol) (Service, error)) {
	newService = f
}

type Service interface {
	gossip.Service
    {{- range .Primitive.Methods }}

    {{- $comments := split .Comment "\n" }}
    {{- range $comment := $comments }}
    {{- if $comment }}
    // {{ $comment | trim }}
    {{- end }}
    {{- end }}

    {{- if or .Type.IsAsync .Response.IsStream }}
    {{ .Name }}(context.Context, *{{ template "type" .Request.Type }}, chan<- {{ template "type" .Response.Type }}) error
    {{- else }}
    {{ .Name }}(context.Context, *{{ template "type" .Request.Type }}) (*{{ template "type" .Response.Type }}, error)
    {{- end }}
    {{- end }}
}