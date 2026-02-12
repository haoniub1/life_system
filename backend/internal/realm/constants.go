package realm

// Realm indices
const (
	RealmFanRen  = 0 // 凡人
	RealmLianQi  = 1 // 炼气
	RealmZhuJi   = 2 // 筑基
	RealmJinDan  = 3 // 金丹
	RealmYuanYing = 4 // 元婴
	RealmHuaShen = 5 // 化神
	RealmHeTi    = 6 // 合体
	RealmDaCheng = 7 // 大乘
	RealmDuJie   = 8 // 渡劫

	MaxRealm = RealmDuJie
)

// Sub-realm indices
const (
	SubRealmChuQi      = 0 // 初期
	SubRealmZhongQi    = 1 // 中期
	SubRealmHouQi      = 2 // 后期
	SubRealmDaYuanMan  = 3 // 大圆满

	MaxSubRealm = SubRealmDaYuanMan
)

// RealmInfo holds display information for a realm.
type RealmInfo struct {
	Index int
	Name  string
}

// Realms lists all 9 realms in order.
var Realms = []RealmInfo{
	{RealmFanRen, "凡人"},
	{RealmLianQi, "炼气"},
	{RealmZhuJi, "筑基"},
	{RealmJinDan, "金丹"},
	{RealmYuanYing, "元婴"},
	{RealmHuaShen, "化神"},
	{RealmHeTi, "合体"},
	{RealmDaCheng, "大乘"},
	{RealmDuJie, "渡劫"},
}

// SubRealmNames maps sub-realm index to Chinese name.
var SubRealmNames = []string{"初期", "中期", "后期", "大圆满"}

// AttrDisplayInfo holds UI display information for an attribute.
type AttrDisplayInfo struct {
	Key     string
	Name    string // Chinese display name
	Emoji   string
	Color   string // Hex color for frontend
	HasRealm bool  // Whether this attribute has the realm/bottleneck system
}

// AttrKeys defines the canonical order of the 6 cultivation attributes + luck.
var AttrKeys = []string{
	"physique",
	"willpower",
	"intelligence",
	"perception",
	"charisma",
	"agility",
}

// AllAttrKeys includes luck.
var AllAttrKeys = []string{
	"physique",
	"willpower",
	"intelligence",
	"perception",
	"charisma",
	"agility",
	"luck",
}

// AttrDisplay maps attribute key to display info.
var AttrDisplay = map[string]AttrDisplayInfo{
	"physique":     {Key: "physique", Name: "体魄", Emoji: "\U0001f4aa", Color: "#ef4444", HasRealm: true},
	"willpower":    {Key: "willpower", Name: "意志", Emoji: "\U0001f9e0", Color: "#8b5cf6", HasRealm: true},
	"intelligence": {Key: "intelligence", Name: "智力", Emoji: "\U0001f4da", Color: "#3b82f6", HasRealm: true},
	"perception":   {Key: "perception", Name: "感知", Emoji: "\U0001f441", Color: "#10b981", HasRealm: true},
	"charisma":     {Key: "charisma", Name: "魅力", Emoji: "\u2728", Color: "#ec4899", HasRealm: true},
	"agility":      {Key: "agility", Name: "敏捷", Emoji: "\U0001f3c3", Color: "#f59e0b", HasRealm: true},
	"luck":         {Key: "luck", Name: "幸运", Emoji: "\U0001f340", Color: "#6366f1", HasRealm: false},
}
