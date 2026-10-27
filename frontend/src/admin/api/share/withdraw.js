import request from '@/admin/utils/request'

// 获取提现列表
export function listWithdraw(query) {
  return request({
    url: '/share/withdraw/list',
    method: 'get',
    params: query
  })
}

// 获取提现详情
export function getWithdrawDetail(id) {
  return request({
    url: `/share/withdraw/detail/${id}`,
    method: 'get'
  })
}

// 通过提现
export function approveWithdraw(id) {
  return request({
    url: `/share/withdraw/approve/${id}`,
    method: 'post'
  })
}

// 拒绝提现
export function rejectWithdraw(data) {
  return request({
    url: '/share/withdraw/reject',
    method: 'post',
    data
  })
}

// 确认打款
export function confirmWithdrawPay(data) {
  return request({
    url: '/share/withdraw/pay',
    method: 'post',
    data
  })
} 