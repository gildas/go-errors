package errors

import (
	"fmt"
	"net/http"
)

// NewSentinel creates a new sentinel.
//
// A sentinel is an Error that hasn't been decorated with a stack trace
//
// Typically, it can be used to create error that can be matched later
func NewSentinel(code int, id, message string) Error {
	return Error{Code: code, ID: id, Text: message}
}

// FromHTTPStatusCode creates a new error of the sentinel that matches the given HTTP status code.
//
// It also records the stack trace at the point it was called.
func FromHTTPStatusCode(code int) error {
	// TODO: We should not have HTTPUnauthorized and Unauthorized... They should be merged.
	switch code {
	case http.StatusBadGateway:
		return HTTPBadGateway.WithStack()
	case http.StatusBadRequest:
		return HTTPBadRequest.WithStack()
	case http.StatusForbidden:
		return HTTPForbidden.WithStack()
	case http.StatusInternalServerError:
		return HTTPInternalServerError.WithStack()
	case http.StatusMethodNotAllowed:
		return HTTPMethodNotAllowed.WithStack()
	case http.StatusNotFound:
		return HTTPNotFound.WithStack()
	case http.StatusNotImplemented:
		return HTTPNotImplemented.WithStack()
	case http.StatusServiceUnavailable:
		return HTTPServiceUnavailable.WithStack()
	case http.StatusUnauthorized:
		return HTTPUnauthorized.WithStack()
	case http.StatusConflict:
		return HTTPStatusConflict.WithStack()
	case http.StatusExpectationFailed:
		return HTTPStatusExpectationFailed.WithStack()
	case http.StatusFailedDependency:
		return HTTPStatusFailedDependency.WithStack()
	case http.StatusGatewayTimeout:
		return HTTPStatusGatewayTimeout.WithStack()
	case http.StatusGone:
		return HTTPStatusGone.WithStack()
	case http.StatusHTTPVersionNotSupported:
		return HTTPStatusHTTPVersionNotSupported.WithStack()
	case http.StatusInsufficientStorage:
		return HTTPStatusInsufficientStorage.WithStack()
	case http.StatusLengthRequired:
		return HTTPStatusLengthRequired.WithStack()
	case http.StatusLocked:
		return HTTPStatusLocked.WithStack()
	case http.StatusLoopDetected:
		return HTTPStatusLoopDetected.WithStack()
	case http.StatusMisdirectedRequest:
		return HTTPStatusMisdirectedRequest.WithStack()
	case http.StatusNetworkAuthenticationRequired:
		return HTTPStatusNetworkAuthenticationRequired.WithStack()
	case http.StatusNotAcceptable:
		return HTTPStatusNotAcceptable.WithStack()
	case http.StatusNotExtended:
		return HTTPStatusNotExtended.WithStack()
	case http.StatusPaymentRequired:
		return HTTPStatusPaymentRequired.WithStack()
	case http.StatusPreconditionFailed:
		return HTTPStatusPreconditionFailed.WithStack()
	case http.StatusPreconditionRequired:
		return HTTPStatusPreconditionRequired.WithStack()
	case http.StatusProxyAuthRequired:
		return HTTPStatusProxyAuthRequired.WithStack()
	case http.StatusRequestEntityTooLarge:
		return HTTPStatusRequestEntityTooLarge.WithStack()
	case http.StatusRequestHeaderFieldsTooLarge:
		return HTTPStatusRequestHeaderFieldsTooLarge.WithStack()
	case http.StatusRequestTimeout:
		return HTTPStatusRequestTimeout.WithStack()
	case http.StatusRequestURITooLong:
		return HTTPStatusRequestURITooLong.WithStack()
	case http.StatusRequestedRangeNotSatisfiable:
		return HTTPStatusRequestedRangeNotSatisfiable.WithStack()
	case http.StatusTeapot:
		return HTTPStatusTeapot.WithStack()
	case http.StatusTooEarly:
		return HTTPStatusTooEarly.WithStack()
	case http.StatusTooManyRequests:
		return HTTPStatusTooManyRequests.WithStack()
	case http.StatusUnavailableForLegalReasons:
		return HTTPStatusUnavailableForLegalReasons.WithStack()
	case http.StatusUnprocessableEntity:
		return HTTPStatusUnprocessableEntity.WithStack()
	case http.StatusUnsupportedMediaType:
		return HTTPStatusUnsupportedMediaType.WithStack()
	case http.StatusUpgradeRequired:
		return HTTPStatusUpgradeRequired.WithStack()
	case http.StatusUseProxy:
		return HTTPStatusUseProxy.WithStack()
	case http.StatusVariantAlsoNegotiates:
		return HTTPStatusVariantAlsoNegotiates.WithStack()
	default:
		return NewSentinel(code, fmt.Sprintf("error.http.%d", code), fmt.Sprintf("HTTP Status %d", code)).WithStack()
	}
}

