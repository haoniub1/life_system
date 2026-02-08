import request from './index'
import type { BindCodeResponse, TgStatus, ApiResponse } from '@/types'

export async function getBindCode(): Promise<ApiResponse<BindCodeResponse>> {
  return request.post('/telegram/bindcode')
}

export async function getStatus(): Promise<ApiResponse<TgStatus>> {
  return request.get('/telegram/status')
}

export async function unbind(): Promise<ApiResponse> {
  return request.delete('/telegram/unbind')
}
