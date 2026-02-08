<template>
  <div v-if="character" class="character-section">
    <n-grid cols="1 s:2" responsive="screen" :x-gap="24" :y-gap="24">
      <!-- Main Character Card -->
      <n-grid-item>
        <n-card class="character-main-card" :segmented="{ content: 'hard', footer: 'hard' }">
          <div class="character-header">
            <div class="title-section">
              <h2 class="character-title">{{ character.title }}</h2>
              <div class="level-badge">
                <span class="level-number">LV.{{ character.level }}</span>
              </div>
            </div>
          </div>

          <!-- HP Bar -->
          <div class="stat-row">
            <div class="stat-label">❤️ 生命值</div>
            <div class="stat-bar-container">
              <n-progress
                type="line"
                :percentage="(character.hp / character.maxHp) * 100"
                :color="getHpColor()"
                show-indicator
              />
            </div>
            <div class="stat-value">{{ character.hp }}/{{ character.maxHp }}</div>
          </div>

          <!-- Experience Bar -->
          <div class="stat-row">
            <div class="stat-label">✨ 经验值</div>
            <div class="stat-bar-container">
              <n-progress
                type="line"
                :percentage="characterStore.expProgress"
                color="#6366f1"
                show-indicator
              />
            </div>
            <div class="stat-value">{{ character.exp }}/{{ characterStore.expForNextLevel }}</div>
          </div>

          <!-- Mental Power Bar -->
          <div class="stat-row">
            <div class="stat-label">🧠 脑力</div>
            <div class="stat-bar-container">
              <n-progress
                type="line"
                :percentage="Math.max(0, character.mentalPower)"
                :color="getPowerColor(character.mentalPower)"
                show-indicator
              />
            </div>
            <div class="stat-value">{{ character.mentalPower }}/100</div>
          </div>

          <!-- Physical Power Bar -->
          <div class="stat-row">
            <div class="stat-label">💪 体力</div>
            <div class="stat-bar-container">
              <n-progress
                type="line"
                :percentage="Math.max(0, character.physicalPower)"
                :color="getPowerColor(character.physicalPower)"
                show-indicator
              />
            </div>
            <div class="stat-value">{{ character.physicalPower }}/100</div>
          </div>

          <!-- Sleep Aid Accumulation -->
          <div class="sleep-aid-section">
            <div class="sleep-aid-item">
              <span class="sleep-aid-label">😴 心理助眠</span>
              <span class="sleep-aid-value">{{ character.mentalSleepAid }}</span>
            </div>
            <div class="sleep-aid-item">
              <span class="sleep-aid-label">🏃 身体助眠</span>
              <span class="sleep-aid-value">{{ character.physicalSleepAid }}</span>
            </div>
          </div>

          <!-- Gold -->
          <div class="gold-section">
            <span class="gold-label">🪙 金币</span>
            <span class="gold-value">{{ character.gold }}</span>
          </div>
        </n-card>
      </n-grid-item>

      <!-- Attributes Card -->
      <n-grid-item>
        <n-card class="attributes-card" :segmented="{ content: 'hard' }">
          <h3 class="attributes-title">属性</h3>
          <n-space vertical :size="16">
            <div class="attribute-item">
              <div class="attribute-header">
                <span class="attribute-label">💪 力量</span>
                <span class="attribute-value">{{ character.strength }}</span>
              </div>
              <n-progress
                type="line"
                :percentage="(character.strength / 100) * 100"
                color="#ef4444"
              />
            </div>

            <div class="attribute-item">
              <div class="attribute-header">
                <span class="attribute-label">🧠 智力</span>
                <span class="attribute-value">{{ character.intelligence }}</span>
              </div>
              <n-progress
                type="line"
                :percentage="(character.intelligence / 100) * 100"
                color="#3b82f6"
              />
            </div>

            <div class="attribute-item">
              <div class="attribute-header">
                <span class="attribute-label">❤️ 体力</span>
                <span class="attribute-value">{{ character.vitality }}</span>
              </div>
              <n-progress
                type="line"
                :percentage="(character.vitality / 100) * 100"
                color="#ec4899"
              />
            </div>

            <div class="attribute-item">
              <div class="attribute-header">
                <span class="attribute-label">✨ 精神</span>
                <span class="attribute-value">{{ character.spirit }}</span>
              </div>
              <n-progress
                type="line"
                :percentage="(character.spirit / 100) * 100"
                color="#8b5cf6"
              />
            </div>
          </n-space>
        </n-card>
      </n-grid-item>
    </n-grid>

    <!-- Stats Grid -->
    <n-grid cols="2 s:4" responsive="screen" :x-gap="16" :y-gap="16" style="margin-top: 24px">
      <n-grid-item>
        <n-card class="stat-card" :segmented="{ content: 'hard' }">
          <div class="stat-card-content">
            <div class="stat-card-value">{{ character.level }}</div>
            <div class="stat-card-label">等级</div>
          </div>
        </n-card>
      </n-grid-item>

      <n-grid-item>
        <n-card class="stat-card" :segmented="{ content: 'hard' }">
          <div class="stat-card-content">
            <div class="stat-card-value">{{ character.gold }}</div>
            <div class="stat-card-label">金币</div>
          </div>
        </n-card>
      </n-grid-item>

      <n-grid-item>
        <n-card class="stat-card" :segmented="{ content: 'hard' }">
          <div class="stat-card-content">
            <div class="stat-card-value">{{ character.strength }}</div>
            <div class="stat-card-label">力量</div>
          </div>
        </n-card>
      </n-grid-item>

      <n-grid-item>
        <n-card class="stat-card" :segmented="{ content: 'hard' }">
          <div class="stat-card-content">
            <div class="stat-card-value">{{ character.intelligence }}</div>
            <div class="stat-card-label">智力</div>
          </div>
        </n-card>
      </n-grid-item>
    </n-grid>
  </div>
  <div v-else class="empty-state">
    <n-spin />
    <p>加载角色信息中...</p>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { NCard, NGrid, NGridItem, NSpace, NProgress, NSpin } from 'naive-ui'
