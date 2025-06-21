<template>
  <el-dialog v-model="visible" title="批量编辑SKU" width="700px" @close="handleClose">
    <div class="batch-editor-content">
      <el-alert
        title="批量编辑说明"
        type="info"
        :closable="false"
        show-icon
      >
        <template #default>
          <p>• 默认已选中所有字段，只有勾选的字段才会被修改，未勾选的字段保持原值不变</p>
          <p>• 可以取消不需要修改的字段，同时修改多个字段，提高编辑效率</p>
        </template>
      </el-alert>
      
      <el-form :model="batchEditForm" label-width="120px" style="margin-top: 20px;">
        <!-- 选择应用范围 -->
        <el-form-item label="应用范围">
          <el-radio-group v-model="batchEditForm.applyType">
            <el-radio label="all">全部SKU ({{ skuList.length }}个)</el-radio>
            <el-radio label="selected" :disabled="selectedSkus.length === 0">
              已选SKU ({{ selectedSkus.length }}个)
            </el-radio>
          </el-radio-group>
        </el-form-item>
        
        <el-divider content-position="left">
          <span>价格设置</span>
          <el-button-group style="margin-left: 16px;">
            <el-button size="small" @click="selectAllFields">全选</el-button>
            <el-button size="small" @click="unselectAllFields">全不选</el-button>
          </el-button-group>
        </el-divider>
        
        <!-- 销售价格 -->
        <el-form-item>
          <template #label>
            <el-checkbox v-model="batchEditForm.fields.price.enabled">销售价格</el-checkbox>
          </template>
          <el-input-number
            v-model="batchEditForm.fields.price.value"
            :disabled="!batchEditForm.fields.price.enabled"
            :min="0"
            :max="999999"
            :precision="2"
            :step="0.01"
            style="width: 200px"
            controls-position="right"
            placeholder="0.00"
          />
          <span class="field-unit">元</span>
        </el-form-item>
        
        <!-- 市场价格 -->
        <el-form-item>
          <template #label>
            <el-checkbox v-model="batchEditForm.fields.originalPrice.enabled">市场价格</el-checkbox>
          </template>
          <el-input-number
            v-model="batchEditForm.fields.originalPrice.value"
            :disabled="!batchEditForm.fields.originalPrice.enabled"
            :min="0"
            :max="999999"
            :precision="2"
            :step="0.01"
            style="width: 200px"
            controls-position="right"
            placeholder="0.00"
          />
          <span class="field-unit">元</span>
        </el-form-item>
        
        <!-- 成本价格 -->
        <el-form-item>
          <template #label>
            <el-checkbox v-model="batchEditForm.fields.costPrice.enabled">成本价格</el-checkbox>
          </template>
          <el-input-number
            v-model="batchEditForm.fields.costPrice.value"
            :disabled="!batchEditForm.fields.costPrice.enabled"
            :min="0"
            :max="999999"
            :precision="2"
            :step="0.01"
            style="width: 200px"
            controls-position="right"
            placeholder="0.00"
          />
          <span class="field-unit">元</span>
        </el-form-item>
        
        <el-divider content-position="left">库存与规格</el-divider>
        
        <!-- 库存数量 -->
        <el-form-item>
          <template #label>
            <el-checkbox v-model="batchEditForm.fields.stock.enabled">库存数量</el-checkbox>
          </template>
          <el-input-number
            v-model="batchEditForm.fields.stock.value"
            :disabled="!batchEditForm.fields.stock.enabled"
            :min="0"
            :max="999999"
            :precision="0"
            style="width: 200px"
            controls-position="right"
            placeholder="0"
          />
          <span class="field-unit">件</span>
        </el-form-item>
        
        <!-- 重量 -->
        <el-form-item>
          <template #label>
            <el-checkbox v-model="batchEditForm.fields.weight.enabled">商品重量</el-checkbox>
          </template>
          <el-input-number
            v-model="batchEditForm.fields.weight.value"
            :disabled="!batchEditForm.fields.weight.enabled"
            :min="0"
            :precision="3"
            :step="0.001"
            style="width: 200px"
            controls-position="right"
            placeholder="0.000"
          />
          <span class="field-unit">kg</span>
        </el-form-item>
        
        <!-- 体积 -->
        <el-form-item>
          <template #label>
            <el-checkbox v-model="batchEditForm.fields.volume.enabled">商品体积</el-checkbox>
          </template>
          <el-input-number
            v-model="batchEditForm.fields.volume.value"
            :disabled="!batchEditForm.fields.volume.enabled"
            :min="0"
            :precision="3"
            :step="0.001"
            style="width: 200px"
            controls-position="right"
            placeholder="0.000"
          />
          <span class="field-unit">m³</span>
        </el-form-item>
        
        <!-- 预览受影响的SKU -->
        <el-form-item label="预览">
          <div class="affected-skus-preview">
            <div class="preview-header">
              <span>将影响以下SKU：</span>
              <el-tag size="small" type="info">
                {{ getAffectedSkus().length }} 个
              </el-tag>
              <span class="enabled-fields">
                修改字段：{{ getEnabledFieldsText() }}
              </span>
            </div>
            <div class="preview-list">
              <el-tag
                v-for="(sku, index) in getAffectedSkus().slice(0, 5)"
                :key="index"
                size="small"
                class="sku-preview-tag"
              >
                {{ getSkuDisplayName(sku) }}
              </el-tag>
              <span v-if="getAffectedSkus().length > 5" class="more-indicator">
                ...还有{{ getAffectedSkus().length - 5 }}个
              </span>
            </div>
          </div>
        </el-form-item>
      </el-form>
    </div>
    
    <template #footer>
      <el-button @click="handleClose">取消</el-button>
      <el-button 
        type="primary" 
        @click="confirmBatchEdit"
        :disabled="getAffectedSkus().length === 0 || !hasEnabledFields()"
      >
        确定修改 ({{ getAffectedSkus().length }}个)
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { reactive, computed, watch } from 'vue'
import { ElMessage } from 'element-plus'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  skuList: {
    type: Array,
    default: () => []
  },
  selectedSkus: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['update:modelValue', 'confirm'])

