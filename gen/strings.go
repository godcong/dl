// Copyright (c) 2024 GodCong. All rights reserved.

// Package gen for Default Loader
package gen

import (
	"strings"
)

const convertPointer = "setup.Pointer(%s)"

func trimSide(src string, trim string) string {
	src = strings.TrimPrefix(src, trim)
	src = strings.TrimSuffix(src, trim)
	return src
}
