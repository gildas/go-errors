package errors_test

import (
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
	err := errors.New("this is the error")
	_, ok := err.(error)
	suite.Require().True(ok, "Object is not an error")
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

		cause := errors.Cause(err)
		suite.Require().NotNil(cause, "Error's cause should not be nil")
		suite.Assert().Equal("Error", reflect.ValueOf(cause).Elem().Type().Name())

		unwrap := errors.Unwrap(err)
		suite.Require().NotNil(unwrap, "Error's unwrap should not be nil")
		suite.Assert().Equal("Error", reflect.ValueOf(cause).Elem().Type().Name())

		ok = errors.Is(err, sentinel)
		suite.Assert().True(ok, "Inner Error should be of the same type as the sentinel")

		var inner *errors.Error
		ok = errors.As(err, &inner)
		suite.Assert().True(ok, "Inner Error should be an errors.Error")
	}
}