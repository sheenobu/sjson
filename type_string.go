// Code generated by "stringer -type=Type"; DO NOT EDIT

package sjson

import "fmt"

const _Type_name = "NumberTypeStringTypeBoolTypeNullTypeObjectTypeArrayTypeMemberTypeEndType"

var _Type_index = [...]uint8{0, 10, 20, 28, 36, 46, 55, 65, 72}

func (i Type) String() string {
	i -= 1
	if i < 0 || i >= Type(len(_Type_index)-1) {
		return fmt.Sprintf("Type(%d)", i+1)
	}
	return _Type_name[_Type_index[i]:_Type_index[i+1]]
}
