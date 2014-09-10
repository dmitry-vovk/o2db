package config

import (
	"io/ioutil"
	"encoding/json"
)

var Config *ConfigType

func Read(configFile string) error {
	c := &ConfigType{}
	file, err := ioutil.ReadFile(configFile)
	if err == nil {
		err = json.Unmarshal(file, c)
		if err == nil {
			Config = c
		}
	}
	return err
}
