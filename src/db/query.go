// Core method that processes parsed message and returns response
package db

import (
	"errors"
	"fmt"
	. "logger"
	"reflect"
	. "types"
)

// This is the main entry for processing queries
func (this *DbCore) ProcessQuery(c *Client, q *Container) Response {
	if q == nil {
		return respond("no message", nil)
	}
	DebugLog.Printf("Payload type: %s", reflect.TypeOf(q.Payload))
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
		case CreateCollection:
			if clientDb, ok := this.databases[c.Db]; ok {
				if _, ok := clientDb.Collections[q.Payload.(CreateCollection).Name]; !ok {
					return respond("Collection created", clientDb.CreateCollection(q.Payload.(CreateCollection)))
				} else {
					return respond("Collection already exists", nil)
				}
			} else {
				return respond("Database not selected", nil)
			}
		case DropCollection:
			if clientDb, ok := this.databases[c.Db]; ok {
				if _, ok := clientDb.Collections[q.Payload.(WriteObject).Collection]; ok {
					return respond("Collection deleted", clientDb.DropCollection(q.Payload.(DropCollection)))
				} else {
					return respond("Collection does not exist", nil)
				}
			} else {
				return respond("Database not selected", nil)
			}
		case WriteObject:
			if clientDb, ok := this.databases[c.Db]; ok {
				collectionKey := hash(q.Payload.(WriteObject).Collection)
				if collection, ok := clientDb.Collections[collectionKey]; ok {
					return respond("Object written", collection.WriteObject(q.Payload.(WriteObject)))
				} else {
					return respond("Collection does not exist", nil)
				}
			} else {
				return respond("Database not selected", nil)
			}
		case ReadObject:
			if clientDb, ok := this.databases[c.Db]; ok {
				collectionKey := hash(q.Payload.(ReadObject).Collection)
				if collection, ok := clientDb.Collections[collectionKey]; ok {
					obj, err := collection.ReadObject(q.Payload.(ReadObject))
					return respond(obj, err)
				} else {
					return respond("Collection does not exist", nil)
				}
			} else {
				return respond("Database not selected", nil)
			}
		default:
			ErrorLog.Printf("Unknown query type [%s]", reflect.TypeOf(q.Payload))
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
