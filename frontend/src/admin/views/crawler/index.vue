<template>
  <div class="app-container">
    <div class="filter-container">
      <el-input
        v-model="queryParams.keyword"
        placeholder="任务名称"
        style="width: 200px;"
        class="filter-item"
        @keyup.enter="handleSearch"
      />
      <el-select
        v-model="queryParams.status"
        placeholder="任务状态"
        clearable
        style="width: 130px"
        class="filter-item"
      >
        <el-option v-for="item in statusOptions" :key="item.value" :label="item.label" :value="item.value" />
      </el-select>
      <el-select
        v-model="queryParams.type"
        placeholder="爬虫类型"
        clearable
        style="width: 130px"
        class="filter-item"
      >
        <el-option v-for="item in crawlerTypes" :key="item.code" :label="item.name" :value="item.code" />
      </el-select>
      <el-button type="primary" class="filter-item" @click="handleSearch">搜索</el-button>
      <el-button type="success" class="filter-item" @click="handleCreate">创建任务</el-button>
    </div>

    <el-table
      v-loading="loading"
      :data="taskList"
      stripe
      border
      style="width: 100%; margin-top: 10px;"
    >
      <el-table-column prop="name" label="任务名称" min-width="150" />
      <el-table-column prop="typeName" label="爬虫类型" width="120" />
      <el-table-column label="任务状态" width="100">
        <template #default="scope">
          <el-tag :type="getStatusType(scope.row.status)">{{ getStatusText(scope.row.status) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="进度" width="200">
        <template #default="scope">
          <el-progress :percentage="scope.row.progress" :status="getProgressStatus(scope.row)" />
        </template>
      </el-table-column>
      <el-table-column prop="createdAt" label="创建时间" width="180" />
      <el-table-column label="操作" width="240" fixed="right">
        <template #default="scope">
          <el-button size="small" @click="handleDetail(scope.row)">详情</el-button>
          <el-button 
            size="small" 
            type="primary" 
            @click="handleDownload(scope.row)"
            :disabled="!canDownload(scope.row)"
          >下载</el-button>
          <el-button 
            size="small" 
            type="danger" 
            @click="handleDelete(scope.row)"
            :disabled="!canDelete(scope.row)"
          >删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <div class="pagination-container">
      <el-pagination
        background
        layout="total, sizes, prev, pager, next, jumper"
        :total="total"
        v-model:page-size="queryParams.pageSize"
        v-model:current-page="queryParams.page"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </div>

    <!-- 创建任务对话框 -->
    <el-dialog
      v-model="createDialogVisible"
      title="创建爬虫任务"
      width="700px"
    >
      <create-task @created="handleTaskCreated" @cancel="createDialogVisible = false" />
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import CreateTask from './create.vue'
import { useRouter } from 'vue-router'
// 导入API
import { getTaskList, deleteTask } from '@/admin/api/crawler/task'
import { getCrawlerTypes } from '@/admin/api/crawler/type'

const router = useRouter()
const loading = ref(false)
const taskList = ref([])
const total = ref(0)
const createDialogVisible = ref(false)

const queryParams = reactive({
  page: 1,
  pageSize: 10,
  keyword: '',
  status: '',
  type: ''
})

// 状态选项
const statusOptions = [
  { value: 'waiting', label: '等待中' },
  { value: 'running', label: '运行中' },
  { value: 'completed', label: '已完成' },
  { value: 'failed', label: '失败' },
  { value: 'cancelled', label: '已取消' }
]

// 爬虫类型
const crawlerTypes = ref([])

// 获取状态样式类型
const getStatusType = (status) => {
  const map = {
    'waiting': 'info',
    'running': 'warning',
    'completed': 'success',
    'failed': 'danger',
    'cancelled': ''
  }
  return map[status] || 'info'
}

// 获取状态文本
const getStatusText = (status) => {
  const map = {
    'waiting': '等待中',
    'running': '运行中',
    'completed': '已完成',
    'failed': '失败',
    'cancelled': '已取消'
  }
  return map[status] || '未知'
}

// 获取进度条状态
const getProgressStatus = (row) => {
  if (row.status === 'failed') return 'exception'
  if (row.status === 'completed') return 'success'
  return ''
}

// 是否可以下载
const canDownload = (row) => {
  return row.status === 'completed' && row.progress === 100
}

// 是否可以删除
const canDelete = (row) => {
  return row.status !== 'running'
}

// 搜索
const handleSearch = () => {
  queryParams.page = 1
  fetchTaskList()
}

// 页面大小变化
const handleSizeChange = (val) => {
  queryParams.pageSize = val
  fetchTaskList()
}

// 页码变化
const handleCurrentChange = (val) => {
  queryParams.page = val
  fetchTaskList()
}

// 创建任务
const handleCreate = () => {
  createDialogVisible.value = true
}

// 任务详情
const handleDetail = (row) => {
  router.push({ path: `/admin/crawler/detail/${row.id}` })
}

// 下载
const handleDownload = (row) => {
  ElMessage.success('开始下载，文件将被保存到您的下载目录')
}

// 删除
const handleDelete = (row) => {
  ElMessageBox.confirm('确定要删除此任务吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    deleteTask(row.id).then(res => {
      ElMessage.success('删除成功')
      fetchTaskList()
    })
  }).catch(() => {
    // 用户取消操作
  })
}

// 任务创建成功处理
const handleTaskCreated = () => {
  createDialogVisible.value = false
  ElMessage.success('任务创建成功')
  fetchTaskList()
}

// 获取任务列表数据
const fetchTaskList = () => {
  loading.value = true
  getTaskList(queryParams).then(res => {
    taskList.value = res.list || []
    total.value = res.total || 0
    loading.value = false
  }).catch(() => {
    loading.value = false
  })
}

// 获取爬虫类型
const fetchCrawlerTypes = () => {
      getCrawlerTypes().then(res => {
      crawlerTypes.value = res || []
    })
}

onMounted(() => {
  fetchTaskList()
  fetchCrawlerTypes()
})
</script>

<style lang="scss" scoped>
.filter-container {
  display: flex;
  align-items: center;
  margin-bottom: 20px;
  flex-wrap: wrap;
  
  .filter-item {
    margin-right: 10px;
    margin-bottom: 10px;
  }
}

.pagination-container {
  text-align: right;
  margin-top: 20px;
}
</style> 