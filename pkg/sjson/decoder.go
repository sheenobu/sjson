package sjson

import (
	"bufio"
	"bytes"
	"io"

	"strings"

	"github.com/pkg/errors"
)

type stateFn func() (stateFn, bool, error)

// A Decoder is a JSON decoder
type Decoder struct {
	r *bufio.Reader

	currentState stateFn
	initialState stateFn

	tok Token
}

// NewDecoder creates a new decoder from the given reader
func NewDecoder(r io.Reader) *Decoder {
	d := &Decoder{
		r: bufio.NewReader(r),
	}
	d.initialState = d.beginState

	return d
}

// newDecoder creates a new decoder from the given reader
func newDecoder(r *bufio.Reader) *Decoder {
	d := &Decoder{
		r: r,
	}
	d.initialState = d.beginState

	return d
}

// Next gets the next token
func (dec *Decoder) Next() (Token, error) {
	var ok bool
	var err error

	dec.currentState = dec.initialState

	for {
		dec.currentState, ok, err = dec.currentState()
		if err != nil {
			return nil, err
		}
		if ok {
			return dec.tok, nil
		}
		if dec.currentState == nil {
			dec.currentState = dec.initialState
		}
	}
}

func (dec *Decoder) beginState() (stateFn, bool, error) {
	d, err := dec.r.Peek(1)
	if err != nil {
		return nil, false, err
	}

	switch {
	case d[0] == '"':
		return dec.stringState, false, nil
	case d[0] >= '0' && d[0] <= '9':
		return dec.numberState, false, nil
	case d[0] == 't' || d[0] == 'f':
		return dec.boolState, false, nil
	case d[0] == 'n':
		return dec.nullState, false, nil
	case d[0] == '{':
		dec.r.Discard(1)
		dec.tok = &objectToken{r: dec.r}
		return nil, true, nil
	case d[0] == ' ' || d[0] == '\t' || d[0] == '\n' || d[0] == '\r':
		dec.r.Discard(1)
		return dec.beginState, false, nil
	case d[0] == '[':
		dec.r.Discard(1)
		dec.tok = &arrayToken{r: dec.r}
		return nil, true, nil
	case d[0] == '}':
		return nil, true, nil
	default:
		return nil, false, errors.Errorf("Unexpected token %s", d)
	}
}

func (dec *Decoder) arrayState() (stateFn, bool, error) {
	d, err := dec.r.Peek(1)
	if err != nil {
		return nil, false, err
	}

	switch {
	case d[0] == '"':
		return dec.stringState, false, nil
	case d[0] >= '0' && d[0] <= '9':
		return dec.numberState, false, nil
	case d[0] == 't' || d[0] == 'f':
		return dec.boolState, false, nil
	case d[0] == 'n':
		return dec.nullState, false, nil
	case d[0] == '{':
		dec.initialState = dec.arrayState
		dec.r.Discard(1)
		dec.tok = &objectToken{r: dec.r}
		return nil, true, nil
	case d[0] == ' ' || d[0] == '\t' || d[0] == '\n' || d[0] == '\r':
		dec.r.Discard(1)
		return dec.arrayState, false, nil
	case d[0] == '[':
		dec.r.Discard(1)
		dec.tok = &arrayToken{r: dec.r}
		return nil, true, nil
	case d[0] == ']':
		dec.r.Discard(1)
		dec.tok = endToken
		return nil, true, nil
	case d[0] == '}':
		dec.tok = endToken
		return nil, true, nil
	case d[0] == ',':
		dec.r.Discard(1)
		return dec.arrayState, false, nil
	default:
		return nil, false, errors.Errorf("Unexpected token %s", d)
	}
}

func (dec *Decoder) boolState() (stateFn, bool, error) {
	var str bytes.Buffer

	for {
		b, _, err := dec.r.ReadRune()
		if err != nil && err != io.EOF {
			return nil, false, err
		}

		str.WriteRune(b)
		v := str.String()
		if strings.HasPrefix("true", v) || strings.HasPrefix("false", v) {
			if v == "true" || v == "false" {
				dec.tok = boolToken(v)
				return nil, true, nil
			}
			continue
		}

		return nil, false, errors.Errorf("Unexpected token %s", string(b))
	}
}

