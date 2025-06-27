import request from '@/utils/request'

// 获取菜单列表
export function getMenuList(params) {
  return request({
    url: '/system/menu/list',
    method: 'get',
    params
  })
}

// 新增菜单
export function addMenu(data) {
  return request({
    url: '/system/menu/add',
    method: 'post',
    data
  })
}

// 修改菜单
export function updateMenu(data) {
  return request({
    url: '/system/menu/update',
    method: 'put',
    data
  })
}

// 删除菜单
export function deleteMenu(id) {
  return request({
    url: `/system/menu/delete/${id}`,
    method: 'delete'
  })
}

// 获取菜单详情
export function getMenuInfo(id) {
  return request({
    url: `/system/menu/info/${id}`,
    method: 'get'
  })
}

// 获取菜单树形数据
export function getMenuTree(params) {
  return request({
    url: '/system/menu/tree',
    method: 'get',
    params
  })
}

// 检查权限标识是否存在
export function checkPermission(params) {
  return request({
    url: '/system/menu/check-permission',
    method: 'get',
    params
  })
}

// 更新菜单排序
export function updateMenuSort(data) {
  return request({
    url: '/system/menu/sort',
    method: 'put',
    data
  })
}

// 更新菜单状态
export function updateMenuStatus(data) {
  return request({
    url: '/system/menu/status',
    method: 'put',
    data
  })
} 