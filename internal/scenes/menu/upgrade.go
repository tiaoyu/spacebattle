package menu

import (
	"fmt"
	"image/color"

	"spacebattle/internal/config"
	"spacebattle/internal/fonts"
	"spacebattle/internal/i18n"
	"spacebattle/internal/progress"
	"spacebattle/internal/scenes/battle"
	"spacebattle/internal/utils"

	"github.com/hajimehoshi/ebiten/v2"
)

// UpgradeScene 战机升级（最小骨架）
type UpgradeScene struct {
	input     *utils.InputManager
	opts      battle.PlayerOptions
	confirmed bool
	sel       int
}

func NewUpgradeScene(opts battle.PlayerOptions) *UpgradeScene {
	// 预填上次保存的加点
	if last, err := progress.GetUpgrades(); err == nil {
		opts.ModFireRateHz = last.ModFireRateHz
		opts.ModBulletsPerShot = last.ModBulletsPerShot
		opts.ModPenetration = last.ModPenetration
		opts.ModBulletDamage = last.ModBulletDamage
		opts.ModSpreadDeltaDeg = last.ModSpreadDeltaDeg
		opts.ModBulletSpeed = last.ModBulletSpeed
		opts.ModBurstChance = last.ModBurstChance
		opts.ModEnableHoming = last.ModEnableHoming
		opts.ModTurnRateRad = last.ModTurnRateRad
	}

	return &UpgradeScene{
		input: utils.NewInputManager(),
		opts:  opts,
		sel:   0,
	}
}

// Confirmed 是否完成
func (u *UpgradeScene) Confirmed() bool { return u.confirmed }

// GetOptions 透传当前选项
func (u *UpgradeScene) GetOptions() battle.PlayerOptions { return u.opts }

func (u *UpgradeScene) Update() error {
	u.input.Update()
	// 选择加点项：0 射速、1 子弹数、2 穿透、3 散射、4 子弹速度、5 连发、6 启用追踪、7 转向
	if u.input.IsKeyJustPressed(ebiten.KeyUp) || u.input.IsKeyJustPressed(ebiten.KeyArrowUp) {
		u.sel--
		if u.sel < 0 {
			u.sel = 7
		}
	}
	if u.input.IsKeyJustPressed(ebiten.KeyDown) || u.input.IsKeyJustPressed(ebiten.KeyArrowDown) {
		u.sel++
		if u.sel > 7 {
			u.sel = 0
		}
	}
	// 左右用于加点/减点（使用功勋，逐级成本翻倍）
	if u.input.IsKeyJustPressed(ebiten.KeyRight) || u.input.IsKeyJustPressed(ebiten.KeyArrowRight) {
		if u.canInc(u.sel) {
			cost := u.nextCost(u.sel)
			if cost > 0 && progress.SpendMerits(cost) {
				u.inc(u.sel)
			}
		}
	}
	if u.input.IsKeyJustPressed(ebiten.KeyLeft) || u.input.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		refund := u.refundCost(u.sel)
		if u.dec(u.sel) && refund > 0 {
			progress.AddMerits(refund)
		}
	}
	// Enter 继续
	if u.input.IsKeyJustPressed(ebiten.KeyEnter) || u.input.IsKeyJustPressed(ebiten.KeySpace) {
		// 保存升级数据到SQLite
		_ = progress.SaveUpgrades(progress.UpgradeData{
			ModFireRateHz:     u.opts.ModFireRateHz,
			ModBulletsPerShot: u.opts.ModBulletsPerShot,
			ModPenetration:    u.opts.ModPenetration,
			ModSpreadDeltaDeg: u.opts.ModSpreadDeltaDeg,
			ModBulletSpeed:    u.opts.ModBulletSpeed,
			ModBulletDamage:   u.opts.ModBulletDamage,
			ModBurstChance:    u.opts.ModBurstChance,
			ModEnableHoming:   u.opts.ModEnableHoming,
			ModTurnRateRad:    u.opts.ModTurnRateRad,
		})
		u.confirmed = true
	}
	return nil
}

