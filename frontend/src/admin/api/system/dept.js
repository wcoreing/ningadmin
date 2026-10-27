import request from '@/admin/utils/request'

// 获取部门列表
export function getDeptList(params) {
  return request({
    url: '/system/dept/list',
    method: 'get',
    params
  })
}

// 获取部门树
export function getDeptTree(params) {
  return request({
    url: '/system/dept/tree',
    method: 'get',
    params
  })
}

// 新增部门
export function addDept(data) {
  return request({
    url: '/system/dept',
    method: 'post',
    data
  })
}

// 修改部门
export function updateDept(data) {
  return request({
    url: '/system/dept',
    method: 'put',
    data
  })
}

// 删除部门
export function deleteDept(deptId) {
  return request({
    url: `/system/dept/${deptId}`,
    method: 'delete'
  })
}

// 获取部门详情
export function getDeptInfo(deptId) {
  return request({
    url: `/system/dept/${deptId}`,
    method: 'get'
  })
}

// 更新部门状态
export function updateDeptStatus(data) {
  return request({
    url: '/system/dept/status',
    method: 'put',
    data
  })
} 