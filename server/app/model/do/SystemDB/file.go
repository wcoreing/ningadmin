package SystemDB

import (
	"time"
	"wnadmin/app/provider/db"

	"gorm.io/gorm"
)

// File 文件信息
type File struct {
	Id           int64     `json:"id" gorm:"primaryKey;autoIncrement;comment:文件ID"`
	CategoryId   int64     `json:"categoryId" gorm:"comment:分类ID"`
	Name         string    `json:"name" gorm:"size:100;comment:文件名"`
	OriginalName string    `json:"originalName" gorm:"size:100;comment:原始文件名"`
	Url          string    `json:"url" gorm:"size:500;comment:文件URL"`
	Path         string    `json:"path" gorm:"size:500;comment:文件路径"`
	Size         int64     `json:"size" gorm:"comment:文件大小(字节)"`
	Type         string    `json:"type" gorm:"size:50;comment:文件类型"`
	Md5          string    `json:"md5" gorm:"size:32;comment:文件MD5"`
	Status       int       `json:"status" gorm:"default:1;comment:状态(0:禁用 1:正常)"`
	CreateBy     string    `json:"createBy" gorm:"size:64;comment:创建者"`
	CreateTime   time.Time `json:"createTime" gorm:"autoCreateTime;comment:创建时间"`
	UpdateTime   time.Time `json:"updateTime" gorm:"autoUpdateTime;comment:更新时间"`
}

// TableName 表名
func (File) TableName() string {
	return "sys_file"
}

// Query 查询构造器
func (f *File) Query() *gorm.DB {
	return db.DB.WnAdmin.Model(f).Where("status = ?", 1)
}

// First 获取第一条记录
func (f *File) First() error {
	return f.Query().Where("id = ?", f.Id).First(f).Error
}

// Create 创建记录
func (f *File) Create() error {
	return db.DB.WnAdmin.Create(f).Error
}

// Update 更新记录
func (f *File) Update() error {
	return db.DB.WnAdmin.Model(f).Updates(f).Error
}

// Delete 删除记录
func (f *File) Delete() error {
	return db.DB.WnAdmin.Model(f).Where("id = ?", f.Id).Update("status", 0).Error
}

// FileCategory 文件分类
type FileCategory struct {
	Id         int64     `json:"id" gorm:"primaryKey;autoIncrement;comment:分类ID"`
	ParentId   int64     `json:"parentId" gorm:"comment:父分类ID"`
	Name       string    `json:"name" gorm:"size:50;comment:分类名称"`
	Code       string    `json:"code" gorm:"size:50;comment:分类编码"`
	Icon       string    `json:"icon" gorm:"size:100;comment:图标"`
	Sort       int       `json:"sort" gorm:"default:0;comment:排序"`
	Status     int       `json:"status" gorm:"default:1;comment:状态(0:禁用 1:正常)"`
	CreateTime time.Time `json:"createTime" gorm:"autoCreateTime;comment:创建时间"`
	UpdateTime time.Time `json:"updateTime" gorm:"autoUpdateTime;comment:更新时间"`
}

// TableName 表名
func (FileCategory) TableName() string {
	return "sys_file_category"
}

// Query 查询构造器
func (fc *FileCategory) Query() *gorm.DB {
	return db.DB.WnAdmin.Model(fc).Where("status = ?", 1)
}

// First 获取第一条记录
func (fc *FileCategory) First() error {
	return fc.Query().Where("id = ?", fc.Id).First(fc).Error
}

// Create 创建记录
func (fc *FileCategory) Create() error {
	return db.DB.WnAdmin.Create(fc).Error
}

// Update 更新记录
func (fc *FileCategory) Update() error {
	return db.DB.WnAdmin.Model(fc).Updates(fc).Error
}

// Delete 删除记录
func (fc *FileCategory) Delete() error {
	return db.DB.WnAdmin.Model(fc).Where("id = ?", fc.Id).Update("status", 0).Error
}
