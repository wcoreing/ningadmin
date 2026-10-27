<template>
  <div class="dept-manage">
    <el-card>
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <el-input
              v-model="searchForm.keyword"
              placeholder="搜索部门名称"
              style="width: 200px"
              clearable
              @keyup.enter="handleSearch"
            >
              <template #prefix>
                <el-icon><Search /></el-icon>
              </template>
            </el-input>
            <el-select
              v-model="searchForm.status"
              placeholder="部门状态"
              clearable
              style="width: 120px; margin-left: 16px"
            >
              <el-option label="正常" value="0" />
              <el-option label="停用" value="1" />
            </el-select>
            <el-button type="primary" @click="handleSearch" style="margin-left: 16px">
              搜索
            </el-button>
            <el-button @click="resetSearch">重置</el-button>
          </div>
          <div class="header-right">
            <el-button @click="toggleExpand" style="margin-right: 16px">
              <el-icon>
                <component :is="isExpanded ? 'Fold' : 'Expand'" />
              </el-icon>
              {{ isExpanded ? '折叠' : '展开' }}
            </el-button>
            <el-button type="primary" @click="handleAdd">
              <el-icon><Plus /></el-icon>
              新增部门
            </el-button>
          </div>
        </div>
      </template>

      <el-table
        ref="deptTableRef"
        v-loading="loading"
        :data="deptList"
        row-key="deptId"
        :tree-props="{ children: 'children', hasChildren: 'hasChildren' }"
        style="width: 100%"
      >
        <el-table-column prop="deptName" label="部门名称" min-width="200">
          <template #default="{ row }">
            <el-icon v-if="row.icon">
              <component :is="iconMap[row.icon]" v-if="iconMap[row.icon]" />
            </el-icon>
            <span style="margin-left: 8px">{{ row.deptName }}</span>
          </template>
        </el-table-column>
        <el-table-column label="排序" width="120" align="center">
          <template #default="{ row }">
            <el-input-number 
              v-model="row.orderNum" 
              :min="0" 
              :max="999"
              size="small"
              controls-position="right"
              @change="handleSortChange(row)"
              style="width: 90px"
            />
          </template>
        </el-table-column>
        <el-table-column prop="leader" label="负责人" width="120" />
        <el-table-column prop="phone" label="联系电话" width="130" />
        <el-table-column prop="status" label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === '0' ? 'success' : 'danger'">
              {{ row.status === '0' ? '正常' : '停用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="创建时间" align="center" prop="createTime" width="200">
          <template #default="{ row }">
            <span>{{ formatDateTime(row.createTime) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleAdd(row)">
              新增
            </el-button>
            <el-button type="primary" link @click="handleEdit(row)">
              编辑
            </el-button>
            <el-button 
              v-if="row.parentId !== 0"
              type="danger" 
              link 
              @click="handleDelete(row)"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 部门表单弹窗 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogType === 'add' ? '新增部门' : '编辑部门'"
      width="600px"
    >
      <el-form
        ref="deptFormRef"
        :model="deptForm"
        :rules="deptRules"
        label-width="100px"
      >
        <el-form-item label="上级部门" v-if="deptForm.parentId !== 0">
          <el-tree-select
            v-model="deptForm.parentId"
            :data="deptOptions"
            :props="{ value: 'deptId', label: 'deptName', children: 'children' }"
            value-key="deptId"
            placeholder="选择上级部门"
            check-strictly
            :render-after-expand="false"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="部门名称" prop="deptName">
          <el-input v-model="deptForm.deptName" placeholder="请输入部门名称" />
        </el-form-item>
        <el-form-item label="显示排序" prop="orderNum">
          <el-input-number v-model="deptForm.orderNum" :min="0" :max="999" />
        </el-form-item>
        <el-form-item label="负责人" prop="leader">
          <el-input v-model="deptForm.leader" placeholder="请输入负责人" maxlength="20" />
        </el-form-item>
        <el-form-item label="联系电话" prop="phone">
          <el-input v-model="deptForm.phone" placeholder="请输入联系电话" maxlength="11" />
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="deptForm.email" placeholder="请输入邮箱" maxlength="50" />
        </el-form-item>
        <el-form-item label="部门状态">
          <el-radio-group v-model="deptForm.status">
            <el-radio value="0">正常</el-radio>
            <el-radio value="1">停用</el-radio>
          </el-radio-group>
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
import { ref, reactive, onMounted, nextTick } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import { getDeptList, getDeptTree, addDept, updateDept, deleteDept, updateDeptStatus } from '@/admin/api/system/dept'
import { Search, Refresh, Plus, Sort, Expand, Fold } from '@element-plus/icons-vue'

// 图标映射
const iconMap = ElementPlusIconsVue

// 搜索表单
const searchForm = reactive({
  keyword: '',
  status: ''
})

// 列表数据
const loading = ref(false)
const deptList = ref([])

// 获取部门列表
const getList = async () => {
  loading.value = true
  try {
    // 如果有搜索条件，使用列表接口
    let res
    if (searchForm.keyword || searchForm.status) {
      res = await getDeptTree(searchForm)
      // 搜索时直接使用返回的树形结构
      deptList.value = Array.isArray(res) ? res : []
    } else {
      // 没有搜索条件，直接使用树形接口
      res = await getDeptTree()
      deptList.value = Array.isArray(res) ? res : []
    }
    // 设置默认展开
    nextTick(() => {
      isExpanded.value = true
      const allNodes = getAllNodes(deptList.value)
      allNodes.forEach(node => {
        deptTableRef.value?.toggleRowExpansion(node, true)
      })
    })
  } catch (error) {
    console.error('获取部门列表失败:', error)
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  getList()
}

// 重置搜索
const resetSearch = () => {
  searchForm.keyword = ''
  searchForm.status = ''
  getList() // 直接调用getList重新获取数据
}

// 部门表单
const dialogVisible = ref(false)
const dialogType = ref('add')
const deptFormRef = ref()
const deptForm = reactive({
  deptId: undefined,
  parentId: 0,
  deptName: '',
  orderNum: 0,
  leader: '',
  phone: '',
  email: '',
  status: '0'
})

// 表单校验规则
const deptRules = {
  deptName: [
    { required: true, message: '部门名称不能为空', trigger: 'blur' },
    { min: 2, max: 50, message: '长度在 2 到 50 个字符', trigger: 'blur' }
  ],
  orderNum: [
    { required: true, message: '显示排序不能为空', trigger: 'blur' }
  ],
  email: [
    {
      type: 'email',
      message: '请输入正确的邮箱地址',
      trigger: ['blur', 'change']
    }
  ],
  phone: [
    {
      pattern: /^1[3-9]\d{9}$/,
      message: '请输入正确的手机号码',
      trigger: 'blur'
    }
  ]
}

// 部门树选项
const deptOptions = ref([])

// 获取部门树数据
const getDeptTreeData = async () => {
  try {
    const res = await getDeptTree()
    deptOptions.value = [
      { deptId: 0, deptName: '顶级部门' },
      ...(Array.isArray(res) ? res : [])
    ]
  } catch (error) {
    console.error('获取部门树失败:', error)
  }
}

// 新增部门
const handleAdd = (row) => {
  dialogType.value = 'add'
  Object.assign(deptForm, {
    deptId: undefined,
    parentId: row?.deptId || 0,
    deptName: '',
    orderNum: 0,
    leader: '',
    phone: '',
    email: '',
    status: '0'
  })
  getDeptTreeData()
  dialogVisible.value = true
}

// 编辑部门
const handleEdit = async (row) => {
  dialogType.value = 'edit'
  try {
    Object.assign(deptForm, { ...row })
    getDeptTreeData()
    dialogVisible.value = true
  } catch (error) {
    console.error('获取部门详情失败:', error)
    ElMessage.error('获取部门详情失败')
  }
}

// 删除部门
const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确认要删除部门"${row.deptName}"吗？`,
      '警告',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await deleteDept(row.deptId)
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
const handleSubmit = async () => {
  if (!deptFormRef.value) return
  
  await deptFormRef.value.validate(async (valid) => {
    if (valid) {
      try {
        if (dialogType.value === 'add') {
          await addDept(deptForm)
        } else {
          await updateDept(deptForm)
        }
        
        ElMessage.success(dialogType.value === 'add' ? '新增成功' : '更新成功')
        dialogVisible.value = false
        getList()
      } catch (error) {
        console.error('提交失败:', error)
        ElMessage.error('提交失败')
      }
    }
  })
}

const deptTableRef = ref(null)
const isExpanded = ref(false)

// 递归获取所有节点
const getAllNodes = (data) => {
  const nodes = []
  const traverse = (items) => {
    items.forEach(item => {
      nodes.push(item)
      if (item.children && item.children.length) {
        traverse(item.children)
      }
    })
  }
  traverse(data)
  return nodes
}

// 切换展开/折叠
const toggleExpand = () => {
  isExpanded.value = !isExpanded.value
  const allNodes = getAllNodes(deptList.value)
  allNodes.forEach(node => {
    deptTableRef.value?.toggleRowExpansion(node, isExpanded.value)
  })
}

// 处理排序变更
const handleSortChange = async (row) => {
  try {
    await updateDept({
      deptId: row.deptId,
      parentId: row.parentId,
      deptName: row.deptName,
      orderNum: row.orderNum,
      leader: row.leader || '',
      phone: row.phone || '',
      email: row.email || '',
      status: row.status
    })
    ElMessage.success('排序更新成功')
  } catch (error) {
    console.error('更新排序失败:', error)
    ElMessage.error('更新排序失败')
    // 恢复原值
    getList()
  }
}

// 格式化日期时间
const formatDateTime = (dateStr) => {
  if (!dateStr) return '-'
  try {
    const date = new Date(dateStr)
    if (isNaN(date.getTime())) return '-'
    return date.toLocaleString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit'
    })
  } catch (error) {
    return '-'
  }
}

onMounted(() => {
  getList()
})
</script>

<style scoped>
.el-table :deep(.el-table__cell) {
  padding: 6px 0;
}

.el-table :deep(.el-input-number.el-input-number--small) {
  margin: 0 auto;
}

.dept-manage {
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
</style> 