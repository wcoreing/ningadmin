import request from '@/utils/request'

// 获取品牌列表
export function getBrandList() {
  return request({
    url: '/shop/brand/list',
    method: 'get'
  })
} 