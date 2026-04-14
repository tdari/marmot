package connection

import (
	"fmt"
	"sync"
)

type Registry struct {
	sources map[string]Source
	meta    map[string]ConnectionTypeMeta
	mu      sync.RWMutex
}

var (
	reg     *Registry
	regOnce sync.Once
)

func GetRegistry() *Registry {
	regOnce.Do(func() {
		reg = &Registry{
			sources: make(map[string]Source),
			meta:    make(map[string]ConnectionTypeMeta),
		}
	})
	return reg
}

func (r *Registry) Register(meta ConnectionTypeMeta, source Source) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.sources[meta.ID]; exists {
		return fmt.Errorf("connection type %s already registered", meta.ID)
	}
	r.sources[meta.ID] = source
	r.meta[meta.ID] = meta
	return nil
}

func (r *Registry) Get(id string) (Source, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	s, ok := r.sources[id]
	if !ok {
		return nil, fmt.Errorf("connection type %s not found", id)
	}
	return s, nil
}

func (r *Registry) GetMeta(id string) (ConnectionTypeMeta, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	m, ok := r.meta[id]
	if !ok {
		return ConnectionTypeMeta{}, fmt.Errorf("connection meta %s not found", id)
	}
	return m, nil
}

func (r *Registry) List() []ConnectionTypeMeta {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]ConnectionTypeMeta, 0, len(r.meta))
	for _, m := range r.meta {
		result = append(result, m)
	}
	return result
}
