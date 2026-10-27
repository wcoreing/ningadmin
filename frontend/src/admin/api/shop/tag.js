import request from '@/admin/utils/request'

// 获取标签列表
export function getTagList() {
  return request({
    url: '/shop/tag/list',
    method: 'get'
  })
} 