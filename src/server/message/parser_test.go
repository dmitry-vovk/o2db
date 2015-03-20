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
		t.Errorf("Failed auth test: %v", err)
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
		t.Errorf("Failed db create test: %v", err)
	} else {
		if msg.Payload.(types.CreateDatabase).Name != "db_0001" {
			t.Error("Failed create db test: wrong db name")
		}
	}
}

func TestOpenDatabase(t *testing.T) {
	if msg, err := Parse([]byte(`{"type":102, "payload": {"name": "db_0001"}}`)); err != nil {
		t.Errorf("Failed db open test: %v", err)
	} else {
		if msg.Payload.(types.OpenDatabase).Name != "db_0001" {
			t.Error("Failed db open test: wrong db name")
		}
	}
}

func TestDropDatabase(t *testing.T) {
	if msg, err := Parse([]byte(`{"type":101, "payload": {"name": "db_0001"}}`)); err != nil {
		t.Errorf("Failed db drop test: %v", err)
	} else {
		if msg.Payload.(types.DropDatabase).Name != "db_0001" {
			t.Error("Failed db drop test: wrong db name")
		}
	}
}

func TestListDatabases(t *testing.T) {
	if msg, err := Parse([]byte(`{"type":103, "payload": {"mask": "db_*"}}`)); err != nil {
		t.Errorf("Failed db list test: %v", err)
	} else {
		if msg.Payload.(types.ListDatabases).Mask != "db_*" {
			t.Error("Failed db list test: wrong mask")
		}
	}
}

func TestCreateCollection(t *testing.T) {
	if msg, err := Parse([]byte(`{"type":200, "payload": {"class": "SampleEntity", "fields": {"id": {"type":"int"}}}}`)); err != nil {
		t.Errorf("Failed collection create test: %v", err)
	} else {
		if msg.Payload.(types.CreateCollection).Name != "SampleEntity" {
			t.Error("Failed collection create test: wrong class name")
		}
		if msg.Payload.(types.CreateCollection).Fields["id"].Type != "int" {
			t.Error("Failed collection create test: wrong type")
		}
	}
}

func TestDropCollection(t *testing.T) {
	if msg, err := Parse([]byte(`{"type":201, "payload": {"class": "SampleEntity"}}`)); err != nil {
		t.Errorf("Failed collection drop test: %v", err)
	} else {
		if msg.Payload.(types.DropCollection).Name != "SampleEntity" {
			t.Error("Failed collection drop test: wrong class name")
		}
	}
}

func TestListCollections(t *testing.T) {
	// nothing to do here, no payload parsing to test
}

func TestObjectWrite(t *testing.T) {
	if msg, err := Parse([]byte(`{"type":301, "payload": {"class": "SampleEntity", "data": {"field": "val"}}}`)); err != nil {
		t.Errorf("Failed object write test: %v", err)
	} else {
		if msg.Payload.(types.WriteObject).Collection != "SampleEntity" {
			t.Error("Failed object write test: wrong collection name")
		}
		if msg.Payload.(types.WriteObject).Data["field"] != "val" {
			t.Error("Failed object write test: wrong field value")
		}
	}
}

func TestObjectGet(t *testing.T) {
	if msg, err := Parse([]byte(`{"type":300, "payload": {"class": "SampleEntity", "data": {"field": "val"}}}`)); err != nil {
		t.Errorf("Failed object read test: %v", err)
	} else {
		if msg.Payload.(types.ReadObject).Collection != "SampleEntity" {
			t.Error("Failed object read test: wrong collection name")
		}
		if msg.Payload.(types.ReadObject).Fields["field"] != "val" {
			t.Error("Failed object read test: wrong field value")
		}
	}
}

func TestObjectVersions(t *testing.T) {
	if msg, err := Parse([]byte(`{"type":303, "payload": {"class": "SampleEntity", "id": 5}}`)); err != nil {
		t.Errorf("Failed get object versions test: %v", err)
	} else {
		if msg.Payload.(types.GetObjectVersions).Collection != "SampleEntity" {
			t.Error("Failed get object versions test: wrong collection name")
		}
		if msg.Payload.(types.GetObjectVersions).Id != 5 {
			t.Error("Failed get object versions test: wrong field value")
		}
	}
}

