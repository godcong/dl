// Copyright (c) 2024 GodCong. All rights reserved.

// Package gen for Default Loader
package gen

import (
	"fmt"
	"log"
)

var debug = false

func Debug() {
	debug = true
}

func debugPrint(a ...any) {
	if debug {
		log.SetFlags(log.Lshortfile | log.LstdFlags)
		_ = log.Output(2, fmt.Sprintf("%v", a))
	}
}
