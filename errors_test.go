package errors_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
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
	suite.Name = strings.TrimSuffix(reflect.TypeOf(suite).Elem().Name(), "Suite")
}

func (suite *ErrorsSuite) TestCanCreate() {
	err := errors.Error{}
	suite.Assert().Equal("runtime error", err.Error())
}

func (suite *ErrorsSuite) TestCanTellIsError() {
	err := errors.NotFound.With("key")
	suite.Require().NotNil(err, "err should not be nil")
	suite.Assert().ErrorIs(err, errors.Error{}, "err should be an errors.Error")
	suite.Assert().True(errors.Is(err, errors.NotFound), "err should match a NotFoundError")
	suite.Assert().True(errors.NotFound.Is(err), "err should match a NotFoundError")

	err = fmt.Errorf("simple error")
	suite.Require().NotNil(err, "err should not be nil")
	suite.Assert().NotErrorIs(err, errors.Error{}, "err should not be an errors.Error")
	suite.Assert().False(errors.Is(err, errors.NotFound), "err should not match an NotFoundError")
	suite.Assert().False(errors.NotFound.Is(err), "err should not match an NotFoundError")
}

func (suite *ErrorsSuite) TestCanConvertToError() {
	err := errors.NotFound.With("key")
	suite.Require().NotNil(err, "err should not be nil")

	var details *errors.Error
	suite.Require().ErrorAs(err, &details, "err should contain an errors.Error")
	suite.Assert().Equal("key", details.What)
}

func (suite *ErrorsSuite) TestCanConvertToSpecificError() {
	err := errors.NotFound.With("key1").(errors.Error).Wrap(errors.ArgumentInvalid.With("key2"))
	suite.Require().NotNil(err, "err should not be nil")

	suite.Assert().ErrorIs(err, errors.Error{}, "err should be an errors.Error")
	suite.Assert().ErrorIs(err, errors.NotFound, "err should match a NotFoundError")
	suite.Assert().ErrorIs(err, errors.ArgumentInvalid, "err should match an ArgumentInvalidError")

	details := errors.NotFound.Clone()
	suite.Require().ErrorAs(err, &details, "err should contain an errors.NotFound")
	suite.Assert().Equal(errors.NotFound.ID, details.ID)
	suite.Assert().Equal("key1", details.What)
	suite.Assert().Len(errors.ArgumentInvalid.What, 0, "ArgumentInvalid should not have changed")

	details = errors.ArgumentInvalid.Clone()
	suite.Require().ErrorAs(err, &details, "err should contain an errors.ArgumentInvalid")
	suite.Assert().Equal(errors.ArgumentInvalid.ID, details.ID)
	suite.Assert().Equal("key2", details.What)
	suite.Assert().Len(errors.ArgumentInvalid.What, 0, "ArgumentInvalid should not have changed")

	details = errors.ArgumentMissing.Clone()
	suite.Require().False(errors.As(err, &details), "err should not contain an errors.ArgumentMissing")
}

func (suite *ErrorsSuite) TestShouldFailConvertingToUrlError() {
	// var err error = &url.Error{Op: "Get", URL: "https://bogus.acme.com", Err: fmt.Errorf("Houston, we have a problem")}
	var err error = errors.NotFound.With("key")

	var details *url.Error
	suite.Assert().False(errors.As(err, &details), "err should not contain an url.Error")
}

func (suite *ErrorsSuite) TestCanWrap() {
	wrapped := errors.NotImplemented.Wrap(errors.Errorf("Houston, we have a problem"))
	suite.Require().NotNil(wrapped)
	suite.Assert().True(errors.Is(wrapped, errors.NotImplemented), "wrapped err should be a Not Implemented Error")
	wrapped = errors.NotImplemented.Wrap(nil)
	suite.Require().Nil(wrapped, "Wrapped error of nil should be nil")
}

func (suite *ErrorsSuite) TestCanWrapNilError() {
	wrapped := errors.WrapErrors(nil, errors.NotImplemented.WithStack())
	suite.Require().Nil(wrapped, "Chained error of nil should be nil")
}

func (suite *ErrorsSuite) TestCanWrapErrorWithNilError() {
	wrapped := errors.WrapErrors(errors.NotImplemented.WithStack(), nil)
	suite.Require().Nil(wrapped, "Chained error of nil should be nil")
}

