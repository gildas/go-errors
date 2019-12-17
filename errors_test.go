package errors_test

import (
	"fmt"
	"encoding/json"
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
	err := errors.NotFoundError.WithWhat("key")
	suite.Require().NotNil(err, "err should not be nil")
	suite.Assert().True(errors.Is(err, errors.NotFoundError), "err should be a NotFoundError")
}

func (suite *ErrorsSuite) TestCanTellContainsAnError() {
	err := errors.NotFoundError.WithWhat("key")
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

func (suite *ErrorsSuite) TestFailsWithNonErrorTarget() {
	suite.Assert().False(errors.NotFoundError.Is(errors.New("Hello")), "Error should not be a pkg/errors.fundamental")
}

func (suite *ErrorsSuite) TestSentinels() {
	sentinels := map[string]errors.Error{
		"ArgumentInvalidError": errors.ArgumentInvalidError,
	}

	for name, sentinel := range sentinels {
		err := sentinel.WithWhatAndValue("test", "value")
		suite.Assert().Equal("", sentinel.What, "Sentinel's what should not have been changed")
		suite.Assert().Nil(sentinel.Value, "Sentinel's value should not have been changed")

		_, ok := err.(error)
		suite.Assert().Truef(ok, "Instance of %s is not an error", name)
		suite.Assert().Equal("withStack", reflect.ValueOf(err).Elem().Type().Name())

		unwrap := errors.Unwrap(err)
		suite.Require().NotNil(unwrap, "Error's unwrap should not be nil")
		suite.Assert().Equal("Error", reflect.ValueOf(unwrap).Elem().Type().Name())

		ok = errors.Is(err, sentinel)
		suite.Assert().True(ok, "Inner Error should be of the same type as the sentinel")

		var inner *errors.Error
		ok = errors.As(err, &inner)
		suite.Assert().True(ok, "Inner Error should be an errors.Error")
	}
}

func (suite *ErrorsSuite) TestWrappers() {
	var err error

	err = errors.New("Hello World")
	suite.Assert().NotNil(err)

	err = errors.Errorf("Hello World")
	suite.Assert().NotNil(err)

	err = errors.Wrap(errors.NotFoundError, "Hello World")
	suite.Assert().NotNil(err)

	err = errors.Wrapf(errors.NotFoundError, "Hello World")
	suite.Assert().NotNil(err)

	err = errors.WithMessage(errors.NotFoundError, "Hello World")
	suite.Assert().NotNil(err)

	err = errors.WithMessagef(errors.NotFoundError, "Hello World")
	suite.Assert().NotNil(err)
}

func ExampleError_WithWhat() {
	err := errors.ArgumentMissingError.WithWhat("key")
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

func ExampleError_WithWhatAndValue() {
	err := errors.ArgumentInvalidError.WithWhatAndValue("key", "value")
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

func ExampleError_WithWhatAndValue_array() {
	err := errors.ArgumentInvalidError.WithWhatAndValue("key", []string{"value1", "value2"})
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
	var payload struct {Value string `json:"value"`}

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