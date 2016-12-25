package sjson

import (
	"fmt"
	"reflect"
)

type stringToken string

func (s stringToken) Type() Type {
	return StringType
}

func (s stringToken) Unmarshal(b interface{}) error {
	reflect.Indirect(reflect.ValueOf(b)).SetString(string(s))
	return nil
}

func (s stringToken) String() string {
	return fmt.Sprintf("%s", string(s))
}
