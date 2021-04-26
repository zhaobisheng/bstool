package SyncMap

import (
	"sync"
)

type BsStruct struct {
	sync.RWMutex
	Val map[interface{}]interface{} //interface{}
}

func NewMap() *BsStruct {
	return &BsStruct{Val: make(map[interface{}]interface{})}
}

func GetMapVal(bs *BsStruct) interface{} {
	bs.RLock()
	defer bs.RUnlock()
	return bs.Val
}

func (m *BsStruct) Get(key interface{}) interface{} {
	m.Lock()
	defer m.Unlock()
	return m.Val[key]
}

func (m *BsStruct) Set(key interface{}, value interface{}) {
	m.Lock()
	defer m.Unlock()
	m.Val[key] = value
}

func (m *BsStruct) Remove(key interface{}) interface{} {
	m.Lock()
	defer m.Unlock()
	if m.Val[key] != nil {
		delete(m.Val, key)
		return key
	}
	return key
}

func (m *BsStruct) Range(fun func(interface{}, interface{}) bool) {
	for key := range m.Val {
		val := m.Get(key)
		fun(key, val)
	}
}
