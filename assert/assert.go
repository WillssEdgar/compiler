package assert

import "fmt"

func assert(condition bool, msg string, args ...interface{}) {
	if !condition {
		panic(fmt.Sprintf("Assertion failed: "+msg, args...))
	}
}
