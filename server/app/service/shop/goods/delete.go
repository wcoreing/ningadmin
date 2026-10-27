package goods

import (
	"fmt"
	"wnadmin/app/model/do/shop"
	"wnadmin/app/provider/db"

	"gorm.io/gorm"
)

// DeleteService 商品删除服务
type DeleteService struct{}

// Execute 执行商品删除
func (s *DeleteService) Execute(id int64) error {
	return db.WnAdmin().Transaction(func(tx *gorm.DB) error {
		// 1. 检查商品是否存在
		var goods shop.GoodsDO
		if err := tx.Where("id = ? AND is_delete = 0", id).First(&goods).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("商品不存在")
			}
			return fmt.Errorf("查询商品失败: %v", err)
		}

		// 2. 删除SKU规格关联
		if err := s.deleteSkuSpecs(tx, id); err != nil {
			return err
		}

		// 3. 删除SKU
		if err := s.deleteSkus(tx, id); err != nil {
			return err
		}

		// 4. 删除规格值
		if err := s.deleteSpecValues(tx, id); err != nil {
			return err
		}

		// 5. 删除规格组
		if err := s.deleteSpecGroups(tx, id); err != nil {
			return err
		}

		// 6. 删除商品图片
		if err := s.deleteImages(tx, id); err != nil {
			return err
		}

		// 7. 删除商品详情
		if err := s.deleteDetail(tx, id); err != nil {
			return err
		}

		// 8. 删除标签关联
		if err := s.deleteTagRelations(tx, id); err != nil {
			return err
		}

		// 9. 最后软删除商品主表
		if err := tx.Model(&shop.GoodsDO{}).Where("id = ?", id).Update("is_delete", 1).Error; err != nil {
			return fmt.Errorf("删除商品失败: %v", err)
		}

		return nil
	})
}

// deleteSkuSpecs 删除SKU规格关联
func (s *DeleteService) deleteSkuSpecs(tx *gorm.DB, goodsId int64) error {
	// 先查询该商品的所有SKU
	var skuIds []int64
	if err := tx.Model(&shop.ShopGoodsSku{}).Where("goods_id = ?", goodsId).Pluck("id", &skuIds).Error; err != nil {
		return fmt.Errorf("查询商品SKU失败: %v", err)
	}

	// 删除这些SKU的规格关联
	if len(skuIds) > 0 {
		if err := tx.Where("sku_id IN ?", skuIds).Delete(&shop.ShopGoodsSkuSpec{}).Error; err != nil {
			return fmt.Errorf("删除SKU规格关联失败: %v", err)
		}
	}

	return nil
}

// deleteSkus 删除SKU
func (s *DeleteService) deleteSkus(tx *gorm.DB, goodsId int64) error {
	if err := tx.Where("goods_id = ?", goodsId).Delete(&shop.ShopGoodsSku{}).Error; err != nil {
		return fmt.Errorf("删除商品SKU失败: %v", err)
	}
	return nil
}

// deleteSpecValues 删除规格值
func (s *DeleteService) deleteSpecValues(tx *gorm.DB, goodsId int64) error {
	if err := tx.Where("goods_id = ?", goodsId).Delete(&shop.ShopGoodsSpecValue{}).Error; err != nil {
		return fmt.Errorf("删除商品规格值失败: %v", err)
	}
	return nil
}

// deleteSpecGroups 删除规格组
func (s *DeleteService) deleteSpecGroups(tx *gorm.DB, goodsId int64) error {
	if err := tx.Where("goods_id = ?", goodsId).Delete(&shop.ShopGoodsSpecGroup{}).Error; err != nil {
		return fmt.Errorf("删除商品规格组失败: %v", err)
	}
	return nil
}

// deleteImages 删除商品图片
func (s *DeleteService) deleteImages(tx *gorm.DB, goodsId int64) error {
	if err := tx.Where("goods_id = ?", goodsId).Delete(&shop.ShopGoodsImage{}).Error; err != nil {
		return fmt.Errorf("删除商品图片失败: %v", err)
	}
	return nil
}

// deleteDetail 删除商品详情
func (s *DeleteService) deleteDetail(tx *gorm.DB, goodsId int64) error {
	if err := tx.Where("goods_id = ?", goodsId).Delete(&shop.ShopGoodsDetail{}).Error; err != nil {
		return fmt.Errorf("删除商品详情失败: %v", err)
	}
	return nil
}

// deleteTagRelations 删除标签关联
func (s *DeleteService) deleteTagRelations(tx *gorm.DB, goodsId int64) error {
	if err := tx.Where("goods_id = ?", goodsId).Delete(&shop.ShopGoodsTagRel{}).Error; err != nil {
		return fmt.Errorf("删除商品标签关联失败: %v", err)
	}
	return nil
}
