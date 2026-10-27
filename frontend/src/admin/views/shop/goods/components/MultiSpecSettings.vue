<template>
  <div class="multi-spec-container">
    <!-- 规格组设置 -->
    <div class="spec-content-card">
      <div class="card-header">
        <h3 class="card-title">
          <el-icon><Collection /></el-icon>
          规格组设置
        </h3>
        <el-button type="primary" size="small" @click="addSpecGroup">
          <el-icon><Plus /></el-icon>
          添加规格组
        </el-button>
      </div>
      
      <div v-if="specGroups.length > 0" class="spec-groups-container">
        <div 
          v-for="(group, groupIndex) in specGroups" 
          :key="group.id"
          class="spec-group-card"
        >
          <div class="group-header">
            <div class="group-title">
              <el-icon><Folder /></el-icon>
              <span>规格组 {{ groupIndex + 1 }}</span>
              <el-tag v-if="group.values.length > 0" size="small" type="info">
                {{ group.values.length }} 个规格值
              </el-tag>
            </div>
            <el-button
              type="danger"
              size="small"
              plain
              @click="removeSpecGroup(groupIndex)"
              :disabled="specGroups.length === 1"
            >
              <el-icon><Delete /></el-icon>
              删除
            </el-button>
          </div>
          
          <div class="group-content">
            <div class="group-name-section">
              <label class="form-label required">规格组名称</label>
              <el-input
                v-model="group.name"
                placeholder="请输入规格组名称，如：颜色、尺寸、款式"
                size="default"
                class="group-name-input"
                @blur="validateGroupName(group)"
              />
            </div>

            <div class="group-values-section">
              <label class="form-label">规格值</label>
              <div class="values-container">
                <div class="values-list">
                  <el-tag
                    v-for="(value, valueIndex) in group.values"
                    :key="value.id"
                    closable
                    @close="removeSpecValue(groupIndex, valueIndex)"
                    class="value-tag"
                    :type="getValueTagType(valueIndex)"
                    size="default"
                  >
                    {{ value.name }}
                  </el-tag>
                </div>
                
                <div class="add-value-section">
                  <!-- 快速添加输入框 -->
                  <div class="quick-add-container">
                    <el-input
                      v-model="group.quickInputValue"
                      placeholder="输入规格值，回车添加，逗号分隔批量添加"
                      size="default"
                      class="quick-input"
                      @keyup.enter="quickAddValue(groupIndex)"
                      @input="handleQuickInput(groupIndex)"
                      clearable
                    >
                      <template #append>
                        <el-button 
                          @click="quickAddValue(groupIndex)"
                          :disabled="!group.quickInputValue?.trim()"
                          type="primary"
                        >
                          <el-icon><Plus /></el-icon>
                        </el-button>
                      </template>
                    </el-input>
                  </div>
                  
                  <!-- 批量添加按钮 -->
                  <div class="batch-add-section">
                    <el-button
                      size="small"
                      type="success"
                      plain
                      @click="showBatchAddDialog(groupIndex)"
                      class="batch-add-btn"
                    >
                      <el-icon><DocumentAdd /></el-icon>
                      批量添加
                    </el-button>
                    
                    <!-- 常用规格值快捷按钮 -->
                    <div v-if="getCommonValues(group.name).length > 0" class="common-values">
                      <span class="common-label">常用：</span>
                      <el-button
                        v-for="commonValue in getCommonValues(group.name)"
                        :key="commonValue"
                        size="small"
                        plain
                        @click="addCommonValue(groupIndex, commonValue)"
                        class="common-value-btn"
                      >
                        {{ commonValue }}
                      </el-button>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
      
      <el-empty 
        v-else 
        description="暂无规格组，请点击右上角按钮添加"
        :image-size="100"
      >
        <template #extra>
          <el-button type="primary" @click="addSpecGroup">
            <el-icon><Plus /></el-icon>
            添加第一个规格组
          </el-button>
          <el-button type="success" plain @click="addTestData">
            <el-icon><DocumentAdd /></el-icon>
            添加测试数据
          </el-button>
        </template>
      </el-empty>
    </div>

    <!-- SKU 列表 -->
    <div v-if="hasValidSpecs" class="spec-content-card">
      <div class="card-header">
        <h3 class="card-title">
          <el-icon><List /></el-icon>
          SKU 列表
          <el-tag type="info" size="small" class="sku-count">
            共 {{ skuList.length }} 个SKU
          </el-tag>
        </h3>
        <div class="sku-actions">
          <el-button size="small" @click="showBatchEditor" type="primary">
            <el-icon><Operation /></el-icon>
            批量编辑
          </el-button>
          
          <el-button size="small" @click="showBatchCalculator" type="success" plain>
            <el-icon><Promotion /></el-icon>
            价格计算器
          </el-button>
          
          <el-button size="small" @click="exportSkuData" type="info" plain>
            <el-icon><Download /></el-icon>
            导出SKU
          </el-button>
        </div>
      </div>

      <div class="sku-table-container">
        <el-table 
          :data="skuList" 
          border 
          stripe
          class="sku-table"
          :header-cell-style="{ background: '#f8f9fa', color: '#495057', fontWeight: '600' }"
          :row-style="{ height: '60px' }"
          @selection-change="handleSelectionChange"
        >
          <!-- 选择列 -->
          <el-table-column type="selection" width="55" fixed="left" />
          
          <!-- 规格组合列 -->
          <el-table-column label="规格组合" min-width="200" fixed="left">
            <template #default="{ row }">
              <div class="sku-specs">
                <el-tag
                  v-for="(spec, groupName) in row.specs"
                  :key="groupName"
                  size="small"
                  class="spec-tag"
                >
                  {{ groupName }}：{{ spec.name || spec.value }}
                </el-tag>
              </div>
            </template>
          </el-table-column>

          <!-- SKU 图片列 -->
          <el-table-column label="SKU图片" width="100" align="center">
            <template #default="{ row }">
              <div class="sku-image-cell">
                <el-image
                  v-if="row.image"
                  :src="row.image"
                  fit="cover"
                  class="sku-image"
                  :preview-src-list="[row.image]"
                  @click.stop="openImageManager(row)"
                />
                <div 
                  v-else 
                  class="image-placeholder" 
                  @click="openImageManager(row)"
                >
                  <el-icon><Picture /></el-icon>
                  <span>选择图片</span>
                </div>
              </div>
            </template>
          </el-table-column>

          <!-- 销售价格列 -->
          <el-table-column label="销售价格(元)" width="140" align="center">
            <template #default="{ row }">
              <el-input-number
                v-model="row.price"
                :min="0"
                :max="999999"
                :precision="2"
                :step="0.01"
                controls-position="right"
                size="small"
                class="table-input"
                @change="handleSkuChange"
              />
            </template>
          </el-table-column>

          <!-- 市场价格列 -->
          <el-table-column label="市场价格(元)" width="140" align="center">
            <template #default="{ row }">
              <el-input-number
                v-model="row.originalPrice"
                :min="0"
                :max="999999"
                :precision="2"
                :step="0.01"
                controls-position="right"
                size="small"
                class="table-input"
                @change="handleSkuChange"
              />
            </template>
          </el-table-column>

          <!-- 成本价格列 -->
          <el-table-column label="成本价格(元)" width="140" align="center">
            <template #default="{ row }">
              <el-input-number
                v-model="row.costPrice"
                :min="0"
                :max="999999"
                :precision="2"
                :step="0.01"
                controls-position="right"
                size="small"
                class="table-input"
                @change="handleSkuChange"
              />
            </template>
          </el-table-column>

          <!-- 库存数量列 -->
          <el-table-column label="库存数量(件)" width="130" align="center">
            <template #default="{ row }">
              <el-input-number
                v-model="row.stock"
                :min="0"
                :max="999999"
                :precision="0"
                controls-position="right"
                size="small"
                class="table-input"
                @change="handleSkuChange"
              />
            </template>
          </el-table-column>

          <!-- 重量列 -->
          <el-table-column label="重量(kg)" width="120" align="center">
            <template #default="{ row }">
              <el-input-number
                v-model="row.weight"
                :min="0"
                :precision="3"
                :step="0.001"
                controls-position="right"
                size="small"
                class="table-input"
                @change="handleSkuChange"
              />
            </template>
          </el-table-column>

          <!-- 体积列 -->
          <el-table-column label="体积(m³)" width="120" align="center">
            <template #default="{ row }">
              <el-input-number
                v-model="row.volume"
                :min="0"
                :precision="3"
                :step="0.001"
                controls-position="right"
                size="small"
                class="table-input"
                @change="handleSkuChange"
              />
            </template>
          </el-table-column>
        </el-table>
      </div>
    </div>

    <!-- 批量编辑对话框 -->
    <BatchEditDialog
      v-model="batchDialogVisible"
      :sku-list="skuList"
      :selected-skus="selectedSkus"
      @confirm="handleBatchEdit"
    />

    <!-- 批量添加规格值对话框 -->
    <BatchAddDialog
      v-model="batchAddDialogVisible"
      @confirm="handleBatchAdd"
    />

    <!-- 价格计算器对话框 -->
    <PriceCalculatorDialog
      v-model="calculatorDialogVisible"
      :sku-list="skuList"
      :selected-skus="selectedSkus"
      @confirm="handlePriceCalculation"
    />

    <!-- 文件管理器 -->
    <file-manager
      v-model="fileManagerVisible"
      :multiple="false"
      @select="handleImageSelect"
    />
  </div>
