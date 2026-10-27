package SystemDB

import (
	"fmt"
	"strconv"
	"strings"
	"wnadmin/app/common/utils"
	"wnadmin/app/provider/db"

	"gorm.io/gorm"
)

// DeptDO 部门数据对象
type DeptDO struct {
	DeptId     int64          `gorm:"column:dept_id;primaryKey;autoIncrement" json:"deptId"`    // 部门ID
	ParentId   int64          `gorm:"column:parent_id;not null;default:0" json:"parentId"`      // 父部门ID
	Ancestors  string         `gorm:"column:ancestors;size:50" json:"ancestors"`                // 祖级列表
	DeptName   string         `gorm:"column:dept_name;size:30;not null" json:"deptName"`        // 部门名称
	OrderNum   int            `gorm:"column:order_num;not null;default:0" json:"orderNum"`      // 显示顺序
	Leader     string         `gorm:"column:leader;size:20" json:"leader"`                      // 负责人
	Phone      string         `gorm:"column:phone;size:11" json:"phone"`                        // 联系电话
	Email      string         `gorm:"column:email;size:50" json:"email"`                        // 邮箱
	Status     string         `gorm:"column:status;size:1;not null;default:0" json:"status"`    // 部门状态（0正常 1停用）
	DelFlag    string         `gorm:"column:del_flag;size:1;not null;default:0" json:"delFlag"` // 删除标志（0代表存在 2代表删除）
	CreateBy   string         `gorm:"column:create_by;size:64" json:"createBy"`                 // 创建者
	CreateTime utils.JSONTime `gorm:"column:create_time;autoCreateTime" json:"createTime"`      // 创建时间
	UpdateBy   string         `gorm:"column:update_by;size:64" json:"updateBy"`                 // 更新者
	UpdateTime utils.JSONTime `gorm:"column:update_time;autoUpdateTime" json:"updateTime"`      // 更新时间
	Children   []DeptDO       `json:"children" gorm:"-"`                                        // 子部门列表
}

// TableName 表名
func (DeptDO) TableName() string {
	return "sys_dept"
}

// Query 查询构造器
func (d *DeptDO) Query() *gorm.DB {
	return db.WnAdmin().Model(d).Where("del_flag = ?", "0")
}

// First 获取第一条记录
func (d *DeptDO) First() error {
	return d.Query().Where("dept_id = ?", d.DeptId).First(d).Error
}

// Create 创建记录
func (d *DeptDO) Create() error {
	return db.WnAdmin().Create(d).Error
}

// Update 更新记录
func (d *DeptDO) Update() error {
	return db.WnAdmin().Model(d).Updates(d).Error
}

// Delete 删除记录
func (d *DeptDO) Delete() error {
	return db.WnAdmin().Model(d).Where("dept_id = ?", d.DeptId).Update("del_flag", "2").Error
}

// BuildTree 构建部门树
func (d *DeptDO) BuildTree() ([]DeptDO, error) {
	// 获取所有部门
	var allDepts []DeptDO
	err := d.Query().Order("order_num asc, dept_id asc").Find(&allDepts).Error
	if err != nil {
		return nil, err
	}

	// 构建树形结构的递归函数
	var buildChildren func(parentId int64) []DeptDO
	buildChildren = func(parentId int64) []DeptDO {
		var children []DeptDO
		for _, dept := range allDepts {
			if dept.ParentId == parentId {
				// 递归构建子节点
				dept.Children = buildChildren(dept.DeptId)
				children = append(children, dept)
			}
		}
		return children
	}

	// 构建根节点树
	return buildChildren(0), nil
}

