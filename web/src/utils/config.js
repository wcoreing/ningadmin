/**
 * 配置管理工具类
 */

// 配置类型列表
export const CONFIG_TYPES = [
  { label: '字符串', value: 'string' },
  { label: '数字', value: 'number' },
  { label: '布尔值', value: 'boolean' },
  { label: '日期', value: 'date' },
  { label: '图片', value: 'image' }
]

// 配置类型标签样式映射
export const CONFIG_TYPE_TAGS = {
  string: '',
  number: 'success',
  boolean: 'warning',
  date: 'info',
  image: 'danger'
}

// 根据配置类型获取标签类型
export function getConfigTypeTag(type) {
  return CONFIG_TYPE_TAGS[type] || ''
}

// 根据配置类型获取显示名称
export function getConfigTypeName(type) {
  const configType = CONFIG_TYPES.find(item => item.value === type)
  return configType ? configType.label : type
}

// 配置值类型转换
export function parseConfigValue(value, type) {
  if (value === null || value === undefined) {
    return value
  }
  
  switch (type) {
    case 'number':
      return Number(value)
    case 'boolean':
      return Boolean(value)
    case 'date':
      return value ? new Date(value).toISOString() : null
    default:
      return String(value)
  }
}

// 配置值格式化（用于显示）
export function formatConfigValue(value, type) {
  if (value === null || value === undefined) {
    return '-'
  }
  
  switch (type) {
    case 'boolean':
      return value ? '是' : '否'
    case 'date':
      return value ? new Date(value).toLocaleString() : '-'
    default:
      return value.toString()
  }
}

// 验证配置值格式
export function validateConfigValue(value, type) {
  if (value === null || value === undefined || value === '') {
    return true // 允许空值
  }
  
  switch (type) {
    case 'number':
      return !isNaN(Number(value))
    case 'boolean':
      return typeof value === 'boolean' || ['true', 'false'].includes(value.toLowerCase())
    case 'date':
      return !isNaN(Date.parse(value))
    case 'image':
      return typeof value === 'string' && value.length > 0
    default:
      return true
  }
} 