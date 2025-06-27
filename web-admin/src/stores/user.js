import { defineStore } from 'pinia'
import { ref } from 'vue'
import { login, getUserInfo as getInfo, logout } from '@/api/auth'
import { asyncRoutes } from '@/router'
import Layout from '@/layout/index.vue'

export const useUserStore = defineStore('user', () => {
  const token = ref('')
  const userInfo = ref(null)
  const menuList = ref([])

  // 设置token
  const setToken = (value) => {
    token.value = value
    localStorage.setItem('token', value)
  }

  // 获取token
  const getToken = () => {
    if (!token.value) {
      token.value = localStorage.getItem('token') || ''
    }
    return token.value
  }

  // 登录
  const loginFn = async (loginForm) => {
    try {
      const res = await login(loginForm)
      if (res?.token) {
        setToken(res.token)
        return true
      }
      return false
    } catch (error) {
      console.error('登录失败:', error)
      return false
    }
  }

  // 获取用户信息
  const getUserInfo = async () => {
    try {
      console.log('开始获取用户信息，当前token:', token.value)
      const res = await getInfo()
      console.log('获取用户信息响应:', res)
      if (!res) {
        throw new Error('获取用户信息失败')
      }
      userInfo.value = res
      menuList.value = res.menu || []
      return res
    } catch (error) {
      console.error('获取用户信息失败，详细错误:', error)
      resetState()
      return null
    }
  }
  // 生成动态路由
  const generateRoutes = async () => {
    const routes = []
    
    console.log('开始生成动态路由, 当前菜单列表:', menuList.value)
    
    // 递归构建路由树
    const buildRouteTree = (menus, parentPath = '') => {
      const routeMap = new Map() // 用于避免重复路由
      
      menus.forEach(menu => {
        // 构建完整的路由路径
        const fullPath = menu.path
        
        // 如果菜单有组件，创建路由
        if (menu.component && menu.component.trim() !== '') {
          // 生成唯一的路由名称，避免冲突
          const routeName = generateUniqueRouteName(menu, fullPath)
          
          const route = {
            path: fullPath,
            component: Layout,
            children: [{
              path: '', // 使用空字符串作为默认子路由
              name: routeName,
              component: () => import(/* @vite-ignore */ menu.component),
              meta: { 
                title: menu.name, 
                icon: menu.icon,
                menuId: menu.id,
                parentId: menu.parentId
              }
            }],
            meta: { 
              title: menu.name, 
              icon: menu.icon,
              menuId: menu.id,
              parentId: menu.parentId
            }
          }
          
          // 避免重复路径
          if (!routeMap.has(fullPath)) {
            routeMap.set(fullPath, route)
            console.log('生成的路由:', route)
          } else {
            console.warn('路由路径冲突，跳过:', fullPath, menu)
          }
        }
        
        // 递归处理子菜单
        if (menu.children && menu.children.length > 0) {
          const childRoutes = buildRouteTree(menu.children, fullPath)
          childRoutes.forEach(childRoute => {
            if (!routeMap.has(childRoute.path)) {
              routeMap.set(childRoute.path, childRoute)
            }
          })
        }
      })
      
      return Array.from(routeMap.values())
    }
    
    // 生成唯一路由名称的辅助函数
    const generateUniqueRouteName = (menu, path) => {
      // 使用路径生成唯一名称，避免同名冲突
      const pathSegments = path.split('/').filter(segment => segment !== '')
      const uniqueName = pathSegments.length > 0 
        ? pathSegments.join('_') + '_' + menu.id
        : menu.name + '_' + menu.id
      
      return uniqueName
    }
    
    // 构建路由树
    const generatedRoutes = buildRouteTree(menuList.value)
    routes.push(...generatedRoutes)
    
    console.log('最终的路由配置:', [...asyncRoutes, ...routes])
    
    return [...asyncRoutes, ...routes]
  }

  // 退出登录
  const logoutFn = async () => {
    try {
      await logout()
      resetState()
      return true
    } catch (error) {
      return false
    }
  }

  // 重置状态
  const resetState = () => {
    token.value = ''
    userInfo.value = null
    menuList.value = []
    localStorage.removeItem('token')
  }

  return {
    token,
    userInfo,
    menuList,
    loginFn,
    getUserInfo,
    logoutFn,
    resetState,
    getToken,
    setToken,
    generateRoutes
  }
}) 