import { useCharacterStore } from '@/stores/character'

const characterStore = useCharacterStore()

const character = computed(() => characterStore.character)

const getHpColor = () => {
  if (!character.value) return '#3b82f6'
  const percentage = (character.value.hp / character.value.maxHp) * 100
  if (percentage > 50) return '#10b981'
  if (percentage > 25) return '#f59e0b'
  return '#ef4444'
}

const getPowerColor = (power: number) => {
  if (power > 60) return '#10b981' // Green
  if (power > 30) return '#f59e0b' // Orange
  return '#ef4444' // Red
}
</script>

<style scoped>
.character-section {
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

.character-main-card,
.attributes-card,
.stat-card {
  background: linear-gradient(135deg, rgba(30, 30, 50, 0.8) 0%, rgba(20, 20, 40, 0.8) 100%);
  border: 1px solid rgba(255, 215, 0, 0.2);
  border-radius: 8px;
}

.character-header {
  margin-bottom: 24px;
  padding-bottom: 16px;
  border-bottom: 1px solid rgba(255, 215, 0, 0.2);
}

.title-section {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.character-title {
  font-size: 28px;
  font-weight: bold;
  background: linear-gradient(135deg, #ffd700, #ffed4e, #d4af37);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin: 0;
}

.level-badge {
  background: linear-gradient(135deg, #ffd700, #ffed4e);
  padding: 8px 16px;
  border-radius: 20px;
  font-weight: bold;
  color: #000;
  box-shadow: 0 4px 12px rgba(255, 215, 0, 0.3);
}

.level-number {
  font-size: 16px;
  font-weight: 700;
}

.stat-row {
  display: grid;
  grid-template-columns: 80px 1fr 100px;
  align-items: center;
  gap: 12px;
  margin-bottom: 16px;
}

.stat-label {
  font-weight: 600;
  color: #d0d0e0;
  white-space: nowrap;
}

.stat-bar-container {
  min-width: 0;
}

:deep(.n-progress) {
  --n-fill-color: #ffd700;
}

.stat-value {
  text-align: right;
  color: #a0a0b0;
  font-size: 12px;
  white-space: nowrap;
}

.sleep-aid-section {
  display: flex;
  gap: 16px;
  margin-bottom: 16px;
  padding: 12px;
  background: rgba(139, 92, 246, 0.1);
  border-radius: 8px;
  border: 1px solid rgba(139, 92, 246, 0.2);
}

.sleep-aid-item {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.sleep-aid-label {
  font-weight: 600;
  color: #c4b5fd;
  font-size: 13px;
}

.sleep-aid-value {
  font-size: 16px;
  font-weight: bold;
  color: #a78bfa;
}

.gold-section {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  background: rgba(255, 215, 0, 0.1);
  border-radius: 8px;
  margin-top: 16px;
}

.gold-label {
  font-weight: 600;
  color: #ffd700;
  font-size: 16px;
}

.gold-value {
  font-size: 20px;
  font-weight: bold;
  color: #ffed4e;
  margin-left: auto;
}

.attributes-title {
  margin-bottom: 16px;
  color: #d0d0e0;
  font-size: 18px;
  font-weight: bold;
}

.attribute-item {
  padding: 12px;
  background: rgba(255, 255, 255, 0.02);
  border-radius: 6px;
  border: 1px solid rgba(255, 215, 0, 0.1);
}

.attribute-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.attribute-label {
  font-weight: 600;
  color: #d0d0e0;
}

.attribute-value {
  font-weight: bold;
  color: #ffd700;
  font-size: 14px;
}

.stat-card {
  height: 100%;
}

.stat-card-content {
  text-align: center;
  padding: 12px 0;
}

.stat-card-value {
  font-size: 28px;
  font-weight: bold;
  background: linear-gradient(135deg, #ffd700, #ffed4e);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin-bottom: 8px;
}

.stat-card-label {
  font-size: 12px;
  color: #a0a0b0;
  font-weight: 600;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  color: #a0a0b0;
  text-align: center;
}

.empty-state p {
  margin-top: 16px;
  font-size: 14px;
}

@media (max-width: 768px) {
  .stat-row {
    grid-template-columns: 60px 1fr 70px;
    gap: 8px;
    margin-bottom: 12px;
  }

  .stat-label {
    font-size: 13px;
  }

  .character-title {
    font-size: 20px;
  }

  .level-badge {
    padding: 4px 10px;
  }

  .level-number {
    font-size: 13px;
  }

  .stat-card-value {
    font-size: 20px;
  }

  .sleep-aid-section {
    flex-direction: column;
    gap: 8px;
  }

  .gold-section {
    padding: 10px;
  }

  .attribute-item {
    padding: 8px;
  }

  .character-header {
    margin-bottom: 16px;
    padding-bottom: 12px;
  }
}
</style>
