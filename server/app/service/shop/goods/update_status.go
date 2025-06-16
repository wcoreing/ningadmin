package goods

import (
	"wnadmin/app/model/do/shop"
	"wnadmin/app/provider/db"
)

// UpdateStatusService 商品状态更新服务
type UpdateStatusService struct{}

// Execute 执行商品状态更新
func (s *UpdateStatusService) Execute(id int64, status int8) error {
	// 更新状态
	return db.WnAdmin().Model(&shop.GoodsDO{}).
		Where("is_delete = 0 AND id = ?", id).
		Update("status", status).Error
}