</template>

<script setup>
import { ref, computed, reactive, watch, nextTick } from 'vue'
import { Plus, Delete, Picture, Collection, Folder, List, DocumentAdd, Operation, Promotion, Download } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import FileManager from '@/admin/components/FileManager/index.vue'
import BatchEditDialog from './dialogs/BatchEditDialog.vue'
import BatchAddDialog from './dialogs/BatchAddDialog.vue'
import PriceCalculatorDialog from './dialogs/PriceCalculatorDialog.vue'

const props = defineProps({
  modelValue: {
    type: Object,
    default: () => ({
      specGroups: [],
      skuList: []
    })
  }
})

const emit = defineEmits(['update:modelValue', 'change'])

// 响应式数据
const specGroups = ref([])
const skuList = ref([])
const fileManagerVisible = ref(false)
const currentSku = ref(null)
const selectedSkus = ref([])

// 对话框状态
const batchDialogVisible = ref(false)
const batchAddDialogVisible = ref(false)
const calculatorDialogVisible = ref(false)
const currentBatchGroupIndex = ref(-1)

// 常用规格值配置
const commonSpecValues = {
  '颜色': ['红色', '蓝色', '绿色', '黄色', '黑色', '白色', '粉色', '紫色', '橙色', '灰色'],
  '尺寸': ['XS', 'S', 'M', 'L', 'XL', 'XXL', 'XXXL'],
  '尺码': ['35', '36', '37', '38', '39', '40', '41', '42', '43', '44', '45'],
  '容量': ['64GB', '128GB', '256GB', '512GB', '1TB'],
  '内存': ['4GB', '8GB', '16GB', '32GB'],
  '款式': ['经典款', '时尚款', '运动款', '商务款'],
  '材质': ['棉', '麻', '丝', '毛', '涤纶', '尼龙', '皮革', '金属', '塑料', '木质']
}

