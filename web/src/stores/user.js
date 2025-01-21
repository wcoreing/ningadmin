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
      const res = await getInfo()
      if (!res) {
        throw new Error('获取用户信息失败')
      }
      userInfo.value = res
      menuList.value = res.menu || []
      return res
    } catch (error) {
      console.error('获取用户信息失败:', error)
      resetState()
      return null
    }
  }

  // 生成动态路由
  const generateRoutes = async () => {
    const routes = []
    
    // 生成路由配置
    menuList.value.forEach(menu => {
      if (menu.children && menu.children.length > 0) {
        // 有子菜单的情况
        const route = {
          path: menu.path,
          component: Layout,
          redirect: `${menu.path}/${menu.children[0].path}`,
          meta: { title: menu.name, icon: menu.icon },
          children: []
        }
        
        menu.children.forEach(child => {
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
        })
        
        routes.push(route)
      } else {
        // 没有子菜单的情况
        routes.push({
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
        })
      }
    })
    
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