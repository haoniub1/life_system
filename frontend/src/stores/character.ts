import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import * as characterApi from '@/api/character'
import type { CharacterStats, CharacterAttribute } from '@/types'

export const useCharacterStore = defineStore('character', () => {
  const character = ref<CharacterStats | null>(null)

  const highestRealm = computed(() => {
    if (!character.value?.attributes?.length) return null
    return character.value.attributes.reduce((max, attr) => {
      if (attr.realm > max.realm || (attr.realm === max.realm && attr.subRealm > max.subRealm)) {
        return attr
      }
      return max
    }, character.value.attributes[0])
  })

  function getAttributeByKey(key: string): CharacterAttribute | undefined {
    return character.value?.attributes?.find(a => a.attrKey === key)
  }

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
    highestRealm,
    getAttributeByKey,
    fetchCharacter,
    updateCharacter
  }
})
