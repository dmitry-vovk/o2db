package message

import (
	"testing"
	"types"
)

func TestParser(t *testing.T) {
	if _, err := Parse([]byte("illegal message")); err != ErrCannotParseMessage {
		t.Error("Failed parsing illegal message")
	}
	if _, err := Parse([]byte(`{"val":1}`)); err != ErrMissingTypeField {
		t.Error("Failed detecting missing type field")
	}
	if _, err := Parse([]byte(`{"type":0}`)); err != ErrMissingPayloadField {
		t.Error("Failed detecting missing payload field")
	}
	if _, err := Parse([]byte(`{"type":false, "payload":[]}`)); err == nil {
		t.Errorf("Failed to parse type")
	}
	if _, err := Parse([]byte(`{"type":99999, "payload":[]}`)); err != ErrUnsupportedMessageType {
		t.Errorf("Failed to detect unsupported message type")
	}
}

func TestAuth(t *testing.T) {
	if msg, err := Parse([]byte(`{"type":0, "payload": {"name": "admin", "password": "secret"}}`)); err != nil {
		t.Error("Failed auth test: %v", err)
	} else {
		if msg.Payload.(types.Authentication).Name != "admin" {
			t.Error("Failed auth test: wrong username")
		}
		if msg.Payload.(types.Authentication).Password != "secret" {
			t.Error("Failed auth test: wrong password")
		}
	}
}

func TestCreateDatabase(t *testing.T) {
	if msg, err := Parse([]byte(`{"type":100, "payload": {"name": "db_0001"}}`)); err != nil {
		t.Error("Failed db create test: %v", err)
	} else {
		if msg.Payload.(types.CreateDatabase).Name != "db_0001" {
			t.Error("Failed create db test: wrong db name")
		}
	}
}
