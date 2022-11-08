package errors_test

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
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
	err := errors.NewSentinel(32123, "error.test.create", "this is the error")
	suite.Require().NotNil(err, "newly created sentinel cannot be nil")
}

func (suite *ErrorsSuite) TestCanFormatStackTrace() {
	err := errors.NotImplemented.WithStack()
	actual, ok := err.(errors.Error)
	suite.Require().True(ok)
	suite.Require().NotEmpty(actual.Stack, "The stack should not be empty")
	suite.Assert().Contains(fmt.Sprintf("%v", actual.Stack), "[errors_test.go:42 value.go")
	suite.Assert().Contains(fmt.Sprintf("%s", actual.Stack), "[errors_test.go value.go")
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
	wrapped := errors.WrapErrors(nil)
	suite.Require().Nil(wrapped, "Chained error of nil should be nil")
}

func (suite *ErrorsSuite) TestCanWrapOneError() {
	error1 := errors.ArgumentMissing.With("key1")
	wrapped := errors.WrapErrors(error1)
	suite.Assert().Equal(error1, wrapped, "Chained errors of one error should be the error")
}

func (suite *ErrorsSuite) TestCanWrapTwoErrors() {
	error1 := errors.ArgumentMissing.With("key1")
	error2 := errors.NotFound.With("key", "key2")

	wrapped := errors.WrapErrors(error1, error2)
	suite.Assert().ErrorIs(wrapped, errors.Error{}, "Chained errors should contain an errors.Error")
	suite.Assert().ErrorIs(wrapped, errors.ArgumentMissing, "Chained errors should match an ArgumentMissingError")
	suite.Assert().ErrorIs(wrapped, errors.NotFound, "Chained errors should match a NotFoundError")

	first, ok := wrapped.(errors.Error)
	suite.Require().True(ok, "The first error should be an errors.Error")
	suite.Assert().Equal(errors.ArgumentMissing.ID, first.ID, "The first error should be error1")
	suite.Require().NotNil(first.Cause, "The first error should have a cause")

	unwrapped := first.Unwrap()
	suite.Require().NotNil(unwrapped, "The first error should have a cause")

	second, ok := unwrapped.(errors.Error)
	suite.Require().True(ok, "The second error should be an errors.Error")
	suite.Assert().Equal(errors.NotFound.ID, second.ID, "The second error should be error2")
	suite.Assert().Nil(second.Cause, "The second error should not have a cause")

	unwrapped = second.Unwrap()
	suite.Require().Nil(unwrapped, "The second error should have a cause")
}

func (suite *ErrorsSuite) TestCanWrapTwoErrorsAndNil() {
	error1 := errors.ArgumentMissing.With("key1")
	error2 := errors.NotFound.With("key", "key2")

	wrapped := errors.WrapErrors(error1, error2, nil)
	suite.Assert().ErrorIs(wrapped, errors.Error{}, "Chained errors should contain an errors.Error")
	suite.Assert().ErrorIs(wrapped, errors.ArgumentMissing, "Chained errors should match an ArgumentMissingError")
	suite.Assert().ErrorIs(wrapped, errors.NotFound, "Chained errors should match a NotFoundError")

	first, ok := wrapped.(errors.Error)
	suite.Require().True(ok, "The first error should be an errors.Error")
	suite.Assert().Equal(errors.ArgumentMissing.ID, first.ID, "The first error should be error1")
	suite.Require().NotNil(first.Cause, "The first error should have a cause")

	unwrapped := first.Unwrap()
	suite.Require().NotNil(unwrapped, "The first error should have a cause")

	second, ok := unwrapped.(errors.Error)
	suite.Require().True(ok, "The second error should be an errors.Error")
	suite.Assert().Equal(errors.NotFound.ID, second.ID, "The second error should be error2")
	suite.Require().Nil(second.Cause, "The second error should not have a cause")
}

func (suite *ErrorsSuite) TestCanWrapThreeErrors() {
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
	suite.Require().Nil(unwrapped, "The third error should have a cause")
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
	suite.Assert().ErrorIs(wrapped, errors.ArgumentMissing, "Chained errors should match an ArgumentMissingError")
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
		fmt.Errorf("some obscure error"),
		nil,
	)
	payload, err := json.Marshal(testerr)
	suite.Require().Nil(err)
	suite.Assert().JSONEq(expected, string(payload))
}

