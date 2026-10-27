package form

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"
	"wnadmin/app/common/utils"
	formdo "wnadmin/app/model/do/form"
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
func (s *FormService) GetTemplateList(query *dto.FormTemplateQuery) ([]formdo.FormTemplate, int64, error) {
	var templates []formdo.FormTemplate
	var total int64

	db := s.db.Model(&formdo.FormTemplate{})
	if query.Name != "" {
		db = db.Where("name LIKE ?", "%"+query.Name+"%")
	}

	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = db.Offset((query.PageNum - 1) * query.PageSize).
		Order("id desc").
		Limit(query.PageSize).
		Find(&templates).Error
	if err != nil {
		return nil, 0, err
	}

	return templates, total, nil
}

// CreateTemplate 创建表单模板
func (s *FormService) CreateTemplate(template *formdo.FormTemplate, username string) (*formdo.FormTemplate, error) {
	now := time.Now()
	template.CreateTime = now
	template.UpdateTime = utils.Ptr(now)
	template.Status = "0"
	template.DelFlag = "0"
	template.CreateBy = username
	err := s.db.Create(template).Error
	if err != nil {
		return nil, errors.New("创建表单模板失败: " + err.Error())
	}
	return template, nil
}

// UpdateTemplate 更新表单模板
func (s *FormService) UpdateTemplate(template *formdo.FormTemplate, username string) (*formdo.FormTemplate, error) {
	now := time.Now()
	updates := map[string]interface{}{
		"name":            template.Name,
		"description":     template.Description,
		"submit_url":      template.SubmitURL,
		"submit_method":   template.SubmitMethod,
		"success_message": template.SuccessMessage,
		"fail_message":    template.FailMessage,
		"form_items":      template.FormItems,
		"config":          template.Config,
		"update_time":     now,
		"update_by":       username,
	}

	err := s.db.Model(template).Updates(updates).Error
	if err != nil {
		return nil, errors.New("更新表单模板失败: " + err.Error())
	}

	// 重新查询完整数据
	err = s.db.First(template, template.ID).Error
	if err != nil {
		return nil, errors.New("获取更新后的表单模板失败: " + err.Error())
	}

	return template, nil
}

// SaveTemplate 保存表单模板
func (s *FormService) SaveTemplate(template *formdo.FormTemplate, username string) (*formdo.FormTemplate, error) {
	if template.ID == 0 {
		return s.CreateTemplate(template, username)
	}
	return s.UpdateTemplate(template, username)
}

// DeleteTemplate 删除表单模板
func (s *FormService) DeleteTemplate(id int64) error {
	// 检查是否有提交记录
	var count int64
	err := s.db.Model(&formdo.FormSubmission{}).Where("template_id = ?", id).Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("该模板已有提交记录，不能删除")
	}

	// 删除分享记录
	err = s.db.Where("template_id = ?", id).Delete(&formdo.FormShare{}).Error
	if err != nil {
		return err
	}

	// 删除模板
	return s.db.Delete(&formdo.FormTemplate{}, id).Error
}

// CreateShare 创建表单分享
func (s *FormService) CreateShare(share *formdo.FormShare) error {
	share.ShareCode = utils.GenerateUUID()
	share.CreateTime = time.Now()
	return s.db.Create(share).Error
}

// GetShareInfo 获取分享表单信息
func (s *FormService) GetShareInfo(shareCode string) (*formdo.FormTemplate, error) {
	var share formdo.FormShare
	err := s.db.Where("share_code = ?", shareCode).First(&share).Error
	if err != nil {
		return nil, err
	}

	if share.ExpireTime != nil && share.ExpireTime.Before(time.Now()) {
		return nil, errors.New("分享链接已过期")
	}

	var template formdo.FormTemplate
	err = s.db.First(&template, share.TemplateID).Error
	if err != nil {
		return nil, err
	}

	return &template, nil
}

// SubmitForm 提交表单
func (s *FormService) SubmitForm(submission *formdo.FormSubmission) error {
	submission.SubmitTime = time.Now()
	return s.db.Create(submission).Error
}

// GetSubmissionList 获取表单提交记录
func (s *FormService) GetSubmissionList(query *dto.FormSubmissionQuery) ([]formdo.FormSubmission, int64, error) {
	var submissions []formdo.FormSubmission
	var total int64

	db := s.db.Model(&formdo.FormSubmission{}).Where("template_id = ?", query.TemplateID)

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

// SaveTemplateFromRequest 从请求数据保存表单模板
func (s *FormService) SaveTemplateFromRequest(reqData map[string]interface{}, username string) (*formdo.FormTemplate, error) {
	template := &formdo.FormTemplate{
		Name:           reqData["name"].(string),
		Description:    reqData["description"].(string),
		SubmitURL:      reqData["submitUrl"].(string),
		SubmitMethod:   reqData["submitMethod"].(string),
		SuccessMessage: reqData["successMessage"].(string),
		FailMessage:    reqData["failMessage"].(string),
		Status:         "0",
		DelFlag:        "0",
	}

	// 处理ID
	if idStr, ok := reqData["id"].(string); ok && idStr != "" {
		if id, err := strconv.ParseInt(idStr, 10, 64); err == nil {
			template.ID = id
		}
	}

	// 处理JSON字段
	if formItems, ok := reqData["formItems"]; ok {
		formItemsBytes, _ := json.Marshal(formItems)
		template.FormItems = formItemsBytes
	}
	if config, ok := reqData["config"]; ok {
		configBytes, _ := json.Marshal(config)
		template.Config = configBytes
	}

	// 设置创建者/更新者
	if template.ID > 0 {
		template.UpdateBy = username
	} else {
		template.CreateBy = username
	}

	savedTemplate, err := s.SaveTemplate(template, username)
	return savedTemplate, err
}

// GetTemplateInfo 获取表单模板详情
func (s *FormService) GetTemplateInfo(id int64) (*formdo.FormTemplate, error) {
	var template formdo.FormTemplate
	err := s.db.First(&template, id).Error
	if err != nil {
		return nil, err
	}
	return &template, nil
}
