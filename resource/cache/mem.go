package cache

import (
	"sync"
)

type MemMap struct {
	mux  sync.RWMutex
	data map[interface{}]interface{}
}

func NewMemMap() ICache {
	m := new(MemMap)
	m.data = make(map[interface{}]interface{})

	return m
}

func (m *MemMap) Delete(key interface{}) {
	m.mux.Lock()
	defer m.mux.Unlock()
	delete(m.data, key)
}

func (m *MemMap) Len() int {
	m.mux.RLock()
	defer m.mux.RUnlock()
	return len(m.data)
}

func (m *MemMap) Get(key interface{}) (value interface{}, ok bool) {
	m.mux.RLock()
	defer m.mux.RUnlock()
	v, ok := m.data[key]
	return v, ok
}

func (m *MemMap) Put(key interface{}, value interface{}) {
	m.mux.Lock()
	defer m.mux.Unlock()
	m.data[key] = value
}
