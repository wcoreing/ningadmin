<template>
  <div class="login-container">
    <!-- 科技感背景 -->
    <div class="tech-background">
      <div class="light-beams">
        <div class="beam"></div>
        <div class="beam"></div>
        <div class="beam"></div>
        <div class="beam"></div>
        <div class="beam"></div>
        <div class="beam"></div>
        <div class="beam"></div>
        <div class="beam"></div>
      </div>
      <div class="shapes">
        <div class="shape"></div>
        <div class="shape"></div>
        <div class="shape"></div>
        <div class="shape"></div>
        <div class="shape"></div>
        <div class="shape"></div>
        <div class="shape"></div>
        <div class="shape"></div>
      </div>
      <div class="lines">
        <div class="line"></div>
        <div class="line"></div>
        <div class="line"></div>
        <div class="line"></div>
      </div>
    </div>
    
    <el-form
      ref="loginFormRef"
      :model="loginForm"
      :rules="loginRules"
      class="login-form"
      autocomplete="on"
      label-position="left"
    >
      <div class="title-container">
        <div class="logo">
          <div class="wn-text">WN</div>
          <h3 class="title">后台管理系统</h3>
        </div>
      </div>

      <div class="form-content">
        <el-form-item prop="username" class="form-item">
          <div class="input-wrapper">
            <el-input
              ref="usernameRef"
              v-model="loginForm.username"
              placeholder="请输入用户名"
              type="text"
              tabindex="1"
              autocomplete="on"
            >
              <template #prefix>
                <el-icon><User /></el-icon>
              </template>
            </el-input>
          </div>
        </el-form-item>

        <el-form-item prop="password" class="form-item">
          <div class="input-wrapper">
            <el-input
              ref="passwordRef"
              v-model="loginForm.password"
              placeholder="请输入密码"
              :type="passwordVisible ? 'text' : 'password'"
              tabindex="2"
              autocomplete="on"
              @keyup.enter="handleLogin"
            >
              <template #prefix>
                <el-icon><Lock /></el-icon>
              </template>
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
          </div>
        </el-form-item>

        <el-button
          :loading="loading"
          type="primary"
          class="login-button"
          @click="handleLogin"
        >
          <span class="button-text">{{ loading ? '登录中...' : '登 录' }}</span>
          <span class="button-effect"></span>
        </el-button>
      </div>
    </el-form>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { User, Lock, View, Hide } from '@element-plus/icons-vue'
import { login } from '@/api/auth'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const userStore = useUserStore()

const loginFormRef = ref()
const loading = ref(false)
const passwordVisible = ref(false)

// 登录表单
const loginForm = reactive({
  username: '',
  password: ''
})

// 表单校验规则
const loginRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'submit' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'submit' }
  ]
}