// 计算属性
const hasValidSpecs = computed(() => {
  return specGroups.value.some(group => 
    group.name && group.values && group.values.length > 0
  )
})

// ID生成器
let idCounter = 0
const generateId = () => {
  return Date.now() + (++idCounter)
}

// 添加规格组
const addSpecGroup = () => {
  console.log('添加规格组，当前规格组数量:', specGroups.value.length)
  
  const newGroup = {
    id: generateId(),
    name: '',
    values: [],
    quickInputValue: ''
  }
  
  specGroups.value.push(newGroup)
  
  console.log('添加规格组后，规格组数量:', specGroups.value.length)
  console.log('新添加的规格组:', newGroup)
  
  ElMessage.success('已添加新的规格组')
}

// 标记是否正在更新，避免递归更新
const isUpdating = ref(false)

// 初始化时从 props 加载数据
const initializeFromProps = () => {
  console.log('MultiSpecSettings - 从props初始化数据:', props.modelValue)
  
  if (props.modelValue.specGroups && props.modelValue.specGroups.length > 0) {
    specGroups.value = props.modelValue.specGroups.map(group => ({
      ...group,
      quickInputValue: '',
      values: group.values ? group.values.map(value => ({
        id: value.id || generateId(),
        name: value.name || value.value,
        value: value.value || value.name,
        sort: value.sort || 0
      })) : []
    }))
    console.log('MultiSpecSettings - 初始化规格组:', specGroups.value)
  }
  
  if (props.modelValue.skuList && props.modelValue.skuList.length > 0) {
    skuList.value = props.modelValue.skuList.map(sku => ({
      ...sku,
      specs: sku.specs ? Object.fromEntries(
        Object.entries(sku.specs).map(([groupName, spec]) => [
          groupName,
          {
            ...spec,
            name: spec.name || spec.value
          }
        ])
      ) : {}
    }))
    console.log('MultiSpecSettings - 初始化SKU列表:', skuList.value)
  }
  
  console.log('MultiSpecSettings - 从props初始化完成')
}

