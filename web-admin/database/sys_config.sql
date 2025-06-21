-- 系统配置表
CREATE TABLE `sys_config` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '配置ID',
  `config_name` varchar(100) NOT NULL COMMENT '配置名称',
  `config_key` varchar(100) NOT NULL COMMENT '配置键名',
  `config_value` varchar(500) NOT NULL COMMENT '配置键值',
  `config_type` varchar(20) NOT NULL COMMENT '配置类型',
  `module_key` varchar(50) NOT NULL COMMENT '所属模块标识',
  `remark` varchar(500) DEFAULT NULL COMMENT '备注',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_module_key` (`module_key`,`config_key`) COMMENT '模块下配置键唯一'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统配置表';

-- 初始数据
INSERT INTO `sys_config` (`config_name`, `config_key`, `config_value`, `config_type`, `module_key`, `remark`) VALUES 
('系统名称', 'sys.name', 'WnAdmin', 'string', 'system', '系统显示名称'),
('系统版本', 'sys.version', '1.0.0', 'string', 'system', '系统版本号'),
('上传大小限制', 'sys.uploadMaxSize', '50', 'number', 'system', '文件上传大小限制，单位MB'),
('SMTP服务器', 'mail.smtp.host', 'smtp.example.com', 'string', 'mail', '邮件服务器地址'),
('SMTP端口', 'mail.smtp.port', '465', 'number', 'mail', '邮件服务器端口'),
('上传路径', 'file.uploadPath', '/data/upload', 'string', 'file', '文件上传保存路径'),
('密码最小长度', 'security.password.minLength', '6', 'number', 'security', '密码最小长度限制'); 