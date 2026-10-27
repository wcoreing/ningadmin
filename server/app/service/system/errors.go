package system

import "errors"

var (
	// 通用错误
	ErrInvalidParams = errors.New("无效的参数")
	ErrNotFound      = errors.New("记录不存在")
	ErrSystemError   = errors.New("系统错误")

	// 用户相关错误
	ErrUserNotFound     = errors.New("用户不存在")
	ErrPasswordInvalid  = errors.New("密码错误")
	ErrUserDisabled     = errors.New("用户已禁用")
	ErrUserExists       = errors.New("用户已存在")
	ErrUserNameExists   = errors.New("用户名已存在")
	ErrOldPasswordWrong = errors.New("原密码错误")

	// 角色相关错误
	ErrRoleNotFound   = errors.New("角色不存在")
	ErrRoleExists     = errors.New("角色已存在")
	ErrRoleKeyExists  = errors.New("角色标识已存在")
	ErrRoleHasUsers   = errors.New("角色已分配用户，不能删除")
	ErrRoleNameExists = errors.New("角色名称已存在")

	// 菜单相关错误
	ErrMenuNotFound     = errors.New("菜单不存在")
	ErrPermissionExists = errors.New("权限标识已存在")
	ErrHasChildren      = errors.New("存在子菜单，不能删除")
	ErrMenuHasRoles     = errors.New("菜单已分配角色，不能删除")

	// 配置相关错误
	ErrConfigKeyExists = errors.New("配置键已存在")
)
