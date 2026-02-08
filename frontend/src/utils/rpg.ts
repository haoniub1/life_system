export function expForLevel(level: number): number {
  return Math.floor(100 * Math.pow(1.5, level - 1))
}

export function getTitleForLevel(level: number): string {
  const titles: { [key: number]: string } = {
    1: '新手冒险者',
    5: '初级战士',
    10: '中级战士',
    15: '高级战士',
    20: '骑士',
    25: '传奇骑士',
    30: '黄金骑士',
    50: '传奇英雄',
    100: '世界英雄'
  }

  const levels = Object.keys(titles).map(Number).sort((a, b) => b - a)
  for (const lv of levels) {
    if (level >= lv) {
      return titles[lv]
    }
  }
  return titles[1]
}

export function getMaxHP(strength: number, vitality: number): number {
  return 100 + strength * 2 + vitality * 3
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

export function expProgressPercentage(currentExp: number, expForNext: number): number {
  if (expForNext === 0) return 0
  return Math.min((currentExp / expForNext) * 100, 100)
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
      color: '#f87171',      // red-400
      background: 'rgba(239, 68, 68, 0.3)'  // red-500/30
    }
  }

  const totalHours = diff / (1000 * 60 * 60)

  if (totalHours > 24) {
    return {
      level: 'safe',
      color: '#4ade80',      // green-400
      background: 'rgba(34, 197, 94, 0.2)'  // green-500/20
    }
  } else if (totalHours > 6) {
    return {
      level: 'normal',
      color: '#fbbf24',      // yellow-400
      background: 'rgba(245, 158, 11, 0.2)' // yellow-600/20
    }
  } else if (totalHours > 1) {
    return {
      level: 'warning',
      color: '#fb923c',      // orange-400
      background: 'rgba(249, 115, 22, 0.3)' // orange-500/30
    }
  } else {
    return {
      level: 'urgent',
      color: '#f87171',      // red-400
      background: 'rgba(239, 68, 68, 0.3)'  // red-500/30
    }
  }
}
