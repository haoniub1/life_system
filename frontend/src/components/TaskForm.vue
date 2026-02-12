<template>
  <n-modal
    v-model:show="showModal"
    preset="card"
    :title="isEditing ? '编辑任务' : '新建任务'"
    size="large"
    class="task-form-modal"
    @close="handleClose"
  >
    <n-form
      ref="formRef"
      :model="formData"
      :rules="rules"
      :label-placement="labelPlacement"
      :label-width="labelWidth"
    >
      <!-- Title -->
      <n-form-item label="标题" path="title">
        <n-input
          v-model:value="formData.title"
          placeholder="输入任务标题"
          clearable
        />
      </n-form-item>

      <!-- Description -->
      <n-form-item label="描述" path="description">
        <n-input
          v-model:value="formData.description"
          placeholder="输入任务描述"
          type="textarea"
          :rows="3"
          clearable
        />
      </n-form-item>

      <!-- Type -->
      <n-form-item label="任务类型" path="type">
        <n-radio-group v-model:value="formData.type" @update:value="handleTypeChange">
          <n-space>
            <n-radio value="once">一次性</n-radio>
            <n-radio value="repeatable">可重复</n-radio>
            <n-radio value="challenge">挑战</n-radio>
          </n-space>
        </n-radio-group>
      </n-form-item>

      <!-- Difficulty (0-5 stars) -->
      <n-form-item label="难度" path="difficulty">
        <div class="difficulty-section">
          <n-rate v-model:value="formData.difficulty" :count="5" clearable :on-update:value="handleDifficultyChange" />
          <div class="difficulty-preview">
            <span class="preview-item">😴 {{ DIFFICULTY_TABLE[formData.difficulty]?.fatigue ?? 0 }}</span>
            <span class="preview-item">💎 {{ DIFFICULTY_TABLE[formData.difficulty]?.spiritStones ?? 0 }}</span>
            <span v-if="formData.difficulty > 0" class="preview-item">📈 +{{ DIFFICULTY_TABLE[formData.difficulty]?.attrBonus ?? 0 }}</span>
          </div>
        </div>
      </n-form-item>

      <!-- Category: 6 attribute dropdowns -->
      <div class="section-title">任务分类（选择相关标签，自动计算属性加成）</div>

      <div class="category-grid">
        <div v-for="cat in categoryDefs" :key="cat.key" class="category-item">
          <div class="category-label" :style="{ color: cat.color }">{{ cat.emoji }} {{ cat.name }}</div>
          <n-select
            v-model:value="selectedCategories[cat.key]"
            :options="cat.options"
            multiple
            clearable
            placeholder="选择..."
            size="small"
            @update:value="applyTemplate"
          />
        </div>
      </div>

      <!-- Deadline (optional for all types) -->
      <n-form-item label="截止时间" path="deadline">
        <n-date-picker
          v-model:value="formData.deadline"
          type="datetime"
          placeholder="选择截止时间（可选）"
          clearable
        />
      </n-form-item>

      <!-- Rewards Section (auto-filled, user editable) -->
      <div class="section-title">奖励设置</div>

      <div class="stats-row">
        <n-form-item label="💎 灵石" path="rewardSpiritStones">
          <n-input-number
            v-model:value="formData.rewardSpiritStones"
            :min="0"
            :step="1"
            placeholder="0"
          />
        </n-form-item>
        <n-form-item label="😴 疲劳" path="fatigueCost">
          <n-input-number
            v-model:value="formData.fatigueCost"
            :min="0"
            :step="1"
            placeholder="0"
          />
        </n-form-item>
      </div>

      <div class="stats-row">
        <n-form-item label="💪 体魄" path="rewardPhysique">
          <n-input-number
            v-model:value="formData.rewardPhysique"
            :min="0"
            :max="10"
            :step="0.1"
            :precision="1"
            placeholder="0"
          />
        </n-form-item>
        <n-form-item label="🧠 意志" path="rewardWillpower">
          <n-input-number
            v-model:value="formData.rewardWillpower"
            :min="0"
            :max="10"
            :step="0.1"
            :precision="1"
            placeholder="0"
          />
        </n-form-item>
      </div>

      <div class="stats-row">
        <n-form-item label="📚 智力" path="rewardIntelligence">
          <n-input-number
            v-model:value="formData.rewardIntelligence"
            :min="0"
            :max="10"
            :step="0.1"
            :precision="1"
            placeholder="0"
          />
        </n-form-item>
        <n-form-item label="👁 感知" path="rewardPerception">
          <n-input-number
            v-model:value="formData.rewardPerception"
            :min="0"
            :max="10"
            :step="0.1"
            :precision="1"
            placeholder="0"
          />
        </n-form-item>
      </div>

      <div class="stats-row">
        <n-form-item label="✨ 魅力" path="rewardCharisma">
          <n-input-number
            v-model:value="formData.rewardCharisma"
            :min="0"
            :max="10"
            :step="0.1"
            :precision="1"
            placeholder="0"
          />
        </n-form-item>
        <n-form-item label="🏃 敏捷" path="rewardAgility">
          <n-input-number
            v-model:value="formData.rewardAgility"
            :min="0"
            :max="10"
            :step="0.1"
            :precision="1"
            placeholder="0"
          />
        </n-form-item>
      </div>

      <!-- Penalty Section (only for challenge type) -->
      <div v-if="formData.type === 'challenge'" class="section-title">惩罚设置</div>

      <div v-if="formData.type === 'challenge'" class="stats-row">
        <n-form-item label="经验扣除" path="penaltyExp">
          <n-input-number
            v-model:value="formData.penaltyExp"
            :min="0"
            :step="10"
            placeholder="0"
          />
        </n-form-item>
        <n-form-item label="灵石扣除" path="penaltySpiritStones">
          <n-input-number
            v-model:value="formData.penaltySpiritStones"
            :min="0"
            :step="1"
            placeholder="0"
          />
        </n-form-item>
      </div>

      <!-- Repeat Limits (only for repeatable type) -->
      <div v-if="formData.type === 'repeatable'" class="section-title">重复限制</div>

      <div v-if="formData.type === 'repeatable'" class="stats-row">
        <n-form-item label="每日限制" path="dailyLimit">
          <n-input-number
            v-model:value="formData.dailyLimit"
            :min="0"
            :step="1"
            placeholder="0 = 无限"
          />
        </n-form-item>
        <n-form-item label="总体限制" path="totalLimit">
          <n-input-number
            v-model:value="formData.totalLimit"
            :min="0"
            :step="1"
            placeholder="0 = 无限"
          />
        </n-form-item>
      </div>

      <!-- Telegram Reminder Section (shows when deadline is set) -->
      <template v-if="formData.deadline">
        <div class="section-title">提醒设置</div>

        <div class="stats-row">
          <n-form-item label="提前提醒">
            <n-select
              v-model:value="formData.remindBefore"
              :options="remindBeforeOptions"
              placeholder="提前提醒时间"
            />
          </n-form-item>
          <n-form-item label="提醒间隔">
            <n-select
              v-model:value="formData.remindInterval"
              :options="remindIntervalOptions"
              placeholder="提醒间隔"
            />
          </n-form-item>
        </div>
      </template>
    </n-form>

    <template #footer>
      <n-space justify="end">
        <n-button @click="handleClose">取消</n-button>
        <n-button type="primary" @click="handleSubmit" :loading="submitting">
          {{ isEditing ? '保存修改' : '创建任务' }}
        </n-button>
      </n-space>
    </template>
  </n-modal>
