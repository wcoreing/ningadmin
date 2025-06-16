<template>
  <div class="app-container">
    <el-card class="file-manager-container">
      <file-manager
        v-model="visible"
        standalone
        :upload-url="uploadUrl"
        :multiple="true"
        :page-size="pageSize"
        @select="handleSelect"
        ref="fileManagerRef"
      />
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import FileManager from '@/components/FileManager/index.vue'

const visible = ref(true)
const uploadUrl = ref((import.meta.env.VITE_API_BASE_URL || '/api') + '/common/upload')
const pageSize = ref(55)
const fileManagerRef = ref(null)

const handleSelect = (files) => {
  console.log('选中的文件:', files)
}

onMounted(() => {
  console.log('File page mounted')
})
</script>

<style lang="scss" scoped>
.file-manager-container {
  min-height: calc(100vh - 140px);
  
  :deep(.el-card__body) {
    padding: 0;
  }
}
</style> 