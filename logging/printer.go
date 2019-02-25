package logging

import (
	"fmt"
)

var (
	verboseFlagEnabled = false
)

func SetVerbose(b bool) {
	verboseFlagEnabled = b
}

func Println(a ...interface{}) {
	if verboseFlagEnabled {
		fmt.Println(a...)
	}
}

func Printlnf(format string, a ...interface{}) {
	if verboseFlagEnabled {
		fmt.Printf(format+"\n", a...)
	}
}
