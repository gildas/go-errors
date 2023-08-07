package errors_test

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"regexp"

	"github.com/gildas/go-errors"
)

func (suite *ErrorsSuite) TestCanFormatStackTrace() {
	err := errors.NotImplemented.WithStack()
	actual, ok := err.(errors.Error)
	suite.Require().True(ok)
	suite.Require().NotEmpty(actual.Stack, "The stack should not be empty")
	suite.Assert().Contains(fmt.Sprintf("%v", actual.Stack), "[stack_test.go:13 value.go")
	suite.Assert().Contains(fmt.Sprintf("%s", actual.Stack), "[stack_test.go value.go")
}

func (suite *ErrorsSuite) TestCanFormatStackFrame() {
	err := errors.NotImplemented.WithStack()
	actual, ok := err.(errors.Error)
	suite.Require().True(ok)
	suite.Require().NotEmpty(actual.Stack, "The stack should not be empty")
	frame := actual.Stack[0]
	suite.Assert().Equal("(*ErrorsSuite).TestCanFormatStackFrame", fmt.Sprintf("%n", frame))
}

func (suite *ErrorsSuite) TestCanMarshalStackTrace() {
	testerr := errors.ArgumentInvalid.With("key", "value")
	_, err := json.Marshal(testerr.(errors.Error).Stack)
	suite.Require().Nil(err)
}

func (suite *ErrorsSuite) TestCanMarshalStackFrameAsText() {
	err := errors.NotImplemented.WithStack()
	actual, ok := err.(errors.Error)
	suite.Require().True(ok)
	suite.Require().NotEmpty(actual.Stack, "The stack should not be empty")
	frame := actual.Stack[0]
	payload, xmlerr := xml.Marshal(frame)
	suite.Require().Nil(xmlerr)
	pattern := regexp.MustCompile(`<StackFrame>.*TestCanMarshalStackFrameAsText .*/stack_test.go:[0-9]+</StackFrame>`)
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
