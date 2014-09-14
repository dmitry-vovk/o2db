package query

import (
	"db"
	"db/auth"
	"fmt"
	"server/client"
	"server/types"
	"reflect"
	"log"
)

// This is the main entry for processing queries
func ProcessQuery(c *client.ClientType, q *types.Container) interface{} {
	if q == nil {
		return "no message"
	}
	log.Printf("Payload type: %s", reflect.TypeOf(q.Payload))
	switch q.Payload.(type) {
	case types.Authenticate:
		if auth.Authenticate(c, q.Payload.(types.Authenticate)) {
			return "Authenticated"
		} else {
			return "Authentication failed"
		}
	}
	if c.Authenticated {
		switch q.Payload.(type) {
		// Database operations
		case types.OpenDatabase:
			dbPtr, err := db.OpenDatabase(q.Payload.(types.OpenDatabase))
			if err == nil {
				c.Db = dbPtr
				return "Database opened"
			}
			return string(fmt.Sprintf("%s", err))
		case types.CreateDatabase:
			err := db.CreateDatabase(q.Payload.(types.CreateDatabase))
			if err == nil {
				return "Database created"
			}
			return string(fmt.Sprintf("%s", err))
		case types.DropDatabase:
			err := db.DropDatabase(q.Payload.(types.DropDatabase))
			if err == nil {
				return "Database deleted"
			}
			return string(fmt.Sprintf("%s", err))
		case types.ListDatabases:
			resp, err := db.ListDatabases(q.Payload.(types.ListDatabases))
			if err == nil {
				return resp
			}
			return string(fmt.Sprintf("%s", err))
		// Collection operations
		case types.CreateCollection:
			err := db.CreateCollection(c, q.Payload.(types.CreateCollection))
			if err == nil {
				return "Collection created"
			}
			return string(fmt.Sprintf("%s", err))
		case types.DropCollection:
			err := db.DropCollection(c, q.Payload.(types.DropCollection))
			if err == nil {
				return "Collection deleted"
			}
			return string(fmt.Sprintf("%s", err))
		// Default stub
		default:
			log.Printf("Unknown query type [%s]", reflect.TypeOf(q.Payload))
			return string(fmt.Sprintf("Unknown query type [%s]", reflect.TypeOf(q.Payload)))
		}
	}
	return "Authentication required"
}
