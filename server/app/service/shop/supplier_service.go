package shop

import (
	"wnadmin/app/model/do/shop"
	"wnadmin/app/model/dto"
)

// SupplierService 供应商服务
type SupplierService struct{}

// GetSupplierList 获取供应商列表
func (s *SupplierService) GetSupplierList() ([]dto.SupplierDTO, error) {
	var suppliers []shop.ShopSupplier

	// 创建查询对象并执行查询
	if err := (&shop.ShopSupplier{
		Status: 1, // 只查询启用状态的供应商
	}).Find(&suppliers); err != nil {
		return nil, err
	}

	var result []dto.SupplierDTO
	for _, supplier := range suppliers {
		result = append(result, dto.SupplierDTO{
			ID:          supplier.ID,
			Name:        supplier.Name,
			Contact:     supplier.Contact,
			Phone:       supplier.Phone,
			Address:     supplier.Address,
			Description: supplier.Description,
			Level:       supplier.Level,
			Sort:        supplier.Sort,
		})
	}
	return result, nil
}
