import { createRouter, createWebHistory } from 'vue-router'
import Layout from '@/layout/index.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/login',
      name: 'Login',
      component: () => import('@/views/login/index.vue'),
      meta: { title: '登录' }
    },
    {
      path: '/register',
      name: 'Register',
      component: () => import('@/views/register/index.vue'),
      meta: { title: '注册' }
    },
    {
      path: '/',
      component: Layout,
      redirect: '/dashboard',
      children: [
        {
          path: 'dashboard',
          name: 'Dashboard',
          component: () => import('@/views/dashboard/index.vue'),
          meta: { title: '工作台' }
        },
        {
          path: 'form-generator',
          name: 'FormGenerator',
          component: () => import('@/views/form-generator/index.vue'),
          meta: { title: '表单设计' }
        },
        {
          path: 'form-data',
          name: 'FormData',
          component: () => import('@/views/form-data/index.vue'),
          meta: { title: '数据管理' }
        },
        {
          path: 'system/user',
          name: 'User',
          component: () => import('@/views/system/user/index.vue'),
          meta: { title: '用户管理' }
        },
        {
          path: 'system/role',
          name: 'Role',
          component: () => import('@/views/system/role/index.vue'),
          meta: { title: '角色管理' }
        },
        {
          path: 'system/menu',
          name: 'Menu',
          component: () => import('@/views/system/menu/index.vue'),
          meta: { title: '菜单管理' }
        },
        {
          path: '/monitor',
          component: () => import('@/views/monitor/index.vue'),
          meta: { title: '服务器状态' }
        },
        {
          path: 'system/config',
          name: 'Config',
          component: () => import('@/views/system/config/index.vue'),
          meta: { title: '配置管理' }
        }
      ]
    },
    {
      path: '/form-view/:id',
      component: () => import('@/views/form-view/index.vue'),
      meta: { title: '表单预览' }
    }
  ]
})

// 路由守卫
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  if (to.path !== '/login' && to.path !== '/register' && !token) {
    next('/login')
  } else {
    next()
  }
})

export default router