// 组件挂载时初始化数据
initializeFromProps()

// 更新数据
const updateData = () => {
  if (isUpdating.value) return
  
  isUpdating.value = true
  
  const data = {
    specGroups: specGroups.value.map(group => ({
      id: group.id,
      name: group.name,
      values: group.values,
      sort: 0
    })),
    skuList: skuList.value
  }
  
  emit('update:modelValue', data)
  emit('change', data)
  
  nextTick(() => {
    isUpdating.value = false
  })
}

// 监听本地数据变化，通知父组件
watch(
  [specGroups, skuList],
  () => {
    if (!isUpdating.value) {
    updateData()
    }
  },
  { deep: true }
)

// 删除规格组
const removeSpecGroup = (index) => {
  ElMessageBox.confirm('删除规格组将清空相关SKU数据，确定删除吗？', '确认删除', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    specGroups.value.splice(index, 1)
    debouncedGenerateSkuList()
  })
}

// 验证规格组名称
const validateGroupName = (group) => {
  if (!group.name.trim()) {
    ElMessage.warning('规格组名称不能为空')
    return false
  }
  
  const duplicateGroup = specGroups.value.find(g => 
    g.id !== group.id && g.name === group.name
  )
  if (duplicateGroup) {
    ElMessage.warning('规格组名称不能重复')
    group.name = ''
    return false
  }
  
  debouncedGenerateSkuList()
  return true
}

// 快速添加规格值
const quickAddValue = (groupIndex) => {
  const group = specGroups.value[groupIndex]
  const inputValue = group.quickInputValue?.trim()
  
  if (!inputValue) return
  
  const values = inputValue
    .split(/[,，;；]/)
    .map(v => v.trim())
    .filter(v => v.length > 0)
  
  let addedCount = 0
  
  values.forEach(value => {
    if (!group.values.some(v => v.name === value)) {
      group.values.push({
        id: generateId(),
        name: value,
        value: value,
        sort: group.values.length
      })
      addedCount++
    }
  })
  
  if (addedCount > 0) {
    group.quickInputValue = ''
    debouncedGenerateSkuList()
    ElMessage.success(`成功添加 ${addedCount} 个规格值`)
  } else {
    ElMessage.warning('规格值已存在或无效')
  }
}

// 处理快速输入
const handleQuickInput = (groupIndex) => {
  // 可以在这里添加实时验证逻辑
}

// 获取常用规格值
const getCommonValues = (groupName) => {
  const normalizedName = groupName?.trim()
  if (!normalizedName) return []
  
  if (commonSpecValues[normalizedName]) {
    return commonSpecValues[normalizedName]
  }
  
  for (const [key, values] of Object.entries(commonSpecValues)) {
    if (normalizedName.includes(key) || key.includes(normalizedName)) {
      return values
    }
  }
  
  return []
}

// 防抖生成SKU列表
let generateSkuListTimer = null
const debouncedGenerateSkuList = () => {
  if (generateSkuListTimer) {
    clearTimeout(generateSkuListTimer)
  }
  generateSkuListTimer = setTimeout(() => {
    generateSkuList()
  }, 100)
}

