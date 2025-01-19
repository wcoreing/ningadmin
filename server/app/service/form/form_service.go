package form

import (
	"errors"
	"time"
	"wnadmin/app/common/utils"
	"wnadmin/app/model/do/form"
	"wnadmin/app/model/dto"
	"wnadmin/app/provider/db"

	"gorm.io/gorm"
)

type FormService struct {
	db *gorm.DB
}

func NewFormService() *FormService {
	return &FormService{
		db: db.DB.WnAdmin,
	}
}

// GetTemplateList 获取表单模板列表
func (s *FormService) GetTemplateList(query *dto.FormTemplateQuery) ([]form.FormTemplate, int64, error) {
	var templates []form.FormTemplate
	var total int64

	db := s.db.Model(&form.FormTemplate{})
	if query.Name != "" {
		db = db.Where("name LIKE ?", "%"+query.Name+"%")
	}

	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = db.Offset((query.PageNum - 1) * query.PageSize).
		Limit(query.PageSize).
		Find(&templates).Error
	if err != nil {
		return nil, 0, err
	}

	return templates, total, nil
}

// SaveTemplate 保存表单模板
func (s *FormService) SaveTemplate(template *form.FormTemplate) error {
	if template.ID > 0 {
		template.UpdateTime = utils.Ptr(time.Now())
		return s.db.Save(template).Error
	}
	template.CreateTime = time.Now()
	return s.db.Create(template).Error
}

// DeleteTemplate 删除表单模板
func (s *FormService) DeleteTemplate(id int64) error {
	// 检查是否有提交记录
	var count int64
	err := s.db.Model(&form.FormSubmission{}).Where("template_id = ?", id).Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("该模板已有提交记录，不能删除")
	}

	// 删除分享记录
	err = s.db.Where("template_id = ?", id).Delete(&form.FormShare{}).Error
	if err != nil {
		return err
	}

	// 删除模板
	return s.db.Delete(&form.FormTemplate{}, id).Error
}

// CreateShare 创建表单分享
func (s *FormService) CreateShare(share *form.FormShare) error {
	share.ShareCode = utils.GenerateUUID()
	share.CreateTime = time.Now()
	return s.db.Create(share).Error
}

// GetShareInfo 获取分享表单信息
func (s *FormService) GetShareInfo(shareCode string) (*form.FormTemplate, error) {
	var share form.FormShare
	err := s.db.Where("share_code = ?", shareCode).First(&share).Error
	if err != nil {
		return nil, err
	}

	if share.ExpireTime != nil && share.ExpireTime.Before(time.Now()) {
		return nil, errors.New("分享链接已过期")
	}

	var template form.FormTemplate
	err = s.db.First(&template, share.TemplateID).Error
	if err != nil {
		return nil, err
	}

	return &template, nil
}

// SubmitForm 提交表单
func (s *FormService) SubmitForm(submission *form.FormSubmission) error {
	submission.SubmitTime = time.Now()
	return s.db.Create(submission).Error
}

// GetSubmissionList 获取表单提交记录
func (s *FormService) GetSubmissionList(query *dto.FormSubmissionQuery) ([]form.FormSubmission, int64, error) {
	var submissions []form.FormSubmission
	var total int64

	db := s.db.Model(&form.FormSubmission{}).Where("template_id = ?", query.TemplateID)

	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = db.Offset((query.PageNum - 1) * query.PageSize).
		Limit(query.PageSize).
		Find(&submissions).Error
	if err != nil {
		return nil, 0, err
	}

	return submissions, total, nil
}
