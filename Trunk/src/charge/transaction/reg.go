package transaction

import "charge/server"

func init() {
	server.RegisterTransaction("ios", NewAppleTransaction)
}
