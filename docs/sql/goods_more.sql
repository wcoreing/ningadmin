-- 补齐产品目录，使前台每页 12 条时能翻页。重复执行不会插入同名产品。
SET NAMES utf8mb4;

INSERT INTO shop_goods (category_id, supplier, brand, type, name, description, unit, spec_type, status, sort, is_delete)
SELECT 525, '', '', 1, '单戒开槽礼盒', '单戒开槽，植绒托按戒圈固定，适合柜台单只陈列。', '个', 1, 1, 10, 0
FROM DUAL WHERE NOT EXISTS (SELECT 1 FROM shop_goods WHERE name='单戒开槽礼盒' AND is_delete=0);

INSERT INTO shop_goods_image (goods_id, position, url, sort, is_delete)
SELECT g.id, 1, '/catalog/goods-box.jpg', 0, 0 FROM shop_goods g
WHERE g.name='单戒开槽礼盒' AND g.is_delete=0
AND NOT EXISTS (SELECT 1 FROM shop_goods_image i WHERE i.goods_id=g.id AND i.position=1 AND i.is_delete=0);

INSERT INTO shop_goods_image (goods_id, position, url, sort, is_delete)
SELECT g.id, 0, '/catalog/goods-box-top.jpg', 1, 0 FROM shop_goods g
WHERE g.name='单戒开槽礼盒' AND g.is_delete=0
AND NOT EXISTS (SELECT 1 FROM shop_goods_image i WHERE i.goods_id=g.id AND i.url='/catalog/goods-box-top.jpg' AND i.is_delete=0);

INSERT INTO shop_goods_detail (goods_id, detail_desc, spec_desc, service_desc)
SELECT g.id, '内托只开一槽，戒圈直径请在询盘里写明。', '材质：纸板+植绒
结构：翻盖
常见尺寸：7×7×4 cm', '可烫金或丝印 logo。' FROM shop_goods g
WHERE g.name='单戒开槽礼盒' AND g.is_delete=0
AND NOT EXISTS (SELECT 1 FROM shop_goods_detail d WHERE d.goods_id=g.id);

INSERT INTO shop_goods (category_id, supplier, brand, type, name, description, unit, spec_type, status, sort, is_delete)
SELECT 525, '', '', 1, '对戒双槽礼盒', '两个戒槽并排，适合对戒和婚礼套装。', '个', 1, 1, 10, 0
FROM DUAL WHERE NOT EXISTS (SELECT 1 FROM shop_goods WHERE name='对戒双槽礼盒' AND is_delete=0);

INSERT INTO shop_goods_image (goods_id, position, url, sort, is_delete)
SELECT g.id, 1, '/catalog/goods-box-top.jpg', 0, 0 FROM shop_goods g
WHERE g.name='对戒双槽礼盒' AND g.is_delete=0
AND NOT EXISTS (SELECT 1 FROM shop_goods_image i WHERE i.goods_id=g.id AND i.position=1 AND i.is_delete=0);

INSERT INTO shop_goods_image (goods_id, position, url, sort, is_delete)
SELECT g.id, 0, '/catalog/goods-box.jpg', 1, 0 FROM shop_goods g
WHERE g.name='对戒双槽礼盒' AND g.is_delete=0
AND NOT EXISTS (SELECT 1 FROM shop_goods_image i WHERE i.goods_id=g.id AND i.url='/catalog/goods-box.jpg' AND i.is_delete=0);

INSERT INTO shop_goods_detail (goods_id, detail_desc, spec_desc, service_desc)
SELECT g.id, '两槽间距可按戒圈调整，盒盖可压一对标记。', '材质：特种纸+植绒
结构：翻盖
常见尺寸：12×7×4 cm', '打样前请给戒圈尺寸。' FROM shop_goods g
WHERE g.name='对戒双槽礼盒' AND g.is_delete=0
AND NOT EXISTS (SELECT 1 FROM shop_goods_detail d WHERE d.goods_id=g.id);

INSERT INTO shop_goods (category_id, supplier, brand, type, name, description, unit, spec_type, status, sort, is_delete)
SELECT 525, '', '', 1, '耳环对装盒', '左右两格放一对耳环，耳针有单独卡位。', '个', 1, 1, 10, 0
FROM DUAL WHERE NOT EXISTS (SELECT 1 FROM shop_goods WHERE name='耳环对装盒' AND is_delete=0);

INSERT INTO shop_goods_image (goods_id, position, url, sort, is_delete)
SELECT g.id, 1, '/catalog/goods-box.jpg', 0, 0 FROM shop_goods g
WHERE g.name='耳环对装盒' AND g.is_delete=0
AND NOT EXISTS (SELECT 1 FROM shop_goods_image i WHERE i.goods_id=g.id AND i.position=1 AND i.is_delete=0);

INSERT INTO shop_goods_detail (goods_id, detail_desc, spec_desc, service_desc)
SELECT g.id, '耳钉、耳钩分开卡位，避免互相刮花。', '材质：灰板+植绒
结构：翻盖
常见尺寸：8×6×3 cm', '请说明是耳钉还是耳钩。' FROM shop_goods g
WHERE g.name='耳环对装盒' AND g.is_delete=0
AND NOT EXISTS (SELECT 1 FROM shop_goods_detail d WHERE d.goods_id=g.id);

INSERT INTO shop_goods (category_id, supplier, brand, type, name, description, unit, spec_type, status, sort, is_delete)
SELECT 525, '', '', 1, '吊坠翻盖盒', '盒身留出链长，吊坠落在植绒槽里。', '个', 1, 1, 10, 0
FROM DUAL WHERE NOT EXISTS (SELECT 1 FROM shop_goods WHERE name='吊坠翻盖盒' AND is_delete=0);

