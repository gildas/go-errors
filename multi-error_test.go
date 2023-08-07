package errors_test

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/gildas/go-errors"
	"github.com/stretchr/testify/suite"
)

type MultiErrorSuite struct {
	suite.Suite
	Name string
}

func TestMultiErrorSuite(t *testing.T) {
	suite.Run(t, new(MultiErrorSuite))
}

func (suite *MultiErrorSuite) SetupSuite() {
	suite.Name = strings.TrimSuffix(reflect.TypeOf(suite).Elem().Name(), "Suite")
}

func (suite *MultiErrorSuite) TestCanCreate() {
	errs := &errors.MultiError{}
	suite.Require().NotNil(errs, "err should not be nil")
	suite.Assert().Nil(errs.AsError(), "err should contain nothing")
	suite.Assert().True(errs.IsEmpty(), "err should be contain nothing")
	suite.Assert().Equal("", errs.Error())
}

func (suite *MultiErrorSuite) TestCanConvertToGoErrorWith1Error() {
	var errs errors.MultiError
	suite.Assert().Nil(errs.AsError(), "err should contain nothing")
	errs.Append(errors.New("this is error 1"))
	suite.Require().Len(errs.Errors, 1, "err should contain one error")
	suite.Assert().NotNil(errs.AsError(), "err should contain something")
}

func (suite *MultiErrorSuite) TestCanConvertToGoErrorWith1ErrorWithoutStack() {
	var errs errors.MultiError
	suite.Assert().Nil(errs.AsError(), "err should contain nothing")
	err := errors.NotImplemented.Clone()
	errs.Append(err)
	suite.Require().Len(errs.Errors, 1, "err should contain one error")
	suite.Assert().NotNil(errs.AsError(), "err should contain something")
}

func (suite *MultiErrorSuite) TestCanConvertToGoErrorWith1PtrError() {
	var errs errors.MultiError
	suite.Assert().Nil(errs.AsError(), "err should contain nothing")
	err := errors.NotImplemented.Clone()
	err.Stack.Initialize()
	errs.Append(err)
	suite.Require().Len(errs.Errors, 1, "err should contain one error")
	suite.Assert().NotNil(errs.AsError(), "err should contain something")
}

func (suite *MultiErrorSuite) TestCanConvertToGoErrorWith1PtrErrorWithoutStack() {
	var errs errors.MultiError
	suite.Assert().Nil(errs.AsError(), "err should contain nothing")
	err := errors.NotImplemented.Clone()
	errs.Append(err)
	suite.Require().Len(errs.Errors, 1, "err should contain one error")
	suite.Assert().NotNil(errs.AsError(), "err should contain something")
}

func (suite *MultiErrorSuite) TestCanConvertToGoErrorWith1FmtError() {
	var errs errors.MultiError
	suite.Assert().Nil(errs.AsError(), "err should contain nothing")
	errs.Append(fmt.Errorf("this is error 1"))
	suite.Require().Len(errs.Errors, 1, "err should contain one error")
	suite.Assert().NotNil(errs.AsError(), "err should contain something")
}

func (suite *MultiErrorSuite) TestCanConvertToGoErrorWith2Errors() {
	var errs errors.MultiError
	suite.Assert().Nil(errs.AsError(), "err should contain nothing")
	errs.Append(errors.New("this is error 1"))
	suite.Require().Len(errs.Errors, 1, "err should contain one error")
	suite.Assert().NotNil(errs.AsError(), "err should contain something")
	errs.Append(errors.New("this is error 2"))
	suite.Require().Len(errs.Errors, 2, "err should contain one error")
	suite.Assert().NotNil(errs.AsError(), "err should contain something")
}

func (suite *MultiErrorSuite) TestCanGetErrorString() {
	var errs errors.MultiError
	suite.Assert().Nil(errs.AsError(), "err should contain nothing")
	errs.Append(errors.New("this is error 1"))
	suite.Require().Len(errs.Errors, 1, "err should contain one error")
	suite.Assert().Equal("this is error 1", errs.Error())
	errs.Append(errors.New("this is error 2"))
	suite.Require().Len(errs.Errors, 2, "err should contain one error")
	suite.Assert().Equal("2 errors:\nthis is error 1\nthis is error 2", errs.Error())
}

func (suite *MultiErrorSuite) TestCanCheckWithErrorIs() {
	var errs errors.MultiError

	errs.Append(errors.New("this is error 1"))
	errs.Append(errors.ArgumentMissing.With("name1"))
	errs.Append(errors.ArgumentInvalid.With("name2", "value2"))
	suite.Assert().ErrorIs(errs.AsError(), &errors.MultiError{})
	suite.Assert().ErrorIs(errs.AsError(), errors.ArgumentMissing)
	suite.Assert().ErrorIs(errs.AsError(), errors.ArgumentInvalid)
}

func (suite *MultiErrorSuite) TestCanConvertToErrorWithErrorAs() {
	var errs errors.MultiError
	errs.Append(errors.ArgumentMissing.With("name1"))
	errs.Append(errors.ArgumentInvalid.With("name2", "value2"))

	var multiDetails *errors.MultiError
	suite.Require().ErrorAs(errs.AsError(), &multiDetails, "err should be a MultiError")
	suite.Assert().Len(multiDetails.Errors, 2, "err should contain two errors")

	details := errors.ArgumentMissing.Clone()
	suite.Require().ErrorAs(errs.AsError(), &details, "should be able to convert to errors.AgumentMissing")
	suite.Assert().Equal(errors.ArgumentMissing.ID, details.ID)
	suite.Assert().Equal("name1", details.What)

	details = errors.ArgumentInvalid.Clone()
	suite.Require().ErrorAs(errs.AsError(), &details, "should be able to convert to errors.ArgumentInvalid")
	suite.Assert().Equal(errors.ArgumentInvalid.ID, details.ID)
	suite.Assert().Equal("name2", details.What)
	value, ok := details.Value.(string)
	suite.Require().True(ok, "value should be a string")
	suite.Assert().Equal("value2", value)
}

func (suite *MultiErrorSuite) TestShouldNotMatchWithOtherErrors() {
	var errs errors.MultiError

	errs.Append(errors.ArgumentMissing.With("name1"))
	errs.Append(errors.ArgumentInvalid.With("name2", "value2"))
	err := errs.AsError()

	suite.Assert().NotErrorIs(err, &os.PathError{})
	suite.Assert().NotErrorIs(err, errors.NotImplemented)
}

func (suite *MultiErrorSuite) TestShouldNotConvertToOtherErrors() {
	var errs errors.MultiError
	errs.Append(errors.ArgumentMissing.With("name1"))
	errs.Append(errors.ArgumentInvalid.With("name2", "value2"))

	details := errors.NotImplemented.Clone()
	suite.Assert().False(errors.As(errs.AsError(), &details), "should not be able to convert to errors.ArgumentInvalid")

	var otherDetails *os.PathError
	suite.Assert().False(errors.As(errs.AsError(), &otherDetails), "should not be able to convert to os.PathError")
}

func ExampleMultiError() {
	var errs errors.MultiError

	errs.Append(
		errors.New("this is the first error"),
		errors.New("this is the second error"),
	)
	fmt.Println(errs.Error())
	// Output:
	// 2 errors:
	// this is the first error
	// this is the second error
}
