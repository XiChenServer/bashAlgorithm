package Map

import (
	"sync"
)

// MutexMap 对于map添加读写锁，可以保证并发安全
type MutexMap struct {
	v map[interface{}]interface{}
	sync.RWMutex
}

// NewMutexMap init
func NewMutexMap() *MutexMap {
	return &MutexMap{}
}
func (m *MutexMap) Put(key, value interface{}) {
	m.Lock()
	defer m.Unlock()
	m.v[key] = value

}
func (m *MutexMap) Get(key interface{}) interface{} {
	m.RLock()
	v := m.v[key]
	m.RUnlock()
	return v
}

func (m *MutexMap) Remove(key interface{}) {
	m.Lock()
	defer m.Unlock()
	delete(m.v, key)
}

//以区块化的形式进行加锁

type SharedMap struct {
}
