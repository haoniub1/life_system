# API 文档

Base URL: `http://localhost:8081`

## 通用说明

### 响应格式

所有接口统一返回 `CommonResp`：

```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

| code | 含义 |
|------|------|
| 0 | 成功 |
| 400 | 参数错误 |
| 401 | 未授权 |
| 429 | 请求过于频繁（限流） |
| 500 | 服务器错误 |

### 认证

除了注册和登录，所有接口需要在 Header 中携带 JWT Token：

```
Authorization: Bearer <token>
```

### 限流

- 同一 IP 每日登录失败上限：10 次（可配置）
- 同一 IP 每日注册上限：10 次（可配置）
- 超出后返回 429，次日重置

---

## 认证

### 注册

```
POST /api/auth/register
```

**请求：**

```json
{
  "username": "testuser",
  "password": "123456"
}
```

**响应 data：**

```json
{
  "token": "eyJhbGciOi...",
  "user": {
    "id": 1,
    "username": "testuser",
    "displayName": "",
    "avatar": "",
    "tgChatId": 0,
    "tgUsername": ""
  }
}
```

### 登录

```
POST /api/auth/login
```

**请求：** 同注册

**响应 data：** 同注册

失败时返回剩余尝试次数：`"用户名或密码错误（今日剩余尝试次数：8）"`

### 登出

```
POST /api/auth/logout
```

### 获取当前用户

```
GET /api/auth/me
```

**响应 data：** `UserInfo` 对象

---

## 用户

### 更新资料

```
PUT /api/user/profile
```

```json
{
  "displayName": "道友",
  "avatar": "/uploads/avatar.jpg"
}
```

### 修改密码

```
PUT /api/user/password
```

```json
{
  "oldPassword": "123456",
  "newPassword": "654321"
}
```

### 上传文件

```
POST /api/upload
Content-Type: multipart/form-data
```

字段名：`file`

**响应 data：**

```json
{
  "url": "/uploads/abc123.jpg"
}
```

---

## 角色

### 获取角色信息

```
GET /api/character
```

**响应 data：**

```json
{
  "userId": 1,
  "spiritStones": 1250,
  "fatigue": 30,
  "fatigueCap": 100,
  "fatigueLevel": 0,
  "overdraftPenalty": 0,
  "title": "炼气初期",
  "lastActivityDate": "2026-02-12",
  "attributes": [
    {
      "attrKey": "physique",
      "displayName": "体魄",
      "emoji": "💪",
      "value": 105.5,
      "realm": 1,
      "realmName": "炼气",
      "subRealm": 0,
      "subRealmName": "初期",
      "realmExp": 3,
      "isBottleneck": false,
      "accumulationPool": 0,
      "attrCap": 200,
      "progressPercent": 5.5,
      "color": "#10b981"
    }
  ]
}
```

**属性 key 列表：**

| key | 名称 | emoji |
|-----|------|-------|
| `physique` | 体魄 | 💪 |
| `willpower` | 意志 | 🧠 |
| `intelligence` | 智力 | 📚 |
| `perception` | 感知 | 👁 |
| `charisma` | 魅力 | ✨ |
| `agility` | 敏捷 | 🏃 |
| `luck` | 幸运 | 🍀 |

---

## 任务

### 获取任务列表

```
GET /api/tasks?type=once&status=active
```

| 参数 | 可选值 | 说明 |
|------|--------|------|
| `type` | `once`, `repeatable`, `challenge` | 不传返回全部 |
| `status` | `active`, `completed`, `failed` | 不传返回非 deleted |

**响应 data：**

```json
{
  "tasks": [TaskResp, ...]
}
```

### 创建任务

```
POST /api/tasks
```

```json
{
  "title": "晨跑30分钟",
  "description": "",
  "category": "",
  "type": "once",
  "difficulty": 2,
  "rewardSpiritStones": 120,
  "rewardPhysique": 0.2,
  "fatigueCost": 10,
  "deadline": "",
  "dailyLimit": 0,
  "totalLimit": 0,
  "remindBefore": 0,
  "remindInterval": 0
}
```

### 更新任务

```
PUT /api/tasks/:id
```

只传需要修改的字段（partial update）。

### 完成任务

```
POST /api/tasks/complete/:id
```

**响应 data：**

```json
{
  "task": { TaskResp },
  "character": { CharacterResp },
  "message": "✅ 任务「晨跑30分钟」已完成！获得 120灵石"
}
```

### 删除任务

```
DELETE /api/tasks/:id
```

### 快速任务（第三方 API 推荐）

```
POST /api/tasks/quick
```

只需传难度和分类，自动按模板填充疲劳/灵石/属性加成。适合 iOS 快捷指令、自动化工具等第三方调用。

```json
{
  "difficulty": 2,
  "categories": ["physique", "intelligence"],
  "title": "晨跑+读书",
  "type": "once",
  "source": "ios-shortcut",
  "dailyLimit": 0,
  "totalLimit": 0,
  "deadline": ""
}
```

| 字段 | 必填 | 说明 |
|------|------|------|
| `difficulty` | 是 | 0-5 星 |
| `categories` | 否 | 属性 key 数组，见下表 |
| `title` | 否 | 不传自动生成 `快速任务 (★2)` |
| `type` | 否 | `once`（默认，立即完成）/ `repeatable` / `challenge` |
| `source` | 否 | 来源标识，默认 `"api"` |
| `dailyLimit` | 否 | repeatable 每日上限，0=不限 |
| `totalLimit` | 否 | repeatable 总上限，0=不限 |
| `deadline` | challenge 必填 | ISO8601 格式，如 `2026-02-15T23:59:59+08:00` |

**难度模板（自动填充）：**

| 星 | 疲劳 | 灵石 | 属性加成 |
|----|------|------|---------|
| 0 | 1 | 10 | 0 |
| 1 | 5 | 50 | 0.1 |
| 2 | 10 | 120 | 0.2 |
| 3 | 20 | 300 | 0.4 |
| 4 | 40 | 800 | 0.7 |
| 5 | 90 | 2500 | 1.0 |

**分类（categories）：**

| key | 名称 | 说明 |
|-----|------|------|
| `physique` | 💪 体魄 | 运动、健康、饮食 |
| `willpower` | 🧠 意志 | 自律、习惯、冥想 |
| `intelligence` | 📚 智力 | 学习、阅读、编程 |
| `perception` | 👁 感知 | 观察、艺术、反思 |
| `charisma` | ✨ 魅力 | 沟通、社交 |
| `agility` | 🏃 敏捷 | 执行力、协调 |

**行为差异：**

- `once`：创建 + 立即完成，返回奖励，`completed: true`
- `repeatable`：仅创建，之后通过 `POST /api/tasks/complete/:id` 反复完成，`completed: false`
- `challenge`：仅创建，有截止时间，过期未完成会扣罚，`completed: false`

**响应 data：**

```json
{
  "task": { TaskResp },
  "character": { CharacterResp },
  "message": "✅ 任务「晨跑+读书」已完成！获得 120灵石",
  "completed": true
}
```

### 任务排序

```
PUT /api/tasks/reorder
```

```json
{
  "taskIds": [3, 1, 5, 2]
}
```

传入所有任务 ID 的有序数组，按数组顺序设置 `sortOrder`。

---

## 商店

### 获取商品列表

```
GET /api/shop/items
```

**响应 data：**

```json
{
  "items": [
    {
      "id": 1,
      "name": "回复丹",
      "description": "恢复 20 点疲劳",
      "price": 100,
      "sellPrice": 0,
      "itemType": "consumable",
      "icon": "💊",
      "image": "",
      "stock": -1
    }
  ]
}
```

`itemType`：`consumable`（消耗品）/ `equipment`（装备）

`stock`：`-1` 表示无限库存

### 创建商品

```
POST /api/shop/items
```

```json
{
  "name": "灵剑",
  "description": "一把好剑",
  "price": 500,
  "sellPrice": 250,
  "itemType": "equipment",
  "icon": "⚔️",
  "image": "",
  "stock": -1
}
```

### 更新商品

```
PUT /api/shop/items/:id
```

Partial update，只传需要修改的字段。

### 删除商品

```
DELETE /api/shop/items/:id
```

### 购买商品

```
POST /api/shop/purchase
```

```json
{
  "itemId": 1,
  "quantity": 1
}
```

**响应 data：**

```json
{
  "success": true,
  "message": "购买成功",
  "remainingSpiritStones": 900
}
```

### 获取背包

```
GET /api/shop/inventory
```

**响应 data：**

```json
{
  "items": [
    {
      "id": 1,
      "itemId": 3,
      "name": "回复丹",
      "description": "恢复 20 点疲劳",
      "itemType": "consumable",
      "sellPrice": 0,
      "icon": "💊",
      "image": "",
      "quantity": 3
    }
  ]
}
```

### 使用消耗品

```
POST /api/shop/use
```

```json
{
  "itemId": 1,
  "quantity": 1
}
```

**响应 data：**

```json
{
  "success": true,
  "message": "使用成功",
  "character": { CharacterResp }
}
```

### 出售装备

```
POST /api/shop/sell
```

```json
{
  "itemId": 3,
  "quantity": 1
}
```

**响应 data：**

```json
{
  "success": true,
  "message": "出售成功，获得 250 灵石",
  "remainingSpiritStones": 1150
}
```

### 购买记录

```
GET /api/shop/history
```

**响应 data：**

```json
{
  "history": [
    {
      "id": 1,
      "itemName": "回复丹",
      "quantity": 1,
      "totalPrice": 100,
      "createdAt": "2026-02-12T10:00:00Z"
    }
  ]
}
```

---

## 动态

### 获取时间线

```
GET /api/timeline
```

**响应 data：**

```json
{
  "events": [
    {
      "id": "task_1",
      "type": "task_complete",
      "title": "完成任务「晨跑」",
      "description": "",
      "rewards": {
        "spiritStones": 120
      },
      "timestamp": "2026-02-12T08:00:00Z"
    }
  ],
  "tasksCompleted": 5,
  "totalExp": 0,
  "totalSpiritStones": 600,
  "sleepRecords": 1
}
```

`type` 可选值：`task_complete`, `task_fail`, `task_delete`, `sleep`, `purchase`

---

## Telegram

### 生成绑定码

```
POST /api/telegram/bindcode
```

**响应 data：**

```json
{
  "code": "ABC123",
  "botUsername": "life_system_bot",
  "expiresIn": 300
}
```

用户在 Telegram 向 Bot 发送绑定码即可绑定。

### 获取绑定状态

```
GET /api/telegram/status
```

**响应 data：**

```json
{
  "bound": true,
  "tgUsername": "myuser",
  "tgChatId": 123456789
}
```

### 解绑

```
DELETE /api/telegram/unbind
```

---

## Bark 推送

### 设置 Bark Key

```
PUT /api/bark/key
```

```json
{
  "barkKey": "your-bark-device-key"
}
```

Bark Key 从 Bark App 中获取，系统直接使用官方服务器 `https://api.day.app` 推送。

