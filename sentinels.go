package errors

import "net/http"

// NewSentinel creates a new sentinel
// a sentinel is an Error that hasn't been decorated with a stack trace
// Typically, it can be used to create error that can be matched later
func NewSentinel(code int, id, message string) *Error {
	return &Error{Code: code, ID: id, Text: message}
}

/*********** Standard Errors ***********************************************************************************************************/

// ArgumentMissingError is used when an argument is missing
var ArgumentMissingError = NewSentinel(http.StatusBadRequest, "error.argument.missing", "Argument %s is missing")

// ArgumentInvalidError is used when an argument has an unexpected value
var ArgumentInvalidError = NewSentinel(http.StatusBadRequest, "error.argument.invalid", "Argument %s is invalid (value: %v)")

// EnvironmentMissingError is used when an argument is missing
var EnvironmentMissingError = NewSentinel(http.StatusBadRequest, "error.environment.missing", "Environment variable %s is missing")

// EnvironmentInvalidError is used when an argument has an unexpected value
var EnvironmentInvalidError = NewSentinel(http.StatusBadRequest, "error.environment.invalid", "Environment variable %s is invalid (value: %v)")

// FoundError is used when something is found but it should not have been
var FoundError = NewSentinel(http.StatusFound, "error.found", "%s %s Found")

// JSONMarshalError is used when data failed to be marshaled into JSON
var JSONMarshalError = NewSentinel(http.StatusBadRequest, "error.json.marshal", "JSON failed to marshal data")

// JSONUnmarshalError is used when JSON data is missing a property
var JSONUnmarshalError = NewSentinel(http.StatusBadRequest, "error.json.unmarshal", "JSON failed to unmarshal data")

// JSONPropertyMissingError is used when JSON data is missing a property
var JSONPropertyMissingError = NewSentinel(http.StatusBadRequest, "error.json.property.missing", "JSON data is missing property %s")

// NotFoundError is used when something is not found
var NotFoundError = NewSentinel(http.StatusNotFound, "error.notfound", "%s %s Not Found")

// NotImplementedError is used when some code/method/func is not written yet
var NotImplementedError = NewSentinel(http.StatusNotImplemented, "error.notimplemented", "Not Implemented")

// TooManyError is used when something is found too many times
var TooManyError = NewSentinel(http.StatusInternalServerError, "error.toomany", "Too Many")

// UnsupportedError is used when something is unsupported by the code
var UnsupportedError = NewSentinel(http.StatusMethodNotAllowed, "error.unsupported", "Unsupported %s: %s")

// UnknownError is used when the code does not know which error it is facing...
var UnknownError = NewSentinel(http.StatusInternalServerError, "error.unknown", "Unknown Error: %s")

/*********** HTTP Errors ***************************************************************************************************************/
// HTTPBadGatewayError is used when an http.Client request fails
var HTTPBadGatewayError = NewSentinel(http.StatusBadGateway, http.StatusText(http.StatusBadGateway), "error.http.gateway")

// HTTPBadRequestError is used when an http.Client request fails
var HTTPBadRequestError = NewSentinel(http.StatusBadRequest, http.StatusText(http.StatusBadRequest)+". %s", "error.http.request")

// HTTPForbiddenError is used when an http.Client request fails
var HTTPForbiddenError = NewSentinel(http.StatusForbidden, http.StatusText(http.StatusForbidden), "error.http.forbidden")

// HTTPInternalServerError is used when an http.Client request fails
var HTTPInternalServerError = NewSentinel(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "error.http.server")

// HTTPMethodNotAllowedError is used when an http.Client request fails
var HTTPMethodNotAllowedError = NewSentinel(http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed), "error.http.notallowed")

// HTTPNotFoundError is used when an http.Client request fails
var HTTPNotFoundError = NewSentinel(http.StatusNotFound, http.StatusText(http.StatusNotFound), "error.http.notfound")

// HTTPNotImplementedError is used when an http.Client request fails
var HTTPNotImplementedError = NewSentinel(http.StatusNotImplemented, http.StatusText(http.StatusNotImplemented), "error.http.notimplemented")

// HTTPServiceUnavailableError is used when an http.Client request fails
var HTTPServiceUnavailableError = NewSentinel(http.StatusServiceUnavailable, http.StatusText(http.StatusServiceUnavailable), "error.http.unavailable")

// HTTPUnauthorizedError is used when an http.Client request fails
var HTTPUnauthorizedError = NewSentinel(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized), "error.http.unauthorized")

