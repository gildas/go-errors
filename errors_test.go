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
	sentinels := map[string]*errors.Error{
		"ArgumentInvalidError": &errors.ArgumentInvalidError,
	}

	for name, sentinel := range sentinels {
		err := sentinel.WithWhatAndValue("test", "value")
		_, ok := err.(error)
		suite.Assert().Truef(ok, "Instance of %s is not an error", name)
		var inner errors.Error
		suite.Assert().True(errors.As(err, &inner), "Inner Error should be an errors.Error")

		value := errors.Unwrap(err)
		unwrapped, ok := value.(errors.Error)
		suite.Assert().True(ok, "Unwrapped error should be an errors.Error")
		if ok {
			suite.Assert().Equal("test", unwrapped.What)
		}

	}
}