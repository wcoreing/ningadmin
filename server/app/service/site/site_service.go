package site

import (
	"errors"
	"strings"
	"wnadmin/app/common/mail"
	"wnadmin/app/model/do/shop"
	sitemodel "wnadmin/app/model/do/site"
	"wnadmin/app/model/dto"
	cmsservice "wnadmin/app/service/cms"

	"gorm.io/gorm"
)

// Service 官网读取与询盘
type Service struct {
	cms *cmsservice.Service
}

// Home 首页聚合
func (s *Service) Home() (*dto.SiteHome, error) {
	profile, err := s.Profile()
	if err != nil {
		return nil, err
	}
	categories, err := s.Categories()
	if err != nil {
		return nil, err
	}
	featured, err := s.goodsCards(0, "", 1, 4)
	if err != nil {
		return nil, err
	}
	articles, err := s.articleCards("", 1, 3)
	if err != nil {
		return nil, err
	}
	pages, err := s.pageCards()
	if err != nil {
		return nil, err
	}
	banners, err := s.banners()
	if err != nil {
		return nil, err
	}
	return &dto.SiteHome{
		Profile:    *profile,
		Banners:    banners,
		Categories: categories,
		Featured:   featured.List,
		Pages:      pages,
		Articles:   articles.List,
	}, nil
}

// Profile 企业资料
func (s *Service) Profile() (*dto.SiteProfile, error) {
	var row sitemodel.Profile
	err := row.Query().Order("id asc").First(&row).Error
	if err != nil {
		return nil, errors.New("企业资料不存在")
	}
	return &dto.SiteProfile{
		Name:           row.Name,
		ShortName:      row.ShortName,
		Phone:          row.Phone,
		ContactName:    row.ContactName,
		Address:        row.Address,
		Email:          row.Email,
		Summary:        row.Summary,
		SeoTitle:       row.SeoTitle,
		SeoKeywords:    row.SeoKeywords,
		SeoDescription: row.SeoDescription,
	}, nil
}

// Page 单页栏目
func (s *Service) Page(code string) (*dto.SitePageDetail, error) {
	channel, content, err := s.cmsService().PublishedPage(code)
	if err != nil {
		return nil, err
	}
	attr := channel.Attr()
	detail := &dto.SitePageDetail{
		SitePageCard: dto.SitePageCard{
			Code:    attr.Code,
			Title:   channel.Name,
			Summary: attr.Summary,
			Cover:   channel.Image,
		},
		SeoTitle:       attr.SeoTitle,
		SeoKeywords:    attr.SeoKeywords,
		SeoDescription: attr.SeoDescription,
	}
	if content != nil {
		if content.Title != "" {
			detail.Title = content.Title
		}
		if content.Summary != "" {
			detail.Summary = content.Summary
		}
		if content.Cover != "" {
			detail.Cover = content.Cover
		}
		detail.Content = content.Body
		if content.SeoTitle != "" {
			detail.SeoTitle = content.SeoTitle
		}
		if content.SeoKeywords != "" {
			detail.SeoKeywords = content.SeoKeywords
		}
		if content.SeoDescription != "" {
			detail.SeoDescription = content.SeoDescription
		}
	}
	return detail, nil
}

// Categories 产品分类
func (s *Service) Categories() ([]dto.SiteCategory, error) {
	var rows []shop.Category
	err := (&shop.Category{}).Query().
		Where("module = ? AND status = ?", "goods", 1).
		Order("sort desc, id asc").
		Find(&rows).Error
	if err != nil {
		return nil, errors.New("查询分类失败")
	}
	list := make([]dto.SiteCategory, 0, len(rows))
	for _, row := range rows {
		list = append(list, dto.SiteCategory{
			ID:    row.ID,
			Name:  row.Name,
			Image: row.Image,
			Sort:  row.Sort,
		})
	}
	return list, nil
}

// Goods 产品列表
func (s *Service) Goods(q *dto.SiteGoodsQuery) (*dto.SiteGoodsPage, error) {
	page, size := normalizePage(q.Page, q.PageSize)
	return s.goodsCards(q.CategoryID, strings.TrimSpace(q.Keyword), page, size)
}

