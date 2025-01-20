package form

import (
	"errors"
	"time"
	formdo "wnadmin/app/model/do/form"
	"wnadmin/app/model/dto"
	"wnadmin/app/provider/db"

	"gorm.io/gorm"
)

// FormListService 表单服务
type FormListService struct {
	db *gorm.DB
}

// NewFormListService 创建表单服务实例
func NewFormListService() *FormListService {
	return &FormListService{
		db: db.DB.WnAdmin,
	}
}

// GetFormList 获取表单列表
func (s *FormListService) GetFormList(query *dto.FormListQuery) ([]formdo.Form, int64, error) {
	var forms []formdo.Form
	var total int64

	// 联表查询表单和模板
	db := s.db.Model(&formdo.Form{}).
		Select("form.*, template.name as template_name").
		Joins("LEFT JOIN form_template template ON form.template_id = template.id").
		Where("form.form_del_flag = ?", "0")

	// 条件查询
	if query.FormName != "" {
		db = db.Where("form.form_name LIKE ?", "%"+query.FormName+"%")
	}
	if query.FormStatus != "" {
		db = db.Where("form.form_status = ?", query.FormStatus)
	}

	// 统计总数
	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	err = db.Offset((query.PageNum - 1) * query.PageSize).
		Limit(query.PageSize).
		Order("form.id DESC").
		Find(&forms).Error
	if err != nil {
		return nil, 0, err
	}

	return forms, total, nil
}

// GetFormDetail 获取表单详情
func (s *FormListService) GetFormDetail(id int64) (*formdo.Form, error) {
	var formInfo formdo.Form
	err := s.db.Where("id = ? AND form_del_flag = ?", id, "0").First(&formInfo).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("表单不存在")
		}
		return nil, err
	}
	return &formInfo, nil
}

// CreateForm 创建表单
func (s *FormListService) CreateForm(form *formdo.Form, username string) (*formdo.Form, error) {
	// 设置基础字段
	now := time.Now()
	form.FormCreateTime = &now
	form.FormUpdateTime = &now
	form.FormCreateBy = username
	form.FormUpdateBy = username
	form.FormStatus = "0"
	form.FormDelFlag = "0"

	// 创建表单
	err := s.db.Create(form).Error
	if err != nil {
		return nil, errors.New("创建表单失败: " + err.Error())
	}

	return form, nil
}

// UpdateForm 更新表单
func (s *FormListService) UpdateForm(form *formdo.Form, username string) (*formdo.Form, error) {
	// 检查表单是否存在
	var existForm formdo.Form
	err := s.db.Where("id = ? AND form_del_flag = ?", form.ID, "0").First(&existForm).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("表单不存在")
		}
		return nil, err
	}

	// 设置更新字段
	now := time.Now()
	updates := map[string]interface{}{
		"template_id":      form.TemplateID,
		"form_name":        form.FormName,
		"form_description": form.FormDescription,
		"form_attributes":  form.FormAttributes,
		"form_status":      form.FormStatus,
		"form_update_time": now,
		"form_update_by":   username,
	}

	// 更新表单
	err = s.db.Model(&existForm).Updates(updates).Error
	if err != nil {
		return nil, errors.New("更新表单失败: " + err.Error())
	}

	// 重新查询完整数据
	err = s.db.First(&existForm, existForm.ID).Error
	if err != nil {
		return nil, errors.New("获取更新后的表单失败: " + err.Error())
	}

	return &existForm, nil
}

// DeleteForm 删除表单
func (s *FormListService) DeleteForm(id int64, username string) error {
	// 检查表单是否存在
	var existForm formdo.Form
	err := s.db.Where("id = ? AND form_del_flag = ?", id, "0").First(&existForm).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("表单不存在")
		}
		return err
	}

	// 软删除表单
	now := time.Now()
	updates := map[string]interface{}{
		"form_del_flag":    "2",
		"form_update_time": now,
		"form_update_by":   username,
	}

	err = s.db.Model(&existForm).Updates(updates).Error
	if err != nil {
		return errors.New("删除表单失败: " + err.Error())
	}

	return nil
}
