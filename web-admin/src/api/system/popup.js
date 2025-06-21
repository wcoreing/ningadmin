import request from '@/utils/request'

// 获取弹窗列表
export function getPopupList(params) {
  return request({
    url: '/system/popup/list',
    method: 'get',
    params
  })
}

// 获取弹窗详情
export function getPopupInfo(id) {
  return request({
    url: `/system/popup/${id}`,
    method: 'get'
  })
}

// 创建弹窗
export function createPopup(data) {
  return request({
    url: '/system/popup',
    method: 'post',
    data
  })
}

// 更新弹窗
export function updatePopup(data) {
  return request({
    url: '/system/popup',
    method: 'put',
    data
  })
}

// 删除弹窗
export function deletePopup(id) {
  return request({
    url: `/system/popup/${id}`,
    method: 'delete'
  })
}

// 更新弹窗状态
export function updatePopupStatus(data) {
  return request({
    url: '/system/popup/status',
    method: 'put',
    data
  })
}

// 更新弹窗统计
export function updatePopupStats(data) {
  return request({
    url: '/system/popup/stats',
    method: 'put',
    data
  })
}

// 获取活跃弹窗
export function getActivePopups(params) {
  return request({
    url: '/system/popup/active',
    method: 'get',
    params
  })
} 