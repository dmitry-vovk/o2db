// Incoming message parser
package message

import (
	"encoding/json"
	"errors"
	"fmt"
	"types"
)

// Parse incoming JSON bytes into Package to be fed into query processor
func Parse(msg []byte) (*types.Container, error) {
	var m map[string]*json.RawMessage
	err := json.Unmarshal(msg, &m)
	if err != nil {
		return nil, errors.New("Cannot parse message")
	}
	if _, ok := m["type"]; !ok {
		return nil, errors.New("Unknown message format: missing type field.")
	}
	if _, ok := m["payload"]; !ok {
		return nil, errors.New("Unknown message format: missing payload field.")
	}
	parsedMessage := &types.Container{}
	err = json.Unmarshal([]byte(*m["type"]), &parsedMessage.Type)
	if err != nil {
		return nil, err
	}
	payload := []byte(*m["payload"])
	switch parsedMessage.Type {
	case types.TypeAuthenticate:
		var p types.Authentication
		err = json.Unmarshal(payload, &p)
		if err == nil {
			parsedMessage.Payload = p
		}
	case types.TypeCreateDatabase:
		var p types.CreateDatabase
		err = json.Unmarshal(payload, &p)
		if err == nil {
			parsedMessage.Payload = p
		}
	case types.TypeOpenDatabase:
		var p types.OpenDatabase
		err = json.Unmarshal(payload, &p)
		if err == nil {
			parsedMessage.Payload = p
		}
	case types.TypeDropDatabase:
		var p types.DropDatabase
		err = json.Unmarshal(payload, &p)
		if err == nil {
			parsedMessage.Payload = p
		}
	case types.TypeListDatabases:
		var p types.ListDatabases
		err = json.Unmarshal(payload, &p)
		if err == nil {
			parsedMessage.Payload = p
		}
	case types.TypeCreateCollection:
		var p types.CreateCollection
		err = json.Unmarshal(payload, &p)
		if err == nil {
			parsedMessage.Payload = p
		}
	case types.TypeDropCollection:
		var p types.DropCollection
		err = json.Unmarshal(payload, &p)
		if err == nil {
			parsedMessage.Payload = p
		}
	case types.TypeObjectWrite:
		var p types.WriteObject
		err = json.Unmarshal(payload, &p)
		if err == nil {
			parsedMessage.Payload = p
		}
	case types.TypeObjectGet:
		var p types.ReadObject
		err = json.Unmarshal(payload, &p)
		if err == nil {
			parsedMessage.Payload = p
		}
	case types.TypeGetObjectVersions:
		var p types.GetObjectVersions
		err = json.Unmarshal(payload, &p)
		if err == nil {
			parsedMessage.Payload = p
		}
	case types.TypeGetObjectDiff:
		var p types.GetObjectDiff
		err = json.Unmarshal(payload, &p)
		if err == nil {
			parsedMessage.Payload = p
		}
	case types.TypeSelectObjects:
		var p types.SelectObjects
		err = json.Unmarshal(payload, &p)
		if err == nil {
			parsedMessage.Payload = p
		}
	default:
		return nil, errors.New(fmt.Sprintf("Unsupported message type: %d", parsedMessage.Type))
	}
	return parsedMessage, err
}
