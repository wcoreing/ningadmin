package shop

import (
	"errors"
	"fmt"
	"wnadmin/app/model/do/shop"
	"wnadmin/app/model/dto"
	"wnadmin/app/provider/db"
)

type BannerService struct{}

// 获取Banner列表
func (s *BannerService) GetBannerList(req *dto.BannerListReq) (*dto.PageResult, error) {
	var list []shop.ShopBanner
	var total int64

	bannerDO := &shop.ShopBanner{}
	db := bannerDO.Query()

	fmt.Printf("查询条件: %+v\n", req)

	// 添加查询条件
	if req.Title != "" {
		db = db.Where("title LIKE ?", "%"+req.Title+"%")
	}
	if req.Position > 0 {
		db = db.Where("position = ?", req.Position)
	}
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}

	// 获取总数
	if err := db.Count(&total).Error; err != nil {
		fmt.Printf("获取总数失败: %v\n", err)
		return nil, err
	}
	fmt.Printf("总数: %d\n", total)

	// 分页查询
	if err := db.Order("sort desc").
		Offset((req.PageNum - 1) * req.PageSize).
		Limit(req.PageSize).
		Find(&list).Error; err != nil {
		fmt.Printf("查询列表失败: %v\n", err)
		return nil, err
	}
	fmt.Printf("查询到的数据: %+v\n", list)

	// 构建响应数据
	var result []dto.BannerListResp
	for _, banner := range list {
		item := dto.BannerListResp{
			Id:       banner.Id,
			Title:    banner.Title,
			Image:    banner.Image,
			Position: banner.Position,
			Type:     banner.Type,
			TargetId: banner.TargetId,
			Url:      banner.Url,
			Sort:     banner.Sort,
			Status:   banner.Status,
		}
		// 如果是商品类型且有目标ID，获取商品信息
		if banner.Type == 1 && banner.TargetId > 0 {
			var goods shop.GoodsDO
			if err := (&shop.GoodsDO{}).Query().Where("id = ?", banner.TargetId).First(&goods).Error; err == nil {
				// 获取商品主图
				var coverImg string
				var image shop.ShopGoodsImage
				if err := (&shop.ShopGoodsImage{}).Query().Where("goods_id = ? AND position = ?", goods.Id, 1).
					First(&image).Error; err == nil {
					coverImg = image.Url
				}

				// 获取商品最低价格（从SKU中获取）
				var minPrice float64
				var sku shop.ShopGoodsSku
				if err := (&shop.ShopGoodsSku{}).Query().Where("goods_id = ? AND status = 1", goods.Id).
					Order("price ASC").First(&sku).Error; err == nil {
					minPrice = sku.Price
				}

				item.GoodsInfo = &dto.GoodsBaseInfo{
					GoodsId:  uint(goods.Id),
					Name:     goods.Name,
					CoverImg: coverImg,
					Price:    minPrice,
				}
			}
		}
		result = append(result, item)
	}

	return &dto.PageResult{
		Total: total,
		List:  result,
	}, nil
}

// 获取指定位置的Banner列表
func (s *BannerService) GetBannersByPosition(position int) ([]*shop.ShopBanner, error) {
	var list []*shop.ShopBanner
	bannerDO := &shop.ShopBanner{}
	err := bannerDO.Query().
		Where("position = ? AND status = ?", position, shop.BannerStatusEnable).
		Order("sort desc, id desc").
		Find(&list).Error
	return list, err
}

// 创建Banner
func (s *BannerService) CreateBanner(req *dto.BannerCreateReq) error {
	// 打印请求参数，方便调试
	fmt.Printf("CreateBanner request: %+v\n", req)

	banner := &shop.ShopBanner{
		Title:    req.Title,
		Image:    req.Image,
		Url:      req.Url,
		Position: req.Position,
		Type:     req.Type,
		TargetId: req.TargetId,
		Sort:     req.Sort,
		Status:   req.Status,
	}
	return banner.Create()
}

// 更新Banner
func (s *BannerService) UpdateBanner(req *dto.BannerUpdateReq) error {
	// 打印请求参数，方便调试
	fmt.Printf("UpdateBanner request: %+v\n", req)

	banner := &shop.ShopBanner{Id: req.Id}
	if err := banner.First(); err != nil {
		return err
	}
	if banner.IsDelete == 1 {
		return errors.New("banner不存在")
	}

	banner.Title = req.Title
	banner.Image = req.Image
	banner.Url = req.Url
	banner.Position = req.Position
	banner.Type = req.Type
	banner.TargetId = req.TargetId
	banner.Sort = req.Sort
	banner.Status = req.Status

	fmt.Printf("更新前的banner数据: %+v\n", banner)
	return banner.Update()
}

// 删除Banner
func (s *BannerService) DeleteBanner(id int64) error {
	banner := &shop.ShopBanner{Id: id}
	return banner.Query().Update("is_delete", 1).Error
}

// 修改Banner状态
func (s *BannerService) ChangeBannerStatus(id int64, status int) error {
	if status != shop.BannerStatusDisable && status != shop.BannerStatusEnable {
		return errors.New("无效的状态值")
	}

	banner := &shop.ShopBanner{Id: id}
	return banner.Query().Update("status", status).Error
}

// GetBannerInfo 获取Banner详情
func (s *BannerService) GetBannerInfo(id int64) (*dto.BannerInfoResp, error) {
	var banner shop.ShopBanner
	if err := db.WnAdmin().Where("id = ?", id).First(&banner).Error; err != nil {
		return nil, err
	}

	info := &dto.BannerInfoResp{
		Id:       banner.Id,
		Title:    banner.Title,
		Image:    banner.Image,
		Url:      banner.Url,
		Position: banner.Position,
		Type:     banner.Type,
		TargetId: banner.TargetId,
		Sort:     banner.Sort,
		Status:   banner.Status,
	}

	// 如果是商品类型且有目标ID，获取商品信息
	if banner.Type == 1 && banner.TargetId > 0 {
		var goods shop.GoodsDO
		if err := (&shop.GoodsDO{}).Query().Where("id = ?", banner.TargetId).First(&goods).Error; err == nil {
			// 获取商品主图
			var coverImg string
			var image shop.ShopGoodsImage
			if err := (&shop.ShopGoodsImage{}).Query().Where("goods_id = ? AND position = ?", goods.Id, 1).
				First(&image).Error; err == nil {
				coverImg = image.Url
			}

			// 获取商品最低价格（从SKU中获取）
			var minPrice float64
			var sku shop.ShopGoodsSku
			if err := (&shop.ShopGoodsSku{}).Query().Where("goods_id = ? AND status = 1", goods.Id).
				Order("price ASC").First(&sku).Error; err == nil {
				minPrice = sku.Price
			}

			info.GoodsInfo = &dto.GoodsBaseInfo{
				GoodsId:  uint(goods.Id),
				Name:     goods.Name,
				CoverImg: coverImg,
				Price:    minPrice,
			}
		}
	}

	return info, nil
}
