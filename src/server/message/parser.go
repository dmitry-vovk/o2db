// Incoming message parser
package message

import (
	"encoding/json"
	"errors"
	"fmt"
	. "types"
)

// Parse incoming JSON bytes into Package to be fed into query processor
func Parse(msg []byte) (*Container, error) {
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
	parsedMessage := &Container{}
	err = json.Unmarshal([]byte(*m["type"]), &parsedMessage.Type)
	if err != nil {
		return nil, err
	}
	payload := []byte(*m["payload"])
	switch parsedMessage.Type {
	case TypeAuthenticate:
		var p Authentication
		err = json.Unmarshal(payload, &p)
		if err == nil {
			parsedMessage.Payload = p
		}
	case TypeCreateDatabase:
		var p CreateDatabase
		err = json.Unmarshal(payload, &p)
		if err == nil {
			parsedMessage.Payload = p
		}
	case TypeOpenDatabase:
		var p OpenDatabase
		err = json.Unmarshal(payload, &p)
		if err == nil {
			parsedMessage.Payload = p
		}
	case TypeDropDatabase:
		var p DropDatabase
		err = json.Unmarshal(payload, &p)
		if err == nil {
			parsedMessage.Payload = p
		}
	case TypeListDatabases:
		var p ListDatabases
		err = json.Unmarshal(payload, &p)
		if err == nil {
			parsedMessage.Payload = p
		}
	case TypeCreateCollection:
		var p CreateCollection
		err = json.Unmarshal(payload, &p)
		if err == nil {
			parsedMessage.Payload = p
		}
	case TypeDropCollection:
		var p DropCollection
		err = json.Unmarshal(payload, &p)
		if err == nil {
			parsedMessage.Payload = p
		}
	case TypeObjectWrite:
		var p WriteObject
		err = json.Unmarshal(payload, &p)
		if err == nil {
			parsedMessage.Payload = p
		}
	default:
		return nil, errors.New(fmt.Sprintf("Unsupported message type: %d", parsedMessage.Type))
	}
	return parsedMessage, err
}