func (suite *ErrorsSuite) TestCanWrapOneError() {
	error1 := errors.ArgumentMissing.With("key1")
	wrapped := errors.WrapErrors(error1)
	suite.Assert().Equal(error1, wrapped, "Chained errors of one error should be the error")
}

func (suite *ErrorsSuite) TestCanWrapErrors() {
	error1 := errors.ArgumentMissing.With("key1")
	error2 := errors.NotFound.With("key", "key2")
	error3 := errors.EnvironmentMissing.With("key3")

	wrapped := errors.WrapErrors(error1, error2, error3)
	suite.Assert().ErrorIs(wrapped, errors.Error{}, "Chained errors should contain an errors.Error")
	suite.Assert().ErrorIs(wrapped, errors.ArgumentMissing, "Chained errors should match an ArgumentMissingError")
	suite.Assert().ErrorIs(wrapped, errors.NotFound, "Chained errors should match a NotFoundError")
	suite.Assert().ErrorIs(wrapped, errors.EnvironmentMissing, "Chained errors should match a EnvironmentMissingError")

	first, ok := wrapped.(errors.Error)
	suite.Require().True(ok, "The first error should be an errors.Error")
	suite.Assert().Equal(errors.ArgumentMissing.ID, first.ID, "The first error should be error1")
	suite.Require().NotNil(first.Cause, "The first error should have a cause")

	unwrapped := first.Unwrap()
	suite.Require().NotNil(unwrapped, "The first error should have a cause")

	second, ok := unwrapped.(errors.Error)
	suite.Require().True(ok, "The second error should be an errors.Error")
	suite.Assert().Equal(errors.NotFound.ID, second.ID, "The second error should be error2")
	suite.Assert().NotNil(second.Cause, "The second error should have a cause")

	unwrapped = second.Unwrap()
	suite.Require().NotNil(unwrapped, "The second error should have a cause")

	third, ok := unwrapped.(errors.Error)
	suite.Require().True(ok, "The third error should be an errors.Error")
	suite.Assert().Equal(errors.EnvironmentMissing.ID, third.ID, "The third error should be error3")
	suite.Assert().Nil(third.Cause, "The third error should not have a cause")

	unwrapped = third.Unwrap()
	suite.Require().Nil(unwrapped, "The third error should not have a cause")
}

func (suite *ErrorsSuite) TestCanWrapErrorsEndingWithNilError() {
	wrapped := errors.WrapErrors(errors.ArgumentMissing.With("key1"), errors.ArgumentMissing.With("key2"), errors.NotImplemented.WithStack(), nil)
	suite.Require().Nil(wrapped, "Chained error ending with nil should be nil")
}

func (suite *ErrorsSuite) TestCanWrapErrorsWithNilError() {
	error1 := errors.ArgumentMissing.With("key1")
	error2 := errors.NotFound.With("key", "key2")
	error3 := errors.EnvironmentMissing.With("key3")

	wrapped := errors.WrapErrors(error1, error2, nil, error3)
	suite.Assert().ErrorIs(wrapped, errors.Error{}, "Chained errors should contain an errors.Error")
	suite.Assert().ErrorIs(wrapped, errors.ArgumentMissing, "Chained errors should match an ArgumentMissingError")
	suite.Assert().ErrorIs(wrapped, errors.NotFound, "Chained errors should match a NotFoundError")
	suite.Assert().ErrorIs(wrapped, errors.EnvironmentMissing, "Chained errors should match a EnvironmentMissingError")

	first, ok := wrapped.(errors.Error)
	suite.Require().True(ok, "The first error should be an errors.Error")
	suite.Assert().Equal(errors.ArgumentMissing.ID, first.ID, "The first error should be error1")
	suite.Require().NotNil(first.Cause, "The first error should have a cause")

	unwrapped := first.Unwrap()
	suite.Require().NotNil(unwrapped, "The first error should have a cause")

	second, ok := unwrapped.(errors.Error)
	suite.Require().True(ok, "The second error should be an errors.Error")
	suite.Assert().Equal(errors.NotFound.ID, second.ID, "The second error should be error2")
	suite.Assert().NotNil(second.Cause, "The second error should have a cause")

	unwrapped = second.Unwrap()
	suite.Require().NotNil(unwrapped, "The second error should have a cause")

	third, ok := unwrapped.(errors.Error)
	suite.Require().True(ok, "The third error should be an errors.Error")
	suite.Assert().Equal(errors.EnvironmentMissing.ID, third.ID, "The third error should be error3")
	suite.Assert().Nil(third.Cause, "The third error should not have a cause")

	unwrapped = third.Unwrap()
	suite.Require().Nil(unwrapped, "The third error should not have a cause")
}

