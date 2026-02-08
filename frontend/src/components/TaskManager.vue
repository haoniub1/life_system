<template>
  <div class="task-manager">
    <div class="task-controls">
      <div class="filter-section">
        <div class="filter-group">
          <span class="filter-label">任务类型:</span>
          <n-radio-group v-model:value="typeFilter" size="small">
            <n-radio-button value="">全部</n-radio-button>
            <n-radio-button value="once">一次性</n-radio-button>
            <n-radio-button value="repeatable">可重复</n-radio-button>
            <n-radio-button value="challenge">挑战</n-radio-button>
          </n-radio-group>
        </div>

        <div class="filter-group">
          <span class="filter-label">状态:</span>
          <n-radio-group v-model:value="statusFilter" size="small">
            <n-radio-button value="">全部</n-radio-button>
            <n-radio-button value="active">进行中</n-radio-button>
            <n-radio-button value="completed">已完成</n-radio-button>
            <n-radio-button value="failed">已失败</n-radio-button>
            <n-radio-button value="deleted">已删除</n-radio-button>
          </n-radio-group>
        </div>
      </div>

      <n-button type="primary" @click="showCreateForm = true">
        + 新建任务
      </n-button>
    </div>

    <!-- Tasks List -->
    <div v-if="filteredTasks.length > 0" class="tasks-list">
      <transition-group name="task" tag="div">
        <div
          v-for="task in filteredTasks"
          :key="task.id"
          class="task-card"
        >
          <n-card :segmented="{ content: 'hard', footer: 'soft' }">
            <template #header>
              <div class="task-header">
                <div class="task-header-left">
                  <h3 class="task-title">{{ task.title }}</h3>
                  <div class="task-badges">
                    <n-tag type="info" :bordered="false" size="small">
                      {{ getTaskTypeName(task.type) }}
                    </n-tag>
                    <n-tag :type="getStatusType(task.status)" :bordered="false" size="small">
                      {{ getStatusName(task.status) }}
                    </n-tag>
                  </div>
                </div>
                <div v-if="task.type === 'challenge' && task.deadline" class="deadline-badge">
                  <span
                    class="deadline-timer"
                    :style="getDeadlineStyle(task.deadline)"
                  >
                    ⏱ {{ formatTimeRemaining(task.deadline) }}
                  </span>
                </div>
              </div>
            </template>

            <div class="task-content">
              <p v-if="task.description" class="task-description">{{ task.description }}</p>

              <div class="task-meta">
                <span v-if="task.category" class="meta-item">
                  📁 {{ task.category }}
                </span>
                <span v-if="task.type === 'repeatable' && task.totalLimit" class="meta-item">
                  🔄 {{ task.completedCount }}/{{ task.totalLimit }}
                </span>
              </div>

              <!-- Rewards Section -->
              <div v-if="hasRewards(task)" class="rewards-section">
                <div class="section-title">奖励</div>
                <n-space>
                  <n-tag v-if="task.rewardExp" type="success" :bordered="false">
                    💫 经验: {{ task.rewardExp }}
                  </n-tag>
                  <n-tag v-if="task.rewardGold" type="warning" :bordered="false">
                    🪙 金币: {{ task.rewardGold }}
                  </n-tag>
                  <n-tag v-if="task.rewardStrength" type="error" :bordered="false">
                    💪 力量: +{{ task.rewardStrength }}
                  </n-tag>
                  <n-tag v-if="task.rewardIntelligence" type="info" :bordered="false">
                    🧠 智力: +{{ task.rewardIntelligence }}
                  </n-tag>
                  <n-tag v-if="task.rewardVitality" type="error" :bordered="false">
                    ❤️ 体力: +{{ task.rewardVitality }}
                  </n-tag>
                  <n-tag v-if="task.rewardSpirit" type="default" :bordered="false">
                    ✨ 精神: +{{ task.rewardSpirit }}
                  </n-tag>
                </n-space>
              </div>

              <!-- Penalty Section -->
              <div v-if="task.type === 'challenge' && (task.penaltyExp || task.penaltyGold)" class="penalty-section">
                <div class="section-title">失败惩罚</div>
                <n-space>
                  <n-tag v-if="task.penaltyExp" type="error" :bordered="false">
                    💫 经验: -{{ task.penaltyExp }}
                  </n-tag>
                  <n-tag v-if="task.penaltyGold" type="error" :bordered="false">
                    🪙 金币: -{{ task.penaltyGold }}
                  </n-tag>
                </n-space>
              </div>
            </div>

            <template #footer>
              <div class="task-actions">
                <n-button
                  v-if="task.status === 'active'"
                  type="success"
                  text
                  size="small"
                  @click="completeTask(task.id)"
                >
                  ✅ 完成
                </n-button>
                <n-button
                  v-if="task.status === 'active'"
                  type="info"
                  text
                  size="small"
                  @click="handleEditTask(task)"
                >
                  ✏️ 编辑
                </n-button>
                <n-popconfirm
                  v-if="task.status === 'active'"
                  positive-text="确定"
                  negative-text="取消"
                  @positive-click="deleteTask(task.id)"
                >
                  <template #trigger>
                    <n-button type="error" text size="small">
                      🗑 删除
                    </n-button>
                  </template>
                  <p>确定要删除这个任务吗？</p>
                </n-popconfirm>
              </div>
            </template>
          </n-card>
        </div>
      </transition-group>
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
  NCard,
  NTag,
  NButton,
  NRadioGroup,
  NRadioButton,
  NSpace,
  NEmpty,
  NPopconfirm
} from 'naive-ui'
import { useTaskStore } from '@/stores/task'
import { useCharacterStore } from '@/stores/character'
import { formatTimeRemaining, getDeadlineUrgency } from '@/utils/rpg'
import TaskForm from '@/components/TaskForm.vue'
import type { Task } from '@/types'

