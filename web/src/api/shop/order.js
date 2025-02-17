import request from '@/utils/request'

const api = {
  OrderList: '/shop/order/list',
  OrderInfo: '/shop/order/info',
  OrderShip: '/shop/order/ship',
  OrderCancel: '/shop/order/cancel',
  OrderLogistics: '/shop/order/logistics'
}

// 获取订单列表
export function getOrderList(params) {
  return request({
    url: api.OrderList,
    method: 'get',
    params
  })
}

// 获取订单详情
export function getOrderInfo(id) {
  return request({
    url: `${api.OrderInfo}/${id}`,
    method: 'get'
  })
}

// 发货
export function shipOrder(data) {
  return request({
    url: api.OrderShip,
    method: 'put',
    data
  })
}

// 取消订单
export function cancelOrder(id) {
  return request({
    url: `${api.OrderCancel}/${id}`,
    method: 'put'
  })
}

// 获取物流轨迹
export function getLogisticsTrack(data) {
  return request({
    url: api.OrderLogistics,
    method: 'post',
    data
  })
} 