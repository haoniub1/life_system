<template>
  <n-layout has-sider class="dashboard-layout">
    <n-layout-sider
      class="sidebar"
      collapse-mode="width"
      :collapsed-width="64"
      :width="240"
      :collapsed="collapsed"
      show-trigger="bar"
      @collapse="collapsed = true"
      @expand="collapsed = false"
    >
      <div class="sidebar-header">
        <h2 class="sidebar-title" v-if="!collapsed">人生RPG</h2>
      </div>

      <n-menu
        :value="activeMenu"
        :options="menuOptions"
        @update:value="handleMenuChange"
      />
    </n-layout-sider>

    <n-layout class="main-layout">
      <n-layout-header class="dashboard-header">
        <div class="header-left">
          <h1 class="page-title">{{ getPageTitle() }}</h1>
        </div>
        <div class="header-right">
          <div class="user-info" @click="showSettings = true">
            <div class="header-avatar">
              <img v-if="userStore.user?.avatar" :src="userStore.user.avatar" class="header-avatar-img" />
              <span v-else class="header-avatar-text">
                {{ (userStore.user?.displayName || userStore.user?.username || '?')[0] }}
              </span>
            </div>
            <span class="username">{{ userStore.user?.displayName || userStore.user?.username }}</span>
            <span class="settings-arrow">▾</span>
          </div>
        </div>
      </n-layout-header>

      <n-layout-content class="dashboard-content">
        <div class="content-wrapper">
          <character-card v-if="activeMenu === 'character'" />
          <task-manager v-else-if="activeMenu === 'task'" />
          <shop v-else-if="activeMenu === 'shop'" />
          <activity-timeline v-else-if="activeMenu === 'timeline'" />
        </div>
      </n-layout-content>
    </n-layout>
  </n-layout>

  <!-- Settings Drawer -->
  <n-drawer v-model:show="showSettings" :width="420" placement="right">
    <n-drawer-content closable>
      <template #header>
        <span class="drawer-title">设置</span>
      </template>

      <n-tabs type="line" animated>
        <n-tab-pane name="profile" tab="个人资料">
          <user-profile @saved="handleProfileSaved" />
        </n-tab-pane>
        <n-tab-pane name="telegram" tab="Telegram">
          <telegram-bind />
        </n-tab-pane>
        <n-tab-pane name="password" tab="修改密码">
          <password-form />
        </n-tab-pane>
      </n-tabs>

      <template #footer>
        <n-button type="error" block @click="handleLogout">
          退出登录
        </n-button>
      </template>
    </n-drawer-content>
  </n-drawer>
</template>

<script setup lang="ts">
import { ref, h, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useMessage } from 'naive-ui'
import {
  NLayout,
  NLayoutSider,
  NLayoutHeader,
  NLayoutContent,
  NMenu,
  NButton,
  NDrawer,
  NDrawerContent,
  NTabs,
  NTabPane,
  type MenuOption
} from 'naive-ui'
import {
  GameController,
  CheckmarkDone,
  Cart,
  Time
} from '@vicons/ionicons5'
import { useUserStore } from '@/stores/user'
import { useCharacterStore } from '@/stores/character'
import CharacterCard from '@/components/CharacterCard.vue'
import TaskManager from '@/components/TaskManager.vue'
import TelegramBind from '@/components/TelegramBind.vue'
import Shop from '@/components/Shop.vue'
import ActivityTimeline from '@/components/ActivityTimeline.vue'
import UserProfile from '@/components/UserProfile.vue'
import PasswordForm from '@/components/PasswordForm.vue'

const router = useRouter()
const message = useMessage()
const userStore = useUserStore()
const characterStore = useCharacterStore()

const collapsed = ref(false)
const activeMenu = ref('character')
const showSettings = ref(false)