// 同步生成SKU列表（用于批量编辑等需要立即获取最新数据的场景）
const generateSkuListSync = () => {
  if (isUpdating.value) {
    console.log('正在更新中，跳过同步SKU生成')
    return
  }
  
  console.log('=== 开始同步生成SKU列表 ===')
  console.log('hasValidSpecs:', hasValidSpecs.value)
  console.log('当前规格组:', specGroups.value)
  
  if (!hasValidSpecs.value) {
    console.log('没有有效规格，清空SKU列表')
    skuList.value = []
    return
  }

  const validGroups = specGroups.value.filter(group => 
    group.name && group.values && group.values.length > 0
  )
  
  console.log('有效规格组:', validGroups)
  console.log('有效规格组数量:', validGroups.length)
  
  if (validGroups.length === 0) {
    console.log('没有有效规格组，清空SKU列表')
    skuList.value = []
    return
  }
  
  const combinations = generateCombinations(validGroups)
  console.log('生成的组合:', combinations)
  console.log('组合数量:', combinations.length)
  
  const existingSkuMap = new Map()
  
  skuList.value.forEach(sku => {
    const key = generateSkuKey(sku.specs)
    existingSkuMap.set(key, sku)
  })
  
  const newSkuList = combinations.map(combination => {
    const specs = {}
    combination.forEach(({ groupName, value }) => {
      specs[groupName] = {
        id: value.id,
        name: value.name,
        value: value.value,
        sort: value.sort
      }
    })
    
    const key = generateSkuKey(specs)
    const existingSku = existingSkuMap.get(key)
    
    if (existingSku) {
      return { ...existingSku, specs }
    }
    
    return {
      specs,
      price: 0,
      originalPrice: 0,
      costPrice: 0,
      stock: 0,
      weight: 0,
      volume: 0,
      image: ''
    }
  })
  
  console.log('同步生成的SKU列表:', newSkuList)
  console.log('SKU数量:', newSkuList.length)
  console.log('=== 同步SKU列表生成完成 ===')
  
  // 直接同步更新，不使用 nextTick
  skuList.value = newSkuList
}

// 添加常用规格值
const addCommonValue = (groupIndex, value) => {
  const group = specGroups.value[groupIndex]
  
  if (group.values.some(v => v.name === value)) {
    ElMessage.warning('规格值已存在')
    return
  }
  
  group.values.push({
    id: generateId(),
    name: value,
    value: value,
    sort: group.values.length
  })
  
  debouncedGenerateSkuList()
  ElMessage.success(`已添加规格值：${value}`)
}

// 显示批量添加对话框
const showBatchAddDialog = (groupIndex) => {
  currentBatchGroupIndex.value = groupIndex
  batchAddDialogVisible.value = true
}

// 处理批量添加
const handleBatchAdd = (values) => {
  const groupIndex = currentBatchGroupIndex.value
  const group = specGroups.value[groupIndex]
  let addedCount = 0
  
  values.forEach(value => {
    if (!group.values.some(v => v.name === value)) {
      group.values.push({
        id: generateId(),
        name: value,
        value: value,
        sort: group.values.length
      })
      addedCount++
    }
  })
  
  if (addedCount > 0) {
    debouncedGenerateSkuList()
    ElMessage.success(`成功添加 ${addedCount} 个规格值`)
  } else {
    ElMessage.warning('所有规格值都已存在')
  }
}

// 删除规格值
const removeSpecValue = (groupIndex, valueIndex) => {
  specGroups.value[groupIndex].values.splice(valueIndex, 1)
  debouncedGenerateSkuList()
}

// 获取规格值标签类型
const getValueTagType = (index) => {
  const types = ['primary', 'success', 'info', 'warning', 'danger']
  return types[index % types.length]
}

