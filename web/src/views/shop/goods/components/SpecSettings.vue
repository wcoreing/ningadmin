<template>
  <div class="spec-settings">
    <el-form-item label="规格设置">
      <el-radio-group v-model="formData.specType">
        <el-radio :label="1">单规格</el-radio>
        <el-radio :label="2">多规格</el-radio>
      </el-radio-group>
    </el-form-item>
    
    <!-- 多规格设置 -->
    <template v-if="formData.specType === 2">
      <div class="spec-container">
        <div class="spec-groups-wrapper">
          <!-- 规格组列表 -->
          <div v-for="(group, groupIndex) in formData.specGroups" :key="groupIndex" class="spec-group-item">
            <div class="spec-group-header">
              <el-input
                v-model="group.name"
                placeholder="请输入规格组名称"
                class="spec-group-name"
                size="small"
              />
              <el-button type="danger" link @click="removeSpecGroup(groupIndex)">
                <el-icon><Delete /></el-icon>
              </el-button>
            </div>
            
            <!-- 规格值列表 -->
            <div class="spec-values">
              <div class="spec-values-list">
                <el-tag
                  v-for="(value, valueIndex) in group.values"
                  :key="valueIndex"
                  closable
                  class="spec-value-tag"
                  size="small"
                  @close="handleRemoveSpecValue(groupIndex, valueIndex)"
                >
                  {{ value.name }}
                </el-tag>
              </div>
              
              <!-- 添加规格值 -->
              <div class="spec-value-input">
                <el-input
                  v-if="group.inputVisible"
                  ref="inputRefs"
                  v-model="group.inputValue"
                  size="small"
                  class="spec-input"
                  placeholder="请输入规格值"
                  @keyup.enter="handleInputConfirm(groupIndex)"
                  @blur="handleInputConfirm(groupIndex)"
                />
                <el-button
                  v-else
                  size="small"
                  type="primary"
                  link
                  @click="showInput(groupIndex)"
                >
                  <el-icon><Plus /></el-icon> 添加规格值
                </el-button>
              </div>
            </div>
          </div>
        </div>
        
        <!-- 添加规格组按钮 -->
        <el-button type="primary" class="add-spec-group" @click="addSpecGroup">
          <el-icon><Plus /></el-icon> 添加规格组
        </el-button>

        <!-- SKU 列表 -->
        <div v-if="hasValidSpecGroups" class="sku-table">
          <el-table :data="formData.skuList" border>
            <!-- 添加 SKU 图片列 -->
            <el-table-column label="规格图片" width="120">
              <template #default="{ row }">
                <div class="sku-image">
                  <el-image
                    v-if="row.image"
                    :src="row.image"
                    fit="cover"
                    class="preview-image"
                    :preview-src-list="[row.image]"
                    @click.stop="openFileManager('sku', row)"
                  />
                  <div v-else class="upload-placeholder" @click="openFileManager('sku', row)">
                    <el-icon><Plus /></el-icon>
                  </div>
                </div>
              </template>
            </el-table-column>

            <el-table-column label="规格组合">
              <template #default="{ row }">
                <div class="sku-specs">
                  <el-tag v-for="(value, key) in row.specs" :key="key" size="small">
                    {{ key }}: {{ value.name }}
                  </el-tag>
                </div>
              </template>
            </el-table-column>
            <el-table-column label="价格" width="200">
              <template #default="{ row }">
                <el-input-number
                  v-model="row.price"
                  :min="0"
                  :precision="2"
                  controls-position="right"
                  style="width: 160px"
                />
              </template>
            </el-table-column>
            <el-table-column label="库存" width="200">
              <template #default="{ row }">
                <el-input-number
                  v-model="row.stock"
                  :min="0"
                  controls-position="right"
                  style="width: 160px"
                />
              </template>
            </el-table-column>
          </el-table>
        </div>
      </div>
    </template>

    <!-- 添加文件管理器组件 -->
    <file-manager
      v-model="fileManagerVisible"
      :upload-url="uploadUrl"
      :multiple="false"
      :max-count="1"
      @select="handleFileSelect"
    />
  </div>
</template>

<script setup>
import { ref, computed, watch, nextTick } from 'vue'
import { Plus, Delete } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import FileManager from '@/components/FileManager/index.vue'

const props = defineProps({
  formData: {
    type: Object,
    required: true
  },
  uploadUrl: {
    type: String,
    required: true
  }
})

const inputRefs = ref([])