func TestGetObjectDiff(t *testing.T) {
	if msg, err := Parse([]byte(`{"type":304, "payload": {"class": "SampleEntity", "id": 5, "from": 6, "to": 8}}`)); err != nil {
		t.Errorf("Failed get object versions test: %v", err)
	} else {
		if msg.Payload.(types.GetObjectDiff).Collection != "SampleEntity" {
			t.Error("Failed get object versions test: wrong collection name")
		}
		if msg.Payload.(types.GetObjectDiff).Id != 5 {
			t.Error("Failed get object versions test: wrong id value")
		}
		if msg.Payload.(types.GetObjectDiff).From != 6 {
			t.Error("Failed get object versions test: wrong from value")
		}
		if msg.Payload.(types.GetObjectDiff).To != 8 {
			t.Error("Failed get object versions test: wrong to value")
		}
	}
}

func TestSelectObjects(t *testing.T) {
	if msg, err := Parse([]byte(`{"type":305, "payload": {"class": "SampleEntity", "query": {"field": "val"}}}`)); err != nil {
		t.Errorf("Failed select objects test: %v", err)
	} else {
		if msg.Payload.(types.SelectObjects).Collection != "SampleEntity" {
			t.Error("Failed select objects test: wrong collection name")
		}
		if msg.Payload.(types.SelectObjects).Query["field"] != "val" {
			t.Errorf("Failed select object test: wrong field value")
		}
	}
}

func TestSubscribe(t *testing.T) {
	if msg, err := Parse([]byte(`{"type":400, "payload": {"class": "SampleEntity", "key": "key_key"}}`)); err != nil {
		t.Errorf("Failed subscribe test: %v", err)
	} else {
		if msg.Payload.(types.Subscribe).Collection != "SampleEntity" {
			t.Error("Failed subscribe test: wrong collection name")
		}
		if msg.Payload.(types.Subscribe).Key != "key_key" {
			t.Errorf("Failed subscribe test: wrong key value")
		}
	}
}

func TestAddSubscription(t *testing.T) {
	if msg, err := Parse([]byte(`{"type":401, "payload": {"class": "SampleEntity", "key": "key_key", "query": {"field": "val"}}}`)); err != nil {
		t.Errorf("Failed add subscription test: %v", err)
	} else {
		if msg.Payload.(types.AddSubscription).Collection != "SampleEntity" {
			t.Error("Failed add subscription test: wrong collection name")
		}
		if msg.Payload.(types.AddSubscription).Key != "key_key" {
			t.Errorf("Failed add subscription test: wrong key value")
		}
		if msg.Payload.(types.AddSubscription).Query["field"] != "val" {
			t.Errorf("Failed add subscription test: wrong field value")
		}
	}
}

func TestListSubscriptions(t *testing.T) {
	if msg, err := Parse([]byte(`{"type":403, "payload": {"classes": ["SampleEntity"]}}`)); err != nil {
		t.Errorf("Failed list subscriptions test: %v", err)
	} else {
		if msg.Payload.(types.ListSubscriptions).Collections[0] != "SampleEntity" {
			t.Error("Failed list subscriptions test: wrong collection name")
		}
	}
}

func TestCancelSubscriptions(t *testing.T) {
	if msg, err := Parse([]byte(`{"type":402, "payload": {"class": "SampleEntity", "key": "key_key"}}`)); err != nil {
		t.Errorf("Failed subscribe test: %v", err)
	} else {
		if msg.Payload.(types.CancelSubscription).Collection != "SampleEntity" {
			t.Error("Failed subscribe test: wrong collection name")
		}
		if msg.Payload.(types.CancelSubscription).Key != "key_key" {
			t.Errorf("Failed subscribe test: wrong key value")
		}
	}
}
