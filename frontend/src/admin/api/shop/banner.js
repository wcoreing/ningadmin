import request from '@/admin/utils/request'

const api = {
  BannerList: '/shop/banner/list',
  BannerInfo: '/shop/banner/info',
  BannerCreate: '/shop/banner/create',
  BannerUpdate: '/shop/banner/update',
  BannerDelete: '/shop/banner/delete',
  BannerStatus: '/shop/banner/status'
}

// 获取Banner列表
export function getBannerList(params) {
  return request({
    url: api.BannerList,
    method: 'get',
    params
  })
}

// 获取Banner详情
export function getBannerInfo(id) {
  return request({
    url: `${api.BannerInfo}/${id}`,
    method: 'get'
  })
}

// 创建Banner
export function createBanner(data) {
  return request({
    url: api.BannerCreate,
    method: 'post',
    data
  })
}

// 更新Banner
export function updateBanner(data) {
  return request({
    url: api.BannerUpdate,
    method: 'put',
    data
  })
}

// 删除Banner
export function deleteBanner(id) {
  return request({
    url: `${api.BannerDelete}/${id}`,
    method: 'delete'
  })
}

// 修改Banner状态
export function changeBannerStatus(id, status) {
  return request({
    url: `${api.BannerStatus}/${id}`,
    method: 'put',
    data: { status }
  })
} 