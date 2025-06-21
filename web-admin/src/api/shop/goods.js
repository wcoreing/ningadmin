import request from '@/utils/request'

const api = {
  GoodsList: '/shop/goods/list',
  GoodsInfo: '/shop/goods/info',
  GoodsUpdate: '/shop/goods/update',
  GoodsDelete: '/shop/goods/delete',
  GoodsStatus: '/shop/goods/status',
  GoodsCreate: '/shop/goods/create',
  GoodsCopy: '/shop/goods/copy',
  GoodsMove: '/shop/goods/move',
  GoodsBatchDelete: '/shop/goods/batch-delete',
  // 套餐相关接口
  ComboConfig: '/shop/goods/combo/config',
  ComboCalculate: '/shop/goods/combo/calculate',
  ComboValidate: '/shop/goods/combo/validate',
  ComboDefault: '/shop/goods/combo/default'
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
    url: '/common/category/list',
    method: 'get',
    params: {
      module: 'goods'
    }
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

// 复制商品
export function copyGoods(id, count) {
  return request({
    url: `${api.GoodsCopy}/${id}`,
    params: { count },
    method: 'post'
  })
}

// 批量移动商品
export function batchMoveGoods(data) {
  return request({
    url: api.GoodsMove,
    method: 'put',
    data
  })
}

// 批量删除商品
export function batchDeleteGoods(goodsIds) {
  return request({
    url: api.GoodsBatchDelete,
    method: 'delete',
    data: { goodsIds }
  })
} 

// ========== 套餐相关接口 ==========

// 获取套餐配置
export function getComboConfig(comboId) {
  return request({
    url: `${api.ComboConfig}/${comboId}`,
    method: 'get'
  })
}

// 计算套餐价格
export function calculateComboPrice(data) {
  return request({
    url: api.ComboCalculate,
    method: 'post',
    data
  })
}

// 验证套餐选择
export function validateComboSelection(data) {
  return request({
    url: api.ComboValidate,
    method: 'post',
    data
  })
}

// 获取默认套餐选择
export function getDefaultComboSelection(comboId) {
  return request({
    url: `${api.ComboDefault}/${comboId}`,
    method: 'get'
  })
}

// 兼容旧接口名称
export const listGoods = getGoodsList
export const getGoods = getGoodsInfo
export const delGoods = deleteGoods 