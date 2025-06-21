<template>
  <div class="role-manage">
    <el-card>
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <el-input
              v-model="searchForm.roleName"
              placeholder="搜索角色名称"
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
              placeholder="角色状态"
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
            <el-button type="primary" @click="handleAdd">
              <el-icon><Plus /></el-icon>
              新增角色
            </el-button>
          </div>
        </div>
      </template>

      <el-table
        v-loading="loading"
        :data="roleList"
        style="width: 100%"
      >
        <el-table-column prop="roleName" label="角色名称" min-width="150" show-overflow-tooltip />
        <el-table-column prop="roleKey" label="角色标识" min-width="150" show-overflow-tooltip />
        <el-table-column prop="status" label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-switch
              v-model="row.status"
              :active-value="'0'"
              :inactive-value="'1'"
              @change="handleStatusChange(row)"
            />
          </template>
        </el-table-column>
        <el-table-column prop="createTime" label="创建时间" width="180" align="center" />
        <el-table-column label="操作" width="250" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">
              编辑
            </el-button>
            <el-button type="primary" link @click="handlePermission(row)">
              权限
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
        layout="total, sizes, prev, pager, next, jumper"
        class="pagination"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </el-card>

    <!-- 角色表单弹窗 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogType === 'add' ? '新增角色' : '编辑角色'"
      width="500px"
    >
      <el-form
        ref="roleFormRef"
        :model="roleForm"
        :rules="roleRules"
        label-width="100px"
      >
        <el-form-item label="角色名称" prop="roleName">
          <el-input v-model="roleForm.roleName" placeholder="请输入角色名称" />
        </el-form-item>
        <el-form-item label="角色标识" prop="roleKey">
          <el-input v-model="roleForm.roleKey" placeholder="请输入角色标识" />
        </el-form-item>
        <el-form-item label="角色状态">
          <el-radio-group v-model="roleForm.status">
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

    <!-- 分配权限弹窗 -->
    <el-dialog
      v-model="permissionDialogVisible"
      title="分配权限"
      width="600px"
      destroy-on-close
      class="permission-dialog"
    >
      <div class="permission-header">
        <el-alert
          title="温馨提示：勾选模块即可分配该角色的权限，取消勾选则收回权限"
          type="info"
          :closable="false"
          show-icon
        />
        <div class="role-info">
          <span class="label">角色名称：</span>
          <span class="value">{{ currentRole?.roleName }}</span>
          <el-tag 
            :type="currentRole?.status === '0' ? 'success' : 'danger'"
            class="status-tag"
          >
            {{ currentRole?.status === '0' ? '正常' : '停用' }}
          </el-tag>
        </div>
      </div>
      
      <div class="permission-content">
        <div class="tree-toolbar">
          <div class="toolbar-left">
            <el-icon><Operation /></el-icon>
            <span class="title">功能权限</span>
          </div>
          <div class="toolbar-right">
            <el-button type="primary" link @click="handleCheckAll(true)">全选</el-button>
            <el-divider direction="vertical" />
            <el-button type="primary" link @click="handleCheckAll(false)">取消全选</el-button>
            <el-divider direction="vertical" />
            <el-button type="primary" link @click="handleExpandAll(true)">展开</el-button>
            <el-divider direction="vertical" />
            <el-button type="primary" link @click="handleExpandAll(false)">折叠</el-button>
          </div>
        </div>
        <el-scrollbar height="400px">
          <el-tree
            ref="permissionTreeRef"
            :data="permissionList"
            :props="{ 
              label: 'name',
              children: 'children',
              id: 'id'
            }"
            show-checkbox
            node-key="id"
            :default-expand-all="true"
            class="permission-tree"
          >
            <template #default="{ node, data }">
              <div class="custom-tree-node">
                <span>{{ data.name }}</span>
                <el-tag 
                  size="small" 
                  :type="data.type === 'menu' ? 'success' : 'info'"
                  style="margin-left: 8px"
                >
                  {{ data.type === 'menu' ? '菜单' : '按钮' }}
                </el-tag>
                <span 
                  v-if="data.permission" 
                  class="permission-code"
                  :title="data.permission"
                >
                  {{ data.permission }}
                </span>
              </div>
            </template>
          </el-tree>
        </el-scrollbar>
      </div>

      <template #footer>
        <div class="dialog-footer">
          <el-button @click="permissionDialogVisible = false">取 消</el-button>
          <el-button type="primary" @click="handlePermissionSubmit">确 定</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, nextTick } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Plus, Operation } from '@element-plus/icons-vue'
