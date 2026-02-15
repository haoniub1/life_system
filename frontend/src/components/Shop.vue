<template>
  <div class="shop-container">
    <n-tabs v-model:value="activeTab" type="segment" animated>
      <n-tab-pane name="shop" tab="🛒 商店">
        <div class="shop-content">
          <!-- Spirit Stones Display with RMB toggle -->
          <div class="spirit-display" @click="showRMB = !showRMB">
            <template v-if="showRMB">
              <span class="spirit-icon">💰</span>
              <span class="spirit-amount">¥{{ character?.spiritStones || 0 }}</span>
              <span class="spirit-label-text">RMB</span>
            </template>
            <template v-else>
              <div class="spirit-stones-breakdown">
                <span v-if="spiritDisplay.supreme > 0" class="stone-chip stone-supreme">🔮 {{ spiritDisplay.supreme }} 极品</span>
                <span v-if="spiritDisplay.high > 0" class="stone-chip stone-high">💠 {{ spiritDisplay.high }} 上品</span>
                <span v-if="spiritDisplay.medium > 0" class="stone-chip stone-medium">💎 {{ spiritDisplay.medium }} 中品</span>
                <span class="stone-chip stone-low">🪨 {{ spiritDisplay.low }} 下品</span>
              </div>
            </template>
          </div>

          <!-- Create Button -->
          <div class="action-bar">
            <n-button type="primary" @click="showCreateForm = true">
              + 创建商品
            </n-button>
          </div>

          <!-- Shop Items -->
          <div v-if="loading" class="loading-state">
            <n-spin />
            <p>加载中...</p>
          </div>

          <div v-else-if="shopItems.length === 0" class="empty-state">
            <n-empty description="还没有商品，点击上方按钮创建" />
          </div>

          <n-grid v-else cols="1 s:2 m:3 l:4" :x-gap="12" :y-gap="12" responsive="screen">
            <n-grid-item v-for="item in shopItems" :key="item.id">
              <n-card class="shop-item-card" hoverable>
                <div class="item-visual">
                  <img v-if="item.image" :src="item.image" class="item-image" />
                  <div v-else class="item-icon">{{ item.icon || '🎁' }}</div>
                </div>
                <div class="item-type-badge" :class="item.itemType === 'equipment' ? 'type-equipment' : 'type-consumable'">
                  {{ item.itemType === 'equipment' ? '装备' : '消耗品' }}
                </div>
                <h3 class="item-name">{{ item.name }}</h3>
                <p class="item-description">{{ item.description }}</p>

                <div class="item-footer">
                  <div class="item-price">
                    <span class="price-icon">{{ showRMB ? '💰' : '💎' }}</span>
                    <span class="price-value">{{ showRMB ? '¥' + item.price : item.price }}</span>
                  </div>
                  <n-space :size="4">
                    <n-button
                      type="primary"
                      size="small"
                      :disabled="(character?.spiritStones || 0) < item.price"
                      :loading="purchasing === item.id"
                      @click="handlePurchase(item)"
                    >
                      购买
                    </n-button>
                    <n-button size="small" @click="handleEdit(item)">编辑</n-button>
                    <n-popconfirm @positive-click="handleDelete(item.id)">
                      <template #trigger>
                        <n-button size="small" type="error">删除</n-button>
                      </template>
                      确定删除「{{ item.name }}」吗？
                    </n-popconfirm>
                  </n-space>
                </div>
              </n-card>
            </n-grid-item>
          </n-grid>
        </div>
      </n-tab-pane>

      <n-tab-pane name="inventory" tab="🎒 背包">
        <div class="inventory-content">
          <div v-if="loadingInventory" class="loading-state">
            <n-spin />
            <p>加载中...</p>
          </div>

          <div v-else-if="inventoryItems.length === 0" class="empty-state">
            <n-empty description="背包是空的">
              <template #extra>
                <n-button type="primary" @click="activeTab = 'shop'">
                  去商店看看
                </n-button>
              </template>
            </n-empty>
          </div>

          <n-grid v-else cols="1 s:2 m:3 l:4" :x-gap="12" :y-gap="12" responsive="screen">
            <n-grid-item v-for="item in inventoryItems" :key="item.id">
              <n-card class="inventory-item-card" hoverable>
                <div class="item-visual">
                  <img v-if="item.image" :src="item.image" class="item-image" />
                  <div v-else class="item-icon">{{ item.icon || '🎁' }}</div>
                </div>
                <div class="item-type-badge" :class="item.itemType === 'equipment' ? 'type-equipment' : 'type-consumable'">
                  {{ item.itemType === 'equipment' ? '装备' : '消耗品' }}
                </div>
                <h3 class="item-name">{{ item.name }}</h3>
                <p class="item-description">{{ item.description }}</p>

                <div class="item-meta">
                  <n-tag type="success" size="small">
                    数量: {{ item.quantity }}
                  </n-tag>
                  <n-tag v-if="item.itemType === 'equipment' && item.sellPrice > 0" type="warning" size="small">
                    售价: {{ showRMB ? '¥' + item.sellPrice : item.sellPrice + '灵石' }}
                  </n-tag>
                </div>

                <div class="inventory-actions">
                  <n-popconfirm v-if="item.itemType === 'consumable'" @positive-click="handleUseItem(item)">
                    <template #trigger>
                      <n-button type="info" size="small" block :loading="usingItem === item.itemId">
                        使用
                      </n-button>
                    </template>
                    确定使用「{{ item.name }}」吗？使用后将从背包中移除
                  </n-popconfirm>
                  <n-popconfirm v-if="item.itemType === 'equipment' && item.sellPrice > 0" @positive-click="handleSellItem(item)">
                    <template #trigger>
                      <n-button type="warning" size="small" block :loading="sellingItem === item.itemId">
                        出售 ({{ showRMB ? '¥' + item.sellPrice : item.sellPrice + '灵石' }})
                      </n-button>
                    </template>
                    确定出售「{{ item.name }}」吗？将获得 {{ item.sellPrice }} 灵石
                  </n-popconfirm>
                </div>
              </n-card>
            </n-grid-item>
          </n-grid>
        </div>
      </n-tab-pane>

      <n-tab-pane name="history" tab="📜 历史">
        <div class="history-content">
          <div v-if="loadingHistory" class="loading-state">
            <n-spin />
            <p>加载中...</p>
          </div>

          <div v-else-if="history.length === 0" class="empty-state">
            <n-empty description="暂无购买记录" />
          </div>

          <n-space v-else vertical :size="12">
            <n-card
              v-for="record in history"
              :key="record.id"
              class="history-card"
              size="small"
            >
              <div class="history-item">
                <div class="history-info">
                  <span class="history-name">{{ record.itemName }}</span>
                  <span class="history-quantity">x{{ record.quantity }}</span>
                </div>
                <div class="history-meta">
                  <span class="history-price">{{ showRMB ? '¥' + record.totalPrice : '💎 ' + record.totalPrice }}</span>
                  <span class="history-date">{{ formatDate(record.createdAt) }}</span>
                </div>
              </div>
            </n-card>
          </n-space>
        </div>
      </n-tab-pane>
    </n-tabs>

    <!-- Create / Edit Modal -->
    <n-modal
      v-model:show="showCreateForm"
      preset="card"
      :title="editingItem ? '编辑商品' : '创建商品'"
      size="large"
      class="shop-form-modal"
      @after-leave="resetForm"
    >
      <n-form :model="itemForm" label-placement="left" label-width="100px">
        <n-form-item label="商品名称">
          <n-input v-model:value="itemForm.name" placeholder="输入商品名称" />
        </n-form-item>

        <n-form-item label="描述">
          <n-input v-model:value="itemForm.description" placeholder="输入商品描述" type="textarea" :rows="2" />
        </n-form-item>

        <n-form-item label="商品类型">
          <n-radio-group v-model:value="itemForm.itemType">
            <n-space>
              <n-radio value="consumable">消耗品</n-radio>
              <n-radio value="equipment">装备</n-radio>
            </n-space>
          </n-radio-group>
        </n-form-item>

        <n-form-item label="图标">
          <n-input v-model:value="itemForm.icon" placeholder="输入emoji，如 💊🥤📚" style="width: 120px" />
        </n-form-item>

        <n-form-item label="商品图片">
          <div class="image-upload-area">
            <div v-if="itemForm.image" class="image-preview">
              <img :src="itemForm.image" />
              <n-button size="tiny" class="remove-image" @click="itemForm.image = ''">✕</n-button>
            </div>
            <n-button v-else size="small" @click="triggerImageUpload" :loading="uploadingImage">
              上传图片
            </n-button>
            <input
              ref="imageInput"
              type="file"
              accept="image/*"
              style="display: none"
              @change="handleImageUpload"
            />
          </div>
        </n-form-item>

        <n-form-item label="购买价格">
          <n-input-number v-model:value="itemForm.price" :min="0" placeholder="灵石数" />
          <span style="margin-left: 8px; color: #808090; font-size: 12px;">= ¥{{ itemForm.price }}</span>
        </n-form-item>

        <n-form-item v-if="itemForm.itemType === 'equipment'" label="出售价格">
          <n-input-number v-model:value="itemForm.sellPrice" :min="0" placeholder="灵石数" />
          <span style="margin-left: 8px; color: #808090; font-size: 12px;">= ¥{{ itemForm.sellPrice }}</span>
        </n-form-item>

        <n-form-item label="库存">
          <n-input-number v-model:value="itemForm.stock" :min="-1" placeholder="-1 为无限" />
          <span style="margin-left: 8px; color: #808090; font-size: 12px;">-1 表示无限库存</span>
        </n-form-item>
      </n-form>

      <template #footer>
        <n-space justify="end">
          <n-button @click="showCreateForm = false">取消</n-button>
          <n-button type="primary" :loading="submitting" @click="handleSubmitItem">
            {{ editingItem ? '保存' : '创建' }}
          </n-button>
        </n-space>
      </template>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useMessage } from 'naive-ui'