// GoodsDetail 产品详情
func (s *Service) GoodsDetail(id int64) (*dto.SiteGoodsDetail, error) {
	var goods shop.GoodsDO
	err := (&shop.GoodsDO{Status: 1}).Query().Where("id = ?", id).First(&goods).Error
	if err != nil {
		return nil, errors.New("产品不存在")
	}
	names := s.categoryNames([]int64{goods.CategoryId})
	covers := s.covers([]int64{goods.Id})
	detail := &dto.SiteGoodsDetail{
		SiteGoodsCard: dto.SiteGoodsCard{
			ID:           goods.Id,
			Name:         goods.Name,
			CategoryID:   goods.CategoryId,
			CategoryName: names[goods.CategoryId],
			Description:  goods.Description,
			Cover:        covers[goods.Id],
			Unit:         goods.Unit,
		},
		Images: []string{},
	}
	if detail.Cover != "" {
		detail.Images = append(detail.Images, detail.Cover)
	}
	var extra shop.ShopGoodsDetail
	if err := extra.Query().Where("goods_id = ?", goods.Id).First(&extra).Error; err == nil {
		detail.SpecDesc = extra.SpecDesc
		detail.DetailDesc = extra.DetailDesc
		detail.ServiceDesc = extra.ServiceDesc
	}
	var images []shop.ShopGoodsImage
	_ = (&shop.ShopGoodsImage{}).Query().
		Where("goods_id = ? AND is_delete = 0 AND position <> 1", goods.Id).
		Order("sort desc, id asc").
		Find(&images).Error
	for _, image := range images {
		detail.Images = append(detail.Images, image.Url)
	}
	return detail, nil
}

// Articles 资讯列表
func (s *Service) Articles(q *dto.SiteArticleQuery) (*dto.SiteArticlePage, error) {
	page, size := normalizePage(q.Page, q.PageSize)
	return s.articleCards(q.Category, page, size)
}

// Article 资讯详情
func (s *Service) Article(id int64) (*dto.SiteArticleDetail, error) {
	item, err := s.cmsService().PublishedContent(id)
	if err != nil {
		return nil, err
	}
	card := toArticleCard(*item)
	return &dto.SiteArticleDetail{
		SiteArticleCard: card,
		Content:         item.Body,
		SeoTitle:        item.SeoTitle,
		SeoKeywords:     item.SeoKeywords,
		SeoDescription:  item.SeoDescription,
	}, nil
}

// CreateInquiry 保存询盘
func (s *Service) CreateInquiry(req *dto.SiteInquiryReq) error {
	req.Name = strings.TrimSpace(req.Name)
	req.Phone = strings.TrimSpace(req.Phone)
	req.Message = strings.TrimSpace(req.Message)
	if req.Name == "" || req.Phone == "" || req.Message == "" {
		return errors.New("请填写联系人、电话和需求说明")
	}
	goodsName := ""
	if req.GoodsID > 0 {
		var goods shop.GoodsDO
		err := (&shop.GoodsDO{Status: 1}).Query().Where("id = ?", req.GoodsID).First(&goods).Error
		if err != nil {
			return errors.New("产品不存在")
		}
		goodsName = goods.Name
	}
	row := &sitemodel.Inquiry{
		GoodsID:   req.GoodsID,
		GoodsName: goodsName,
		Name:      req.Name,
		Company:   strings.TrimSpace(req.Company),
		Phone:     req.Phone,
		Email:     strings.TrimSpace(req.Email),
		Country:   strings.TrimSpace(req.Country),
		Quantity:  strings.TrimSpace(req.Quantity),
		Message:   req.Message,
		Status:    0,
	}
	if err := row.Create(); err != nil {
		return errors.New("提交失败")
	}
	s.notifyInquiry(row)
	return nil
}

