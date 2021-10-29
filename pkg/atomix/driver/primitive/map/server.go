// Code generated by atomix-go-sdk. DO NOT EDIT.
package _map

import (
	"context"
	_map "github.com/atomix/atomix-api/go/atomix/primitive/map"
	"github.com/atomix/atomix-go-sdk/pkg/atomix/driver/env"
	"github.com/atomix/atomix-go-sdk/pkg/atomix/errors"
	"github.com/atomix/atomix-go-sdk/pkg/atomix/logging"
)

var log = logging.GetLogger("atomix", "map")

// NewProxyServer creates a new ProxyServer
func NewProxyServer(registry *ProxyRegistry, env env.DriverEnv) _map.MapServiceServer {
	return &ProxyServer{
		registry: registry,
		env:      env,
	}
}

type ProxyServer struct {
	registry *ProxyRegistry
	env      env.DriverEnv
}

func (s *ProxyServer) Size(ctx context.Context, request *_map.SizeRequest) (*_map.SizeResponse, error) {
	if request.Headers.PrimitiveID.Namespace == "" {
		request.Headers.PrimitiveID.Namespace = s.env.Namespace
	}
	proxy, err := s.registry.GetProxy(request.Headers.PrimitiveID)
	if err != nil {
		log.Warnf("SizeRequest %+v failed: %v", request, err)
		if errors.IsNotFound(err) {
			return nil, errors.NewUnavailable(err.Error())
		}
		return nil, err
	}
	return proxy.Size(ctx, request)
}

func (s *ProxyServer) Put(ctx context.Context, request *_map.PutRequest) (*_map.PutResponse, error) {
	if request.Headers.PrimitiveID.Namespace == "" {
		request.Headers.PrimitiveID.Namespace = s.env.Namespace
	}
	proxy, err := s.registry.GetProxy(request.Headers.PrimitiveID)
	if err != nil {
		log.Warnf("PutRequest %+v failed: %v", request, err)
		if errors.IsNotFound(err) {
			return nil, errors.NewUnavailable(err.Error())
		}
		return nil, err
	}
	return proxy.Put(ctx, request)
}

func (s *ProxyServer) Get(ctx context.Context, request *_map.GetRequest) (*_map.GetResponse, error) {
	if request.Headers.PrimitiveID.Namespace == "" {
		request.Headers.PrimitiveID.Namespace = s.env.Namespace
	}
	proxy, err := s.registry.GetProxy(request.Headers.PrimitiveID)
	if err != nil {
		log.Warnf("GetRequest %+v failed: %v", request, err)
		if errors.IsNotFound(err) {
			return nil, errors.NewUnavailable(err.Error())
		}
		return nil, err
	}
	return proxy.Get(ctx, request)
}

func (s *ProxyServer) Remove(ctx context.Context, request *_map.RemoveRequest) (*_map.RemoveResponse, error) {
	if request.Headers.PrimitiveID.Namespace == "" {
		request.Headers.PrimitiveID.Namespace = s.env.Namespace
	}
	proxy, err := s.registry.GetProxy(request.Headers.PrimitiveID)
	if err != nil {
		log.Warnf("RemoveRequest %+v failed: %v", request, err)
		if errors.IsNotFound(err) {
			return nil, errors.NewUnavailable(err.Error())
		}
		return nil, err
	}
	return proxy.Remove(ctx, request)
}

func (s *ProxyServer) Clear(ctx context.Context, request *_map.ClearRequest) (*_map.ClearResponse, error) {
	if request.Headers.PrimitiveID.Namespace == "" {
		request.Headers.PrimitiveID.Namespace = s.env.Namespace
	}
	proxy, err := s.registry.GetProxy(request.Headers.PrimitiveID)
	if err != nil {
		log.Warnf("ClearRequest %+v failed: %v", request, err)
		if errors.IsNotFound(err) {
			return nil, errors.NewUnavailable(err.Error())
		}
		return nil, err
	}
	return proxy.Clear(ctx, request)
}

func (s *ProxyServer) Events(request *_map.EventsRequest, srv _map.MapService_EventsServer) error {
	if request.Headers.PrimitiveID.Namespace == "" {
		request.Headers.PrimitiveID.Namespace = s.env.Namespace
	}
	proxy, err := s.registry.GetProxy(request.Headers.PrimitiveID)
	if err != nil {
		log.Warnf("EventsRequest %+v failed: %v", request, err)
		if errors.IsNotFound(err) {
			return errors.NewUnavailable(err.Error())
		}
		return err
	}
	return proxy.Events(request, srv)
}

func (s *ProxyServer) Entries(request *_map.EntriesRequest, srv _map.MapService_EntriesServer) error {
	if request.Headers.PrimitiveID.Namespace == "" {
		request.Headers.PrimitiveID.Namespace = s.env.Namespace
	}
	proxy, err := s.registry.GetProxy(request.Headers.PrimitiveID)
	if err != nil {
		log.Warnf("EntriesRequest %+v failed: %v", request, err)
		if errors.IsNotFound(err) {
			return errors.NewUnavailable(err.Error())
		}
		return err
	}
	return proxy.Entries(request, srv)
}
