import request from '@/admin/utils/request'

// 获取分享申请列表
export function listShareApply(query) {
  return request({
    url: '/share/apply/list',
    method: 'get',
    params: query
  })
}

// 审核通过
export function approveShareApply(id) {
  return request({
    url: `/share/apply/approve/${id}`,
    method: 'post'
  })
}

// 审核拒绝
export function rejectShareApply(data) {
  return request({
    url: '/share/apply/reject',
    method: 'post',
    data: data
  })
}

// 获取申请详情
export function getShareApplyDetail(id) {
  return request({
    url: `/share/apply/detail/${id}`,
    method: 'get'
  })
} 