import {
  NTabs,
  NTabPane,
  NCard,
  NButton,
  NGrid,
  NGridItem,
  NTag,
  NEmpty,
  NSpin,
  NSpace,
  NModal,
  NForm,
  NFormItem,
  NInput,
  NInputNumber,
  NPopconfirm,
  NRadioGroup,
  NRadio
} from 'naive-ui'
import { useCharacterStore } from '@/stores/character'
import { decomposeSpiritStones } from '@/utils/rpg'
import type { ShopItem, InventoryItem, PurchaseRecord } from '@/types'
import { shopApi, uploadFile } from '@/api/shop'

const message = useMessage()
const characterStore = useCharacterStore()

const character = computed(() => characterStore.character)
const spiritDisplay = computed(() => decomposeSpiritStones(character.value?.spiritStones || 0))
const activeTab = ref('shop')
const loading = ref(false)
const loadingInventory = ref(false)
const loadingHistory = ref(false)
const purchasing = ref<number | null>(null)
const usingItem = ref<number | null>(null)
const sellingItem = ref<number | null>(null)
const showRMB = ref(false)

const shopItems = ref<ShopItem[]>([])
const inventoryItems = ref<InventoryItem[]>([])
const history = ref<PurchaseRecord[]>([])

// Form state
const showCreateForm = ref(false)
const editingItem = ref<ShopItem | null>(null)
const submitting = ref(false)
const uploadingImage = ref(false)
const imageInput = ref<HTMLInputElement | null>(null)

