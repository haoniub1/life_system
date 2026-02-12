<template>
  <div v-if="character" class="character-section">
    <!-- Main Info Card -->
    <n-card class="character-main-card" :segmented="{ content: 'soft', footer: 'soft' }">
      <div class="character-header">
        <div class="title-section">
          <h2 class="character-title">{{ character.title }}</h2>
          <div v-if="highestRealmAttr" class="realm-badge">
            <span class="realm-name">{{ highestRealmAttr.realmName }}{{ highestRealmAttr.subRealmName }}</span>
          </div>
        </div>
      </div>

      <!-- Overdraft warning -->
      <div v-if="character.overdraftPenalty > 0" class="overdraft-warning">
        <span>&#x26A0;&#xFE0F; 透支惩罚: -{{ character.overdraftPenalty }}% 效率</span>
      </div>

      <!-- Spirit Stones Breakdown (detailed view in card) -->
      <div class="spirit-stones-section">
        <span class="spirit-label">灵石明细</span>
        <div class="spirit-stones-grid">
          <div v-if="spiritDisplay.supreme > 0" class="stone-item stone-supreme">
            <span class="stone-icon">🔮</span>
            <span class="stone-count">{{ spiritDisplay.supreme }}</span>
            <span class="stone-type">极品</span>
          </div>
          <div v-if="spiritDisplay.high > 0" class="stone-item stone-high">
            <span class="stone-icon">💠</span>
            <span class="stone-count">{{ spiritDisplay.high }}</span>
            <span class="stone-type">上品</span>
          </div>
          <div v-if="spiritDisplay.medium > 0" class="stone-item stone-medium">
            <span class="stone-icon">💎</span>
            <span class="stone-count">{{ spiritDisplay.medium }}</span>
            <span class="stone-type">中品</span>
          </div>
          <div class="stone-item stone-low">
            <span class="stone-icon">🪨</span>
            <span class="stone-count">{{ spiritDisplay.low }}</span>
            <span class="stone-type">下品</span>
          </div>
        </div>
      </div>
    </n-card>

    <!-- 6 Attributes Grid -->
    <div class="attributes-grid">
      <div
        v-for="attr in mainAttributes"
        :key="attr.attrKey"
        class="attribute-card"
        :class="{ 'attr-bottleneck': attr.isBottleneck }"
        :style="{ '--attr-color': attr.color }"
      >
        <div class="attr-top">
          <span class="attr-emoji">{{ attr.emoji }}</span>
          <span class="attr-name">{{ attr.displayName }}</span>
          <n-popover trigger="click" placement="bottom">
            <template #trigger>
              <span class="attr-info-icon">&#x24D8;</span>
            </template>
            <div class="attr-info-popover">
              <div class="attr-info-title">{{ attr.emoji }} {{ attr.displayName }}</div>
              <p class="attr-info-desc">{{ getAttrDescription(attr.attrKey) }}</p>
              <div class="attr-info-stats">
                <div class="attr-info-row"><span>当前值</span><span>{{ attr.value.toFixed(1) }}</span></div>
                <div class="attr-info-row"><span>境界</span><span>{{ attr.realmName }}·{{ attr.subRealmName }}</span></div>
                <div class="attr-info-row"><span>属性上限</span><span>{{ attr.attrCap }}</span></div>
                <div class="attr-info-row"><span>境界进度</span><span>{{ attr.progressPercent.toFixed(1) }}%</span></div>
                <div v-if="attr.isBottleneck" class="attr-info-row warning"><span>状态</span><span>瓶颈中 - 需要突破</span></div>
                <div v-if="attr.accumulationPool > 0" class="attr-info-row"><span>积累池</span><span>{{ attr.accumulationPool.toFixed(1) }}</span></div>
                <div v-if="attr.realmExp > 0" class="attr-info-row"><span>境界经验</span><span>{{ attr.realmExp }}</span></div>
              </div>
            </div>
          </n-popover>
          <span class="attr-value">{{ attr.value.toFixed(0) }}</span>
        </div>
        <div class="attr-realm-badge">
          {{ attr.realmName }}{{ attr.subRealmName }}
        </div>
        <div class="attr-progress-bar">
          <div
            class="attr-progress-fill"
            :style="{ width: attr.progressPercent + '%', backgroundColor: attr.color }"
          ></div>
        </div>
        <div class="attr-bottom">
          <span class="attr-progress-text">{{ attr.progressPercent.toFixed(0) }}%</span>
          <span v-if="attr.isBottleneck" class="bottleneck-badge">&#x1F6A7; 瓶颈</span>
          <span v-if="attr.accumulationPool > 0" class="accumulation-badge">&#x1F4E6; 积累 {{ attr.accumulationPool.toFixed(1) }}</span>
        </div>
      </div>
    </div>

  </div>
  <div v-else class="empty-state">
    <n-spin />
    <p>加载角色信息中...</p>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { NCard, NPopover, NSpin } from 'naive-ui'
