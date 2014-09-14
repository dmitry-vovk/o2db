package message

import (
	"encoding/json"
	"server/types"
)

func Parse(msg []byte) (*types.Container, error) {
	var m map[string]*json.RawMessage
	err := json.Unmarshal(msg, &m)
	if err != nil {
		return nil, err
	}
	var parsedMessage types.Container
	err = json.Unmarshal(*m["type"], &parsedMessage.Type)
	switch parsedMessage.Type {
	case types.TypeAuthenticate:
		var p types.Authenticate
		err = json.Unmarshal(*m["payload"], &p)
		if err == nil {
			parsedMessage.Payload = p
		}
	case types.TypeCreateDatabase:
		var p types.CreateDatabase
		err = json.Unmarshal(*m["payload"], &p)
		if err == nil {
			parsedMessage.Payload = p
		}
	case types.TypeOpenDatabase:
		var p types.OpenDatabase
		err = json.Unmarshal(*m["payload"], &p)
		if err == nil {
			parsedMessage.Payload = p
		}
	case types.TypeCreateCollection:
		var p types.CreateCollection
		err = json.Unmarshal(*m["payload"], &p)
		if err == nil {
			parsedMessage.Payload = p
		}
	case types.TypeDropCollection:
		var p types.DropCollection
		err = json.Unmarshal(*m["payload"], &p)
		if err == nil {
			parsedMessage.Payload = p
		}
	default:
		err = json.Unmarshal(*m["payload"], &parsedMessage.Payload)
	}
	return &parsedMessage, err
}
