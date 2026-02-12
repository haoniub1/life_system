# 任务详情功能改进方案

## 需求
1. 任务列表中添加详情按钮（ℹ️图标）
2. 点击详情按钮显示任务完整信息（不跳转）
3. 修改任务表单中显示分类信息

## 修改方案

### 1. TaskManager.vue 修改

#### 1.1 添加详情按钮到操作菜单
在第 230 行附近，修改操作菜单：
```typescript
const menuOptions = [
  { label: 'ℹ️ 详情', key: 'detail' },  // 新增
  { label: '✏️ 编辑', key: 'edit' },
  { label: '🗑️ 删除', key: 'delete' }
]
```

#### 1.2 添加详情状态和抽屉组件
在 `<script setup>` 中添加：
```typescript
const showDetailDrawer = ref(false)
const detailTask = ref<Task | null>(null)
```

#### 1.3 处理详情点击事件
在 `handleMenuSelect` 函数中添加：
```typescript
function handleMenuSelect(key: string, task: Task) {
  if (key === 'detail') {
    detailTask.value = task
    showDetailDrawer.value = true
  } else if (key === 'edit') {
    // 现有代码...
```

#### 1.4 添加详情抽屉组件（在 `</template>` 前添加）
```vue
<!-- Task Detail Drawer -->
<n-drawer
  v-model:show="showDetailDrawer"
  :width="400"
  placement="right"
>
  <n-drawer-content title="任务详情" closable>
    <div v-if="detailTask" class="task-detail">
      <div class="detail-section">
        <h4>📝 基本信息</h4>
        <div class="detail-item">
          <span class="label">标题：</span>
          <span class="value">{{ detailTask.title }}</span>
        </div>
        <div v-if="detailTask.description" class="detail-item">
          <span class="label">描述：</span>
          <span class="value">{{ detailTask.description }}</span>
        </div>
        <div class="detail-item">
          <span class="label">类型：</span>
          <span class="value">
            {{ getTaskTypeIcon(detailTask.type) }}
            {{ detailTask.type === 'once' ? '一次性' : detailTask.type === 'repeatable' ? '重复任务' : '挑战任务' }}
          </span>
        </div>
        <div class="detail-item">
          <span class="label">分类：</span>
          <span class="value">{{ detailTask.category || '无' }}</span>
        </div>
        <div v-if="detailTask.primaryAttribute" class="detail-item">
          <span class="label">主属性：</span>
          <span class="value">
            {{ ATTR_DISPLAY[detailTask.primaryAttribute]?.emoji }}
            {{ ATTR_DISPLAY[detailTask.primaryAttribute]?.name }}
          </span>
        </div>
      </div>

      <div class="detail-section">
        <h4>💎 奖励</h4>
        <div v-if="detailTask.rewardSpiritStones" class="detail-item">
          <span class="label">灵石：</span>
          <span class="value">💎 {{ detailTask.rewardSpiritStones }}</span>
        </div>
        <div v-if="detailTask.rewardExp" class="detail-item">
          <span class="label">经验：</span>
          <span class="value">⭐ {{ detailTask.rewardExp }}</span>
        </div>
        <div v-if="detailTask.rewardPhysique" class="detail-item">
          <span class="label">体质：</span>
          <span class="value">💪 +{{ detailTask.rewardPhysique }}</span>
        </div>
        <div v-if="detailTask.rewardWillpower" class="detail-item">
          <span class="label">意志：</span>
          <span class="value">🧠 +{{ detailTask.rewardWillpower }}</span>
        </div>
        <div v-if="detailTask.rewardIntelligence" class="detail-item">
          <span class="label">智慧：</span>
          <span class="value">📚 +{{ detailTask.rewardIntelligence }}</span>
        </div>
        <div v-if="detailTask.rewardPerception" class="detail-item">
          <span class="label">悟性：</span>
          <span class="value">👁 +{{ detailTask.rewardPerception }}</span>
        </div>
        <div v-if="detailTask.rewardCharisma" class="detail-item">
          <span class="label">魅力：</span>
          <span class="value">✨ +{{ detailTask.rewardCharisma }}</span>
        </div>
        <div v-if="detailTask.rewardAgility" class="detail-item">
          <span class="label">敏捷：</span>
          <span class="value">🏃 +{{ detailTask.rewardAgility }}</span>
        </div>
      </div>

      <div class="detail-section">
        <h4>⚡ 消耗 & 限制</h4>
        <div class="detail-item">
          <span class="label">疲劳消耗：</span>
          <span class="value">⚡ {{ detailTask.fatigueCost }}</span>
        </div>
        <div class="detail-item">
          <span class="label">难度：</span>
          <span class="value">
            <span v-for="i in detailTask.difficulty" :key="i">⭐</span>
            ({{ detailTask.difficulty }}星)
          </span>
        </div>
        <div v-if="detailTask.type === 'repeatable' && detailTask.dailyLimit" class="detail-item">
          <span class="label">每日限制：</span>
          <span class="value">{{ detailTask.todayCompletionCount }} / {{ detailTask.dailyLimit }}</span>
        </div>
        <div v-if="detailTask.deadline" class="detail-item">
          <span class="label">截止时间：</span>
          <span class="value">{{ formatDeadline(detailTask.deadline) }}</span>
        </div>
      </div>

      <div class="detail-section">
        <h4>📊 统计</h4>
        <div class="detail-item">
          <span class="label">完成次数：</span>
          <span class="value">{{ detailTask.completedCount }} 次</span>
        </div>
        <div class="detail-item">
          <span class="label">状态：</span>
          <span class="value" :class="detailTask.status">
            {{ detailTask.status === 'active' ? '进行中' : detailTask.status === 'completed' ? '已完成' : '已失败' }}
          </span>
        </div>
      </div>
    </div>
  </n-drawer-content>
</n-drawer>
```

#### 1.5 添加详情样式（在 `<style>` 中添加）
```css
.task-detail {
  padding: 16px 0;
}

.detail-section {
  margin-bottom: 24px;
}

.detail-section h4 {
  font-size: 14px;
  font-weight: 600;
  margin-bottom: 12px;
  color: var(--text-color-1);
  border-bottom: 1px solid var(--border-color);
  padding-bottom: 8px;
}

.detail-item {
  display: flex;
  justify-content: space-between;
  padding: 8px 0;
  border-bottom: 1px dashed var(--border-color);
}

.detail-item:last-child {
  border-bottom: none;
}

.detail-item .label {
  font-weight: 500;
  color: var(--text-color-2);
  min-width: 80px;
}

.detail-item .value {
  flex: 1;
  text-align: right;
  color: var(--text-color-1);
}

.detail-item .value.active {
  color: #18a058;
}

.detail-item .value.completed {
  color: #909399;
}

.detail-item .value.failed {
  color: #d03050;
}
```

### 2. 快速实现建议

由于文件较大，建议使用以下命令快速修改：

1. 添加详情按钮到操作菜单
2. 添加详情抽屉组件
3. 添加样式

完整的修改脚本见下方。