const message = useMessage()
const taskStore = useTaskStore()
const characterStore = useCharacterStore()

const typeFilter = ref('')
const statusFilter = ref('active')
const showCreateForm = ref(false)
const editingTask = ref<Task | null>(null)

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

const getTaskTypeName = (type: string): string => {
  const names: { [key: string]: string } = {
    once: '一次性',
    repeatable: '可重复',
    challenge: '挑战'
  }
  return names[type] || type
}

const getStatusName = (status: string): string => {
  const names: { [key: string]: string } = {
    active: '进行中',
    completed: '已完成',
    failed: '已失败',
    deleted: '已删除'
  }
  return names[status] || status
}

const getStatusType = (status: string) => {
  const types: { [key: string]: string } = {
    active: 'info',
    completed: 'success',
    failed: 'error',
    deleted: 'default'
  }
  return types[status] || 'default'
}

const hasRewards = (task: Task): boolean => {
  return !!(
    task.rewardExp ||
    task.rewardGold ||
    task.rewardStrength ||
    task.rewardIntelligence ||
    task.rewardVitality ||
    task.rewardSpirit
  )
}

const completeTask = async (id: number) => {
  try {
    const result = await taskStore.completeTask(id)

    // Update character stats with the returned character data
    if (result?.character) {
      characterStore.character = result.character
    } else {
      // Fallback: refresh character from server
      await characterStore.fetchCharacter()
    }

    // Show success message with rewards
    if (result?.message) {
      message.success(result.message, {
        duration: 4000
      })
    } else {
      message.success('任务已完成！获得奖励')
    }

    // Refresh task list
    setTimeout(() => fetchCurrentTasks(), 500)
  } catch (error: any) {
    const errorMsg = error?.response?.data?.message || error?.message || '完成失败'
    message.error(errorMsg)
    console.error('Complete task error:', error)
  }
}

const deleteTask = async (id: number) => {
  try {
    await taskStore.deleteTask(id)
    message.success('任务已删除')
  } catch (error: any) {
    const errorMsg = error?.response?.data?.message || error?.message || '删除失败'
    message.error(errorMsg)
  }
}

