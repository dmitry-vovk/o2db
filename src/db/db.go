package db

import (
	"config"
	"errors"
	"os"
	"server/message"
)

type Database struct {

}

var (
	databases = make(map[string]*Database)
)

func CreateDatabase(p message.Payload) error {
	if p["name"] == "" {
		return errors.New("Cannot create database with empty name")
	}
	var dbPath = config.Config.DataDir + string(os.PathSeparator) + p["name"]
	if _, err := os.Stat(dbPath); !os.IsNotExist(err) {
		return errors.New("Database already exists")
	}
	return os.Mkdir(dbPath, 0700)
}

func OpenDatabase(p message.Payload) (*Database, error) {
	if p["name"] == "" {
		return nil, errors.New("Database name cannot be empty")
	}
	if db, has := databases[p["name"]]; has {
		return db, nil
	}
	// TODO Open DB here and put pointer to databases map
	return nil, nil
}
