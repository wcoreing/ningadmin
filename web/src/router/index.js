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
  },
  {
    path: '/form/view/:id',
    name: 'FormView',
    component: () => import('@/views/form/view/index.vue'),
    meta: { title: '表单详情' }
  }
]

// 动态路由
export const asyncRoutes = [
  {
    path: '/',
    component: Layout,
    redirect: '/dashboard',
    meta: { title: 'wnadmin', icon: 'dashboard' },
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/dashboard/index.vue'),
        meta: { title: '首页' }
      }
    ]
  },

  {
    path: '/form',
    component: Layout,
    meta: { title: '表单管理', icon: 'form' },
    children: [
      {
        path: 'list',
        name: 'FormList',
        component: () => import('@/views/form/form/index.vue'),
        meta: { title: '表单列表' }
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
  console.log('路由守卫 - 当前token:', hasToken)
  console.log('路由守卫 - 目标路由:', to.path)

  if (hasToken) {
    if (to.path === '/login') {
      next({ path: '/dashboard' })
    } else {
      const hasUserInfo = userStore.userInfo && Object.keys(userStore.userInfo).length > 0
      console.log("hasUserInfo",hasUserInfo)
      if (hasUserInfo) {
        next()
      } else {
        try {
          await userStore.getUserInfo()
          const accessRoutes = await userStore.generateRoutes()
          
          // 修改这里的路由注册逻辑
          accessRoutes.forEach(route => {
            router.addRoute(route)
          })
          
          // 这里改用 replace: true 确保导航不会留下历史记录
          next({ ...to, replace: true })
        } catch (error) {
          console.error('路由守卫错误:', error)
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
