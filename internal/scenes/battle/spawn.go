package battle

import (
	"math"
	"math/rand"
	"time"

	"spacebattle/internal/config"
)

// 刷怪与 Boss 生成逻辑（从 shooter_game.go 拆分）

func (sg *ShooterGame) spawnEnemies() {
	// Boss在场：只更新Boss移动
	if sg.bossActive {
		sg.boss.X += sg.boss.VX
		sg.boss.Y += sg.boss.VY
		if sg.boss.X < 50 || sg.boss.X+sg.boss.Width > 750 {
			sg.boss.VX = -sg.boss.VX
		}
		if sg.boss.Y < 20 || sg.boss.Y+sg.boss.Height > 300 {
			sg.boss.VY = -sg.boss.VY
		}
		return
	}

	// 时间制推进：小怪阶段结束后直接生成Boss（不等待清场）
	elapsed := time.Since(sg.startTime)
	if elapsed >= sg.smallPhaseDuration {
		if !sg.bossSpawned {
			sg.spawnBoss()
		}
		return
	}

	// 计算当前时间波次
	idx := int(elapsed / sg.waveLength)
	if idx < 0 {
		idx = 0
	}
	if idx >= sg.waveCount {
		idx = sg.waveCount - 1
	}
	sg.waveIndex = idx

	// 难度因子（>1 提升强度；<1 降低强度）
	diff := sg.difficultyMul
	if diff <= 0 {
		diff = 1
	}
	cfg := config.DefaultConfig()

	// 基础冷却
	sg.enemyDelay = sg.waveMinIntervals[idx]
	// 冷却按难度缩放：>1 时除以 (1+K*log10(diff))，<1 时除以 max(0.5, diff)
	delayFactor := 1.0
	if diff > 1 {
		k := cfg.DiffDelayLogK
		if k <= 0 {
			k = 1
		}
		delayFactor = 1.0 / math.Max(1.0, 1.0+k*math.Log10(diff))
	} else {
		delayFactor = 1.0 / math.Max(0.5, diff)
	}
	sg.enemyDelay = time.Duration(float64(sg.enemyDelay) * delayFactor)
	// 生成间隔下限
	if cfg.EnemyDelayFloor > 0 && sg.enemyDelay < cfg.EnemyDelayFloor {
		sg.enemyDelay = cfg.EnemyDelayFloor
	}
	if time.Since(sg.lastEnemy) < sg.enemyDelay {
		return
	}
	// 同屏上限按难度缩放（温和）：>1 × (1+K*log10(diff))；<1 × max(0.5, diff)
	capScale := 1.0
	if diff > 1 {
		k := cfg.DiffCapLogK
		capScale = 1.0 + k*math.Log10(diff)
	} else {
		capScale = math.Max(0.5, diff)
	}
	effMax := int(float64(sg.maxSimultaneous) * capScale)
	if cfg.MaxSimultaneousCap > 0 && effMax > cfg.MaxSimultaneousCap {
		effMax = cfg.MaxSimultaneousCap
	}
	if effMax < 1 {
		effMax = 1
	}
	space := effMax - sg.countActiveEnemies()
	if space <= 0 {
		return
	}

	// 按批量与剩余空间生成（批量轻微放大）
	batchScale := 1.0
	if diff > 1 {
		k := cfg.DiffBatchLogK
		batchScale = math.Min(2.5, 1.0+k*math.Log10(diff))
	}
	batch := min(max(int(float64(sg.batchSize)*batchScale), 1), space)
	if batch <= 0 {
		return
	}

	// 敌机生命：基础随波次，上限 5；再按难度缩放
	enemyHP := min(1+sg.waveIndex/2, 5)
	hpScale := 1.0
	if diff > 1 {
		k := cfg.DiffHpLogK
		if k <= 0 {
			k = 1
		}
		hpScale = 1.0 + k*math.Log2(diff)
	} else {
		hpScale = math.Max(0.5, diff)
	}
	enemyHP = max(int(math.Ceil(float64(enemyHP)*hpScale)), 1)

	// 速度缩放
	speedScale := 1.0
	if diff > 1 {
		k := cfg.DiffSpeedLogK
		speedScale = 1.0 + k*math.Log10(diff)
	} else {
		speedScale = math.Max(0.6, diff)
	}

	for range batch {
		enemy := EnemyShip{
			X:         float64(rand.Intn(760)),
			Y:         -30,
			VX:        (float64(rand.Intn(3) - 1)) * speedScale,
			VY:        (2 + float64(rand.Intn(2))) * speedScale,
			Width:     30,
			Height:    25,
			Active:    true,
			Health:    enemyHP,
			MaxHealth: enemyHP,
		}
		sg.enemies = append(sg.enemies, enemy)
		sg.spawnedEnemyCount++
	}
	sg.lastEnemy = time.Now()
}

func (sg *ShooterGame) spawnBoss() {
	// 简单的Boss：体积大、血量高、移动缓慢
	cfg := config.DefaultConfig()
	sg.boss = EnemyShip{
		X:         350,
		Y:         60,
		VX:        1.2,
		VY:        0.8,
		Width:     100,
		Height:    60,
		Active:    true,
		Health:    int(math.Ceil(60 * math.Max(1.0, 1.0+cfg.DiffHpLogK*math.Log10(math.Max(1.0, sg.difficultyMul))))),
		MaxHealth: int(math.Ceil(60 * math.Max(1.0, 1.0+cfg.DiffHpLogK*math.Log10(math.Max(1.0, sg.difficultyMul))))),
	}
	sg.bossSpawned = true
	sg.bossActive = true
}

func (sg *ShooterGame) countActiveEnemies() int {
	cnt := 0
	for _, e := range sg.enemies {
		if e.Active {
			cnt++
		}
	}
	return cnt
}
