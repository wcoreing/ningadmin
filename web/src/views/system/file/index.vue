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
      />
    </el-card>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import FileManager from '@/components/FileManager/index.vue'

const visible = ref(true)
const uploadUrl = ref(import.meta.env.VITE_API_URL + '/common/upload')
const pageSize = ref(55)

const handleSelect = (files) => {
  console.log('选中的文件:', files)
}
</script>

<style lang="scss" scoped>
.file-manager-container {
  min-height: calc(100vh - 140px);
  
  :deep(.el-card__body) {
    padding: 0;
  }
}
</style> 