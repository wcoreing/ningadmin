import request from '@/utils/request'

// 登录
export function login(data) {
  return request({
    url: '/auth/login',
    method: 'post',
    data
  })
}

// 注册
export function register(data) {
  return request({
    url: '/auth/register',
    method: 'post',
    data
  })
}

// 获取用户信息
export function getUserInfo() {
  return request({
    url: '/user/info',
    method: 'get'
  })
}

// 登出
export function logout() {
  return request({
    url: '/auth/logout',
    method: 'post'
  })
}

// 获取用户列表
export function getUserList(params) {
  return request({
    url: '/user/list',
    method: 'get',
    params
  })
}

// 获取用户详情
export function getUserDetail(id) {
  return request({
    url: `/user/${id}`,
    method: 'get'
  })
}

// 新增用户
export function createUser(data) {
  return request({
    url: '/user',
    method: 'post',
    data
  })
}

// 更新用户
export function updateUser(data) {
  return request({
    url: '/user',
    method: 'put',
    data
  })
}

// 删除用户
export function deleteUser(id) {
  return request({
    url: `/user/${id}`,
    method: 'delete'
  })
}

// 更新用户状态
export function updateUserStatus(data) {
  return request({
    url: '/user/status',
    method: 'put',
    data
  })
}

// 重置用户密码
export function resetUserPassword(data) {
  return request({
    url: '/user/password',
    method: 'put',
    data
  })
}

// 更新用户余额
export function updateUserBalance(data) {
  return request({
    url: '/user/balance',
    method: 'put',
    data
  })
}

// 更新用户积分
export function updateUserPoints(data) {
  return request({
    url: '/user/points',
    method: 'put',
    data
  })
}

// 获取用户变动记录列表
export function getUserChangeLogList(params) {
  return request({
    url: '/user/change-log',
    method: 'get',
    params
  })
}

// 获取指定用户的变动记录
export function getUserChangeLog(userId, params) {
  return request({
    url: `/user/${userId}/change-log`,
    method: 'get',
    params
  })
} 