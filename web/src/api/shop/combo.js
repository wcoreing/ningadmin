import request from '@/utils/request'

// 获取套餐列表
export function getComboList(params) {
  return request({
    url: '/shop/combo/list',
    method: 'get',
    params
  })
}

// 获取套餐详情
export function getComboInfo(id) {
  return request({
    url: `/shop/combo/info/${id}`,
    method: 'get'
  })
}

// 创建套餐
export function createCombo(data) {
  return request({
    url: '/shop/combo/create',
    method: 'post',
    data
  })
}

// 更新套餐
export function updateCombo(data) {
  return request({
    url: '/shop/combo/update',
    method: 'put',
    data
  })
}

// 删除套餐
export function deleteCombo(id) {
  return request({
    url: `/shop/combo/delete/${id}`,
    method: 'delete'
  })
}

// 批量删除套餐
export function batchDeleteCombo(ids) {
  return request({
    url: '/shop/combo/batch-delete',
    method: 'post',
    data: { ids }
  })
}

// 修改套餐状态
export function changeComboStatus(data) {
  return request({
    url: '/shop/combo/status',
    method: 'put',
    data
  })
}

// 复制套餐
export function copyCombo(id) {
  return request({
    url: `/shop/combo/copy/${id}`,
    method: 'post'
  })
}

// 获取套餐组列表
export function getComboGroups(comboId) {
  return request({
    url: `/shop/combo/${comboId}/groups`,
    method: 'get'
  })
}

// 保存套餐组配置
export function saveComboGroups(comboId, groups) {
  return request({
    url: `/shop/combo/${comboId}/groups`,
    method: 'post',
    data: { groups }
  })
} 