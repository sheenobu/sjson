package sjson

var endToken _endToken

type _endToken struct{}

func (n _endToken) Type() Type {
	return EndType
}

func (n _endToken) String() string {
	return "endToken"
}
