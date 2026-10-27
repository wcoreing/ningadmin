package SystemDB

import (
	"strings"
	"time"
	"wnadmin/app/provider/db"

	"gorm.io/gorm"
)

// ConfigDO 系统配置
type ConfigDO struct {
	ID          int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`           // 配置ID
	ParentId    int64      `gorm:"column:parent_id;not null;default:0" json:"parentId"`    // 父级ID
	ModuleKey   string     `gorm:"column:module_key;size:100;not null" json:"moduleKey"`   // 模块标识
	ModuleName  string     `gorm:"column:module_name;size:100;not null" json:"moduleName"` // 模块标识
	ConfigName  string     `gorm:"column:config_name;size:100;not null" json:"configName"` // 配置名称
	ConfigKey   string     `gorm:"column:config_key;size:100;not null" json:"configKey"`   // 配置键名
	ConfigType  string     `gorm:"column:config_type;size:10;not null" json:"configType"`  // 配置类型
	ConfigValue string     `gorm:"column:config_value;type:text" json:"configValue"`       // 配置值
	Level       int        `gorm:"column:level;not null;default:1" json:"level"`           // 层级
	Sort        int        `gorm:"column:sort;not null;default:0" json:"sort"`             // 排序
	Status      int        `gorm:"column:status;not null;default:1" json:"status"`         // 状态（0禁用 1启用）
	Remark      string     `gorm:"column:remark;size:500" json:"remark"`                   // 备注
	DelFlag     string     `gorm:"column:del_flag;size:1;default:0" json:"delFlag"`        // 删除标志（0代表存在 2代表删除）
	CreateTime  time.Time  `gorm:"column:create_time;autoCreateTime" json:"createTime"`    // 创建时间
	UpdateTime  time.Time  `gorm:"column:update_time;autoUpdateTime" json:"updateTime"`    // 更新时间
	Children    []ConfigDO `json:"children" gorm:"-"`                                      // 子配置列表
}

// TableName 表名
func (ConfigDO) TableName() string {
	return "sys_config"
}

// Query 查询构造器
func (c *ConfigDO) Query() *gorm.DB {
	return db.WnAdmin().Model(c).Where("del_flag = ?", "0")
}

// First 获取第一条记录
func (c *ConfigDO) First() error {
	return c.Query().Where("id = ?", c.ID).First(c).Error
}

// Create 创建记录
func (c *ConfigDO) Create() error {
	return db.WnAdmin().Create(c).Error
}

// Update 更新记录
func (c *ConfigDO) Update() error {
	return db.WnAdmin().Model(c).Updates(c).Error
}

// Delete 删除记录
func (c *ConfigDO) Delete() error {
	return db.WnAdmin().Model(c).Where("id = ?", c.ID).Update("del_flag", "2").Error
}

// BuildTree 构建配置树
func (c *ConfigDO) BuildTree() ([]ConfigDO, error) {
	// 获取所有配置
	var allConfigs []ConfigDO
	err := c.Query().Order("sort desc, id asc").Find(&allConfigs).Error
	if err != nil {
		return nil, err
	}

	// 构建树形结构的递归函数
	var buildChildren func(parentId int64) []ConfigDO
	buildChildren = func(parentId int64) []ConfigDO {
		var children []ConfigDO
		for _, config := range allConfigs {
			if config.ParentId == parentId {
				// 递归构建子节点
				config.Children = buildChildren(config.ID)
				children = append(children, config)
			}
		}
		return children
	}

	// 构建根节点树
	return buildChildren(0), nil
}

// BuildTreeWithFilter 构建带搜索条件的配置树
func (c *ConfigDO) BuildTreeWithFilter(keyword, moduleKey, status string) ([]ConfigDO, error) {
	// 获取所有配置
	var allConfigs []ConfigDO
	err := c.Query().Order("sort desc, id asc").Find(&allConfigs).Error
	if err != nil {
		return nil, err
	}

	// 如果没有搜索条件，直接返回完整树
	if keyword == "" && moduleKey == "" && status == "" {
		return c.BuildTree()
	}

	// 辅助函数：检查是否包含关键字
	containsKeyword := func(text, keyword string) bool {
		if keyword == "" {
			return true
		}
		return strings.Contains(text, keyword)
	}

	// 根据条件筛选匹配的配置
	var matchedConfigIDs []int64
	for _, config := range allConfigs {
		isMatched := true

		// 关键字匹配
		if keyword != "" {
			isMatched = false
			if containsKeyword(config.ConfigName, keyword) || containsKeyword(config.ConfigKey, keyword) {
				isMatched = true
			}
		}

		// 模块匹配
		if moduleKey != "" && config.ModuleKey != moduleKey {
			isMatched = false
		}

		// 状态匹配
		if status != "" {
			statusInt := 0
			if status == "1" {
				statusInt = 1
			}
			if config.Status != statusInt {
				isMatched = false
			}
		}

		if isMatched {
			matchedConfigIDs = append(matchedConfigIDs, config.ID)
			// 添加所有父级ID
			parentID := config.ParentId
			for parentID != 0 {
				if !contains(matchedConfigIDs, parentID) {
					matchedConfigIDs = append(matchedConfigIDs, parentID)
				}
				// 查找父级的父级
				for _, m := range allConfigs {
					if m.ID == parentID {
						parentID = m.ParentId
						break
					}
				}
			}
		}
	}

	// 筛选出匹配的配置
	var filteredConfigs []ConfigDO
	for _, config := range allConfigs {
		if contains(matchedConfigIDs, config.ID) {
			filteredConfigs = append(filteredConfigs, config)
		}
	}

	// 构建树形结构的递归函数
	var buildChildren func(parentId int64) []ConfigDO
	buildChildren = func(parentId int64) []ConfigDO {
		var children []ConfigDO
		for _, config := range filteredConfigs {
			if config.ParentId == parentId {
				// 递归构建子节点
				config.Children = buildChildren(config.ID)
				children = append(children, config)
			}
		}
		return children
	}

	// 构建根节点树
	return buildChildren(0), nil
}

// contains 检查切片是否包含指定元素
func contains(slice []int64, item int64) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
