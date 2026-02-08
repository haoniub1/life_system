<template>
  <div class="profile-section">
    <!-- Avatar -->
    <div class="avatar-section">
      <div class="avatar-wrapper" @click="triggerAvatarUpload">
        <img v-if="avatarUrl" :src="avatarUrl" class="avatar-img" />
        <div v-else class="avatar-placeholder">
          {{ (userStore.user?.displayName || userStore.user?.username || '?')[0] }}
        </div>
        <div class="avatar-overlay">
          <span>更换头像</span>
        </div>
      </div>
      <input
        ref="avatarInput"
        type="file"
        accept="image/*"
        style="display: none"
        @change="handleAvatarChange"
      />
      <n-spin v-if="uploadingAvatar" size="small" style="margin-top: 8px" />
    </div>

    <!-- Display Name -->
    <n-form-item label="昵称" style="margin-top: 16px">
      <n-input
        v-model:value="displayName"
        placeholder="输入昵称"
        clearable
      />
    </n-form-item>

    <n-button
      type="primary"
      block
      :loading="savingProfile"
      @click="handleSaveProfile"
    >
      保存资料
    </n-button>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useMessage } from 'naive-ui'
import { NFormItem, NInput, NButton, NSpin } from 'naive-ui'
import { useUserStore } from '@/stores/user'
import { updateProfile } from '@/api/user'
import { uploadFile } from '@/api/shop'

const emit = defineEmits<{
  saved: []
}>()

const message = useMessage()
const userStore = useUserStore()

const displayName = ref('')
const avatarUrl = ref('')
const savingProfile = ref(false)
const uploadingAvatar = ref(false)
const avatarInput = ref<HTMLInputElement | null>(null)

onMounted(() => {
  if (userStore.user) {
    displayName.value = userStore.user.displayName || ''
    avatarUrl.value = userStore.user.avatar || ''
  }
})

const triggerAvatarUpload = () => {
  avatarInput.value?.click()
}

const handleAvatarChange = async (e: Event) => {
  const target = e.target as HTMLInputElement
  const file = target.files?.[0]
  if (!file) return

  if (file.size > 5 * 1024 * 1024) {
    message.error('图片大小不能超过5MB')
    return
  }

  try {
    uploadingAvatar.value = true
    const response = await uploadFile(file) as any
    if (response.data?.url) {
      avatarUrl.value = response.data.url
      message.success('头像上传成功')
    }
  } catch (error: any) {
    message.error(error?.message || '上传失败')
  } finally {
    uploadingAvatar.value = false
    target.value = ''
  }
}

const handleSaveProfile = async () => {
  try {
    savingProfile.value = true
    const response = await updateProfile({
      displayName: displayName.value,
      avatar: avatarUrl.value
    }) as any

    if (response.data) {
      userStore.user = response.data
      message.success('资料已更新')
      emit('saved')
    }
  } catch (error: any) {
    message.error(error?.message || '保存失败')
  } finally {
    savingProfile.value = false
  }
}
</script>

<style scoped>
.profile-section {
  padding: 8px 0;
}

.avatar-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-bottom: 8px;
}

.avatar-wrapper {
  width: 90px;
  height: 90px;
  border-radius: 50%;
  overflow: hidden;
  cursor: pointer;
  position: relative;
  border: 3px solid rgba(255, 215, 0, 0.4);
  transition: all 0.3s ease;
}

.avatar-wrapper:hover {
  border-color: #ffd700;
  box-shadow: 0 0 20px rgba(255, 215, 0, 0.3);
}

.avatar-wrapper:hover .avatar-overlay {
  opacity: 1;
}

.avatar-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.avatar-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #ffd700, #ffed4e);
  color: #000;
  font-size: 32px;
  font-weight: bold;
}

.avatar-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.6);
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-size: 12px;
  opacity: 0;
  transition: opacity 0.3s ease;
}

:deep(.n-button--primary) {
  background: linear-gradient(135deg, #ffd700, #ffed4e) !important;
  color: #000 !important;
  border: none !important;
}
</style>
