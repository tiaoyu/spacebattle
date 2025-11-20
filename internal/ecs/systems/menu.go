package systems

import (
	"fmt"
	"image/color"

	"spacebattle/internal/config"
	"spacebattle/internal/ecs/components"
	"spacebattle/internal/fonts"
	"spacebattle/internal/i18n"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
)

// MenuSystem 菜单系统
type MenuSystem struct{}

// NewMenuSystem 创建菜单系统
func NewMenuSystem() *MenuSystem {
	return &MenuSystem{}
}

// DrawMainMenu 绘制主菜单
func (s *MenuSystem) DrawMainMenu(w donburi.World, screen *ebiten.Image) {
	cfg := config.DefaultConfig()
	// 绘制背景
	screen.Fill(cfg.UIBackgroundColor)

	// 获取菜单状态
	var menuState *components.MenuStateData
	query.NewQuery(filter.Contains(components.MenuState)).Each(w, func(entry *donburi.Entry) {
		menuState = components.MenuState.Get(entry)
	})

	if menuState == nil {
		return
	}

	// 绘制标题
	fonts.DrawTextCenteredLarge(screen, i18n.T("menu.title"), 0, 150, 800, color.White)
	fonts.DrawTextCentered(screen, "===================", 0, 180, 800, cfg.UIHintColor)

	// 绘制菜单选项
	options := []string{i18n.T("menu.start"), i18n.T("menu.exit")}
	for i, option := range options {
		y := 250 + i*50
		if i == menuState.SelectedIndex {
			fonts.DrawText(screen, "> "+option, 300, y, cfg.UIHighlightColor)
		} else {
			fonts.DrawText(screen, "  "+option, 300, y, color.White)
		}
	}

	// 绘制说明
	fonts.DrawTextCentered(screen, i18n.T("menu.instructions"), 0, 500, 800, cfg.UIHintColor)
	fonts.DrawTextCentered(screen, i18n.T("menu.esc_hint"), 0, 530, 800, cfg.UIHintColor)

	// 绘制语言切换提示
	langText := "Language: " + string(i18n.GetCurrentLanguage()) + " (Press L to switch)"
	fonts.DrawTextCentered(screen, langText, 0, 560, 800, cfg.UIGreyTextColor)
}

// DrawShipSelect 绘制战机选择
func (s *MenuSystem) DrawShipSelect(w donburi.World, screen *ebiten.Image) {
	cfg := config.DefaultConfig()
	// 绘制背景
	screen.Fill(cfg.UIBackgroundColor)

	// 获取菜单状态
	var menuState *components.MenuStateData
	query.NewQuery(filter.Contains(components.MenuState)).Each(w, func(entry *donburi.Entry) {
		menuState = components.MenuState.Get(entry)
	})

	if menuState == nil || len(menuState.ShipTemplates) == 0 {
		return
	}

	// 绘制标题
	title := i18n.T("select.title")
	fonts.DrawTextCenteredLarge(screen, title, 0, 80, 800, color.White)

	hint := i18n.T("select.hint")
	fonts.DrawTextCentered(screen, hint, 0, 120, 800, cfg.UIHintColor)

	// 显示当前选中模板信息
	template := menuState.ShipTemplates[menuState.SelectedIndex]
	name := i18n.T(template.NameKey)
	fonts.DrawTextCentered(screen, fmt.Sprintf("[%s]", name), 0, 200, 800, color.White)

	y := 240
	fonts.DrawTextCentered(screen, fmt.Sprintf("%s: %.1f", i18n.T("select.speed"), template.Speed), 0, y, 800, color.White)
	y += 24
	fonts.DrawTextCentered(screen, fmt.Sprintf("%s: %.0f%%", i18n.T("select.size"), template.SizeScale*100), 0, y, 800, color.White)
	y += 24
	fonts.DrawTextCentered(screen, fmt.Sprintf("%s: %d", i18n.T("select.lives"), template.Lives), 0, y, 800, color.White)
	y += 24
	fonts.DrawTextCentered(screen, fmt.Sprintf("%s: %s", i18n.T("select.passive"), i18n.T(template.PassiveKey)), 0, y, 800, color.White)

	// 左右指示
	fonts.DrawTextCentered(screen, "<   >", 0, 320, 800, cfg.UIHighlightColor)
}

