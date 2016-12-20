package sjson

import "bufio"

type arrayToken struct {
	r   *bufio.Reader
	dec *Decoder
}

func (a *arrayToken) Type() Type {
	return ArrayType
}

func (a *arrayToken) Next() (Token, error) {
	if a.dec == nil {
		a.dec = newDecoder(a.r)
		a.dec.initialState = a.dec.arrayState
	}

	return a.dec.Next()
}

func (a *arrayToken) String() string {
	return "arrayToken"
}
