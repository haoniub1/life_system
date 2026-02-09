import request from './index'
import type { BarkStatus, ApiResponse } from '@/types'

export interface SetBarkKeyRequest {
  barkKey: string
}

export interface TestBarkRequest {
  title?: string
  body?: string
}

export async function setBarkKey(data: SetBarkKeyRequest): Promise<ApiResponse> {
  return request.put('/bark/key', data)
}

export async function getBarkStatus(): Promise<ApiResponse<BarkStatus>> {
  return request.get('/bark/status')
}

export async function testBark(data?: TestBarkRequest): Promise<ApiResponse> {
  return request.post('/bark/test', data || {})
}

export async function deleteBarkKey(): Promise<ApiResponse> {
  return request.delete('/bark/key')
}
