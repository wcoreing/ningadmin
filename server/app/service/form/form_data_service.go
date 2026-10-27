package form

import (
	"errors"
	"time"
	formdo "wnadmin/app/model/do/form"
	"wnadmin/app/model/dto"
	"wnadmin/app/provider/db"

	"gorm.io/gorm"
)

// FormDataService 表单数据服务
type FormDataService struct {
	db *gorm.DB
}

// NewFormDataService 创建表单数据服务实例
func NewFormDataService() *FormDataService {
	return &FormDataService{
		db: db.DB.WnAdmin,
	}
}

// GetFormDataList 获取表单数据列表
func (s *FormDataService) GetFormDataList(query *dto.FormDataQuery) ([]dto.FormDataResponse, int64, error) {
	var list []dto.FormDataResponse
	var total int64

	// 构建查询
	db := s.db.Table("form_data fd").
		Select("fd.*, f.form_name").
		Joins("LEFT JOIN form f ON fd.form_id = f.id")

	// 条件查询
	if query.FormName != "" {
		db = db.Where("f.form_name LIKE ?", "%"+query.FormName+"%")
	}
	if !query.StartTime.IsZero() {
		db = db.Where("fd.create_time >= ?", query.StartTime)
	}
	if !query.EndTime.IsZero() {
		db = db.Where("fd.create_time <= ?", query.EndTime)
	}

	// 统计总数
	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	err = db.Offset((query.PageNum - 1) * query.PageSize).
		Limit(query.PageSize).
		Order("fd.id DESC").
		Find(&list).Error
	if err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

// GetFormDataDetail 获取表单数据详情
func (s *FormDataService) GetFormDataDetail(id int64) (*dto.FormDataResponse, error) {
	var data dto.FormDataResponse

	err := s.db.Table("form_data fd").
		Select("fd.*, f.form_name").
		Joins("LEFT JOIN form f ON fd.form_id = f.id").
		Where("fd.id = ?", id).
		First(&data).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("表单数据不存在")
		}
		return nil, err
	}

	return &data, nil
}

// DeleteFormData 删除表单数据
func (s *FormDataService) DeleteFormData(id int64) error {
	// 检查数据是否存在
	var count int64
	err := s.db.Model(&formdo.FormData{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("表单数据不存在")
	}

	// 删除数据
	err = s.db.Delete(&formdo.FormData{}, id).Error
	if err != nil {
		return errors.New("删除表单数据失败")
	}

	return nil
}

// SubmitForm 提交表单
func (s *FormDataService) SubmitForm(req *dto.FormSubmitRequest, ip string, username string) error {
	// 检查表单是否存在且可用
	var form formdo.Form
	err := s.db.Where("id = ? AND form_status = '0' AND form_del_flag = '0'", req.FormID).First(&form).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("表单不存在或已停用")
		}
		return err
	}

	// 创建表单数据记录
	formData := &formdo.FormData{
		FormID:     req.FormID,
		Data:       req.Data,
		Submitter:  req.Submitter,
		IP:         ip,
		Status:     "0", // 0-成功
		CreateBy:   username,
		CreateTime: time.Now(),
	}

	// 保存数据
	err = s.db.Create(formData).Error
	if err != nil {
		return errors.New("提交表单失败: " + err.Error())
	}

	return nil
}