const itemForm = ref<{
  name: string
  description: string
  icon: string
  image: string
  itemType: 'consumable' | 'equipment'
  price: number
  sellPrice: number
  stock: number
}>({
  name: '',
  description: '',
  icon: '',
  image: '',
  itemType: 'consumable',
  price: 0,
  sellPrice: 0,
  stock: -1
})

const formatDate = (dateStr: string): string => {
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-CN') + ' ' + date.toLocaleTimeString('zh-CN', {
    hour: '2-digit',
    minute: '2-digit'
  })
}

const resetForm = () => {
  editingItem.value = null
  itemForm.value = {
    name: '',
    description: '',
    icon: '',
    image: '',
    itemType: 'consumable',
    price: 0,
    sellPrice: 0,
    stock: -1
  }
}

const handleEdit = (item: ShopItem) => {
  editingItem.value = item
  itemForm.value = {
    name: item.name,
    description: item.description,
    icon: item.icon,
    image: item.image || '',
    itemType: item.itemType || 'consumable',
    price: item.price,
    sellPrice: item.sellPrice || 0,
    stock: item.stock
  }
  showCreateForm.value = true
}

const triggerImageUpload = () => {
  imageInput.value?.click()
}

const handleImageUpload = async (e: Event) => {
  const target = e.target as HTMLInputElement
  const file = target.files?.[0]
  if (!file) return

  if (file.size > 5 * 1024 * 1024) {
    message.error('图片大小不能超过5MB')
    return
  }

  try {
    uploadingImage.value = true
    const response = await uploadFile(file) as any
    if (response.data?.url) {
      itemForm.value.image = response.data.url
      message.success('图片上传成功')
    }
  } catch (error: any) {
    message.error(error?.message || '上传失败')
  } finally {
    uploadingImage.value = false
    target.value = ''
  }
}

