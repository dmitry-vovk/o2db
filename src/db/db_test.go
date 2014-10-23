package db

import (
	"testing"
	"types"
)

func TestCreateCollection(t *testing.T) {
	db := &Database{}
	fields := make(map[string]types.Field)
	fields["f1"] = types.Field{
		Type:  "int",
		Index: "",
	}
	var cc types.CreateCollection
	cc.Name = "new / collection"
	cc.Fields = fields
	err := db.CreateCollection(cc)
	if err != nil {
		t.Fatalf("Error creating collection: %s", err)
	}
}
