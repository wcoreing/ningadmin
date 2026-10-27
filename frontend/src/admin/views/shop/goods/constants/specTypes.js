/**
 * 商品规格类型常量
 */
export const SPEC_TYPES = {
  SINGLE: 1,    // 单规格
  MULTI: 2,     // 多规格
  COMBO: 3      // 套餐规格
}

/**
 * 规格类型配置
 */
export const SPEC_TYPE_CONFIG = {
  [SPEC_TYPES.SINGLE]: {
    name: '单规格',
    description: '商品只有一种规格，统一价格和库存',
    icon: 'Document',
    fields: ['price', 'originalPrice', 'costPrice', 'stock', 'weight', 'volume']
  },
  [SPEC_TYPES.MULTI]: {
    name: '多规格',
    description: '商品有多种规格组合，如颜色、尺寸等',
    icon: 'Grid',
    fields: ['specGroups', 'skuList']
  },
  [SPEC_TYPES.COMBO]: {
    name: '套餐规格',
    description: '商品为套餐形式，包含多个组和选项',
    icon: 'SetUp',
    fields: ['basePrice', 'specGroups', 'skuList', 'comboGroups']
  }
}

/**
 * 获取规格类型名称
 * @param {number} type 规格类型
 * @returns {string} 规格类型名称
 */
export const getSpecTypeName = (type) => {
  return SPEC_TYPE_CONFIG[type]?.name || '未知规格'
}

/**
 * 获取规格类型描述
 * @param {number} type 规格类型
 * @returns {string} 规格类型描述
 */
export const getSpecTypeDescription = (type) => {
  return SPEC_TYPE_CONFIG[type]?.description || ''
}

/**
 * 获取规格类型相关字段
 * @param {number} type 规格类型
 * @returns {string[]} 字段列表
 */
export const getSpecTypeFields = (type) => {
  return SPEC_TYPE_CONFIG[type]?.fields || []
}

/**
 * 验证规格类型是否有效
 * @param {number} type 规格类型
 * @returns {boolean} 是否有效
 */
export const isValidSpecType = (type) => {
  return Object.values(SPEC_TYPES).includes(type)
} 