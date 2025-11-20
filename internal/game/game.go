package game

import (
	"spacebattle/internal/ecs/scenes"
	"spacebattle/internal/utils"

	"github.com/hajimehoshi/ebiten/v2"
)

// Game 结构体实现 ebiten.Game 接口
type Game struct {
	sceneManager *scenes.SceneManager
	input        *utils.InputManager
}

// NewGame 创建新的游戏实例
func NewGame() *Game {
	return &Game{
		sceneManager: scenes.NewSceneManager(),
		input:        utils.NewInputManager(),
	}
}

// Update 更新游戏逻辑
func (g *Game) Update() error {
	g.input.Update()

	// 检查ESC键返回主菜单（在战斗场景中）
	if g.input.IsKeyJustPressed(ebiten.KeyEscape) {
		if g.sceneManager.GetCurrentSceneType() == scenes.SceneTypeBattle {
			g.sceneManager.SwitchToMainMenu()
			return nil
		}
	}

	// 处理退出（在主菜单中）
	if g.sceneManager.GetCurrentSceneType() == scenes.SceneTypeMainMenu {
		if g.input.IsKeyJustPressed(ebiten.KeyEnter) || g.input.IsKeyJustPressed(ebiten.KeySpace) {
			// 这里需要检查是否选择了退出选项
			// 由于场景管理器封装了场景，我们需要另一种方式处理
			// 暂时保持简单实现
		}
	}

	return g.sceneManager.Update()
}

// Draw 绘制游戏画面
func (g *Game) Draw(screen *ebiten.Image) {
	g.sceneManager.Draw(screen)
}

// Layout 返回游戏窗口布局
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 800, 600
}
