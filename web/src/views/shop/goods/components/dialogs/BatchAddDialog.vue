<template>
  <el-dialog v-model="visible" title="批量添加规格值" width="500px" @close="handleClose">
    <div class="batch-add-content">
      <el-alert
        title="批量添加说明"
        type="info"
        :closable="false"
        show-icon
      >
        <template #default>
          <p>• 每行一个规格值，或用逗号、分号分隔</p>
          <p>• 系统会自动去重和去除空值</p>
          <p>• 支持从Excel复制粘贴</p>
        </template>
      </el-alert>
      
      <el-form :model="batchAddForm" label-width="80px" style="margin-top: 20px;">
        <el-form-item label="规格值">
          <el-input
            v-model="batchAddForm.values"
            type="textarea"
            :rows="8"
            placeholder="请输入规格值，支持以下格式：&#10;红色&#10;蓝色&#10;绿色&#10;&#10;或者：红色,蓝色,绿色&#10;或者：红色;蓝色;绿色"
            class="batch-textarea"
          />
        </el-form-item>
        
        <el-form-item label="预览">
          <div class="preview-container">
            <el-tag
              v-for="(value, index) in parsedValues"
              :key="index"
              size="small"
              class="preview-tag"
              :type="getValueTagType(index)"
            >
              {{ value }}
            </el-tag>
            <span v-if="parsedValues.length === 0" class="empty-preview">
              暂无有效规格值
            </span>
          </div>
        </el-form-item>
      </el-form>
    </div>
    
    <template #footer>
      <el-button @click="handleClose">取消</el-button>
      <el-button 
        type="primary" 
        @click="confirmBatchAdd"
        :disabled="parsedValues.length === 0"
      >
        添加 {{ parsedValues.length }} 个规格值
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { reactive, computed, watch } from 'vue'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:modelValue', 'confirm'])

const visible = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

// 批量添加表单
const batchAddForm = reactive({
  values: ''
})

// 解析批量添加的规格值
const parsedValues = computed(() => {
  if (!batchAddForm.values.trim()) return []
  
  // 支持多种分隔符：换行、逗号、分号
  const values = batchAddForm.values
    .split(/[\n,;，；]/)
    .map(v => v.trim())
    .filter(v => v.length > 0)
  
  // 去重
  return [...new Set(values)]
})

// 监听对话框打开，重置表单
watch(visible, (newValue) => {
  if (newValue) {
    resetForm()
  }
})

// 重置表单
const resetForm = () => {
  batchAddForm.values = ''
}

// 获取规格值标签类型
const getValueTagType = (index) => {
  const types = ['primary', 'success', 'info', 'warning', 'danger']
  return types[index % types.length]
}

// 确认批量添加
const confirmBatchAdd = () => {
  if (parsedValues.value.length === 0) return
  
  emit('confirm', parsedValues.value)
  handleClose()
}

// 关闭对话框
const handleClose = () => {
  visible.value = false
}
</script>

<style lang="scss" scoped>
.batch-add-content {
  .batch-textarea {
    :deep(.el-textarea__inner) {
      font-family: 'Courier New', monospace;
      line-height: 1.6;
    }
  }
  
  .preview-container {
    max-height: 120px;
    overflow-y: auto;
    padding: 12px;
    background: #f8f9fa;
    border: 1px solid #e4e7ed;
    border-radius: 6px;
    
    .preview-tag {
      margin: 2px 4px 2px 0;
    }
    
    .empty-preview {
      color: #909399;
      font-size: 14px;
      font-style: italic;
    }
  }
}
</style> 