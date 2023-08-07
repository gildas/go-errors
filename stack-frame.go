package errors

import (
	"encoding/json"
	"fmt"
	"io"
	"path"
	"runtime"
	"strconv"
	"strings"
)

/*
Imported from https://github.com/pkg/errors/blob/master/stack.go
*/

type StackFrame uintptr

func (frame StackFrame) pc() uintptr {
	return uintptr(frame) - 1
}

func (frame StackFrame) Filepath() string {
	function := runtime.FuncForPC(frame.pc())
	if function == nil {
		return "unknown"
	}
	file, _ := function.FileLine(frame.pc())
	return file
}

func (frame StackFrame) Line() int {
	function := runtime.FuncForPC(frame.pc())
	if function == nil {
		return 0
	}
	_, line := function.FileLine(frame.pc())
	return line

}

func (frame StackFrame) FuncName() string {
	function := runtime.FuncForPC(frame.pc())
	if function == nil {
		return "unknown"
	}
	return function.Name()
}

func (frame StackFrame) MarshalText() ([]byte, error) {
	funcName := frame.FuncName()
	if funcName == "unknown" {
		return []byte(funcName), nil
	}
	return []byte(fmt.Sprintf("%s %s:%d", funcName, frame.Filepath(), frame.Line())), nil
}

// Format formats the frame according to the fmt.Formatter interface.
//
//	%s    source file
//	%d    source line
//	%n    function name
//	%v    equivalent to %s:%d
//
// Format accepts flags that alter the printing of some verbs, as follows:
//
//	%+s   function name and path of source file relative to the compile time
//	      GOPATH separated by \n\t (<funcname>\n\t<path>)
//	%+v   equivalent to %+s:%d
func (frame StackFrame) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		switch {
		case s.Flag('+'):
			_, _ = io.WriteString(s, frame.FuncName())
			_, _ = io.WriteString(s, "\n\t")
			_, _ = io.WriteString(s, frame.Filepath())
		default:
			_, _ = io.WriteString(s, path.Base(frame.Filepath()))
		}
	case 'd':
		_, _ = io.WriteString(s, strconv.Itoa(frame.Line()))
	case 'n':
		_, _ = io.WriteString(s, funcname(frame.FuncName()))
	case 'v':
		frame.Format(s, 's')
		_, _ = io.WriteString(s, ":")
		frame.Format(s, 'd')
	}
}

// MarshalJSON marshals this into JSON
func (frame StackFrame) MarshalJSON() ([]byte, error) {
	data, err := json.Marshal(struct {
		FuncName string `json:"func"`
		FuncLine int    `json:"line"`
		FuncPath string `json:"path"`
	}{
		FuncName: frame.FuncName(),
		FuncLine: frame.Line(),
		FuncPath: frame.Filepath(),
	})
	return data, JSONMarshalError.Wrap(err)
}

// funcname removes the path prefix component of a function's name
func funcname(name string) string {
	i := strings.LastIndex(name, "/")
	name = name[i+1:]
	i = strings.Index(name, ".")
	return name[i+1:]
}
