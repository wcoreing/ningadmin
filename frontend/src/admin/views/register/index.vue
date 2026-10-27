<template>
  <div class="register-container">
    <el-card class="register-card">
      <template #header>
        <div class="card-header">
          <span class="logo">集</span>
          <h2>集亿玖泽</h2>
        </div>
      </template>
      
      <el-form 
        ref="registerFormRef"
        :model="registerForm"
        :rules="registerRules"
        label-position="top"
      >
        <el-form-item label="用户名" prop="username">
          <el-input 
            v-model="registerForm.username"
            placeholder="请输入用户名"
            :prefix-icon="User"
          />
        </el-form-item>
        
        <el-form-item label="密码" prop="password">
          <el-input 
            v-model="registerForm.password"
            type="password"
            placeholder="请输入密码"
            :prefix-icon="Lock"
            show-password
          />
        </el-form-item>

        <el-form-item label="确认密码" prop="confirmPassword">
          <el-input 
            v-model="registerForm.confirmPassword"
            type="password"
            placeholder="请再次输入密码"
            :prefix-icon="Lock"
            show-password
          />
        </el-form-item>

        <el-form-item>
          <el-button 
            type="primary" 
            :loading="loading"
            class="register-button"
            @click="handleRegister"
          >
            注册
          </el-button>
          <div class="login-link">
            已有账号？<el-link type="primary" @click="router.push('/admin/login')">去登录</el-link>
          </div>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/admin/stores/user'
import { User, Lock } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'
import { ElMessage } from 'element-plus'

const router = useRouter()
const userStore = useUserStore()
const registerFormRef = ref<FormInstance>()
const loading = ref(false)

const registerForm = reactive({
  username: '',
  password: '',
  confirmPassword: ''
})

const validatePass = (rule: any, value: string, callback: any) => {
  if (value === '') {
    callback(new Error('请输入密码'))
  } else {
    if (registerForm.confirmPassword !== '') {
      registerFormRef.value?.validateField('confirmPassword')
    }
    callback()
  }
}

const validatePass2 = (rule: any, value: string, callback: any) => {
  if (value === '') {
    callback(new Error('请再次输入密码'))
  } else if (value !== registerForm.password) {
    callback(new Error('两次输入密码不一致!'))
  } else {
    callback()
  }
}

const registerRules = reactive<FormRules>({
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '长度在 3 到 20 个字符', trigger: 'blur' }
  ],
  password: [
    { required: true, validator: validatePass, trigger: 'blur' },
    { min: 6, message: '密码长度不能小于6位', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, validator: validatePass2, trigger: 'blur' }
  ]
})

const handleRegister = async () => {
  if (!registerFormRef.value) return
  
  try {
    const valid = await registerFormRef.value.validate()
    if (valid) {
      loading.value = true
      const { confirmPassword, ...registerData } = registerForm
      await userStore.register(registerData)
      ElMessage.success('注册成功，请登录')
      router.push('/admin/login')
    }
  } catch (error) {
    console.error('注册失败:', error)
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.register-container {
  min-height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 20px;
  background-color: #fff;
}

.register-card {
  width: 100%;
  max-width: 480px;
  margin: auto;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.05);
}

.card-header {
  text-align: center;
  padding: 20px 0;
}

.logo {
  width: 60px;
  height: 60px;
  margin: 0 auto 15px;
  background: #e1251b;
  color: #fff;
  display: grid;
  place-items: center;
  border-radius: 8px;
  font-size: 28px;
  font-weight: 700;
}

.register-button {
  width: 100%;
  height: 40px;
  font-size: 16px;
}

.login-link {
  text-align: center;
  margin-top: 15px;
}

:deep(.el-form-item__label) {
  font-size: 15px;
}

:deep(.el-input__wrapper) {
  height: 42px;
}

:deep(.el-card__header) {
  border-bottom: none;
}

:deep(.el-card__body) {
  padding: 0 40px 40px;
}
</style> 