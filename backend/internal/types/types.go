package types

// Auth
type RegisterReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserInfo struct {
	ID          int64  `json:"id"`
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
	Avatar      string `json:"avatar"`
	TgChatID    int64  `json:"tgChatId"`
	TgUsername  string `json:"tgUsername"`
}

type UpdateProfileReq struct {
	DisplayName string `json:"displayName"`
	Avatar      string `json:"avatar"`
}

type ChangePasswordReq struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

type AuthResp struct {
	Token string   `json:"token"`
	User  UserInfo `json:"user"`
}

// Character
type CharacterResp struct {
	UserID           int64           `json:"userId"`
	SpiritStones     int             `json:"spiritStones"`
	Fatigue          int             `json:"fatigue"`
	FatigueCap       int             `json:"fatigueCap"`
	FatigueLevel     int             `json:"fatigueLevel"`
	OverdraftPenalty float64         `json:"overdraftPenalty"`
	Title            string          `json:"title"`
	LastActivityDate string          `json:"lastActivityDate"`
	Attributes       []AttributeResp `json:"attributes"`
}

type AttributeResp struct {
	AttrKey          string  `json:"attrKey"`
	DisplayName      string  `json:"displayName"`
	Emoji            string  `json:"emoji"`
	Value            float64 `json:"value"`
	TodayGain        float64 `json:"todayGain"`
	Realm            int     `json:"realm"`
	RealmName        string  `json:"realmName"`
	SubRealm         int     `json:"subRealm"`
	SubRealmName     string  `json:"subRealmName"`
	RealmExp         int     `json:"realmExp"`
	IsBottleneck     bool    `json:"isBottleneck"`
	AccumulationPool float64 `json:"accumulationPool"`
	AttrCap          float64 `json:"attrCap"`
	ProgressPercent  float64 `json:"progressPercent"`
	Color            string  `json:"color"`
}

type SpiritStoneDisplay struct {
	Total   int `json:"total"`
	Supreme int `json:"supreme"` // floor(n/1000000)
	High    int `json:"high"`    // floor(n%1000000/10000)
	Medium  int `json:"medium"`  // floor(n%10000/100)
	Low     int `json:"low"`     // n%100
}

// Task
type CreateTaskReq struct {
	Title              string  `json:"title"`
	Description        string  `json:"description"`
	Category           string  `json:"category"`
	Type               string  `json:"type"`
	Deadline           string  `json:"deadline"`
	PrimaryAttribute   string  `json:"primaryAttribute"`
	Difficulty         int     `json:"difficulty"`
	RewardExp          int     `json:"rewardExp"`
	RewardSpiritStones int     `json:"rewardSpiritStones"`
	RewardPhysique     float64 `json:"rewardPhysique"`
	RewardWillpower    float64 `json:"rewardWillpower"`
	RewardIntelligence float64 `json:"rewardIntelligence"`
	RewardPerception   float64 `json:"rewardPerception"`
	RewardCharisma     float64 `json:"rewardCharisma"`
	RewardAgility      float64 `json:"rewardAgility"`
	PenaltyExp         int     `json:"penaltyExp"`
	PenaltySpiritStones int    `json:"penaltySpiritStones"`
	FatigueCost        int     `json:"fatigueCost"`
	DailyLimit         int     `json:"dailyLimit"`
	TotalLimit         int     `json:"totalLimit"`
	RemindBefore       int     `json:"remindBefore"`
	RemindInterval     int     `json:"remindInterval"`
}

type UpdateTaskReq struct {
	Title              *string  `json:"title,omitempty"`
	Description        *string  `json:"description,omitempty"`
	Category           *string  `json:"category,omitempty"`
	Type               *string  `json:"type,omitempty"`
	Deadline           *string  `json:"deadline,omitempty"`
	PrimaryAttribute   *string  `json:"primaryAttribute,omitempty"`
	Difficulty         *int     `json:"difficulty,omitempty"`
	RewardExp          *int     `json:"rewardExp,omitempty"`
	RewardSpiritStones *int     `json:"rewardSpiritStones,omitempty"`
	RewardPhysique     *float64 `json:"rewardPhysique,omitempty"`
	RewardWillpower    *float64 `json:"rewardWillpower,omitempty"`
	RewardIntelligence *float64 `json:"rewardIntelligence,omitempty"`
	RewardPerception   *float64 `json:"rewardPerception,omitempty"`
	RewardCharisma     *float64 `json:"rewardCharisma,omitempty"`
	RewardAgility      *float64 `json:"rewardAgility,omitempty"`
	PenaltyExp         *int     `json:"penaltyExp,omitempty"`
	PenaltySpiritStones *int    `json:"penaltySpiritStones,omitempty"`
	FatigueCost        *int     `json:"fatigueCost,omitempty"`
	DailyLimit         *int     `json:"dailyLimit,omitempty"`
	TotalLimit         *int     `json:"totalLimit,omitempty"`
	RemindBefore       *int     `json:"remindBefore,omitempty"`
	RemindInterval     *int     `json:"remindInterval,omitempty"`
}