func (suite *ErrorsSuite) TestCanWrapWithBasicErrors() {
	error1 := errors.ArgumentMissing.With("key1")

	wrapped := errors.WrapErrors(fmt.Errorf("basic error"), error1)
	suite.Assert().ErrorIs(wrapped, errors.Error{}, "Chained errors should contain an errors.Error")
	suite.Assert().ErrorIs(wrapped, errors.ArgumentMissing, "Chained errors should match an ArgumentMissingError")

	first, ok := wrapped.(errors.Error)
	suite.Require().True(ok, "The first error should be an errors.Error")
	suite.Assert().Equal(http.StatusInternalServerError, first.Code, "The first error should have a 500 code")
	suite.Assert().Equal("error.runtime", first.ID, "The first error should have a runtime error id")
	suite.Assert().Equal("basic error", first.Text, "The first error should have the basic error as its Text")
	suite.Assert().Equal("basic error", first.Error())
	suite.Require().NotNil(first.Cause, "The first error should have a cause")

	unwrapped := first.Unwrap()
	suite.Require().NotNil(unwrapped, "The first error should have a cause")

	second, ok := unwrapped.(errors.Error)
	suite.Require().True(ok, "The second error should be an errors.Error")
	suite.Assert().Equal(errors.ArgumentMissing.ID, second.ID, "The second error should be error1")
	suite.Require().Nil(second.Cause, "The second error should not have a cause")
}

func (suite *ErrorsSuite) TestCanWrapWithURLErrors() {
	error1 := errors.ArgumentMissing.With("key1")
	urlError := &url.Error{
		Op:  "Get",
		URL: "https://example.com/",
		Err: &net.OpError{
			Op:  "remote error",
			Net: "",
			Err: fmt.Errorf("tls handshake failure"),
		},
	}

	wrapped := errors.WrapErrors(error1, urlError)
	suite.Assert().ErrorIs(wrapped, errors.Error{}, "Chained errors should contain an errors.Error")
	suite.Assert().ErrorIs(wrapped, errors.ArgumentMissing, "Chained errors should match an ArgumentMissing")
	suite.Assert().ErrorIs(wrapped, urlError, "Chained errors should match an url.Error")

	first, ok := wrapped.(errors.Error)
	suite.Require().True(ok, "The first error should be an errors.Error")
	suite.Assert().Equal(errors.ArgumentMissing.ID, first.ID, "The first error should be error1")
	suite.Require().NotNil(first.Cause, "The first error should have a cause")

	unwrapped := first.Unwrap()
	suite.Require().NotNil(unwrapped, "The first error should have a cause")

	second, ok := unwrapped.(*url.Error)
	suite.Require().True(ok, "The second error should be an *url.Error")
	suite.Assert().Equal("Get", second.Op, "The second error should be a GET operation")

	unwrapped = second.Unwrap()
	suite.Require().NotNil(unwrapped, "The first error should have a cause")

	third, ok := unwrapped.(*net.OpError)
	suite.Require().True(ok, "The second error should be an *net.OpError")
	suite.Assert().Equal("remote error", third.Op, "The second error should be a remote error")
}

func (suite *ErrorsSuite) TestShouldAddStackOnlyOnce() {
	err := errors.NotImplemented
	suite.Assert().Empty(err.Stack, "The error should not have a stack")

	err1 := errors.WithStack(err)
	suite.Assert().NotEmpty(err1, "The error should have a stack")
	stack1 := err1.(errors.Error).Stack

	err2 := errors.WithStack(err1)
	suite.Assert().NotEmpty(err2, "The error should have a stack")
	stack2 := err2.(errors.Error).Stack
	suite.Assert().Equal(stack1, stack2, "The error should have the same stack")
}

