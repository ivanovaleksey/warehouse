package config

import (
	"log"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type Config interface {
	GetConfig(path string, obj any) error
}

func NewConfig() Config {
	return NewConfigFromFile("config/config.local.yaml")
}

func NewConfigFromFile(f string) Config {
	k := koanf.New(".")
	if err := k.Load(file.Provider(f), yaml.Parser()); err != nil {
		log.Fatalf("error loading config file: %v", err)
	}
	return NewLoader(k)
}

func NewLoader(k *koanf.Koanf) *Loader {
	return &Loader{k: k}
}

type Loader struct {
	k *koanf.Koanf
}

func (l *Loader) GetConfig(path string, obj any) error {
	return l.k.Unmarshal(path, obj)
}