INSERT INTO shop_goods_image (goods_id, position, url, sort, is_delete)
SELECT g.id, 1, '/catalog/goods-drawer.jpg', 0, 0 FROM shop_goods g
WHERE g.name='吊坠翻盖盒' AND g.is_delete=0
AND NOT EXISTS (SELECT 1 FROM shop_goods_image i WHERE i.goods_id=g.id AND i.position=1 AND i.is_delete=0);

INSERT INTO shop_goods_detail (goods_id, detail_desc, spec_desc, service_desc)
SELECT g.id, '链子盘在槽外，吊坠单独固定。', '材质：纸板+植绒
结构：翻盖
常见尺寸：8×8×3.5 cm', '请提供链长和坠子尺寸。' FROM shop_goods g
WHERE g.name='吊坠翻盖盒' AND g.is_delete=0
AND NOT EXISTS (SELECT 1 FROM shop_goods_detail d WHERE d.goods_id=g.id);

INSERT INTO shop_goods (category_id, supplier, brand, type, name, description, unit, spec_type, status, sort, is_delete)
SELECT 525, '', '', 1, '手链长条盒', '长条盒放手链和手镯，避免硬折。', '个', 1, 1, 10, 0
FROM DUAL WHERE NOT EXISTS (SELECT 1 FROM shop_goods WHERE name='手链长条盒' AND is_delete=0);

INSERT INTO shop_goods_image (goods_id, position, url, sort, is_delete)
SELECT g.id, 1, '/catalog/goods-drawer.jpg', 0, 0 FROM shop_goods g
WHERE g.name='手链长条盒' AND g.is_delete=0
AND NOT EXISTS (SELECT 1 FROM shop_goods_image i WHERE i.goods_id=g.id AND i.position=1 AND i.is_delete=0);

INSERT INTO shop_goods_image (goods_id, position, url, sort, is_delete)
SELECT g.id, 0, '/catalog/goods-box.jpg', 1, 0 FROM shop_goods g
WHERE g.name='手链长条盒' AND g.is_delete=0
AND NOT EXISTS (SELECT 1 FROM shop_goods_image i WHERE i.goods_id=g.id AND i.url='/catalog/goods-box.jpg' AND i.is_delete=0);

INSERT INTO shop_goods_detail (goods_id, detail_desc, spec_desc, service_desc)
SELECT g.id, '内长按手链实际长度留空，不共用戒指槽。', '材质：灰板+特种纸
结构：翻盖
常见尺寸：20×5×3 cm', '请写手链长度。' FROM shop_goods g
WHERE g.name='手链长条盒' AND g.is_delete=0
AND NOT EXISTS (SELECT 1 FROM shop_goods_detail d WHERE d.goods_id=g.id);

INSERT INTO shop_goods (category_id, supplier, brand, type, name, description, unit, spec_type, status, sort, is_delete)
SELECT 525, '', '', 1, '天地盖礼盒', '盖和底分开，适合套装和节日礼盒。', '个', 1, 1, 10, 0
FROM DUAL WHERE NOT EXISTS (SELECT 1 FROM shop_goods WHERE name='天地盖礼盒' AND is_delete=0);

INSERT INTO shop_goods_image (goods_id, position, url, sort, is_delete)
SELECT g.id, 1, '/catalog/goods-box.jpg', 0, 0 FROM shop_goods g
WHERE g.name='天地盖礼盒' AND g.is_delete=0
AND NOT EXISTS (SELECT 1 FROM shop_goods_image i WHERE i.goods_id=g.id AND i.position=1 AND i.is_delete=0);

INSERT INTO shop_goods_image (goods_id, position, url, sort, is_delete)
SELECT g.id, 0, '/catalog/goods-box-top.jpg', 1, 0 FROM shop_goods g
WHERE g.name='天地盖礼盒' AND g.is_delete=0
AND NOT EXISTS (SELECT 1 FROM shop_goods_image i WHERE i.goods_id=g.id AND i.url='/catalog/goods-box-top.jpg' AND i.is_delete=0);

INSERT INTO shop_goods_detail (goods_id, detail_desc, spec_desc, service_desc)
SELECT g.id, '盖面可烫金，内托按套装件数开槽。', '材质：灰板+特种纸
结构：天地盖
常见尺寸：15×15×6 cm', '请说明套装里有几件。' FROM shop_goods g
WHERE g.name='天地盖礼盒' AND g.is_delete=0
AND NOT EXISTS (SELECT 1 FROM shop_goods_detail d WHERE d.goods_id=g.id);

INSERT INTO shop_goods (category_id, supplier, brand, type, name, description, unit, spec_type, status, sort, is_delete)
SELECT 525, '', '', 1, '书型首饰盒', '像一本书一样翻开，适合柜台展示一套饰品。', '个', 1, 1, 10, 0
FROM DUAL WHERE NOT EXISTS (SELECT 1 FROM shop_goods WHERE name='书型首饰盒' AND is_delete=0);

INSERT INTO shop_goods_image (goods_id, position, url, sort, is_delete)
SELECT g.id, 1, '/catalog/goods-box-top.jpg', 0, 0 FROM shop_goods g
WHERE g.name='书型首饰盒' AND g.is_delete=0
AND NOT EXISTS (SELECT 1 FROM shop_goods_image i WHERE i.goods_id=g.id AND i.position=1 AND i.is_delete=0);

INSERT INTO shop_goods_detail (goods_id, detail_desc, spec_desc, service_desc)
SELECT g.id, '打开后饰品平铺，方便顾客一次看全。', '材质：纸板+PU
结构：书型翻盖
常见尺寸：16×12×4 cm', '封面颜色可按潘通号对。' FROM shop_goods g
WHERE g.name='书型首饰盒' AND g.is_delete=0
AND NOT EXISTS (SELECT 1 FROM shop_goods_detail d WHERE d.goods_id=g.id);

