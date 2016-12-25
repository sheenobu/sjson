package sjson

import (
	"bufio"
	"fmt"
	"io"
	"testing"

	"strings"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDecoder(t *testing.T) {

	Convey(`", should fail to parse"`, t, func() {
		r := strings.NewReader(`,`)
		dec := NewDecoder(r)

		t, err := dec.Next()
		So(err, ShouldNotBeNil)
		So(t, ShouldBeNil)
	})

	Convey(`{,} should fail to parse"`, t, func() {
		r := strings.NewReader(`{,}`)
		dec := NewDecoder(r)

		t, err := dec.Next()
		So(err, ShouldBeNil)
		So(t, ShouldNotBeNil)

		ct := t.(ComplexToken)
		t, err = ct.Next()
		So(err, ShouldNotBeNil)
		So(t, ShouldBeNil)
	})

	Convey(`{"key":"value",} should fail to parse"`, t, func() {
		r := strings.NewReader(`{"key":"value",}`)
		dec := NewDecoder(r)

		t, err := dec.Next()
		So(err, ShouldBeNil)
		So(t, ShouldNotBeNil)

		ct := t.(ComplexToken)
		t, err = ct.Next()
		So(err, ShouldBeNil)
		So(t, ShouldNotBeNil)

		t, err = ct.Next()
		So(err, ShouldNotBeNil)
		So(t, ShouldBeNil)
	})

	Convey(`{"key","value"} should fail to parse"`, t, func() {
		r := strings.NewReader(`{"key","value"}`)
		dec := NewDecoder(r)

		t, err := dec.Next()
		So(err, ShouldBeNil)
		So(t, ShouldNotBeNil)

		ct := t.(ComplexToken)
		t, err = ct.Next()
		So(err, ShouldNotBeNil)
		So(t, ShouldBeNil)

	})

	Convey("A JSON string should decode to a String token", t, func() {
		r := strings.NewReader(`"hello"`)
		dec := NewDecoder(r)

		t, err := dec.Next()
		So(err, ShouldBeNil)
		So(t, ShouldNotBeNil)
		So(t.Type(), ShouldEqual, StringType)

		var b string
		err = t.(SimpleToken).Unmarshal(&b)
		So(err, ShouldBeNil)
		So(b, ShouldEqual, "hello")

		_, err = dec.Next()
		So(err, ShouldEqual, io.EOF)
	})

	Convey("A JSON int number should decode to a Number token", t, func() {
		r := strings.NewReader(`42`)
		dec := NewDecoder(r)

		t, err := dec.Next()
		So(err, ShouldBeNil)
		So(t, ShouldNotBeNil)
		So(t.Type(), ShouldEqual, NumberType)

		var b int64
		err = t.(SimpleToken).Unmarshal(&b)
		So(err, ShouldBeNil)
		So(b, ShouldEqual, 42)

		_, err = dec.Next()
		So(err, ShouldEqual, io.EOF)
	})

	Convey("A JSON float number should decode to a Number token", t, func() {
		r := strings.NewReader(`42.3`)
		dec := NewDecoder(r)

		t, err := dec.Next()
		So(err, ShouldBeNil)
		So(t, ShouldNotBeNil)
		So(t.Type(), ShouldEqual, NumberType)

		var b float64
		err = t.(SimpleToken).Unmarshal(&b)
		So(err, ShouldBeNil)
		So(b, ShouldEqual, 42.3)

		_, err = dec.Next()
		So(err, ShouldEqual, io.EOF)
	})

	Convey("A JSON float64 number should Marshal into a float32 pointer", t, func() {
		r := strings.NewReader(`42.3`)
		dec := NewDecoder(r)

		t, err := dec.Next()
		So(err, ShouldBeNil)
		So(t, ShouldNotBeNil)
		So(t.Type(), ShouldEqual, NumberType)

		var b float32
		err = t.(SimpleToken).Unmarshal(&b)
		So(err, ShouldBeNil)
		So(b, ShouldEqual, 42.3)

		_, err = dec.Next()
		So(err, ShouldEqual, io.EOF)
	})

	Convey("A JSON true should decode into a Bool token", t, func() {
		r := strings.NewReader(`true`)
		dec := NewDecoder(r)

		t, err := dec.Next()
		So(err, ShouldBeNil)
		So(t, ShouldNotBeNil)
		So(t.Type(), ShouldEqual, BoolType)

		var b bool
		err = t.(SimpleToken).Unmarshal(&b)
		So(err, ShouldBeNil)
		So(b, ShouldEqual, true)

		_, err = dec.Next()
		So(err, ShouldEqual, io.EOF)
	})

	Convey("A JSON false should decode into a Bool token", t, func() {
		r := strings.NewReader(`false`)
		dec := NewDecoder(r)

		t, err := dec.Next()
		So(err, ShouldBeNil)
		So(t, ShouldNotBeNil)
		So(t.Type(), ShouldEqual, BoolType)

		var b bool
		err = t.(SimpleToken).Unmarshal(&b)
		So(err, ShouldBeNil)
		So(b, ShouldEqual, false)

		_, err = dec.Next()
		So(err, ShouldEqual, io.EOF)
	})

	Convey("A JSON null should decode into a Null token", t, func() {
		r := strings.NewReader(`null`)
		dec := NewDecoder(r)

		t, err := dec.Next()
		So(err, ShouldBeNil)
		So(t, ShouldNotBeNil)
		So(t.Type(), ShouldEqual, NullType)

		_, err = dec.Next()
		So(err, ShouldEqual, io.EOF)
	})

	Convey("A JSON Object should decode into an Object token", t, func() {
		r := strings.NewReader(`{}`)
		dec := NewDecoder(r)

		t, err := dec.Next()
		So(err, ShouldBeNil)
		So(t, ShouldNotBeNil)
		So(t.Type(), ShouldEqual, ObjectType)

		ctoken := t.(ComplexToken)
		t, err = ctoken.Next()

		So(err, ShouldEqual, nil)
		So(t.Type(), ShouldEqual, EndType)
	})

	Convey("A JSON Object with keys should decode into an Object token", t, func() {
		r := strings.NewReader(`{"key":"val"}`)
		dec := NewDecoder(r)

		t, err := dec.Next()
		So(err, ShouldBeNil)
		So(t, ShouldNotBeNil)
		So(t.Type(), ShouldEqual, ObjectType)

		ctoken := t.(ComplexToken)
		t, err = ctoken.Next()

		So(err, ShouldBeNil)
		So(t, ShouldNotBeNil)
		So(t.Type(), ShouldEqual, MemberType)

		mt := t.(*MemberToken)
		So(mt.Key, ShouldEqual, "key")
		So(mt.Value.Type(), ShouldEqual, StringType)
	})

	Convey("A JSON Object with keys and whitespace should decode into an Object token", t, func() {
		r := strings.NewReader(`{
			"key":
			"val"}`)
		dec := NewDecoder(r)

		t, err := dec.Next()
		So(err, ShouldBeNil)
		So(t, ShouldNotBeNil)
		So(t.Type(), ShouldEqual, ObjectType)

		ctoken := t.(ComplexToken)
		t, err = ctoken.Next()

		So(err, ShouldBeNil)
		So(t, ShouldNotBeNil)
		So(t.Type(), ShouldEqual, MemberType)

		mt := t.(*MemberToken)
		So(mt.Key, ShouldEqual, "key")
		So(mt.Value.Type(), ShouldEqual, StringType)
	})

	Convey("A JSON Object with multiple keys and whitespace should decode into an Object token", t, func() {
		r := strings.NewReader(`{
			"key":
			"val",
			"key2": "val2",
			"key3": {}
		}`)
		dec := NewDecoder(r)

		t, err := dec.Next()
		So(err, ShouldBeNil)
		So(t, ShouldNotBeNil)
		So(t.Type(), ShouldEqual, ObjectType)

		ctoken := t.(ComplexToken)

		// read key 1
		t, err = ctoken.Next()

		So(err, ShouldBeNil)
		So(t, ShouldNotBeNil)
		So(t.Type(), ShouldEqual, MemberType)

		mt := t.(*MemberToken)
		So(mt.Key, ShouldEqual, "key")
		So(mt.Value.Type(), ShouldEqual, StringType)
		So(toString(mt.Value), ShouldEqual, "val")

		// read key 2
		t, err = ctoken.Next()

		So(err, ShouldBeNil)
		So(t, ShouldNotBeNil)
		So(t.Type(), ShouldEqual, MemberType)

		mt = t.(*MemberToken)
		So(mt.Key, ShouldEqual, "key2")
		So(mt.Value.Type(), ShouldEqual, StringType)
		So(toString(mt.Value), ShouldEqual, "val2")

		// read key 3
		t, err = ctoken.Next()

		So(err, ShouldBeNil)
		So(t, ShouldNotBeNil)
		So(t.Type(), ShouldEqual, MemberType)

		mt = t.(*MemberToken)
		So(mt.Key, ShouldEqual, "key3")
		So(mt.Value.Type(), ShouldEqual, ObjectType)

		ct2 := mt.Value.(ComplexToken)
		t, err = ct2.Next()

		So(err, ShouldBeNil)
		So(t.Type(), ShouldEqual, EndType)
	})

	Convey("A JSON Array should decode into an Array token", t, func() {
		r := strings.NewReader(`[]`)
		dec := NewDecoder(r)

		t, err := dec.Next()
		So(err, ShouldBeNil)
		So(t, ShouldNotBeNil)
		So(t.Type(), ShouldEqual, ArrayType)
	})

	Convey("A JSON Array of numbers should decode into an Array token", t, func() {
		r := strings.NewReader(`[1,2]`)
		dec := NewDecoder(r)

		t, err := dec.Next()
		So(err, ShouldBeNil)
		So(t, ShouldNotBeNil)
		So(t.Type(), ShouldEqual, ArrayType)

		ct := t.(ComplexToken)

		t, err = ct.Next()
		So(err, ShouldBeNil)
		So(t, ShouldNotBeNil)
		So(t.Type(), ShouldEqual, NumberType)

		t, err = ct.Next()
		So(err, ShouldBeNil)
		So(t, ShouldNotBeNil)
		So(t.Type(), ShouldEqual, NumberType)
	})

	Convey("A JSON Array of numbers and strings should decode into an Array token", t, func() {
		r := strings.NewReader(`["1",2]`)
		dec := NewDecoder(r)

		t, err := dec.Next()
		So(err, ShouldBeNil)
		So(t, ShouldNotBeNil)
		So(t.Type(), ShouldEqual, ArrayType)

		ct := t.(ComplexToken)

		t, err = ct.Next()
		So(err, ShouldBeNil)
		So(t, ShouldNotBeNil)
		So(t.Type(), ShouldEqual, StringType)

		t, err = ct.Next()
		So(err, ShouldBeNil)
		So(t, ShouldNotBeNil)
		So(t.Type(), ShouldEqual, NumberType)
	})

	Convey("A JSON Array of numbers and strings and objects should decode into an Array token", t, func() {

		var r = bufio.NewReader(strings.NewReader(`[{"key":  12 },   1 ]`))
		dec := NewDecoder(r)

		t, err := dec.Next()
		So(err, ShouldBeNil)
		So(t, ShouldNotBeNil)
		So(t.Type(), ShouldEqual, ArrayType)

		ct := t.(ComplexToken)

		t, err = ct.Next()
		So(err, ShouldBeNil)
		So(t, ShouldNotBeNil)
		So(t.Type(), ShouldEqual, ObjectType)

		// read key 1
		ct2 := t.(ComplexToken)
		t, err = ct2.Next()
		So(err, ShouldBeNil)
		So(t, ShouldNotBeNil)
		So(t.Type(), ShouldEqual, MemberType)

		// inspect key 1
		mt := t.(*MemberToken)
		So(mt.Key, ShouldEqual, "key")
		So(mt.Value.Type(), ShouldEqual, NumberType)
		So(toString(mt.Value), ShouldEqual, "12")

		t, err = ct2.Next()
		So(err, ShouldBeNil)
		So(t, ShouldNotBeNil)
		So(t.Type(), ShouldEqual, EndType)

		t, err = ct.Next()
		So(err, ShouldBeNil)
		So(t, ShouldNotBeNil)
		So(t.Type(), ShouldEqual, NumberType)

	})

	Convey("A JSON Object with an array", t, func() {

		var r io.Reader = strings.NewReader(
			`
{
	"key": "value",
	"key2": [1,3,4]
}
`)

		dec := NewDecoder(r)

		t, err := dec.Next()
		So(err, ShouldBeNil)
		So(t, ShouldNotBeNil)
		So(t.Type(), ShouldEqual, ObjectType)

		ct := t.(ComplexToken)

		t, err = ct.Next()
		So(err, ShouldBeNil)
		So(t, ShouldNotBeNil)
		So(t.Type(), ShouldEqual, MemberType)

		t, err = ct.Next()
		So(err, ShouldBeNil)
		So(t, ShouldNotBeNil)
		So(t.Type(), ShouldEqual, MemberType)

		mt := t.(*MemberToken)
		So(mt.Value.Type(), ShouldEqual, ArrayType)

		ct2 := mt.Value.(ComplexToken)

		t, err = ct2.Next()
		So(err, ShouldBeNil)
		So(t, ShouldNotBeNil)
		So(t.Type(), ShouldEqual, NumberType)

		t, err = ct2.Next()
		So(err, ShouldBeNil)
		So(t, ShouldNotBeNil)
		So(t.Type(), ShouldEqual, NumberType)

		t, err = ct2.Next()
		So(err, ShouldBeNil)
		So(t, ShouldNotBeNil)
		So(t.Type(), ShouldEqual, NumberType)

		t, err = ct2.Next()
		So(err, ShouldBeNil)
		So(t.Type(), ShouldEqual, EndType)
	})
}

func toString(t Token) string {
	return fmt.Sprintf("%s", t)
}
