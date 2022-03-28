package errors

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
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
	// Value contains the value that was wrong for errors that need it, like ArgumentInvalidError
	// TODO: use structpb
	Value interface{} `json:"value"`
	// Cause contains the error that caused this error (to wrap a json error in a JSONMarshalError, for example)
	Cause error `json:"cause,omitempty"`
	// stack contains the StackTrace when this Error is instanciated
	Stack StackTrace `json:"-"`
}

// New creates a new instance of this error.
//
// New also records the stack trace at the point it was called.
func (e Error) New() error {
	final := e
	if final.Code == 0 {
		final.Code = http.StatusInternalServerError
	}
	if len(final.ID) == 0 {
		final.ID = "error.runtime"
	}
	final.Stack.Initialize()
	return final
}

// GetID tells the ID of this Error
func (e Error) GetID() string {
	return e.ID
}

// Is tells if this error matches the target.
//
// implements errors.Is interface (package "errors").
func (e Error) Is(target error) bool {
	if actual, ok := target.(Error); ok {
		return e.ID == actual.ID
	}
	return false
}

// Extract extracts an Error with the same ID as this Error from the error chain
func (e Error) Extract(err error) (extracted Error, found bool) {
	for err != nil {
		if identifiable, ok := err.(interface{ GetID() string }); ok && identifiable.GetID() == e.GetID() {
			extracted := Error{}
			if As(err, &extracted) {
				return extracted, true
			}
		}
		err = Unwrap(err)
	}
	return Error{}, false
}

// Wrap wraps the given error in this Error.
//
// If err is nil, Wrap returns nil.
//
// Wrap also records the stack trace at the point it was called.
func (e Error) Wrap(err error) error {
	if err == nil {
		return nil
	}
	final := e
	final.Cause = err
	final.Stack.Initialize()
	return final
}

// Unwrap gives the Cause of this Error, if any.
//
// implements errors.Unwrap interface (package "errors").
func (e Error) Unwrap() error {
	return e.Cause
}

// With creates a new Error from a given sentinel telling "what" is wrong and eventually their value.
//
// With also records the stack trace at the point it was called.
func (e Error) With(what string, values ...interface{}) error {
	final := e
	final.What = what
	if len(values) > 0 {
		final.Value = values[0]
	}
	final.Stack.Initialize()
	return final
}

// WithStack creates a new error from a given Error and records its stack.
func (e Error) WithStack() error {
	final := e
	final.Stack.Initialize()
	return final
}

// WithoutStack creates a new error from a given Error and records its stack.
func (e Error) WithoutStack() error {
	final := e
	final.Stack = StackTrace{}
	return final
}

// Error returns the string version of this error.
//
// implements error interface.
func (e Error) Error() string {
	var sb strings.Builder

	switch strings.Count(e.Text, "%") {
	case 0:
		sb.WriteString(e.Text)
	case 1:
		fmt.Fprintf(&sb, e.Text, e.What)
	default:
		fmt.Fprintf(&sb, e.Text, e.What, e.Value)
	}
	if e.Cause != nil {
		if len(e.Text) > 0 {
			sb.WriteString(": ")
		}
		sb.WriteString(e.Cause.Error())
	}
	return sb.String()
}

// GoString returns the Go syntax of this Error
//
// implements fmt.GoStringer
func (e Error) GoString() string {
	var sb strings.Builder

	_, _ = sb.WriteString("errors.Error{Code:")
	_, _ = sb.WriteString(strconv.Itoa(e.Code))
	_, _ = sb.WriteString(", ID:\"")
	_, _ = sb.WriteString(e.ID)
	_, _ = sb.WriteString("\", Text:\"")
	_, _ = sb.WriteString(e.Text)
	_, _ = sb.WriteString("\", What:\"")
	_, _ = sb.WriteString(e.What)
	_, _ = sb.WriteString("\", Value:")
	_, _ = sb.WriteString(fmt.Sprintf("%#v", e.Value))
	_, _ = sb.WriteString(", Cause:")
	_, _ = sb.WriteString(fmt.Sprintf("%#v", e.Cause))
	_, _ = sb.WriteString(", Stack:")
	_, _ = sb.WriteString(fmt.Sprintf("%#v", e.Stack))
	_, _ = sb.WriteString("}")
	return sb.String()
}

// Format interprets fmt State and rune to generate an output for fmt.Sprintf, etc
//
// implements fmt.Formatter
func (e Error) Format(state fmt.State, verb rune) {
	switch verb {
	case 'v':
		if state.Flag('+') {
			_, _ = io.WriteString(state, e.Error())
			e.Stack.Format(state, verb)
			return
		}
		if state.Flag('#') {
			_, _ = io.WriteString(state, e.GoString())
			return
		}
		fallthrough
	case 's':
		_, _ = io.WriteString(state, e.Error())
	case 'q':
		_, _ = fmt.Fprintf(state, "%q", e.Error())
	}
}

// MarshalJSON marshals this into JSON
func (e Error) MarshalJSON() ([]byte, error) {
	type surrogate Error
	data, err := json.Marshal(struct {
		Type string `json:"type"`
		surrogate
	}{
		Type:      "error",
		surrogate: surrogate(e),
	})
	return data, JSONMarshalError.Wrap(err)
}

// UnmarshalJSON decodes JSON
func (e *Error) UnmarshalJSON(payload []byte) (err error) {
	type surrogate Error
	var inner struct {
		Type string `json:"type"`
		surrogate
	}
	if err = json.Unmarshal(payload, &inner); err != nil {
		return JSONUnmarshalError.Wrap(err)
	}
	if inner.Type != "error" {
		return JSONUnmarshalError.Wrap(InvalidType.With("error", inner.Type))
	}
	*e = Error(inner.surrogate)
	return nil
}
