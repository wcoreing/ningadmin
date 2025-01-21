export const defaultMenus = [
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: '@/views/dashboard/index.vue',
    meta: { title: '工作台', icon: 'House' }
  },
  {
    meta: { title: '表单管理', icon: 'Tickets' },
    children: [
      {
        path: '/form-generator',
        name: 'FormGenerator',
        component: '@/views/form-generator/index.vue',
        meta: { title: '表单设计' }
      },
      {
        path: '/form-data',
        name: 'FormData',
        component: '@/views/form-data/index.vue',
        meta: { title: '数据管理' }
      }
    ]
  },
  // ... 其他菜单
] 