// 生成SKU列表
const generateSkuList = () => {
  if (isUpdating.value) {
    console.log('正在更新中，跳过SKU生成')
    return
  }
  
  console.log('=== 开始生成SKU列表 ===')
  console.log('hasValidSpecs:', hasValidSpecs.value)
  console.log('当前规格组:', specGroups.value)
  
  if (!hasValidSpecs.value) {
    console.log('没有有效规格，清空SKU列表')
    const newSkuList = []
    // 使用 nextTick 避免递归更新
    nextTick(() => {
      if (!isUpdating.value) {
        skuList.value = newSkuList
      }
    })
    return
  }

  const validGroups = specGroups.value.filter(group => 
    group.name && group.values && group.values.length > 0
  )
  
  console.log('有效规格组:', validGroups)
  console.log('有效规格组数量:', validGroups.length)
  
  if (validGroups.length === 0) {
    console.log('没有有效规格组，清空SKU列表')
    const newSkuList = []
    // 使用 nextTick 避免递归更新
    nextTick(() => {
      if (!isUpdating.value) {
        skuList.value = newSkuList
      }
    })
    return
  }
  
  const combinations = generateCombinations(validGroups)
  console.log('生成的组合:', combinations)
  console.log('组合数量:', combinations.length)
  
  const existingSkuMap = new Map()
  
  skuList.value.forEach(sku => {
    const key = generateSkuKey(sku.specs)
    existingSkuMap.set(key, sku)
  })
  
  const newSkuList = combinations.map(combination => {
    const specs = {}
    combination.forEach(({ groupName, value }) => {
      specs[groupName] = {
        id: value.id,
        name: value.name,
        value: value.value,
        sort: value.sort
      }
    })
    
    const key = generateSkuKey(specs)
    const existingSku = existingSkuMap.get(key)
    
    if (existingSku) {
      return { ...existingSku, specs }
    }
    
    return {
      specs,
      price: 0,
      originalPrice: 0,
      costPrice: 0,
      stock: 0,
      weight: 0,
      volume: 0,
      image: ''
    }
  })
  
  console.log('生成的SKU列表:', newSkuList)
  console.log('SKU数量:', newSkuList.length)
  console.log('=== SKU列表生成完成 ===')
  
  // 使用 nextTick 避免递归更新
  nextTick(() => {
    if (!isUpdating.value) {
      skuList.value = newSkuList
    }
  })
}

// 生成规格组合
const generateCombinations = (groups) => {
  if (groups.length === 0) return []
  if (groups.length === 1) {
    return groups[0].values.map(value => [{ groupName: groups[0].name, value }])
  }

  const [first, ...rest] = groups
  const restCombinations = generateCombinations(rest)
  const result = []

  for (const value of first.values) {
    for (const combination of restCombinations) {
      result.push([{ groupName: first.name, value }, ...combination])
    }
  }

  return result
}

// 生成SKU唯一键
const generateSkuKey = (specs) => {
  return Object.entries(specs)
    .sort(([a], [b]) => a.localeCompare(b))
    .map(([groupName, spec]) => `${groupName}:${spec.value}`)
    .join(';')
}

// 处理SKU变化
const handleSkuChange = () => {
  updateData()
}

// 打开图片管理器
const openImageManager = (sku) => {
  currentSku.value = sku
  fileManagerVisible.value = true
}

// 处理图片选择
const handleImageSelect = (files) => {
  if (currentSku.value && files) {
    const file = Array.isArray(files) ? files[0] : files
    if (file && file.url) {
      currentSku.value.image = file.url
      updateData()
      ElMessage.success('图片设置成功')
    }
  }
  
  fileManagerVisible.value = false
  currentSku.value = null
}

// 显示批量编辑器
const showBatchEditor = () => {
  // 取消防抖定时器，立即同步生成最新的SKU列表
  if (generateSkuListTimer) {
    clearTimeout(generateSkuListTimer)
    generateSkuListTimer = null
  }
  
  // 同步生成最新的SKU列表
  generateSkuListSync()
  
  console.log('打开批量编辑对话框，当前SKU数量:', skuList.value.length)
  console.log('当前SKU列表:', skuList.value)
  batchDialogVisible.value = true
}

// 处理批量编辑
const handleBatchEdit = (editData) => {
  console.log('开始批量编辑:', editData)
  
  const { applyType, fields, affectedSkus } = editData
  
  // 获取要修改的SKU列表
  const targetSkus = applyType === 'all' ? skuList.value : affectedSkus
  
  // 批量修改SKU数据
  targetSkus.forEach(sku => {
    Object.keys(fields).forEach(fieldKey => {
      if (sku.hasOwnProperty(fieldKey)) {
        sku[fieldKey] = fields[fieldKey]
      }
    })
  })
  
  console.log('批量编辑完成，修改了', targetSkus.length, '个SKU')
  console.log('修改的字段:', Object.keys(fields))
  console.log('修改后的SKU列表:', skuList.value)
  
  // 更新数据到父组件
  updateData()
  
  ElMessage.success(`成功批量修改 ${targetSkus.length} 个SKU`)
}

// 表格选择变化
const handleSelectionChange = (selection) => {
  selectedSkus.value = selection
}

