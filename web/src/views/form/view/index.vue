<template>
  <div class="form-view">
    <div v-loading="loading">
      <div v-if="formDetail && !loading" 
        class="preview-form"
        :class="{
          'theme-dark': formDetail.formAttributes?.theme === 'dark',
          'theme-light': formDetail.formAttributes?.theme === 'light',
          'theme-default': !formDetail.formAttributes?.theme || formDetail.formAttributes?.theme === 'default',
          'is-dialog': isDialog
        }"
      >
        <div class="form-header">
          <h3 class="form-title">{{ formDetail.formName }}</h3>
          <p v-if="formDetail.formDescription" class="form-description">{{ formDetail.formDescription }}</p>
        </div>

        <el-form
          ref="formRef"
          :model="formData"
          label-width="120px"
          :validate-on-rule-change="false"
          :style="{
            width: formDetail.formAttributes?.width || '100%',
            fontSize: formDetail.formAttributes?.fontSize || '14px'
          }"
        >
          <template v-for="(item, index) in formItems" :key="index">
            <!-- 输入框 -->
            <el-form-item 
              v-if="item.type === 'input'"
              :label="item.label"
              :prop="item.prop"
              :rules="item.required ? { required: true, message: `请输入${item.label}`, trigger: 'blur' } : []"
            >
              <el-input 
                v-model="formData[item.prop]"
                :placeholder="item.placeholder"
              />
            </el-form-item>

            <!-- 开关 -->
            <el-form-item 
              v-if="item.type === 'switch'"
              :label="item.label"
              :prop="item.prop"
              :rules="item.required ? { required: true, message: `请选择${item.label}`, trigger: 'change' } : []"
            >
              <el-switch
                v-model="formData[item.prop]"
              />
            </el-form-item>

            <!-- 文本域 -->
            <el-form-item 
              v-if="item.type === 'textarea'"
              :label="item.label"
              :prop="item.prop"
              :rules="item.required ? { required: true, message: `请输入${item.label}`, trigger: 'blur' } : []"
            >
              <el-input 
                v-model="formData[item.prop]"
                type="textarea"
                :rows="3"
                :placeholder="item.placeholder"
              />
            </el-form-item>

            <!-- 单选框 -->
            <el-form-item 
              v-if="item.type === 'radio'"
              :label="item.label"
              :prop="item.prop"
              :rules="item.required ? { required: true, message: `请选择${item.label}`, trigger: 'change' } : []"
            >
              <el-radio-group 
                v-model="formData[item.prop]"
              >
                <el-radio 
                  v-for="option in item.options"
                  :key="option.value"
                  :label="option.value"
                >
                  {{ option.label }}
                </el-radio>
              </el-radio-group>
            </el-form-item>

            <!-- 复选框 -->
            <el-form-item 
              v-if="item.type === 'checkbox'"
              :label="item.label"
              :prop="item.prop"
              :rules="item.required ? { required: true, message: `请选择${item.label}`, trigger: 'change' } : []"
            >
              <el-checkbox-group 
                v-model="formData[item.prop]"
              >
                <el-checkbox 
                  v-for="option in item.options"
                  :key="option.value"
                  :label="option.value"
                >
                  {{ option.label }}
                </el-checkbox>
              </el-checkbox-group>
            </el-form-item>

            <!-- 下拉选择 -->
            <el-form-item 
              v-if="item.type === 'select'"
              :label="item.label"
              :prop="item.prop"
              :rules="item.required ? { required: true, message: `请选择${item.label}`, trigger: 'change' } : []"
            >
              <el-select 
                v-model="formData[item.prop]"
                :placeholder="item.placeholder"
                style="width: 100%"
              >
                <el-option
                  v-for="option in item.options"
                  :key="option.value"
                  :label="option.label"
                  :value="option.value"
                />
              </el-select>
            </el-form-item>

            <!-- 日期选择 -->
            <el-form-item 
              v-if="item.type === 'date'"
              :label="item.label"
              :prop="item.prop"
              :rules="item.required ? { required: true, message: `请选择${item.label}`, trigger: 'change' } : []"
            >
              <el-date-picker
                v-model="formData[item.prop]"
                :type="item.dateType"
                :format="item.format"
                :placeholder="item.placeholder"
                style="width: 100%"
                value-format="YYYY-MM-DD"
              />
            </el-form-item>

            <!-- 时间选择 -->
            <el-form-item 
              v-if="item.type === 'time'"
              :label="item.label"
              :prop="item.prop"
              :rules="item.required ? { required: true, message: `请选择${item.label}`, trigger: 'change' } : []"
            >
              <el-time-picker
                v-model="formData[item.prop]"
                :placeholder="item.placeholder || `请选择${item.label}`"
                style="width: 100%"
                value-format="HH:mm:ss"
              />
            </el-form-item>

            <!-- 日期时间选择 -->
            <el-form-item 
              v-if="item.type === 'datetime'"
              :label="item.label"
              :prop="item.prop"
              :rules="item.required ? { required: true, message: `请选择${item.label}`, trigger: 'change' } : []"
            >
              <el-date-picker
                v-model="formData[item.prop]"
                type="datetime"
                :placeholder="item.placeholder || `请选择${item.label}`"
                style="width: 100%"
                value-format="YYYY-MM-DD HH:mm:ss"
              />
            </el-form-item>

            <!-- 图片上传 -->
            <el-form-item 
              v-if="item.type === 'upload'"
              :label="item.label"
              :prop="item.prop"
              :rules="item.required ? { required: true, message: `请上传${item.label}`, trigger: 'change' } : []"
            >
              <el-upload
                v-model:file-list="formData[item.prop]"
                class="upload-demo"
                :action="item.action"
                :limit="item.limit"
                :multiple="item.multiple"
                :accept="item.accept"
                :list-type="item.listType"
                :auto-upload="item.autoUpload"
                :on-success="(response, file) => handleUploadSuccess(response, file, item.prop)"
                :on-error="(error, file) => handleUploadError(error, file)"
                :before-upload="handleBeforeUpload"
              >
                <el-icon><Plus /></el-icon>
                <template #tip>
                  <div class="el-upload__tip">
                    只能上传{{ item.accept.replace('image/*', '图片') }}文件，且不超过2MB
                  </div>
                </template>
              </el-upload>
            </el-form-item>
          </template>
        </el-form>

        <div class="form-actions">
          <el-button type="primary" @click="handleSubmit" :loading="submitting">提交</el-button>
        </div>
      </div>
      
      <el-empty v-else-if="!loading" description="表单不存在或已被删除" />
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { getFormDetail, getFormTemplateDetail } from '@/api/form'
import { submitFormData } from '@/api/form-data'
const props = defineProps({
  formId: {
    type: [String, Number],
    default: ''
  },
  isDialog: {
    type: Boolean,
    default: false
  }
})