const menuOptions: MenuOption[] = [
  {
    label: '角色总览',
    key: 'character',
    icon: () => h(GameController)
  },
  {
    label: '任务管理',
    key: 'task',
    icon: () => h(CheckmarkDone)
  },
  {
    label: '奖励商店',
    key: 'shop',
    icon: () => h(Cart)
  },
  {
    label: '活动时间线',
    key: 'timeline',
    icon: () => h(Time)
  }
]

const getPageTitle = (): string => {
  const titles: { [key: string]: string } = {
    character: '角色总览',
    task: '任务管理',
    shop: '奖励商店',
    timeline: '活动时间线'
  }
  return titles[activeMenu.value] || '人生RPG'
}

const handleMenuChange = (key: string | number) => {
  activeMenu.value = key as string
}

const handleProfileSaved = () => {
  // Profile updated, user store already refreshed inside UserProfile
}

const handleLogout = async () => {
  try {
    await userStore.logout()
    message.success('已退出登录')
    showSettings.value = false
    await router.push('/login')
  } catch (error: any) {
    const errorMsg = error?.response?.data?.message || error?.message || '退出失败'
    message.error(errorMsg)
  }
}

onMounted(async () => {
  try {
    await characterStore.fetchCharacter()
  } catch (error) {
    console.error('Failed to fetch character:', error)
  }
})
</script>

<style scoped>
.dashboard-layout {
  width: 100%;
  height: 100vh;
  background-color: #1a1a1a;
}

.sidebar {
  background: linear-gradient(180deg, rgba(20, 20, 35, 0.95) 0%, rgba(15, 15, 25, 0.95) 100%);
  border-right: 1px solid rgba(255, 215, 0, 0.2);
}

.sidebar-header {
  padding: 16px;
  border-bottom: 1px solid rgba(255, 215, 0, 0.2);
}

.sidebar-title {
  font-size: 18px;
  font-weight: bold;
  background: linear-gradient(135deg, #ffd700, #ffed4e);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin: 0;
  padding: 8px 0;
}

:deep(.n-menu) {
  background: transparent !important;
}

:deep(.n-menu-item) {
  color: #d0d0e0 !important;
}

:deep(.n-menu-item--selected) {
  background: rgba(255, 215, 0, 0.15) !important;
  color: #ffd700 !important;
}

:deep(.n-menu-item:hover) {
  background: rgba(255, 215, 0, 0.1) !important;
}

.main-layout {
  display: flex;
  flex-direction: column;
}

.dashboard-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 24px;
  background: linear-gradient(90deg, rgba(20, 20, 35, 0.8) 0%, rgba(30, 30, 50, 0.8) 100%);
  border-bottom: 1px solid rgba(255, 215, 0, 0.2);
  height: 70px;
}

.header-left {
  flex: 1;
}

.page-title {
  font-size: 24px;
  font-weight: bold;
  background: linear-gradient(135deg, #ffd700, #ffed4e);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin: 0;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 10px;
  cursor: pointer;
  padding: 6px 12px;
  border-radius: 8px;
  transition: all 0.2s ease;
}

.user-info:hover {
  background: rgba(255, 215, 0, 0.1);
}

.header-avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  overflow: hidden;
  border: 2px solid rgba(255, 215, 0, 0.4);
  flex-shrink: 0;
}

.header-avatar-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.header-avatar-text {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #ffd700, #ffed4e);
  color: #000;
  font-size: 16px;
  font-weight: bold;
}

.username {
  color: #d0d0e0;
  font-weight: 500;
  white-space: nowrap;
}

.settings-arrow {
  color: #808090;
  font-size: 12px;
}

.drawer-title {
  font-size: 18px;
  font-weight: bold;
  color: #ffd700;
}

.dashboard-content {
  flex: 1;
  overflow: auto;
  padding: 24px;
}

.content-wrapper {
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

@media (max-width: 768px) {
  .dashboard-header {
    padding: 12px 16px;
    height: auto;
  }

  .page-title {
    font-size: 18px;
  }

  .sidebar {
    display: none;
  }
}
</style>
