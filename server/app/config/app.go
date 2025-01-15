package config

type AppConfig struct {
	ENV  string `yaml:"ENV"`
	Host string `yaml:"Host"`
	Port int    `yaml:"Port"`
	Name string `yaml:"Name"`
}

var App AppConfig
