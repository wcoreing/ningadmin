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
    
    // 递归生成路由配置
    const generateRoute = (menu) => {
      // 递归查找第一个叶子节点
      const findFirstLeaf = (menu) => {
        if (!menu.children || menu.children.length === 0) {
          return menu
        }
        return findFirstLeaf(menu.children[0])
      }

      // 递归收集所有叶子节点
      const collectLeafNodes = (menu, parentPath = '') => {
        const currentPath = parentPath ? `${parentPath}/${menu.path}` : menu.path
        
        // 为顶级菜单（无论是否有子菜单）创建 Layout 包装
        if (menu.parentId == 0) {
          const route = {
            path: currentPath,
            component: Layout,
            meta: { title: menu.name, icon: menu.icon }
          }
          
          // 如果没有子菜单，创建一个子路由指向组件
          if (!menu.children || menu.children.length === 0) {
            route.redirect = `${currentPath}/index`
            route.children = [{
              path: 'index',
              name: menu.name,
              component: () => import(/* @vite-ignore */ menu.component),
              meta: { title: menu.name, icon: menu.icon }
            }]
            return route
          }
          
          // 如果有子菜单，处理子菜单
          const firstLeaf = findFirstLeaf(menu)
          const redirectPath = `${currentPath}/${firstLeaf.path}`
          route.redirect = redirectPath
          route.children = []
          
          // 递归处理子节点
          menu.children.forEach(child => {
            const childRoute = collectLeafNodes(child, currentPath)
            route.children.push(childRoute)
          })
          
          return route
        }
        
        // 非顶级菜单的处理逻辑
        if (!menu.children || menu.children.length === 0) {
          // 叶子节点
          return {
            path: menu.path,
            name: menu.name,
            component: () => import(/* @vite-ignore */ menu.component),
            meta: { title: menu.name, icon: menu.icon }
          }
        }
        
        // 非顶级且有子菜单的节点
        const route = {
          path: menu.path,
          meta: { title: menu.name, icon: menu.icon },
          children: []
        }
        
        const firstLeaf = findFirstLeaf(menu)
        const redirectPath = `${menu.path}/${firstLeaf.path}`
        route.redirect = redirectPath
        
        // 递归处理子节点
        menu.children.forEach(child => {
          const childRoute = collectLeafNodes(child, menu.path)
          route.children.push(childRoute)
        })
        
        return route
      }

      return collectLeafNodes(menu)
    }

    // 遍历顶层菜单生成路由
    menuList.value.forEach(menu => {
      const route = generateRoute(menu)
      console.log('生成的路由:', route)
      routes.push(route)
    })
    
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