const handleEditTask = (task: Task) => {
  editingTask.value = task
  showCreateForm.value = true
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

// Auto-refresh timer for deadline countdown
const currentTime = ref(Date.now())

const getDeadlineStyle = (deadline: string) => {
  // Reference currentTime to ensure reactivity
  currentTime.value // eslint-disable-line
  const urgency = getDeadlineUrgency(deadline)
  return {
    color: urgency.color,
    background: urgency.background
  }
}

let timerInterval: number | null = null

onMounted(async () => {
  await fetchCurrentTasks()

  // Update current time every minute to refresh deadline display
  timerInterval = window.setInterval(() => {
    currentTime.value = Date.now()
  }, 60000) // 60 seconds
})

onUnmounted(() => {
  if (timerInterval) {
    clearInterval(timerInterval)
  }
})
</script>

<style scoped>
.task-manager {
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

.task-controls {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  flex-wrap: wrap;
  gap: 16px;
}

.filter-section {
  display: flex;
  gap: 24px;
  flex-wrap: wrap;
}

.filter-group {
  display: flex;
  align-items: center;
  gap: 12px;
}

.filter-label {
  font-weight: 600;
  color: #d0d0e0;
}

:deep(.n-radio-group) {
  background: transparent !important;
}

:deep(.n-radio-button) {
  background-color: rgba(255, 255, 255, 0.05) !important;
  color: #d0d0e0 !important;
  border-color: rgba(255, 215, 0, 0.2) !important;
}

:deep(.n-radio-button--checked) {
  background-color: rgba(255, 215, 0, 0.3) !important;
  color: #ffd700 !important;
  border-color: #ffd700 !important;
}

.tasks-list {
  display: grid;
  gap: 16px;
}

.task-card {
  animation: slideIn 0.3s ease-out;
}

@keyframes slideIn {
  from {
    opacity: 0;
    transform: translateX(-20px);
  }
  to {
    opacity: 1;
    transform: translateX(0);
  }
}

:deep(.n-card) {
  background: linear-gradient(135deg, rgba(30, 30, 50, 0.8) 0%, rgba(20, 20, 40, 0.8) 100%);
  border: 1px solid rgba(255, 215, 0, 0.2);
  border-radius: 8px;
}

.task-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 16px;
  width: 100%;
}

.task-header-left {
  flex: 1;
}

.task-title {
  font-size: 18px;
  font-weight: bold;
  color: #d0d0e0;
  margin: 0 0 8px 0;
}

.task-badges {
  display: flex;
  gap: 8px;
}

:deep(.n-tag) {
  background: rgba(255, 215, 0, 0.1) !important;
  color: #ffd700 !important;
}

.deadline-badge {
  white-space: nowrap;
}

.deadline-timer {
  padding: 6px 12px;
  border-radius: 6px;
  font-weight: 600;
  font-size: 12px;
  transition: all 0.3s ease;
}

.task-content {
  margin: 16px 0;
}

.task-description {
  color: #a0a0b0;
  margin: 0 0 12px 0;
  line-height: 1.5;
}

.task-meta {
  display: flex;
  gap: 16px;
  margin-bottom: 12px;
  flex-wrap: wrap;
}

.meta-item {
  font-size: 12px;
  color: #a0a0b0;
}

.rewards-section,
.penalty-section {
  margin: 16px 0;
  padding: 12px;
  background: rgba(255, 255, 255, 0.02);
  border-radius: 6px;
  border-left: 3px solid rgba(255, 215, 0, 0.3);
}

.section-title {
  font-size: 12px;
  font-weight: 600;
  color: #ffd700;
  margin-bottom: 8px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.task-actions {
  display: flex;
  gap: 8px;
}

.task-enter-active,
.task-leave-active {
  transition: all 0.3s ease;
}

.task-enter-from {
  opacity: 0;
  transform: translateX(-30px);
}

.task-leave-to {
  opacity: 0;
  transform: translateX(30px);
}

@media (max-width: 768px) {
  .task-controls {
    flex-direction: column;
    align-items: stretch;
  }

  .filter-section {
    flex-direction: column;
    gap: 12px;
  }

  .task-header {
    flex-direction: column;
  }

  .deadline-badge {
    width: 100%;
  }
}
</style>
