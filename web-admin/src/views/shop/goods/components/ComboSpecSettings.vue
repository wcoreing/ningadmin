<template>
  <div class="combo-spec-container">
    <!-- 基础价格设置 -->
    <div class="spec-content-card">
      <div class="card-header">
        <h3 class="card-title">
          <el-icon><Money /></el-icon>
          套餐基础设置
        </h3>
      </div>
      
      <div class="base-settings">
        <div class="form-group">
          <label class="form-label required">套餐价格3</label>
          <div class="input-with-unit">
            <el-input-number
              v-model="comboSpecPrice"
              :precision="2"
              :step="0.01"
              :min="0"
              :max="999999"
              placeholder="0.00"
              controls-position="right"
              class="price-input"
              @change="handleDataChange"
            />
            <span class="unit">元</span>
          </div>
          <div class="form-help">套餐的基础价格，选项加价将在此基础上计算</div>
        </div>
      </div>
    </div>

    <!-- 套餐组设置 -->
    <div class="spec-content-card">
      <div class="card-header">
        <h3 class="card-title">
          <el-icon><SetUp /></el-icon>
          套餐组设置
        </h3>
        <el-button type="primary" size="small" @click="addComboGroup">
          <el-icon><Plus /></el-icon>
          添加套餐组
        </el-button>
      </div>
      
      <div v-if="comboGroups.length > 0" class="combo-groups-container">
        <div 
          v-for="(group, groupIndex) in comboGroups" 
          :key="group.id"
          class="combo-group-card"
        >
          <div class="group-header">
            <div class="group-title">
              <el-icon><Box /></el-icon>
              <span>{{ group.name || `套餐组 ${groupIndex + 1}` }}</span>
              <el-tag v-if="group.options.length > 0" size="small" type="success">
                {{ group.options.length }} 个选项
              </el-tag>
            </div>
            <el-button
              type="danger"
              size="small"
              plain
              @click="removeComboGroup(groupIndex)"
              :disabled="comboGroups.length === 1"
            >
              <el-icon><Delete /></el-icon>
              删除
            </el-button>
          </div>
          
          <div class="group-content">
            <!-- 基本信息 -->
            <div class="group-basic-info">
              <el-row :gutter="16">
                <el-col :span="12">
                  <div class="form-group">
                    <label class="form-label required">组名称</label>
                    <el-input
                      v-model="group.name"
                      placeholder="如：主食选择、饮料选择"
                      @blur="validateGroupName(group)"
                    />
                  </div>
                </el-col>
                <el-col :span="12">
                  <div class="form-group">
                    <label class="form-label">组描述</label>
                    <el-input
                      v-model="group.description"
                      placeholder="选择说明（可选）"
                    />
                  </div>
                </el-col>
              </el-row>
              
              <el-row :gutter="16">
                <el-col :span="8">
                  <div class="form-group">
                    <label class="form-label required">最少选择</label>
                    <el-input-number
                      v-model="group.minSelect"
                      :min="0"
                      :max="group.maxSelect || 99"
                      :precision="0"
                      controls-position="right"
                    />
                  </div>
                </el-col>
                <el-col :span="8">
                  <div class="form-group">
                    <label class="form-label required">最多选择</label>
                    <el-input-number
                      v-model="group.maxSelect"
                      :min="group.minSelect || 1"
                      :max="99"
                      :precision="0"
                      controls-position="right"
                    />
                  </div>
                </el-col>
                <el-col :span="8">
                  <div class="form-group">
                    <label class="form-label">是否必选</label>
                    <el-switch
                      v-model="group.isRequired"
                      active-text="必选"
                      inactive-text="可选"
                    />
                  </div>
                </el-col>
              </el-row>
            </div>

            <!-- 选项列表 -->
            <div class="group-options-section">
              <div class="section-header">
                <label class="form-label">套餐选项</label>
                <el-button size="small" type="success" plain @click="addComboOption(groupIndex)">
                  <el-icon><Plus /></el-icon>
                  添加选项
                </el-button>
              </div>
              
              <div v-if="group.options.length > 0" class="options-list">
                <div 
                  v-for="(option, optionIndex) in group.options"
                  :key="option.id"
                  class="option-card"
                >
                  <div class="option-header">
                    <div class="option-title">
                      <span>选项 {{ optionIndex + 1 }}</span>
                      <el-tag v-if="option.isDefault" size="small" type="warning">默认</el-tag>
                    </div>
                    <el-button
                      type="danger"
                      size="small"
                      plain
                      @click="removeComboOption(groupIndex, optionIndex)"
                    >
                      <el-icon><Delete /></el-icon>
                    </el-button>
                  </div>
                  
                  <div class="option-content">
                    <el-row :gutter="12">
                      <el-col :span="8">
                        <div class="form-group">
                          <label class="form-label required">选项名称</label>
                          <el-input
                            v-model="option.name"
                            placeholder="如：汉堡、薯条"
                          />
                        </div>
                      </el-col>
                      <el-col :span="6">
                        <div class="form-group">
                          <label class="form-label">加价金额</label>
                          <div class="input-with-unit">
                            <el-input-number
                              v-model="option.addPrice"
                              :precision="2"
                              :step="0.01"
                              :min="0"
                              placeholder="0.00"
                              controls-position="right"
                              size="small"
                            />
                            <span class="unit">元</span>
                          </div>
                        </div>
                      </el-col>
                      <el-col :span="4">
                        <div class="form-group">
                          <label class="form-label">数量</label>
                          <el-input-number
                            v-model="option.quantity"
                            :min="1"
                            :max="99"
                            :precision="0"
                            controls-position="right"
                            size="small"
                          />
                        </div>
                      </el-col>
                      <el-col :span="6">
                        <div class="form-group">
                          <label class="form-label">设置</label>
                          <div class="option-settings">
                            <el-checkbox v-model="option.isDefault" size="small">
                              默认选择
                            </el-checkbox>
                          </div>
                        </div>
                      </el-col>
                    </el-row>
                    
                    <el-row :gutter="12">
                      <el-col :span="12">
                        <div class="form-group">
                          <label class="form-label">关联商品</label>
                          <el-select
                            v-model="option.optionGoodsId"
                            placeholder="选择关联的商品（可选）"
                            clearable
                            filterable
                            size="small"
                            style="width: 100%"
                          >
                            <el-option
                              v-for="goods in availableGoods"
                              :key="goods.id"
                              :label="goods.name"
                              :value="goods.id"
                            />
                          </el-select>
                        </div>
                      </el-col>
                      <el-col :span="12">
                        <div class="form-group">
                          <label class="form-label">选项描述</label>
                          <el-input
                            v-model="option.description"
                            placeholder="选项说明（可选）"
                            size="small"
                          />
                        </div>
                      </el-col>
                    </el-row>
                  </div>
                </div>
              </div>
              
              <el-empty 
                v-else 
                description="暂无套餐选项，请点击添加选项"
                :image-size="80"
              />
            </div>
          </div>
        </div>
      </div>
      
      <el-empty 
        v-else 
        description="暂无套餐组，请点击右上角按钮添加"
        :image-size="100"
      />
    </div>

    <!-- 套餐预览 -->
    <div v-if="hasValidCombo" class="spec-content-card">
      <div class="card-header">
        <h3 class="card-title">
          <el-icon><View /></el-icon>
          套餐预览
        </h3>
        <div class="preview-actions">
          <el-button size="small" @click="calculatePreviewPrice" type="primary">
            <el-icon><Operation /></el-icon>
            计算价格
          </el-button>
        </div>
      </div>
      
      <div class="combo-preview">
        <div class="preview-header">
          <div class="base-price">
            <span class="label">基础价格：</span>
            <span class="price">¥{{ comboSpecPrice.toFixed(2) }}</span>
          </div>
          <div v-if="previewPrice > comboSpecPrice" class="add-price">
            <span class="label">选项加价：</span>
            <span class="price">+¥{{ (previewPrice - comboSpecPrice).toFixed(2) }}</span>
          </div>
          <div class="total-price">
            <span class="label">总价格：</span>
            <span class="price total">¥{{ previewPrice.toFixed(2) }}</span>
          </div>
        </div>
        
        <div class="preview-groups">
          <div 
            v-for="group in comboGroups"
            :key="group.id"
            class="preview-group"
          >
            <div class="group-info">
              <h4>{{ group.name }}</h4>
              <span class="group-rule">
                {{ group.isRequired ? '必选' : '可选' }} · 
                选择 {{ group.minSelect }}-{{ group.maxSelect }} 项
              </span>
            </div>
            <div class="group-options">
              <div 
                v-for="option in group.options"
                :key="option.id"
                class="preview-option"
                :class="{ default: option.isDefault }"
              >
                <span class="option-name">{{ option.name }}</span>
                <span v-if="option.addPrice > 0" class="option-price">+¥{{ option.addPrice.toFixed(2) }}</span>
                <span v-if="option.quantity > 1" class="option-quantity">×{{ option.quantity }}</span>
                <el-tag v-if="option.isDefault" size="small" type="warning">默认</el-tag>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, reactive, watch, nextTick } from 'vue'