// HTTPStatusConflictError reports HTTP Error StatusConflict
var HTTPStatusConflictError = NewSentinel(http.StatusConflict, http.StatusText(http.StatusConflict), "error.http.conflict")

// HTTPStatusExpectationFailedError reports HTTP Error StatusExpectationFailed
var HTTPStatusExpectationFailedError = NewSentinel(http.StatusExpectationFailed, http.StatusText(http.StatusExpectationFailed), "error.http.expectation.failed")

// HTTPStatusFailedDependencyError reports HTTP Error StatusFailedDependency
var HTTPStatusFailedDependencyError = NewSentinel(http.StatusFailedDependency, http.StatusText(http.StatusFailedDependency), "error.http.failed.dependency")

// HTTPStatusGatewayTimeoutError reports HTTP Error StatusGatewayTimeout
var HTTPStatusGatewayTimeoutError = NewSentinel(http.StatusGatewayTimeout, http.StatusText(http.StatusGatewayTimeout), "error.http.gateway.timeout")

// HTTPStatusGoneError reports HTTP Error StatusGone
var HTTPStatusGoneError = NewSentinel(http.StatusGone, http.StatusText(http.StatusGone), "error.http.gone")

// HTTPStatusHTTPVersionNotSupportedError reports HTTP Error StatusHTTPVersionNotSupported
var HTTPStatusHTTPVersionNotSupportedError = NewSentinel(http.StatusHTTPVersionNotSupported, http.StatusText(http.StatusHTTPVersionNotSupported), "error.http.unsupported.version")

// HTTPStatusInsufficientStorageError reports HTTP Error StatusInsufficientStorage
var HTTPStatusInsufficientStorageError = NewSentinel(http.StatusInsufficientStorage, http.StatusText(http.StatusInsufficientStorage), "error.http.storage.insufficient")

// HTTPStatusLengthRequiredError reports HTTP Error StatusLengthRequired
var HTTPStatusLengthRequiredError = NewSentinel(http.StatusLengthRequired, http.StatusText(http.StatusLengthRequired), "error.http.length.required")

// HTTPStatusLockedError reports HTTP Error StatusLocked
var HTTPStatusLockedError = NewSentinel(http.StatusLocked, http.StatusText(http.StatusLocked), "error.http.locked")

// HTTPStatusLoopDetectedError reports HTTP Error StatusLoopDetected
var HTTPStatusLoopDetectedError = NewSentinel(http.StatusLoopDetected, http.StatusText(http.StatusLoopDetected), "error.http.loop.detected")

// HTTPStatusMisdirectedRequestError reports HTTP Error StatusMisdirectedRequest
var HTTPStatusMisdirectedRequestError = NewSentinel(http.StatusMisdirectedRequest, http.StatusText(http.StatusMisdirectedRequest), "error.http.misdirect.request")

// HTTPStatusNetworkAuthenticationRequiredError reports HTTP Error StatusNetworkAuthenticationRequired
var HTTPStatusNetworkAuthenticationRequiredError = NewSentinel(http.StatusNetworkAuthenticationRequired, http.StatusText(http.StatusNetworkAuthenticationRequired), "error.http.network.authentication.required")

// HTTPStatusNotAcceptableError reports HTTP Error StatusNotAcceptable
var HTTPStatusNotAcceptableError = NewSentinel(http.StatusNotAcceptable, http.StatusText(http.StatusNotAcceptable), "error.http.notacceptable")

// HTTPStatusNotExtendedError reports HTTP Error StatusNotExtended
var HTTPStatusNotExtendedError = NewSentinel(http.StatusNotExtended, http.StatusText(http.StatusNotExtended), "error.http.notextended")

// HTTPStatusPaymentRequiredError reports HTTP Error StatusPaymentRequired
var HTTPStatusPaymentRequiredError = NewSentinel(http.StatusPaymentRequired, http.StatusText(http.StatusPaymentRequired), "error.http.payment.required")

// HTTPStatusPreconditionFailedError reports HTTP Error StatusPreconditionFailed
var HTTPStatusPreconditionFailedError = NewSentinel(http.StatusPreconditionFailed, http.StatusText(http.StatusPreconditionFailed), "error.http.precondition.failed")

// HTTPStatusPreconditionRequiredError reports HTTP Error StatusPreconditionRequired
var HTTPStatusPreconditionRequiredError = NewSentinel(http.StatusPreconditionRequired, http.StatusText(http.StatusPreconditionRequired), "error.precondition.required")

