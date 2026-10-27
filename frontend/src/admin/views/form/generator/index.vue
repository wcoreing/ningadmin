<template>
  <div class="form-generator">
    <el-container style="height: 100%">
      <component-list @clone="handleClone" />
      <form-canvas
        v-model:form-items="formItems"
        v-model:current-item="currentItem"
        :form-data="formData"
        @preview="previewVisible = true"
        @generate="codeVisible = true"
        class="form-canvas-container"
      >
        <template #header-actions>
          <div class="template-actions">
            <el-button-group>
              <el-button type="primary" @click="handleAdd">
                <el-icon><Plus /></el-icon>
                新增模板
              </el-button>
           
              <el-button type="success" @click="handleSave">
                <el-icon><Document /></el-icon>
                保存模板
              </el-button>
              <el-button type="primary" @click="openTemplateDialog">
                <el-icon><List /></el-icon>
                模板管理
              </el-button>
            </el-button-group>
          </div>
        </template>
      </form-canvas>
      <props-panel
        v-model:current-item="currentItem"
        v-if="currentItem"
      />
    </el-container>

    <preview-dialog
      v-model:visible="previewVisible"
      :form-items="formItems"
      :form-data="formData"
      @submit="handleSubmit"
    />

    <code-dialog
      v-model:visible="codeVisible"
      :form-items="formItems"
    />

    <el-dialog
      v-model="saveDialogVisible"
      title="保存表单模板"
      width="500px"
    >
      <el-form
        ref="saveFormRef"
        :model="saveForm"
        :rules="saveFormRules"
        label-width="100px"
      >
        <el-form-item label="模板名称" prop="name">
          <el-input v-model="saveForm.name" placeholder="请输入模板名称" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input
            v-model="saveForm.description"
            type="textarea"
            placeholder="请输入模板描述"
          />
        </el-form-item>
        <el-divider>接口配置</el-divider>
        <el-form-item label="提交地址" prop="submitUrl">
          <el-input v-model="saveForm.submitUrl" placeholder="请输入接口地址" />
        </el-form-item>
        <el-form-item label="提交方法" prop="submitMethod">
          <el-select v-model="saveForm.submitMethod" style="width: 100%">
            <el-option label="POST" value="POST" />
            <el-option label="PUT" value="PUT" />
            <el-option label="GET" value="GET" />
          </el-select>
        </el-form-item>
        <el-form-item label="成功提示" prop="successMessage">
          <el-input v-model="saveForm.successMessage" placeholder="请输入成功提示语" />
        </el-form-item>
        <el-form-item label="失败提示" prop="failMessage">
          <el-input v-model="saveForm.failMessage" placeholder="请输入失败提示语" />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="saveForm.status">
            <el-radio label="0">正常</el-radio>
            <el-radio label="1">停用</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="saveDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="saveTemplate">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog
      v-model="templateDialogVisible"
      title="模板管理"
      width="1000px"
      destroy-on-close
      @closed="handleTemplateDialogClosed"
    >
      <div class="template-header">
        <el-input
          v-model="templateSearch"
          placeholder="搜索模板名称"
          style="width: 200px"
          clearable
          @keyup.enter="handleTemplateSearch"
          @clear="handleTemplateSearch"
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
      </div>
      
      <el-table
        v-loading="templateLoading"
        :data="templates"
        style="width: 100%"
        border
      >
        <template #empty>
          <el-empty description="暂无数据" />
        </template>
        <el-table-column prop="name" label="模板名称" min-width="120" />
        <el-table-column prop="description" label="描述" min-width="150" show-overflow-tooltip />
      
        <el-table-column prop="status" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === '0' ? 'success' : 'info'" size="small">
              {{ row.status === '0' ? '正常' : '停用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="useTemplate(row)">
              使用
            </el-button>
            <el-button type="success" link @click="editTemplate(row)">
              编辑
            </el-button>
            <el-button type="danger" link @click="deleteTemplate(row)">
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      
      <el-pagination
        v-model:current-page="templateCurrentPage"
        v-model:page-size="templatePageSize"
        :total="total"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next"
        class="pagination"
        @size-change="handleTemplatePageSizeChange"
        @current-change="handleTemplatePageChange"
      />
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Document, Plus, List, Search } from '@element-plus/icons-vue'
import { formatTime } from '@/admin/utils/date'
import ComponentList from '@/admin/components/form-generator/ComponentList.vue'
import FormCanvas from '@/admin/components/form-generator/FormCanvas.vue'
import PropsPanel from '@/admin/components/form-generator/PropsPanel.vue'
import PreviewDialog from '@/admin/components/form-generator/PreviewDialog.vue'
import CodeDialog from '@/admin/components/form-generator/CodeDialog.vue'
import { getFormTemplateList, saveFormTemplate, deleteFormTemplate } from '@/admin/api/form'

const formItems = ref([])
const formData = reactive({})
const currentItem = ref(null)
const previewVisible = ref(false)
const codeVisible = ref(false)

const templates = ref([])
const templateDialogVisible = ref(false)
const templateLoading = ref(false)
const templateSearch = ref('')
const templateCurrentPage = ref(1)
const templatePageSize = ref(10)
const total = ref(0)

const saveDialogVisible = ref(false)
const saveFormRef = ref(null)
const saveForm = reactive({
  name: '',
  description: '',
  submitUrl: '',
  submitMethod: 'POST',
  successMessage: '提交成功',
  failMessage: '提交失败',
  status: '0'
})

const saveFormRules = {
  name: [{ required: true, message: '请输入模板名称', trigger: 'blur' }],
  submitUrl: [{ required: true, message: '请输入提交地址', trigger: 'blur' }],
  submitMethod: [{ required: true, message: '请选择提交方法', trigger: 'change' }]
}

const currentTemplateId = ref('')

const getTemplates = async () => {
  console.log('开始获取模板列表')
  templateLoading.value = true
  try {
    const params = {
      name: templateSearch.value || undefined,
      pageNum: templateCurrentPage.value,
      pageSize: templatePageSize.value
    }
    console.log('请求参数:', params)
    const res = await getFormTemplateList(params)
    console.log('获取模板列表响应:', res)
    
    // 处理模板列表数据
    templates.value = (res.list || []).map(item => {
      // 如果form_items是加密字符串，需要先解码
      let formItemsStr = item.form_items || '[]'
      try {
        // 尝试解码Base64
        formItemsStr = decodeURIComponent(escape(atob(formItemsStr)))
      } catch (e) {
        console.warn('form_items解码失败，尝试直接解析:', e)
      }
      
      // 尝试解析JSON
      let formItems = []
      try {
        formItems = JSON.parse(formItemsStr)
      } catch (e) {
        console.error('form_items JSON解析失败:', e)
      }
      
      return {
        ...item,
        formItems
      }
    })
    
    total.value = res.total || 0
    console.log('设置模板列表数据:', templates.value)
    return true
  } catch (error) {
    console.error('获取模板列表失败:', error)
    ElMessage.error('获取模板列表失败')
    return false
  } finally {
    templateLoading.value = false
  }
}

const useTemplate = async (template) => {
  if (template) {
    try {
      console.log('使用模板数据:', template)
      console.log('模板formItems:', template.formItems)
      
      // 确保每个formItem都有必要的属性
      const processedFormItems = (template.formItems || []).map(item => ({
        type: item.type || 'el-input',
        key: item.key || `field_${Date.now()}`,
        props: {
          label: item.props?.label || '未命名',
          placeholder: item.props?.placeholder || '请输入',
          ...item.props
        },
        rules: item.rules || [],
        ...item
      }))
      
      formItems.value = processedFormItems
      console.log('处理后的formItems:', formItems.value)
      
      // 清空并初始化表单数据
      Object.keys(formData).forEach(key => delete formData[key])
      formItems.value.forEach(item => {
        if (item.type === 'el-upload' || item.type === 'el-checkbox-group') {
          formData[item.key] = []
        } else if (item.type === 'el-date-picker') {
          formData[item.key] = formatTime(new Date())
        } else {
          formData[item.key] = ''
        }
      })
      console.log('初始化后的formData:', formData)
      
      // 保存模板的基本信息
      currentTemplateId.value = template.id
      Object.assign(saveForm, {
        name: template.name,
        description: template.description,
        submitUrl: template.submit_url || '',
        submitMethod: template.submit_method || 'POST',
        successMessage: template.success_message || '提交成功',
        failMessage: template.fail_message || '提交失败',
        status: template.status || '0'
      })
      
      templateDialogVisible.value = false
      ElMessage.success('模板加载成功')
    } catch (error) {
      console.error('加载模板失败:', error)
      ElMessage.error('加载模板失败')
    }
  }
}

const editTemplate = (template) => {
  try {
    console.log('编辑模板数据:', template)
    console.log('模板formItems:', template.formItems)
    
    // 确保每个formItem都有必要的属性
    const processedFormItems = (template.formItems || []).map(item => ({
      type: item.type || 'el-input',
      key: item.key || `field_${Date.now()}`,
      props: {
        label: item.props?.label || '未命名',
        placeholder: item.props?.placeholder || '请输入',
        ...item.props
      },
      rules: item.rules || [],
      ...item
    }))
    
    formItems.value = processedFormItems
    console.log('处理后的formItems:', formItems.value)
    
    currentTemplateId.value = template.id
    Object.assign(saveForm, {
      name: template.name,
      description: template.description,
      submitUrl: template.submit_url || '',
      submitMethod: template.submit_method || 'POST',
      successMessage: template.success_message || '提交成功',
      failMessage: template.fail_message || '提交失败',
      status: template.status || '0'
    })
    saveDialogVisible.value = true
    templateDialogVisible.value = false
  } catch (error) {
    console.error('加载模板失败:', error)
    ElMessage.error('加载模板失败')
  }
}

const deleteTemplate = async (template) => {
  try {
    await ElMessageBox.confirm(
      `确认要删除模板"${template.name}"吗？`,
      '警告',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await deleteFormTemplate(template.id)
    if (template.id === currentTemplateId.value) {
      clearCurrentTemplate()
    }
    ElMessage.success('删除成功')
    getTemplates()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除模板失败:', error)
      ElMessage.error('删除模板失败')
    }
  }
}

const handleTemplateSearch = () => {
  templateCurrentPage.value = 1
  getTemplates()
}

const handleTemplatePageSizeChange = (val) => {
  templatePageSize.value = val
  templateCurrentPage.value = 1
  getTemplates()
}

const handleTemplatePageChange = (val) => {
  templateCurrentPage.value = val
  getTemplates()
}

const handleSave = () => {
  if (formItems.value.length === 0) {
    ElMessage.warning('请先添加表单项')
    return
  }
  
  // 如果不是编辑现有模板，则重置表单
  if (!currentTemplateId.value) {
    Object.assign(saveForm, {
      name: '',
      description: '',
      submitUrl: '',
      submitMethod: 'POST',
      successMessage: '提交成功',
      failMessage: '提交失败',
      status: '0'
    })
  }
  
  saveDialogVisible.value = true
}

const saveTemplate = async () => {
  if (!saveFormRef.value) return
  
  await saveFormRef.value.validate(async (valid) => {
    if (valid) {
      try {
        const templateData = {
          id: currentTemplateId.value,
          name: saveForm.name,
          description: saveForm.description,
          submit_url: saveForm.submitUrl,
          submit_method: saveForm.submitMethod,
          success_message: saveForm.successMessage,
          fail_message: saveForm.failMessage,
          form_items: formItems.value,
          status: saveForm.status
        }
        
        const res = await saveFormTemplate(templateData)
        currentTemplateId.value = res.id
        ElMessage.success(currentTemplateId.value ? '更新成功' : '保存成功')
        saveDialogVisible.value = false
        await getTemplates()
      } catch (error) {
        console.error('保存模板失败:', error)
        ElMessage.error('保存模板失败')
      }
    }
  })
}

const clearCurrentTemplate = () => {
  currentTemplateId.value = ''
  formItems.value = []
  Object.keys(formData).forEach(key => delete formData[key])
}

const handleClone = (cloned) => {
  if (formItems.value.some(item => item.key === cloned.key)) {
    return
  }

  if (cloned.type === 'el-upload' || cloned.type === 'el-checkbox-group') {
    formData[cloned.key] = []
  } else if (cloned.type === 'el-date-picker') {
    formData[cloned.key] = formatTime(new Date())
  } else {
    formData[cloned.key] = ''
  }
  formItems.value.push(cloned)
}

const handleSubmit = () => {
  console.log('表单数据:', formData)
  previewVisible.value = false
}

const openTemplateDialog = async () => {
  templateDialogVisible.value = true
  await getTemplates()
}

const handleTemplateDialogClosed = () => {
  templateSearch.value = ''
  templateCurrentPage.value = 1
  templates.value = []
  total.value = 0
}

const handleAdd = () => {
  // 清空当前模板ID
  currentTemplateId.value = ''
  
  // 清空表单项
  formItems.value = []
  
  // 清空表单数据
  Object.keys(formData).forEach(key => delete formData[key])
  
  // 重置保存表单的数据
  Object.assign(saveForm, {
    name: '',
    description: '',
    submitUrl: '',
    submitMethod: 'POST',
    successMessage: '提交成功',
    failMessage: '提交失败',
    status: '0'
  })
  
  // 清空当前选中的表单项
  currentItem.value = null
  
  ElMessage.success('已切换到新表单')
}

onMounted(() => {
  getTemplates()
})
</script>

<style scoped>
.form-generator {
  height: 100%;
}

.form-canvas-container {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.template-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.el-select :deep(.el-select-dropdown__item) {
  display: flex;
  align-items: center;
}

.template-header {
  margin-bottom: 16px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style> 