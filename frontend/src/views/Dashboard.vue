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
        <div class="header-center" v-if="characterStore.character">
          <div class="global-spirit-stones" @click="showRMB = !showRMB">
            <template v-if="showRMB">
              <span class="global-icon">💰</span>
              <span class="global-value">¥{{ characterStore.character.spiritStones }}</span>
            </template>
            <template v-else>
              <span v-if="spiritDisplay.supreme > 0" class="stone-chip stone-supreme">🔮{{ spiritDisplay.supreme }}极品</span>
              <span v-if="spiritDisplay.high > 0" class="stone-chip stone-high">💠{{ spiritDisplay.high }}上品</span>
              <span v-if="spiritDisplay.medium > 0" class="stone-chip stone-medium">💎{{ spiritDisplay.medium }}中品</span>
              <span class="stone-chip stone-low">🪨{{ spiritDisplay.low }}下品</span>
            </template>
          </div>
          <div class="global-fatigue" @click="toggleFatigueMode">
            <span class="global-icon">{{ showActivity ? '⚡' : '😴' }}</span>
            <div class="global-bar-wrapper">
              <div class="global-bar">
                <div
                  class="global-bar-fill"
                  :style="{ width: fatigueBarPercent + '%', background: fatigueBarColor }"
                ></div>
              </div>
              <span class="global-bar-text">{{ fatigueBarLabel }}</span>
            </div>
          </div>
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

  <!-- Mobile Bottom Tab Bar -->
  <div class="mobile-tab-bar">
    <div
      v-for="tab in mobileTabs"
      :key="tab.key"
      class="tab-item"
      :class="{ active: activeMenu === tab.key }"
      @click="activeMenu = tab.key"
    >
      <span class="tab-icon">{{ tab.icon }}</span>
      <span class="tab-label">{{ tab.label }}</span>
    </div>
  </div>

  <!-- Settings Drawer -->
  <n-drawer v-model:show="showSettings" :width="drawerWidth" placement="right">
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
        <n-tab-pane name="bark" tab="Bark 推送">
          <bark-bind />
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
import { ref, h, computed, onMounted, onUnmounted } from 'vue'
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
import { decomposeSpiritStones } from '@/utils/rpg'
import CharacterCard from '@/components/CharacterCard.vue'
import TaskManager from '@/components/TaskManager.vue'
import TelegramBind from '@/components/TelegramBind.vue'
import BarkBind from '@/components/BarkBind.vue'
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
const showActivity = ref(true) // true=活跃度, false=疲劳值
const showRMB = ref(false) // true=RMB, false=灵石

const spiritDisplay = computed(() => {
  return decomposeSpiritStones(characterStore.character?.spiritStones || 0)
})
const windowWidth = ref(window.innerWidth)

const toggleFatigueMode = () => {
  showActivity.value = !showActivity.value
}

const fatigueBarPercent = computed(() => {
  const c = characterStore.character
  if (!c || c.fatigueCap === 0) return showActivity.value ? 100 : 0
  if (showActivity.value) {
    return Math.max(0, Math.min(100, (1 - c.fatigue / c.fatigueCap) * 100))
  }
  return Math.max(0, Math.min(100, (c.fatigue / c.fatigueCap) * 100))
})

const fatigueBarColor = computed(() => {
  if (showActivity.value) {
    const pct = fatigueBarPercent.value
    if (pct > 60) return '#10b981'
    if (pct > 30) return '#f59e0b'
    return '#ef4444'
  }
  const pct = fatigueBarPercent.value
  if (pct < 40) return '#10b981'
  if (pct < 70) return '#f59e0b'
  return '#ef4444'
})

const fatigueBarLabel = computed(() => {
  const c = characterStore.character
  if (!c) return ''
  if (showActivity.value) {
    const pct = Math.max(0, Math.min(100, (1 - c.fatigue / c.fatigueCap) * 100))
    return `${pct.toFixed(0)}%`
  }
  return `${c.fatigue}/${c.fatigueCap}`
})

