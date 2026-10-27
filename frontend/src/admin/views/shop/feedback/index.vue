<template>
  <div class="app-container">
    <!-- 搜索区域 -->
    <el-card class="search-card">
      <el-form :model="queryParams" ref="queryForm" :inline="true">
        <el-form-item label="反馈内容" prop="keyword">
          <el-input
            v-model="queryParams.keyword"
            placeholder="请输入关键字"
            clearable
            style="width: 240px"
            @keyup.enter="handleQuery"
          />
        </el-form-item>
        <el-form-item label="用户ID" prop="userId">
          <el-input
            v-model="queryParams.userId"
            placeholder="请输入用户ID"
            clearable
            style="width: 200px"
            @keyup.enter="handleQuery"
          />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-select v-model="queryParams.status" placeholder="请选择状态" clearable style="width: 150px">
            <el-option
              v-for="dict in statusOptions"
              :key="dict.value"
              :label="dict.label"
              :value="dict.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleQuery">查询</el-button>
          <el-button @click="resetQuery">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 操作按钮区域 -->
    <el-card class="box-card">
      <template #header>
        <div class="card-header">
          <el-button type="primary" @click="handleAdd">新增反馈</el-button>
        </div>
      </template>

      <!-- 表格区域 -->
      <el-table
        v-loading="loading"
        :data="feedbackList"
        row-key="id"
        border
      >
        <el-table-column label="ID" prop="id" width="80" align="center" />
        <el-table-column label="用户ID" prop="userId" width="100" align="center" />
        <el-table-column label="反馈内容" prop="content" min-width="200" show-overflow-tooltip />
        <el-table-column label="联系方式" prop="contact" width="120" align="center" />
        <el-table-column label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="回复内容" prop="reply" min-width="150" show-overflow-tooltip>
          <template #default="{ row }">
            <span v-if="row.reply">{{ row.reply }}</span>
            <span v-else class="text-gray">暂无回复</span>
          </template>
        </el-table-column>
        <el-table-column label="回复时间" prop="replyTime" width="160" align="center">
          <template #default="{ row }">
            <span v-if="row.replyTime">{{ formatTime(row.replyTime) }}</span>
            <span v-else class="text-gray">-</span>
          </template>
        </el-table-column>
        <el-table-column label="创建时间" prop="createTime" width="160" align="center">
          <template #default="{ row }">
            {{ formatTime(row.createTime) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" align="center" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="handleView(row)">查看</el-button>
            <el-button 
              v-if="row.status === 0" 
              link 
              type="success" 
              @click="handleReply(row)"
            >
              回复
            </el-button>
            <el-button link type="danger" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页区域 -->
      <el-pagination
        v-show="total > 0"
        v-model:current-page="queryParams.pageNum"
        v-model:page-size="queryParams.pageSize"
        :total="total"
        :page-sizes="[10, 20, 30, 50]"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </el-card>

    <!-- 新增/编辑对话框 -->
    <el-dialog
      :title="dialog.title"
      v-model="dialog.visible"
      width="600px"
      append-to-body
      destroy-on-close
      :close-on-click-modal="false"
    >
      <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
        <el-form-item label="用户ID" prop="userId">
          <el-input v-model="form.userId" placeholder="请输入用户ID" type="number" />
        </el-form-item>
        <el-form-item label="反馈内容" prop="content">
          <el-input 
            v-model="form.content" 
            type="textarea" 
            :rows="4"
            placeholder="请输入反馈内容" 
          />
        </el-form-item>
        <el-form-item label="联系方式" prop="contact">
          <el-input v-model="form.contact" placeholder="请输入联系方式（可选）" />
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button type="primary" @click="submitForm">确 定</el-button>
          <el-button @click="cancel">取 消</el-button>
        </div>
      </template>
    </el-dialog>

    <!-- 查看详情对话框 -->
    <el-dialog
      title="反馈详情"
      v-model="detailVisible"
      width="700px"
      append-to-body
    >
      <div v-if="currentFeedback" class="feedback-detail">
        <el-descriptions border :column="2">
          <el-descriptions-item label="反馈ID">{{ currentFeedback.id }}</el-descriptions-item>
          <el-descriptions-item label="用户ID">{{ currentFeedback.userId }}</el-descriptions-item>
          <el-descriptions-item label="联系方式">
            {{ currentFeedback.contact || '无' }}
          </el-descriptions-item>
          <el-descriptions-item label="状态">
            <el-tag :type="getStatusType(currentFeedback.status)">
              {{ getStatusText(currentFeedback.status) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="创建时间" :span="2">
            {{ formatTime(currentFeedback.createTime) }}
          </el-descriptions-item>
        </el-descriptions>

        <div class="content-section">
          <h4>反馈内容</h4>
          <div class="content-box">
            {{ currentFeedback.content }}
          </div>
        </div>

        <div v-if="currentFeedback.reply" class="reply-section">
          <h4>回复内容</h4>
          <div class="reply-box">
            {{ currentFeedback.reply }}
          </div>
          <div class="reply-time">
            回复时间：{{ formatTime(currentFeedback.replyTime) }}
          </div>
        </div>
      </div>
    </el-dialog>

    <!-- 回复对话框 -->
    <el-dialog
      title="回复反馈"
      v-model="replyVisible"
      width="600px"
      append-to-body
      :close-on-click-modal="false"
    >
      <el-form ref="replyFormRef" :model="replyForm" :rules="replyRules" label-width="80px">
        <el-form-item label="回复内容" prop="reply">
          <el-input 
            v-model="replyForm.reply" 
            type="textarea" 
            :rows="5"
            placeholder="请输入回复内容" 
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button type="primary" @click="submitReply">确 定</el-button>
          <el-button @click="cancelReply">取 消</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { 
  getFeedbackList, 
  createFeedback, 
  deleteFeedback, 
  getFeedbackInfo,
  replyFeedback,
  updateFeedbackStatus
} from '@/admin/api/shop/feedback'

// 响应式数据
const loading = ref(false)
const total = ref(0)
const feedbackList = ref([])

// 查询参数
const queryParams = ref({
  pageNum: 1,
  pageSize: 10,
  keyword: '',
  userId: '',
  status: null
})

// 表单数据
const form = ref({
  userId: '',
  content: '',
  contact: ''
})

// 对话框状态
const dialog = ref({
  visible: false,
  title: ''
})

const detailVisible = ref(false)
const replyVisible = ref(false)
const currentFeedback = ref(null)

// 回复表单
const replyForm = ref({
  id: null,
  reply: ''
})

// 表单引用
const formRef = ref()
const replyFormRef = ref()

// 状态选项
const statusOptions = [
  { label: '待处理', value: 0 },
  { label: '已处理', value: 1 },
  { label: '已关闭', value: 2 }
]

// 表单验证规则
const rules = {
  userId: [{ required: true, message: '请输入用户ID', trigger: 'blur' }],
  content: [{ required: true, message: '请输入反馈内容', trigger: 'blur' }]
}

const replyRules = {
  reply: [{ required: true, message: '请输入回复内容', trigger: 'blur' }]
}

// 获取反馈列表
const getList = async () => {
  loading.value = true
  try {
    const params = { ...queryParams.value }
    // 清空无效的参数
    Object.keys(params).forEach(key => {
      if (params[key] === '' || params[key] === null || params[key] === undefined) {
        delete params[key]
      }
    })
    
    const res = await getFeedbackList(params)
    feedbackList.value = res.list || []
    total.value = res.total || 0
  } catch (error) {
    console.error('获取反馈列表失败:', error)
    ElMessage.error('获取反馈列表失败')
  } finally {
    loading.value = false
  }
}

// 重置表单
const reset = () => {
  form.value = {
    userId: '',
    content: '',
    contact: ''
  }
  formRef.value?.resetFields()
}

// 查询
const handleQuery = () => {
  queryParams.value.pageNum = 1
  getList()
}

// 重置查询
const resetQuery = () => {
  queryParams.value = {
    pageNum: 1,
    pageSize: 10,
    keyword: '',
    userId: '',
    status: null
  }
  getList()
}

// 新增
const handleAdd = () => {
  reset()
  dialog.value = {
    visible: true,
    title: '新增反馈'
  }
}

// 查看详情
const handleView = async (row) => {
  try {
    const res = await getFeedbackInfo(row.id)
    currentFeedback.value = res
    detailVisible.value = true
  } catch (error) {
    console.error('获取反馈详情失败:', error)
    ElMessage.error('获取反馈详情失败')
  }
}

// 回复反馈
const handleReply = (row) => {
  replyForm.value = {
    id: row.id,
    reply: ''
  }
  replyVisible.value = true
}

// 删除
const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm('确定要删除这条反馈吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    await deleteFeedback(row.id)
    ElMessage.success('删除成功')
    getList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除失败:', error)
      ElMessage.error('删除失败')
    }
  }
}

// 提交表单
const submitForm = async () => {
  const valid = await formRef.value.validate()
  if (!valid) return
  
  try {
    await createFeedback(form.value)
    ElMessage.success('新增成功')
    dialog.value.visible = false
    getList()
  } catch (error) {
    console.error('新增失败:', error)
    ElMessage.error('新增失败')
  }
}

// 提交回复
const submitReply = async () => {
  const valid = await replyFormRef.value.validate()
  if (!valid) return
  
  try {
    await replyFeedback(replyForm.value.id, replyForm.value.reply)
    ElMessage.success('回复成功')
    replyVisible.value = false
    getList()
  } catch (error) {
    console.error('回复失败:', error)
    ElMessage.error('回复失败')
  }
}

// 取消
const cancel = () => {
  dialog.value.visible = false
}

const cancelReply = () => {
  replyVisible.value = false
}

// 分页
const handleSizeChange = (val) => {
  queryParams.value.pageSize = val
  queryParams.value.pageNum = 1
  getList()
}

const handleCurrentChange = (val) => {
  queryParams.value.pageNum = val
  getList()
}

// 获取状态类型
const getStatusType = (status) => {
  const typeMap = {
    0: '',           // 待处理 - 默认
    1: 'success',    // 已处理
    2: 'info'        // 已关闭
  }
  return typeMap[status] || ''
}

// 获取状态文本
const getStatusText = (status) => {
  const textMap = {
    0: '待处理',
    1: '已处理',
    2: '已关闭'
  }
  return textMap[status] || '未知'
}

// 格式化时间
const formatTime = (time) => {
  if (!time) return ''
  return new Date(time).toLocaleString('zh-CN')
}

// 初始化
onMounted(() => {
  getList()
})
</script>

<style lang="scss" scoped>
.search-card {
  margin-bottom: 15px;

  :deep(.el-form-item) {
    margin-bottom: 12px;
    margin-right: 12px;
  }
}

.feedback-detail {
  .content-section, .reply-section {
    margin-top: 20px;
    
    h4 {
      margin: 0 0 10px 0;
      color: #303133;
      font-size: 14px;
      font-weight: 600;
    }
  }
  
  .content-box, .reply-box {
    padding: 12px;
    background: #f5f7fa;
    border-radius: 4px;
    border: 1px solid #e4e7ed;
    white-space: pre-wrap;
    line-height: 1.6;
  }
  
  .reply-box {
    background: #f0f9ff;
    border-color: #c6f7ff;
  }
  
  .reply-time {
    margin-top: 8px;
    font-size: 12px;
    color: #909399;
    text-align: right;
  }
}

.text-gray {
  color: #909399;
}
</style> 