/*********** Standard Errors ***********************************************************************************************************/

// ArgumentMissing is used when an argument is missing.
var ArgumentMissing = NewSentinel(http.StatusBadRequest, "error.argument.missing", "Argument %s is missing")

// ArgumentExpected is used when an argument is expected and another was set.
var ArgumentExpected = NewSentinel(http.StatusBadRequest, "error.argument.expected", "Argument %s is invalid (value: %v, expected: %v)")

// ArgumentInvalid is used when an argument has an unexpected value.
var ArgumentInvalid = NewSentinel(http.StatusBadRequest, "error.argument.invalid", "Argument %s is invalid (value: %v)")

// CreationFailed is used when something was not created properly.
var CreationFailed = NewSentinel(http.StatusInternalServerError, "error.creation.failed", "Failed Creating %s")

// Empty is used when something is empty whereas it should not.
var Empty = NewSentinel(http.StatusBadRequest, "error.empty", "%s is empty")

// EnvironmentMissing is used when an argument is missing.
var EnvironmentMissing = NewSentinel(http.StatusBadRequest, "error.environment.missing", "Environment variable %s is missing")

// EnvironmentInvalid is used when an argument has an unexpected value.
var EnvironmentInvalid = NewSentinel(http.StatusBadRequest, "error.environment.invalid", "Environment variable %s is invalid (value: %v)")

// DuplicateFound is used when something is found but it should not have been.
var DuplicateFound = NewSentinel(http.StatusFound, "error.found", "%s %s Found")

// Invalid is used when something is not valid.
var Invalid = NewSentinel(http.StatusBadRequest, "error.invalid", "Invalid %s (value: %v, expected: %v)")

// InvalidType is used when a type is not valid.
var InvalidType = NewSentinel(http.StatusBadRequest, "error.type.invalid", "Invalid Type %s, expected: %s")

// InvalidURL is used when a URL is not valid.
var InvalidURL = NewSentinel(http.StatusBadRequest, "error.url.invalid", "Invalid URL %s")

// JSONMarshalError is used when data failed to be marshaled into JSON.
var JSONMarshalError = NewSentinel(http.StatusBadRequest, "error.json.marshal", "JSON failed to marshal data")

// JSONUnmarshalError is used when JSON data is missing a property.
var JSONUnmarshalError = NewSentinel(http.StatusBadRequest, "error.json.unmarshal", "JSON failed to unmarshal data")

// JSONPropertyMissing is used when JSON data is missing a property.
var JSONPropertyMissing = NewSentinel(http.StatusBadRequest, "error.json.property.missing", "JSON data is missing property %s")

// Missing is used when something is missing.
var Missing = NewSentinel(http.StatusBadRequest, "error.missing", "%s is missing")

// NotConnected is used when some socket, client is not connected to its server.
var NotConnected = NewSentinel(http.StatusGone, "error.client.not_connected", "%s Not Connected")

// NotInitialized is used when something is not yet initialized.
var NotInitialized = NewSentinel(http.StatusBadRequest, "error.notinitialized", "%s is not yet initialized")

