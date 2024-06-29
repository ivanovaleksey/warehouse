package config

import (
	"flag"
	"log"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

var cfgFile = flag.String("config", "", "Path to config file")

type Config interface {
	GetConfig(path string, obj any) error
}

func NewConfig() Config {
	return NewConfigFromFile("config/config.local.yaml")
}

func NewConfigFromFile(fl string) Config {
	flag.Parse()
	if *cfgFile == "" {
		*cfgFile = fl
	}
	k := koanf.New(".")
	if err := k.Load(file.Provider(*cfgFile), yaml.Parser()); err != nil {
		log.Fatalf("error loading config file: %v", err)
	}
	return &Loader{k: k}
}

type Loader struct {
	k *koanf.Koanf
}

func (l *Loader) GetConfig(path string, obj any) error {
	return l.k.Unmarshal(path, obj)
}
