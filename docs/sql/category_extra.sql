-- 内容分类的编码、类型、摘要和 SEO 放进 category.extra，然后去掉 code 列。

SET @has_extra = (
  SELECT COUNT(*) FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'category' AND COLUMN_NAME = 'extra'
);
SET @add_extra = IF(
  @has_extra = 0,
  'ALTER TABLE category ADD COLUMN extra json NULL COMMENT ''模块扩展属性，内容分类存放编码、类型、摘要和SEO'' AFTER module',
  'SELECT 1'
);
PREPARE stmt FROM @add_extra;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @has_code = (
  SELECT COUNT(*) FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'category' AND COLUMN_NAME = 'code'
);
SET @fill_extra = IF(
  @has_code = 1,
  'UPDATE category SET extra = JSON_OBJECT(
      ''code'', code,
      ''type'', IF(code IN (''news'',''company'',''industry''), ''list'', ''page''),
      ''summary'', '''',
      ''seoTitle'', '''',
      ''seoKeywords'', '''',
      ''seoDescription'', ''''
    ) WHERE module = ''cms''',
  'SELECT 1'
);
PREPARE stmt FROM @fill_extra;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

UPDATE category c
INNER JOIN cms_content t ON t.channel_id = c.id AND t.is_delete = 0 AND t.title = c.name
SET c.extra = JSON_SET(
  c.extra,
  '$.summary', t.summary,
  '$.seoTitle', IFNULL(t.seo_title, ''),
  '$.seoKeywords', IFNULL(t.seo_keywords, ''),
  '$.seoDescription', IFNULL(t.seo_description, '')
)
WHERE c.module = 'cms'
  AND JSON_UNQUOTE(JSON_EXTRACT(c.extra, '$.type')) = 'page';

SET @has_code = (
  SELECT COUNT(*) FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'category' AND COLUMN_NAME = 'code'
);
SET @drop_code = IF(@has_code = 1, 'ALTER TABLE category DROP COLUMN code', 'SELECT 1');
PREPARE stmt FROM @drop_code;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
