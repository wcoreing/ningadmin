import request from '@/admin/utils/request'

// 获取配置列表
export function getConfigList(params) {
  return request({
    url: '/system/config/list',
    method: 'get',
    params
  })
}

// 新增配置
export function addConfig(data) {
  return request({
    url: '/system/config/add',
    method: 'post',
    data
  })
}

// 修改配置
export function updateConfig(data) {
  return request({
    url: '/system/config/update',
    method: 'put',
    data
  })
}

// 删除配置
export function deleteConfig(id) {
  return request({
    url: `/system/config/delete/${id}`,
    method: 'delete'
  })
}

// 获取配置详情
export function getConfigInfo(id) {
  return request({
    url: `/system/config/info/${id}`,
    method: 'get'
  })
}

// 检查配置键是否存在
export function checkConfigKey(params) {
  return request({
    url: '/system/config/check-key',
    method: 'get',
    params
  })
} 