// 登录方法
const handleLogin = async () => {
  if (!loginFormRef.value) return
  
  try {
    loading.value = true
    // 手动触发表单验证
    const valid = await loginFormRef.value.validate()
    if (valid) {
      // 调用登录接口
      const res = await login(loginForm)
      // 保存token
      userStore.setToken(res.token)
      // 获取用户信息
      await userStore.getUserInfo()
      // 登录成功提示
      ElMessage.success('登录成功')
      // 跳转到首页
      location.href = '/dashboard'
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
  width: 100%;
  background: #000;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
}

.tech-background {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  overflow: hidden;
  z-index: 0;
  background: radial-gradient(circle at center,
    #ffffff 0%,
    rgba(240, 248, 255, 0.95) 8%,
    rgba(176, 224, 255, 0.9) 16%,
    rgba(100, 181, 246, 0.7) 30%,
    rgba(30, 136, 229, 0.5) 50%,
    rgba(21, 101, 192, 0.4) 70%,
    rgba(13, 71, 161, 0.3) 100%
  );
}

/* 中心光晕 */
.tech-background::before {
  content: '';
  position: absolute;
  width: 80vw;
  height: 80vh;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  background: radial-gradient(
    circle at center,
    rgba(255, 255, 255, 0.1) 0%,
    rgba(255, 255, 255, 0.06) 20%,
    rgba(255, 255, 255, 0.03) 40%,
    transparent 60%
  );
  filter: blur(40px);
  animation: pulseLight 12s ease-in-out infinite;
}

/* 柔和光线 */
.tech-background::after {
  content: '';
  position: absolute;
  width: 200vw;
  height: 200vh;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  background: conic-gradient(
    from 0deg at 50% 50%,
    transparent 0deg,
    rgba(255, 255, 255, 0.02) 20deg,
    transparent 40deg,
    rgba(255, 255, 255, 0.02) 60deg,
    transparent 80deg
  );
  animation: rotateLight 60s linear infinite;
}

/* 移除之前的光束和形状动画 */
.light-beams, .shapes, .lines {
  display: none;
}

.beam {
  position: absolute;
  top: 50%;
  left: 50%;
  width: 2px;
  height: 100vh;
  background: linear-gradient(to top, 
    transparent,
    rgba(52, 152, 219, 0.24),
    rgba(52, 152, 219, 0.4),
    rgba(52, 152, 219, 0.24),
    transparent
  );
  transform-origin: center top;
}

.beam:nth-child(1) { transform: rotate(0deg); animation: beamGlow 3s infinite; }
.beam:nth-child(2) { transform: rotate(45deg); animation: beamGlow 3s infinite 0.5s; }
.beam:nth-child(3) { transform: rotate(90deg); animation: beamGlow 3s infinite 1s; }
.beam:nth-child(4) { transform: rotate(135deg); animation: beamGlow 3s infinite 1.5s; }
.beam:nth-child(5) { transform: rotate(180deg); animation: beamGlow 3s infinite 2s; }
.beam:nth-child(6) { transform: rotate(225deg); animation: beamGlow 3s infinite 2.5s; }
.beam:nth-child(7) { transform: rotate(270deg); animation: beamGlow 3s infinite 1.2s; }
.beam:nth-child(8) { transform: rotate(315deg); animation: beamGlow 3s infinite 0.8s; }

@keyframes beamGlow {
  0%, 100% {
    opacity: 0.24;
    filter: blur(5px);
  }
  50% {
    opacity: 0.56;
    filter: blur(3px);
  }
}

/* 几何图形 */
.shapes {
  position: absolute;
  width: 100%;
  height: 100%;
  top: 0;
  left: 0;
  opacity: 0.4;
}

.shape {
  position: absolute;
  animation: shapeFloat 10s infinite ease-in-out;
}

/* 方形 */
.shape:nth-child(1) {
  top: 15%;
  left: 15%;
  width: 60px;
  height: 60px;
  border: 2px solid rgba(80, 161, 226, 0.24);
  transform: rotate(45deg);
  border-radius: 15px;
  animation-delay: 0s;
}

/* 五边形 */
.shape:nth-child(2) {
  top: 25%;
  right: 20%;
  width: 40px;
  height: 40px;
  background: rgba(80, 161, 226, 0.08);
  clip-path: polygon(50% 0%, 100% 38%, 82% 100%, 18% 100%, 0% 38%);
  animation-delay: 2s;
}

/* 三角形 */
.shape:nth-child(3) {
  bottom: 20%;
  left: 25%;
  width: 50px;
  height: 50px;
  background: rgba(124, 182, 229, 0.08);
  clip-path: polygon(50% 0%, 0% 100%, 100% 100%);
  animation-delay: 4s;
}

/* 圆形 */
.shape:nth-child(4) {
  bottom: 30%;
  right: 25%;
  width: 45px;
  height: 45px;
  border-radius: 50%;
  border: 2px dashed rgba(80, 161, 226, 0.24);
  animation-delay: 6s;
}

/* 六边形 */
.shape:nth-child(5) {
  top: 45%;
  left: 10%;
  width: 35px;
  height: 40px;
  background: rgba(80, 161, 226, 0.08);
  clip-path: polygon(50% 0%, 100% 25%, 100% 75%, 50% 100%, 0% 75%, 0% 25%);
  animation-delay: 1s;
}

/* 八边形 */
.shape:nth-child(6) {
  top: 60%;
  right: 15%;
  width: 40px;
  height: 40px;
  background: rgba(124, 182, 229, 0.08);
  clip-path: polygon(30% 0%, 70% 0%, 100% 30%, 100% 70%, 70% 100%, 30% 100%, 0% 70%, 0% 30%);
  animation-delay: 3s;
}

/* 菱形 */
.shape:nth-child(7) {
  top: 75%;
  left: 50%;
  width: 30px;
  height: 30px;
  border: 2px solid rgba(80, 161, 226, 0.24);
  transform: rotate(45deg);
  animation-delay: 5s;
}

/* 半圆 */
.shape:nth-child(8) {
  top: 20%;
  left: 40%;
  width: 50px;
  height: 25px;
  background: rgba(124, 182, 229, 0.08);
  border-radius: 25px 25px 0 0;
  animation-delay: 7s;
}

@keyframes shapeFloat {
  0%, 100% {
    transform: translate(0, 0) rotate(0deg);
  }
  25% {
    transform: translate(-15px, 15px) rotate(5deg);
  }
  50% {
    transform: translate(15px, -15px) rotate(-5deg);
  }
  75% {
    transform: translate(15px, 15px) rotate(5deg);
  }
}

/* 添加流动的线条 */
.lines {
  position: absolute;
  width: 100%;
  height: 100%;
  top: 0;
  left: 0;
  opacity: 0.24;
}

.line {
  position: absolute;
  width: 1px;
  height: 50vh;
  background: linear-gradient(to bottom,
    transparent,
    rgba(80, 161, 226, 0.4),
    rgba(124, 182, 229, 0.4),
    transparent
  );
  filter: blur(1px);
  animation: lineMove 8s infinite;
}

.line:nth-child(1) { left: 20%; animation-delay: 0s; }
.line:nth-child(2) { left: 40%; animation-delay: 2s; }
.line:nth-child(3) { left: 60%; animation-delay: 4s; }
.line:nth-child(4) { left: 80%; animation-delay: 6s; }

@keyframes lineMove {
  0% {
    transform: translateY(-100%) translateX(-50px) rotate(-5deg);
    opacity: 0;
  }
  50% {
    transform: translateY(100%) translateX(50px) rotate(5deg);
    opacity: 1;
  }
  100% {
    transform: translateY(200%) translateX(-50px) rotate(-5deg);
    opacity: 0;
  }
}

.login-form {
  position: relative;
  width: 360px;
  padding: 0;
  background: rgba(255, 255, 255, 0.85);
  border-radius: 24px;
  box-shadow: 
    0 10px 30px rgba(0, 0, 0, 0.1),
    0 20px 60px rgba(0, 0, 0, 0.05),
    0 4px 8px rgba(0, 0, 0, 0.04),
    0 3px 6px rgba(219, 234, 254, 0.4),
    inset 0 2px 4px rgba(255, 255, 255, 0.8),
    inset 0 -2px 4px rgba(0, 0, 0, 0.05);
  backdrop-filter: blur(10px);
  z-index: 1;
  animation: formAppear 0.6s ease-out;
  overflow: hidden;
  border: 1px solid rgba(255, 255, 255, 0.3);
  transform: perspective(1000px) translateZ(0);
  transition: all 0.3s ease;
}

.login-form:hover {
  transform: perspective(1000px) translateZ(10px);
  box-shadow: 
    0 15px 35px rgba(0, 0, 0, 0.12),
    0 25px 65px rgba(0, 0, 0, 0.07),
    0 5px 10px rgba(0, 0, 0, 0.05),
    0 4px 8px rgba(219, 234, 254, 0.5),
    inset 0 2px 4px rgba(255, 255, 255, 0.9),
    inset 0 -2px 4px rgba(0, 0, 0, 0.06);
}

.login-form::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: linear-gradient(
    135deg,
    rgba(255, 255, 255, 0.4) 0%,
    rgba(255, 255, 255, 0.2) 50%,
    rgba(255, 255, 255, 0.1) 100%
  );
  opacity: 0.5;
  pointer-events: none;
}

.login-form::after {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: radial-gradient(
    circle at 30% 30%,
    rgba(255, 255, 255, 0.4) 0%,
    transparent 70%
  );
  opacity: 0.6;
  pointer-events: none;
  animation: lightMove 8s infinite alternate ease-in-out;
}

@keyframes lightMove {
  0% {
    opacity: 0.4;
    transform: translate(-5%, -5%) scale(0.95);
  }
  100% {
    opacity: 0.6;
    transform: translate(5%, 5%) scale(1.05);
  }
}

.title-container {
  position: relative;
  padding: 15px;
  text-align: center;
  background: linear-gradient(135deg, 
    rgba(255, 255, 255, 0.95),
    rgba(219, 234, 254, 0.9)
  );
  border-radius: 28px 28px 0 0;
  overflow: hidden;

  &::before {
    content: '';
    position: absolute;
    top: -50%;
    left: -50%;
    width: 200%;
    height: 200%;
    background: linear-gradient(
      to right,
      rgba(255, 255, 255, 0) 0%,
      rgba(255, 255, 255, 0.3) 50%,
      rgba(255, 255, 255, 0) 100%
    );
    transform: rotate(45deg);
    animation: shine 3s infinite;
  }

  &::after {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background: radial-gradient(
      circle at 30% 30%,
      rgba(255, 255, 255, 0.4) 0%,
      transparent 60%
    );
    animation: titleGlow 6s ease-in-out infinite;
  }

  .wn-text {
    position: relative;
    font-size: 24px;
    font-weight: bold;
    background: linear-gradient(135deg, #1976d2, #2196f3);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    text-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    letter-spacing: 2px;
    z-index: 1;
  }
}

@keyframes shine {
  0% {
    left: -50%;
    opacity: 0;
  }
  50% {
    opacity: 1;
  }
  100% {
    left: 150%;
    opacity: 0;
  }
}

.logo {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 4px;
}

.wn-text {
  font-size: 28px;
  font-weight: 700;
  background: linear-gradient(135deg, #1976d2, #2196f3);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  letter-spacing: -1px;
  position: relative;
}

.wn-text::after {
  content: '';
  position: absolute;
  bottom: -4px;
  left: 0;
  width: 100%;
  height: 2px;
  background: linear-gradient(90deg, transparent, #3498db, transparent);
}

.title {
  font-size: 14px;
  color: #1976d2;
  font-weight: 500;
  margin-top: 6px;
  opacity: 0.9;
}

.form-content {
  padding: 35px 35px 30px;
}

.input-wrapper {
  margin-bottom: 0;
  width: 100%;
}

.login-form :deep(.el-form-item) {
  margin-bottom: 25px;
  width: 100%;
}

.login-form :deep(.el-input) {
  width: 100%;
}

.login-form :deep(.el-input__wrapper) {
  background: rgba(255, 255, 255, 0.95);
  border: 1px solid rgba(200, 215, 235, 0.6);
  border-radius: 16px;
  height: 48px;
  padding: 0 16px;
  width: 100%;
  box-sizing: border-box;
  backdrop-filter: blur(10px);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  box-shadow: 
    0 2px 6px rgba(0, 0, 0, 0.03),
    0 1px 3px rgba(0, 0, 0, 0.05),
    inset 0 1px 1px rgba(255, 255, 255, 0.8);
}

.login-form :deep(.el-input__wrapper.is-focus) {
  border-color: #1976d2;
  background: #ffffff;
  box-shadow: 
    0 0 0 1px rgba(25, 118, 210, 0.3),
    0 4px 12px rgba(25, 118, 210, 0.08) !important;
  transform: translateY(-1px);
}

.login-form :deep(.el-input__wrapper:hover) {
  background: #ffffff;
  border-color: rgba(25, 118, 210, 0.5);
  box-shadow: 
    0 2px 8px rgba(0, 0, 0, 0.05),
    0 1px 4px rgba(0, 0, 0, 0.1);
  transform: translateY(-1px);
}

.login-form :deep(.el-input__inner) {
  height: 48px;
  font-size: 15px;
  color: #2c3e50;
  letter-spacing: 0.5px;
  font-weight: 500;
  padding: 0 8px;
}

.login-form :deep(.el-input__inner::placeholder) {
  color: #94a3b8;
  font-weight: 400;
}

.login-form :deep(.el-input__prefix) {
  padding: 0 8px;
}

.login-form :deep(.el-input__prefix-inner) {
  display: flex;
  align-items: center;
}

.login-form :deep(.el-input__prefix .el-icon) {
  font-size: 18px;
  color: #64748b;
  transition: all 0.3s ease;
}

.login-form :deep(.el-input__wrapper.is-focus .el-input__prefix .el-icon) {
  color: #1976d2;
}

.login-form :deep(.el-input__suffix) {
  padding: 0 8px;
}

.login-form :deep(.el-input__suffix .el-icon) {
  font-size: 18px;
  color: #64748b;
  transition: all 0.3s ease;
  cursor: pointer;
}

.login-form :deep(.el-input__suffix .el-icon:hover) {
  color: #1976d2;
}

/* 移除自动填充时的背景色 */
.login-form input:-webkit-autofill,
.login-form input:-webkit-autofill:hover,
.login-form input:-webkit-autofill:focus,
.login-form input:-webkit-autofill:active {
  -webkit-box-shadow: 0 0 0 30px rgba(255, 255, 255, 0.9) inset !important;
  -webkit-text-fill-color: #2c3e50 !important;
  transition: background-color 5000s ease-in-out 0s;
}

.form-content {
  padding: 35px 35px 30px;
}

.form-item {
  margin-bottom: 25px;
}

.input-icon {
  font-size: 18px;
  color: #666;
  transition: all 0.3s ease;
}

.login-form :deep(.el-input__wrapper.is-focus) .input-icon {
  color: #3498db;
}

.login-button {
  width: 100%;
  height: 42px;
  margin-top: 10px;
  border: none;
  border-radius: 12px;
  background: linear-gradient(135deg, #1976d2, #2196f3);
  font-size: 16px;
  font-weight: 600;
  color: #ffffff;
  position: relative;
  overflow: hidden;
  transition: all 0.3s ease;
  box-shadow: 
    0 4px 15px rgba(33, 150, 243, 0.3),
    0 2px 8px rgba(25, 118, 210, 0.2);
}

.login-button:hover {
  transform: translateY(-2px);
  background: linear-gradient(135deg, #1565c0, #1976d2);
  box-shadow: 
    0 8px 25px rgba(33, 150, 243, 0.4),
    0 4px 12px rgba(25, 118, 210, 0.3);
}

.login-button:active {
  transform: translateY(0);
  background: linear-gradient(135deg, #0d47a1, #1565c0);
}

.button-text {
  position: relative;
  z-index: 1;
  letter-spacing: 4px;
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.2);
}

.button-effect {
  position: absolute;
  top: 0;
  left: 0;
  width: 200%;
  height: 100%;
  background: linear-gradient(90deg, 
    transparent,
    rgba(255, 255, 255, 0.3),
    transparent
  );
  transform: translateX(-100%);
  animation: buttonShine 4s infinite;
}

@keyframes buttonShine {
  0% {
    transform: translateX(-100%);
  }
  20%, 100% {
    transform: translateX(100%);
  }
}

@media (max-width: 576px) {
  .login-form {
    width: 320px;
  }
  
  .title-container {
    padding: 16px;
  }
  
  .form-content {
    padding: 0 25px 25px;
  }
  
  .wn-text {
    font-size: 28px;
  }
  
  .title {
    font-size: 14px;
  }
}

@keyframes formAppear {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes formGlow {
  0%, 100% {
    transform: translate(-50%, -50%) scale(1);
    opacity: 0.5;
  }
  50% {
    transform: translate(-50%, -50%) scale(1.1);
    opacity: 0.3;
  }
}

@keyframes titleGlow {
  0%, 100% {
    opacity: 0.4;
    transform: scale(1);
  }
  50% {
    opacity: 0.6;
    transform: scale(1.1);
  }
}
</style> 