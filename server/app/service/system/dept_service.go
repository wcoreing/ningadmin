package system

import (
	"errors"
	"time"
	"wnadmin/app/common/utils"
	"wnadmin/app/model/do/SystemDB"
	"wnadmin/app/model/dto"
)

// DeptService 部门服务
type DeptService struct{}

// GetDeptList 获取部门列表(树形结构)
func (s *DeptService) GetDeptList(query *dto.DeptQuery) ([]SystemDB.DeptDO, int64, error) {
	deptDO := &SystemDB.DeptDO{}

	// 如果是树形查询，返回树形结构
	if query.PageNum == 0 && query.PageSize == 0 {
		depts, err := deptDO.BuildTreeWithFilter(query.Keyword, query.Status)
		return depts, int64(len(depts)), err
	}

	// 否则返回分页列表
	var depts []SystemDB.DeptDO
	db := deptDO.Query()

	if query.Keyword != "" {
		db = db.Where("dept_name LIKE ? OR leader LIKE ?", "%"+query.Keyword+"%", "%"+query.Keyword+"%")
	}
	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}
	if query.ParentId != nil {
		db = db.Where("parent_id = ?", *query.ParentId)
	}

	var total int64
	db.Count(&total)

	offset := (query.PageNum - 1) * query.PageSize
	err := db.Offset(offset).Limit(query.PageSize).Order("order_num asc, dept_id asc").Find(&depts).Error

	return depts, total, err
}

// CreateDept 创建部门
func (s *DeptService) CreateDept(req *dto.CreateDeptReq, createBy string) error {
	now := time.Now()

	// 检查部门名称是否已存在
	var count int64
	deptDO := &SystemDB.DeptDO{}
	err := deptDO.Query().Where("dept_name = ? AND parent_id = ?", req.DeptName, req.ParentId).Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("同级部门名称已存在")
	}

	dept := &SystemDB.DeptDO{
		ParentId:   req.ParentId,
		DeptName:   req.DeptName,
		OrderNum:   req.OrderNum,
		Leader:     req.Leader,
		Phone:      req.Phone,
		Email:      req.Email,
		Status:     req.Status,
		DelFlag:    "0",
		CreateBy:   createBy,
		CreateTime: utils.JSONTime{Time: now},
		UpdateTime: utils.JSONTime{Time: now},
	}

	// 设置默认状态
	if dept.Status == "" {
		dept.Status = "0"
	}

	// 更新祖级列表
	err = dept.UpdateAncestors()
	if err != nil {
		return err
	}

	return dept.Create()
}

// UpdateDept 更新部门
func (s *DeptService) UpdateDept(req *dto.UpdateDeptReq, updateBy string) error {
	dept := &SystemDB.DeptDO{}
	err := dept.Query().Where("dept_id = ?", req.DeptId).First(dept).Error
	if err != nil {
		return err
	}

	// 检查是否将部门设为自己的子部门
	if req.ParentId == req.DeptId {
		return errors.New("不能将部门设为自己的子部门")
	}

	// 检查新的父部门是否是当前部门的子部门
	if req.ParentId != dept.ParentId && req.ParentId != 0 {
		var parent SystemDB.DeptDO
		err := parent.Query().Where("dept_id = ?", req.ParentId).First(&parent).Error
		if err != nil {
			return errors.New("父部门不存在")
		}

		// 检查是否存在循环引用
		if parent.Ancestors != "" && parent.Ancestors != "0" {
			ancestors := parent.Ancestors + "," + string(rune(req.DeptId))
			if parent.Ancestors == ancestors {
				return errors.New("不能将部门设为其子部门的子部门")
			}
		}
	}

	// 检查部门名称是否已存在（排除自己）
	var count int64
	err = dept.Query().Where("dept_name = ? AND parent_id = ? AND dept_id != ?", req.DeptName, req.ParentId, req.DeptId).Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("同级部门名称已存在")
	}

	oldParentId := dept.ParentId
	dept.ParentId = req.ParentId
	dept.DeptName = req.DeptName
	dept.OrderNum = req.OrderNum
	dept.Leader = req.Leader
	dept.Phone = req.Phone
	dept.Email = req.Email
	dept.Status = req.Status
	dept.UpdateBy = updateBy
	dept.UpdateTime = utils.JSONTime{Time: time.Now()}

	// 如果父部门发生变化，需要更新祖级列表
	if oldParentId != req.ParentId {
		err = dept.UpdateAncestors()
		if err != nil {
			return err
		}
	}

	err = dept.Update()
	if err != nil {
		return err
	}

	// 如果父部门发生变化，需要更新所有子部门的祖级列表
	if oldParentId != req.ParentId {
		err = dept.UpdateChildrenAncestors()
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteDept 删除部门
func (s *DeptService) DeleteDept(deptId int64) error {
	// 检查是否有子部门
	var count int64
	deptDO := &SystemDB.DeptDO{}
	err := deptDO.Query().Where("parent_id = ?", deptId).Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("请先删除子部门")
	}

	// 检查是否有用户关联到该部门
	// TODO: 这里需要根据实际的用户表结构来检查
	// var userCount int64
	// err = db.WnAdmin().Model(&SystemDB.UserDO{}).Where("dept_id = ?", deptId).Count(&userCount).Error
	// if err != nil {
	// 	return err
	// }
	// if userCount > 0 {
	// 	return errors.New("该部门下还有用户，不能删除")
	// }

	dept := &SystemDB.DeptDO{DeptId: deptId}
	return dept.Delete()
}

// GetDeptInfo 获取部门详情
func (s *DeptService) GetDeptInfo(deptId int64) (*SystemDB.DeptDO, error) {
	dept := &SystemDB.DeptDO{DeptId: deptId}
	err := dept.First()
	return dept, err
}

// UpdateDeptStatus 更新部门状态
func (s *DeptService) UpdateDeptStatus(deptId int64, status string) error {
	dept := &SystemDB.DeptDO{DeptId: deptId}
	err := dept.Query().Where("dept_id = ?", deptId).First(dept).Error
	if err != nil {
		return err
	}

	dept.Status = status
	dept.UpdateTime = utils.JSONTime{Time: time.Now()}
	return dept.Update()
}

// GetDeptTree 获取部门树
func (s *DeptService) GetDeptTree() ([]SystemDB.DeptDO, error) {
	deptDO := &SystemDB.DeptDO{}
	return deptDO.BuildTree()
}

// GetDeptOptions 获取部门选项列表（用于下拉选择）
func (s *DeptService) GetDeptOptions() ([]dto.DeptOption, error) {
	deptDO := &SystemDB.DeptDO{}

	// 构建树形结构
	tree, err := deptDO.BuildTree()
	if err != nil {
		return nil, err
	}

	// 转换为选项格式
	var options []dto.DeptOption
	s.buildDeptOptions(tree, &options, "")

	return options, nil
}

// buildDeptOptions 递归构建部门选项
func (s *DeptService) buildDeptOptions(depts []SystemDB.DeptDO, options *[]dto.DeptOption, prefix string) {
	for _, dept := range depts {
		*options = append(*options, dto.DeptOption{
			Value: dept.DeptId,
			Label: prefix + dept.DeptName,
		})

		if len(dept.Children) > 0 {
			s.buildDeptOptions(dept.Children, options, prefix+"├─")
		}
	}
}
