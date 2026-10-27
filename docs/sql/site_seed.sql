CREATE TABLE IF NOT EXISTS `site_profile` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `name` varchar(100) NOT NULL COMMENT '企业全称',
  `short_name` varchar(50) NOT NULL COMMENT '简称',
  `phone` varchar(30) NOT NULL DEFAULT '' COMMENT '联系电话',
  `contact_name` varchar(50) NOT NULL DEFAULT '' COMMENT '联系人',
  `address` varchar(255) NOT NULL DEFAULT '' COMMENT '地址',
  `email` varchar(100) NOT NULL DEFAULT '' COMMENT '邮箱',
  `summary` varchar(500) NOT NULL DEFAULT '' COMMENT '一句话介绍',
  `seo_title` varchar(120) NOT NULL DEFAULT '' COMMENT '站点标题',
  `seo_keywords` varchar(255) NOT NULL DEFAULT '' COMMENT '关键词',
  `seo_description` varchar(300) NOT NULL DEFAULT '' COMMENT '描述',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='官网企业资料';

CREATE TABLE IF NOT EXISTS `cms_content` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '内容ID',
  `channel_id` bigint(20) NOT NULL COMMENT '栏目ID',
  `title` varchar(150) NOT NULL COMMENT '标题',
  `summary` varchar(500) NOT NULL DEFAULT '' COMMENT '摘要',
  `cover` varchar(255) NOT NULL DEFAULT '' COMMENT '封面',
  `body` mediumtext COMMENT '正文',
  `seo_title` varchar(120) NOT NULL DEFAULT '' COMMENT 'SEO标题',
  `seo_keywords` varchar(255) NOT NULL DEFAULT '' COMMENT 'SEO关键词',
  `seo_description` varchar(300) NOT NULL DEFAULT '' COMMENT 'SEO描述',
  `sort` int(11) NOT NULL DEFAULT 0 COMMENT '排序',
  `status` tinyint(4) NOT NULL DEFAULT 1 COMMENT '状态(0草稿 1发布)',
  `is_delete` tinyint(4) NOT NULL DEFAULT 0,
  `publish_time` datetime NOT NULL COMMENT '发布时间',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_channel_id` (`channel_id`),
  KEY `idx_publish_time` (`publish_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='CMS内容';