func (suite *ErrorsSuite) TestCanMarshalStackTrace() {
	testerr := errors.ArgumentInvalid.With("key", "value")
	_, err := json.Marshal(testerr.(errors.Error).Stack)
	suite.Require().Nil(err)
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

func (suite *ErrorsSuite) TestCanFormatStackFrame() {
	err := errors.NotImplemented.WithStack()
	actual, ok := err.(errors.Error)
	suite.Require().True(ok)
	suite.Require().NotEmpty(actual.Stack, "The stack should not be empty")
	frame := actual.Stack[0]
	suite.Assert().Equal("(*ErrorsSuite).TestCanFormatStackFrame", fmt.Sprintf("%n", frame))
}

func (suite *ErrorsSuite) TestCanMarshalStackFrameAsText() {
	err := errors.NotImplemented.WithStack()
	actual, ok := err.(errors.Error)
	suite.Require().True(ok)
	suite.Require().NotEmpty(actual.Stack, "The stack should not be empty")
	frame := actual.Stack[0]
	payload, xmlerr := xml.Marshal(frame)
	suite.Require().Nil(xmlerr)
	pattern := regexp.MustCompile(`<StackFrame>.*TestCanMarshalStackFrameAsText .*/errors_test.go:[0-9]+</StackFrame>`)
	suite.Assert().Regexp(pattern, string(payload))
}

func (suite *ErrorsSuite) TestCanUseInvalidStackFrame() {
	frame := errors.StackFrame(0)
	suite.Assert().Equal("unknown", frame.FuncName())
	suite.Assert().Equal(0, frame.Line())
	suite.Assert().Equal("unknown", frame.Filepath())
	payload, xmlerr := xml.Marshal(frame)
	suite.Require().Nil(xmlerr)
	pattern := regexp.MustCompile(`<StackFrame>unknown</StackFrame>`)
	suite.Assert().Regexp(pattern, string(payload))
}

func (suite *ErrorsSuite) TestWrappers() {
	var err error

	err = errors.New("Hello World")
	suite.Assert().NotNil(err)
	suite.Assert().Equal("Hello World", fmt.Sprintf("%s", err))
	if actual, ok := err.(errors.Error); ok {
		suite.Assert().Equal(http.StatusInternalServerError, actual.Code)
		suite.Assert().Equal("error.runtime", actual.ID)
	} else {
		suite.Assert().Fail("Error should be an errors.Error")
	}

	err = errors.Errorf("Hello World")
	suite.Assert().NotNil(err)
	suite.Assert().Equal("Hello World", fmt.Sprintf("%s", err))

	err = errors.WithStack(fmt.Errorf("Hello World"))
	suite.Assert().NotNil(err)
	suite.Assert().Equal("error.runtime\nCaused by:\n\tHello World", fmt.Sprintf("%s", err))
	suite.Assert().Nil(errors.WithStack(nil))

	err = errors.WithStack(errors.NotImplemented)
	suite.Assert().NotNil(err)
	suite.Assert().Equal("Not Implemented", fmt.Sprintf("%s", err))

	err = errors.WithMessage(errors.NotFound.With("greetings", "hi"), "Hello World")
	suite.Assert().NotNil(err)
	suite.Assert().Equal("Hello World\nCaused by:\n\tgreetings hi Not Found", fmt.Sprintf("%s", err))
	suite.Assert().Nil(errors.WithMessage(nil, "Hello World"))

	err = errors.WithMessagef(errors.NotFound.With("greetings", "hi"), "Hello %s", "World")
	suite.Assert().NotNil(err)
	suite.Assert().Equal("Hello World\nCaused by:\n\tgreetings hi Not Found", fmt.Sprintf("%s", err))
	suite.Assert().Nil(errors.WithMessagef(nil, "Hello %s", "World"))

	err = errors.Wrap(errors.NotFound.With("greetings", "hi"), "Hello World")
	suite.Assert().NotNil(err)
	suite.Assert().Equal("Hello World\nCaused by:\n\tgreetings hi Not Found", fmt.Sprintf("%s", err))
	suite.Assert().Nil(errors.Wrap(nil, "Hello World"))

	err = errors.Wrapf(errors.NotFound.With("greetings", "hi"), "Hello %s", "World")
	suite.Assert().NotNil(err)
	suite.Assert().Equal("Hello World\nCaused by:\n\tgreetings hi Not Found", fmt.Sprintf("%s", err))
	suite.Assert().Nil(errors.Wrapf(nil, "Hello %s", "World"))

	unwrapped := errors.Unwrap(err)
	suite.Assert().NotNil(unwrapped)

	suite.Assert().True(errors.Is(err, errors.NotFound), "err should be of the same type as NotFoundError")

	var inner *errors.Error
	suite.Assert().True(errors.As(err, &inner), "Inner Error should be an errors.Error")
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
	// errors.Error{Code: 400, ID: "error.argument.invalid", Text: "Argument %s is invalid (value: %v)", What: "key", Value: "value", Cause: errors.Error{Code: 400, ID: "error.argument.missing", Text: "Argument %s is missing", What: "key", Stack: []errors.StackFrame{errors_test.go, errors_test.go, errors_test.go, run_example.go, example.go, testing.go, _testmain.go, proc.go}}, Stack: []errors.StackFrame{wrappers.go, errors_test.go, errors_test.go, errors_test.go, run_example.go, example.go, testing.go, _testmain.go, proc.go}}
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
	// errors.Error{Code: 400, ID: "error.argument.invalid", Text: "Argument %s is invalid (value: %v)", What: "key", Value: "value", Cause: "unknown error", Stack: []errors.StackFrame{wrappers.go, errors_test.go, errors_test.go, errors_test.go, run_example.go, example.go, testing.go, _testmain.go, proc.go}}
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
