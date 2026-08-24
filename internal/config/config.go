package config

import (
	"os"
	"strconv"
)

type Config struct {
	DBPath  string
	Port    int
	DataDir string
}

func Load() Config {
	port, _ := strconv.Atoi(os.Getenv("GYM_PORT"))
	if port == 0 {
		port = 8080
	}
	path := os.Getenv("GYM_DB")
	if path == "" {
		path = "gymrecommend.db"
	}
	dir := os.Getenv("GYM_DATA")
	if dir == "" {
		dir = "data"
	}
	return Config{DBPath: path, Port: port, DataDir: dir}
}
func (c Config) Address() string { return ":" + strconv.Itoa(c.Port) }
