package system

import "errors"

var (
	// ErrConfigKeyExists 配置键已存在
	ErrConfigKeyExists = errors.New("配置键已存在")

	// ErrUserNameExists 用户名已存在
	ErrUserNameExists = errors.New("用户名已存在")
)
