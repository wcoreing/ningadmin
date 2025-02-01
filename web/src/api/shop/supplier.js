import request from '@/utils/request'

// 获取供应商列表
export function getSupplierList() {
  return request({
    url: '/shop/supplier/list',
    method: 'get'
  })
} 