// 保存原始 SKU 数据的 Map
const originalSkuMap = ref(new Map())

// 添加文件管理器相关的响应式变量
const fileManagerVisible = ref(false)
const currentSku = ref(null)

// 修改处理文件选择的方法
const handleFileSelect = (files) => {
  if (currentSku.value) {
    // 确保 files 是数组，并且取第一个文件的 URL
    const fileUrl = Array.isArray(files) ? files[0]?.url : files?.url
    if (fileUrl) {
      currentSku.value.image = fileUrl
    }
  }
  fileManagerVisible.value = false
}

// 修改打开文件管理器的方法
const openFileManager = (type, sku) => {
  if (!sku) return
  currentSku.value = sku
  fileManagerVisible.value = true
}

// 添加规格组
const addSpecGroup = () => {
  const newGroup = {
    id: Date.now(),
    name: '',
    values: [],
    inputVisible: false,
    inputValue: ''
  }
  props.formData.specGroups.push(newGroup)
}

// 显示输入框
const showInput = (groupIndex) => {
  props.formData.specGroups[groupIndex].inputVisible = true
  // 等待DOM更新后聚焦
  nextTick(() => {
    inputRefs.value[groupIndex].focus()
  })
}

// 处理输入框失焦事件
const handleInputConfirm = (groupIndex) => {
  const group = props.formData.specGroups[groupIndex]
  const inputValue = group.inputValue?.trim()
  if (inputValue) {
    if (group.values && group.values.some(item => item.name === inputValue)) {
      ElMessage.warning('该规格值已存在')
      return
    }
    
    // 确保 values 数组已初始化
    if (!Array.isArray(group.values)) {
      group.values = []
    }
    group.values.push({
      id: Date.now(),
      name: inputValue,
      value: inputValue,
      sort: group.values.length + 1
    })
    // 触发 SKU 列表更新
    if (props.formData.specType === 2) {
      updateSkuList()
    }
  }
  group.inputVisible = false
  group.inputValue = ''
  
  // 更新规格组合
  generateSpecCombinations()
}

// 移除规格值
const handleRemoveSpecValue = (groupIndex, valueIndex) => {
  props.formData.specGroups[groupIndex].values.splice(valueIndex, 1)
}

// 移除规格组
const removeSpecGroup = (groupIndex) => {
  ElMessageBox.confirm('确定要删除该规格组吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    props.formData.specGroups.splice(groupIndex, 1)
  })
}

// 计算有效的规格组
const hasValidSpecGroups = computed(() => {
  return props.formData.specGroups?.some(group => 
    group.values?.length > 0 && group.name
  )
})

// 将 cartesian 函数重命名为 getCombinations
const getCombinations = (...arrays) => {
  return arrays.reduce((acc, curr) => {
    return acc.flatMap(a => curr.map(c => [...a, c]))
  }, [[]])
}

// 修改 initSkuList 方法
const initSkuList = () => {
  originalSkuMap.value.clear()
  
  if (!props.formData.skuList || !Array.isArray(props.formData.skuList)) {
    props.formData.skuList = []
    return
  }

  props.formData.skuList.forEach(sku => {
    if (!sku || !sku.specs) return

    try {
      const key = Object.entries(sku.specs)
        .sort(([a], [b]) => a.localeCompare(b))
        .map(([groupName, spec]) => `${groupName}:${spec.value}`)
        .join(';')
      
      originalSkuMap.value.set(key, {
        specs: sku.specs,
        price: Number(sku.price || 0),
        stock: Number(sku.stock || 0),
        image: sku.image || ''
      })
    } catch (error) {
      console.error('处理 SKU 数据出错:', error, sku)
    }
  })
}