</template>

<script setup lang="ts">
import { ref, reactive, watch, computed, onMounted, onUnmounted } from 'vue'
import { useMessage } from 'naive-ui'
import {
  NModal,
  NForm,
  NFormItem,
  NInput,
  NInputNumber,
  NButton,
  NSpace,
  NSelect,
  NRadioGroup,
  NRadio,
  NDatePicker,
  NRate,
  type FormInst
} from 'naive-ui'
import type { Task } from '@/types'

const props = defineProps<{
  task?: Task
}>()

const emit = defineEmits<{
  submit: [data: Partial<Task>]
  close: []
}>()

const isEditing = computed(() => !!props.task)

const message = useMessage()
const formRef = ref<FormInst | null>(null)
const submitting = ref(false)
const showModal = ref(true)
const windowWidth = ref(window.innerWidth)
const onResize = () => { windowWidth.value = window.innerWidth }
const labelPlacement = computed(() => windowWidth.value <= 768 ? 'top' : 'left')
const labelWidth = computed(() => windowWidth.value <= 768 ? 'auto' : '120px')

// Difficulty table
const DIFFICULTY_TABLE: Record<number, { fatigue: number, spiritStones: number, attrBonus: number }> = {
  0: { fatigue: 1, spiritStones: 10, attrBonus: 0 },
  1: { fatigue: 5, spiritStones: 50, attrBonus: 0.1 },
  2: { fatigue: 10, spiritStones: 120, attrBonus: 0.2 },
  3: { fatigue: 20, spiritStones: 300, attrBonus: 0.4 },
  4: { fatigue: 40, spiritStones: 800, attrBonus: 0.7 },
  5: { fatigue: 90, spiritStones: 2500, attrBonus: 1.0 },
}

