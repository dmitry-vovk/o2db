// Database definition and methods to work with database collections
package db

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/ioutil"
	. "logger"
	"os"
	. "types"
)

type Database struct {
	DataDir     string
	Collections map[string]*Collection
}

// Creates new empty collection
func (d *Database) CreateCollection(p CreateCollection) error {
	var collectionPath = d.DataDir + string(os.PathSeparator) + hash(p.Name)
	if _, err := os.Stat(collectionPath); !os.IsNotExist(err) {
		return errors.New("Collection already exists")
	}
	err := os.Mkdir(collectionPath, os.FileMode(0700))
	if err != nil {
		return err
	}
	DebugLog.Printf("Creating collection %s in %s", p.Name, collectionPath)
	var schema []byte
	schema, err = json.MarshalIndent(p.Fields, "", "  ")
	if err != nil {
		return err
	}
	collectionNameHash := hash(p.Name)
	d.Collections[collectionNameHash] = &Collection{
		Name:             p.Name,
		Objects:          make(map[int]ObjectPointer),
		IndexPointerFile: &DbFile{},
		ObjectIndexFlush: make(chan (bool), 100),
	}
	d.Collections[collectionNameHash].IndexPointerFile.Touch()
	basePath := collectionPath + string(os.PathSeparator)
	err = ioutil.WriteFile(basePath+"schema.json", schema, os.FileMode(0600))
	if err != nil {
		return err
	}
	d.Collections[collectionNameHash].DataFile = &DbFile{
		FileName: basePath + dataFileName,
	}
	d.Collections[collectionNameHash].DataFile.Touch()
	d.Collections[collectionNameHash].IndexFile = make(map[string]*DbFile)
	d.Collections[collectionNameHash].IndexFile["primary"] = &DbFile{
		FileName: basePath + primaryIndexFileName,
	}
	d.Collections[collectionNameHash].IndexFile["primary"].Touch()
	go d.Collections[collectionNameHash].objectIndexFlusher()
	return nil
}

// Deletes collection
func (d *Database) DropCollection(p DropCollection) error {
	var collectionPath = d.DataDir + string(os.PathSeparator) + hash(p.Name)
	if _, err := os.Stat(collectionPath); os.IsNotExist(err) {
		return errors.New("Collection does not exist")
	}
	return os.RemoveAll(collectionPath)
}

// Shorthand to get SHA1 string
func hash(s string) string {
	sh := sha1.New()
	sh.Write([]byte(s))
	return hex.EncodeToString(sh.Sum(nil))
}
