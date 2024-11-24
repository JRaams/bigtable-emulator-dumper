package config

import (
	"log"
	"sync"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Address          string `envconfig:"API_ADDRESS"       required:"true"`
	BigtableInstance string `envconfig:"BIGTABLE_INSTANCE" required:"true"`
	BigtableProject  string `envconfig:"BIGTABLE_PROJECT"  required:"true"`
}

var config *Config
var once sync.Once

func Load() *Config {
	once.Do(func() {
		load()
	})

	return config
}

func load() {
	config = new(Config)
	err := envconfig.Process("", config)
	if err != nil {
		log.Fatal(err)
	}
}
