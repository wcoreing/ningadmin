INSERT INTO sys_permission (
  permission_id, parent_id, permission_type, permission_name, icon,
  route_path, component_path, permission_code, sort, status, del_flag, create_time, create_by
)
SELECT 303, 300, 'menu', '内容分类', 'Collection', '/cms/category', '/views/cms/category/index.vue', 'cms:category', 30, '0', '0', NOW(), 'admin'
FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM sys_permission WHERE permission_id = 303);

INSERT INTO sys_role_permission (role_id, permission_id)
SELECT 1, 303
FROM DUAL
WHERE NOT EXISTS (
  SELECT 1 FROM sys_role_permission WHERE role_id = 1 AND permission_id = 303
);
