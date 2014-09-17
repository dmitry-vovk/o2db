package db

import (
	. "types"
)

type Package struct {
	Container *Container
	Client    *ClientType
	RespChan  chan []byte
}

type DbCore struct {
	databases map[string]*Database
	Input     chan *Package
}

var (
	Core DbCore
)

func (this *DbCore) Processor() {
	for {
		pkg := <-this.Input
		pkg.RespChan <- this.ProcessQuery(pkg.Client, pkg.Container)
	}
}
