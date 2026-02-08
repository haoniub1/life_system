<template>
  <div class="timeline-container">
    <h2 class="section-title">活动时间线</h2>

    <!-- Stats Overview -->
    <n-grid cols="2 s:4" :x-gap="16" :y-gap="16" style="margin-bottom: 24px">
      <n-grid-item>
        <n-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon">✅</div>
            <div class="stat-value">{{ stats.tasksCompleted }}</div>
            <div class="stat-label">已完成任务</div>
          </div>
        </n-card>
      </n-grid-item>

      <n-grid-item>
        <n-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon">⭐</div>
            <div class="stat-value">{{ stats.totalExp }}</div>
            <div class="stat-label">总经验</div>
          </div>
        </n-card>
      </n-grid-item>

      <n-grid-item>
        <n-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon">🪙</div>
            <div class="stat-value">{{ stats.totalGold }}</div>
            <div class="stat-label">总金币</div>
          </div>
        </n-card>
      </n-grid-item>

      <n-grid-item>
        <n-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon">🌙</div>
            <div class="stat-value">{{ stats.sleepRecords }}</div>
            <div class="stat-label">睡眠记录</div>
          </div>
        </n-card>
      </n-grid-item>
    </n-grid>

    <!-- Filter Tabs -->
    <n-tabs v-model:value="activeFilter" type="segment" animated style="margin-bottom: 24px">
      <n-tab-pane name="all" tab="全部" />
      <n-tab-pane name="tasks" tab="任务" />
      <n-tab-pane name="sleep" tab="睡眠" />
      <n-tab-pane name="shop" tab="商店" />
    </n-tabs>

    <!-- Timeline -->
    <div v-if="loading" class="loading-state">
      <n-spin />
      <p>加载中...</p>
    </div>

    <div v-else-if="filteredEvents.length === 0" class="empty-state">
      <n-empty description="暂无活动记录" />
    </div>

    <n-timeline v-else>
      <n-timeline-item
        v-for="event in filteredEvents"
        :key="event.id"
        :type="getEventColor(event.type)"
        :title="event.title"
        :time="formatTime(event.timestamp)"
      >
        <template #icon>
          <span class="timeline-icon">{{ getEventIcon(event.type) }}</span>
        </template>
        <div class="activity-content">
          <p class="activity-description">{{ event.description }}</p>
          <div v-if="event.rewards" class="activity-rewards">
            <n-tag v-if="event.rewards.exp" type="success" size="small">
              +{{ event.rewards.exp }} EXP
            </n-tag>
            <n-tag v-if="event.rewards.gold && event.rewards.gold > 0" type="warning" size="small">
              +{{ event.rewards.gold }} 🪙
            </n-tag>
            <n-tag v-if="event.rewards.gold && event.rewards.gold < 0" type="error" size="small">
              {{ event.rewards.gold }} 🪙
            </n-tag>
            <n-tag v-if="event.rewards.energy" type="info" size="small">
              +{{ event.rewards.energy }} ⚡
            </n-tag>
          </div>
        </div>
      </n-timeline-item>
    </n-timeline>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useMessage } from 'naive-ui'
import {
  NCard,
  NGrid,
  NGridItem,
  NTabs,
  NTabPane,
  NTimeline,
  NTimelineItem,
  NTag,
  NEmpty,
  NSpin
} from 'naive-ui'
import request from '@/api/index'

interface TimelineRewards {
  exp?: number
  gold?: number
  energy?: number
}

interface TimelineEvent {
  id: string
  type: string
  title: string
  description: string
  timestamp: string
  rewards?: TimelineRewards
}

const message = useMessage()
const loading = ref(false)
const activeFilter = ref('all')
const events = ref<TimelineEvent[]>([])
const stats = ref({
  tasksCompleted: 0,
  totalExp: 0,
  totalGold: 0,
  sleepRecords: 0
})

const filteredEvents = computed(() => {
  if (activeFilter.value === 'all') {
    return events.value
  }

  const typeMap: { [key: string]: string[] } = {
    tasks: ['task_complete', 'task_fail', 'task_delete'],
    sleep: ['sleep'],
    shop: ['purchase']
  }

  const types = typeMap[activeFilter.value] || []
  return events.value.filter(e => types.includes(e.type))
})

const getEventIcon = (type: string): string => {
  const icons: { [key: string]: string } = {
    task_complete: '✅',
    task_fail: '❌',
    task_delete: '🗑',
    sleep: '🌙',
    purchase: '🛒'
  }
  return icons[type] || '📝'
}

const getEventColor = (type: string): any => {
  const colors: { [key: string]: any } = {
    task_complete: 'success',
    task_fail: 'error',
    task_delete: 'default',
    sleep: 'info',
    purchase: 'warning'
  }
  return colors[type] || 'default'
}

const formatTime = (timestamp: string): string => {
  const date = new Date(timestamp.replace(' ', 'T') + 'Z')
  if (isNaN(date.getTime())) {
    return timestamp
  }
  const now = new Date()
  const diff = now.getTime() - date.getTime()

  const minutes = Math.floor(diff / 60000)
  const hours = Math.floor(diff / 3600000)
  const days = Math.floor(diff / 86400000)

  if (minutes < 1) return '刚刚'
  if (minutes < 60) return `${minutes} 分钟前`
  if (hours < 24) return `${hours} 小时前`
  if (days < 7) return `${days} 天前`
  return date.toLocaleDateString('zh-CN')
}

const fetchTimeline = async () => {
  try {
    loading.value = true
    const response = await request.get('/timeline') as any
    if (response.data) {
      events.value = response.data.events || []
      stats.value = {
        tasksCompleted: response.data.tasksCompleted || 0,
        totalExp: response.data.totalExp || 0,
        totalGold: response.data.totalGold || 0,
        sleepRecords: response.data.sleepRecords || 0
      }
    }
  } catch (error: any) {
    console.error('Failed to fetch timeline:', error)
    message.error(error?.message || '获取时间线失败')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchTimeline()
})
</script>

<style scoped>
.timeline-container {
  animation: fadeIn 0.3s ease-in-out;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.section-title {
  font-size: 24px;
  font-weight: bold;
  background: linear-gradient(135deg, #ffd700, #ffed4e);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin: 0 0 24px 0;
}

.stat-card {
  background: linear-gradient(135deg, rgba(30, 30, 50, 0.8) 0%, rgba(20, 20, 40, 0.8) 100%);
  border: 1px solid rgba(255, 215, 0, 0.2);
}

.stat-content {
  text-align: center;
  padding: 12px 0;
}

.stat-icon {
  font-size: 32px;
  margin-bottom: 8px;
}

.stat-value {
  font-size: 24px;
  font-weight: bold;
  color: #ffd700;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 12px;
  color: #a0a0b0;
}

:deep(.n-tabs .n-tabs-nav) {
  background: rgba(30, 30, 50, 0.5);
  border-radius: 8px;
  padding: 4px;
}

:deep(.n-timeline) {
  padding-left: 20px;
}

.timeline-icon {
  font-size: 20px;
}

.activity-content {
  margin-top: 4px;
}

.activity-description {
  color: #d0d0e0;
  margin: 0 0 12px 0;
  font-size: 14px;
}

.activity-rewards {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.loading-state,
.empty-state {
  text-align: center;
  padding: 60px 20px;
  color: #a0a0b0;
}

@media (max-width: 768px) {
  .stat-icon {
    font-size: 24px;
  }

  .stat-value {
    font-size: 20px;
  }
}
</style>
