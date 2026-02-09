<template>
  <div class="bark-section">
    <n-card class="bark-card" :segmented="{ content: true }">
      <div v-if="!barkStatus?.enabled" class="unbound-section">
        <div class="section-header">
          <h2 class="section-title">🔔 Bark 推送绑定</h2>
          <p class="section-description">
            绑定 Bark 后，任务提醒会以闹钟模式推送到你的 iPhone，响铃 30 秒确保不会错过！
          </p>
        </div>

        <div class="setup-steps">
          <p class="section-subtitle">设置步骤：</p>
          <ol class="steps-list">
            <li>在 App Store 搜索并下载 <strong>Bark</strong></li>
            <li>打开 Bark App，复制推送 URL</li>
            <li>URL 格式：<code class="url-example">https://api.day.app/<span class="highlight">你的Key</span>/</code></li>
            <li>将 Key 部分粘贴到下方输入框</li>
          </ol>
        </div>

        <div class="input-section">
          <n-input
            v-model:value="barkKeyInput"
            placeholder="输入你的 Bark Key（如：z3i8rTvmNcLTtbUxzB4SQd）"
            size="large"
            :disabled="loading"
          >
            <template #prefix>
              <span style="color: #a0a0b0">🔑</span>
            </template>
          </n-input>
        </div>

        <div class="action-area">
          <n-button
            type="primary"
            size="large"
            :loading="loading"
            :disabled="!barkKeyInput.trim()"
            @click="saveBarkKey"
          >
            保存并测试
          </n-button>
        </div>
      </div>

      <div v-else class="bound-section">
        <div class="success-section">
          <div class="success-icon">✅</div>
          <div class="success-text">
            <h3 class="success-title">已绑定</h3>
            <p class="success-desc">Bark 推送已配置，任务提醒会以闹钟模式推送</p>
          </div>
        </div>

        <n-space vertical :size="16">
          <div class="bound-info">
            <span class="info-label">Bark Key：</span>
            <span class="info-value">{{ barkStatus?.barkKey }}</span>
          </div>
          <div class="bound-info">
            <span class="info-label">推送模式：</span>
            <span class="info-value">🔊 闹钟模式（alarm 铃声 30秒）</span>
          </div>
        </n-space>

        <div class="action-area" style="margin-top: 24px">
          <n-button type="info" :loading="testLoading" @click="testPush">
            📱 发送测试推送
          </n-button>
          <n-popconfirm
            positive-text="确定解除"
            negative-text="取消"
            @positive-click="unbindBark"
          >
            <template #trigger>
              <n-button type="error">
                🔗 解除绑定
              </n-button>
            </template>
            <p>确定要解除 Bark 绑定吗？解除后将无法通过 Bark 接收任务提醒。</p>
          </n-popconfirm>
        </div>
      </div>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useMessage } from 'naive-ui'
import {
  NCard,
  NButton,
  NSpace,
  NInput,
  NPopconfirm
} from 'naive-ui'
import * as barkApi from '@/api/bark'
import type { BarkStatus } from '@/types'

const message = useMessage()

const barkStatus = ref<BarkStatus | null>(null)
const barkKeyInput = ref('')
const loading = ref(false)
const testLoading = ref(false)

const fetchStatus = async () => {
  try {
    const response = await barkApi.getBarkStatus()
    if (response.data) {
      barkStatus.value = response.data
    }
  } catch (error: any) {
    console.error('获取 Bark 状态失败:', error)
  }
}

const saveBarkKey = async () => {
  const key = barkKeyInput.value.trim()
  if (!key) {
    message.warning('请输入 Bark Key')
    return
  }

  try {
    loading.value = true
    
    // 保存 key
    await barkApi.setBarkKey({ barkKey: key })
    
    // 发送测试推送
    await barkApi.testBark({
      title: '🎉 绑定成功！',
      body: 'Life System Bark 推送已配置完成'
    })
    
    message.success('Bark 绑定成功！请检查手机是否收到测试推送')
    barkKeyInput.value = ''
    await fetchStatus()
  } catch (error: any) {
    const errorMsg = error?.response?.data?.message || error?.message || '保存失败'
    message.error(errorMsg)
  } finally {
    loading.value = false
  }
}

