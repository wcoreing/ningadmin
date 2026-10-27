<template>
  <div class="app-container">
    <!-- 页面标题 -->
    <div class="page-header">
      <h2>配置管理</h2>
    </div>
    
    <!-- 搜索工具栏 -->
    <div class="filter-container">
      <el-input
        v-model="searchForm.keyword"
        placeholder="搜索配置名称/键"
        style="width: 200px"
        class="filter-item"
        @keyup.enter="handleSearch"
      >
        <template #prefix>
          <el-icon><Search /></el-icon>
        </template>
      </el-input>
      <el-select
        v-model="searchForm.moduleKey"
        placeholder="选择模块"
        clearable
        style="width: 120px; margin-left: 16px"
        class="filter-item"
      >
        <el-option
          v-for="item in moduleList"
          :key="item.moduleKey"
          :label="item.moduleName"
          :value="item.moduleKey"
        />
      </el-select>
      <el-select
        v-model="searchForm.status"
        placeholder="状态"
        clearable
        class="filter-item"
        style="width: 130px; margin-left: 16px"
      >
        <el-option label="启用" :value="1" />
        <el-option label="禁用" :value="0" />
      </el-select>
      <el-button
        class="filter-item"
        type="primary"
        :icon="Search"
        @click="handleSearch"
        style="margin-left: 16px"
      >
        搜索
      </el-button>
      <el-button @click="resetSearch">重置</el-button>
      <el-button
        class="filter-item"
        type="primary"
        :icon="Plus"
        @click="handleAdd"
        style="margin-left: 16px"
      >
        新增配置
      </el-button>
    </div>

    <!-- 配置表格 -->
    <el-table
      v-loading="loading"
      :data="configList"
      row-key="id"
      border
      default-expand-all
      :tree-props="{ children: 'children', hasChildren: 'hasChildren' }"
    >
      <template #empty>
        <span>暂无数据</span>
      </template>
      <el-table-column prop="configName" label="配置名称" min-width="200" />
      <el-table-column prop="configKey" label="配置键" min-width="250" show-overflow-tooltip />
      <el-table-column prop="configValue" label="配置值" min-width="150" show-overflow-tooltip>
        <template #default="{ row }">
          <div v-if="row.configType === 'image' && row.configValue">
            <el-image 
              :src="getImageUrl(row.configValue)" 
              fit="cover" 
              style="width: 40px; height: 40px; border-radius: 4px; cursor: pointer;"
              :preview-src-list="[getImageUrl(row.configValue)]"
              preview-teleported
            />
          </div>
          <span v-else>{{ row.configValue || '-' }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="moduleName" label="所属模块" width="120">
        <template #default="{ row }">
          <el-tag>{{ row.moduleName }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="configType" label="配置类型" width="100">
        <template #default="{ row }">
          <el-tag :type="getConfigTypeTag(row.configType)">
            {{ getConfigTypeName(row.configType) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="sort" label="排序" width="80" align="center" />
      <el-table-column prop="status" label="状态" width="100" align="center">
        <template #default="{ row }">
          <el-tag :type="row.status === 1 ? 'success' : 'info'">
            {{ row.status === 1 ? '启用' : '禁用' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column
        label="操作"
        align="center"
        width="280"
        class-name="small-padding fixed-width"
      >
        <template #default="{ row }">
          <el-button
            link
            type="primary"
            :icon="Plus"
            @click="handleAdd(row)"
          >
            添加子配置
          </el-button>
          <el-button
            link
            type="primary"
            :icon="Edit"
            @click="handleEdit(row)"
          >
            编辑
          </el-button>
          <el-button
            link
            type="danger"
            :icon="Delete"
            @click="handleDelete(row)"
          >
            删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 配置表单弹窗 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogType === 'add' ? '新增配置' : '编辑配置'"
      width="500px"
    >
      <el-form
        ref="configFormRef"
        :model="configForm"
        :rules="configRules"
        label-width="100px"
      >
        <el-form-item label="上级配置">
          <el-tree-select
            v-model="configForm.parentId"
            :data="configOptions"
            :props="{ label: 'configName', value: 'id' }"
            value-key="id"
            placeholder="选择上级配置"
            check-strictly
            clearable
          />
        </el-form-item>
        <el-form-item label="所属模块" prop="moduleKey">
          <el-select v-model="configForm.moduleKey" style="width: 100%">
            <el-option
              v-for="item in moduleList"
              :key="item.moduleKey"
              :label="item.moduleName"
              :value="item.moduleKey"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="配置名称" prop="configName">
          <el-input v-model="configForm.configName" placeholder="请输入配置名称" />
        </el-form-item>
        <el-form-item label="配置键" prop="configKey">
          <el-input 
            v-model="configForm.configKey" 
            placeholder="可选，用于程序中获取配置值"
            :disabled="dialogType === 'edit'"
          />
          <div class="form-tip">
            <el-text size="small" type="info">
              配置键为可选项，用于程序中获取配置值，如不填写则仅作为分组使用
            </el-text>
          </div>
        </el-form-item>
        <el-form-item label="配置类型" prop="configType">
          <el-select v-model="configForm.configType" style="width: 100%" placeholder="可选，默认为字符串">
            <el-option label="字符串" value="string" />
            <el-option label="数字" value="number" />
            <el-option label="布尔值" value="boolean" />
            <el-option label="日期" value="date" />
            <el-option label="图片" value="image" />
          </el-select>
        </el-form-item>
        <el-form-item label="配置值" prop="configValue">
          <!-- 图片上传组件 -->
          <div v-if="configForm.configType === 'image'" class="image-upload-container">
            <div v-if="configForm.configValue" class="image-preview">
              <el-image 
                :src="getImageUrl(configForm.configValue)" 
                fit="cover" 
                style="width: 100px; height: 100px; border-radius: 6px;"
              />
              <div class="image-actions">
                <el-button 
                  type="primary" 
                  size="small" 
                  @click="openFileManager"
                  style="margin-right: 8px;"
                >
                  更换图片
                </el-button>
                <el-button 
                  type="danger" 
                  size="small" 
                  @click="removeImage"
                >
                  删除图片
                </el-button>
              </div>
            </div>
            <div v-else class="upload-placeholder" @click="openFileManager">
              <el-icon size="40" color="#c0c4cc"><Plus /></el-icon>
              <div style="margin-top: 8px; color: #c0c4cc;">点击上传图片</div>
            </div>
          </div>
          <!-- 其他类型的输入组件 -->
          <component 
            v-else
            :is="getConfigValueComponent" 
            v-model="configForm.configValue" 
            placeholder="可选，具体的配置内容" 
          />
          <div class="form-tip">
            <el-text size="small" type="info">
              配置值为可选项，如不填写则仅作为分组使用
            </el-text>
          </div>
        </el-form-item>
        <el-form-item label="排序" prop="sort">
          <el-input-number v-model="configForm.sort" :min="0" :max="999" />
        </el-form-item>
        <el-form-item label="状态">
          <el-radio-group v-model="configForm.status">
            <el-radio :value="1">启用</el-radio>
            <el-radio :value="0">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="备注" prop="remark">
          <el-input
            v-model="configForm.remark"
            type="textarea"
            placeholder="请输入备注信息"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>

    <!-- 文件管理器 -->
    <file-manager
      v-model="fileManagerVisible"
      :multiple="false"
      @select="handleImageSelect"
    />
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Plus, Edit, Delete } from '@element-plus/icons-vue'
import modules from '@/admin/config/modules'
import FileManager from '@/admin/components/FileManager/index.vue'
import { 
  getConfigList as getConfigListApi,
  addConfig,
  updateConfig,
  deleteConfig,
  getConfigInfo
} from '@/admin/api/system/config'
import { 
  getConfigTypeTag,
  getConfigTypeName,
  parseConfigValue,
  formatConfigValue,
  validateConfigValue
} from '@/admin/utils/config'

// 列表数据
const loading = ref(false)
const configList = ref([])
const configOptions = ref([])
const moduleList = modules

// 搜索表单
const searchForm = reactive({
  keyword: '',
  moduleKey: '',
  status: undefined
})

// 配置表单
const dialogVisible = ref(false)
const dialogType = ref('add')
const configFormRef = ref()
const configForm = reactive({
  id: undefined,
  parentId: 0,
  moduleKey: '',
  configName: '',
  configKey: '',
  configType: 'string',
  configValue: '',
  sort: 0,
  status: 1,
  remark: ''
})

// 表单校验规则
const configRules = {
  moduleKey: [
    { required: true, message: '请选择所属模块', trigger: 'change' }
  ],
  configName: [
    { required: true, message: '请输入配置名称', trigger: 'blur' },
    { min: 2, max: 50, message: '长度在 2 到 50 个字符', trigger: 'blur' }
  ],
  configKey: [
    { 
      min: 2, 
      max: 50, 
      message: '长度在 2 到 50 个字符', 
      trigger: 'blur',
      validator: (rule, value, callback) => {
        if (value && value.length > 0 && value.length < 2) {
          callback(new Error('配置键长度在 2 到 50 个字符'))
        } else {
          callback()
        }
      }
    }
  ],
  configValue: [
    { 
      validator: (rule, value, callback) => {
        if (value && configForm.configType && !validateConfigValue(value, configForm.configType)) {
          callback(new Error('配置值格式不正确'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ]
}

// 文件管理器
const fileManagerVisible = ref(false)

// 根据配置类型获取对应的输入组件
const getConfigValueComponent = computed(() => {
  switch (configForm.configType) {
    case 'boolean':
      return 'el-radio-group'
    case 'date':
      return 'el-date-picker'
    case 'number':
      return 'el-input-number'
    case 'image':
      return null // 图片类型使用自定义组件
    default:
      return 'el-input'
  }
})

// 获取图片URL
const getImageUrl = (url) => {
  if (!url) return ''
  if (url.startsWith('http')) return url
  // 如果是相对路径，拼接基础URL
  const baseUrl = import.meta.env.VITE_IMG_URL || import.meta.env.VITE_API_BASE_URL || window.location.origin
  return baseUrl + url
}

// 打开文件管理器
const openFileManager = () => {
  fileManagerVisible.value = true
}

// 处理图片选择
const handleImageSelect = (file) => {
  if (file && file.url) {
    configForm.configValue = file.url
    fileManagerVisible.value = false
  }
}

// 删除图片
const removeImage = () => {
  configForm.configValue = ''
}

// 获取配置列表
const getConfigList = async () => {
  loading.value = true
  try {
    // 获取树形结构数据
    const res = await getConfigListApi({
      ...searchForm,
      pageNum: 0,
      pageSize: 0
    })
    console.log('配置列表响应:', res)
    
    // 确保数据是数组
    if (Array.isArray(res.list)) {
      configList.value = res.list
      configOptions.value = [{ 
        id: 0, 
        configName: '顶级配置', 
        children: res.list
      }]
    } else {
      configList.value = []
      configOptions.value = [{ 
        id: 0, 
        configName: '顶级配置', 
        children: []
      }]
    }
  } catch (error) {
    console.error('获取配置列表失败:', error)
    configList.value = []
    configOptions.value = [{ 
      id: 0, 
      configName: '顶级配置', 
      children: []
    }]
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  getConfigList()
}

// 重置搜索
const resetSearch = () => {
  searchForm.keyword = ''
  searchForm.moduleKey = ''
  searchForm.status = undefined
  handleSearch()
}

// 重置表单
const resetForm = () => {
  configForm.id = undefined
  configForm.parentId = 0
  configForm.moduleKey = ''
  configForm.configName = ''
  configForm.configKey = ''
  configForm.configType = 'string'
  configForm.configValue = ''
  configForm.sort = 0
  configForm.status = 1
  configForm.remark = ''
  configFormRef.value?.resetFields()
}

// 新增配置
const handleAdd = (row) => {
  dialogType.value = 'add'
  resetForm()
  if (row) {
    configForm.parentId = row.id
  }
  dialogVisible.value = true
}

// 编辑配置
const handleEdit = async (row) => {
  dialogType.value = 'edit'
  try {
    const res = await getConfigInfo(row.id)
    Object.assign(configForm, res)
    dialogVisible.value = true
  } catch (error) {
    console.error('获取配置详情失败:', error)
    ElMessage.error('获取配置详情失败')
  }
}

// 删除配置
const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确认要删除配置"${row.configName}"吗？`,
      '警告',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await deleteConfig(row.id)
    ElMessage.success('删除成功')
    getConfigList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除失败:', error)
      ElMessage.error('删除失败')
    }
  }
}

// 提交表单
const handleSubmit = async () => {
  if (!configFormRef.value) return
  
  await configFormRef.value.validate(async (valid) => {
    if (valid) {
      try {
        // 获取当前选中模块的名称
        const currentModule = moduleList.find(m => m.moduleKey === configForm.moduleKey)
        const moduleName = currentModule ? currentModule.moduleName : ''

        // 将配置值转换为字符串
        const configValue = String(configForm.configValue)

        if (dialogType.value === 'add') {
          const addData = {
            parentId: configForm.parentId,
            configName: configForm.configName,
            configKey: configForm.configKey,
            configValue: configValue,
            configType: configForm.configType,
            moduleKey: configForm.moduleKey,
            moduleName: moduleName,
            sort: configForm.sort,
            status: configForm.status,
            remark: configForm.remark
          }
          await addConfig(addData)
        } else {
          const updateData = {
            id: configForm.id,
            parentId: configForm.parentId,
            configName: configForm.configName,
            configKey: configForm.configKey,
            configValue: configValue,
            configType: configForm.configType,
            moduleKey: configForm.moduleKey,
            moduleName: moduleName,
            sort: configForm.sort,
            status: configForm.status,
            remark: configForm.remark
          }
          await updateConfig(updateData)
        }
        
        ElMessage.success(dialogType.value === 'add' ? '新增成功' : '更新成功')
        dialogVisible.value = false
        getConfigList()
      } catch (error) {
        console.error('提交失败:', error)
        ElMessage.error('提交失败')
      }
    }
  })
}



onMounted(() => {
  getConfigList()
})
</script>

<style scoped>
.page-header {
  margin-bottom: 20px;
}

.page-header h2 {
  margin: 0;
  color: #303133;
  font-size: 20px;
  font-weight: 500;
}

.filter-container {
  padding-bottom: 10px;
}

.filter-item {
  margin-right: 10px;
}

.form-tip {
  margin-top: 4px;
}

.image-upload-container {
  .image-preview {
    display: flex;
    align-items: center;
    gap: 16px;
  }
  
  .image-actions {
    display: flex;
    flex-direction: column;
    gap: 8px;
  }
  
  .upload-placeholder {
    width: 100px;
    height: 100px;
    border: 2px dashed #d9d9d9;
    border-radius: 6px;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    transition: border-color 0.3s;
    
    &:hover {
      border-color: #409eff;
    }
  }
}
</style> 