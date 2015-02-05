package transfer_test

import (
	"github.com/elos/data"
	. "github.com/elos/transfer"

	"errors"

	"github.com/elos/server/models"
	"github.com/elos/server/models/user"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NullConnection", func() {

	var (
		u models.User
	)

	BeforeEach(func() {
		u, _ = user.New(data.NewNullStoreWithType("mongo"))
	})

	Describe("NewNullConnection", func() {
		c := NewNullConnection(u)

		It("Allocates and returns a new NullConnection", func() {
			Expect(c).NotTo(BeNil())
			Expect(c).To(BeAssignableToTypeOf(&NullConnection{}))
		})

		It("Allocates Writes and Reads maps", func() {
			Expect(c.Writes).NotTo(BeNil())
			Expect(c.Reads).NotTo(BeNil())
			Expect(c.Writes).To(BeEmpty())
			Expect(c.Reads).To(BeEmpty())
		})
	})

	Describe("Reset()", func() {
		It("Properly resets necessary fields", func() {
			c := NewNullConnection(u)
			x := 1
			y := 2

			c.ReadJSON(x)
			c.WriteJSON(y)
			c.Close()
			c.Error = errors.New("error")

			cc := c.Reset()

			By("Returns itself")
			Expect(cc).To(Equal(c))
			By("Writes -> empty")
			Expect(c.Writes).To(BeEmpty())
			By("Reads -> empty")
			Expect(c.Reads).To(BeEmpty())
			By("Closed -> false")
			Expect(c.Closed).To(BeFalse())
			By("Error -> nil")
			Expect(c.Error).To(BeNil())
		})
	})

	Describe("SetError()", func() {
		It("sets the error field", func() {
			c := NewNullConnection(u)
			e := errors.New("er")
			c.SetError(e)
			Expect(c.Error).To(Equal(e))
		})
	})

	Describe("conn.Connection Implementation", func() {
		Describe("WriteJSON", func() {
			Context("Error", func() {
				Describe("Custom Error", func() {
					It("Returns its Error field", func() {
						c := NewNullConnection(u)

						e := errors.New("ealdjf")
						c.SetError(e)

						x := 142312
						Expect(c.WriteJSON(x)).To(Equal(e))
					})
				})
				Describe("ConnectionClosedError", func() {
					It("Returns ConnectionClosedError", func() {
						c := NewNullConnection(u)

						x := 122421

						c.Close()
						Expect(c.WriteJSON(x)).To(Equal(ConnectionClosedError))
					})
				})
			})

			Context("Success", func() {
				It("Adds the value to Writes field", func() {
					c := NewNullConnection(u)

					x := 21
					err := c.WriteJSON(x)
					Expect(err).ToNot(HaveOccurred())
					Expect(c.Writes).To(HaveKeyWithValue(x, true))
				})
			})
		})

		Describe("ReadJSON", func() {
			Context("Error", func() {
				Describe("Custom Error", func() {
					It("Returns its Error Field", func() {
						c := NewNullConnection(u)

						e := errors.New("er")
						c.SetError(e)
						Expect(c.ReadJSON(u)).To(Equal(e))
					})
				})

				Describe("ConnectionClosedError", func() {
					It("Returns ConnectionClosedError", func() {
						c := NewNullConnection(u)
						c.Close()
						Expect(c.ReadJSON(u)).To(Equal(ConnectionClosedError))
					})
				})
			})

			Context("Success", func() {
				It("Adds the value to Reads field", func() {
					c := NewNullConnection(u)
					x := 11231221 //anything, should, in theory be a struct
					err := c.ReadJSON(x)
					Expect(err).ToNot(HaveOccurred())
					Expect(c.Reads).To(HaveKeyWithValue(x, true))
				})
			})
		})

		Describe("Close()", func() {
			It("Succesfully closes the connection", func() {
				c := NewNullConnection(u)

				Expect(c.WriteJSON(u)).To(BeNil())
				Expect(c.ReadJSON(u)).To(BeNil())
				Expect(c.Closed).To(BeFalse())

				c.Close()
				defer c.Reset()
				By("Setting Closed field to true")
				Expect(c.Closed).To(BeTrue())
				By("Returning ConnectionClosedError on ReadJSON()")
				Expect(c.WriteJSON(u)).To(Equal(ConnectionClosedError))
				By("Returning ConnectionClosedError on WriteJSON()")
				Expect(c.ReadJSON(u)).To(Equal(ConnectionClosedError))
			})
		})

		Describe("Agent", func() {
			It("Returns the agent it was intialized with", func() {
				c := NewNullConnection(u)

				Expect(c.Agent()).To(Equal(u))
			})
		})
	})
})
