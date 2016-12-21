package jdyn

import (
	"reflect"
	"sync"
	"unicode"
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

			n := []rune(tx.Field(idx).Name)[0]
			if !unicode.IsLower(n) {
				if fieldType.Kind() == reflect.Ptr {
					reflect.Indirect(ixt).Field(idx).Set(reflect.New(fieldType.Elem()))
				}
			}
		}

		ix = ixt.Interface()

		return
	}
	t.tLock.Unlock()
}
