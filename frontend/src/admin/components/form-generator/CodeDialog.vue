<template>
  <el-dialog
    v-model="dialogVisible"
    title="生成代码"
    width="800px"
  >
    <el-tabs v-model="activeTab">
      <el-tab-pane label="Vue Template" name="template">
        <pre><code>{{ generateTemplate() }}</code></pre>
      </el-tab-pane>
      <el-tab-pane label="JavaScript" name="script">
        <pre><code>{{ generateScript() }}</code></pre>
      </el-tab-pane>
    </el-tabs>
    <template #footer>
      <el-button @click="dialogVisible = false">关闭</el-button>
      <el-button type="primary" @click="handleCopy">复制代码</el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, computed } from 'vue'
import { ElMessage } from 'element-plus'

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  formItems: {
    type: Array,
    required: true
  }
})

const emit = defineEmits(['update:visible'])

const dialogVisible = computed({
  get: () => props.visible,
  set: (val) => emit('update:visible', val)
})

const activeTab = ref('template')

const generateTemplate = () => {
  // 先生成表单项
  const formItems = props.formItems.map(item => {
    let rules = item.props.required ? 
      `\n      :rules="{ required: true, message: '请输入${item.props.label}', trigger: 'blur' }"` : ''
    
    let placeholder = item.props.placeholder ? 
      `\n      placeholder="${item.props.placeholder}"` : ''
    
    let content = ''
    if (item.type === 'el-select') {
      content = `
        <el-option
          v-for="option in ${item.key}_options"
          :key="option.value"
          :label="option.label"
          :value="option.value"
        />`
    }

    return `    <el-form-item
      label="${item.props.label}"
      prop="${item.key}"${rules}>
      <${item.type}
        v-model="formData.${item.key}"${placeholder}
        style="width: 100%"
      >${content}
      </${item.type}>
    </el-form-item>`
  }).join('\n\n')

  // 生成完整的模板代码
  return `<template>
  <el-form
    :model="formData"
    label-width="100px"
    @submit.prevent="handleSubmit"
  >
${formItems}

    <el-form-item>
      <el-button type="primary" @click="handleSubmit">提交</el-button>
      <el-button @click="resetForm">重置</el-button>
    </el-form-item>
  </el-form>
</template>`
}

const generateScript = () => {
  // 生成表单数据初始值
  const formDataInit = props.formItems.reduce((acc, curr) => {
    acc[curr.key] = curr.type === 'el-upload' ? [] : ''
    return acc
  }, {})

  // 生成选择框选项数据
  const selectOptions = props.formItems
    .filter(item => item.type === 'el-select')
    .reduce((acc, curr) => {
      acc[`${curr.key}_options`] = curr.props.options
      return acc
    }, {})

  // 生成选项常量定义
  const optionsCode = Object.entries(selectOptions)
    .map(([key, value]) => `const ${key} = ${JSON.stringify(value, null, 2)}`)
    .join('\n\n')

  // 生成完整的脚本代码
  return `<script setup>
import { reactive } from 'vue'
import { ElMessage } from 'element-plus'

const formData = reactive(${JSON.stringify(formDataInit, null, 2)})

${optionsCode}

const handleSubmit = async () => {
  try {
    // 这里添加表单提交逻辑
    console.log('表单数据:', formData)
    ElMessage.success('提交成功')
  } catch (error) {
    console.error('提交失败:', error)
    ElMessage.error('提交失败')
  }
}

const resetForm = () => {
  Object.keys(formData).forEach(key => {
    formData[key] = Array.isArray(formData[key]) ? [] : ''
  })
}
<\/script>

<style scoped>
.el-form {
  max-width: 600px;
  margin: 0 auto;
  padding: 20px;
}
<\/style>`
}

const handleCopy = () => {
  const code = activeTab.value === 'template' ? generateTemplate() : generateScript()
  navigator.clipboard.writeText(code).then(() => {
    ElMessage.success('代码已复制到剪贴板')
  }).catch(() => {
    ElMessage.error('复制失败')
  })
}
</script>

<style scoped>
pre {
  background-color: #f5f5f5;
  padding: 16px;
  border-radius: 4px;
  overflow: auto;
}

code {
  font-family: Consolas, Monaco, 'Andale Mono', monospace;
  font-size: 14px;
}
</style> 