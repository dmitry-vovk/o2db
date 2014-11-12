package db

import (
	"config"
	"errors"
	"os"
	. "types"
)

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
