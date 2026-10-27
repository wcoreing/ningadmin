package shop

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
	"wnadmin/app/provider/db"

	"gorm.io/gorm"
)

// Extra 模块自己的扩展属性。商品分类用不到，内容分类把编码和发布方式放在这里。
type Extra json.RawMessage

// Scan 读取 JSON，空值当成空对象。
func (e *Extra) Scan(value interface{}) error {
	if value == nil {
		*e = Extra("{}")
		return nil
	}
	switch v := value.(type) {
	case []byte:
		if len(v) == 0 {
			*e = Extra("{}")
			return nil
		}
		*e = Extra(append([]byte(nil), v...))
		return nil
	case string:
		if v == "" {
			*e = Extra("{}")
			return nil
		}
		*e = Extra(v)
		return nil
	default:
		return fmt.Errorf("分类扩展属性类型不正确")
	}
}

// Value 写入 JSON。
func (e Extra) Value() (driver.Value, error) {
	if len(e) == 0 {
		return []byte("{}"), nil
	}
	return []byte(e), nil
}

// MarshalJSON 按对象输出，避免被编成字符串。
func (e Extra) MarshalJSON() ([]byte, error) {
	if len(e) == 0 {
		return []byte("{}"), nil
	}
	return []byte(e), nil
}

// UnmarshalJSON 接收对象。
func (e *Extra) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" {
		*e = Extra("{}")
		return nil
	}
	if !json.Valid(data) {
		return fmt.Errorf("分类扩展属性不是合法 JSON")
	}
	*e = Extra(append([]byte(nil), data...))
	return nil
}

// Category 通用分类表
type Category struct {
	ID         int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`         // 分类ID
	ParentId   int64      `gorm:"column:parent_id;not null;default:0" json:"parent_id"` // 父级ID
	Name       string     `gorm:"column:name;not null" json:"name"`                     // 分类名称
	Icon       string     `gorm:"column:icon" json:"icon"`                              // 分类图标
	Image      string     `gorm:"column:image" json:"image"`                            // 分类图片
	Level      int        `gorm:"column:level;not null;default:1" json:"level"`         // 层级
	Sort       int        `gorm:"column:sort;not null;default:0" json:"sort"`           // 排序
	Status     int        `gorm:"column:status;not null;default:1" json:"status"`       // 状态
	Module     string     `gorm:"column:module;not null;default:goods" json:"module"`   // 模块标识
	Extra      Extra      `gorm:"column:extra;type:json" json:"extra"`                  // 模块扩展属性
	IsDelete   int        `gorm:"column:is_delete;not null;default:0" json:"is_delete"` // 是否删除
	CreateTime time.Time  `gorm:"column:create_time;autoCreateTime" json:"create_time"` // 创建时间
	UpdateTime time.Time  `gorm:"column:update_time;autoUpdateTime" json:"update_time"` // 更新时间
	Children   []Category `json:"children" gorm:"-"`                                    // 子分类列表
}

// TableName 表名
func (c *Category) TableName() string {
	return "category"
}

// Query 查询构造器
func (c *Category) Query() *gorm.DB {
	return db.WnAdmin().Model(c).Where("is_delete = 0")
}

// BuildTree 构建分类树
func (c *Category) BuildTree(module string) ([]Category, error) {
	var categories []Category

	// 查询指定模块的所有分类
	if err := c.Query().Where("module = ?", module).Order("sort DESC, id ASC").Find(&categories).Error; err != nil {
		return nil, err
	}

	nodes := make(map[int64]*Category, len(categories))
	for i := range categories {
		categories[i].Children = nil
		nodes[categories[i].ID] = &categories[i]
	}

	var attach func(id int64) Category
	attach = func(id int64) Category {
		node := *nodes[id]
		node.Children = nil
		for _, row := range categories {
			if row.ParentId == id {
				node.Children = append(node.Children, attach(row.ID))
			}
		}
		return node
	}

	rootCategories := make([]Category, 0)
	for _, row := range categories {
		if row.ParentId == 0 {
			rootCategories = append(rootCategories, attach(row.ID))
		}
	}
	return rootCategories, nil
}

// GetByModule 根据模块获取分类列表
func (c *Category) GetByModule(module string) ([]Category, error) {
	var categories []Category
	err := c.Query().Where("module = ?", module).Order("sort DESC, id ASC").Find(&categories).Error
	return categories, err
}
