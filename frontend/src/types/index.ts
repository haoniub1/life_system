export interface User {
  id: number
  username: string
  displayName: string
  avatar: string
  tgChatId: number
  tgUsername: string
}

export interface CharacterStats {
  userId: number
  spiritStones: number
  fatigue: number
  fatigueCap: number
  fatigueLevel: number
  overdraftPenalty: number
  title: string
  lastActivityDate: string
  attributes: CharacterAttribute[]
}

export interface CharacterAttribute {
  attrKey: string
  displayName: string
  emoji: string
  value: number
  todayGain: number
  realm: number
  realmName: string
  subRealm: number
  subRealmName: string
  realmExp: number
  isBottleneck: boolean
  accumulationPool: number
  attrCap: number
  progressPercent: number
  color: string
}

export interface SpiritStoneDisplay {
  total: number
  supreme: number  // 极品
  high: number     // 上品
  medium: number   // 中品
  low: number      // 下品
}

export interface SleepRecord {
  id: number
  userId: number
  sleepStart: string
  sleepEnd: string
  durationHours: number
  quality: 'poor' | 'fair' | 'good' | 'excellent'
  energyGained: number
  createdAt: string
}

export interface ShopItem {
  id: number
  name: string
  description: string
  price: number
  sellPrice: number
  itemType: 'consumable' | 'equipment'
  icon: string
  image: string
  stock: number
}

export interface InventoryItem {
  id: number
  itemId: number
  name: string
  description: string
  itemType: string
  sellPrice: number
  icon: string
  image: string
  quantity: number
}

export interface PurchaseRecord {
  id: number
  itemName: string
  quantity: number
  totalPrice: number
  createdAt: string
}

export interface Task {
  id: number
  userId: number
  title: string
  description: string
  category: string
  type: 'once' | 'repeatable' | 'challenge'
  status: 'active' | 'completed' | 'failed' | 'deleted'
  deadline: string | null
  primaryAttribute: string
  difficulty: number
  rewardExp: number
  rewardSpiritStones: number
  rewardPhysique: number
  rewardWillpower: number
  rewardIntelligence: number
  rewardPerception: number
  rewardCharisma: number
  rewardAgility: number
  penaltyExp: number
  penaltySpiritStones: number
  fatigueCost: number
  dailyLimit: number
  totalLimit: number
  completedCount: number
  todayCompletionCount: number
  lastCompletedDate: string
  remindBefore: number
  remindInterval: number
  lastRemindedAt?: string | null
  completedAt?: string | null
  createdAt: string
  updatedAt: string
}

export interface ApiResponse<T = any> {
  code: number
  message: string
  data?: T
}

export interface BindCodeResponse {
  code: string
  botUsername: string
  expiresIn: number
}

export interface TgStatus {
  bound: boolean
  tgUsername: string
  tgChatId: number
}

export interface BarkStatus {
  enabled: boolean
  barkKey: string  // Masked key (first 8 chars + ***)
}

export interface CompleteTaskResult {
  task: Task
  character: CharacterStats
  message: string
}
