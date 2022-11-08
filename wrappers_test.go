package errors_test

import (
	"fmt"
	"net/http"

	"github.com/gildas/go-errors"
)

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
