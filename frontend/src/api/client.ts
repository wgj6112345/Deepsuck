import axios from 'axios'
import type { AxiosInstance, AxiosError } from 'axios'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'

const client: AxiosInstance = axios.create({
  baseURL: API_BASE_URL,
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// 请求拦截器
client.interceptors.request.use(
  (config) => {
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器
client.interceptors.response.use(
  (response) => {
    return response
  },
  (error: AxiosError) => {
    if (error.response) {
      // 服务器返回错误
      const status = error.response.status
      switch (status) {
        case 402:
          console.error('Token 欠费，请检查配置')
          break
        case 500:
          console.error('服务器内部错误')
          break
        default:
          console.error(`请求失败: ${status}`)
      }
    } else if (error.request) {
      // 请求已发送但无响应
      console.error('网络异常，请检查网络连接')
    } else {
      // 请求配置错误
      console.error('请求配置错误:', error.message)
    }
    return Promise.reject(error)
  }
)

export default client