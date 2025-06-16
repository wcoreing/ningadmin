package goods

import (
	"fmt"
	"wnadmin/app/common/utils"
	"wnadmin/app/model/do/shop"
	"wnadmin/app/model/dto"
	"wnadmin/app/provider/db"

	"gorm.io/gorm"
)

// GetInfoService 商品详情查询服务
type GetInfoService struct{}

// Execute 执行商品详情查询
func (s *GetInfoService) Execute(id int64) (*dto.GoodsInfoResp, error) {
	// 查询商品基本信息
	goods, err := s.getGoodsById(id)
	if err != nil {
		return nil, err
	}

	// 查询商品详情
	detail, err := s.getGoodsDetail(id)
	if err != nil {
		return nil, err
	}

	// 构建响应对象
	resp := s.buildGoodsInfoResp(goods, detail)

	// 处理图片信息
	s.processImages(goods, resp)

	// 加载所有规格类型的数据（新版本）
	s.loadAllSpecData(goods.Id, resp)

	// 添加调试日志
	fmt.Printf("返回的图片数据: CoverImg=%s, Images=%+v\n", resp.CoverImg, resp.Images)

	return resp, nil
}

// loadAllSpecData 加载所有规格类型的数据
func (s *GetInfoService) loadAllSpecData(goodsId int64, resp *dto.GoodsInfoResp) {
	multiSpecService := &MultiSpecService{}
	allSpecData, err := multiSpecService.LoadAllSpecData(goodsId)
	if err != nil {
		// 如果加载失败，记录错误但不影响主流程
		fmt.Printf("加载多规格数据失败: %v\n", err)
		return
	}

	// 设置多规格数据到响应中
	resp.SingleSpecData = allSpecData.SingleSpecData
	resp.MultiSpecData = allSpecData.MultiSpecData
	resp.ComboSpecData = allSpecData.ComboSpecData

	// 不设置根级别的价格字段，让前端从各自的规格数据中读取
	// 这样可以避免混淆，每种规格类型都有自己独立的价格数据
}

// GetGoodsById 公开方法：根据ID获取商品（兼容方法）
func (s *GetInfoService) GetGoodsById(id int64) (*shop.GoodsDO, error) {
	return s.getGoodsById(id)
}

// getGoodsById 根据ID查询商品
func (s *GetInfoService) getGoodsById(id int64) (*shop.GoodsDO, error) {
	goods := &shop.GoodsDO{}

	// 查询商品基本信息
	if err := db.WnAdmin().First(goods, id).Error; err != nil {
		return nil, err
	}

	// 查询商品图片
	if err := s.loadImages(goods, id); err != nil {
		return nil, err
	}

	// 查询商品标签
	if err := s.loadTags(goods, id); err != nil {
		return nil, err
	}

	return goods, nil
}

// loadImages 加载商品图片
func (s *GetInfoService) loadImages(goods *shop.GoodsDO, id int64) error {
	var images []*shop.ShopGoodsImage
	if err := db.WnAdmin().Where("goods_id = ?", id).Find(&images).Error; err != nil {
		return err
	}
	goods.Images = images
	return nil
}

// loadTags 加载商品标签
func (s *GetInfoService) loadTags(goods *shop.GoodsDO, id int64) error {
	var tagRels []*shop.ShopGoodsTagRel
	if err := db.WnAdmin().Where("goods_id = ?", id).Find(&tagRels).Error; err != nil {
		return err
	}
	var tagNames []string
	for _, rel := range tagRels {
		if rel.TagName != "" {
			tagNames = append(tagNames, rel.TagName)
		}
	}
	goods.TagNames = tagNames
	return nil
}

// getGoodsDetail 查询商品详情
func (s *GetInfoService) getGoodsDetail(id int64) (*shop.ShopGoodsDetail, error) {
	var detail shop.ShopGoodsDetail
	if err := db.WnAdmin().Where("goods_id = ?", id).First(&detail).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return nil, err
		}
	}
	return &detail, nil
}

// buildGoodsInfoResp 构建商品详情响应
func (s *GetInfoService) buildGoodsInfoResp(goods *shop.GoodsDO, detail *shop.ShopGoodsDetail) *dto.GoodsInfoResp {
	fmt.Printf("商品图片数据: %+v\n", goods.Images)

	return &dto.GoodsInfoResp{
		GoodsId:     goods.Id,
		GoodsName:   goods.Name,
		CategoryId:  goods.CategoryId,
		Brand:       goods.Brand,
		Supplier:    goods.Supplier,
		Type:        goods.Type,
		Unit:        goods.Unit,
		Sort:        goods.Sort,
		Status:      goods.Status,
		Description: goods.Description,
		DetailDesc:  detail.DetailDesc,
		SpecType:    goods.SpecType,
		TagNames:    goods.TagNames,
	}
}

// processImages 处理图片信息
func (s *GetInfoService) processImages(goods *shop.GoodsDO, resp *dto.GoodsInfoResp) {
	if len(goods.Images) > 0 {
		// 设置主图
		for _, img := range goods.Images {
			if img.Position == shop.ImagePositionMain {
				resp.CoverImg = utils.AddDomainToURL(img.Url)
				break
			}
		}

		// 设置其他图片列表
		var images []*dto.GoodsImageUpdateItem
		for _, img := range goods.Images {
			if img.Position != shop.ImagePositionMain {
				images = append(images, &dto.GoodsImageUpdateItem{
					Url:      utils.AddDomainToURL(img.Url),
					Position: img.Position,
					Sort:     img.Sort,
				})
			}
		}
		resp.Images = images
	} else {
		resp.Images = make([]*dto.GoodsImageUpdateItem, 0) // 确保返回空数组而不是 nil
	}
}
