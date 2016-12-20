package sjson

/*
   A JSON value MUST be an object, array, number, or string, or one of
   the following three literal names: false, true, null
*/

// Type represents the type of objects
type Type int64

//go:generate stringer -type=Type

// The type enumerations
const (
	// simple types
	NumberType Type = 1
	StringType Type = 2
	BoolType   Type = 3
	NullType   Type = 4

	// complex types
	ObjectType Type = 5
	ArrayType  Type = 6

	// special types
	MemberType Type = 7
	EndType    Type = 8
)

// A Token is one of the primary objects
type Token interface {
	Type() Type
}

// A SimpleToken is a token that can be unmarshalled into a pointer
type SimpleToken interface {
	Unmarshal(b interface{}) error
}

// A ComplexToken is a token that is made up of other tokens
type ComplexToken interface {
	Next() (Token, error)
}
