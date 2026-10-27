import axios from 'axios'

const request = axios.create({
  baseURL: '/api',
  timeout: 15000
})

request.interceptors.response.use(
  response => {
    const res = response.data
    if (res.code === 200) {
      return res.data
    }
    return Promise.reject(new Error(res.msg || '请求失败'))
  },
  error => Promise.reject(error)
)

export default request
