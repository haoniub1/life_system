<template>
  <div class="telegram-section">
    <n-card class="telegram-card" :segmented="{ content: 'hard' }">
      <div v-if="!tgStatus?.bound" class="unbound-section">
        <div class="section-header">
          <h2 class="section-title">📱 Telegram 绑定</h2>
          <p class="section-description">
            绑定 Telegram 后，你可以通过 Telegram 机器人接收任务提醒，完成或删除任务。
          </p>
        </div>

        <div class="action-area">
          <n-button
            type="primary"
            size="large"
            :loading="loadingCode"
            @click="generateBindCode"
          >
            获取绑定码
          </n-button>
        </div>

        <div v-if="bindCode" class="bind-code-section">
          <div class="code-box">
            <div class="code-label">绑定码</div>
            <div class="code-display">{{ bindCode.code }}</div>
            <n-button text size="small" @click="copyCode" style="margin-top: 8px">
              📋 复制
            </n-button>
          </div>

          <div class="bot-link-section">
            <p class="section-subtitle">步骤：</p>
            <ol class="steps-list">
              <li>
                点击下方链接打开 Telegram Bot
                <n-button text type="primary" @click="openBotLink" class="bot-link-btn">
                  打开 @{{ bindCode.botUsername }}
                </n-button>
              </li>
              <li>发送以下命令：<code class="command-text">/start {{ bindCode.code }}</code></li>
              <li>绑定成功后刷新此页面或点击下方按钮检查状态</li>
            </ol>
          </div>

          <div v-if="codeExpiry > 0" class="expiry-section">
            <span class="expiry-text">绑定码有效期还剩</span>
            <span class="expiry-timer">{{ formatExpiry() }}</span>
          </div>

          <div class="action-area" style="margin-top: 16px">
            <n-button @click="checkStatus" :loading="loadingStatus">
              🔄 刷新状态
            </n-button>
          </div>
        </div>
      </div>

      <div v-else class="bound-section">
        <div class="success-section">
          <div class="success-icon">✅</div>
          <div class="success-text">
            <h3 class="success-title">已绑定</h3>
            <p class="success-desc">你的 Telegram 账号已成功绑定</p>
          </div>
        </div>

        <n-space vertical :size="16">
          <div class="bound-info">
            <span class="info-label">Telegram 用户名：</span>
            <span class="info-value">@{{ tgStatus?.tgUsername }}</span>
          </div>
          <div class="bound-info">
            <span class="info-label">Chat ID：</span>
            <span class="info-value">{{ tgStatus?.tgChatId }}</span>
          </div>
        </n-space>

        <n-popconfirm
          positive-text="确定解除"
          negative-text="取消"
          @positive-click="unbindTelegram"
        >
          <template #trigger>
            <n-button type="error" style="margin-top: 24px">
              🔗 解除绑定
            </n-button>
          </template>
          <p>确定要解除 Telegram 绑定吗？解除后将无法通过 Telegram 接收任务提醒。</p>
        </n-popconfirm>
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
  NPopconfirm
} from 'naive-ui'
import * as telegramApi from '@/api/telegram'
import type { BindCodeResponse, TgStatus } from '@/types'

const message = useMessage()

const tgStatus = ref<TgStatus | null>(null)
const bindCode = ref<BindCodeResponse | null>(null)
const loadingCode = ref(false)
const loadingStatus = ref(false)
const codeExpiry = ref(0)
let expiryInterval: ReturnType<typeof setInterval> | null = null

const formatExpiry = (): string => {
  const minutes = Math.floor(codeExpiry.value / 60)
  const seconds = codeExpiry.value % 60
  return `${minutes}分${seconds}秒`
}

const generateBindCode = async () => {
  try {
    loadingCode.value = true
    const response = await telegramApi.getBindCode()
    if (response.data) {
      bindCode.value = response.data
      codeExpiry.value = response.data.expiresIn

      if (expiryInterval) clearInterval(expiryInterval)
      expiryInterval = setInterval(() => {
        codeExpiry.value--
        if (codeExpiry.value <= 0) {
          if (expiryInterval) clearInterval(expiryInterval)
          message.warning('绑定码已过期，请重新获取')
          bindCode.value = null
        }
      }, 1000)

      message.success('绑定码已生成')
    }
  } catch (error: any) {
    const errorMsg = error?.response?.data?.message || error?.message || '获取绑定码失败'
    message.error(errorMsg)
  } finally {
    loadingCode.value = false
  }
}