const route = useRoute()
const loading = ref(false)
const submitting = ref(false)
const formDetail = ref(null)
const formItems = ref([])
const formData = reactive({})
const formRef = ref(null)

// 获取表单数据
const getFormData = async (id) => {
  if (!id) {
    ElMessage.error('表单ID不能为空')
    return
  }

  loading.value = true
  try {
    // 获取表单详情
    const detail = await getFormDetail(id)
    
    // 处理表单属性
    try {
      if (typeof detail.formAttributes === 'string') {
        detail.formAttributes = JSON.parse(detail.formAttributes)
      }
    } catch (e) {
      console.warn('解析表单属性失败，使用默认值', e)
      detail.formAttributes = {
        theme: 'default',
        fontSize: '14px',
        width: '100%'
      }
    }
    
    formDetail.value = detail
    
    // 获取表单模板详情
    const templateDetail = await getFormTemplateDetail(detail.templateId)
    
    // 处理表单项数据
    formItems.value = templateDetail.formItems.map(item => {
      // 处理选项数据
      let options = []
      if (item.type === 'el-select' || item.type === 'el-radio' || item.type === 'el-checkbox') {
        try {
          if (item.props && item.props.options) {
            options = item.props.options.map(opt => ({
              label: opt.label,
              value: opt.value
            }))
          }
        } catch (e) {
          console.warn('解析选项数据失败:', e)
          options = []
        }
      }

      // 处理日期选择器的格式
      let dateFormat = ''
      if (item.type === 'el-date-picker') {
        dateFormat = item.props?.format || 'YYYY-MM-DD'
      }

      // 映射组件类型
      const typeMap = {
        'el-input': item.props?.type === 'textarea' ? 'textarea' : 'input',
        'el-select': 'select',
        'el-date-picker': 'date',
        'el-radio': 'radio',
        'el-checkbox': 'checkbox',
        'el-switch': 'switch',
        'el-time-picker': 'time',
        'el-upload': 'upload'
      }

      // 处理上传组件的特殊配置
      let uploadConfig = {}
      if (item.type === 'el-upload') {
        uploadConfig = {
          action: item.props?.action || '/api/upload',
          limit: item.props?.limit || 1,
          multiple: item.props?.multiple || false,
          accept: item.props?.accept || 'image/*',
          listType: item.props?.listType || 'picture-card',
          autoUpload: item.props?.autoUpload !== false
        }
      }

      return {
        type: typeMap[item.type] || item.type,
        label: item.props?.label || item.label,
        prop: item.key,
        placeholder: item.props?.placeholder || `请输入${item.label}`,
        options: options,
        defaultValue: item.defaultValue,
        required: item.required || false,
        dateType: item.props?.type || 'date',
        format: dateFormat,
        ...uploadConfig
      }
    })
    
    // 初始化表单数据
    formItems.value.forEach(item => {
      if (item.prop) {
        // 如果有默认值且是字符串类型的JSON，尝试解析
        let defaultValue = item.defaultValue
        if (typeof defaultValue === 'string' && (item.type === 'select' || item.type === 'radio' || item.type === 'checkbox')) {
          try {
            defaultValue = JSON.parse(defaultValue)
          } catch (e) {
            console.warn('解析默认值失败:', e)
          }
        }

        switch (item.type) {
          case 'checkbox':
            formData[item.prop] = defaultValue || []
            break
          case 'switch':
            formData[item.prop] = defaultValue === 'true' || defaultValue === true || false
            break
          case 'radio':
          case 'select':
            formData[item.prop] = defaultValue || (item.options?.[0]?.value ?? '')
            break
          case 'date':
            formData[item.prop] = defaultValue || ''
            break
          case 'time':
            formData[item.prop] = defaultValue || ''
            break
          case 'datetime':
            formData[item.prop] = defaultValue || ''
            break
          case 'upload':
            // 处理上传文件的默认值
            if (defaultValue) {
              try {
                if (typeof defaultValue === 'string') {
                  defaultValue = JSON.parse(defaultValue)
                }
                formData[item.prop] = Array.isArray(defaultValue) ? defaultValue : [defaultValue]
              } catch (e) {
                console.warn('解析上传文件默认值失败:', e)
                formData[item.prop] = []
              }
            } else {
              formData[item.prop] = []
            }
            break
          default:
            formData[item.prop] = defaultValue || ''
        }
      }
    })
  } catch (error) {
    console.error('获取表单数据失败:', error)
    ElMessage.error('获取表单数据失败：' + error.message)
  } finally {
    loading.value = false
  }
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return
  
  try {
    await formRef.value.validate()
    submitting.value = true
    
    // 构建提交数据，包含字段名称
    const formDataWithLabels = {}
    formItems.value.forEach(item => {
      formDataWithLabels[item.prop] = {
        label: item.label,
        value: formData[item.prop]
      }
    })
    
    await submitFormData({
      formId: formDetail.value.id,
      data: formDataWithLabels,
      submitter: {
        name: '', // 这里可以从用户状态获取
        email: ''
      }
    })
    
    ElMessage.success('提交成功')
    // 触发提交成功事件
    emit('submit-success')
  } catch (error) {
    console.error('提交失败:', error)
    ElMessage.error('提交失败：' + error.message)
  } finally {
    submitting.value = false
  }
}

