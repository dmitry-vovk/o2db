package config

import "testing"

func TestConfig(t *testing.T) {
	if err := Read("good-config.json"); err != nil {
		t.Errorf("Failed to read good config: %s", err)
	}
	if err := Read("bad-config.json"); err == nil {
		t.Error("Succeeded to read bad config")
	}
}
