package errors

import "net/http"

import "fmt"

// NewSentinel creates a new sentinel.
// A sentinel is an Error that hasn't been decorated with a stack trace
//
// Typically, it can be used to create error that can be matched later
func NewSentinel(code int, id, message string) *Error {
	return &Error{Code: code, ID: id, Text: message}
}

// FromHTTPStatusCode creates a new error of the sentinel that matches the given HTTP status code.
// It also records the stack trace at the point it was called.
func FromHTTPStatusCode(code int) error {
	switch code {
	case http.StatusBadGateway:
		return WithStack(HTTPBadGatewayError)
	case http.StatusBadRequest:
		return WithStack(HTTPBadRequestError)
	case http.StatusForbidden:
		return WithStack(HTTPForbiddenError)
	case http.StatusInternalServerError:
		return WithStack(HTTPInternalServerError)
	case http.StatusMethodNotAllowed:
		return WithStack(HTTPMethodNotAllowedError)
	case http.StatusNotFound:
		return WithStack(HTTPNotFoundError)
	case http.StatusNotImplemented:
		return WithStack(HTTPNotImplementedError)
	case http.StatusServiceUnavailable:
		return WithStack(HTTPServiceUnavailableError)
	case http.StatusUnauthorized:
		return WithStack(HTTPUnauthorizedError)
	case http.StatusConflict:
		return WithStack(HTTPStatusConflictError)
	case http.StatusExpectationFailed:
		return WithStack(HTTPStatusExpectationFailedError)
	case http.StatusFailedDependency:
		return WithStack(HTTPStatusFailedDependencyError)
	case http.StatusGatewayTimeout:
		return WithStack(HTTPStatusGatewayTimeoutError)
	case http.StatusGone:
		return WithStack(HTTPStatusGoneError)
	case http.StatusHTTPVersionNotSupported:
		return WithStack(HTTPStatusHTTPVersionNotSupportedError)
	case http.StatusInsufficientStorage:
		return WithStack(HTTPStatusInsufficientStorageError)
	case http.StatusLengthRequired:
		return WithStack(HTTPStatusLengthRequiredError)
	case http.StatusLocked:
		return WithStack(HTTPStatusLockedError)
	case http.StatusLoopDetected:
		return WithStack(HTTPStatusLoopDetectedError)
	case http.StatusMisdirectedRequest:
		return WithStack(HTTPStatusMisdirectedRequestError)
	case http.StatusNetworkAuthenticationRequired:
		return WithStack(HTTPStatusNetworkAuthenticationRequiredError)
	case http.StatusNotAcceptable:
		return WithStack(HTTPStatusNotAcceptableError)
	case http.StatusNotExtended:
		return WithStack(HTTPStatusNotExtendedError)
	case http.StatusPaymentRequired:
		return WithStack(HTTPStatusPaymentRequiredError)
	case http.StatusPreconditionFailed:
		return WithStack(HTTPStatusPreconditionFailedError)
	case http.StatusPreconditionRequired:
		return WithStack(HTTPStatusPreconditionRequiredError)
	case http.StatusProxyAuthRequired:
		return WithStack(HTTPStatusProxyAuthRequiredError)
	case http.StatusRequestEntityTooLarge:
		return WithStack(HTTPStatusRequestEntityTooLargeError)
	case http.StatusRequestHeaderFieldsTooLarge:
		return WithStack(HTTPStatusRequestHeaderFieldsTooLargeError)
	case http.StatusRequestTimeout:
		return WithStack(HTTPStatusRequestTimeoutError)
	case http.StatusRequestURITooLong:
		return WithStack(HTTPStatusRequestURITooLongError)
	case http.StatusRequestedRangeNotSatisfiable:
		return WithStack(HTTPStatusRequestedRangeNotSatisfiableError)
	case http.StatusTeapot:
		return WithStack(HTTPStatusTeapotError)
	case http.StatusTooEarly:
		return WithStack(HTTPStatusTooEarlyError)
	case http.StatusTooManyRequests:
		return WithStack(HTTPStatusTooManyRequestsError)
	case http.StatusUnavailableForLegalReasons:
		return WithStack(HTTPStatusUnavailableForLegalReasonsError)
	case http.StatusUnprocessableEntity:
		return WithStack(HTTPStatusUnprocessableEntityError)
	case http.StatusUnsupportedMediaType:
		return WithStack(HTTPStatusUnsupportedMediaTypeError)
	case http.StatusUpgradeRequired:
		return WithStack(HTTPStatusUpgradeRequiredError)
	case http.StatusUseProxy:
		return WithStack(HTTPStatusUseProxyError)
	case http.StatusVariantAlsoNegotiates:
		return WithStack(HTTPStatusVariantAlsoNegotiatesError)
	default:
		return WithStack(NewSentinel(code, fmt.Sprintf("error.http.%d", code), fmt.Sprintf("HTTP Status %d", code)))
	}
}

