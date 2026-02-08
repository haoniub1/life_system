import { defineStore } from 'pinia'
import { ref } from 'vue'
import * as taskApi from '@/api/task'
import type { Task, CompleteTaskResult } from '@/types'

export const useTaskStore = defineStore('task', () => {
  const tasks = ref<Task[]>([])
  const loading = ref(false)

  async function fetchTasks(type?: 'once' | 'repeatable' | 'challenge', status?: 'active' | 'completed' | 'failed' | 'deleted'): Promise<void> {
    loading.value = true
    try {
      const response = await taskApi.getTasks({ type, status })
      if (response.data) {
        // Backend returns {tasks: Task[]}, not Task[] directly
        if (Array.isArray(response.data)) {
          tasks.value = response.data
        } else if ((response.data as any).tasks && Array.isArray((response.data as any).tasks)) {
          tasks.value = (response.data as any).tasks
        } else {
          tasks.value = []
        }
      }
    } catch (error) {
      throw error
    } finally {
      loading.value = false
    }
  }

  async function createTask(data: Partial<Task>): Promise<Task> {
    try {
      const response = await taskApi.createTask(data)
      if (response.data) {
        tasks.value.push(response.data)
        return response.data
      }
      throw new Error('Failed to create task')
    } catch (error) {
      throw error
    }
  }

  async function completeTask(id: number): Promise<CompleteTaskResult | undefined> {
    try {
      console.log('📝 Completing task:', id)
      const response = await taskApi.completeTask(id)
      console.log('✅ Complete task response:', response)

      const index = tasks.value.findIndex(t => t.id === id)
      console.log('📊 Task index in array:', index)

      if (index !== -1 && response.data?.task) {
        // Update the task in the store with the completed task data
        tasks.value[index] = response.data.task
        console.log('✅ Task updated in store:', response.data.task)
      }

      // Return the complete result (includes task, character, and message)
      return response.data
    } catch (error) {
      console.error('❌ Complete task error:', error)
      throw error
    }
  }

  async function deleteTask(id: number): Promise<void> {
    try {
      await taskApi.deleteTask(id)
      tasks.value = tasks.value.filter(t => t.id !== id)
    } catch (error) {
      throw error
    }
  }

  async function updateTask(id: number, data: Partial<Task>): Promise<void> {
    try {
      const response = await taskApi.updateTask(id, data)
      const index = tasks.value.findIndex(t => t.id === id)
      if (index !== -1 && response.data) {
        tasks.value[index] = response.data
      }
    } catch (error) {
      throw error
    }
  }

  return {
    tasks,
    loading,
    fetchTasks,
    createTask,
    completeTask,
    deleteTask,
    updateTask
  }
})
