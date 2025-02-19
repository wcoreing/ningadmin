import request from '@/utils/request'

// 获取分销设置
export function getDistributionSetting() {
  return request({
    url: '/distribution/setting',
    method: 'get'
  })
}

// 更新分销设置
export function updateDistributionSetting(data) {
  return request({
    url: '/distribution/setting',
    method: 'put',
    data
  })
} 