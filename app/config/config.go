package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	DB Mysql `json:"db" yaml:"db"`
}

type Mysql struct {
	User     string `json:"user" yaml:"user"`
	Password string `json:"password" yaml:"password"`
	Host     string `json:"host" yaml:"host"`
	Port     int    `json:"port" yaml:"port"`
	Dbname   string `json:"dbname" yaml:"dbname"`
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
