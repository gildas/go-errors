package errors

import (
	"fmt"
	"strings"

	pkerrors "github.com/pkg/errors"
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
}

// New returns an error with the supplied message.
// New also records the stack trace at the point it was called.
func New(message string) error {
	return pkerrors.New(message)
}

func (e Error) Error() string {
	// implements error interface
	switch strings.Count(e.Text, "%") {
	case 0:  return e.Text
	case 1:  return fmt.Sprintf(e.Text, e.What)
	default: return fmt.Sprintf(e.Text, e.What, e.Value)
	}
}

func (e Error) Is(target error) bool {
	inner, ok := target.(Error)
	if !ok {
		return false
	}
	return e.ID == inner.ID
}