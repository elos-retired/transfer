package transfer_test

import (
	"encoding/json"
	. "github.com/elos/transfer"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Json", func() {

	type TestStruct struct {
		Hey string
	}

	x := TestStruct{Hey: "World"}

	Describe("ToJSON", func() {
		It("Properly marshals a struct", func() {
			// Not tested too rigorously because this basically relies on std lib's json pkg

			y := TestStruct{}

			bytesArray, err := ToJSON(x)
			Expect(err).ToNot(HaveOccurred())

			err = json.Unmarshal(bytesArray, &y)
			Expect(err).ToNot(HaveOccurred())

			Expect(y).To(Equal(x))
		})
	})

	Describe("SetContentJson", func() {
		It("Properly sets content header", func() {
			w := httptest.NewRecorder()
			SetContentJSON(w)
			Expect(w.HeaderMap["Content-Type"][0]).To(Equal("application/json; charset=utf-8"))
		})
	})

	Describe("WriteJSON", func() {
		It("Properly writes JSON", func() {
			w := httptest.NewRecorder()
			WriteJSON(w, x)

			bytesArray, err := ToJSON(x)
			Expect(err).ToNot(HaveOccurred())

			Expect(w.Body.String()).To(Equal(string(bytesArray)))
		})
	})
})
