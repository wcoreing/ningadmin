import request from '@/admin/utils/request'

// 获取反馈列表
export function getFeedbackList(params) {
  return request({
    url: '/shop/feedback/list',
    method: 'get',
    params
  })
}

// 获取反馈详情
export function getFeedbackInfo(id) {
  return request({
    url: `/shop/feedback/${id}`,
    method: 'get'
  })
}

// 新增反馈
export function createFeedback(data) {
  return request({
    url: '/shop/feedback',
    method: 'post',
    data
  })
}

// 更新反馈
export function updateFeedback(data) {
  return request({
    url: '/shop/feedback',
    method: 'put',
    data
  })
}

// 删除反馈
export function deleteFeedback(id) {
  return request({
    url: `/shop/feedback/${id}`,
    method: 'delete'
  })
}

// 回复反馈
export function replyFeedback(id, reply) {
  return request({
    url: `/shop/feedback/${id}/reply`,
    method: 'post',
    data: { reply }
  })
}

// 更新反馈状态
export function updateFeedbackStatus(id, status) {
  return request({
    url: `/shop/feedback/${id}/status`,
    method: 'put',
    data: { status }
  })
} 