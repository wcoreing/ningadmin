// 规格工具函数 - 支持多规格数据同时存储

/**
 * 从表单数据中提取所有规格类型的数据
 * @param {Object} formData - 表单数据
 * @returns {Object} 包含所有规格类型数据的对象
 */
export function extractAllSpecData(formData) {
  return {
    // 单规格数据
    singleSpecData: extractSingleSpecData(formData),
    // 多规格数据  
    multiSpecData: extractMultiSpecData(formData),
    // 套餐规格数据
    comboSpecData: extractComboSpecData(formData)
  }
}

/**
 * 提取单规格数据
 * @param {Object} formData - 表单数据
 * @returns {Object|null} 单规格数据
 */
export function extractSingleSpecData(formData) {
  // 如果有单规格相关数据，提取它
  if (formData.price !== undefined || formData.stock !== undefined) {
    return {
      price: formData.price || 0,
      originalPrice: formData.originalPrice || 0,
      costPrice: formData.costPrice || 0,
      stock: formData.stock || 0,
      weight: formData.weight || 0,
      volume: formData.volume || 0
    }
  }
  return null
}

/**
 * 提取多规格数据
 * @param {Object} formData - 表单数据
 * @returns {Object|null} 多规格数据
 */
export function extractMultiSpecData(formData) {
  if (formData.specGroups && formData.specGroups.length > 0) {
    return {
      specGroups: formData.specGroups || [],
      skuList: formData.skuList || []
    }
  }
  return null
}

/**
 * 提取套餐规格数据
 * @param {Object} formData - 表单数据
 * @returns {Object|null} 套餐规格数据
 */
export function extractComboSpecData(formData) {
  if (formData.comboGroups && formData.comboGroups.length > 0) {
    return {
      comboGroups: formData.comboGroups || []
    }
  }
  return null
}

/**
 * 根据规格类型获取对应的数据
 * @param {Object} allSpecData - 所有规格数据
 * @param {number} specType - 规格类型 1-单规格 2-多规格 3-套餐规格
 * @returns {Object} 对应规格类型的数据
 */
export function getSpecDataByType(allSpecData, specType) {
  switch (specType) {
    case 1: // 单规格
      return {
        price: allSpecData.singleSpecData?.price || 0,
        originalPrice: allSpecData.singleSpecData?.originalPrice || 0,
        costPrice: allSpecData.singleSpecData?.costPrice || 0,
        stock: allSpecData.singleSpecData?.stock || 0,
        weight: allSpecData.singleSpecData?.weight || 0,
        volume: allSpecData.singleSpecData?.volume || 0
      }
    case 2: // 多规格
      return {
        specGroups: allSpecData.multiSpecData?.specGroups || [],
        skuList: allSpecData.multiSpecData?.skuList || []
      }
    case 3: // 套餐规格
      return {
        comboGroups: allSpecData.comboSpecData?.comboGroups || []
      }
    default:
      return {}
  }
}

/**
 * 更新指定规格类型的数据
 * @param {Object} allSpecData - 所有规格数据
 * @param {number} specType - 规格类型
 * @param {Object} newData - 新数据
 * @returns {Object} 更新后的所有规格数据
 */
export function updateSpecDataByType(allSpecData, specType, newData) {
  const updated = { ...allSpecData }
  
  switch (specType) {
    case 1: // 单规格
      updated.singleSpecData = {
        price: newData.price || 0,
        originalPrice: newData.originalPrice || 0,
        costPrice: newData.costPrice || 0,
        stock: newData.stock || 0,
        weight: newData.weight || 0,
        volume: newData.volume || 0
      }
      break
    case 2: // 多规格
      updated.multiSpecData = {
        specGroups: newData.specGroups || [],
        skuList: newData.skuList || []
      }
      break
    case 3: // 套餐规格
      updated.comboSpecData = {
        comboGroups: newData.comboGroups || []
      }
      break
  }
  
  return updated
}

/**
 * 将后端返回的数据转换为前端格式
 * @param {Object} backendData - 后端返回的数据
 * @returns {Object} 前端格式的数据
 */
export function convertBackendToFrontend(backendData) {
  const result = {
    // 基础商品信息
    goodsId: backendData.goodsId,
    goodsName: backendData.goodsName,
    categoryId: backendData.categoryId,
    brand: backendData.brand,
    supplier: backendData.supplier,
    type: backendData.type,
    sort: backendData.sort,
    status: backendData.status,
    description: backendData.description,
    coverImg: backendData.coverImg,
    tagIds: backendData.tagIds || [],
    images: backendData.images || [],
    specType: backendData.specType,
    detailDesc: backendData.detailDesc,
    
    // 多规格数据存储
    allSpecData: {
      singleSpecData: backendData.singleSpecData,
      multiSpecData: backendData.multiSpecData,
      comboSpecData: backendData.comboSpecData
    }
  }

  // 兼容旧版本：根据当前specType设置对应的数据到根级别
  const currentSpecData = getSpecDataByType(result.allSpecData, backendData.specType)
  Object.assign(result, currentSpecData)

  return result
}