// Inquiries 后台询盘列表
func (s *Service) Inquiries(q *dto.SiteInquiryQuery) (*dto.SiteInquiryPage, error) {
	page, size := normalizePage(q.Page, q.PageSize)
	rows, total, err := (&sitemodel.Inquiry{}).Page(q.Status, page, size)
	if err != nil {
		return nil, errors.New("查询询盘失败")
	}
	list := make([]dto.SiteInquiryItem, 0, len(rows))
	for _, row := range rows {
		list = append(list, dto.SiteInquiryItem{
			ID:         row.ID,
			GoodsName:  row.GoodsName,
			Name:       row.Name,
			Company:    row.Company,
			Phone:      row.Phone,
			Email:      row.Email,
			Country:    row.Country,
			Quantity:   row.Quantity,
			Message:    row.Message,
			Status:     row.Status,
			CreateTime: row.CreateTime.Format("2006-01-02 15:04"),
		})
	}
	return &dto.SiteInquiryPage{List: list, Total: total}, nil
}

// UpdateInquiryStatus 标记询盘是否已联系
func (s *Service) UpdateInquiryStatus(id int64, status int) error {
	if id <= 0 || (status != 0 && status != 1) {
		return errors.New("状态不正确")
	}
	if err := (&sitemodel.Inquiry{}).Mark(id, status); err != nil {
		return errors.New("更新失败")
	}
	return nil
}

// SaveProfile 保存企业资料和询盘收件邮箱
func (s *Service) SaveProfile(req *dto.SiteProfile) error {
	req.Name = strings.TrimSpace(req.Name)
	req.Phone = strings.TrimSpace(req.Phone)
	if req.Name == "" || req.Phone == "" {
		return errors.New("请填写公司名称和电话")
	}
	var row sitemodel.Profile
	if err := row.Query().Order("id asc").First(&row).Error; err != nil {
		return errors.New("企业资料不存在")
	}
	err := row.UpdateFields(map[string]any{
		"name":            req.Name,
		"short_name":      strings.TrimSpace(req.ShortName),
		"phone":           req.Phone,
		"contact_name":    strings.TrimSpace(req.ContactName),
		"address":         strings.TrimSpace(req.Address),
		"email":           strings.TrimSpace(req.Email),
		"summary":         strings.TrimSpace(req.Summary),
		"seo_title":       strings.TrimSpace(req.SeoTitle),
		"seo_keywords":    strings.TrimSpace(req.SeoKeywords),
		"seo_description": strings.TrimSpace(req.SeoDescription),
	})
	if err != nil {
		return errors.New("保存失败")
	}
	return nil
}

func (s *Service) notifyInquiry(row *sitemodel.Inquiry) {
	var profile sitemodel.Profile
	if err := profile.Query().Order("id asc").First(&profile).Error; err != nil {
		return
	}
	body := "联系人：" + row.Name + "\n电话：" + row.Phone + "\n公司：" + row.Company +
		"\n邮箱：" + row.Email + "\n市场：" + row.Country + "\n数量：" + row.Quantity +
		"\n产品：" + row.GoodsName + "\n需求：\n" + row.Message
	_ = mail.Send(profile.Email, "集亿玖泽新询盘", body)
}

func (s *Service) banners() ([]dto.SiteBanner, error) {
	var rows []shop.ShopBanner
	err := (&shop.ShopBanner{}).Query().
		Where("position = ? AND status = ?", shop.BannerPositionHome, shop.BannerStatusEnable).
		Order("sort desc, id desc").
		Find(&rows).Error
	if err != nil {
		return nil, errors.New("查询轮播失败")
	}
	list := make([]dto.SiteBanner, 0, len(rows))
	for _, row := range rows {
		list = append(list, dto.SiteBanner{
			ID:    row.Id,
			Title: row.Title,
			Image: row.Image,
			URL:   row.Url,
		})
	}
	return list, nil
}

func (s *Service) pageCards() ([]dto.SitePageCard, error) {
	rows, err := s.cmsService().PageChannels()
	if err != nil {
		return nil, err
	}
	list := make([]dto.SitePageCard, 0, len(rows))
	for _, row := range rows {
		attr := row.Attr()
		card := dto.SitePageCard{
			Code:    attr.Code,
			Title:   row.Name,
			Summary: attr.Summary,
			Cover:   row.Image,
		}
		_, content, _ := s.cmsService().PublishedPage(attr.Code)
		if content != nil {
			card.Summary = content.Summary
			if content.Cover != "" {
				card.Cover = content.Cover
			}
		}
		list = append(list, card)
	}
	return list, nil
}

