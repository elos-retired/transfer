package transfer

import (
	"encoding/json"
	"net/http"
)

func ToJSON(v interface{}) ([]byte, error) {
	// Always pretty-print JSON
	return json.MarshalIndent(v, "", "    ")
}

func SetContentJSON(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
}

/*
	Helper function that writes an interface as JSON
	- Takes care of nominal things such as setting the content header
*/
func WriteJSON(w http.ResponseWriter, resource interface{}) error {
	SetContentJSON(w)

	bytes, err := ToJSON(resource)
	if err != nil {
		return err
	}

	_, err = w.Write(bytes)
	return err
}
