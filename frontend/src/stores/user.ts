import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import * as authApi from '@/api/auth'
import type { User } from '@/types'

export const useUserStore = defineStore('user', () => {
  const user = ref<User | null>(null)
  const isLoggedIn = computed(() => user.value !== null)

  async function login(username: string, password: string): Promise<void> {
    try {
      const response = await authApi.login(username, password) as any
      // Backend returns AuthResp { token, user } in data
      if (response.data?.token) {
        // Save token to localStorage
        localStorage.setItem('token', response.data.token)
        console.log('✅ Token saved to localStorage')
      }
      if (response.data?.user) {
        user.value = response.data.user
      } else if (response.data) {
        user.value = response.data
      }
    } catch (error) {
      throw error
    }
  }

  async function register(username: string, password: string): Promise<void> {
    try {
      const response = await authApi.register(username, password) as any
      if (response.data?.token) {
        // Save token to localStorage
        localStorage.setItem('token', response.data.token)
        console.log('✅ Token saved to localStorage')
      }
      if (response.data?.user) {
        user.value = response.data.user
      } else if (response.data) {
        user.value = response.data
      }
    } catch (error) {
      throw error
    }
  }

  async function logout(): Promise<void> {
    try {
      await authApi.logout()
      user.value = null
      // Clear token from localStorage
      localStorage.removeItem('token')
      console.log('✅ Token removed from localStorage')
    } catch (error) {
      user.value = null
      localStorage.removeItem('token')
      throw error
    }
  }

  async function fetchMe(): Promise<void> {
    try {
      const response = await authApi.getMe()
      if (response.data) {
        user.value = response.data
      }
    } catch (error) {
      user.value = null
      throw error
    }
  }

  async function initAuth(): Promise<void> {
    try {
      await fetchMe()
    } catch (error) {
      user.value = null
    }
  }

  return {
    user,
    isLoggedIn,
    login,
    register,
    logout,
    fetchMe,
    initAuth
  }
})