/**
 * 将前端数据转换为后端格式
 * @param {Object} frontendData - 前端数据
 * @returns {Object} 后端格式的数据
 */
export function convertFrontendToBackend(frontendData) {
  const result = {
    // 基础商品信息
    goodsId: frontendData.goodsId,
    goodsName: frontendData.goodsName,
    categoryId: frontendData.categoryId,
    brand: frontendData.brand,
    supplier: frontendData.supplier,
    type: frontendData.type,
    sort: frontendData.sort,
    status: frontendData.status,
    description: frontendData.description,
    coverImg: frontendData.coverImg,
    tagIds: frontendData.tagIds || [],
    images: frontendData.images || [],
    specType: frontendData.specType,
    detailDesc: frontendData.detailDesc
  }

  // 如果有allSpecData，使用新格式
  if (frontendData.allSpecData) {
    result.singleSpecData = frontendData.allSpecData.singleSpecData
    result.multiSpecData = frontendData.allSpecData.multiSpecData
    result.comboSpecData = frontendData.allSpecData.comboSpecData
  } else {
    // 兼容旧版本：从根级别提取数据
    result.singleSpecData = extractSingleSpecData(frontendData)
    result.multiSpecData = extractMultiSpecData(frontendData)
    result.comboSpecData = extractComboSpecData(frontendData)
  }

  // 兼容旧版本字段
  if (frontendData.specType === 1) {
    Object.assign(result, {
      price: frontendData.price || 0,
      originalPrice: frontendData.originalPrice || 0,
      costPrice: frontendData.costPrice || 0,
      stock: frontendData.stock || 0,
      weight: frontendData.weight || 0,
      volume: frontendData.volume || 0
    })
  } else if (frontendData.specType === 2) {
    result.specGroups = frontendData.specGroups || []
    result.skuList = frontendData.skuList || []
  } else if (frontendData.specType === 3) {
    result.price = frontendData.price || 0
    result.comboGroups = frontendData.comboGroups || []
  }

  return result
}

/**
 * 验证规格数据的完整性
 * @param {Object} specData - 规格数据
 * @param {number} specType - 规格类型
 * @returns {Object} 验证结果 {valid: boolean, errors: string[]}
 */
export function validateSpecData(specData, specType) {
  const errors = []
  
  switch (specType) {
    case 1: // 单规格
      if (!specData.singleSpecData) {
        errors.push('缺少单规格数据')
      } else {
        if (specData.singleSpecData.price <= 0) {
          errors.push('单规格价格必须大于0')
        }
        if (specData.singleSpecData.stock < 0) {
          errors.push('单规格库存不能小于0')
        }
      }
      break
      
    case 2: // 多规格
      if (!specData.multiSpecData) {
        errors.push('缺少多规格数据')
      } else {
        if (!specData.multiSpecData.specGroups || specData.multiSpecData.specGroups.length === 0) {
          errors.push('多规格必须至少有一个规格组')
        }
        if (!specData.multiSpecData.skuList || specData.multiSpecData.skuList.length === 0) {
          errors.push('多规格必须至少有一个SKU')
        }
      }
      break
      
    case 3: // 套餐规格
      if (!specData.comboSpecData) {
        errors.push('缺少套餐规格数据')
      } else {
        if (!specData.comboSpecData.comboGroups || specData.comboSpecData.comboGroups.length === 0) {
          errors.push('套餐规格必须至少有一个套餐组')
        }
        // 套餐价格验证移到外层price字段
      }
      break
      
    default:
      errors.push('无效的规格类型')
  }
  
  return {
    valid: errors.length === 0,
    errors
  }
}

/**
 * 获取规格数据的摘要信息
 * @param {Object} allSpecData - 所有规格数据
 * @returns {Object} 摘要信息
 */
export function getSpecDataSummary(allSpecData) {
  const summary = {
    singleSpec: {
      hasData: !!allSpecData.singleSpecData,
      price: allSpecData.singleSpecData?.price || 0,
      stock: allSpecData.singleSpecData?.stock || 0
    },
    multiSpec: {
      hasData: !!allSpecData.multiSpecData,
      groupCount: allSpecData.multiSpecData?.specGroups?.length || 0,
      skuCount: allSpecData.multiSpecData?.skuList?.length || 0
    },
    comboSpec: {
      hasData: !!allSpecData.comboSpecData,
      groupCount: allSpecData.comboSpecData?.comboGroups?.length || 0
    }
  }
  
  return summary
}

/**
 * 清空指定规格类型的数据
 * @param {Object} allSpecData - 所有规格数据
 * @param {number} specType - 要清空的规格类型
 * @returns {Object} 清空后的数据
 */
export function clearSpecDataByType(allSpecData, specType) {
  const updated = { ...allSpecData }
  
  switch (specType) {
    case 1:
      updated.singleSpecData = null
      break
    case 2:
      updated.multiSpecData = null
      break
    case 3:
      updated.comboSpecData = null
      break
  }
  
  return updated
}

/**
 * 深度克隆规格数据
 * @param {Object} specData - 规格数据
 * @returns {Object} 克隆后的数据
 */
export function cloneSpecData(specData) {
  return JSON.parse(JSON.stringify(specData))
} 