// Category definitions with tag options per attribute
const categoryDefs = [
  {
    key: 'physique' as const,
    name: '体魄',
    emoji: '💪',
    color: '#ef4444',
    options: [
      { label: '运动', value: '运动' },
      { label: '体能', value: '体能' },
      { label: '健康', value: '健康' },
      { label: '饮食', value: '饮食' },
      { label: '睡眠', value: '睡眠' },
    ]
  },
  {
    key: 'intelligence' as const,
    name: '智力',
    emoji: '📚',
    color: '#3b82f6',
    options: [
      { label: '学习', value: '学习' },
      { label: '思考', value: '思考' },
      { label: '知识积累', value: '知识积累' },
      { label: '阅读', value: '阅读' },
      { label: '编程', value: '编程' },
    ]
  },
  {
    key: 'charisma' as const,
    name: '魅力',
    emoji: '✨',
    color: '#ec4899',
    options: [
      { label: '沟通', value: '沟通' },
      { label: '人脉', value: '人脉' },
      { label: '表达能力', value: '表达能力' },
      { label: '形象管理', value: '形象管理' },
      { label: '社交', value: '社交' },
    ]
  },
  {
    key: 'willpower' as const,
    name: '意志',
    emoji: '🧠',
    color: '#8b5cf6',
    options: [
      { label: '自律', value: '自律' },
      { label: '专注', value: '专注' },
      { label: '抗压', value: '抗压' },
      { label: '冥想', value: '冥想' },
      { label: '戒瘾', value: '戒瘾' },
    ]
  },
  {
    key: 'agility' as const,
    name: '敏捷',
    emoji: '🏃',
    color: '#f59e0b',
    options: [
      { label: '执行效率', value: '执行效率' },
      { label: '反应速度', value: '反应速度' },
      { label: '手工技巧', value: '手工技巧' },
      { label: '整理收纳', value: '整理收纳' },
      { label: '打字', value: '打字' },
    ]
  },
  {
    key: 'perception' as const,
    name: '感知',
    emoji: '👁',
    color: '#10b981',
    options: [
      { label: '审美', value: '审美' },
      { label: '艺术', value: '艺术' },
      { label: '观察力', value: '观察力' },
      { label: '想象力', value: '想象力' },
      { label: '音乐', value: '音乐' },
    ]
  },
]

type AttrKey = 'physique' | 'intelligence' | 'charisma' | 'willpower' | 'agility' | 'perception'

const selectedCategories = reactive<Record<AttrKey, string[]>>({
  physique: [],
  intelligence: [],
  charisma: [],
  willpower: [],
  agility: [],
  perception: [],
})

const formData = ref<any>({
  title: '',
  description: '',
  category: '',
  type: 'once',
  deadline: null,
  difficulty: 1,
  rewardExp: 0,
  rewardSpiritStones: 50,
  rewardPhysique: 0,
  rewardWillpower: 0,
  rewardIntelligence: 0,
  rewardPerception: 0,
  rewardCharisma: 0,
  rewardAgility: 0,
  penaltyExp: 0,
  penaltySpiritStones: 0,
  fatigueCost: 5,
  dailyLimit: 0,
  totalLimit: 0,
  remindBefore: 30,
  remindInterval: 60
})

