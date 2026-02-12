package realm

import "math"

// AttrCap returns the attribute cap for the given realm index.
// Formula: 100 * 2^(N+1), where N is the realm index.
// 凡人(0)→200, 炼气(1)→400, 筑基(2)→800, ...
func AttrCap(realmIndex int) float64 {
	return 100.0 * math.Pow(2, float64(realmIndex+1))
}

// AttrMin returns the minimum attribute value for the given realm.
// This equals the previous realm's cap, or 0 for 凡人.
func AttrMin(realmIndex int) float64 {
	if realmIndex <= 0 {
		return 0
	}
	return AttrCap(realmIndex - 1)
}

// BreakthroughExpRequired returns the experience needed to break through the bottleneck
// at the given realm index. Formula: 1000 * 2^N.
func BreakthroughExpRequired(realmIndex int) int {
	return int(1000.0 * math.Pow(2, float64(realmIndex)))
}

// FatigueCapForLevel returns the fatigue cap for the given fatigue level.
// Formula: 100 * 2^N.
func FatigueCapForLevel(level int) int {
	return int(100.0 * math.Pow(2, float64(level)))
}

// GetRealmForValue returns the realm index for the given attribute value.
// The realm is determined by which cap range the value falls in.
func GetRealmForValue(value float64) int {
	for i := MaxRealm; i >= 0; i-- {
		if value >= AttrMin(i) {
			return i
		}
	}
	return 0
}

// IsBottleneck returns true if the attribute value has reached or exceeded
// the cap for the given realm.
func IsBottleneck(value float64, realmIndex int) bool {
	return value >= AttrCap(realmIndex)
}

// GetRealmName returns the Chinese name for the given realm index.
func GetRealmName(realmIndex int) string {
	if realmIndex < 0 || realmIndex >= len(Realms) {
		return "未知"
	}
	return Realms[realmIndex].Name
}

// GetSubRealmName returns the Chinese name for the given sub-realm index.
func GetSubRealmName(subRealm int) string {
	if subRealm < 0 || subRealm >= len(SubRealmNames) {
		return "未知"
	}
	return SubRealmNames[subRealm]
}

// GetFullRealmName returns e.g. "筑基·中期".
func GetFullRealmName(realmIndex, subRealm int) string {
	return GetRealmName(realmIndex) + "·" + GetSubRealmName(subRealm)
}

// AttrGainResult holds the output of ProcessAttrGain.
type AttrGainResult struct {
	NewValue       float64
	NewAccPool     float64
	NewRealmExp    int
	NewIsBottleneck bool
}

// ProcessAttrGain applies an attribute gain, handling bottleneck and accumulation pool logic.
//
// If NOT bottleneck: add gain to value, then check if value >= cap → set bottleneck.
// If IS bottleneck: gain goes to accumulation_pool, and also adds to realm_exp.
func ProcessAttrGain(currentValue, gain float64, realmIndex, realmExp int, isBottleneck bool, accPool float64) AttrGainResult {
	if !isBottleneck {
		newValue := currentValue + gain
		cap := AttrCap(realmIndex)
		newBottleneck := newValue >= cap
		if newBottleneck {
			newValue = cap // clamp to cap
		}
		return AttrGainResult{
			NewValue:        newValue,
			NewAccPool:      accPool,
			NewRealmExp:     realmExp,
			NewIsBottleneck: newBottleneck,
		}
	}

	// Bottleneck: gain goes to accumulation pool and realm exp
	newAccPool := accPool + gain
	newRealmExp := realmExp + int(gain)
	return AttrGainResult{
		NewValue:        currentValue,
		NewAccPool:      newAccPool,
		NewRealmExp:     newRealmExp,
		NewIsBottleneck: true,
	}
}
