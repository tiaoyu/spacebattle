package scenes

import (
	"fmt"

	"spacebattle/internal/config"
	"spacebattle/internal/ecs"
	"spacebattle/internal/ecs/components"
	"spacebattle/internal/ecs/systems"
	"spacebattle/internal/i18n"
	"spacebattle/internal/progress"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

// UpgradeScene 升级场景
type UpgradeScene struct {
	world         *ecs.World
	inputSystem   *systems.InputSystem
	menuSystem    *systems.MenuSystem
	playerOptions PlayerOptions
	selectedIndex int
}

// NewUpgradeScene 创建升级场景
func NewUpgradeScene(opts PlayerOptions) *UpgradeScene {
	world := ecs.NewWorld()

	// 预填上次保存的加点
	if last, err := progress.GetUpgrades(); err == nil {
		opts.ModFireRateHz = last.ModFireRateHz
		opts.ModBulletsPerShot = last.ModBulletsPerShot
		opts.ModPenetration = last.ModPenetration
		opts.ModSpreadDeltaDeg = last.ModSpreadDeltaDeg
		opts.ModBulletSpeed = last.ModBulletSpeed
		opts.ModBulletDamage = last.ModBulletDamage
		opts.ModBurstChance = last.ModBurstChance
		opts.ModEnableHoming = last.ModEnableHoming
		opts.ModTurnRateRad = last.ModTurnRateRad
	}

	scene := &UpgradeScene{
		world:         world,
		inputSystem:   systems.NewInputSystem(),
		menuSystem:    systems.NewMenuSystem(),
		playerOptions: opts,
		selectedIndex: 0,
	}

	// 获取当前功勋
	merits := progress.GetMerits()

	// 创建菜单状态
	menuState := world.ECS.World.Entry(world.ECS.World.Create(components.MenuState))
	components.MenuState.Set(menuState, &components.MenuStateData{
		SelectedIndex:   0,
		OptionCount:     9, // 9 个升级选项
		Confirmed:       false,
		AvailableMerits: merits,
		Upgrades:        make(map[string]int),
	})

	return scene
}

// Update 更新升级场景
func (s *UpgradeScene) Update() error {
	s.inputSystem.Update(s.world.ECS.World)

	// 获取菜单状态
	var menuState *components.MenuStateData
	components.MenuState.Each(s.world.ECS.World, func(entry *donburi.Entry) {
		menuState = components.MenuState.Get(entry)
	})

	if menuState == nil {
		return nil
	}

	// 上下选择
	if s.inputSystem.IsGMUpPressed() {
		menuState.SelectedIndex--
		if menuState.SelectedIndex < 0 {
			menuState.SelectedIndex = 8
		}
	}
	if s.inputSystem.IsGMDownPressed() {
		menuState.SelectedIndex++
		if menuState.SelectedIndex > 8 {
			menuState.SelectedIndex = 0
		}
	}

	// 左右加点/减点
	if s.inputSystem.IsGMRightPressed() {
		if s.canIncrease(menuState.SelectedIndex) {
			cost := s.nextCost(menuState.SelectedIndex)
			if cost > 0 && progress.SpendMerits(cost) {
				s.increase(menuState.SelectedIndex)
				menuState.AvailableMerits = progress.GetMerits()
			}
		}
	}
	if s.inputSystem.IsGMLeftPressed() {
		refund := s.refundCost(menuState.SelectedIndex)
		if s.decrease(menuState.SelectedIndex) && refund > 0 {
			progress.AddMerits(refund)
			menuState.AvailableMerits = progress.GetMerits()
		}
	}

	// Enter 确认
	if s.inputSystem.IsConfirmed() {
		// 保存升级数据
		_ = progress.SaveUpgrades(progress.UpgradeData{
			ModFireRateHz:     s.playerOptions.ModFireRateHz,
			ModBulletsPerShot: s.playerOptions.ModBulletsPerShot,
			ModPenetration:    s.playerOptions.ModPenetration,
			ModSpreadDeltaDeg: s.playerOptions.ModSpreadDeltaDeg,
			ModBulletSpeed:    s.playerOptions.ModBulletSpeed,
			ModBulletDamage:   s.playerOptions.ModBulletDamage,
			ModBurstChance:    s.playerOptions.ModBurstChance,
			ModEnableHoming:   s.playerOptions.ModEnableHoming,
			ModTurnRateRad:    s.playerOptions.ModTurnRateRad,
		})
		menuState.Confirmed = true
	}

	return nil
}

// Draw 绘制升级场景
func (s *UpgradeScene) Draw(screen *ebiten.Image) {
	// 获取菜单状态
	var menuState *components.MenuStateData
	components.MenuState.Each(s.world.ECS.World, func(entry *donburi.Entry) {
		menuState = components.MenuState.Get(entry)
	})

	if menuState == nil {
		return
	}

	// 准备升级项信息
	items := s.buildUpgradeItems()

	// 使用详细渲染
	s.menuSystem.DrawUpgradeWithDetails(screen, menuState, items)
}

// buildUpgradeItems 构建升级项列表
func (s *UpgradeScene) buildUpgradeItems() []systems.UpgradeItem {
	items := []systems.UpgradeItem{
		{
			Key:   "upgrade.fire_rate",
			Value: formatFloat(s.playerOptions.ModFireRateHz, 1),
			Cost:  s.nextCost(0),
		},
		{
			Key:   "upgrade.bullets_per_shot",
			Value: formatInt(s.playerOptions.ModBulletsPerShot),
			Cost:  s.nextCost(1),
		},
		{
			Key:   "upgrade.penetration",
			Value: formatInt(s.playerOptions.ModPenetration),
			Cost:  s.nextCost(2),
		},
		{
			Key:   "upgrade.spread_narrow",
			Value: formatFloat(s.playerOptions.ModSpreadDeltaDeg, 0),
			Cost:  s.nextCost(3),
		},
		{
			Key:   "upgrade.bullet_speed",
			Value: formatFloat(s.playerOptions.ModBulletSpeed, 1),
			Cost:  s.nextCost(4),
		},
		{
			Key:   "upgrade.bullet_damage",
			Value: formatInt(s.playerOptions.ModBulletDamage),
			Cost:  s.nextCost(5),
		},
		{
			Key:   "upgrade.burst_chance",
			Value: formatFloat(s.playerOptions.ModBurstChance, 2),
			Cost:  s.nextCost(6),
		},
		{
			Key:   "upgrade.enable_homing",
			Value: formatBool(s.playerOptions.ModEnableHoming),
			Cost:  s.nextCost(7),
		},
		{
			Key:   "upgrade.turn_rate",
			Value: formatFloat(s.playerOptions.ModTurnRateRad, 2),
			Cost:  s.nextCost(8),
		},
	}
	return items
}

// 格式化辅助函数
func formatInt(v int) string {
	return fmt.Sprintf("%d", v)
}

func formatFloat(v float64, prec int) string {
	return fmt.Sprintf("%.*f", prec, v)
}

func formatBool(v bool) string {
	if v {
		return i18n.T("common.on")
	}
	return i18n.T("common.off")
}

// IsConfirmed 检查是否已确认
func (s *UpgradeScene) IsConfirmed() bool {
	var confirmed bool
	components.MenuState.Each(s.world.ECS.World, func(entry *donburi.Entry) {
		state := components.MenuState.Get(entry)
		confirmed = state.Confirmed
	})
	return confirmed
}

// GetOptions 获取更新后的玩家选项
func (s *UpgradeScene) GetOptions() PlayerOptions {
	return s.playerOptions
}

// 计算当前等级
func (s *UpgradeScene) levelOf(idx int) int {
	switch idx {
	case 0:
		return int(s.playerOptions.ModFireRateHz / 0.5)
	case 1:
		return s.playerOptions.ModBulletsPerShot
	case 2:
		return s.playerOptions.ModPenetration
	case 3:
		v := s.playerOptions.ModSpreadDeltaDeg
		if v < 0 {
			v = -v
		}
		return int(v / 5)
	case 4:
		return int(s.playerOptions.ModBulletSpeed / 0.5)
	case 5:
		return s.playerOptions.ModBulletDamage
	case 6:
		return int(s.playerOptions.ModBurstChance / 0.05)
	case 7:
		if s.playerOptions.ModEnableHoming {
			return 1
		}
		return 0
	case 8:
		return int(s.playerOptions.ModTurnRateRad / 0.02)
	default:
		return 0
	}
}

// 计算下一级成本
func (s *UpgradeScene) nextCost(idx int) int {
	cfg := config.DefaultConfig()
	bases := []int{
		cfg.UpgradeCostFireRate,
		cfg.UpgradeCostBulletsPerShot,
		cfg.UpgradeCostPenetration,
		cfg.UpgradeCostSpreadNarrow,
		cfg.UpgradeCostBulletSpeed,
		cfg.UpgradeCostBulletDamage,
		cfg.UpgradeCostBurstChance,
		cfg.UpgradeCostEnableHoming,
		cfg.UpgradeCostTurnRate,
	}
	level := s.levelOf(idx)
	if level < 0 {
		level = 0
	}
	return bases[idx] * (1 << level)
}

// 计算退款金额
func (s *UpgradeScene) refundCost(idx int) int {
	cfg := config.DefaultConfig()
	bases := []int{
		cfg.UpgradeCostFireRate,
		cfg.UpgradeCostBulletsPerShot,
		cfg.UpgradeCostPenetration,
		cfg.UpgradeCostSpreadNarrow,
		cfg.UpgradeCostBulletSpeed,
		cfg.UpgradeCostBulletDamage,
		cfg.UpgradeCostBurstChance,
		cfg.UpgradeCostEnableHoming,
		cfg.UpgradeCostTurnRate,
	}
	level := s.levelOf(idx) - 1
	if level < 0 {
		return 0
	}
	return bases[idx] * (1 << level)
}

// 增加属性
func (s *UpgradeScene) increase(idx int) {
	switch idx {
	case 0:
		s.playerOptions.ModFireRateHz += 0.5
	case 1:
		s.playerOptions.ModBulletsPerShot++
	case 2:
		s.playerOptions.ModPenetration++
	case 3:
		s.playerOptions.ModSpreadDeltaDeg += 5
	case 4:
		s.playerOptions.ModBulletSpeed += 0.5
	case 5:
		s.playerOptions.ModBulletDamage++
	case 6:
		s.playerOptions.ModBurstChance += 0.05
	case 7:
		s.playerOptions.ModEnableHoming = true
	case 8:
		s.playerOptions.ModTurnRateRad += 0.02
	}
}

// 减少属性
func (s *UpgradeScene) decrease(idx int) bool {
	switch idx {
	case 0:
		if s.playerOptions.ModFireRateHz >= 0.5 {
			s.playerOptions.ModFireRateHz -= 0.5
			return true
		}
	case 1:
		if s.playerOptions.ModBulletsPerShot > 0 {
			s.playerOptions.ModBulletsPerShot--
			return true
		}
	case 2:
		if s.playerOptions.ModPenetration > 0 {
			s.playerOptions.ModPenetration--
			return true
		}
	case 3:
		if s.playerOptions.ModSpreadDeltaDeg >= 5 {
			s.playerOptions.ModSpreadDeltaDeg -= 5
			return true
		}
	case 4:
		if s.playerOptions.ModBulletSpeed >= 0.5 {
			s.playerOptions.ModBulletSpeed -= 0.5
			return true
		}
	case 5:
		if s.playerOptions.ModBulletDamage > 0 {
			s.playerOptions.ModBulletDamage--
			return true
		}
	case 6:
		if s.playerOptions.ModBurstChance >= 0.05 {
			s.playerOptions.ModBurstChance -= 0.05
			return true
		}
	case 7:
		if s.playerOptions.ModEnableHoming {
			s.playerOptions.ModEnableHoming = false
			return true
		}
	case 8:
		if s.playerOptions.ModTurnRateRad >= 0.02 {
			s.playerOptions.ModTurnRateRad -= 0.02
			return true
		}
	}
	return false
}

// 检查是否可以增加
func (s *UpgradeScene) canIncrease(idx int) bool {
	cfg := config.DefaultConfig()
	switch idx {
	case 0:
		return (5.0 + s.playerOptions.ModFireRateHz + 0.5) <= cfg.MaxFireRateHz
	case 1:
		return (1 + s.playerOptions.ModBulletsPerShot + 1) <= cfg.MaxBulletsPerShot
	case 2:
		return (0 + s.playerOptions.ModPenetration + 1) <= cfg.MaxPenetration
	case 3:
		return (2 + s.playerOptions.ModSpreadDeltaDeg + 5) <= cfg.MaxSpreadDeg
	case 4:
		return (8.0 + s.playerOptions.ModBulletSpeed + 0.5) <= cfg.MaxBulletSpeed
	case 5:
		return true // 伤害没有上限
	case 6:
		return (0.0 + s.playerOptions.ModBurstChance + 0.05) <= cfg.MaxBurstChance
	case 7:
		return !s.playerOptions.ModEnableHoming
	case 8:
		return (0.01 + s.playerOptions.ModTurnRateRad + 0.02) <= cfg.MaxTurnRateRad
	default:
		return true
	}
}
