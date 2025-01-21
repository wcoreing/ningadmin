<template>
  <div class="form-list">
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
        <el-form-item label="状态" prop="formStatus">
          <el-select v-model="queryParams.formStatus" placeholder="表单状态" clearable>
            <el-option label="正常" value="0" />
            <el-option label="停用" value="1" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :icon="Search" @click="handleQuery">搜索</el-button>
          <el-button :icon="Refresh" @click="resetQuery">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 操作按钮区域 -->
    <el-card class="table-card">
      <template #header>
        <el-button type="primary" :icon="Plus" @click="handleAdd">新增</el-button>
      </template>

      <!-- 表格区域 -->
      <el-table v-loading="loading" :data="formList">
        <el-table-column type="index" label="序号" width="60" />
        <el-table-column prop="formName" label="表单名称" min-width="120" show-overflow-tooltip />
        <el-table-column prop="formDescription" label="表单描述" min-width="150" show-overflow-tooltip />
        <el-table-column prop="templateName" label="模板名称" min-width="120" show-overflow-tooltip />
       
        <el-table-column prop="formStatus" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.formStatus === '0' ? 'success' : 'info'" size="small">
              {{ row.formStatus === '0' ? '正常' : '停用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="formCreateTime" label="创建时间" width="160" />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
            <el-button type="success" link @click="handleView(row)">查看</el-button>
            <el-button 
              type="danger" 
              link 
              @click="handleDelete(row)"
              v-if="row.formStatus !== '1'"
            >删除</el-button>
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

    <!-- 表单对话框 -->
    <el-dialog
      v-model="dialog.visible"
      :title="dialog.title"
      width="600px"
      append-to-body
      @closed="handleDialogClosed"
    >
      <el-form
        ref="formRef"
        :model="formData"
        :rules="rules"
        label-width="100px"
      >
        <el-form-item label="表单模板" prop="templateId">
          <el-select v-model="formData.templateId" placeholder="请选择表单模板" style="width: 100%">
            <el-option
              v-for="item in templateList"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            >
              <span>{{ item.name }}</span>
              <span style="color: #999; font-size: 13px; margin-left: 8px">{{ item.description }}</span>
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="表单名称" prop="formName">
          <el-input v-model="formData.formName" placeholder="请输入表单名称" />
        </el-form-item>
        <el-form-item label="表单描述" prop="formDescription">
          <el-input
            v-model="formData.formDescription"
            type="textarea"
            placeholder="请输入表单描述"
          />
        </el-form-item>
        <el-divider>表单属性</el-divider>
        <el-form-item label="主题">
          <el-select v-model="formData.formAttributes.theme" style="width: 100%">
            <el-option label="默认主题" value="default" />
            <el-option label="暗色主题" value="dark" />
            <el-option label="浅色主题" value="light" />
          </el-select>
        </el-form-item>
        <el-form-item label="字体大小">
          <el-select v-model="formData.formAttributes.fontSize" style="width: 100%">
            <el-option label="小" value="12px" />
            <el-option label="中" value="14px" />
            <el-option label="大" value="16px" />
          </el-select>
        </el-form-item>
        <el-form-item label="表单宽度">
          <el-select v-model="formData.formAttributes.width" style="width: 100%">
            <el-option label="100%" value="100%" />
            <el-option label="75%" value="75%" />
            <el-option label="50%" value="50%" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态" prop="formStatus">
          <el-radio-group v-model="formData.formStatus">
            <el-radio label="0">正常</el-radio>
            <el-radio label="1">停用</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialog.visible = false">取消</el-button>
        <el-button type="primary" @click="submitForm">确定</el-button>
      </template>
    </el-dialog>

    <!-- 添加预览对话框 -->
    <el-dialog
      v-model="viewDialog.visible"
      :width="viewDialog.width"
      append-to-body
      destroy-on-close
      :close-on-click-modal="false"
      class="form-preview-dialog"
    >
      <template #header>
        <div class="dialog-header">
          <span class="dialog-title">表单详情</span>
          <el-button
            type="primary"
            link
            @click="handleOpenInNewWindow"
          >
            <el-icon><Share /></el-icon>
            新窗口打开
          </el-button>
        </div>
      </template>

      <form-view 
        v-if="viewDialog.visible" 
        :form-id="viewDialog.formId" 
        @submit-success="handleViewSubmitSuccess"
      />
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Refresh, Plus, Share } from '@element-plus/icons-vue'
import { useRouter } from 'vue-router'
import { getFormTemplateList, addForm, updateForm, deleteForm, getFormList, getFormTemplateDetail, getFormDetail } from '@/api/form'
import FormView from '@/views/form/view/index.vue'

const router = useRouter()

// 查询参数
const queryParams = reactive({
  pageNum: 1,
  pageSize: 10,
  formName: '',
  formStatus: ''
})

// 表单数据
const formList = ref([])
const total = ref(0)
const loading = ref(false)

// 对话框控制
const dialog = reactive({
  visible: false,
  title: ''
})

// 添加模板列表数据
const templateList = ref([])

// 表单对象
const formData = reactive({
  id: undefined,
  templateId: undefined,
  formName: '',
  formDescription: '',
  formAttributes: {
    theme: 'default',
    fontSize: '14px',
    width: '100%'
  },
  formStatus: '0'
})

// 表单校验规则
const rules = {
  templateId: [
    { required: true, message: '请选择表单模板', trigger: 'change' }
  ],
  formName: [
    { required: true, message: '请输入表单名称', trigger: 'blur' },
    { min: 2, max: 100, message: '长度在 2 到 100 个字符', trigger: 'blur' }
  ],
  formDescription: [
    { max: 500, message: '长度不能超过 500 个字符', trigger: 'blur' }
  ],
  formStatus: [
    { required: true, message: '请选择状态', trigger: 'change' }
  ]
}

// 引用对象
const queryRef = ref()
const formRef = ref()

// 添加预览对话框控制
const viewDialog = reactive({
  visible: false,
  formId: null,
  width: '100%'  // 添加宽度控制
})

// 获取表单列表
const getList = async () => {
  loading.value = true
  try {
    const res = await getFormList(queryParams)
    // 直接使用返回的数据，不进行 JSON 解析
    formList.value = res.list || []
    total.value = res.total || 0
  } catch (error) {
    console.error('获取表单列表失败:', error)
    ElMessage.error('获取表单列表失败')
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

// 获取模板列表
const getTemplateList = async () => {
  try {
    const res = await getFormTemplateList({
      pageSize: 100,
      status: '0' // 只获取正常状态的模板
    })
    templateList.value = res.list || []
  } catch (error) {
    console.error('获取模板列表失败:', error)
  }
}

// 新增表单
const handleAdd = () => {
  dialog.title = '新增表单'
  dialog.visible = true
  Object.assign(formData, {
    id: undefined,
    templateId: undefined,
    formName: '',
    formDescription: '',
    formAttributes: {
      theme: 'default',
      fontSize: '14px',
      width: '100%'
    },
    formStatus: '0'
  })
  // 获取模板列表
  getTemplateList()
}

// 修改表单
const handleEdit = (row) => {
  dialog.title = '修改表单'
  dialog.visible = true
  
  // 先重置表单数据
  formData.id = undefined
  formData.templateId = undefined
  formData.formName = ''
  formData.formDescription = ''
  formData.formStatus = '0'
  formData.formAttributes = {
    theme: 'default',
    fontSize: '14px',
    width: '100%'
  }

  // 直接使用表单属性
  formData.id = row.id
  formData.templateId = row.templateId
  formData.formName = row.formName
  formData.formDescription = row.formDescription
  formData.formStatus = row.formStatus
  formData.formAttributes = row.formAttributes

  // 获取模板列表
  getTemplateList()
}

// 查看表单
const handleView = async (row) => {
  try {
    // 获取表单详情
    const detail = await getFormDetail(row.id)
    // 设置弹窗宽度为表单配置的宽度
    let formWidth = detail.formAttributes?.width || '100%'
    // 如果是百分比，转换为具体的像素值
    if (formWidth.endsWith('%')) {
      const percent = parseInt(formWidth) / 100
      viewDialog.width = `${Math.min(1200, window.innerWidth * percent)}px`
    } else {
      viewDialog.width = formWidth
    }
    viewDialog.formId = row.id
    viewDialog.visible = true
  } catch (error) {
    console.error('获取表单详情失败:', error)
    ElMessage.error('获取表单详情失败')
  }
}

// 删除表单
const handleDelete = async (row) => {
  ElMessageBox.confirm(
    `确认要删除表单"${row.formName}"吗？`,
    '警告',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(async () => {
    try {
      await deleteForm(row.id)
      ElMessage.success('删除成功')
      getList()
    } catch (error) {
      console.error('删除表单失败:', error)
    }
  }).catch(() => {})
}

// 提交表单
const submitForm = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (valid) {
      try {
        // 直接提交数据，不进行 JSON 转换
        const submitData = {
          ...formData
        }

        if (formData.id) {
          await updateForm(submitData)
          ElMessage.success('修改成功')
        } else {
          await addForm(submitData)
          ElMessage.success('新增成功')
        }
        dialog.visible = false
        getList()
      } catch (error) {
        console.error('保存表单失败:', error)
        ElMessage.error('保存失败：' + (error.message || '未知错误'))
      }
    }
  })
}

// 处理对话框关闭
const handleDialogClosed = () => {
  formRef.value?.resetFields()
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

// 处理表单预览提交成功
const handleViewSubmitSuccess = () => {
  viewDialog.visible = false
  getList() // 刷新列表
}

// 在新窗口打开
const handleOpenInNewWindow = () => {
  window.open(`/form/view/${viewDialog.formId}`, '_blank')
}

onMounted(() => {
  getList()
})
</script>

<style scoped>
.form-list {
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

/* 添加预览对话框样式 */
:deep(.form-preview-dialog) {
  .el-dialog__body {
    padding: 0;
    max-height: calc(90vh - 120px);
    overflow-y: auto;
  }

  .dialog-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 20px;
  }

  .dialog-title {
    font-size: 16px;
    font-weight: 500;
    color: var(--el-text-color-primary);
  }
}
</style>
