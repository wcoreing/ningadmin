import request from '@/utils/request'

// 获取分享收益列表
export function getIncomeList(query) {
  return request({
    url: '/share/income/list',
    method: 'get',
    params: query
  })
}

// 解冻分享收益
export function unfreezeIncome(data) {
  return request({
    url: '/share/income/unfreeze',
    method: 'post',
    data
  })
}

// 导出分享收益数据
export function exportIncome(query) {
  return request({
    url: '/share/income/export',
    method: 'get',
    params: query,
    responseType: 'blob'
  })
} 