// 处理上传成功
const handleUploadSuccess = (response, file, prop) => {
  if (!formData[prop]) {
    formData[prop] = []
  }
  // 如果返回的数据中包含文件URL，使用返回的URL
  const fileUrl = response.data?.url || response.url || URL.createObjectURL(file.raw)
  formData[prop].push({
    name: file.name,
    url: fileUrl
  })
}

// 处理上传失败
const handleUploadError = (error, file) => {
  ElMessage.error('上传失败：' + (error.message || '未知错误'))
}

// 上传前检查
const handleBeforeUpload = (file) => {
  const isImage = file.type.startsWith('image/')
  const isLt2M = file.size / 1024 / 1024 < 2

  if (!isImage) {
    ElMessage.error('只能上传图片文件!')
    return false
  }
  if (!isLt2M) {
    ElMessage.error('图片大小不能超过 2MB!')
    return false
  }
  return true
}

// 监听 formId 变化
watch(() => props.formId, (newId) => {
  if (newId) {
    getFormData(newId)
  }
}, { immediate: true })

// 如果是独立页面，从路由获取ID
onMounted(() => {
  if (!props.formId && route.params.id) {
    getFormData(route.params.id)
  }
})

// 定义事件
const emit = defineEmits(['submit-success'])
</script>

<style>
/* 暗色主题基础样式 */
.preview-form.theme-dark {
  background-color: #1a1a1a;  /* 深灰背景 */
  color: #e5e5e5;            /* 浅灰文字，不用纯白色更护眼 */
  border: 1px solid #333333; /* 深色边框 */
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15); /* 添加阴影增加层次感 */
}

