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
          <el-option label="待结算" value="0" />
          <el-option label="已结算" value="1" />
          <el-option label="已取消" value="2" />
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
        <div class="header-actions">
          <el-button type="success" @click="handleBatchSettle" :disabled="selectedRows.length === 0">
            批量结算
          </el-button>
          <el-button type="primary" @click="handleExport">导出数据</el-button>
        </div>
      </div>

      <!-- 数据表格 -->
      <el-table 
        v-loading="loading" 
        :data="incomeList"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" align="center" />
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
        <el-table-column label="操作" align="center" width="120">
          <template #default="scope">
            <el-button 
              v-if="scope.row.status === 0"
              type="primary" 
              size="small" 
              @click="handleSettle(scope.row)"
            >
              结算
            </el-button>
            <span v-else class="text-muted">已处理</span>
          </template>
        </el-table-column>
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
import { getIncomeList, settleIncome, batchSettleIncome } from '@/api/share/income'
import { ElMessage, ElMessageBox } from 'element-plus'

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

// 选中的行数据
const selectedRows = ref([])

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
    0: 'warning', // 待结算
    1: 'success', // 已结算
    2: 'danger'   // 已取消
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

// 选择行变化
const handleSelectionChange = (selection) => {
  selectedRows.value = selection
}

// 单个结算
const handleSettle = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确认要结算用户 ${row.userName} 的收益 ${row.amount} 元吗？`,
      '结算确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    await settleIncome({ id: row.id })
    ElMessage.success('结算成功')
    getList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('结算失败:', error)
      ElMessage.error('结算失败')
    }
  }
}

// 批量结算
const handleBatchSettle = async () => {
  if (selectedRows.value.length === 0) {
    ElMessage.warning('请选择要结算的收益记录')
    return
  }

  // 过滤出可结算的记录
  const settlableRows = selectedRows.value.filter(row => row.status === 0)
  if (settlableRows.length === 0) {
    ElMessage.warning('所选记录中没有可结算的收益')
    return
  }

  const totalAmount = settlableRows.reduce((sum, row) => sum + row.amount, 0)

  try {
    await ElMessageBox.confirm(
      `确认要批量结算 ${settlableRows.length} 条记录，总金额 ${totalAmount.toFixed(2)} 元吗？`,
      '批量结算确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    const ids = settlableRows.map(row => row.id)
    await batchSettleIncome({ ids })
    ElMessage.success('批量结算成功')
    getList()
    selectedRows.value = []
  } catch (error) {
    if (error !== 'cancel') {
      console.error('批量结算失败:', error)
      ElMessage.error('批量结算失败')
    }
  }
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

    .header-actions {
      display: flex;
      gap: 12px;
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

.text-muted {
  color: #999;
  font-size: 12px;
}
</style> 