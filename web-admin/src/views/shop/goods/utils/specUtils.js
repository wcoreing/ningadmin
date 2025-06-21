import { SPEC_TYPES } from '../constants/specTypes.js'

/**
 * 初始化表单数据
 * @returns {Object} 初始化的表单数据
 */
export const initFormData = () => ({
  goodsId: undefined,
  goodsName: '',
  categoryId: undefined,
  brand: '',
  supplier: '',
  type: 1,
  basePrice: 0,
  price: 0,
  originalPrice: 0,
  costPrice: 0,
  stock: 0,
  weight: 0,
  volume: 0,
  sort: 0,
  status: 0,
  description: '',
  coverImg: '',
  tagNames: [],
  images: [],
  specType: SPEC_TYPES.SINGLE,
  specGroups: [],
  skuList: [],
  comboGroups: [],
  detailDesc: ''
})

/**
 * 提取单规格数据
 * @param {Object} formData 表单数据
 * @returns {Object} 单规格数据
 */
export const extractSingleSpecData = (formData) => ({
  price: formData.price || 0,
  originalPrice: formData.originalPrice || 0,
  costPrice: formData.costPrice || 0,
  stock: formData.stock || 0,
  weight: formData.weight || 0,
  volume: formData.volume || 0
})

/**
 * 提取多规格数据
 * @param {Object} formData 表单数据
 * @returns {Object} 多规格数据
 */
export const extractMultiSpecData = (formData) => ({
  specGroups: Array.isArray(formData.specGroups) ? formData.specGroups : [],
  skuList: Array.isArray(formData.skuList) ? formData.skuList : []
})

/**
 * 提取套餐规格数据
 * @param {Object} formData 表单数据
 * @returns {Object} 套餐规格数据
 */
export const extractComboSpecData = (formData) => {
  // 只从comboSpecData中获取数据，不再使用根级别的price字段
  const comboData = formData.comboSpecData || {}
  
  console.log('extractComboSpecData - formData:', formData)
  console.log('extractComboSpecData - formData.comboSpecData:', formData.comboSpecData)
  console.log('extractComboSpecData - comboData:', comboData)
  console.log('extractComboSpecData - comboData.price:', comboData.price)
  
  // 套餐价格只从comboSpecData.price获取
  const comboPrice = comboData.price || 0
  
  const result = {
    price: comboPrice,
    comboGroups: comboData.comboGroups || []
  }
  
  console.log('extractComboSpecData - 返回结果:', result)
  console.log('extractComboSpecData - 最终price值:', result.price)
  return result
}

/**
 * 更新单规格数据到表单
 * @param {Object} formData 表单数据
 * @param {Object} data 单规格数据
 */
export const updateSingleSpecToForm = (formData, data) => {
  Object.keys(data).forEach(key => {
    if (data[key] !== undefined) {
      formData[key] = data[key]
    }
  })
}

/**
 * 更新多规格数据到表单
 * @param {Object} formData 表单数据
 * @param {Object} data 多规格数据
 */
export const updateMultiSpecToForm = (formData, data) => {
  if (data.specGroups !== undefined) {
    formData.specGroups = data.specGroups
  }
  if (data.skuList !== undefined) {
    formData.skuList = data.skuList
  }
}

/**
 * 更新套餐规格数据到表单
 * @param {Object} formData 表单数据
 * @param {Object} data 套餐规格数据
 */
export const updateComboSpecToForm = (formData, data) => {
  // 确保comboSpecData对象存在
  if (!formData.comboSpecData) {
    formData.comboSpecData = {}
  }
  
  // 更新套餐规格数据
  if (data.price !== undefined) {
    formData.comboSpecData.price = data.price
    formData.price = data.price // 同时更新根级别的price字段
  }
  if (data.comboGroups !== undefined) {
    formData.comboSpecData.comboGroups = data.comboGroups
    formData.comboGroups = data.comboGroups // 兼容性
  }
}

/**
 * 验证单规格数据
 * @param {Object} formData 表单数据
 * @returns {Object} 验证结果 { valid: boolean, message: string }
 */
export const validateSingleSpec = (formData) => {
  if (!formData.price || formData.price <= 0) {
    return { valid: false, message: '请设置商品价格，价格必须大于0' }
  }
  if (formData.stock < 0 || !Number.isInteger(Number(formData.stock))) {
    return { valid: false, message: '请设置商品库存，库存必须为非负整数' }
  }
  return { valid: true, message: '' }
}

/**
 * 验证多规格数据
 * @param {Object} formData 表单数据
 * @returns {Object} 验证结果 { valid: boolean, message: string }
 */
export const validateMultiSpec = (formData) => {
  if (!formData.specGroups || formData.specGroups.length === 0) {
    return { valid: false, message: '多规格商品必须设置规格组' }
  }

  const hasValidGroups = formData.specGroups.some(group => 
    group.name && group.values && group.values.length > 0
  )
  if (!hasValidGroups) {
    return { valid: false, message: '请完善规格组设置，至少需要一个完整的规格组' }
  }

  if (!formData.skuList || formData.skuList.length === 0) {
    return { valid: false, message: '多规格商品必须有SKU数据' }
  }

  const invalidSku = formData.skuList.find(sku => 
    !sku.price || sku.price < 0 || !Number.isFinite(Number(sku.price)) ||
    sku.stock < 0 || !Number.isInteger(Number(sku.stock))
  )
  if (invalidSku) {
    return { valid: false, message: '请检查SKU的价格、库存设置，价格必须大于0，库存必须为非负整数' }
  }

  return { valid: true, message: '' }
}

/**
 * 验证套餐规格数据
 * @param {Object} formData 表单数据
 * @returns {Object} 验证结果 { valid: boolean, message: string }
 */
export const validateComboSpec = (formData) => {
  if (formData.basePrice === undefined || formData.basePrice < 0) {
    return { valid: false, message: '请设置套餐基础价格，价格必须大于等于0' }
  }

  if (!formData.comboGroups || formData.comboGroups.length === 0) {
    return { valid: false, message: '套餐规格必须设置套餐组' }
  }

  const hasValidGroups = formData.comboGroups.some(group => 
    group.name && group.options && group.options.length > 0
  )
  if (!hasValidGroups) {
    return { valid: false, message: '请完善套餐组设置，至少需要一个完整的套餐组' }
  }

  return { valid: true, message: '' }
} 