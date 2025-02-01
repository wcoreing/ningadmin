import request from '@/utils/request'

// 获取商品列表
export function getGoodsList(params) {
  return request({
    url: '/shop/goods/list',
    method: 'get',
    params
  })
}

// 删除商品
export function deleteGoods(id) {
  return request({
    url: `/shop/goods/delete/${id}`,
    method: 'delete'
  })
}

// 更新商品状态
export function updateGoodsStatus(data) {
  return request({
    url: '/shop/goods/status',
    method: 'put',
    data
  })
}

// 获取商品分类列表
export function getCategoryList() {
  return request({
    url: '/shop/category/list',
    method: 'get'
  })
}

// 获取供应商列表
export function getSupplierList() {
  return request({
    url: '/shop/supplier/list',
    method: 'get'
  })
}

// 获取品牌列表
export function getBrandList() {
  return request({
    url: '/shop/brand/list',
    method: 'get'
  })
}

// 获取标签列表
export function getTagList() {
  return request({
    url: '/shop/tag/list',
    method: 'get'
  })
} 