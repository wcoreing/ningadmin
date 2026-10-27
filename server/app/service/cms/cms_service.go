package cms

import (
	"errors"
	"strings"
	"time"
	cmsmodel "wnadmin/app/model/do/cms"
	"wnadmin/app/model/dto"

	"gorm.io/gorm"
)

// Service 内容发布。栏目是内容分类，继承通用分类。
type Service struct{}

func (s *Service) contentQuery(q *dto.CmsContentQuery) *gorm.DB {
	query := (&cmsmodel.Content{}).Query()
	if q.ChannelID > 0 {
		query = query.Where("channel_id = ?", q.ChannelID)
	}
	if keyword := strings.TrimSpace(q.Keyword); keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("title LIKE ? OR summary LIKE ?", like, like)
	}
	if q.Status == "0" || q.Status == "1" {
		query = query.Where("status = ?", q.Status)
	}
	return query
}

// ContentList 内容列表
func (s *Service) ContentList(q *dto.CmsContentQuery) ([]dto.CmsContentItem, int64, error) {
	page, size := pageSize(q.Page, q.PageSize)
	var total int64
	if err := s.contentQuery(q).Count(&total).Error; err != nil {
		return nil, 0, errors.New("查询内容失败")
	}
	var rows []cmsmodel.Content
	err := s.contentQuery(q).Order("publish_time desc, id desc").Offset((page - 1) * size).Limit(size).Find(&rows).Error
	if err != nil {
		return nil, 0, errors.New("查询内容失败")
	}
	return s.toItems(rows, false), total, nil
}

// ContentInfo 内容详情
func (s *Service) ContentInfo(id int64) (*dto.CmsContentItem, error) {
	var row cmsmodel.Content
	if err := row.Query().Where("id = ?", id).First(&row).Error; err != nil {
		return nil, errors.New("内容不存在")
	}
	items := s.toItems([]cmsmodel.Content{row}, true)
	return &items[0], nil
}

// SaveContent 新增或修改内容
func (s *Service) SaveContent(req *dto.CmsContentReq) error {
	req.Title = strings.TrimSpace(req.Title)
	if req.Title == "" || req.ChannelID <= 0 {
		return errors.New("请选择栏目并填写标题")
	}
	var channel cmsmodel.Category
	if err := channel.Query().Where("id = ?", req.ChannelID).First(&channel).Error; err != nil {
		return errors.New("栏目不存在")
	}
	publishAt := time.Now()
	if strings.TrimSpace(req.PublishTime) != "" {
		parsed, err := time.ParseInLocation("2006-01-02 15:04:05", req.PublishTime, time.Local)
		if err != nil {
			parsed, err = time.ParseInLocation("2006-01-02", req.PublishTime, time.Local)
		}
		if err != nil {
			return errors.New("发布时间格式不正确")
		}
		publishAt = parsed
	}
	now := time.Now()
	if req.ID == 0 {
		row := &cmsmodel.Content{
			ChannelID:      req.ChannelID,
			Title:          req.Title,
			Summary:        req.Summary,
			Cover:          req.Cover,
			Body:           req.Body,
			SeoTitle:       req.SeoTitle,
			SeoKeywords:    req.SeoKeywords,
			SeoDescription: req.SeoDescription,
			Sort:           req.Sort,
			Status:         req.Status,
			PublishTime:    publishAt,
			CreateTime:     now,
			UpdateTime:     now,
		}
		if err := row.Create(); err != nil {
			return errors.New("保存内容失败")
		}
		return nil
	}
	var row cmsmodel.Content
	if err := row.Query().Where("id = ?", req.ID).First(&row).Error; err != nil {
		return errors.New("内容不存在")
	}
	row.ChannelID = req.ChannelID
	row.Title = req.Title
	row.Summary = req.Summary
	row.Cover = req.Cover
	row.Body = req.Body
	row.SeoTitle = req.SeoTitle
	row.SeoKeywords = req.SeoKeywords
	row.SeoDescription = req.SeoDescription
	row.Sort = req.Sort
	row.Status = req.Status
	row.PublishTime = publishAt
	row.UpdateTime = now
	if err := row.Save(); err != nil {
		return errors.New("保存内容失败")
	}
	return nil
}

// DeleteContent 删除内容
func (s *Service) DeleteContent(id int64) error {
	var row cmsmodel.Content
	if err := row.Query().Where("id = ?", id).First(&row).Error; err != nil {
		return errors.New("内容不存在")
	}
	row.IsDelete = 1
	row.UpdateTime = time.Now()
	if err := row.Save(); err != nil {
		return errors.New("删除失败")
	}
	return nil
}

// PublishedPage 已发布的单页。
func (s *Service) PublishedPage(code string) (*cmsmodel.Category, *cmsmodel.Content, error) {
	var channel cmsmodel.Category
	err := channel.Query().
		Where("status = ? AND JSON_UNQUOTE(JSON_EXTRACT(extra, '$.code')) = ? AND JSON_UNQUOTE(JSON_EXTRACT(extra, '$.type')) = ?", 1, code, cmsmodel.TypePage).
		First(&channel).Error
	if err != nil {
		return nil, nil, errors.New("栏目不存在")
	}
	var content cmsmodel.Content
	err = content.Query().Where("channel_id = ? AND status = ?", channel.ID, 1).
		Order("publish_time desc, id desc").First(&content).Error
	if err != nil {
		return &channel, nil, nil
	}
	return &channel, &content, nil
}

