package transfer_test

import (
	"github.com/elos/data"
	"github.com/elos/models/user"
	. "github.com/elos/transfer"

	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("WebSocket", func() {

	It("Defines the WebSocketProtocolHeader", func() {
		Expect(WebSocketProtocolHeader).NotTo(BeNil())
	})

	It("Defines a Default websocket upgrader", func() {
		Expect(DefaultWebSocketUpgrader).NotTo(BeNil())
	})

	// NullUpgrader {{{
	Describe("NullUpgrader", func() {
		var (
			a data.Identifiable
			c *NullConnection
			u *NullUpgrader

			w http.ResponseWriter
			r *http.Request
		)

		BeforeEach(func() {
			a, _ = user.New(data.NewNullStore())
			c = NewNullConnection(a)
			u = NewNullUpgrader(c)

			w = httptest.NewRecorder()
			r = &http.Request{}
		})

		Describe("NewNullUpgrader()", func() {
			It("Allocates and returns a new *NullUpgrader", func() {
				By("Returning a non-nil pointer")
				Expect(u).NotTo(BeNil())
				By("Returns a type of *NullUpgrader")
				Expect(u).To(BeAssignableToTypeOf(&NullUpgrader{}))
				By("Allocating Upgraded")
				Expect(u.Upgraded).NotTo(BeNil())
				Expect(u.Upgraded).To(BeEmpty())
				By("Setting Connection to connection passed in")
				Expect(u.Connection).To(Equal(c))
			})
		})

		Describe("Reset()", func() {
			It("resets all crucial fields", func() {
				// Example use, both formally tested elsewhere
				u.Upgrade(w, r, a)
				u.SetError(errors.New("asdf"))

				// Reset it
				cu := u.Reset()
				By("Returns itself")
				Expect(cu).To(Equal(u))
				By("Upgraded -> empty")
				Expect(cu.Upgraded).To(BeEmpty())
				By("Error -> nil")
				Expect(cu.Error).To(BeNil())
			})
		})

		Describe("SetError", func() {
			It("Sets the Error field to the error supplied", func() {
				e := errors.New("asdf")
				u.SetError(e)
				By("Setting Error field to the error")
				Expect(u.Error).To(Equal(e))
			})
		})

		Describe("WebSocketUpgrader Implmentation", func() {
			Describe("Upgrade", func() {
				Context("Error", func() {
					It("returns the error instructed to", func() {
						e := errors.New("example error")
						u.SetError(e)
						connection, err := u.Upgrade(w, r, a)
						By("Returning a nil connection")
						Expect(connection).To(BeNil())
						By("Returning the error")
						Expect(err).To(HaveOccurred())
						Expect(err).To(Equal(e))
					})
				})

				Context("Success", func() {
					It("Upgrades the request and returns a connection", func() {
						connection, err := u.Upgrade(w, r, a)
						By("Error returned being nil")
						Expect(err).ToNot(HaveOccurred())
						By("Adding the request to Upgraded")
						Expect(u.Upgraded).To(HaveKeyWithValue(r, true))
						By("Returning the connection initialized with")
						Expect(connection).To(Equal(c))
					})
				})
			})
		})
	})
	// NullUpgrader }}}

	// GorillaUpgrader {{{
	Describe("GorillaUpgrader", func() {
		var (
			ReadBufferSize  int
			WriteBufferSize int
			CheckOrigin     bool
			url             *url.URL
			r               *http.Request
			w               *httptest.ResponseRecorder
			u               *GorillaUpgrader
			a               data.Identifiable
		)

		BeforeEach(func() {
			a, _ = user.New(data.NewNullStore())
			ReadBufferSize = 1024
			WriteBufferSize = 2014
			CheckOrigin = true

			var err error
			url, err = url.Parse("http://localhost:8000/v1/upgrade")
			if err != nil {
				Fail("Couldn't parse example URL")
			}

			r = &http.Request{URL: url}
			w = httptest.NewRecorder()

			u = NewGorillaUpgrader(ReadBufferSize, WriteBufferSize, CheckOrigin)
		})

		Describe("NewGorillaUpgrader", func() {
			It("Allocates and returns a new *GorillaUpgrader", func() {
				Expect(u).NotTo(BeNil(), "New gorilla upgrader pointer was nil")
				Expect(u).To(BeAssignableToTypeOf(&GorillaUpgrader{}))
				Expect(u.Upgrader).NotTo(BeNil(), "GorillaUpgrader has a nil Upgrader")
			})
		})

		Describe("WebSocketUpgrader Implmentation", func() {
			Describe("Upgrade", func() {
				It("Forwards the upgrade attempty to a gorilla upgrader", func() {
					wc := httptest.NewRecorder()
					u.Upgrade(w, r, a)
					u.Upgrader.Upgrade(wc, r, ExtractProtocolHeader(r))
					Expect(wc.Body).To(Equal(w.Body))
					Expect(wc.Code).To(Equal(w.Code))
				})
			})
		})
	})
	// GorillaUpgrader }}}

	// ExtractProtocolHeader {{{
	Describe("ExtractProtocolHeader", func() {
		p := "askldfjasdjfkjalsdfljkasdjkfasdflkaf"

		header := http.Header{}
		header.Add(WebSocketProtocolHeader, p)

		r := &http.Request{
			Header: header,
		}

		It("Extracts the websocket protocol header into a new header objects", func() {
			h := ExtractProtocolHeader(r)
			Expect(h).NotTo(BeNil())
			Expect(h).To(BeAssignableToTypeOf(http.Header{}))
			Expect(h.Get(WebSocketProtocolHeader)).To(Equal(p))
		})
	})
	// ExtractProtocolHeader }}}

})
