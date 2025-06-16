<template>
  <div class="spec-settings-container">
    <!-- 规格类型选择卡片 -->
    <div v-if="showTypeSelector" class="spec-type-card">
      <div class="spec-type-selector">
        <div 
          class="spec-type-option"
          :class="{ active: formData.specType === 1 }"
          @click="switchSpecType(1)"
        >
          <div class="option-icon">
            <el-icon><Document /></el-icon>
          </div>
          <div class="option-content">
            <h4>单规格</h4>
            <p>商品只有一种规格，统一价格和库存</p>
          </div>
          <div class="option-check">
            <el-icon v-if="formData.specType === 1"><Check /></el-icon>
          </div>
        </div>
        
        <div 
          class="spec-type-option"
          :class="{ active: formData.specType === 2 }"
          @click="switchSpecType(2)"
        >
          <div class="option-icon">
            <el-icon><Grid /></el-icon>
          </div>
          <div class="option-content">
            <h4>多规格</h4>
            <p>商品有多种规格组合，如颜色、尺寸等</p>
          </div>
          <div class="option-check">
            <el-icon v-if="formData.specType === 2"><Check /></el-icon>
      </div>
    </div>

        <div 
          class="spec-type-option"
          :class="{ active: formData.specType === 3 }"
          @click="switchSpecType(3)"
        >
          <div class="option-icon">
            <el-icon><SetUp /></el-icon>
              </div>
          <div class="option-content">
            <h4>套餐规格</h4>
            <p>商品为套餐形式，包含多个组和选项</p>
            </div>
          <div class="option-check">
            <el-icon v-if="formData.specType === 3"><Check /></el-icon>
              </div>
              </div>
          </div>
        </div>
        
    <!-- 单规格设置 -->
    <SingleSpecSettings
      v-if="formData.specType === 1"
      :model-value="singleSpecData"
      :form-data="formData"
      @update:modelValue="updateSingleSpec"
    />
    
    <!-- 多规格设置 -->
    <MultiSpecSettings
      v-if="formData.specType === 2"
      :model-value="multiSpecData"
      @update:modelValue="updateMultiSpec"
    />
    
    <!-- 套餐规格设置 -->
    <ComboSpecSettings
      v-if="formData.specType === 3"
      :model-value="comboSpecData"
      :form-data="formData"
      :available-goods="availableGoods"
      @update:modelValue="updateComboSpec"
    />
        </div>
</template>

<script setup>
import { computed } from 'vue'
import { Document, Grid, Check, SetUp } from '@element-plus/icons-vue'
import SingleSpecSettings from './SingleSpecSettings.vue'
import MultiSpecSettings from './MultiSpecSettings.vue'
import ComboSpecSettings from './ComboSpecSettings.vue'
import { 
  extractSingleSpecData, 
  extractMultiSpecData, 
  extractComboSpecData,
  updateSingleSpecToForm,
  updateMultiSpecToForm,
  updateComboSpecToForm
} from '../utils/specUtils.js'

const props = defineProps({
  formData: {
    type: Object,
    required: true
  },
  uploadUrl: {
    type: String,
    required: true
  },
  showTypeSelector: {
    type: Boolean,
    default: true
  },
  availableGoods: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['spec-type-change'])

// 计算属性：使用工具函数提取数据
const singleSpecData = computed(() => extractSingleSpecData(props.formData))
const multiSpecData = computed(() => extractMultiSpecData(props.formData))
const comboSpecData = computed(() => extractComboSpecData(props.formData))

// 切换规格类型
const switchSpecType = (type) => {
  if (props.formData.specType === type) return
  
  props.formData.specType = type
  emit('spec-type-change', type)
}

// 更新数据：使用工具函数
const updateSingleSpec = (data) => updateSingleSpecToForm(props.formData, data)
const updateMultiSpec = (data) => updateMultiSpecToForm(props.formData, data)
const updateComboSpec = (data) => updateComboSpecToForm(props.formData, data)

// 获取完整数据
const getAllData = () => ({
  specType: props.formData.specType,
  ...extractSingleSpecData(props.formData),
  ...extractMultiSpecData(props.formData),
  ...extractComboSpecData(props.formData)
})

// 暴露方法
defineExpose({ getAllData })
</script>

<style lang="scss" scoped>
.spec-settings-container {
  padding: 0;
  
  .spec-type-card {
    background: #fff;
    border-radius: 12px;
    border: 1px solid #e4e7ed;
    margin-bottom: 20px;
    overflow: hidden;
    
    .spec-type-selector {
      padding: 24px;
    display: grid;
      grid-template-columns: 1fr 1fr 1fr;
      gap: 16px;
      
      .spec-type-option {
        display: flex;
        align-items: center;
        padding: 20px;
        border: 2px solid #e4e7ed;
        border-radius: 12px;
        cursor: pointer;
        transition: all 0.3s ease;
        background: #fafbfc;
        
        &:hover {
          border-color: #409eff;
          background: #f0f9ff;
          transform: translateY(-2px);
          box-shadow: 0 4px 12px rgba(64, 158, 255, 0.15);
        }
        
        &.active {
          border-color: #409eff;
          background: #f0f9ff;
          box-shadow: 0 0 0 3px rgba(64, 158, 255, 0.1);
        }
        
        .option-icon {
          margin-right: 16px;
          font-size: 24px;
          color: #409eff;
        }
        
        .option-content {
          flex: 1;
          
          h4 {
            margin: 0 0 4px 0;
            font-size: 16px;
            font-weight: 600;
            color: #303133;
          }
          
          p {
            margin: 0;
            font-size: 14px;
            color: #606266;
            line-height: 1.4;
          }
        }
        
        .option-check {
          font-size: 20px;
          color: #67c23a;
        }
      }
    }
  }
}
</style> 