// PageChannels 对外展示的单页：启用的顶级单页，隐私说明只在页脚进入。
func (s *Service) PageChannels() ([]cmsmodel.Category, error) {
	rows, err := s.cmsCategories()
	if err != nil {
		return nil, err
	}
	list := make([]cmsmodel.Category, 0)
	for _, row := range rows {
		attr := row.Attr()
		if row.Status != 1 || row.ParentId != 0 || attr.Type != cmsmodel.TypePage || attr.Code == "privacy" {
			continue
		}
		list = append(list, row)
	}
	return list, nil
}

// PublishedContents 某个列表栏目及其子栏目下已发布的内容
func (s *Service) PublishedContents(rootCode, childCode string, page, size int) ([]dto.CmsContentItem, int64, error) {
	ids, channels, err := s.listChannelIDs(rootCode, childCode)
	if err != nil {
		return nil, 0, err
	}
	if len(ids) == 0 {
		return []dto.CmsContentItem{}, 0, nil
	}
	page, size = pageSize(page, size)
	query := (&cmsmodel.Content{}).Query().Where("channel_id IN ? AND status = ?", ids, 1)
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, errors.New("查询内容失败")
	}
	var rows []cmsmodel.Content
	err = (&cmsmodel.Content{}).Query().
		Where("channel_id IN ? AND status = ?", ids, 1).
		Order("publish_time desc, id desc").
		Offset((page - 1) * size).
		Limit(size).
		Find(&rows).Error
	if err != nil {
		return nil, 0, errors.New("查询内容失败")
	}
	return fillChannels(rows, channels, false), total, nil
}

// PublishedContent 已发布内容
func (s *Service) PublishedContent(id int64) (*dto.CmsContentItem, error) {
	var row cmsmodel.Content
	if err := row.Query().Where("id = ? AND status = ?", id, 1).First(&row).Error; err != nil {
		return nil, errors.New("内容不存在")
	}
	channels, _ := s.channelMap([]int64{row.ChannelID})
	items := fillChannels([]cmsmodel.Content{row}, channels, true)
	return &items[0], nil
}

func (s *Service) cmsCategories() ([]cmsmodel.Category, error) {
	var rows []cmsmodel.Category
	err := (&cmsmodel.Category{}).Query().Order("sort desc, id asc").Find(&rows).Error
	if err != nil {
		return nil, errors.New("查询栏目失败")
	}
	if rows == nil {
		rows = []cmsmodel.Category{}
	}
	return rows, nil
}

func (s *Service) listChannelIDs(rootCode, childCode string) ([]int64, map[int64]cmsmodel.Category, error) {
	all, err := s.cmsCategories()
	if err != nil {
		return nil, nil, err
	}
	byID := map[int64]cmsmodel.Category{}
	byCode := map[string]cmsmodel.Category{}
	for _, row := range all {
		if row.Status != 1 {
			continue
		}
		byID[row.ID] = row
		byCode[row.Attr().Code] = row
	}
	code := childCode
	if code == "" {
		code = rootCode
	}
	root, ok := byCode[code]
	if !ok {
		return nil, byID, nil
	}
	ids := []int64{root.ID}
	if childCode == "" {
		for _, row := range all {
			if row.Status == 1 && row.ParentId == root.ID {
				ids = append(ids, row.ID)
			}
		}
	}
	return ids, byID, nil
}

func (s *Service) toItems(rows []cmsmodel.Content, withBody bool) []dto.CmsContentItem {
	ids := make([]int64, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, row.ChannelID)
	}
	channels, _ := s.channelMap(ids)
	return fillChannels(rows, channels, withBody)
}

func (s *Service) channelMap(ids []int64) (map[int64]cmsmodel.Category, error) {
	result := map[int64]cmsmodel.Category{}
	if len(ids) == 0 {
		return result, nil
	}
	var rows []cmsmodel.Category
	err := (&cmsmodel.Category{}).Query().Where("id IN ?", ids).Find(&rows).Error
	for _, row := range rows {
		result[row.ID] = row
	}
	return result, err
}

func fillChannels(rows []cmsmodel.Content, channels map[int64]cmsmodel.Category, withBody bool) []dto.CmsContentItem {
	list := make([]dto.CmsContentItem, 0, len(rows))
	for _, row := range rows {
		item := dto.CmsContentItem{
			ID:          row.ID,
			ChannelID:   row.ChannelID,
			Title:       row.Title,
			Summary:     row.Summary,
			Cover:       row.Cover,
			Sort:        row.Sort,
			Status:      row.Status,
			PublishTime: row.PublishTime.Format("2006-01-02 15:04:05"),
		}
		if channel, ok := channels[row.ChannelID]; ok {
			item.ChannelName = channel.Name
			item.ChannelCode = channel.Attr().Code
		}
		if withBody {
			item.Body = row.Body
			item.SeoTitle = row.SeoTitle
			item.SeoKeywords = row.SeoKeywords
			item.SeoDescription = row.SeoDescription
		}
		list = append(list, item)
	}
	return list
}

func pageSize(page, size int) (int, int) {
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 10
	}
	return page, size
}
