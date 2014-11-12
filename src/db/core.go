// The file contains types for package and
// core database object with methods to handle databases
package db

import (
	"config"
	"encoding/json"
	"errors"
	_ "github.com/kr/pretty"
	"os"
	"path/filepath"
	"strings"
	. "types"
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
