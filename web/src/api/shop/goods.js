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