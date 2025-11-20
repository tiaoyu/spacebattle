package components

import (
	"time"

	"github.com/yohamta/donburi"
)

// EnemyAIData 敌机AI行为数据
type EnemyAIData struct {
	EnemyType string // "basic", "shooter", "zigzag", "tank"

	// Shooter - 射击型敌机
	// 每 2-3 秒向玩家位置发射 1 发子弹
	ShootInterval time.Duration // 射击间隔
	LastShotTime  time.Time     // 上次射击时间

	// Zigzag - 之字型敌机
	// 横向摆动下降，按正弦波移动
	ZigzagPhase  float64 // 摆动相位（0-2π）
	ZigzagSpeed  float64 // 相位变化速度
	ZigzagPeriod float64 // 摆动周期（秒）
}

// EnemyAI 敌机AI组件
var EnemyAI = donburi.NewComponentType[EnemyAIData]()

