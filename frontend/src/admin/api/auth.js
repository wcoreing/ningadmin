import request from '@/admin/utils/request'

// 用户登录
export function login(data) {
  return request({
    url: '/auth/login',
    method: 'post',
    data
  })
}

// 获取用户信息
export function getUserInfo() {
  return request({
    url: '/auth/info',
    method: 'get'
  })
}

// 退出登录
export function logout() {
  return request({
    url: '/auth/logout',
    method: 'post'
  })
} 