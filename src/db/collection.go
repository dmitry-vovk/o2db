// Collection definition and methods to work with collection objects
package db

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"io/ioutil"
	"os"
	. "types"
)

const (
	FIELD_ID      = "id"
	FIELD_VERSION = "__version__"
)

// Object instance
type ObjectVersion struct {
	Offset int // Offset from the beginning of object data file
	Len    int // Number of bytes to read
}

// List of object instances. index is version number.
type ObjectPointer map[int]ObjectVersion

type Collection struct {
	BaseDir          string                // Where collection files are
	Name             string                // Collection/class name
	Objects          map[int]ObjectPointer // Objects. map index is object ID
	Indices          map[string]FieldIndex // collection of indices
	DataFile         *DbFile               // Objects storage
	IndexFile        map[string]*DbFile    // List of indices
	freeSlotOffset   int
	IndexPointerFile string
	ObjectIndexFlush chan (bool)
	Schema           map[string]Field
	Subscriptions    map[string]*Subscription
}

// Returns pointer to the start of unallocated file space
func (c *Collection) getFreeSpaceOffset() int {
	return c.freeSlotOffset
}

// Write collection schema to JSON and GOB files
func (c *Collection) DumpSchema() error {
	// JSON (human readable)
	schema, err := json.MarshalIndent(c.Schema, "", "  ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(c.BaseDir+CollectionSchema+".json", schema, os.FileMode(0600))
	if err != nil {
		return err
	}
	// GOB
	var b bytes.Buffer
	gob.NewEncoder(&b).Encode(c.Schema)
	return ioutil.WriteFile(c.BaseDir+CollectionSchema, b.Bytes(), os.FileMode(0600))
}

// Read schema from GOB file
func (c *Collection) ReadSchema() error {
	handler, err := os.Open(c.BaseDir + CollectionSchema)
	if err != nil {
		return err
	}
	defer handler.Close()
	dec := gob.NewDecoder(handler)
	err = dec.Decode(&c.Schema)
	if err != nil {
		return err
	}
	return nil
}
