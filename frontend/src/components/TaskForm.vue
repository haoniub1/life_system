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
      label-placement="left"
      label-width="120px"
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

      <!-- Category -->
      <n-form-item label="分类" path="category">
        <n-select
          v-model:value="formData.category"
          :options="categoryOptions"
          placeholder="选择分类"
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

      <!-- Deadline (optional for all types) -->
      <n-form-item label="截止时间" path="deadline">
        <n-date-picker
          v-model:value="formData.deadline"
          type="datetime"
          placeholder="选择截止时间（可选）"
          clearable
        />
      </n-form-item>

      <!-- Rewards Section -->
      <div class="section-title">奖励设置</div>

      <n-form-item label="经验值" path="rewardExp">
        <n-input-number
          v-model:value="formData.rewardExp"
          :min="0"
          :step="10"
          placeholder="0"
        />
      </n-form-item>

      <n-form-item label="金币" path="rewardGold">
        <n-input-number
          v-model:value="formData.rewardGold"
          :min="0"
          :step="10"
          placeholder="0"
        />
      </n-form-item>

      <div class="stats-row">
        <n-form-item label="力量" path="rewardStrength">
          <n-input-number
            v-model:value="formData.rewardStrength"
            :min="0"
            :step="0.5"
            placeholder="0"
          />
        </n-form-item>
        <n-form-item label="智力" path="rewardIntelligence">
          <n-input-number
            v-model:value="formData.rewardIntelligence"
            :min="0"
            :step="0.5"
            placeholder="0"
          />
        </n-form-item>
      </div>

      <div class="stats-row">
        <n-form-item label="体力" path="rewardVitality">
          <n-input-number
            v-model:value="formData.rewardVitality"
            :min="0"
            :step="0.5"
            placeholder="0"
          />
        </n-form-item>
        <n-form-item label="精神" path="rewardSpirit">
          <n-input-number
            v-model:value="formData.rewardSpirit"
            :min="0"
            :step="0.5"
            placeholder="0"
          />
        </n-form-item>
      </div>

      <!-- Energy Cost Section -->
      <div class="section-title">能量消耗</div>

      <div class="stats-row">
        <n-form-item label="脑力消耗" path="costMental">
          <n-input-number
            v-model:value="formData.costMental"
            :min="0"
            :step="5"
            placeholder="0"
          />
        </n-form-item>
        <n-form-item label="体力消耗" path="costPhysical">
          <n-input-number
            v-model:value="formData.costPhysical"
            :min="0"
            :step="5"
            placeholder="0"
          />
        </n-form-item>
      </div>

      <!-- Penalty Section (only for challenge type) -->
      <div v-if="formData.type === 'challenge'" class="section-title">惩罚设置</div>

      <n-form-item v-if="formData.type === 'challenge'" label="经验扣除" path="penaltyExp">
        <n-input-number
          v-model:value="formData.penaltyExp"
          :min="0"
          :step="10"
          placeholder="0"
        />
      </n-form-item>

      <n-form-item v-if="formData.type === 'challenge'" label="金币扣除" path="penaltyGold">
        <n-input-number
          v-model:value="formData.penaltyGold"
          :min="0"
          :step="10"
          placeholder="0"
        />
      </n-form-item>

      <!-- Repeat Limits (only for repeatable type) -->
      <div v-if="formData.type === 'repeatable'" class="section-title">重复限制</div>

      <n-form-item v-if="formData.type === 'repeatable'" label="每日限制" path="dailyLimit">
        <n-input-number
          v-model:value="formData.dailyLimit"
          :min="0"
          :step="1"
          placeholder="0 表示无限制"
        />
      </n-form-item>

      <n-form-item v-if="formData.type === 'repeatable'" label="总体限制" path="totalLimit">
        <n-input-number
          v-model:value="formData.totalLimit"
          :min="0"
          :step="1"
          placeholder="0 表示无限制"
        />
      </n-form-item>

      <!-- Telegram Reminder Section (shows when deadline is set) -->
      <template v-if="formData.deadline">
        <div class="section-title">Telegram提醒设置</div>

        <n-form-item label="提前提醒">
          <n-select
            v-model:value="formData.remindBefore"
            :options="remindBeforeOptions"
            placeholder="选择提前提醒时间"
          />
        </n-form-item>

        <n-form-item label="提醒间隔">
          <n-select
            v-model:value="formData.remindInterval"
            :options="remindIntervalOptions"
            placeholder="选择提醒间隔"
          />
        </n-form-item>

        <n-alert type="warning" style="margin-bottom: 16px">
          需要先绑定 Telegram 才能收到提醒
        </n-alert>
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
import { ref, watch, computed, onMounted } from 'vue'
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
  NAlert,
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

const formData = ref<any>({
  title: '',
  description: '',
  category: '',
  type: 'once',
  deadline: null,
  rewardExp: 0,
  rewardGold: 0,
  rewardStrength: 0,
  rewardIntelligence: 0,
  rewardVitality: 0,
  rewardSpirit: 0,
  penaltyExp: 0,
  penaltyGold: 0,
  dailyLimit: 0,
  totalLimit: 0,
  remindBefore: 30,
  remindInterval: 60,
  costMental: 10,
  costPhysical: 10
})

