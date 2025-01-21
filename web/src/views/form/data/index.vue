<template>
  <div class="form-data">
    <!-- 搜索区域 -->
    <el-card class="search-card">
      <el-form :model="queryParams" ref="queryRef" :inline="true">
        <el-form-item label="表单名称" prop="formName">
          <el-input
            v-model="queryParams.formName"
            placeholder="请输入表单名称"
            clearable
            @keyup.enter="handleQuery"
          />
        </el-form-item>
        <el-form-item label="提交时间" prop="timeRange">
          <el-date-picker
            v-model="queryParams.timeRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            value-format="YYYY-MM-DD"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :icon="Search" @click="handleQuery">搜索</el-button>
          <el-button :icon="Refresh" @click="resetQuery">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 表格区域 -->
    <el-card class="table-card">
      <el-table v-loading="loading" :data="dataList">
        <el-table-column type="index" label="序号" width="60" />
        <el-table-column prop="formName" label="表单名称" min-width="120" show-overflow-tooltip />
        <el-table-column prop="submitter.name" label="提交人" min-width="100" show-overflow-tooltip>
          <template #default="{ row }">
            {{ row.submitter?.name || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="ip" label="IP地址" width="130" show-overflow-tooltip />
        <el-table-column prop="status" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === '0' ? 'success' : row.status === '2' ? 'warning' : 'danger'" size="small">
              {{ row.status === '0' ? '成功' : row.status === '2' ? '处理中' : '失败' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="createTime" label="提交时间" width="160">
          <template #default="{ row }">
            {{ formatDateTime(row.createTime) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="120" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleView(row)">查看</el-button>
            <el-button type="danger" link @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页区域 -->
      <el-pagination
        v-model:current-page="queryParams.pageNum"
        v-model:page-size="queryParams.pageSize"
        :total="total"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next"
        class="pagination"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </el-card>

    <!-- 查看对话框 -->
    <el-dialog
      v-model="viewDialog.visible"
      title="表单数据详情"
      width="800px"
      append-to-body
      destroy-on-close
    >
      <el-descriptions :column="2" border>
        <el-descriptions-item label="表单名称" :span="2">
          {{ currentData.formName }}
        </el-descriptions-item>
        <el-descriptions-item label="提交人">
          {{ currentData.submitter?.name || '-' }}
        </el-descriptions-item>
        <el-descriptions-item label="提交人邮箱">
          {{ currentData.submitter?.email || '-' }}
        </el-descriptions-item>
        <el-descriptions-item label="IP地址">
          {{ currentData.ip }}
        </el-descriptions-item>
        <el-descriptions-item label="提交时间">
          {{ formatDateTime(currentData.createTime) }}
        </el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="currentData.status === '0' ? 'success' : currentData.status === '2' ? 'warning' : 'danger'" size="small">
            {{ currentData.status === '0' ? '成功' : currentData.status === '2' ? '处理中' : '失败' }}
          </el-tag>
        </el-descriptions-item>
      </el-descriptions>

      <!-- 表单数据展示 -->
      <div class="form-data-content">
        <el-divider>表单数据</el-divider>
        <el-descriptions :column="1" border>
          <template v-if="currentData.data">
            <el-descriptions-item 
              v-for="(field, key) in currentData.data" 
              :key="key"
              :label="field.label"
            >
              <!-- 处理不同类型的数据展示 -->
              <template v-if="Array.isArray(field.value) && field.value[0]?.url">
                <!-- 图片类型 -->
                <div class="image-list">
                  <el-image
                    v-for="(item, index) in field.value"
                    :key="index"
                    :src="item.url"
                    :preview-src-list="field.value.map(img => img.url)"
                    fit="cover"
                    class="preview-image"
                  >
                    <template #error>
                      <div class="image-error">
                        <el-icon><Picture /></el-icon>
                      </div>
                    </template>
                  </el-image>
                </div>
              </template>
              <template v-else-if="Array.isArray(field.value)">
                <!-- 数组类型（如多选） -->
                {{ field.value.join(', ') }}
              </template>
              <template v-else>
                <!-- 普通文本 -->
                {{ field.value }}
              </template>
            </el-descriptions-item>
          </template>
        </el-descriptions>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Refresh, Picture } from '@element-plus/icons-vue'
import { getFormDataList, deleteFormData } from '@/api/form-data'

// 查询参数
const queryParams = reactive({
  pageNum: 1,
  pageSize: 10,
  formName: '',
  timeRange: []
})

// 数据列表
const dataList = ref([])
const total = ref(0)
const loading = ref(false)
const queryRef = ref()

// 查看对话框
const viewDialog = reactive({
  visible: false
})

// 当前查看的数据
const currentData = ref({})

// 获取数据列表
const getList = async () => {
  loading.value = true
  try {
    // 处理时间范围
    const params = {
      ...queryParams,
      startTime: queryParams.timeRange?.[0],
      endTime: queryParams.timeRange?.[1]
    }
    delete params.timeRange

    const res = await getFormDataList(params)
    // 直接使用返回的数据，不需要解构
    dataList.value = res.list || []
    total.value = res.total || 0
  } catch (error) {
    console.error('获取数据失败:', error)
    ElMessage.error('获取数据失败：' + (error.message || '未知错误'))
  } finally {
    loading.value = false
  }
}

// 查询操作
const handleQuery = () => {
  queryParams.pageNum = 1
  getList()
}

// 重置查询
const resetQuery = () => {
  queryRef.value?.resetFields()
  handleQuery()
}

// 查看详情
const handleView = (row) => {
  currentData.value = row
  viewDialog.visible = true
}

// 删除数据
const handleDelete = (row) => {
  ElMessageBox.confirm(
    '确认要删除这条数据吗？',
    '警告',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(async () => {
    try {
      await deleteFormData(row.id)
      ElMessage.success('删除成功')
      getList()
    } catch (error) {
      console.error('删除失败:', error)
      ElMessage.error('删除失败：' + error.message)
    }
  }).catch(() => {})
}

// 处理分页大小变化
const handleSizeChange = (val) => {
  queryParams.pageSize = val
  getList()
}

// 处理页码变化
const handleCurrentChange = (val) => {
  queryParams.pageNum = val
  getList()
}

// 格式化时间
const formatDateTime = (time) => {
  if (!time) return ''
  const date = new Date(time)
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  const seconds = String(date.getSeconds()).padStart(2, '0')
  return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`
}

onMounted(() => {
  getList()
})
</script>

<style scoped>
.form-data {
  padding: 20px;
}

.search-card {
  margin-bottom: 20px;
}

.table-card {
  margin-bottom: 20px;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.form-data-content {
  margin-top: 20px;
}

.image-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.preview-image {
  width: 100px;
  height: 100px;
  border-radius: 4px;
  overflow: hidden;
}

.image-error {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #f5f7fa;
  color: #909399;
}

:deep(.el-descriptions__label) {
  width: 120px;
  justify-content: flex-end;
}

:deep(.el-dialog__body) {
  padding: 20px;
  max-height: calc(90vh - 120px);
  overflow-y: auto;
}
</style> 