func (suite *ErrorsSuite) TestCanRemoveStackTrace() {
	suite.Assert().Nil(errors.WithoutStack(nil))

	err := errors.WithoutStack(errors.NotImplemented.WithStack())
	suite.Require().NotNil(err)
	check, ok := err.(errors.Error)
	suite.Require().True(ok, "err should be an errors.Error")
	suite.Assert().Empty(check.Stack)

	err = errors.WithoutStack(fmt.Errorf("simple error"))
	suite.Require().NotNil(err)
}

func (suite *ErrorsSuite) TestCanUnwrap() {
	err := errors.JSONUnmarshalError.Wrap(errors.New("Houston, we have a problem"))
	unwrapped := errors.Unwrap(err)
	suite.Assert().Equal("Houston, we have a problem", unwrapped.Error())
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
	suite.Assert().Equal("JSON failed to unmarshal data\nCaused by:\n\tunexpected end of JSON input", err.Error())

	cause := errors.Unwrap(err)
	suite.Assert().Equal("unexpected end of JSON input", cause.Error())
}

func (suite *ErrorsSuite) TestFailsWithNonErrorTarget() {
	suite.Assert().False(errors.NotFound.Is(errors.New("Hello")), "Error should not be a pkg/errors.fundamental")
}

func (suite *ErrorsSuite) TestCanMarshalError() {
	expected := `{"type": "error", "id": "error.argument.invalid", "code": 400, "text": "Argument %s is invalid (value: %v)", "what": "key", "value": "value"}`
	testerr := errors.ArgumentInvalid.With("key", "value")
	payload, err := json.Marshal(testerr)
	suite.Require().Nil(err)
	suite.Assert().JSONEq(expected, string(payload))
}

func (suite *ErrorsSuite) TestCanMarshalErrorWithoutValue() {
	expected := `{"type": "error", "id": "error.argument.invalid", "code": 400, "text": "Argument %s is invalid (value: %v)", "what": "key"}`
	testerr := errors.ArgumentInvalid.With("key")
	payload, err := json.Marshal(testerr)
	suite.Require().Nil(err)
	suite.Assert().JSONEq(expected, string(payload))
}

func (suite *ErrorsSuite) TestCanMarshalErrorWithCause() {
	expected := `{
		"type": "error",
		"id": "error.argument.invalid",
		"code": 400,
		"text": "Argument %s is invalid (value: %v)",
		"what": "key",
		"value": "value",
		"cause": {
			"type": "error",
			"code": 400,
			"id": "error.http.request",
			"text": "Bad Request. %s"
		}
	}`
	var cause error = errors.FromHTTPStatusCode(400)
	testerr := errors.WrapErrors(errors.ArgumentInvalid.With("key", "value"), cause)
	payload, err := json.Marshal(testerr)
	suite.Require().Nil(err)
	suite.Assert().JSONEq(expected, string(payload))
}

func (suite *ErrorsSuite) TestCanMarshalErrorWithURLErrorCause01() {
	expected := `{
		"type": "error",
		"id": "error.argument.invalid",
		"code": 400,
		"text": "Argument %s is invalid (value: %v)",
		"what": "key",
		"value": "value",
		"cause": {
			"type": "error",
			"code": 500,
			"id": "error.runtime.url.Error",
			"text": "Get \"https://example.com/\": remote error: tls handshake failure"
		}
	}`
	var cause error = &url.Error{
		Op:  "Get",
		URL: "https://example.com/",
		Err: &net.OpError{
			Op:  "remote error",
			Net: "",
			Err: fmt.Errorf("tls handshake failure"),
		},
	}
	testerr := errors.WrapErrors(errors.ArgumentInvalid.With("key", "value"), cause)
	payload, err := json.Marshal(testerr)
	suite.Require().Nil(err)
	suite.Assert().JSONEq(expected, string(payload))
}

