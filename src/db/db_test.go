package db

import (
	. "dbtest"
	"os"
	"testing"
	"types"
)

func TestCreateCollection(t *testing.T) {
	os.RemoveAll(TestDataDir)
	err := os.MkdirAll(TestDataDir, 0700)
	if err != nil {
		t.Fatalf("Error creating temp dir: %s", err)
	}
	db := &Database{}
	db.DataDir = TestDataDir
	db.Collections = make(map[string]*Collection)
	fields := make(map[string]types.Field)
	fields["f1"] = types.Field{
		Type:  "int",
		Index: "",
	}
	var cc types.CreateCollection
	cc.Name = TestCollectionName
	cc.Fields = fields
	err = db.CreateCollection(cc)
	if err != nil {
		t.Errorf("Error creating collection: %s", err)
	}
	if _, ok := db.Collections[TestCollectionHash]; !ok {
		t.Error("Collection not in database")
	}
	if db.Collections[TestCollectionHash].Name != TestCollectionName {
		t.Error("Collection name is incorrect")
	}
	if len(db.Collections[TestCollectionHash].Indices) != len(fields) {
		t.Error("Collection index count is incorrect")
	}
	if db.Collections[TestCollectionHash].DataFile == nil {
		t.Error("Collection data file failed to initialize")
	}
	if db.Collections[TestCollectionHash].DataFile.FileName != TestDataFileName {
		t.Error("Collection data file name is incorrect")
	}
	if db.Collections[TestCollectionHash].IndexFile == nil {
		t.Error("Collection index file failed to initialize")
	}
	if db.Collections[TestCollectionHash].IndexPointerFile != TestIndexFileName {
		t.Error("Collection index file name is incorrect")
	}
	//t.Logf("%# v", pretty.Formatter(db.Collections[testCollectionHash]))
}
