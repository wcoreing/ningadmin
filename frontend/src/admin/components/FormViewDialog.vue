<template>
  <el-dialog
    v-model="dialogVisible"
    :title="title"
    width="800px"
    :close-on-click-modal="false"
    :close-on-press-escape="false"
    @close="handleClose"
  >
    <form-view
      :form-id="formId"
      :is-dialog="true"
      @submit-success="handleSubmitSuccess"
    />
  </el-dialog>
</template>

<script setup>
import { ref, watch } from 'vue'
import FormView from '@/admin/views/form/view/index.vue'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  formId: {
    type: [String, Number],
    default: ''
  },
  title: {
    type: String,
    default: '查看表单'
  }
})

const emit = defineEmits(['update:modelValue', 'close', 'submit-success'])

const dialogVisible = ref(false)

watch(() => props.modelValue, (val) => {
  dialogVisible.value = val
})

watch(() => dialogVisible.value, (val) => {
  emit('update:modelValue', val)
})

const handleClose = () => {
  emit('close')
}

const handleSubmitSuccess = () => {
  emit('submit-success')
  dialogVisible.value = false
}
</script>

<style scoped>
:deep(.el-dialog__body) {
  padding: 0;
}
</style> 