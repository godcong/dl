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
		log.Output(2, fmt.Sprintf("%v", a))
	}
}