// 修改 generateSkuList 方法
const generateSkuList = () => {
  const groups = props.formData.specGroups
  if (!groups || groups.length === 0) return []

  const validGroups = groups.filter(group => group.name && group.values?.length > 0)
  if (validGroups.length === 0) return []

  const specArrays = validGroups.map(group => {
    return {
      groupName: group.name,
      values: group.values
    }
  })

  const generateCombinations = (groups) => {
    if (groups.length === 0) return []
    if (groups.length === 1) {
      return groups[0].values.map(value => ({
        [groups[0].groupName]: value
      }))
    }

    const [first, ...rest] = groups
    const restCombinations = generateCombinations(rest)
    const result = []

    for (const value of first.values) {
      for (const combination of restCombinations) {
        result.push({
          [first.groupName]: value,
          ...combination
        })
      }
    }

    return result
  }

  const combinations = generateCombinations(specArrays)

  const newList = combinations.map(combination => {
    const specs = {}
    Object.entries(combination).forEach(([groupName, value]) => {
      specs[groupName] = {
        id: value.id,
        name: value.name,
        value: value.value,
        sort: value.sort
      }
    })

    const key = Object.entries(specs)
      .sort(([a], [b]) => a.localeCompare(b))
      .map(([groupName, spec]) => `${groupName}:${spec.value}`)
      .join(';')
    
    const originalSku = originalSkuMap.value.get(key)
    
    if (originalSku) {
      return {
        specs,
        price: Number(originalSku.price || 0),
        stock: Number(originalSku.stock || 0),
        image: originalSku.image || ''
      }
    }

    return {
      specs,
      price: 0,
      stock: 0,
      image: ''
    }
  })

  return newList
}

// 修改 updateSkuList 方法
const updateSkuList = () => {
  if (props.formData.specType === 2 && hasValidSpecGroups.value) {
    const currentSkuList = props.formData.skuList || []
    const newSkuList = generateSkuList()
    
    props.formData.skuList = newSkuList.map(newSku => {
      const existingSku = currentSkuList.find(sku => {
        if (!sku.specs) return false
        return Object.entries(newSku.specs).every(([groupName, spec]) => {
          const existingSpec = sku.specs[groupName]
          return existingSpec && existingSpec.value === spec.value
        })
      })

      if (existingSku) {
        return {
          specs: newSku.specs,
          price: Number(existingSku.price || 0),
          stock: Number(existingSku.stock || 0),
          image: existingSku.image || ''
        }
      }

      return {
        specs: newSku.specs,
        price: 0,
        stock: 0,
        image: ''
      }
    })
  }
}

// 监听规格类型变化
watch(
  () => props.formData.specType,
  (newType) => {
    if (newType === 2) {
      // 切换到多规格时
      if (!props.formData.specGroups?.length) {
        props.formData.specGroups = []
        addSpecGroup()
      }
      // 初始化 SKU 列表
      initSkuList()
    } else {
      // 切换到单规格时，清空规格组和 SKU 列表
      props.formData.specGroups = []
      props.formData.skuList = []
      originalSkuMap.value.clear()
    }
  },
  { immediate: true }
)

// 监听规格变化
watch(
  () => [
    props.formData.specGroups?.length,
    ...props.formData.specGroups?.map(g => [g.name, g.values?.length])
  ],
  () => {
    updateSkuList()
  },
  { deep: true }
)
</script>

<style lang="scss" scoped>
.spec-container {
  padding: 10px;
  background: #f5f7fa;
  border-radius: 4px;

  .spec-groups-wrapper {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 10px;
  }

  .spec-group-item {
    padding: 10px;
    background: #fff;
    border: 1px solid #dcdfe6;
    border-radius: 4px;

    .spec-group-header {
      display: flex;
      align-items: center;
      margin-bottom: 10px;
      gap: 8px;

      .spec-group-name {
        width: 200px;
      }
    }

    .spec-values {
      padding-left: 0;

      .spec-values-list {
        display: flex;
        flex-wrap: wrap;
        gap: 4px;
        margin-bottom: 8px;
      }

      .spec-value-input {
        display: flex;
        align-items: center;
        gap: 8px;

        .spec-input {
          width: 200px;
        }
      }
    }
  }

  .add-spec-group {
    width: 100%;
    height: 40px;
    margin-top: 10px;
    border: 1px dashed #dcdfe6;
    
    &:hover {
      border-color: var(--el-color-primary);
      color: var(--el-color-primary);
    }
  }
}

.sku-table {
  margin-top: 20px;
  
  .sku-specs {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
  }

  :deep(.el-input-number) {
    width: 160px;
  }

  .sku-image {
    width: 60px;
    height: 60px;
    border: 1px dashed #dcdfe6;
    border-radius: 4px;
    overflow: hidden;
    cursor: pointer;

    .preview-image {
      width: 100%;
      height: 100%;
      object-fit: cover;
    }

    .upload-placeholder {
      width: 100%;
      height: 100%;
      display: flex;
      align-items: center;
      justify-content: center;
      color: #909399;
      font-size: 20px;

      &:hover {
        color: var(--el-color-primary);
        border-color: var(--el-color-primary);
      }
    }
  }
}
</style> 