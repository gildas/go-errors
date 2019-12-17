package errors

import "net/http"

// NewSentinel creates a new sentinel
// a sentinel is an Error that hasn't been decorated with a stack trace
// Typically, it can be used to create error that can be matched later
func NewSentinel(code int, id, message string) *Error {
	return &Error{ Code: code, ID: id, Text: message}
}

// WithWhat creates a new error from a given sentinal telling "What" is wrong
func (e *Error) WithWhat(what string) error {
	final := *e
	final.What = what
	return WithStack(&final)
}

// WithWhatAndValue creates a new error from a given sentinal telling "What" is wrong and the wrong value
func (e *Error) WithWhatAndValue(what string, value interface{}) error {
	final := *e
	final.What  = what
	final.Value = value
	return WithStack(&final)
}

/*********** Standard Errors ***********************************************************************************************************/

// ArgumentMissingError is used when an argument is missing
var ArgumentMissingError = NewSentinel(http.StatusBadRequest, "error.argument.missing", "Argument %s is missing")

// ArgumentInvalidError is used when an argument has an unexpected value
var ArgumentInvalidError = NewSentinel(http.StatusBadRequest, "error.argument.invalid", "Argument %s is invalid (value: %v)")

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
var HTTPBadRequestError = NewSentinel(http.StatusBadRequest, http.StatusText(http.StatusBadRequest) + ". %s", "error.http.request")

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
/*
		http.StatusAlreadyReported: HTTPGenericError,
		http.StatusBadGateway: HTTPBadGatewayError,
		http.StatusBadRequest: HTTPBadRequestError,
		http.StatusConflict: HTTPGenericError,
		http.StatusContinue: HTTPGenericError,
		http.StatusCreated: HTTPGenericError,
		http.StatusExpectationFailed: HTTPGenericError,
		http.StatusFailedDependency: HTTPGenericError,
		http.StatusForbidden: HTTPForbiddenError,
		http.StatusFound: HTTPGenericError,
		http.StatusGatewayTimeout: HTTPGenericError,
		http.StatusGone: HTTPGenericError,
		http.StatusHTTPVersionNotSupported: HTTPGenericError,
		http.StatusIMUsed: HTTPGenericError,
		http.StatusInsufficientStorage: HTTPGenericError,
		http.StatusInternalServerError: HTTPInternalServerError,
		http.StatusLengthRequired: HTTPGenericError,
		http.StatusLocked: HTTPGenericError,
		http.StatusLoopDetected: HTTPGenericError,
		http.StatusMethodNotAllowed: HTTPMethodNotAllowedError,
		http.StatusMisdirectedRequest: HTTPGenericError,
		http.StatusMovedPermanently: HTTPGenericError,
		http.StatusMultiStatus: HTTPGenericError,
		http.StatusMultipleChoices: HTTPGenericError,
		http.StatusNetworkAuthenticationRequired: HTTPGenericError,
		http.StatusNoContent: HTTPGenericError,
		http.StatusNonAuthoritativeInfo: HTTPGenericError,
		http.StatusNotAcceptable: HTTPGenericError,
		http.StatusNotExtended: HTTPGenericError,
		http.StatusNotFound: HTTPNotFoundError,
		http.StatusNotImplemented: HTTPNotImplementedError,
		http.StatusNotModified: HTTPGenericError,
		http.StatusOK: HTTPGenericError,
		http.StatusPartialContent: HTTPGenericError,
		http.StatusPaymentRequired: HTTPGenericError,
		http.StatusPermanentRedirect: HTTPGenericError,
		http.StatusPreconditionFailed: HTTPGenericError,
		http.StatusPreconditionRequired: HTTPGenericError,
		http.StatusProcessing: HTTPGenericError,
		http.StatusProxyAuthRequired: HTTPGenericError,
		http.StatusRequestEntityTooLarge: HTTPGenericError,
		http.StatusRequestHeaderFieldsTooLarge: HTTPGenericError,
		http.StatusRequestTimeout: HTTPGenericError,
		http.StatusRequestURITooLong: HTTPGenericError,
		http.StatusRequestedRangeNotSatisfiable: HTTPGenericError,
		http.StatusResetContent: HTTPGenericError,
		http.StatusSeeOther: HTTPGenericError,
		http.StatusServiceUnavailable: HTTPServiceUnavailableError,
		http.StatusSwitchingProtocols: HTTPGenericError,
		http.StatusTeapot: HTTPGenericError,
		http.StatusTemporaryRedirect: HTTPGenericError,
		http.StatusTooEarly: HTTPGenericError,
		http.StatusTooManyRequests: HTTPGenericError,
		http.StatusUnauthorized: HTTPUnauthorizedError,
		http.StatusUnavailableForLegalReasons: HTTPGenericError,
		http.StatusUnprocessableEntity: HTTPGenericError,
		http.StatusUnsupportedMediaType: HTTPGenericError,
		http.StatusUpgradeRequired: HTTPGenericError,
		http.StatusUseProxy: HTTPGenericError,
		http.StatusVariantAlsoNegotiates: HTTPGenericError,
*/