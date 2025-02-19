import request from '@/utils/request'

// 获取分销员申请列表
export function listDistributionApply(query) {
  return request({
    url: '/distribution/apply/list',
    method: 'get',
    params: query
  })
}

// 审核通过
export function approveDistributionApply(id) {
  return request({
    url: `/distribution/apply/approve/${id}`,
    method: 'post'
  })
}

// 审核拒绝
export function rejectDistributionApply(data) {
  return request({
    url: '/distribution/apply/reject',
    method: 'post',
    data: data
  })
}

// 获取申请详情
export function getDistributionApplyDetail(id) {
  return request({
    url: `/distribution/apply/detail/${id}`,
    method: 'get'
  })
} 