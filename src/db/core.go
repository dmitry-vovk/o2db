// The file contains types for package and
// core database object with methods to handle databases
package db

import (
	. "types"
)

type Package struct {
	Container *Container
	Client    *Client
	RespChan  chan Response
}

type DbCore struct {
	databases map[string]*Database
	Input     chan *Package
}

// Goroutine that handles queries asynchronously
func (с *DbCore) Processor() {
	с.databases = make(map[string]*Database)
	for {
		pkg := <-с.Input
		pkg.RespChan <- с.ProcessRequest(pkg.Client, pkg.Container)
	}
}
