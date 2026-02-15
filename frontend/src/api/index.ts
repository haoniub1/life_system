import axios from 'axios'
import type { ApiResponse } from '@/types'

const instance = axios.create({
  baseURL: '/api',
  withCredentials: true,
  timeout: 10000
})

// Request interceptor: add token from localStorage and log requests
instance.interceptors.request.use(
  (config) => {
    // Get token from localStorage and add to Authorization header
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
      console.log('🔑 Added token to Authorization header')
    }

    console.log('📤 API Request:', config.method?.toUpperCase(), config.url)
    console.log('  Headers:', config.headers)
    console.log('  withCredentials:', config.withCredentials)
    return config
  },
  (error) => {
    console.error('❌ Request Error:', error)
    return Promise.reject(error)
  }
)

// Response interceptor: unwrap CommonResp and handle auth errors
instance.interceptors.response.use(
  (response): any => {
    console.log('📥 API Response:', response.config.url, response.status)
    console.log('  Data:', response.data)
    console.log('  Headers:', response.headers)

    const data = response.data as ApiResponse
    // Backend returns HTTP 200 with code in body
    if (data.code === 401) {
      console.error('❌ Auth Error: 401 Unauthorized')
      // Don't redirect here - let the calling code handle it
      // This prevents infinite reload loops during initAuth()
      return Promise.reject(new Error(data.message || 'unauthorized'))
    }
    if (data.code !== 0) {
      console.error('❌ API Error:', data.code, data.message)
      return Promise.reject(new Error(data.message || 'request failed'))
    }
    // Return the full CommonResp so callers can access .data
    return data
  },
  (error) => {
    console.error('❌ Response Error:', error)
    // Don't redirect here either - let router guards handle auth
    return Promise.reject(error)
  }
)

export default instance
