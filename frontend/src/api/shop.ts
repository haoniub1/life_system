import request from './index'
import type { ApiResponse, ShopItem, InventoryItem, PurchaseRecord } from '@/types'

export const shopApi = {
  // Get shop items
  getItems(): Promise<ApiResponse<{ items: ShopItem[] }>> {
    return request.get('/shop/items')
  },

  // Create shop item
  createItem(data: Partial<ShopItem>): Promise<ApiResponse<ShopItem>> {
    return request.post('/shop/items', data)
  },

  // Update shop item
  updateItem(id: number, data: Partial<ShopItem>): Promise<ApiResponse<ShopItem>> {
    return request.put(`/shop/items/${id}`, data)
  },

  // Delete shop item
  deleteItem(id: number): Promise<ApiResponse> {
    return request.delete(`/shop/items/${id}`)
  },

  // Purchase item
  purchase(data: { itemId: number; quantity: number }): Promise<ApiResponse<any>> {
    return request.post('/shop/purchase', data)
  },

  // Get inventory
  getInventory(): Promise<ApiResponse<{ items: InventoryItem[] }>> {
    return request.get('/shop/inventory')
  },

  // Use item (consumable)
  useItem(data: { itemId: number; quantity: number }): Promise<ApiResponse<any>> {
    return request.post('/shop/use', data)
  },

  // Sell item (equipment)
  sellItem(data: { itemId: number; quantity: number }): Promise<ApiResponse<any>> {
    return request.post('/shop/sell', data)
  },

  // Get purchase history
  getHistory(): Promise<ApiResponse<{ history: PurchaseRecord[] }>> {
    return request.get('/shop/history')
  }
}

// Upload file
export function uploadFile(file: File): Promise<ApiResponse<{ url: string }>> {
  const formData = new FormData()
  formData.append('file', file)
  return request.post('/upload', formData, {
    headers: { 'Content-Type': 'multipart/form-data' }
  })
}