import { Plus, Delete, SetUp, Box, Money, View, Operation } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'

const props = defineProps({
  modelValue: {
    type: Object,
    default: () => ({
      price: 0,
      specGroups: []
    })
  },
  availableGoods: {
    type: Array,
    default: () => []
  },
  formData: {
    type: Object,
    required: true
  }
})

const emit = defineEmits(['update:modelValue', 'change'])

// 套餐规格价格计算属性 - 直接操作allSpecData.comboSpecData.price
const comboSpecPrice = computed({
  get() {
    return props.formData.allSpecData?.comboSpecData?.price || 0
  },
  set(newPrice) {
    if (!props.formData.allSpecData) {
      props.formData.allSpecData = {}
    }
    if (!props.formData.allSpecData.comboSpecData) {
      props.formData.allSpecData.comboSpecData = { price: 0, comboGroups: [] }
    }
    props.formData.allSpecData.comboSpecData.price = newPrice
    handleDataChange()
  }
})

// 响应式数据
const formData = reactive({
  // 删除本地price字段，使用计算属性
})

const comboGroups = ref([])
const previewPrice = ref(0)

// 计算属性
const hasValidCombo = computed(() => {
  return comboGroups.value.some(group => 
    group.name && group.options && group.options.length > 0
  )
})

// ID生成器
let idCounter = 0
const generateId = () => {
  return Date.now() + (++idCounter)
}

