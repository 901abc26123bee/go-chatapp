package errors

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// ErrCode defines error codes of server
type ErrCode int

// defines general errors
const (
	Unknown ErrCode = iota + 1000
	NotFound
	InternalServerError
	InvalidArgument
	RequestTimeout
	ParamIncorrect
	// ...
)

// defines authorization errors
const (
	TokenExpired ErrCode = iota + 2000
	TokenEmpty
	TokenInvalid
	// ...
)

// ErrContent defines customized errors
type ErrContent struct {
	StatusCode int           `json:"status_code"`
	Code       ErrCode       `json:"code"`
	Message    string        `json:"message"`
	Detail     []interface{} `json:"detail,omitempty"`
}

type ErrContents []*ErrContent

var (
	// errors for 1000(general) prefix
	errUnknown             = *newErrContent(http.StatusBadRequest, Unknown, "unknown error")
	errNotFound            = *newErrContent(http.StatusNotFound, NotFound, "resource not found")
	errInternalServerError = *newErrContent(http.StatusInternalServerError, InternalServerError, "internal server error")
	errInvalidArgument     = *newErrContent(http.StatusInternalServerError, InvalidArgument, "invalid argument")
	errRequestTimeout      = *newErrContent(http.StatusRequestTimeout, RequestTimeout, "request timeout")
	errParamIncorrect      = *newErrContent(http.StatusBadRequest, ParamIncorrect, "parameter incorrect")

	// errors for 2000(authorization) prefix
	errTokenExpired = *newErrContent(http.StatusUnauthorized, TokenExpired, "token is expired")
	errTokenEmpty   = *newErrContent(http.StatusUnauthorized, TokenEmpty, "token is empty")
	errTokenInvalid = *newErrContent(http.StatusUnauthorized, TokenInvalid, "token is invalid")
)

var errorPool = map[ErrCode]ErrContent{
	// errors for 1000(general) prefix
	Unknown:             errUnknown,
	NotFound:            errNotFound,
	InternalServerError: errInternalServerError,
	InvalidArgument:     errInvalidArgument,
	RequestTimeout:      errRequestTimeout,
	ParamIncorrect:      errParamIncorrect,

	// errors for 2000(authorization) prefix
	TokenEmpty:   errTokenEmpty,
	TokenExpired: errTokenExpired,
	TokenInvalid: errTokenInvalid,
}

func newErrContent(statusCode int, code ErrCode, msg string) *ErrContent {
	return &ErrContent{StatusCode: statusCode, Code: code, Message: msg}
}

// NewError return an ErrContent, default with un-know error
func NewError(code ErrCode, details ...any) *ErrContent {
	// find pre-defined error by code
	err, ok := errorPool[code]
	if !ok {
		err = errUnknown
	}

	// attach the details
	err.Detail = details

	return &err
}

// NewErrorf return an ErrContent with formatted string
func NewErrorf(code ErrCode, format string, args ...any) *ErrContent {
	return NewError(code, fmt.Sprintf(format, args...))
}

// New returns an unknown error with message.
func New(msg string) *ErrContent {
	return NewError(Unknown, msg)
}

// Errorf returns a error with formatted string.
// If one of args is with http code, the returned error will be inherited.
// If all args are without http code, the returned error will be unknown error.
func Errorf(format string, args ...any) *ErrContent {
	// iterate with all args to check whether it contains an ErrContent.
	code := Unknown
	for _, arg := range args {
		if ec, ok := arg.(*ErrContent); ok {
			// only use http code for first ErrContent here.
			code = ec.Code
			break
		}
	}

	// construct the error by detected http code in the args.
	return NewErrorf(code, format, args...)
}

func (e *ErrContent) Error() string {
	return fmt.Sprintf("statuscode: %d, code: %d, message: %s, detail: %v",
		e.StatusCode, e.Code, e.Message, e.Detail)
}

// MarshalJSON marshal ErrContents to set the first error as main info and others as extra_errors
func (e ErrContents) MarshalJSON() ([]byte, error) {
	if len(e) == 0 {
		return []byte(""), nil
	}

	return json.Marshal(&struct {
		*ErrContent
		ExtraErrors []*ErrContent `json:"extra_errors,omitempty"`
	}{
		ErrContent:  e[0],
		ExtraErrors: e[1:],
	})
}