const handleSubmitItem = async () => {
  if (!itemForm.value.name) {
    message.error('请输入商品名称')
    return
  }

  try {
    submitting.value = true

    if (editingItem.value) {
      await shopApi.updateItem(editingItem.value.id, itemForm.value)
      message.success('商品已更新')
    } else {
      await shopApi.createItem(itemForm.value)
      message.success('商品已创建')
    }

    showCreateForm.value = false
    await fetchShopItems()
  } catch (error: any) {
    message.error(error?.message || '操作失败')
  } finally {
    submitting.value = false
  }
}

const handleDelete = async (id: number) => {
  try {
    await shopApi.deleteItem(id)
    message.success('商品已删除')
    await fetchShopItems()
  } catch (error: any) {
    message.error(error?.message || '删除失败')
  }
}

const fetchShopItems = async () => {
  try {
    loading.value = true
    const response = await shopApi.getItems()
    if (response.data?.items) {
      shopItems.value = response.data.items
    } else {
      shopItems.value = []
    }
  } catch (error: any) {
    console.error('Failed to fetch shop items:', error)
    message.error(error?.message || '获取商品列表失败')
  } finally {
    loading.value = false
  }
}

const fetchInventory = async () => {
  try {
    loadingInventory.value = true
    const response = await shopApi.getInventory()
    if (response.data?.items) {
      inventoryItems.value = response.data.items
    }
  } catch (error: any) {
    console.error('Failed to fetch inventory:', error)
    message.error(error?.message || '获取背包失败')
  } finally {
    loadingInventory.value = false
  }
}

const fetchHistory = async () => {
  try {
    loadingHistory.value = true
    const response = await shopApi.getHistory()
    if (response.data?.history) {
      history.value = response.data.history
    }
  } catch (error: any) {
    console.error('Failed to fetch history:', error)
    message.error(error?.message || '获取历史记录失败')
  } finally {
    loadingHistory.value = false
  }
}

const handlePurchase = async (item: ShopItem) => {
  try {
    purchasing.value = item.id
    await shopApi.purchase({ itemId: item.id, quantity: 1 })
    message.success('购买成功')
    await characterStore.fetchCharacter()
    await fetchShopItems()
    await fetchInventory()
    await fetchHistory()
  } catch (error: any) {
    message.error(error?.message || '购买失败')
  } finally {
    purchasing.value = null
  }
}

const handleUseItem = async (item: InventoryItem) => {
  try {
    usingItem.value = item.itemId
    await shopApi.useItem({ itemId: item.itemId, quantity: 1 })
    message.success(`已使用「${item.name}」`)
    await characterStore.fetchCharacter()
    await fetchInventory()
  } catch (error: any) {
    message.error(error?.message || '使用失败')
  } finally {
    usingItem.value = null
  }
}

const handleSellItem = async (item: InventoryItem) => {
  try {
    sellingItem.value = item.itemId
    await shopApi.sellItem({ itemId: item.itemId, quantity: 1 })
    message.success(`已出售「${item.name}」，获得 ${item.sellPrice} 灵石`)
    await characterStore.fetchCharacter()
    await fetchInventory()
  } catch (error: any) {
    message.error(error?.message || '出售失败')
  } finally {
    sellingItem.value = null
  }
}

onMounted(async () => {
  await fetchShopItems()
  await fetchInventory()
  await fetchHistory()
})
</script>

<style scoped>
.shop-container {
  animation: fadeIn 0.3s ease-in-out;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}

:deep(.n-tabs .n-tabs-nav) {
  background: rgba(30, 30, 50, 0.5);
  border-radius: 8px;
  padding: 4px;
  margin-bottom: 24px;
}

.spirit-display {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 16px;
  background: linear-gradient(135deg, rgba(139, 92, 246, 0.1) 0%, rgba(99, 102, 241, 0.05) 100%);
  border: 2px solid rgba(139, 92, 246, 0.3);
  border-radius: 12px;
  margin-bottom: 16px;
  cursor: pointer;
  transition: background 0.2s;
  user-select: none;
}

.spirit-display:hover {
  background: linear-gradient(135deg, rgba(139, 92, 246, 0.15) 0%, rgba(99, 102, 241, 0.1) 100%);
}

