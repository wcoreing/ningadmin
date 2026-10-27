package system

import (
	"errors"
	"time"
	"wnadmin/app/model/do/SystemDB"
	"wnadmin/app/model/dto"
)

// PopupService 弹窗服务
type PopupService struct{}

// GetPopupList 获取弹窗列表
func (s *PopupService) GetPopupList(query *dto.PopupQuery) ([]dto.PopupVO, int64, error) {
	popupDO := &SystemDB.PopupDO{}
	db := popupDO.Query()

	// 添加查询条件
	if query.Title != "" {
		db = db.Where("title LIKE ?", "%"+query.Title+"%")
	}
	if query.TriggerType != "" {
		db = db.Where("trigger_type = ?", query.TriggerType)
	}
	if query.Status != nil {
		db = db.Where("status = ?", *query.Status)
	}

	// 统计总数
	var total int64
	db.Count(&total)

	// 分页查询
	var popups []SystemDB.PopupDO
	offset := (query.Page - 1) * query.PageSize
	err := db.Offset(offset).Limit(query.PageSize).Order("priority desc, create_time desc").Find(&popups).Error
	if err != nil {
		return nil, 0, err
	}

	// 转换为VO
	var result []dto.PopupVO
	for _, popup := range popups {
		vo := dto.PopupVO{
			ID:           popup.ID,
			Title:        popup.Title,
			Content:      popup.Content,
			TitleStyle:   popup.GetTitleStyleMap(),
			Buttons:      popup.GetButtonsMap(),
			TriggerType:  popup.TriggerType,
			TriggerRule:  popup.GetTriggerRuleMap(),
			Priority:     popup.Priority,
			AutoClose:    popup.AutoClose,
			Delay:        popup.Delay,
			MaskClosable: popup.MaskClosable,
			StartTime:    popup.StartTime,
			EndTime:      popup.EndTime,
			Status:       popup.Status,
			ShowCount:    popup.ShowCount,
			ClickCount:   popup.ClickCount,
			CreateTime:   popup.CreateTime,
			UpdateTime:   popup.UpdateTime,
		}
		result = append(result, vo)
	}

	return result, total, nil
}

// CreatePopup 创建弹窗
func (s *PopupService) CreatePopup(req *dto.CreatePopupReq) error {
	now := time.Now()

	popup := &SystemDB.PopupDO{
		Title:        req.Title,
		Content:      req.Content,
		TriggerType:  req.TriggerType,
		Priority:     req.Priority,
		AutoClose:    req.AutoClose,
		Delay:        req.Delay,
		MaskClosable: req.MaskClosable,
		StartTime:    req.StartTime,
		EndTime:      req.EndTime,
		Status:       req.Status,
		CreateTime:   now,
		UpdateTime:   now,
	}

	// 设置JSON字段
	if err := popup.SetTitleStyle(req.TitleStyle); err != nil {
		return err
	}
	if err := popup.SetButtons(req.Buttons); err != nil {
		return err
	}
	if err := popup.SetTriggerRule(req.TriggerRule); err != nil {
		return err
	}

	// 设置默认值
	if popup.Status == 0 {
		popup.Status = 1
	}
	if popup.MaskClosable == 0 {
		popup.MaskClosable = 1
	}

	return popup.Create()
}

// UpdatePopup 更新弹窗
func (s *PopupService) UpdatePopup(req *dto.UpdatePopupReq) error {
	popup := &SystemDB.PopupDO{}
	err := popup.Query().Where("id = ?", req.ID).First(popup).Error
	if err != nil {
		return err
	}

	// 更新字段
	popup.Title = req.Title
	popup.Content = req.Content
	popup.TriggerType = req.TriggerType
	popup.Priority = req.Priority
	popup.AutoClose = req.AutoClose
	popup.Delay = req.Delay
	popup.MaskClosable = req.MaskClosable
	popup.StartTime = req.StartTime
	popup.EndTime = req.EndTime
	popup.Status = req.Status
	popup.UpdateTime = time.Now()

	// 设置JSON字段
	if err := popup.SetTitleStyle(req.TitleStyle); err != nil {
		return err
	}
	if err := popup.SetButtons(req.Buttons); err != nil {
		return err
	}
	if err := popup.SetTriggerRule(req.TriggerRule); err != nil {
		return err
	}

	return popup.Update()
}

