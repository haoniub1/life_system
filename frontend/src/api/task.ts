import request from './index'
import type { Task, ApiResponse, CompleteTaskResult } from '@/types'

export interface GetTasksParams {
  type?: 'once' | 'repeatable' | 'challenge'
  status?: 'active' | 'completed' | 'failed' | 'deleted'
}

export async function getTasks(params?: GetTasksParams): Promise<ApiResponse<Task[]>> {
  return request.get('/tasks', { params })
}

export async function createTask(data: Partial<Task>): Promise<ApiResponse<Task>> {
  return request.post('/tasks', data)
}

export async function updateTask(id: number, data: Partial<Task>): Promise<ApiResponse<Task>> {
  return request.put(`/tasks/${id}`, data)
}

export async function completeTask(id: number): Promise<ApiResponse<CompleteTaskResult>> {
  return request.post(`/tasks/complete/${id}`)
}

export async function deleteTask(id: number): Promise<ApiResponse> {
  return request.delete(`/tasks/${id}`)
}
