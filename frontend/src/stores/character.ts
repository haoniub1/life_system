import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import * as characterApi from '@/api/character'
import { expForLevel } from '@/utils/rpg'
import type { CharacterStats } from '@/types'

export const useCharacterStore = defineStore('character', () => {
  const character = ref<CharacterStats | null>(null)

  const expForNextLevel = computed(() => {
    if (!character.value) return 0
    return expForLevel(character.value.level + 1)
  })

  const expProgress = computed(() => {
    if (!character.value) return 0
    const totalExpForNext = expForNextLevel.value
    if (totalExpForNext === 0) return 0
    return Math.min((character.value.exp / totalExpForNext) * 100, 100)
  })

  async function fetchCharacter(): Promise<void> {
    try {
      const response = await characterApi.getCharacter()
      if (response.data) {
        character.value = response.data
      }
    } catch (error) {
      throw error
    }
  }

  async function updateCharacter(data: Partial<CharacterStats>): Promise<void> {
    try {
      const response = await characterApi.updateCharacter(data)
      if (response.data) {
        character.value = response.data
      }
    } catch (error) {
      throw error
    }
  }

  return {
    character,
    expForNextLevel,
    expProgress,
    fetchCharacter,
    updateCharacter
  }
})