/*********** Standard Errors ***********************************************************************************************************/

// ArgumentMissingError is used when an argument is missing.
var ArgumentMissingError = NewSentinel(http.StatusBadRequest, "error.argument.missing", "Argument %s is missing")

// ArgumentInvalidError is used when an argument has an unexpected value.
var ArgumentInvalidError = NewSentinel(http.StatusBadRequest, "error.argument.invalid", "Argument %s is invalid (value: %v)")

// EnvironmentMissingError is used when an argument is missing.
var EnvironmentMissingError = NewSentinel(http.StatusBadRequest, "error.environment.missing", "Environment variable %s is missing")

// EnvironmentInvalidError is used when an argument has an unexpected value.
var EnvironmentInvalidError = NewSentinel(http.StatusBadRequest, "error.environment.invalid", "Environment variable %s is invalid (value: %v)")

// FoundError is used when something is found but it should not have been.
var FoundError = NewSentinel(http.StatusFound, "error.found", "%s %s Found")

// JSONMarshalError is used when data failed to be marshaled into JSON.
var JSONMarshalError = NewSentinel(http.StatusBadRequest, "error.json.marshal", "JSON failed to marshal data")

// JSONUnmarshalError is used when JSON data is missing a property.
var JSONUnmarshalError = NewSentinel(http.StatusBadRequest, "error.json.unmarshal", "JSON failed to unmarshal data")

// JSONPropertyMissingError is used when JSON data is missing a property.
var JSONPropertyMissingError = NewSentinel(http.StatusBadRequest, "error.json.property.missing", "JSON data is missing property %s")

// NotConnectedError is used when some socket, client is not connected to its server.
var NotConnectedError = NewSentinel(http.StatusGone, "error.client.not_connected", "%s Not Connected")

// NotFoundError is used when something is not found.
var NotFoundError = NewSentinel(http.StatusNotFound, "error.notfound", "%s %s Not Found")

// NotImplementedError is used when some code/method/func is not written yet.
var NotImplementedError = NewSentinel(http.StatusNotImplemented, "error.notimplemented", "Not Implemented")

// RuntimeError is used when the code failed executing something.
var RuntimeError = NewSentinel(http.StatusInternalServerError, "error.runtime", "Runtime Error")

// TooManyError is used when something is found too many times.
var TooManyError = NewSentinel(http.StatusInternalServerError, "error.toomany", "Too Many")

// UnsupportedError is used when something is unsupported by the code.
var UnsupportedError = NewSentinel(http.StatusMethodNotAllowed, "error.unsupported", "Unsupported %s: %s")

// UnknownError is used when the code does not know which error it is facing...
var UnknownError = NewSentinel(http.StatusInternalServerError, "error.unknown", "Unknown Error: %s")

/*********** HTTP Errors ***************************************************************************************************************/
// HTTPBadGatewayError is used when an http.Client request fails.
var HTTPBadGatewayError = NewSentinel(http.StatusBadGateway, "error.http.gateway", http.StatusText(http.StatusBadGateway))

// HTTPBadRequestError is used when an http.Client request fails.
var HTTPBadRequestError = NewSentinel(http.StatusBadRequest, "error.http.request", http.StatusText(http.StatusBadRequest)+". %s")

// HTTPForbiddenError is used when an http.Client request fails.
var HTTPForbiddenError = NewSentinel(http.StatusForbidden, "error.http.forbidden", http.StatusText(http.StatusForbidden))

// HTTPInternalServerError is used when an http.Client request fails.
var HTTPInternalServerError = NewSentinel(http.StatusInternalServerError, "error.http.server", http.StatusText(http.StatusInternalServerError))

// HTTPMethodNotAllowedError is used when an http.Client request fails.
var HTTPMethodNotAllowedError = NewSentinel(http.StatusMethodNotAllowed, "error.http.notallowed", http.StatusText(http.StatusMethodNotAllowed))

// HTTPNotFoundError is used when an http.Client request fails.
var HTTPNotFoundError = NewSentinel(http.StatusNotFound, "error.http.notfound", http.StatusText(http.StatusNotFound))

// HTTPNotImplementedError is used when an http.Client request fails.
var HTTPNotImplementedError = NewSentinel(http.StatusNotImplemented, "error.http.notimplemented", http.StatusText(http.StatusNotImplemented))

// HTTPServiceUnavailableError is used when an http.Client request fails.
var HTTPServiceUnavailableError = NewSentinel(http.StatusServiceUnavailable, "error.http.unavailable", http.StatusText(http.StatusServiceUnavailable))

