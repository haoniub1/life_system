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
	UserID           int64   `json:"userId"`
	Level            int     `json:"level"`
	Exp              int     `json:"exp"`
	Strength         float64 `json:"strength"`
	Intelligence     float64 `json:"intelligence"`
	Vitality         float64 `json:"vitality"`
	Spirit           float64 `json:"spirit"`
	HP               int     `json:"hp"`
	MaxHP            int     `json:"maxHp"`
	Gold             int     `json:"gold"`
	Title            string  `json:"title"`
	LastActivityDate string  `json:"lastActivityDate"`
	Energy           int     `json:"energy"`
	MaxEnergy        int     `json:"maxEnergy"`
	MentalPower      int     `json:"mentalPower"`
	PhysicalPower    int     `json:"physicalPower"`
	MentalSleepAid   int     `json:"mentalSleepAid"`
	PhysicalSleepAid int     `json:"physicalSleepAid"`
}

type UpdateCharacterReq struct {
	DisplayName  *string  `json:"displayName,omitempty"`
	Strength     *float64 `json:"strength,omitempty"`
	Intelligence *float64 `json:"intelligence,omitempty"`
	Vitality     *float64 `json:"vitality,omitempty"`
	Spirit       *float64 `json:"spirit,omitempty"`
}

// Task
type CreateTaskReq struct {
	Title              string  `json:"title"`
	Description        string  `json:"description"`
	Category           string  `json:"category"`
	Type               string  `json:"type"`     // once, repeatable, challenge
	Deadline           string  `json:"deadline"` // ISO8601 format
	RewardExp          int     `json:"rewardExp"`
	RewardGold         int     `json:"rewardGold"`
	RewardStrength     float64 `json:"rewardStrength"`
	RewardIntelligence float64 `json:"rewardIntelligence"`
	RewardVitality     float64 `json:"rewardVitality"`
	RewardSpirit       float64 `json:"rewardSpirit"`
	PenaltyExp         int     `json:"penaltyExp"`
	PenaltyGold        int     `json:"penaltyGold"`
	DailyLimit         int     `json:"dailyLimit"`
	TotalLimit         int     `json:"totalLimit"`
	RemindBefore       int     `json:"remindBefore"`   // minutes
	RemindInterval     int     `json:"remindInterval"` // minutes
	CostMental         int     `json:"costMental"`
	CostPhysical       int     `json:"costPhysical"`
}

type UpdateTaskReq struct {
	Title              *string  `json:"title,omitempty"`
	Description        *string  `json:"description,omitempty"`
	Category           *string  `json:"category,omitempty"`
	Type               *string  `json:"type,omitempty"`
	Deadline           *string  `json:"deadline,omitempty"` // ISO8601 format
	RewardExp          *int     `json:"rewardExp,omitempty"`
	RewardGold         *int     `json:"rewardGold,omitempty"`
	RewardStrength     *float64 `json:"rewardStrength,omitempty"`
	RewardIntelligence *float64 `json:"rewardIntelligence,omitempty"`
	RewardVitality     *float64 `json:"rewardVitality,omitempty"`
	RewardSpirit       *float64 `json:"rewardSpirit,omitempty"`
	PenaltyExp         *int     `json:"penaltyExp,omitempty"`
	PenaltyGold        *int     `json:"penaltyGold,omitempty"`
	DailyLimit         *int     `json:"dailyLimit,omitempty"`
	TotalLimit         *int     `json:"totalLimit,omitempty"`
	RemindBefore       *int     `json:"remindBefore,omitempty"`
	RemindInterval     *int     `json:"remindInterval,omitempty"`
	CostMental         *int     `json:"costMental,omitempty"`
	CostPhysical       *int     `json:"costPhysical,omitempty"`
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
	RewardExp            int     `json:"rewardExp"`
	RewardGold           int     `json:"rewardGold"`
	RewardStrength       float64 `json:"rewardStrength"`
	RewardIntelligence   float64 `json:"rewardIntelligence"`
	RewardVitality       float64 `json:"rewardVitality"`
	RewardSpirit         float64 `json:"rewardSpirit"`
	PenaltyExp           int     `json:"penaltyExp"`
	PenaltyGold          int     `json:"penaltyGold"`
	DailyLimit           int     `json:"dailyLimit"`
	TotalLimit           int     `json:"totalLimit"`
	CompletedCount       int     `json:"completedCount"`
	TodayCompletionCount int     `json:"todayCompletionCount"`
	LastCompletedDate    string  `json:"lastCompletedDate"`
	RemindBefore         int     `json:"remindBefore"`
	RemindInterval       int     `json:"remindInterval"`
	LastRemindedAt       *string `json:"lastRemindedAt"`
	CreatedAt            string  `json:"createdAt"`
	UpdatedAt            string  `json:"updatedAt"`
	CostMental           int     `json:"costMental"`
	CostPhysical         int     `json:"costPhysical"`
}

type TaskListResp struct {
	Tasks []TaskResp `json:"tasks"`
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
	Icon        string `json:"icon"`
	Image       string `json:"image"`
	Stock       int    `json:"stock"`
}

type CreateShopItemReq struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	Icon        string `json:"icon"`
	Image       string `json:"image"`
	Stock       int    `json:"stock"`
}

type UpdateShopItemReq struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Price       *int    `json:"price,omitempty"`
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
	Success       bool   `json:"success"`
	Message       string `json:"message"`
	RemainingGold int    `json:"remainingGold"`
}

type InventoryItemResp struct {
	ID          int64  `json:"id"`
	ItemID      int64  `json:"itemId"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Image       string `json:"image"`
	Quantity    int    `json:"quantity"`
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
	Timestamp   string           `json:"timestamp"`
	Rewards     *TimelineRewards `json:"rewards,omitempty"`
}

type TimelineRewards struct {
	Exp    int `json:"exp,omitempty"`
	Gold   int `json:"gold,omitempty"`
	Energy int `json:"energy,omitempty"`
}

type TimelineResp struct {
	Events         []TimelineEvent `json:"events"`
	TasksCompleted int             `json:"tasksCompleted"`
	TotalExp       int             `json:"totalExp"`
	TotalGold      int             `json:"totalGold"`
	SleepRecords   int             `json:"sleepRecords"`
}

// Common
type CommonResp struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
