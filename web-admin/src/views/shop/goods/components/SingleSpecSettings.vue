<template>
  <div class="single-spec-container">
    <div class="spec-content-card">
      <div class="card-header">
        <h3 class="card-title">
          <el-icon><Money /></el-icon>
          价格库存设置
        </h3>
      </div>
      
      <div class="form-grid">
        <div class="form-group">
          <label class="form-label required">销售价格</label>
          <div class="input-with-unit">
            <el-input-number
              v-model="singleSpecPrice"
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
        </div>
        
        <div class="form-group">
          <label class="form-label">市场价格</label>
          <div class="input-with-unit">
            <el-input-number
              v-model="singleSpecOriginalPrice"
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
        </div>
        
        <div class="form-group">
          <label class="form-label">成本价格</label>
          <div class="input-with-unit">
            <el-input-number
              v-model="singleSpecCostPrice"
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
        </div>
        
        <div class="form-group">
          <label class="form-label required">商品库存</label>
          <div class="input-with-unit">
            <el-input-number
              v-model="singleSpecStock"
              :min="0"
              :max="999999"
              :precision="0"
              placeholder="0"
              controls-position="right"
              class="stock-input"
              @change="handleDataChange"
            />
            <span class="unit">件</span>
          </div>
        </div>
        
        <div class="form-group">
          <label class="form-label">商品重量</label>
          <div class="input-with-unit">
            <el-input-number
              v-model="singleSpecWeight"
              :precision="3"
              :step="0.001"
              :min="0"
              placeholder="0.000"
              controls-position="right"
              class="weight-input"
              @change="handleDataChange"
            />
            <span class="unit">kg</span>
          </div>
        </div>
        
        <div class="form-group">
          <label class="form-label">商品体积</label>
          <div class="input-with-unit">
            <el-input-number
              v-model="singleSpecVolume"
              :precision="3"
              :step="0.001"
              :min="0"
              placeholder="0.000"
              controls-position="right"
              class="volume-input"
              @change="handleDataChange"
            />
            <span class="unit">m³</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, nextTick, computed } from 'vue'
import { Money } from '@element-plus/icons-vue'

const props = defineProps({
  modelValue: {
    type: Object,
    default: () => ({
      price: 0,
      originalPrice: 0,
      costPrice: 0,
      stock: 0,
      weight: 0,
      volume: 0
    })
  },
  formData: {
    type: Object,
    required: true
  }
})

const emit = defineEmits(['update:modelValue', 'change'])

// 单规格价格计算属性 - 直接操作allSpecData.singleSpecData.price
const singleSpecPrice = computed({
  get() {
    return props.formData.allSpecData?.singleSpecData?.price || 0
  },
  set(newPrice) {
    if (!props.formData.allSpecData) {
      props.formData.allSpecData = {}
    }
    if (!props.formData.allSpecData.singleSpecData) {
      props.formData.allSpecData.singleSpecData = { price: 0, originalPrice: 0, costPrice: 0, stock: 0, weight: 0, volume: 0 }
    }
    props.formData.allSpecData.singleSpecData.price = newPrice
    handleDataChange()
  }
})

// 计算属性：直接操作allSpecData.singleSpecData中的字段
const singleSpecOriginalPrice = computed({
  get() {
    return props.formData.allSpecData?.singleSpecData?.originalPrice || 0
  },
  set(newValue) {
    if (!props.formData.allSpecData) {
      props.formData.allSpecData = {}
    }
    if (!props.formData.allSpecData.singleSpecData) {
      props.formData.allSpecData.singleSpecData = { price: 0, originalPrice: 0, costPrice: 0, stock: 0, weight: 0, volume: 0 }
    }
    props.formData.allSpecData.singleSpecData.originalPrice = newValue
    handleDataChange()
  }
})

const singleSpecCostPrice = computed({
  get() {
    return props.formData.allSpecData?.singleSpecData?.costPrice || 0
  },
  set(newValue) {
    if (!props.formData.allSpecData.singleSpecData) {
      props.formData.allSpecData.singleSpecData = { price: 0, originalPrice: 0, costPrice: 0, stock: 0, weight: 0, volume: 0 }
    }
    props.formData.allSpecData.singleSpecData.costPrice = newValue
    handleDataChange()
  }
})

const singleSpecStock = computed({
  get() {
    return props.formData.allSpecData?.singleSpecData?.stock || 0
  },
  set(newValue) {
    if (!props.formData.allSpecData.singleSpecData) {
      props.formData.allSpecData.singleSpecData = { price: 0, originalPrice: 0, costPrice: 0, stock: 0, weight: 0, volume: 0 }
    }
    props.formData.allSpecData.singleSpecData.stock = newValue
    handleDataChange()
  }
})

const singleSpecWeight = computed({
  get() {
    return props.formData.allSpecData?.singleSpecData?.weight || 0
  },
  set(newValue) {
    if (!props.formData.allSpecData.singleSpecData) {
      props.formData.allSpecData.singleSpecData = { price: 0, originalPrice: 0, costPrice: 0, stock: 0, weight: 0, volume: 0 }
    }
    props.formData.allSpecData.singleSpecData.weight = newValue
    handleDataChange()
  }
})

const singleSpecVolume = computed({
  get() {
    return props.formData.allSpecData?.singleSpecData?.volume || 0
  },
  set(newValue) {
    if (!props.formData.allSpecData.singleSpecData) {
      props.formData.allSpecData.singleSpecData = { price: 0, originalPrice: 0, costPrice: 0, stock: 0, weight: 0, volume: 0 }
    }
    props.formData.allSpecData.singleSpecData.volume = newValue
    handleDataChange()
  }
})

// 标记是否正在更新，避免递归更新
const isUpdating = ref(false)

// 处理数据变化
const handleDataChange = () => {
  if (isUpdating.value) return
  
  isUpdating.value = true
  const updatedData = { 
    price: singleSpecPrice.value,
    originalPrice: singleSpecOriginalPrice.value,
    costPrice: singleSpecCostPrice.value,
    stock: singleSpecStock.value,
    weight: singleSpecWeight.value,
    volume: singleSpecVolume.value
  }
  emit('update:modelValue', updatedData)
  emit('change', updatedData)
  
  nextTick(() => {
    isUpdating.value = false
  })
}

// 不再需要监听逻辑，因为使用计算属性直接操作数据源
</script>

<style lang="scss" scoped>
.single-spec-container {
  .spec-content-card {
    background: #fff;
    border-radius: 12px;
    border: 1px solid #e4e7ed;
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
    }
  }
  
  .form-grid {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 24px;
    padding: 24px;
    
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
      
      :deep(.el-input-number) {
        width: 100%;
      }
    }
  }
}
</style> 