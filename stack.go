package errors

import (
	"fmt"
	"io"
	"runtime"
)

/*
Imported from https://github.com/pkg/errors/blob/master/stack.go
*/

type StackTrace []StackFrame

// Initialize initializes the StackTrace with the callers of the current func
func (st *StackTrace) Initialize() {
	const depth = 32
	var counters [depth]uintptr
	count := runtime.Callers(3, counters[:]) // skip extern.go, this func, Error.func
	*st = make(StackTrace, count)
	for i := 0; i < count; i++ {
		(*st)[i] = StackFrame(counters[i])
	}
}

// Format formats the stack of Frames according to the fmt.Formatter interface.
//
//    %s	lists source files for each Frame in the stack
//    %v	lists the source file and line number for each Frame in the stack
//
// Format accepts flags that alter the printing of some verbs, as follows:
//
//    %+v   Prints filename, function, and line number for each Frame in the stack.
func (st StackTrace) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		switch {
		case s.Flag('+'):
			for _, f := range st {
				fmt.Fprintf(s, "\n%+v", f)
			}
		case s.Flag('#'):
			fmt.Fprintf(s, "%#v", []StackFrame(st))
		default:
			st.formatSlice(s, verb)
		}
	case 's':
		st.formatSlice(s, verb)
	}
}

// formatSlice will format this StackTrace into the given buffer as a slice of
// Frame, only valid when called with '%s' or '%v'.
func (st StackTrace) formatSlice(s fmt.State, verb rune) {
	_, _ = io.WriteString(s, "[")
	for i, f := range st {
		if i > 0 {
			_, _ = io.WriteString(s, " ")
		}
		f.Format(s, verb)
	}
	_, _ = io.WriteString(s, "]")
}
