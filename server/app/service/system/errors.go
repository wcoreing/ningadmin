package system

import "errors"

var (
	// ErrConfigKeyExists 配置键已存在
	ErrConfigKeyExists = errors.New("配置键已存在")
)