INSERT INTO shop_goods (category_id, supplier, brand, type, name, description, unit, spec_type, status, sort, is_delete)
SELECT 525, '', '', 1, '三件套礼盒', '戒指、耳环、项链分三格，一套一起出货。', '个', 1, 1, 10, 0
FROM DUAL WHERE NOT EXISTS (SELECT 1 FROM shop_goods WHERE name='三件套礼盒' AND is_delete=0);

INSERT INTO shop_goods_image (goods_id, position, url, sort, is_delete)
SELECT g.id, 1, '/catalog/goods-box.jpg', 0, 0 FROM shop_goods g
WHERE g.name='三件套礼盒' AND g.is_delete=0
AND NOT EXISTS (SELECT 1 FROM shop_goods_image i WHERE i.goods_id=g.id AND i.position=1 AND i.is_delete=0);

INSERT INTO shop_goods_image (goods_id, position, url, sort, is_delete)
SELECT g.id, 0, '/catalog/goods-drawer.jpg', 1, 0 FROM shop_goods g
WHERE g.name='三件套礼盒' AND g.is_delete=0
AND NOT EXISTS (SELECT 1 FROM shop_goods_image i WHERE i.goods_id=g.id AND i.url='/catalog/goods-drawer.jpg' AND i.is_delete=0);

INSERT INTO shop_goods_detail (goods_id, detail_desc, spec_desc, service_desc)
SELECT g.id, '三格尺寸不同，项链格单独留长。', '材质：灰板+植绒
结构：翻盖
常见尺寸：22×8×4 cm', '请列出三件的尺寸。' FROM shop_goods g
WHERE g.name='三件套礼盒' AND g.is_delete=0
AND NOT EXISTS (SELECT 1 FROM shop_goods_detail d WHERE d.goods_id=g.id);

INSERT INTO shop_goods (category_id, supplier, brand, type, name, description, unit, spec_type, status, sort, is_delete)
SELECT 525, '', '', 1, '圆盖首饰盒', '圆盖小盒，适合单只戒指或一对耳钉。', '个', 1, 1, 10, 0
FROM DUAL WHERE NOT EXISTS (SELECT 1 FROM shop_goods WHERE name='圆盖首饰盒' AND is_delete=0);

INSERT INTO shop_goods_image (goods_id, position, url, sort, is_delete)
SELECT g.id, 1, '/catalog/goods-box-top.jpg', 0, 0 FROM shop_goods g
WHERE g.name='圆盖首饰盒' AND g.is_delete=0
AND NOT EXISTS (SELECT 1 FROM shop_goods_image i WHERE i.goods_id=g.id AND i.position=1 AND i.is_delete=0);

INSERT INTO shop_goods_image (goods_id, position, url, sort, is_delete)
SELECT g.id, 0, '/catalog/goods-box.jpg', 1, 0 FROM shop_goods g
WHERE g.name='圆盖首饰盒' AND g.is_delete=0
AND NOT EXISTS (SELECT 1 FROM shop_goods_image i WHERE i.goods_id=g.id AND i.url='/catalog/goods-box.jpg' AND i.is_delete=0);

INSERT INTO shop_goods_detail (goods_id, detail_desc, spec_desc, service_desc)
SELECT g.id, '盖面可压纹，内托仍按饰品开槽。', '材质：纸板
结构：圆盖
常见直径：6 cm', '小批量也可以询盘。' FROM shop_goods g
WHERE g.name='圆盖首饰盒' AND g.is_delete=0
AND NOT EXISTS (SELECT 1 FROM shop_goods_detail d WHERE d.goods_id=g.id);

INSERT INTO shop_goods (category_id, supplier, brand, type, name, description, unit, spec_type, status, sort, is_delete)
SELECT 525, '', '', 1, '植绒空托礼盒', '盒身做好，内托按订单再开槽，方便先看盒型。', '个', 1, 1, 10, 0
FROM DUAL WHERE NOT EXISTS (SELECT 1 FROM shop_goods WHERE name='植绒空托礼盒' AND is_delete=0);

INSERT INTO shop_goods_image (goods_id, position, url, sort, is_delete)
SELECT g.id, 1, '/catalog/goods-box.jpg', 0, 0 FROM shop_goods g
WHERE g.name='植绒空托礼盒' AND g.is_delete=0
AND NOT EXISTS (SELECT 1 FROM shop_goods_image i WHERE i.goods_id=g.id AND i.position=1 AND i.is_delete=0);

INSERT INTO shop_goods_detail (goods_id, detail_desc, spec_desc, service_desc)
SELECT g.id, '先确认盒型和颜色，再按实物开槽。', '材质：纸板+植绒
结构：翻盖
内托：按单开槽', '打样请寄一枚饰品或给尺寸。' FROM shop_goods g
WHERE g.name='植绒空托礼盒' AND g.is_delete=0
AND NOT EXISTS (SELECT 1 FROM shop_goods_detail d WHERE d.goods_id=g.id);

INSERT INTO shop_goods (category_id, supplier, brand, type, name, description, unit, spec_type, status, sort, is_delete)
SELECT 525, '', '', 1, '柜台陈列盒', '盒盖可立住，打开后适合放在柜台里陈列。', '个', 1, 1, 10, 0
FROM DUAL WHERE NOT EXISTS (SELECT 1 FROM shop_goods WHERE name='柜台陈列盒' AND is_delete=0);

INSERT INTO shop_goods_image (goods_id, position, url, sort, is_delete)
SELECT g.id, 1, '/catalog/goods-drawer.jpg', 0, 0 FROM shop_goods g
WHERE g.name='柜台陈列盒' AND g.is_delete=0
AND NOT EXISTS (SELECT 1 FROM shop_goods_image i WHERE i.goods_id=g.id AND i.position=1 AND i.is_delete=0);