/* 表单标签和文本 */
.preview-form.theme-dark .el-form-item__label,
.preview-form.theme-dark .el-radio__label,
.preview-form.theme-dark .el-checkbox__label {
  color: #e5e5e5 !important;
}

/* 输入框样式 */
.preview-form.theme-dark .el-input__wrapper,
.preview-form.theme-dark .el-textarea__wrapper {
  background-color: #2a2a2a !important;  /* 比背景稍浅 */
  box-shadow: 0 0 0 1px #404040 inset !important;  /* 更柔和的边框 */
  transition: box-shadow 0.2s;
}

.preview-form.theme-dark .el-input__wrapper:hover,
.preview-form.theme-dark .el-textarea__wrapper:hover {
  box-shadow: 0 0 0 1px #505050 inset !important;  /* 悬停时边框变亮 */
}

.preview-form.theme-dark .el-input__wrapper.is-focus,
.preview-form.theme-dark .el-textarea__wrapper.is-focus {
  box-shadow: 0 0 0 1px #0960bd inset !important;  /* 聚焦时使用主题蓝 */
}

/* 输入框文本 */
.preview-form.theme-dark .el-input__inner,
.preview-form.theme-dark .el-textarea__inner {
  background-color: transparent !important;
  color: #e5e5e5 !important;
}

/* 单选框和复选框 */
.preview-form.theme-dark .el-radio__inner,
.preview-form.theme-dark .el-checkbox__inner {
  background-color: #2a2a2a !important;
  border-color: #404040 !important;
}

.preview-form.theme-dark .el-radio__input.is-checked .el-radio__inner,
.preview-form.theme-dark .el-checkbox__input.is-checked .el-checkbox__inner {
  background-color: #0960bd !important;
  border-color: #0960bd !important;
}

/* 按钮样式 */
.preview-form.theme-dark .el-button {
  background-color: #2a2a2a;
  border-color: #404040;
  color: #e5e5e5;
  transition: all 0.3s;
}

.preview-form.theme-dark .el-button:hover {
  background-color: #333333;
  border-color: #505050;
}

.preview-form.theme-dark .el-button--primary {
  background-color: #0960bd;  /* 更深的蓝色 */
  border-color: #0960bd;
  color: #ffffff;
}

.preview-form.theme-dark .el-button--primary:hover {
  background-color: #1976d2;  /* 悬停时颜色变亮 */
  border-color: #1976d2;
}

/* 下拉选择器 */
.preview-form.theme-dark .el-select .el-input__wrapper {
  background-color: #2a2a2a !important;
}

.el-select-dropdown.is-dark,
.el-picker__popper.is-dark {
  background-color: #2a2a2a !important;
  border-color: #404040 !important;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15) !important;
}