CREATE TABLE IF NOT EXISTS `site_inquiry` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `goods_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '关联商品',
  `goods_name` varchar(100) NOT NULL DEFAULT '' COMMENT '商品名称快照',
  `name` varchar(50) NOT NULL COMMENT '联系人',
  `company` varchar(100) NOT NULL DEFAULT '' COMMENT '公司',
  `phone` varchar(30) NOT NULL COMMENT '电话',
  `email` varchar(100) NOT NULL DEFAULT '' COMMENT '邮箱',
  `country` varchar(50) NOT NULL DEFAULT '' COMMENT '目标市场',
  `quantity` varchar(50) NOT NULL DEFAULT '' COMMENT '数量',
  `message` varchar(1000) NOT NULL COMMENT '需求说明',
  `status` tinyint(4) NOT NULL DEFAULT 0 COMMENT '状态(0未读 1已读)',
  `is_delete` tinyint(4) NOT NULL DEFAULT 0,
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_goods_id` (`goods_id`),
  KEY `idx_create_time` (`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='官网询盘';

DELETE i FROM shop_goods_image i INNER JOIN shop_goods g ON g.id = i.goods_id WHERE g.brand = '集亿';
DELETE d FROM shop_goods_detail d INNER JOIN shop_goods g ON g.id = d.goods_id WHERE g.brand = '集亿';
DELETE s FROM shop_goods_sku s INNER JOIN shop_goods g ON g.id = s.goods_id WHERE g.brand = '集亿';
DELETE FROM shop_goods WHERE brand = '集亿';
DELETE FROM shop_banner WHERE title LIKE '集亿·%';
DELETE FROM category WHERE module = 'goods' AND name IN ('珠宝首饰盒','珠宝首饰袋','化妆盒','化妆包','电脑包','背包','钱包','旅行包');
DELETE FROM site_profile;
DELETE c FROM cms_content c INNER JOIN category h ON h.id = c.channel_id WHERE h.module = 'cms' AND JSON_UNQUOTE(JSON_EXTRACT(h.extra, '$.code')) IN ('about','culture','honor','history','privacy','news','company','industry');
DELETE FROM category WHERE module = 'cms' AND JSON_UNQUOTE(JSON_EXTRACT(extra, '$.code')) IN ('about','culture','honor','history','privacy','company','industry','news');

INSERT INTO site_profile (name, short_name, phone, contact_name, address, email, summary, seo_title, seo_keywords, seo_description)
VALUES (
  '深圳市集亿玖泽贸易有限公司',
  '集亿玖泽',
  '15989589813',
  '韦桂新',
  '深圳市龙岗区平湖街道禾花社区平吉大道321号B栋616',
  '',
  '供应珠宝首饰盒、首饰袋、化妆盒、化妆包、电脑包、背包、钱包和旅行包，承接国内贸易与进出口配套。',
  '集亿玖泽｜首饰包装与箱包',
  '首饰盒,首饰袋,化妆包,电脑包,背包,钱包,旅行包,集亿玖泽',
  '深圳市集亿玖泽贸易有限公司，供应珠宝包装与出行箱包，支持看样和询盘。'
);

INSERT INTO category (module, parent_id, name, icon, image, level, sort, status, is_delete) VALUES
('goods', 0, '珠宝首饰盒', '', '/catalog/cat-box.svg', 1, 80, 1, 0);
SET @c_box = LAST_INSERT_ID();
INSERT INTO category (module, parent_id, name, icon, image, level, sort, status, is_delete) VALUES
('goods', 0, '珠宝首饰袋', '', '/catalog/cat-pouch.svg', 1, 70, 1, 0);
SET @c_pouch = LAST_INSERT_ID();
INSERT INTO category (module, parent_id, name, icon, image, level, sort, status, is_delete) VALUES
('goods', 0, '化妆盒', '', '/catalog/cat-case.svg', 1, 60, 1, 0);
SET @c_case = LAST_INSERT_ID();
INSERT INTO category (module, parent_id, name, icon, image, level, sort, status, is_delete) VALUES
('goods', 0, '化妆包', '', '/catalog/cat-makeup.svg', 1, 50, 1, 0);
SET @c_makeup = LAST_INSERT_ID();
INSERT INTO category (module, parent_id, name, icon, image, level, sort, status, is_delete) VALUES
('goods', 0, '电脑包', '', '/catalog/cat-laptop.svg', 1, 40, 1, 0);
SET @c_laptop = LAST_INSERT_ID();
INSERT INTO category (module, parent_id, name, icon, image, level, sort, status, is_delete) VALUES
('goods', 0, '背包', '', '/catalog/cat-backpack.svg', 1, 30, 1, 0);
SET @c_backpack = LAST_INSERT_ID();
INSERT INTO category (module, parent_id, name, icon, image, level, sort, status, is_delete) VALUES
('goods', 0, '钱包', '', '/catalog/cat-wallet.svg', 1, 20, 1, 0);
SET @c_wallet = LAST_INSERT_ID();
INSERT INTO category (module, parent_id, name, icon, image, level, sort, status, is_delete) VALUES
('goods', 0, '旅行包', '', '/catalog/cat-travel.svg', 1, 10, 1, 0);
SET @c_travel = LAST_INSERT_ID();

INSERT INTO shop_goods (category_id, supplier, brand, type, name, description, unit, spec_type, status, sort, is_delete) VALUES
(@c_box, '深圳市集亿玖泽贸易有限公司', '集亿', 1, '绒面首饰礼盒', '翻盖结构，植绒内托，适合戒指、耳环和吊坠的柜台陈列。', '个', 1, 1, 100, 0);
SET @g1 = LAST_INSERT_ID();
INSERT INTO shop_goods_image (goods_id, position, url, sort, is_delete) VALUES (@g1, 1, '/catalog/p-box.svg', 1, 0);
INSERT INTO shop_goods_detail (goods_id, detail_desc, spec_desc, service_desc) VALUES
(@g1, '盒身可做亮光或哑光，内托按饰品尺寸开槽。适合珠宝柜台和礼品套装。', '材质：纸板+PU\n内里：植绒\n常见尺寸：9×9×4.5 cm\n工艺：烫金、丝印', '支持看样。数量、颜色和 logo 请在询盘里说明。');

INSERT INTO shop_goods (category_id, supplier, brand, type, name, description, unit, spec_type, status, sort, is_delete) VALUES
(@c_box, '深圳市集亿玖泽贸易有限公司', '集亿', 1, '抽屉式项链盒', '抽屉抽出，适合项链和手链，方便柜台取放。', '个', 1, 1, 90, 0);
SET @g2 = LAST_INSERT_ID();
INSERT INTO shop_goods_image (goods_id, position, url, sort, is_delete) VALUES (@g2, 1, '/catalog/p-drawer.svg', 1, 0);
INSERT INTO shop_goods_detail (goods_id, detail_desc, spec_desc, service_desc) VALUES
(@g2, '抽屉滑轨顺滑，盒盖可压 logo。适合套装和节日礼盒。', '材质：灰板+特种纸\n结构：抽屉\n常见尺寸：20×12×4 cm\n工艺：压纹、烫金', '按图打样，确认后再安排批量。');

INSERT INTO shop_goods (category_id, supplier, brand, type, name, description, unit, spec_type, status, sort, is_delete) VALUES
(@c_pouch, '深圳市集亿玖泽贸易有限公司', '集亿', 1, '绒布抽绳首饰袋', '柔软绒布，抽绳收口，用于耳环和吊坠的独立包装。', '个', 1, 1, 80, 0);
SET @g3 = LAST_INSERT_ID();
INSERT INTO shop_goods_image (goods_id, position, url, sort, is_delete) VALUES (@g3, 1, '/catalog/p-pouch.svg', 1, 0);
INSERT INTO shop_goods_detail (goods_id, detail_desc, spec_desc, service_desc) VALUES
(@g3, '颜色可按潘通号对色，袋口抽绳可换色。', '材质：植绒布\n结构：抽绳\n常见尺寸：8×10 cm\n工艺：丝印、织标', '小批量也可以询盘。');

INSERT INTO shop_goods (category_id, supplier, brand, type, name, description, unit, spec_type, status, sort, is_delete) VALUES
(@c_case, '深圳市集亿玖泽贸易有限公司', '集亿', 1, '磁吸化妆盒', '磁吸开合，内置镜面和隔层，适合彩妆套装。', '个', 1, 1, 70, 0);
SET @g4 = LAST_INSERT_ID();
INSERT INTO shop_goods_image (goods_id, position, url, sort, is_delete) VALUES (@g4, 1, '/catalog/p-case.svg', 1, 0);
INSERT INTO shop_goods_detail (goods_id, detail_desc, spec_desc, service_desc) VALUES
(@g4, '盒盖磁铁力度可调，内格按粉饼、口红尺寸划分。', '材质：PU+纸板\n结构：磁吸翻盖\n常见尺寸：18×12×5 cm\n工艺：压纹、烫金', '请提供产品尺寸，方便排内格。');

INSERT INTO shop_goods (category_id, supplier, brand, type, name, description, unit, spec_type, status, sort, is_delete) VALUES
(@c_makeup, '深圳市集亿玖泽贸易有限公司', '集亿', 1, '手提化妆包', '双拉链分层，适合旅行携带刷具和瓶罐。', '个', 1, 1, 60, 0);
SET @g5 = LAST_INSERT_ID();
INSERT INTO shop_goods_image (goods_id, position, url, sort, is_delete) VALUES (@g5, 1, '/catalog/p-makeup.svg', 1, 0);
INSERT INTO shop_goods_detail (goods_id, detail_desc, spec_desc, service_desc) VALUES
(@g5, '外层耐磨，内层可水洗。手柄可改成腕带。', '材质：PU\n结构：双层拉链\n常见尺寸：22×15×10 cm\n五金：枪色或金色', '可按品牌色做面料。');

INSERT INTO shop_goods (category_id, supplier, brand, type, name, description, unit, spec_type, status, sort, is_delete) VALUES
(@c_laptop, '深圳市集亿玖泽贸易有限公司', '集亿', 1, '防震电脑内胆包', '四周加厚，适合 14 寸和 16 寸笔记本单独携带。', '个', 1, 1, 50, 0);
SET @g6 = LAST_INSERT_ID();
INSERT INTO shop_goods_image (goods_id, position, url, sort, is_delete) VALUES (@g6, 1, '/catalog/p-laptop.svg', 1, 0);
INSERT INTO shop_goods_detail (goods_id, detail_desc, spec_desc, service_desc) VALUES
(@g6, '内胆独立，可放进背包或单独手提。', '材质：尼龙+珍珠棉\n适用：14 / 16 寸\n结构：顶部拉链\n工艺：织唛、车线', '请注明屏幕尺寸。');

INSERT INTO shop_goods (category_id, supplier, brand, type, name, description, unit, spec_type, status, sort, is_delete) VALUES
(@c_backpack, '深圳市集亿玖泽贸易有限公司', '集亿', 1, '商务双肩包', '前袋快取，后背有电脑隔层，适合通勤。', '个', 1, 1, 40, 0);
SET @g7 = LAST_INSERT_ID();
INSERT INTO shop_goods_image (goods_id, position, url, sort, is_delete) VALUES (@g7, 1, '/catalog/p-backpack.svg', 1, 0);
INSERT INTO shop_goods_detail (goods_id, detail_desc, spec_desc, service_desc) VALUES
(@g7, '肩带可调，背面有隐藏电脑仓。', '材质：帆布或尼龙\n容量：约 18 L\n电脑仓：15.6 寸\n工艺：织标、压胶', '颜色和五金可在询盘里指定。');

INSERT INTO shop_goods (category_id, supplier, brand, type, name, description, unit, spec_type, status, sort, is_delete) VALUES
(@c_wallet, '深圳市集亿玖泽贸易有限公司', '集亿', 1, '短款拉链钱包', '卡位清楚，中间有零钱袋，适合礼品和日常。', '个', 1, 1, 30, 0);
SET @g8 = LAST_INSERT_ID();
INSERT INTO shop_goods_image (goods_id, position, url, sort, is_delete) VALUES (@g8, 1, '/catalog/p-wallet.svg', 1, 0);
INSERT INTO shop_goods_detail (goods_id, detail_desc, spec_desc, service_desc) VALUES
(@g8, '皮料可选十字纹或荔枝纹，内页卡位数量可调。', '材质：PU\n结构：拉链\n常见尺寸：11×9×2 cm\n工艺：压印、烫金', '可配礼盒一起询盘。');

INSERT INTO shop_goods (category_id, supplier, brand, type, name, description, unit, spec_type, status, sort, is_delete) VALUES
(@c_travel, '深圳市集亿玖泽贸易有限公司', '集亿', 1, '短途旅行包', '仓位简单，适合两到三天的出行。', '个', 1, 1, 20, 0);
SET @g9 = LAST_INSERT_ID();
INSERT INTO shop_goods_image (goods_id, position, url, sort, is_delete) VALUES (@g9, 1, '/catalog/p-travel.svg', 1, 0);
INSERT INTO shop_goods_detail (goods_id, detail_desc, spec_desc, service_desc) VALUES
(@g9, '侧袋放水杯，主仓开口大，方便取放。', '材质：牛津布\n容量：约 30 L\n结构：手提+肩带\n工艺：车线、织标', '请说明是否需要拉杆。');

INSERT INTO shop_banner (title, image, url, position, type, sort, status, is_delete) VALUES
('集亿·首饰包装', '/catalog/banner-box.svg', '/products', 1, 3, 30, 1, 0),
('集亿·出行箱包', '/catalog/banner-bag.svg', '/products', 1, 3, 20, 1, 0),
('集亿·直接询盘', '/catalog/banner-ask.svg', '/inquiry', 1, 3, 10, 1, 0);

INSERT INTO category (module, parent_id, name, image, level, sort, status, is_delete, extra) VALUES
('cms', 0, '公司简介', '/catalog/page-about.svg', 1, 40, 1, 0, '{"code":"about","type":"page","summary":"集亿玖泽在深圳，做珠宝包装和出行箱包的供应与贸易。","seoTitle":"公司简介｜集亿玖泽","seoKeywords":"集亿玖泽,公司简介,箱包,首饰包装","seoDescription":"集亿玖泽公司简介，业务覆盖首饰包装与出行箱包。"}'),
('cms', 0, '做事方式', '/catalog/page-culture.svg', 1, 30, 1, 0, '{"code":"culture","type":"page","summary":"先看清楚尺寸和用途，再谈材料和数量。","seoTitle":"做事方式｜集亿玖泽","seoKeywords":"打样,定制包装,集亿玖泽","seoDescription":"集亿玖泽按尺寸、数量和目的地确认订单。"}'),
('cms', 0, '服务范围', '/catalog/page-honor.svg', 1, 20, 1, 0, '{"code":"honor","type":"page","summary":"国内供应，以及进出口业务需要的包装和箱包配套。","seoTitle":"服务范围｜集亿玖泽","seoKeywords":"进出口,包装供应,箱包","seoDescription":"集亿玖泽的产品与服务范围。"}'),
('cms', 0, '现在在做的事', '/catalog/page-history.svg', 1, 10, 1, 0, '{"code":"history","type":"page","summary":"把八条产品线放在同一个目录里，方便采购直接询盘。","seoTitle":"产品目录说明｜集亿玖泽","seoKeywords":"产品目录,询盘,集亿玖泽","seoDescription":"集亿玖泽按用途整理产品目录。"}'),
('cms', 0, '隐私说明', '', 1, 0, 1, 0, '{"code":"privacy","type":"page","summary":"询盘信息只用于回复这次采购需求。","seoTitle":"隐私说明｜集亿玖泽","seoKeywords":"隐私,询盘","seoDescription":"集亿玖泽询盘信息的使用说明。"}'),
('cms', 0, '资讯', '', 1, 50, 1, 0, '{"code":"news","type":"list","summary":"公司动态和行业资料","seoTitle":"资讯｜集亿玖泽","seoKeywords":"资讯,集亿玖泽","seoDescription":"集亿玖泽资讯"}');
SET @news = LAST_INSERT_ID();
INSERT INTO category (module, parent_id, name, image, level, sort, status, is_delete, extra) VALUES
('cms', @news, '公司动态', '', 2, 20, 1, 0, '{"code":"company","type":"list","summary":"","seoTitle":"","seoKeywords":"","seoDescription":""}'),
('cms', @news, '行业资料', '', 2, 10, 1, 0, '{"code":"industry","type":"list","summary":"","seoTitle":"","seoKeywords":"","seoDescription":""}');

INSERT INTO cms_content (channel_id, title, summary, cover, body, seo_title, seo_keywords, seo_description, sort, status, publish_time)
SELECT id, name,
  CASE JSON_UNQUOTE(JSON_EXTRACT(extra, '$.code'))
    WHEN 'about' THEN '集亿玖泽在深圳，做珠宝包装和出行箱包的供应与贸易。'
    WHEN 'culture' THEN '先看清楚尺寸和用途，再谈材料和数量。'
    WHEN 'honor' THEN '国内供应，以及进出口业务需要的包装和箱包配套。'
    WHEN 'history' THEN '把八条产品线放在同一个目录里，方便采购直接询盘。'
    ELSE '询盘信息只用于回复这次采购需求。'
  END,
  image, '',
  CASE JSON_UNQUOTE(JSON_EXTRACT(extra, '$.code'))
    WHEN 'about' THEN '公司简介｜集亿玖泽'
    WHEN 'culture' THEN '做事方式｜集亿玖泽'
    WHEN 'honor' THEN '服务范围｜集亿玖泽'
    WHEN 'history' THEN '产品目录说明｜集亿玖泽'
    ELSE '隐私说明｜集亿玖泽'
  END,
  CASE JSON_UNQUOTE(JSON_EXTRACT(extra, '$.code'))
    WHEN 'about' THEN '集亿玖泽,公司简介,箱包,首饰包装'
    WHEN 'culture' THEN '打样,定制包装,集亿玖泽'
    WHEN 'honor' THEN '进出口,包装供应,箱包'
    WHEN 'history' THEN '产品目录,询盘,集亿玖泽'
    ELSE '隐私,询盘'
  END,
  CASE JSON_UNQUOTE(JSON_EXTRACT(extra, '$.code'))
    WHEN 'about' THEN '集亿玖泽公司简介，业务覆盖首饰包装与出行箱包。'
    WHEN 'culture' THEN '集亿玖泽按尺寸、数量和目的地确认订单。'
    WHEN 'honor' THEN '集亿玖泽的产品与服务范围。'
    WHEN 'history' THEN '集亿玖泽按用途整理产品目录。'
    ELSE '集亿玖泽询盘信息的使用说明。'
  END,
  sort, 1, NOW()
FROM category WHERE module = 'cms' AND JSON_UNQUOTE(JSON_EXTRACT(extra, '$.code')) IN ('about','culture','honor','history','privacy');

UPDATE cms_content c JOIN category h ON h.id = c.channel_id SET c.body = '<p>深圳市集亿玖泽贸易有限公司位于龙岗区平湖。业务覆盖珠宝首饰盒、首饰袋、化妆盒、化妆包、电脑包、背包、钱包和旅行包。</p><p>我们同时做国内贸易，以及货物与技术进出口相关的配套供应。网站上的产品用于看款和询盘，具体规格以确认样为准。</p>' WHERE h.module = 'cms' AND JSON_UNQUOTE(JSON_EXTRACT(h.extra, '$.code')) = 'about';
UPDATE cms_content c JOIN category h ON h.id = c.channel_id SET c.body = '<p>包装和箱包最容易在尺寸、五金和颜色上返工。所以询盘时我们会先问三件事：装什么、做多少、发到哪里。</p><p>材料按用途选，不堆不必要的工艺。能打样的先打样，确认后再安排批量。</p>' WHERE h.module = 'cms' AND JSON_UNQUOTE(JSON_EXTRACT(h.extra, '$.code')) = 'culture';
UPDATE cms_content c JOIN category h ON h.id = c.channel_id SET c.body = '<p>目前可以对接的方向：</p><p>珠宝与美妆包装：首饰盒、首饰袋、化妆盒、化妆包。</p><p>出行产品：电脑包、背包、钱包、旅行包。</p><p>资质文件和合作案例由企业补充后展示在这里，不使用未经确认的证书。</p>' WHERE h.module = 'cms' AND JSON_UNQUOTE(JSON_EXTRACT(h.extra, '$.code')) = 'honor';
UPDATE cms_content c JOIN category h ON h.id = c.channel_id SET c.body = '<p>产品目录按用途分开：珠宝包装、美妆包装、数码与通勤、出行。</p><p>每一款都写了常见材质和尺寸。你看到的是可询盘的款式，不是网店库存。</p><p>公司沿革的具体年份，等企业资料确认后再补进这一页。</p>' WHERE h.module = 'cms' AND JSON_UNQUOTE(JSON_EXTRACT(h.extra, '$.code')) = 'history';
UPDATE cms_content c JOIN category h ON h.id = c.channel_id SET c.body = '<p>询盘表单会收集联系人、电话、公司、邮箱、目标市场、数量和需求说明。如果从某款产品进入，还会记下产品名称。</p><p>这些信息只用于和你联系、确认样品或报价，不用于与本次询盘无关的用途。</p><p>如需更正或删除，请拨打页面上的联系电话。</p>' WHERE h.module = 'cms' AND JSON_UNQUOTE(JSON_EXTRACT(h.extra, '$.code')) = 'privacy';

INSERT INTO cms_content (channel_id, title, summary, cover, body, seo_title, seo_keywords, seo_description, sort, status, publish_time)
SELECT id, '集亿玖泽开放产品目录询盘', '八类产品可以在网站上按分类查看，并直接留下采购需求。', '/catalog/news-3.svg', '网站现已按珠宝首饰盒、首饰袋、化妆盒、化妆包、电脑包、背包、钱包、旅行包分类。\n\n每款页面都有常见材质和尺寸。看到合适的款式，可以直接发起询盘，我们按你留下的电话回复。', '集亿玖泽开放产品目录询盘', '集亿玖泽,产品目录', '集亿玖泽产品目录已可在线询盘。', 10, 1, '2026-10-20 09:00:00' FROM category WHERE module = 'cms' AND JSON_UNQUOTE(JSON_EXTRACT(extra, '$.code')) = 'company';
INSERT INTO cms_content (channel_id, title, summary, cover, body, seo_title, seo_keywords, seo_description, sort, status, publish_time)
SELECT id, '出口包装先确认三件事', '装什么、做多少、发往哪里，这三件事决定材料和打样。', '/catalog/news-1.svg', '采购首饰盒或箱包时，建议先写清内装物尺寸、预计数量和目的地。\n\n尺寸决定内托和箱体，数量决定是否开模或换材料，目的地决定五金和包装方式。这三项写进询盘，回复会快很多。', '出口包装先确认三件事', '询盘,打样,包装', '询盘时建议写清尺寸、数量和目的地。', 30, 1, '2026-08-12 09:00:00' FROM category WHERE module = 'cms' AND JSON_UNQUOTE(JSON_EXTRACT(extra, '$.code')) = 'industry';
INSERT INTO cms_content (channel_id, title, summary, cover, body, seo_title, seo_keywords, seo_description, sort, status, publish_time)
SELECT id, '植绒内托适合哪些饰品', '戒指、耳环和吊坠更适合植绒开槽，项链要注意长度。', '/catalog/news-2.svg', '植绒内托能托住小件饰品，避免运输中互相摩擦。\n\n戒指和耳环按颗开槽即可。项链如果较长，需要加挂钩或加长盒身，不能只按吊坠尺寸开槽。询盘时附一张摆放草图会更准确。', '植绒内托适合哪些饰品', '植绒,首饰盒', '不同饰品对内托的要求不一样。', 20, 1, '2026-09-03 09:00:00' FROM category WHERE module = 'cms' AND JSON_UNQUOTE(JSON_EXTRACT(extra, '$.code')) = 'industry';
