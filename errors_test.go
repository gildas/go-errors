package errors_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"regexp"
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
	err := errors.NotFound.With("key")
	suite.Require().NotNil(err, "err should not be nil")
	suite.Assert().True(errors.Is(err, errors.NotFound), "err should match a NotFoundError (pointer)")

	var details errors.Error
	suite.Require().True(errors.As(err, &details), "err should contain an errors.Error")
	suite.Assert().Equal("key", details.What)

	err = errors.ArgumentMissing.With("key")
	suite.Require().NotNil(err, "err should not be nil")
	suite.Assert().True(errors.Is(err, errors.ArgumentMissing), "err should match a ArgumentMissingError (pointer)")
	suite.Require().True(errors.As(err, &details), "err should contain an errors.Error")
	suite.Assert().Equal("key", details.What)
}

func (suite *ErrorsSuite) TestCanTellContainsAnError() {
	err := errors.NotFound.With("key")
	suite.Require().NotNil(err, "err should not be nil")
	var details errors.Error
	suite.Assert().True(errors.As(err, &details), "err should contain an errors.Error")
}

func (suite *ErrorsSuite) TestCanTellDoesNotContainAnError() {
	err := errors.Errorf("Houston, we have a problem")
	suite.Require().NotNil(err, "err should not be nil")
	suite.Assert().False(errors.Is(err, errors.Error{}), "err should not contain an errors.Error")
	var inner errors.Error
	suite.Assert().False(errors.As(err, &inner), "err should not contain an errors.Error")
}

func (suite *ErrorsSuite) TestCanWrap() {
	err := errors.Errorf("Houston, we have a problem")
	wrapped := errors.NotImplemented.Wrap(err)
	suite.Require().NotNil(wrapped)
	suite.Assert().True(errors.Is(wrapped, errors.NotImplemented), "wrapped err should be a Not Implemented Error")
	wrapped = errors.NotImplemented.Wrap(nil)
	suite.Require().Nil(wrapped, "Wrapped error of nil should be nil")
}

func (suite *ErrorsSuite) TestCanUnwrap() {
	err := errors.JSONUnmarshalError.Wrap(errors.New("Houston, we have a problem"))
	unwrapped := errors.Unwrap(err)
	suite.Assert().Equal("Houston, we have a problem", unwrapped.Error())
}

func (suite *ErrorsSuite) TestCanExtractError() {
	err := errors.JSONUnmarshalError.Wrap(errors.Unsupported.With("genre", "funky"))

	details, found := errors.Unsupported.Extract(err)
	suite.Require().True(found, "Error does not contain an Unsupported Error")
	suite.Require().Equal(errors.Unsupported.ID, details.ID)
	suite.Assert().Equal("genre", details.What)
	suite.Require().NotNil(details.Value, "Error Value should not be nil")
	value, ok := details.Value.(string)
	suite.Require().True(ok, "Value should be a string")
	suite.Assert().Equal("funky", value)
}

func (suite *ErrorsSuite) TestCanUnwrapJSONError() {
	var payload struct {
		Value string `json:"value"`
	}

	jsonerr := json.Unmarshal([]byte(`{"value": 0`), &payload)
	suite.Require().NotNil(jsonerr)
	suite.Assert().Equal("unexpected end of JSON input", jsonerr.Error())

	err := errors.JSONUnmarshalError.Wrap(jsonerr)
	suite.Require().NotNil(err)
	suite.Assert().Equal("JSON failed to unmarshal data: unexpected end of JSON input", err.Error())

	cause := errors.Unwrap(err)
	suite.Assert().Equal("unexpected end of JSON input", cause.Error())
}

func (suite *ErrorsSuite) TestFailsWithNonErrorTarget() {
	suite.Assert().False(errors.NotFound.Is(errors.New("Hello")), "Error should not be a pkg/errors.fundamental")
}

