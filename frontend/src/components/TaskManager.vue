<template>
  <div class="task-manager">
    <div class="task-controls">
      <div class="filter-section">
        <div class="filter-group">
          <n-radio-group v-model:value="typeFilter" size="small">
            <n-radio-button value="">全部</n-radio-button>
            <n-radio-button value="once">一次性</n-radio-button>
            <n-radio-button value="repeatable">重复</n-radio-button>
            <n-radio-button value="challenge">挑战</n-radio-button>
          </n-radio-group>
        </div>

        <div class="filter-group">
          <n-radio-group v-model:value="statusFilter" size="small">
            <n-radio-button value="active">进行中</n-radio-button>
            <n-radio-button value="completed">已完成</n-radio-button>
            <n-radio-button value="">全部</n-radio-button>
          </n-radio-group>
        </div>
      </div>

      <n-button type="primary" size="small" @click="showCreateForm = true">
        + 新建
      </n-button>
    </div>

    <!-- Compact Tasks List -->
    <div v-if="filteredTasks.length > 0" class="tasks-list">
      <div
        v-for="task in filteredTasks"
        :key="task.id"
        class="task-item"
        :class="{ 
          'task-completing': completingTaskId === task.id,
          'task-completed': task.status === 'completed',
          'task-failed': task.status === 'failed'
        }"
        @mousedown="task.status === 'active' && startComplete(task.id, $event)"
        @mouseup="cancelComplete"
        @mouseleave="cancelComplete"
        @touchstart="task.status === 'active' && startComplete(task.id, $event)"
        @touchend="cancelComplete"
        @touchcancel="cancelComplete"
      >
        <!-- Progress overlay for long press -->
        <div 
          class="complete-progress" 
          :style="{ width: (completingTaskId === task.id || completedTaskId === task.id) ? '100%' : '0%' }"
        ></div>
        
        <div class="task-content">
          <div class="task-main">
            <span class="task-title">{{ task.title }}</span>
            <div class="task-rewards">
              <span v-if="task.rewardGold" class="reward-badge gold">🪙{{ task.rewardGold }}</span>
              <span v-if="task.rewardExp" class="reward-badge exp">⭐{{ task.rewardExp }}</span>
              <span v-if="task.rewardIntelligence" class="reward-badge attr">🧠+{{ task.rewardIntelligence }}</span>
              <span v-if="task.rewardVitality" class="reward-badge attr">💪+{{ task.rewardVitality }}</span>
            </div>
          </div>
          
          <div class="task-meta">
            <span class="task-type" :class="task.type">{{ getTaskTypeIcon(task.type) }}</span>
            <span v-if="task.type === 'repeatable' && task.dailyLimit" class="task-count">
              {{ task.todayCompletionCount }}/{{ task.dailyLimit }}
            </span>
            <span v-if="task.type === 'challenge' && task.deadline" class="task-deadline">
              ⏱{{ formatTimeRemaining(task.deadline) }}
            </span>
          </div>
        </div>

        <!-- More actions menu -->
        <n-dropdown
          v-if="task.status === 'active'"
          trigger="click"
          :options="getTaskActions(task)"
          @select="(key) => handleAction(key, task)"
        >
          <div 
            class="more-btn" 
            @click.stop 
            @mousedown.stop 
            @touchstart.stop
          >
            <span>⋮</span>
          </div>
        </n-dropdown>
        
        <div v-else class="status-badge">
          {{ task.status === 'completed' ? '✓' : '✗' }}
        </div>
      </div>
    </div>

    <!-- Empty State -->
    <n-empty v-else description="暂无任务" style="margin-top: 40px" />

    <!-- Create / Edit Task Modal -->
    <task-form
      v-if="showCreateForm"
      :task="editingTask || undefined"
      @submit="handleSubmitTask"
      @close="showCreateForm = false; editingTask = null"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import { useMessage } from 'naive-ui'
import {
  NButton,
  NRadioGroup,
  NRadioButton,
  NEmpty,
  NDropdown
} from 'naive-ui'
import { useTaskStore } from '@/stores/task'
import { useCharacterStore } from '@/stores/character'
import { formatTimeRemaining } from '@/utils/rpg'
import TaskForm from '@/components/TaskForm.vue'
import type { Task } from '@/types'

const message = useMessage()
const taskStore = useTaskStore()
const characterStore = useCharacterStore()

const typeFilter = ref('')
const statusFilter = ref('active')
const showCreateForm = ref(false)
const editingTask = ref<Task | null>(null)

