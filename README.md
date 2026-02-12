# 人生修炼系统 - Life System V4.1

一个以修仙境界为核心的游戏化人生管理系统。完成现实任务获得灵石、提升属性、突破境界。

**不是记录人生，而是运行人生。**

## 核心系统

### 修仙境界体系

6 个独立属性，各自修炼进阶：

| 属性 | 说明 |
|------|------|
| 💪 体魄 | 运动、健康、饮食 |
| 🧠 意志 | 自律、习惯、冥想 |
| 📚 智力 | 学习、阅读、编程 |
| 👁 感知 | 观察、艺术、反思 |
| ✨ 魅力 | 沟通、社交 |
| 🏃 敏捷 | 执行力、协调 |

9 大境界 x 4 小境界（初期/中期/后期/大圆满）：

凡人 → 炼气 → 筑基 → 金丹 → 元婴 → 化神 → 合体 → 大乘 → 渡劫

### 灵石体系

1 下品灵石 = 1 RMB，分级显示：

🪨 下品 | 💎 中品(x100) | 💠 上品(x10000) | 🔮 极品(x1000000)

### 任务系统

- **一次性** / **可重复** / **挑战**（限时+惩罚）
- 0-5 星难度，自动匹配疲劳/灵石/属性加成
- 分类标签关联属性，1 星+ 强制选分类
- 拖拽排序

### 商店系统

- **消耗品**：使用后消耗，可带游戏效果
- **装备**：持久物品，可按售价回收灵石

### 通知

- **Telegram Bot**：任务提醒、截止通知
- **Bark 推送**：用户自行绑定 Key，无需服务端配置

## 快速开始

### 后端

```bash
cd backend
cp etc/config.example.yaml etc/config.yaml  # 首次需修改 Auth.Secret
make run
```

### 前端

```bash
cd frontend
npm install
npm run dev
```

### Docker

```bash
docker-compose up -d    # 后端 :8081, 前端 :8082
```

## 第三方 API

支持通过 API 快速完成任务，适合 iOS 快捷指令等自动化场景：

```bash
curl -X POST http://localhost:8081/api/tasks/quick \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"difficulty": 2, "categories": ["physique"]}'
```

完整 API 文档见 [doc/API.md](doc/API.md)

## 技术栈

**后端**: Go-Zero + SQLite + JWT + Telegram Bot API + Bark

**前端**: Vue 3 + TypeScript + Vite + Naive UI + Pinia

## License

MIT