// DeletePopup 删除弹窗
func (s *PopupService) DeletePopup(id int64) error {
	popup := &SystemDB.PopupDO{ID: id}
	return popup.Delete()
}

// GetPopupInfo 获取弹窗详情
func (s *PopupService) GetPopupInfo(id int64) (*dto.PopupVO, error) {
	popup := &SystemDB.PopupDO{ID: id}
	err := popup.First()
	if err != nil {
		return nil, err
	}

	vo := &dto.PopupVO{
		ID:           popup.ID,
		Title:        popup.Title,
		Content:      popup.Content,
		TitleStyle:   popup.GetTitleStyleMap(),
		Buttons:      popup.GetButtonsMap(),
		TriggerType:  popup.TriggerType,
		TriggerRule:  popup.GetTriggerRuleMap(),
		Priority:     popup.Priority,
		AutoClose:    popup.AutoClose,
		Delay:        popup.Delay,
		MaskClosable: popup.MaskClosable,
		StartTime:    popup.StartTime,
		EndTime:      popup.EndTime,
		Status:       popup.Status,
		ShowCount:    popup.ShowCount,
		ClickCount:   popup.ClickCount,
		CreateTime:   popup.CreateTime,
		UpdateTime:   popup.UpdateTime,
	}

	return vo, nil
}

// UpdatePopupStatus 更新弹窗状态
func (s *PopupService) UpdatePopupStatus(req *dto.PopupStatusReq) error {
	popup := &SystemDB.PopupDO{}
	err := popup.Query().Where("id = ?", req.ID).First(popup).Error
	if err != nil {
		return err
	}

	popup.Status = req.Status
	popup.UpdateTime = time.Now()

	return popup.Update()
}

// UpdatePopupStats 更新弹窗统计
func (s *PopupService) UpdatePopupStats(req *dto.PopupStatsReq) error {
	popup := &SystemDB.PopupDO{}
	err := popup.Query().Where("id = ?", req.ID).First(popup).Error
	if err != nil {
		return err
	}

	switch req.Type {
	case "show":
		popup.ShowCount++
	case "click":
		popup.ClickCount++
	default:
		return errors.New("无效的统计类型")
	}

	popup.UpdateTime = time.Now()
	return popup.Update()
}

// GetActivePopups 获取活跃的弹窗（前端调用）
func (s *PopupService) GetActivePopups(req *dto.GetActivePopupsReq) ([]dto.PopupVO, error) {
	popupDO := &SystemDB.PopupDO{}
	db := popupDO.Query().Where("status = ?", 1)

	// 时间范围过滤
	now := time.Now()
	db = db.Where("(start_time IS NULL OR start_time <= ?) AND (end_time IS NULL OR end_time >= ?)", now, now)

	// 触发类型过滤
	if req.TriggerType != "" {
		db = db.Where("trigger_type = ?", req.TriggerType)
	}

	var popups []SystemDB.PopupDO
	err := db.Order("priority desc, create_time desc").Find(&popups).Error
	if err != nil {
		return nil, err
	}

	// 转换为VO
	var result []dto.PopupVO
	for _, popup := range popups {
		vo := dto.PopupVO{
			ID:           popup.ID,
			Title:        popup.Title,
			Content:      popup.Content,
			TitleStyle:   popup.GetTitleStyleMap(),
			Buttons:      popup.GetButtonsMap(),
			TriggerType:  popup.TriggerType,
			TriggerRule:  popup.GetTriggerRuleMap(),
			Priority:     popup.Priority,
			AutoClose:    popup.AutoClose,
			Delay:        popup.Delay,
			MaskClosable: popup.MaskClosable,
			StartTime:    popup.StartTime,
			EndTime:      popup.EndTime,
			Status:       popup.Status,
			ShowCount:    popup.ShowCount,
			ClickCount:   popup.ClickCount,
			CreateTime:   popup.CreateTime,
			UpdateTime:   popup.UpdateTime,
		}
		result = append(result, vo)
	}

	return result, nil
}