func (s *Service) goodsQuery(categoryID int64, keyword string) *gorm.DB {
	query := (&shop.GoodsDO{Status: 1}).Query()
	if categoryID > 0 {
		query = query.Where("category_id = ?", categoryID)
	}
	if keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("name LIKE ? OR description LIKE ?", like, like)
	}
	return query
}

func (s *Service) goodsCards(categoryID int64, keyword string, page, size int) (*dto.SiteGoodsPage, error) {
	var total int64
	if err := s.goodsQuery(categoryID, keyword).Count(&total).Error; err != nil {
		return nil, errors.New("查询产品失败")
	}
	var rows []shop.GoodsDO
	err := s.goodsQuery(categoryID, keyword).Order("sort desc, id desc").Offset((page - 1) * size).Limit(size).Find(&rows).Error
	if err != nil {
		return nil, errors.New("查询产品失败")
	}
	ids := make([]int64, 0, len(rows))
	categoryIDs := make([]int64, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, row.Id)
		categoryIDs = append(categoryIDs, row.CategoryId)
	}
	names := s.categoryNames(categoryIDs)
	covers := s.covers(ids)
	list := make([]dto.SiteGoodsCard, 0, len(rows))
	for _, row := range rows {
		list = append(list, dto.SiteGoodsCard{
			ID:           row.Id,
			Name:         row.Name,
			CategoryID:   row.CategoryId,
			CategoryName: names[row.CategoryId],
			Description:  row.Description,
			Cover:        covers[row.Id],
			Unit:         row.Unit,
		})
	}
	return &dto.SiteGoodsPage{List: list, Total: total}, nil
}

func (s *Service) articleCards(category string, page, size int) (*dto.SiteArticlePage, error) {
	rows, total, err := s.cmsService().PublishedContents("news", category, page, size)
	if err != nil {
		return nil, err
	}
	list := make([]dto.SiteArticleCard, 0, len(rows))
	for _, row := range rows {
		list = append(list, toArticleCard(row))
	}
	return &dto.SiteArticlePage{List: list, Total: total}, nil
}

func (s *Service) cmsService() *cmsservice.Service {
	if s.cms == nil {
		s.cms = &cmsservice.Service{}
	}
	return s.cms
}

func (s *Service) categoryNames(ids []int64) map[int64]string {
	names := map[int64]string{}
	if len(ids) == 0 {
		return names
	}
	var rows []shop.Category
	_ = (&shop.Category{}).Query().Where("id IN ?", ids).Find(&rows).Error
	for _, row := range rows {
		names[row.ID] = row.Name
	}
	return names
}

func (s *Service) covers(goodsIDs []int64) map[int64]string {
	covers := map[int64]string{}
	if len(goodsIDs) == 0 {
		return covers
	}
	var rows []shop.ShopGoodsImage
	_ = (&shop.ShopGoodsImage{}).Query().
		Where("goods_id IN ? AND is_delete = 0", goodsIDs).
		Order("position desc, sort desc, id asc").
		Find(&rows).Error
	for _, row := range rows {
		if _, ok := covers[row.GoodsId]; ok && row.Position != shop.ImagePositionMain {
			continue
		}
		if row.Position == shop.ImagePositionMain || covers[row.GoodsId] == "" {
			covers[row.GoodsId] = row.Url
		}
	}
	return covers
}

func toArticleCard(row dto.CmsContentItem) dto.SiteArticleCard {
	publish := row.PublishTime
	if len(publish) >= 10 {
		publish = publish[:10]
	}
	return dto.SiteArticleCard{
		ID:           row.ID,
		Title:        row.Title,
		Summary:      row.Summary,
		Cover:        row.Cover,
		Category:     row.ChannelCode,
		CategoryName: row.ChannelName,
		PublishTime:  publish,
	}
}

func normalizePage(page, size int) (int, int) {
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 48 {
		size = 12
	}
	return page, size
}
