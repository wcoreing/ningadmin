<template>
  <div class="app-container">
    <!-- 搜索表单 -->
    <el-form :model="queryParams" ref="queryForm" :inline="true" class="search-form">
      <el-form-item label="用户ID" prop="userId">
        <el-input v-model="queryParams.userId" placeholder="请输入用户ID" clearable />
      </el-form-item>
      <el-form-item label="用户昵称" prop="nickname">
        <el-input v-model="queryParams.nickname" placeholder="请输入用户昵称" clearable />
      </el-form-item>
      <el-form-item label="收益状态" prop="status" style="width: 250px;">
        <el-select v-model="queryParams.status" placeholder="请选择状态" clearable>
          <el-option label="冻结中" value="0" />
          <el-option label="已解冻" value="1" />
          <el-option label="已失效" value="2" />
          <el-option label="已入账" value="3" />
        </el-select>
      </el-form-item>
      <el-form-item label="时间范围" prop="timeRange">
        <el-date-picker
          v-model="dateRange"
          type="daterange"
          range-separator="至"
          start-placeholder="开始日期"
          end-placeholder="结束日期"
          value-format="YYYY-MM-DD"
        />
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="handleQuery">搜索</el-button>
        <el-button @click="resetQuery">重置</el-button>
      </el-form-item>
    </el-form>

    <!-- 分销收益明细 -->
    <div class="income-detail">
      <div class="detail-header">
        <div class="title">分销收益明细</div>
        <el-button type="primary" @click="handleExport">导出数据</el-button>
      </div>

      <!-- 数据表格 -->
      <el-table v-loading="loading" :data="incomeList">
        <el-table-column prop="userId" label="用户ID" align="center" />
        <el-table-column prop="userName" label="用户昵称" align="center" />
        <el-table-column prop="orderNo" label="订单编号" align="center" width="180" />
        <el-table-column label="收益金额" align="center">
          <template #default="scope">
            <span class="amount-text">+{{ scope.row.amount }}</span>
          </template>
        </el-table-column>
        <el-table-column label="收益类型" align="center">
          <template #default="scope">
            {{ getLevelText(scope.row.level) }}
          </template>
        </el-table-column>
        <el-table-column label="状态" align="center">
          <template #default="scope">
            <el-tag :type="getStatusType(scope.row.status)">
              {{ scope.row.statusText }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="createTime" label="创建时间" align="center" width="160" />
      </el-table>

      <!-- 分页 -->
      <el-pagination
        v-model:current-page="queryParams.pageNum"
        v-model:page-size="queryParams.pageSize"
        :total="total"
        :page-sizes="[10, 20, 30, 50]"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { getIncomeList } from '@/api/share/income'
import { ElMessage } from 'element-plus'

// 加载状态
const loading = ref(false)

// 查询参数
const queryParams = reactive({
  userId: '',
  nickname: '',
  status: '',
  startTime: '',
  endTime: '',
  pageNum: 1,
  pageSize: 10
})

// 日期范围
const dateRange = ref([])

// 收益列表数据
const incomeList = ref([])
const total = ref(0)

// 获取收益类型文本
const getLevelText = (level) => {
  const levelMap = {
    1: '一级分销',
    2: '二级分销',
    3: '三级分销'
  }
  return levelMap[level] || '-'
}

// 获取状态标签类型
const getStatusType = (status) => {
  const statusMap = {
    0: 'info',    // 冻结中
    1: 'success', // 已解冻
    2: 'danger',  // 已失效
    3: 'warning'  // 已入账
  }
  return statusMap[status] || ''
}

// 查询列表
const getList = async () => {
  loading.value = true
  try {
    // 处理时间范围
    if (dateRange.value && dateRange.value.length === 2) {
      queryParams.startTime = dateRange.value[0]
      queryParams.endTime = dateRange.value[1]
    }

    const res = await getIncomeList(queryParams)
    incomeList.value = res.list
    total.value = res.total
  } catch (error) {
    console.error('获取列表失败:', error)
    ElMessage.error('获取列表失败')
  } finally {
    loading.value = false
  }
}

// 搜索按钮操作
const handleQuery = () => {
  queryParams.pageNum = 1
  getList()
}

// 重置按钮操作
const resetQuery = () => {
  dateRange.value = []
  Object.assign(queryParams, {
    userId: '',
    nickname: '',
    status: '',
    startTime: '',
    endTime: '',
    pageNum: 1,
    pageSize: 10
  })
  handleQuery()
}

// 分页大小改变
const handleSizeChange = (val) => {
  queryParams.pageSize = val
  getList()
}

// 页码改变
const handleCurrentChange = (val) => {
  queryParams.pageNum = val
  getList()
}

// 导出数据
const handleExport = () => {
  ElMessage.success('导出功能待实现')
}

// 初始化加载数据
getList()
</script>

<style lang="scss" scoped>
.income-detail {
  background: #fff;
  padding: 20px;
  border-radius: 4px;

  .detail-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;

    .title {
      font-size: 16px;
      font-weight: bold;
    }
  }
}

.amount-text {
  color: #f56c6c;
  font-weight: bold;
}

.search-form {
  background: #fff;
  padding: 20px;
  margin-bottom: 20px;
  border-radius: 4px;
}
</style> 