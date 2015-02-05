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

func NewHTTPConnection(w http.ResponseWriter, r *http.Request, a data.Identifiable) Connection {
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
