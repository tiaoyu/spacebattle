package scenes

import (
	"time"

	"spacebattle/internal/balance"
	"spacebattle/internal/config"
	"spacebattle/internal/ecs"
	"spacebattle/internal/ecs/components"
	"spacebattle/internal/ecs/systems"
	"spacebattle/internal/progress"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

// DeployScene 出征场景（连续可调难度）
type DeployScene struct {
	world          *ecs.World
	inputSystem    *systems.InputSystem
	menuSystem     *systems.MenuSystem
	playerOptions  PlayerOptions
	difficulty     float64
	minDifficulty  float64
	maxDifficulty  float64
	leftHoldStart  time.Time
	rightHoldStart time.Time
	lastAdjust     time.Time
}

// NewDeployScene 创建出征场景
func NewDeployScene(opts PlayerOptions) *DeployScene {
	world := ecs.NewWorld()
	cfg := config.DefaultConfig()

	scene := &DeployScene{
		world:         world,
		inputSystem:   systems.NewInputSystem(),
		menuSystem:    systems.NewMenuSystem(),
		playerOptions: opts,
		difficulty:    cfg.DifficultyMin,
		minDifficulty: cfg.DifficultyMin,
		maxDifficulty: cfg.DifficultyMax,
	}

	// 创建菜单状态（用于显示）
	menuState := world.ECS.World.Entry(world.ECS.World.Create(components.MenuState))
	components.MenuState.Set(menuState, &components.MenuStateData{
		SelectedIndex:   0,
		OptionCount:     0,
		Confirmed:       false,
		DifficultyMul:   cfg.DifficultyMin,
		AvailableMerits: progress.GetMerits(),
	})

	return scene
}

// Update 更新出征场景
func (s *DeployScene) Update() error {
	s.inputSystem.Update(s.world.ECS.World)
	now := time.Now()

	// 获取当前功勋并计算可支付的最大难度
	currentMerits := progress.GetMerits()
	affordableMaxDiff := balance.MaxAffordableDifficulty(currentMerits)
	
	// 动态限制最大难度为可支付的难度
	effectiveMaxDiff := min(s.maxDifficulty, affordableMaxDiff)

	// 计算步进（根据当前难度动态调整）
	step := s.calculateStep()
	repeat := 150 * time.Millisecond

	// 右键（增加难度）
	if s.inputSystem.IsGMRightPressed() {
		s.rightHoldStart = now
		newDiff := clampFloat(s.difficulty+step, s.minDifficulty, effectiveMaxDiff)
		s.difficulty = newDiff
		s.lastAdjust = now
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		hold := now.Sub(s.rightHoldStart)
		if hold > 500*time.Millisecond {
			repeat = 30 * time.Millisecond
		}
		if now.Sub(s.lastAdjust) >= repeat {
			newDiff := clampFloat(s.difficulty+step, s.minDifficulty, effectiveMaxDiff)
			s.difficulty = newDiff
			s.lastAdjust = now
		}
	}

	// 左键（快速设置为最大可支付难度，作为"极限挑战"快捷键）
	if s.inputSystem.IsGMLeftPressed() {
		s.leftHoldStart = now
		// 如果当前难度不是最大可支付难度，直接跳到最大可支付难度
		if s.difficulty < effectiveMaxDiff {
			s.difficulty = effectiveMaxDiff
		} else {
			// 如果已经是最大难度，则重置为最小难度
			s.difficulty = s.minDifficulty
		}
		s.lastAdjust = now
	} else if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		// 持续按住左键时，在最小和最大可支付难度之间切换
		if now.Sub(s.lastAdjust) >= repeat {
			if s.difficulty < effectiveMaxDiff {
				s.difficulty = effectiveMaxDiff
			} else {
				s.difficulty = s.minDifficulty
			}
			s.lastAdjust = now
		}
	}

	// 更新菜单状态
	components.MenuState.Each(s.world.ECS.World, func(entry *donburi.Entry) {
		state := components.MenuState.Get(entry)
		state.DifficultyMul = s.difficulty
		state.AvailableMerits = progress.GetMerits()
	})

	// Enter 确认
	if s.inputSystem.IsConfirmed() {
		cost := balance.DifficultyCost(s.difficulty)
		var confirmed bool
		components.MenuState.Each(s.world.ECS.World, func(entry *donburi.Entry) {
			state := components.MenuState.Get(entry)
			if progress.SpendMerits(cost) {
				s.playerOptions.DifficultyMultiplier = s.difficulty
				state.Confirmed = true
				confirmed = true
			}
		})
		if confirmed {
			// 更新功勋显示
			components.MenuState.Each(s.world.ECS.World, func(entry *donburi.Entry) {
				state := components.MenuState.Get(entry)
				state.AvailableMerits = progress.GetMerits()
			})
		}
	}

	return nil
}

// Draw 绘制出征场景
func (s *DeployScene) Draw(screen *ebiten.Image) {
	// 获取菜单状态
	var menuState *components.MenuStateData
	components.MenuState.Each(s.world.ECS.World, func(entry *donburi.Entry) {
		menuState = components.MenuState.Get(entry)
	})

	if menuState == nil {
		return
	}

	// 计算成本
	cost := balance.DifficultyCost(s.difficulty)

	// 使用详细渲染
	s.menuSystem.DrawDeployWithDetails(screen, menuState, s.difficulty, cost)
}

// IsConfirmed 检查是否已确认
func (s *DeployScene) IsConfirmed() bool {
	var confirmed bool
	components.MenuState.Each(s.world.ECS.World, func(entry *donburi.Entry) {
		state := components.MenuState.Get(entry)
		confirmed = state.Confirmed
	})
	return confirmed
}

// GetOptions 获取最终的玩家选项
func (s *DeployScene) GetOptions() PlayerOptions {
	return s.playerOptions
}

// calculateStep 计算难度调整步进
func (s *DeployScene) calculateStep() float64 {
	v := s.difficulty
	switch {
	case v < 2:
		return 0.1
	case v < 10:
		return 0.5
	case v < 100:
		return 1
	case v < 1000:
		return 10
	case v < 10000:
		return 100
	case v < 100000:
		return 1000
	case v < 1000000:
		return 10000
	default:
		return 100000
	}
}

// clampFloat 限制浮点数范围
func clampFloat(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

