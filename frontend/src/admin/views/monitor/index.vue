<template>
  <div class="monitor">
    <el-row :gutter="20">
      <!-- 系统信息 -->
      <el-col :span="24">
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <span>系统信息</span>
            </div>
          </template>
          <el-descriptions :column="3" border>
            <el-descriptions-item label="操作系统">
              {{ systemInfo.platform }}
            </el-descriptions-item>
            <el-descriptions-item label="浏览器">
              {{ systemInfo.browser }}
            </el-descriptions-item>
            <el-descriptions-item label="屏幕分辨率">
              {{ systemInfo.screen }}
            </el-descriptions-item>
            <el-descriptions-item label="语言">
              {{ systemInfo.language }}
            </el-descriptions-item>
            <el-descriptions-item label="时区">
              {{ systemInfo.timezone }}
            </el-descriptions-item>
            <el-descriptions-item label="在线时长">
              {{ systemInfo.uptime }}
            </el-descriptions-item>
          </el-descriptions>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { Monitor } from '@element-plus/icons-vue'

// 系统信息
const systemInfo = ref({
  platform: `${navigator.platform} (${navigator.userAgent.match(/\([^)]+\)/)[0]})`,
  browser: navigator.userAgent.split(' ').pop(),
  screen: `${window.screen.width}x${window.screen.height}`,
  language: navigator.language,
  timezone: Intl.DateTimeFormat().resolvedOptions().timeZone,
  uptime: '0:00:00'
})

// 更新在线时长
const updateUptime = () => {
  const start = Date.now()
  setInterval(() => {
    const diff = Date.now() - start
    const hours = Math.floor(diff / 3600000)
    const minutes = Math.floor((diff % 3600000) / 60000)
    const seconds = Math.floor((diff % 60000) / 1000)
    systemInfo.value.uptime = `${hours}:${minutes.toString().padStart(2, '0')}:${seconds.toString().padStart(2, '0')}`
  }, 1000)
}

onMounted(() => {
  updateUptime()
})
</script>

<style scoped>
.monitor {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.el-card {
  margin-bottom: 20px;
}

:deep(.el-descriptions__label) {
  width: 120px;
}
</style> 