// Build category string from selected tags
function buildCategoryString(): string {
  const tags: string[] = []
  for (const cat of categoryDefs) {
    const selected = selectedCategories[cat.key]
    if (selected && selected.length > 0) {
      tags.push(...selected)
    }
  }
  return tags.join(',')
}

// Apply template: set attribute bonuses based on selected categories + difficulty
function applyTemplate() {
  if (isEditing.value) return

  const diff = formData.value.difficulty
  const bonus = DIFFICULTY_TABLE[diff]?.attrBonus ?? 0

  const attrMap: Record<AttrKey, string> = {
    physique: 'rewardPhysique',
    willpower: 'rewardWillpower',
    intelligence: 'rewardIntelligence',
    perception: 'rewardPerception',
    charisma: 'rewardCharisma',
    agility: 'rewardAgility',
  }

  for (const key of Object.keys(attrMap) as AttrKey[]) {
    const hasSelection = selectedCategories[key] && selectedCategories[key].length > 0
    formData.value[attrMap[key]] = hasSelection ? bonus : 0
  }

  formData.value.category = buildCategoryString()
}

function handleDifficultyChange(val: number | null) {
  formData.value.difficulty = val ?? 0
  const entry = DIFFICULTY_TABLE[formData.value.difficulty]
  if (entry && !isEditing.value) {
    formData.value.fatigueCost = entry.fatigue
    formData.value.rewardSpiritStones = entry.spiritStones
  }
  applyTemplate()
}

onMounted(() => {
  window.addEventListener('resize', onResize)
  if (props.task) {
    formData.value = {
      title: props.task.title || '',
      description: props.task.description || '',
      category: props.task.category || '',
      type: props.task.type || 'once',
      deadline: props.task.deadline ? new Date(props.task.deadline).getTime() : null,
      difficulty: props.task.difficulty || 1,
      rewardExp: props.task.rewardExp || 0,
      rewardSpiritStones: props.task.rewardSpiritStones || 0,
      rewardPhysique: props.task.rewardPhysique || 0,
      rewardWillpower: props.task.rewardWillpower || 0,
      rewardIntelligence: props.task.rewardIntelligence || 0,
      rewardPerception: props.task.rewardPerception || 0,
      rewardCharisma: props.task.rewardCharisma || 0,
      rewardAgility: props.task.rewardAgility || 0,
      penaltyExp: props.task.penaltyExp || 0,
      penaltySpiritStones: props.task.penaltySpiritStones || 0,
      fatigueCost: props.task.fatigueCost ?? 2,
      dailyLimit: props.task.dailyLimit || 0,
      totalLimit: props.task.totalLimit || 0,
      remindBefore: props.task.remindBefore || 30,
      remindInterval: props.task.remindInterval || 60
    }
  } else {
    // Set defaults from difficulty 1
    handleDifficultyChange(1)
  }
})

const remindBeforeOptions = [
  { label: '5分钟', value: 5 },
  { label: '10分钟', value: 10 },
  { label: '15分钟', value: 15 },
  { label: '30分钟', value: 30 },
  { label: '1小时', value: 60 },
  { label: '2小时', value: 120 },
  { label: '1天', value: 1440 }
]

const remindIntervalOptions = [
  { label: '5分钟', value: 5 },
  { label: '10分钟', value: 10 },
  { label: '15分钟', value: 15 },
  { label: '30分钟', value: 30 },
  { label: '1小时', value: 60 }
]

const rules = {
  title: [
    { required: true, message: '请输入任务标题', trigger: 'blur' },
    { min: 2, max: 100, message: '标题长度2-100字', trigger: 'blur' }
  ]
}

const handleTypeChange = () => {
  if (formData.value.type === 'challenge') {
    formData.value.deadline = null
  }
}

watch(
  () => formData.value.deadline,
  (newVal) => {
    if (!newVal) {
      formData.value.remindBefore = 30
      formData.value.remindInterval = 60
    }
  }
)