func (suite *ErrorsSuite) TestWrappers() {
	var err error

	err = errors.New("Hello World")
	suite.Assert().NotNil(err)

	err = errors.Errorf("Hello World")
	suite.Assert().NotNil(err)

	err = errors.WithMessage(errors.NotFound, "Hello World")
	suite.Assert().NotNil(err)

	err = errors.WithMessagef(errors.NotFound, "Hello World")
	suite.Assert().NotNil(err)

	err = errors.Wrap(errors.NotFound, "Hello World")
	suite.Assert().NotNil(err)

	err = errors.Wrapf(errors.NotFound, "Hello World")
	suite.Assert().NotNil(err)

	unwrapped := errors.Unwrap(err)
	suite.Assert().NotNil(unwrapped)

	suite.Assert().True(errors.Is(err, errors.NotFound), "err should be of the same type as NotFoundError")

	var inner errors.Error
	suite.Assert().True(errors.As(err, &inner), "Inner Error should be an errors.Error")
}

func ExampleError() {
	sentinel := errors.NewSentinel(500, "error.test.custom", "Test Error")
	err := sentinel.New()
	if err != nil {
		fmt.Println(err)

		var details errors.Error
		if errors.As(err, &details) {
			fmt.Println(details.ID)
		}
	}
	// Output:
	// Test Error
	// error.test.custom
}

func ExampleError_Format_default() {
	err := errors.NotImplemented.WithStack()
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	// Output:
	// Not Implemented
}

func ExampleError_Format_quoted() {
	err := errors.NotImplemented.WithStack()
	if err != nil {
		fmt.Printf("%q", err)
	}
	// Output:
	// "Not Implemented"
}

func ExampleError_Format_withStack() {
	output := CaptureStdout(func() {
		err := errors.NotImplemented.WithStack()
		if err != nil {
			fmt.Printf("%+v", err)
		}
	})
	// remove the path of each file and line numbers as they change for each Go deployments
	lines := strings.Split(output, "\n")
	simplifier := regexp.MustCompile(`\s*(.*/)?(.*):[0-9]+`)
	for i := 0; i < len(lines); i++ {
		line := lines[i]
	  lines[i] = simplifier.ReplaceAllString(line, "  ${2}")
	}
	// we also do not care about last line that is machine dependent
	fmt.Println(strings.Join(lines[0:len(lines)-1], "\n"))
	// Output:
	// Not Implemented
  // github.com/gildas/go-errors_test.ExampleError_Format_withStack.func1
  //   errors_test.go
  // github.com/gildas/go-errors_test.CaptureStdout
  //   errors_test.go
  // github.com/gildas/go-errors_test.ExampleError_Format_withStack
  //   errors_test.go
  // testing.runExample
  //   run_example.go
  // testing.runExamples
  //   example.go
  // testing.(*M).Run
  //   testing.go
  // main.main
  //   _testmain.go
  // runtime.main
  //   proc.go
  // runtime.goexit
}

func ExampleError_Format_gosyntax() {
	output := CaptureStdout(func() {
		err := errors.NotImplemented.WithStack()
		if err != nil {
			fmt.Printf("%#v\n", err)
		}
	})
	// remove the line numbers from the stack trace as they change when the code is changed
	simplifier := regexp.MustCompile(`\.go:[0-9]+`)
	// we also do not care about the last file which is machine dependent
	noasm := regexp.MustCompile(`, asm_.*.s:[0-9]+`)
	fmt.Println(noasm.ReplaceAllString(simplifier.ReplaceAllString(output, ".go"), ""))
	// Output:
	// errors.Error{Code:501, ID:"error.notimplemented", Text:"Not Implemented", What:"", Value:<nil>, Cause:<nil>, Stack:[]errors.StackFrame{errors_test.go, errors_test.go, errors_test.go, run_example.go, example.go, testing.go, _testmain.go, proc.go}}
}

