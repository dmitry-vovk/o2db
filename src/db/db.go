package db

import (
	"config"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	. "types"
)

func (this *DbCore) CreateDatabase(p CreateDatabase) error {
	if p.Name == "" {
		return errors.New("Cannot create database with empty name")
	}
	var dbPath = config.Config.DataDir + string(os.PathSeparator) + p.Name
	if _, err := os.Stat(dbPath); !os.IsNotExist(err) {
		return errors.New("Database already exists")
	}
	return os.Mkdir(dbPath, os.FileMode(0700))
}

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

func (this *DbCore) OpenDatabase(p OpenDatabase) (*Database, error) {
	if p.Name == "" {
		return nil, errors.New("Database name cannot be empty")
	}
	if db, has := this.databases[p.Name]; has {
		return db, nil
	}
	err := this.openDatabase(p.Name)
	if err == nil {
		return this.databases[p.Name], nil
	}
	return nil, err
}

func (this *DbCore) openDatabase(dbName string) error {
	var dbPath = config.Config.DataDir + string(os.PathSeparator) + dbName
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return errors.New("Database does not exists")
	}
	this.databases[dbName] = &Database{
		DataDir:     dbPath,
		Collections: make(map[string]Collection),
	}
	return nil
}

func (this *DbCore) CreateCollection(c *ClientType, p CreateCollection) error {
	if c.Db == nil {
		return errors.New("Database not selected")
	}
	var collectionPath = c.Db.DataDir + string(os.PathSeparator) + getHash(p.Name)
	if _, err := os.Stat(collectionPath); !os.IsNotExist(err) {
		return errors.New("Collection already exists")
	}
	err := os.Mkdir(collectionPath, os.FileMode(0700))
	if err != nil {
		return err
	}
	log.Printf("Creating collection %s in %s", p.Name, collectionPath)
	var schema []byte
	schema, err = json.MarshalIndent(p.Fields, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(collectionPath+string(os.PathSeparator)+"schema.json", schema, os.FileMode(0600))
}

func (this *DbCore) DropCollection(c *ClientType, p DropCollection) error {
	if c.Db == nil {
		return errors.New("Database not selected")
	}
	var collectionPath = c.Db.DataDir + string(os.PathSeparator) + getHash(p.Name)
	if _, err := os.Stat(collectionPath); os.IsNotExist(err) {
		return errors.New("Collection does not exist")
	}
	return os.RemoveAll(collectionPath)
}

func getHash(s string) string {
	return hex.EncodeToString(sha1.New().Sum([]byte(s)))
}
