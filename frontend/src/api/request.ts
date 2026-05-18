import axios from 'axios'

export interface ApiResponse<T> {
  code: number
  message: string
  data: T
  traceId?: string
}

export const request = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api/v1',
  timeout: 15000,
})

request.interceptors.request.use((config) => {
  const token = localStorage.getItem('access_token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

request.interceptors.response.use((response) => {
  const body = response.data as ApiResponse<unknown>
  if (body.code === 0) {
    return body.data
  }
  return Promise.reject(body)
})
