package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
	. "types"
)

// This is the main entry for processing queries
func (this *DbCore) ProcessQuery(c *ClientType, q *Container) []byte {
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
			dbPtr, err := this.OpenDatabase(q.Payload.(OpenDatabase))
			if err == nil {
				c.Db = dbPtr
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
			return respond("Collection created", this.CreateCollection(c, q.Payload.(CreateCollection)))
		case DropCollection:
			return respond("Collection deleted", this.DropCollection(c, q.Payload.(DropCollection)))
		// Default stub
		default:
			log.Printf("Unknown query type [%s]", reflect.TypeOf(q.Payload))
			return respond(nil, errors.New(fmt.Sprintf("Unknown query type [%s]", reflect.TypeOf(q.Payload))))
		}
	}
	return respond("Authentication required", nil)
}

// Wraps response structure and error into JSON
func respond(r interface{}, e error) []byte {
	resp := Response{}
	if e == nil {
		resp.Result = true
		resp.Response = r
	} else {
		resp.Result = false
		resp.Response = fmt.Sprintf("%s", e)
	}
	out, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Error encoding response: %s", err)
		return []byte{}
	}
	log.Printf("Response: %s", out)
	return out
}
