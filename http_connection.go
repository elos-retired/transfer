package transfer

import (
	"log"
	"net/http"

	"github.com/elos/data"
)

type HTTPConnection struct {
	w http.ResponseWriter
	r *http.Request
	a data.Identifiable
}

func NewHTTPConnection(w http.ResponseWriter, r *http.Request, a data.Identifiable) *HTTPConnection {
	return &HTTPConnection{
		w: w,
		r: r,
		a: a,
	}
}

func (c *HTTPConnection) WriteJSON(v interface{}) error {
	c.w.Header().Set("Content-Type", "application/json; charset=utf-8")

	bytes, err := ToJSON(v)
	if err != nil {
		return err
	}

	_, err = c.w.Write(bytes)

	log.Print("Hasdfasdfd")
	return err
}

func (c *HTTPConnection) Agent() data.Identifiable {
	return c.a
}

func (c *HTTPConnection) IsAuthenticated() bool {
	return c.Agent() != nil
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