// BuildTreeWithFilter 构建带搜索条件的部门树
func (d *DeptDO) BuildTreeWithFilter(keyword, status string) ([]DeptDO, error) {
	// 获取所有部门
	var allDepts []DeptDO
	err := d.Query().Order("order_num asc, dept_id asc").Find(&allDepts).Error
	if err != nil {
		return nil, err
	}

	// 如果没有搜索条件，直接返回完整树
	if keyword == "" && status == "" {
		return d.BuildTree()
	}

	// 辅助函数：检查是否包含关键字
	containsKeyword := func(text, keyword string) bool {
		if keyword == "" {
			return true
		}
		return strings.Contains(text, keyword)
	}

	// 根据条件筛选匹配的部门
	var matchedDeptIDs []int64
	for _, dept := range allDepts {
		isMatched := true

		// 关键字匹配
		if keyword != "" {
			isMatched = false
			if containsKeyword(dept.DeptName, keyword) || containsKeyword(dept.Leader, keyword) {
				isMatched = true
			}
		}

		// 状态匹配
		if status != "" && dept.Status != status {
			isMatched = false
		}

		if isMatched {
			matchedDeptIDs = append(matchedDeptIDs, dept.DeptId)
			// 添加所有父级ID
			parentID := dept.ParentId
			for parentID != 0 {
				if !containsDeptId(matchedDeptIDs, parentID) {
					matchedDeptIDs = append(matchedDeptIDs, parentID)
				}
				// 查找父级的父级
				for _, m := range allDepts {
					if m.DeptId == parentID {
						parentID = m.ParentId
						break
					}
				}
			}
		}
	}

	// 筛选出匹配的部门
	var filteredDepts []DeptDO
	for _, dept := range allDepts {
		if containsDeptId(matchedDeptIDs, dept.DeptId) {
			filteredDepts = append(filteredDepts, dept)
		}
	}

	// 构建树形结构的递归函数
	var buildChildren func(parentId int64) []DeptDO
	buildChildren = func(parentId int64) []DeptDO {
		var children []DeptDO
		for _, dept := range filteredDepts {
			if dept.ParentId == parentId {
				// 递归构建子节点
				dept.Children = buildChildren(dept.DeptId)
				children = append(children, dept)
			}
		}
		return children
	}

	// 构建根节点树
	return buildChildren(0), nil
}

// UpdateAncestors 更新祖级列表
func (d *DeptDO) UpdateAncestors() error {
	if d.ParentId == 0 {
		d.Ancestors = "0"
	} else {
		// 获取父部门信息
		var parent DeptDO
		err := parent.Query().Where("dept_id = ?", d.ParentId).First(&parent).Error
		if err != nil {
			return err
		}
		d.Ancestors = parent.Ancestors + "," + strconv.FormatInt(d.ParentId, 10)
	}
	return nil
}

// UpdateChildrenAncestors 更新子部门的祖级列表
func (d *DeptDO) UpdateChildrenAncestors() error {
	// 查找所有子部门
	var children []DeptDO
	err := d.Query().Where("parent_id = ?", d.DeptId).Find(&children).Error
	if err != nil {
		return err
	}

	// 更新每个子部门的祖级列表
	for _, child := range children {
		oldAncestors := child.Ancestors
		newAncestors := d.Ancestors + "," + strconv.FormatInt(d.DeptId, 10)

		// 更新当前子部门
		err = db.WnAdmin().Model(&child).Where("dept_id = ?", child.DeptId).Update("ancestors", newAncestors).Error
		if err != nil {
			return err
		}

		// 递归更新所有后代部门的祖级列表
		err = d.updateDescendantsAncestors(oldAncestors, newAncestors)
		if err != nil {
			return err
		}
	}

	return nil
}

// updateDescendantsAncestors 递归更新后代部门的祖级列表
func (d *DeptDO) updateDescendantsAncestors(oldAncestors, newAncestors string) error {
	// 使用SQL批量更新所有后代部门
	sql := fmt.Sprintf("UPDATE %s SET ancestors = REPLACE(ancestors, ?, ?) WHERE ancestors LIKE ? AND del_flag = '0'", d.TableName())
	return db.WnAdmin().Exec(sql, oldAncestors, newAncestors, oldAncestors+"%").Error
}

// containsDeptId 检查切片是否包含指定元素
func containsDeptId(slice []int64, item int64) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