type TaskResp struct {
	ID                   int64   `json:"id"`
	UserID               int64   `json:"userId"`
	Title                string  `json:"title"`
	Description          string  `json:"description"`
	Category             string  `json:"category"`
	Type                 string  `json:"type"`
	Status               string  `json:"status"`
	Deadline             *string `json:"deadline"`
	PrimaryAttribute     string  `json:"primaryAttribute"`
	Difficulty           int     `json:"difficulty"`
	RewardExp            int     `json:"rewardExp"`
	RewardSpiritStones   int     `json:"rewardSpiritStones"`
	RewardPhysique       float64 `json:"rewardPhysique"`
	RewardWillpower      float64 `json:"rewardWillpower"`
	RewardIntelligence   float64 `json:"rewardIntelligence"`
	RewardPerception     float64 `json:"rewardPerception"`
	RewardCharisma       float64 `json:"rewardCharisma"`
	RewardAgility        float64 `json:"rewardAgility"`
	FatigueCost          int     `json:"fatigueCost"`
	PenaltyExp           int     `json:"penaltyExp"`
	PenaltySpiritStones  int     `json:"penaltySpiritStones"`
	DailyLimit           int     `json:"dailyLimit"`
	TotalLimit           int     `json:"totalLimit"`
	CompletedCount       int     `json:"completedCount"`
	TodayCompletionCount int     `json:"todayCompletionCount"`
	LastCompletedDate    string  `json:"lastCompletedDate"`
	RemindBefore         int     `json:"remindBefore"`
	RemindInterval       int     `json:"remindInterval"`
	LastRemindedAt       *string `json:"lastRemindedAt"`
	SortOrder            int     `json:"sortOrder"`
	CreatedAt            string  `json:"createdAt"`
	UpdatedAt            string  `json:"updatedAt"`
}

type TaskListResp struct {
	Tasks []TaskResp `json:"tasks"`
}

type ReorderTasksReq struct {
	TaskIDs []int64 `json:"taskIds"`
}

// Quick Task (API shortcut)
// POST /api/tasks/quick
//
// Difficulty template (auto-filled):
//   0★: fatigue=1,  spiritStones=10,   attrBonus=0
//   1★: fatigue=5,  spiritStones=50,   attrBonus=0.1
//   2★: fatigue=10, spiritStones=120,  attrBonus=0.2
//   3★: fatigue=20, spiritStones=300,  attrBonus=0.4
//   4★: fatigue=40, spiritStones=800,  attrBonus=0.7
//   5★: fatigue=90, spiritStones=2500, attrBonus=1.0
//
// Categories (attribute keys, each selected one gets attrBonus):
//   "physique"     - 体魄 💪 (exercise, health, diet)
//   "willpower"    - 意志 🧠 (discipline, habits, meditation)
//   "intelligence" - 智力 📚 (study, reading, coding)
//   "perception"   - 感知 👁 (observation, art, reflection)
//   "charisma"     - 魅力 ✨ (communication, networking)
//   "agility"      - 敏捷 🏃 (speed, execution, coordination)
//
// Task types:
//   "once"       - (default) Create + complete in one shot, immediate rewards
//   "repeatable" - Create only, stays active for repeated completion via POST /api/tasks/complete/:id
//   "challenge"  - Create only, has deadline, penalties on failure
type QuickTaskReq struct {
	Title      string   `json:"title"`      // Optional, auto-generated if empty
	Difficulty int      `json:"difficulty"`  // 0-5 stars
	Categories []string `json:"categories"` // Attribute keys
	Type       string   `json:"type"`       // once (default), repeatable, challenge
	DailyLimit int      `json:"dailyLimit"` // For repeatable: max completions per day (0=unlimited)
	TotalLimit int      `json:"totalLimit"` // For repeatable: max total completions (0=unlimited)
	Deadline   string   `json:"deadline"`   // For challenge: ISO8601 deadline
	Source     string   `json:"source"`     // e.g. "ios-shortcut", "api"
}

// Telegram
type BindCodeResp struct {
	Code        string `json:"code"`
	BotUsername string `json:"botUsername"`
	ExpiresIn   int    `json:"expiresIn"` // seconds
}

