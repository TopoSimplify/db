package db

import (
	"github.com/naoina/toml"
	"github.com/intdxdt/fileutil"
)

type Config struct {
	Ignore         bool `toml:"ignore"`
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

func NewConfig(fileName string) Config {
	var cfg Config
	var txt, err = fileutil.ReadAllOfFile(fileName)
	if err != nil {
		panic(err)
	}
	err = toml.Unmarshal([]byte(txt), &cfg)
	if err != nil {
		panic(err)
	}
	return cfg
}
