package db

import "net"

type Database struct {
	DataDir     string
	Collections map[string]Collection
}

type ClientType struct {
	Conn          net.Conn
	Authenticated bool
	Db            *Database
}
