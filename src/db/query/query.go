package query

import (
	"db"
	"db/auth"
	"fmt"
	"server/client"
	"server/types"
	"reflect"
	"log"
	"encoding/json"
	"errors"
)

// This is the main entry for processing queries
func ProcessQuery(c *client.ClientType, q *types.Container) []byte {
	if q == nil {
		return respond("no message", nil)
	}
	log.Printf("Payload type: %s", reflect.TypeOf(q.Payload))
	switch q.Payload.(type) {
	case types.Authenticate:
		if auth.Authenticate(c, q.Payload.(types.Authenticate)) {
			return respond("Authenticated", nil)
		} else {
			return respond("Authentication failed", nil)
		}
	}
	if c.Authenticated {
		switch q.Payload.(type) {
		// Database operations
		case types.OpenDatabase:
			dbPtr, err := db.OpenDatabase(q.Payload.(types.OpenDatabase))
			if err == nil {
				c.Db = dbPtr
			}
			return respond(types.ResponseMessage{Message:"Database opened"}, err)
		case types.CreateDatabase:
			return respond("Database created", db.CreateDatabase(q.Payload.(types.CreateDatabase)))
		case types.DropDatabase:
			return respond("Database deleted", db.DropDatabase(q.Payload.(types.DropDatabase)))
		case types.ListDatabases:
			resp, err := db.ListDatabases(q.Payload.(types.ListDatabases))
			return respond(resp, err)
		// Collection operations
		case types.CreateCollection:
			return respond("Collection created", db.CreateCollection(c, q.Payload.(types.CreateCollection)))
		case types.DropCollection:
			return respond("Collection deleted", db.DropCollection(c, q.Payload.(types.DropCollection)))
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
	resp := types.Response{}
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
