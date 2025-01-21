<template>
  <div class="config-manage">
    <el-card>
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <el-input
              v-model="searchForm.keyword"
              placeholder="搜索配置名称/键"
              style="width: 200px"
              clearable
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
            >
              <el-option
                v-for="item in moduleList"
                :key="item.moduleKey"
                :label="item.moduleName"
                :value="item.moduleKey"
              />
            </el-select>
            <el-button type="primary" @click="handleSearch" style="margin-left: 16px">
              搜索
            </el-button>
            <el-button @click="resetSearch">重置</el-button>
          </div>
          <div class="header-right">
            <el-button type="primary" @click="handleAdd">
              <el-icon><Plus /></el-icon>
              新增配置
            </el-button>
          </div>
        </div>
      </template>

      <el-table
        v-loading="loading"
        :data="configList"
        style="width: 100%"
      >
        <el-table-column prop="configName" label="配置名称" min-width="150" />
        <el-table-column prop="configKey" label="配置键" min-width="150" show-overflow-tooltip />
        <el-table-column prop="configValue" label="配置值" min-width="200" show-overflow-tooltip />
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
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">
              编辑
            </el-button>
            <el-button type="danger" link @click="handleDelete(row)">
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="searchForm.pageNum"
        v-model:page-size="searchForm.pageSize"
        :total="total"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next"
        class="pagination"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </el-card>

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
            placeholder="请输入配置键"
            :disabled="dialogType === 'edit'"
          />
        </el-form-item>
        <el-form-item label="配置类型" prop="configType">
          <el-select v-model="configForm.configType" style="width: 100%">
            <el-option label="字符串" value="string" />
            <el-option label="数字" value="number" />
            <el-option label="布尔值" value="boolean" />
            <el-option label="日期" value="date" />
          </el-select>
        </el-form-item>
        <el-form-item label="配置值" prop="configValue">
          <component :is="getConfigValueComponent" v-model="configForm.configValue" />
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
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Plus } from '@element-plus/icons-vue'
import modules from '@/config/modules'
import { 
  getConfigList as getConfigListApi,
  addConfig,
  updateConfig,
  deleteConfig,
  getConfigInfo,
  checkConfigKey
} from '@/api/system/config'
import { 
  getConfigTypeTag,
  getConfigTypeName,
  parseConfigValue,
  formatConfigValue,
  validateConfigValue
} from '@/utils/config'

// 列表数据
const loading = ref(false)
const configList = ref([])
const total = ref(0)
const moduleList = modules

// 分页
const searchForm = reactive({
  keyword: '',
  moduleKey: '',
  pageNum: 1,
  pageSize: 10
})

// 配置表单
const dialogVisible = ref(false)
const dialogType = ref('add')
const configFormRef = ref()
const configForm = reactive({
  id: undefined,
  moduleKey: '',
  configName: '',
  configKey: '',
  configType: 'string',
  configValue: '',
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
    { required: true, message: '请输入配置键', trigger: 'blur' },
    { min: 2, max: 50, message: '长度在 2 到 50 个字符', trigger: 'blur' },
  ],
  configType: [
    { required: true, message: '请选择配置类型', trigger: 'change' }
  ],
  configValue: [
    { required: true, message: '请输入配置值', trigger: 'blur' },
    { 
      validator: (rule, value, callback) => {
        if (!validateConfigValue(value, configForm.configType)) {
          callback(new Error('配置值格式不正确'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ]
}

// 根据配置类型获取对应的输入组件
const getConfigValueComponent = computed(() => {
  switch (configForm.configType) {
    case 'boolean':
      return 'el-radio-group'
    case 'date':
      return 'el-date-picker'
    case 'number':
      return 'el-input-number'
    default:
      return 'el-input'
  }
})

// 获取配置列表
const getConfigList = async () => {
  loading.value = true
  try {
    const res = await getConfigListApi(searchForm)
    console.log(res)
    configList.value = res.list
    total.value = res.total
  } catch (error) {
    console.error('获取配置列表失败:', error)
    ElMessage.error('获取配置列表失败')
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  searchForm.pageNum = 1
  getConfigList()
}

// 重置搜索
const resetSearch = () => {
  searchForm.keyword = ''
  searchForm.moduleKey = ''
  searchForm.pageNum = 1
  handleSearch()
}

// 新增配置
const handleAdd = () => {
  dialogType.value = 'add'
  Object.assign(configForm, {
    moduleKey: '',
    configName: '',
    configKey: '',
    configType: 'string',
    configValue: '',
    remark: ''
  })
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
            configName: configForm.configName,
            configKey: configForm.configKey,
            configValue: configValue,
            configType: configForm.configType,
            moduleKey: configForm.moduleKey,
            moduleName: moduleName,
            remark: configForm.remark
          }
          const addRes = await addConfig(addData)
        } else {
          const updateData = {
            id: configForm.id,
            configName: configForm.configName,
            configKey: configForm.configKey,
            configValue: configValue,
            configType: configForm.configType,
            moduleKey: configForm.moduleKey,
            moduleName: moduleName,
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

// 分页相关方法
const handleSizeChange = (val) => {
  searchForm.pageSize = val
  searchForm.pageNum = 1
  getConfigList()
}

const handleCurrentChange = (val) => {
  searchForm.pageNum = val
  getConfigList()
}

onMounted(() => {
  getConfigList()
})
</script>

<style scoped>
.config-manage {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-left {
  display: flex;
  align-items: center;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style> 