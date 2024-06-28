package config

import (
	"github.com/knadh/koanf/v2"
)

type Config interface {
	GetConfig(path string, obj any) error
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
