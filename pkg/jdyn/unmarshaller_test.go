package jdyn

import (
	"fmt"
	"io"
	"strings"
	"testing"
)

type myClass struct {
	Field1 string      `json:"field1"`
	Field2 string      `json:"field2"`
	Field3 interface{} `json:"field3"`
}

type myClass2 struct {
	Field4 string `json:"field4"`
}

func (mc *myClass) String() string {
	return fmt.Sprintf("myClass{%s, %s, %s}", mc.Field1, mc.Field2, mc.Field3)
}

func (mc *myClass2) String() string {
	return fmt.Sprintf("myClass2{%s}", mc.Field4)
}

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

func TestReadClassTree(t *testing.T) {
	um := NewUnmarshaller("type")
	um.Register("myClass", &myClass{})
	um.Register("myClass2", &myClass2{})

	i, err := um.Unmarshal(strings.NewReader(body))
	if err != nil && err != io.EOF {
		t.Errorf("Error: %v\n", err)
	}

	if fmt.Sprintf("%s", i) != "myClass{x1, x2, myClass2{x3}}" {
		t.Errorf("Serialization failed, mismatch")
	}

}
