import request from './index'
import type { CharacterStats, ApiResponse } from '@/types'

export async function getCharacter(): Promise<ApiResponse<CharacterStats>> {
  return request.get('/character')
}

export async function updateCharacter(data: Partial<CharacterStats>): Promise<ApiResponse<CharacterStats>> {
  return request.put('/character', data)
}
