package jdyn

import (
	"reflect"
	"sync"
)

type constructor func() interface{}

// Types contains the type mappings
type Types struct {
	types map[string]constructor
	tLock sync.RWMutex
}

func (t *Types) Register(name string, i interface{}) {

	t.tLock.Lock()
	t.types[name] = func() (ix interface{}) {
		tx := reflect.TypeOf(i).Elem()

		ixt := reflect.New(tx)

		for idx := tx.NumField() - 1; idx >= 0; idx-- {
			fieldType := tx.Field(idx).Type
			if fieldType.Kind() == reflect.Ptr {
				reflect.Indirect(ixt).Field(idx).Set(reflect.New(fieldType.Elem()))
			}
		}

		ix = ixt.Interface()

		return
	}
	t.tLock.Unlock()
}
