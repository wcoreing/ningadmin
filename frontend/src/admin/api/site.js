import request from '@/admin/utils/request'

export function getSiteProfile() {
  return request({ url: '/site-admin/profile', method: 'get' })
}

export function saveSiteProfile(data) {
  return request({ url: '/site-admin/profile', method: 'put', data })
}

export function getInquiries(params) {
  return request({ url: '/site-admin/inquiries', method: 'get', params })
}

export function updateInquiryStatus(id, status) {
  return request({ url: `/site-admin/inquiries/${id}/status`, method: 'put', data: { status } })
}
