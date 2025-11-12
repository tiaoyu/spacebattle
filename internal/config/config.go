package config

import "time"

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
		DifficultyMax:      100000.0,
		DifficultyCostBase: 15,
		EnemyDelayFloor:    80 * time.Millisecond,
		MaxSimultaneousCap: 60,
		DiffDelayLogK:      1.0,
		DiffCapLogK:        0.2,
		DiffBatchLogK:      0.15,
		DiffSpeedLogK:      0.2,
		DiffHpLogK:         2.0,
	}
}
