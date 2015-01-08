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
		return respond("no message", RNoError, nil)
	}
	DebugLog.Printf("Payload type: %s", reflect.TypeOf(query.Payload))
	switch query.Payload.(type) {
	case Authentication:
		if core.Authenticate(client, query.Payload.(Authentication)) {
			return respond("Authenticated", RAuthenticated, nil)
		} else {
			return respond("Authentication failed", RNotAuthenticated, nil)
		}
	}
	if subscribe, ok := query.Payload.(Subscribe); ok {
		dbName, err := core.OpenDatabase(OpenDatabase{subscribe.Database})
		if err != nil {
			return respond(nil, RDatabaseNotSelected, err)
		}
		client.Db = dbName
		collection, code, err := core.getCollection(client, subscribe.Collection)
		if err != nil {
			return respond(nil, code, err)
		}
		return respond(collection.Subscribe(subscribe, client))
	} else if client.Authenticated {
		switch query.Payload.(type) {
		case OpenDatabase:
			dbName, err := core.OpenDatabase(query.Payload.(OpenDatabase))
			if err == nil {
				client.Db = dbName
			}
			return respond("Database opened", RDatabaseSelected, err)
		case CreateDatabase:
			return respond("Database created", RDatabaseCreated, core.CreateDatabase(query.Payload.(CreateDatabase)))
		case DropDatabase:
			return respond("Database deleted", RDatabaseDeleted, core.DropDatabase(query.Payload.(DropDatabase)))
		case ListDatabases:
			resp, err := core.ListDatabases(query.Payload.(ListDatabases))
			return respond(resp, RDatabaseList, err)
		case CreateCollection:
			if clientDb, ok := core.databases[client.Db]; ok {
				if _, ok := clientDb.Collections[query.Payload.(CreateCollection).Name]; !ok {
					return respond("Collection created", RCollectionCreated, clientDb.CreateCollection(query.Payload.(CreateCollection)))
				} else {
					return respond("Collection already exists", RCollectionAlreadyExists, nil)
				}
			} else {
				return respond("Database not selected", RDatabaseNotSelected, nil)
			}
		case DropCollection:
			_, code, err := core.getCollection(client, query.Payload.(DropCollection).Name)
			if err != nil {
				return respond(nil, code, err)
			}
			return respond("Collection deleted", RCollectionDeleted, core.databases[client.Db].DropCollection(query.Payload.(DropCollection)))
		case WriteObject:
			collection, code, err := core.getCollection(client, query.Payload.(WriteObject).Collection)
			if err != nil {
				return respond(nil, code, err)
			}
			if code, err := collection.WriteObject(query.Payload.(WriteObject)); err == nil {
				data := query.Payload.(WriteObject).Data
				go collection.SubscriptionDispatcher(&data)
				return respond("Object written", RObjectWritten, nil)
			} else {
				return respond(nil, code, err)
			}
		case ReadObject:
			collection, code, err := core.getCollection(client, query.Payload.(ReadObject).Collection)
			if err != nil {
				return respond(nil, code, err)
			}
			return respond(collection.ReadObject(query.Payload.(ReadObject)))
		case GetObjectVersions:
			collection, code, err := core.getCollection(client, query.Payload.(GetObjectVersions).Collection)
			if err != nil {
				return respond(nil, code, err)
			}
			return respond(collection.GetObjectVersions(query.Payload.(GetObjectVersions)))
		case GetObjectDiff:
			collection, code, err := core.getCollection(client, query.Payload.(GetObjectDiff).Collection)
			if err != nil {
				return respond(nil, code, err)
			}
			return respond(collection.GetObjectDiff(query.Payload.(GetObjectDiff)))
		case SelectObjects:
			collection, code, err := core.getCollection(client, query.Payload.(SelectObjects).Collection)
			if err != nil {
				return respond(nil, code, err)
			}
			return respond(collection.SelectObjects(query.Payload.(SelectObjects)))
		case AddSubscription:
			collection, code, err := core.getCollection(client, query.Payload.(AddSubscription).Collection)
			if err != nil {
				return respond(nil, code, err)
			}
			return respond(collection.AddSubscription(query.Payload.(AddSubscription)))
		case CancelSubscription:
			collection, code, err := core.getCollection(client, query.Payload.(CancelSubscription).Collection)
			if err != nil {
				return respond(nil, code, err)
			}
			return respond(collection.CancelSubscription(query.Payload.(CancelSubscription)))
		case ListSubscriptions:
			var subscriptions []SubscriptionItem
			for _, collectionName := range query.Payload.(ListSubscriptions).Collections {
				collection, code, err := core.getCollection(client, collectionName)
				if err != nil {
					return respond(nil, code, err)
				}
				subscriptions = append(subscriptions, collection.ListSubscriptions()...)
			}
			return respond(subscriptions, RSubscriptionList, nil)
		default:
			ErrorLog.Printf("Unknown query type [%s]", reflect.TypeOf(query.Payload))
			return respond(nil, RUnknownQueryType, errors.New(fmt.Sprintf("Unknown query type [%s]", reflect.TypeOf(query.Payload))))
		}
	}
	return respond("Authentication required", RAuthenticationRequired, nil)
}

func (core *DbCore) getCollection(client *Client, collectionName string) (*Collection, uint, error) {
	if clientDb, ok := core.databases[client.Db]; ok {
		collectionKey := hash(collectionName)
		if collection, ok := clientDb.Collections[collectionKey]; ok {
			return collection, RNoError, nil
		} else {
			return nil, RCollectionDoesNotExist, errors.New("Collection " + collectionName + " does not exist")
		}
	} else {
		return nil, RDatabaseNotSelected, errors.New("Database not selected")
	}
}

// Wraps response structure and error into JSON
func respond(r interface{}, code uint, e error) Response {
	if e == nil {
		return Response{
			Result:   true,
			Code:     code,
			Response: r,
		}
	} else {
		return Response{
			Result:   false,
			Code:     code,
			Response: fmt.Sprintf("%s", e),
		}
	}
}
