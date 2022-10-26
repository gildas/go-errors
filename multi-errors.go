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
	text := strings.Builder{}
	for _, err := range me.Errors {
		text.WriteString(err.Error())
		text.WriteString("\n")
	}
	return fmt.Sprintf("%d Errors:\n%s", len(me.Errors), text.String())
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
func (me *MultiError) AsError() error {
	if me == nil || len(me.Errors) == 0 {
		return nil
	}
	return WithStack(me)
}
