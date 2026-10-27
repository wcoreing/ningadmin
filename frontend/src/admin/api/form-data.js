import request from '@/admin/utils/request'

// 获取表单数据列表
export function getFormDataList(params) {
  return request({
    url: '/form/data/list',
    method: 'get',
    params
  })
}

// 获取表单数据详情
export function getFormDataDetail(id) {
  return request({
    url: `/form/data/${id}`,
    method: 'get'
  })
}

// 删除表单数据
export function deleteFormData(id) {
  return request({
    url: `/form/data/${id}`,
    method: 'delete'
  })
}

// 提交表单数据
export function submitFormData(data) {
  return request({
    url: '/form/data/submit',
    method: 'post',
    data
  })
} 