// HTTPUnauthorizedError is used when an http.Client request fails.
var HTTPUnauthorizedError = NewSentinel(http.StatusUnauthorized, "error.http.unauthorized", http.StatusText(http.StatusUnauthorized))

// HTTPStatusConflictError reports HTTP Error StatusConflict.
var HTTPStatusConflictError = NewSentinel(http.StatusConflict, "error.http.conflict", http.StatusText(http.StatusConflict))

// HTTPStatusExpectationFailedError reports HTTP Error StatusExpectationFailed.
var HTTPStatusExpectationFailedError = NewSentinel(http.StatusExpectationFailed, "error.http.expectation.failed", http.StatusText(http.StatusExpectationFailed))

// HTTPStatusFailedDependencyError reports HTTP Error StatusFailedDependency.
var HTTPStatusFailedDependencyError = NewSentinel(http.StatusFailedDependency, "error.http.failed.dependency", http.StatusText(http.StatusFailedDependency))

// HTTPStatusGatewayTimeoutError reports HTTP Error StatusGatewayTimeout.
var HTTPStatusGatewayTimeoutError = NewSentinel(http.StatusGatewayTimeout, "error.http.gateway.timeout", http.StatusText(http.StatusGatewayTimeout))

// HTTPStatusGoneError reports HTTP Error StatusGone.
var HTTPStatusGoneError = NewSentinel(http.StatusGone, "error.http.gone", http.StatusText(http.StatusGone))

// HTTPStatusHTTPVersionNotSupportedError reports HTTP Error StatusHTTPVersionNotSupported.
var HTTPStatusHTTPVersionNotSupportedError = NewSentinel(http.StatusHTTPVersionNotSupported, "error.http.unsupported.version", http.StatusText(http.StatusHTTPVersionNotSupported))

// HTTPStatusInsufficientStorageError reports HTTP Error StatusInsufficientStorage.
var HTTPStatusInsufficientStorageError = NewSentinel(http.StatusInsufficientStorage, "error.http.storage.insufficient", http.StatusText(http.StatusInsufficientStorage))

// HTTPStatusLengthRequiredError reports HTTP Error StatusLengthRequired.
var HTTPStatusLengthRequiredError = NewSentinel(http.StatusLengthRequired, "error.http.length.required", http.StatusText(http.StatusLengthRequired))

// HTTPStatusLockedError reports HTTP Error StatusLocked.
var HTTPStatusLockedError = NewSentinel(http.StatusLocked, "error.http.locked", http.StatusText(http.StatusLocked))

// HTTPStatusLoopDetectedError reports HTTP Error StatusLoopDetected.
var HTTPStatusLoopDetectedError = NewSentinel(http.StatusLoopDetected, "error.http.loop.detected", http.StatusText(http.StatusLoopDetected))

// HTTPStatusMisdirectedRequestError reports HTTP Error StatusMisdirectedRequest.
var HTTPStatusMisdirectedRequestError = NewSentinel(http.StatusMisdirectedRequest, "error.http.misdirect.request", http.StatusText(http.StatusMisdirectedRequest))

// HTTPStatusNetworkAuthenticationRequiredError reports HTTP Error StatusNetworkAuthenticationRequired.
var HTTPStatusNetworkAuthenticationRequiredError = NewSentinel(http.StatusNetworkAuthenticationRequired, "error.http.network.authentication.required", http.StatusText(http.StatusNetworkAuthenticationRequired))

// HTTPStatusNotAcceptableError reports HTTP Error StatusNotAcceptable.
var HTTPStatusNotAcceptableError = NewSentinel(http.StatusNotAcceptable, "error.http.notacceptable", http.StatusText(http.StatusNotAcceptable))

// HTTPStatusNotExtendedError reports HTTP Error StatusNotExtended.
var HTTPStatusNotExtendedError = NewSentinel(http.StatusNotExtended, "error.http.notextended", http.StatusText(http.StatusNotExtended))

// HTTPStatusPaymentRequiredError reports HTTP Error StatusPaymentRequired.
var HTTPStatusPaymentRequiredError = NewSentinel(http.StatusPaymentRequired, "error.http.payment.required", http.StatusText(http.StatusPaymentRequired))

// HTTPStatusPreconditionFailedError reports HTTP Error StatusPreconditionFailed.
var HTTPStatusPreconditionFailedError = NewSentinel(http.StatusPreconditionFailed, "error.http.precondition.failed", http.StatusText(http.StatusPreconditionFailed))

