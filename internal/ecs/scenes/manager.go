package scenes

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Scene 场景接口
type Scene interface {
	Update() error
	Draw(screen *ebiten.Image)
}

// SceneType 场景类型
type SceneType int

const (
	SceneTypeMainMenu SceneType = iota
	SceneTypeShipSelect
	SceneTypeUpgrade
	SceneTypeDeploy
	SceneTypeBattle
)

// SceneManager 场景管理器
type SceneManager struct {
	currentScene Scene
	sceneType    SceneType
}

// NewSceneManager 创建场景管理器
func NewSceneManager() *SceneManager {
	return &SceneManager{
		currentScene: NewMainMenuScene(),
		sceneType:    SceneTypeMainMenu,
	}
}

// Update 更新当前场景
func (sm *SceneManager) Update() error {
	if sm.currentScene == nil {
		return nil
	}

	err := sm.currentScene.Update()
	if err != nil {
		return err
	}

	// 处理场景切换逻辑
	sm.handleSceneTransition()

	return nil
}

// Draw 绘制当前场景
func (sm *SceneManager) Draw(screen *ebiten.Image) {
	if sm.currentScene != nil {
		sm.currentScene.Draw(screen)
	}
}

// handleSceneTransition 处理场景切换
func (sm *SceneManager) handleSceneTransition() {
	switch sm.sceneType {
	case SceneTypeMainMenu:
		if mainMenu, ok := sm.currentScene.(*MainMenuScene); ok {
			if mainMenu.IsConfirmed() {
				selectedOption := mainMenu.GetSelectedOption()
				if selectedOption == 0 {
					// 开始游戏 -> 战机选择
					sm.currentScene = NewShipSelectScene()
					sm.sceneType = SceneTypeShipSelect
				}
				// 退出选项由 ebiten.Termination 处理
			}
		}

	case SceneTypeShipSelect:
		if shipSelect, ok := sm.currentScene.(*ShipSelectScene); ok {
			if shipSelect.IsConfirmed() {
				opts := shipSelect.GetOptions()
				sm.currentScene = NewUpgradeScene(opts)
				sm.sceneType = SceneTypeUpgrade
			}
		}

	case SceneTypeUpgrade:
		if upgrade, ok := sm.currentScene.(*UpgradeScene); ok {
			if upgrade.IsConfirmed() {
				opts := upgrade.GetOptions()
				sm.currentScene = NewDeployScene(opts)
				sm.sceneType = SceneTypeDeploy
			}
		}

	case SceneTypeDeploy:
		if deploy, ok := sm.currentScene.(*DeployScene); ok {
			if deploy.IsConfirmed() {
				opts := deploy.GetOptions()
				sm.currentScene = NewBattleScene(opts)
				sm.sceneType = SceneTypeBattle
			}
		}

	case SceneTypeBattle:
		// 战斗场景中可以通过 ESC 返回主菜单
		// 这个逻辑可以在 Update 中单独处理
	}
}

// SwitchToMainMenu 切换到主菜单
func (sm *SceneManager) SwitchToMainMenu() {
	sm.currentScene = NewMainMenuScene()
	sm.sceneType = SceneTypeMainMenu
}

// GetCurrentSceneType 获取当前场景类型
func (sm *SceneManager) GetCurrentSceneType() SceneType {
	return sm.sceneType
}

