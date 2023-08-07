package errors_test

import (
	"net/http"

	"github.com/gildas/go-errors"
)

func (suite *ErrorsSuite) TestCanCreateSentinel() {
	err := errors.NewSentinel(32123, "error.test.create", "this is the error")
	suite.Require().NotNil(err, "newly created sentinel cannot be nil")
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
