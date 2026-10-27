package dto

// BrandDTO 品牌DTO
type BrandDTO struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Logo        string `json:"logo"`
	Description string `json:"description"`
	Website     string `json:"website"`
	Sort        int    `json:"sort"`
}

// SupplierDTO 供应商DTO
type SupplierDTO struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Contact     string `json:"contact"`
	Phone       string `json:"phone"`
	Address     string `json:"address"`
	Description string `json:"description"`
	Level       int8   `json:"level"`
	Sort        int    `json:"sort"`
}

// CategoryDTO 商品分类DTO
type CategoryDTO struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	ParentID int64  `json:"parentId"`
	Level    int8   `json:"level"`
	Sort     int    `json:"sort"`
}

// TagDTO 商品标签DTO
type TagDTO struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Sort int    `json:"sort"`
}
