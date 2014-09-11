package db

import (
	"config"
	"errors"
	"os"
	"server/message"
)

func Create(p message.Payload) error {
	if p["name"] == "" {
		return errors.New("Cannot create database with empty name")
	}
	var dbPath = config.Config.DataDir + string(os.PathSeparator) + p["name"]
	if _, err := os.Stat(dbPath); !os.IsNotExist(err) {
		return errors.New("Database already exists")
	}
	return os.Mkdir(dbPath, 0700)
}