func (suite *ErrorsSuite) TestCanMarshalErrorWithURLErrorCause02() {
	expected := `{
		"type": "error",
		"id": "error.argument.invalid",
		"code": 400,
		"text": "Argument %s is invalid (value: %v)",
		"what": "key",
		"value": "value",
		"cause": {
			"type": "error",
			"code": 500,
			"id": "error.runtime.url.Error",
			"text": "Get \"https://bogus.example.com/\": Dial tcp: lookup bogus.example.com on 208.67.222.222:53: no such host"
		}
	}`
	var cause error = &url.Error{
		Op:  "Get",
		URL: "https://bogus.example.com/",
		Err: &net.OpError{
			Op:     "Dial",
			Net:    "tcp",
			Source: nil,
			Addr:   nil,
			Err: &net.DNSError{
				Err:         "no such host",
				Name:        "bogus.example.com",
				Server:      "208.67.222.222:53",
				IsTimeout:   false,
				IsTemporary: false,
				IsNotFound:  true,
			},
		},
	}
	testerr := errors.WrapErrors(errors.ArgumentInvalid.With("key", "value"), cause)
	payload, err := json.Marshal(testerr)
	suite.Require().Nil(err)
	suite.Assert().JSONEq(expected, string(payload))
}

func (suite *ErrorsSuite) TestCanMarshalErrorWithManyCauses() {
	expected := `{
		"type": "error",
		"id": "error.argument.invalid",
		"code": 400,
		"text": "Argument %s is invalid (value: %v)",
		"what": "key",
		"value": "value",
		"cause": {
			"type": "error",
			"id": "error.argument.missing",
			"code": 400,
			"text": "Argument %s is missing",
			"what": "key",
			"cause": {
				"type": "error",
				"id": "error.runtime",
				"code": 500,
				"text": "some obscure error"
			}
		}
	}`
	testerr := errors.WrapErrors(
		errors.ArgumentInvalid.With("key", "value"),
		errors.ArgumentMissing.With("key"),
		nil,
		fmt.Errorf("some obscure error"),
	)
	payload, err := json.Marshal(testerr)
	suite.Require().Nil(err)
	suite.Assert().JSONEq(expected, string(payload))
}

func (suite *ErrorsSuite) TestCanUnmarshalError() {
	payload := `{"type": "error", "id": "error.argument.invalid", "code": 400, "text": "Argument %s is invalid (value: %v)", "what": "key", "value": "value"}`
	testerr := errors.Error{}
	err := json.Unmarshal([]byte(payload), &testerr)
	suite.Require().Nil(err)
	suite.Assert().Equal(400, testerr.Code)
	suite.Assert().Equal("error.argument.invalid", testerr.ID)
}

func (suite *ErrorsSuite) TestCanUnmarshalErrorWithErrorCause() {
	payload := `{
		"type": "error",
		"id": "error.argument.invalid",
		"code": 400,
		"text": "Argument %s is invalid (value: %v)",
		"what": "key",
		"value": "value",
		"cause": {
			"type": "error",
			"code": 400,
			"id": "error.http.request",
			"text": "Bad Request. %s"
		}
	}`
	testerr := errors.Error{}
	err := json.Unmarshal([]byte(payload), &testerr)
	suite.Require().Nil(err)
	suite.Assert().Equal(400, testerr.Code)
	suite.Assert().Equal("error.argument.invalid", testerr.ID)
	suite.Assert().Equal("Argument %s is invalid (value: %v)", testerr.Text)
	suite.Require().NotNil(testerr.Cause, "error should have a cause")
	var cause *errors.Error
	suite.Require().ErrorAs(testerr.Cause, &cause, "cause should be an errors.Error")
	suite.Assert().Equal("error.http.request", cause.ID)
	suite.Assert().Equal(400, cause.Code)
	suite.Assert().Equal("Bad Request. %s", cause.Text)
}

func (suite *ErrorsSuite) TestCanUnmarshalErrorWithTextCause() {
	payload := `{
		"type": "error",
		"id": "error.argument.invalid",
		"code": 400,
		"text": "Argument %s is invalid (value: %v)",
		"what": "key",
		"value": "value",
		"cause": {
			"type": "error",
			"code": 500,
			"id": "error.runtime.url.Error",
			"text": "Get \"https://bogus.example.com/\": Dial tcp: lookup bogus.example.com on 208.67.222.222:53: no such host"
		}
	}`
	testerr := errors.Error{}
	err := json.Unmarshal([]byte(payload), &testerr)
	suite.Require().Nil(err)
	suite.Assert().Equal(400, testerr.Code)
	suite.Assert().Equal("error.argument.invalid", testerr.ID)
	suite.Assert().Equal("Argument %s is invalid (value: %v)", testerr.Text)
	suite.Require().NotNil(testerr.Cause, "error should have a cause")
	var cause *errors.Error
	suite.Require().ErrorAs(testerr.Cause, &cause, "causes[0] should be an errors.Error")
	suite.Assert().Equal("error.runtime.url.Error", cause.ID)
	suite.Assert().Equal(500, cause.Code)
	suite.Assert().Equal(`Get "https://bogus.example.com/": Dial tcp: lookup bogus.example.com on 208.67.222.222:53: no such host`, cause.Text)
}

