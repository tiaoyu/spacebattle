package battle

import (
	"time"
)

// ShooterGame 相关的基础类型与通用工具函数

// PlayerShip 玩家飞船
type PlayerShip struct {
	X, Y   float64         // 位置
	Width  float64         // 宽度
	Height float64         // 高度
	Speed  float64         // 速度
	Fire   FireSkillConfig // 发射技能配置
}

// PlayerOptions 用于外部场景定制玩家初始属性
type PlayerOptions struct {
	Speed                float64
	SizeScale            float64 // 1.0 为原始大小
	Lives                int
	DifficultyMultiplier float64
	PassiveKey           string
	// 升级加成
	ModFireRateHz     float64 // 叠加到 FireRateHz
	ModBulletsPerShot int     // 叠加到 BulletsPerShot
	ModPenetration    int     // 叠加到 PenetrationCount
	ModSpreadDeltaDeg float64 // 负值缩小散射
	ModBulletSpeed    float64
	ModBulletDamage   int     // 叠加到子弹伤害
	ModBurstChance    float64 // 0..1 范围叠加
	ModEnableHoming   bool
	ModTurnRateRad    float64
}

// Bullet 子弹
type Bullet struct {
	X, Y                 float64 // 位置
	VX, VY               float64 // 方向
	Width                float64 // 宽度
	Height               float64 // 高度
	Active               bool    // 是否激活
	Speed                float64 // 速度
	RemainingPenetration int     // 剩余穿透次数
	Homing               bool    // 是否追踪
	HomingTurnRate       float64 // 追踪转向速率
}

// EnemyShip 敌机
type EnemyShip struct {
	X, Y      float64 // 位置
	VX, VY    float64 // 方向
	Width     float64 // 宽度
	Height    float64 // 高度
	Active    bool    // 是否激活
	Health    int     // 生命值
	MaxHealth int     // 最大生命值
}

// Explosion 爆炸效果
type Explosion struct {
	X, Y      float64 // 位置
	Radius    float64 // 半径
	MaxRadius float64 // 最大半径
	Active    bool    // 是否激活
	Timer     float64 // 计时器
}

// Star 星星
type Star struct {
	X, Y  float64 // 位置
	Speed float64 // 速度
	Size  float64 // 大小
}

// FireSkillConfig 发射技能配置
type FireSkillConfig struct {
	FireRateHz        float64       // 每秒发射次数
	BulletsPerShot    int           // 每次发射子弹数
	SpreadDeg         float64       // 散射总角度（度）
	BulletSpeed       float64       // 子弹速度
	BurstChance       float64       // 概率连续发射（0-1）
	PenetrationCount  int           // 可穿透敌人数
	EnableHoming      bool          // 是否追踪
	HomingTurnRateRad float64       // 每帧最大转向弧度
	BurstInterval     time.Duration // 连射间隔
}
