package db

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	. "types"
)

// This is the main entry for processing queries
func (this *DbCore) ProcessQuery(c *Client, q *Container) Response {
	if q == nil {
		return respond("no message", nil)
	}
	log.Printf("Payload type: %s", reflect.TypeOf(q.Payload))
	switch q.Payload.(type) {
	case Authentication:
		if this.Authenticate(c, q.Payload.(Authentication)) {
			return respond("Authenticated", nil)
		} else {
			return respond("Authentication failed", nil)
		}
	}
	if c.Authenticated {
		switch q.Payload.(type) {
		// Database operations
		case OpenDatabase:
			dbName, err := this.OpenDatabase(q.Payload.(OpenDatabase))
			if err == nil {
				c.Db = dbName
			}
			return respond("Database opened", err)
		case CreateDatabase:
			return respond("Database created", this.CreateDatabase(q.Payload.(CreateDatabase)))
		case DropDatabase:
			return respond("Database deleted", this.DropDatabase(q.Payload.(DropDatabase)))
		case ListDatabases:
			resp, err := this.ListDatabases(q.Payload.(ListDatabases))
			return respond(resp, err)
			// Collection operations
		case CreateCollection:
			return respond("Collection created", this.databases[c.Db].CreateCollection(q.Payload.(CreateCollection)))
		case DropCollection:
			return respond("Collection deleted", this.databases[c.Db].DropCollection(q.Payload.(DropCollection)))
			// Object operations
		case WriteObject:
			return respond("Object written", this.databases[c.Db].Collections[q.Payload.(WriteObject).Collection].WriteObject(q.Payload.(WriteObject)))
			// Default stub
		default:
			log.Printf("Unknown query type [%s]", reflect.TypeOf(q.Payload))
			return respond(nil, errors.New(fmt.Sprintf("Unknown query type [%s]", reflect.TypeOf(q.Payload))))
		}
	}
	return respond("Authentication required", nil)
}

// Wraps response structure and error into JSON
func respond(r interface{}, e error) Response {
	if e == nil {
		return Response{
			Result:   true,
			Response: r,
		}
	} else {
		return Response{
			Result:   false,
			Response: fmt.Sprintf("%s", e),
		}
	}
}
