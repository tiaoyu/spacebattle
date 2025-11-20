package ecs

import (
	"math"
	"math/rand"
	"time"

	"spacebattle/internal/config"
	"spacebattle/internal/ecs/components"
	"spacebattle/internal/ecs/tags"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

// World ECS World 包装器
type World struct {
	*ecs.ECS
}

// NewWorld 创建新的 ECS World
func NewWorld() *World {
	return &World{
		ECS: ecs.NewECS(donburi.NewWorld()),
	}
}

// CreatePlayer 创建玩家实体
func (w *World) CreatePlayer(x, y, width, height, speed float64, fireConfig components.FireSkillData) *donburi.Entry {
	cfg := config.DefaultConfig()
	player := w.ECS.World.Entry(w.ECS.World.Create(
		tags.Player,
		components.Position,
		components.Velocity,
		components.Size,
		components.PlayerInput,
		components.FireSkill,
		components.Health,
		components.ShipAbility,
		components.Sprite,
	))

	components.Position.Set(player, &components.PositionData{X: x, Y: y})
	components.Velocity.Set(player, &components.VelocityData{VX: 0, VY: 0})
	components.Size.Set(player, &components.SizeData{Width: width, Height: height})
	components.PlayerInput.Set(player, &components.PlayerInputData{Speed: speed})
	components.FireSkill.Set(player, &fireConfig)
	components.Sprite.Set(player, &components.SpriteData{
		Color: cfg.PlayerColor,
		Shape: "rect",
	})

	return player
}

// CreateBullet 创建子弹实体
func (w *World) CreateBullet(x, y, vx, vy, speed float64, damage, penetration int, homing bool, homingTurnRate float64) *donburi.Entry {
	cfg := config.DefaultConfig()
	comps := []donburi.IComponentType{
		tags.Bullet,
		components.Position,
		components.Velocity,
		components.Size,
		components.Sprite,
		components.Damage,
	}

	if penetration > 0 {
		comps = append(comps, components.Penetration)
	}

	if homing {
		comps = append(comps, components.Homing)
	}

	bullet := w.ECS.World.Entry(w.ECS.World.Create(comps...))

	components.Position.Set(bullet, &components.PositionData{X: x, Y: y})
	components.Velocity.Set(bullet, &components.VelocityData{VX: vx, VY: vy})
	components.Size.Set(bullet, &components.SizeData{Width: 4, Height: 10})
	components.Sprite.Set(bullet, &components.SpriteData{
		Color: cfg.BulletColor,
		Shape: "rect",
	})
	components.Damage.Set(bullet, &components.DamageData{Value: damage})

	if penetration > 0 {
		components.Penetration.Set(bullet, &components.PenetrationData{Remaining: penetration})
	}

	if homing {
		components.Homing.Set(bullet, &components.HomingData{
			TurnRate:         homingTurnRate,
			Speed:            speed,
			TargetEntity:     0, // 初始无目标，系统会自动分配
			LastRetargetTime: time.Now(),
			RetargetInterval: 2 * time.Second, // 2秒重新评估一次目标
		})
	}

	return bullet
}

// CreateEnemy 创建敌机实体
func (w *World) CreateEnemy(x, y, vx, vy, width, height float64, health int) *donburi.Entry {
	enemy := w.ECS.World.Entry(w.ECS.World.Create(
		tags.Enemy,
		components.Position,
		components.Velocity,
		components.Size,
		components.Health,
		components.Sprite,
	))

	components.Position.Set(enemy, &components.PositionData{X: x, Y: y})
	components.Velocity.Set(enemy, &components.VelocityData{VX: vx, VY: vy})
	components.Size.Set(enemy, &components.SizeData{Width: width, Height: height})
	components.Health.Set(enemy, &components.HealthData{Current: health, Max: health})
	components.Sprite.Set(enemy, &components.SpriteData{
		Color: config.DefaultConfig().EnemyColor,
		Shape: "rect",
	})

	return enemy
}

// CreateBoss 创建 Boss 实体
func (w *World) CreateBoss(x, y, vx, vy, width, height float64, health int) *donburi.Entry {
	boss := w.ECS.World.Entry(w.ECS.World.Create(
		tags.Boss,
		components.Position,
		components.Velocity,
		components.Size,
		components.Health,
		components.Sprite,
	))

	components.Position.Set(boss, &components.PositionData{X: x, Y: y})
	components.Velocity.Set(boss, &components.VelocityData{VX: vx, VY: vy})
	components.Size.Set(boss, &components.SizeData{Width: width, Height: height})
	components.Health.Set(boss, &components.HealthData{Current: health, Max: health})
	components.Sprite.Set(boss, &components.SpriteData{
		Color: config.DefaultConfig().BossColor,
		Shape: "rect",
	})

	return boss
}

// CreateExplosion 创建爆炸效果实体
func (w *World) CreateExplosion(x, y, maxRadius float64) *donburi.Entry {
	cfg := config.DefaultConfig()
	explosion := w.ECS.World.Entry(w.ECS.World.Create(
		tags.Explosion,
		components.Position,
		components.Lifetime,
		components.Sprite,
	))

	components.Position.Set(explosion, &components.PositionData{X: x, Y: y})
	components.Lifetime.Set(explosion, &components.LifetimeData{
		Timer:     0,
		MaxRadius: maxRadius,
		Radius:    0,
	})
	components.Sprite.Set(explosion, &components.SpriteData{
		Color: cfg.ExplosionColor,
		Shape: "circle",
	})

	return explosion
}

// CreateStar 创建星星背景实体
func (w *World) CreateStar(x, y, speed, size float64) *donburi.Entry {
	cfg := config.DefaultConfig()
	star := w.ECS.World.Entry(w.ECS.World.Create(
		tags.Star,
		components.Position,
		components.Star,
		components.Sprite,
	))

	components.Position.Set(star, &components.PositionData{X: x, Y: y})
	components.Star.Set(star, &components.StarData{Speed: speed, Size: size})
	components.Sprite.Set(star, &components.SpriteData{
		Color: cfg.StarColor,
		Shape: "circle",
	})

	return star
}

// CreateGameState 创建游戏状态实体（单例）
func (w *World) CreateGameState(lives int, difficultyMul float64) *donburi.Entry {
	cfg := config.DefaultConfig()

	gameState := w.ECS.World.Entry(w.ECS.World.Create(components.GameState))

	components.GameState.Set(gameState, &components.GameStateData{
		Score:              0,
		Lives:              lives,
		GameOver:           false,
		Victory:            false,
		Settled:            false,
		KilledEnemyCount:   0,
		SpawnedCount:       0,
		BossSpawned:        false,
		BossKilled:         false,
		RewardCached:       0,
		DifficultyMul:      difficultyMul,
		StartTime:          time.Now(),
		TotalDuration:      cfg.TotalDuration,
		SmallPhaseDuration: cfg.SmallPhaseDuration,
		WaveLength:         cfg.WaveLength,
		WaveCount:          cfg.WaveCount,
		WaveIndex:          0,
		WaveMinIntervals:   cfg.WaveMinIntervals,
		MaxSimultaneous:    cfg.MaxSimultaneous,
		BatchSize:          cfg.BatchSize,
		LastEnemyTime:      time.Now(),
		GMOpen:             false,
		GMIndex:            0,
		GMTab:              0,
	})

	return gameState
}

// InitializeStars 初始化背景星星
func (w *World) InitializeStars(count int) {
	for range count {
		x := rand.Float64() * 800
		y := rand.Float64() * 600
		speed := 1 + rand.Float64()*2 // 1-3
		size := 1 + rand.Float64()    // 1-2
		w.CreateStar(x, y, speed, size)
	}
}

// ComputeShotDelay 计算射击冷却时间
func ComputeShotDelay(fireRateHz float64) time.Duration {
	if fireRateHz <= 0 {
		return 300 * time.Millisecond
	}
	perShotSeconds := 1.0 / fireRateHz
	return time.Duration(perShotSeconds * float64(time.Second))
}

// ClampFloat64 限制浮点数范围
func ClampFloat64(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// WrapAngle 角度归一化到 [-π, π]
func WrapAngle(angle float64) float64 {
	for angle > math.Pi {
		angle -= 2 * math.Pi
	}
	for angle < -math.Pi {
		angle += 2 * math.Pi
	}
	return angle
}

// CreateEnemyShooter 创建射击型敌机
func (w *World) CreateEnemyShooter(x, y, vx, vy, width, height float64, health int) *donburi.Entry {
	cfg := config.DefaultConfig()
	enemy := w.ECS.World.Entry(w.ECS.World.Create(
		tags.EnemyShooter,
		components.Position,
		components.Velocity,
		components.Size,
		components.Health,
		components.Sprite,
		components.EnemyAI,
	))

	components.Position.Set(enemy, &components.PositionData{X: x, Y: y})
	components.Velocity.Set(enemy, &components.VelocityData{VX: vx, VY: vy})
	components.Size.Set(enemy, &components.SizeData{Width: width, Height: height})
	components.Health.Set(enemy, &components.HealthData{Current: health, Max: health})
	components.Sprite.Set(enemy, &components.SpriteData{
		Color: cfg.EnemyShooterColor,
		Shape: "rect",
	})
	components.EnemyAI.Set(enemy, &components.EnemyAIData{
		EnemyType:     "shooter",
		ShootInterval: cfg.EnemyShooterShootInterval,
		LastShotTime:  time.Now(),
	})

	return enemy
}

// CreateEnemyZigzag 创建之字型敌机
func (w *World) CreateEnemyZigzag(x, y, vx, vy, width, height float64, health int) *donburi.Entry {
	cfg := config.DefaultConfig()
	enemy := w.ECS.World.Entry(w.ECS.World.Create(
		tags.EnemyZigzag,
		components.Position,
		components.Velocity,
		components.Size,
		components.Health,
		components.Sprite,
		components.EnemyAI,
	))

	components.Position.Set(enemy, &components.PositionData{X: x, Y: y})
	components.Velocity.Set(enemy, &components.VelocityData{VX: vx, VY: vy})
	components.Size.Set(enemy, &components.SizeData{Width: width, Height: height})
	components.Health.Set(enemy, &components.HealthData{Current: health, Max: health})
	components.Sprite.Set(enemy, &components.SpriteData{
		Color: cfg.EnemyZigzagColor,
		Shape: "rect",
	})
	components.EnemyAI.Set(enemy, &components.EnemyAIData{
		EnemyType:    "zigzag",
		ZigzagPhase:  rand.Float64() * 2 * math.Pi,
		ZigzagSpeed:  2 * math.Pi / cfg.EnemyZigzagPeriod,
		ZigzagPeriod: cfg.EnemyZigzagPeriod,
	})

	return enemy
}

// CreateEnemyTank 创建肉盾型敌机
func (w *World) CreateEnemyTank(x, y, vx, vy, width, height float64, health int) *donburi.Entry {
	enemy := w.ECS.World.Entry(w.ECS.World.Create(
		tags.EnemyTank,
		components.Position,
		components.Velocity,
		components.Size,
		components.Health,
		components.Sprite,
		components.EnemyAI,
	))

	components.Position.Set(enemy, &components.PositionData{X: x, Y: y})
	components.Velocity.Set(enemy, &components.VelocityData{VX: vx, VY: vy})
	components.Size.Set(enemy, &components.SizeData{Width: width, Height: height})
	components.Health.Set(enemy, &components.HealthData{Current: health, Max: health})
	components.Sprite.Set(enemy, &components.SpriteData{
		Color: config.DefaultConfig().EnemyTankColor,
		Shape: "rect",
	})
	components.EnemyAI.Set(enemy, &components.EnemyAIData{
		EnemyType: "tank",
	})

	return enemy
}

// CreateEnemyBullet 创建敌机子弹
func (w *World) CreateEnemyBullet(x, y, vx, vy float64) *donburi.Entry {
	cfg := config.DefaultConfig()
	bullet := w.ECS.World.Entry(w.ECS.World.Create(
		tags.EnemyBullet,
		components.Position,
		components.Velocity,
		components.Size,
		components.Sprite,
	))

	components.Position.Set(bullet, &components.PositionData{X: x, Y: y})
	components.Velocity.Set(bullet, &components.VelocityData{VX: vx, VY: vy})
	components.Size.Set(bullet, &components.SizeData{Width: 5, Height: 5})
	components.Sprite.Set(bullet, &components.SpriteData{
		Color: cfg.EnemyBulletColor,
		Shape: "circle",
	})

	return bullet
}

// CreateParticle 创建粒子
func (w *World) CreateParticle(x, y, vx, vy, size float64, r, g, b uint8) *donburi.Entry {
	cfg := config.DefaultConfig()
	particle := w.ECS.World.Entry(w.ECS.World.Create(
		tags.Particle,
		components.Position,
		components.Particle,
	))

	components.Position.Set(particle, &components.PositionData{X: x, Y: y})
	components.Particle.Set(particle, &components.ParticleData{
		VX:        vx,
		VY:        vy,
		Life:      cfg.ParticleLifetime,
		MaxLife:   cfg.ParticleLifetime,
		Size:      size,
		DecayRate: 1.0 / cfg.ParticleLifetime,
		ColorR:    r,
		ColorG:    g,
		ColorB:    b,
		Alpha:     255,
	})

	return particle
}

// CreateScreenShake 创建屏幕震动效果
func (w *World) CreateScreenShake(intensity, duration float64) *donburi.Entry {
	shake := w.ECS.World.Entry(w.ECS.World.Create(
		components.ScreenShake,
	))

	components.ScreenShake.Set(shake, &components.ScreenShakeData{
		Intensity: intensity,
		Duration:  duration,
		Elapsed:   0,
		OffsetX:   0,
		OffsetY:   0,
	})

	return shake
}
