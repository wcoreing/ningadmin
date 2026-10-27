package user

import (
	"errors"
	"time"
	"wnadmin/app/model/do"
	"wnadmin/app/model/dto"

	"golang.org/x/crypto/bcrypt"
)

// Service 用户服务
type Service struct{}

// NewService 创建用户服务实例
func NewService() *Service {
	return &Service{}
}

// GetUserList 获取用户列表
func (s *Service) GetUserList(query *dto.UserQuery) (*dto.UserListResp, error) {
	// 设置默认值
	if query.PageNum <= 0 {
		query.PageNum = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 10
	}

	// 计算偏移量
	offset := (query.PageNum - 1) * query.PageSize

	// 构建查询条件
	conditions := make(map[string]interface{})
	if query.Keyword != "" {
		conditions["keyword"] = query.Keyword
	}
	if query.Status != nil {
		conditions["status"] = *query.Status
	}
	if query.Level != nil {
		conditions["level"] = *query.Level
	}
	if query.Gender != nil {
		conditions["gender"] = *query.Gender
	}

	// 查询数据
	userModel := &do.User{}
	users, total, err := userModel.GetList(offset, query.PageSize, conditions)
	if err != nil {
		return nil, err
	}

	// 转换为VO
	userVOs := make([]dto.UserVO, 0, len(users))
	for _, user := range users {
		userVO := dto.UserVO{
			ID:            user.ID,
			Openid:        user.Openid,
			Username:      user.Username,
			Nickname:      user.Nickname,
			Gender:        user.Gender,
			GenderText:    dto.GetGenderText(user.Gender),
			Avatar:        user.Avatar,
			Phone:         user.Phone,
			Email:         user.Email,
			Birthday:      user.Birthday,
			Balance:       user.Balance,
			Points:        user.Points,
			Level:         user.Level,
			LevelText:     dto.GetLevelText(user.Level),
			TotalAmount:   user.TotalAmount,
			LastLoginTime: user.LastLoginTime,
			LastLoginIP:   user.LastLoginIP,
			Status:        user.Status,
			StatusText:    dto.GetStatusText(user.Status),
			CreateTime:    user.CreateTime,
			UpdateTime:    user.UpdateTime,
		}
		userVOs = append(userVOs, userVO)
	}

	return &dto.UserListResp{
		List:  userVOs,
		Total: total,
	}, nil
}

// CreateUser 创建用户
func (s *Service) CreateUser(req *dto.CreateUserReq) error {
	// 检查用户名是否存在
	userModel := &do.User{}
	exists, err := userModel.CheckUsernameExists(req.Username, 0)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("用户名已存在")
	}

	// 检查手机号是否存在（如果提供了手机号）
	if req.Phone != "" {
		exists, err := userModel.CheckPhoneExists(req.Phone, 0)
		if err != nil {
			return err
		}
		if exists {
			return errors.New("手机号已存在")
		}
	}

	// 检查邮箱是否存在（如果提供了邮箱）
	if req.Email != "" {
		exists, err := userModel.CheckEmailExists(req.Email, 0)
		if err != nil {
			return err
		}
		if exists {
			return errors.New("邮箱已存在")
		}
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 创建用户
	user := &do.User{
		Username:   req.Username,
		Password:   string(hashedPassword),
		Nickname:   req.Nickname,
		Gender:     req.Gender,
		Avatar:     req.Avatar,
		Phone:      req.Phone,
		Email:      req.Email,
		Birthday:   req.Birthday,
		Level:      req.Level,
		Status:     req.Status,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}

	return user.Create()
}

