import request from '@/admin/utils/request'

// 使用Mock数据
import { getCrawlerTypes as mockGetCrawlerTypes } from '@/admin/mock/crawler/task'

const USE_MOCK = true // 是否使用Mock数据

// 获取爬虫类型列表
export function getCrawlerTypes() {
  if (USE_MOCK) {
    return Promise.resolve(mockGetCrawlerTypes())
  }

  return request({
    url: '/crawler/types',
    method: 'get'
  })
} 