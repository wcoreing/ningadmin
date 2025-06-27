<template>
  <div class="operation-log">
    <el-card>
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <el-input
              v-model="searchForm.username"
              placeholder="搜索用户名"
              style="width: 200px"
              clearable
              @keyup.enter="handleSearch"
            >
              <template #prefix>
                <el-icon><Search /></el-icon>
              </template>
            </el-input>
            <el-select
              v-model="searchForm.module"
              placeholder="操作模块"
              clearable
              style="width: 120px; margin-left: 16px"
            >
              <el-option label="用户管理" value="用户管理" />
              <el-option label="角色管理" value="角色管理" />
              <el-option label="菜单管理" value="菜单管理" />
              <el-option label="部门管理" value="部门管理" />
              <el-option label="系统配置" value="系统配置" />
              <el-option label="系统登录" value="系统登录" />
            </el-select>
            <el-select
              v-model="searchForm.operationType"
              placeholder="操作类型"
              clearable
              style="width: 120px; margin-left: 16px"
            >
              <el-option label="查询" value="SELECT" />
              <el-option label="新增" value="CREATE" />
              <el-option label="修改" value="UPDATE" />
              <el-option label="删除" value="DELETE" />
              <el-option label="登录" value="LOGIN" />
              <el-option label="登出" value="LOGOUT" />
            </el-select>
            <el-date-picker
              v-model="dateRange"
              type="daterange"
              range-separator="至"
              start-placeholder="开始日期"
              end-placeholder="结束日期"
              style="width: 240px; margin-left: 16px"
              @change="handleDateChange"
            />
            <el-button type="primary" @click="handleSearch" style="margin-left: 16px">
              <el-icon><Search /></el-icon>
              搜索
            </el-button>
            <el-button @click="handleReset">
              <el-icon><Refresh /></el-icon>
              重置
            </el-button>
          </div>
          <div class="header-right">
            <el-button type="danger" @click="handleBatchDelete" :disabled="!selectedRows.length">
              <el-icon><Delete /></el-icon>
              批量删除
            </el-button>
            <el-button type="warning" @click="handleClear">
              <el-icon><Delete /></el-icon>
              清空日志
            </el-button>
            <el-button type="success" @click="handleExport">
              <el-icon><Download /></el-icon>
              导出
            </el-button>
          </div>
        </div>
      </template>

      <el-table
        :data="tableData"
        v-loading="loading"
        @selection-change="handleSelectionChange"
        style="width: 100%"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="username" label="用户名" width="100" />
        <el-table-column prop="module" label="操作模块" width="100" />
        <el-table-column prop="operationName" label="操作名称" width="150" />
        <el-table-column prop="operationType" label="操作类型" width="80">
          <template #default="{ row }">
            <el-tag
              :type="getOperationTypeColor(row.operationType)"
              size="small"
            >
              {{ getOperationTypeName(row.operationType) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="method" label="请求方法" width="80">
          <template #default="{ row }">
            <el-tag
              :type="getMethodColor(row.method)"
              size="small"
            >
              {{ row.method }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="ipAddress" label="IP地址" width="120" />
        <el-table-column prop="status" label="状态" width="80">
          <template #default="{ row }">
            <el-tag
              :type="row.status === '1' ? 'success' : 'danger'"
              size="small"
            >
              {{ row.status === '1' ? '成功' : '失败' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="costTime" label="耗时(ms)" width="100" />
        <el-table-column prop="createTime" label="操作时间" width="180" />
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button
              link
              type="primary"
              @click="handleView(row)"
            >
              查看详情
            </el-button>
            <el-button
              link
              type="danger"
              @click="handleDelete(row)"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :page-sizes="[10, 20, 50, 100]"
        :total="pagination.total"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
        style="margin-top: 20px; text-align: right"
      />
    </el-card>

    <!-- 详情对话框 -->
    <el-dialog
      v-model="detailVisible"
      title="操作日志详情"
      width="800px"
    >
      <el-descriptions :column="2" border>
        <el-descriptions-item label="用户名">{{ currentRow.username }}</el-descriptions-item>
        <el-descriptions-item label="操作模块">{{ currentRow.module }}</el-descriptions-item>
        <el-descriptions-item label="操作名称">{{ currentRow.operationName }}</el-descriptions-item>
        <el-descriptions-item label="操作类型">
          <el-tag :type="getOperationTypeColor(currentRow.operationType)" size="small">
            {{ getOperationTypeName(currentRow.operationType) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="请求方法">
          <el-tag :type="getMethodColor(currentRow.method)" size="small">
            {{ currentRow.method }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="请求URL">{{ currentRow.requestUrl }}</el-descriptions-item>
        <el-descriptions-item label="IP地址">{{ currentRow.ipAddress }}</el-descriptions-item>
        <el-descriptions-item label="用户代理" :span="2">{{ currentRow.userAgent }}</el-descriptions-item>
        <el-descriptions-item label="操作状态">
          <el-tag :type="currentRow.status === '1' ? 'success' : 'danger'" size="small">
            {{ currentRow.status === '1' ? '成功' : '失败' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="耗时">{{ currentRow.costTime }}ms</el-descriptions-item>
        <el-descriptions-item label="操作时间" :span="2">{{ currentRow.createTime }}</el-descriptions-item>
      </el-descriptions>
      
      <el-divider content-position="left">请求参数</el-divider>
      <el-input
        v-model="currentRow.requestParams"
        type="textarea"
        :rows="4"
        readonly
        placeholder="无请求参数"
      />
      
      <el-divider content-position="left">响应数据</el-divider>
      <el-input
        v-model="currentRow.responseData"
        type="textarea"
        :rows="6"
        readonly
        placeholder="无响应数据"
      />
      
      <el-divider content-position="left" v-if="currentRow.errorMsg">错误信息</el-divider>
      <el-input
        v-if="currentRow.errorMsg"
        v-model="currentRow.errorMsg"
        type="textarea"
        :rows="3"
        readonly
      />
    </el-dialog>

    <!-- 清空日志对话框 -->
    <el-dialog
      v-model="clearVisible"
      title="清空操作日志"
      width="400px"
    >
      <p>请选择要保留的日志天数：</p>
      <el-input-number
        v-model="clearDays"
        :min="1"
        :max="365"
        controls-position="right"
        style="width: 100%"
      />
      <p style="margin-top: 10px; color: #999; font-size: 12px;">
        将删除 {{ clearDays }} 天前的所有操作日志，此操作不可恢复！
      </p>
      <template #footer>
        <el-button @click="clearVisible = false">取消</el-button>
        <el-button type="danger" @click="confirmClear">确认清空</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Refresh, Delete, Download } from '@element-plus/icons-vue'
import {
  getOperationLogList,
  deleteOperationLog,
  batchDeleteOperationLogs,
  clearOperationLogs
} from '@/api/system/operation-log'

// 响应式数据
const loading = ref(false)
const tableData = ref([])
const selectedRows = ref([])
const detailVisible = ref(false)
const clearVisible = ref(false)
const currentRow = ref({})
const dateRange = ref([])
const clearDays = ref(30)

// 搜索表单
const searchForm = reactive({
  username: '',
  module: '',
  operationType: '',
  startTime: '',
  endTime: ''
})

// 分页信息
const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

// 获取操作类型颜色
const getOperationTypeColor = (type) => {
  const colors = {
    'SELECT': 'info',
    'CREATE': 'success',
    'UPDATE': 'warning',
    'DELETE': 'danger',
    'LOGIN': 'success',
    'LOGOUT': 'info'
  }
  return colors[type] || 'info'
}

// 获取操作类型名称
const getOperationTypeName = (type) => {
  const names = {
    'SELECT': '查询',
    'CREATE': '新增',
    'UPDATE': '修改',
    'DELETE': '删除',
    'LOGIN': '登录',
    'LOGOUT': '登出'
  }
  return names[type] || type
}

// 获取请求方法颜色
const getMethodColor = (method) => {
  const colors = {
    'GET': 'info',
    'POST': 'success',
    'PUT': 'warning',
    'DELETE': 'danger'
  }
  return colors[method] || 'info'
}

// 获取列表数据
const getList = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      page_size: pagination.pageSize,
      ...searchForm
    }
    const response = await getOperationLogList(params)
    console.log(response)
    tableData.value = response.list || []
    pagination.total = response.total || 0
  } catch (error) {
    console.error('获取操作日志列表失败:', error)
    ElMessage.error('获取操作日志列表失败')
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  getList()
}

// 重置
const handleReset = () => {
  Object.keys(searchForm).forEach(key => {
    searchForm[key] = ''
  })
  dateRange.value = []
  pagination.page = 1
  getList()
}

// 日期范围变化
const handleDateChange = (dates) => {
  if (dates && dates.length === 2) {
    searchForm.startTime = dates[0]
    searchForm.endTime = dates[1]
  } else {
    searchForm.startTime = ''
    searchForm.endTime = ''
  }
}

// 选择变化
const handleSelectionChange = (selection) => {
  selectedRows.value = selection
}

// 查看详情
const handleView = (row) => {
  currentRow.value = { ...row }
  detailVisible.value = true
}

// 删除
const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确认删除用户"${row.username}"的操作记录吗？`,
      '警告',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await deleteOperationLog(row.id)
    ElMessage.success('删除成功')
    getList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除失败:', error)
      ElMessage.error('删除失败')
    }
  }
}

// 批量删除
const handleBatchDelete = async () => {
  try {
    await ElMessageBox.confirm(
      `确认删除选中的 ${selectedRows.value.length} 条操作记录吗？`,
      '警告',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    const ids = selectedRows.value.map(row => row.id)
    await batchDeleteOperationLogs(ids)
    ElMessage.success('批量删除成功')
    getList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('批量删除失败:', error)
      ElMessage.error('批量删除失败')
    }
  }
}

// 清空日志
const handleClear = () => {
  clearVisible.value = true
}

// 确认清空
const confirmClear = async () => {
  try {
    const result = await clearOperationLogs(clearDays.value)
    ElMessage.success(`成功清空 ${result.data.count} 条日志记录`)
    clearVisible.value = false
    getList()
  } catch (error) {
    console.error('清空日志失败:', error)
    ElMessage.error('清空日志失败')
  }
}

// 导出
const handleExport = () => {
  ElMessage.info('导出功能开发中...')
}

// 分页大小变化
const handleSizeChange = (size) => {
  pagination.pageSize = size
  pagination.page = 1
  getList()
}

// 当前页变化
const handleCurrentChange = (page) => {
  pagination.page = page
  getList()
}

// 初始化
onMounted(() => {
  getList()
})
</script>

<style scoped>
.operation-log {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 16px;
}

.header-left {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 8px;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 8px;
}

:deep(.el-descriptions__label) {
  width: 120px;
}

:deep(.el-table .el-table__cell) {
  padding: 8px 0;
}
</style>
