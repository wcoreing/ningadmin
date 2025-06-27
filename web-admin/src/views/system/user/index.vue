<template>
  <div class="user-manage">
    <el-card class="box-card">
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <el-form :inline="true" :model="searchForm">
              <el-form-item label="用户账号">
                <el-input
                  v-model="searchForm.userName"
                  placeholder="请输入用户账号"
                  clearable
                  @keyup.enter="handleSearch"
                />
              </el-form-item>
              <el-form-item label="部门">
                <el-select v-model="searchForm.deptId" placeholder="请选择部门" clearable style="width: 200px">
                  <el-option
                    v-for="item in deptOptions"
                    :key="item.value"
                    :label="item.label"
                    :value="item.value"
                  />
                </el-select>
              </el-form-item>
              <el-form-item label="状态">
                <el-select v-model="searchForm.status" placeholder="用户状态" clearable  style="width: 120px; margin-left: 16px">
                  <el-option label="正常" value="0" />
                  <el-option label="停用" value="1" />
                </el-select>
              </el-form-item>
              <el-form-item>
                <el-button type="primary" :icon="Search" @click="handleSearch">搜索</el-button>
                <el-button :icon="Refresh" @click="resetSearch">重置</el-button>
              </el-form-item>
            </el-form>
          </div>
          <el-button type="primary" :icon="Plus" @click="handleAdd">新增用户</el-button>
        </div>
      </template>

      <el-table v-loading="loading" :data="userList" style="width: 100%">
        <el-table-column prop="userName" label="用户账号" />
        <el-table-column prop="nickName" label="用户昵称" />
        <el-table-column prop="email" label="邮箱" />
        <el-table-column prop="phonenumber" label="手机号码" />
        <el-table-column prop="deptName" label="部门" width="150">
          <template #default="scope">
            <span>{{ scope.row.deptNames ? scope.row.deptNames.join(', ') : scope.row.deptName || '' }}</span>
          </template>
        </el-table-column>
        <el-table-column label="角色">
          <template #default="{ row }">
            {{ row.role?.roleName }}
          </template>
        </el-table-column>
        <el-table-column label="数据权限">
          <template #default="{ row }">
            {{ getDataScopeText(row.dataScope) }}
          </template>
        </el-table-column>
        <el-table-column prop="sex" label="性别">
          <template #default="{ row }">
            {{ row.sex === '0' ? '男' : row.sex === '1' ? '女' : '未知' }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态">
          <template #default="{ row }">
            <el-switch
              v-model="row.status"
              :active-value="'0'"
              :inactive-value="'1'"
              @change="handleStatusChange(row)"
            />
          </template>
        </el-table-column>
        <el-table-column prop="createTime" label="创建时间" />
        <el-table-column fixed="right" label="操作" width="180">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
            <el-button type="primary" link @click="handleResetPwd(row)">重置密码</el-button>
            <el-button type="danger" link @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination">
        <el-pagination
          v-model:current-page="searchForm.pageNum"
          v-model:page-size="searchForm.pageSize"
          :page-sizes="[10, 20, 30, 50]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>

    <!-- 用户表单弹窗 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogType === 'add' ? '新增用户' : '编辑用户'"
      width="600px"
    >
      <el-form
        ref="userFormRef"
        :model="userForm"
        :rules="userRules"
        label-width="100px"
      >
        <el-form-item label="用户账号" prop="userName">
          <el-input v-model="userForm.userName" placeholder="请输入用户账号" />
        </el-form-item>
        <el-form-item 
          label="密码" 
          prop="password"
          v-if="dialogType === 'add'"
        >
          <el-input
            v-model="userForm.password"
            type="password"
            placeholder="请输入密码"
            show-password
          />
        </el-form-item>
        <el-form-item label="用户昵称" prop="nickName">
          <el-input v-model="userForm.nickName" placeholder="请输入用户昵称" />
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="userForm.email" placeholder="请输入邮箱" />
        </el-form-item>
        <el-form-item label="手机号码" prop="phonenumber">
          <el-input v-model="userForm.phonenumber" placeholder="请输入手机号码" />
        </el-form-item>
        <el-form-item label="性别">
          <el-radio-group v-model="userForm.sex">
            <el-radio value="0">男</el-radio>
            <el-radio value="1">女</el-radio>
            <el-radio value="2">未知</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="状态">
          <el-radio-group v-model="userForm.status">
            <el-radio value="0">正常</el-radio>
            <el-radio value="1">停用</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="部门权限">
          <div class="dept-permission-container">
            <el-button type="primary" size="small" @click="addDeptPermission">添加部门</el-button>
            <div v-for="(item, index) in deptPermissions" :key="index" class="dept-permission-item">
              <el-select v-model="item.deptId" placeholder="请选择部门" style="width: 200px; margin-right: 10px;">
                <el-option
                  v-for="dept in deptOptions"
                  :key="dept.value"
                  :label="dept.label"
                  :value="dept.value"
                />
              </el-select>
              <el-select v-model="item.dataScope" placeholder="请选择数据权限" style="width: 200px; margin-right: 10px;">
                <el-option label="全部数据权限" value="1" />
                <el-option label="自定数据权限" value="2" />
                <el-option label="本部门数据权限" value="3" />
                <el-option label="本部门及以下数据权限" value="4" />
                <el-option label="仅本人数据权限" value="5" />
              </el-select>
              <el-button type="danger" size="small" @click="removeDeptPermission(index)">删除</el-button>
            </div>
          </div>
        </el-form-item>
        <el-form-item label="角色" prop="roleId">
          <el-select v-model="userForm.roleId" placeholder="请选择角色" clearable>
            <el-option
              v-for="item in roleOptions"
              :key="item.roleId"
              :label="item.roleName"
              :value="item.roleId"
            />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取 消</el-button>
        <el-button type="primary" @click="handleSubmit">确 定</el-button>
      </template>
    </el-dialog>

    <!-- 重置密码弹窗 -->
    <el-dialog
      v-model="resetPwdDialogVisible"
      title="重置密码"
      width="500px"
    >
      <el-form
        ref="resetPwdFormRef"
        :model="resetPwdForm"
        :rules="resetPwdRules"
        label-width="100px"
      >
        <el-form-item label="新密码" prop="password">
          <el-input
            v-model="resetPwdForm.password"
            type="password"
            placeholder="请输入新密码"
            show-password
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="resetPwdDialogVisible = false">取 消</el-button>
        <el-button type="primary" @click="handleResetPwdSubmit">确 定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Plus, Refresh } from '@element-plus/icons-vue'
