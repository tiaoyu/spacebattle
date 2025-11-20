package components

import (
	"time"

	"github.com/yohamta/donburi"
)

// ShipAbilityData 战机被动技能数据
type ShipAbilityData struct {
	AbilityType string // "harvest", "speed_frenzy", "dodge_master", "energy_shield"

	// Harvest (Alpha) - 战场收割
	// 每击杀 5 个敌机回复 1 点生命值
	KillCounter int

	// SpeedFrenzy (Beta) - 速度狂热
	// 连续击杀敌机叠加射速 buff（每层 +5% 射速，最多 5 层，3 秒未击杀清空）
	FrenzyStacks int
	LastKillTime time.Time

	// DodgeMaster (Gamma) - 闪避大师
	// 受到伤害后获得 2 秒无敌时间（冷却 10 秒）
	InvulnTime       float64 // 剩余无敌时间（秒）
	InvulnCooldown   float64 // 剩余冷却时间（秒）
	IsInvulnerable   bool    // 是否处于无敌状态
	LastDamageTaken  time.Time

	// EnergyShield (Delta) - 能量护盾
	// 额外护盾值等于最大生命值，护盾先于生命承受伤害，5 秒不受伤害自动回复护盾
	ShieldCurrent  int
	ShieldMax      int
	LastDamageTime time.Time
}

// ShipAbility 战机被动技能组件
var ShipAbility = donburi.NewComponentType[ShipAbilityData]()

