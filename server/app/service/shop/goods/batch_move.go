package goods

import (
	"wnadmin/app/model/do/shop"
	"wnadmin/app/model/dto"
	"wnadmin/app/provider/db"
)

// BatchMoveService 批量移动服务
type BatchMoveService struct{}

// Execute 执行批量移动
func (s *BatchMoveService) Execute(req *dto.GoodsBatchMoveReq) error {
	return db.WnAdmin().Model(&shop.GoodsDO{}).
		Where("id IN ? AND is_delete = 0", req.GoodsIds).
		Update("category_id", req.CategoryId).Error
}
