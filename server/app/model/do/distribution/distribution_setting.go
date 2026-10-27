package distribution

import (
	"time"
)

// DistributionSetting 分销设置
type DistributionSetting struct {
	Id                     int64     `gorm:"primaryKey;column:id" json:"id"`                                // 主键
	Enabled                bool      `gorm:"column:enabled" json:"enabled"`                                 // 分销功能启用状态
	Level                  int8      `gorm:"column:level" json:"level"`                                     // 分销层级
	Mode                   string    `gorm:"column:mode" json:"mode"`                                       // 分销模式
	ApplyEnabled           bool      `gorm:"column:apply_enabled" json:"applyEnabled"`                      // 分销权限申请开关
	UserType               string    `gorm:"column:user_type" json:"userType"`                              // 绑定用户类型
	BindMode               string    `gorm:"column:bind_mode" json:"bindMode"`                              // 绑定模式
	CenterType             string    `gorm:"column:center_type" json:"centerType"`                          // 分销中心类型
	PosterImages           string    `gorm:"column:poster_images" json:"posterImages"`                      // 分销海报图片
	CalcType               string    `gorm:"column:calc_type" json:"calcType"`                              // 佣金计算方式
	Level1Rate             float64   `gorm:"column:level1_rate" json:"level1Rate"`                          // 一级返佣比例
	Level2Rate             float64   `gorm:"column:level2_rate" json:"level2Rate"`                          // 二级返佣比例
	FreezeDays             int       `gorm:"column:freeze_days" json:"freezeDays"`                          // 佣金冻结天数
	SelfBuyEnabled         bool      `gorm:"column:self_buy_enabled" json:"selfBuyEnabled"`                 // 自购返佣开关
	PromotionEnabled       bool      `gorm:"column:promotion_enabled" json:"promotionEnabled"`              // 推广用户返佣开关
	MaxWithdraw            float64   `gorm:"column:max_withdraw" json:"maxWithdraw"`                        // 提现最高金额
	MinWithdraw            float64   `gorm:"column:min_withdraw" json:"minWithdraw"`                        // 提现最低金额
	WithdrawFeeRate        float64   `gorm:"column:withdraw_fee_rate" json:"withdrawFeeRate"`               // 提现手续费率
	BankWithdrawEnabled    bool      `gorm:"column:bank_withdraw_enabled" json:"bankWithdrawEnabled"`       // 银行卡提现开关
	WechatWithdrawEnabled  bool      `gorm:"column:wechat_withdraw_enabled" json:"wechatWithdrawEnabled"`   // 微信提现开关
	WechatPayMode          string    `gorm:"column:wechat_pay_mode" json:"wechatPayMode"`                   // 微信到账方式
	AlipayWithdrawEnabled  bool      `gorm:"column:alipay_withdraw_enabled" json:"alipayWithdrawEnabled"`   // 支付宝提现开关
	BalanceWithdrawEnabled bool      `gorm:"column:balance_withdraw_enabled" json:"balanceWithdrawEnabled"` // 余额提现开关
	CreateTime             time.Time `gorm:"column:create_time" json:"createTime"`                          // 创建时间
	UpdateTime             time.Time `gorm:"column:update_time" json:"updateTime"`                          // 更新时间
	CreateBy               string    `gorm:"column:create_by" json:"createBy"`                              // 创建者
	UpdateBy               string    `gorm:"column:update_by" json:"updateBy"`                              // 更新者
	Remark                 string    `gorm:"column:remark" json:"remark"`                                   // 备注
}

func (m *DistributionSetting) TableName() string {
	return "distribution_setting"
}
