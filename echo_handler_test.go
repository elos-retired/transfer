package transfer_test

import (
	. "github.com/elos/transfer"

	"github.com/elos/data"
	"github.com/elos/server/models/user"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("EchoHandler", func() {
	It("Responds with the exact envelope given it", func() {
		a := user.New()
		c := conn.NewNullConnection(a)
		db := test.NewDB()
		id := data.NewObjectID()

		info := map[data.Kind]map[string]interface{}{
			user.Kind: {
				"id": id.Hex(),
			},
		}

		e := data.NewEnvelope(POST, info)
		EchoHandler(e, db, c)

		Expect(c.LastWrite).To(Equal(e))
	})
})
