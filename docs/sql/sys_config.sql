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

 Date: 11/06/2025 13:44:26
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for sys_config
-- ----------------------------
DROP TABLE IF EXISTS `sys_config`;
CREATE TABLE `sys_config`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '配置ID',
  `parent_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '父级ID',
  `config_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '配置名称',
  `config_key` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '配置键名',
  `config_value` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '配置键值',
  `config_type` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '配置类型',
  `module_key` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '所属模块标识',
  `module_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '模块名称',
  `level` tinyint(4) NOT NULL DEFAULT 1 COMMENT '层级(1:一级 2:二级 3:三级)',
  `sort` int(11) NOT NULL DEFAULT 0 COMMENT '排序',
  `status` tinyint(3) UNSIGNED NOT NULL DEFAULT 1 COMMENT '状态(0:禁用 1:启用)',
  `remark` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '备注',
  `del_flag` char(1) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT '0' COMMENT '删除标志（0代表存在 2代表删除）',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_parent_id`(`parent_id` ASC) USING BTREE,
  INDEX `idx_module_key`(`module_key` ASC) USING BTREE,
  INDEX `idx_status`(`status` ASC) USING BTREE,
  INDEX `uk_module_key`(`module_key` ASC, `config_key` ASC) USING BTREE COMMENT '模块下配置键唯一'
) ENGINE = InnoDB AUTO_INCREMENT = 39 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '系统配置表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_config
-- ----------------------------
INSERT INTO `sys_config` VALUES (16, 0, '运费', '', '', 'string', 'shopConfig', '商城配置', 1, 0, 1, '', '0', '2025-06-08 10:25:42', '2025-06-08 10:25:42');
INSERT INTO `sys_config` VALUES (17, 16, '顺丰', 'shop.yunfei.shunfeng', '5', 'number', 'shopConfig', '商城配置', 2, 0, 1, '', '0', '2025-06-08 10:28:20', '2025-06-08 15:27:18');
INSERT INTO `sys_config` VALUES (19, 0, '客服', '', '', 'string', 'shopConfig', '商城配置', 1, 0, 1, '', '0', '2025-06-09 10:51:27', '2025-06-09 10:51:27');
INSERT INTO `sys_config` VALUES (20, 19, '客服电话', 'shop.service.phone', '15001764308', 'string', 'shopConfig', '商城配置', 2, 0, 1, '', '0', '2025-06-09 18:58:13', '2025-06-09 18:58:13');
INSERT INTO `sys_config` VALUES (21, 19, '客服二维码', 'shop.service.qrcode', 'http://localhost:8080/uploads/a5c1c332-4ccf-4b21-8ee1-c4a789e1ca31.png', 'image', 'shopConfig', '商城配置', 2, 0, 1, '', '0', '2025-06-09 18:58:41', '2025-06-09 19:05:35');
INSERT INTO `sys_config` VALUES (22, 0, '分享功能', '', '', 'group', 'shareConfig', '分享配置', 1, 100, 1, '分享功能相关配置', '0', '2025-06-11 13:39:36', NULL);
INSERT INTO `sys_config` VALUES (23, 22, '功能启用', 'share.enabled', '1', 'boolean', 'shareConfig', '分享配置', 2, 1, 1, '分享功能总开关', '0', '2025-06-11 13:39:36', NULL);
INSERT INTO `sys_config` VALUES (24, 22, '注册佣金', '', '', 'group', 'shareConfig', '分享配置', 2, 2, 1, '注册佣金相关配置', '0', '2025-06-11 13:39:36', NULL);
INSERT INTO `sys_config` VALUES (25, 22, '购买佣金', '', '', 'group', 'shareConfig', '分享配置', 2, 3, 1, '购买佣金相关配置', '0', '2025-06-11 13:39:36', NULL);
INSERT INTO `sys_config` VALUES (26, 22, '佣金管理', '', '', 'group', 'shareConfig', '分享配置', 2, 4, 1, '佣金管理相关配置', '0', '2025-06-11 13:39:36', NULL);
INSERT INTO `sys_config` VALUES (27, 22, '提现设置', '', '', 'group', 'shareConfig', '分享配置', 2, 5, 1, '提现相关配置', '0', '2025-06-11 13:39:36', NULL);
INSERT INTO `sys_config` VALUES (28, 24, '注册佣金开关', 'share.register.enabled', '1', 'boolean', 'shareConfig', '分享配置', 3, 1, 1, '是否启用注册佣金', '0', '2025-06-11 13:39:36', NULL);
INSERT INTO `sys_config` VALUES (29, 24, '注册佣金金额', 'share.register.amount', '5.00', 'number', 'shareConfig', '分享配置', 3, 2, 1, '新用户注册奖励金额(元)', '0', '2025-06-11 13:39:36', NULL);
INSERT INTO `sys_config` VALUES (30, 25, '购买佣金开关', 'share.purchase.enabled', '1', 'boolean', 'shareConfig', '分享配置', 3, 1, 1, '是否启用购买佣金', '0', '2025-06-11 13:39:36', NULL);
INSERT INTO `sys_config` VALUES (31, 25, '一级返佣比例', 'share.purchase.level1_rate', '10.00', 'number', 'shareConfig', '分享配置', 3, 2, 1, '一级分销员返佣比例(%)', '0', '2025-06-11 13:39:36', NULL);
INSERT INTO `sys_config` VALUES (32, 25, '二级返佣比例', 'share.purchase.level2_rate', '5.00', 'number', 'shareConfig', '分享配置', 3, 3, 1, '二级分销员返佣比例(%)', '0', '2025-06-11 13:39:36', NULL);
INSERT INTO `sys_config` VALUES (33, 26, '佣金冻结天数', 'share.commission.freeze_days', '7', 'number', 'shareConfig', '分享配置', 3, 1, 1, '佣金冻结天数，0表示不冻结', '0', '2025-06-11 13:39:36', NULL);
INSERT INTO `sys_config` VALUES (34, 27, '最低提现金额', 'share.withdraw.min_amount', '10.00', 'number', 'shareConfig', '分享配置', 3, 1, 1, '最低提现金额(元)', '0', '2025-06-11 13:39:36', NULL);
INSERT INTO `sys_config` VALUES (35, 27, '最高提现金额', 'share.withdraw.max_amount', '5000.00', 'number', 'shareConfig', '分享配置', 3, 2, 1, '最高提现金额(元)', '0', '2025-06-11 13:39:36', NULL);
INSERT INTO `sys_config` VALUES (36, 27, '提现手续费率', 'share.withdraw.fee_rate', '0.00', 'number', 'shareConfig', '分享配置', 3, 3, 1, '提现手续费率(%)', '0', '2025-06-11 13:39:36', NULL);
INSERT INTO `sys_config` VALUES (37, 27, '微信提现开关', 'share.withdraw.wechat_enabled', '1', 'boolean', 'shareConfig', '分享配置', 3, 4, 1, '是否启用微信提现', '0', '2025-06-11 13:39:36', NULL);
INSERT INTO `sys_config` VALUES (38, 27, '支付宝提现开关', 'share.withdraw.alipay_enabled', '0', 'boolean', 'shareConfig', '分享配置', 3, 5, 1, '是否启用支付宝提现', '0', '2025-06-11 13:39:36', NULL);

SET FOREIGN_KEY_CHECKS = 1;