// Long press completion
const completingTaskId = ref<number | null>(null)
const completedTaskId = ref<number | null>(null)  // 防止回退
const completeTimer = ref<number | null>(null)
const COMPLETE_DURATION = 2000 // 2 seconds

// 完成音效 - 使用用户提供的音效
const playCompleteSound = () => {
  const audio = new Audio('/complete.mp3')
  audio.volume = 0.7
  audio.play().catch(() => {})
}

// 撒花动画
const createConfetti = () => {
  const colors = ['#ffd700', '#10b981', '#818cf8', '#f59e0b', '#ec4899']
  const container = document.createElement('div')
  container.className = 'confetti-container'
  document.body.appendChild(container)

  for (let i = 0; i < 50; i++) {
    const confetti = document.createElement('div')
    confetti.className = 'confetti'
    confetti.style.left = Math.random() * 100 + 'vw'
    confetti.style.backgroundColor = colors[Math.floor(Math.random() * colors.length)]
    confetti.style.animationDelay = Math.random() * 0.5 + 's'
    confetti.style.animationDuration = (Math.random() * 1 + 1.5) + 's'
    container.appendChild(confetti)
  }

  setTimeout(() => container.remove(), 3000)
}

const filteredTasks = computed(() => taskStore.tasks)

const fetchCurrentTasks = async () => {
  try {
    await taskStore.fetchTasks(
      (typeFilter.value || undefined) as 'once' | 'repeatable' | 'challenge' | undefined,
      (statusFilter.value || undefined) as 'active' | 'completed' | 'failed' | 'deleted' | undefined
    )
  } catch (error) {
    console.error('Failed to fetch tasks:', error)
  }
}

watch([typeFilter, statusFilter], () => {
  fetchCurrentTasks()
})

const getTaskTypeIcon = (type: string): string => {
  const icons: { [key: string]: string } = {
    once: '📌',
    repeatable: '🔄',
    challenge: '⚔️'
  }
  return icons[type] || '📋'
}

const getTaskActions = (task: Task) => {
  return [
    { label: '✏️ 编辑', key: 'edit' },
    { label: '🗑️ 删除', key: 'delete' }
  ]
}

const handleAction = async (key: string, task: Task) => {
  if (key === 'edit') {
    editingTask.value = task
    showCreateForm.value = true
  } else if (key === 'delete') {
    try {
      await taskStore.deleteTask(task.id)
      message.success('任务已删除')
    } catch (error: any) {
      const errorMsg = error?.response?.data?.message || error?.message || '删除失败'
      message.error(errorMsg)
    }
  }
}

const startComplete = (taskId: number, event: MouseEvent | TouchEvent) => {
  event.preventDefault()
  // 如果已经完成，不再触发
  if (completedTaskId.value === taskId) return
  
  completingTaskId.value = taskId
  
  completeTimer.value = window.setTimeout(async () => {
    // 标记为已完成，防止进度条回退
    completedTaskId.value = taskId
    
    // 播放音效 + 撒花
    playCompleteSound()
    createConfetti()
    
    // 延迟一点再调用完成接口，让动画更明显
    setTimeout(async () => {
      await completeTask(taskId)
      completingTaskId.value = null
      completedTaskId.value = null
    }, 500)
  }, COMPLETE_DURATION)
}

const cancelComplete = () => {
  if (completeTimer.value) {
    clearTimeout(completeTimer.value)
    completeTimer.value = null
  }
  // 立即回退，但用快速动画
  if (!completedTaskId.value) {
    completingTaskId.value = null
  }
}

const completeTask = async (id: number) => {
  try {
    const result = await taskStore.completeTask(id)

    if (result?.character) {
      characterStore.character = result.character
    } else {
      await characterStore.fetchCharacter()
    }

    if (result?.message) {
      message.success(result.message, { duration: 3000 })
    } else {
      message.success('✅ 任务完成！')
    }

    setTimeout(() => fetchCurrentTasks(), 300)
  } catch (error: any) {
    const errorMsg = error?.response?.data?.message || error?.message || '完成失败'
    message.error(errorMsg)
  }
}

const handleSubmitTask = async (taskData: Partial<Task>) => {
  try {
    if (editingTask.value) {
      await taskStore.updateTask(editingTask.value.id, taskData)
      message.success('任务已更新')
    } else {
      await taskStore.createTask(taskData)
      message.success('任务创建成功！')
    }
    showCreateForm.value = false
    editingTask.value = null
    await fetchCurrentTasks()
  } catch (error: any) {
    const errorMsg = error?.response?.data?.message || error?.message || '操作失败'
    message.error(errorMsg)
  }
}

onMounted(async () => {
  await fetchCurrentTasks()
})

