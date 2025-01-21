import request from '@/utils/request'

// 获取表单列表
export function getFormList(params) {
  return request({
    url: '/form/form/list',
    method: 'get',
    params
  })
}

// 获取表单详情
export function getFormDetail(id) {
  return request({
    url: `/form/form/${id}`,
    method: 'get'
  })
}

// 新增表单
export function addForm(data) {
  return request({
    url: '/form/form',
    method: 'post',
    data
  })
}

// 修改表单
export function updateForm(data) {
  return request({
    url: '/form/form',
    method: 'put',
    data
  })
}

// 删除表单
export function deleteForm(id) {
  return request({
    url: `/form/form/${id}`,
    method: 'delete'
  })
}

// 获取表单模板列表
export function getFormTemplateList(params) {
  return request({
    url: '/form/template/list',
    method: 'get',
    params
  })
}

// 保存表单模板
export function saveFormTemplate(data) {
  return request({
    url: '/form/template',
    method: data.id ? 'put' : 'post',
    data
  })
}

// 删除表单模板
export function deleteFormTemplate(id) {
  return request({
    url: `/form/template/${id}`,
    method: 'delete'
  })
}


// 获取表单模板详情
export function getFormTemplateDetail(id) {
  return request({
    url: `/form/template/${id}`,
    method: 'get'
  })
} 