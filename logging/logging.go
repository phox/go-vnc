/*
Package logging provides common logging functionality.
*/
package logging

import (
	"runtime"
	"strings"
)

const (
	FlowLevel   = 2
	FnDeclLevel = 3
	ResultLevel = 4
	SpamLevel   = 5
)

// FnName returns the calling function name, e.g. "SomeFunction()". This is
// useful for logging a function name.
func FnName() string {
	pc := make([]uintptr, 10) // At least 1 entry needed.
	runtime.Callers(2, pc)
	name := runtime.FuncForPC(pc[0]).Name()
	return name[strings.LastIndex(name, ".")+1:] + "()"
}

// FnNameWithArgs returns the calling function name, with argument values,
// e.g. "SomeFunction(arg1, arg2)". This is useful for logging function calls
// with arguments.
func FnNameWithArgs(args ...string) string {
	pc := make([]uintptr, 10) // At least 1 entry needed.
	runtime.Callers(2, pc)
	name := runtime.FuncForPC(pc[0]).Name()
	argstr := ""
	for _, arg := range args {
		if len(argstr) > 0 {
			argstr += ", "
		}
		argstr += arg
	}
	return name[strings.LastIndex(name, ".")+1:] + "(" + argstr + ")"
}
