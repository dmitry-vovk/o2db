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
func (this *Database) CreateCollection(p CreateCollection) error {
	var collectionPath = this.DataDir + string(os.PathSeparator) + hash(p.Name)
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
	this.Collections[p.Name] = &Collection{
		Name: p.Name,
	}
	basePath := collectionPath + string(os.PathSeparator)
	err = ioutil.WriteFile(basePath+"schema.json", schema, os.FileMode(0600))
	if err != nil {
		return err
	}
	this.Collections[p.Name].DataFile = &DbFile{
		FileName: basePath + dataFileName,
	}
	this.Collections[p.Name].IndexFile = make(map[string]*DbFile)
	this.Collections[p.Name].IndexFile["primary"] = &DbFile{
		FileName: basePath + primaryIndexFileName,
	}
	return nil
}

// Deletes collection
func (this *Database) DropCollection(p DropCollection) error {
	var collectionPath = this.DataDir + string(os.PathSeparator) + hash(p.Name)
	if _, err := os.Stat(collectionPath); os.IsNotExist(err) {
		return errors.New("Collection does not exist")
	}
	return os.RemoveAll(collectionPath)
}

// Shorthand to get SHA1 string
func hash(s string) string {
	return hex.EncodeToString(sha1.New().Sum([]byte(s)))
}
