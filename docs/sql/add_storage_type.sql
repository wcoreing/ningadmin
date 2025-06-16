-- 添加文件存储类型字段
-- 执行时间：请根据实际情况替换
-- 作者：System

-- 为sys_file表添加storage_type字段
ALTER TABLE `sys_file` ADD COLUMN `storage_type` varchar(20) NOT NULL DEFAULT 'local' COMMENT '存储类型(local:本地 oss:阿里云OSS)' AFTER `type`;

-- 更新历史数据的storage_type字段
-- 根据URL格式判断存储类型
UPDATE `sys_file` SET `storage_type` = 'oss' WHERE `url` LIKE 'http%';
UPDATE `sys_file` SET `storage_type` = 'local' WHERE `url` NOT LIKE 'http%';

-- 添加索引以提高查询性能
ALTER TABLE `sys_file` ADD INDEX `idx_storage_type` (`storage_type`);

-- 验证数据更新结果
SELECT 
    storage_type,
    COUNT(*) as count,
    GROUP_CONCAT(DISTINCT LEFT(url, 50) SEPARATOR '; ') as sample_urls
FROM `sys_file` 
WHERE `status` = 1 
GROUP BY storage_type; 