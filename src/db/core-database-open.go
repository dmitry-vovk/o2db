package db

import (
	"config"
	"errors"
	"github.com/kr/pretty"
	"logger"
	"os"
	"path/filepath"
	"strings"
	. "types"
)

// Open existing database
func (с *DbCore) OpenDatabase(p OpenDatabase) (string, error) {
	if p.Name == "" {
		return "", errors.New("Database name cannot be empty")
	}
	if _, has := с.databases[p.Name]; has {
		return p.Name, nil
	}
	err := с.openDatabase(p.Name)
	if err != nil {
		return "", err
	}
	return p.Name, nil
}

// Low level database opener
func (с *DbCore) openDatabase(dbName string) error {
	var dbPath = config.Config.DataDir + string(os.PathSeparator) + dbName
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return errors.New("Database does not exists")
	}
	с.databases[dbName] = &Database{
		DataDir:     dbPath,
		Collections: make(map[string]*Collection),
	}
	return с.populateCollections(с.databases[dbName])
}

// Scans directories under database data directory
// and creates collection structures from found files
func (c *DbCore) populateCollections(d *Database) error {
	files, err := filepath.Glob(d.DataDir + string(os.PathSeparator) + "*")
	if err != nil {
		return err
	}
	// Iterate through all directories under database data dir
	for _, dir := range files {
		fi, err := os.Stat(dir)
		if err != nil {
			return err
		}
		if fi.IsDir() {
			// get directory name
			collectionHashedName := strings.Replace(dir, d.DataDir+string(os.PathSeparator), "", 1)
			// full path to collection directory
			collectionDir := d.DataDir + string(os.PathSeparator) + collectionHashedName + string(os.PathSeparator)
			// check if primary index file is present
			primaryIndexFile, err := os.Stat(collectionDir + PrimaryIndexFileName)
			if err != nil {
				logger.ErrorLog.Printf("No primary index file found in %s", collectionDir)
				continue
			}
			// create collection object
			d.Collections[collectionHashedName] = &Collection{
				Name:    collectionHashedName,
				Objects: make(map[int]ObjectPointer),
				Indices: make(map[string]FieldIndex),
				DataFile: &DbFile{
					FileName: collectionDir + DataFileName,
				},
				IndexFile:        make(map[string]*DbFile),
				IndexPointerFile: collectionDir + ObjectIndexFileName,
				ObjectIndexFlush: make(chan (bool), 100),
				BaseDir:          collectionDir,
			}
			// Open object storage
			d.Collections[collectionHashedName].DataFile.Open()
			// Add primary index
			d.Collections[collectionHashedName].IndexFile["primary"] = &DbFile{
				FileName: primaryIndexFile.Name(),
			}
			// Run goroutine that will flush objects index
			go d.Collections[collectionHashedName].objectIndexFlusher()
			// Try to read existing objects index (it may not exist for empty collection)
			if err := d.Collections[collectionHashedName].readObjectIndex(); err != nil {
				logger.ErrorLog.Printf("Could not read object index: %s", err)
				return err
			}
			if err := d.Collections[collectionHashedName].ReadSchema(); err != nil {
				logger.ErrorLog.Printf("Could not read schema: %s", err)
				return err
			}
			for k, v := range d.Collections[collectionHashedName].Schema {
				indexFileName := d.Collections[collectionHashedName].BaseDir + hash(k) + ".index"
				switch v.Type {
				case "string":
					d.Collections[collectionHashedName].Indices[k] = NewStringIndex(indexFileName)
				case "int":
					d.Collections[collectionHashedName].Indices[k] = NewIntIndex(indexFileName)
				case "float":
					d.Collections[collectionHashedName].Indices[k] = NewFloatIndex(indexFileName)
				}
			}
		}
	}
	logger.ErrorLog.Printf("%# v", pretty.Formatter(c.databases))
	return nil
}
