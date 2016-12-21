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
			types: make(map[string]constructor),
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

	var bl []func()
	i, bl, err = um.loop(ch, errChan, nil, nil, 0)

	for _, b := range bl {
		b()
	}

	return
}

func (um *Unmarshaller) inject(i interface{}, mt *sjson.MemberToken, childI interface{}) (err error) {
	tx := reflect.TypeOf(i)

	if tx.Kind() == reflect.Ptr {
		tx = tx.Elem()
	}

	for j := tx.NumField(); j != 0; j-- {
		f := tx.Field(j - 1)
		js := f.Tag.Get(um.TagName)
		keyName := strings.Split(js, ",")[0]
		if mt.Key == keyName {
			if st, ok := mt.Value.(sjson.SimpleToken); ok {
				tmp := reflect.Indirect(reflect.ValueOf(i)).Field(j - 1)
				if tmp.CanAddr() {
					if err = st.Unmarshal(tmp.Addr().Interface()); err != nil {
						return
					}
				}
			} else if childI != nil {
				reflect.Indirect(reflect.ValueOf(i)).Field(j - 1).Set(reflect.ValueOf(childI))
			} else if childI == nil {
				tmp := reflect.Indirect(reflect.ValueOf(i)).Field(j - 1)

				if tmp.Kind() == reflect.Ptr && tmp.IsNil() {
					t := reflect.Indirect(reflect.ValueOf(i)).Field(j - 1).Type()
					if t.Kind() == reflect.Ptr {
						t = t.Elem()
						childI = reflect.New(t).Interface()
					} else {
						childI = reflect.Zero(t).Interface()
					}

					if reflect.ValueOf(childI).IsValid() {
						reflect.Indirect(reflect.ValueOf(i)).Field(j - 1).Set(reflect.ValueOf(childI))
					}
				}
			}
		}
	}

	return
}

func (um *Unmarshaller) loop(ch chan sjson.Token, errChan chan error, mt *sjson.MemberToken, parent interface{}, loopLevel int) (i interface{}, backlog []func(), err error) {

	if parent != nil && mt != nil {
		for idx := reflect.TypeOf(parent).Elem().NumField() - 1; idx >= 0; idx-- {
			f := reflect.TypeOf(parent).Elem().Field(idx)
			ti := f.Tag.Get(um.TagName)
			t := strings.Split(ti, ",")[0]
			if mt.Key == t {
				i = reflect.Indirect(reflect.ValueOf(parent)).Field(idx).Interface()
			}
		}
	}

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

			if stackCount == 1 {
				if t.Type() == sjson.MemberType {
					mt := t.(*sjson.MemberToken)
					if mt.Key == um.FieldName {
						i = um.types[fmt.Sprintf("%s", mt.Value)]()
					}

					if mt.Value.Type() == sjson.ObjectType {
						var childI interface{}
						var bl []func()
						childI, bl, err = um.loop(ch, errChan, mt, i, loopLevel+1)
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
						backlog = append(backlog, bl...)
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

	return
}
