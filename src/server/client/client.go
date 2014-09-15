package client

import (
	"db/schema"
	"net"
)

type ClientType struct {
	Conn          net.Conn
	Authenticated bool
	Db            *schema.Database
}
