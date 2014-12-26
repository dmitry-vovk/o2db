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
func (d *DbCore) ProcessRequest(c *Client, q *Container) Response {
	if q == nil {
		return respond("no message", nil)
	}
	DebugLog.Printf("Payload type: %s", reflect.TypeOf(q.Payload))
	switch q.Payload.(type) {
	case Authentication:
		if d.Authenticate(c, q.Payload.(Authentication)) {
			return respond("Authenticated", nil)
		} else {
			return respond("Authentication failed", nil)
		}
	}
	if subscribe, ok := q.Payload.(Subscribe); ok {
		collection, err := d.getCollection(c, subscribe.Collection)
		if err != nil {
			return respond(nil, err)
		}
		return respond(collection.Subscribe(subscribe, c))
	} else if c.Authenticated {
		switch q.Payload.(type) {
		case OpenDatabase:
			dbName, err := d.OpenDatabase(q.Payload.(OpenDatabase))
			if err == nil {
				c.Db = dbName
			}
			return respond("Database opened", err)
		case CreateDatabase:
			return respond("Database created", d.CreateDatabase(q.Payload.(CreateDatabase)))
		case DropDatabase:
			return respond("Database deleted", d.DropDatabase(q.Payload.(DropDatabase)))
		case ListDatabases:
			resp, err := d.ListDatabases(q.Payload.(ListDatabases))
			return respond(resp, err)
		case CreateCollection:
			if clientDb, ok := d.databases[c.Db]; ok {
				if _, ok := clientDb.Collections[q.Payload.(CreateCollection).Name]; !ok {
					return respond("Collection created", clientDb.CreateCollection(q.Payload.(CreateCollection)))
				} else {
					return respond("Collection already exists", nil)
				}
			} else {
				return respond("Database not selected", nil)
			}
		case DropCollection:
			_, err := d.getCollection(c, q.Payload.(DropCollection).Name)
			if err != nil {
				return respond(nil, err)
			}
			return respond("Collection deleted", d.databases[c.Db].DropCollection(q.Payload.(DropCollection)))
		case WriteObject:
			collection, err := d.getCollection(c, q.Payload.(WriteObject).Collection)
			if err != nil {
				return respond(nil, err)
			}
			err = collection.WriteObject(q.Payload.(WriteObject))
			if err == nil {
				data := q.Payload.(WriteObject).Data
				collection.SubscriptionDispatcher(&data)
			}
			return respond("Object written", err)
		case ReadObject:
			collection, err := d.getCollection(c, q.Payload.(ReadObject).Collection)
			if err != nil {
				return respond(nil, err)
			}
			return respond(collection.ReadObject(q.Payload.(ReadObject)))
		case GetObjectVersions:
			collection, err := d.getCollection(c, q.Payload.(GetObjectVersions).Collection)
			if err != nil {
				return respond(nil, err)
			}
			return respond(collection.GetObjectVersions(q.Payload.(GetObjectVersions)))
		case GetObjectDiff:
			collection, err := d.getCollection(c, q.Payload.(GetObjectDiff).Collection)
			if err != nil {
				return respond(nil, err)
			}
			return respond(collection.GetObjectDiff(q.Payload.(GetObjectDiff)))
		case SelectObjects:
			collection, err := d.getCollection(c, q.Payload.(SelectObjects).Collection)
			if err != nil {
				return respond(nil, err)
			}
			return respond(collection.SelectObjects(q.Payload.(SelectObjects)))
		case AddSubscription:
			collection, err := d.getCollection(c, q.Payload.(AddSubscription).Collection)
			if err != nil {
				return respond(nil, err)
			}
			return respond(collection.AddSubscription(q.Payload.(AddSubscription)))
		case CancelSubscription:
			collection, err := d.getCollection(c, q.Payload.(CancelSubscription).Collection)
			if err != nil {
				return respond(nil, err)
			}
			return respond(collection.CancelSubscription(q.Payload.(CancelSubscription)))
		default:
			ErrorLog.Printf("Unknown query type [%s]", reflect.TypeOf(q.Payload))
			return respond(nil, errors.New(fmt.Sprintf("Unknown query type [%s]", reflect.TypeOf(q.Payload))))
		}
	}
	return respond("Authentication required", nil)
}

func (d *DbCore) getCollection(c *Client, collectionName string) (*Collection, error) {
	if clientDb, ok := d.databases[c.Db]; ok {
		collectionKey := hash(collectionName)
		if collection, ok := clientDb.Collections[collectionKey]; ok {
			return collection, nil
		} else {
			return nil, errors.New("Collection " + collectionName + " does not exist")
		}
	} else {
		return nil, errors.New("Database not selected")
	}
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
