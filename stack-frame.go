package errors

import (
	"fmt"
	"runtime"
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