// NotFound is used when something is not found.
var NotFound = NewSentinel(http.StatusNotFound, "error.notfound", "%s %s Not Found")

// NotImplemented is used when some code/method/func is not written yet.
var NotImplemented = NewSentinel(http.StatusNotImplemented, "error.notimplemented", "Not Implemented")

// IndexOutOfBounds is used when an index is out of bounds.
var IndexOutOfBounds = NewSentinel(http.StatusBadRequest, "error.index.outofbounds", "Index %s is out of bounds (value: %v)")

// RuntimeError is used when the code failed executing something.
var RuntimeError = NewSentinel(http.StatusInternalServerError, "error.runtime", "Runtime Error")

// TooManyErrors is used when something is found too many times.
var TooManyErrors = NewSentinel(http.StatusInternalServerError, "error.toomany", "Too Many")

// Unauthorized is used when some credentials failed some authentication process.
var Unauthorized = NewSentinel(http.StatusUnauthorized, "error.unauthorized", "Invalid Credentials")

// Unsupported is used when something is unsupported by the code.
var Unsupported = NewSentinel(http.StatusMethodNotAllowed, "error.unsupported", "Unsupported %s: %s")

// UnknownError is used when the code does not know which error it is facing...
var UnknownError = NewSentinel(http.StatusInternalServerError, "error.unknown", "Unknown Error: %s")

/*********** HTTP Errors ***************************************************************************************************************/
// HTTPBadGateway is used when an http.Client request fails.
var HTTPBadGateway = NewSentinel(http.StatusBadGateway, "error.http.gateway", http.StatusText(http.StatusBadGateway))

// HTTPBadRequest is used when an http.Client request fails.
var HTTPBadRequest = NewSentinel(http.StatusBadRequest, "error.http.request", http.StatusText(http.StatusBadRequest)+". %s")

// HTTPForbidden is used when an http.Client request fails.
var HTTPForbidden = NewSentinel(http.StatusForbidden, "error.http.forbidden", http.StatusText(http.StatusForbidden))

// HTTPInternalServerError is used when an http.Client request fails.
var HTTPInternalServerError = NewSentinel(http.StatusInternalServerError, "error.http.server", http.StatusText(http.StatusInternalServerError))

// HTTPMethodNotAllowed is used when an http.Client request fails.
var HTTPMethodNotAllowed = NewSentinel(http.StatusMethodNotAllowed, "error.http.notallowed", http.StatusText(http.StatusMethodNotAllowed))

// HTTPNotFound is used when an http.Client request fails.
var HTTPNotFound = NewSentinel(http.StatusNotFound, "error.http.notfound", http.StatusText(http.StatusNotFound))

// HTTPNotImplemented is used when an http.Client request fails.
var HTTPNotImplemented = NewSentinel(http.StatusNotImplemented, "error.http.notimplemented", http.StatusText(http.StatusNotImplemented))

// HTTPServiceUnavailable is used when an http.Client request fails.
var HTTPServiceUnavailable = NewSentinel(http.StatusServiceUnavailable, "error.http.unavailable", http.StatusText(http.StatusServiceUnavailable))

// HTTPUnauthorized is used when an http.Client request fails.
var HTTPUnauthorized = NewSentinel(http.StatusUnauthorized, "error.http.unauthorized", http.StatusText(http.StatusUnauthorized))

// HTTPStatusConflict reports HTTP Error StatusConflict.
var HTTPStatusConflict = NewSentinel(http.StatusConflict, "error.http.conflict", http.StatusText(http.StatusConflict))

// HTTPStatusExpectationFailed reports HTTP Error StatusExpectationFailed.
var HTTPStatusExpectationFailed = NewSentinel(http.StatusExpectationFailed, "error.http.expectation.failed", http.StatusText(http.StatusExpectationFailed))

// HTTPStatusFailedDependency reports HTTP Error StatusFailedDependency.
var HTTPStatusFailedDependency = NewSentinel(http.StatusFailedDependency, "error.http.failed.dependency", http.StatusText(http.StatusFailedDependency))

