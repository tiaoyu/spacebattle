package scenes

import (
	"time"

	"spacebattle/internal/balance"
	"spacebattle/internal/ecs"
	"spacebattle/internal/ecs/components"
	"spacebattle/internal/ecs/systems"
	"spacebattle/internal/progress"
	"spacebattle/internal/sound"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

// BattleScene 战斗场景
type BattleScene struct {
	world             *ecs.World
	inputSystem       *systems.InputSystem
	shipAbilitySystem *systems.ShipAbilitySystem
	enemyAISystem     *systems.EnemyAISystem
	movementSystem    *systems.MovementSystem
	fireSystem        *systems.FireSystem
	homingSystem      *systems.HomingSystem
	collisionSystem   *systems.CollisionSystem
	spawnSystem       *systems.SpawnSystem
	lifetimeSystem    *systems.LifetimeSystem
	particleSystem    *systems.ParticleSystem
	shakeSystem       *systems.ScreenShakeSystem
	renderSystem      *systems.RenderSystem
	initialOptions    PlayerOptions
}

// PlayerOptions 玩家配置选项
type PlayerOptions struct {
	Speed                float64
	SizeScale            float64
	Lives                int
	DifficultyMultiplier float64
	PassiveKey           string
	// 升级加成
	ModFireRateHz     float64
	ModBulletsPerShot int
	ModPenetration    int
	ModSpreadDeltaDeg float64
	ModBulletSpeed    float64
	ModBulletDamage   int
	ModBurstChance    float64
	ModEnableHoming   bool
	ModTurnRateRad    float64
}

// NewBattleScene 创建战斗场景
func NewBattleScene(opts PlayerOptions) *BattleScene {
	// 初始化音频
	sound.Init()

	// 创建 ECS World
	world := ecs.NewWorld()

	// 创建系统（按执行顺序）
	shipAbilitySystem := systems.NewShipAbilitySystem()
	enemyAISystem := systems.NewEnemyAISystem(world)
	particleSystem := systems.NewParticleSystem(world)
	shakeSystem := systems.NewScreenShakeSystem(world)
	
	scene := &BattleScene{
		world:             world,
		inputSystem:       systems.NewInputSystem(),
		shipAbilitySystem: shipAbilitySystem,
		enemyAISystem:     enemyAISystem,
		movementSystem:    systems.NewMovementSystem(),
		fireSystem:        systems.NewFireSystem(world),
		homingSystem:      systems.NewHomingSystem(),
		collisionSystem:   systems.NewCollisionSystem(world, shipAbilitySystem, particleSystem, shakeSystem),
		spawnSystem:       systems.NewSpawnSystem(world),
		lifetimeSystem:    systems.NewLifetimeSystem(),
		particleSystem:    particleSystem,
		shakeSystem:       shakeSystem,
		renderSystem:      systems.NewRenderSystem(),
		initialOptions:    opts,
	}

	// 初始化场景
	scene.initialize(opts)

	return scene
}

// initialize 初始化场景
func (s *BattleScene) initialize(opts PlayerOptions) {
	// 应用默认值
	if opts.Speed == 0 {
		opts.Speed = 5.0
	}
	if opts.SizeScale == 0 {
		opts.SizeScale = 1.0
	}
	if opts.Lives == 0 {
		opts.Lives = 3
	}
	if opts.DifficultyMultiplier == 0 {
		opts.DifficultyMultiplier = 1.0
	}

	// 计算玩家尺寸
	baseWidth := 40.0
	baseHeight := 30.0
	playerWidth := baseWidth * opts.SizeScale
	playerHeight := baseHeight * opts.SizeScale

	// 应用被动技能
	switch opts.PassiveKey {
	case "passive.speed":
		opts.Speed *= 1.5
	case "passive.small":
		opts.SizeScale *= 0.8
		playerWidth = baseWidth * opts.SizeScale
		playerHeight = baseHeight * opts.SizeScale
	case "passive.life":
		opts.Lives++
	}

	// 配置射击技能
	fireConfig := components.FireSkillData{
		FireRateHz:        5.0 + opts.ModFireRateHz,
		BulletsPerShot:    1 + opts.ModBulletsPerShot,
		SpreadDeg:         2.0 + opts.ModSpreadDeltaDeg,
		BulletSpeed:       8.0 + opts.ModBulletSpeed,
		BulletDamage:      1 + opts.ModBulletDamage,
		BurstChance:       0.0 + opts.ModBurstChance,
		PenetrationCount:  0 + opts.ModPenetration,
		EnableHoming:      opts.ModEnableHoming,
		HomingTurnRateRad: 0.01 + opts.ModTurnRateRad,
		BurstInterval:     60 * time.Millisecond,
	}

	// 初始化射击技能
	systems.InitializeFireSkill(&fireConfig)

	// 创建玩家实体
	playerEntry := s.world.CreatePlayer(400, 500, playerWidth, playerHeight, opts.Speed, fireConfig)

	// 初始化生命值
	components.Health.SetValue(playerEntry, components.HealthData{
		Current: opts.Lives,
		Max:     opts.Lives,
	})

	// 添加战机被动技能
	abilityType := "harvest" // 默认Alpha
	if opts.PassiveKey != "" {
		switch opts.PassiveKey {
		case "alpha":
			abilityType = "harvest"
		case "beta":
			abilityType = "speed_frenzy"
		case "gamma":
			abilityType = "dodge_master"
		case "delta":
			abilityType = "energy_shield"
		}
	}
	
	// 初始化ShipAbility组件
	components.ShipAbility.SetValue(playerEntry, components.ShipAbilityData{
		AbilityType:   abilityType,
		ShieldMax:     opts.Lives, // Delta的护盾上限
		ShieldCurrent: opts.Lives,
	})

	// 创建游戏状态实体
	s.world.CreateGameState(opts.Lives, opts.DifficultyMultiplier)

	// 初始化背景星星
	s.world.InitializeStars(100)
}

// Update 更新战斗场景
func (s *BattleScene) Update() error {
	// 更新输入
	s.inputSystem.Update(s.world.ECS.World)

	// 获取游戏状态
	gameState := s.getGameState()
	if gameState == nil {
		return nil
	}

	// 如果游戏结束或胜利，处理结算和重开逻辑
	if gameState.GameOver || gameState.Victory {
		if s.inputSystem.IsRestartPressed() {
			// 重开游戏
			*s = *NewBattleScene(s.initialOptions)
			return nil
		}

		// 结算功勋（一次性）
		if !gameState.Settled {
			spawned := gameState.SpawnedCount
			if gameState.BossSpawned {
				spawned++
			}
			kills := gameState.KilledEnemyCount
			if gameState.BossKilled {
				kills++
			}
			elapsed := time.Since(gameState.StartTime)

			// 计算详细奖励
			breakdown := balance.ComputeDetailedReward(
				gameState.DifficultyMul,
				kills,
				spawned,
				elapsed,
				gameState.TotalDuration,
				gameState.Victory,
			)

			// 失败时给予基础奖励的 1/3 作为安慰奖
			if !gameState.Victory && breakdown.BaseReward > 0 {
				gameState.RewardCached = breakdown.BaseReward / 3
				breakdown.TotalReward = gameState.RewardCached
				breakdown.DifficultyBonus = 0
				breakdown.KillBonus = 0
				breakdown.SpeedBonus = 0
				breakdown.PerfectBonus = 0
				breakdown.BossBonus = 0
			} else {
				gameState.RewardCached = breakdown.TotalReward
			}

			// 保存奖励分解信息
			gameState.RewardBreakdown = components.RewardBreakdownData{
				BaseReward:       breakdown.BaseReward,
				DifficultyBonus:  breakdown.DifficultyBonus,
				KillBonus:        breakdown.KillBonus,
				SpeedBonus:       breakdown.SpeedBonus,
				PerfectBonus:     breakdown.PerfectBonus,
				BossBonus:        breakdown.BossBonus,
				TotalReward:      breakdown.TotalReward,
				PerformanceScore: breakdown.PerformanceScore,
			}

			if gameState.RewardCached > 0 {
				progress.AddMerits(gameState.RewardCached)
			}
			gameState.Settled = true
		}

		return nil
	}

	// 硬收束：总时长达到后若未胜利，直接结算为胜利
	if time.Since(gameState.StartTime) >= gameState.TotalDuration && !gameState.Victory {
		gameState.Victory = true
		return nil
	}

	// 处理玩家输入（移动）
	s.inputSystem.ProcessPlayerInput(s.world.ECS.World)

	// 更新战机被动技能状态
	dt := 1.0 / 60.0 // 假设60 FPS
	s.shipAbilitySystem.Update(s.world.ECS.World, dt)

	// 更新敌机AI行为
	s.enemyAISystem.Update(s.world.ECS.World, dt)

	// 更新移动系统
	s.movementSystem.Update(s.world.ECS.World)

	// 处理射击
	firePressed := s.inputSystem.IsFirePressed()
	s.fireSystem.Update(s.world.ECS.World, firePressed)

	// 更新追踪系统
	s.homingSystem.Update(s.world.ECS.World)

	// 更新碰撞检测（会触发粒子和震动）
	s.collisionSystem.Update(s.world.ECS.World)

	// 更新生成系统
	s.spawnSystem.Update(s.world.ECS.World)

	// 清理超出屏幕的实体
	s.movementSystem.CleanOutOfBounds(s.world.ECS.World)

	// 更新生命周期系统
	s.lifetimeSystem.Update(s.world.ECS.World)

	// 更新粒子系统
	s.particleSystem.Update(s.world.ECS.World, dt)

	// 更新屏幕震动系统
	s.shakeSystem.Update(s.world.ECS.World, dt)

	return nil
}

// Draw 绘制战斗场景
func (s *BattleScene) Draw(screen *ebiten.Image) {
	// 绘制主要场景
	s.renderSystem.Draw(s.world.ECS.World, screen)
	
	// 绘制粒子效果
	s.particleSystem.Draw(s.world.ECS.World, screen)
}

// getGameState 获取游戏状态
func (s *BattleScene) getGameState() *components.GameStateData {
	var gameState *components.GameStateData
	components.GameState.Each(s.world.ECS.World, func(entry *donburi.Entry) {
		gameState = components.GameState.Get(entry)
	})
	return gameState
}

// IsGameOver 检查游戏是否结束
func (s *BattleScene) IsGameOver() bool {
	gameState := s.getGameState()
	return gameState != nil && gameState.GameOver
}

// IsVictory 检查是否胜利
func (s *BattleScene) IsVictory() bool {
	gameState := s.getGameState()
	return gameState != nil && gameState.Victory
}