.spirit-icon { font-size: 32px; }
.spirit-amount { font-size: 28px; font-weight: bold; color: #a78bfa; }
.spirit-label-text { font-size: 14px; color: #d0d0e0; }

.spirit-stones-breakdown {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
  justify-content: center;
}

.stone-chip {
  font-size: 15px;
  font-weight: 700;
  padding: 4px 10px;
  border-radius: 8px;
}

.stone-supreme {
  color: #ffd700;
  background: rgba(255, 215, 0, 0.15);
}

.stone-high {
  color: #c084fc;
  background: rgba(168, 85, 247, 0.15);
}

.stone-medium {
  color: #60a5fa;
  background: rgba(59, 130, 246, 0.15);
}

.stone-low {
  color: #9ca3af;
  background: rgba(156, 163, 175, 0.1);
}

.action-bar {
  margin-bottom: 16px;
  display: flex;
  justify-content: flex-end;
}

.shop-item-card,
.inventory-item-card,
.history-card {
  background: linear-gradient(135deg, rgba(30, 30, 50, 0.8) 0%, rgba(20, 20, 40, 0.8) 100%);
  border: 1px solid rgba(255, 215, 0, 0.2);
  transition: all 0.3s ease;
}

.shop-item-card:hover,
.inventory-item-card:hover {
  border-color: rgba(255, 215, 0, 0.5);
  transform: translateY(-4px);
  box-shadow: 0 8px 16px rgba(255, 215, 0, 0.2);
}

.item-visual {
  text-align: center;
  margin-bottom: 12px;
}

.item-image {
  width: 64px;
  height: 64px;
  object-fit: cover;
  border-radius: 8px;
  border: 1px solid rgba(255, 215, 0, 0.2);
}

.item-icon { font-size: 48px; }

.item-type-badge {
  text-align: center;
  font-size: 11px;
  font-weight: 600;
  padding: 2px 8px;
  border-radius: 10px;
  display: inline-block;
  margin-bottom: 8px;
}

.type-consumable {
  background: rgba(59, 130, 246, 0.15);
  color: #60a5fa;
  border: 1px solid rgba(59, 130, 246, 0.3);
}

.type-equipment {
  background: rgba(245, 158, 11, 0.15);
  color: #fbbf24;
  border: 1px solid rgba(245, 158, 11, 0.3);
}

.item-name { font-size: 18px; font-weight: bold; color: #ffd700; margin: 0 0 8px 0; text-align: center; }
.item-description { font-size: 14px; color: #a0a0b0; margin: 0 0 12px 0; text-align: center; min-height: 40px; }

.item-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
  gap: 8px;
}

.inventory-actions {
  padding-top: 8px;
  border-top: 1px solid rgba(255, 215, 0, 0.15);
}

.item-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 12px;
  border-top: 1px solid rgba(255, 215, 0, 0.2);
  flex-wrap: wrap;
  gap: 8px;
}

.item-price { display: flex; align-items: center; gap: 4px; font-size: 18px; font-weight: bold; color: #a78bfa; }

.loading-state,
.empty-state {
  text-align: center;
  padding: 60px 20px;
  color: #a0a0b0;
}

.history-item { display: flex; justify-content: space-between; align-items: center; }
.history-info { display: flex; align-items: center; gap: 12px; }
.history-name { font-weight: 600; color: #d0d0e0; }
.history-quantity { color: #a0a0b0; font-size: 14px; }
.history-meta { display: flex; flex-direction: column; align-items: flex-end; gap: 4px; }
.history-price { color: #a78bfa; font-weight: 600; }
.history-date { font-size: 12px; color: #808080; }

/* Form modal styles */
:deep(.n-modal) {
  background: linear-gradient(135deg, rgba(30, 30, 50, 0.95) 0%, rgba(20, 20, 40, 0.95) 100%);
  border: 1px solid rgba(255, 215, 0, 0.2);
}

:deep(.n-radio) {
  color: #d0d0e0 !important;
}

:deep(.n-radio__label) {
  color: #d0d0e0 !important;
}

.image-upload-area {
  display: flex;
  align-items: center;
  gap: 12px;
}

.image-preview {
  position: relative;
  display: inline-block;
}

.image-preview img {
  width: 64px;
  height: 64px;
  object-fit: cover;
  border-radius: 8px;
  border: 1px solid rgba(255, 215, 0, 0.3);
}

.remove-image {
  position: absolute;
  top: -6px;
  right: -6px;
  border-radius: 50% !important;
  min-width: 20px !important;
  padding: 0 !important;
}

:deep(.n-button--primary) {
  background: linear-gradient(135deg, #ffd700, #ffed4e) !important;
  color: #000 !important;
  border: none !important;
}

@media (max-width: 768px) {
  .item-description { min-height: auto; }
  .item-footer { flex-direction: column; align-items: stretch; }
  .history-item { flex-direction: column; align-items: flex-start; gap: 8px; }
  .history-meta { align-items: flex-start; }
  .action-bar { justify-content: stretch; }
  .action-bar .n-button { width: 100%; }
  .spirit-amount { font-size: 20px; }
  .item-name { font-size: 15px; }
}
</style>
