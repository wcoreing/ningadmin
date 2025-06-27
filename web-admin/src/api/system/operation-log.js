import request from '@/utils/request'

// 获取操作日志列表
export function getOperationLogList(params) {
  return request({
    url: '/system/operation-log/list',
    method: 'get',
    params
  })
}

// 删除操作日志
export function deleteOperationLog(id) {
  return request({
    url: `/system/operation-log/delete/${id}`,
    method: 'delete'
  })
}

// 批量删除操作日志
export function batchDeleteOperationLogs(ids) {
  return request({
    url: '/system/operation-log/batch-delete',
    method: 'delete',
    data: { ids }
  })
}

// 清空操作日志
export function clearOperationLogs(days) {
  return request({
    url: '/system/operation-log/clear',
    method: 'delete',
    data: { days }
  })
}

// 导出操作日志
export function exportOperationLogs(params) {
  return request({
    url: '/system/operation-log/export',
    method: 'get',
    params,
    responseType: 'blob'
  })
} 