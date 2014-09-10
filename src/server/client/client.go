package client

import "net"

type ClientType struct {
	Conn          net.Conn
	Authenticated bool
}
