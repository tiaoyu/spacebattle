package systems

import (
	"math"
	"math/rand"
	"time"

	"spacebattle/internal/config"
	"spacebattle/internal/ecs"
	"spacebattle/internal/ecs/components"
	"spacebattle/internal/ecs/tags"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
)

// SpawnSystem 生成系统
type SpawnSystem struct {
	world *ecs.World
	cfg   *config.Config
}

// NewSpawnSystem 创建生成系统
func NewSpawnSystem(world *ecs.World) *SpawnSystem {
	return &SpawnSystem{
		world: world,
		cfg:   config.DefaultConfig(),
	}
}

// Update 更新生成逻辑
func (s *SpawnSystem) Update(w donburi.World) {
	// 获取游戏状态
	var gameState *components.GameStateData
	query.NewQuery(filter.Contains(components.GameState)).Each(w, func(entry *donburi.Entry) {
		gameState = components.GameState.Get(entry)
	})

	if gameState == nil {
		return
	}

	// 检查是否已经有 Boss
	bossExists := false
	query.NewQuery(filter.Contains(tags.Boss)).Each(w, func(entry *donburi.Entry) {
		bossExists = true
	})

	// Boss 阶段
	if bossExists {
		return
	}

	// 检查是否到达 Boss 生成时间
	elapsed := time.Since(gameState.StartTime)
	if elapsed >= gameState.SmallPhaseDuration {
		if !gameState.BossSpawned {
			s.SpawnBoss(w, gameState)
			gameState.BossSpawned = true
		}
		return
	}

	// 小怪生成阶段
	s.SpawnEnemies(w, gameState, elapsed)
}

// SpawnEnemies 生成敌机
func (s *SpawnSystem) SpawnEnemies(w donburi.World, gameState *components.GameStateData, elapsed time.Duration) {
	// 计算当前波次
	waveIndex := int(elapsed / gameState.WaveLength)
	if waveIndex < 0 {
		waveIndex = 0
	}
	if waveIndex >= gameState.WaveCount {
		waveIndex = gameState.WaveCount - 1
	}
	gameState.WaveIndex = waveIndex

	// 难度系数
	diff := gameState.DifficultyMul
	if diff <= 0 {
		diff = 1
	}

	// 计算生成间隔
	enemyDelay := gameState.WaveMinIntervals[waveIndex]

	// 根据难度调整间隔
	delayFactor := 1.0
	if diff > 1 {
		k := s.cfg.DiffDelayLogK
		if k <= 0 {
			k = 1
		}
		delayFactor = 1.0 / math.Max(1.0, 1.0+k*math.Log10(diff))
	} else {
		delayFactor = 1.0 / math.Max(0.5, diff)
	}
	enemyDelay = time.Duration(float64(enemyDelay) * delayFactor)

	// 应用最小间隔
	if s.cfg.EnemyDelayFloor > 0 && enemyDelay < s.cfg.EnemyDelayFloor {
		enemyDelay = s.cfg.EnemyDelayFloor
	}

	// 检查是否到达生成时间
	if time.Since(gameState.LastEnemyTime) < enemyDelay {
		return
	}

	// 计算同屏上限
	capScale := 1.0
	if diff > 1 {
		k := s.cfg.DiffCapLogK
		capScale = 1.0 + k*math.Log10(diff)
	} else {
		capScale = math.Max(0.5, diff)
	}
	effMax := int(float64(gameState.MaxSimultaneous) * capScale)
	if s.cfg.MaxSimultaneousCap > 0 && effMax > s.cfg.MaxSimultaneousCap {
		effMax = s.cfg.MaxSimultaneousCap
	}
	if effMax < 1 {
		effMax = 1
	}

	// 计算当前敌机数量
	activeCount := s.CountActiveEnemies(w)
	space := effMax - activeCount
	if space <= 0 {
		return
	}

	// 计算批量
	batchScale := 1.0
	if diff > 1 {
		k := s.cfg.DiffBatchLogK
		batchScale = math.Min(2.5, 1.0+k*math.Log10(diff))
	}
	batch := min(max(int(float64(gameState.BatchSize)*batchScale), 1), space)
	if batch <= 0 {
		return
	}

	// 计算敌机生命值
	enemyHP := min(1+waveIndex/2, 5)
	hpScale := 1.0
	if diff > 1 {
		k := s.cfg.DiffHpLogK
		if k <= 0 {
			k = 1
		}
		hpScale = 1.0 + k*math.Log2(diff)
	} else {
		hpScale = math.Max(0.5, diff)
	}
	enemyHP = max(int(math.Ceil(float64(enemyHP)*hpScale)), 1)

	// 计算速度缩放
	speedScale := 1.0
	if diff > 1 {
		k := s.cfg.DiffSpeedLogK
		speedScale = 1.0 + k*math.Log10(diff)
	} else {
		speedScale = math.Max(0.6, diff)
	}

	// 生成敌机（根据波次权重随机选择类型）
	for range batch {
		x := float64(rand.Intn(760))
		y := -30.0
		vx := (float64(rand.Intn(3) - 1)) * speedScale
		vy := (2 + float64(rand.Intn(2))) * speedScale

		// 根据波次确定敌机类型
		enemyType := s.selectEnemyType(waveIndex)

		switch enemyType {
		case "shooter":
			// 射击型：生命更高，速度稍慢
			hp := enemyHP + waveIndex/2
			s.world.CreateEnemyShooter(x, y, vx*0.6, vy*0.6, 35, 30, hp)
		case "zigzag":
			// 之字型：生命较低，速度正常
			hp := enemyHP
			if waveIndex >= 4 {
				hp += 1
			}
			s.world.CreateEnemyZigzag(x, y, vx, vy, 25, 20, hp)
		case "tank":
			// 肉盾型：高生命，慢速
			cfg := s.cfg
			hp := int(float64(enemyHP) * cfg.EnemyTankHealthMultiplier)
			s.world.CreateEnemyTank(x, y, 0, vy*0.5, 45, 35, hp)
		default:
			// 基础型：标准属性
			s.world.CreateEnemy(x, y, vx, vy, 30, 25, enemyHP)
		}

		gameState.SpawnedCount++
	}

	gameState.LastEnemyTime = time.Now()
}

