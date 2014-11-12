package db

import (
	"config"
	"errors"
	"os"
	. "types"
)

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
