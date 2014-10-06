// The file contains types for package and
// core database object with methods to handle databases
package db

import (
	"config"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	. "types"
)

const (
	dataFileName         = "objects.data"
	primaryIndexFileName = "primary.index"
	objectIndexFileName  = "object.index"
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
func (this *DbCore) Processor() {
	this.databases = make(map[string]*Database)
	for {
		pkg := <-this.Input
		pkg.RespChan <- this.ProcessQuery(pkg.Client, pkg.Container)
	}
}

// Creates new database
func (this *DbCore) CreateDatabase(p CreateDatabase) error {
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
	this.databases[p.Name] = &Database{
		DataDir:     p.Name,
		Collections: make(map[string]*Collection),
	}
	return nil
}

// Deletes existing database
func (this *DbCore) DropDatabase(p DropDatabase) error {
	if p.Name == "" {
		return errors.New("Database name cannot be empty")
	}
	var dbPath = config.Config.DataDir + string(os.PathSeparator) + p.Name
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return errors.New("Database does not exists")
	}
	if _, has := this.databases[p.Name]; has {
		delete(this.databases, p.Name)
	}
	return os.RemoveAll(dbPath)
}

// Returns the list of existing databases
func (this *DbCore) ListDatabases(p ListDatabases) (string, error) {
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
func (this *DbCore) OpenDatabase(p OpenDatabase) (string, error) {
	if p.Name == "" {
		return "", errors.New("Database name cannot be empty")
	}
	if _, has := this.databases[p.Name]; has {
		return p.Name, nil
	}
	err := this.openDatabase(p.Name)
	if err != nil {
		return "", err
	}
	return p.Name, nil
}

// Low level database opener
func (this *DbCore) openDatabase(dbName string) error {
	var dbPath = config.Config.DataDir + string(os.PathSeparator) + dbName
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return errors.New("Database does not exists")
	}
	this.databases[dbName] = &Database{
		DataDir:     dbPath,
		Collections: make(map[string]*Collection),
	}
	return this.populateCollections(this.databases[dbName])
}

// Scans directories under database data directory
func (this *DbCore) populateCollections(d *Database) error {
	files, err := filepath.Glob(d.DataDir + string(os.PathSeparator) + "*")
	if err != nil {
		return err
	}
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
			// create collection object
			d.Collections[collectionHashedName] = &Collection{
				Name:    collectionHashedName,
				Objects: make(map[int]ObjectPointer),
				Indices: make(map[string]ObjectIndex),
				DataFile: &DbFile{
					FileName: collectionDir + dataFileName,
				},
				IndexFile: make(map[string]*DbFile),
				IndexPointerFile: &DbFile{
					FileName: collectionDir + objectIndexFileName,
				},
			}
			// check if primary index file is present
			if primaryIndexFile, err := os.Stat(collectionDir + primaryIndexFileName); err == nil {
				d.Collections[collectionHashedName].IndexFile["primary"].FileName = primaryIndexFile.Name()
			}
		}
	}
	return nil
}