// UpdateUser 更新用户
func (s *Service) UpdateUser(req *dto.UpdateUserReq) error {
	// 检查用户是否存在
	userModel := &do.User{}
	err := userModel.GetByID(req.ID)
	if err != nil {
		return errors.New("用户不存在")
	}

	// 检查用户名是否存在（排除当前用户）
	exists, err := userModel.CheckUsernameExists(req.Username, req.ID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("用户名已存在")
	}

	// 检查手机号是否存在（如果提供了手机号）
	if req.Phone != "" {
		exists, err := userModel.CheckPhoneExists(req.Phone, req.ID)
		if err != nil {
			return err
		}
		if exists {
			return errors.New("手机号已存在")
		}
	}

	// 检查邮箱是否存在（如果提供了邮箱）
	if req.Email != "" {
		exists, err := userModel.CheckEmailExists(req.Email, req.ID)
		if err != nil {
			return err
		}
		if exists {
			return errors.New("邮箱已存在")
		}
	}

	// 更新用户信息
	userModel.Username = req.Username
	userModel.Nickname = req.Nickname
	userModel.Gender = req.Gender
	userModel.Avatar = req.Avatar
	userModel.Phone = req.Phone
	userModel.Email = req.Email
	userModel.Birthday = req.Birthday
	userModel.Level = req.Level
	userModel.Status = req.Status
	userModel.UpdateTime = time.Now()

	return userModel.Update()
}

// DeleteUser 删除用户
func (s *Service) DeleteUser(id int64) error {
	userModel := &do.User{}
	err := userModel.GetByID(id)
	if err != nil {
		return errors.New("用户不存在")
	}

	return userModel.Delete()
}

// GetUserInfo 获取用户详情
func (s *Service) GetUserInfo(id int64) (*dto.UserVO, error) {
	userModel := &do.User{}
	err := userModel.GetByID(id)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	userVO := &dto.UserVO{
		ID:            userModel.ID,
		Openid:        userModel.Openid,
		Username:      userModel.Username,
		Nickname:      userModel.Nickname,
		Gender:        userModel.Gender,
		GenderText:    dto.GetGenderText(userModel.Gender),
		Avatar:        userModel.Avatar,
		Phone:         userModel.Phone,
		Email:         userModel.Email,
		Birthday:      userModel.Birthday,
		Balance:       userModel.Balance,
		Points:        userModel.Points,
		Level:         userModel.Level,
		LevelText:     dto.GetLevelText(userModel.Level),
		TotalAmount:   userModel.TotalAmount,
		LastLoginTime: userModel.LastLoginTime,
		LastLoginIP:   userModel.LastLoginIP,
		Status:        userModel.Status,
		StatusText:    dto.GetStatusText(userModel.Status),
		CreateTime:    userModel.CreateTime,
		UpdateTime:    userModel.UpdateTime,
	}

	return userVO, nil
}

// UpdateUserStatus 更新用户状态
func (s *Service) UpdateUserStatus(id int64, status int) error {
	userModel := &do.User{}
	err := userModel.GetByID(id)
	if err != nil {
		return errors.New("用户不存在")
	}

	return userModel.UpdateStatus(id, status)
}