// 标记是否正在更新，避免递归更新
const isUpdating = ref(false)

// 初始化时从 props 加载数据
const initializeFromProps = () => {
  console.log('=== ComboSpecSettings 初始化开始 ===')
  console.log('ComboSpecSettings - 从props初始化数据:', props.modelValue)
  console.log('ComboSpecSettings - props.modelValue.price:', props.modelValue.price)
  console.log('ComboSpecSettings - props.modelValue类型:', typeof props.modelValue.price)
  
  // 确保comboSpecData结构存在并设置价格
  if (!props.formData.allSpecData) {
    props.formData.allSpecData = {}
  }
  if (!props.formData.allSpecData.comboSpecData) {
    props.formData.allSpecData.comboSpecData = { price: props.modelValue.price || 0, comboGroups: [] }
  } else if (props.modelValue.price !== undefined) {
    props.formData.allSpecData.comboSpecData.price = props.modelValue.price
  }
  
  console.log('ComboSpecSettings - 设置基础价格:', comboSpecPrice.value)
  console.log('=== ComboSpecSettings 初始化结束 ===')
  
  // 强制触发页面更新
  setTimeout(() => {
    console.log('ComboSpecSettings - 延迟检查价格:', comboSpecPrice.value)
  }, 100)
  
  // 优先从specGroups加载数据，也支持从comboGroups加载作为兼容
  let sourceGroups = []
  if (props.modelValue.specGroups && props.modelValue.specGroups.length > 0) {
    sourceGroups = props.modelValue.specGroups
  } else if (props.modelValue.comboGroups && props.modelValue.comboGroups.length > 0) {
    sourceGroups = props.modelValue.comboGroups
  }
  
  if (sourceGroups.length > 0) {
    comboGroups.value = sourceGroups.map(group => ({
      id: group.id || generateId(),
      name: group.name || '',
      description: group.description || '',
      minSelect: group.minSelect || 1,
      maxSelect: group.maxSelect || 1,
      isRequired: group.isRequired !== false,
      options: (group.values || group.options || []).map(value => ({
        id: value.id || generateId(),
        name: value.name || '',
        value: value.value || value.name || '',
        addPrice: value.addPrice || 0,
        quantity: value.quantity || 1,
        isDefault: value.isDefault || false,
        optionGoodsId: value.optionGoodsId || null,
        description: value.description || ''
      }))
    }))
    console.log('ComboSpecSettings - 初始化套餐组:', comboGroups.value)
  }
  
  if (comboGroups.value.length === 0) {
    console.log('ComboSpecSettings - 没有套餐组，添加默认套餐组')
    addComboGroup()
  }
  
  calculatePreviewPrice()
}

