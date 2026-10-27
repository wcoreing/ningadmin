package dto

// CmsContentQuery 内容查询
type CmsContentQuery struct {
	ChannelID int64  `form:"channelId"`
	Keyword   string `form:"keyword"`
	Status    string `form:"status"`
	Page      int    `form:"page"`
	PageSize  int    `form:"pageSize"`
}

// CmsContentReq 保存内容
type CmsContentReq struct {
	ID             int64  `json:"id"`
	ChannelID      int64  `json:"channelId" binding:"required"`
	Title          string `json:"title" binding:"required"`
	Summary        string `json:"summary"`
	Cover          string `json:"cover"`
	Body           string `json:"body"`
	SeoTitle       string `json:"seoTitle"`
	SeoKeywords    string `json:"seoKeywords"`
	SeoDescription string `json:"seoDescription"`
	Sort           int    `json:"sort"`
	Status         int    `json:"status"`
	PublishTime    string `json:"publishTime"`
}

// CmsContentItem 内容列表项
type CmsContentItem struct {
	ID             int64  `json:"id"`
	ChannelID      int64  `json:"channelId"`
	ChannelName    string `json:"channelName"`
	ChannelCode    string `json:"channelCode"`
	Title          string `json:"title"`
	Summary        string `json:"summary"`
	Cover          string `json:"cover"`
	Body           string `json:"body,omitempty"`
	SeoTitle       string `json:"seoTitle,omitempty"`
	SeoKeywords    string `json:"seoKeywords,omitempty"`
	SeoDescription string `json:"seoDescription,omitempty"`
	Sort           int    `json:"sort"`
	Status         int    `json:"status"`
	PublishTime    string `json:"publishTime"`
}
