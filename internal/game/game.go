package game

import (
	"spacebattle/internal/scenes/battle"
	"spacebattle/internal/scenes/menu"
	"spacebattle/internal/utils"

	"github.com/hajimehoshi/ebiten/v2"
)

// Game 结构体实现 ebiten.Game 接口
type Game struct {
	// 游戏状态
	scene Scene
	input *utils.InputManager
}

// Scene 接口定义游戏场景
type Scene interface {
	Update() error
	Draw(screen *ebiten.Image)
}

// NewGame 创建新的游戏实例
func NewGame() *Game {
	return &Game{
		// 初始化默认场景
		scene: menu.NewMainMenuScene(),
		input: utils.NewInputManager(),
	}
}

// Update 更新游戏逻辑
func (g *Game) Update() error {
	g.input.Update()

	// 检查ESC键返回主菜单
	if g.input.IsKeyJustPressed(ebiten.KeyEscape) {
		// 如果当前不是主菜单，则返回主菜单
		if _, ok := g.scene.(*menu.MainMenuScene); !ok {
			g.scene = menu.NewMainMenuScene()
			return nil
		}
	}

	// 处理主菜单的选择
	if mainMenu, ok := g.scene.(*menu.MainMenuScene); ok {
		if g.input.IsKeyJustPressed(ebiten.KeyEnter) || g.input.IsKeyJustPressed(ebiten.KeySpace) {
			switch mainMenu.GetSelectedOption() {
			case 0: // 太空射击游戏 → 战机选择
				g.scene = menu.NewShipSelectScene()
			case 1: // 退出
				return ebiten.Termination
			}
			return nil
		}
	}

	// 战机选择完成后进入升级
	if selectScene, ok := g.scene.(*menu.ShipSelectScene); ok {
		if selectScene.Confirmed() {
			opts := selectScene.GetOptions()
			g.scene = menu.NewUpgradeScene(opts)
			return nil
		}
	}

	// 升级完成后进入出征
	if upScene, ok := g.scene.(*menu.UpgradeScene); ok {
		if upScene.Confirmed() {
			opts := upScene.GetOptions()
			g.scene = menu.NewDeployScene(opts)
			return nil
		}
	}

	// 出征确认后进入战斗
	if depScene, ok := g.scene.(*menu.DeployScene); ok {
		if depScene.Confirmed() {
			opts := depScene.GetOptions()
			g.scene = battle.NewShooterGameWithOptions(opts)
			return nil
		}
	}

	return g.scene.Update()
}

// Draw 绘制游戏画面
func (g *Game) Draw(screen *ebiten.Image) {
	g.scene.Draw(screen)
}

// Layout 返回游戏窗口布局
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 800, 600
}