const onResize = () => { windowWidth.value = window.innerWidth }
const isMobile = computed(() => windowWidth.value <= 768)
const drawerWidth = computed(() => isMobile.value ? '100%' : 420)

const mobileTabs = [
  { key: 'character', icon: '🎮', label: '角色' },
  { key: 'task', icon: '✅', label: '任务' },
  { key: 'shop', icon: '🛒', label: '商店' },
  { key: 'timeline', icon: '📅', label: '时间线' }
]

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
  window.addEventListener('resize', onResize)
  try {
    await characterStore.fetchCharacter()
  } catch (error) {
    console.error('Failed to fetch character:', error)
  }
})

onUnmounted(() => {
  window.removeEventListener('resize', onResize)
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

/* Global Stats in Header */
.header-center {
  display: flex;
  align-items: center;
  gap: 16px;
}

.global-spirit-stones {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 4px 10px;
  background: rgba(139, 92, 246, 0.12);
  border: 1px solid rgba(139, 92, 246, 0.25);
  border-radius: 16px;
  cursor: pointer;
  transition: background 0.2s;
  user-select: none;
}

.global-spirit-stones:hover {
  background: rgba(139, 92, 246, 0.2);
}

.stone-chip {
  font-size: 12px;
  font-weight: 600;
  padding: 1px 5px;
  border-radius: 4px;
}

.stone-supreme {
  color: #ffd700;
  background: rgba(255, 215, 0, 0.15);
}

.stone-high {
  color: #c084fc;
  background: rgba(168, 85, 247, 0.15);
}

.stone-medium {
  color: #60a5fa;
  background: rgba(59, 130, 246, 0.15);
}

.stone-low {
  color: #9ca3af;
  background: rgba(156, 163, 175, 0.1);
}

.global-icon {
  font-size: 14px;
  line-height: 1;
}

.global-value {
  font-size: 13px;
  font-weight: 700;
  color: #c4b5fd;
}

.global-fatigue {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 4px 10px;
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 16px;
  cursor: pointer;
  transition: background 0.2s;
  user-select: none;
}

.global-fatigue:hover {
  background: rgba(255, 255, 255, 0.08);
}

.global-bar-wrapper {
  display: flex;
  align-items: center;
  gap: 6px;
}

.global-bar {
  width: 60px;
  height: 6px;
  background: rgba(255, 255, 255, 0.08);
  border-radius: 3px;
  overflow: hidden;
}

.global-bar-fill {
  height: 100%;
  border-radius: 3px;
  transition: width 0.3s ease, background 0.3s ease;
}

.global-bar-text {
  font-size: 11px;
  color: #a0a0b0;
  min-width: 36px;
  white-space: nowrap;
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

/* Mobile Bottom Tab Bar */
.mobile-tab-bar {
  display: none;
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  height: 56px;
  background: rgba(20, 20, 35, 0.98);
  border-top: 1px solid rgba(255, 215, 0, 0.2);
  z-index: 100;
  grid-template-columns: repeat(4, 1fr);
  align-items: center;
}

.tab-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 2px;
  padding: 6px 0;
  color: #808090;
  cursor: pointer;
  transition: color 0.2s;
  -webkit-tap-highlight-color: transparent;
}

.tab-item.active {
  color: #ffd700;
}

.tab-icon {
  font-size: 20px;
  line-height: 1;
}

.tab-label {
  font-size: 11px;
  font-weight: 500;
}

@media (max-width: 768px) {
  .mobile-tab-bar {
    display: grid;
  }

  .dashboard-header {
    padding: 8px 12px;
    height: 48px;
  }

  .page-title {
    font-size: 16px;
  }

  .header-center {
    gap: 6px;
  }

  .stone-chip {
    font-size: 10px;
    padding: 0 3px;
  }

  .global-bar {
    width: 40px;
  }

  .username {
    display: none;
  }

  .settings-arrow {
    display: none;
  }

  .header-avatar {
    width: 30px;
    height: 30px;
  }

  .sidebar {
    display: none;
  }

  .dashboard-content {
    padding: 12px;
    padding-bottom: 72px;
  }
}
</style>