onMounted(() => {
  if (props.task) {
    formData.value = {
      title: props.task.title || '',
      description: props.task.description || '',
      category: props.task.category || '',
      type: props.task.type || 'once',
      deadline: props.task.deadline ? new Date(props.task.deadline).getTime() : null,
      rewardExp: props.task.rewardExp || 0,
      rewardGold: props.task.rewardGold || 0,
      rewardStrength: props.task.rewardStrength || 0,
      rewardIntelligence: props.task.rewardIntelligence || 0,
      rewardVitality: props.task.rewardVitality || 0,
      rewardSpirit: props.task.rewardSpirit || 0,
      penaltyExp: props.task.penaltyExp || 0,
      penaltyGold: props.task.penaltyGold || 0,
      dailyLimit: props.task.dailyLimit || 0,
      totalLimit: props.task.totalLimit || 0,
      remindBefore: props.task.remindBefore || 30,
      remindInterval: props.task.remindInterval || 60,
      costMental: props.task.costMental ?? 10,
      costPhysical: props.task.costPhysical ?? 10
    }
  }
})

const categoryOptions = [
  { label: '通用', value: '通用' },
  { label: '学习', value: '学习' },
  { label: '运动', value: '运动' },
  { label: '工作', value: '工作' },
  { label: '生活', value: '生活' }
]

// Auto-preset energy cost defaults based on category (only for new tasks)
const categoryCostDefaults: Record<string, { mental: number; physical: number }> = {
  '学习': { mental: 20, physical: 0 },
  '运动': { mental: 0, physical: 25 },
  '工作': { mental: 15, physical: 10 },
  '通用': { mental: 10, physical: 10 },
  '生活': { mental: 10, physical: 10 }
}

watch(
  () => formData.value.category,
  (newCategory) => {
    if (!isEditing.value && newCategory && categoryCostDefaults[newCategory]) {
      formData.value.costMental = categoryCostDefaults[newCategory].mental
      formData.value.costPhysical = categoryCostDefaults[newCategory].physical
    }
  }
)

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
    {
      required: true,
      message: '请输入任务标题',
      trigger: 'blur'
    },
    {
      min: 2,
      max: 100,
      message: '标题长度2-100字',
      trigger: 'blur'
    }
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

    submitting.value = true

    const submitData: Partial<Task> = {
      title: formData.value.title,
      description: formData.value.description,
      category: formData.value.category,
      type: formData.value.type as any,
      deadline: formData.value.deadline ? new Date(formData.value.deadline).toISOString() : null,
      rewardExp: formData.value.rewardExp,
      rewardGold: formData.value.rewardGold,
      rewardStrength: formData.value.rewardStrength,
      rewardIntelligence: formData.value.rewardIntelligence,
      rewardVitality: formData.value.rewardVitality,
      rewardSpirit: formData.value.rewardSpirit,
      penaltyExp: formData.value.penaltyExp,
      penaltyGold: formData.value.penaltyGold,
      dailyLimit: formData.value.dailyLimit,
      totalLimit: formData.value.totalLimit,
      remindBefore: formData.value.remindBefore,
      remindInterval: formData.value.remindInterval,
      costMental: formData.value.costMental,
      costPhysical: formData.value.costPhysical
    }

    emit('submit', submitData)
    showModal.value = false
  } catch (error: any) {
    message.error(error?.message || '验证失败')
  } finally {
    submitting.value = false
  }
}

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

:deep(.n-input__input),
:deep(.n-input-number__input) {
  background-color: rgba(255, 255, 255, 0.05) !important;
  color: #e0e0e0 !important;
  border-color: rgba(255, 215, 0, 0.2) !important;
}

:deep(.n-input__input::placeholder),
:deep(.n-input-number__input::placeholder) {
  color: #707080 !important;
}

:deep(.n-input__input:focus),
:deep(.n-input-number__input:focus) {
  border-color: #ffd700 !important;
  box-shadow: 0 0 8px rgba(255, 215, 0, 0.3) !important;
}

:deep(.n-select) {
  --n-color: rgba(255, 255, 255, 0.05);
  --n-border-color: rgba(255, 215, 0, 0.2);
}

:deep(.n-select__input) {
  color: #e0e0e0 !important;
}

:deep(.n-radio-group) {
  background: transparent;
}

:deep(.n-radio) {
  color: #d0d0e0 !important;
}

:deep(.n-radio__label) {
  color: #d0d0e0 !important;
}

:deep(.n-date-picker) {
  --n-color: rgba(255, 255, 255, 0.05);
  --n-border-color: rgba(255, 215, 0, 0.2);
}

:deep(.n-date-picker__input) {
  color: #e0e0e0 !important;
}

.section-title {
  font-size: 14px;
  font-weight: 600;
  color: #ffd700;
  margin: 20px 0 16px 0;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  border-bottom: 1px solid rgba(255, 215, 0, 0.2);
  padding-bottom: 8px;
}

.stats-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

:deep(.n-alert) {
  background: rgba(245, 158, 11, 0.1) !important;
  border-color: rgba(245, 158, 11, 0.3) !important;
  color: #fbbf24 !important;
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
  }
}
</style>
