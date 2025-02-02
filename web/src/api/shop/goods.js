import request from '@/utils/request'

const api = {
  GoodsList: '/shop/goods/list',
  GoodsInfo: '/shop/goods/info',
  GoodsUpdate: '/shop/goods/update',
  GoodsDelete: '/shop/goods/delete',
  GoodsStatus: '/shop/goods/status',
  GoodsCreate: '/shop/goods/create'
}

// 获取商品列表
export function getGoodsList(params) {
  return request({
    url: api.GoodsList,
    method: 'get',
    params
  })
}

// 删除商品
export function deleteGoods(id) {
  return request({
    url: `${api.GoodsDelete}/${id}`,
    method: 'delete'
  })
}

// 更新商品状态
export function updateGoodsStatus(data) {
  return request({
    url: api.GoodsStatus,
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

// 获取商品详情
export function getGoodsInfo(id) {
  return request({
    url: `${api.GoodsInfo}/${id}`,
    method: 'get'
  })
}

// 更新商品
export function updateGoods(data) {
  return request({
    url: api.GoodsUpdate,
    method: 'put',
    data
  })
}

// 创建商品
export function createGoods(data) {
  return request({
    url: api.GoodsCreate,
    method: 'post',
    data
  })
} 