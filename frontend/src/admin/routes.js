export const asyncRoutes = []

export const adminChildren = [
  {
    path: 'login',
    name: 'Login',
    component: () => import('@/admin/views/login/index.vue'),
    meta: { title: '登录' }
  },
  {
    path: 'form/view/:id',
    name: 'FormView',
    component: () => import('@/admin/views/form/view/index.vue'),
    meta: { title: '表单详情' }
  }
]