const handleSubmit = async () => {
  try {
    await formRef.value?.validate()

    // 1星及以上必须选择至少一个分类
    if (formData.value.difficulty >= 1) {
      const hasCategory = Object.values(selectedCategories).some(arr => arr.length > 0)
      if (!hasCategory) {
        message.warning('1星及以上任务请至少选择一个分类标签')
        return
      }
    }

    submitting.value = true

    const submitData: Partial<Task> = {
      title: formData.value.title,
      description: formData.value.description,
      category: formData.value.category || buildCategoryString(),
      type: formData.value.type as any,
      deadline: formData.value.deadline ? new Date(formData.value.deadline).toISOString() : null,
      primaryAttribute: '',
      difficulty: formData.value.difficulty,
      rewardExp: formData.value.rewardExp,
      rewardSpiritStones: formData.value.rewardSpiritStones,
      rewardPhysique: formData.value.rewardPhysique,
      rewardWillpower: formData.value.rewardWillpower,
      rewardIntelligence: formData.value.rewardIntelligence,
      rewardPerception: formData.value.rewardPerception,
      rewardCharisma: formData.value.rewardCharisma,
      rewardAgility: formData.value.rewardAgility,
      penaltyExp: formData.value.penaltyExp,
      penaltySpiritStones: formData.value.penaltySpiritStones,
      fatigueCost: formData.value.fatigueCost,
      dailyLimit: formData.value.dailyLimit,
      totalLimit: formData.value.totalLimit,
      remindBefore: formData.value.remindBefore,
      remindInterval: formData.value.remindInterval
    }

    emit('submit', submitData)
    showModal.value = false
  } catch (error: any) {
    message.error(error?.message || '验证失败')
  } finally {
    submitting.value = false
  }
}

onUnmounted(() => {
  window.removeEventListener('resize', onResize)
})

const handleClose = () => {
  showModal.value = false
  emit('close')
}
</script>

<style scoped>
:deep(.n-modal) {
  background: linear-gradient(135deg, rgba(30, 30, 50, 0.95) 0%, rgba(20, 20, 40, 0.95) 100%);
  border: 1px solid rgba(255, 215, 0, 0.2);
}

:deep(.n-modal-mask) {
  background: rgba(0, 0, 0, 0.7);
}

:deep(.n-form-item__label) {
  color: #d0d0e0 !important;
  font-weight: 500;
}

:deep(.n-select) {
  --n-color: rgba(255, 255, 255, 0.05);
  --n-border-color: rgba(255, 215, 0, 0.2);
}

:deep(.n-radio) {
  color: #d0d0e0 !important;
}

:deep(.n-radio__label) {
  color: #d0d0e0 !important;
}

:deep(.n-rate) {
  --n-item-size: 22px;
}

.section-title {
  font-size: 13px;
  font-weight: 600;
  color: #ffd700;
  margin: 16px 0 12px 0;
  letter-spacing: 0.5px;
  border-bottom: 1px solid rgba(255, 215, 0, 0.2);
  padding-bottom: 6px;
}

.stats-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

.difficulty-section {
  display: flex;
  align-items: center;
  gap: 16px;
  flex-wrap: wrap;
}

.difficulty-preview {
  display: flex;
  gap: 12px;
  font-size: 12px;
  color: #a0a0b0;
}

.preview-item {
  white-space: nowrap;
}

/* Category grid */
.category-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 10px;
  margin-bottom: 16px;
}

.category-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.category-label {
  font-size: 12px;
  font-weight: 600;
}

:deep(.n-button--primary) {
  background: linear-gradient(135deg, #ffd700, #ffed4e) !important;
  color: #000 !important;
  border: none !important;
}

:deep(.n-button--primary:hover) {
  box-shadow: 0 4px 16px rgba(255, 215, 0, 0.4) !important;
}

@media (max-width: 768px) {
  :deep(.n-modal) {
    width: 95vw !important;
  }

  .stats-row {
    grid-template-columns: 1fr;
    gap: 0;
  }

  .category-grid {
    grid-template-columns: 1fr;
  }

  .section-title {
    font-size: 12px;
    margin: 10px 0 6px;
  }
}
</style>
