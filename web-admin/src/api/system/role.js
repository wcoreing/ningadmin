import request from '@/utils/request'

// 获取权限树
export function getPermissionTree() {
  return request({
    url: '/system/permission/tree',
    method: 'get'
  })
}

// 获取角色列表
export function getRoleList(params) {
  return request({
    url: '/system/role/list',
    method: 'get',
    params
  })
}

// 获取角色详情
export function getRoleInfo(roleId) {
  return request({
    url: `/system/role/${roleId}`,
    method: 'get'
  })
}

// 新增角色
export function createRole(data) {
  return request({
    url: '/system/role',
    method: 'post',
    data
  })
}

// 修改角色
export function updateRole(data) {
  return request({
    url: '/system/role',
    method: 'put',
    data
  })
}

// 删除角色
export function deleteRole(roleId) {
  return request({
    url: `/system/role/${roleId}`,
    method: 'delete'
  })
}

// 修改角色状态
export function updateRoleStatus(data) {
  return request({
    url: '/system/role/status',
    method: 'put',
    data
  })
}

// 获取角色权限
export function getRolePermissions(roleId) {
  return request({
    url: `/system/role/permission/${roleId}`,
    method: 'get'
  })
}

// 更新角色权限
export function updateRolePermissions(data) {
  return request({
    url: '/system/role/permission',
    method: 'put',
    data
  })
}

// 获取用户角色
export function getUserRoles(userId) {
  return request({
    url: `/system/user/${userId}/roles`,
    method: 'get'
  })
} 