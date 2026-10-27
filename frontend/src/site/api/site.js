import request from './request'

export function getHome() {
  return request({ url: '/site/home', method: 'get' })
}

export function getProfile() {
  return request({ url: '/site/profile', method: 'get' })
}

export function getPage(code) {
  return request({ url: `/site/pages/${code}`, method: 'get' })
}

export function getCategories() {
  return request({ url: '/site/categories', method: 'get' })
}

export function getGoods(params) {
  return request({ url: '/site/goods', method: 'get', params })
}

export function getGoodsDetail(id) {
  return request({ url: `/site/goods/${id}`, method: 'get' })
}

export function getArticles(params) {
  return request({ url: '/site/articles', method: 'get', params })
}

export function getArticle(id) {
  return request({ url: `/site/articles/${id}`, method: 'get' })
}

export function createInquiry(data) {
  return request({ url: '/site/inquiries', method: 'post', data })
}
