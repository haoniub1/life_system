import request from './index'
import type { ApiResponse, User } from '@/types'

export function updateProfile(data: { displayName?: string; avatar?: string }): Promise<ApiResponse<User>> {
  return request.put('/user/profile', data)
}

export function changePassword(data: { oldPassword: string; newPassword: string }): Promise<ApiResponse> {
  return request.put('/user/password', data)
}
