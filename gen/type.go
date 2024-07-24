package gen

import (
	"reflect"
)

type StructTag = reflect.StructTag

// NewStructTag creates a new StructTag by trimming the backticks from the input string.
func NewStructTag(v string) StructTag {
	return StructTag(trimSide(v, "`"))
}
