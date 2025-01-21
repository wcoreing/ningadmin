import axios from 'axios'
import { ElMessage } from 'element-plus'

// 创建 axios 实例
const request = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api', // 从环境变量获取基础URL
  timeout: 15000, // 请求超时时间
  headers: {
    'Content-Type': 'application/json'
  }
})

// 请求拦截器
request.interceptors.request.use(
  config => {
    // 从 localStorage 获取 token
    const token = localStorage.getItem('token')
    if (token) {
      config.headers['Authorization'] = `Bearer ${token}`
    }
    return config
  },
  error => {
    console.error('请求错误：', error)
    return Promise.reject(error)
  }
)

// 响应拦截器
request.interceptors.response.use(
  response => {
    const res = response.data
    
    // 这里可以根据后端的响应结构进行调整
    if (res.code === 200) {
      return res.data
    }
    
    // 如果是无效令牌,跳转到登录页
    if (res.code === 401) {
      ElMessage.error('登录已过期,请重新登录')
      localStorage.removeItem('token')
      window.location.href = '/login'
      return Promise.reject(new Error('登录已过期'))
    }
    
    ElMessage.error(res.msg || '请求失败')
    return Promise.reject(new Error(res.msg || '请求失败'))
  },
  error => {
    console.error('响应错误：', error)
    
    // 如果是401错误,跳转到登录页
    if (error.response?.code === 401) {
      ElMessage.error('登录已过期,请重新登录')
      localStorage.removeItem('token')
      window.location.href = '/login'
      return Promise.reject(error)
    }
    
    const message = error.response?.data?.msg || '服务器错误'
    ElMessage.error(message)
    return Promise.reject(error)
  }
)

export default request 