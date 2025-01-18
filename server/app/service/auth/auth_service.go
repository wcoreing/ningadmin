package auth

import (
	"errors"
	"wnadmin/app/common/utils"
	"wnadmin/app/model/do/SystemDB"
	"wnadmin/app/model/dto"
	"wnadmin/app/service/system"
)

// AuthService 认证服务
type AuthService struct {
	menuService *system.MenuService
}

// NewAuthService 创建认证服务
func NewAuthService() *AuthService {
	return &AuthService{
		menuService: &system.MenuService{},
	}
}

// Login 用户登录
func (s *AuthService) Login(req *dto.LoginReq) (*dto.LoginResp, error) {
	// 1. 根据用户名查询用户
	user := &SystemDB.UserDO{}
	err := user.Query().Where("user_name = ? AND del_flag = ?", req.UserName, "0").First(user).Error
	if err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	// 2. 检查用户状态
	if user.Status == "1" {
		return nil, errors.New("账号已被停用")
	}

	// 3. 验证密码
	if !utils.ComparePassword(user.Password, req.Password) {
		return nil, errors.New("用户名或密码错误")
	}

	// 4. 生成令牌
	token, err := utils.GenerateToken(user.ID, user.UserName)
	if err != nil {
		return nil, errors.New("生成令牌失败")
	}

	// 5. 构建响应数据
	userInfo := dto.UserInfo{
		UserID:      user.ID,
		UserName:    user.UserName,
		NickName:    user.NickName,
		Email:       user.Email,
		Phonenumber: user.Phonenumber,
		Sex:         user.Sex,
		Avatar:      user.Avatar,
		Status:      user.Status,
		CreateTime:  user.CreateTime.Format("2006-01-02 15:04:05"),
	}

	return &dto.LoginResp{
		Token:    token,
		UserInfo: userInfo,
	}, nil
}

// GetUserInfo 获取用户信息
func (s *AuthService) GetUserInfo(userID int64) (*dto.UserInfo, error) {
	// 查询用户信息
	user := &SystemDB.UserDO{ID: userID}
	if err := user.First(); err != nil {
		return nil, err
	}

	// 构建用户信息
	userInfo := &dto.UserInfo{
		UserID:      user.ID,
		UserName:    user.UserName,
		NickName:    user.NickName,
		Email:       user.Email,
		Phonenumber: user.Phonenumber,
		Sex:         user.Sex,
		Avatar:      user.Avatar,
		Status:      user.Status,
		CreateTime:  user.CreateTime.Format("2006-01-02 15:04:05"),
	}

	return userInfo, nil
}

// GetUserMenuTree 获取用户的权限菜单树
func (s *AuthService) GetUserMenuTree(userID int64) ([]dto.MenuTree, error) {
	// 获取用户角色
	user := &SystemDB.UserDO{ID: userID}
	role, err := user.GetUserRole()
	if err != nil {
		return nil, err
	}

	// 获取角色的权限菜单树
	return s.menuService.GetRoleMenuTree(role.RoleID)
}

// Logout 用户登出
func (s *AuthService) Logout() error {
	// 这里可以添加一些清理工作，比如清除缓存中的用户信息等
	return nil
}
