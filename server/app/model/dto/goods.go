package dto

// GoodsListReq 商品列表请求参数
type GoodsListReq struct {
	PageNum     int     `json:"pageNum" form:"pageNum"`       // 页码
	PageSize    int     `json:"pageSize" form:"pageSize"`     // 每页数量
	Name        string  `json:"name" form:"name"`             // 商品名称
	CategoryId  int64   `json:"categoryId" form:"categoryId"` // 分类ID
	CategoryIds []int64 `json:"-" form:"-"`                   // 分类ID列表(内部使用)
	Brand       string  `json:"brand" form:"brand"`           // 品牌名称
	TagId       int64   `json:"tagId" form:"tagId"`           // 标签ID
	Type        int8    `json:"type" form:"type"`             // 商品类型：1-实物商品 2-虚拟商品 3-套餐商品
	GoodsId     int64   `json:"-" form:"-"`                   // 商品ID(内部使用)
	Status      string  `form:"status"`
}

// GoodsDTO 商品信息DTO
type GoodsDTO struct {
	GoodsId      int64    `json:"goodsId"`
	CategoryId   int64    `json:"categoryId"`
	CategoryName string   `json:"categoryName"`
	Brand        string   `json:"brand"`
	Name         string   `json:"name"`
	Status       int8     `json:"status"`
	Sort         int      `json:"sort"`
	TagNames     []string `json:"tagNames"` // 标签名称列表
	CoverImg     string   `json:"coverImg"`
	// 价格和库存信息从SKU中获取
	MinPrice   float64 `json:"minPrice"`   // 最低价格（从SKU中计算）
	MaxPrice   float64 `json:"maxPrice"`   // 最高价格（从SKU中计算）
	TotalStock int     `json:"totalStock"` // 总库存（从SKU中计算）
}

// GoodsStatusReq 商品状态修改请求
type GoodsStatusReq struct {
	GoodsId int64 `json:"goodsId" binding:"required"`
	Status  int   `json:"status" binding:"gte=0,lte=1"`
}

// GoodsImageUpdateItem 商品图片更新项
type GoodsImageUpdateItem struct {
	Url      string `json:"url"`      // 图片URL
	Position uint8  `json:"position"` // 图片位置
	Sort     int    `json:"sort"`     // 排序
}

// GoodsUpdateReq 商品更新请求
type GoodsUpdateReq struct {
	GoodsId     int64                   `json:"goodsId,omitempty"`
	GoodsName   string                  `json:"goodsName" binding:"required"`
	CategoryId  int64                   `json:"categoryId" binding:"required"`
	Brand       string                  `json:"brand" binding:"required"`
	Supplier    string                  `json:"supplier" binding:"required"`
	Type        int8                    `json:"type"` // 商品类型：1-实物商品 2-虚拟商品 3-套餐商品
	Unit        string                  `json:"unit"` // 商品单位
	Sort        int                     `json:"sort"`
	Status      int                     `json:"status" binding:"oneof=0 1"`
	Description string                  `json:"description"`
	CoverImg    string                  `json:"coverImg"`
	TagNames    []string                `json:"tagNames"`
	Images      []*GoodsImageUpdateItem `json:"images"`
	SpecType    int                     `json:"specType"`   // 规格类型：1-单规格，2-多规格，3-套餐规格
	DetailDesc  string                  `json:"detailDesc"` // 商品详情

	// 多规格数据存储 - 支持同时保存多种规格类型的数据
	SingleSpecData *SingleSpecData `json:"singleSpecData"` // 单规格数据
	MultiSpecData  *MultiSpecData  `json:"multiSpecData"`  // 多规格数据
	ComboSpecData  *ComboSpecData  `json:"comboSpecData"`  // 套餐规格数据

	// 兼容旧版本的字段（优先级低于上面的分类数据）
	SpecGroups  []SpecGroup      `json:"specGroups"`  // 规格组（兼容）
	SkuList     []Sku            `json:"skuList"`     // SKU列表（兼容）
	ComboGroups []ComboSpecGroup `json:"comboGroups"` // 套餐组（兼容）

	// 单规格商品的价格库存字段（兼容旧版本）
	Price         float64 `json:"price"`         // 价格
	OriginalPrice float64 `json:"originalPrice"` // 原价
	CostPrice     float64 `json:"costPrice"`     // 成本价
	Stock         int     `json:"stock"`         // 库存
	Weight        float64 `json:"weight"`        // 重量
	Volume        float64 `json:"volume"`        // 体积
}

// SingleSpecData 单规格数据
type SingleSpecData struct {
	Price         float64 `json:"price"`         // 价格
	OriginalPrice float64 `json:"originalPrice"` // 原价
	CostPrice     float64 `json:"costPrice"`     // 成本价
	Stock         int     `json:"stock"`         // 库存
	Weight        float64 `json:"weight"`        // 重量
	Volume        float64 `json:"volume"`        // 体积
}

// MultiSpecData 多规格数据
type MultiSpecData struct {
	SpecGroups []SpecGroup `json:"specGroups"` // 规格组
	SkuList    []Sku       `json:"skuList"`    // SKU列表
}

// ComboSpecData 套餐规格数据
type ComboSpecData struct {
	Price       float64          `json:"price"`       // 套餐价格
	ComboGroups []ComboSpecGroup `json:"comboGroups"` // 套餐组
}

