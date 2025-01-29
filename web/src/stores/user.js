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
    
    // 添加调试日志
    console.log('开始生成动态路由, 当前菜单列表:', menuList.value)
    
    // 递归生成路由配置
    const generateRoute = (menu) => {
      if (menu.children && menu.children.length > 0) {
        // 有子菜单的情况
        const route = {
          path: menu.path,
          component: Layout,
          redirect: `${menu.path}/${menu.children[0].path}`,
          meta: { title: menu.name, icon: menu.icon },
          children: []
        }
        
        // 递归处理子菜单
        menu.children.forEach(child => {
          if (child.children && child.children.length > 0) {
            // 递归处理多级子菜单
            route.children.push(generateRoute(child))
          } else {
            route.children.push({
              path: child.path,
              name: child.name,
              component: () => {
                return import(/* @vite-ignore */ child.component)
              },
              meta: { 
                title: child.name,
                icon: child.icon
              }
            })
          }
        })
        
        return route
      } else {
        // 没有子菜单的情况
        return {
          path: menu.path,
          component: Layout,
          children: [
            {
              path: '',
              name: menu.name,
              component: () => {
                return import(/* @vite-ignore */ menu.component)
              },
              meta: { 
                title: menu.name,
                icon: menu.icon
              }
            }
          ]
        }
      }
    }

    // 遍历顶层菜单生成路由
    menuList.value.forEach(menu => {
      routes.push(generateRoute(menu))
    })
    
    // 添加调试日志
    console.log('生成的路由配置:', [...asyncRoutes, ...routes])
    
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