type TgStatusResp struct {
	Bound      bool   `json:"bound"`
	TgUsername string `json:"tgUsername"`
	TgChatID   int64  `json:"tgChatId"`
}

// Bark Push Notification
type SetBarkKeyReq struct {
	BarkKey string `json:"barkKey"`
}

type BarkStatusResp struct {
	Enabled bool   `json:"enabled"`
	BarkKey string `json:"barkKey"` // Masked for security
}

type TestBarkReq struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

// Sleep
type RecordSleepReq struct {
	SleepStart string `json:"sleepStart"` // ISO8601 format
	SleepEnd   string `json:"sleepEnd"`   // ISO8601 format
	Quality    string `json:"quality"`    // poor, fair, good, excellent
}

type SleepRecordResp struct {
	ID            int64   `json:"id"`
	UserID        int64   `json:"userId"`
	SleepStart    string  `json:"sleepStart"`
	SleepEnd      string  `json:"sleepEnd"`
	DurationHours float64 `json:"durationHours"`
	Quality       string  `json:"quality"`
	EnergyGained  int     `json:"energyGained"`
	CreatedAt     string  `json:"createdAt"`
}

type SleepRecordListResp struct {
	Records []SleepRecordResp `json:"records"`
}

// Shop
type ShopItemResp struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	SellPrice   int    `json:"sellPrice"`
	ItemType    string `json:"itemType"`
	Icon        string `json:"icon"`
	Image       string `json:"image"`
	Stock       int    `json:"stock"`
}

type CreateShopItemReq struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	SellPrice   int    `json:"sellPrice"`
	ItemType    string `json:"itemType"`
	Icon        string `json:"icon"`
	Image       string `json:"image"`
	Stock       int    `json:"stock"`
}

type UpdateShopItemReq struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Price       *int    `json:"price,omitempty"`
	SellPrice   *int    `json:"sellPrice,omitempty"`
	ItemType    *string `json:"itemType,omitempty"`
	Icon        *string `json:"icon,omitempty"`
	Image       *string `json:"image,omitempty"`
	Stock       *int    `json:"stock,omitempty"`
}

type ShopItemListResp struct {
	Items []ShopItemResp `json:"items"`
}

type PurchaseItemReq struct {
	ItemID   int64 `json:"itemId"`
	Quantity int   `json:"quantity"`
}

type PurchaseResult struct {
	Success              bool   `json:"success"`
	Message              string `json:"message"`
	RemainingSpiritStones int    `json:"remainingSpiritStones"`
}

type InventoryItemResp struct {
	ID          int64  `json:"id"`
	ItemID      int64  `json:"itemId"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ItemType    string `json:"itemType"`
	SellPrice   int    `json:"sellPrice"`
	Icon        string `json:"icon"`
	Image       string `json:"image"`
	Quantity    int    `json:"quantity"`
}

type SellItemReq struct {
	ItemID   int64 `json:"itemId"`
	Quantity int   `json:"quantity"`
}

type SellItemResult struct {
	Success              bool   `json:"success"`
	Message              string `json:"message"`
	RemainingSpiritStones int    `json:"remainingSpiritStones"`
}

type InventoryListResp struct {
	Items []InventoryItemResp `json:"items"`
}

type UseItemReq struct {
	ItemID   int64 `json:"itemId"`
	Quantity int   `json:"quantity"`
}

type UseItemResult struct {
	Success   bool          `json:"success"`
	Message   string        `json:"message"`
	Character CharacterResp `json:"character"`
}

type PurchaseRecordResp struct {
	ID         int64  `json:"id"`
	ItemName   string `json:"itemName"`
	Quantity   int    `json:"quantity"`
	TotalPrice int    `json:"totalPrice"`
	CreatedAt  string `json:"createdAt"`
}

type PurchaseHistoryResp struct {
	History []PurchaseRecordResp `json:"history"`
}

// Timeline
type TimelineEvent struct {
	ID          string           `json:"id"`
	Type        string           `json:"type"` // task_complete, task_fail, task_delete, sleep, purchase
	Title       string           `json:"title"`
	Description string           `json:"description"`
	Rewards     *TimelineRewards `json:"rewards,omitempty"`
	Timestamp   string           `json:"timestamp"`
}

type TimelineRewards struct {
	Exp          int `json:"exp,omitempty"`
	SpiritStones int `json:"spiritStones,omitempty"`
}

type TimelineResp struct {
	Events            []TimelineEvent `json:"events"`
	TasksCompleted    int             `json:"tasksCompleted"`
	TotalExp          int             `json:"totalExp"`
	TotalSpiritStones int             `json:"totalSpiritStones"`
	SleepRecords      int             `json:"sleepRecords"`
}

// Common
type CommonResp struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
