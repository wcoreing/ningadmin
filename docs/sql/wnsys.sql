/*
 Navicat Premium Data Transfer

 Source Server         : xpmysql
 Source Server Type    : MySQL
 Source Server Version : 80012 (8.0.12)
 Source Host           : localhost:3306
 Source Schema         : wnsys

 Target Server Type    : MySQL
 Target Server Version : 80012 (8.0.12)
 File Encoding         : 65001

 Date: 06/02/2025 13:19:34
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for data_city
-- ----------------------------
DROP TABLE IF EXISTS `data_city`;
CREATE TABLE `data_city`  (
  `id` int(11) NOT NULL,
  `city_id` int(11) NOT NULL DEFAULT 0 COMMENT '城市id',
  `level` int(11) NOT NULL DEFAULT 0 COMMENT '省市级别',
  `parent_id` int(11) NOT NULL DEFAULT 0 COMMENT '父级id',
  `area_code` varchar(30) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '区号',
  `name` varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '名称',
  `merger_name` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '合并名称',
  `lng` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '经度',
  `lat` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '纬度',
  `is_show` tinyint(1) NOT NULL DEFAULT 1 COMMENT '是否展示',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '城市表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for form
-- ----------------------------
DROP TABLE IF EXISTS `form`;
CREATE TABLE `form`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '表单ID',
  `template_id` bigint(20) NOT NULL COMMENT '模板ID',
  `form_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '表单名称',
  `form_description` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '表单描述',
  `form_attributes` json NULL COMMENT '表单属性（可以存储一些自定义的表单属性，如表单的主题颜色、字体大小等）',
  `form_status` char(1) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT '0' COMMENT '表单状态（0正常 1停用）',
  `form_del_flag` char(1) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT '0' COMMENT '表单删除标志（0存在 2删除）',
  `form_create_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '表单创建者',
  `form_create_time` datetime NULL DEFAULT NULL COMMENT '表单创建时间',
  `form_update_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '表单更新者',
  `form_update_time` datetime NULL DEFAULT NULL COMMENT '表单更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `fk_form_template`(`template_id` ASC) USING BTREE,
  CONSTRAINT `fk_form_template` FOREIGN KEY (`template_id`) REFERENCES `form_template` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE = InnoDB AUTO_INCREMENT = 3 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '表单表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for form_data
-- ----------------------------
DROP TABLE IF EXISTS `form_data`;
CREATE TABLE `form_data`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '提交数据的唯一ID',
  `form_id` bigint(20) NOT NULL COMMENT '表单的ID，关联到相应的表单',
  `data` json NOT NULL COMMENT '用户提交的表单数据，以JSON格式存储，方便存储各种不同结构的数据',
  `submitter` json NULL COMMENT '提交者的信息，以JSON格式存储，可包含用户ID、用户名、邮箱等',
  `ip` varchar(45) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '提交者的IP地址，可用于记录提交者的网络位置信息',
  `status` char(1) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT '0' COMMENT '提交状态，可用于标记提交是否成功或其他状态，例如：0-成功，1-失败，2-处理中',
  `create_by` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '提交人的id',
  `create_time` datetime NOT NULL COMMENT '提交时间，用于记录用户提交表单的时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 6 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '表单提交数据存储表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for form_share
-- ----------------------------
DROP TABLE IF EXISTS `form_share`;
CREATE TABLE `form_share`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '分享ID',
  `template_id` bigint(20) NOT NULL COMMENT '模板ID',
  `share_code` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '分享码',
  `expire_time` datetime NULL DEFAULT NULL COMMENT '过期时间',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  `create_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '创建者',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_share_code`(`share_code` ASC) USING BTREE,
  INDEX `idx_template_id`(`template_id` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '表单分享记录' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for form_template
-- ----------------------------
DROP TABLE IF EXISTS `form_template`;
CREATE TABLE `form_template`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '模板ID',
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '模板名称',
  `description` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '模板描述',
  `submit_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '提交地址',
  `submit_method` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '提交方法(POST/GET)',
  `success_message` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '成功提示',
  `fail_message` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '失败提示',
  `form_items` json NULL COMMENT '表单项配置',
  `config` json NULL COMMENT '表单配置',
  `status` char(1) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT '0' COMMENT '状态（0正常 1停用）',
  `del_flag` char(1) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT '0' COMMENT '删除标志（0存在 2删除）',
  `create_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '创建者',
  `create_time` datetime NULL DEFAULT NULL COMMENT '创建时间',
  `update_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '更新者',
  `update_time` datetime NULL DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 13 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '表单模板表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for shop_activity
-- ----------------------------
DROP TABLE IF EXISTS `shop_activity`;
CREATE TABLE `shop_activity`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '活动ID',
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '活动名称',
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '活动描述',
  `image` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '活动图片',
  `start_time` datetime NOT NULL COMMENT '开始时间',
  `end_time` datetime NOT NULL COMMENT '结束时间',
  `type` tinyint(4) NOT NULL COMMENT '活动类型(1:限时特惠 2:新品首发 3:主题活动)',
  `status` tinyint(4) NOT NULL DEFAULT 0 COMMENT '活动状态(0:未开始 1:进行中 2:已结束)',
  `sort` int(11) NOT NULL DEFAULT 0 COMMENT '排序',
  `is_delete` tinyint(4) NOT NULL DEFAULT 0 COMMENT '是否删除',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 5 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '活动表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for shop_activity_goods
-- ----------------------------
DROP TABLE IF EXISTS `shop_activity_goods`;
CREATE TABLE `shop_activity_goods`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '活动商品ID',
  `activity_id` bigint(20) NOT NULL COMMENT '活动ID',
  `goods_id` bigint(20) NOT NULL COMMENT '商品ID',
  `activity_price` decimal(10, 2) NOT NULL COMMENT '活动价格',
  `activity_stock` int(11) NOT NULL COMMENT '活动库存',
  `sort` int(11) NOT NULL DEFAULT 0 COMMENT '排序',
  `is_delete` tinyint(4) NOT NULL DEFAULT 0 COMMENT '是否删除',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_activity_id`(`activity_id` ASC) USING BTREE,
  INDEX `idx_goods_id`(`goods_id` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 14 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '活动商品表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for shop_banner
-- ----------------------------
DROP TABLE IF EXISTS `shop_banner`;
CREATE TABLE `shop_banner`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '轮播图ID',
  `title` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '标题',
  `image` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '图片地址',
  `url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '跳转链接',
  `position` tinyint(4) NOT NULL DEFAULT 1 COMMENT '位置(1:首页 2:分类页 3:活动页)',
  `type` tinyint(4) NOT NULL DEFAULT 1 COMMENT '类型(1:商品 2:活动 3:外链)',
  `target_id` bigint(20) NULL DEFAULT NULL COMMENT '目标ID',
  `sort` int(11) NOT NULL DEFAULT 0 COMMENT '排序',
  `status` tinyint(4) NOT NULL DEFAULT 1 COMMENT '状态(0:禁用 1:启用)',
  `is_delete` tinyint(4) NOT NULL DEFAULT 0 COMMENT '是否删除',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 6 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '轮播图表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for shop_brand
-- ----------------------------
DROP TABLE IF EXISTS `shop_brand`;
CREATE TABLE `shop_brand`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '品牌ID',
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '品牌名称',
  `logo` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT '' COMMENT '品牌logo',
  `description` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT '' COMMENT '品牌描述',
  `website` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT '' COMMENT '品牌官网',
  `sort` int(11) NOT NULL DEFAULT 0 COMMENT '排序',
  `status` tinyint(4) NOT NULL DEFAULT 1 COMMENT '状态(0:禁用 1:启用)',
  `is_delete` tinyint(4) NOT NULL DEFAULT 0 COMMENT '是否删除',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 6 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '商品品牌表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for shop_cart
-- ----------------------------
DROP TABLE IF EXISTS `shop_cart`;
CREATE TABLE `shop_cart`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '购物车ID',
  `user_id` bigint(20) NOT NULL COMMENT '用户ID',
  `goods_id` bigint(20) NOT NULL COMMENT '商品ID',
  `sku_code` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT 'sku编码',
  `specs_json` json NOT NULL COMMENT '规格JSON',
  `quantity` int(11) NOT NULL COMMENT '数量',
  `price` decimal(10, 2) NULL DEFAULT NULL COMMENT '加入购物车时的价格',
  `selected` tinyint(11) NULL DEFAULT NULL COMMENT '是否选中',
  `is_delete` tinyint(4) NOT NULL DEFAULT 0 COMMENT '是否删除',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_user_id`(`user_id` ASC) USING BTREE,
  INDEX `idx_goods_id`(`goods_id` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 4 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '购物车表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for shop_category
-- ----------------------------
DROP TABLE IF EXISTS `shop_category`;
CREATE TABLE `shop_category`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '分类ID',
  `parent_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '父级ID',
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '分类名称',
  `icon` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '分类图标',
  `image` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '分类图片',
  `level` tinyint(4) NOT NULL DEFAULT 1 COMMENT '层级(1:一级 2:二级 3:三级)',
  `sort` int(11) NOT NULL DEFAULT 0 COMMENT '排序',
  `status` tinyint(3) UNSIGNED NOT NULL DEFAULT 1 COMMENT '状态(0:禁用 1:启用)',
  `is_delete` tinyint(4) NOT NULL DEFAULT 0 COMMENT '是否删除',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_parent_id`(`parent_id` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 505 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '商品分类表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for shop_goods
-- ----------------------------
DROP TABLE IF EXISTS `shop_goods`;
CREATE TABLE `shop_goods`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '商品ID',
  `category_id` bigint(20) NOT NULL COMMENT '分类ID',
  `supplier_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '供应商ID',
  `brand_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '品牌ID',
  `type` tinyint(4) NOT NULL DEFAULT 1 COMMENT '商品类型1普通商品2虚拟商品',
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '商品名称',
  `detail_desc` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '商品详情',
  `price` decimal(10, 2) NOT NULL COMMENT '售价',
  `original_price` decimal(10, 2) NOT NULL COMMENT '原价',
  `cost_price` decimal(10, 2) NOT NULL COMMENT '成本价',
  `stock` int(11) NOT NULL DEFAULT 0 COMMENT '库存',
  `warn_stock` int(11) NOT NULL DEFAULT 0 COMMENT '库存预警值',
  `unit` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '单位',
  `weight` decimal(10, 2) NOT NULL DEFAULT 0.00 COMMENT '重量(kg)',
  `volume` decimal(10, 2) NOT NULL DEFAULT 0.00 COMMENT '体积(m³)',
  `sales` int(11) NOT NULL DEFAULT 0 COMMENT '销量',
  `virtual_sales` int(11) NOT NULL DEFAULT 0 COMMENT '虚拟销量',
  `spec_type` tinyint(4) NOT NULL DEFAULT 1 COMMENT '规格类型：1-单规格 2-多规格',
  `status` tinyint(4) NOT NULL DEFAULT 0 COMMENT '状态(0:下架 1:上架)',
  `sort` int(11) NOT NULL DEFAULT 0 COMMENT '排序',
  `is_delete` tinyint(4) NOT NULL DEFAULT 0 COMMENT '是否删除',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_category_id`(`category_id` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 5006 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '商品表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for shop_goods_collect
-- ----------------------------
DROP TABLE IF EXISTS `shop_goods_collect`;
CREATE TABLE `shop_goods_collect`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '收藏ID',
  `goods_id` bigint(20) NOT NULL COMMENT '商品ID',
  `user_id` bigint(20) NOT NULL COMMENT '用户ID',
  `is_delete` tinyint(4) NOT NULL DEFAULT 0 COMMENT '是否删除',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_user_goods`(`user_id` ASC, `goods_id` ASC) USING BTREE,
  INDEX `idx_goods_id`(`goods_id` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '商品收藏表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for shop_goods_comment
-- ----------------------------
DROP TABLE IF EXISTS `shop_goods_comment`;
CREATE TABLE `shop_goods_comment`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '评论ID',
  `goods_id` bigint(20) NOT NULL COMMENT '商品ID',
  `user_id` bigint(20) NOT NULL COMMENT '用户ID',
  `nickname` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '用户昵称',
  `avatar` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '用户头像',
  `rate` tinyint(4) NOT NULL COMMENT '评分(1-5)',
  `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '评论内容',
  `images` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '评论图片，JSON数组',
  `is_delete` tinyint(4) NOT NULL DEFAULT 0 COMMENT '是否删除',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_goods_id`(`goods_id` ASC) USING BTREE,
  INDEX `idx_user_id`(`user_id` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 41 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '商品评论表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for shop_goods_detail
-- ----------------------------
DROP TABLE IF EXISTS `shop_goods_detail`;
CREATE TABLE `shop_goods_detail`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `goods_id` bigint(20) NOT NULL,
  `detail_desc` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `video_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT '' COMMENT '主视频',
  `spec_desc` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '规格说明',
  `service_desc` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '服务说明',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_goods_id`(`goods_id` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 3 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '商品详情表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for shop_goods_image
-- ----------------------------
DROP TABLE IF EXISTS `shop_goods_image`;
CREATE TABLE `shop_goods_image`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '图片ID',
  `goods_id` bigint(20) NOT NULL COMMENT '商品ID',
  `position` tinyint(3) UNSIGNED NOT NULL DEFAULT 0 COMMENT '图片位置(0:普通图片 1:主图 2:列表图 3:详情主图 4:分享图)',
  `url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '图片地址',
  `sort` int(11) NOT NULL DEFAULT 0 COMMENT '排序',
  `is_delete` tinyint(4) NOT NULL DEFAULT 0 COMMENT '是否删除',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_goods_id`(`goods_id` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 298 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '商品图片表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for shop_goods_sku
-- ----------------------------
DROP TABLE IF EXISTS `shop_goods_sku`;
CREATE TABLE `shop_goods_sku`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'SKU ID',
  `goods_id` bigint(20) NOT NULL COMMENT '商品ID',
  `sku_code` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT 'sku编码',
  `specs_json` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '规格值字符串',
  `image` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT 'SKU图片',
  `price` decimal(10, 2) NOT NULL COMMENT '售价',
  `original_price` decimal(10, 2) NOT NULL COMMENT '原价',
  `cost_price` decimal(10, 2) NOT NULL COMMENT '成本价',
  `stock` int(11) NOT NULL COMMENT '库存',
  `code` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '商品编码',
  `weight` decimal(10, 2) NULL DEFAULT NULL COMMENT '重量(kg)',
  `volume` decimal(10, 2) NULL DEFAULT NULL COMMENT '体积(m³)',
  `sort` int(11) NOT NULL COMMENT '排序',
  `status` tinyint(4) NOT NULL COMMENT '状态(0:禁用 1:启用)',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  `update_time` datetime NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = MyISAM AUTO_INCREMENT = 57 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '商品SKU表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for shop_goods_spec
-- ----------------------------
DROP TABLE IF EXISTS `shop_goods_spec`;
CREATE TABLE `shop_goods_spec`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '规格ID',
  `goods_id` bigint(20) NOT NULL COMMENT '商品ID',
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '规格名称',
  `sort` int(11) NOT NULL DEFAULT 0 COMMENT '排序',
  `is_delete` tinyint(4) NOT NULL DEFAULT 0 COMMENT '是否删除',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_goods_id`(`goods_id` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 362 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '商品规格表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for shop_goods_spec_value
-- ----------------------------
DROP TABLE IF EXISTS `shop_goods_spec_value`;
CREATE TABLE `shop_goods_spec_value`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `goods_id` bigint(20) NOT NULL COMMENT '商品ID',
  `spec_id` bigint(20) NOT NULL COMMENT '规格组ID',
  `value` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '规格值',
  `sort` int(11) NULL DEFAULT 0 COMMENT '排序',
  `created_at` datetime NULL DEFAULT NULL,
  `updated_at` datetime NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_goods_id`(`goods_id` ASC) USING BTREE,
  INDEX `idx_spec_group_id`(`spec_id` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 326 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '商品规格值' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for shop_goods_tag
-- ----------------------------
DROP TABLE IF EXISTS `shop_goods_tag`;
CREATE TABLE `shop_goods_tag`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '标签ID',
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '标签名称',
  `sort` int(11) NOT NULL DEFAULT 0 COMMENT '排序',
  `is_delete` tinyint(4) NOT NULL DEFAULT 0 COMMENT '是否删除',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 5 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '商品标签表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for shop_goods_tag_relation
-- ----------------------------
DROP TABLE IF EXISTS `shop_goods_tag_relation`;
CREATE TABLE `shop_goods_tag_relation`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `goods_id` int(10) UNSIGNED NOT NULL COMMENT '商品ID',
  `tag_id` int(10) UNSIGNED NOT NULL COMMENT '标签ID',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_goods_tag`(`goods_id` ASC, `tag_id` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 177 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '商品标签关联表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for shop_order
-- ----------------------------
DROP TABLE IF EXISTS `shop_order`;
CREATE TABLE `shop_order`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '订单ID',
  `order_no` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '订单编号',
  `user_id` bigint(20) NOT NULL COMMENT '用户ID',
  `total_amount` decimal(10, 2) NOT NULL COMMENT '订单总金额',
  `pay_amount` decimal(10, 2) NOT NULL COMMENT '实付金额',
  `freight_amount` decimal(10, 2) NOT NULL COMMENT '运费金额',
  `discount_amount` decimal(10, 2) NOT NULL DEFAULT 0.00 COMMENT '优惠金额',
  `pay_type` tinyint(4) NULL DEFAULT NULL COMMENT '支付方式(1:微信 2:支付宝 3:余额)',
  `pay_time` datetime NULL DEFAULT NULL COMMENT '支付时间',
  `ship_time` datetime NULL DEFAULT NULL COMMENT '发货时间',
  `receive_time` datetime NULL DEFAULT NULL COMMENT '收货时间',
  `finish_time` datetime NULL DEFAULT NULL COMMENT '完成时间',
  `cancel_time` datetime NULL DEFAULT NULL COMMENT '取消时间',
  `status` tinyint(4) NOT NULL DEFAULT 0 COMMENT '订单状态(0:待付款 1:待发货 2:待收货 3:已完成 4:已取消)',
  `cancel_reason` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '取消原因',
  `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '订单备注',
  `is_delete` tinyint(4) NOT NULL DEFAULT 0 COMMENT '是否删除',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_order_no`(`order_no` ASC) USING BTREE,
  INDEX `idx_user_id`(`user_id` ASC) USING BTREE,
  INDEX `idx_status`(`status` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 20 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '订单表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for shop_order_goods
-- ----------------------------
DROP TABLE IF EXISTS `shop_order_goods`;
CREATE TABLE `shop_order_goods`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '订单商品ID',
  `order_id` bigint(20) NOT NULL COMMENT '订单ID',
  `goods_id` bigint(20) NOT NULL COMMENT '商品ID',
  `goods_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '商品名称',
  `goods_image` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '商品图片',
  `sku_code` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT 'SKU编码',
  `specs` json NOT NULL COMMENT '规格JSON',
  `price` decimal(10, 2) NOT NULL COMMENT '商品单价',
  `quantity` int(11) NOT NULL COMMENT '购买数量',
  `total_amount` decimal(10, 2) NOT NULL COMMENT '商品总金额',
  `is_comment` tinyint(4) NOT NULL DEFAULT 0 COMMENT '是否已评价',
  `is_delete` tinyint(4) NOT NULL DEFAULT 0 COMMENT '是否删除',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_order_id`(`order_id` ASC) USING BTREE,
  INDEX `idx_goods_id`(`goods_id` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 26 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '订单商品表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for shop_order_log
-- ----------------------------
DROP TABLE IF EXISTS `shop_order_log`;
CREATE TABLE `shop_order_log`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '日志ID',
  `order_id` bigint(20) NOT NULL COMMENT '订单ID',
  `order_no` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '订单编号',
  `operator` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '操作人',
  `operator_type` tinyint(4) NOT NULL COMMENT '操作人类型(1:用户 2:管理员)',
  `operator_id` bigint(20) NOT NULL COMMENT '操作人ID',
  `action` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '操作行为',
  `status` tinyint(4) NOT NULL COMMENT '订单状态',
  `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '备注',
  `ip` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '操作IP',
  `is_delete` tinyint(4) NOT NULL DEFAULT 0 COMMENT '是否删除',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_order_id`(`order_id` ASC) USING BTREE,
  INDEX `idx_order_no`(`order_no` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '订单日志表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for shop_order_payment
-- ----------------------------
DROP TABLE IF EXISTS `shop_order_payment`;
CREATE TABLE `shop_order_payment`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '支付ID',
  `order_id` bigint(20) NOT NULL COMMENT '订单ID',
  `order_no` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '订单编号',
  `payment_no` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '支付单号',
  `transaction_id` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '第三方支付流水号',
  `pay_type` tinyint(4) NOT NULL COMMENT '支付方式(1:微信 2:支付宝 3:余额)',
  `pay_amount` decimal(10, 2) NOT NULL COMMENT '支付金额',
  `pay_status` tinyint(4) NOT NULL DEFAULT 0 COMMENT '支付状态(0:待支付 1:支付成功 2:支付失败)',
  `pay_time` datetime NULL DEFAULT NULL COMMENT '支付时间',
  `is_delete` tinyint(4) NOT NULL DEFAULT 0 COMMENT '是否删除',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_payment_no`(`payment_no` ASC) USING BTREE,
  INDEX `idx_order_id`(`order_id` ASC) USING BTREE,
  INDEX `idx_order_no`(`order_no` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '订单支付表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for shop_order_receiver
-- ----------------------------
DROP TABLE IF EXISTS `shop_order_receiver`;
CREATE TABLE `shop_order_receiver`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '地址ID',
  `order_id` bigint(20) NOT NULL COMMENT '订单ID',
  `receiver` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '收货人',
  `phone` varchar(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '联系电话',
  `province` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '省份',
  `city` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '城市',
  `district` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '区/县',
  `detail` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '详细地址',
  `is_delete` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否删除',
  `logistics_company` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '物流公司',
  `logistics_no` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '物流单号',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_order_id`(`order_id` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 15 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '订单地址表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for shop_payment_log
-- ----------------------------
DROP TABLE IF EXISTS `shop_payment_log`;
CREATE TABLE `shop_payment_log`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '支付记录ID',
  `order_id` bigint(20) NOT NULL COMMENT '订单ID',
  `order_no` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '订单编号',
  `pay_type` tinyint(4) NOT NULL COMMENT '支付方式(1:微信 2:支付宝 3:余额)',
  `trade_no` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '支付平台交易号',
  `pay_amount` decimal(10, 2) NOT NULL COMMENT '支付金额',
  `pay_status` tinyint(4) NOT NULL COMMENT '支付状态(0:未支付 1:支付成功 2:支付失败)',
  `pay_time` datetime NULL DEFAULT NULL COMMENT '支付时间',
  `notify_time` datetime NULL DEFAULT NULL COMMENT '回调时间',
  `notify_data` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '回调数据',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_order_no`(`order_no` ASC) USING BTREE,
  INDEX `idx_trade_no`(`trade_no` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '支付记录表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for shop_supplier
-- ----------------------------
DROP TABLE IF EXISTS `shop_supplier`;
CREATE TABLE `shop_supplier`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '供应商ID',
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '供应商名称',
  `contact` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '联系人',
  `phone` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '联系电话',
  `email` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT '' COMMENT '邮箱',
  `address` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '地址',
  `business_license` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT '' COMMENT '营业执照',
  `description` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT '' COMMENT '供应商描述',
  `level` tinyint(4) NOT NULL DEFAULT 1 COMMENT '供应商等级(1:普通 2:金牌 3:钻石)',
  `sort` int(11) NOT NULL DEFAULT 0 COMMENT '排序',
  `status` tinyint(4) NOT NULL DEFAULT 1 COMMENT '状态(0:禁用 1:启用)',
  `is_delete` tinyint(4) NOT NULL DEFAULT 0 COMMENT '是否删除',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_level`(`level` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 6 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '供应商表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for shop_user
-- ----------------------------
DROP TABLE IF EXISTS `shop_user`;
CREATE TABLE `shop_user`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '用户ID',
  `username` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '用户名',
  `password` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '密码',
  `nickname` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '昵称',
  `avatar` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '头像',
  `phone` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '手机号',
  `email` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '邮箱',
  `balance` decimal(10, 2) NOT NULL DEFAULT 0.00 COMMENT '余额',
  `points` int(11) NOT NULL DEFAULT 0 COMMENT '积分',
  `level` tinyint(4) NOT NULL DEFAULT 1 COMMENT '会员等级(1:普通会员 2:黄金会员 3:钻石会员)',
  `total_amount` decimal(10, 2) NOT NULL DEFAULT 0.00 COMMENT '累计消费金额',
  `last_login_time` datetime NULL DEFAULT NULL COMMENT '最后登录时间',
  `last_login_ip` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '最后登录IP',
  `status` tinyint(4) NOT NULL DEFAULT 1 COMMENT '状态(0:禁用 1:正常)',
  `is_delete` tinyint(4) NOT NULL DEFAULT 0 COMMENT '是否删除',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_username`(`username` ASC) USING BTREE,
  UNIQUE INDEX `uk_phone`(`phone` ASC) USING BTREE,
  UNIQUE INDEX `uk_email`(`email` ASC) USING BTREE,
  INDEX `idx_level`(`level` ASC) USING BTREE,
  INDEX `idx_points`(`points` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '用户表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for shop_user_address
-- ----------------------------
DROP TABLE IF EXISTS `shop_user_address`;
CREATE TABLE `shop_user_address`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '地址ID',
  `user_id` bigint(20) NOT NULL COMMENT '用户ID',
  `receiver_name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '收货人姓名',
  `receiver_mobile` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '收货人手机号',
  `province` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '省份',
  `city` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '城市',
  `district` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '区县',
  `detail` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '详细地址',
  `is_default` tinyint(4) NOT NULL DEFAULT 0 COMMENT '是否默认',
  `is_delete` tinyint(4) NOT NULL DEFAULT 0 COMMENT '是否删除',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_user_id`(`user_id` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 4 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '用户地址表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for sys_config
-- ----------------------------
DROP TABLE IF EXISTS `sys_config`;
CREATE TABLE `sys_config`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '配置ID',
  `config_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '配置名称',
  `config_key` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '配置键名',
  `config_value` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '配置键值',
  `config_type` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '配置类型',
  `module_key` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '所属模块标识',
  `module_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '模块名称',
  `remark` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '备注',
  `del_flag` char(1) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT '0' COMMENT '删除标志（0代表存在 2代表删除）',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_module_key`(`module_key` ASC, `config_key` ASC) USING BTREE COMMENT '模块下配置键唯一'
) ENGINE = InnoDB AUTO_INCREMENT = 16 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '系统配置表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for sys_file
-- ----------------------------
DROP TABLE IF EXISTS `sys_file`;
CREATE TABLE `sys_file`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '文件ID',
  `category_id` bigint(20) NULL DEFAULT NULL COMMENT '分类ID',
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '文件名',
  `original_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '原始文件名',
  `url` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '文件URL',
  `path` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '文件路径',
  `size` bigint(20) NOT NULL COMMENT '文件大小(字节)',
  `type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '文件类型',
  `md5` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '文件MD5',
  `status` tinyint(1) NULL DEFAULT 1 COMMENT '状态(0:禁用 1:正常)',
  `create_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '创建者',
  `create_time` datetime NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_category_id`(`category_id` ASC) USING BTREE,
  INDEX `idx_md5`(`md5` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 20 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '文件信息表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for sys_file_category
-- ----------------------------
DROP TABLE IF EXISTS `sys_file_category`;
CREATE TABLE `sys_file_category`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '分类ID',
  `parent_id` bigint(20) NULL DEFAULT NULL COMMENT '父分类ID',
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '分类名称',
  `code` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '分类编码',
  `icon` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '图标',
  `sort` int(11) NULL DEFAULT 0 COMMENT '排序',
  `status` tinyint(1) NULL DEFAULT 1 COMMENT '状态(0:禁用 1:正常)',
  `create_time` datetime NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_code`(`code` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 44 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '文件分类表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for sys_permission
-- ----------------------------
DROP TABLE IF EXISTS `sys_permission`;
CREATE TABLE `sys_permission`  (
  `permission_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '权限ID',
  `parent_id` bigint(20) NULL DEFAULT NULL COMMENT '上级权限ID',
  `permission_type` varchar(10) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '权限类型（菜单、按钮等）',
  `permission_name` varchar(30) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '权限名称',
  `icon` varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '图标',
  `route_path` varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '路由路径',
  `component_path` varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '组件路径',
  `permission_code` varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '权限标识',
  `sort` int(11) NULL DEFAULT NULL COMMENT '排序',
  `status` char(1) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT '0' COMMENT '状态（0显示 1隐藏）',
  `del_flag` char(1) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT '0' COMMENT '删除标志（0代表存在 2代表删除）',
  `create_time` datetime NULL DEFAULT NULL COMMENT '创建时间',
  `create_by` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '创建者',
  `update_time` datetime NULL DEFAULT NULL COMMENT '更新时间',
  `update_by` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '更新者',
  PRIMARY KEY (`permission_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 22 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '权限信息表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for sys_role
-- ----------------------------
DROP TABLE IF EXISTS `sys_role`;
CREATE TABLE `sys_role`  (
  `role_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '角色ID',
  `role_name` varchar(30) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '角色名称',
  `role_key` varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '角色标识',
  `status` char(1) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT '0' COMMENT '角色状态（0正常 1停用）',
  `del_flag` char(1) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT '0' COMMENT '删除标志（0代表存在 2代表删除）',
  `create_time` datetime NULL DEFAULT NULL COMMENT '创建时间',
  `create_by` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '创建者',
  `update_time` datetime NULL DEFAULT NULL COMMENT '更新时间',
  `update_by` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '更新者',
  PRIMARY KEY (`role_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 7 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '角色信息表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for sys_role_permission
-- ----------------------------
DROP TABLE IF EXISTS `sys_role_permission`;
CREATE TABLE `sys_role_permission`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '关联ID',
  `role_id` bigint(20) NOT NULL COMMENT '角色ID',
  `permission_id` bigint(20) NOT NULL COMMENT '权限ID',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uniq_role_permission`(`role_id` ASC, `permission_id` ASC) USING BTREE,
  INDEX `fk_role_permission_permission`(`permission_id` ASC) USING BTREE,
  CONSTRAINT `fk_role_permission_permission` FOREIGN KEY (`permission_id`) REFERENCES `sys_permission` (`permission_id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
  CONSTRAINT `fk_role_permission_role` FOREIGN KEY (`role_id`) REFERENCES `sys_role` (`role_id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE = InnoDB AUTO_INCREMENT = 266 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '角色权限关联表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for sys_user
-- ----------------------------
DROP TABLE IF EXISTS `sys_user`;
CREATE TABLE `sys_user`  (
  `user_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '用户ID',
  `user_name` varchar(30) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '用户账号',
  `password` varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '密码',
  `nick_name` varchar(30) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '用户昵称',
  `email` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '用户邮箱',
  `phonenumber` varchar(11) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '手机号码',
  `sex` char(1) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT '0' COMMENT '用户性别（0男 1女 2未知）',
  `avatar` varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '头像地址',
  `status` char(1) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT '0' COMMENT '帐号状态（0正常 1停用）',
  `del_flag` char(1) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT '0' COMMENT '删除标志（0代表存在 2代表删除）',
  `create_time` datetime NULL DEFAULT NULL COMMENT '创建时间',
  `create_by` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '创建者',
  `update_time` datetime NULL DEFAULT NULL COMMENT '更新时间',
  `update_by` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '更新者',
  PRIMARY KEY (`user_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 6 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '用户信息表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for sys_user_role
-- ----------------------------
DROP TABLE IF EXISTS `sys_user_role`;
CREATE TABLE `sys_user_role`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '关联ID',
  `user_id` bigint(20) NOT NULL COMMENT '用户ID',
  `role_id` bigint(20) NOT NULL COMMENT '角色ID',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uniq_user_role`(`user_id`, `role_id`) USING BTREE
) ENGINE = MyISAM AUTO_INCREMENT = 14 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Fixed;

INSERT INTO `sys_permission` 
(`permission_id`, `parent_id`, `permission_type`, `permission_name`, `icon`, `route_path`, `component_path`, `permission_code`, `sort`, `status`, `del_flag`) 
VALUES 
(NULL, (SELECT permission_id FROM (SELECT permission_id FROM sys_permission WHERE permission_code = 'system') AS p), 
'menu', '文件管理', 'Folder', 'file', 'system/file/index', 'system:file:list', 6, '0', '0');

SET FOREIGN_KEY_CHECKS = 1;