// DrawUpgrade 绘制升级界面（需要从场景传递额外信息）
func (s *MenuSystem) DrawUpgrade(w donburi.World, screen *ebiten.Image) {
	cfg := config.DefaultConfig()
	// 绘制背景
	screen.Fill(cfg.UIBackgroundColor)

	// 获取菜单状态
	var menuState *components.MenuStateData
	query.NewQuery(filter.Contains(components.MenuState)).Each(w, func(entry *donburi.Entry) {
		menuState = components.MenuState.Get(entry)
	})

	if menuState == nil {
		return
	}

	// 绘制标题
	title := i18n.T("upgrade.title")
	fonts.DrawTextCenteredLarge(screen, title, 0, 80, 800, color.White)

	// 显示功勋
	meritText := fmt.Sprintf("%s: %d", i18n.T("upgrade.merits"), menuState.AvailableMerits)
	fonts.DrawTextCentered(screen, meritText, 0, 120, 800, cfg.UIMeritColor)

	// 提示文字
	fonts.DrawTextCentered(screen, i18n.T("upgrade.hint"), 0, 160, 800, cfg.UIHintColor)

	// 显示升级选项（简化版本，只显示选项名称）
	upgradeOptions := []string{
		"upgrade.fire_rate",
		"upgrade.bullets_per_shot",
		"upgrade.penetration",
		"upgrade.spread_narrow",
		"upgrade.bullet_speed",
		"upgrade.burst_chance",
		"upgrade.enable_homing",
		"upgrade.turn_rate",
	}

	y := 220
	for i, optKey := range upgradeOptions {
		optText := i18n.T(optKey)
		if i == menuState.SelectedIndex {
			fonts.DrawTextCentered(screen, "> "+optText, 0, y, 800, cfg.UIHighlightColor)
		} else {
			fonts.DrawTextCentered(screen, "  "+optText, 0, y, 800, color.White)
		}
		y += 30
	}

	// 控制提示
	fonts.DrawTextCentered(screen, i18n.T("upgrade.controls"), 0, 500, 800, cfg.UIHintColor)
	fonts.DrawTextCentered(screen, i18n.T("upgrade.skip"), 0, 530, 800, cfg.UIHintColor)
}

// DrawUpgradeWithDetails 绘制详细的升级界面（包含属性值和成本）
func (s *MenuSystem) DrawUpgradeWithDetails(screen *ebiten.Image, menuState *components.MenuStateData, items []UpgradeItem) {
	cfg := config.DefaultConfig()
	// 绘制背景
	screen.Fill(cfg.UIBackgroundColor)

	// 绘制标题
	title := i18n.T("upgrade.title")
	fonts.DrawTextCenteredLarge(screen, title, 0, 80, 800, color.White)

	// 显示功勋
	meritText := fmt.Sprintf("%s: %d", i18n.T("common.merits"), menuState.AvailableMerits)
	fonts.DrawTextCentered(screen, meritText, 0, 120, 800, cfg.UIMeritColor)

	// 提示文字
	fonts.DrawTextCentered(screen, i18n.T("upgrade.hint"), 0, 160, 800, cfg.UIHintColor)

	// 显示升级选项
	y := 220
	for i, item := range items {
		prefix := "  "
		col := cfg.UITextColor
		if i == menuState.SelectedIndex {
			prefix = "> "
			col = cfg.UIHighlightColor
		}

		line := fmt.Sprintf("%s%s: %s  (%s: %d)", prefix, i18n.T(item.Key), item.Value, i18n.T("upgrade.cost"), item.Cost)
		fonts.DrawTextCentered(screen, line, 0, y, 800, col)
		y += 28
	}

	// 控制提示
	fonts.DrawTextCentered(screen, i18n.T("common.confirm"), 0, 500, 800, cfg.UIHintColor)
}

