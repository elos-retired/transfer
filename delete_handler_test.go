package transfer_test

import (
	. "github.com/elos/transfer"

	"github.com/elos/data"
	"github.com/elos/server/models/user"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

/* Still needs to test:
* Bad envelope form
*Multiple data kinds ?
 */

var _ = Describe("DeleteHandler", func() {
	var (
		a    data.Identifiable
		e    *Envelope
		db   *data.TestDB
		c    *NullConnection
		info map[data.Kind]map[string]interface{}
		id   data.ID
	)

	BeforeEach(func() {

		a = user.New()

		info = map[data.Kind]map[string]interface{}{
			user.Kind: {
				"id": id.Hex(),
			},
		}

		e = data.NewEnvelope(DELETE, info)

		db = test.NewDB()
		id = db.NewObjectID()

		c = NewNullConnection(a)
	})

	Context("Successful removal", func() {
		BeforeEach(func() {
			DeleteHandler(e, db, c)
		})

		It("Attempts to remove the record from the database", func() {
			Expect(db.Deleted).To(HaveLen(1))
			Expect(db.Deleted[0]).To(BeAssignableToTypeOf(user.New()))
			Expect(db.Deleted[0].GetID()).To(Equal(id))
		})

		It("Writes nothing to the connection", func() {
			Expect(c.Writes).NotTo(BeEmpty(), "The handler should have written to the connection")

			lastWrite := c.LastWrite // get that last write

			Expect(lastWrite).NotTo(BeNil(), "Our null connection has a non-nil last write")

			pack, ok := c.LastWrite.(*data.Package) // typ asser to a data.Package pointer
			Expect(ok).To(BeTrue(), "The last write should have been of type *data.Package")
			Expect(pack.Action).To(Equal(DELETE), "The package's action should be DELETE")
			Expect(pack.Data[user.Kind].GetID()).To(Equal(id), "The package's model should have the id of the model we requested to delete")
		})
	})

	Context("Unsuccessful removal", func() {
		BeforeEach(func() {
			db.Error = true
			DeleteHandler(e, db, c)
		})

		It("Attempts, but fails, to remove the record from the database", func() {
			Expect(db.Deleted).To(BeEmpty())
		})

		It("Writes a POST back to the client with the data", func() {
			Expect(c.Writes).NotTo(BeEmpty(), "The handler should have written to the connection")

			lastWrite := c.LastWrite // get that last write
			Expect(lastWrite).NotTo(BeNil(), "The last write should be some non-nil value")

			pack, ok := lastWrite.(*data.Envelope)
			Expect(ok).To(BeTrue(), "The last write should have the type *data.Envelope")
			Expect(pack.Action).To(Equal(POST))
			Expect(pack.Data).To(Equal(info), "The handler parroted our delete request, simply changing DELETE->POST")
		})
	})
})
