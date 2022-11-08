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

	}
}

// Append appends a new error
//
// If the error is nil, nothing is added
func (me *MultiError) Append(err error) *MultiError {
	if err != nil {
		me.Errors = append(me.Errors, err)
	}
	return me
}

// AsError returns this if it contains errors, nil otherwise
//
// If this contains only one error, that error is returned
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
