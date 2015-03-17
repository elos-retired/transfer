package transfer

import (
	"net/http"

	"github.com/elos/data"
)

type HTTPConnection struct {
	w http.ResponseWriter
	r *http.Request
	data.Access
}

func NewHTTPConnection(w http.ResponseWriter, r *http.Request, a data.Access) *HTTPConnection {
	return &HTTPConnection{
		w:      w,
		r:      r,
		Access: a,
	}
}

func (c *HTTPConnection) WriteJSON(v interface{}) error {
	c.w.Header().Set("Content-Type", "application/json; charset=utf-8")

	bytes, err := ToJSON(v)
	if err != nil {
		return err
	}

	_, err = c.w.Write(bytes)

	return err
}

func (c *HTTPConnection) NotFound() {
	c.WriteErrorResponse(ErrNotFound)
}

func (c *HTTPConnection) ServerError(err error) {
	c.WriteErrorResponse(NewServerErrorWithError(err))
}

func (c *HTTPConnection) InvalidMethod() {
	c.WriteErrorResponse(NewMethodNotAllowedError("one you gave me"))
}

func (c *HTTPConnection) Unauthorized() {
	c.WriteErrorResponse(ErrAuth)
}

func (c *HTTPConnection) WebSocketFailed() {
	c.WriteErrorResponse(ErrWebSocketFailed)
}

func (c *HTTPConnection) CustomError(w http.ResponseWriter, status int, code int, message string, dMessage string) {
	c.WriteErrorResponse(NewError(status, code, message, dMessage))
}

func (c *HTTPConnection) WriteErrorResponse(apiError *Error) {
	c.WriteResourceResponse(c.w, apiError.Status, apiError)
}

func (c *HTTPConnection) WriteResourceResponse(w http.ResponseWriter, status int, resource interface{}) {
	w.WriteHeader(status)
	WriteJSON(w, resource)
}

func (c *HTTPConnection) Upgrade(u WebSocketUpgrader) (SocketConnection, error) {
	return u.Upgrade(c.w, c.r, c.Client())
}

func (c *HTTPConnection) ResponseWriter() http.ResponseWriter {
	return c.w
}

func (c *HTTPConnection) Request() *http.Request {
	return c.r
}
