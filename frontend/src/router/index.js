import { createRouter, createWebHistory } from 'vue-router'
import { adminChildren } from '@/admin/routes'
import { firstMenuPath, useUserStore } from '@/admin/stores/user'
import { ensureModuleUi } from '@/plugins/ui'
import { adminPath } from '@/router/base'

const siteRoutes = [
  { path: '', name: 'home', component: () => import('@/site/views/Home.vue'), meta: { title: '集亿玖泽' } },
  { path: 'products', name: 'products', component: () => import('@/site/views/Products.vue'), meta: { title: '产品' } },
  { path: 'products/:id', name: 'product', component: () => import('@/site/views/ProductDetail.vue'), meta: { title: '产品详情' } },
  { path: 'page/:code', name: 'page', component: () => import('@/site/views/Page.vue'), meta: { title: '栏目' } },
  { path: 'news', name: 'news', component: () => import('@/site/views/News.vue'), meta: { title: '资讯' } },
  { path: 'news/:id', name: 'article', component: () => import('@/site/views/NewsDetail.vue'), meta: { title: '资讯' } },
  { path: 'inquiry', name: 'inquiry', component: () => import('@/site/views/Inquiry.vue'), meta: { title: '询盘' } },
  { path: 'privacy', redirect: '/page/privacy' }
]

const cmsRoutes = [
  { path: '', name: 'cms-home', component: () => import('@/cms/views/Home.vue'), meta: { title: 'WnAdmin - 现代化全栈管理系统' } },
  { path: 'features', name: 'cms-features', component: () => import('@/cms/views/Features.vue'), meta: { title: '功能特性 - WnAdmin' } },
  { path: 'technology', name: 'cms-technology', component: () => import('@/cms/views/Technology.vue'), meta: { title: '技术架构 - WnAdmin' } },
  { path: 'demo', name: 'cms-demo', component: () => import('@/cms/views/Demo.vue'), meta: { title: '在线演示 - WnAdmin' } },
  { path: 'docs', name: 'cms-docs', component: () => import('@/cms/views/Docs.vue'), meta: { title: '开发文档 - WnAdmin' } },
  { path: 'contact', name: 'cms-contact', component: () => import('@/cms/views/Contact.vue'), meta: { title: '联系我们 - WnAdmin' } }
]

const router = createRouter({
  history: createWebHistory(),
  scrollBehavior() {
    return { top: 0 }
  },
  routes: [
    {
      path: '/',
      component: () => import('@/site/App.vue'),
      children: siteRoutes
    },
    {
      path: '/admin',
      name: 'admin',
      component: () => import('@/admin/pass.vue'),
      children: adminChildren
    },
    {
      path: '/cms',
      component: () => import('@/cms/pass.vue'),
      children: cmsRoutes
    }
  ]
})

function inGroup(path, base) {
  return path === base || path.startsWith(`${base}/`)
}

router.beforeEach(async (to, from, next) => {
  if (to.meta?.title) document.title = to.meta.title
  if (inGroup(to.path, '/admin') || inGroup(to.path, '/cms')) {
    await ensureModuleUi()
  }
  if (!inGroup(to.path, '/admin')) {
    next()
    return
  }

  const userStore = useUserStore()
  const hasToken = userStore.getToken()
  if (!hasToken) {
    if (to.path === '/admin/login') {
      next()
      return
    }
    next(`/admin/login?redirect=${encodeURIComponent(to.path)}`)
    return
  }
  if (to.path === '/admin/login') {
    next({ path: adminPath('/dashboard') })
    return
  }

  const hasUserInfo = userStore.userInfo && Object.keys(userStore.userInfo).length > 0
  if (!hasUserInfo) {
    try {
      const info = await userStore.getUserInfo()
      if (!info) {
        next(`/admin/login?redirect=${encodeURIComponent(to.path)}`)
        return
      }
      const accessRoutes = await userStore.generateRoutes()
      accessRoutes.forEach(route => router.addRoute('admin', route))
    } catch (error) {
      console.error('路由守卫错误:', error)
      await userStore.resetState()
      next(`/admin/login?redirect=${encodeURIComponent(to.path)}`)
      return
    }
  }
  if (to.path === '/admin') {
    next({ path: firstMenuPath(userStore.menuList), replace: true })
    return
  }
  if (!hasUserInfo) {
    next({ ...to, replace: true })
    return
  }
  next()
})

export default router
