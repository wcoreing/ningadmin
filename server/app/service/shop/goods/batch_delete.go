package goods

import (
	"fmt"
	"wnadmin/app/model/do/shop"
	"wnadmin/app/model/dto"
	"wnadmin/app/provider/db"

	"gorm.io/gorm"
)

// BatchDeleteService 批量删除服务
type BatchDeleteService struct{}

// Execute 执行批量删除
func (s *BatchDeleteService) Execute(req *dto.GoodsBatchDeleteReq) error {
	return db.WnAdmin().Transaction(func(tx *gorm.DB) error {
		// 1. 检查商品是否存在
		var existingIds []int64
		if err := tx.Model(&shop.GoodsDO{}).
			Where("id IN ? AND is_delete = 0", req.GoodsIds).
			Pluck("id", &existingIds).Error; err != nil {
			return fmt.Errorf("查询商品失败: %v", err)
		}

		if len(existingIds) == 0 {
			return fmt.Errorf("没有找到可删除的商品")
		}

		// 2. 删除SKU规格关联
		if err := s.deleteSkuSpecs(tx, existingIds); err != nil {
			return err
		}

		// 3. 删除SKU
		if err := s.deleteSkus(tx, existingIds); err != nil {
			return err
		}

		// 4. 删除规格值
		if err := s.deleteSpecValues(tx, existingIds); err != nil {
			return err
		}

		// 5. 删除规格组
		if err := s.deleteSpecGroups(tx, existingIds); err != nil {
			return err
		}

		// 6. 删除商品图片
		if err := s.deleteImages(tx, existingIds); err != nil {
			return err
		}

		// 7. 删除商品详情
		if err := s.deleteDetail(tx, existingIds); err != nil {
			return err
		}

		// 8. 删除标签关联
		if err := s.deleteTagRelations(tx, existingIds); err != nil {
			return err
		}

		// 9. 最后软删除商品主表
		if err := tx.Model(&shop.GoodsDO{}).
			Where("id IN ?", existingIds).
			Update("is_delete", 1).Error; err != nil {
			return fmt.Errorf("批量删除商品失败: %v", err)
		}

		return nil
	})
}

// deleteSkuSpecs 删除SKU规格关联
func (s *BatchDeleteService) deleteSkuSpecs(tx *gorm.DB, goodsIds []int64) error {
	// 先查询这些商品的所有SKU
	var skuIds []int64
	if err := tx.Model(&shop.ShopGoodsSku{}).Where("goods_id IN ?", goodsIds).Pluck("id", &skuIds).Error; err != nil {
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
func (s *BatchDeleteService) deleteSkus(tx *gorm.DB, goodsIds []int64) error {
	if err := tx.Where("goods_id IN ?", goodsIds).Delete(&shop.ShopGoodsSku{}).Error; err != nil {
		return fmt.Errorf("删除商品SKU失败: %v", err)
	}
	return nil
}

// deleteSpecValues 删除规格值
func (s *BatchDeleteService) deleteSpecValues(tx *gorm.DB, goodsIds []int64) error {
	if err := tx.Where("goods_id IN ?", goodsIds).Delete(&shop.ShopGoodsSpecValue{}).Error; err != nil {
		return fmt.Errorf("删除商品规格值失败: %v", err)
	}
	return nil
}

// deleteSpecGroups 删除规格组
func (s *BatchDeleteService) deleteSpecGroups(tx *gorm.DB, goodsIds []int64) error {
	if err := tx.Where("goods_id IN ?", goodsIds).Delete(&shop.ShopGoodsSpecGroup{}).Error; err != nil {
		return fmt.Errorf("删除商品规格组失败: %v", err)
	}
	return nil
}

// deleteImages 删除商品图片
func (s *BatchDeleteService) deleteImages(tx *gorm.DB, goodsIds []int64) error {
	if err := tx.Where("goods_id IN ?", goodsIds).Delete(&shop.ShopGoodsImage{}).Error; err != nil {
		return fmt.Errorf("删除商品图片失败: %v", err)
	}
	return nil
}

// deleteDetail 删除商品详情
func (s *BatchDeleteService) deleteDetail(tx *gorm.DB, goodsIds []int64) error {
	if err := tx.Where("goods_id IN ?", goodsIds).Delete(&shop.ShopGoodsDetail{}).Error; err != nil {
		return fmt.Errorf("删除商品详情失败: %v", err)
	}
	return nil
}

// deleteTagRelations 删除标签关联
func (s *BatchDeleteService) deleteTagRelations(tx *gorm.DB, goodsIds []int64) error {
	if err := tx.Where("goods_id IN ?", goodsIds).Delete(&shop.ShopGoodsTagRel{}).Error; err != nil {
		return fmt.Errorf("删除商品标签关联失败: %v", err)
	}
	return nil
}
