package errors

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
)

// Error describes an augmented implementation of Go's error interface
type Error struct {
	// Code is an numerical code, like an HTTP Status Code
	Code int `json:"code,omitempty"`
	// ID is the string identifier, like: "error.argument.invalid"
	ID string `json:"id,omitempty"`
	// Text is the human readable error message
	Text string `json:"text,omitempty"`
	// What contains what element is wrong for errors that need it, like NotFoundError
	What string `json:"what,omitempty"`
	// Value contains the value that was wrong for errors that need it, like ArgumentInvalidError
	// TODO: use structpb
	Value interface{} `json:"value,omitempty"`
	// Origin contains the real error from another package, if any
	Origin error `json:"-"`
	// Cause contains the error that caused this error
	Cause error `json:"-"`
	// stack contains the StackTrace when this Error is instanciated
	Stack StackTrace `json:"-"`
}

// Clone creates an exact copy of this Error
func (e Error) Clone() *Error {
	final := e
	return &final
}

// GetID tells the ID of this Error
func (e Error) GetID() string {
	return e.ID
}

// Is tells if this error matches the target.
//
// implements errors.Is interface (package "errors").
//
// To check if an error is an errors.Error, simply write:
//
//	if errors.Is(err, errors.Error{}) {
//	  // do something with err
//	}
func (e Error) Is(target error) bool {
	if actual, ok := target.(Error); ok {
		if len(actual.ID) == 0 {
			return true // no ID means any error is a match
		}
		return e.ID == actual.ID
	}
	if actual, ok := target.(*Error); ok && actual != nil {
		if len(actual.ID) == 0 {
			return true // no ID means any error is a match
		}
		return e.ID == actual.ID
	}
	if e.Origin != nil {
		return Is(e.Origin, target)
	}
	return false
}

// As attempts to convert the given error into the given target
//
// As returns true if the conversion was successful and the target is now populated.
//
// Example:
//
//	target := errors.ArgumentInvalid.Clone()
//	if errors.As(err, &target) {
//	  // do something with target
//	}
func (e Error) As(target interface{}) bool {
	if actual, ok := target.(**Error); ok {
		if *actual != nil && (*actual).GetID() != e.ID {
			return false
		}
		copy := e
		*actual = &copy
		return true
	}
	if e.Origin != nil {
		return As(e.Origin, target)
	}
	return false
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
	if len(final.Stack) == 0 {
		final.Stack.Initialize()
	}
	return final
}

// Unwrap gives the first Cause of this Error, if any.
//
// implements errors.Unwrap interface (package "errors").
func (e Error) Unwrap() error {
	if e.Cause == nil {
		return nil
	}
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
	// At some point this should be a pointer receiver
	// https://github.com/jszwec/csvutil/issues/39
	// But when it is, it breaks the errors.As() as it cannot find sentinel errors anymore:
	// Line wrap.go:92 is always true so line wrap.go:96 is never reached and Error.As never called.
	// https://cs.opensource.google/go/go/+/refs/tags/go1.19.3:src/errors/wrap.go;drc=2580d0e08d5e9f979b943758d3c49877fb2324cb;l=92
	if e.Origin != nil {
		return e.Origin.Error()
	}
	var sb strings.Builder

	switch strings.Count(e.Text, "%") - strings.Count(e.Text, "%%") {
	case 0:
		if len(e.Text) > 0 {
			_, _ = sb.WriteString(e.Text)
		} else if len(e.ID) > 0 {
			_, _ = sb.WriteString(e.ID)
		} else {
			_, _ = sb.WriteString("runtime error")
		}
	case 1:
		_, _ = fmt.Fprintf(&sb, e.Text, e.What)
	default:
		_, _ = fmt.Fprintf(&sb, e.Text, e.What, e.Value)
	}
	if e.Cause != nil {
		_, _ = sb.WriteString("\nCaused by:")
		_, _ = sb.WriteString("\n\t")
		_, _ = sb.WriteString(e.Cause.Error())
	}
	return sb.String()
}

// GoString returns the Go syntax of this Error
//
// implements fmt.GoStringer
func (e Error) GoString() string {
	var sb strings.Builder

	_, _ = fmt.Fprintf(&sb, `errors.Error{Code: %d, ID: "%s", Text: "%s"`, e.Code, e.ID, e.Text)
	if len(e.What) > 0 {
		_, _ = fmt.Fprintf(&sb, `, What: "%s"`, e.What)
	}
	if e.Value != nil {
		_, _ = fmt.Fprintf(&sb, `, Value: %#v`, e.Value)
	}
	if e.Cause != nil {
		_, _ = sb.WriteString(", Cause: ")
		if gostringer, ok := e.Cause.(fmt.GoStringer); ok {
			_, _ = sb.WriteString(gostringer.GoString())
		} else {
			_, _ = sb.WriteString(`"`)
			_, _ = sb.WriteString(e.Cause.Error())
			_, _ = sb.WriteString(`"`)
		}
	}
	if len(e.Stack) > 0 {
		_, _ = fmt.Fprintf(&sb, `, Stack: %#v`, e.Stack)
	}
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
	var payload interface{}
	var cause *Error

	if e.Cause != nil {
		if value, ok := e.Cause.(Error); ok {
			cause = &value
		} else {
			var id strings.Builder
			causeType := reflect.TypeOf(e.Cause)
			if causeType.Kind() == reflect.Ptr {
				causeType = causeType.Elem()
			}
			_, _ = id.WriteString("error.runtime")
			if causeType.PkgPath() != "errors" || causeType.Name() != "errorString" {
				_, _ = id.WriteString(".")
				_, _ = id.WriteString(causeType.String())
			}
			cause = &Error{Code: http.StatusInternalServerError, ID: id.String(), Text: e.Cause.Error()}
		}
	}

	payload = struct {
		Type string `json:"type"`
		surrogate
		Cause *Error `json:"cause,omitempty"`
	}{
		Type:      "error",
		surrogate: surrogate(e),
		Cause:     cause,
	}
	data, err := json.Marshal(payload)
	return data, JSONMarshalError.Wrap(err)
}

// UnmarshalJSON decodes JSON
func (e *Error) UnmarshalJSON(payload []byte) (err error) {
	type surrogate Error
	var inner struct {
		Type string `json:"type"`
		surrogate
		Cause *Error `json:"cause,omitempty"`
	}
	if err = json.Unmarshal(payload, &inner); err != nil {
		return JSONUnmarshalError.Wrap(err)
	}
	if inner.Type != "error" {
		return JSONUnmarshalError.Wrap(InvalidType.With("error", inner.Type))
	}
	*e = Error(inner.surrogate)
	if inner.Cause != nil {
		e.Cause = *inner.Cause
	}
	return nil
}
