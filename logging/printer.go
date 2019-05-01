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

func Verbose(a ...interface{}) {
	if verboseFlagEnabled {
		fmt.Println(a...)
	}
}

func Verbosef(format string, a ...interface{}) {
	if verboseFlagEnabled {
		fmt.Printf(format+"\n", a...)
	}
}

func Inlinef(format string, a ...interface{}) {
	if !verboseFlagEnabled {
		fmt.Printf("\r"+format, a...)
	}
}
