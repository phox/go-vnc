/*
Package logging provides common logging functionality.
*/
package logging

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/golang/glog"
)

const (
	FlowLevel      = 2
	FnDeclLevel    = 3
	ResultLevel    = 4
	SpamLevel      = 5
	CrazySpamLevel = 6
)

func V(level glog.Level) glog.Verbose {
	return glog.V(level)
}

// FnName returns the calling function name, e.g. "SomeFunction()". This is
// useful for logging a function name.
//
// Example:
//   if glog.V(logging.FnDeclLevel) {
//     glog.Info(logging.FnName())
//   }
func FnName() string {
	pc := make([]uintptr, 10) // At least 1 entry needed.
	runtime.Callers(2, pc)
	name := runtime.FuncForPC(pc[0]).Name()
	return name[strings.LastIndex(name, ".")+1:] + "()"
}

// FnNameWithArgs returns the calling function name, with argument values,
// e.g. "SomeFunction(arg1, arg2)". This is useful for logging function calls
// with arguments.
//
// Example:
//   if glog.V(logging.FnDeclLevel) {
//     glog.Info(logging.FnNameWithArgs(arg1, arg2, arg3))
//   }
func FnNameWithArgs(format string, args ...interface{}) string {
	pc := make([]uintptr, 10) // At least 1 entry needed.
	runtime.Callers(2, pc)
	name := runtime.FuncForPC(pc[0]).Name()
	a := []interface{}{name[strings.LastIndex(name, ".")+1:]}
	a = append(a, args...)
	return fmt.Sprintf("%s("+format+")", a...)
}
