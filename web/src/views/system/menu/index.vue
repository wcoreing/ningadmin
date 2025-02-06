<template>
  <div class="menu-manage">
    <el-card>
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <el-input
              v-model="searchForm.keyword"
              placeholder="搜索菜单名称"
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
              placeholder="菜单状态"
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
              新增菜单
            </el-button>
          </div>
        </div>
      </template>

      <el-table
        ref="menuTableRef"
        v-loading="loading"
        :data="menuList"
        row-key="id"
        :tree-props="{ children: 'children' }"
        style="width: 100%"
      >
        <el-table-column prop="name" label="菜单名称" min-width="200">
          <template #default="{ row }">
            <el-icon v-if="row.icon">
              <component :is="iconMap[row.icon]" v-if="iconMap[row.icon]" />
            </el-icon>
            <span style="margin-left: 8px">{{ row.name }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="type" label="类型" width="100">
          <template #default="{ row }">
            <el-tag :type="row.type === 'menu' ? '' : 'info'">
              {{ row.type === 'menu' ? '菜单' : '按钮' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="routePath" label="路由路径" min-width="150" show-overflow-tooltip />
        <el-table-column prop="permission" label="权限标识" min-width="150" show-overflow-tooltip />
        <el-table-column label="排序" width="120" align="center">
          <template #default="{ row }">
            <el-input-number 
              v-model="row.sort" 
              :min="0" 
              :max="999"
              size="small"
              controls-position="right"
              @change="handleSortChange(row)"
              style="width: 90px"
            />
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === '0' ? 'success' : 'danger'">
              {{ row.status === '0' ? '正常' : '停用' }}
            </el-tag>
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
            <el-button type="danger" link @click="handleDelete(row)">
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 菜单表单弹窗 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogType === 'add' ? '新增菜单' : '编辑菜单'"
      width="600px"
    >
      <el-form
        ref="menuFormRef"
        :model="menuForm"
        :rules="menuRules"
        label-width="100px"
      >
        <el-form-item label="上级菜单">
          <el-tree-select
            v-model="menuForm.parentId"
            :data="menuOptions"
            :props="{ label: 'name', value: 'id' }"
            placeholder="选择上级菜单"
            style="width: 100%"
            clearable
          />
        </el-form-item>
        <el-form-item label="菜单类型" prop="type">
          <el-radio-group v-model="menuForm.type">
            <el-radio label="menu">菜单</el-radio>
            <el-radio label="button">按钮</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="菜单名称" prop="name">
          <el-input v-model="menuForm.name" placeholder="请输入菜单名称" />
        </el-form-item>
        <el-form-item label="菜单图标" v-if="menuForm.type === 'menu'">
          <el-input v-model="menuForm.icon" placeholder="请选择图标">
            <template #append>
              <el-button @click="showIconSelect = true">
                选择图标
              </el-button>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item label="路由路径" prop="routePath" v-if="menuForm.type === 'menu'">
          <el-input v-model="menuForm.routePath" placeholder="请输入路由路径" />
        </el-form-item>
        <el-form-item label="组件路径" v-if="menuForm.type === 'menu'">
          <el-input v-model="menuForm.componentPath" placeholder="请输入组件路径" />
        </el-form-item>
        <el-form-item label="权限标识" prop="permission">
          <el-input v-model="menuForm.permission" placeholder="请输入权限标识" />
        </el-form-item>
        <el-form-item label="显示排序" prop="sort">
          <el-input-number v-model="menuForm.sort" :min="0" :max="999" />
        </el-form-item>
        <el-form-item label="菜单状态">
          <el-radio-group v-model="menuForm.status">
            <el-radio label="0">正常</el-radio>
            <el-radio label="1">停用</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>

    <!-- 图标选择器 -->
    <icon-select 
      v-model="showIconSelect" 
      :selected-icon="menuForm.icon"
      @select="handleIconSelect"
    />
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, nextTick } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import IconSelect from '@/components/IconSelect/index.vue'
import {
  getMenuList,
  addMenu,
  updateMenu,
  updateMenuSort,
  deleteMenu,
  getMenuInfo,
  getMenuTree,
  checkPermission
} from '@/api/system/menu'
import { Expand, Fold } from '@element-plus/icons-vue'

// 图标映射
const iconMap = ElementPlusIconsVue

// 搜索表单
const searchForm = reactive({
  keyword: '',
  status: ''
})

// 列表数据
const loading = ref(false)
const menuList = ref([])

// 获取菜单列表
const getList = async () => {
  loading.value = true
  try {
    // 如果有搜索条件，使用列表接口
    let res
    if (searchForm.keyword || searchForm.status) {
      res = await getMenuTree(searchForm)
      // 搜索时直接使用返回的树形结构
      menuList.value = Array.isArray(res) ? res : []
    } else {
      // 没有搜索条件，直接使用树形接口
      res = await getMenuTree()
      menuList.value = Array.isArray(res) ? res : []
    }
    // 设置默认展开
    nextTick(() => {
      isExpanded.value = true
      const allNodes = getAllNodes(menuList.value)
      allNodes.forEach(node => {
        menuTableRef.value?.toggleRowExpansion(node, true)
      })
    })
  } catch (error) {
    console.error('获取菜单列表失败:', error)
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

// 菜单表单
const dialogVisible = ref(false)
const dialogType = ref('add')
const menuFormRef = ref()
const menuForm = reactive({
  id: undefined,
  parentId: 0,
  name: '',
  icon: '',
  type: 'menu',
  routePath: '',
  componentPath: '',
  permission: '',
  sort: 0,
  status: '0'
})

// 表单校验规则
const menuRules = {
  name: [
    { required: true, message: '请输入菜单名称', trigger: 'blur' },
    { min: 2, max: 50, message: '长度在 2 到 50 个字符', trigger: 'blur' }
  ],
  permission: [
    { required: true, message: '请输入权限标识', trigger: 'blur' },
    { pattern: /^[a-z][a-z0-9:]*$/, message: '以小写字母开头，只能包含小写字母、数字和冒号', trigger: 'blur' }
  ]
}

// 菜单树选项
const menuOptions = ref([])

// 获取菜单树数据
const getMenuTreeData = async () => {
  try {
    const res = await getMenuTree()
    menuOptions.value = [
      { id: 0, name: '顶级菜单' },
      ...(Array.isArray(res) ? res : [])
    ]
  } catch (error) {
    console.error('获取菜单树失败:', error)
  }
}

// 新增菜单
const handleAdd = (row) => {
  dialogType.value = 'add'
  Object.assign(menuForm, {
    id: undefined,
    parentId: row?.id || 0,
    name: '',
    icon: '',
    type: 'menu',
    routePath: '',
    componentPath: '',
    permission: '',
    sort: 0,
    status: '0'
  })
  dialogVisible.value = true
}

// 编辑菜单
const handleEdit = async (row) => {
  dialogType.value = 'edit'
  try {
    const res = await getMenuInfo(row.id)
    Object.assign(menuForm, res)
    dialogVisible.value = true
  } catch (error) {
    console.error('获取菜单详情失败:', error)
    ElMessage.error('获取菜单详情失败')
  }
}

// 删除菜单
const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确认要删除菜单"${row.name}"吗？`,
      '警告',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await deleteMenu(row.id)
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
  if (!menuFormRef.value) return
  
  await menuFormRef.value.validate(async (valid) => {
    if (valid) {
      try {
        if (dialogType.value === 'add') {
          // 检查权限标识是否存在
          const checkRes = await checkPermission({
            permission: menuForm.permission
          })
          if (checkRes === true) {
            ElMessage.error('权限标识已存在')
            return
          }
          await addMenu(menuForm)
        } else {
          await updateMenu(menuForm)
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

// 图标选择器状态
const showIconSelect = ref(false)

// 选择图标
const handleIconSelect = (icon) => {
  menuForm.icon = icon
}

const menuTableRef = ref(null)
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
  const allNodes = getAllNodes(menuList.value)
  allNodes.forEach(node => {
    menuTableRef.value?.toggleRowExpansion(node, isExpanded.value)
  })
}

// 处理排序变更
const handleSortChange = async (row) => {
  try {
    await updateMenuSort({
      id: row.id,
      sort: row.sort
    })
    ElMessage.success('排序更新成功')
  } catch (error) {
    console.error('更新排序失败:', error)
    ElMessage.error('更新排序失败')
    // 恢复原值
    getList()
  }
}

onMounted(() => {
  getList()
  getMenuTreeData()
})
</script>

<style scoped>
.el-table :deep(.el-table__cell) {
  padding: 6px 0;
}

.el-table :deep(.el-input-number.el-input-number--small) {
  margin: 0 auto;
}

.menu-manage {
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