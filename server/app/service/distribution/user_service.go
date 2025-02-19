package distribution

type UserService struct {
}

// 用户申请成为分销员
func (s *UserService) ApplyDistributor(userId int64, parentId int64) error {
	// 1. 检查分销功能是否开启
	// 2. 检查用户是否已是分销员
	// 3. 检查上级是否合法
	// 4. 创建分销用户记录
	// 5. 如果是指定分销模式，需要等待审核
	// 6. 如果是满额分销，检查用户消费金额
	// 7. 建立分销关系
	return nil
}

// 绑定上级分销员
func (s *UserService) BindParent(userId int64, parentId int64) error {
	// 1. 检查绑定模式(永久/临时/有效期)
	// 2. 检查用户类型(新用户/所有用户)
	// 3. 建立分销关系
	// 4. 更新团队人数
	return nil
}

// 计算订单佣金
func (s *UserService) CalculateCommission(orderId int64) error {
	// 1. 获取订单信息
	// 2. 获取分销设置
	// 3. 计算一级佣金
	// 4. 计算二级佣金
	// 5. 创建佣金记录
	// 6. 更新分销员余额
	return nil
}

// 佣金解冻
func (s *UserService) UnfreezeCommission(incomeId int64) error {
	// 1. 检查冻结期是否到期
	// 2. 更新佣金状态
	// 3. 更新可用余额
	return nil
}

// 提现申请
func (s *UserService) ApplyWithdraw(userId int64, amount float64, withdrawType string) error {
	// 1. 检查提现设置
	// 2. 检查提现金额
	// 3. 检查可用余额
	// 4. 创建提现记录
	// 5. 扣减可用余额
	return nil
}
