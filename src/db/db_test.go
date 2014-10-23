package db

import (
	"github.com/kr/pretty"
	"os"
	"testing"
	"types"
)

const (
	testDataDir        = "/tmp/dbtest"
	testCollectionName = "new / collection"
	testCollectionHash = "e2b34345c73ebfae7ae71c940dca37746a987ca2"
	testDataFileName   = testDataDir + "/" + testCollectionHash + "/" + dataFileName
	testIndexFileName  = testDataDir + "/" + testCollectionHash + "/" + objectIndexFileName
)

func TestCreateCollection(t *testing.T) {
	os.RemoveAll(testDataDir)
	err := os.MkdirAll(testDataDir, 0700)
	if err != nil {
		t.Fatalf("Error creating temp dir: %s", err)
	}
	db := &Database{}
	db.DataDir = testDataDir
	db.Collections = make(map[string]*Collection)
	fields := make(map[string]types.Field)
	fields["f1"] = types.Field{
		Type:  "int",
		Index: "",
	}
	var cc types.CreateCollection
	cc.Name = testCollectionName
	cc.Fields = fields
	err = db.CreateCollection(cc)
	if err != nil {
		t.Errorf("Error creating collection: %s", err)
	}
	if _, ok := db.Collections[testCollectionHash]; !ok {
		t.Error("Collection not in database")
	}
	if db.Collections[testCollectionHash].Name != testCollectionName {
		t.Error("Collection name is incorrect")
	}
	if len(db.Collections[testCollectionHash].Indices) != len(fields) {
		t.Error("Collection index count is incorrect")
	}
	if db.Collections[testCollectionHash].DataFile == nil {
		t.Error("Collection data file failed to initialize")
	}
	if db.Collections[testCollectionHash].DataFile.FileName != testDataFileName {
		t.Error("Collection data file name is incorrect")
	}
	if db.Collections[testCollectionHash].IndexFile == nil {
		t.Error("Collection index file failed to initialize")
	}
	if db.Collections[testCollectionHash].IndexPointerFile != testIndexFileName {
		t.Error("Collection index file name is incorrect")
	}
	t.Logf("%# v", pretty.Formatter(db.Collections[testCollectionHash]))
}
