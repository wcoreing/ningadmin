-- 菜单权限表
CREATE TABLE `sys_menu` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '菜单ID',
  `parent_id` bigint DEFAULT '0' COMMENT '父菜单ID',
  `name` varchar(50) NOT NULL COMMENT '菜单名称',
  `icon` varchar(100) DEFAULT NULL COMMENT '菜单图标',
  `type` varchar(10) NOT NULL COMMENT '菜单类型（menu菜单 button按钮）',
  `route_path` varchar(200) DEFAULT NULL COMMENT '路由地址',
  `component_path` varchar(255) DEFAULT NULL COMMENT '组件路径',
  `permission` varchar(100) DEFAULT NULL COMMENT '权限标识',
  `sort` int DEFAULT '0' COMMENT '显示顺序',
  `status` char(1) DEFAULT '0' COMMENT '菜单状态（0正常 1停用）',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='菜单权限表';

-- 初始数据
INSERT INTO `sys_menu` 
(`parent_id`, `name`, `icon`, `type`, `route_path`, `component_path`, `permission`, `sort`) 
VALUES 
-- 一级菜单
(0, '工作台', 'Monitor', 'menu', '/dashboard', '@/views/dashboard/index.vue', 'system', 1),
(0, '表单管理', 'Tickets', 'menu', '/form', NULL, 'form', 2),
(0, '系统监控', 'Monitor', 'menu', '/monitor', NULL, 'monitor', 3),
(0, '权限系统', 'Setting', 'menu', '/system', NULL, 'system', 4),

-- 表单管理子菜单
(2, '表单设计', NULL, 'menu', '/form-generator', '@/views/form-generator/index.vue', 'form:design', 1),
(2, '表单数据', NULL, 'menu', '/form-data', '@/views/form-data/index.vue', 'form:data', 2),

-- 系统监控子菜单  
(3, '服务器状态', NULL, 'menu', '/monitor', '@/views/monitor/index.vue', 'monitor:server', 1),

-- 权限系统子菜单
(4, '用户管理', NULL, 'menu', '/system/user', '@/views/system/user/index.vue', 'system:user', 1),
(4, '角色管理', NULL, 'menu', '/system/role', '@/views/system/role/index.vue', 'system:role', 2),
(4, '菜单管理', NULL, 'menu', '/system/menu', '@/views/system/menu/index.vue', 'system:menu', 3),
(4, '配置管理', NULL, 'menu', '/system/config', '@/views/system/config/index.vue', 'system:config', 4); 