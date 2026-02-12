import type { SpiritStoneDisplay } from '@/types'

export const REALMS = [
  { name: '凡人', subRealms: ['初期','中期','后期','大圆满'] },
  { name: '炼气', subRealms: ['初期','中期','后期','大圆满'] },
  { name: '筑基', subRealms: ['初期','中期','后期','大圆满'] },
  { name: '金丹', subRealms: ['初期','中期','后期','大圆满'] },
  { name: '元婴', subRealms: ['初期','中期','后期','大圆满'] },
  { name: '化神', subRealms: ['初期','中期','后期','大圆满'] },
  { name: '合体', subRealms: ['初期','中期','后期','大圆满'] },
  { name: '大乘', subRealms: ['初期','中期','后期','大圆满'] },
  { name: '渡劫', subRealms: ['初期','中期','后期','大圆满'] },
]

export const ATTR_KEYS = ['physique','willpower','intelligence','perception','charisma','agility'] as const

export const ATTR_DISPLAY: Record<string, {emoji: string, name: string, color: string, description: string}> = {
  physique:     { emoji: '\u{1F4AA}', name: '体魄', color: '#ef4444', description: '代表身体素质与体能。通过运动、锻炼、健康饮食等方式提升。影响你的耐力、力量和身体健康。' },
  willpower:    { emoji: '\u{1F9E0}', name: '意志', color: '#8b5cf6', description: '代表意志力与自律能力。通过坚持习惯、克服困难、冥想等方式提升。影响你的专注力和抗压能力。' },
  intelligence: { emoji: '\u{1F4DA}', name: '智力', color: '#3b82f6', description: '代表学习能力与知识储备。通过阅读、学习、思考、解题等方式提升。影响你的分析能力和创造力。' },
  perception:   { emoji: '\u{1F441}', name: '感知', color: '#10b981', description: '代表观察力与洞察力。通过冥想、艺术鉴赏、自我反思等方式提升。影响你的直觉和对细节的敏感度。' },
  charisma:     { emoji: '\u2728', name: '魅力', color: '#ec4899', description: '代表社交能力与个人魅力。通过社交活动、演讲、形象管理等方式提升。影响你的人际关系和影响力。' },
  agility:      { emoji: '\u{1F3C3}', name: '敏捷', color: '#f59e0b', description: '代表反应速度与灵活性。通过协调训练、竞技运动、快速决策等方式提升。影响你的执行效率和应变能力。' },
  luck:         { emoji: '\u{1F340}', name: '幸运', color: '#6366f1', description: '隐藏属性，代表运气与机遇。无法直接提升，由系统随机波动。' },
}

export function decomposeSpiritStones(total: number): SpiritStoneDisplay {
  return {
    total,
    supreme: Math.floor(total / 1000000),
    high: Math.floor((total % 1000000) / 10000),
    medium: Math.floor((total % 10000) / 100),
    low: total % 100,
  }
}

export function formatSpiritStones(display: SpiritStoneDisplay): string {
  const parts: string[] = []
  if (display.supreme > 0) parts.push(`${display.supreme}极品`)
  if (display.high > 0) parts.push(`${display.high}上品`)
  if (display.medium > 0) parts.push(`${display.medium}中品`)
  if (display.low > 0) parts.push(`${display.low}下品`)
  return parts.length > 0 ? parts.join(' ') : '0下品'
}

export function formatTimeRemaining(deadline: string): string {
  const now = new Date()
  const end = new Date(deadline)
  const diff = end.getTime() - now.getTime()

  if (diff <= 0) {
    return '已截止'
  }

  const totalSeconds = Math.floor(diff / 1000)
  const days = Math.floor(totalSeconds / (24 * 3600))
  const hours = Math.floor((totalSeconds % (24 * 3600)) / 3600)
  const minutes = Math.floor((totalSeconds % 3600) / 60)

  if (days > 0) {
    return `${days}天${hours}小时`
  } else if (hours > 0) {
    return `${hours}小时${minutes}分钟`
  } else {
    return `${minutes}分钟`
  }
}

export interface DeadlineUrgency {
  level: 'safe' | 'normal' | 'warning' | 'urgent' | 'expired'
  color: string
  background: string
}

export function getDeadlineUrgency(deadline: string): DeadlineUrgency {
  const now = new Date()
  const end = new Date(deadline)
  const diff = end.getTime() - now.getTime()

  if (diff <= 0) {
    return {
      level: 'expired',
      color: '#f87171',
      background: 'rgba(239, 68, 68, 0.3)'
    }
  }

  const totalHours = diff / (1000 * 60 * 60)

  if (totalHours > 24) {
    return {
      level: 'safe',
      color: '#4ade80',
      background: 'rgba(34, 197, 94, 0.2)'
    }
  } else if (totalHours > 6) {
    return {
      level: 'normal',
      color: '#fbbf24',
      background: 'rgba(245, 158, 11, 0.2)'
    }
  } else if (totalHours > 1) {
    return {
      level: 'warning',
      color: '#fb923c',
      background: 'rgba(249, 115, 22, 0.3)'
    }
  } else {
    return {
      level: 'urgent',
      color: '#f87171',
      background: 'rgba(239, 68, 68, 0.3)'
    }
  }
}
