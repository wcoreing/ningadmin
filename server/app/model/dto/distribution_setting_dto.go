package dto

// DistributionSettingDTO 分销设置DTO
type DistributionSettingDTO struct {
	// 分销模式
	Enabled      bool   `json:"enabled"`      // 分销功能启用状态
	Level        int8   `json:"level"`        // 分销层级
	Mode         string `json:"mode"`         // 分销模式
	ApplyEnabled bool   `json:"applyEnabled"` // 分销权限申请开关
	UserType     string `json:"userType"`     // 绑定用户类型
	BindMode     string `json:"bindMode"`     // 绑定模式
	CenterType   string `json:"centerType"`   // 分销中心类型
	PosterImages string `json:"posterImages"` // 分销海报图片

	// 返佣设置
	CalcType         string  `json:"calcType"`         // 佣金计算方式
	Level1Rate       float64 `json:"level1Rate"`       // 一级返佣比例
	Level2Rate       float64 `json:"level2Rate"`       // 二级返佣比例
	FreezeDays       int     `json:"freezeDays"`       // 佣金冻结天数
	SelfBuyEnabled   bool    `json:"selfBuyEnabled"`   // 自购返佣开关
	PromotionEnabled bool    `json:"promotionEnabled"` // 推广用户返佣开关

	// 提现设置
	MaxWithdraw            float64 `json:"maxWithdraw"`            // 提现最高金额
	MinWithdraw            float64 `json:"minWithdraw"`            // 提现最低金额
	WithdrawFeeRate        float64 `json:"withdrawFeeRate"`        // 提现手续费率
	BankWithdrawEnabled    bool    `json:"bankWithdrawEnabled"`    // 银行卡提现开关
	WechatWithdrawEnabled  bool    `json:"wechatWithdrawEnabled"`  // 微信提现开关
	WechatPayMode          string  `json:"wechatPayMode"`          // 微信到账方式
	AlipayWithdrawEnabled  bool    `json:"alipayWithdrawEnabled"`  // 支付宝提现开关
	BalanceWithdrawEnabled bool    `json:"balanceWithdrawEnabled"` // 余额提现开关
}
