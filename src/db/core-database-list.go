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
			dirs = append(dirs, strings.Replace(dir, config.Config.DataDir+string(os.PathSeparator), "", 1))
		}
	}
	response, err := json.Marshal(dirs)
	if err != nil {
		return "", err
	}
	return string(response), nil
}
