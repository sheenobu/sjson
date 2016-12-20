package sjson

import "bufio"

type objectToken struct {
	r   *bufio.Reader
	dec *Decoder
}

func (o *objectToken) Type() Type {
	return ObjectType
}

func (o *objectToken) Next() (Token, error) {
	if o.dec == nil {
		o.dec = newDecoder(o.r)
		o.dec.initialState = o.dec.objectBeginState
	}

	return o.dec.Next()
}

func (o *objectToken) String() string {
	return "objectToken"
}
