package client

import (
	"net"
	"db/schema"
)

type ClientType struct {
	Conn          net.Conn
	Authenticated bool
	Db            *schema.Database
}
