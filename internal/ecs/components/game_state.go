package components

import (
	"time"

	"github.com/yohamta/donburi"
)

// GameStateData 游戏状态数据（全局单例）
type GameStateData struct {
	Score            int
	Lives            int
	GameOver         bool
	Victory          bool
	Settled          bool
	KilledEnemyCount int
	SpawnedCount     int
	BossSpawned      bool
	BossKilled       bool
	RewardCached     int
	DifficultyMul    float64
	StartTime        time.Time
	// 波次配置
	TotalDuration      time.Duration
	SmallPhaseDuration time.Duration
	WaveLength         time.Duration
	WaveCount          int
	WaveIndex          int
	WaveMinIntervals   []time.Duration
	// 生成配置
	MaxSimultaneous int
	BatchSize       int
	LastEnemyTime   time.Time
	// GM 调试
	GMOpen  bool
	GMIndex int
	GMTab   int
	// 奖励分解信息
	RewardBreakdown RewardBreakdownData
}

// RewardBreakdownData 功勋奖励分解信息
type RewardBreakdownData struct {
	BaseReward       int
	DifficultyBonus  int
	KillBonus        int
	SpeedBonus       int
	PerfectBonus     int
	BossBonus        int
	TotalReward      int
	PerformanceScore float64
}

// GameState 游戏状态组件
var GameState = donburi.NewComponentType[GameStateData]()

