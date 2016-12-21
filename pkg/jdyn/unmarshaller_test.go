package jdyn

import (
	"fmt"
	"io"
	"strings"
	"testing"
)

type myClassWithPointerChild struct {
	Field1 string    `json:"field1"`
	Field2 string    `json:"field2"`
	Field3 *myClass2 `json:"field3"`
}

func (mc *myClassWithPointerChild) String() string {
	return fmt.Sprintf("myClass{%s, %s, %v}", mc.Field1, mc.Field2, mc.Field3)
}

type myClassWithInterface struct {
	Field1 string      `json:"field1"`
	Field2 string      `json:"field2"`
	Field3 interface{} `json:"field3"`
}

func (mc *myClassWithInterface) String() string {
	return fmt.Sprintf("myClass{%s, %s, %v}", mc.Field1, mc.Field2, mc.Field3)
}

type myClass2 struct {
	Field4 string `json:"field4"`
}

func (mc *myClass2) String() string {
	return fmt.Sprintf("myClass2{%s}", mc.Field4)
}

func TestReadClassTreeNestedField(t *testing.T) {

	var body = `
{
	"type": "myClass",
    "field2": "x2",
    "field1": "x1",
    "field3": {
        "field4": "x3"
    }
}`

	um := NewUnmarshaller("type")
	um.Register("myClass", &myClassWithPointerChild{})

	i, err := um.Unmarshal(strings.NewReader(body))
	if err != nil && err != io.EOF {
		t.Errorf("Error: %v\n", err)
	}

	if fmt.Sprintf("%s", i) != "myClass{x1, x2, myClass2{x3}}" {
		t.Errorf("Serialization failed, mismatch")
	}
}

func TestReadClassTreeNestedTypeField(t *testing.T) {
	var body = `
{
	"type": "myClass",
    "field2": "x2",
    "field1": "x1",
    "field3": {
        "field4": "x3",
		"type": "myClass2"
    }
}`

	um := NewUnmarshaller("type")
	um.Register("myClass", &myClassWithInterface{})
	um.Register("myClass2", &myClass2{})

	i, err := um.Unmarshal(strings.NewReader(body))
	if err != nil && err != io.EOF {
		t.Errorf("Error: %v\n", err)
	}

	if fmt.Sprintf("%s", i) != "myClass{x1, x2, myClass2{x3}}" {
		t.Errorf("Serialization failed, mismatch")
	}
}

func TestReadClassTreeNestedTypeFieldReverse(t *testing.T) {
	var body = `
{
    "field2": "x2",
    "field1": "x1",
    "field3": {
        "field4": "x3",
		"type": "myClass2"
    },
	"type": "myClass"
}`

	um := NewUnmarshaller("type")
	um.Register("myClass", &myClassWithInterface{})
	um.Register("myClass2", &myClass2{})

	i, err := um.Unmarshal(strings.NewReader(body))
	if err != nil && err != io.EOF {
		t.Errorf("Error: %v\n", err)
	}

	if fmt.Sprintf("%s", i) != "myClass{x1, x2, myClass2{x3}}" {
		t.Errorf("Serialization failed, mismatch")
	}
}

func TestReadClassTreeNestedNoTypeField(t *testing.T) {
	var body = `
{
    "field2": "x2",
    "field1": "x1",
    "field3": {
        "field4": "x3"
    },
	"type": "myClass"
}`

	um := NewUnmarshaller("type")
	um.Register("myClass", &myClassWithInterface{})
	um.Register("myClass2", &myClass2{})

	i, err := um.Unmarshal(strings.NewReader(body))
	if err != nil && err != io.EOF {
		t.Errorf("Error: %v\n", err)
	}

	if fmt.Sprintf("%s", i) != "myClass{x1, x2, <nil>}" {
		t.Errorf("Serialization failed, mismatch: '%v'", fmt.Sprintf("%s", i))
	}
}
