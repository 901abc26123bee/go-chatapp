package errors

import (
	"encoding/json"
	e "errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"
)

type errSuite struct {
	suite.Suite
	errs      ErrContents     // expected ErrContents
	errMsgs   []string        // expected error messages
	errFormat []string        // input format string for NewErrorf
	errArgs   [][]interface{} // input args for NewErrorf
	jsons     []string        // exptected json string
}

func TestErrSuite(t *testing.T) {
	suite.Run(t, &errSuite{
		errs: ErrContents{
			{
				StatusCode: http.StatusRequestTimeout,
				Code:       RequestTimeout,
				Message:    "request timeout",
				Detail:     []interface{}{"failed to query user by id: context cancel"},
			},
			{
				StatusCode: http.StatusNotFound,
				Code:       NotFound,
				Message:    "resource not found",
				Detail:     []interface{}{"failed to find id: 1"},
			},
			{
				StatusCode: http.StatusUnauthorized,
				Code:       TokenEmpty,
				Message:    "token is empty",
				Detail:     []interface{}{"failed to parse token: empty token"},
			},
		},
		errMsgs: []string{
			"statuscode: 408, code: 1004, message: request timeout, detail: [failed to query user by id: context cancel]",
			"statuscode: 404, code: 1001, message: resource not found, detail: [failed to find id: 1]",
			"statuscode: 401, code: 2001, message: token is empty, detail: [failed to parse token: empty token]",
		},
		errFormat: []string{
			"failed to query user by id: %s",
			"failed to find id: %d",
			"failed to parse token: %s",
		},
		errArgs: [][]interface{}{
			{
				"context cancel",
			},
			{
				1,
			},
			{
				"empty token",
			},
		},
		jsons: []string{
			`{"status_code":408,"code":1004,"message":"request timeout","detail":["failed to query user by id: context cancel"]}`,
			`{"status_code":408,"code":1004,"message":"request timeout","detail":["failed to query user by id: context cancel"],"extra_errors":[{"status_code":404,"code":1001,"message":"resource not found","detail":["failed to find id: 1"]}]}`,
			`{"status_code":408,"code":1004,"message":"request timeout","detail":["failed to query user by id: context cancel"],"extra_errors":[{"status_code":404,"code":1001,"message":"resource not found","detail":["failed to find id: 1"]},{"status_code":401,"code":2001,"message":"token is empty","detail":["failed to parse token: empty token"]}]}`,
		},
	})
}

func (s *errSuite) TestNewError() {
	for i, expected := range s.errs {
		err := NewError(expected.Code, expected.Detail...)

		s.EqualValues(expected, err)
		s.EqualValues(s.errMsgs[i], err.Error())
	}
}

func (s *errSuite) TestNewErrorf() {
	for i, expected := range s.errs {
		s.EqualValues(expected, NewErrorf(expected.Code, s.errFormat[i], s.errArgs[i]...))
	}
}

func (s *errSuite) TestNew() {
	err := New("something wrong")
	s.EqualValues(Unknown, err.Code)
	s.EqualValues(http.StatusBadRequest, err.StatusCode)
	s.EqualError(err, "statuscode: 400, code: 1000, message: unknown error, detail: [something wrong]")
}

func (s *errSuite) TestErrorfWithoutError() {
	err := Errorf("something wrong: %d", 0)
	s.EqualValues(Unknown, err.Code)
	s.EqualValues(http.StatusBadRequest, err.StatusCode)
	s.EqualError(err, "statuscode: 400, code: 1000, message: unknown error, detail: [something wrong: 0]")
}

func (s *errSuite) TestErrorfWithNoCodeError() {
	err := Errorf("something wrong: %v", e.New("msg"))
	s.EqualValues(Unknown, err.Code)
	s.EqualValues(http.StatusBadRequest, err.StatusCode)
	s.EqualError(err, "statuscode: 400, code: 1000, message: unknown error, detail: [something wrong: msg]")
}

func (s *errSuite) TestErrorfWithCodeError() {
	err := Errorf("something wrong: %v", NewError(RequestTimeout, "msg"))
	s.EqualValues(RequestTimeout, err.Code)
	s.EqualValues(http.StatusRequestTimeout, err.StatusCode)
	s.EqualError(err, "statuscode: 408, code: 1004, message: request timeout, detail: [something wrong: statuscode: 408, code: 1004, message: request timeout, detail: [msg]]")
}

func (s *errSuite) TestMarshalJSON() {
	errors := ErrContents{}
	for i, e := range s.errs {
		errors = append(errors, e)
		bytes, err := json.Marshal(&errors)
		s.NoError(err)
		s.Equal(s.jsons[i], string(bytes))
	}
}
