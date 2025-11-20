package systems

import (
	"spacebattle/internal/ecs/components"
	"spacebattle/internal/ecs/tags"
	"spacebattle/internal/i18n"
	"spacebattle/internal/utils"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
)

// InputSystem 输入处理系统
type InputSystem struct {
	inputManager *utils.InputManager
}

// NewInputSystem 创建输入系统
func NewInputSystem() *InputSystem {
	return &InputSystem{
		inputManager: utils.NewInputManager(),
	}
}

// Update 更新输入
func (s *InputSystem) Update(w donburi.World) {
	s.inputManager.Update()
}

// ProcessPlayerInput 处理玩家输入（战斗场景）
func (s *InputSystem) ProcessPlayerInput(w donburi.World) {
	// 查找玩家实体
	playerQuery := query.NewQuery(
		filter.Contains(tags.Player, components.Position, components.Velocity, components.PlayerInput),
	)

	playerQuery.Each(w, func(entry *donburi.Entry) {
		pos := components.Position.Get(entry)
		vel := components.Velocity.Get(entry)
		input := components.PlayerInput.Get(entry)
		size := components.Size.Get(entry)

		// 重置速度
		vel.VX = 0
		vel.VY = 0

		// 处理移动输入
		if ebiten.IsKeyPressed(ebiten.KeyArrowUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
			vel.VY = -input.Speed
		}
		if ebiten.IsKeyPressed(ebiten.KeyArrowDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
			vel.VY = input.Speed
		}
		if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
			vel.VX = -input.Speed
		}
		if ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
			vel.VX = input.Speed
		}

		// 边界限制
		if pos.X < 0 {
			pos.X = 0
		}
		if pos.X+size.Width > 800 {
			pos.X = 800 - size.Width
		}
		if pos.Y < 0 {
			pos.Y = 0
		}
		if pos.Y+size.Height > 600 {
			pos.Y = 600 - size.Height
		}
	})
}

// IsFirePressed 检查是否按下射击键
func (s *InputSystem) IsFirePressed() bool {
	return ebiten.IsKeyPressed(ebiten.KeySpace)
}

// IsGMTogglePressed 检查是否按下 GM 面板切换键
func (s *InputSystem) IsGMTogglePressed() bool {
	return s.inputManager.IsKeyJustPressed(ebiten.KeyG)
}

// IsGMTabPressed 检查是否按下 GM 标签切换键
func (s *InputSystem) IsGMTabPressed() bool {
	return s.inputManager.IsKeyJustPressed(ebiten.KeyTab)
}

// IsGMUpPressed 检查是否按下 GM 向上键
func (s *InputSystem) IsGMUpPressed() bool {
	return s.inputManager.IsKeyJustPressed(ebiten.KeyArrowUp)
}

// IsGMDownPressed 检查是否按下 GM 向下键
func (s *InputSystem) IsGMDownPressed() bool {
	return s.inputManager.IsKeyJustPressed(ebiten.KeyArrowDown)
}

// IsGMLeftPressed 检查是否按下 GM 向左键
func (s *InputSystem) IsGMLeftPressed() bool {
	return s.inputManager.IsKeyJustPressed(ebiten.KeyArrowLeft)
}

// IsGMRightPressed 检查是否按下 GM 向右键
func (s *InputSystem) IsGMRightPressed() bool {
	return s.inputManager.IsKeyJustPressed(ebiten.KeyArrowRight)
}

// IsRestartPressed 检查是否按下重开键
func (s *InputSystem) IsRestartPressed() bool {
	return s.inputManager.IsKeyJustPressed(ebiten.KeyR)
}

// IsEscapePressed 检查是否按下 ESC 键
func (s *InputSystem) IsEscapePressed() bool {
	return s.inputManager.IsKeyJustPressed(ebiten.KeyEscape)
}

// ProcessMenuInput 处理菜单输入
func (s *InputSystem) ProcessMenuInput(w donburi.World) {
	// 查找菜单状态实体
	query.NewQuery(filter.Contains(components.MenuState)).Each(w, func(entry *donburi.Entry) {
		state := components.MenuState.Get(entry)

		if s.inputManager.IsKeyJustPressed(ebiten.KeyArrowUp) {
			state.SelectedIndex--
			if state.SelectedIndex < 0 {
				state.SelectedIndex = state.OptionCount - 1
			}
		}

		if s.inputManager.IsKeyJustPressed(ebiten.KeyArrowDown) {
			state.SelectedIndex++
			if state.SelectedIndex >= state.OptionCount {
				state.SelectedIndex = 0
			}
		}

		if s.inputManager.IsKeyJustPressed(ebiten.KeyArrowLeft) {
			// 用于战机选择等场景的左右导航
			state.SelectedIndex--
			if state.SelectedIndex < 0 {
				state.SelectedIndex = state.OptionCount - 1
			}
		}

		if s.inputManager.IsKeyJustPressed(ebiten.KeyArrowRight) {
			// 用于战机选择等场景的左右导航
			state.SelectedIndex++
			if state.SelectedIndex >= state.OptionCount {
				state.SelectedIndex = 0
			}
		}

		if s.inputManager.IsKeyJustPressed(ebiten.KeyEnter) || s.inputManager.IsKeyJustPressed(ebiten.KeySpace) {
			state.Confirmed = true
		}

		// 处理语言切换
		if s.inputManager.IsKeyJustPressed(ebiten.KeyL) {
			currentLang := i18n.GetCurrentLanguage()
			switch currentLang {
			case i18n.Chinese:
				i18n.SetLanguage(i18n.English)
			case i18n.English:
				i18n.SetLanguage(i18n.Russian)
			case i18n.Russian:
				i18n.SetLanguage(i18n.Chinese)
			}
		}
	})
}

// IsConfirmed 检查是否确认（用于菜单）
func (s *InputSystem) IsConfirmed() bool {
	return s.inputManager.IsKeyJustPressed(ebiten.KeyEnter) || s.inputManager.IsKeyJustPressed(ebiten.KeySpace)
}

