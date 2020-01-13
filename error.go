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

// New creates a new instance of this error.
// New also records the stack trace at the point it was called.
func (e Error) New() error {
	final := e
	return WithStack(&final)
}

// WithMessage annotates a new instance of this error with a new message.
// If err is nil, WithMessage returns nil.
//
// WithMessage also records the stack trace at the point it was called.
func (e Error) WithMessage(message string) error {
	final := e
	return WithMessage(&final, message)
}

// WithMessagef annotates a new instance of this error with a new message.
// The message is a format with eventually some arguments
// If err is nil, WithMessage returns nil.
//
// WithMessage also records the stack trace at the point it was called.
func (e Error) WithMessagef(format string, args ...interface{}) error {
	final := e
	return WithMessagef(&final, format, args...)
}

// Error returns the string version of this error.
func (e Error) Error() string {
	// implements error interface
	var sb strings.Builder

	switch strings.Count(e.Text, "%") {
	case 0:  sb.WriteString(e.Text)
	case 1:  fmt.Fprintf(&sb, e.Text, e.What)
	default: fmt.Fprintf(&sb, e.Text, e.What, e.Value)
	}
	if e.Cause != nil {
		sb.WriteString(": ")
		sb.WriteString(e.Cause.Error())
	}
	return sb.String()
}

// Is tells if this error matches the target.
func (e Error) Is(target error) bool {
	// implements errors.Is interface (package "errors")
	if pactual, ok := target.(*Error); ok {
		return e.ID == pactual.ID
	}
	if actual, ok := target.(Error); ok {
		return e.ID == actual.ID
	}
	return false
}

// Wrap wraps the given error in this Error.
// If err is nil, Wrap returns nil.
func (e Error) Wrap(err error) error {
	if err == nil {
		return nil
	}
	final := e
	final.Cause = err
	return WithStack(&final)
}

// Unwrap gives the Cause of this Error, if any.
func (e Error) Unwrap() error {
	// implements errors.Unwrap interface (package "errors")
	return e.Cause
}

// With creates a new Error from a given sentinel telling "what" is wrong and eventually their value.
func (e *Error) With(what string, values ...interface{}) Error {
	final := *e
	final.What = what
	if len(values) > 0 {
		final.Value = values[0]
	}
	return final
}

// WithStack creates a new error from a given Error and records its stack.
func (e Error) WithStack() error {
	return WithStack(&e)
}