package system

import "errors"

var (
	// ErrConfigKeyExists 配置键已存在
	ErrConfigKeyExists = errors.New("配置键已存在")

	// ErrUserNameExists 用户名已存在
	ErrUserNameExists = errors.New("用户名已存在")

	// ErrRoleKeyExists 角色标识已存在
	ErrRoleKeyExists = errors.New("角色标识已存在")

	// ErrRoleNotFound 角色不存在
	ErrRoleNotFound = errors.New("角色不存在")

	// ErrRoleHasUsers 角色已分配用户
	ErrRoleHasUsers = errors.New("角色已分配用户，不能删除")
)
