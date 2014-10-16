// The file contains types for package and
// core database object with methods to handle databases
package db

import (
	"config"
	"encoding/json"
	"errors"
	"logger"
	"os"
	"path/filepath"
	"strings"
	. "types"
)

const (
	dataFileName         = "objects.data"
	primaryIndexFileName = "primary.index"
	objectIndexFileName  = "objects.index"
)

type Package struct {
	Container *Container
	Client    *Client
	RespChan  chan Response
}

type DbCore struct {
	databases map[string]*Database
	Input     chan *Package
}

var (
	Core DbCore
)

// Goroutine that handles queries asynchronously
func (с *DbCore) Processor() {
	с.databases = make(map[string]*Database)
	for {
		pkg := <-с.Input
		pkg.RespChan <- с.ProcessQuery(pkg.Client, pkg.Container)
	}
}

// Creates new database
func (с *DbCore) CreateDatabase(p CreateDatabase) error {
	if p.Name == "" {
		return errors.New("Cannot create database with empty name")
	}
	var dbPath = config.Config.DataDir + string(os.PathSeparator) + p.Name
	if _, err := os.Stat(dbPath); !os.IsNotExist(err) {
		return errors.New("Database already exists")
	}
	if err := os.Mkdir(dbPath, os.FileMode(0700)); err != nil {
		return err
	}
	с.databases[p.Name] = &Database{
		DataDir:     p.Name,
		Collections: make(map[string]*Collection),
	}
	return nil
}

// Deletes existing database
func (с *DbCore) DropDatabase(p DropDatabase) error {
	if p.Name == "" {
		return errors.New("Database name cannot be empty")
	}
	var dbPath = config.Config.DataDir + string(os.PathSeparator) + p.Name
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return errors.New("Database does not exists")
	}
	if _, has := с.databases[p.Name]; has {
		delete(с.databases, p.Name)
	}
	return os.RemoveAll(dbPath)
}

// Returns the list of existing databases
func (с *DbCore) ListDatabases(p ListDatabases) (string, error) {
	if p.Mask == "" {
		return "", errors.New("Mask cannot be empty")
	}
	files, err := filepath.Glob(config.Config.DataDir + string(os.PathSeparator) + p.Mask)
	if err != nil {
		return "", err
	}
	var dirs []string
	for _, dir := range files {
		fi, err := os.Stat(dir)
		if err == nil && fi.IsDir() {
			// TODO add more sophisticated check for database presence besides being a directory
			dirs = append(dirs, strings.Replace(dir, config.Config.DataDir+string(os.PathSeparator), "", 1))
		}
	}
	response, err := json.Marshal(dirs)
	if err != nil {
		return "", err
	}
	return string(response), nil
}

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
			primaryIndexFile, err := os.Stat(collectionDir + primaryIndexFileName)
			if err != nil {
				logger.ErrorLog.Printf("No primary index file found in %s", collectionDir)
				continue
			}
			// create collection object
			d.Collections[collectionHashedName] = &Collection{
				Name:    collectionHashedName,
				Objects: make(map[int]ObjectPointer),
				Indices: make(map[string]ObjectIndex),
				DataFile: &DbFile{
					FileName: collectionDir + dataFileName,
				},
				IndexFile:        make(map[string]*DbFile),
				IndexPointerFile: collectionDir + objectIndexFileName,
				ObjectIndexFlush: make(chan (bool), 100),
			}
			// Add primary index
			d.Collections[collectionHashedName].IndexFile["primary"] = &DbFile{
				FileName: primaryIndexFile.Name(),
			}
			go d.Collections[collectionHashedName].objectIndexFlusher()
		}
	}
	return nil
}
