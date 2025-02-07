import request from '@/utils/request'

// 获取分类列表
export function getCategoryList(params) {
  return request({
    url: '/shop/categories',
    method: 'get',
    params
  })
}

// 创建分类
export function createCategory(data) {
  return request({
    url: '/shop/categories',
    method: 'post',
    data
  })
}

// 更新分类
export function updateCategory(id, data) {
  return request({
    url: `/shop/categories/${id}`,
    method: 'put',
    data
  })
}

// 删除分类
export function deleteCategory(id) {
  return request({
    url: `/shop/categories/${id}`,
    method: 'delete'
  })
} 