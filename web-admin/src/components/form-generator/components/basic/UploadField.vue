<template>
  <el-upload
    v-model:file-list="fileList"
    v-bind="$attrs"
    :on-change="handleChange"
    :http-request="customUpload"
    :on-exceed="handleExceed"
  >
    <el-icon><Plus /></el-icon>
    <template #file="{ file }">
      <div class="preview-item">
        <img :src="file.url" class="preview-image" />
        <span class="preview-actions">
          <el-icon color="#fff" @click.stop="handlePreview(file)"><ZoomIn /></el-icon>
          <el-icon color="#fff" @click.stop="handleRemove(file)"><Delete /></el-icon>
        </span>
      </div>
    </template>
    <template #tip>
      <div class="el-upload__tip">
        支持多选，每张图片不超过2MB
      </div>
    </template>
  </el-upload>

  <!-- 图片预览对话框 -->
  <el-dialog v-model="previewVisible" width="800px">
    <img :src="previewUrl" style="width: 100%" />
  </el-dialog>
</template>

<script setup>
import { ref, watch } from 'vue'
import { Plus, Delete, ZoomIn } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

const props = defineProps({
  modelValue: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['update:modelValue'])

const fileList = ref([])
const previewVisible = ref(false)
const previewUrl = ref('')

// 监听外部值变化
watch(() => props.modelValue, (val) => {
  fileList.value = val
}, { deep: true })

// 处理文件变化
const handleChange = (uploadFile) => {
  if (uploadFile.status === 'ready') {
    const reader = new FileReader()
    reader.onload = (e) => {
      uploadFile.url = e.target.result
      emit('update:modelValue', fileList.value)
    }
    reader.readAsDataURL(uploadFile.raw)
  }
}

// 自定义上传
const customUpload = () => {
  // 不实际上传，只做本地预览
  return false
}

// 移除文件
const handleRemove = (file) => {
  fileList.value = fileList.value.filter(item => item.uid !== file.uid)
  emit('update:modelValue', fileList.value)
}

// 处理超出限制
const handleExceed = () => {
  ElMessage.warning(`最多只能上传 ${props.$attrs.limit} 张图片`)
}

// 预览图片
const handlePreview = (file) => {
  previewUrl.value = file.url
  previewVisible.value = true
}
</script>

<style scoped>
.preview-item {
  position: relative;
  width: 148px;
  height: 148px;
  border: 1px solid #d9d9d9;
  border-radius: 6px;
  overflow: hidden;
}

.preview-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.preview-actions {
  position: absolute;
  top: 0;
  right: 0;
  bottom: 0;
  left: 0;
  display: flex;
  justify-content: center;
  align-items: center;
  background: rgba(0, 0, 0, 0.5);
  opacity: 0;
  transition: opacity 0.3s;
  cursor: pointer;
}

.preview-item:hover .preview-actions {
  opacity: 1;
}

.preview-actions .el-icon {
  font-size: 20px;
  margin: 0 6px;
  cursor: pointer;
}

.preview-actions .el-icon:hover {
  transform: scale(1.2);
}
</style> 