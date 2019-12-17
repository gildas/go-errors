package errors

import (
	"fmt"
	"strings"
)

// Error describes an augmented implementation of Go's error interface
type Error struct {
	// Code is an numerical code, like an HTTP Status Code
	Code int `json:"code"`
	// ID is the string identifier, like: "error.argument.invalid"
	ID string `json:"id"`
	// Text is the human readable error message
	Text string `json:"text"`
	// What contains what element is wrong for errors that need it, like NotFoundError
	What string `json:"what,omitempty"`
	// Value contains the value that was wrong for errros that need it, like ArgumentInvalidError
	Value interface{} `json:"value"`
	// Cause contains the error that caused this error (to wrap a json error in a JSONMarshalError, for example)
	Cause error
}

// New creates a new instance of this error
// New also records the stack trace at the point it was called.
func (e Error) New() error {
	final := e
	return WithStack(&final)
}

// WithMessage annotates a new instance of this error with a new message.
// If err is nil, WithMessage returns nil.
// WithMessage also records the stack trace at the point it was called.
func (e Error) WithMessage(message string) error {
	final := e
	return WithMessage(&final, message)
}

// Error returns the string version of this error
// implements error interface
func (e Error) Error() string {
	switch strings.Count(e.Text, "%") {
	case 0:  return e.Text
	case 1:  return fmt.Sprintf(e.Text, e.What)
	default: return fmt.Sprintf(e.Text, e.What, e.Value)
	}
}

// Is tells if this error matches the target
// implements errors.Is interface (package "errors")
func (e Error) Is(target error) bool {
	inner, ok := target.(Error)
	if !ok {
		return false
	}
	return e.ID == inner.ID
}

func (e Error) Wrap(err error) error {
	final := e
	e.Cause = err
	return WithStack(&final)
}