import {
  getUserList,
  getUserInfo,
  createUser,
  updateUser,
  deleteUser,
  updateUserStatus,
  resetUserPassword,
  getDeptOptions
} from '@/api/system/user'
import { getRoleList, getUserRoles } from '@/api/system/role'

// 搜索表单
const searchForm = reactive({
  pageNum: 1,
  pageSize: 10,
  userName: '',
  status: undefined,
  deptId: undefined
})

// 列表数据
const loading = ref(false)
const userList = ref([])
const total = ref(0)

// 获取用户列表
const getList = async () => {
  loading.value = true
  try {
    // 构建查询参数，过滤掉空值
    const params = {}
    Object.keys(searchForm).forEach(key => {
      if (searchForm[key] !== '' && searchForm[key] !== undefined) {
        params[key] = searchForm[key]
      }
    })
    
    const res = await getUserList(params)
    console.log('用户列表数据:', res)
    
    // 直接处理返回的数据
    if (Array.isArray(res)) {
      userList.value = res
      total.value = res.length
    } else if (res) {
      userList.value = res.list || res.records || []
      total.value = res.total || 0
    } else {
      userList.value = []
      total.value = 0
    }
  } catch (error) {
    console.error('获取用户列表失败:', error)
    ElMessage.error('获取用户列表失败')
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
  searchForm.userName = ''
  searchForm.status = undefined
  searchForm.deptId = undefined
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

// 用户表单
const dialogVisible = ref(false)
const dialogType = ref('add')
const userFormRef = ref()
const userForm = reactive({
  userId: undefined,
  userName: '',
  password: '',
  nickName: '',
  email: '',
  phonenumber: '',
  sex: '0',
  status: '0',
  deptIds: [],
  dataScopes: [],
  roleId: undefined
})

// 表单校验规则
const userRules = {
  userName: [
    { required: true, message: '请输入用户账号', trigger: 'blur' },
    { min: 2, max: 20, message: '长度在 2 到 20 个字符', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' }
  ],
  nickName: [
    { required: true, message: '请输入用户昵称', trigger: 'blur' }
  ],
  email: [
    { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
  ],
  phonenumber: [
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号码', trigger: 'blur' }
  ],
  roleId: [
    { required: true, message: '请选择角色', trigger: 'change' }
  ]
}

// 新增用户
const handleAdd = () => {
  dialogType.value = 'add'
  userForm.userId = undefined
  userForm.userName = ''
  userForm.password = ''
  userForm.nickName = ''
  userForm.email = ''
  userForm.phonenumber = ''
  userForm.sex = '0'
  userForm.status = '0'
  userForm.deptIds = []
  userForm.dataScopes = []
  userForm.roleId = undefined
  deptPermissions.value = []
  getRoles()
  dialogVisible.value = true
}

// 编辑用户
const handleEdit = async (row) => {
  dialogType.value = 'edit'
  try {
    const userInfo = await getUserInfo(row.id || row.userId)
    // 先清空表单
    Object.keys(userForm).forEach(key => {
      userForm[key] = undefined
    })
    // 设置用户信息
    Object.assign(userForm, {
      ...userInfo,
      userId: userInfo.id || userInfo.userId
    })
    // 设置角色ID
    if (userInfo.role) {
      userForm.roleId = userInfo.role.roleId
    }
    
    // 设置部门权限信息
    deptPermissions.value = []
    if (userInfo.deptIds && userInfo.deptIds.length > 0) {
      for (let i = 0; i < userInfo.deptIds.length; i++) {
        deptPermissions.value.push({
          deptId: userInfo.deptIds[i],
          dataScope: userInfo.dataScopes && userInfo.dataScopes[i] ? userInfo.dataScopes[i] : '3'
        })
      }
    }
    
    await getRoles()
    dialogVisible.value = true
  } catch (error) {
    console.error('获取用户信息失败:', error)
    ElMessage.error('获取用户信息失败')
  }
}

// 提交表单
const handleSubmit = async () => {
  if (!userFormRef.value) return
  
  await userFormRef.value.validate(async (valid) => {
    if (valid) {
      try {
        const submitData = { ...userForm }
        if (dialogType.value === 'edit') {
          delete submitData.password
        }
        
        // 从部门权限数组中提取部门ID和数据权限
        submitData.deptIds = deptPermissions.value.map(item => item.deptId).filter(id => id !== undefined)
        submitData.dataScopes = deptPermissions.value.map(item => item.dataScope)
        
        await (dialogType.value === 'add' ? createUser(submitData) : updateUser(submitData))
        ElMessage.success(dialogType.value === 'add' ? '新增成功' : '修改成功')
        dialogVisible.value = false
        getList()
      } catch (error) {
        console.error('提交失败:', error)
      }
    }
  })
}

// 删除用户
const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确认要删除用户"${row.userName}"吗？`,
      '警告',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await deleteUser(row.userId)
    ElMessage.success('删除成功')
    getList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除失败:', error)
    }
  }
}

// 修改用户状态
const handleStatusChange = async (row) => {
  try {
    await updateUserStatus({
      userId: row.userId,
      status: row.status
    })
    ElMessage.success('修改成功')
  } catch (error) {
    console.error('修改状态失败:', error)
    // 恢复状态
    row.status = row.status === '0' ? '1' : '0'
  }
}

// 重置密码
const resetPwdDialogVisible = ref(false)
const resetPwdFormRef = ref()
const resetPwdForm = reactive({
  userId: undefined,
  password: ''
})

const resetPwdRules = {
  password: [
    { required: true, message: '请输入新密码', trigger: 'blur' }
  ]
}

const handleResetPwd = (row) => {
  resetPwdForm.userId = row.userId
  resetPwdForm.password = ''
  resetPwdDialogVisible.value = true
}

const handleResetPwdSubmit = async () => {
  if (!resetPwdFormRef.value) return
  
  await resetPwdFormRef.value.validate(async (valid) => {
    if (valid) {
      try {
        await resetUserPassword(resetPwdForm)
        ElMessage.success('重置密码成功')
        resetPwdDialogVisible.value = false
      } catch (error) {
        console.error('重置密码失败:', error)
      }
    }
  })
}

const roleOptions = ref([])
const deptOptions = ref([])

// 部门权限管理
const deptPermissions = ref([])

const addDeptPermission = () => {
  deptPermissions.value.push({
    deptId: undefined,
    dataScope: '3' // 默认本部门数据权限
  })
}

const removeDeptPermission = (index) => {
  deptPermissions.value.splice(index, 1)
}

const getRoles = async () => {
  try {
    const res = await getRoleList({
      pageNum: 1,
      pageSize: 100,
      status: '0'
    })
    roleOptions.value = res?.list || res?.records || []
  } catch (error) {
    console.error('获取角色列表失败:', error)
  }
}

// 获取部门选项
const getDeptOptionsList = async () => {
  try {
    const res = await getDeptOptions()
    deptOptions.value = res || []
  } catch (error) {
    console.error('获取部门选项失败:', error)
  }
}

// 获取数据权限文本
const getDataScopeText = (dataScope) => {
  const dataScopeMap = {
    '1': '全部数据权限',
    '2': '自定数据权限',
    '3': '本部门数据权限',
    '4': '本部门及以下数据权限',
    '5': '仅本人数据权限'
  }
  return dataScopeMap[dataScope] || '未知'
}

onMounted(() => {
  getList()
  getDeptOptionsList()
})
</script>

<style scoped>
.user-manage {
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

.dept-permission-container {
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  padding: 15px;
  background-color: #fafafa;
}

.dept-permission-item {
  display: flex;
  align-items: center;
  margin-top: 10px;
}

.dept-permission-item:first-child {
  margin-top: 0;
}
</style> 