func (u *UpgradeScene) Draw(screen *ebiten.Image) {
	fonts.DrawTextCenteredLarge(screen, i18n.T("upgrade.title"), 0, 120, 800, color.White)
	fonts.DrawTextCentered(screen, i18n.T("upgrade.hint"), 0, 180, 800, color.RGBA{R: 200, G: 200, B: 200, A: 255})
	// 显示可用功勋
	fonts.DrawTextCentered(screen, i18n.T("common.merits")+": "+fmtInt(progress.GetMerits()), 0, 210, 800, color.White)

	// 动态成本显示：base * 2^(level)
	cfg := config.DefaultConfig()
	levels := []int{
		int(u.opts.ModFireRateHz / 0.5),
		u.opts.ModBulletsPerShot,
		u.opts.ModPenetration,
		int(absFloat(u.opts.ModSpreadDeltaDeg) / 5),
		int(u.opts.ModBulletSpeed / 0.5),
		int(u.opts.ModBurstChance / 0.05),
		ternaryInt(u.opts.ModEnableHoming, 1, 0),
		int(u.opts.ModTurnRateRad / 0.02),
	}
	bases := []int{cfg.UpgradeCostFireRate, cfg.UpgradeCostBulletsPerShot, cfg.UpgradeCostPenetration, cfg.UpgradeCostSpreadNarrow, cfg.UpgradeCostBulletSpeed, cfg.UpgradeCostBurstChance, cfg.UpgradeCostEnableHoming, cfg.UpgradeCostTurnRate}
	costs := make([]int, len(bases))
	for i := range bases {
		lvl := levels[i]
		if lvl < 0 {
			lvl = 0
		}
		costs[i] = bases[i] * (1 << lvl)
	}

	items := []struct {
		key   string
		value string
		cost  int
	}{
		{"upgrade.fire_rate", fmtFloat(u.opts.ModFireRateHz, 1), costs[0]},
		{"upgrade.bullets_per_shot", fmtInt(u.opts.ModBulletsPerShot), costs[1]},
		{"upgrade.penetration", fmtInt(u.opts.ModPenetration), costs[2]},
		{"upgrade.spread_narrow", fmtFloat(u.opts.ModSpreadDeltaDeg, 0), costs[3]},
		{"upgrade.bullet_speed", fmtFloat(u.opts.ModBulletSpeed, 1), costs[4]},
		{"upgrade.burst_chance", fmtFloat(u.opts.ModBurstChance, 2), costs[5]},
		{"upgrade.enable_homing", ternary(u.opts.ModEnableHoming, i18n.T("common.on"), i18n.T("common.off")), costs[6]},
		{"upgrade.turn_rate", fmtFloat(u.opts.ModTurnRateRad, 2), costs[7]},
	}
	y := 250
	for i, it := range items {
		prefix := "  "
		col := color.RGBA{R: 255, G: 255, B: 255, A: 255}
		if i == u.sel {
			prefix = "> "
			col = color.RGBA{R: 255, G: 255, B: 0, A: 255}
		}
		line := fmt.Sprintf("%s%s: %s  (%s: %d)", prefix, i18n.T(it.key), it.value, i18n.T("upgrade.cost"), it.cost)
		fonts.DrawTextCentered(screen, line, 0, y, 800, col)
		y += 28
	}
}

// 小工具
func fmtInt(v int) string                 { return fmt.Sprintf("%d", v) }
func fmtFloat(v float64, prec int) string { return fmt.Sprintf("%.*f", prec, v) }
func ternary(cond bool, a, b string) string {
	if cond {
		return a
	}
	return b
}
func ternaryInt(cond bool, a, b int) int {
	if cond {
		return a
	}
	return b
}
func absFloat(v float64) float64 {
	if v < 0 {
		return -v
	}
	return v
}

// 计算成本与等级
func (u *UpgradeScene) levelOf(idx int) int {
	switch idx {
	case 0:
		return int(u.opts.ModFireRateHz / 0.5)
	case 1:
		return u.opts.ModBulletsPerShot
	case 2:
		return u.opts.ModPenetration
	case 3:
		return int(absFloat(u.opts.ModSpreadDeltaDeg) / 5)
	case 4:
		return int(u.opts.ModBulletSpeed / 0.5)
	case 5:
		return int(u.opts.ModBurstChance / 0.05)
	case 6:
		if u.opts.ModEnableHoming {
			return 1
		}
		return 0
	case 7:
		return int(u.opts.ModTurnRateRad / 0.02)
	default:
		return 0
	}
}

func (u *UpgradeScene) nextCost(idx int) int {
	cfg := config.DefaultConfig()
	base := []int{cfg.UpgradeCostFireRate, cfg.UpgradeCostBulletsPerShot, cfg.UpgradeCostPenetration, cfg.UpgradeCostSpreadNarrow, cfg.UpgradeCostBulletSpeed, cfg.UpgradeCostBurstChance, cfg.UpgradeCostEnableHoming, cfg.UpgradeCostTurnRate}
	level := max(u.levelOf(idx), 0)
	return base[idx] * (1 << level)
}