func (suite *ErrorsSuite) TestCanUnmarshalErrorWithManyCauses() {
	payload := `{
		"type": "error",
		"id": "error.argument.invalid",
		"code": 400,
		"text": "Argument %s is invalid (value: %v)",
		"what": "key",
		"value": "value",
		"cause": {
			"type": "error",
			"id": "error.argument.missing",
			"code": 400,
			"text": "Argument %s is missing",
			"what": "key",
			"cause": {
				"type": "error",
				"id": "error.runtime",
				"code": 500,
				"text": "some obscure error"
			}
		}
	}`
	var testerr errors.Error

	err := json.Unmarshal([]byte(payload), &testerr)
	suite.Require().Nil(err)
	suite.Require().NotNil(testerr.Cause, "error should have a cause")

	// Unwrap the causes with the standard errors.As()/errors.Unwrap() functions
	cause0 := errors.ArgumentInvalid.Clone()
	suite.Require().ErrorAs(testerr, &cause0, "cause0 should be an errors.Error")
	suite.Assert().Equal("error.argument.invalid", cause0.ID)
	suite.Assert().Equal(400, cause0.Code)
	suite.Assert().Equal("Argument %s is invalid (value: %v)", cause0.Text)
	suite.Assert().Equal("key", cause0.What)
	suite.Assert().Equal("value", cause0.Value)

	cause1 := errors.ArgumentMissing.Clone()
	suite.Require().ErrorAs(testerr, &cause1, "causes1 should be an errors.Error")
	suite.Assert().Equal("error.argument.missing", cause1.ID)
	suite.Assert().Equal(400, cause1.Code)
	suite.Assert().Equal("Argument %s is missing", cause1.Text)
	suite.Assert().Equal("key", cause1.What)

	cause2 := errors.RuntimeError.Clone()
	suite.Require().ErrorAs(testerr, &cause2, "causes2 should be an errors.Error")
	suite.Assert().Equal("error.runtime", cause2.ID)
	suite.Assert().Equal(500, cause2.Code)
	suite.Assert().Equal("some obscure error", cause2.Text)

	// Now let's verify the chain of causes
	suite.Assert().Equal(errors.ArgumentInvalid.ID, testerr.ID, "the main errir should be an ArgumentInvalid")

	c1, ok := testerr.Cause.(errors.Error)
	suite.Require().True(ok, "the first cause should be an errors.Error")
	suite.Assert().Equal(errors.ArgumentMissing.ID, c1.ID, "the first cause should be an ArgumentMissing")

	c2, ok := c1.Cause.(errors.Error)
	suite.Require().True(ok, "the second cause should be an errors.Error")
	suite.Assert().Equal(errors.RuntimeError.ID, c2.ID, "the second cause should be a Runtime")
	suite.Assert().Equal("some obscure error", c2.Text)
}

func (suite *ErrorsSuite) TestFailsUnmarshallErrorWithWrongPayload() {
	payload := `{"type": "error", "id": 1000, "code": 400, "text": "Argument %s is invalid (value: %v)", "what": "key", "value": "value"}`
	testerr := errors.Error{}
	err := json.Unmarshal([]byte(payload), &testerr)
	suite.Require().NotNil(err)
	suite.Assert().True(errors.Is(err, errors.JSONUnmarshalError), "Error should be a JSONUnmarshalError")
	suite.Assert().Equal("json: cannot unmarshal number into Go struct field .id of type string", errors.Unwrap(err).Error())
}

