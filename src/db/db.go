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

func (this *Database) CreateCollection(p CreateCollection) error {
	var collectionPath = this.DataDir + string(os.PathSeparator) + getHash(p.Name)
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
	return ioutil.WriteFile(collectionPath+string(os.PathSeparator)+"schema.json", schema, os.FileMode(0600))
}

func (this *Database) DropCollection(p DropCollection) error {
	var collectionPath = this.DataDir + string(os.PathSeparator) + getHash(p.Name)
	if _, err := os.Stat(collectionPath); os.IsNotExist(err) {
		return errors.New("Collection does not exist")
	}
	return os.RemoveAll(collectionPath)
}

func getHash(s string) string {
	return hex.EncodeToString(sha1.New().Sum([]byte(s)))
}
