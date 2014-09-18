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

type Package struct {
	Container *Container
	Client    *ClientType
	RespChan  chan Response
}

type DbCore struct {
	databases map[string]*Database
	Input     chan *Package
}

var (
	Core DbCore
)

func (this *DbCore) Processor() {
	this.databases = make(map[string]*Database)
	for {
		pkg := <-this.Input
		pkg.RespChan <- this.ProcessQuery(pkg.Client, pkg.Container)
	}
}

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
