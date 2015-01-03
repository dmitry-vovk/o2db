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
func (core *DbCore) ProcessRequest(client *Client, query *Container) Response {
	if query == nil {
		return respond("no message", nil)
	}
	DebugLog.Printf("Payload type: %s", reflect.TypeOf(query.Payload))
	switch query.Payload.(type) {
	case Authentication:
		if core.Authenticate(client, query.Payload.(Authentication)) {
			return respond("Authenticated", nil)
		} else {
			return respond("Authentication failed", nil)
		}
	}
	if subscribe, ok := query.Payload.(Subscribe); ok {
		collection, err := core.getCollection(client, subscribe.Collection)
		if err != nil {
			return respond(nil, err)
		}
		return respond(collection.Subscribe(subscribe, client))
	} else if client.Authenticated {
		switch query.Payload.(type) {
		case OpenDatabase:
			dbName, err := core.OpenDatabase(query.Payload.(OpenDatabase))
			if err == nil {
				client.Db = dbName
			}
			return respond("Database opened", err)
		case CreateDatabase:
			return respond("Database created", core.CreateDatabase(query.Payload.(CreateDatabase)))
		case DropDatabase:
			return respond("Database deleted", core.DropDatabase(query.Payload.(DropDatabase)))
		case ListDatabases:
			resp, err := core.ListDatabases(query.Payload.(ListDatabases))
			return respond(resp, err)
		case CreateCollection:
			if clientDb, ok := core.databases[client.Db]; ok {
				if _, ok := clientDb.Collections[query.Payload.(CreateCollection).Name]; !ok {
					return respond("Collection created", clientDb.CreateCollection(query.Payload.(CreateCollection)))
				} else {
					return respond("Collection already exists", nil)
				}
			} else {
				return respond("Database not selected", nil)
			}
		case DropCollection:
			_, err := core.getCollection(client, query.Payload.(DropCollection).Name)
			if err != nil {
				return respond(nil, err)
			}
			return respond("Collection deleted", core.databases[client.Db].DropCollection(query.Payload.(DropCollection)))
		case WriteObject:
			collection, err := core.getCollection(client, query.Payload.(WriteObject).Collection)
			if err != nil {
				return respond(nil, err)
			}
			err = collection.WriteObject(query.Payload.(WriteObject))
			if err == nil {
				data := query.Payload.(WriteObject).Data
				collection.SubscriptionDispatcher(&data)
			}
			return respond("Object written", err)
		case ReadObject:
			collection, err := core.getCollection(client, query.Payload.(ReadObject).Collection)
			if err != nil {
				return respond(nil, err)
			}
			return respond(collection.ReadObject(query.Payload.(ReadObject)))
		case GetObjectVersions:
			collection, err := core.getCollection(client, query.Payload.(GetObjectVersions).Collection)
			if err != nil {
				return respond(nil, err)
			}
			return respond(collection.GetObjectVersions(query.Payload.(GetObjectVersions)))
		case GetObjectDiff:
			collection, err := core.getCollection(client, query.Payload.(GetObjectDiff).Collection)
			if err != nil {
				return respond(nil, err)
			}
			return respond(collection.GetObjectDiff(query.Payload.(GetObjectDiff)))
		case SelectObjects:
			collection, err := core.getCollection(client, query.Payload.(SelectObjects).Collection)
			if err != nil {
				return respond(nil, err)
			}
			return respond(collection.SelectObjects(query.Payload.(SelectObjects)))
		case AddSubscription:
			collection, err := core.getCollection(client, query.Payload.(AddSubscription).Collection)
			if err != nil {
				return respond(nil, err)
			}
			return respond(collection.AddSubscription(query.Payload.(AddSubscription)))
		case CancelSubscription:
			collection, err := core.getCollection(client, query.Payload.(CancelSubscription).Collection)
			if err != nil {
				return respond(nil, err)
			}
			return respond(collection.CancelSubscription(query.Payload.(CancelSubscription)))
		case ListSubscriptions:
			var subscriptions []SubscriptionItem
			var err error
			for _, collectionName := range query.Payload.(ListSubscriptions).Collections {
				collection, err := core.getCollection(client, collectionName)
				if err != nil {
					return respond(nil, err)
				}
				subscriptions = append(subscriptions, collection.ListSubscriptions()...)
			}
			return respond(subscriptions, err)
		default:
			ErrorLog.Printf("Unknown query type [%s]", reflect.TypeOf(query.Payload))
			return respond(nil, errors.New(fmt.Sprintf("Unknown query type [%s]", reflect.TypeOf(query.Payload))))
		}
	}
	return respond("Authentication required", nil)
}

func (core *DbCore) getCollection(client *Client, collectionName string) (*Collection, error) {
	if clientDb, ok := core.databases[client.Db]; ok {
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
