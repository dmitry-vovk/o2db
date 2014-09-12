package client

import (
	"net"
	"db"
)

type ClientType struct {
	Conn          net.Conn
	Authenticated bool
	Db            *db.Database
}
