import { createRouter, createWebHistory } from 'vue-router'
import Home from '@/views/Home.vue'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: Home,
    meta: { title: 'WnAdmin - 现代化全栈管理系统' }
  },
  {
    path: '/features',
    name: 'Features',
    component: () => import('@/views/Features.vue'),
    meta: { title: '功能特性 - WnAdmin' }
  },
  {
    path: '/technology',
    name: 'Technology',
    component: () => import('@/views/Technology.vue'),
    meta: { title: '技术架构 - WnAdmin' }
  },
  {
    path: '/demo',
    name: 'Demo',
    component: () => import('@/views/Demo.vue'),
    meta: { title: '在线演示 - WnAdmin' }
  },
  {
    path: '/docs',
    name: 'Docs',
    component: () => import('@/views/Docs.vue'),
    meta: { title: '开发文档 - WnAdmin' }
  },
  {
    path: '/contact',
    name: 'Contact',
    component: () => import('@/views/Contact.vue'),
    meta: { title: '联系我们 - WnAdmin' }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior(to, from, savedPosition) {
    if (savedPosition) {
      return savedPosition
    } else {
      return { top: 0 }
    }
  }
})

router.beforeEach((to, from, next) => {
  document.title = to.meta.title || 'WnAdmin'
  next()
})

export default router 