INSERT INTO shop_goods_image (goods_id, position, url, sort, is_delete)
SELECT g.id, 0, '/catalog/goods-box-top.jpg', 1, 0 FROM shop_goods g
WHERE g.name='柜台陈列盒' AND g.is_delete=0
AND NOT EXISTS (SELECT 1 FROM shop_goods_image i WHERE i.goods_id=g.id AND i.url='/catalog/goods-box-top.jpg' AND i.is_delete=0);

INSERT INTO shop_goods_detail (goods_id, detail_desc, spec_desc, service_desc)
SELECT g.id, '盖子撑开后不倒，方便店员拿取。', '材质：特种纸+植绒
结构：翻盖可立
常见尺寸：18×12×5 cm', '请说明是玻璃柜还是开放架。' FROM shop_goods g
WHERE g.name='柜台陈列盒' AND g.is_delete=0
AND NOT EXISTS (SELECT 1 FROM shop_goods_detail d WHERE d.goods_id=g.id);

INSERT INTO shop_goods (category_id, supplier, brand, type, name, description, unit, spec_type, status, sort, is_delete)
SELECT 526, '', '', 1, '绒布按扣袋', '袋口用按扣，比抽绳更不容易散开。', '个', 1, 1, 10, 0
FROM DUAL WHERE NOT EXISTS (SELECT 1 FROM shop_goods WHERE name='绒布按扣袋' AND is_delete=0);

INSERT INTO shop_goods_image (goods_id, position, url, sort, is_delete)
SELECT g.id, 1, '/catalog/goods-pouch.jpg', 0, 0 FROM shop_goods g
WHERE g.name='绒布按扣袋' AND g.is_delete=0
AND NOT EXISTS (SELECT 1 FROM shop_goods_image i WHERE i.goods_id=g.id AND i.position=1 AND i.is_delete=0);

INSERT INTO shop_goods_detail (goods_id, detail_desc, spec_desc, service_desc)
SELECT g.id, '按扣颜色可跟袋身分开选。', '材质：植绒布
结构：按扣
常见尺寸：8×10 cm', '适合单件饰品出货。' FROM shop_goods g
WHERE g.name='绒布按扣袋' AND g.is_delete=0
AND NOT EXISTS (SELECT 1 FROM shop_goods_detail d WHERE d.goods_id=g.id);

INSERT INTO shop_goods (category_id, supplier, brand, type, name, description, unit, spec_type, status, sort, is_delete)
SELECT 526, '', '', 1, '缎面束口袋', '缎面束口，适合礼盒里再套一只小袋。', '个', 1, 1, 10, 0
FROM DUAL WHERE NOT EXISTS (SELECT 1 FROM shop_goods WHERE name='缎面束口袋' AND is_delete=0);

INSERT INTO shop_goods_image (goods_id, position, url, sort, is_delete)
SELECT g.id, 1, '/catalog/goods-pouch.jpg', 0, 0 FROM shop_goods g
WHERE g.name='缎面束口袋' AND g.is_delete=0
AND NOT EXISTS (SELECT 1 FROM shop_goods_image i WHERE i.goods_id=g.id AND i.position=1 AND i.is_delete=0);

INSERT INTO shop_goods_detail (goods_id, detail_desc, spec_desc, service_desc)
SELECT g.id, '绳色可换，袋身可丝印小标记。', '材质：缎面
结构：束口
常见尺寸：9×12 cm', '请给参考色号。' FROM shop_goods g
WHERE g.name='缎面束口袋' AND g.is_delete=0
AND NOT EXISTS (SELECT 1 FROM shop_goods_detail d WHERE d.goods_id=g.id);

INSERT INTO shop_goods (category_id, supplier, brand, type, name, description, unit, spec_type, status, sort, is_delete)
SELECT 526, '', '', 1, '耳环独立小袋', '一只耳环一个小袋，避免一对互相勾住。', '个', 1, 1, 10, 0
FROM DUAL WHERE NOT EXISTS (SELECT 1 FROM shop_goods WHERE name='耳环独立小袋' AND is_delete=0);

INSERT INTO shop_goods_image (goods_id, position, url, sort, is_delete)
SELECT g.id, 1, '/catalog/goods-pouch.jpg', 0, 0 FROM shop_goods g
WHERE g.name='耳环独立小袋' AND g.is_delete=0
AND NOT EXISTS (SELECT 1 FROM shop_goods_image i WHERE i.goods_id=g.id AND i.position=1 AND i.is_delete=0);

INSERT INTO shop_goods_detail (goods_id, detail_desc, spec_desc, service_desc)
SELECT g.id, '尺寸小，适合耳钉和耳钩分开装。', '材质：绒布
结构：抽绳
常见尺寸：5×7 cm', '可跟礼盒一起询盘。' FROM shop_goods g
WHERE g.name='耳环独立小袋' AND g.is_delete=0
AND NOT EXISTS (SELECT 1 FROM shop_goods_detail d WHERE d.goods_id=g.id);

INSERT INTO shop_goods (category_id, supplier, brand, type, name, description, unit, spec_type, status, sort, is_delete)
SELECT 527, '', '', 1, '粉饼磁吸盒', '圆格按粉饼直径做，盖子磁吸。', '个', 1, 1, 10, 0
FROM DUAL WHERE NOT EXISTS (SELECT 1 FROM shop_goods WHERE name='粉饼磁吸盒' AND is_delete=0);

INSERT INTO shop_goods_image (goods_id, position, url, sort, is_delete)
SELECT g.id, 1, '/catalog/goods-case.jpg', 0, 0 FROM shop_goods g
WHERE g.name='粉饼磁吸盒' AND g.is_delete=0
AND NOT EXISTS (SELECT 1 FROM shop_goods_image i WHERE i.goods_id=g.id AND i.position=1 AND i.is_delete=0);