func (dec *Decoder) nullState() (stateFn, bool, error) {
	var str bytes.Buffer

	for {
		b, _, err := dec.r.ReadRune()
		if err != nil && err != io.EOF {
			return nil, false, err
		}

		str.WriteRune(b)
		v := str.String()
		if strings.HasPrefix("null", v) {
			if v == "null" {
				dec.tok = nullToken
				return nil, true, nil
			}
			continue
		}

		return nil, false, errors.Errorf("Unexpected token %s", string(b))
	}
}

func (dec *Decoder) stringState() (stateFn, bool, error) {
	dec.r.Discard(1)

	str, err := dec.r.ReadString('"')
	if err != nil {
		return nil, false, err
	}

	dec.tok = stringToken(str[0 : len(str)-1])

	return nil, true, nil
}

func (dec *Decoder) numberState() (stateFn, bool, error) {
	var str bytes.Buffer
	var isFloat bool

	for {
		b, _, err := dec.r.ReadRune()
		if err != nil && err != io.EOF {
			return nil, false, err
		}
		if (b < '0' || b > '9' || err == io.EOF) && b != '.' {
			dec.r.UnreadRune()
			if !isFloat {
				dec.tok = numberToken(str.String())
			} else {
				dec.tok = floatToken(str.String())
			}
			return nil, true, nil
		}

		if b == '.' {
			isFloat = true
		}
		str.WriteRune(b)
	}
}

func (dec *Decoder) objectBeginState() (stateFn, bool, error) {
	d, err := dec.r.Peek(1)
	if err != nil {
		return nil, false, err
	}

	switch {
	case d[0] == '"':
		return dec.objectMemberState, false, nil
	case d[0] == '}':
		dec.r.Discard(1)
		dec.tok = endToken
		return nil, true, nil
	case d[0] == ' ' || d[0] == '\t' || d[0] == '\n' || d[0] == '\r':
		dec.r.Discard(1)
		return dec.objectBeginState, false, nil
	default:
		return nil, false, errors.Errorf("Unexpected token %s (objectBeginState)", d)
	}
}

func (dec *Decoder) objectState() (stateFn, bool, error) {
	d, err := dec.r.Peek(1)
	if err != nil {
		return nil, false, err
	}

	switch {
	case d[0] == ',':
		dec.r.Discard(1)
		return dec.objectMemberState, false, nil
	case d[0] == '"':
		return dec.objectMemberState, false, nil
	case d[0] == '}':
		dec.r.Discard(1)
		dec.tok = endToken
		return nil, true, nil
	case d[0] == ' ' || d[0] == '\t' || d[0] == '\n' || d[0] == '\r':
		dec.r.Discard(1)
		return dec.objectState, false, nil
	default:
		return nil, false, errors.Errorf("Unexpected token %s  (objectBeginState)", d)
	}
}

func (dec *Decoder) objectMemberState() (stateFn, bool, error) {
	for {
		r, _, err := dec.r.ReadRune()
		if err != nil {
			return nil, false, err
		}
		if r == ' ' || r == '\n' || r == '\r' || r == '\t' {
			continue
		}
		if r != '"' {
			return nil, false, errors.Errorf("Unexpected token %s", string(r))
		}
		break
	}
	dec.r.UnreadRune()

	dec.r.Discard(1)

	key, err := dec.r.ReadString('"')
	if err != nil {
		return nil, false, err
	}

	// discard whitespace
	for {
		r, _, err := dec.r.ReadRune()
		if err != nil {
			return nil, false, err
		}
		if r == ' ' || r == '\n' || r == '\r' || r == '\t' {
			continue
		}
		if r != ':' {
			return nil, false, errors.Errorf("Unexpected token %s", string(r))
		}
		break
	}

	dr := newDecoder(dec.r)
	tok, err := dr.Next()
	if err != nil {
		return nil, false, err
	}

	dec.tok = &MemberToken{Key: key[0 : len(key)-1], Value: tok}
	dec.initialState = dec.objectState
	return dec.objectState, true, nil
}