// HTTPStatusProxyAuthRequiredError reports HTTP Error StatusProxyAuthRequired
var HTTPStatusProxyAuthRequiredError = NewSentinel(http.StatusProxyAuthRequired, http.StatusText(http.StatusProxyAuthRequired), "error.http.proxy.authentication.required")

// HTTPStatusRequestEntityTooLargeError reports HTTP Error StatusRequestEntityTooLarge
var HTTPStatusRequestEntityTooLargeError = NewSentinel(http.StatusRequestEntityTooLarge, http.StatusText(http.StatusRequestEntityTooLarge), "error.http.request.entity.toolarge")

// HTTPStatusRequestHeaderFieldsTooLargeError reports HTTP Error StatusRequestHeaderFieldsTooLarge
var HTTPStatusRequestHeaderFieldsTooLargeError = NewSentinel(http.StatusRequestHeaderFieldsTooLarge, http.StatusText(http.StatusRequestHeaderFieldsTooLarge), "error.http.request.fields.toolarge")

// HTTPStatusRequestTimeoutError reports HTTP Error StatusRequestTimeout
var HTTPStatusRequestTimeoutError = NewSentinel(http.StatusRequestTimeout, http.StatusText(http.StatusRequestTimeout), "error.http.request.timeout")

// HTTPStatusRequestURITooLongError reports HTTP Error StatusRequestURITooLong
var HTTPStatusRequestURITooLongError = NewSentinel(http.StatusRequestURITooLong, http.StatusText(http.StatusRequestURITooLong), "error.http.request.uri.toolong")

// HTTPStatusRequestedRangeNotSatisfiableError reports HTTP Error StatusRequestedRangeNotSatisfiable
var HTTPStatusRequestedRangeNotSatisfiableError = NewSentinel(http.StatusRequestedRangeNotSatisfiable, http.StatusText(http.StatusRequestedRangeNotSatisfiable), "error.http.request.range.notstatisfiable")

// HTTPStatusTeapotError reports HTTP Error StatusTeapot
var HTTPStatusTeapotError = NewSentinel(http.StatusTeapot, http.StatusText(http.StatusTeapot), "error.http.teapot")

// HTTPStatusTooEarlyError reports HTTP Error StatusTooEarly
var HTTPStatusTooEarlyError = NewSentinel(http.StatusTooEarly, http.StatusText(http.StatusTooEarly), "error.http.tooearly")

// HTTPStatusTooManyRequestsError reports HTTP Error StatusTooManyRequests
var HTTPStatusTooManyRequestsError = NewSentinel(http.StatusTooManyRequests, http.StatusText(http.StatusTooManyRequests), "error.http.request.toomany")

// HTTPStatusUnavailableForLegalReasonsError reports HTTP Error StatusUnavailableForLegalReasons
var HTTPStatusUnavailableForLegalReasonsError = NewSentinel(http.StatusUnavailableForLegalReasons, http.StatusText(http.StatusUnavailableForLegalReasons), "error.http.unavailable")

// HTTPStatusUnprocessableEntityError reports HTTP Error StatusUnprocessableEntity
var HTTPStatusUnprocessableEntityError = NewSentinel(http.StatusUnprocessableEntity, http.StatusText(http.StatusUnprocessableEntity), "error.http.entity.unprocessable")

// HTTPStatusUnsupportedMediaTypeError reports HTTP Error StatusUnsupportedMediaType
var HTTPStatusUnsupportedMediaTypeError = NewSentinel(http.StatusUnsupportedMediaType, http.StatusText(http.StatusUnsupportedMediaType), "error.http.mediatype.unsupported")

// HTTPStatusUpgradeRequiredError reports HTTP Error StatusUpgradeRequired
var HTTPStatusUpgradeRequiredError = NewSentinel(http.StatusUpgradeRequired, http.StatusText(http.StatusUpgradeRequired), "error.http.upgrade.required")

// HTTPStatusUseProxyError reports HTTP Error StatusUseProxy
var HTTPStatusUseProxyError = NewSentinel(http.StatusUseProxy, http.StatusText(http.StatusUseProxy), "error.http.proxy.required")

// HTTPStatusVariantAlsoNegotiatesError reports HTTP Error StatusVariantAlsoNegotiates
var HTTPStatusVariantAlsoNegotiatesError = NewSentinel(http.StatusVariantAlsoNegotiates, http.StatusText(http.StatusVariantAlsoNegotiates), "error.http.variant.alsonegotiate")