INSERT INTO shop_goods_detail (goods_id, detail_desc, spec_desc, service_desc)
SELECT g.id, '粉饼直径不同，格位不要通用。', '材质：PU+纸板
结构：磁吸
常见直径：8 cm', '请提供粉饼直径。' FROM shop_goods g
WHERE g.name='粉饼磁吸盒' AND g.is_delete=0
AND NOT EXISTS (SELECT 1 FROM shop_goods_detail d WHERE d.goods_id=g.id);

INSERT INTO shop_goods (category_id, supplier, brand, type, name, description, unit, spec_type, status, sort, is_delete)
SELECT 527, '', '', 1, '口红收纳盒', '口红竖着插，一排的数量可以改。', '个', 1, 1, 10, 0
FROM DUAL WHERE NOT EXISTS (SELECT 1 FROM shop_goods WHERE name='口红收纳盒' AND is_delete=0);

INSERT INTO shop_goods_image (goods_id, position, url, sort, is_delete)
SELECT g.id, 1, '/catalog/goods-case.jpg', 0, 0 FROM shop_goods g
WHERE g.name='口红收纳盒' AND g.is_delete=0
AND NOT EXISTS (SELECT 1 FROM shop_goods_image i WHERE i.goods_id=g.id AND i.position=1 AND i.is_delete=0);

INSERT INTO shop_goods_detail (goods_id, detail_desc, spec_desc, service_desc)
SELECT g.id, '管径按口红实物留孔，避免太松。', '材质：纸板+植绒
结构：翻盖
常见：6 支或 12 支', '请写支数和管径。' FROM shop_goods g
WHERE g.name='口红收纳盒' AND g.is_delete=0
AND NOT EXISTS (SELECT 1 FROM shop_goods_detail d WHERE d.goods_id=g.id);

INSERT INTO shop_goods (category_id, supplier, brand, type, name, description, unit, spec_type, status, sort, is_delete)
SELECT 527, '', '', 1, '眼影书本盒', '像书一样翻开，盘面平放，盖里可以放镜子。', '个', 1, 1, 10, 0
FROM DUAL WHERE NOT EXISTS (SELECT 1 FROM shop_goods WHERE name='眼影书本盒' AND is_delete=0);

INSERT INTO shop_goods_image (goods_id, position, url, sort, is_delete)
SELECT g.id, 1, '/catalog/goods-case.jpg', 0, 0 FROM shop_goods g
WHERE g.name='眼影书本盒' AND g.is_delete=0
AND NOT EXISTS (SELECT 1 FROM shop_goods_image i WHERE i.goods_id=g.id AND i.position=1 AND i.is_delete=0);

INSERT INTO shop_goods_detail (goods_id, detail_desc, spec_desc, service_desc)
SELECT g.id, '内腔按眼影盘厚度留高。', '材质：灰板+PU
结构：书型
常见尺寸：14×10×2 cm', '请提供盘的长宽厚。' FROM shop_goods g
WHERE g.name='眼影书本盒' AND g.is_delete=0
AND NOT EXISTS (SELECT 1 FROM shop_goods_detail d WHERE d.goods_id=g.id);

INSERT INTO shop_goods (category_id, supplier, brand, type, name, description, unit, spec_type, status, sort, is_delete)
SELECT 528, '', '', 1, '便携刷具包', '刷具隔层分开，笔杆不跟粉饼挤在一起。', '个', 1, 1, 10, 0
FROM DUAL WHERE NOT EXISTS (SELECT 1 FROM shop_goods WHERE name='便携刷具包' AND is_delete=0);

INSERT INTO shop_goods_image (goods_id, position, url, sort, is_delete)
SELECT g.id, 1, '/catalog/goods-makeup.jpg', 0, 0 FROM shop_goods g
WHERE g.name='便携刷具包' AND g.is_delete=0
AND NOT EXISTS (SELECT 1 FROM shop_goods_image i WHERE i.goods_id=g.id AND i.position=1 AND i.is_delete=0);

INSERT INTO shop_goods_detail (goods_id, detail_desc, spec_desc, service_desc)
SELECT g.id, '刷具位按支数车线，数量可改。', '材质：PU
结构：拉链
常见尺寸：20×8×6 cm', '请说明刷具大概支数。' FROM shop_goods g
WHERE g.name='便携刷具包' AND g.is_delete=0
AND NOT EXISTS (SELECT 1 FROM shop_goods_detail d WHERE d.goods_id=g.id);

INSERT INTO shop_goods (category_id, supplier, brand, type, name, description, unit, spec_type, status, sort, is_delete)
SELECT 528, '', '', 1, '方形化妆包', '方型站得住，瓶罐不容易倒。', '个', 1, 1, 10, 0
FROM DUAL WHERE NOT EXISTS (SELECT 1 FROM shop_goods WHERE name='方形化妆包' AND is_delete=0);

INSERT INTO shop_goods_image (goods_id, position, url, sort, is_delete)
SELECT g.id, 1, '/catalog/goods-makeup.jpg', 0, 0 FROM shop_goods g
WHERE g.name='方形化妆包' AND g.is_delete=0
AND NOT EXISTS (SELECT 1 FROM shop_goods_image i WHERE i.goods_id=g.id AND i.position=1 AND i.is_delete=0);

INSERT INTO shop_goods_detail (goods_id, detail_desc, spec_desc, service_desc)
SELECT g.id, '底部加硬，打开后瓶罐还站着。', '材质：尼龙
结构：双拉链
常见尺寸：18×12×8 cm', '颜色可按品牌色做。' FROM shop_goods g
WHERE g.name='方形化妆包' AND g.is_delete=0
AND NOT EXISTS (SELECT 1 FROM shop_goods_detail d WHERE d.goods_id=g.id);