.el-select-dropdown.is-dark .el-select-dropdown__item {
  color: #e5e5e5 !important;
}

.el-select-dropdown.is-dark .el-select-dropdown__item.hover,
.el-select-dropdown.is-dark .el-select-dropdown__item:hover {
  background-color: #333333 !important;
}

.el-select-dropdown.is-dark .el-select-dropdown__item.selected {
  background-color: #0960bd !important;
  color: #ffffff !important;
}

/* 日期选择器 */
.el-picker__popper.is-dark .el-date-picker {
  background-color: #2a2a2a !important;
  border-color: #404040 !important;
  color: #e5e5e5 !important;
}

.el-picker__popper.is-dark .el-date-picker__header {
  color: #e5e5e5 !important;
}

.el-picker__popper.is-dark .el-date-table th {
  color: #808080 !important;  /* 表头使用灰色 */
}

.el-picker__popper.is-dark .el-date-table td.next-month .el-date-table-cell__text,
.el-picker__popper.is-dark .el-date-table td.prev-month .el-date-table-cell__text {
  color: #666666 !important;  /* 非当前月份日期颜色 */
}

.el-picker__popper.is-dark .el-date-table td.available:not(.current):not(.today) .el-date-table-cell__text {
  color: #e5e5e5 !important;
}

.el-picker__popper.is-dark .el-date-table td.today .el-date-table-cell__text {
  color: #0960bd !important;  /* 今天日期使用主题色 */
}

.el-picker__popper.is-dark .el-date-table td.current:not(.disabled) .el-date-table-cell__text {
  background-color: #0960bd !important;
  color: #ffffff !important;
}

.el-picker__popper.is-dark .el-date-table td.available:hover {
  color: #1976d2 !important;
}

/* 弹出层箭头 */
.el-popper.is-dark .el-popper__arrow::before {
  background-color: #2a2a2a !important;
  border-color: #404040 !important;
}

/* 表单验证消息 */
.preview-form.theme-dark .el-form-item__error {
  color: #ff4d4f !important;  /* 使用更鲜艳的红色 */
}

/* 表单标题和描述 */
.preview-form.theme-dark .form-title {
  color: #ffffff !important;  /* 标题使用纯白色 */
}

.preview-form.theme-dark .form-description {
  color: #999999 !important;  /* 描述文字使用灰色 */
}

.preview-form.theme-light {
  background-color: #ffffff;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.preview-form.theme-default {
  background-color: var(--el-bg-color-overlay);
  border: 1px solid var(--el-border-color);
}

/* 日期选择器样式 */
.preview-form.theme-dark .el-date-editor .el-input__wrapper {
  background-color: var(--el-bg-color-overlay) !important;
}

.preview-form.theme-dark .el-date-editor .el-input__inner {
  color: var(--el-text-color-primary) !important;
}

/* 下拉选择器样式 */
.preview-form.theme-dark .el-select .el-input__wrapper {
  background-color: #1d1e1f !important;
}

.preview-form.theme-dark .el-select .el-input__inner {
  color: #ffffff !important;
}

/* 下拉选项样式 */
.el-select-dropdown.is-dark,
.el-picker__popper.is-dark {
  background-color: #1d1e1f !important;
  border-color: #363637 !important;
}

.el-select-dropdown.is-dark .el-select-dropdown__item,
.el-picker__popper.is-dark .el-date-picker {
  color: #ffffff !important;
}

.el-select-dropdown.is-dark .el-select-dropdown__item.hover,
.el-select-dropdown.is-dark .el-select-dropdown__item:hover {
  background-color: #18222c !important;
}

.el-select-dropdown.is-dark .el-select-dropdown__item.selected {
  background-color: #409eff !important;
  color: #ffffff !important;
}

/* 日期选择器样式 */
.el-picker__popper.is-dark .el-date-picker {
  background-color: #1d1e1f !important;
}

.el-picker__popper.is-dark .el-date-picker__header {
  color: #ffffff !important;
}

.el-picker__popper.is-dark .el-date-table th {
  color: #909399 !important;
}

