package sjson

import "fmt"

// A MemberToken is a token that is a member of an object
type MemberToken struct {
	Key   string
	Value Token
}

// Type returns the type of the member token
func (mt *MemberToken) Type() Type {
	return MemberType
}

func (mt *MemberToken) String() string {
	return fmt.Sprintf("%s=%s", mt.Key, mt.Value)
}
