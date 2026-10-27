<template>
  <div class="login-container">
    <div class="login-box">
      <!-- 左侧 LOGO 区域 -->
      <div class="logo-section">
        <div class="logo-content">
          <div class="brand-area">
            <span class="logo">集</span>
            <div class="brand-text">
              <h1 class="title">集亿玖泽</h1>
              <p class="subtitle">首饰包装与箱包</p>
            </div>
          </div>
          
          <div class="feature-list">
            <div class="feature-column">
              <div class="feature-item">
                <el-icon><Monitor /></el-icon>
                <span>产品目录</span>
              </div>
              <div class="feature-item">
                <el-icon><Operation /></el-icon>
                <span>内容资料</span>
              </div>
            </div>
            <div class="feature-column">
              <div class="feature-item">
                <el-icon><DataLine /></el-icon>
                <span>询盘跟进</span>
              </div>
              <div class="feature-item">
                <el-icon><Connection /></el-icon>
                <span>账号权限</span>
              </div>
            </div>
          </div>
      </div>
      </div>

      <!-- 右侧登录表单区域 -->
      <div class="form-section">
        <div class="form-header">
          <h2>账号登录</h2>
          <div class="form-decoration">
            <span class="line"></span>
            <span class="dot"></span>
            <span class="line"></span>
      </div>
    </div>
    
    <el-form
      ref="loginFormRef"
      :model="loginForm"
      :rules="loginRules"
      class="login-form"
        >
          <el-form-item prop="username">
            <el-input
              ref="usernameRef"
              v-model="loginForm.username"
              placeholder="请输入用户名"
              :prefix-icon="User"
            />
        </el-form-item>

          <el-form-item prop="password">
            <el-input
              ref="passwordRef"
              v-model="loginForm.password"
              :type="passwordVisible ? 'text' : 'password'"
              placeholder="请输入密码"
              :prefix-icon="Lock"
              @keyup.enter="handleLogin"
            >
              <template #suffix>
                <el-icon 
                  class="cursor-pointer"
                  @click="passwordVisible = !passwordVisible"
                >
                  <View v-if="passwordVisible" />
                  <Hide v-else />
                </el-icon>
              </template>
            </el-input>
        </el-form-item>

        <el-button
          :loading="loading"
          type="primary"
          class="login-button"
          @click="handleLogin"
        >
          <span class="button-text">{{ loading ? '登录中...' : '登 录' }}</span>
            <div class="button-effect"></div>
        </el-button>
        </el-form>

        <div class="form-footer">
          <p class="copyright">深圳市集亿玖泽贸易有限公司</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { 
  User, 
  Lock, 
  View, 
  Hide, 
  Monitor, 
  Operation, 
  DataLine, 
  Connection 
} from '@element-plus/icons-vue'
import { login } from '@/admin/api/auth'
import { useUserStore } from '@/admin/stores/user'

const router = useRouter()
const userStore = useUserStore()

const loginFormRef = ref()
const loading = ref(false)
const passwordVisible = ref(false)

const loginForm = reactive({
  username: '',
  password: ''
})

const loginRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' }
  ]
}

const handleLogin = async () => {
  if (!loginFormRef.value) return
  
  try {
    loading.value = true
    const valid = await loginFormRef.value.validate()
    if (valid) {
      const res = await login(loginForm)
      userStore.setToken(res.token)
      await userStore.getUserInfo()
      ElMessage.success('登录成功')
      location.href = '/admin/dashboard'
    }
  } catch (error) {
    console.error('登录失败:', error)
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  /* 使用更丰富的渐变背景 */
  background: linear-gradient(135deg, 
    #0c1f4d 0%,
    #1a3a7e 30%,
    #2c4f9c 50%,
    #3666b8 70%,
    #4b7fd4 100%
  );
  background-size: 200% 200%;
  animation: gradientBG 15s ease infinite;
  padding: 20px;
}

/* 背景动画 */
@keyframes gradientBG {
  0% { background-position: 0% 50% }
  50% { background-position: 100% 50% }
  100% { background-position: 0% 50% }
}

.login-box {
  width: 1000px;
  height: 560px;
  background: rgba(255, 255, 255, 0.95);
  border-radius: 16px;
  display: flex;
  overflow: hidden;
  box-shadow: 
    0 10px 40px rgba(0, 0, 0, 0.2),
    0 0 100px rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.1);
}

