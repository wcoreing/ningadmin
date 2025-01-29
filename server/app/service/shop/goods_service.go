package shop

import (
	"wnadmin/app/model/do/shop"
	"wnadmin/app/model/dto"
	"wnadmin/app/model/vo"
)

// GoodsService 商品服务
type GoodsService struct{}

// GetGoodsList 获取商品列表
func (s *GoodsService) GetGoodsList(query *dto.GoodsListReq) (*vo.GoodsListVO, error) {
	goodsDO := &shop.GoodsDO{}

	goods, total, err := goodsDO.GetGoodsList(query)
	if err != nil {
		return nil, err
	}

	// 转换为响应结构
	list := make([]vo.GoodsVO, 0, len(goods))
	for _, g := range goods {
		list = append(list, vo.GoodsVO{
			GoodsId:       g.GoodsId,
			CategoryId:    g.CategoryId,
			Type:          g.Type,
			Name:          g.Name,
			Description:   g.Description,
			Price:         g.Price,
			OriginalPrice: g.OriginalPrice,
			Stock:         g.Stock,
			Sales:         g.Sales,
			Status:        g.Status,
			Sort:          g.Sort,
		})
	}

	return &vo.GoodsListVO{
		List:  list,
		Total: total,
	}, nil
}

// DeleteGoods 删除商品
func (s *GoodsService) DeleteGoods(id int64) error {
	goodsDO := &shop.GoodsDO{GoodsId: id}
	return goodsDO.Delete()
}

// UpdateGoodsStatus 更新商品状态
func (s *GoodsService) UpdateGoodsStatus(id int64, status int8) error {
	goodsDO := &shop.GoodsDO{
		GoodsId: id,
		Status:  status,
	}
	return goodsDO.UpdateStatus()
}
