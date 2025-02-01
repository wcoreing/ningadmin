package shop

import (
	"time"
	"wnadmin/app/provider/db"

	"gorm.io/gorm"
)

// ShopSupplier 供应商
type ShopSupplier struct {
	ID              int64     `gorm:"primaryKey;column:id" json:"id"`
	Name            string    `gorm:"column:name" json:"name"`                        // 供应商名称
	Contact         string    `gorm:"column:contact" json:"contact"`                  // 联系人
	Phone           string    `gorm:"column:phone" json:"phone"`                      // 联系电话
	Email           string    `gorm:"column:email" json:"email"`                      // 邮箱
	Address         string    `gorm:"column:address" json:"address"`                  // 地址
	BusinessLicense string    `gorm:"column:business_license" json:"businessLicense"` // 营业执照
	Description     string    `gorm:"column:description" json:"description"`          // 供应商描述
	Level           int8      `gorm:"column:level" json:"level"`                      // 供应商等级
	Sort            int       `gorm:"column:sort" json:"sort"`                        // 排序
	Status          int8      `gorm:"column:status" json:"status"`                    // 状态
	IsDelete        int8      `gorm:"column:is_delete" json:"isDelete"`               // 是否删除
	CreateTime      time.Time `gorm:"column:create_time" json:"createTime"`           // 创建时间
	UpdateTime      time.Time `gorm:"column:update_time" json:"updateTime"`           // 更新时间
}

// TableName 表名
func (s *ShopSupplier) TableName() string {
	return "shop_supplier"
}

// Query 获取查询构建器
func (s *ShopSupplier) Query() *gorm.DB {
	query := db.DB.WnAdmin.Model(s).Where("is_delete = ?", 0)

	// 供应商名称
	if s.Name != "" {
		query = query.Where("name LIKE ?", "%"+s.Name+"%")
	}

	// 状态
	if s.Status > 0 {
		query = query.Where("status = ?", s.Status)
	}

	// 供应商等级
	if s.Level > 0 {
		query = query.Where("level = ?", s.Level)
	}

	return query.Order("sort DESC, id DESC")
}

// Find 查询供应商列表
func (s *ShopSupplier) Find(dest interface{}) error {
	return s.Query().Find(dest).Error
}

// First 获取单个供应商
func (s *ShopSupplier) First() error {
	return s.Query().First(s).Error
}