### 获取 Bark 状态

```
GET /api/bark/status
```

**响应 data：**

```json
{
  "enabled": true,
  "barkKey": "abcdefgh***"
}
```

Key 会脱敏显示（前 8 位 + ***）。

### 测试推送

```
POST /api/bark/test
```

```json
{
  "title": "测试",
  "body": "推送测试消息"
}
```

不传 title/body 会使用默认测试消息。

### 删除 Bark Key

```
DELETE /api/bark/key
```

---

## 调用示例

### cURL - 快速完成任务

```bash
# 登录获取 token
TOKEN=$(curl -s -X POST http://localhost:8081/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"test","password":"123456"}' | jq -r '.data.token')

# 快速完成一个 2 星任务（体魄+智力）
curl -X POST http://localhost:8081/api/tasks/quick \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "difficulty": 2,
    "categories": ["physique", "intelligence"],
    "title": "晨跑+刷题"
  }'
```

### iOS 快捷指令

1. 使用「获取 URL 内容」操作
2. URL: `http://your-server:8081/api/tasks/quick`
3. 方法: POST
4. Headers: `Authorization: Bearer <你的token>`
5. Body (JSON):
   ```json
   {
     "difficulty": 1,
     "categories": ["physique"],
     "source": "ios-shortcut"
   }
   ```
