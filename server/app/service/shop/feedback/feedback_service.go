package feedback

import (
	"time"
	"wnadmin/app/model/do/shop"
	shopDto "wnadmin/app/model/dto/shop"
)

// Service 反馈服务
type Service struct{}

// NewService 创建服务实例
func NewService() *Service {
	return &Service{}
}

// GetFeedbackList 获取反馈列表
func (s *Service) GetFeedbackList(query *shopDto.FeedbackQuery) ([]shopDto.FeedbackVO, int64, error) {
	// 设置分页默认值
	if query.PageNum <= 0 {
		query.PageNum = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 10
	}

	// 构建查询条件
	conditions := make(map[string]interface{})
	if query.Keyword != "" {
		conditions["keyword"] = query.Keyword
	}
	if query.UserID > 0 {
		conditions["user_id"] = query.UserID
	}
	if query.Status != nil {
		conditions["status"] = *query.Status
	}

	// 获取数据
	feedbackModel := &shop.Feedback{}
	offset := (query.PageNum - 1) * query.PageSize
	feedbacks, total, err := feedbackModel.GetList(offset, query.PageSize, conditions)
	if err != nil {
		return nil, 0, err
	}

	// 转换为VO
	var feedbackVOs []shopDto.FeedbackVO
	for _, feedback := range feedbacks {
		feedbackVOs = append(feedbackVOs, s.convertToVO(&feedback))
	}

	return feedbackVOs, total, nil
}

// CreateFeedback 创建反馈
func (s *Service) CreateFeedback(req *shopDto.CreateFeedbackReq) error {
	feedback := &shop.Feedback{
		UserID:     req.UserID,
		Content:    req.Content,
		Contact:    req.Contact,
		Status:     0, // 默认待处理
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}

	return feedback.Create()
}

// UpdateFeedback 更新反馈
func (s *Service) UpdateFeedback(req *shopDto.UpdateFeedbackReq) error {
	feedback := &shop.Feedback{}
	if err := feedback.GetByID(req.ID); err != nil {
		return err
	}

	// 更新字段
	feedback.Status = req.Status
	feedback.Reply = req.Reply
	if !req.ReplyTime.IsZero() {
		feedback.ReplyTime = &req.ReplyTime
	}
	feedback.UpdateTime = time.Now()

	return feedback.Update()
}

// DeleteFeedback 删除反馈
func (s *Service) DeleteFeedback(id int64) error {
	feedback := &shop.Feedback{ID: id}
	return feedback.Delete()
}

// GetFeedbackInfo 获取反馈详情
func (s *Service) GetFeedbackInfo(id int64) (*shopDto.FeedbackVO, error) {
	feedback := &shop.Feedback{}
	if err := feedback.GetByID(id); err != nil {
		return nil, err
	}

	feedbackVO := s.convertToVO(feedback)
	return &feedbackVO, nil
}

// ReplyFeedback 回复反馈
func (s *Service) ReplyFeedback(id int64, reply string) error {
	feedback := &shop.Feedback{}
	if err := feedback.GetByID(id); err != nil {
		return err
	}

	now := time.Now()
	feedback.Reply = reply
	feedback.ReplyTime = &now
	feedback.Status = 1 // 已处理
	feedback.UpdateTime = now

	return feedback.Update()
}

// UpdateFeedbackStatus 更新反馈状态
func (s *Service) UpdateFeedbackStatus(id int64, status int) error {
	feedback := &shop.Feedback{}
	if err := feedback.GetByID(id); err != nil {
		return err
	}

	feedback.Status = status
	feedback.UpdateTime = time.Now()

	return feedback.Update()
}

// convertToVO 转换为VO对象
func (s *Service) convertToVO(feedback *shop.Feedback) shopDto.FeedbackVO {
	statusText := s.getStatusText(feedback.Status)

	return shopDto.FeedbackVO{
		ID:         feedback.ID,
		UserID:     feedback.UserID,
		Content:    feedback.Content,
		Contact:    feedback.Contact,
		Status:     feedback.Status,
		StatusText: statusText,
		Reply:      feedback.Reply,
		ReplyTime:  feedback.ReplyTime,
		CreateTime: feedback.CreateTime,
		UpdateTime: feedback.UpdateTime,
	}
}

// getStatusText 获取状态文本
func (s *Service) getStatusText(status int) string {
	switch status {
	case 0:
		return "待处理"
	case 1:
		return "已处理"
	case 2:
		return "已关闭"
	default:
		return "未知"
	}
}
