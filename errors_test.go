package errors_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/gildas/go-errors"
)

type ErrorsSuite struct {
	suite.Suite
	Name string
}

func TestErrorsSuite(t *testing.T) {
	suite.Run(t, new(ErrorsSuite))
}

func (suite *ErrorsSuite) SetupSuite() {
	suite.Name = strings.TrimSuffix(reflect.TypeOf(*suite).Name(), "Suite")
}

func (suite *ErrorsSuite) TestCanCreate() {
	err := errors.NewSentinel(32123, "error.test.create", "this is the error")
	suite.Require().NotNil(err, "newly created sentinel cannot be nil")
}

func (suite *ErrorsSuite) TestCanTellIsError() {
	err := errors.NotFoundError.With("key").WithStack()
	suite.Require().NotNil(err, "err should not be nil")
	suite.Assert().True(errors.Is(err, errors.NotFoundError), "err should match a NotFoundError (pointer)")
	suite.Assert().True(errors.Is(err, *errors.NotFoundError), "err should match a NotFoundError (object)")
}

func (suite *ErrorsSuite) TestCanTellContainsAnError() {
	err := errors.NotFoundError.With("key").WithStack()
	suite.Require().NotNil(err, "err should not be nil")
	var inner *errors.Error
	suite.Assert().True(errors.As(err, &inner), "err should contain an errors.Error")
}

func (suite *ErrorsSuite) TestCanTellDoesNotContainAnError() {
	err := errors.Errorf("Houston, we have a problem")
	suite.Require().NotNil(err, "err should not be nil")
	suite.Assert().False(errors.Is(err, errors.Error{}), "err should not contain an errors.Error")
	var inner *errors.Error
	suite.Assert().False(errors.As(err, &inner), "err should not contain an errors.Error")
}

func (suite *ErrorsSuite) TestCanWrap() {
	err := errors.Errorf("Houston, we have a problem")
	wrapped := errors.NotImplementedError.Wrap(err)
	suite.Require().NotNil(wrapped)
	suite.Assert().True(errors.Is(wrapped, errors.NotImplementedError), "wrapped err should be a Not Implemented Error")
	wrapped = errors.NotImplementedError.Wrap(nil)
	suite.Require().Nil(wrapped, "Wrapped error of nil should be nil")
}

func (suite *ErrorsSuite) TestFailsWithNonErrorTarget() {
	suite.Assert().False(errors.NotFoundError.Is(errors.New("Hello")), "Error should not be a pkg/errors.fundamental")
}

func (suite *ErrorsSuite) TestWrappers() {
	var err error

	err = errors.New("Hello World")
	suite.Assert().NotNil(err)

	err = errors.Errorf("Hello World")
	suite.Assert().NotNil(err)

	err = errors.WithMessage(errors.NotFoundError, "Hello World")
	suite.Assert().NotNil(err)

	err = errors.WithMessagef(errors.NotFoundError, "Hello World")
	suite.Assert().NotNil(err)

	err = errors.Wrap(errors.NotFoundError, "Hello World")
	suite.Assert().NotNil(err)

	err = errors.Wrapf(errors.NotFoundError, "Hello World")
	suite.Assert().NotNil(err)

	unwrapped := errors.Unwrap(err)
	suite.Assert().NotNil(unwrapped)

	suite.Assert().True(errors.Is(err, errors.NotFoundError), "err should be of the same type as NotFoundError")

	var inner *errors.Error
	suite.Assert().True(errors.As(err, &inner), "Inner Error should be an errors.Error")
}

func ExampleError() {
	sentinel := errors.NewSentinel(500, "error.test.custom", "Test Error")
	err := sentinel.New()
	if err != nil {
		fmt.Println(err)

		var details *errors.Error
		if errors.As(err, &details) {
			fmt.Println(details.ID)
		}
	}
	// Output:
	// Test Error
	// error.test.custom
}

func ExampleError_WithMessage() {
	sentinel := errors.NewSentinel(500, "error.test.custom", "Test Error")
	err := sentinel.WithMessage("hmmm... this is bad")
	if err != nil {
		fmt.Println(err)

		var details *errors.Error
		if errors.As(err, &details) {
			fmt.Println(details.ID)
		}
	}
	// Output:
	// hmmm... this is bad: Test Error
	// error.test.custom
}

func ExampleError_WithMessagef() {
	sentinel := errors.NewSentinel(500, "error.test.custom", "Test Error")
	err := sentinel.WithMessagef("hmmm... this is bad %s", "stuff")
	if err != nil {
		fmt.Println(err)

		var details *errors.Error
		if errors.As(err, &details) {
			fmt.Println(details.ID)
		}
	}
	// Output:
	// hmmm... this is bad stuff: Test Error
	// error.test.custom
}