// HTTPStatusGatewayTimeout reports HTTP Error StatusGatewayTimeout.
var HTTPStatusGatewayTimeout = NewSentinel(http.StatusGatewayTimeout, "error.http.gateway.timeout", http.StatusText(http.StatusGatewayTimeout))

// HTTPStatusGone reports HTTP Error StatusGone.
var HTTPStatusGone = NewSentinel(http.StatusGone, "error.http.gone", http.StatusText(http.StatusGone))

// HTTPStatusHTTPVersionNotSupported reports HTTP Error StatusHTTPVersionNotSupported.
var HTTPStatusHTTPVersionNotSupported = NewSentinel(http.StatusHTTPVersionNotSupported, "error.http.unsupported.version", http.StatusText(http.StatusHTTPVersionNotSupported))

// HTTPStatusInsufficientStorage reports HTTP Error StatusInsufficientStorage.
var HTTPStatusInsufficientStorage = NewSentinel(http.StatusInsufficientStorage, "error.http.storage.insufficient", http.StatusText(http.StatusInsufficientStorage))

// HTTPStatusLengthRequired reports HTTP Error StatusLengthRequired.
var HTTPStatusLengthRequired = NewSentinel(http.StatusLengthRequired, "error.http.length.required", http.StatusText(http.StatusLengthRequired))

// HTTPStatusLocked reports HTTP Error StatusLocked.
var HTTPStatusLocked = NewSentinel(http.StatusLocked, "error.http.locked", http.StatusText(http.StatusLocked))

// HTTPStatusLoopDetected reports HTTP Error StatusLoopDetected.
var HTTPStatusLoopDetected = NewSentinel(http.StatusLoopDetected, "error.http.loop.detected", http.StatusText(http.StatusLoopDetected))

// HTTPStatusMisdirectedRequest reports HTTP Error StatusMisdirectedRequest.
var HTTPStatusMisdirectedRequest = NewSentinel(http.StatusMisdirectedRequest, "error.http.misdirect.request", http.StatusText(http.StatusMisdirectedRequest))

// HTTPStatusNetworkAuthenticationRequired reports HTTP Error StatusNetworkAuthenticationRequired.
var HTTPStatusNetworkAuthenticationRequired = NewSentinel(http.StatusNetworkAuthenticationRequired, "error.http.network.authentication.required", http.StatusText(http.StatusNetworkAuthenticationRequired))

// HTTPStatusNotAcceptable reports HTTP Error StatusNotAcceptable.
var HTTPStatusNotAcceptable = NewSentinel(http.StatusNotAcceptable, "error.http.notacceptable", http.StatusText(http.StatusNotAcceptable))

// HTTPStatusNotExtended reports HTTP Error StatusNotExtended.
var HTTPStatusNotExtended = NewSentinel(http.StatusNotExtended, "error.http.notextended", http.StatusText(http.StatusNotExtended))

// HTTPStatusPaymentRequired reports HTTP Error StatusPaymentRequired.
var HTTPStatusPaymentRequired = NewSentinel(http.StatusPaymentRequired, "error.http.payment.required", http.StatusText(http.StatusPaymentRequired))

// HTTPStatusPreconditionFailed reports HTTP Error StatusPreconditionFailed.
var HTTPStatusPreconditionFailed = NewSentinel(http.StatusPreconditionFailed, "error.http.precondition.failed", http.StatusText(http.StatusPreconditionFailed))

// HTTPStatusPreconditionRequired reports HTTP Error StatusPreconditionRequired.
var HTTPStatusPreconditionRequired = NewSentinel(http.StatusPreconditionRequired, "error.precondition.required", http.StatusText(http.StatusPreconditionRequired))

// HTTPStatusProxyAuthRequired reports HTTP Error StatusProxyAuthRequired.
var HTTPStatusProxyAuthRequired = NewSentinel(http.StatusProxyAuthRequired, "error.http.proxy.authentication.required", http.StatusText(http.StatusProxyAuthRequired))