const visible = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

// 批量编辑表单
const batchEditForm = reactive({
  applyType: 'all',
  fields: {
    price: { enabled: true, value: 0 },        // 默认选中
    originalPrice: { enabled: true, value: 0 }, // 默认选中
    costPrice: { enabled: true, value: 0 },     // 默认选中
    stock: { enabled: true, value: 0 },         // 默认选中
    weight: { enabled: true, value: 0 },        // 默认选中
    volume: { enabled: true, value: 0 }         // 默认选中
  }
})

// 监听对话框打开，重置表单
watch(visible, (newValue) => {
  if (newValue) {
    resetForm()
  }
})

// 重置表单
const resetForm = () => {
  batchEditForm.applyType = 'all'
  Object.keys(batchEditForm.fields).forEach(key => {
    batchEditForm.fields[key].enabled = true  // 默认选中所有字段
    batchEditForm.fields[key].value = 0
  })
}

// 全选字段
const selectAllFields = () => {
  Object.keys(batchEditForm.fields).forEach(key => {
    batchEditForm.fields[key].enabled = true
  })
}

// 全不选字段
const unselectAllFields = () => {
  Object.keys(batchEditForm.fields).forEach(key => {
    batchEditForm.fields[key].enabled = false
  })
}

// 获取受影响的SKU
const getAffectedSkus = () => {
  return batchEditForm.applyType === 'all' ? props.skuList : props.selectedSkus
}

// 检查是否有启用的字段
const hasEnabledFields = () => {
  return Object.values(batchEditForm.fields).some(field => field.enabled)
}

// 获取启用字段的文本描述
const getEnabledFieldsText = () => {
  const fieldLabels = {
    price: '销售价格',
    originalPrice: '市场价格',
    costPrice: '成本价格',
    stock: '库存数量',
    weight: '重量',
    volume: '体积'
  }
  
  const enabledFields = Object.keys(batchEditForm.fields)
    .filter(key => batchEditForm.fields[key].enabled)
    .map(key => fieldLabels[key])
  
  return enabledFields.length > 0 ? enabledFields.join('、') : '无'
}

// 获取SKU显示名称
const getSkuDisplayName = (sku) => {
  return Object.values(sku.specs || {}).map(spec => spec.name || spec.value).join(' / ')
}

// 确认批量编辑
const confirmBatchEdit = () => {
  const affectedSkus = getAffectedSkus()
  if (affectedSkus.length === 0) {
    ElMessage.warning('没有选中的SKU')
    return
  }
  
  if (!hasEnabledFields()) {
    ElMessage.warning('请至少选择一个要修改的字段')
    return
  }
  
  const enabledFields = Object.keys(batchEditForm.fields).filter(key => batchEditForm.fields[key].enabled)
  
  const editData = {
    applyType: batchEditForm.applyType,
    fields: enabledFields.reduce((acc, key) => {
      acc[key] = batchEditForm.fields[key].value
      return acc
    }, {}),
    affectedSkus
  }
  
  emit('confirm', editData)
  handleClose()
  
  const fieldsText = getEnabledFieldsText()
  ElMessage.success(`成功修改 ${affectedSkus.length} 个SKU的${fieldsText}`)
}

// 关闭对话框
const handleClose = () => {
  visible.value = false
}
</script>

<style lang="scss" scoped>
.batch-editor-content {
  .field-unit {
    margin-left: 8px;
    color: #909399;
    font-size: 14px;
  }
  
  :deep(.el-form-item__label) {
    .el-checkbox {
      margin-right: 8px;
    }
  }
  
  .affected-skus-preview {
    .preview-header {
      display: flex;
      align-items: center;
      gap: 8px;
      margin-bottom: 8px;
      font-size: 14px;
      color: #606266;
      flex-wrap: wrap;
      
      .enabled-fields {
        color: #67c23a;
        font-weight: 500;
        margin-left: auto;
      }
    }
    
    .preview-list {
      display: flex;
      flex-wrap: wrap;
      gap: 6px;
      padding: 12px;
      background: #f8f9fa;
      border: 1px solid #e4e7ed;
      border-radius: 6px;
      min-height: 40px;
      align-items: center;
      
      .sku-preview-tag {
        font-size: 12px;
      }
      
      .more-indicator {
        color: #909399;
        font-size: 12px;
        font-style: italic;
      }
    }
  }
}
</style> 