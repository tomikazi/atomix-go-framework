// Code generated by atomix-go-framework. DO NOT EDIT.
package {{ .Package.Name }}

{{ $server := printf "%sProxyServer" .Generator.Prefix }}
{{- $registry := printf "%sProxyRegistry" .Generator.Prefix }}
{{- $service := printf "%s.%sServer" .Primitive.Type.Package.Alias .Primitive.Type.Name }}

{{- define "type" }}{{ printf "%s.%s" .Package.Alias .Name }}{{ end }}

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

import (
	"context"
	"github.com/atomix/atomix-go-framework/pkg/atomix/errors"
	"github.com/atomix/atomix-go-framework/pkg/atomix/logging"
	"github.com/atomix/atomix-go-framework/pkg/atomix/driver/env"
	{{- $package := .Package }}
	{{- range .Imports }}
	{{ .Alias }} {{ .Path | quote }}
	{{- end }}
)

var log = logging.GetLogger("atomix", {{ .Primitive.Name | lower | quote }})

// New{{ $server }} creates a new {{ $server }}
func New{{ $server }}(registry *{{ $registry }}, env env.DriverEnv) {{ $service }} {
	return &{{ $server }}{
		registry: registry,
		env:      env,
	}
}

{{- $primitive := .Primitive }}
type {{ $server }} struct {
	registry *{{ $registry }}
	env      env.DriverEnv
}

{{- range .Primitive.Methods }}
{{- $method := . }}
{{ if and .Request.IsUnary .Response.IsUnary }}
func (s *{{ $server }}) {{ .Name }}(ctx context.Context, request *{{ template "type" .Request.Type }}) (*{{ template "type" .Response.Type }}, error) {
	if request.Headers.PrimitiveID.Namespace == "" {
		request.Headers.PrimitiveID.Namespace = s.env.Namespace
	}
	proxy, err := s.registry.GetProxy(request{{ template "field" .Request.Headers }}.PrimitiveID)
	if err != nil {
	    log.Warnf("{{ .Request.Type.Name }} %+v failed: %v", request, err)
	    if errors.IsNotFound(err) {
	        return nil, errors.NewUnavailable(err.Error())
	    }
		return nil, err
	}
	return proxy.{{ .Name }}(ctx, request)
}
{{ else if .Response.IsStream }}
func (s *{{ $server }}) {{ .Name }}(request *{{ template "type" .Request.Type }}, srv {{ template "type" $primitive.Type }}_{{ .Name }}Server) error {
	if request.Headers.PrimitiveID.Namespace == "" {
		request.Headers.PrimitiveID.Namespace = s.env.Namespace
	}
	proxy, err := s.registry.GetProxy(request{{ template "field" .Request.Headers }}.PrimitiveID)
	if err != nil {
	    log.Warnf("{{ .Request.Type.Name }} %+v failed: %v", request, err)
	    if errors.IsNotFound(err) {
	        return errors.NewUnavailable(err.Error())
	    }
		return err
	}
	return proxy.{{ .Name }}(request, srv)
}
{{ end }}
{{- end }}
