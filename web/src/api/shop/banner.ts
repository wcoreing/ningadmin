import request from '@/utils/request'

// 获取Banner列表
export function getBannerList(params) {
  return request({
    url: '/shop/banner/list',
    method: 'get',
    params
  })
}

// 获取指定位置的Banner列表
export function getBannersByPosition(position) {
  return request({
    url: '/shop/banner/position',
    method: 'get',
    params: { position }
  })
}

// 创建Banner
export function createBanner(data) {
  return request({
    url: '/shop/banner/create',
    method: 'post',
    data
  })
}

// 更新Banner
export function updateBanner(data) {
  return request({
    url: '/shop/banner/update',
    method: 'put',
    data
  })
}

// 删除Banner
export function deleteBanner(id) {
  return request({
    url: `/shop/banner/delete/${id}`,
    method: 'delete'
  })
}

// 修改Banner状态
export function changeBannerStatus(id, status) {
  return request({
    url: `/shop/banner/status/${id}`,
    method: 'put',
    params: { status }
  })
} 