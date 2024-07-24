// Copyright (c) 2024 GodCong. All rights reserved.

// Package gen for Default Loader
package gen

import (
	"reflect"
)

type StructTag = reflect.StructTag

// NewStructTag creates a new StructTag by trimming the backticks from the input string.
func NewStructTag(v string) StructTag {
	return StructTag(trimSide(v, "`"))
}
