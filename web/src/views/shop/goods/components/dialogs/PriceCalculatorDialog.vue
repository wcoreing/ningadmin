<template>
  <el-dialog v-model="visible" title="价格计算器" width="600px" @close="handleClose">
    <div class="calculator-content">
      <el-alert
        title="价格计算器"
        type="success"
        :closable="false"
        show-icon
      >
        <template #default>
          <p>根据成本价自动计算销售价格和市场价格</p>
        </template>
      </el-alert>
      
      <el-form :model="calculatorForm" label-width="120px" style="margin-top: 20px;">
        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="成本价基准">
              <el-radio-group v-model="calculatorForm.costBase">
                <el-radio value="manual">手动输入</el-radio>
                <el-radio value="current">使用当前成本价</el-radio>
              </el-radio-group>
            </el-form-item>
            
            <el-form-item v-if="calculatorForm.costBase === 'manual'" label="成本价格">
              <el-input-number
                v-model="calculatorForm.baseCost"
                :precision="2"
                :step="0.01"
                :min="0"
                style="width: 100%"
                placeholder="0.00"
              />
              <span class="field-unit">元</span>
            </el-form-item>
          </el-col>
          
          <el-col :span="12">
            <el-form-item label="销售价利润率">
              <el-input-number
                v-model="calculatorForm.profitRate"
                :precision="1"
                :step="1"
                :min="0"
                :max="1000"
                style="width: 100%"
                placeholder="30.0"
              />
              <span class="field-unit">%</span>
            </el-form-item>
            
            <el-form-item label="市场价利润率">
              <el-input-number
                v-model="calculatorForm.marketRate"
                :precision="1"
                :step="1"
                :min="0"
                :max="1000"
                style="width: 100%"
                placeholder="50.0"
              />
              <span class="field-unit">%</span>
            </el-form-item>
          </el-col>
        </el-row>
        
        <el-form-item label="应用范围">
          <el-radio-group v-model="calculatorForm.applyType">
            <el-radio value="all">全部SKU ({{ skuList.length }}个)</el-radio>
            <el-radio value="selected" :disabled="selectedSkus.length === 0">
              已选SKU ({{ selectedSkus.length }}个)
            </el-radio>
          </el-radio-group>
        </el-form-item>
        
        <!-- 计算预览 -->
        <el-form-item label="计算预览">
          <div class="calculation-preview">
            <div class="preview-item">
              <span class="label">成本价：</span>
              <span class="value">{{ getPreviewCost() }} 元</span>
            </div>
            <div class="preview-item">
              <span class="label">销售价：</span>
              <span class="value">{{ getPreviewPrice() }} 元</span>
              <span class="rate">(+{{ calculatorForm.profitRate }}%)</span>
            </div>
            <div class="preview-item">
              <span class="label">市场价：</span>
              <span class="value">{{ getPreviewMarketPrice() }} 元</span>
              <span class="rate">(+{{ calculatorForm.marketRate }}%)</span>
            </div>
          </div>
        </el-form-item>
      </el-form>
    </div>
    
    <template #footer>
      <el-button @click="handleClose">取消</el-button>
      <el-button type="primary" @click="applyCalculation">
        应用计算结果
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

// 价格计算器表单
const calculatorForm = reactive({
  costBase: 'current',
  baseCost: 0,
  profitRate: 30,
  marketRate: 50,
  applyType: 'all'
})

// 监听对话框打开，重置表单
watch(visible, (newValue) => {
  if (newValue) {
    resetForm()
  }
})

// 重置表单
const resetForm = () => {
  calculatorForm.costBase = 'current'
  calculatorForm.baseCost = 0
  calculatorForm.profitRate = 30
  calculatorForm.marketRate = 50
  calculatorForm.applyType = 'all'
}

// 计算器预览
const getPreviewCost = () => {
  if (calculatorForm.costBase === 'manual') {
    return calculatorForm.baseCost.toFixed(2)
  }
  // 使用当前成本价的平均值作为预览
  if (props.skuList.length === 0) return '0.00'
  const avgCost = props.skuList.reduce((sum, sku) => sum + (sku.costPrice || 0), 0) / props.skuList.length
  return avgCost.toFixed(2)
}

const getPreviewPrice = () => {
  const cost = parseFloat(getPreviewCost())
  const price = cost * (1 + calculatorForm.profitRate / 100)
  return price.toFixed(2)
}

const getPreviewMarketPrice = () => {
  const cost = parseFloat(getPreviewCost())
  const marketPrice = cost * (1 + calculatorForm.marketRate / 100)
  return marketPrice.toFixed(2)
}

// 应用计算结果
const applyCalculation = () => {
  const affectedSkus = calculatorForm.applyType === 'all' ? props.skuList : props.selectedSkus
  
  if (affectedSkus.length === 0) {
    ElMessage.warning('没有选中的SKU')
    return
  }
  
  const calculationData = {
    costBase: calculatorForm.costBase,
    baseCost: calculatorForm.baseCost,
    profitRate: calculatorForm.profitRate,
    marketRate: calculatorForm.marketRate,
    applyType: calculatorForm.applyType,
    affectedSkus
  }
  
  emit('confirm', calculationData)
  handleClose()
  ElMessage.success(`成功计算 ${affectedSkus.length} 个SKU的价格`)
}

// 关闭对话框
const handleClose = () => {
  visible.value = false
}
</script>

<style lang="scss" scoped>
.calculator-content {
  .field-unit {
    margin-left: 8px;
    color: #909399;
    font-size: 14px;
  }
  
  .calculation-preview {
    padding: 16px;
    background: #f8f9fa;
    border: 1px solid #e4e7ed;
    border-radius: 6px;
    
    .preview-item {
      display: flex;
      align-items: center;
      margin-bottom: 8px;
      
      &:last-child {
        margin-bottom: 0;
      }
      
      .label {
        font-weight: 500;
        color: #606266;
        min-width: 80px;
      }
      
      .value {
        font-weight: 600;
        color: #303133;
        margin-right: 8px;
      }
      
      .rate {
        font-size: 12px;
        color: #67c23a;
      }
    }
  }
}
</style> 