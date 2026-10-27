-- 用实拍风格图片和正文填官网。不写虚假证书和成立年份。

UPDATE category SET image = '/catalog/goods-box.jpg' WHERE id = 525;
UPDATE category SET image = '/catalog/goods-pouch.jpg' WHERE id = 526;
UPDATE category SET image = '/catalog/goods-case.jpg' WHERE id = 527;
UPDATE category SET image = '/catalog/goods-makeup.jpg' WHERE id = 528;
UPDATE category SET image = '/catalog/goods-laptop.jpg' WHERE id = 529;
UPDATE category SET image = '/catalog/goods-backpack.jpg' WHERE id = 530;
UPDATE category SET image = '/catalog/goods-wallet.jpg' WHERE id = 531;
UPDATE category SET image = '/catalog/goods-travel.jpg' WHERE id = 532;

UPDATE shop_goods_image SET url = '/catalog/goods-box.jpg' WHERE goods_id = 17 AND position = 1 AND is_delete = 0;
UPDATE shop_goods_image SET url = '/catalog/goods-drawer.jpg' WHERE goods_id = 18 AND position = 1 AND is_delete = 0;
UPDATE shop_goods_image SET url = '/catalog/goods-pouch.jpg' WHERE goods_id = 19 AND position = 1 AND is_delete = 0;
UPDATE shop_goods_image SET url = '/catalog/goods-case.jpg' WHERE goods_id = 20 AND position = 1 AND is_delete = 0;
UPDATE shop_goods_image SET url = '/catalog/goods-makeup.jpg' WHERE goods_id = 21 AND position = 1 AND is_delete = 0;
UPDATE shop_goods_image SET url = '/catalog/goods-laptop.jpg' WHERE goods_id = 22 AND position = 1 AND is_delete = 0;
UPDATE shop_goods_image SET url = '/catalog/goods-backpack.jpg' WHERE goods_id = 23 AND position = 1 AND is_delete = 0;
UPDATE shop_goods_image SET url = '/catalog/goods-wallet.jpg' WHERE goods_id = 24 AND position = 1 AND is_delete = 0;
UPDATE shop_goods_image SET url = '/catalog/goods-travel.jpg' WHERE goods_id = 25 AND position = 1 AND is_delete = 0;

INSERT INTO shop_goods_image (goods_id, position, url, sort, is_delete)
SELECT 17, 0, '/catalog/goods-box-top.jpg', 10, 0
FROM DUAL
WHERE NOT EXISTS (
  SELECT 1 FROM shop_goods_image WHERE goods_id = 17 AND url = '/catalog/goods-box-top.jpg' AND is_delete = 0
);

UPDATE shop_banner SET image = '/catalog/banner-box.jpg', title = '首饰包装', url = '/products?categoryId=525' WHERE id = 20;
UPDATE shop_banner SET image = '/catalog/banner-bag.jpg', title = '出行箱包', url = '/products?categoryId=530' WHERE id = 21;
UPDATE shop_banner SET image = '/catalog/banner-ask.jpg', title = '留下采购需求', url = '/inquiry' WHERE id = 22;

UPDATE shop_goods SET description = '翻盖礼盒，植绒内托按戒指、耳环或吊坠开槽，适合柜台陈列和套装出货。' WHERE id = 17;
UPDATE shop_goods SET description = '抽屉抽出取放项链和手链，盒盖可压商标，适合柜台和礼盒套装。' WHERE id = 18;
UPDATE shop_goods SET description = '绒布抽绳袋，一件饰品一个袋子，颜色可按潘通号对色。' WHERE id = 19;
UPDATE shop_goods SET description = '磁吸开合，盒盖带镜，内格按粉饼和口红划分，适合彩妆套装。' WHERE id = 20;
UPDATE shop_goods SET description = '双拉链分层，外层放刷具，内层放瓶罐，适合短途携带。' WHERE id = 21;
UPDATE shop_goods SET description = '四周加厚的电脑内胆，14 寸和 16 寸可单独手提，也可放进背包。' WHERE id = 22;
UPDATE shop_goods SET description = '前袋快取，后背有电脑隔层，适合每天通勤。' WHERE id = 23;
UPDATE shop_goods SET description = '短款拉链钱包，卡位和零钱袋分开，皮纹和卡位数量可以改。' WHERE id = 24;
UPDATE shop_goods SET description = '两到三天的短途包，主仓开口大，侧袋可以放水杯。' WHERE id = 25;