INSERT INTO shop_goods (category_id, supplier, brand, type, name, description, unit, spec_type, status, sort, is_delete)
SELECT 528, '', '', 1, '挂式洗漱包', '背后有挂钩，打开能挂在门后。', '个', 1, 1, 10, 0
FROM DUAL WHERE NOT EXISTS (SELECT 1 FROM shop_goods WHERE name='挂式洗漱包' AND is_delete=0);

INSERT INTO shop_goods_image (goods_id, position, url, sort, is_delete)
SELECT g.id, 1, '/catalog/goods-makeup.jpg', 0, 0 FROM shop_goods g
WHERE g.name='挂式洗漱包' AND g.is_delete=0
AND NOT EXISTS (SELECT 1 FROM shop_goods_image i WHERE i.goods_id=g.id AND i.position=1 AND i.is_delete=0);

INSERT INTO shop_goods_detail (goods_id, detail_desc, spec_desc, service_desc)
SELECT g.id, '挂钩承重按洗漱瓶考虑，不是空包。', '材质：PU
结构：挂式
常见尺寸：24×18×8 cm', '请说明大概要放几瓶。' FROM shop_goods g
WHERE g.name='挂式洗漱包' AND g.is_delete=0
AND NOT EXISTS (SELECT 1 FROM shop_goods_detail d WHERE d.goods_id=g.id);

INSERT INTO shop_goods (category_id, supplier, brand, type, name, description, unit, spec_type, status, sort, is_delete)
SELECT 529, '', '', 1, '13寸轻薄内胆', '给轻薄本用的内胆，可放进背包或单独拿。', '个', 1, 1, 10, 0
FROM DUAL WHERE NOT EXISTS (SELECT 1 FROM shop_goods WHERE name='13寸轻薄内胆' AND is_delete=0);

INSERT INTO shop_goods_image (goods_id, position, url, sort, is_delete)
SELECT g.id, 1, '/catalog/goods-laptop.jpg', 0, 0 FROM shop_goods g
WHERE g.name='13寸轻薄内胆' AND g.is_delete=0
AND NOT EXISTS (SELECT 1 FROM shop_goods_image i WHERE i.goods_id=g.id AND i.position=1 AND i.is_delete=0);

INSERT INTO shop_goods_detail (goods_id, detail_desc, spec_desc, service_desc)
SELECT g.id, '厚度按超薄本留，不做成 15 寸的厚度。', '材质：尼龙+珍珠棉
适用：13 寸
结构：顶部拉链', '请注明屏幕尺寸。' FROM shop_goods g
WHERE g.name='13寸轻薄内胆' AND g.is_delete=0
AND NOT EXISTS (SELECT 1 FROM shop_goods_detail d WHERE d.goods_id=g.id);

INSERT INTO shop_goods (category_id, supplier, brand, type, name, description, unit, spec_type, status, sort, is_delete)
SELECT 529, '', '', 1, '15.6寸手提电脑包', '有提手的电脑包，15.6 寸可单独拿，不必再套背包。', '个', 1, 1, 10, 0
FROM DUAL WHERE NOT EXISTS (SELECT 1 FROM shop_goods WHERE name='15.6寸手提电脑包' AND is_delete=0);

INSERT INTO shop_goods_image (goods_id, position, url, sort, is_delete)
SELECT g.id, 1, '/catalog/goods-laptop.jpg', 0, 0 FROM shop_goods g
WHERE g.name='15.6寸手提电脑包' AND g.is_delete=0
AND NOT EXISTS (SELECT 1 FROM shop_goods_image i WHERE i.goods_id=g.id AND i.position=1 AND i.is_delete=0);

INSERT INTO shop_goods_detail (goods_id, detail_desc, spec_desc, service_desc)
SELECT g.id, '前袋放电源，主仓只放电脑。', '材质：牛津布
适用：15.6 寸
结构：手提', '请确认是 15.6 还是 16 寸。' FROM shop_goods g
WHERE g.name='15.6寸手提电脑包' AND g.is_delete=0
AND NOT EXISTS (SELECT 1 FROM shop_goods_detail d WHERE d.goods_id=g.id);

INSERT INTO shop_goods (category_id, supplier, brand, type, name, description, unit, spec_type, status, sort, is_delete)
SELECT 529, '', '', 1, '平板内胆包', '按平板尺寸做的薄包，可放进手提包。', '个', 1, 1, 10, 0
FROM DUAL WHERE NOT EXISTS (SELECT 1 FROM shop_goods WHERE name='平板内胆包' AND is_delete=0);

INSERT INTO shop_goods_image (goods_id, position, url, sort, is_delete)
SELECT g.id, 1, '/catalog/goods-laptop.jpg', 0, 0 FROM shop_goods g
WHERE g.name='平板内胆包' AND g.is_delete=0
AND NOT EXISTS (SELECT 1 FROM shop_goods_image i WHERE i.goods_id=g.id AND i.position=1 AND i.is_delete=0);

INSERT INTO shop_goods_detail (goods_id, detail_desc, spec_desc, service_desc)
SELECT g.id, '四角加垫，避免和充电器挤在一起。', '材质：尼龙
适用：11 寸平板
结构：拉链', '请给平板型号或尺寸。' FROM shop_goods g
WHERE g.name='平板内胆包' AND g.is_delete=0
AND NOT EXISTS (SELECT 1 FROM shop_goods_detail d WHERE d.goods_id=g.id);

INSERT INTO shop_goods (category_id, supplier, brand, type, name, description, unit, spec_type, status, sort, is_delete)
SELECT 530, '', '', 1, '轻便通勤包', '比商务包更轻，仍留一个电脑隔层。', '个', 1, 1, 10, 0
FROM DUAL WHERE NOT EXISTS (SELECT 1 FROM shop_goods WHERE name='轻便通勤包' AND is_delete=0);

