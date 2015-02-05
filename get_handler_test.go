package transfer_test

import (
	. "github.com/elos/transfer"

	"github.com/elos/data"
	"github.com/elos/server/models/user"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GetHandler", func() {
	var (
		a    data.Agent
		e    *data.Envelope
		db   *test.TestDB
		c    *conn.NullConnection
		info map[data.Kind]map[string]interface{}
		id   data.ID
	)

	BeforeEach(func() {
		id = data.NewObjectID()

		a = user.New()

		info = map[data.Kind]map[string]interface{}{
			user.Kind: {
				"id": id.Hex(),
			},
		}

		e = data.NewEnvelope(DELETE, info)

		db = test.NewDB()

		c = conn.NewNullConnection(a)
	})

	Context("Success", func() {
		It("hey", func() {
			Expect(true).To(BeTrue())
		})
	})

	Context("Failure", func() {
	})
})
