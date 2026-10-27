package file

import (
	"errors"
	"mime/multipart"
	"path"
	"strings"
	"time"
	"wnadmin/app/common/utils"
	"wnadmin/app/config"
	"wnadmin/app/model/do/SystemDB"
	"wnadmin/app/model/dto/common"
	"wnadmin/app/provider/db"
)

// Service 文件服务
type Service struct{}

// NewService 创建文件服务实例
func NewService() *Service {
	return &Service{}
}

// Upload 上传文件
func (s *Service) Upload(file *multipart.FileHeader, categoryId string) (*common.FileVO, error) {
	// 1. 验证文件类型
	if !utils.IsValidFileType(file.Filename) {
		return nil, errors.New("不支持的文件类型")
	}

	// 2. 验证文件大小
	if !utils.IsValidFileSize(file.Size) {
		return nil, errors.New("文件大小超出限制")
	}

	var relativePath, actualPath string
	var storageType string

	// 3. 根据配置选择上传方式
	if strings.ToLower(config.App.Upload.Type) == "oss" {
		// OSS上传
		ossUploader, err := utils.NewOSSUploader()
		if err != nil {
			return nil, err
		}
		relativePath, err = ossUploader.UploadFile(file)
		if err != nil {
			return nil, err
		}
		actualPath = strings.TrimPrefix(relativePath, "/") // OSS对象键
		storageType = "oss"
	} else {
		// 本地上传
		ext := path.Ext(file.Filename)
		fileName := utils.GenerateUUID() + ext
		actualPath = "uploads/" + fileName
		relativePath = "/" + actualPath

		if err := utils.SaveUploadedFile(file, actualPath); err != nil {
			return nil, err
		}
		storageType = "local"
	}

	// 4. 保存文件信息到数据库
	now := time.Now()

	fileModel := &SystemDB.File{
		CategoryId:   utils.StrToInt64(categoryId),
		Name:         path.Base(actualPath),
		OriginalName: file.Filename,
		Url:          relativePath, // 统一存储相对路径
		Path:         actualPath,   // 存储实际路径（本地路径或OSS对象键）
		Size:         file.Size,
		Type:         file.Header.Get("Content-Type"),
		StorageType:  storageType, // 存储类型
		Status:       1,
		CreateTime:   now,
		UpdateTime:   now,
	}

	if err := fileModel.Create(); err != nil {
		return nil, err
	}

	return s.convertToVO(fileModel), nil
}

// GetList 获取文件列表
func (s *Service) GetList(query *common.FileQuery) ([]common.FileVO, int64, error) {
	fileDO := &SystemDB.File{}
	db := fileDO.Query()

	if query.CategoryId > 0 {
		db = db.Where("category_id = ?", query.CategoryId)
	}
	if query.Keyword != "" {
		db = db.Where("name LIKE ?", "%"+query.Keyword+"%")
	}

	var total int64
	db.Count(&total)

	var files []SystemDB.File
	offset := (query.Page - 1) * query.PageSize
	err := db.Offset(offset).Limit(query.PageSize).Order("create_time DESC").Find(&files).Error
	if err != nil {
		return nil, 0, err
	}

	items := make([]common.FileVO, 0)
	for _, file := range files {
		items = append(items, *s.convertToVO(&file))
	}

	return items, total, nil
}

// Delete 删除文件
func (s *Service) Delete(id string) error {
	file := &SystemDB.File{}
	file.Id = utils.StrToInt64(id)

	if err := file.First(); err != nil {
		return err
	}

	// 1. 删除物理文件
	if file.StorageType == "oss" {
		// OSS删除
		ossUploader, err := utils.NewOSSUploader()
		if err != nil {
			return err
		}
		if err := ossUploader.DeleteFile(file.Url); err != nil {
			return err
		}
	} else {
		// 本地删除
		if err := utils.RemoveFile(file.Path); err != nil {
			return err
		}
	}

	// 2. 删除数据库记录
	return file.Delete()
}

