package cms

import (
	"encoding/json"
	"wnadmin/app/model/do/shop"
	"wnadmin/app/provider/db"

	"gorm.io/gorm"
)

const (
	// Module 内容分类在通用分类里的模块标识
	Module = "cms"
	// TypePage 单页，对外展示最新一篇内容
	TypePage = "page"
	// TypeList 列表，按栏目聚合多篇内容
	TypeList = "list"
)

// Category 内容分类。树、名称、图片和排序继承通用分类，内容专有属性放 extra。
type Category struct {
	shop.Category
}

// TableName 与通用分类同一张表
func (Category) TableName() string {
	return "category"
}

// Query 只取内容分类
func (c *Category) Query() *gorm.DB {
	return db.WnAdmin().Model(c).Where("is_delete = 0 AND module = ?", Module)
}

// Attr 内容分类扩展属性
type Attr struct {
	Code           string `json:"code"`
	Type           string `json:"type"`
	Summary        string `json:"summary"`
	SeoTitle       string `json:"seoTitle"`
	SeoKeywords    string `json:"seoKeywords"`
	SeoDescription string `json:"seoDescription"`
}

// Attr 读出扩展属性。缺省字段保持零值。
func (c Category) Attr() Attr {
	var attr Attr
	if len(c.Extra) > 0 {
		_ = json.Unmarshal(c.Extra, &attr)
	}
	return attr
}
