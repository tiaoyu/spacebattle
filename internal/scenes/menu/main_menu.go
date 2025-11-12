package menu

import (
	"image/color"

	"spacebattle/internal/fonts"
	"spacebattle/internal/i18n"
	"spacebattle/internal/utils"

	"github.com/hajimehoshi/ebiten/v2"
)

// MainMenuScene 主菜单场景
type MainMenuScene struct {
	// 场景状态
	selectedOption int
	options        []string
	input          *utils.InputManager
}

// NewMainMenuScene 创建主菜单场景
func NewMainMenuScene() *MainMenuScene {
	return &MainMenuScene{
		selectedOption: 0,
		options:        []string{i18n.T("menu.start"), i18n.T("menu.exit")},
		input:          utils.NewInputManager(),
	}
}

// GetSelectedOption 获取当前选中的选项
func (m *MainMenuScene) GetSelectedOption() int {
	return m.selectedOption
}

// Update 更新主菜单逻辑
func (m *MainMenuScene) Update() error {
	m.input.Update()

	// 处理键盘输入
	if m.input.IsKeyJustPressed(ebiten.KeyArrowUp) {
		m.selectedOption--
		if m.selectedOption < 0 {
			m.selectedOption = len(m.options) - 1
		}
	}

	if m.input.IsKeyJustPressed(ebiten.KeyArrowDown) {
		m.selectedOption++
		if m.selectedOption >= len(m.options) {
			m.selectedOption = 0
		}
	}

	// 处理语言切换 (L键)
	if m.input.IsKeyJustPressed(ebiten.KeyL) {
		currentLang := i18n.GetCurrentLanguage()
		switch currentLang {
		case i18n.Chinese:
			i18n.SetLanguage(i18n.English)
		case i18n.English:
			i18n.SetLanguage(i18n.Russian)
		case i18n.Russian:
			i18n.SetLanguage(i18n.Chinese)
		}
		// 重新加载菜单选项
		m.options = []string{i18n.T("menu.start"), i18n.T("menu.exit")}
	}

	// 处理选择
	if m.input.IsKeyJustPressed(ebiten.KeyEnter) || m.input.IsKeyJustPressed(ebiten.KeySpace) {
		switch m.selectedOption {
		case 1: // 退出
			return ebiten.Termination
		}
		// 其他选项由外部处理
	}

	return nil
}

// Draw 绘制主菜单
func (m *MainMenuScene) Draw(screen *ebiten.Image) {
	// 绘制背景
	screen.Fill(color.RGBA{R: 20, G: 20, B: 60, A: 255})

	// 绘制标题
	fonts.DrawTextCenteredLarge(screen, i18n.T("menu.title"), 0, 150, 800, color.RGBA{R: 255, G: 255, B: 255, A: 255})
	fonts.DrawTextCentered(screen, "===================", 0, 180, 800, color.RGBA{R: 200, G: 200, B: 200, A: 255})

	// 绘制菜单选项
	for i, option := range m.options {
		y := 250 + i*50
		if i == m.selectedOption {
			fonts.DrawText(screen, "> "+option, 300, y, color.RGBA{R: 255, G: 255, B: 0, A: 255})
		} else {
			fonts.DrawText(screen, "  "+option, 300, y, color.RGBA{R: 255, G: 255, B: 255, A: 255})
		}
	}

	// 绘制说明
	fonts.DrawTextCentered(screen, i18n.T("menu.instructions"), 0, 500, 800, color.RGBA{R: 200, G: 200, B: 200, A: 255})
	fonts.DrawTextCentered(screen, i18n.T("menu.esc_hint"), 0, 530, 800, color.RGBA{R: 200, G: 200, B: 200, A: 255})

	// 绘制语言切换提示
	langText := "Language: " + string(i18n.GetCurrentLanguage()) + " (Press L to switch)"
	fonts.DrawTextCentered(screen, langText, 0, 560, 800, color.RGBA{R: 150, G: 150, B: 150, A: 255})
}