// 组件挂载时初始化数据
initializeFromProps()

// 更新数据
const updateData = () => {
  if (isUpdating.value) return
  
  isUpdating.value = true
  
  const data = {
    price: comboSpecPrice.value,
    specGroups: comboGroups.value.map(group => ({
      id: group.id,
      name: group.name,
      description: group.description,
      minSelect: group.minSelect,
      maxSelect: group.maxSelect,
      isRequired: group.isRequired,
      sort: 0,
      values: group.options.map(option => ({
        id: option.id,
        name: option.name,
        value: option.value,
        addPrice: option.addPrice,
        quantity: option.quantity,
        isDefault: option.isDefault,
        optionGoodsId: option.optionGoodsId,
        description: option.description,
        sort: 0
      }))
    }))
  }
  
  emit('update:modelValue', data)
  emit('change', data)
  
  nextTick(() => {
    isUpdating.value = false
  })
}

// 处理数据变化
const handleDataChange = () => {
  updateData()
}

// 监听props变化，更新本地数据
watch(
  () => props.modelValue,
  (newValue) => {
    if (newValue && !isUpdating.value) {
      initializeFromProps()
    }
  },
  { deep: true }
)

// 监听数据变化，通知父组件
watch(
  [comboGroups],
  () => {
    if (!isUpdating.value) {
    updateData()
    calculatePreviewPrice()
    }
  },
  { deep: true }
)

// 添加套餐组
const addComboGroup = () => {
  const newGroup = {
    id: generateId(),
    name: '',
    description: '',
    minSelect: 1,
    maxSelect: 1,
    isRequired: true,
    options: []
  }
  comboGroups.value.push(newGroup)
}

// 删除套餐组
const removeComboGroup = (index) => {
  ElMessageBox.confirm('删除套餐组将清空相关选项数据，确定删除吗？', '确认删除', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    comboGroups.value.splice(index, 1)
  })
}

// 验证组名称
const validateGroupName = (group) => {
  if (!group.name.trim()) {
    ElMessage.warning('套餐组名称不能为空')
    return false
  }
  
  const duplicateGroup = comboGroups.value.find(g => 
    g.id !== group.id && g.name === group.name
  )
  if (duplicateGroup) {
    ElMessage.warning('套餐组名称不能重复')
    group.name = ''
    return false
  }
  
  return true
}

// 添加套餐选项
const addComboOption = (groupIndex) => {
  const group = comboGroups.value[groupIndex]
  const newOption = {
    id: generateId(),
    name: '',
    value: '',
    addPrice: 0,
    quantity: 1,
    isDefault: false,
    optionGoodsId: null,
    description: ''
  }
  
  group.options.push(newOption)
}

// 删除套餐选项
const removeComboOption = (groupIndex, optionIndex) => {
  const group = comboGroups.value[groupIndex]
  group.options.splice(optionIndex, 1)
}

// 计算预览价格
const calculatePreviewPrice = () => {
  let totalPrice = comboSpecPrice.value
  
  comboGroups.value.forEach(group => {
    const defaultOptions = group.options.filter(option => option.isDefault)
    defaultOptions.forEach(option => {
      totalPrice += option.addPrice * option.quantity
    })
  })
  
  previewPrice.value = totalPrice
}

// 暴露方法
defineExpose({
  addComboGroup,
  calculatePreviewPrice
})
</script>

