import request from '@/utils/request'

// 获取分销收益列表
export function getIncomeList(query) {
  return request({
    url: '/distribution/income/list',
    method: 'get',
    params: query
  })
}

// 解冻分销收益
export function unfreezeIncome(data) {
  return request({
    url: '/distribution/income/unfreeze',
    method: 'post',
    data
  })
}

// 导出分销收益数据
export function exportIncome(query) {
  return request({
    url: '/distribution/income/export',
    method: 'get',
    params: query,
    responseType: 'blob'
  })
} 