.el-picker__popper.is-dark .el-date-table td.next-month .el-date-table-cell__text,
.el-picker__popper.is-dark .el-date-table td.prev-month .el-date-table-cell__text {
  color: #606266 !important;
}

.el-picker__popper.is-dark .el-date-table td.available:not(.current):not(.today) .el-date-table-cell__text {
  color: #ffffff !important;
}

.el-picker__popper.is-dark .el-date-table td.today .el-date-table-cell__text {
  color: #409eff !important;
}

.el-picker__popper.is-dark .el-date-table td.current:not(.disabled) .el-date-table-cell__text {
  background-color: #409eff !important;
  color: #ffffff !important;
}

.el-picker__popper.is-dark .el-date-table td.available:hover {
  color: #409eff !important;
}

/* 弹出层箭头样式 */
.el-popper.is-dark .el-popper__arrow::before {
  background-color: #1d1e1f !important;
  border-color: #363637 !important;
}

/* 表单项验证消息 */
.preview-form.theme-dark .el-form-item__error {
  color: #f56c6c !important;
}

/* 上传组件样式 */
.preview-form.theme-dark .el-upload {
  border-color: #404040 !important;
  background-color: #2a2a2a !important;
}

.preview-form.theme-dark .el-upload:hover {
  border-color: #0960bd !important;
}

.preview-form.theme-dark .el-upload-list__item {
  background-color: #2a2a2a !important;
  border-color: #404040 !important;
  color: #e5e5e5 !important;
}

.preview-form.theme-dark .el-upload__tip {
  color: #909399 !important;
}
</style>

<style scoped>
.form-view {
  padding: 0;
}

.form-header {
  margin-bottom: 24px;
  text-align: center;
}

.form-title {
  margin: 0 0 8px;
  font-size: 20px;
  font-weight: 500;
  color: var(--el-text-color-primary);
}

.form-description {
  margin: 0;
  font-size: 14px;
  color: var(--el-text-color-secondary);
  line-height: 1.5;
}

.preview-form {
  padding: 24px;
  border-radius: 4px;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.preview-form .el-form {
  margin: 0 auto;
}

.form-actions {
  margin-top: 32px;
  display: flex;
  justify-content: center;
  gap: 16px;
  width: 100%;
}

:deep(.el-form-item) {
  margin-bottom: 22px;
  width: 100%;
}

:deep(.el-form-item__label) {
  font-weight: 500;
  padding-right: 12px;
  line-height: 32px;  /* 调整行高与输入框一致 */
  height: 32px;       /* 固定高度 */
  display: flex;      /* 使用flex布局 */
  align-items: center; /* 垂直居中 */
}

:deep(.el-input__wrapper),
:deep(.el-textarea__wrapper) {
  padding: 4px 11px;
  line-height: 1.5;
  min-height: 32px;  /* 设置最小高度 */
}

:deep(.el-radio__label),
:deep(.el-checkbox__label) {
  padding-left: 8px;
  line-height: 32px;  /* 调整行高 */
  height: 32px;       /* 固定高度 */
}

:deep(.el-radio),
:deep(.el-checkbox) {
  display: flex;      /* 使用flex布局 */
  align-items: center; /* 垂直居中 */
  height: 32px;       /* 固定高度 */
}

:deep(.el-select) {
  width: 100%;
  line-height: 32px;  /* 调整行高 */
}

:deep(.el-select-dropdown__item) {
  padding: 0 12px;
  line-height: 32px;  /* 调整行高 */
  height: 32px;       /* 固定高度 */
}

:deep(.el-date-editor) {
  width: 100% !important;
  line-height: 32px;  /* 调整行高 */
}

.preview-form.is-dialog {
  padding: 16px;
  margin: 0;
  box-shadow: none;
  border: none;
}

.preview-form.is-dialog .form-header {
  margin-bottom: 16px;
}

.preview-form.is-dialog .form-title {
  font-size: 18px;
}

.preview-form.is-dialog .form-actions {
  margin-top: 24px;
}

.upload-demo {
  width: 100%;
}

.upload-demo :deep(.el-upload-list) {
  width: 100%;
}
</style> 