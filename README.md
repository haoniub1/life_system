# 人生RPG系统 - Life System

一个将人生gamify化的任务管理系统，通过RPG游戏机制激励自我提升。

## 🎮 核心功能

### ✅ 任务系统
- **三种任务类型**：一次性、可重复、挑战
- **奖励惩罚机制**：经验、金币、属性提升/扣除
- **能量系统**：完成任务消耗能量，睡眠恢复

### 🎯 角色成长
- **等级系统**：经验升级
- **四维属性**：力量、智力、体力、精神
- **用进废退**：不活跃会导致属性衰减（1%/天）

### 😴 睡眠追踪
- 根据时长和质量恢复能量
- ≥8小时: 100% | 6-8小时: 80% | 4-6小时: 50%

### 🛒 奖励商店
- 消耗品：生命药水、能量饮料
- 永久道具：属性精华

### 📱 Telegram集成
- 任务提醒 | 过期通知 | 衰减警告

## 🚀 快速开始

### 后端
```bash
cd backend
make build
./bin/life-system-backend -f etc/config.yaml
```

### 前端
```bash
cd frontend
npm install
npm run dev
```

## ⚙️ 配置

编辑 `backend/etc/config.yaml`:
```yaml
Name: life-system-backend
Host: 0.0.0.0
Port: 8080

Database:
  Path: ./data/life-system.db

Auth:
  Secret: your-secret

Telegram:
  Enabled: true
  BotToken: your-bot-token
```

## 📊 已完成功能

✅ 用户认证 | ✅ 角色系统 | ✅ 任务管理 | ✅ 睡眠记录
✅ 商店系统 | ✅ Telegram Bot | ✅ 属性衰减 | ✅ 每日重置

## 🏗️ 技术栈

**后端**: Go-Zero + SQLite + Telegram Bot API
**前端**: Vue 3 + TypeScript + Naive UI + Pinia

## 📝 License

MIT

---

**开始你的人生RPG之旅吧！** 🎮✨
