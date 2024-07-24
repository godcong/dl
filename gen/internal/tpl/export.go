// Copyright (c) 2024 GodCong. All rights reserved.

// Package tpl for Default Loader
package tpl

import (
	_ "embed"
)

//go:embed default.go.tpl
var StructTemplate string
