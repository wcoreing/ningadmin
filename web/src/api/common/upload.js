import request from '@/utils/request'

// 服务器域名配置
const SERVER_URL = 'http://localhost:8080'  // 根据实际环境配置

// 处理文件 URL，添加服务器域名
function processFileUrl(file) {
  if (!file) return file
  return {
    ...file,
    url: file.url.startsWith('http') ? file.url : `${SERVER_URL}${file.url}`
  }
}

// 获取文件列表
export function getFileList(params) {
  return request({
    url: '/files',
    method: 'get',
    params: {
      page: params.page || 1,
      page_size: params.page_size || 12,
      keyword: params.keyword || '',
      category_id: params.category_id || ''
    }
  }).then(res => {
    // 处理列表中的文件 URL
    return {
      list: (res.list || []).map(processFileUrl),
      total: res.total
    }
  })
}

// 上传文件
export function uploadFile(data) {
  return request({
    url: '/files/upload',
    method: 'post',
    data,
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  }).then(res => {
    // 处理上传返回的文件 URL
    return processFileUrl(res)
  })
}

// 获取文件分类
export function getFileCategories() {
  return request({
    url: '/files/categories',
    method: 'get'
  })
}

// 删除文件
export function deleteFile(id) {
  return request({
    url: `/files/${id}`,
    method: 'delete'
  })
}

// 移动文件到分类
export function moveFiles(data) {
  return request({
    url: '/files/move',
    method: 'post',
    data: {
      file_ids: data.fileIds,
      category_id: data.categoryId
    }
  })
} 