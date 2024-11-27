package kvcacher

import (
	"context"
	"sync"
)

// NewMock creates new mock cacher
func NewMock(prefix string) *Mock {
	return &Mock{
		cache:  make(map[string][]byte),
		prefix: prefix,
	}
}

// Mock represents in-memory cacher
type Mock struct {
	sync.Mutex
	cache  map[string][]byte
	prefix string
}

// Get item from cache
func (m *Mock) Get(ctx context.Context, key string) ([]byte, error) {
	m.Lock()
	defer m.Unlock()
	b, ok := m.cache[m.prefix+key]
	if !ok {
		return nil, nil
	}
	return b, nil
}

// Set item into cache
func (m *Mock) Set(ctx context.Context, key string, b []byte) error {
	m.Lock()
	defer m.Unlock()
	m.cache[m.prefix+key] = b
	return nil
}

// Clear the cache
func (m *Mock) Clear(ctx context.Context) error {
	m.Lock()
	defer m.Unlock()
	m.cache = make(map[string][]byte)
	return nil
}
