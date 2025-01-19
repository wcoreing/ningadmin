import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/stores/user'
import Layout from '@/layout/index.vue'

// 公共路由
const publicRoutes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/login/index.vue'),
    meta: { title: '登录' }
  }
]

// 动态路由
export const asyncRoutes = [
  {
    path: '/',
    component: Layout,
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/dashboard/index.vue'),
        meta: { title: '首页' }
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes: publicRoutes
})

// 路由守卫
router.beforeEach(async (to, from, next) => {
  const userStore = useUserStore()
  const hasToken = userStore.getToken()

  if (hasToken) {
    if (to.path === '/login') {
      next({ path: '/' })
    } else {
      // 判断是否已获取用户信息
      const hasUserInfo = userStore.userInfo && Object.keys(userStore.userInfo).length > 0
      
      if (hasUserInfo) {
        next()
      } else {
        try {
          // 获取用户信息
          await userStore.getUserInfo()
          // 生成动态路由
          const accessRoutes = await userStore.generateRoutes()
          // 添加动态路由
          accessRoutes.forEach(route => {
            router.addRoute(route)
          })
          // 继续路由
          next({ ...to, replace: true })
        } catch (error) {
          // 重置token并跳转登录页
          await userStore.resetState()
          next(`/login?redirect=${to.path}`)
        }
      }
    }
  } else {
    if (to.path === '/login') {
      next()
    } else {
      next(`/login?redirect=${to.path}`)
    }
  }
})

export default router
