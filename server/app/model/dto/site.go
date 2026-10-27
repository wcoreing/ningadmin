package dto

// SiteGoodsQuery 官网商品查询
type SiteGoodsQuery struct {
	CategoryID int64  `form:"categoryId"`
	Keyword    string `form:"keyword"`
	Page       int    `form:"page"`
	PageSize   int    `form:"pageSize"`
}

// SiteArticleQuery 官网资讯查询
type SiteArticleQuery struct {
	Category string `form:"category"`
	Page     int    `form:"page"`
	PageSize int    `form:"pageSize"`
}

// SiteInquiryQuery 后台询盘查询
type SiteInquiryQuery struct {
	Status   string `form:"status"`
	Page     int    `form:"page"`
	PageSize int    `form:"pageSize"`
}

// SiteInquiryPage 询盘分页
type SiteInquiryPage struct {
	List  []SiteInquiryItem `json:"list"`
	Total int64             `json:"total"`
}

// SiteInquiryItem 询盘记录
type SiteInquiryItem struct {
	ID         int64  `json:"id"`
	GoodsName  string `json:"goodsName"`
	Name       string `json:"name"`
	Company    string `json:"company"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
	Country    string `json:"country"`
	Quantity   string `json:"quantity"`
	Message    string `json:"message"`
	Status     int    `json:"status"`
	CreateTime string `json:"createTime"`
}

// SiteInquiryStatusReq 更新询盘状态
type SiteInquiryStatusReq struct {
	Status int `json:"status"`
}

// SiteInquiryReq 提交询盘
type SiteInquiryReq struct {
	GoodsID  int64  `json:"goodsId"`
	Name     string `json:"name" binding:"required"`
	Company  string `json:"company"`
	Phone    string `json:"phone" binding:"required"`
	Email    string `json:"email"`
	Country  string `json:"country"`
	Quantity string `json:"quantity"`
	Message  string `json:"message" binding:"required"`
}

// SiteProfile 企业资料
type SiteProfile struct {
	Name           string `json:"name"`
	ShortName      string `json:"shortName"`
	Phone          string `json:"phone"`
	ContactName    string `json:"contactName"`
	Address        string `json:"address"`
	Email          string `json:"email"`
	Summary        string `json:"summary"`
	SeoTitle       string `json:"seoTitle"`
	SeoKeywords    string `json:"seoKeywords"`
	SeoDescription string `json:"seoDescription"`
}

// SiteBanner 首页轮播
type SiteBanner struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
	Image string `json:"image"`
	URL   string `json:"url"`
}

// SiteCategory 产品分类
type SiteCategory struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image"`
	Sort  int    `json:"sort"`
}

// SiteGoodsCard 产品卡片
type SiteGoodsCard struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	CategoryID   int64  `json:"categoryId"`
	CategoryName string `json:"categoryName"`
	Description  string `json:"description"`
	Cover        string `json:"cover"`
	Unit         string `json:"unit"`
}

// SiteGoodsDetail 产品详情
type SiteGoodsDetail struct {
	SiteGoodsCard
	SpecDesc    string   `json:"specDesc"`
	DetailDesc  string   `json:"detailDesc"`
	ServiceDesc string   `json:"serviceDesc"`
	Images      []string `json:"images"`
}

// SitePageCard 栏目摘要
type SitePageCard struct {
	Code    string `json:"code"`
	Title   string `json:"title"`
	Summary string `json:"summary"`
	Cover   string `json:"cover"`
}

// SitePageDetail 栏目正文
type SitePageDetail struct {
	SitePageCard
	Content        string `json:"content"`
	SeoTitle       string `json:"seoTitle"`
	SeoKeywords    string `json:"seoKeywords"`
	SeoDescription string `json:"seoDescription"`
}

// SiteArticleCard 资讯卡片
type SiteArticleCard struct {
	ID           int64  `json:"id"`
	Title        string `json:"title"`
	Summary      string `json:"summary"`
	Cover        string `json:"cover"`
	Category     string `json:"category"`
	CategoryName string `json:"categoryName"`
	PublishTime  string `json:"publishTime"`
}

// SiteArticleDetail 资讯正文
type SiteArticleDetail struct {
	SiteArticleCard
	Content        string `json:"content"`
	SeoTitle       string `json:"seoTitle"`
	SeoKeywords    string `json:"seoKeywords"`
	SeoDescription string `json:"seoDescription"`
}

// SiteHome 首页数据
type SiteHome struct {
	Profile    SiteProfile     `json:"profile"`
	Banners    []SiteBanner    `json:"banners"`
	Categories []SiteCategory  `json:"categories"`
	Featured   []SiteGoodsCard `json:"featured"`
	Pages      []SitePageCard  `json:"pages"`
	Articles   []SiteArticleCard `json:"articles"`
}

// SiteGoodsPage 产品分页
type SiteGoodsPage struct {
	List  []SiteGoodsCard `json:"list"`
	Total int64           `json:"total"`
}

// SiteArticlePage 资讯分页
type SiteArticlePage struct {
	List  []SiteArticleCard `json:"list"`
	Total int64             `json:"total"`
}