/* 左侧 LOGO 区域 */
.logo-section {
  width: 60%;
  background: linear-gradient(135deg, #e1251b 0%, #8d1510 100%);
  padding: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  position: relative;
  overflow: hidden;
}

/* 添加动态效果 */
.logo-section::before {
  content: '';
  position: absolute;
  top: 0;
  left: -50%;
  width: 200%;
  height: 100%;
  background: linear-gradient(
    to right,
    transparent,
    rgba(255, 255, 255, 0.1),
    transparent
  );
  transform: skewX(-45deg);
  animation: shine 3s infinite;
}

@keyframes shine {
  0% { transform: skewX(-45deg) translateX(-150%) }
  100% { transform: skewX(-45deg) translateX(150%) }
}

.logo-content {
  width: 100%;
  position: relative;
  z-index: 1;
}

.brand-area {
  display: flex;
  align-items: center;
  margin-bottom: 60px;
}

.logo {
  width: 90px;
  height: 90px;
  border-radius: 20px;
  box-shadow: 0 12px 30px rgba(0, 0, 0, 0.3);
  margin-right: 24px;
  flex-shrink: 0;
  background: #fff;
  color: #e1251b;
  display: grid;
  place-items: center;
  font-size: 42px;
  font-weight: 700;
}

.brand-text {
  text-align: left;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.title {
  font-size: 42px;
  font-weight: 600;
  margin-bottom: 8px;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
  line-height: 1.2;
  white-space: nowrap;
}

.subtitle {
  font-size: 20px;
  opacity: 0.9;
  font-weight: 500;
  letter-spacing: 1px;
  white-space: nowrap;
}

.feature-list {
  display: flex;
  justify-content: space-between;
  gap: 30px;
  margin-top: 20px;
}

.feature-column {
  flex: 1;
}

.feature-item {
  display: flex;
  align-items: center;
  margin-bottom: 30px;
  font-size: 18px;
  opacity: 0.9;
  transition: all 0.3s ease;
  padding: 14px;
  border-radius: 14px;
  background: rgba(255, 255, 255, 0.1);
}

.feature-item:hover {
  transform: translateY(-3px);
  opacity: 1;
  background: rgba(255, 255, 255, 0.15);
}

.feature-item .el-icon {
  font-size: 28px;
  margin-right: 14px;
  background: rgba(255, 255, 255, 0.2);
  padding: 10px;
  border-radius: 10px;
}

/* 右侧表单区域 */
.form-section {
  width: 40%;
  padding: 40px;
  display: flex;
  flex-direction: column;
  background: linear-gradient(135deg, 
    rgba(255, 255, 255, 0.95) 0%,
    rgba(255, 255, 255, 0.98) 100%
  );
  position: relative;
  overflow: hidden;
}

.form-section::before {
  content: '';
  position: absolute;
  top: 0;
  right: 0;
  width: 100%;
  height: 100%;
  background: linear-gradient(45deg,
    transparent 0%,
    rgba(0, 198, 251, 0.03) 50%,
    rgba(0, 91, 234, 0.03) 100%
  );
  z-index: 0;
}

.form-header {
  text-align: center;
  margin-bottom: 40px;
  position: relative;
  z-index: 1;
}

.form-header h2 {
  font-size: 32px;
  background: linear-gradient(135deg, #00c6fb, #005bea);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
  margin-bottom: 15px;
  font-weight: 600;
}

.form-decoration {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  margin-top: 10px;
}

.form-decoration .line {
  width: 40px;
  height: 2px;
  background: linear-gradient(90deg, #00c6fb, #005bea);
  border-radius: 2px;
}

.form-decoration .dot {
  width: 6px;
  height: 6px;
  background: #005bea;
  border-radius: 50%;
}

.login-form {
  position: relative;
  z-index: 1;
  width: 100%;
  max-width: 320px;
  margin: 0 auto;

  .el-form-item {
  margin-bottom: 25px;
  }

  :deep(.el-input__wrapper) {
    background: rgba(255, 255, 255, 0.9);
    border: 1px solid rgba(0, 198, 251, 0.2);
    border-radius: 12px;
  padding: 0 16px;
    height: 50px;
    box-shadow: 0 4px 12px rgba(0, 91, 234, 0.1) !important;
    transition: all 0.3s ease;

    &:hover, &:focus-within {
      border-color: rgba(0, 198, 251, 0.5);
      box-shadow: 0 6px 16px rgba(0, 91, 234, 0.15) !important;
  transform: translateY(-1px);
}
  }

  :deep(.el-input__inner) {
    height: 50px;
  font-size: 15px;
  color: #2c3e50;
  }

  :deep(.el-input__prefix) {
    font-size: 20px;
    color: #00c6fb;
  }

  :deep(.el-input__suffix) {
    color: #00c6fb;
    font-size: 20px;
  }
}

.login-button {
  width: 100%;
  height: 50px;
  border-radius: 12px;
  font-size: 18px;
  font-weight: 600;
  letter-spacing: 4px;
  background: linear-gradient(45deg, #00c6fb, #005bea);
  border: none;
  margin-top: 30px;
  position: relative;
  overflow: hidden;
  transition: all 0.3s ease;

  &:hover {
  transform: translateY(-2px);
    box-shadow: 0 8px 25px rgba(0, 91, 234, 0.3);
}

  &:active {
  transform: translateY(0);
}

.button-text {
  position: relative;
  z-index: 1;
}

.button-effect {
  position: absolute;
  top: 0;
  left: 0;
  width: 200%;
  height: 100%;
  background: linear-gradient(90deg, 
      transparent 0%,
      rgba(255, 255, 255, 0.2) 50%,
      transparent 100%
    );
    transform: translateX(-100%);
    animation: buttonEffect 3s infinite;
  }
}

@keyframes buttonEffect {
  0% { transform: translateX(-100%) }
  50% { transform: translateX(100%) }
  100% { transform: translateX(100%) }
}

.form-footer {
  margin-top: auto;
  text-align: center;
  position: relative;
  z-index: 1;
}

.copyright {
  color: #7f8c8d;
  font-size: 14px;
  opacity: 0.8;
}

/* 响应式处理 */
@media (max-width: 768px) {
  .login-box {
    width: 100%;
    max-width: 400px;
    height: auto;
    flex-direction: column;
  }

  .logo-section {
    width: 100%;
    padding: 40px;
  }

  .form-section {
    width: 100%;
    padding: 30px;
  }

  .brand-area {
    flex-direction: column;
    text-align: center;
    margin-bottom: 40px;
  }

  .logo {
    width: 100px;
    height: 100px;
    margin-right: 0;
    margin-bottom: 15px;
  }

  .brand-text {
    text-align: center;
  }

  .feature-list {
    flex-direction: column;
    gap: 15px;
  }

  .feature-item {
    margin-bottom: 20px;
  }

  .feature-item .el-icon {
    font-size: 24px;
    padding: 8px;
  }

  .title {
    font-size: 36px;
  }

  .subtitle {
    font-size: 20px;
  }
}
</style> 