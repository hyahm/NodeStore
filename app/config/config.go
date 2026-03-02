package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	DB string `json:"db" yaml:"db"`
}

var Cfg Config

func InitConfig(file string) error {
	b, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	Cfg = Config{}
	return yaml.Unmarshal(b, &Cfg)
}
