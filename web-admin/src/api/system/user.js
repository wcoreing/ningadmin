import request from '@/utils/request'

// 获取用户列表
export function getUserList(params) {
  return request({
    url: '/system/user/list',
    method: 'get',
    params
  })
}

// 获取用户详情
export function getUserInfo(userId) {
  return request({
    url: `/system/user/info/${userId}`,
    method: 'get'
  })
}

// 新增用户
export function createUser(data) {
  return request({
    url: '/system/user',
    method: 'post',
    data
  })
}

// 修改用户
export function updateUser(data) {
  return request({
    url: '/system/user',
    method: 'put',
    data
  })
}

// 删除用户
export function deleteUser(userId) {
  return request({
    url: `/system/user/delete/${userId}`,
    method: 'delete'
  })
}

// 修改用户状态
export function updateUserStatus(data) {
  return request({
    url: '/system/user/status',
    method: 'put',
    data
  })
}

// 重置用户密码
export function resetUserPassword(data) {
  return request({
    url: '/system/user/password',
    method: 'put',
    data
  })
}

// 获取部门选项
export function getDeptOptions() {
  return request({
    url: '/system/dept/options',
    method: 'get'
  })
} 