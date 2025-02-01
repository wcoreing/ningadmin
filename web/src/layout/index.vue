<template>
  <div class="layout">
    <el-container>
      <!-- 侧边栏 -->
      <el-aside :width="isCollapse ? '64px' : '200px'" class="aside">
        <div class="logo">
          <img :src="logoUrl" alt="logo">
          <span v-show="!isCollapse">WnAdmin</span>
        </div>
        <el-menu
          :default-active="route.path"
          :collapse="isCollapse"
          :unique-opened="true"
          router
          class="menu"
          background-color="#001529"
          text-color="#fff"
          active-text-color="#409EFF"
        >
          <!-- 遍历菜单数据生成菜单项 -->
          <template v-for="menu in userStore.menuList" :key="menu.id">
            <!-- 有子菜单 -->
            <el-sub-menu v-if="menu.children && menu.children.length" :index="String(menu.id)">
              <template #title>
                <el-icon><component :is="menu.icon" /></el-icon>
                <span>{{ menu.name }}</span>
              </template>
              <el-menu-item 
                v-for="child in menu.children" 
                :key="child.id"
                :index="child.path"
              >
                {{ child.name }}
              </el-menu-item>
            </el-sub-menu>
            
            <!-- 无子菜单 -->
            <el-menu-item v-else :index="menu.path">
              <el-icon><component :is="menu.icon" /></el-icon>
              <template #title>{{ menu.name }}</template>
            </el-menu-item>
          </template>
        </el-menu>
      </el-aside>

      <el-container>
        <!-- 顶栏 -->
        <el-header class="header">
          <div class="header-left">
            <el-icon class="collapse-btn" @click="toggleCollapse">
              <Fold v-if="!isCollapse" />
              <Expand v-else />
            </el-icon>
            <el-breadcrumb>
              <el-breadcrumb-item :to="{ path: '/dashboard' }">首页</el-breadcrumb-item>
              <el-breadcrumb-item>{{ route.meta.title }}</el-breadcrumb-item>
            </el-breadcrumb>
          </div>
          <div class="header-right">
            <el-dropdown @command="handleCommand">
              <span class="user-info">
                <el-avatar :size="30" src="https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png" />
                <span class="username">admin</span>
              </span>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="profile">个人信息</el-dropdown-item>
                  <el-dropdown-item command="logout">退出登录</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>
        </el-header>

        <!-- 主要内容区 -->
        <el-main class="main">
          <router-view></router-view>
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { Monitor, Setting, Fold, Expand, Operation } from '@element-plus/icons-vue'
import logoUrl from '@/assets/logo.jpg'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const isCollapse = ref(false)

const toggleCollapse = () => {
  isCollapse.value = !isCollapse.value
}

const handleCommand = async (command) => {
  if (command === 'logout') {
    await userStore.logoutFn()
    router.push('/login')
  } else if (command === 'profile') {
    // 跳转到个人信息页
    router.push('/profile')
  }
}
</script>

<style scoped>
.layout {
  height: 100vh;
  background-color: #f0f2f5;
}

.el-container {
  height: 100%;
}

.aside {
  background-color: #001529;
  transition: width 0.3s;
  overflow: hidden;
  border-right: 1px solid #e6e6e6;
  height: 100%;
}

.logo {
  height: 64px;
  padding: 0 16px;
  display: flex;
  align-items: center;
  background-color: #001529;
  overflow: hidden;
  border-bottom: 1px solid #1a1a1a;
}

.logo img {
  height: 32px;
  width: 32px;
  object-fit: cover;
  border-radius: 4px;
}

.logo span {
  color: #fff;
  font-size: 18px;
  font-weight: 600;
  margin-left: 12px;
  white-space: nowrap;
}

.menu {
  border-right: none;
  background-color: #001529;
}

:deep(.el-menu) {
  border-right: none;
}

:deep(.el-menu--collapse) {
  width: 64px;
}

:deep(.el-menu-item),
:deep(.el-sub-menu__title) {
  color: #fff;
  height: 50px;
  line-height: 50px;
}

:deep(.el-menu-item.is-active) {
  color: #409EFF !important;
  background-color: #0c2135 !important;
  border-right: 3px solid #409EFF;
}

:deep(.el-menu-item:hover),
:deep(.el-sub-menu__title:hover) {
  color: #409EFF !important;
  background-color: #0c2135 !important;
}

:deep(.el-sub-menu.is-active .el-sub-menu__title) {
  color: #409EFF !important;
}

:deep(.el-menu-item .el-icon),
:deep(.el-sub-menu .el-icon) {
  color: inherit;
  font-size: 18px;
  margin-right: 10px;
}

:deep(.el-sub-menu__title) {
  &:hover {
    background-color: #0c2135 !important;
  }
}

:deep(.el-menu--popup) {
  background-color: #000c17 !important;
  
  .el-menu-item {
    background-color: #000c17;
    
    &:hover {
      background-color: #0c2135 !important;
    }
    
    &.is-active {
      background-color: #0c2135 !important;
    }
  }
}

.header {
  background-color: #fff;
  border-bottom: 1px solid #e6e6e6;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
}

.header-left {
  display: flex;
  align-items: center;
}

.collapse-btn {
  font-size: 20px;
  margin-right: 20px;
  cursor: pointer;
  color: #333;
}

.header-right {
  display: flex;
  align-items: center;
}

.user-info {
  display: flex;
  align-items: center;
  cursor: pointer;
}

.username {
  margin-left: 8px;
  font-size: 14px;
  color: #333;
}

.main {
  background-color: #f0f2f5;
  padding: 20px;
  height: 100%;
  overflow-y: auto;
}

:deep(.el-container.is-vertical) {
  height: 100%;
}
</style> 