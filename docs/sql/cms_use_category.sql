-- 内容栏目并入通用分类 category（module=cms），然后去掉 cms_channel。

SET @has_code = (
  SELECT COUNT(*) FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'category' AND COLUMN_NAME = 'code'
);
SET @add_code = IF(
  @has_code = 0,
  'ALTER TABLE category ADD COLUMN code varchar(50) NOT NULL DEFAULT '''' COMMENT ''编码，内容栏目用于对外地址'' AFTER name',
  'SELECT 1'
);
PREPARE stmt FROM @add_code;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

DROP PROCEDURE IF EXISTS move_cms_channel_to_category;
DELIMITER $$
CREATE PROCEDURE move_cms_channel_to_category()
BEGIN
  IF EXISTS (
    SELECT 1 FROM information_schema.tables
    WHERE table_schema = DATABASE() AND table_name = 'cms_channel'
  ) THEN
    INSERT INTO category (module, parent_id, name, code, image, level, sort, status, is_delete)
    SELECT 'cms', 0, name, code, cover, 1, sort, status, 0
    FROM cms_channel
    WHERE is_delete = 0
      AND NOT EXISTS (
        SELECT 1 FROM category c
        WHERE c.module = 'cms' AND c.code = cms_channel.code AND c.is_delete = 0
      );

    UPDATE category child
    INNER JOIN cms_channel ch ON ch.code = child.code AND ch.is_delete = 0
    INNER JOIN cms_channel parent ON parent.id = ch.parent_id AND parent.is_delete = 0
    INNER JOIN category p ON p.module = 'cms' AND p.code = parent.code AND p.is_delete = 0
    SET child.parent_id = p.id, child.level = 2
    WHERE child.module = 'cms' AND ch.parent_id <> 0;

    UPDATE cms_content ct
    INNER JOIN cms_channel ch ON ch.id = ct.channel_id
    INNER JOIN category c ON c.module = 'cms' AND c.code = ch.code AND c.is_delete = 0
    SET ct.channel_id = c.id;

    DROP TABLE cms_channel;
  END IF;
END$$
DELIMITER ;

CALL move_cms_channel_to_category();
DROP PROCEDURE IF EXISTS move_cms_channel_to_category;

DELETE FROM sys_role_permission WHERE permission_id = 301;
DELETE FROM sys_permission WHERE permission_id = 301;