INSERT INTO shop_goods_image (goods_id, position, url, sort, is_delete)
SELECT g.id, 1, '/catalog/goods-backpack.jpg', 0, 0 FROM shop_goods g
WHERE g.name='轻便通勤包' AND g.is_delete=0
AND NOT EXISTS (SELECT 1 FROM shop_goods_image i WHERE i.goods_id=g.id AND i.position=1 AND i.is_delete=0);

INSERT INTO shop_goods_detail (goods_id, detail_desc, spec_desc, service_desc)
SELECT g.id, '肩带可收，适合每天地铁通勤。', '材质：尼龙
容量：约 16 L
电脑仓：14 寸', '颜色和织标可改。' FROM shop_goods g
WHERE g.name='轻便通勤包' AND g.is_delete=0
AND NOT EXISTS (SELECT 1 FROM shop_goods_detail d WHERE d.goods_id=g.id);

INSERT INTO shop_goods (category_id, supplier, brand, type, name, description, unit, spec_type, status, sort, is_delete)
SELECT 530, '', '', 1, '防泼水双肩包', '面料做了防泼水，短时下雨不用另套袋。', '个', 1, 1, 10, 0
FROM DUAL WHERE NOT EXISTS (SELECT 1 FROM shop_goods WHERE name='防泼水双肩包' AND is_delete=0);

INSERT INTO shop_goods_image (goods_id, position, url, sort, is_delete)
SELECT g.id, 1, '/catalog/goods-backpack.jpg', 0, 0 FROM shop_goods g
WHERE g.name='防泼水双肩包' AND g.is_delete=0
AND NOT EXISTS (SELECT 1 FROM shop_goods_image i WHERE i.goods_id=g.id AND i.position=1 AND i.is_delete=0);

INSERT INTO shop_goods_detail (goods_id, detail_desc, spec_desc, service_desc)
SELECT g.id, '缝线压胶，侧袋可以放伞。', '材质：防泼水尼龙
容量：约 20 L
电脑仓：15.6 寸', '请说明是否要反光条。' FROM shop_goods g
WHERE g.name='防泼水双肩包' AND g.is_delete=0
AND NOT EXISTS (SELECT 1 FROM shop_goods_detail d WHERE d.goods_id=g.id);

INSERT INTO shop_goods (category_id, supplier, brand, type, name, description, unit, spec_type, status, sort, is_delete)
SELECT 530, '', '', 1, '学生双肩包', '主仓能放书和文具，侧袋放水杯。', '个', 1, 1, 10, 0
FROM DUAL WHERE NOT EXISTS (SELECT 1 FROM shop_goods WHERE name='学生双肩包' AND is_delete=0);

INSERT INTO shop_goods_image (goods_id, position, url, sort, is_delete)
SELECT g.id, 1, '/catalog/goods-backpack.jpg', 0, 0 FROM shop_goods g
WHERE g.name='学生双肩包' AND g.is_delete=0
AND NOT EXISTS (SELECT 1 FROM shop_goods_image i WHERE i.goods_id=g.id AND i.position=1 AND i.is_delete=0);

INSERT INTO shop_goods_detail (goods_id, detail_desc, spec_desc, service_desc)
SELECT g.id, '背负调过，适合每天带书。', '材质：帆布
容量：约 18 L
电脑仓：14 寸', '颜色可按学校或品牌指定。' FROM shop_goods g
WHERE g.name='学生双肩包' AND g.is_delete=0
AND NOT EXISTS (SELECT 1 FROM shop_goods_detail d WHERE d.goods_id=g.id);

INSERT INTO shop_goods (category_id, supplier, brand, type, name, description, unit, spec_type, status, sort, is_delete)
SELECT 531, '', '', 1, '长款拉链钱包', '长款能放纸币和较多卡，拉链把一圈都封住。', '个', 1, 1, 10, 0
FROM DUAL WHERE NOT EXISTS (SELECT 1 FROM shop_goods WHERE name='长款拉链钱包' AND is_delete=0);

INSERT INTO shop_goods_image (goods_id, position, url, sort, is_delete)
SELECT g.id, 1, '/catalog/goods-wallet.jpg', 0, 0 FROM shop_goods g
WHERE g.name='长款拉链钱包' AND g.is_delete=0
AND NOT EXISTS (SELECT 1 FROM shop_goods_image i WHERE i.goods_id=g.id AND i.position=1 AND i.is_delete=0);

INSERT INTO shop_goods_detail (goods_id, detail_desc, spec_desc, service_desc)
SELECT g.id, '皮纹可选十字纹或荔枝纹。', '材质：PU
结构：长款拉链
常见尺寸：19×10×2 cm', '可跟礼盒一起做。' FROM shop_goods g
WHERE g.name='长款拉链钱包' AND g.is_delete=0
AND NOT EXISTS (SELECT 1 FROM shop_goods_detail d WHERE d.goods_id=g.id);

INSERT INTO shop_goods (category_id, supplier, brand, type, name, description, unit, spec_type, status, sort, is_delete)
SELECT 531, '', '', 1, '多卡位卡包', '只做卡位，不放大额纸币，方便口袋携带。', '个', 1, 1, 10, 0
FROM DUAL WHERE NOT EXISTS (SELECT 1 FROM shop_goods WHERE name='多卡位卡包' AND is_delete=0);

INSERT INTO shop_goods_image (goods_id, position, url, sort, is_delete)
SELECT g.id, 1, '/catalog/goods-wallet.jpg', 0, 0 FROM shop_goods g
WHERE g.name='多卡位卡包' AND g.is_delete=0
AND NOT EXISTS (SELECT 1 FROM shop_goods_image i WHERE i.goods_id=g.id AND i.position=1 AND i.is_delete=0);