const copyCode = () => {
  if (bindCode.value) {
    navigator.clipboard.writeText(bindCode.value.code).then(() => {
      message.success('绑定码已复制到剪贴板')
    })
  }
}

const openBotLink = () => {
  if (bindCode.value) {
    window.open(`https://t.me/${bindCode.value.botUsername}?start=${bindCode.value.code}`, '_blank')
  }
}

const checkStatus = async () => {
  try {
    loadingStatus.value = true
    const response = await telegramApi.getStatus()
    if (response.data) {
      tgStatus.value = response.data
      if (response.data.bound) {
        message.success('Telegram 已绑定')
        bindCode.value = null
        if (expiryInterval) clearInterval(expiryInterval)
      } else {
        message.info('Telegram 未绑定')
      }
    }
  } catch (error: any) {
    const errorMsg = error?.response?.data?.message || error?.message || '获取状态失败'
    message.error(errorMsg)
  } finally {
    loadingStatus.value = false
  }
}

const unbindTelegram = async () => {
  try {
    await telegramApi.unbind()
    message.success('已解除 Telegram 绑定')
    tgStatus.value = null
    bindCode.value = null
    if (expiryInterval) clearInterval(expiryInterval)
  } catch (error: any) {
    const errorMsg = error?.response?.data?.message || error?.message || '解除绑定失败'
    message.error(errorMsg)
  }
}

onMounted(async () => {
  await checkStatus()
})
</script>

<style scoped>
.telegram-section {
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

.telegram-card {
  background: linear-gradient(135deg, rgba(30, 30, 50, 0.8) 0%, rgba(20, 20, 40, 0.8) 100%);
  border: 1px solid rgba(255, 215, 0, 0.2);
  border-radius: 8px;
  max-width: 600px;
}

.section-header {
  margin-bottom: 24px;
}

.section-title {
  font-size: 24px;
  font-weight: bold;
  background: linear-gradient(135deg, #ffd700, #ffed4e);
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

.action-area {
  display: flex;
  gap: 12px;
}

.bind-code-section {
  margin-top: 24px;
  padding: 20px;
  background: rgba(255, 255, 255, 0.02);
  border-radius: 8px;
  border: 1px solid rgba(255, 215, 0, 0.2);
}

.code-box {
  background: rgba(0, 0, 0, 0.3);
  padding: 16px;
  border-radius: 8px;
  text-align: center;
  margin-bottom: 16px;
  border: 1px solid rgba(255, 215, 0, 0.2);
}

.code-label {
  font-size: 12px;
  color: #a0a0b0;
  text-transform: uppercase;
  letter-spacing: 1px;
  margin-bottom: 8px;
}

.code-display {
  font-size: 28px;
  font-weight: bold;
  color: #ffd700;
  letter-spacing: 4px;
  font-family: 'Courier New', monospace;
}

.bot-link-section {
  margin: 16px 0;
}

.steps-list {
  color: #d0d0e0;
  line-height: 1.8;
  padding-left: 20px;
  margin: 12px 0;
}

.steps-list li {
  margin-bottom: 12px;
}

.bot-link-btn {
  margin-left: 8px;
}

.command-text {
  background: rgba(0, 0, 0, 0.3);
  padding: 2px 8px;
  border-radius: 4px;
  font-family: 'Courier New', monospace;
  color: #ffd700;
  font-weight: 500;
}

.expiry-section {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px;
  background: rgba(239, 68, 68, 0.1);
  border-radius: 6px;
  margin: 12px 0;
}

.expiry-text {
  color: #fca5a5;
  font-size: 12px;
}

.expiry-timer {
  font-weight: bold;
  color: #fca5a5;
  font-family: 'Courier New', monospace;
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
  border: 1px solid rgba(255, 215, 0, 0.1);
}

.info-label {
  font-weight: 600;
  color: #d0d0e0;
}

.info-value {
  color: #ffd700;
  font-family: 'Courier New', monospace;
  font-weight: 500;
}

:deep(.n-button--primary) {
  background: linear-gradient(135deg, #ffd700, #ffed4e);
  color: #000 !important;
  border: none !important;
}

:deep(.n-button--primary:hover) {
  box-shadow: 0 4px 16px rgba(255, 215, 0, 0.4) !important;
}

@media (max-width: 768px) {
  .telegram-card {
    max-width: 100%;
  }

  .section-title {
    font-size: 20px;
  }

  .code-display {
    font-size: 24px;
  }

  .steps-list {
    font-size: 14px;
  }
}
</style>
