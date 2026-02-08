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
  level: number
  exp: number
  strength: number
  intelligence: number
  vitality: number
  spirit: number
  hp: number
  maxHp: number
  gold: number
  title: string
  lastActivityDate: string
  energy: number
  maxEnergy: number
  mentalPower: number
  physicalPower: number
  mentalSleepAid: number
  physicalSleepAid: number
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
  itemType: 'consumable' | 'permanent'
  effect: string
  effectValue: number
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
  effect: string
  effectValue: number
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
  rewardExp: number
  rewardGold: number
  rewardStrength: number
  rewardIntelligence: number
  rewardVitality: number
  rewardSpirit: number
  penaltyExp: number
  penaltyGold: number
  dailyLimit: number
  totalLimit: number
  completedCount: number
  todayCompletionCount: number
  lastCompletedDate: string
  remindBefore: number
  remindInterval: number
  lastRemindedAt?: string | null
  createdAt: string
  updatedAt: string
  costMental: number
  costPhysical: number
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

export interface CompleteTaskResult {
  task: Task
  character: CharacterStats
  message: string
}
