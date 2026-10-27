import request from '@/admin/utils/request'

export function getContentList(params) {
  return request({
    url: '/cms/content/list',
    method: 'get',
    params
  })
}

export function getContentInfo(id) {
  return request({
    url: `/cms/content/info/${id}`,
    method: 'get'
  })
}

export function saveContent(data) {
  return request({
    url: '/cms/content',
    method: data.id ? 'put' : 'post',
    data
  })
}

export function deleteContent(id) {
  return request({
    url: `/cms/content/${id}`,
    method: 'delete'
  })
}