// HTTPStatusPreconditionRequiredError reports HTTP Error StatusPreconditionRequired.
var HTTPStatusPreconditionRequiredError = NewSentinel(http.StatusPreconditionRequired, "error.precondition.required", http.StatusText(http.StatusPreconditionRequired))

// HTTPStatusProxyAuthRequiredError reports HTTP Error StatusProxyAuthRequired.
var HTTPStatusProxyAuthRequiredError = NewSentinel(http.StatusProxyAuthRequired, "error.http.proxy.authentication.required", http.StatusText(http.StatusProxyAuthRequired))

// HTTPStatusRequestEntityTooLargeError reports HTTP Error StatusRequestEntityTooLarge.
var HTTPStatusRequestEntityTooLargeError = NewSentinel(http.StatusRequestEntityTooLarge, "error.http.request.entity.toolarge", http.StatusText(http.StatusRequestEntityTooLarge))

// HTTPStatusRequestHeaderFieldsTooLargeError reports HTTP Error StatusRequestHeaderFieldsTooLarge.
var HTTPStatusRequestHeaderFieldsTooLargeError = NewSentinel(http.StatusRequestHeaderFieldsTooLarge, "error.http.request.fields.toolarge", http.StatusText(http.StatusRequestHeaderFieldsTooLarge))

// HTTPStatusRequestTimeoutError reports HTTP Error StatusRequestTimeout.
var HTTPStatusRequestTimeoutError = NewSentinel(http.StatusRequestTimeout, "error.http.request.timeout", http.StatusText(http.StatusRequestTimeout))

// HTTPStatusRequestURITooLongError reports HTTP Error StatusRequestURITooLong.
var HTTPStatusRequestURITooLongError = NewSentinel(http.StatusRequestURITooLong, "error.http.request.uri.toolong", http.StatusText(http.StatusRequestURITooLong))

// HTTPStatusRequestedRangeNotSatisfiableError reports HTTP Error StatusRequestedRangeNotSatisfiable.
var HTTPStatusRequestedRangeNotSatisfiableError = NewSentinel(http.StatusRequestedRangeNotSatisfiable, "error.http.request.range.notstatisfiable", http.StatusText(http.StatusRequestedRangeNotSatisfiable))

// HTTPStatusTeapotError reports HTTP Error StatusTeapot.
var HTTPStatusTeapotError = NewSentinel(http.StatusTeapot, "error.http.teapot", http.StatusText(http.StatusTeapot))

// HTTPStatusTooEarlyError reports HTTP Error StatusTooEarly.
var HTTPStatusTooEarlyError = NewSentinel(http.StatusTooEarly, "error.http.tooearly", http.StatusText(http.StatusTooEarly))

// HTTPStatusTooManyRequestsError reports HTTP Error StatusTooManyRequests.
var HTTPStatusTooManyRequestsError = NewSentinel(http.StatusTooManyRequests, "error.http.request.toomany", http.StatusText(http.StatusTooManyRequests))

// HTTPStatusUnavailableForLegalReasonsError reports HTTP Error StatusUnavailableForLegalReasons.
var HTTPStatusUnavailableForLegalReasonsError = NewSentinel(http.StatusUnavailableForLegalReasons, "error.http.unavailable", http.StatusText(http.StatusUnavailableForLegalReasons))

// HTTPStatusUnprocessableEntityError reports HTTP Error StatusUnprocessableEntity.
var HTTPStatusUnprocessableEntityError = NewSentinel(http.StatusUnprocessableEntity, "error.http.entity.unprocessable", http.StatusText(http.StatusUnprocessableEntity))

// HTTPStatusUnsupportedMediaTypeError reports HTTP Error StatusUnsupportedMediaType.
var HTTPStatusUnsupportedMediaTypeError = NewSentinel(http.StatusUnsupportedMediaType, "error.http.mediatype.unsupported", http.StatusText(http.StatusUnsupportedMediaType))

// HTTPStatusUpgradeRequiredError reports HTTP Error StatusUpgradeRequired.
var HTTPStatusUpgradeRequiredError = NewSentinel(http.StatusUpgradeRequired, "error.http.upgrade.required", http.StatusText(http.StatusUpgradeRequired))

// HTTPStatusUseProxyError reports HTTP Error StatusUseProxy.
var HTTPStatusUseProxyError = NewSentinel(http.StatusUseProxy, "error.http.proxy.required", http.StatusText(http.StatusUseProxy))

// HTTPStatusVariantAlsoNegotiatesError reports HTTP Error StatusVariantAlsoNegotiates.
var HTTPStatusVariantAlsoNegotiatesError = NewSentinel(http.StatusVariantAlsoNegotiates, "error.http.variant.alsonegotiate", http.StatusText(http.StatusVariantAlsoNegotiates))