func ExampleError_WithStack() {
	err := errors.NotImplementedError.WithStack()
	if err != nil {
		fmt.Println(err)

		var details *errors.Error
		if errors.As(err, &details) {
			fmt.Println(details.ID)
		}
	}
	// Output:
	// Not Implemented
	// error.notimplemented
}

func ExampleError_With() {
	err := errors.ArgumentMissingError.With("key").WithStack()
	if err != nil {
		fmt.Println(err)

		var details *errors.Error
		if errors.As(err, &details) {
			fmt.Println(details.ID)
		}
	}
	// Output:
	// Argument key is missing
	// error.argument.missing
}

func ExampleError_With_value() {
	err := errors.ArgumentInvalidError.With("key", "value").WithStack()
	if err != nil {
		fmt.Println(err)

		var details *errors.Error
		if errors.As(err, &details) {
			fmt.Println(details.ID)
		}
	}
	// Output:
	// Argument key is invalid (value: value)
	// error.argument.invalid
}

func ExampleError_With_array() {
	err := errors.ArgumentInvalidError.With("key", []string{"value1", "value2"}).WithStack()
	if err != nil {
		fmt.Println(err)

		var details *errors.Error
		if errors.As(err, &details) {
			fmt.Println(details.ID)
		}
	}
	// Output:
	// Argument key is invalid (value: [value1 value2])
	// error.argument.invalid
}

func ExampleError_Wrap() {
	var payload struct {
		Value string `json:"value"`
	}

	err := json.Unmarshal([]byte(`{"value": 0`), &payload)
	if err != nil {
		finalerr := errors.JSONMarshalError.Wrap(err)
		var details *errors.Error
		if errors.As(finalerr, &details) {
			fmt.Println(details.ID)
		}

		fmt.Println(finalerr)

		cause := details.Unwrap()
		if cause != nil {
			fmt.Println(cause)
		}
	}
	// Output:
	// error.json.marshal
	// JSON failed to marshal data: unexpected end of JSON input
	// unexpected end of JSON input
}

func (suite *ErrorsSuite) TestCanCreateMultiError() {
	err := &errors.MultiError{}
	suite.Require().NotNil(err, "err should not be nil")
	suite.Assert().Nil(err.AsError(), "err should contain nothing")
	suite.Assert().Equal("", err.Error())
	_ = err.Append(errors.New("this is an error"))
	suite.Assert().NotNil(err.AsError(), "err should contain something")
}

func ExampleMultiError() {
	err := &errors.MultiError{}
	_ = err.Append(errors.New("this is the first error"))
	_ = err.Append(errors.New("this is the second error"))
	fmt.Println(err)
	// Output:
	// 2 Errors:
	// this is the first error
	// this is the second error
}