// UpgradeItem 升级项信息
type UpgradeItem struct {
	Key   string
	Value string
	Cost  int
}

// DrawDeploy 绘制出征界面
func (s *MenuSystem) DrawDeploy(w donburi.World, screen *ebiten.Image) {
	cfg := config.DefaultConfig()
	// 绘制背景
	screen.Fill(cfg.UIBackgroundColor)

	// 获取菜单状态
	var menuState *components.MenuStateData
	query.NewQuery(filter.Contains(components.MenuState)).Each(w, func(entry *donburi.Entry) {
		menuState = components.MenuState.Get(entry)
	})

	if menuState == nil {
		return
	}

	// 绘制标题
	title := i18n.T("deploy.title")
	fonts.DrawTextCenteredLarge(screen, title, 0, 80, 800, color.White)

	// 显示难度
	diffText := fmt.Sprintf("%s: %.1fx", i18n.T("deploy.difficulty"), menuState.DifficultyMul)
	fonts.DrawTextCentered(screen, diffText, 0, 150, 800, color.White)

	// 难度选项
	difficulties := []string{
		i18n.T("deploy.easy"),
		i18n.T("deploy.normal"),
		i18n.T("deploy.hard"),
	}

	y := 220
	for i, diff := range difficulties {
		if i == menuState.SelectedIndex {
			fonts.DrawTextCentered(screen, "> "+diff, 0, y, 800, cfg.UIHighlightColor)
		} else {
			fonts.DrawTextCentered(screen, "  "+diff, 0, y, 800, color.White)
		}
		y += 40
	}

	// 提示
	fonts.DrawTextCentered(screen, i18n.T("deploy.hint"), 0, 450, 800, cfg.UIHintColor)
}

// DrawDeployWithDetails 绘制详细的出征界面（连续可调难度）
func (s *MenuSystem) DrawDeployWithDetails(screen *ebiten.Image, menuState *components.MenuStateData, difficulty float64, cost int) {
	cfg := config.DefaultConfig()
	// 绘制背景
	screen.Fill(cfg.UIBackgroundColor)

	// 绘制标题
	title := i18n.T("deploy.title")
	fonts.DrawTextCenteredLarge(screen, title, 0, 100, 800, color.White)

	// 提示
	fonts.DrawTextCentered(screen, i18n.T("deploy.hint"), 0, 160, 800, cfg.UIHintColor)

	// 显示当前功勋
	meritText := fmt.Sprintf("%s: %d", i18n.T("common.merits"), menuState.AvailableMerits)
	fonts.DrawTextCentered(screen, meritText, 0, 200, 800, cfg.UIMeritColor)

	// 显示当前难度与成本
	diffText := fmt.Sprintf("x%.2f  (%s: %d)", difficulty, i18n.T("upgrade.cost"), cost)
	var diffColor color.Color = color.White
	// 如果成本超过可用功勋，显示为红色警告
	if cost > menuState.AvailableMerits {
		diffColor = cfg.UIGameOverColor
		diffText = fmt.Sprintf("x%.2f  (%s: %d) - INSUFFICIENT MERITS!", difficulty, i18n.T("upgrade.cost"), cost)
	}
	fonts.DrawTextCenteredLarge(screen, diffText, 0, 260, 800, diffColor)

	// 控制提示
	fonts.DrawTextCentered(screen, "Right: Increase  |  Left: MAX Challenge", 0, 330, 800, cfg.UIHighlightColor)
	fonts.DrawTextCentered(screen, "(Difficulty capped by available merits)", 0, 355, 800, cfg.UIGreyTextColor)

	// 确认提示
	fonts.DrawTextCentered(screen, i18n.T("common.confirm"), 0, 500, 800, cfg.UIHintColor)
}
