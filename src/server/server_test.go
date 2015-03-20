package server

import (
	"testing"
	"config"
)

func TestServerNew(t *testing.T) {
	config := &config.ConfigType{}
	s := CreateNew(config)
	if s.Config != config {
		t.Error("Failed config check")
	}
}