const testPush = async () => {
  try {
    testLoading.value = true
    await barkApi.testBark({
      title: '📱 测试推送',
      body: '如果你看到这条消息，说明 Bark 推送正常工作！'
    })
    message.success('测试推送已发送，请检查手机')
  } catch (error: any) {
    const errorMsg = error?.response?.data?.message || error?.message || '发送失败'
    message.error(errorMsg)
  } finally {
    testLoading.value = false
  }
}

const unbindBark = async () => {
  try {
    await barkApi.deleteBarkKey()
    message.success('已解除 Bark 绑定')
    barkStatus.value = null
  } catch (error: any) {
    const errorMsg = error?.response?.data?.message || error?.message || '解除绑定失败'
    message.error(errorMsg)
  }
}

onMounted(async () => {
  await fetchStatus()
})
</script>

<style scoped>
.bark-section {
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

.bark-card {
  background: linear-gradient(135deg, rgba(30, 30, 50, 0.8) 0%, rgba(20, 20, 40, 0.8) 100%);
  border: 1px solid rgba(255, 140, 0, 0.2);
  border-radius: 8px;
  max-width: 600px;
}

.section-header {
  margin-bottom: 24px;
}

.section-title {
  font-size: 24px;
  font-weight: bold;
  background: linear-gradient(135deg, #ff8c00, #ffa500);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin: 0 0 12px 0;
}

.section-description {
  color: #a0a0b0;
  margin: 0;
  line-height: 1.6;
}

.section-subtitle {
  font-size: 14px;
  font-weight: 600;
  color: #d0d0e0;
  margin: 0 0 12px 0;
}

.setup-steps {
  margin-bottom: 24px;
  padding: 16px;
  background: rgba(255, 255, 255, 0.02);
  border-radius: 8px;
  border: 1px solid rgba(255, 140, 0, 0.2);
}

.steps-list {
  color: #d0d0e0;
  line-height: 1.8;
  padding-left: 20px;
  margin: 12px 0 0 0;
}

.steps-list li {
  margin-bottom: 8px;
}

.url-example {
  background: rgba(0, 0, 0, 0.3);
  padding: 2px 8px;
  border-radius: 4px;
  font-family: 'Courier New', monospace;
  color: #a0a0b0;
  font-size: 13px;
}

.highlight {
  color: #ff8c00;
  font-weight: bold;
}

.input-section {
  margin-bottom: 20px;
}

.action-area {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.success-section {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 24px;
  background: rgba(16, 185, 129, 0.1);
  border-radius: 8px;
  border: 1px solid rgba(16, 185, 129, 0.3);
  margin-bottom: 24px;
}

.success-icon {
  font-size: 36px;
  line-height: 1;
}

.success-text {
  flex: 1;
}

.success-title {
  font-size: 18px;
  font-weight: bold;
  color: #10b981;
  margin: 0 0 4px 0;
}

.success-desc {
  color: #6ee7b7;
  margin: 0;
  font-size: 14px;
}

.bound-info {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px;
  background: rgba(255, 255, 255, 0.02);
  border-radius: 6px;
  border: 1px solid rgba(255, 140, 0, 0.1);
}

.info-label {
  font-weight: 600;
  color: #d0d0e0;
}

.info-value {
  color: #ff8c00;
  font-family: 'Courier New', monospace;
  font-weight: 500;
}

:deep(.n-button--primary) {
  background: linear-gradient(135deg, #ff8c00, #ffa500);
  color: #000 !important;
  border: none !important;
}

:deep(.n-button--primary:hover) {
  box-shadow: 0 4px 16px rgba(255, 140, 0, 0.4) !important;
}

:deep(.n-input) {
  --n-border: 1px solid rgba(255, 140, 0, 0.3);
  --n-border-focus: 1px solid rgba(255, 140, 0, 0.6);
}

@media (max-width: 768px) {
  .bark-card {
    max-width: 100%;
  }

  .section-title {
    font-size: 20px;
  }

  .steps-list {
    font-size: 14px;
    padding-left: 16px;
  }

  .bound-info {
    flex-direction: column;
    align-items: flex-start;
    gap: 4px;
  }
}
</style>