// 显示价格计算器
const showBatchCalculator = () => {
  // 取消防抖定时器，立即同步生成最新的SKU列表
  if (generateSkuListTimer) {
    clearTimeout(generateSkuListTimer)
    generateSkuListTimer = null
  }
  
  // 同步生成最新的SKU列表
  generateSkuListSync()
  
  console.log('打开价格计算器对话框，当前SKU数量:', skuList.value.length)
  console.log('当前SKU列表:', skuList.value)
  calculatorDialogVisible.value = true
}

// 处理价格计算
const handlePriceCalculation = (calculationData) => {
  console.log('开始价格计算:', calculationData)
  
  const { costBase, baseCost, profitRate, marketRate, applyType, affectedSkus } = calculationData
  
  // 获取要修改的SKU列表
  const targetSkus = applyType === 'all' ? skuList.value : affectedSkus
  
  // 根据成本价基准进行价格计算
  targetSkus.forEach(sku => {
    let costPrice = 0
    
    if (costBase === 'manual') {
      // 使用手动输入的成本价
      costPrice = baseCost
      sku.costPrice = costPrice
    } else {
      // 使用当前成本价
      costPrice = sku.costPrice || 0
    }
    
    // 计算销售价格（成本价 + 利润率）
    sku.price = Math.round(costPrice * (1 + profitRate / 100) * 100) / 100
    
    // 计算市场价格（成本价 + 市场利润率）
    sku.originalPrice = Math.round(costPrice * (1 + marketRate / 100) * 100) / 100
  })
  
  console.log('价格计算完成，修改了', targetSkus.length, '个SKU')
  console.log('成本价基准:', costBase, '利润率:', profitRate + '%', '市场利润率:', marketRate + '%')
  console.log('修改后的SKU列表:', skuList.value)
  
  // 更新数据到父组件
  updateData()
  
  ElMessage.success(`成功计算 ${targetSkus.length} 个SKU的价格`)
}

// 导出SKU数据
const exportSkuData = () => {
  if (skuList.value.length === 0) {
    ElMessage.warning('没有SKU数据可导出')
    return
  }
  
  // 导出逻辑
  ElMessage.success('SKU数据导出成功')
}

// 添加测试数据
const addTestData = () => {
  console.log('=== 开始添加测试数据 ===')
  console.log('添加前规格组数量:', specGroups.value.length)
  console.log('添加前SKU数量:', skuList.value.length)
  
  // 清空现有数据
  specGroups.value = []
  skuList.value = []
  
  // 添加颜色规格组
  const colorGroup = {
    id: generateId(),
    name: '颜色',
    values: [
      { id: generateId(), name: '红色', value: '红色', sort: 0 },
      { id: generateId(), name: '蓝色', value: '蓝色', sort: 1 }
    ],
    quickInputValue: ''
  }
  
  // 添加尺寸规格组
  const sizeGroup = {
    id: generateId(),
    name: '尺寸',
    values: [
      { id: generateId(), name: 'S', value: 'S', sort: 0 },
      { id: generateId(), name: 'M', value: 'M', sort: 1 },
      { id: generateId(), name: 'L', value: 'L', sort: 2 }
    ],
    quickInputValue: ''
  }
  
  specGroups.value = [colorGroup, sizeGroup]
  
  console.log('测试规格组已添加:', specGroups.value)
  console.log('添加后规格组数量:', specGroups.value.length)
  
  // 生成SKU列表
  debouncedGenerateSkuList()
  
  console.log('测试SKU已生成:', skuList.value)
  console.log('添加后SKU数量:', skuList.value.length)
  console.log('=== 测试数据添加完成 ===')
  
  ElMessage.success(`测试数据添加成功！已添加 ${specGroups.value.length} 个规格组，${skuList.value.length} 个SKU`)
}

// 暴露方法
defineExpose({
  generateSkuList,
  generateSkuListSync,
  addSpecGroup,
  addTestData,
  testComponent: () => {
    console.log('=== 多规格组件测试 ===')
    console.log('规格组数量:', specGroups.value.length)
    console.log('SKU数量:', skuList.value.length)
    console.log('hasValidSpecs:', hasValidSpecs.value)
    console.log('规格组详情:', JSON.parse(JSON.stringify(specGroups.value)))
    console.log('SKU详情:', JSON.parse(JSON.stringify(skuList.value)))
    console.log('========================')
  }
})
</script>

