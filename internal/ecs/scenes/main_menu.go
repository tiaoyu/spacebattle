package scenes

import (
	"spacebattle/internal/ecs"
	"spacebattle/internal/ecs/components"
	"spacebattle/internal/ecs/systems"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

// MainMenuScene 主菜单场景
type MainMenuScene struct {
	world       *ecs.World
	inputSystem *systems.InputSystem
	menuSystem  *systems.MenuSystem
}

// NewMainMenuScene 创建主菜单场景
func NewMainMenuScene() *MainMenuScene {
	world := ecs.NewWorld()

	scene := &MainMenuScene{
		world:       world,
		inputSystem: systems.NewInputSystem(),
		menuSystem:  systems.NewMenuSystem(),
	}

	// 创建菜单状态
	menuState := world.ECS.World.Entry(world.ECS.World.Create(components.MenuState))
	components.MenuState.Set(menuState, &components.MenuStateData{
		SelectedIndex: 0,
		OptionCount:   2, // 开始游戏、退出
		Confirmed:     false,
	})

	return scene
}

// Update 更新主菜单
func (s *MainMenuScene) Update() error {
	s.inputSystem.Update(s.world.ECS.World)
	s.inputSystem.ProcessMenuInput(s.world.ECS.World)
	return nil
}

// Draw 绘制主菜单
func (s *MainMenuScene) Draw(screen *ebiten.Image) {
	s.menuSystem.DrawMainMenu(s.world.ECS.World, screen)
}

// GetSelectedOption 获取选中的选项
func (s *MainMenuScene) GetSelectedOption() int {
	var selected int
	components.MenuState.Each(s.world.ECS.World, func(entry *donburi.Entry) {
		state := components.MenuState.Get(entry)
		selected = state.SelectedIndex
	})
	return selected
}

// IsConfirmed 检查是否已确认
func (s *MainMenuScene) IsConfirmed() bool {
	var confirmed bool
	components.MenuState.Each(s.world.ECS.World, func(entry *donburi.Entry) {
		state := components.MenuState.Get(entry)
		confirmed = state.Confirmed
	})
	return confirmed
}

