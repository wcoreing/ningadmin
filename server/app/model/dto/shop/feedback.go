package shop

import "time"

// FeedbackQuery 反馈查询参数
type FeedbackQuery struct {
	PageNum  int    `form:"pageNum"`
	PageSize int    `form:"pageSize"`
	Keyword  string `form:"keyword"`
	UserID   int64  `form:"userId"`
	Status   *int   `form:"status"`
}

// CreateFeedbackReq 创建反馈请求
type CreateFeedbackReq struct {
	UserID  int64  `json:"userId" binding:"required"`
	Content string `json:"content" binding:"required"`
	Contact string `json:"contact"`
}

// UpdateFeedbackReq 更新反馈请求
type UpdateFeedbackReq struct {
	ID        int64     `json:"id" binding:"required"`
	Status    int       `json:"status"`
	Reply     string    `json:"reply"`
	ReplyTime time.Time `json:"replyTime"`
}

// FeedbackVO 反馈视图对象
type FeedbackVO struct {
	ID         int64      `json:"id"`
	UserID     int64      `json:"userId"`
	Content    string     `json:"content"`
	Contact    string     `json:"contact"`
	Status     int        `json:"status"`
	StatusText string     `json:"statusText"`
	Reply      string     `json:"reply"`
	ReplyTime  *time.Time `json:"replyTime"`
	CreateTime time.Time  `json:"createTime"`
	UpdateTime time.Time  `json:"updateTime"`
} 