import {
  getRoleList,
  createRole,
  updateRole,
  deleteRole,
  updateRoleStatus,
  getRolePermissions,
  updateRolePermissions,
  getPermissionTree
} from '@/api/system/role'

// 搜索表单
const searchForm = reactive({
  pageNum: 1,
  pageSize: 10,
  roleName: '',
  status: ''
})

// 列表数据
const loading = ref(false)
const roleList = ref([])
const total = ref(0)

// 获取角色列表
const getList = async () => {
  loading.value = true
  try {
    const res = await getRoleList(searchForm)
    roleList.value = res.list
    total.value = res.total
  } catch (error) {
    console.error('获取角色列表失败:', error)
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  searchForm.pageNum = 1
  getList()
}

// 重置搜索
const resetSearch = () => {
  searchForm.roleName = ''
  searchForm.status = ''
  handleSearch()
}

// 分页
const handleSizeChange = (val) => {
  searchForm.pageSize = val
  getList()
}

const handleCurrentChange = (val) => {
  searchForm.pageNum = val
  getList()
}

// 角色表单
const dialogVisible = ref(false)
const dialogType = ref('add')
const roleFormRef = ref()
const roleForm = reactive({
  roleId: undefined,
  roleName: '',
  roleKey: '',
  status: '0'
})

// 表单校验规则
const roleRules = {
  roleName: [
    { required: true, message: '请输入角色名称', trigger: 'blur' },
    { min: 2, max: 30, message: '长度在 2 到 30 个字符', trigger: 'blur' }
  ],
  roleKey: [
    { required: true, message: '请输入角色标识', trigger: 'blur' },
    { pattern: /^[a-z][a-z0-9_]*$/, message: '以小写字母开头，只能包含小写字母、数字和下划线', trigger: 'blur' }
  ]
}

// 新增角色
const handleAdd = () => {
  dialogType.value = 'add'
  roleForm.roleId = undefined
  roleForm.roleName = ''
  roleForm.roleKey = ''
  roleForm.status = '0'
  dialogVisible.value = true
}

// 编辑角色
const handleEdit = (row) => {
  dialogType.value = 'edit'
  Object.assign(roleForm, row)
  dialogVisible.value = true
}

// 提交表单
const handleSubmit = async () => {
  if (!roleFormRef.value) return
  
  await roleFormRef.value.validate(async (valid) => {
    if (valid) {
      try {
        if (dialogType.value === 'add') {
          await createRole(roleForm)
          ElMessage.success('新增成功')
        } else {
          await updateRole(roleForm)
          ElMessage.success('修改成功')
        }
        dialogVisible.value = false
        getList()
      } catch (error) {
        console.error('提交失败:', error)
      }
    }
  })
}

// 删除角色
const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确认要删除角色"${row.roleName}"吗？`,
      '警告',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await deleteRole(row.roleId)
    ElMessage.success('删除成功')
    getList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除失败:', error)
    }
  }
}

// 修改角色状态
const handleStatusChange = async (row) => {
  try {
    await updateRoleStatus({
      roleId: row.roleId,
      status: row.status
    })
    ElMessage.success('修改成功')
  } catch (error) {
    console.error('修改状态失败:', error)
    // 恢复状态
    row.status = row.status === '0' ? '1' : '0'
  }
}

// 权限相关
const permissionDialogVisible = ref(false)
const permissionTreeRef = ref()
const permissionList = ref([])
const currentRoleId = ref()
const currentRole = ref(null)

// 获取权限树数据
const getPermissionTreeData = async (roleId) => {
  try {
    const res = await getRolePermissions(roleId)
    // 检查返回的数据结构
    if (res && Array.isArray(res)) {
      permissionList.value = res
      return true
    } else {
      console.error('权限树数据格式错误:', res)
      ElMessage.error('获取权限树失败：数据格式错误')
      return false
    }
  } catch (error) {
    console.error('获取权限树失败:', error)
    ElMessage.error('获取权限树失败：' + (error.message || '未知错误'))
    return false
  }
}

// 打开权限弹窗
const handlePermission = async (row) => {
  currentRole.value = row
  currentRoleId.value = row.roleId
  permissionDialogVisible.value = true
  
  try {
    // 获取角色的权限树
    const success = await getPermissionTreeData(row.roleId)
    if (success) {
      // 根据isChecked字段设置选中状态
      nextTick(() => {
        if (permissionTreeRef.value) {
          // 先清空所有选中状态
          permissionTreeRef.value.setCheckedKeys([])
          // 遍历设置节点状态
          const setNodeStatus = (data) => {
            data.forEach(item => {
              // 如果节点已选中，则勾选该节点
              if (item.isChecked) {
                permissionTreeRef.value.setChecked(item.id, true, false)
              }
              if (item.children && item.children.length) {
                setNodeStatus(item.children)
              }
            })
          }
          setNodeStatus(permissionList.value)
        }
      })
    }
  } catch (error) {
    console.error('获取权限失败:', error)
    ElMessage.error('获取权限失败：' + (error.message || '未知错误'))
  }
}

// 全选/取消全选
const handleCheckAll = (checked) => {
  if (permissionTreeRef.value) {
    // 设置树节点的勾选状态
    permissionTreeRef.value.setCheckedNodes(checked ? permissionList.value : [])
    // 同步更新数据的isChecked状态
    const updateCheckedStatus = (data) => {
      data.forEach(item => {
        const node = permissionTreeRef.value.getNode(item.id)
        if (node) {
          item.isChecked = node.checked
        }
        if (item.children && item.children.length) {
          updateCheckedStatus(item.children)
        }
      })
    }
    updateCheckedStatus(permissionList.value)
  }
}

// 展开/折叠
const handleExpandAll = (expanded) => {
  const setExpanded = (data) => {
    data.forEach(item => {
      item.isExpand = expanded
      if (item.children && item.children.length) {
        setExpanded(item.children)
      }
    })
  }
  setExpanded(permissionList.value)
}

// 提交权限
const handlePermissionSubmit = async () => {
  if (!permissionTreeRef.value || !currentRoleId.value) {
    ElMessage.error('数据未准备就绪，请重试')
    return
  }

  try {
    // 获取当前树的选中状态
    const checkedKeys = permissionTreeRef.value.getCheckedKeys() || []
    const halfCheckedKeys = permissionTreeRef.value.getHalfCheckedKeys() || []
    
    // 更新数据的isChecked状态
    const updateCheckedStatus = (data) => {
      data.forEach(item => {
        item.isChecked = checkedKeys.includes(item.id) || halfCheckedKeys.includes(item.id)
        if (item.children && item.children.length) {
          updateCheckedStatus(item.children)
        }
      })
    }
    updateCheckedStatus(permissionList.value)

    // 提交更新
    const res = await updateRolePermissions({
      roleId: currentRoleId.value,
      permissionIds: [...checkedKeys, ...halfCheckedKeys]
    })
    

      ElMessage.success('权限分配成功')
      permissionDialogVisible.value = false
      // 刷新角色列表
      getList()
  } catch (error) {
    console.error('分配权限失败:', error)
    ElMessage.error('权限分配失败：' + (error.message || '未知错误'))
  }
}

onMounted(() => {
  getList()
})
</script>

<style scoped>
.role-manage {
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

.permission-dialog :deep(.el-dialog__body) {
  padding: 0;
}

.permission-header {
  padding: 15px 20px;
  border-bottom: 1px solid var(--el-border-color-light);
}

.role-info {
  margin-top: 15px;
  display: flex;
  align-items: center;
}

.role-info .label {
  color: var(--el-text-color-secondary);
  margin-right: 8px;
}

.role-info .value {
  font-weight: bold;
  margin-right: 12px;
}

.status-tag {
  font-size: 12px;
}

.permission-content {
  padding: 10px 20px;
}

.tree-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
  padding: 8px 0;
}

.toolbar-left {
  display: flex;
  align-items: center;
  gap: 8px;
}

.toolbar-left .title {
  font-size: 15px;
  font-weight: bold;
}

.toolbar-right {
  display: flex;
  align-items: center;
}

.permission-tree {
  margin: 8px 0;
}

.permission-tree :deep(.el-tree-node__content) {
  height: 36px;
}

.permission-tree :deep(.el-tree-node.is-current > .el-tree-node__content) {
  background-color: var(--el-color-primary-light-9);
}

.dialog-footer {
  padding: 15px 20px;
  border-top: 1px solid var(--el-border-color-light);
}

.custom-tree-node {
  flex: 1;
  display: flex;
  align-items: center;
  font-size: 14px;
}

.permission-code {
  color: var(--el-text-color-secondary);
  font-size: 12px;
  margin-left: 8px;
  max-width: 200px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style> 