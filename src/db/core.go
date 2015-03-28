// The file contains types for package and
// core database object with methods to handle databases
package db

import (
	"client"
	"fmt"
	. "logger"
	"server/message"
	. "types"
)

type Package struct {
	RawInput []byte
	Client   *client.Client
	RespChan chan Response
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
		pkg.Client.QueryCounter++
		container, err := с.parse(pkg.RawInput)
		if err == nil {
			pkg.RespChan <- с.ProcessRequest(pkg.Client, container)
		} else {
			ErrorLog.Printf("Parse error: %s", err)
			ErrorLog.Printf("Message was: %s", pkg.RawInput)
			pkg.Client.Respond(
				Response{
					Result:   false,
					Code:     RQueryParseError,
					Response: fmt.Sprintf("%s", err),
				},
			)
		}
	}
}

// Convert raw bytes into message container
func (c *DbCore) parse(msg []byte) (*Container, error) {
	container, err := message.Parse(msg)
	if err != nil {
		return nil, err
	}
	return container, nil
}
