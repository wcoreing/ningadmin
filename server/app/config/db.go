package config

import "wnadmin/app/provider/db/model"

type DBConfig struct {
	Mysql map[string]model.Mysql `json:"mysql" yaml:"mysql" mapstructure:"mysql"`
	Redis map[string]model.Redis `json:"redis" yaml:"redis" mapstructure:"redis"`
}

var DB DBConfig