// UpdateUserPassword 更新用户密码
func (s *Service) UpdateUserPassword(id int64, password string) error {
	userModel := &do.User{}
	err := userModel.GetByID(id)
	if err != nil {
		return errors.New("用户不存在")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return userModel.UpdatePassword(id, string(hashedPassword))
}

// UpdateUserBalance 更新用户余额
func (s *Service) UpdateUserBalance(req *dto.UpdateUserBalanceReq) error {
	userModel := &do.User{}
	err := userModel.GetByID(req.ID)
	if err != nil {
		return errors.New("用户不存在")
	}

	// 记录变动前的余额
	beforeBalance := userModel.Balance
	afterBalance := *req.Balance

	// 更新余额
	err = userModel.UpdateBalance(req.ID, *req.Balance)
	if err != nil {
		return err
	}

	// 计算变动类型和数量
	var actionType int
	var amount float64
	if afterBalance >= beforeBalance {
		actionType = do.ActionTypeIncome // 收入
		amount = afterBalance - beforeBalance
	} else {
		actionType = do.ActionTypeExpend // 支出
		amount = beforeBalance - afterBalance
	}

	// 创建变动记录
	changeLog := &do.UserChangeLog{
		UserID:       req.ID,
		ChangeType:   do.ChangeTypeBalance,
		ActionType:   actionType,
		Amount:       amount,
		BeforeValue:  beforeBalance,
		AfterValue:   afterBalance,
		SourceType:   do.SourceTypeAdmin, // 管理员调整
		Remark:       req.Remark,
		OperatorName: "管理员", // TODO: 从上下文获取操作员信息
		CreateTime:   time.Now(),
	}

	return changeLog.Create()
}

// UpdateUserPoints 更新用户积分
func (s *Service) UpdateUserPoints(req *dto.UpdateUserPointsReq) error {
	userModel := &do.User{}
	err := userModel.GetByID(req.ID)
	if err != nil {
		return errors.New("用户不存在")
	}

	// 记录变动前的积分
	beforePoints := float64(userModel.Points)
	afterPoints := float64(*req.Points)

	// 更新积分
	err = userModel.UpdatePoints(req.ID, *req.Points)
	if err != nil {
		return err
	}

	// 计算变动类型和数量
	var actionType int
	var amount float64
	if afterPoints >= beforePoints {
		actionType = do.ActionTypeIncome // 增加
		amount = afterPoints - beforePoints
	} else {
		actionType = do.ActionTypeExpend // 减少
		amount = beforePoints - afterPoints
	}

	// 创建变动记录
	changeLog := &do.UserChangeLog{
		UserID:       req.ID,
		ChangeType:   do.ChangeTypePoints,
		ActionType:   actionType,
		Amount:       amount,
		BeforeValue:  beforePoints,
		AfterValue:   afterPoints,
		SourceType:   do.SourceTypeAdmin, // 管理员调整
		Remark:       req.Remark,
		OperatorName: "管理员", // TODO: 从上下文获取操作员信息
		CreateTime:   time.Now(),
	}

	return changeLog.Create()
}

// UpdateLastLogin 更新最后登录信息
func (s *Service) UpdateLastLogin(id int64, loginIP string) error {
	userModel := &do.User{}
	return userModel.UpdateLastLogin(id, time.Now(), loginIP)
}

// GetUserChangeLogList 获取用户变动记录列表
func (s *Service) GetUserChangeLogList(query *dto.UserChangeLogQuery) (*dto.UserChangeLogListResp, error) {
	changeLogModel := &do.UserChangeLog{}

	// 构建查询条件
	conditions := make(map[string]interface{})
	if query.UserID != nil {
		conditions["userId"] = *query.UserID
	}
	if query.ChangeType != nil {
		conditions["changeType"] = *query.ChangeType
	}
	if query.ActionType != nil {
		conditions["actionType"] = *query.ActionType
	}
	if query.SourceType != nil {
		conditions["sourceType"] = *query.SourceType
	}
	if query.StartTime != "" {
		conditions["startTime"] = query.StartTime
	}
	if query.EndTime != "" {
		conditions["endTime"] = query.EndTime
	}

	// 分页参数
	offset := (query.PageNum - 1) * query.PageSize
	limit := query.PageSize

	// 获取数据
	logs, total, err := changeLogModel.GetList(offset, limit, conditions)
	if err != nil {
		return nil, err
	}

	// 转换为VO
	var voList []dto.UserChangeLogVO
	for _, log := range logs {
		// 获取用户信息
		userModel := &do.User{}
		userModel.GetByID(log.UserID)

		vo := dto.UserChangeLogVO{
			ID:             log.ID,
			UserID:         log.UserID,
			Username:       userModel.Username,
			Nickname:       userModel.Nickname,
			ChangeType:     log.ChangeType,
			ChangeTypeText: log.GetChangeTypeText(),
			ActionType:     log.ActionType,
			ActionTypeText: log.GetActionTypeText(),
			Amount:         log.Amount,
			BeforeValue:    log.BeforeValue,
			AfterValue:     log.AfterValue,
			SourceType:     log.SourceType,
			SourceTypeText: log.GetSourceTypeText(),
			SourceID:       log.SourceID,
			Remark:         log.Remark,
			OperatorID:     log.OperatorID,
			OperatorName:   log.OperatorName,
			ExpireTime:     log.ExpireTime,
			CreateTime:     log.CreateTime,
		}
		voList = append(voList, vo)
	}

	return &dto.UserChangeLogListResp{
		List:  voList,
		Total: total,
	}, nil
}
