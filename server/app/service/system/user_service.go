package system

import (
	"strings"
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

	// 按部门搜索
	if query.DeptId > 0 {
		// 通过用户部门关联表过滤
		db = db.Where("user_id IN (SELECT user_id FROM sys_user_dept WHERE dept_id = ?)", query.DeptId)
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

	// 获取用户角色和部门信息
	for _, user := range users {
		role, err := user.GetUserRole()
		if err != nil {
			return nil, err
		}

		// 获取用户部门信息和数据权限
		deptRelations, err := user.GetUserDepts()
		if err != nil {
			return nil, err
		}

		var deptIds []int64
		var deptNames []string
		var dataScopes []string

		for _, relation := range deptRelations {
			deptIds = append(deptIds, relation.DeptId)
			dataScopes = append(dataScopes, relation.DataScope)

			// 获取部门名称
			deptDO := &SystemDB.DeptDO{DeptId: relation.DeptId}
			if err := deptDO.First(); err == nil {
				deptNames = append(deptNames, deptDO.DeptName)
			}
		}

		// 将部门名称用逗号连接
		deptName := strings.Join(deptNames, ", ")

		// 获取主要数据权限（第一个部门的权限）
		var dataScope string
		if len(dataScopes) > 0 {
			dataScope = dataScopes[0]
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
			DeptIds:     deptIds,
			DeptNames:   deptNames,
			DeptName:    deptName,
			DataScopes:  dataScopes,
			DataScope:   dataScope,
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

	// 获取用户部门信息
	deptRelations, err := user.GetUserDepts()
	if err != nil {
		return nil, err
	}

	var deptIds []int64
	var deptNames []string
	var dataScopes []string

	for _, relation := range deptRelations {
		deptIds = append(deptIds, relation.DeptId)
		dataScopes = append(dataScopes, relation.DataScope)

		// 获取部门名称
		deptDO := &SystemDB.DeptDO{DeptId: relation.DeptId}
		if err := deptDO.First(); err == nil {
			deptNames = append(deptNames, deptDO.DeptName)
		}
	}

	// 将部门名称用逗号连接
	deptName := strings.Join(deptNames, ", ")

	// 获取主要数据权限（第一个部门的权限）
	var dataScope string
	if len(dataScopes) > 0 {
		dataScope = dataScopes[0]
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
		DeptIds:     deptIds,
		DeptNames:   deptNames,
		DeptName:    deptName,
		DataScopes:  dataScopes,
		DataScope:   dataScope,
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
