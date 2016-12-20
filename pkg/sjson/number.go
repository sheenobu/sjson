package sjson

import (
	"fmt"
	"reflect"
	"strconv"
)

type numberToken string

func (n numberToken) Type() Type {
	return NumberType
}

func (n numberToken) Unmarshal(b interface{}) error {
	i, err := strconv.ParseInt(string(n), 10, 64)
	if err != nil {
		return err
	}

	reflect.Indirect(reflect.ValueOf(b)).SetInt(i)

	return nil
}

func (n numberToken) String() string {
	return fmt.Sprintf("%s", string(n))
}

type floatToken string

func (n floatToken) Type() Type {
	return NumberType
}

func (n floatToken) Unmarshal(b interface{}) error {
	i, err := strconv.ParseFloat(string(n), 64)
	if err != nil {
		return err
	}

	reflect.Indirect(reflect.ValueOf(b)).SetFloat(i)

	return nil
}

func (n floatToken) String() string {
	return fmt.Sprintf("%s", string(n))
}
