<template>
  <div class="popup-management">
    <!-- 搜索区域 -->
    <el-card class="search-card">
      <el-form :model="queryParams" inline>
        <el-form-item label="标题">
          <el-input
            v-model="queryParams.title"
            placeholder="请输入弹窗标题"
            clearable
            style="width: 200px"
          />
        </el-form-item>
        <el-form-item label="触发类型">
          <el-select v-model="queryParams.trigger_type" placeholder="请选择触发类型" clearable>
            <el-option label="手动触发" value="manual" />
            <el-option label="自动触发" value="auto" />
            <el-option label="定时触发" value="schedule" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="queryParams.status" placeholder="请选择状态" clearable>
            <el-option label="启用" :value="1" />
            <el-option label="禁用" :value="0" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="getList" :icon="Search">搜索</el-button>
          <el-button @click="resetQuery" :icon="Refresh">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 操作区域 -->
    <el-card class="table-card">
      <div class="table-header">
        <el-button type="primary" @click="handleAdd" :icon="Plus">新增弹窗</el-button>
      </div>

      <!-- 表格 -->
      <el-table v-loading="loading" :data="popupList" style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="title" label="标题" min-width="150" />
        <el-table-column prop="trigger_type" label="触发类型" width="100">
          <template #default="{ row }">
            <el-tag :type="getTriggerTypeTag(row.trigger_type)">
              {{ getTriggerTypeName(row.trigger_type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="priority" label="优先级" width="80" />
        <el-table-column prop="status" label="状态" width="80">
          <template #default="{ row }">
            <el-switch
              v-model="row.status"
              :active-value="1"
              :inactive-value="0"
              @change="handleStatusChange(row)"
            />
          </template>
        </el-table-column>
        <el-table-column prop="show_count" label="显示次数" width="100" />
        <el-table-column prop="click_count" label="点击次数" width="100" />
        <el-table-column prop="start_time" label="开始时间" width="160">
          <template #default="{ row }">
            {{ row.start_time ? formatDate(row.start_time) : '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="end_time" label="结束时间" width="160">
          <template #default="{ row }">
            {{ row.end_time ? formatDate(row.end_time) : '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="create_time" label="创建时间" width="160">
          <template #default="{ row }">
            {{ formatDate(row.create_time) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
            <el-button type="success" link @click="handlePreview(row)">预览</el-button>
            <el-button type="danger" link @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination">
        <el-pagination
          v-model:current-page="queryParams.page"
          v-model:page-size="queryParams.page_size"
          :total="total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="getList"
          @current-change="getList"
        />
      </div>
    </el-card>

    <!-- 弹窗表单 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="800px"
      :close-on-click-modal="false"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="120px"
      >
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="弹窗标题" prop="title">
              <el-input v-model="form.title" placeholder="请输入弹窗标题" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="触发类型" prop="trigger_type">
              <el-select v-model="form.trigger_type" placeholder="请选择触发类型">
                <el-option label="手动触发" value="manual" />
                <el-option label="自动触发" value="auto" />
                <el-option label="定时触发" value="schedule" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="弹窗内容" prop="content">
          <el-input
            v-model="form.content"
            type="textarea"
            :rows="4"
            placeholder="请输入弹窗内容"
          />
        </el-form-item>

        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item label="优先级" prop="priority">
              <el-input-number v-model="form.priority" :min="0" :max="999" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="自动关闭时间">
              <el-input-number v-model="form.auto_close" :min="0" placeholder="毫秒，0为不自动关闭" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="延迟显示时间">
              <el-input-number v-model="form.delay" :min="0" placeholder="毫秒" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="开始时间">
              <el-date-picker
                v-model="form.start_time"
                type="datetime"
                placeholder="选择开始时间"
                format="YYYY-MM-DD HH:mm:ss"
                value-format="YYYY-MM-DD HH:mm:ss"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="结束时间">
              <el-date-picker
                v-model="form.end_time"
                type="datetime"
                placeholder="选择结束时间"
                format="YYYY-MM-DD HH:mm:ss"
                value-format="YYYY-MM-DD HH:mm:ss"
              />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="点击遮罩关闭">
              <el-switch v-model="form.mask_closable" :active-value="1" :inactive-value="0" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="状态">
              <el-switch v-model="form.status" :active-value="1" :inactive-value="0" />
            </el-form-item>
          </el-col>
        </el-row>

        <!-- 高级配置 -->
        <el-divider>高级配置</el-divider>
        
        <el-form-item label="标题样式">
          <el-input
            v-model="titleStyleStr"
            type="textarea"
            :rows="2"
            placeholder='例如: {"color": "#333", "fontSize": "18px"}'
          />
        </el-form-item>

        <el-form-item label="按钮配置">
          <el-input
            v-model="buttonsStr"
            type="textarea"
            :rows="3"
            placeholder='例如: [{"text": "确定", "type": "primary", "action": "confirm"}, {"text": "取消", "type": "default", "action": "cancel"}]'
          />
        </el-form-item>

        <el-form-item label="触发规则">
          <el-input
            v-model="triggerRuleStr"
            type="textarea"
            :rows="3"
            placeholder='例如: {"pages": ["/home"], "userType": "new", "frequency": "once"}'
          />
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitForm">确定</el-button>
      </template>
    </el-dialog>

    <!-- 预览弹窗 -->
    <el-dialog
      v-model="previewVisible"
      :title="previewData.title"
      width="500px"
      :close-on-click-modal="previewData.mask_closable === 1"
    >
      <div v-html="previewData.content"></div>
      <template #footer v-if="previewData.buttons && previewData.buttons.length > 0">
        <div>
          <el-button
            v-for="(button, index) in previewData.buttons"
            :key="index"
            :type="button.type || 'default'"
            @click="handlePreviewButtonClick(button)"
          >
            {{ button.text }}
          </el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Refresh, Plus } from '@element-plus/icons-vue'
import {
  getPopupList,
  getPopupInfo,
  createPopup,
  updatePopup,
  deletePopup,
  updatePopupStatus
} from '@/admin/api/system/popup'

// 响应式数据
const loading = ref(false)
const popupList = ref([])
const total = ref(0)
const dialogVisible = ref(false)
const dialogTitle = ref('')
const previewVisible = ref(false)
const previewData = ref({})
const formRef = ref()

// 查询参数
const queryParams = reactive({
  page: 1,
  page_size: 10,
  title: '',
  trigger_type: '',
  status: undefined
})

// 表单数据
const form = reactive({
  id: undefined,
  title: '',
  content: '',
  trigger_type: 'manual',
  priority: 0,
  auto_close: 0,
  delay: 0,
  mask_closable: 1,
  start_time: '',
  end_time: '',
  status: 1,
  title_style: null,
  buttons: null,
  trigger_rule: null
})

// JSON字符串形式的字段
const titleStyleStr = ref('')
const buttonsStr = ref('')
const triggerRuleStr = ref('')

// 表单验证规则
const rules = {
  title: [
    { required: false, message: '请输入弹窗标题', trigger: 'blur' },
    { max: 100, message: '标题长度不能超过100个字符', trigger: 'blur' }
  ],
  content: [
    { required: true, message: '请输入弹窗内容', trigger: 'blur' }
  ],
  trigger_type: [
    { required: true, message: '请选择触发类型', trigger: 'change' }
  ]
}

// 获取列表
const getList = async () => {
  loading.value = true
  try {
    const res = await getPopupList(queryParams)
    popupList.value = res.list || []
    total.value = res.total || 0
  } catch (error) {
    console.error('获取弹窗列表失败:', error)
  } finally {
    loading.value = false
  }
}

// 重置查询
const resetQuery = () => {
  Object.assign(queryParams, {
    page: 1,
    page_size: 10,
    title: '',
    trigger_type: '',
    status: undefined
  })
  getList()
}

// 新增
const handleAdd = () => {
  resetForm()
  dialogTitle.value = '新增弹窗'
  dialogVisible.value = true
}

// 编辑
const handleEdit = async (row) => {
  try {
    const res = await getPopupInfo(row.id)
    Object.assign(form, res)
    
    // 处理JSON字段
    titleStyleStr.value = res.title_style ? JSON.stringify(res.title_style, null, 2) : ''
    buttonsStr.value = res.buttons ? JSON.stringify(res.buttons, null, 2) : ''
    triggerRuleStr.value = res.trigger_rule ? JSON.stringify(res.trigger_rule, null, 2) : ''
    
    dialogTitle.value = '编辑弹窗'
    dialogVisible.value = true
  } catch (error) {
    console.error('获取弹窗详情失败:', error)
  }
}

// 删除
const handleDelete = (row) => {
  ElMessageBox.confirm(`确认要删除弹窗"${row.title}"吗？`, '警告', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await deletePopup(row.id)
      ElMessage.success('删除成功')
      getList()
    } catch (error) {
      console.error('删除失败:', error)
    }
  })
}

// 状态切换
const handleStatusChange = async (row) => {
  try {
    await updatePopupStatus({
      id: row.id,
      status: row.status
    })
    ElMessage.success('状态更新成功')
  } catch (error) {
    console.error('状态更新失败:', error)
    // 回滚状态
    row.status = row.status === 1 ? 0 : 1
  }
}

// 预览
const handlePreview = (row) => {
  previewData.value = { ...row }
  previewVisible.value = true
}

// 预览按钮点击
const handlePreviewButtonClick = (button) => {
  ElMessage.info(`点击了按钮: ${button.text}`)
  if (button.action === 'cancel') {
    previewVisible.value = false
  }
}

// 提交表单
const submitForm = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (valid) {
      try {
        // 处理JSON字段
        const submitData = { ...form }
        
        // 解析JSON字符串
        if (titleStyleStr.value.trim()) {
          try {
            submitData.title_style = JSON.parse(titleStyleStr.value)
          } catch (e) {
            ElMessage.error('标题样式JSON格式错误')
            return
          }
        }
        
        if (buttonsStr.value.trim()) {
          try {
            submitData.buttons = JSON.parse(buttonsStr.value)
          } catch (e) {
            ElMessage.error('按钮配置JSON格式错误')
            return
          }
        }
        
        if (triggerRuleStr.value.trim()) {
          try {
            submitData.trigger_rule = JSON.parse(triggerRuleStr.value)
          } catch (e) {
            ElMessage.error('触发规则JSON格式错误')
            return
          }
        }

        if (form.id) {
          await updatePopup(submitData)
          ElMessage.success('更新成功')
        } else {
          await createPopup(submitData)
          ElMessage.success('创建成功')
        }
        
        dialogVisible.value = false
        getList()
      } catch (error) {
        console.error('保存失败:', error)
      }
    }
  })
}

// 重置表单
const resetForm = () => {
  Object.assign(form, {
    id: undefined,
    title: '',
    content: '',
    trigger_type: 'manual',
    priority: 0,
    auto_close: 0,
    delay: 0,
    mask_closable: 1,
    start_time: '',
    end_time: '',
    status: 1,
    title_style: null,
    buttons: null,
    trigger_rule: null
  })
  titleStyleStr.value = ''
  buttonsStr.value = ''
  triggerRuleStr.value = ''
  formRef.value?.resetFields()
}

// 工具函数
const getTriggerTypeName = (type) => {
  const map = {
    manual: '手动触发',
    auto: '自动触发',
    schedule: '定时触发'
  }
  return map[type] || type
}

const getTriggerTypeTag = (type) => {
  const map = {
    manual: '',
    auto: 'success',
    schedule: 'warning'
  }
  return map[type] || ''
}

const formatDate = (date) => {
  if (!date) return '-'
  return new Date(date).toLocaleString('zh-CN')
}

// 初始化
onMounted(() => {
  getList()
})
</script>

<style scoped>
.popup-management {
  padding: 20px;
}

.search-card {
  margin-bottom: 20px;
}

.table-card {
  margin-bottom: 20px;
}

.table-header {
  margin-bottom: 20px;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

:deep(.el-dialog__body) {
  max-height: 70vh;
  overflow-y: auto;
}
</style> 