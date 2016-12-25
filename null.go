package sjson

import "errors"

var nullToken _nullToken

type _nullToken struct{}

func (n _nullToken) Type() Type {
	return NullType
}

func (n _nullToken) Unmarshal(i interface{}) error {
	return errors.New("Unable to unmarshal a null JSON value")
}