// SpecGroup 规格组
type SpecGroup struct {
	ID           int64       `json:"id"`           // 规格组ID
	Name         string      `json:"name"`         // 规格组名称
	Description  string      `json:"description"`  // 组描述（套餐组用）
	MinSelect    int         `json:"minSelect"`    // 最少选择数量（套餐组用）
	MaxSelect    int         `json:"maxSelect"`    // 最多选择数量（套餐组用）
	IsRequired   bool        `json:"isRequired"`   // 是否必选（套餐组用）
	Values       []SpecValue `json:"values"`       // 规格值列表
	Sort         int         `json:"sort"`         // 排序
	InputVisible bool        `json:"inputVisible"` // 输入框是否可见
	InputValue   string      `json:"inputValue"`   // 输入框的值
}

// ComboSpecGroup 套餐规格组
type ComboSpecGroup struct {
	ID          int64             `json:"id"`          // 规格组ID
	Name        string            `json:"name"`        // 规格组名称
	Description string            `json:"description"` // 组描述
	MinSelect   int               `json:"minSelect"`   // 最少选择数量
	MaxSelect   int               `json:"maxSelect"`   // 最多选择数量
	IsRequired  bool              `json:"isRequired"`  // 是否必选
	Sort        int               `json:"sort"`        // 排序
	Options     []ComboSpecOption `json:"options"`     // 套餐选项列表
}

// ComboSpecOption 套餐规格选项
type ComboSpecOption struct {
	ID            int64   `json:"id"`            // 规格值ID
	Name          string  `json:"name"`          // 选项名称
	Value         string  `json:"value"`         // 选项值
	GoodsId       int64   `json:"goodsId"`       // 关联商品ID
	GoodsName     string  `json:"goodsName"`     // 关联商品名称
	GoodsImage    string  `json:"goodsImage"`    // 关联商品图片
	SkuId         *int64  `json:"skuId"`         // 关联SKU ID
	AddPrice      float64 `json:"addPrice"`      // 加价金额
	Quantity      int     `json:"quantity"`      // 数量
	IsDefault     bool    `json:"isDefault"`     // 是否默认选择
	Sort          int     `json:"sort"`          // 排序
	OriginalPrice float64 `json:"originalPrice"` // 原商品价格
	Stock         int     `json:"stock"`         // 库存
}

// Sku SKU信息
type Sku struct {
	ID            int64                `json:"id,omitempty"`
	Specs         map[string]SpecValue `json:"specs"`         // 规格信息，key为规格组名称
	Price         float64              `json:"price"`         // 价格
	OriginalPrice float64              `json:"originalPrice"` // 原价
	CostPrice     float64              `json:"costPrice"`     // 成本价
	Stock         int                  `json:"stock"`         // 库存
	Weight        float64              `json:"weight"`        // 重量
	Volume        float64              `json:"volume"`        // 体积
	Image         string               `json:"image"`         // SKU图片
}

// SpecValue 规格值
type SpecValue struct {
	ID            int64   `json:"id"`            // 规格值ID
	Name          string  `json:"name"`          // 规格值名称
	Value         string  `json:"value"`         // 规格值
	Sort          int     `json:"sort"`          // 排序
	OptionGoodsId *int64  `json:"optionGoodsId"` // 关联商品ID（套餐选项用）
	OptionSkuId   *int64  `json:"optionSkuId"`   // 关联SKU ID（套餐选项用）
	AddPrice      float64 `json:"addPrice"`      // 加价金额（套餐选项用）
	Quantity      int     `json:"quantity"`      // 数量（套餐选项用）
	IsDefault     bool    `json:"isDefault"`     // 是否默认选择（套餐选项用）
}

// GoodsInfoResp 商品详情响应
type GoodsInfoResp struct {
	GoodsId     int64                   `json:"goodsId"`
	GoodsName   string                  `json:"goodsName"`
	CategoryId  int64                   `json:"categoryId"`
	Brand       string                  `json:"brand"`
	Supplier    string                  `json:"supplier"`
	Type        int8                    `json:"type"`
	Unit        string                  `json:"unit"`
	Sort        int                     `json:"sort"`
	Status      int8                    `json:"status"`
	Description string                  `json:"description"`
	CoverImg    string                  `json:"coverImg"`
	TagNames    []string                `json:"tagNames"`
	Images      []*GoodsImageUpdateItem `json:"images"`
	SpecType    int8                    `json:"specType"`   // 规格类型：1-单规格，2-多规格，3-套餐规格
	DetailDesc  string                  `json:"detailDesc"` // 商品详情

	// 多规格数据存储 - 支持同时返回多种规格类型的数据
	SingleSpecData *SingleSpecData `json:"singleSpecData"` // 单规格数据
	MultiSpecData  *MultiSpecData  `json:"multiSpecData"`  // 多规格数据
	ComboSpecData  *ComboSpecData  `json:"comboSpecData"`  // 套餐规格数据

	// 兼容旧版本的字段（根据当前specType返回对应数据）
	SpecGroups  []SpecGroup      `json:"specGroups"`  // 规格组（兼容）
	SkuList     []Sku            `json:"skuList"`     // SKU列表（兼容）
	ComboGroups []ComboSpecGroup `json:"comboGroups"` // 套餐组（兼容）

}

// GoodsBatchMoveReq 商品批量移动请求
type GoodsBatchMoveReq struct {
	GoodsIds   []int64 `json:"goodsIds" binding:"required"`
	CategoryId int64   `json:"categoryId" binding:"required"`
}

// GoodsBatchDeleteReq 商品批量删除请求
type GoodsBatchDeleteReq struct {
	GoodsIds []int64 `json:"goodsIds" binding:"required"`
}
