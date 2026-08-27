package config

import "testing"

func TestLoadDefaults(t *testing.T) {
	c := Load()
	if c.Port == 0 || c.DBPath == "" {
		t.Fatal(c)
	}
}
