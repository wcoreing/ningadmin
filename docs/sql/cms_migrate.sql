-- 把官网栏目和资讯收进通用 CMS，然后删掉旧表。可重复执行。

CREATE TABLE IF NOT EXISTS `cms_channel` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '栏目ID',
  `parent_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '父栏目',
  `code` varchar(50) NOT NULL COMMENT '栏目编码',
  `name` varchar(100) NOT NULL COMMENT '栏目名称',
  `type` varchar(20) NOT NULL DEFAULT 'list' COMMENT 'page单页 list列表',
  `summary` varchar(500) NOT NULL DEFAULT '' COMMENT '摘要',
  `cover` varchar(255) NOT NULL DEFAULT '' COMMENT '封面',
  `seo_title` varchar(120) NOT NULL DEFAULT '' COMMENT 'SEO标题',
  `seo_keywords` varchar(255) NOT NULL DEFAULT '' COMMENT 'SEO关键词',
  `seo_description` varchar(300) NOT NULL DEFAULT '' COMMENT 'SEO描述',
  `sort` int(11) NOT NULL DEFAULT 0 COMMENT '排序',
  `status` tinyint(4) NOT NULL DEFAULT 1 COMMENT '状态(0隐藏 1显示)',
  `is_delete` tinyint(4) NOT NULL DEFAULT 0,
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_parent_id` (`parent_id`),
  KEY `idx_code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='CMS栏目';

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

DROP PROCEDURE IF EXISTS migrate_site_to_cms;
DELIMITER $$
CREATE PROCEDURE migrate_site_to_cms()
BEGIN
  IF EXISTS (
    SELECT 1 FROM information_schema.tables
    WHERE table_schema = DATABASE() AND table_name = 'site_page'
  ) THEN
    INSERT INTO cms_channel (
      parent_id, code, name, type, summary, cover,
      seo_title, seo_keywords, seo_description, sort, status, is_delete, create_time, update_time
    )
    SELECT
      0, p.code, p.title, 'page', p.summary, p.cover,
      p.seo_title, p.seo_keywords, p.seo_description, p.sort, p.status, 0, p.create_time, p.update_time
    FROM site_page p
    WHERE p.is_delete = 0
      AND NOT EXISTS (
        SELECT 1 FROM cms_channel c WHERE c.code = p.code AND c.is_delete = 0
      );

    INSERT INTO cms_content (
      channel_id, title, summary, cover, body,
      seo_title, seo_keywords, seo_description, sort, status, is_delete,
      publish_time, create_time, update_time
    )
    SELECT
      c.id, p.title, p.summary, p.cover, p.content,
      p.seo_title, p.seo_keywords, p.seo_description, p.sort, p.status, 0,
      IFNULL(p.update_time, NOW()), p.create_time, p.update_time
    FROM site_page p
    INNER JOIN cms_channel c ON c.code = p.code AND c.is_delete = 0 AND c.type = 'page'
    WHERE p.is_delete = 0
      AND NOT EXISTS (
        SELECT 1 FROM cms_content x WHERE x.channel_id = c.id AND x.is_delete = 0
      );

    DROP TABLE site_page;
  END IF;

  IF NOT EXISTS (SELECT 1 FROM cms_channel WHERE code = 'news' AND is_delete = 0) THEN
    INSERT INTO cms_channel (parent_id, code, name, type, summary, sort, status)
    VALUES (0, 'news', '资讯', 'list', '公司动态和行业资料', 50, 1);
  END IF;

  INSERT INTO cms_channel (parent_id, code, name, type, sort, status)
  SELECT n.id, 'company', '公司动态', 'list', 20, 1
  FROM cms_channel n
  WHERE n.code = 'news' AND n.is_delete = 0
    AND NOT EXISTS (SELECT 1 FROM cms_channel WHERE code = 'company' AND is_delete = 0)
  LIMIT 1;

  INSERT INTO cms_channel (parent_id, code, name, type, sort, status)
  SELECT n.id, 'industry', '行业资料', 'list', 10, 1
  FROM cms_channel n
  WHERE n.code = 'news' AND n.is_delete = 0
    AND NOT EXISTS (SELECT 1 FROM cms_channel WHERE code = 'industry' AND is_delete = 0)
  LIMIT 1;

  IF EXISTS (
    SELECT 1 FROM information_schema.tables
    WHERE table_schema = DATABASE() AND table_name = 'site_article'
  ) THEN
    INSERT INTO cms_content (
      channel_id, title, summary, cover, body,
      seo_title, seo_keywords, seo_description, sort, status, is_delete,
      publish_time, create_time, update_time
    )
    SELECT
      ch.id, a.title, a.summary, a.cover, a.content,
      a.seo_title, a.seo_keywords, a.seo_description, a.sort, a.status, 0,
      a.publish_time, a.create_time, a.update_time
    FROM site_article a
    INNER JOIN cms_channel ch ON ch.code = a.category AND ch.is_delete = 0
    WHERE a.is_delete = 0
      AND NOT EXISTS (
        SELECT 1 FROM cms_content x
        WHERE x.channel_id = ch.id AND x.title = a.title AND x.is_delete = 0
      );

    DROP TABLE site_article;
  END IF;
END$$
DELIMITER ;

CALL migrate_site_to_cms();
DROP PROCEDURE IF EXISTS migrate_site_to_cms;

INSERT INTO sys_permission (
  permission_id, parent_id, permission_type, permission_name, icon,
  route_path, component_path, permission_code, sort, status, del_flag, create_time, create_by
)
SELECT 300, 0, 'menu', '内容', 'Notebook', '/cms', '', 'cms', 95, '0', '0', NOW(), 'admin'
FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM sys_permission WHERE permission_id = 300);

INSERT INTO sys_permission (
  permission_id, parent_id, permission_type, permission_name, icon,
  route_path, component_path, permission_code, sort, status, del_flag, create_time, create_by
)
SELECT 301, 300, 'menu', '栏目', 'Collection', '/cms/channel', '/views/cms/channel/index.vue', 'cms:channel', 20, '0', '0', NOW(), 'admin'
FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM sys_permission WHERE permission_id = 301);

INSERT INTO sys_permission (
  permission_id, parent_id, permission_type, permission_name, icon,
  route_path, component_path, permission_code, sort, status, del_flag, create_time, create_by
)
SELECT 302, 300, 'menu', '内容', 'Document', '/cms/content', '/views/cms/content/index.vue', 'cms:content', 10, '0', '0', NOW(), 'admin'
FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM sys_permission WHERE permission_id = 302);

INSERT INTO sys_role_permission (role_id, permission_id)
SELECT 1, p.permission_id
FROM sys_permission p
WHERE p.permission_id IN (300, 301, 302)
  AND NOT EXISTS (
    SELECT 1 FROM sys_role_permission rp
    WHERE rp.role_id = 1 AND rp.permission_id = p.permission_id
  );