INSERT INTO shop_goods_detail (goods_id, detail_desc, spec_desc, service_desc)
SELECT g.id, '卡位数量可以加减。', '材质：PU
结构：卡包
常见卡位：8 个', '请写需要几个卡位。' FROM shop_goods g
WHERE g.name='多卡位卡包' AND g.is_delete=0
AND NOT EXISTS (SELECT 1 FROM shop_goods_detail d WHERE d.goods_id=g.id);

INSERT INTO shop_goods (category_id, supplier, brand, type, name, description, unit, spec_type, status, sort, is_delete)
SELECT 531, '', '', 1, '两折零钱钱包', '对折后有零钱袋，卡位在另一面。', '个', 1, 1, 10, 0
FROM DUAL WHERE NOT EXISTS (SELECT 1 FROM shop_goods WHERE name='两折零钱钱包' AND is_delete=0);

INSERT INTO shop_goods_image (goods_id, position, url, sort, is_delete)
SELECT g.id, 1, '/catalog/goods-wallet.jpg', 0, 0 FROM shop_goods g
WHERE g.name='两折零钱钱包' AND g.is_delete=0
AND NOT EXISTS (SELECT 1 FROM shop_goods_image i WHERE i.goods_id=g.id AND i.position=1 AND i.is_delete=0);

INSERT INTO shop_goods_detail (goods_id, detail_desc, spec_desc, service_desc)
SELECT g.id, '零钱袋和卡位分开，互不干扰。', '材质：PU
结构：两折
常见尺寸：11×9×2 cm', '皮色可按样对。' FROM shop_goods g
WHERE g.name='两折零钱钱包' AND g.is_delete=0
AND NOT EXISTS (SELECT 1 FROM shop_goods_detail d WHERE d.goods_id=g.id);

INSERT INTO shop_goods (category_id, supplier, brand, type, name, description, unit, spec_type, status, sort, is_delete)
SELECT 532, '', '', 1, '登机手提包', '能放进座位下方的短途包，有手提和肩带。', '个', 1, 1, 10, 0
FROM DUAL WHERE NOT EXISTS (SELECT 1 FROM shop_goods WHERE name='登机手提包' AND is_delete=0);

INSERT INTO shop_goods_image (goods_id, position, url, sort, is_delete)
SELECT g.id, 1, '/catalog/goods-travel.jpg', 0, 0 FROM shop_goods g
WHERE g.name='登机手提包' AND g.is_delete=0
AND NOT EXISTS (SELECT 1 FROM shop_goods_image i WHERE i.goods_id=g.id AND i.position=1 AND i.is_delete=0);

INSERT INTO shop_goods_detail (goods_id, detail_desc, spec_desc, service_desc)
SELECT g.id, '主仓开口大，电脑可单独放一侧。', '材质：牛津布
容量：约 28 L
结构：手提+肩带', '请说明是否需要鞋仓。' FROM shop_goods g
WHERE g.name='登机手提包' AND g.is_delete=0
AND NOT EXISTS (SELECT 1 FROM shop_goods_detail d WHERE d.goods_id=g.id);

INSERT INTO shop_goods (category_id, supplier, brand, type, name, description, unit, spec_type, status, sort, is_delete)
SELECT 532, '', '', 1, '大容量旅行袋', '两到四天的行李袋，主仓可以摊开装。', '个', 1, 1, 10, 0
FROM DUAL WHERE NOT EXISTS (SELECT 1 FROM shop_goods WHERE name='大容量旅行袋' AND is_delete=0);

INSERT INTO shop_goods_image (goods_id, position, url, sort, is_delete)
SELECT g.id, 1, '/catalog/goods-travel.jpg', 0, 0 FROM shop_goods g
WHERE g.name='大容量旅行袋' AND g.is_delete=0
AND NOT EXISTS (SELECT 1 FROM shop_goods_image i WHERE i.goods_id=g.id AND i.position=1 AND i.is_delete=0);

INSERT INTO shop_goods_detail (goods_id, detail_desc, spec_desc, service_desc)
SELECT g.id, '底部加厚，拖着走时不磨破。', '材质：牛津布
容量：约 45 L
结构：双提手', '请说明要不要拉杆套。' FROM shop_goods g
WHERE g.name='大容量旅行袋' AND g.is_delete=0
AND NOT EXISTS (SELECT 1 FROM shop_goods_detail d WHERE d.goods_id=g.id);

INSERT INTO shop_goods (category_id, supplier, brand, type, name, description, unit, spec_type, status, sort, is_delete)
SELECT 532, '', '', 1, '行李收纳袋', '放进旅行箱里分装衣物，一套多只。', '个', 1, 1, 10, 0
FROM DUAL WHERE NOT EXISTS (SELECT 1 FROM shop_goods WHERE name='行李收纳袋' AND is_delete=0);

INSERT INTO shop_goods_image (goods_id, position, url, sort, is_delete)
SELECT g.id, 1, '/catalog/goods-travel.jpg', 0, 0 FROM shop_goods g
WHERE g.name='行李收纳袋' AND g.is_delete=0
AND NOT EXISTS (SELECT 1 FROM shop_goods_image i WHERE i.goods_id=g.id AND i.position=1 AND i.is_delete=0);

INSERT INTO shop_goods_detail (goods_id, detail_desc, spec_desc, service_desc)
SELECT g.id, '一套里大小不同，衣物和鞋分开。', '材质：尼龙
结构：拉链
常见：三件套', '可按套询盘。' FROM shop_goods g
WHERE g.name='行李收纳袋' AND g.is_delete=0
AND NOT EXISTS (SELECT 1 FROM shop_goods_detail d WHERE d.goods_id=g.id);
