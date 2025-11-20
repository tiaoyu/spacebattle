package config

import (
	"image/color"
	"time"
)

// Config 游戏配置
type Config struct {
	// 窗口设置
	WindowWidth  int
	WindowHeight int
	WindowTitle  string

	// 游戏设置
	FPS int

	// 调试设置
	Debug bool

	// —— 关卡/波次（时间制）配置 ——
	// 总时长（含Boss）与小怪阶段时长
	TotalDuration      time.Duration
	SmallPhaseDuration time.Duration
	// 波次数与每波长度
	WaveCount  int
	WaveLength time.Duration
	// 每波最小生成间隔（逐波缩短）
	WaveMinIntervals []time.Duration
	// 生成限制
	MaxSimultaneous int
	BatchSize       int
	// Boss 收束（软收束阈值）
	BossSoftEnrageStart   time.Duration // 例如 10s 后提高伤害倍率
	BossTripleDamageStart time.Duration // 例如 14s 后再次提高

	// —— 升级消耗（功勋） ——
	UpgradeCostFireRate       int
	UpgradeCostBulletsPerShot int
	UpgradeCostPenetration    int
	UpgradeCostSpreadNarrow   int
	UpgradeCostBulletSpeed    int
	UpgradeCostBulletDamage   int
	UpgradeCostBurstChance    int
	UpgradeCostEnableHoming   int
	UpgradeCostTurnRate       int

	// —— 战机属性上限（用于 clamp） ——
	MaxFireRateHz     float64
	MaxBulletsPerShot int
	MaxPenetration    int
	MaxSpreadDeg      float64
	MaxBulletSpeed    float64
	MaxBurstChance    float64
	MaxTurnRateRad    float64
	MaxLives          int

	// —— 出征难度可配 ——
	DifficultyMin      float64
	DifficultyMax      float64
	DifficultyCostBase int           // 每+1.0难度的基础成本（线性）
	EnemyDelayFloor    time.Duration // 生成间隔下限
	MaxSimultaneousCap int           // 同屏上限最大cap
	DiffDelayLogK      float64       // 冷却缩放 log 系数
	DiffCapLogK        float64       // 同屏上限缩放 log 系数
	DiffBatchLogK      float64       // 批量缩放 log 系数
	DiffSpeedLogK      float64       // 敌机速度缩放 log 系数
	DiffHpLogK         float64       // 敌机生命缩放 log 系数

	// —— 敌机类型配置 ——
	EnemyShooterShootInterval time.Duration // 射击型敌机射击间隔
	EnemyZigzagPeriod         float64       // 之字型敌机摆动周期（秒）
	EnemyTankHealthMultiplier float64       // 肉盾型敌机生命倍率

	// —— 战机被动配置 ——
	HarvestKillsRequired int     // Alpha - 回血所需击杀数
	FrenzyStackDuration  float64 // Beta - buff持续时间（秒）
	FrenzySpeedBonus     float64 // Beta - 每层射速加成（比例）
	FrenzyMaxStacks      int     // Beta - 最大层数
	DodgeInvulnDuration  float64 // Gamma - 无敌时间（秒）
	DodgeInvulnCooldown  float64 // Gamma - 冷却时间（秒）
	ShieldRegenDelay     float64 // Delta - 护盾回复延迟（秒）

	// —— 粒子效果配置 ——
	ParticleCountOnKill int     // 击杀敌机生成的粒子数量
	ParticleLifetime    float64 // 粒子生命周期（秒）
	ParticleMaxCount    int     // 同屏最大粒子数

	// —— 颜色配置 ——
	BackgroundColor   color.RGBA
	PlayerColor       color.RGBA // 玩家
	BulletColor       color.RGBA // 玩家子弹
	EnemyColor        color.RGBA // 普通敌机
	EnemyShooterColor color.RGBA // 射击型敌机
	EnemyZigzagColor  color.RGBA // 之字型敌机
	EnemyTankColor    color.RGBA // 肉盾型敌机
	EnemyBulletColor  color.RGBA // 敌机子弹
	BossColor         color.RGBA // Boss
	BossHPBarBg       color.RGBA // Boss血条背景
	BossHPBarFg       color.RGBA // Boss血条前景
	ExplosionColor    color.RGBA // 爆炸效果
	StarColor         color.RGBA // 星星背景
	ParticleColor     color.RGBA // 粒子效果

	// —— UI 颜色配置 ——
	UIBackgroundColor  color.RGBA // UI背景色
	UITextColor        color.RGBA // 普通文本
	UIHighlightColor   color.RGBA // 高亮文本（选中项）
	UIHintColor        color.RGBA // 提示文本
	UIMeritColor       color.RGBA // 功勋数字颜色
	UIVictoryColor     color.RGBA // 胜利文本颜色
	UIGameOverColor    color.RGBA // 失败文本颜色
	UIOverlayColor     color.RGBA // 半透明遮罩
	UIBossWarningColor color.RGBA // Boss警告文本
	UIGreyTextColor    color.RGBA // 灰色文本
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		WindowWidth:  800,
		WindowHeight: 600,
		WindowTitle:  "spacebattle",
		FPS:          60,
		Debug:        true,

		// 时间制波次默认：与设计案一致
		TotalDuration:      60 * time.Second,
		SmallPhaseDuration: 45 * time.Second,
		WaveCount:          5,
		WaveLength:         9 * time.Second,
		WaveMinIntervals: []time.Duration{
			600 * time.Millisecond,
			500 * time.Millisecond,
			400 * time.Millisecond,
			320 * time.Millisecond,
			250 * time.Millisecond,
		},
		MaxSimultaneous: 120,
		BatchSize:       2,

		BossSoftEnrageStart:   10 * time.Second,
		BossTripleDamageStart: 14 * time.Second,

		// 升级消耗默认
		UpgradeCostFireRate:       1,
		UpgradeCostBulletsPerShot: 3,
		UpgradeCostPenetration:    2,
		UpgradeCostSpreadNarrow:   1,
		UpgradeCostBulletSpeed:    1,
		UpgradeCostBulletDamage:   2,
		UpgradeCostBurstChance:    2,
		UpgradeCostEnableHoming:   2147483647,
		UpgradeCostTurnRate:       2,

		// 战机属性上限默认
		MaxFireRateHz:     30.0,
		MaxBulletsPerShot: 20,
		MaxPenetration:    10,
		MaxSpreadDeg:      180.0,
		MaxBulletSpeed:    30.0,
		MaxBurstChance:    1.0,
		MaxTurnRateRad:    1.0,
		MaxLives:          9,

		// 出征难度默认
		DifficultyMin:      1.0,
		DifficultyMax:      10000000.0,
		DifficultyCostBase: 15,
		EnemyDelayFloor:    80 * time.Millisecond,
		MaxSimultaneousCap: 60,
		DiffDelayLogK:      1.0,
		DiffCapLogK:        0.2,
		DiffBatchLogK:      0.15,
		DiffSpeedLogK:      0.2,
		DiffHpLogK:         2.0,

		// 敌机类型配置默认
		EnemyShooterShootInterval: 2500 * time.Millisecond, // 2.5秒
		EnemyZigzagPeriod:         1.8,                     // 1.8秒摆动周期
		EnemyTankHealthMultiplier: 3.0,                     // 3倍生命

		// 战机被动配置默认
		HarvestKillsRequired: 5,    // 5次击杀回血
		FrenzyStackDuration:  3.0,  // 3秒buff持续时间
		FrenzySpeedBonus:     0.05, // 5%射速加成
		FrenzyMaxStacks:      5,    // 最多5层
		DodgeInvulnDuration:  2.0,  // 2秒无敌
		DodgeInvulnCooldown:  10.0, // 10秒冷却
		ShieldRegenDelay:     5.0,  // 5秒后回复

		// 粒子效果配置默认
		ParticleCountOnKill: 10,  // 每次击杀10个粒子
		ParticleLifetime:    0.5, // 0.5秒生命周期
		ParticleMaxCount:    200, // 同屏最多200个粒子

		// 颜色配置默认
		BackgroundColor:   color.RGBA{R: 10, G: 16, B: 36, A: 255},    // 深蓝色背景
		PlayerColor:       color.RGBA{R: 100, G: 200, B: 255, A: 255}, // 浅蓝色玩家
		BulletColor:       color.RGBA{R: 255, G: 255, B: 100, A: 255}, // 黄色子弹
		EnemyColor:        color.RGBA{R: 255, G: 100, B: 100, A: 255}, // 浅红色（与背景对比明显）
		EnemyShooterColor: color.RGBA{R: 255, G: 165, B: 0, A: 255},   // 橙色（与背景对比明显）
		EnemyZigzagColor:  color.RGBA{R: 100, G: 255, B: 200, A: 255}, // 青绿色（与背景对比明显）
		EnemyTankColor:    color.RGBA{R: 200, G: 50, B: 50, A: 255},   // 深红色（与背景对比明显）
		EnemyBulletColor:  color.RGBA{R: 255, G: 0, B: 0, A: 255},     // 红色子弹
		BossColor:         color.RGBA{R: 200, G: 50, B: 200, A: 255},  // 紫色Boss（与背景对比明显）
		BossHPBarBg:       color.RGBA{R: 100, G: 100, B: 100, A: 255}, // 灰色血条背景
		BossHPBarFg:       color.RGBA{R: 255, G: 0, B: 0, A: 255},     // 红色血条前景
		ExplosionColor:    color.RGBA{R: 255, G: 150, B: 0, A: 255},   // 橙色爆炸
		StarColor:         color.RGBA{R: 200, G: 200, B: 200, A: 255}, // 灰白色星星
		ParticleColor:     color.RGBA{R: 255, G: 100, B: 0, A: 255},   // 橙红色粒子

		// UI 颜色配置默认
		UIBackgroundColor:  color.RGBA{R: 20, G: 20, B: 60, A: 255},    // UI深蓝色背景
		UITextColor:        color.RGBA{R: 255, G: 255, B: 255, A: 255}, // 白色文本
		UIHighlightColor:   color.RGBA{R: 255, G: 255, B: 0, A: 255},   // 黄色高亮
		UIHintColor:        color.RGBA{R: 200, G: 200, B: 200, A: 255}, // 灰白色提示
		UIMeritColor:       color.RGBA{R: 255, G: 215, B: 0, A: 255},   // 金色功勋
		UIVictoryColor:     color.RGBA{R: 100, G: 255, B: 100, A: 255}, // 绿色胜利
		UIGameOverColor:    color.RGBA{R: 255, G: 100, B: 100, A: 255}, // 红色失败
		UIOverlayColor:     color.RGBA{R: 0, G: 0, B: 0, A: 180},       // 半透明黑色遮罩
		UIBossWarningColor: color.RGBA{R: 255, G: 0, B: 0, A: 255},     // 红色Boss警告
		UIGreyTextColor:    color.RGBA{R: 150, G: 150, B: 150, A: 255}, // 灰色文本
	}
}
