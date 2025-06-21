import request from '@/utils/request'

// 使用Mock数据
import { 
  getTaskList as mockGetTaskList,
  getTaskDetail as mockGetTaskDetail,
  getTaskData as mockGetTaskData,
  createTask as mockCreateTask,
  deleteTask as mockDeleteTask 
} from '@/mock/crawler/task'

const USE_MOCK = true // 是否使用Mock数据

// 获取任务列表
export function getTaskList(params) {
  if (USE_MOCK) {
    return Promise.resolve(mockGetTaskList(params))
  }

  return request({
    url: '/crawler/tasks',
    method: 'get',
    params
  })
}

// 获取任务详情
export function getTaskDetail(id) {
  if (USE_MOCK) {
    return Promise.resolve(mockGetTaskDetail(id))
  }

  return request({
    url: `/crawler/tasks/${id}`,
    method: 'get'
  })
}

// 获取任务数据
export function getTaskData(taskId, params) {
  if (USE_MOCK) {
    return Promise.resolve(mockGetTaskData(taskId, params))
  }

  return request({
    url: `/crawler/tasks/${taskId}/data`,
    method: 'get',
    params
  })
}

// 创建任务
export function createTask(data) {
  if (USE_MOCK) {
    return Promise.resolve(mockCreateTask(data))
  }

  return request({
    url: '/crawler/tasks',
    method: 'post',
    data
  })
}

// 删除任务
export function deleteTask(id) {
  if (USE_MOCK) {
    return Promise.resolve(mockDeleteTask(id))
  }

  return request({
    url: `/crawler/tasks/${id}`,
    method: 'delete'
  })
}

// 取消任务
export function cancelTask(id) {
  if (USE_MOCK) {
    return Promise.resolve({
      code: 200,
      message: '取消成功'
    })
  }

  return request({
    url: `/crawler/tasks/${id}/cancel`,
    method: 'post'
  })
}

// 重试任务
export function retryTask(id) {
  if (USE_MOCK) {
    return Promise.resolve({
      code: 200,
      message: '重试成功'
    })
  }

  return request({
    url: `/crawler/tasks/${id}/retry`,
    method: 'post'
  })
}

// 下载任务结果（所有文件打包）
export function downloadTask(id) {
  if (USE_MOCK) {
    return Promise.resolve({
      code: 200,
      message: '开始下载'
    })
  }

  return request({
    url: `/crawler/tasks/${id}/download`,
    method: 'get',
    responseType: 'blob'
  })
}

// 下载单个文件
export function downloadTaskItem(taskId, itemId) {
  if (USE_MOCK) {
    return Promise.resolve({
      code: 200,
      message: '开始下载文件'
    })
  }

  return request({
    url: `/crawler/data/${itemId}/download`,
    method: 'get',
    responseType: 'blob'
  })
} 