// selectEnemyType 根据波次选择敌机类型
func (s *SpawnSystem) selectEnemyType(waveIndex int) string {
	roll := rand.Float64() * 100

	// 调整为更早出现不同类型（方便测试）
	switch waveIndex {
	case 0:
		// 第 1 波：基础型 70%, 之字型 30%
		if roll < 70 {
			return "basic"
		}
		return "zigzag"
	case 1:
		// 第 2 波：基础型 50%, 之字型 25%, 射击型 25%
		if roll < 50 {
			return "basic"
		} else if roll < 75 {
			return "zigzag"
		}
		return "shooter"
	case 2:
		// 第 3 波：基础型 40%, 之字型 20%, 射击型 20%, 肉盾型 20%
		if roll < 40 {
			return "basic"
		} else if roll < 60 {
			return "zigzag"
		} else if roll < 80 {
			return "shooter"
		}
		return "tank"
	default:
		// 第 4-5 波：所有类型均衡
		if roll < 30 {
			return "basic"
		} else if roll < 50 {
			return "zigzag"
		} else if roll < 75 {
			return "shooter"
		}
		return "tank"
	}
}

// SpawnBoss 生成 Boss
func (s *SpawnSystem) SpawnBoss(w donburi.World, gameState *components.GameStateData) {
	diff := gameState.DifficultyMul
	if diff <= 0 {
		diff = 1
	}

	// Boss 血量随难度缩放
	bossHP := int(math.Ceil(60 * math.Max(1.0, 1.0+s.cfg.DiffHpLogK*math.Log10(math.Max(1.0, diff)))))

	s.world.CreateBoss(350, 60, 1.2, 0.8, 100, 60, bossHP)
}

// CountActiveEnemies 计算当前活跃敌机数量（所有类型）
func (s *SpawnSystem) CountActiveEnemies(w donburi.World) int {
	count := 0
	// 统计基础型
	query.NewQuery(filter.Contains(tags.Enemy)).Each(w, func(entry *donburi.Entry) {
		count++
	})
	// 统计射击型
	query.NewQuery(filter.Contains(tags.EnemyShooter)).Each(w, func(entry *donburi.Entry) {
		count++
	})
	// 统计之字型
	query.NewQuery(filter.Contains(tags.EnemyZigzag)).Each(w, func(entry *donburi.Entry) {
		count++
	})
	// 统计肉盾型
	query.NewQuery(filter.Contains(tags.EnemyTank)).Each(w, func(entry *donburi.Entry) {
		count++
	})
	return count
}
