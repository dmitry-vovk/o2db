// Database definition and methods to work with database collections
package db

import (
	"encoding/json"
	"errors"
	_ "github.com/kr/pretty"
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
		IndexPointerFile: collectionPath + string(os.PathSeparator) + ObjectIndexFileName,
		ObjectIndexFlush: make(chan (bool), 100),
	}
	// TODO test how indices are created
	d.Collections[collectionNameHash].Indices = make(map[string]ObjectIndex)
	for k, _ := range p.Fields {
		d.Collections[collectionNameHash].Indices[k] = make(map[Hash][]int)
	}
	basePath := collectionPath + string(os.PathSeparator)
	err = ioutil.WriteFile(basePath+"schema.json", schema, os.FileMode(0600))
	if err != nil {
		return err
	}
	d.Collections[collectionNameHash].DataFile = &DbFile{
		FileName: basePath + DataFileName,
	}
	d.Collections[collectionNameHash].DataFile.Touch()
	d.Collections[collectionNameHash].IndexFile = make(map[string]*DbFile)
	d.Collections[collectionNameHash].IndexFile["primary"] = &DbFile{
		FileName: basePath + PrimaryIndexFileName,
	}
	d.Collections[collectionNameHash].IndexFile["primary"].Touch()
	go d.Collections[collectionNameHash].objectIndexFlusher()
	d.Collections[collectionNameHash].ObjectIndexFlush <- true
	//ErrorLog.Printf("%# v", pretty.Formatter(d.Collections[collectionNameHash]))
	return nil
}

// Deletes collection
func (d *Database) DropCollection(p DropCollection) error {
	var hashedName = hash(p.Name)
	var collectionPath = d.DataDir + string(os.PathSeparator) + hashedName
	// Check if collection exists
	if _, err := os.Stat(collectionPath); os.IsNotExist(err) {
		return errors.New("Collection does not exist")
	}
	// Close all related files
	d.Collections[hashedName].DataFile.Close()
	for _, f := range d.Collections[hashedName].IndexFile {
		if f != nil {
			f.Close()
		}
	}
	// Delete collection reference from database
	delete(d.Collections, hashedName)
	// Delete all related files
	return os.RemoveAll(collectionPath)
}