// GetCategoryTree 获取文件分类树
func (s *Service) GetCategoryTree() ([]map[string]interface{}, error) {
	categoryDO := &SystemDB.FileCategory{}
	var categories []SystemDB.FileCategory
	var totalCount int64

	// 先获取所有文件的总数
	if err := (&SystemDB.File{}).Query().
		Where("status = ?", 1).
		Count(&totalCount).Error; err != nil {
		return nil, err
	}

	err := categoryDO.Query().Order("sort asc").Find(&categories).Error
	if err != nil {
		return nil, err
	}

	// 构建分类ID到文件数量的映射
	categoryCountMap := make(map[int64]int64)
	// 获取每个分类的文件数量
	for _, category := range categories {
		var count int64
		if err := (&SystemDB.File{}).Query().
			Where("category_id = ? AND status = ?", category.Id, 1).
			Count(&count).Error; err != nil {
			return nil, err
		}
		categoryCountMap[category.Id] = count
	}

	// 构建树形结构
	var buildTree func(parentId int64) []map[string]interface{}
	buildTree = func(parentId int64) []map[string]interface{} {
		var nodes []map[string]interface{}
		for _, category := range categories {
			if category.ParentId == parentId {
				node := map[string]interface{}{
					"id":        category.Id,
					"name":      category.Name,
					"parent_id": category.ParentId,
					"count":     categoryCountMap[category.Id],
					"icon":      category.Icon,
					"sort":      category.Sort,
				}
				children := buildTree(category.Id)
				if len(children) > 0 {
					node["children"] = children
				}
				nodes = append(nodes, node)
			}
		}
		return nodes
	}

	// 生成树形结构
	result := buildTree(0)

	// 添加"全部文件"根节点，包含总数
	allFiles := map[string]interface{}{
		"id":        "all",
		"name":      "全部文件",
		"parent_id": 0,
		"count":     totalCount,
		"icon":      "Folder",
		"sort":      0,
	}

	// 将"全部文件"节点放在最前面
	return append([]map[string]interface{}{allFiles}, result...), nil
}

// convertToVO 转换为VO对象
func (s *Service) convertToVO(file *SystemDB.File) *common.FileVO {
	// 使用统一的工具方法构建最终的访问URL
	finalURL := utils.BuildFileURL(file.Url, file.StorageType)

	return &common.FileVO{
		Id:           file.Id,
		CategoryId:   file.CategoryId,
		Name:         file.Name,
		OriginalName: file.OriginalName,
		Url:          finalURL,
		Size:         file.Size,
		Type:         file.Type,
		CreateTime:   file.CreateTime.Format("2006-01-02 15:04:05"),
	}
}

// buildTree 构建树形结构
func buildTree(categories []SystemDB.FileCategory) []common.FileCategoryTreeVO {
	// 构建父子关系map
	categoryMap := make(map[int64][]SystemDB.FileCategory)
	for _, category := range categories {
		categoryMap[category.ParentId] = append(categoryMap[category.ParentId], category)
	}

	// 递归构建树
	var build func(parentId int64) []common.FileCategoryTreeVO
	build = func(parentId int64) []common.FileCategoryTreeVO {
		children := categoryMap[parentId]
		if len(children) == 0 {
			return nil
		}

		var nodes []common.FileCategoryTreeVO
		for _, child := range children {
			node := common.FileCategoryTreeVO{
				Id:       child.Id,
				ParentId: child.ParentId,
				Name:     child.Name,
				Code:     child.Code,
				Icon:     child.Icon,
				Sort:     child.Sort,
				Children: build(child.Id),
			}
			nodes = append(nodes, node)
		}
		return nodes
	}

	return build(0)
}

// MoveFiles 批量移动文件到指定分类
func (s *Service) MoveFiles(req *common.MoveFilesReq) error {
	// 检查分类是否存在
	var category SystemDB.FileCategory
	if err := db.WnAdmin().Where("id = ?", req.CategoryId).First(&category).Error; err != nil {
		return errors.New("分类不存在")
	}

	// 批量更新文件分类
	err := db.WnAdmin().Model(&SystemDB.File{}).
		Where("id IN ?", req.FileIds).
		Update("category_id", req.CategoryId).Error

	return err
}
