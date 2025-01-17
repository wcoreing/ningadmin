package system

import (
	"time"
	"wnadmin/app/common/utils"
	"wnadmin/app/model/do/SystemDB"
	"wnadmin/app/model/dto"
)

// UserService 用户服务
type UserService struct{}

// GetUserList 获取用户列表
func (s *UserService) GetUserList(query *dto.UserQuery) ([]SystemDB.UserDO, int64, error) {
	var users []SystemDB.UserDO
	userDO := &SystemDB.UserDO{}
	db := userDO.Query()

	if query.Keyword != "" {
		db = db.Where("user_name LIKE ? OR nick_name LIKE ?", "%"+query.Keyword+"%", "%"+query.Keyword+"%")
	}
	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}

	var total int64
	db.Count(&total)

	offset := (query.PageNum - 1) * query.PageSize
	err := db.Offset(offset).Limit(query.PageSize).Order("user_id desc").Find(&users).Error

	return users, total, err
}

// CreateUser 创建用户
func (s *UserService) CreateUser(req *dto.CreateUserReq, createBy string) error {
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

	return user.Create()
}

// UpdateUser 更新用户
func (s *UserService) UpdateUser(req *dto.UpdateUserReq, updateBy string) error {
	user := &SystemDB.UserDO{}
	err := user.Query().Where("user_id = ?", req.ID).First(user).Error
	if err != nil {
		return err
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

	user.UserName = req.UserName
	user.NickName = req.NickName
	user.Email = req.Email
	user.Phonenumber = req.Phonenumber
	user.Sex = req.Sex
	user.Status = req.Status
	user.UpdateBy = updateBy
	user.UpdateTime = time.Now()

	return user.Update()
}

// DeleteUser 删除用户
func (s *UserService) DeleteUser(id int64) error {
	user := &SystemDB.UserDO{ID: id}
	return user.Delete()
}

// GetUserInfo 获取用户详情
func (s *UserService) GetUserInfo(id int64) (*SystemDB.UserDO, error) {
	user := &SystemDB.UserDO{ID: id}
	err := user.First()
	return user, err
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
