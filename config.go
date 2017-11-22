package db


import (
	"github.com/naoina/toml"
	"github.com/intdxdt/fileutil"
)

type Config struct {
	Host     string
	Password string
	Database string
}

func ReadDBConfig(fname string) Config{
	var cfg Config
	var txt, err = fileutil.ReadAllOfFile(fname)
	if err != nil {
		panic(err)
	}
	toml.Unmarshal([]byte(txt), &cfg)
	return cfg
}