<style lang="scss" scoped>
.combo-spec-container {
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
      }

      .preview-actions {
        display: flex;
        gap: 8px;
      }
    }
  }
  
  .base-settings {
    padding: 24px;
    
    .form-group {
      max-width: 400px;
      
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
      
      .input-with-unit {
        display: flex;
        align-items: center;
        
        .unit {
          margin-left: 8px;
          color: #909399;
          font-size: 14px;
          min-width: 20px;
        }
      }
      
      .form-help {
        margin-top: 4px;
        font-size: 12px;
        color: #909399;
        line-height: 1.4;
      }
      
      :deep(.el-input-number) {
        width: 100%;
      }
    }
  }
  
  .combo-groups-container {
    padding: 24px;
    
    .combo-group-card {
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
        
        .group-basic-info {
          margin-bottom: 24px;
          
          .form-group {
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
          }
        }
        
        .group-options-section {
          .section-header {
            display: flex;
            align-items: center;
            justify-content: space-between;
            margin-bottom: 16px;
            
            .form-label {
              margin: 0;
              font-weight: 500;
              color: #606266;
              font-size: 14px;
            }
          }
          
          .options-list {
            .option-card {
              background: #fff;
              border: 1px solid #e4e7ed;
              border-radius: 6px;
              margin-bottom: 12px;
              
              .option-header {
                display: flex;
                align-items: center;
                justify-content: space-between;
                padding: 12px 16px;
                background: #f8f9fa;
                border-bottom: 1px solid #e4e7ed;
                
                .option-title {
                  display: flex;
                  align-items: center;
                  gap: 8px;
                  font-weight: 500;
                  color: #303133;
                }
              }
              
              .option-content {
                padding: 16px;
                
                .form-group {
                  margin-bottom: 12px;
                  
                  &:last-child {
                    margin-bottom: 0;
                  }
                  
                  .form-label {
                    display: block;
                    margin-bottom: 6px;
                    font-weight: 500;
                    color: #606266;
                    font-size: 13px;
                    
                    &.required::after {
                      content: ' *';
                      color: #f56c6c;
                    }
                  }
                  
                  .input-with-unit {
                    display: flex;
                    align-items: center;
                    
                    .unit {
                      margin-left: 8px;
                      color: #909399;
                      font-size: 12px;
                      min-width: 16px;
                    }
                  }
                  
                  .option-settings {
                    display: flex;
                    align-items: center;
                    gap: 12px;
                    margin-top: 4px;
                  }
                }
              }
            }
          }
        }
      }
    }
  }
  
  .combo-preview {
    padding: 24px;
    
    .preview-header {
      display: flex;
      align-items: center;
      gap: 24px;
      padding: 16px;
      background: #f8f9fa;
      border-radius: 8px;
      margin-bottom: 20px;
      
      .base-price,
      .add-price,
      .total-price {
        display: flex;
        align-items: center;
        gap: 4px;
        
        .label {
          font-size: 14px;
          color: #606266;
        }
        
        .price {
          font-weight: 600;
          font-size: 16px;
          
          &.total {
            color: #e6a23c;
            font-size: 18px;
          }
        }
      }
      
      .add-price .price {
        color: #67c23a;
      }
    }
    
    .preview-groups {
      .preview-group {
        margin-bottom: 20px;
        padding: 16px;
        border: 1px solid #e4e7ed;
        border-radius: 8px;
        
        .group-info {
          margin-bottom: 12px;
          
          h4 {
            margin: 0 0 4px 0;
            font-size: 16px;
            font-weight: 600;
            color: #303133;
          }
          
          .group-rule {
            font-size: 12px;
            color: #909399;
          }
        }
        
        .group-options {
          display: flex;
          flex-wrap: wrap;
          gap: 8px;
          
          .preview-option {
            display: flex;
            align-items: center;
            gap: 6px;
            padding: 8px 12px;
            background: #f0f2f5;
            border-radius: 4px;
            font-size: 13px;
            
            &.default {
              background: #fdf6ec;
              border: 1px solid #f5dab1;
            }
            
            .option-name {
              color: #303133;
              font-weight: 500;
            }
            
            .option-price {
              color: #67c23a;
              font-weight: 600;
            }
            
            .option-quantity {
              color: #909399;
              font-size: 12px;
            }
          }
        }
      }
    }
  }
}

:deep(.el-empty) {
  padding: 40px 0;
}
</style> 