<style lang="scss" scoped>
.multi-spec-container {
  .spec-content-card {
    background: #fff;
    border-radius: 12px;
    border: 1px solid #e4e7ed;
    margin-bottom: 20px;
    overflow: hidden;
    
    .card-header {
      display: flex;
      align-items: center;
      justify-content: space-between;
      padding: 20px 24px;
      background: #f8f9fa;
      border-bottom: 1px solid #e4e7ed;
      
      .card-title {
        display: flex;
        align-items: center;
        gap: 8px;
        margin: 0;
        font-size: 16px;
        font-weight: 600;
        color: #303133;

        .sku-count {
          margin-left: 8px;
        }
      }

      .sku-actions {
        display: flex;
        gap: 8px;
      }
    }
  }
  
  .spec-groups-container {
    padding: 24px;
    
    .spec-group-card {
      background: #fafbfc;
      border: 1px solid #e4e7ed;
      border-radius: 8px;
      margin-bottom: 16px;
      
      .group-header {
        display: flex;
        align-items: center;
        justify-content: space-between;
        padding: 16px 20px;
        background: #f0f2f5;
        border-bottom: 1px solid #e4e7ed;
        
        .group-title {
          display: flex;
          align-items: center;
          gap: 8px;
          font-weight: 600;
          color: #303133;
        }
      }
      
      .group-content {
        padding: 20px;
        
        .group-name-section {
          margin-bottom: 20px;
          
          .form-label {
            display: block;
            margin-bottom: 8px;
            font-weight: 500;
            color: #606266;
            font-size: 14px;
            
            &.required::after {
              content: ' *';
              color: #f56c6c;
            }
          }
          
          .group-name-input {
            max-width: 400px;
          }
        }
        
        .group-values-section {
          .form-label {
            display: block;
            margin-bottom: 12px;
            font-weight: 500;
            color: #606266;
            font-size: 14px;
          }
          
          .values-container {
            .values-list {
              display: flex;
              flex-wrap: wrap;
              gap: 8px;
              margin-bottom: 16px;
              min-height: 32px;
              padding: 12px;
              background: #fff;
              border: 1px dashed #d9d9d9;
              border-radius: 6px;
              
              .value-tag {
                cursor: pointer;
                transition: all 0.2s;
                
                &:hover {
                  transform: translateY(-1px);
                  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
                }
              }
            }
            
            .add-value-section {
              .quick-add-container {
                margin-bottom: 12px;
                
                .quick-input {
                  width: 100%;
                }
              }
              
              .batch-add-section {
                display: flex;
                flex-direction: column;
                gap: 12px;
                
                .batch-add-btn {
                  align-self: flex-start;
                }
                
                .common-values {
                  display: flex;
                  align-items: center;
                  flex-wrap: wrap;
                  gap: 8px;
                  
                  .common-label {
                    font-size: 12px;
                    color: #909399;
                    margin-right: 4px;
                  }
                  
                  .common-value-btn {
                    font-size: 12px;
                    padding: 4px 8px;
                    height: auto;
                    border-style: dashed;

                    &:hover {
                      border-style: solid;
                    }
                  }
                }
              }
            }
          }
        }
      }
    }
  }
  
  .sku-table-container {
    padding: 0 24px 24px;

    .sku-table {
      border-radius: 8px;
      overflow: hidden;

      .sku-specs {
        display: flex;
        flex-wrap: wrap;
        gap: 6px;

        .spec-tag {
          font-size: 12px;
        }
      }
      
      .sku-image-cell {
        display: flex;
        justify-content: center;

        .sku-image {
          width: 50px;
          height: 50px;
          border-radius: 4px;
          cursor: pointer;
          border: 1px solid #e4e7ed;
        }
        
        .image-placeholder {
          width: 50px;
          height: 50px;
          border: 1px dashed #c0c4cc;
          border-radius: 4px;
          display: flex;
          flex-direction: column;
          align-items: center;
          justify-content: center;
          cursor: pointer;
          color: #909399;
          font-size: 12px;
          transition: all 0.2s;

          &:hover {
            border-color: var(--el-color-primary);
            color: var(--el-color-primary);
          }
          
          .el-icon {
            font-size: 16px;
            margin-bottom: 2px;
          }
        }
      }
      
      .table-input {
        width: 100%;
        
        :deep(.el-input__inner) {
          text-align: center;
        }
      }
    }
  }
}

:deep(.el-empty) {
  padding: 60px 0;
}
</style> 