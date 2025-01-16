package model

import "net/url"

type Mysql struct {
	GeneralDB `yaml:",inline" mapstructure:",squash"`
}

func (m Mysql) Dsn() string {
	return m.Username + ":" + url.QueryEscape(m.Password) + "@tcp(" + m.Path + ":" + m.Port + ")/" + m.Dbname + "?" + m.Config
}

func (m Mysql) GetLogMode() string {
	return m.LogMode
}
