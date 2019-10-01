package config

import (
	"github.com/kelseyhightower/envconfig"
	"log"
)

type PathConfig struct {
	SchemasDir   string `envconfig:"SCHEMA_DIR" default:"schemas/schemas"`
	SchemasURL   string `envconfig:"SCHEMA_URL" default:"/"`
	ExamplesDir  string `envconfig:"EXAMPLE_DIR" default:"schemas/examples"`
	ExamplesURL  string `envconfig:"EXAMPLE_URL" default:"/examples/"`
	ServerAddr   string `envconfig:"SERVER_ADDR" default:":8080"`
	WarnsAllowed int    `envconfig:"WARNS_ALLOWED" default:"0"`
	LogLevel     string `envconfig:"LOG_LEVEL" default:"INFO"`
}

func Parse() PathConfig {
	var r PathConfig
	err := envconfig.Process("", &r)
	if err != nil {
		log.Fatal(err.Error())
	}

	return r
}
