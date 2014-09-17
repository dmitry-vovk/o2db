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
func ProcessQuery(c *ClientType, q *Container) []byte {
	if q == nil {
		return respond("no message", nil)
	}
	log.Printf("Payload type: %s", reflect.TypeOf(q.Payload))
	switch q.Payload.(type) {
	case TAuthenticate:
		if Authenticate(c, q.Payload.(TAuthenticate)) {
			return respond("Authenticated", nil)
		} else {
			return respond("Authentication failed", nil)
		}
	}
	if c.Authenticated {
		switch q.Payload.(type) {
		// Database operations
		case OpenDatabase:
			dbPtr, err := OpenDatabase(q.Payload.(OpenDatabase))
			if err == nil {
				c.Db = dbPtr
			}
			return respond("Database opened", err)
		case CreateDatabase:
			return respond("Database created", CreateDatabase(q.Payload.(CreateDatabase)))
		case DropDatabase:
			return respond("Database deleted", DropDatabase(q.Payload.(DropDatabase)))
		case ListDatabases:
			resp, err := ListDatabases(q.Payload.(ListDatabases))
			return respond(resp, err)
		// Collection operations
		case CreateCollection:
			return respond("Collection created", CreateCollection(c, q.Payload.(CreateCollection)))
		case DropCollection:
			return respond("Collection deleted", DropCollection(c, q.Payload.(DropCollection)))
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
