package indexedmap

import (
	primitiveapi "github.com/atomix/api/go/atomix/primitive"
	indexedmap "github.com/atomix/api/go/atomix/primitive/indexedmap"
	"github.com/atomix/go-framework/pkg/atomix/errors"
	"sync"
)

// NewProxyRegistry creates a new ProxyRegistry
func NewProxyRegistry() *ProxyRegistry {
	return &ProxyRegistry{
		proxies: make(map[primitiveapi.PrimitiveId]indexedmap.IndexedMapServiceServer),
	}
}

type ProxyRegistry struct {
	proxies map[primitiveapi.PrimitiveId]indexedmap.IndexedMapServiceServer
	mu      sync.RWMutex
}

func (r *ProxyRegistry) AddProxy(id primitiveapi.PrimitiveId, server indexedmap.IndexedMapServiceServer) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.proxies[id]; ok {
		return errors.NewAlreadyExists("proxy '%s' already exists", id)
	}
	r.proxies[id] = server
	return nil
}

func (r *ProxyRegistry) RemoveProxy(id primitiveapi.PrimitiveId) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.proxies[id]; !ok {
		return errors.NewNotFound("proxy '%s' not found", id)
	}
	delete(r.proxies, id)
	return nil
}

func (r *ProxyRegistry) GetProxy(id primitiveapi.PrimitiveId) (indexedmap.IndexedMapServiceServer, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	proxy, ok := r.proxies[id]
	if !ok {
		return nil, errors.NewNotFound("proxy '%s' not found", id)
	}
	return proxy, nil
}