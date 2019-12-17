package errors

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

// ArgumentMissingError is used when an argument is missing
var ArgumentMissingError = NewSentinel(0, "error.argument.missing", "Argument %s is missing")

// ArgumentInvalidError is used when an argument has an unexpected value
var ArgumentInvalidError = NewSentinel(0, "error.argument.invalid", "Argument %s is invalid (value: %v)")

// FoundError is used when something is found but it should not have been
var FoundError = NewSentinel(0, "error.found", "%s %s Found")

// JSONMarshalError is used when data failed to be marshaled into JSON
var JSONMarshalError = NewSentinel(0, "error.json.marshal", "JSON failed to marshal data")

// JSONUnmarshalError is used when JSON data is missing a property
var JSONUnmarshalError = NewSentinel(0, "error.json.unmarshal", "JSON failed to unmarshal data")

// JSONPropertyMissingError is used when JSON data is missing a property
var JSONPropertyMissingError = NewSentinel(0, "error.json.property.missing", "JSON data is missing property %s")

// NotFoundError is used when something is not found
var NotFoundError = NewSentinel(0, "error.notfound", "%s %s Not Found")

// NotImplementedError is used when some code/method/func is not written yet
var NotImplementedError = NewSentinel(0, "error.notimplemented", "Not Implemented")

// TooManyError is used when something is found too many times
var TooManyError = NewSentinel(0, "error.toomany", "Too Many")

// UnsupportedError is used when something is unsupported by the code
var UnsupportedError = NewSentinel(0, "error.unsupported", "Unsupported %s: %s")

// UnknownError is used when the code does not know which error it is facing...
var UnknownError = NewSentinel(0, "error.unknown", "Unknown Error: %s")
