package components

import (
	"time"

	"github.com/yohamta/donburi"
)

// FireSkillData 射击技能配置数据
type FireSkillData struct {
	FireRateHz        float64       // 每秒发射次数
	BulletsPerShot    int           // 每次发射子弹数
	SpreadDeg         float64       // 散射总角度（度）
	BulletSpeed       float64       // 子弹速度
	BulletDamage      int           // 子弹伤害值
	BurstChance       float64       // 概率连续发射（0-1）
	PenetrationCount  int           // 可穿透敌人数
	EnableHoming      bool          // 是否追踪
	HomingTurnRateRad float64       // 每帧最大转向弧度
	BurstInterval     time.Duration // 连射间隔
	LastShot          time.Time     // 上次射击时间
	ShotDelay         time.Duration // 射击冷却
	ScheduledShots    []time.Time   // 计划中的连射
}

// FireSkill 射击技能组件
var FireSkill = donburi.NewComponentType[FireSkillData]()