UPDATE shop_goods_detail SET detail_desc = '盒身可选亮光或哑光纸，内托按饰品实际尺寸开槽。起订和打样按具体尺寸报价，询盘时请写清装什么、做多少。' WHERE goods_id = 17;

UPDATE cms_content SET cover = '/catalog/news-catalog.jpg',
  summary = '集亿玖泽在深圳，供应珠宝包装和出行箱包，采购可以直接按分类询盘。',
  body = '<p>深圳市集亿玖泽贸易有限公司在龙岗区平湖街道。目录里现在有珠宝首饰盒、珠宝首饰袋、化妆盒、化妆包、电脑包、背包、钱包和旅行包。</p><p>国内订单和进出口配套都接。看中某一款，在产品页留下数量、尺寸和发往哪里，我们按电话回复。</p><p>联系人韦桂新，电话 15989589813。地址是深圳市龙岗区平湖街道禾花社区平吉大道 321 号 B 栋 616。</p>'
WHERE id = 1;

UPDATE cms_content SET cover = '/catalog/news-velvet.jpg',
  summary = '先确认装什么、做多少、发往哪里，再定材料和打样。',
  body = '<p>包装和箱包不是先选一个好看的样子。先写清里面放什么、大概数量，以及货发到国内还是海外。</p><p>首饰盒要看是戒指、耳环还是项链，内托开槽不一样。化妆盒要看粉饼和口红的尺寸。箱包要看电脑尺寸和出行天数。</p><p>这三件事清楚之后，再谈纸张、植绒、PU 或尼龙，以及要不要打样。</p>'
WHERE id = 2;

UPDATE cms_content SET cover = '/catalog/banner-box.jpg',
  summary = '当前供应八类包装和箱包，并配合国内贸易与进出口业务。',
  body = '<p>网站上可以直接询盘的，是珠宝首饰盒、珠宝首饰袋、化妆盒、化妆包、电脑包、背包、钱包和旅行包。</p><p>公司经营范围包括国内贸易、货物及技术进出口。进出口订单需要的包装和箱包，可以按目的地一起问规格和数量。</p><p>具体资质文件按订单需要提供，网站上不列举未核验的证书编号。</p>'
WHERE id = 3;

UPDATE cms_content SET cover = '/catalog/banner-bag.jpg',
  summary = '八类产品放在同一个目录里，采购按分类查看后直接询盘。',
  body = '<p>现在把首饰包装和出行箱包放在同一本目录里。左边选分类，中间看图，看中了就留下电话和数量。</p><p>还没有对外公布的成立年份和获奖记录，就不写进页面。有新的产品或规格，在后台内容里更新。</p>'
WHERE id = 4;

UPDATE cms_content SET
  summary = '询盘里的联系方式和采购需求，只用于回复这一次询盘。',
  body = '<p>提交询盘时，我们会保存联系人、电话、公司、邮箱、目标市场、数量和需求说明。</p><p>这些信息只用于回复这次采购，不用于别的推广。后台可以把它标成已联系。</p><p>如果填写了询盘邮箱，并且服务器已经配置发信账号，新询盘会同时发到该邮箱。</p>'
WHERE id = 5;

UPDATE cms_content SET cover = '/catalog/news-export.jpg',
  summary = '装什么、做多少、发往哪里，这三件事决定材料和打样。',
  body = '<p>出口包装最容易返工的，是内托尺寸和数量没写清。戒指盒和项链盒不能共用一个开槽。</p><p>询盘时写上饰品类型、单款数量和目的地。国内仓和海外客户，对盒型和运输方式的要求不一样。</p><p>打样前先确认这三件事，比先做一批成品再改更省时间。</p>'
