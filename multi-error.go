package errors

import (
	"fmt"
	"strings"
)

// MultiError is used to collect errors, like during a loop
type MultiError struct {
	Errors []error `json:"errors"`
}

// Error returns the string version of this error
//
// implements error.Error interface
func (me *MultiError) Error() string {
	if len(me.Errors) == 0 {
		return ""
	}
	if len(me.Errors) == 1 {
		return me.Errors[0].Error()
	}
	text := strings.Builder{}
	for _, err := range me.Errors {
		text.WriteString("\n")
		text.WriteString(err.Error())
	}
	return fmt.Sprintf("%d errors:%s", len(me.Errors), text.String())
}

// IsEmpty returns true if this MultiError contains no errors
func (me *MultiError) IsEmpty() bool {
	return me == nil || len(me.Errors) == 0
}

// Append appends new errors
//
// If an error is nil, it is not added
func (me *MultiError) Append(errs ...error) {
	for _, err := range errs {
		if err != nil {
			me.Errors = append(me.Errors, err)
		}
	}
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
func (e MultiError) Is(target error) bool {
	if _, ok := target.(*MultiError); ok {
		return true
	}
	for _, err := range e.Errors {
		if Is(err, target) {
			return true
		}
	}
	return false
}

// As attempts to convert the given error into the given target
//
// The first error to match the target is returned
func (e MultiError) As(target interface{}) bool {
	for _, err := range e.Errors {
		if As(err, target) {
			return true
		}
	}
	return false
}

// AsError returns this if it contains errors, nil otherwise
//
// If this contains only one error, that error is returned.
//
// AsError also records the stack trace at the point it was called.
func (me *MultiError) AsError() error {
	if me == nil || len(me.Errors) == 0 {
		return nil
	}
	if len(me.Errors) == 1 {
		if err, ok := me.Errors[0].(*Error); ok {
			if len(err.Stack) == 0 {
				return err.WithStack()
			}
			return err
		}
		return WithStack(me.Errors[0])
	}
	return WithStack(me)
}
