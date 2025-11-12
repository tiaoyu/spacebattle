package menu

import (
	"fmt"
	"image/color"

	"spacebattle/internal/fonts"
	"spacebattle/internal/i18n"
	"spacebattle/internal/progress"
	"spacebattle/internal/scenes/battle"
	"spacebattle/internal/utils"

	"github.com/hajimehoshi/ebiten/v2"
)

// ShipTemplate 战机模板（最小可用）
type ShipTemplate struct {
	NameKey    string
	Speed      float64
	SizeScale  float64
	Lives      int
	PassiveKey string
}

// ShipSelectScene 战机选择场景
type ShipSelectScene struct {
	input     *utils.InputManager
	selected  int
	templates []ShipTemplate
	confirmed bool
}

func NewShipSelectScene() *ShipSelectScene {
	return &ShipSelectScene{
		input:    utils.NewInputManager(),
		selected: 0,
		templates: []ShipTemplate{
			{NameKey: "ship.alpha", Speed: 5.0, SizeScale: 1.0, Lives: 3, PassiveKey: "passive.none"},
			{NameKey: "ship.beta", Speed: 6.5, SizeScale: 1.0, Lives: 3, PassiveKey: "passive.speed"},
			{NameKey: "ship.gamma", Speed: 5.0, SizeScale: 0.8, Lives: 3, PassiveKey: "passive.small"},
			{NameKey: "ship.delta", Speed: 5.0, SizeScale: 1.0, Lives: 4, PassiveKey: "passive.life"},
		},
	}
}

// Confirmed 是否已确认选择
func (s *ShipSelectScene) Confirmed() bool { return s.confirmed }

// GetOptions 返回战斗需要的玩家选项
func (s *ShipSelectScene) GetOptions() battle.PlayerOptions {
	t := s.templates[s.selected]
	return battle.PlayerOptions{
		Speed:      t.Speed,
		SizeScale:  t.SizeScale,
		Lives:      t.Lives,
		PassiveKey: t.PassiveKey,
	}
}

// Update 处理输入
func (s *ShipSelectScene) Update() error {
	s.input.Update()

	if s.input.IsKeyJustPressed(ebiten.KeyArrowLeft) || s.input.IsKeyJustPressed(ebiten.KeyLeft) {
		s.selected--
		if s.selected < 0 {
			s.selected = len(s.templates) - 1
		}
	}
	if s.input.IsKeyJustPressed(ebiten.KeyArrowRight) || s.input.IsKeyJustPressed(ebiten.KeyRight) {
		s.selected++
		if s.selected >= len(s.templates) {
			s.selected = 0
		}
	}
	if s.input.IsKeyJustPressed(ebiten.KeyEnter) || s.input.IsKeyJustPressed(ebiten.KeySpace) {
		// 加载上次升级加点
		if up, err := progress.GetUpgrades(); err == nil {
			// 应用到即将开始的选项中（将在 UpgradeScene 可继续调整）
			// 这里只是标记，真正合并在进入战斗时以 opts 为准
			_ = up // 预留：如需在此预览
		}
		s.confirmed = true
	}
	return nil
}

// Draw 简单UI
func (s *ShipSelectScene) Draw(screen *ebiten.Image) {
	title := i18n.T("select.title")
	fonts.DrawTextCenteredLarge(screen, title, 0, 80, 800, color.White)

	hint := i18n.T("select.hint")
	fonts.DrawTextCentered(screen, hint, 0, 120, 800, color.RGBA{R: 200, G: 200, B: 200, A: 255})

	// 显示当前选中模板信息
	t := s.templates[s.selected]
	name := i18n.T(t.NameKey)
	fonts.DrawTextCentered(screen, fmt.Sprintf("[%s]", name), 0, 200, 800, color.White)

	y := 240
	fonts.DrawTextCentered(screen, fmt.Sprintf("%s: %.1f", i18n.T("select.speed"), t.Speed), 0, y, 800, color.White)
	y += 24
	fonts.DrawTextCentered(screen, fmt.Sprintf("%s: %.0f%%", i18n.T("select.size"), t.SizeScale*100), 0, y, 800, color.White)
	y += 24
	fonts.DrawTextCentered(screen, fmt.Sprintf("%s: %d", i18n.T("select.lives"), t.Lives), 0, y, 800, color.White)
	y += 24
	fonts.DrawTextCentered(screen, fmt.Sprintf("%s: %s", i18n.T("select.passive"), i18n.T(t.PassiveKey)), 0, y, 800, color.White)

	// 左右指示
	fonts.DrawTextCentered(screen, "<   >", 0, 320, 800, color.RGBA{R: 220, G: 220, B: 80, A: 255})
}
