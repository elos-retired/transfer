package transfer

import (
	"errors"
	"fmt"
	"net/http"
)

var ErrConnectionClosed = errors.New("SocketConnection is closed")

type Error struct {
	Status           int    `json:"status"`
	Code             int    `json:"code"`
	Message          string `json:"message"`
	DeveloperMessage string `json:"developer_message"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("Transfer Error %d: %s", e.Status, e.Message)
}

func NewError(status int, code int, msg string, dmsg string) *Error {
	return &Error{
		Status:           status,
		Code:             code,
		Message:          msg,
		DeveloperMessage: dmsg,
	}
}

var ErrGeneric = &Error{
	Status:           400,
	Code:             400,
	Message:          "This is the generic error",
	DeveloperMessage: "Complain to nick about returning the generic error",
}

var ErrNotFound = &Error{
	Status:           http.StatusNotFound,
	Code:             http.StatusNotFound,
	Message:          "Not Found",
	DeveloperMessage: "Perhaps you have an incorrect id?",
}

var ErrBadMethod = &Error{
	Status:  http.StatusMethodNotAllowed,
	Code:    http.StatusMethodNotAllowed,
	Message: http.StatusText(http.StatusMethodNotAllowed),
}

var ErrAuth = &Error{
	Status:           401,
	Code:             401,
	Message:          "You are not authenticated",
	DeveloperMessage: "You need to supply auth, every, single, time.",
}

var ErrAccess = &Error{
	Status:           401,
	Code:             401,
	Message:          "You have been denied access",
	DeveloperMessage: "You can only access data on behalf of your use",
}

var ErrWebSocketFailed = &Error{
	Status:           http.StatusBadRequest,
	Code:             http.StatusBadRequest,
	Message:          http.StatusText(http.StatusBadRequest),
	DeveloperMessage: "We were unable to process your websocket request, perhaps it was not spec-valid?",
}

func NewServerError(devMsg string) *Error {
	return &Error{
		Status:           http.StatusInternalServerError,
		Code:             http.StatusInternalServerError,
		Message:          http.StatusText(http.StatusInternalServerError),
		DeveloperMessage: devMsg,
	}
}

func NewServerErrorWithError(err error) *Error {
	return &Error{
		Status:           http.StatusInternalServerError,
		Code:             http.StatusInternalServerError,
		Message:          http.StatusText(http.StatusInternalServerError),
		DeveloperMessage: err.Error(),
	}
}

func NewMethodNotAllowedError(got Action) *Error {
	return &Error{
		Status:           http.StatusMethodNotAllowed,
		Code:             http.StatusMethodNotAllowed,
		Message:          http.StatusText(http.StatusMethodNotAllowed),
		DeveloperMessage: fmt.Sprintf("The endpoint you requested does not handle the %s action", got),
	}
}
