package db

import (
	"errors"
	"logger"
	"os"
	. "types"
)

// Creates new empty collection
func (d *Database) CreateCollection(p CreateCollection) error {
	var collectionPath = d.DataDir + string(os.PathSeparator) + hash(p.Name)
	if _, err := os.Stat(collectionPath); !os.IsNotExist(err) {
		return errors.New("Collection already exists")
	}
	err := os.Mkdir(collectionPath, os.FileMode(0700))
	if err != nil {
		logger.ErrorLog.Printf("Error creating collection dir: %s", err)
		return err
	}
	logger.DebugLog.Printf("Creating collection %s in %s", p.Name, collectionPath)
	collectionNameHash := hash(p.Name)
	d.Collections[collectionNameHash] = &Collection{
		BaseDir:          collectionPath + string(os.PathSeparator),
		Name:             p.Name,
		Objects:          make(map[int]ObjectPointer),
		IndexPointerFile: collectionPath + string(os.PathSeparator) + ObjectIndexFileName,
		ObjectIndexFlush: make(chan (bool), 100),
		Schema:           p.Fields,
	}
	// Save schema
	if err := d.Collections[collectionNameHash].DumpSchema(); err != nil {
		return err
	}
	d.Collections[collectionNameHash].DataFile = &DbFile{
		FileName: d.Collections[collectionNameHash].BaseDir + DataFileName,
	}
	d.Collections[collectionNameHash].DataFile.Touch()
	d.Collections[collectionNameHash].IndexFile = make(map[string]*DbFile)
	d.Collections[collectionNameHash].IndexFile["primary"] = &DbFile{
		FileName: d.Collections[collectionNameHash].BaseDir + PrimaryIndexFileName,
	}
	d.Collections[collectionNameHash].IndexFile["primary"].Touch()
	go d.Collections[collectionNameHash].objectIndexFlusher()
	d.Collections[collectionNameHash].ObjectIndexFlush <- true
	d.Collections[collectionNameHash].CreateIndices(p.Fields)
	return nil
}
