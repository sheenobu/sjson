package jdyn

import (
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/sheenobu/sjson/pkg/sjson"
)

// Unmarshaller is a JSON unmarshaller that supports
// dynamic type inference via a fieldname.
type Unmarshaller struct {
	Types
	FieldName string
	TagName   string
}

// NewUnmarshaller creates a new unmarshaller
func NewUnmarshaller(fieldName string) *Unmarshaller {
	return &Unmarshaller{
		Types: Types{
			types: make(map[string]reflect.Type),
		},
		FieldName: fieldName,
		TagName:   "json",
	}
}

// Unmarshal reads the JSON body and returns the type
func (um *Unmarshaller) Unmarshal(r io.Reader) (i interface{}, err error) {
	ch := make(chan sjson.Token)
	errChan := make(chan error, 1)

	go func() {
		err := sjson.ReadAll(r, ch)
		if err != nil {
			errChan <- err
		}
	}()

	i, err = um.loop(ch, errChan, 0)
	return
}

func (um *Unmarshaller) inject(i interface{}, mt *sjson.MemberToken, childI interface{}) (err error) {
	tx := reflect.TypeOf(i).Elem()

	for j := tx.NumField(); j != 0; j-- {
		f := tx.Field(j - 1)
		js := f.Tag.Get(um.TagName)
		keyName := strings.Split(js, ",")[0]
		if mt.Key == keyName {
			if st, ok := mt.Value.(sjson.SimpleToken); ok {
				if err = st.Unmarshal(reflect.Indirect(reflect.ValueOf(i)).Field(j - 1).Addr().Interface()); err != nil {
					return
				}
			} else if childI != nil {
				reflect.Indirect(reflect.ValueOf(i)).Field(j - 1).Set(reflect.ValueOf(childI))
			}
		}
	}

	return
}

func (um *Unmarshaller) loop(ch chan sjson.Token, errChan chan error, loopLevel int) (i interface{}, err error) {
	var backlog []func()

	var stackCount = 0

L:
	for {
		select {
		case err = <-errChan:
			break L
		case t := <-ch:
			if t.Type() == sjson.ObjectType {
				stackCount++
			}
			if t.Type() == sjson.EndType {
				stackCount--
				if stackCount == 0 {
					break L
				}
			}

			if i == nil && stackCount == 1 {
				if t.Type() == sjson.MemberType {
					mt := t.(*sjson.MemberToken)
					if mt.Key == um.FieldName {
						i = reflect.New(um.types[fmt.Sprintf("%s", mt.Value)]).Interface()
					}

					if mt.Value.Type() == sjson.ObjectType {
						var childI interface{}
						childI, err = um.loop(ch, errChan, loopLevel+1)
						if err != nil {
							break L
						}

						backlog = append(backlog, func(t sjson.Token) func() {
							return func() {
								if t.Type() == sjson.MemberType && i != nil {
									mt := t.(*sjson.MemberToken)
									um.inject(i, mt, childI)
								}
							}
						}(t))
					}
				}

				backlog = append(backlog, func(t sjson.Token) func() {
					return func() {
						if t.Type() == sjson.MemberType && i != nil {
							mt := t.(*sjson.MemberToken)
							um.inject(i, mt, nil)
						}
					}
				}(t))
			}

			if t.Type() == sjson.MemberType && i != nil {
				mt := t.(*sjson.MemberToken)
				um.inject(i, mt, nil)
			}
		}
	}

	for _, t := range backlog {
		t()
	}

	return
}
