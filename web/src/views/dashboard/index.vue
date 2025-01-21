<template>
  <div class="dashboard">
    <el-row :gutter="20">
      <el-col :span="6">
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <span>用户总数</span>
            </div>
          </template>
          <div class="card-value">1,234</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <span>角色总数</span>
            </div>
          </template>
          <div class="card-value">5</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <span>菜单总数</span>
            </div>
          </template>
          <div class="card-value">24</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <span>在线用户</span>
            </div>
          </template>
          <div class="card-value">12</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <span>表单总数</span>
            </div>
          </template>
          <div class="card-value">{{ stats.formCount }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <span>提交总数</span>
            </div>
          </template>
          <div class="card-value">{{ stats.submitCount }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <span>今日提交</span>
            </div>
          </template>
          <div class="card-value">{{ stats.todaySubmit }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <span>模板数量</span>
            </div>
          </template>
          <div class="card-value">{{ stats.templateCount }}</div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { formatTime } from '@/utils/date'

const stats = ref({
  formCount: 0,
  submitCount: 0,
  todaySubmit: 0,
  templateCount: 0
})

// 获取统计数据
const getStats = () => {
  // 从 localStorage 获取表单数据
  const forms = JSON.parse(localStorage.getItem('forms') || '[]')
  const templates = JSON.parse(localStorage.getItem('formTemplates') || '[]')
  const submissions = JSON.parse(localStorage.getItem('formSubmissions') || '{}')
  
  // 计算统计数据
  stats.value.formCount = forms.length
  stats.value.templateCount = templates.length
  
  // 计算提交总数和今日提交数
  let totalSubmit = 0
  let todaySubmit = 0
  const today = new Date().toDateString()
  
  Object.values(submissions).forEach(formSubmits => {
    totalSubmit += formSubmits.length
    todaySubmit += formSubmits.filter(submit => 
      new Date(submit.submitTime).toDateString() === today
    ).length
  })
  
  stats.value.submitCount = totalSubmit
  stats.value.todaySubmit = todaySubmit
}

onMounted(() => {
  getStats()
})
</script>

<style scoped>
.dashboard {
  padding: 20px;
}

.card-header {
  font-size: 16px;
  font-weight: bold;
}

.card-value {
  font-size: 24px;
  font-weight: bold;
  color: #409eff;
  text-align: center;
  padding: 20px 0;
}

.el-card {
  transition: all 0.3s;
}

.el-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}
</style> 