WHERE id = 8;

UPDATE cms_content SET cover = '/catalog/news-velvet.jpg',
  summary = '戒指、耳环和吊坠适合植绒开槽，项链要单独留出长度。',
  body = '<p>植绒内托适合固定小件饰品。戒指、耳环和吊坠可以按颗开槽，陈列时不容易晃。</p><p>项链和手链更适合抽屉盒，槽要按链长留空，不要硬塞进戒指盒。</p><p>颜色可以按潘通号对，询盘时写上参考色和大概尺寸即可。</p>'
WHERE id = 9;

UPDATE cms_content SET cover = '/catalog/news-catalog.jpg',
  summary = '八类产品可以按分类查看，并直接留下采购需求。',
  body = '<p>产品目录已经按八个分类打开。每一类都有封面、规格和询盘按钮。</p><p>采购不用先注册。看中某一款，点「询盘这款」，电话和数量写上就能提交。</p><p>公司动态和行业资料会继续放在资讯里。</p>'
WHERE id = 10;

INSERT INTO site_inquiry (goods_id, goods_name, name, company, phone, email, country, quantity, message, status, is_delete)
SELECT 17, '绒面首饰礼盒', '陈采购', '南宁示例饰品', '13800001111', 'buyer@example.com', '国内', '500 个', '想做戒指礼盒，内托按单戒开槽，先问起订和打样。', 0, 0
FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM site_inquiry WHERE phone = '13800001111');

INSERT INTO site_inquiry (goods_id, goods_name, name, company, phone, email, country, quantity, message, status, is_delete)
SELECT 23, '商务双肩包', '林小姐', '广州示例箱包行', '13900002222', 'lin@example.com', '东南亚', '300 个', '通勤背包要有 15.6 寸电脑仓，想先看现有结构和可改的地方。', 1, 0
FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM site_inquiry WHERE phone = '13900002222');

INSERT INTO sys_permission (
  permission_id, parent_id, permission_type, permission_name, icon,
  route_path, component_path, permission_code, sort, status, del_flag, create_time, create_by
)
SELECT 310, 0, 'menu', '官网', 'OfficeBuilding', '/site', '', 'site', 92, '0', '0', NOW(), 'admin'
FROM DUAL WHERE NOT EXISTS (SELECT 1 FROM sys_permission WHERE permission_id = 310);

INSERT INTO sys_permission (
  permission_id, parent_id, permission_type, permission_name, icon,
  route_path, component_path, permission_code, sort, status, del_flag, create_time, create_by
)
SELECT 311, 310, 'menu', '企业资料', 'EditPen', '/site/profile', '/views/site/profile/index.vue', 'site:profile', 20, '0', '0', NOW(), 'admin'
FROM DUAL WHERE NOT EXISTS (SELECT 1 FROM sys_permission WHERE permission_id = 311);

INSERT INTO sys_permission (
  permission_id, parent_id, permission_type, permission_name, icon,
  route_path, component_path, permission_code, sort, status, del_flag, create_time, create_by
)
SELECT 312, 310, 'menu', '询盘记录', 'ChatDotRound', '/site/inquiry', '/views/site/inquiry/index.vue', 'site:inquiry', 10, '0', '0', NOW(), 'admin'
FROM DUAL WHERE NOT EXISTS (SELECT 1 FROM sys_permission WHERE permission_id = 312);

INSERT INTO sys_role_permission (role_id, permission_id)
SELECT 1, 310 FROM DUAL WHERE NOT EXISTS (SELECT 1 FROM sys_role_permission WHERE role_id = 1 AND permission_id = 310);
INSERT INTO sys_role_permission (role_id, permission_id)
SELECT 1, 311 FROM DUAL WHERE NOT EXISTS (SELECT 1 FROM sys_role_permission WHERE role_id = 1 AND permission_id = 311);
INSERT INTO sys_role_permission (role_id, permission_id)
SELECT 1, 312 FROM DUAL WHERE NOT EXISTS (SELECT 1 FROM sys_role_permission WHERE role_id = 1 AND permission_id = 312);