// HTTPStatusRequestEntityTooLarge reports HTTP Error StatusRequestEntityTooLarge.
var HTTPStatusRequestEntityTooLarge = NewSentinel(http.StatusRequestEntityTooLarge, "error.http.request.entity.toolarge", http.StatusText(http.StatusRequestEntityTooLarge))

// HTTPStatusRequestHeaderFieldsTooLarge reports HTTP Error StatusRequestHeaderFieldsTooLarge.
var HTTPStatusRequestHeaderFieldsTooLarge = NewSentinel(http.StatusRequestHeaderFieldsTooLarge, "error.http.request.fields.toolarge", http.StatusText(http.StatusRequestHeaderFieldsTooLarge))

// HTTPStatusRequestTimeout reports HTTP Error StatusRequestTimeout.
var HTTPStatusRequestTimeout = NewSentinel(http.StatusRequestTimeout, "error.http.request.timeout", http.StatusText(http.StatusRequestTimeout))

// HTTPStatusRequestURITooLong reports HTTP Error StatusRequestURITooLong.
var HTTPStatusRequestURITooLong = NewSentinel(http.StatusRequestURITooLong, "error.http.request.uri.toolong", http.StatusText(http.StatusRequestURITooLong))

// HTTPStatusRequestedRangeNotSatisfiable reports HTTP Error StatusRequestedRangeNotSatisfiable.
var HTTPStatusRequestedRangeNotSatisfiable = NewSentinel(http.StatusRequestedRangeNotSatisfiable, "error.http.request.range.notstatisfiable", http.StatusText(http.StatusRequestedRangeNotSatisfiable))

// HTTPStatusTeapot reports HTTP Error StatusTeapot.
var HTTPStatusTeapot = NewSentinel(http.StatusTeapot, "error.http.teapot", http.StatusText(http.StatusTeapot))

// HTTPStatusTooEarly reports HTTP Error StatusTooEarly.
var HTTPStatusTooEarly = NewSentinel(http.StatusTooEarly, "error.http.tooearly", http.StatusText(http.StatusTooEarly))

// HTTPStatusTooManyRequests reports HTTP Error StatusTooManyRequests.
var HTTPStatusTooManyRequests = NewSentinel(http.StatusTooManyRequests, "error.http.request.toomany", http.StatusText(http.StatusTooManyRequests))

// HTTPStatusUnavailableForLegalReasons reports HTTP Error StatusUnavailableForLegalReasons.
var HTTPStatusUnavailableForLegalReasons = NewSentinel(http.StatusUnavailableForLegalReasons, "error.http.unavailable", http.StatusText(http.StatusUnavailableForLegalReasons))

// HTTPStatusUnprocessableEntity reports HTTP Error StatusUnprocessableEntity.
var HTTPStatusUnprocessableEntity = NewSentinel(http.StatusUnprocessableEntity, "error.http.entity.unprocessable", http.StatusText(http.StatusUnprocessableEntity))

// HTTPStatusUnsupportedMediaType reports HTTP Error StatusUnsupportedMediaType.
var HTTPStatusUnsupportedMediaType = NewSentinel(http.StatusUnsupportedMediaType, "error.http.mediatype.unsupported", http.StatusText(http.StatusUnsupportedMediaType))

// HTTPStatusUpgradeRequired reports HTTP Error StatusUpgradeRequired.
var HTTPStatusUpgradeRequired = NewSentinel(http.StatusUpgradeRequired, "error.http.upgrade.required", http.StatusText(http.StatusUpgradeRequired))

// HTTPStatusUseProxy reports HTTP Error StatusUseProxy.
var HTTPStatusUseProxy = NewSentinel(http.StatusUseProxy, "error.http.proxy.required", http.StatusText(http.StatusUseProxy))

// HTTPStatusVariantAlsoNegotiates reports HTTP Error StatusVariantAlsoNegotiates.
var HTTPStatusVariantAlsoNegotiates = NewSentinel(http.StatusVariantAlsoNegotiates, "error.http.variant.alsonegotiate", http.StatusText(http.StatusVariantAlsoNegotiates))
