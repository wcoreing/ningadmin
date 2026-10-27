<template>
  <div class="user-manage">
    <el-card class="box-card">
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <el-form :inline="true" :model="searchForm">
              <el-form-item label="关键词">
                <el-input
                  v-model="searchForm.keyword"
                  placeholder="请输入用户名、昵称、手机号或邮箱"
                  clearable
                  style="width: 250px"
                  @keyup.enter="handleSearch"
                />
              </el-form-item>
              <el-form-item label="状态">
                <el-select v-model="searchForm.status" placeholder="状态" clearable style="width: 120px">
                  <el-option label="正常" :value="1" />
                  <el-option label="禁用" :value="0" />
                </el-select>
              </el-form-item>
              <el-form-item label="等级">
                <el-select v-model="searchForm.level" placeholder="会员等级" clearable style="width: 120px">
                  <el-option label="普通会员" :value="1" />
                  <el-option label="黄金会员" :value="2" />
                  <el-option label="钻石会员" :value="3" />
                </el-select>
              </el-form-item>
              <el-form-item label="性别">
                <el-select v-model="searchForm.gender" placeholder="性别" clearable style="width: 120px">
                  <el-option label="男" :value="1" />
                  <el-option label="女" :value="0" />
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
        <el-table-column prop="id" label="ID" width="80" align="center" />
        <el-table-column label="头像" width="80" align="center">
          <template #default="{ row }">
            <img v-if="row.avatar" :src="row.avatar" style="width: 40px; height: 40px; border-radius: 50%;">
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="username" label="用户名" width="120" />
        <el-table-column prop="nickname" label="昵称" width="120" />
        <el-table-column label="性别" width="80" align="center">
          <template #default="{ row }">
            {{ row.gender === 1 ? '男' : row.gender === 0 ? '女' : '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="phone" label="手机号" width="120" />
        <el-table-column prop="email" label="邮箱" width="150" />
        <el-table-column label="等级" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="getLevelTagType(row.level)">{{ getLevelText(row.level) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="余额" width="120" align="center">
          <template #default="{ row }">
            <div>
              <span>￥{{ row.balance || 0 }}</span>
              <el-button type="primary" link size="small" @click="handleBalanceAdjust(row)">调整</el-button>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="积分" width="120" align="center">
          <template #default="{ row }">
            <div>
              <span>{{ row.points || 0 }}</span>
              <el-button type="primary" link size="small" @click="handlePointsAdjust(row)">调整</el-button>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-switch
              v-model="row.status"
              :active-value="1"
              :inactive-value="0"
              @change="(value) => handleStatusChange(row, value)"
            />
          </template>
        </el-table-column>
        <el-table-column label="最后登录时间" width="150" align="center">
          <template #default="{ row }">
            <span>{{ formatTime(row.lastLoginTime) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="创建时间" width="150" align="center">
          <template #default="{ row }">
            <span>{{ formatTime(row.createTime) }}</span>
          </template>
        </el-table-column>
        <el-table-column fixed="right" label="操作" width="280">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
            <el-button type="primary" link @click="handleResetPwd(row)">重置密码</el-button>
            <el-button type="info" link @click="handleViewChangeLog(row)">变动记录</el-button>
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
        <el-form-item label="用户名" prop="username">
          <el-input v-model="userForm.username" placeholder="请输入用户名" />
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
        <el-form-item label="昵称" prop="nickname">
          <el-input v-model="userForm.nickname" placeholder="请输入昵称" />
        </el-form-item>
        <el-form-item label="性别">
          <el-radio-group v-model="userForm.gender">
            <el-radio :value="1">男</el-radio>
            <el-radio :value="0">女</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="手机号" prop="phone">
          <el-input v-model="userForm.phone" placeholder="请输入手机号" />
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="userForm.email" placeholder="请输入邮箱" />
        </el-form-item>
        <el-form-item label="生日">
          <el-input v-model="userForm.birthday" placeholder="请输入生日（如：1990-01-01）" />
        </el-form-item>
        <el-form-item label="会员等级">
          <el-select v-model="userForm.level" placeholder="请选择会员等级" style="width: 100%">
            <el-option label="普通会员" :value="1" />
            <el-option label="黄金会员" :value="2" />
            <el-option label="钻石会员" :value="3" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-radio-group v-model="userForm.status">
            <el-radio :value="1">正常</el-radio>
            <el-radio :value="0">禁用</el-radio>
          </el-radio-group>
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

    <!-- 余额调整弹窗 -->
    <el-dialog
      v-model="balanceDialogVisible"
      title="调整用户余额"
      width="500px"
    >
      <el-form
        ref="balanceFormRef"
        :model="balanceForm"
        :rules="balanceRules"
        label-width="100px"
      >
        <el-form-item label="当前余额">
          <span style="color: #409EFF; font-weight: bold;">￥{{ balanceForm.currentBalance }}</span>
        </el-form-item>
        <el-form-item label="调整类型" prop="type">
          <el-radio-group v-model="balanceForm.type">
            <el-radio value="add">增加</el-radio>
            <el-radio value="subtract">减少</el-radio>
            <el-radio value="set">设置</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="调整金额" prop="amount">
          <el-input
            v-model="balanceForm.amount"
            type="number"
            :placeholder="balanceForm.type === 'set' ? '请输入要设置的余额' : '请输入调整金额'"
            min="0"
            step="0.01"
          />
        </el-form-item>
        <el-form-item label="调整原因" prop="remark">
          <el-input
            v-model="balanceForm.remark"
            type="textarea"
            placeholder="请输入调整原因"
            :rows="3"
          />
        </el-form-item>
        <el-form-item label="调整后余额" v-if="balanceForm.type && balanceForm.amount">
          <span style="color: #67C23A; font-weight: bold;">￥{{ calculateNewBalance() }}</span>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="balanceDialogVisible = false">取 消</el-button>
        <el-button type="primary" @click="handleBalanceSubmit">确 定</el-button>
      </template>
    </el-dialog>

    <!-- 积分调整弹窗 -->
    <el-dialog
      v-model="pointsDialogVisible"
      title="调整用户积分"
      width="500px"
    >
      <el-form
        ref="pointsFormRef"
        :model="pointsForm"
        :rules="pointsRules"
        label-width="100px"
      >
        <el-form-item label="当前积分">
          <span style="color: #409EFF; font-weight: bold;">{{ pointsForm.currentPoints }} 分</span>
        </el-form-item>
        <el-form-item label="调整类型" prop="type">
          <el-radio-group v-model="pointsForm.type">
            <el-radio value="add">增加</el-radio>
            <el-radio value="subtract">减少</el-radio>
            <el-radio value="set">设置</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="调整积分" prop="amount">
          <el-input
            v-model="pointsForm.amount"
            type="number"
            :placeholder="pointsForm.type === 'set' ? '请输入要设置的积分' : '请输入调整积分'"
            min="0"
          />
        </el-form-item>
        <el-form-item label="调整原因" prop="remark">
          <el-input
            v-model="pointsForm.remark"
            type="textarea"
            placeholder="请输入调整原因"
            :rows="3"
          />
        </el-form-item>
        <el-form-item label="调整后积分" v-if="pointsForm.type && pointsForm.amount">
          <span style="color: #67C23A; font-weight: bold;">{{ calculateNewPoints() }} 分</span>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="pointsDialogVisible = false">取 消</el-button>
        <el-button type="primary" @click="handlePointsSubmit">确 定</el-button>
      </template>
    </el-dialog>

    <!-- 变动记录弹窗 -->
    <el-dialog
      v-model="changeLogDialogVisible"
      title="用户变动记录"
      width="1000px"
      top="5vh"
    >
      <div class="change-log-header">
        <el-form :inline="true" :model="changeLogSearchForm">
          <el-form-item label="记录类型">
            <el-select v-model="changeLogSearchForm.changeType" placeholder="全部" clearable style="width: 120px">
              <el-option label="余额" :value="1" />
              <el-option label="积分" :value="2" />
            </el-select>
          </el-form-item>
          <el-form-item label="变动行为">
            <el-select v-model="changeLogSearchForm.actionType" placeholder="全部" clearable style="width: 120px">
              <el-option label="收入/增加" :value="1" />
              <el-option label="支出/减少" :value="2" />
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="getChangeLogList">查询</el-button>
            <el-button @click="resetChangeLogSearch">重置</el-button>
          </el-form-item>
        </el-form>
      </div>

      <el-table :data="changeLogList" style="width: 100%" v-loading="changeLogLoading">
        <el-table-column prop="id" label="ID" width="80" align="center" />
        <el-table-column label="记录类型" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="row.changeType === 1 ? 'success' : 'warning'">
              {{ row.changeTypeText }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="变动行为" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="row.actionType === 1 ? 'success' : 'danger'">
              {{ row.actionTypeText }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="变动数量" width="120" align="center">
          <template #default="{ row }">
            <span v-if="row.changeType === 1">￥{{ row.amount }}</span>
            <span v-else>{{ row.amount }} 分</span>
          </template>
        </el-table-column>
        <el-table-column label="变动前" width="120" align="center">
          <template #default="{ row }">
            <span v-if="row.changeType === 1">￥{{ row.beforeValue }}</span>
            <span v-else>{{ row.beforeValue }} 分</span>
          </template>
        </el-table-column>
        <el-table-column label="变动后" width="120" align="center">
          <template #default="{ row }">
            <span v-if="row.changeType === 1">￥{{ row.afterValue }}</span>
            <span v-else>{{ row.afterValue }} 分</span>
          </template>
        </el-table-column>
        <el-table-column prop="sourceTypeText" label="来源类型" width="120" align="center" />
        <el-table-column prop="remark" label="备注" min-width="150" />
        <el-table-column label="创建时间" width="150" align="center">
          <template #default="{ row }">
            <span>{{ formatTime(row.createTime) }}</span>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination" style="margin-top: 20px;">
        <el-pagination
          v-model:current-page="changeLogSearchForm.pageNum"
          v-model:page-size="changeLogSearchForm.pageSize"
          :page-sizes="[10, 20, 30, 50]"
          :total="changeLogTotal"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleChangeLogSizeChange"
          @current-change="handleChangeLogCurrentChange"
        />
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Plus, Refresh } from '@element-plus/icons-vue'
import {
  getUserList,
  getUserDetail,
  createUser,
  updateUser,
  deleteUser,
  updateUserStatus,
  resetUserPassword,
  updateUserBalance,
  updateUserPoints,
  getUserChangeLog
} from '@/admin/api/user'

// 搜索表单
const searchForm = reactive({
  pageNum: 1,
  pageSize: 20,
  keyword: '',
  status: undefined,
  level: undefined,
  gender: undefined
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
  searchForm.keyword = ''
  searchForm.status = undefined
  searchForm.level = undefined
  searchForm.gender = undefined
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
  id: undefined,
  username: '',
  password: '',
  nickname: '',
  gender: null,
  phone: '',
  email: '',
  birthday: '',
  level: 1,
  status: 1
})

// 表单校验规则
const userRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 2, max: 20, message: '长度在 2 到 20 个字符', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于6位', trigger: 'blur' }
  ],
  nickname: [
    { required: true, message: '请输入昵称', trigger: 'blur' }
  ],
  email: [
    { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
  ],
  phone: [
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号码', trigger: 'blur' }
  ]
}

// 新增用户
const handleAdd = () => {
  dialogType.value = 'add'
  userForm.id = undefined
  userForm.username = ''
  userForm.password = ''
  userForm.nickname = ''
  userForm.gender = null
  userForm.phone = ''
  userForm.email = ''
  userForm.birthday = ''
  userForm.level = 1
  userForm.status = 1
  dialogVisible.value = true
}

// 编辑用户
const handleEdit = async (row) => {
  dialogType.value = 'edit'
  try {
    const userInfo = await getUserDetail(row.id)
    // 先清空表单
    Object.keys(userForm).forEach(key => {
      userForm[key] = undefined
    })
    // 设置用户信息
    Object.assign(userForm, userInfo)
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
      `确认要删除用户"${row.username}"吗？`,
      '警告',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await deleteUser(row.id)
    ElMessage.success('删除成功')
    getList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除失败:', error)
    }
  }
}

// 状态切换
const handleStatusChange = async (row, newStatus) => {
  const statusText = newStatus === 1 ? '启用' : '禁用'
  
  try {
    await updateUserStatus({
      id: row.id,
      status: newStatus
    })
    
    ElMessage.success(`${statusText}成功`)
  } catch (error) {
    console.error('修改状态失败:', error)
    ElMessage.error('修改状态失败')
    // 如果失败，恢复原状态
    row.status = row.status === 1 ? 0 : 1
  }
}

// 重置密码
const resetPwdDialogVisible = ref(false)
const resetPwdFormRef = ref()
const resetPwdForm = reactive({
  id: undefined,
  password: ''
})

const resetPwdRules = {
  password: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于6位', trigger: 'blur' }
  ]
}

const handleResetPwd = (row) => {
  resetPwdForm.id = row.id
  resetPwdForm.password = ''
  resetPwdDialogVisible.value = true
}

const handleResetPwdSubmit = async () => {
  if (!resetPwdFormRef.value) return
  
  await resetPwdFormRef.value.validate(async (valid) => {
    if (valid) {
      try {
        await resetUserPassword({
          id: resetPwdForm.id,
          password: resetPwdForm.password
        })
        ElMessage.success('重置密码成功')
        resetPwdDialogVisible.value = false
      } catch (error) {
        console.error('重置密码失败:', error)
      }
    }
  })
}

// 余额调整
const balanceDialogVisible = ref(false)
const balanceFormRef = ref()
const balanceForm = reactive({
  id: undefined,
  currentBalance: 0,
  type: 'add',
  amount: '',
  remark: ''
})

const balanceRules = {
  type: [
    { required: true, message: '请选择调整类型', trigger: 'change' }
  ],
  amount: [
    { required: true, message: '请输入调整金额', trigger: 'blur' },
    { pattern: /^\d+(\.\d{1,2})?$/, message: '请输入正确的金额格式', trigger: 'blur' }
  ],
  remark: [
    { required: true, message: '请输入调整原因', trigger: 'blur' }
  ]
}

const handleBalanceAdjust = (row) => {
  balanceForm.id = row.id
  balanceForm.currentBalance = row.balance || 0
  balanceForm.type = 'add'
  balanceForm.amount = ''
  balanceForm.remark = ''
  balanceDialogVisible.value = true
}

const calculateNewBalance = () => {
  const current = parseFloat(balanceForm.currentBalance) || 0
  const amount = parseFloat(balanceForm.amount) || 0
  
  switch (balanceForm.type) {
    case 'add':
      return (current + amount).toFixed(2)
    case 'subtract':
      return Math.max(0, current - amount).toFixed(2)
    case 'set':
      return amount.toFixed(2)
    default:
      return current.toFixed(2)
  }
}

const handleBalanceSubmit = async () => {
  if (!balanceFormRef.value) return
  
  await balanceFormRef.value.validate(async (valid) => {
    if (valid) {
      try {
        // 计算最终余额
        const finalBalance = parseFloat(calculateNewBalance())
        
        await updateUserBalance({
          id: balanceForm.id,
          balance: finalBalance,
          remark: balanceForm.remark
        })
        
        ElMessage.success('余额调整成功')
        balanceDialogVisible.value = false
        getList() // 刷新列表
      } catch (error) {
        console.error('余额调整失败:', error)
        ElMessage.error('余额调整失败')
      }
    }
  })
}

// 积分调整
const pointsDialogVisible = ref(false)
const pointsFormRef = ref()
const pointsForm = reactive({
  id: undefined,
  currentPoints: 0,
  type: 'add',
  amount: '',
  remark: ''
})

const pointsRules = {
  type: [
    { required: true, message: '请选择调整类型', trigger: 'change' }
  ],
  amount: [
    { required: true, message: '请输入调整积分', trigger: 'blur' },
    { pattern: /^\d+$/, message: '请输入正确的积分数值', trigger: 'blur' }
  ],
  remark: [
    { required: true, message: '请输入调整原因', trigger: 'blur' }
  ]
}

const handlePointsAdjust = (row) => {
  pointsForm.id = row.id
  pointsForm.currentPoints = row.points || 0
  pointsForm.type = 'add'
  pointsForm.amount = ''
  pointsForm.remark = ''
  pointsDialogVisible.value = true
}

const calculateNewPoints = () => {
  const current = parseInt(pointsForm.currentPoints) || 0
  const amount = parseInt(pointsForm.amount) || 0
  
  switch (pointsForm.type) {
    case 'add':
      return current + amount
    case 'subtract':
      return Math.max(0, current - amount)
    case 'set':
      return amount
    default:
      return current
  }
}

const handlePointsSubmit = async () => {
  if (!pointsFormRef.value) return
  
  await pointsFormRef.value.validate(async (valid) => {
    if (valid) {
      try {
        // 计算最终积分
        const finalPoints = calculateNewPoints()
        
        await updateUserPoints({
          id: pointsForm.id,
          points: finalPoints,
          remark: pointsForm.remark
        })
        
        ElMessage.success('积分调整成功')
        pointsDialogVisible.value = false
        getList() // 刷新列表
      } catch (error) {
        console.error('积分调整失败:', error)
        ElMessage.error('积分调整失败')
      }
    }
  })
}

// 工具函数
const getLevelTagType = (level) => {
  const map = {
    1: '',
    2: 'warning',
    3: 'danger'
  }
  return map[level] || ''
}

const getLevelText = (level) => {
  const map = {
    1: '普通会员',
    2: '黄金会员',
    3: '钻石会员'
  }
  return map[level] || '-'
}

// 变动记录相关
const changeLogDialogVisible = ref(false)
const changeLogLoading = ref(false)
const changeLogList = ref([])
const changeLogTotal = ref(0)
const currentUserForLog = ref(null)

const changeLogSearchForm = reactive({
  pageNum: 1,
  pageSize: 20,
  changeType: undefined,
  actionType: undefined
})

// 查看变动记录
const handleViewChangeLog = (row) => {
  currentUserForLog.value = row
  changeLogDialogVisible.value = true
  changeLogSearchForm.pageNum = 1
  changeLogSearchForm.changeType = undefined
  changeLogSearchForm.actionType = undefined
  getChangeLogList()
}

// 获取变动记录列表
const getChangeLogList = async () => {
  if (!currentUserForLog.value) return
  
  changeLogLoading.value = true
  try {
    const params = {
      pageNum: changeLogSearchForm.pageNum,
      pageSize: changeLogSearchForm.pageSize
    }
    
    if (changeLogSearchForm.changeType !== undefined) {
      params.changeType = changeLogSearchForm.changeType
    }
    if (changeLogSearchForm.actionType !== undefined) {
      params.actionType = changeLogSearchForm.actionType
    }
    
    const res = await getUserChangeLog(currentUserForLog.value.id, params)
    console.log('变动记录数据:', res)
    
    if (Array.isArray(res)) {
      changeLogList.value = res
      changeLogTotal.value = res.length
    } else if (res && res.list) {
      changeLogList.value = res.list || []
      changeLogTotal.value = res.total || 0
    } else {
      changeLogList.value = []
      changeLogTotal.value = 0
    }
  } catch (error) {
    console.error('获取变动记录失败:', error)
    ElMessage.error('获取变动记录失败')
  } finally {
    changeLogLoading.value = false
  }
}

// 重置变动记录搜索
const resetChangeLogSearch = () => {
  changeLogSearchForm.changeType = undefined
  changeLogSearchForm.actionType = undefined
  changeLogSearchForm.pageNum = 1
  getChangeLogList()
}

// 变动记录分页
const handleChangeLogSizeChange = (val) => {
  changeLogSearchForm.pageSize = val
  getChangeLogList()
}

const handleChangeLogCurrentChange = (val) => {
  changeLogSearchForm.pageNum = val
  getChangeLogList()
}

const formatTime = (time) => {
  if (!time) return '-'
  return new Date(time).toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

onMounted(() => {
  getList()
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
</style> 