package system

import (
	"time"
	"wnadmin/app/common/utils"
	"wnadmin/app/model/do/SystemDB"
	"wnadmin/app/model/dto"
	"wnadmin/app/provider/db"

	"gorm.io/gorm"
)

// UserService 用户服务
type UserService struct{}

// GetUserList 获取用户列表
func (s *UserService) GetUserList(query *dto.SystemUserQuery) (*dto.SystemUserListResp, error) {
	var users []SystemDB.UserDO
	userDO := &SystemDB.UserDO{}
	db := userDO.Query()

	if query.Keyword != "" {
		db = db.Where("user_name LIKE ? OR nick_name LIKE ?", "%"+query.Keyword+"%", "%"+query.Keyword+"%")
	}
	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}

	// 查询总数
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}

	// 分页查询
	offset := (query.PageNum - 1) * query.PageSize
	if err := db.Offset(offset).Limit(query.PageSize).Order("create_time DESC").Find(&users).Error; err != nil {
		return nil, err
	}

	// 构建响应数据
	resp := &dto.SystemUserListResp{
		Total: total,
		Rows:  make([]dto.SystemUserInfoResp, 0, len(users)),
	}

	// 获取用户角色
	for _, user := range users {
		role, err := user.GetUserRole()
		if err != nil {
			return nil, err
		}

		resp.Rows = append(resp.Rows, dto.SystemUserInfoResp{
			ID:          user.ID,
			UserName:    user.UserName,
			NickName:    user.NickName,
			Email:       user.Email,
			Phonenumber: user.Phonenumber,
			Sex:         user.Sex,
			Status:      user.Status,
			CreateTime:  user.CreateTime.Format("2006-01-02 15:04:05"),
			CreateBy:    user.CreateBy,
			UpdateTime:  user.UpdateTime.Format("2006-01-02 15:04:05"),
			UpdateBy:    user.UpdateBy,
			Role:        role,
		})
	}

	return resp, nil
}

// CreateUser 创建用户
func (s *UserService) CreateUser(req *dto.CreateSystemUserReq, createBy string) error {
	// 检查用户名是否已存在
	userDO := &SystemDB.UserDO{}
	exists, err := userDO.CheckUserNameExists(req.UserName, 0)
	if err != nil {
		return err
	}
	if exists {
		return ErrUserNameExists
	}

	// 加密密码
	encryptedPassword, err := utils.EncryptPassword(req.Password)
	if err != nil {
		return err
	}

	now := time.Now()
	user := &SystemDB.UserDO{
		UserName:    req.UserName,
		Password:    encryptedPassword,
		NickName:    req.NickName,
		Email:       req.Email,
		Phonenumber: req.Phonenumber,
		Sex:         req.Sex,
		Status:      req.Status,
		CreateBy:    createBy,
		CreateTime:  now,
		UpdateTime:  now,
		DelFlag:     "0",
	}

	// 设置默认值
	if user.Status == "" {
		user.Status = "0"
	}
	if user.Sex == "" {
		user.Sex = "0"
	}

	// 开启事务
	return db.WnAdmin().Transaction(func(tx *gorm.DB) error {
		// 创建用户
		if err := user.Create(); err != nil {
			return err
		}

		// 分配角色
		if err := user.UpdateUserRole(req.RoleID); err != nil {
			return err
		}

		return nil
	})
}

// UpdateUser 更新用户
func (s *UserService) UpdateUser(req *dto.UpdateSystemUserReq, updateBy string) error {
	// 检查用户是否存在
	user := &SystemDB.UserDO{ID: req.ID}
	if err := user.First(); err != nil {
		return ErrUserNotFound
	}

	// 检查用户名是否已存在（排除自身）
	if req.UserName != user.UserName {
		exists, err := user.CheckUserNameExists(req.UserName, req.ID)
		if err != nil {
			return err
		}
		if exists {
			return ErrUserNameExists
		}
	}

	// 更新用户信息
	user.UserName = req.UserName
	user.NickName = req.NickName
	user.Email = req.Email
	user.Phonenumber = req.Phonenumber
	user.Sex = req.Sex
	user.Status = req.Status
	user.UpdateBy = updateBy
	user.UpdateTime = time.Now()

	// 开启事务
	return db.WnAdmin().Transaction(func(tx *gorm.DB) error {
		// 更新用户信息
		if err := user.Update(); err != nil {
			return err
		}

		// 更新角色
		if err := user.UpdateUserRole(req.RoleID); err != nil {
			return err
		}

		return nil
	})
}

// DeleteUser 删除用户
func (s *UserService) DeleteUser(id int64) error {
	user := &SystemDB.UserDO{ID: id}
	return user.Delete()
}

// GetUserInfo 获取用户详情
func (s *UserService) GetUserInfo(id int64) (*dto.SystemUserInfoResp, error) {
	user := &SystemDB.UserDO{ID: id}
	if err := user.First(); err != nil {
		return nil, ErrUserNotFound
	}

	// 获取用户角色
	role, err := user.GetUserRole()
	if err != nil {
		return nil, err
	}

	return &dto.SystemUserInfoResp{
		ID:          user.ID,
		UserName:    user.UserName,
		NickName:    user.NickName,
		Email:       user.Email,
		Phonenumber: user.Phonenumber,
		Sex:         user.Sex,
		Status:      user.Status,
		CreateTime:  user.CreateTime.Format("2006-01-02 15:04:05"),
		CreateBy:    user.CreateBy,
		UpdateTime:  user.UpdateTime.Format("2006-01-02 15:04:05"),
		UpdateBy:    user.UpdateBy,
		Role:        role,
	}, nil
}

// UpdateUserStatus 更新用户状态
func (s *UserService) UpdateUserStatus(id int64, status string) error {
	user := &SystemDB.UserDO{
		ID:     id,
		Status: status,
	}
	return user.UpdateStatus()
}

// UpdateUserPassword 更新用户密码
func (s *UserService) UpdateUserPassword(id int64, password string) error {
	// 加密密码
	encryptedPassword, err := utils.EncryptPassword(password)
	if err != nil {
		return err
	}

	user := &SystemDB.UserDO{
		ID:       id,
		Password: encryptedPassword,
	}
	return user.UpdatePassword()
}
