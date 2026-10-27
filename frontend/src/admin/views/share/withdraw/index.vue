<template>
  <div class="app-container">
    <!-- 搜索区域 -->
    <el-card class="search-card">
      <el-form :model="queryParams" ref="queryForm" :inline="true">
        <el-form-item label="申请人" prop="realName">
          <el-input
            v-model="queryParams.realName"
            placeholder="请输入申请人姓名"
            clearable
            @keyup.enter="handleQuery"
          />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-select v-model="queryParams.status" placeholder="请选择状态" clearable>
            <el-option label="待审核" :value="0" />
            <el-option label="已通过" :value="1" />
            <el-option label="已拒绝" :value="2" />
            <el-option label="已打款" :value="3" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="Search" @click="handleQuery">搜索</el-button>
          <el-button icon="Refresh" @click="resetQuery">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 列表区域 -->
    <el-card class="table-card">
      <el-table
        v-loading="loading"
        :data="withdrawList"
        style="width: 100%"
      >
        <el-table-column label="申请人" prop="realName" />
        <el-table-column label="提现金额" prop="amount">
          <template #default="{ row }">
            ¥{{ row.amount }}
          </template>
        </el-table-column>
        <el-table-column label="提现方式" prop="method">
          <template #default="{ row }">
            {{ methodMap[row.method] }}
          </template>
        </el-table-column>
        <el-table-column label="收款账号" prop="account" />
        <el-table-column label="申请时间" prop="createTime" width="160" />
        <el-table-column label="状态" prop="status" width="100">
          <template #default="{ row }">
            <el-tag :type="statusTypeMap[row.status]">
              {{ statusMap[row.status] }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button
              v-if="row.status === 0"
              size="small"
              type="success"
              @click="handleApprove(row)"
            >通过</el-button>
            <el-button
              v-if="row.status === 0"
              size="small"
              type="danger"
              @click="handleReject(row)"
            >拒绝</el-button>
            <el-button
              v-if="row.status === 1"
              size="small"
              type="primary"
              @click="handlePay(row)"
            >打款</el-button>
            <el-button
              size="small"
              type="primary"
              link
              @click="handleDetail(row)"
            >详情</el-button>
          </template>
        </el-table-column>
      </el-table>
      
      <!-- 分页 -->
      <pagination
        v-show="total > 0"
        :total="total"
        :page="queryParams.pageNum"
        :limit="queryParams.pageSize"
        @pagination="handlePagination"
      />
    </el-card>

    <!-- 详情弹窗 -->
    <el-dialog
      title="提现详情"
      v-model="detailDialogVisible"
      width="500px"
      append-to-body
    >
      <el-descriptions :column="2" border>
        <el-descriptions-item label="申请人">{{ detailData.realName }}</el-descriptions-item>
        <el-descriptions-item label="提现金额">¥{{ detailData.amount }}</el-descriptions-item>
        <el-descriptions-item label="提现方式">{{ methodMap[detailData.method] }}</el-descriptions-item>
        <el-descriptions-item label="收款账号">{{ detailData.account }}</el-descriptions-item>
        <el-descriptions-item label="申请时间">{{ detailData.createTime }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="statusTypeMap[detailData.status]">
            {{ statusMap[detailData.status] }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="审核时间">{{ detailData.auditTime || '-' }}</el-descriptions-item>
        <el-descriptions-item label="审核备注">{{ detailData.auditRemark || '-' }}</el-descriptions-item>
        <el-descriptions-item label="打款时间">{{ detailData.payTime || '-' }}</el-descriptions-item>
        <el-descriptions-item label="打款备注">{{ detailData.payRemark || '-' }}</el-descriptions-item>
      </el-descriptions>
    </el-dialog>

    <!-- 拒绝弹窗 -->
    <el-dialog
      title="拒绝提现"
      v-model="rejectDialogVisible"
      width="500px"
      append-to-body
    >
      <el-form ref="rejectFormRef" :model="rejectForm" :rules="rejectRules">
        <el-form-item label="拒绝原因" prop="remark">
          <el-input
            v-model="rejectForm.remark"
            type="textarea"
            :rows="3"
            placeholder="请输入拒绝原因"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="rejectDialogVisible = false">取 消</el-button>
        <el-button type="primary" @click="confirmReject">确 定</el-button>
      </template>
    </el-dialog>

    <!-- 打款弹窗 -->
    <el-dialog
      title="确认打款"
      v-model="payDialogVisible"
      width="500px"
      append-to-body
    >
      <el-form ref="payFormRef" :model="payForm" :rules="payRules">
        <el-form-item label="打款备注" prop="remark">
          <el-input
            v-model="payForm.remark"
            type="textarea"
            :rows="3"
            placeholder="请输入打款备注"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="payDialogVisible = false">取 消</el-button>
        <el-button type="primary" @click="confirmPay">确 定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { listWithdraw, getWithdrawDetail, approveWithdraw, rejectWithdraw, confirmWithdrawPay } from '@/admin/api/share/withdraw'

// 状态映射
const statusMap = {
  0: '待审核',
  1: '已通过',
  2: '已拒绝',
  3: '已打款'
}

const statusTypeMap = {
  0: 'warning',
  1: 'primary',
  2: 'danger',
  3: 'success'
}

// 提现方式映射
const methodMap = {
  'bank': '银行卡',
  'wechat': '微信',
  'alipay': '支付宝'
}

// 查询参数
const queryParams = reactive({
  pageNum: 1,
  pageSize: 10,
  realName: '',
  status: undefined
})

// 数据相关
const loading = ref(false)
const withdrawList = ref([])
const total = ref(0)

// 详情相关
const detailDialogVisible = ref(false)
const detailData = ref({})

// 拒绝相关
const rejectDialogVisible = ref(false)
const rejectFormRef = ref()
const rejectForm = reactive({
  id: '',
  remark: ''
})
const rejectRules = {
  remark: [
    { required: true, message: '请输入拒绝原因', trigger: 'blur' }
  ]
}

// 打款相关
const payDialogVisible = ref(false)
const payFormRef = ref()
const payForm = reactive({
  id: '',
  remark: ''
})
const payRules = {
  remark: [
    { required: true, message: '请输入打款备注', trigger: 'blur' }
  ]
}

// 获取列表数据
const getList = async () => {
  try {
    loading.value = true
    const res = await listWithdraw(queryParams)
    withdrawList.value = res.list
    total.value = res.total
  } catch (error) {
    console.error('获取列表失败:', error)
  } finally {
    loading.value = false
  }
}

// 搜索
const handleQuery = () => {
  queryParams.pageNum = 1
  getList()
}

// 重置
const resetQuery = () => {
  queryParams.realName = ''
  queryParams.status = undefined
  handleQuery()
}

// 查看详情
const handleDetail = async (row) => {
  try {
    const res = await getWithdrawDetail(row.id)
    detailData.value = res
    detailDialogVisible.value = true
  } catch (error) {
    console.error('获取详情失败:', error)
  }
}

// 通过
const handleApprove = (row) => {
  ElMessageBox.confirm('确认通过该提现申请吗？', '提示', {
    type: 'warning'
  }).then(async () => {
    try {
      await approveWithdraw(row.id)
      ElMessage.success('操作成功')
      getList()
    } catch (error) {
      console.error('操作失败:', error)
    }
  })
}

// 拒绝
const handleReject = (row) => {
  rejectForm.id = row.id
  rejectForm.remark = ''
  rejectDialogVisible.value = true
}

// 确认拒绝
const confirmReject = () => {
  rejectFormRef.value.validate(async (valid) => {
    if (valid) {
      try {
        await rejectWithdraw(rejectForm)
        ElMessage.success('操作成功')
        rejectDialogVisible.value = false
        getList()
      } catch (error) {
        console.error('操作失败:', error)
      }
    }
  })
}

// 打款
const handlePay = (row) => {
  payForm.id = row.id
  payForm.remark = ''
  payDialogVisible.value = true
}

// 确认打款
const confirmPay = () => {
  payFormRef.value.validate(async (valid) => {
    if (valid) {
      try {
        await confirmWithdrawPay(payForm)
        ElMessage.success('操作成功')
        payDialogVisible.value = false
        getList()
      } catch (error) {
        console.error('操作失败:', error)
      }
    }
  })
}

// 添加分页处理函数
const handlePagination = ({ page, limit }) => {
  queryParams.pageNum = page
  queryParams.pageSize = limit
  getList()
}

// 初始化
getList()
</script>

<style scoped>
.search-card {
  margin-bottom: 15px;
}

.table-card {
  margin-bottom: 15px;
}
</style> 