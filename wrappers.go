package errors

import (
	goerrors "errors"
	"fmt"
	"net/http"
)

// New returns a new error with the supplied message.
//
// New also records the stack trace at the point it was called.
func New(message string) error {
	return Error{Code: http.StatusInternalServerError, ID: "error.runtime", Text: message}.WithStack()
}

// Errorf formats according to a format specifier and returns the string
// as a value that satisfies error.
//
// Errorf also records the stack trace at the point it was called.
func Errorf(format string, args ...interface{}) error {
	return Error{Code: http.StatusInternalServerError, ID: "error.runtime", Text: fmt.Sprintf(format, args...)}.WithStack()
}

// WithStack annotates err with a stack trace at the point WithStack was called.
//
// If err is nil, WithStack returns nil.
//
// If err is already annotated with a stack trace, WithStack returns err.
func WithStack(err error) error {
	if err == nil {
		return nil
	}
	if _err, ok := err.(Error); ok {
		if len(_err.Stack) == 0 {
			return _err.WithStack()
		}
		return _err
	}
	return Error{Code: http.StatusInternalServerError, ID: "error.runtime"}.Wrap(err)
}

// WithoutStack removes the stack trace from the current error
//
// If err is nil, WithStack returns nil.
func WithoutStack(err error) error {
	if err == nil {
		return nil
	}
	if err0, ok := err.(Error); ok {
		return err0.WithoutStack()
	}
	return err
}

// Wrap returns an error annotating err with a stack trace
// at the point Wrap is called, and the supplied message.
//
// If err is nil, Wrap returns nil.
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	return Error{Code: http.StatusInternalServerError, ID: "error.runtime", Text: message}.Wrap(err)
}

// Wrapf returns an error annotating err with a stack trace
// at the point Wrapf is called, and the format specifier.
//
// If err is nil, Wrapf returns nil.
func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return Error{Code: http.StatusInternalServerError, ID: "error.runtime", Text: fmt.Sprintf(format, args...)}.Wrap(err)
}

// WrapErrors returns an error wrapping given errors
//
// If err is nil, WrapErrors returns nil.
//
// If no errors are given, WrapErrors returns err.
func WrapErrors(err error, errors ...error) error {
	if err == nil {
		return nil
	}
	if len(errors) == 0 {
		return err
	}
	container, ok := err.(Error)
	if !ok {
		container = Error{Code: http.StatusInternalServerError, ID: "error.runtime", Text: err.Error()}
	}
	if len(errors) == 1 {
		if errors[0] != nil {
			return container.Wrap(errors[0])
		}
		return container
	}
	return container.Wrap(WrapErrors(errors[0], errors[1:]...))
}

// WithMessage annotates err with a new message.
//
// If err is nil, WithMessage returns nil.
func WithMessage(err error, message string) error {
	if err == nil {
		return nil
	}
	return Error{Code: http.StatusInternalServerError, ID: "error.runtime", Text: message}.Wrap(err)
}

// WithMessagef annotates err with the format specifier.
//
// If err is nil, WithMessagef returns nil.
func WithMessagef(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return Error{Code: http.StatusInternalServerError, ID: "error.runtime", Text: fmt.Sprintf(format, args...)}.Wrap(err)
}

//***************** goerrors

// Is reports whether any error in err's chain matches target.
//
// The chain consists of err itself followed by the sequence of errors obtained by
// repeatedly calling Unwrap.
//
// An error is considered to match a target if it is equal to that target or if
// it implements a method Is(error) bool such that Is(target) returns true.
func Is(err, target error) bool {
	return goerrors.Is(err, target)
}

// Unwrap returns the result of calling the Unwrap method on err, if err's
// type contains an Unwrap method returning error.
//
// Otherwise, Unwrap returns nil.
func Unwrap(err error) error {
	return goerrors.Unwrap(err)
}

// As finds the first error in err's chain that matches target, and if so, sets
// target to that error value and returns true.
//
// The chain consists of err itself followed by the sequence of errors obtained by
// repeatedly calling Unwrap.
//
// An error matches target if the error's concrete value is assignable to the value
// pointed to by target, or if the error has a method As(interface{}) bool such that
// As(target) returns true. In the latter case, the As method is responsible for
// setting target.
//
// As will panic if target is not a non-nil pointer to either a type that implements
// error, or to any interface type. As returns false if err is nil.
func As(err error, target interface{}) bool {
	return goerrors.As(err, target)
}
