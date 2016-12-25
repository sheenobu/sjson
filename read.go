package sjson

import "io"

// ReadAll reads the entire body and sends each token to the channel
func ReadAll(r io.Reader, ch chan<- Token) (err error) {
	dec := NewDecoder(r)
	err = readAll(dec, ch)

	// hide EOFs
	if err == io.EOF {
		err = nil
	}

	return
}

func readAll(ct ComplexToken, ch chan<- Token) error {
	t, err := ct.Next()

	if err != nil || t == nil {
		return err
	}

	ch <- t

	if t.Type() == EndType {
		return nil
	}

	// simple type
	if t.Type() < 5 {
		return readAll(ct, ch)
	}

	// complex type
	if t.Type() < 7 {
		err := readAll(t.(ComplexToken), ch)
		if err != nil {
			return err
		}
	}

	if t.Type() == MemberType {

		// member type
		mt := t.(*MemberToken)
		t2 := mt.Value

		// simple type
		if t2.Type() < 5 {
			//ch <- t2
			return readAll(ct, ch)
		}

		// complex type
		if t2.Type() < 7 {
			ch <- t2
			err := readAll(t2.(ComplexToken), ch)
			if err != nil {
				return err
			}
		}
	}

	return readAll(ct, ch)
}