func ExampleError_With() {
	err := errors.ArgumentMissing.With("key")
	if err != nil {
		fmt.Println(err)

		var details errors.Error
		if errors.As(err, &details) {
			fmt.Println(details.ID)
		}
	}
	// Output:
	// Argument key is missing
	// error.argument.missing
}

func ExampleError_With_value() {
	err := errors.ArgumentInvalid.With("key", "value")
	if err != nil {
		fmt.Println(err)

		var details errors.Error
		if errors.As(err, &details) {
			fmt.Println(details.ID)
		}
	}
	// Output:
	// Argument key is invalid (value: value)
	// error.argument.invalid
}

func ExampleError_With_array() {
	err := errors.ArgumentInvalid.With("key", []string{"value1", "value2"})
	if err != nil {
		fmt.Println(err)

		var details errors.Error
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
		var details errors.Error
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

	err = errors.FromHTTPStatusCode(http.StatusBadGateway)
	suite.Assert().Truef(errors.Is(err, errors.HTTPBadGateway), "err should match a %s", errors.HTTPBadGateway.ID)

	err = errors.FromHTTPStatusCode(http.StatusBadRequest)
	suite.Assert().Truef(errors.Is(err, errors.HTTPBadRequest), "err should match a %s", errors.HTTPBadRequest.ID)

	err = errors.FromHTTPStatusCode(http.StatusForbidden)
	suite.Assert().Truef(errors.Is(err, errors.HTTPForbidden), "err should match a %s", errors.HTTPForbidden.ID)

	err = errors.FromHTTPStatusCode(http.StatusInternalServerError)
	suite.Assert().Truef(errors.Is(err, errors.HTTPInternalServerError), "err should match a %s", errors.HTTPInternalServerError.ID)

	err = errors.FromHTTPStatusCode(http.StatusMethodNotAllowed)
	suite.Assert().Truef(errors.Is(err, errors.HTTPMethodNotAllowed), "err should match a %s", errors.HTTPMethodNotAllowed.ID)

	err = errors.FromHTTPStatusCode(http.StatusNotFound)
	suite.Assert().Truef(errors.Is(err, errors.HTTPNotFound), "err should match a %s", errors.HTTPNotFound.ID)

	err = errors.FromHTTPStatusCode(http.StatusNotImplemented)
	suite.Assert().Truef(errors.Is(err, errors.HTTPNotImplemented), "err should match a %s", errors.HTTPNotImplemented.ID)

	err = errors.FromHTTPStatusCode(http.StatusServiceUnavailable)
	suite.Assert().Truef(errors.Is(err, errors.HTTPServiceUnavailable), "err should match a %s", errors.HTTPServiceUnavailable.ID)

	err = errors.FromHTTPStatusCode(http.StatusUnauthorized)
	suite.Assert().Truef(errors.Is(err, errors.HTTPUnauthorized), "err should match a %s", errors.HTTPUnauthorized.ID)

	err = errors.FromHTTPStatusCode(http.StatusConflict)
	suite.Assert().Truef(errors.Is(err, errors.HTTPStatusConflict), "err should match a %s", errors.HTTPStatusConflict.ID)

	err = errors.FromHTTPStatusCode(http.StatusExpectationFailed)
	suite.Assert().Truef(errors.Is(err, errors.HTTPStatusExpectationFailed), "err should match a %s", errors.HTTPStatusExpectationFailed.ID)

	err = errors.FromHTTPStatusCode(http.StatusFailedDependency)
	suite.Assert().Truef(errors.Is(err, errors.HTTPStatusFailedDependency), "err should match a %s", errors.HTTPStatusFailedDependency.ID)

	err = errors.FromHTTPStatusCode(http.StatusGatewayTimeout)
	suite.Assert().Truef(errors.Is(err, errors.HTTPStatusGatewayTimeout), "err should match a %s", errors.HTTPStatusGatewayTimeout.ID)

	err = errors.FromHTTPStatusCode(http.StatusGone)
	suite.Assert().Truef(errors.Is(err, errors.HTTPStatusGone), "err should match a %s", errors.HTTPStatusGone.ID)

	err = errors.FromHTTPStatusCode(http.StatusHTTPVersionNotSupported)
	suite.Assert().Truef(errors.Is(err, errors.HTTPStatusHTTPVersionNotSupported), "err should match a %s", errors.HTTPStatusHTTPVersionNotSupported.ID)

	err = errors.FromHTTPStatusCode(http.StatusInsufficientStorage)
	suite.Assert().Truef(errors.Is(err, errors.HTTPStatusInsufficientStorage), "err should match a %s", errors.HTTPStatusInsufficientStorage.ID)

	err = errors.FromHTTPStatusCode(http.StatusLengthRequired)
	suite.Assert().Truef(errors.Is(err, errors.HTTPStatusLengthRequired), "err should match a %s", errors.HTTPStatusLengthRequired.ID)

	err = errors.FromHTTPStatusCode(http.StatusLocked)
	suite.Assert().Truef(errors.Is(err, errors.HTTPStatusLocked), "err should match a %s", errors.HTTPStatusLocked.ID)

	err = errors.FromHTTPStatusCode(http.StatusLoopDetected)
	suite.Assert().Truef(errors.Is(err, errors.HTTPStatusLoopDetected), "err should match a %s", errors.HTTPStatusLoopDetected.ID)

	err = errors.FromHTTPStatusCode(http.StatusMisdirectedRequest)
	suite.Assert().Truef(errors.Is(err, errors.HTTPStatusMisdirectedRequest), "err should match a %s", errors.HTTPStatusMisdirectedRequest.ID)

	err = errors.FromHTTPStatusCode(http.StatusNetworkAuthenticationRequired)
	suite.Assert().Truef(errors.Is(err, errors.HTTPStatusNetworkAuthenticationRequired), "err should match a %s", errors.HTTPStatusNetworkAuthenticationRequired.ID)

	err = errors.FromHTTPStatusCode(http.StatusNotAcceptable)
	suite.Assert().Truef(errors.Is(err, errors.HTTPStatusNotAcceptable), "err should match a %s", errors.HTTPStatusNotAcceptable.ID)

	err = errors.FromHTTPStatusCode(http.StatusNotExtended)
	suite.Assert().Truef(errors.Is(err, errors.HTTPStatusNotExtended), "err should match a %s", errors.HTTPStatusNotExtended.ID)

	err = errors.FromHTTPStatusCode(http.StatusPaymentRequired)
	suite.Assert().Truef(errors.Is(err, errors.HTTPStatusPaymentRequired), "err should match a %s", errors.HTTPStatusPaymentRequired.ID)

	err = errors.FromHTTPStatusCode(http.StatusPreconditionFailed)
	suite.Assert().Truef(errors.Is(err, errors.HTTPStatusPreconditionFailed), "err should match a %s", errors.HTTPStatusPreconditionFailed.ID)

	err = errors.FromHTTPStatusCode(http.StatusPreconditionRequired)
	suite.Assert().Truef(errors.Is(err, errors.HTTPStatusPreconditionRequired), "err should match a %s", errors.HTTPStatusPreconditionRequired.ID)

	err = errors.FromHTTPStatusCode(http.StatusProxyAuthRequired)
	suite.Assert().Truef(errors.Is(err, errors.HTTPStatusProxyAuthRequired), "err should match a %s", errors.HTTPStatusProxyAuthRequired.ID)

	err = errors.FromHTTPStatusCode(http.StatusRequestEntityTooLarge)
	suite.Assert().Truef(errors.Is(err, errors.HTTPStatusRequestEntityTooLarge), "err should match a %s", errors.HTTPStatusRequestEntityTooLarge.ID)

	err = errors.FromHTTPStatusCode(http.StatusRequestHeaderFieldsTooLarge)
	suite.Assert().Truef(errors.Is(err, errors.HTTPStatusRequestHeaderFieldsTooLarge), "err should match a %s", errors.HTTPStatusRequestHeaderFieldsTooLarge.ID)

	err = errors.FromHTTPStatusCode(http.StatusRequestTimeout)
	suite.Assert().Truef(errors.Is(err, errors.HTTPStatusRequestTimeout), "err should match a %s", errors.HTTPStatusRequestTimeout.ID)

	err = errors.FromHTTPStatusCode(http.StatusRequestURITooLong)
	suite.Assert().Truef(errors.Is(err, errors.HTTPStatusRequestURITooLong), "err should match a %s", errors.HTTPStatusRequestURITooLong.ID)

	err = errors.FromHTTPStatusCode(http.StatusRequestedRangeNotSatisfiable)
	suite.Assert().Truef(errors.Is(err, errors.HTTPStatusRequestedRangeNotSatisfiable), "err should match a %s", errors.HTTPStatusRequestedRangeNotSatisfiable.ID)

	err = errors.FromHTTPStatusCode(http.StatusTeapot)
	suite.Assert().Truef(errors.Is(err, errors.HTTPStatusTeapot), "err should match a %s", errors.HTTPStatusTeapot.ID)

	err = errors.FromHTTPStatusCode(http.StatusTooEarly)
	suite.Assert().Truef(errors.Is(err, errors.HTTPStatusTooEarly), "err should match a %s", errors.HTTPStatusTooEarly.ID)

	err = errors.FromHTTPStatusCode(http.StatusTooManyRequests)
	suite.Assert().Truef(errors.Is(err, errors.HTTPStatusTooManyRequests), "err should match a %s", errors.HTTPStatusTooManyRequests.ID)

	err = errors.FromHTTPStatusCode(http.StatusUnavailableForLegalReasons)
	suite.Assert().Truef(errors.Is(err, errors.HTTPStatusUnavailableForLegalReasons), "err should match a %s", errors.HTTPStatusUnavailableForLegalReasons.ID)

	err = errors.FromHTTPStatusCode(http.StatusUnprocessableEntity)
	suite.Assert().Truef(errors.Is(err, errors.HTTPStatusUnprocessableEntity), "err should match a %s", errors.HTTPStatusUnprocessableEntity.ID)

	err = errors.FromHTTPStatusCode(http.StatusUnsupportedMediaType)
	suite.Assert().Truef(errors.Is(err, errors.HTTPStatusUnsupportedMediaType), "err should match a %s", errors.HTTPStatusUnsupportedMediaType.ID)

	err = errors.FromHTTPStatusCode(http.StatusUpgradeRequired)
	suite.Assert().Truef(errors.Is(err, errors.HTTPStatusUpgradeRequired), "err should match a %s", errors.HTTPStatusUpgradeRequired.ID)

	err = errors.FromHTTPStatusCode(http.StatusUseProxy)
	suite.Assert().Truef(errors.Is(err, errors.HTTPStatusUseProxy), "err should match a %s", errors.HTTPStatusUseProxy.ID)

	err = errors.FromHTTPStatusCode(http.StatusVariantAlsoNegotiates)
	suite.Assert().Truef(errors.Is(err, errors.HTTPStatusVariantAlsoNegotiates), "err should match a %s", errors.HTTPStatusVariantAlsoNegotiates.ID)

	err = errors.FromHTTPStatusCode(1234)
	suite.Assert().True(errors.Is(err, errors.Error{ID: "error.http.1234"}), "err should match a error.http.1234")

	var details errors.Error
	suite.Require().True(errors.As(err, &details), "error should be a error.Error")
	suite.Assert().Equal(1234, details.Code, "Error code should be 1234")
}

func CaptureStdout(f func()) string {
	reader, writer, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	stdout := os.Stdout
	os.Stdout = writer
	defer func() {
		os.Stdout = stdout
	}()

	f()
	writer.Close()

	output := bytes.Buffer{}
	_, _ = io.Copy(&output, reader)
	return output.String()
}