import { useCharacterStore } from '@/stores/character'
import { decomposeSpiritStones, ATTR_DISPLAY } from '@/utils/rpg'

const characterStore = useCharacterStore()

const character = computed(() => characterStore.character)
const highestRealmAttr = computed(() => characterStore.highestRealm)

const mainAttributes = computed(() => {
  if (!character.value?.attributes) return []
  return character.value.attributes.filter(a => a.attrKey !== 'luck')
})

function getAttrDescription(key: string): string {
  return ATTR_DISPLAY[key]?.description || ''
}

const spiritDisplay = computed(() => {
  return decomposeSpiritStones(character.value?.spiritStones || 0)
})
</script>

<style scoped>
.character-section {
  animation: fadeIn 0.3s ease-in-out;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}

.character-main-card {
  background: linear-gradient(135deg, rgba(30, 30, 50, 0.8) 0%, rgba(20, 20, 40, 0.8) 100%);
  border: 1px solid rgba(255, 215, 0, 0.2);
  border-radius: 8px;
  margin-bottom: 16px;
}

.character-header {
  margin-bottom: 20px;
  padding-bottom: 16px;
  border-bottom: 1px solid rgba(255, 215, 0, 0.2);
}

.title-section {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.character-title {
  font-size: 24px;
  font-weight: bold;
  background: linear-gradient(135deg, #ffd700, #ffed4e, #d4af37);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin: 0;
}

.realm-badge {
  background: linear-gradient(135deg, #8b5cf6, #a78bfa);
  padding: 6px 14px;
  border-radius: 20px;
  font-weight: bold;
  color: #fff;
  font-size: 13px;
  box-shadow: 0 4px 12px rgba(139, 92, 246, 0.3);
}

.overdraft-warning {
  padding: 8px 12px;
  background: rgba(239, 68, 68, 0.15);
  border: 1px solid rgba(239, 68, 68, 0.3);
  border-radius: 6px;
  color: #f87171;
  font-size: 13px;
  margin-bottom: 16px;
}

.spirit-stones-section {
  padding: 14px;
  background: linear-gradient(135deg, rgba(99, 102, 241, 0.1) 0%, rgba(139, 92, 246, 0.1) 100%);
  border-radius: 8px;
  border: 1px solid rgba(139, 92, 246, 0.2);
}

.spirit-label {
  font-weight: 600;
  color: #a78bfa;
  font-size: 15px;
  display: block;
  margin-bottom: 10px;
}

.spirit-stones-grid {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}

.stone-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 8px 14px;
  border-radius: 8px;
  min-width: 60px;
}

.stone-icon {
  font-size: 20px;
  line-height: 1;
}

.stone-count {
  font-size: 18px;
  font-weight: bold;
}

.stone-type {
  font-size: 11px;
  margin-top: 2px;
}

.stone-supreme {
  background: rgba(255, 215, 0, 0.15);
  border: 1px solid rgba(255, 215, 0, 0.3);
  color: #ffd700;
}

.stone-high {
  background: rgba(168, 85, 247, 0.15);
  border: 1px solid rgba(168, 85, 247, 0.3);
  color: #c084fc;
}

.stone-medium {
  background: rgba(59, 130, 246, 0.15);
  border: 1px solid rgba(59, 130, 246, 0.3);
  color: #60a5fa;
}

.stone-low {
  background: rgba(156, 163, 175, 0.1);
  border: 1px solid rgba(156, 163, 175, 0.2);
  color: #9ca3af;
}

/* Attributes Grid */
.attributes-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
  margin-bottom: 12px;
}

.attribute-card {
  background: linear-gradient(135deg, rgba(30, 30, 50, 0.8) 0%, rgba(20, 20, 40, 0.8) 100%);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 10px;
  padding: 12px;
  transition: all 0.2s ease;
}

.attribute-card:hover {
  border-color: var(--attr-color, rgba(255, 215, 0, 0.3));
  box-shadow: 0 0 12px color-mix(in srgb, var(--attr-color, #ffd700) 20%, transparent);
}

.attr-bottleneck {
  border-color: rgba(239, 68, 68, 0.4) !important;
  box-shadow: 0 0 12px rgba(239, 68, 68, 0.15);
  animation: bottleneckPulse 2s ease-in-out infinite;
}

@keyframes bottleneckPulse {
  0%, 100% { box-shadow: 0 0 8px rgba(239, 68, 68, 0.15); }
  50% { box-shadow: 0 0 16px rgba(239, 68, 68, 0.3); }
}

.attr-top {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-bottom: 6px;
}

.attr-emoji { font-size: 18px; }
.attr-name { font-size: 13px; color: #d0d0e0; font-weight: 600; }
.attr-value {
  margin-left: auto;
  font-size: 16px;
  font-weight: bold;
  color: #ffd700;
}

.attr-realm-badge {
  display: inline-block;
  font-size: 11px;
  padding: 2px 8px;
  border-radius: 10px;
  background: rgba(139, 92, 246, 0.15);
  color: #c4b5fd;
  margin-bottom: 8px;
}

.attr-progress-bar {
  height: 4px;
  background: rgba(255, 255, 255, 0.06);
  border-radius: 2px;
  overflow: hidden;
  margin-bottom: 6px;
}

.attr-progress-fill {
  height: 100%;
  border-radius: 2px;
  transition: width 0.5s ease;
}

.attr-bottom {
  display: flex;
  align-items: center;
  gap: 6px;
  flex-wrap: wrap;
}

.attr-progress-text {
  font-size: 11px;
  color: #808090;
}

.bottleneck-badge {
  font-size: 10px;
  color: #f87171;
  background: rgba(239, 68, 68, 0.1);
  padding: 1px 6px;
  border-radius: 4px;
}

.accumulation-badge {
  font-size: 10px;
  color: #60a5fa;
  background: rgba(59, 130, 246, 0.1);
  padding: 1px 6px;
  border-radius: 4px;
}

/* Attr Info Icon */
.attr-info-icon {
  font-size: 14px;
  color: #606070;
  cursor: pointer;
  transition: color 0.2s;
  line-height: 1;
  user-select: none;
}

.attr-info-icon:hover {
  color: #a0a0b0;
}

.attr-info-popover {
  max-width: 260px;
}

.attr-info-title {
  font-size: 15px;
  font-weight: bold;
  margin-bottom: 8px;
}

.attr-info-desc {
  font-size: 13px;
  color: #999;
  line-height: 1.5;
  margin: 0 0 10px 0;
}

.attr-info-stats {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.attr-info-row {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
  padding: 3px 0;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
}

.attr-info-row span:first-child {
  color: #888;
}

.attr-info-row span:last-child {
  font-weight: 600;
}

.attr-info-row.warning span:last-child {
  color: #f87171;
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
  .character-title { font-size: 18px; }

  .realm-badge {
    padding: 4px 10px;
    font-size: 11px;
  }

  .attributes-grid {
    grid-template-columns: 1fr;
  }

  .spirit-stones-grid {
    gap: 6px;
  }

  .stone-item {
    padding: 6px 10px;
    min-width: 50px;
  }

  .stone-count { font-size: 15px; }

  .character-header {
    margin-bottom: 14px;
    padding-bottom: 12px;
  }
}
</style>