func (u *UpgradeScene) refundCost(idx int) int {
	cfg := config.DefaultConfig()
	base := []int{cfg.UpgradeCostFireRate, cfg.UpgradeCostBulletsPerShot, cfg.UpgradeCostPenetration, cfg.UpgradeCostSpreadNarrow, cfg.UpgradeCostBulletSpeed, cfg.UpgradeCostBurstChance, cfg.UpgradeCostEnableHoming, cfg.UpgradeCostTurnRate}
	level := u.levelOf(idx) - 1
	if level < 0 {
		return 0
	}
	return base[idx] * (1 << level)
}

// 属性增减
func (u *UpgradeScene) inc(idx int) {
	switch idx {
	case 0:
		// 上限由配置控制：在应用处 clamp，这里只做增量
		u.opts.ModFireRateHz += 0.5
	case 1:
		u.opts.ModBulletsPerShot += 1
	case 2:
		u.opts.ModPenetration += 1
	case 3:
		u.opts.ModSpreadDeltaDeg += 5
	case 4:
		u.opts.ModBulletSpeed += 0.5
	case 5:
		u.opts.ModBurstChance += 0.001
	case 6:
		u.opts.ModEnableHoming = true
	case 7:
		u.opts.ModTurnRateRad += 0.02
	}
}

// canInc 判断是否未超过配置上限
func (u *UpgradeScene) canInc(idx int) bool {
	cfg := config.DefaultConfig()
	switch idx {
	case 0:
		// FireRateHz 基础 5.0 + ModFireRateHz，限制最大值
		return (5.0 + u.opts.ModFireRateHz + 0.5) <= cfg.MaxFireRateHz
	case 1:
		// BulletsPerShot 基础 1 + ModBulletsPerShot
		return (1 + u.opts.ModBulletsPerShot + 1) <= cfg.MaxBulletsPerShot
	case 2:
		// Penetration 基础 0 + ModPenetration
		return (0 + u.opts.ModPenetration + 1) <= cfg.MaxPenetration
	case 3:
		// SpreadDeg 基础 2 + ModSpreadDeltaDeg（可正负，正向增加按 +5）
		return (2 + u.opts.ModSpreadDeltaDeg + 5) <= cfg.MaxSpreadDeg
	case 4:
		// BulletSpeed 基础 8.0 + ModBulletSpeed
		return (8.0 + u.opts.ModBulletSpeed + 0.5) <= cfg.MaxBulletSpeed
	case 5:
		// BurstChance 基础 0.0 + ModBurstChance
		return (0.0 + u.opts.ModBurstChance + 0.001) <= cfg.MaxBurstChance
	case 6:
		// EnableHoming 布尔，无上限；若已启用不可继续
		return !u.opts.ModEnableHoming
	case 7:
		// TurnRate 基础 0.01 + ModTurnRateRad
		return (0.01 + u.opts.ModTurnRateRad + 0.02) <= cfg.MaxTurnRateRad
	default:
		return true
	}
}

func (u *UpgradeScene) dec(idx int) bool {
	switch idx {
	case 0:
		if u.opts.ModFireRateHz >= 0.5 {
			u.opts.ModFireRateHz -= 0.5
			return true
		}
	case 1:
		if u.opts.ModBulletsPerShot > 0 {
			u.opts.ModBulletsPerShot--
			return true
		}
	case 2:
		if u.opts.ModPenetration > 0 {
			u.opts.ModPenetration--
			return true
		}
	case 3:
		if u.opts.ModSpreadDeltaDeg >= 5 {
			u.opts.ModSpreadDeltaDeg -= 5
			return true
		}
	case 4:
		if u.opts.ModBulletSpeed >= 0.5 {
			u.opts.ModBulletSpeed -= 0.5
			return true
		}
	case 5:
		if u.opts.ModBurstChance >= 0.05 {
			u.opts.ModBurstChance -= 0.05
			return true
		}
	case 6:
		if u.opts.ModEnableHoming {
			u.opts.ModEnableHoming = false
			return true
		}
	case 7:
		if u.opts.ModTurnRateRad >= 0.02 {
			u.opts.ModTurnRateRad -= 0.02
			return true
		}
	}
	return false
}
