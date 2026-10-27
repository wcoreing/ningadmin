-- 为shop_goods_sku表添加spec_type字段的迁移脚本
-- 执行时间：2025-01-24

-- 1. 添加spec_type字段
ALTER TABLE `shop_goods_sku` 
ADD COLUMN `spec_type` tinyint(4) DEFAULT '1' COMMENT '规格类型：1-单规格 2-多规格 3-套餐规格' 
AFTER `sku_code`;

-- 2. 根据商品的spec_type更新SKU的spec_type
-- 更新单规格商品的SKU
UPDATE `shop_goods_sku` sku 
INNER JOIN `shop_goods` goods ON sku.goods_id = goods.id 
SET sku.spec_type = 1 
WHERE goods.spec_type = 1;

-- 更新多规格商品的SKU
UPDATE `shop_goods_sku` sku 
INNER JOIN `shop_goods` goods ON sku.goods_id = goods.id 
SET sku.spec_type = 2 
WHERE goods.spec_type = 2;

-- 更新套餐商品的SKU
UPDATE `shop_goods_sku` sku 
INNER JOIN `shop_goods` goods ON sku.goods_id = goods.id 
SET sku.spec_type = 3 
WHERE goods.spec_type = 3;

-- 3. 为没有对应商品记录的SKU设置默认值（防止数据不一致）
UPDATE `shop_goods_sku` 
SET spec_type = 1 
WHERE spec_type IS NULL; 