package sjson

import "reflect"

type boolToken string

func (bt boolToken) Type() Type {
	return BoolType
}

func (bt boolToken) Unmarshal(b interface{}) error {
	switch string(bt) {
	case "true":
		reflect.Indirect(reflect.ValueOf(b)).SetBool(true)
	case "false":
		reflect.Indirect(reflect.ValueOf(b)).SetBool(false)
	}

	return nil
}

func (bt boolToken) String() string {
	return string(bt)
}
