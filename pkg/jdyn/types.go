package jdyn

import (
	"reflect"
	"sync"
)

// Types contains the type mappings
type Types struct {
	types map[string]reflect.Type
	tLock sync.RWMutex
}

func (t *Types) Register(name string, i interface{}) {
	t.tLock.Lock()
	t.types[name] = reflect.TypeOf(i).Elem()
	t.tLock.Unlock()
}
