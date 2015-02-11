package transfer

import (
	"encoding/json"
	"net/http/httptest"
	"testing"
)

func TestToJSON(t *testing.T) {

	x := struct{ Hey string }{Hey: "World"}

	// Not tested too rigorously because this basically relies on std lib's json pkg
	y := new(struct{ Hey string })

	bytesArray, err := ToJSON(x)
	if err != nil {
		t.Errorf("ToJSON returned an error: %s", err)
	}

	err = json.Unmarshal(bytesArray, &y)
	if err != nil {
		t.Errorf("json.Unmarshal returned an error: %s", err)
	}

	if *y != x {
		t.Errorf("ToJSON != json.Unmarshal")
	}
}

func TestSetContentJSON(t *testing.T) {
	w := httptest.NewRecorder()
	SetContentJSON(w)
	if w.HeaderMap[ContentTypeHeader][0] != JSONContentType {
		t.Errorf("SetContentJSON failed to set the content type correctly")
	}
}

func TestWriteJSON(t *testing.T) {
	x := struct{ Hey string }{Hey: "World"}
	w := httptest.NewRecorder()
	WriteJSON(w, x)

	bytesArray, err := ToJSON(x)
	if err != nil {
		t.Errorf("ToJSON failed")
	}

	if w.Body.String() != string(bytesArray) {
		t.Errorf("WriteJSON Failed")
	}
}