func (suite *ErrorsSuite) TestFailsUnmarshallErrorWithWrongType() {
	payload := `{"type": "blob", "id": "error.argument.invalid", "code": 400, "text": "Argument %s is invalid (value: %v)", "what": "key", "value": "value"}`
	testerr := errors.Error{}
	err := json.Unmarshal([]byte(payload), &testerr)
	suite.Require().NotNil(err)
	suite.Assert().True(errors.Is(err, errors.JSONUnmarshalError), "Error should be a JSONUnmarshalError")
	suite.Assert().True(errors.Is(err, errors.InvalidType), "Error should be an InvalidType")
	details := errors.InvalidType.Clone()
	suite.Require().ErrorAs(err, &details, "err should contain an errors.InvalidType")
	suite.Assert().Equal("error", details.What)
	value, ok := details.Value.(string)
	suite.Require().True(ok, "details.Value should be a string")
	suite.Assert().Equal("blob", value)
}

func ExampleError() {
	err := errors.NewSentinel(500, "error.test.custom", "Test Error").Clone()
	fmt.Println(err)

	var details *errors.Error
	if errors.As(err, &details) {
		fmt.Println(details.ID)
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

func ExampleError_Format_gosyntax_01() {
	output := CaptureStdout(func() {
		err := errors.WrapErrors(errors.ArgumentInvalid.With("key", "value"), errors.ArgumentMissing.With("key"))
		if err != nil {
			fmt.Printf("%#v\n", err)
		}
	})
	// remove the line numbers from the stack trace as they change when the code is changed
	simplifier := regexp.MustCompile(`\.go:[0-9]+`)
	// we also do not care about the last file which is machine dependent
	noasm := regexp.MustCompile(`, asm_.[^\.]+.s:[0-9]+`)
	fmt.Println(noasm.ReplaceAllString(simplifier.ReplaceAllString(output, ".go"), ""))
	// Output:
	// errors.Error{Code: 400, ID: "error.argument.invalid", Text: "Argument %s is invalid (value: %v)", What: "key", Value: "value", Cause: errors.Error{Code: 400, ID: "error.argument.missing", Text: "Argument %s is missing", What: "key", Stack: []errors.StackFrame{errors_test.go, errors_test.go, errors_test.go, run_example.go, example.go, testing.go, _testmain.go, proc.go}}, Stack: []errors.StackFrame{errors_test.go, errors_test.go, errors_test.go, run_example.go, example.go, testing.go, _testmain.go, proc.go}}
}

func ExampleError_Format_gosyntax_02() {
	output := CaptureStdout(func() {
		err := errors.WrapErrors(errors.ArgumentInvalid.With("key", "value"), fmt.Errorf("unknown error"))
		if err != nil {
			fmt.Printf("%#v\n", err)
		}
	})
	// remove the line numbers from the stack trace as they change when the code is changed
	simplifier := regexp.MustCompile(`\.go:[0-9]+`)
	// we also do not care about the last file which is machine dependent
	noasm := regexp.MustCompile(`, asm_.[^\.]+.s:[0-9]+`)
	fmt.Println(noasm.ReplaceAllString(simplifier.ReplaceAllString(output, ".go"), ""))
	// Output:
	// errors.Error{Code: 400, ID: "error.argument.invalid", Text: "Argument %s is invalid (value: %v)", What: "key", Value: "value", Cause: "unknown error", Stack: []errors.StackFrame{errors_test.go, errors_test.go, errors_test.go, run_example.go, example.go, testing.go, _testmain.go, proc.go}}
}

func ExampleError_With() {
	err := errors.ArgumentMissing.With("key")
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
	err := errors.ArgumentInvalid.With("key", "value")
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
	err := errors.ArgumentInvalid.With("key", []string{"value1", "value2"})
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
	// JSON failed to marshal data
	// Caused by:
	// 	unexpected end of JSON input
	// unexpected end of JSON input
}

func ExampleWrapErrors() {
	err := errors.WrapErrors(
		errors.ArgumentInvalid.With("key", "value"),
		errors.ArgumentMissing.With("key"),
		&url.Error{
			Op:  "GET",
			URL: "https://example.com",
			Err: fmt.Errorf("connection refused"),
		},
	)
	fmt.Println(err)
	// Output:
	// Argument key is invalid (value: value)
	// Caused by:
	// 	Argument key is missing
	// Caused by:
	// 	GET "https://example.com": connection refused
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
