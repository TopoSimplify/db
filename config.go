package db

import (
	"github.com/naoina/toml"
)

type Config struct {
	Ignore         bool   `toml:"ignore"`
	Host           string `toml:"host"`
	Password       string `toml:"password"`
	Database       string `toml:"database"`
	User           string `toml:"user"`
	Table          string `toml:"table"`
	GeometryColumn string `toml:"geometrycolumn"`
	IdColumn       string `toml:"idcolumn"`
}

func (cfg Config) Clone() Config {
	return cfg
}

func NewConfig(txtToml string) Config {
	var cfg Config
	var err = toml.Unmarshal([]byte(txtToml), &cfg)
	if err != nil {
		panic(err)
	}
	return cfg
}