onUnmounted(() => {
  if (completeTimer.value) {
    clearTimeout(completeTimer.value)
  }
})
</script>

<style scoped>
.task-manager {
  animation: fadeIn 0.3s ease-in-out;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}

.task-controls {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  flex-wrap: wrap;
  gap: 12px;
}

.filter-section {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.filter-group {
  display: flex;
  align-items: center;
}

:deep(.n-radio-group) {
  background: transparent !important;
}

:deep(.n-radio-button) {
  background-color: rgba(255, 255, 255, 0.05) !important;
  color: #a0a0b0 !important;
  border-color: rgba(255, 215, 0, 0.2) !important;
  padding: 0 10px !important;
  font-size: 12px !important;
}

:deep(.n-radio-button--checked) {
  background-color: rgba(255, 215, 0, 0.2) !important;
  color: #ffd700 !important;
  border-color: #ffd700 !important;
}

/* Compact task list */
.tasks-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.task-item {
  position: relative;
  display: flex;
  align-items: center;
  padding: 12px 16px;
  background: linear-gradient(135deg, rgba(30, 30, 50, 0.9) 0%, rgba(25, 25, 45, 0.9) 100%);
  border: 1px solid rgba(255, 215, 0, 0.15);
  border-radius: 8px;
  cursor: pointer;
  user-select: none;
  overflow: hidden;
  transition: all 0.2s ease;
}

.task-item:hover {
  border-color: rgba(255, 215, 0, 0.4);
  transform: translateX(2px);
}

.task-item:active {
  transform: scale(0.99);
}

.task-completed {
  opacity: 0.5;
  cursor: default;
}

.task-failed {
  opacity: 0.5;
  border-color: rgba(239, 68, 68, 0.3);
  cursor: default;
}

/* Long press progress animation */
.complete-progress {
  position: absolute;
  left: 0;
  top: 0;
  height: 100%;
  background: linear-gradient(90deg, rgba(16, 185, 129, 0.3) 0%, rgba(16, 185, 129, 0.5) 100%);
  transition: width 1s ease-out;  /* 回退：1秒（快速） */
  pointer-events: none;
  z-index: 0;
}

.task-completing {
  border-color: rgba(16, 185, 129, 0.6) !important;
}

.task-completing .complete-progress {
  width: 100% !important;
  transition: width 2s linear;  /* 填充：2秒 */
}

.task-content {
  flex: 1;
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
  z-index: 1;
  min-width: 0;
}

.task-main {
  flex: 1;
  min-width: 0;
}

.task-title {
  font-size: 14px;
  font-weight: 500;
  color: #e0e0f0;
  display: block;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.task-rewards {
  display: flex;
  gap: 8px;
  margin-top: 4px;
}

.reward-badge {
  font-size: 11px;
  padding: 2px 6px;
  border-radius: 4px;
  white-space: nowrap;
}

.reward-badge.gold {
  background: rgba(255, 193, 7, 0.2);
  color: #ffc107;
}

.reward-badge.exp {
  background: rgba(16, 185, 129, 0.2);
  color: #10b981;
}

.reward-badge.attr {
  background: rgba(99, 102, 241, 0.2);
  color: #818cf8;
}

.task-meta {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-shrink: 0;
}

.task-type {
  font-size: 14px;
}

.task-count {
  font-size: 11px;
  color: #a0a0b0;
  background: rgba(255, 255, 255, 0.05);
  padding: 2px 6px;
  border-radius: 4px;
}

.task-deadline {
  font-size: 11px;
  color: #f59e0b;
  background: rgba(245, 158, 11, 0.1);
  padding: 2px 6px;
  border-radius: 4px;
}

/* More button (three dots) */
.more-btn {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 6px;
  color: #808090;
  font-size: 18px;
  font-weight: bold;
  transition: all 0.2s;
  z-index: 2;
}

.more-btn:hover {
  background: rgba(255, 215, 0, 0.1);
  color: #ffd700;
}

.status-badge {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
  color: #808090;
}

/* Dropdown menu styling */
:deep(.n-dropdown-option) {
  font-size: 13px !important;
}

/* Mobile responsive */
@media (max-width: 768px) {
  .task-controls {
    flex-direction: column;
    align-items: stretch;
  }

  .filter-section {
    flex-direction: column;
    gap: 8px;
  }

  .filter-group {
    overflow-x: auto;
    -webkit-overflow-scrolling: touch;
  }

  .task-item {
    padding: 10px 12px;
  }

  .task-title {
    font-size: 13px;
  }

  .task-rewards {
    flex-wrap: wrap;
  }

  .reward-badge {
    font-size: 10px;
  }
}
</style>
