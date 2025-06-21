import request from '@/utils/request'

// 获取分类列表（通用接口，支持不同模块）
export function getCategoryList(params) {
  // 如果没有指定模块，默认为商品分类
  const module = params.module || 'goods'
  return request({
    url: '/common/category/list',
    method: 'get',
    params: {
      module,
      ...params
    }
  })
}

// 获取分类树（通用接口，支持不同模块）
export function getCategoryTree(params) {
  // 如果没有指定模块，默认为商品分类
  const module = params.module || 'goods'
  return request({
    url: '/common/category/tree',
    method: 'get',
    params: {
      module,
      ...params
    }
  })
}

// 获取商品分类列表（兼容旧接口）
export function getGoodsCategoryList(params) {
  return getCategoryList({ module: 'goods', ...params })
}

// 获取套餐分类列表
export function getComboCategoryList(params) {
  return getCategoryList({ module: 'combo', ...params })
}

// 获取表单分类列表
export function getFormCategoryList(params) {
  return getCategoryList({ module: 'form', ...params })
}

// 创建分类（通用接口）
export function createCategory(data) {
  // 根据模块类型选择不同的API端点
  if (data.module === 'goods') {
    return request({
      url: '/shop/categories',
      method: 'post',
      data
    })
  } else {
    // 其他模块使用通用接口
    return request({
      url: '/common/category',
      method: 'post',
      data
    })
  }
}

// 更新分类（通用接口）
export function updateCategory(id, data) {
  // 根据模块类型选择不同的API端点
  if (data.module === 'goods') {
    return request({
      url: `/shop/categories/${id}`,
      method: 'put',
      data
    })
  } else {
    // 其他模块使用通用接口
    return request({
      url: `/common/category/${id}`,
      method: 'put',
      data
    })
  }
}

// 删除分类（通用接口）
export function deleteCategory(id, module = 'goods') {
  // 根据模块类型选择不同的API端点
  if (module === 'goods') {
    return request({
      url: `/shop/categories/${id}`,
      method: 'delete'
    })
  } else {
    // 其他模块使用通用接口
    return request({
      url: `/common/category/${id}`,
      method: 'delete'
    })
  }
} 