func (suite *ErrorsSuite) TestCanCreateFromHTTPStatus() {
	var err error
	var details *errors.Error

	err = errors.FromHTTPStatusCode(http.StatusBadGateway)
	suite.Assert().Truef(errors.Is(err, errors.HTTPBadGatewayError), "err should match a %s", errors.HTTPBadGatewayError.ID)

	err = errors.FromHTTPStatusCode(http.StatusBadRequest)
	suite.Assert().Truef(errors.Is(err, errors.HTTPBadRequestError), "err should match a %s", errors.HTTPBadRequestError.ID)

	err = errors.FromHTTPStatusCode(http.StatusForbidden)
	suite.Assert().Truef(errors.Is(err, errors.HTTPForbiddenError), "err should match a %s", errors.HTTPForbiddenError.ID)

	err = errors.FromHTTPStatusCode(http.StatusInternalServerError)
	suite.Assert().Truef(errors.Is(err, errors.HTTPInternalServerError), "err should match a %s", errors.HTTPInternalServerError.ID)

	err = errors.FromHTTPStatusCode(http.StatusMethodNotAllowed)
	suite.Assert().Truef(errors.Is(err, errors.HTTPMethodNotAllowedError), "err should match a %s", errors.HTTPMethodNotAllowedError.ID)

	err = errors.FromHTTPStatusCode(http.StatusNotFound)
	suite.Assert().Truef(errors.Is(err, errors.HTTPNotFoundError), "err should match a %s", errors.HTTPNotFoundError.ID)

	err = errors.FromHTTPStatusCode(http.StatusNotImplemented)
	suite.Assert().Truef(errors.Is(err, errors.HTTPNotImplementedError), "err should match a %s", errors.HTTPNotImplementedError.ID)

	err = errors.FromHTTPStatusCode(http.StatusServiceUnavailable)
	suite.Assert().Truef(errors.Is(err, errors.HTTPServiceUnavailableError), "err should match a %s", errors.HTTPServiceUnavailableError.ID)

	err = errors.FromHTTPStatusCode(http.StatusUnauthorized)
	suite.Assert().Truef(errors.Is(err, errors.HTTPUnauthorizedError), "err should match a %s", errors.HTTPUnauthorizedError.ID)

	err = errors.FromHTTPStatusCode(http.StatusConflict)
	suite.Assert().Truef(errors.Is(err, errors.HTTPStatusConflictError), "err should match a %s", errors.HTTPStatusConflictError.ID)

	err = errors.FromHTTPStatusCode(http.StatusExpectationFailed)
	suite.Assert().Truef(errors.Is(err, errors.HTTPStatusExpectationFailedError), "err should match a %s", errors.HTTPStatusExpectationFailedError.ID)

	err = errors.FromHTTPStatusCode(http.StatusFailedDependency)
	suite.Assert().Truef(errors.Is(err, errors.HTTPStatusFailedDependencyError), "err should match a %s", errors.HTTPStatusFailedDependencyError.ID)

	err = errors.FromHTTPStatusCode(http.StatusGatewayTimeout)
	suite.Assert().Truef(errors.Is(err, errors.HTTPStatusGatewayTimeoutError), "err should match a %s", errors.HTTPStatusGatewayTimeoutError.ID)

	err = errors.FromHTTPStatusCode(http.StatusGone)
	suite.Assert().Truef(errors.Is(err, errors.HTTPStatusGoneError), "err should match a %s", errors.HTTPStatusGoneError.ID)

	err = errors.FromHTTPStatusCode(http.StatusHTTPVersionNotSupported)
	suite.Assert().Truef(errors.Is(err, errors.HTTPStatusHTTPVersionNotSupportedError), "err should match a %s", errors.HTTPStatusHTTPVersionNotSupportedError.ID)

	err = errors.FromHTTPStatusCode(http.StatusInsufficientStorage)
	suite.Assert().Truef(errors.Is(err, errors.HTTPStatusInsufficientStorageError), "err should match a %s", errors.HTTPStatusInsufficientStorageError.ID)

	err = errors.FromHTTPStatusCode(http.StatusLengthRequired)
	suite.Assert().Truef(errors.Is(err, errors.HTTPStatusLengthRequiredError), "err should match a %s", errors.HTTPStatusLengthRequiredError.ID)

	err = errors.FromHTTPStatusCode(http.StatusLocked)
	suite.Assert().Truef(errors.Is(err, errors.HTTPStatusLockedError), "err should match a %s", errors.HTTPStatusLockedError.ID)

	err = errors.FromHTTPStatusCode(http.StatusLoopDetected)
	suite.Assert().Truef(errors.Is(err, errors.HTTPStatusLoopDetectedError), "err should match a %s", errors.HTTPStatusLoopDetectedError.ID)

	err = errors.FromHTTPStatusCode(http.StatusMisdirectedRequest)
	suite.Assert().Truef(errors.Is(err, errors.HTTPStatusMisdirectedRequestError), "err should match a %s", errors.HTTPStatusMisdirectedRequestError.ID)

	err = errors.FromHTTPStatusCode(http.StatusNetworkAuthenticationRequired)
	suite.Assert().Truef(errors.Is(err, errors.HTTPStatusNetworkAuthenticationRequiredError), "err should match a %s", errors.HTTPStatusNetworkAuthenticationRequiredError.ID)

	err = errors.FromHTTPStatusCode(http.StatusNotAcceptable)
	suite.Assert().Truef(errors.Is(err, errors.HTTPStatusNotAcceptableError), "err should match a %s", errors.HTTPStatusNotAcceptableError.ID)

	err = errors.FromHTTPStatusCode(http.StatusNotExtended)
	suite.Assert().Truef(errors.Is(err, errors.HTTPStatusNotExtendedError), "err should match a %s", errors.HTTPStatusNotExtendedError.ID)

	err = errors.FromHTTPStatusCode(http.StatusPaymentRequired)
	suite.Assert().Truef(errors.Is(err, errors.HTTPStatusPaymentRequiredError), "err should match a %s", errors.HTTPStatusPaymentRequiredError.ID)

	err = errors.FromHTTPStatusCode(http.StatusPreconditionFailed)
	suite.Assert().Truef(errors.Is(err, errors.HTTPStatusPreconditionFailedError), "err should match a %s", errors.HTTPStatusPreconditionFailedError.ID)

	err = errors.FromHTTPStatusCode(http.StatusPreconditionRequired)
	suite.Assert().Truef(errors.Is(err, errors.HTTPStatusPreconditionRequiredError), "err should match a %s", errors.HTTPStatusPreconditionRequiredError.ID)

	err = errors.FromHTTPStatusCode(http.StatusProxyAuthRequired)
	suite.Assert().Truef(errors.Is(err, errors.HTTPStatusProxyAuthRequiredError), "err should match a %s", errors.HTTPStatusProxyAuthRequiredError.ID)

	err = errors.FromHTTPStatusCode(http.StatusRequestEntityTooLarge)
	suite.Assert().Truef(errors.Is(err, errors.HTTPStatusRequestEntityTooLargeError), "err should match a %s", errors.HTTPStatusRequestEntityTooLargeError.ID)

	err = errors.FromHTTPStatusCode(http.StatusRequestHeaderFieldsTooLarge)
	suite.Assert().Truef(errors.Is(err, errors.HTTPStatusRequestHeaderFieldsTooLargeError), "err should match a %s", errors.HTTPStatusRequestHeaderFieldsTooLargeError.ID)

	err = errors.FromHTTPStatusCode(http.StatusRequestTimeout)
	suite.Assert().Truef(errors.Is(err, errors.HTTPStatusRequestTimeoutError), "err should match a %s", errors.HTTPStatusRequestTimeoutError.ID)

	err = errors.FromHTTPStatusCode(http.StatusRequestURITooLong)
	suite.Assert().Truef(errors.Is(err, errors.HTTPStatusRequestURITooLongError), "err should match a %s", errors.HTTPStatusRequestURITooLongError.ID)

	err = errors.FromHTTPStatusCode(http.StatusRequestedRangeNotSatisfiable)
	suite.Assert().Truef(errors.Is(err, errors.HTTPStatusRequestedRangeNotSatisfiableError), "err should match a %s", errors.HTTPStatusRequestedRangeNotSatisfiableError.ID)

	err = errors.FromHTTPStatusCode(http.StatusTeapot)
	suite.Assert().Truef(errors.Is(err, errors.HTTPStatusTeapotError), "err should match a %s", errors.HTTPStatusTeapotError.ID)

	err = errors.FromHTTPStatusCode(http.StatusTooEarly)
	suite.Assert().Truef(errors.Is(err, errors.HTTPStatusTooEarlyError), "err should match a %s", errors.HTTPStatusTooEarlyError.ID)

	err = errors.FromHTTPStatusCode(http.StatusTooManyRequests)
	suite.Assert().Truef(errors.Is(err, errors.HTTPStatusTooManyRequestsError), "err should match a %s", errors.HTTPStatusTooManyRequestsError.ID)

	err = errors.FromHTTPStatusCode(http.StatusUnavailableForLegalReasons)
	suite.Assert().Truef(errors.Is(err, errors.HTTPStatusUnavailableForLegalReasonsError), "err should match a %s", errors.HTTPStatusUnavailableForLegalReasonsError.ID)

	err = errors.FromHTTPStatusCode(http.StatusUnprocessableEntity)
	suite.Assert().Truef(errors.Is(err, errors.HTTPStatusUnprocessableEntityError), "err should match a %s", errors.HTTPStatusUnprocessableEntityError.ID)

	err = errors.FromHTTPStatusCode(http.StatusUnsupportedMediaType)
	suite.Assert().Truef(errors.Is(err, errors.HTTPStatusUnsupportedMediaTypeError), "err should match a %s", errors.HTTPStatusUnsupportedMediaTypeError.ID)

	err = errors.FromHTTPStatusCode(http.StatusUpgradeRequired)
	suite.Assert().Truef(errors.Is(err, errors.HTTPStatusUpgradeRequiredError), "err should match a %s", errors.HTTPStatusUpgradeRequiredError.ID)

	err = errors.FromHTTPStatusCode(http.StatusUseProxy)
	suite.Assert().Truef(errors.Is(err, errors.HTTPStatusUseProxyError), "err should match a %s", errors.HTTPStatusUseProxyError.ID)

	err = errors.FromHTTPStatusCode(http.StatusVariantAlsoNegotiates)
	suite.Assert().Truef(errors.Is(err, errors.HTTPStatusVariantAlsoNegotiatesError), "err should match a %s", errors.HTTPStatusVariantAlsoNegotiatesError.ID)

	err = errors.FromHTTPStatusCode(1234)
	suite.Assert().True(errors.Is(err, errors.Error{ID: "error.http.1234"}), "err should match a error.http.1234")
	suite.Require().True(errors.As(err, &details), "error should be a error.Error")
	suite.Assert().Equal(1234, details.Code, "Error code should be 1234")
}
