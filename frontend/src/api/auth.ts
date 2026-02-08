import request from './index'
import type { User, ApiResponse } from '@/types'

export async function login(username: string, password: string): Promise<ApiResponse<User>> {
  return request.post('/auth/login', {
    username,
    password
  })
}

export async function register(username: string, password: string): Promise<ApiResponse<User>> {
  return request.post('/auth/register', {
    username,
    password
  })
}

export async function logout(): Promise<ApiResponse> {
  return request.post('/auth/logout')
}

export async function getMe(): Promise<ApiResponse<User>> {
  return request.get('/auth/me')
}
