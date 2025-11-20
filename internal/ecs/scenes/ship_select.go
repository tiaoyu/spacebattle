package scenes

import (
	"spacebattle/internal/ecs"
	"spacebattle/internal/ecs/components"
	"spacebattle/internal/ecs/systems"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

// ShipSelectScene 战机选择场景
type ShipSelectScene struct {
	world       *ecs.World
	inputSystem *systems.InputSystem
	menuSystem  *systems.MenuSystem
}

// NewShipSelectScene 创建战机选择场景
func NewShipSelectScene() *ShipSelectScene {
	world := ecs.NewWorld()

	scene := &ShipSelectScene{
		world:       world,
		inputSystem: systems.NewInputSystem(),
		menuSystem:  systems.NewMenuSystem(),
	}

	// 创建战机模板
	templates := []components.ShipTemplate{
		{NameKey: "ship.alpha", Speed: 5.0, SizeScale: 1.0, Lives: 3, PassiveKey: "passive.none"},
		{NameKey: "ship.beta", Speed: 6.5, SizeScale: 1.0, Lives: 3, PassiveKey: "passive.speed"},
		{NameKey: "ship.gamma", Speed: 5.0, SizeScale: 0.8, Lives: 3, PassiveKey: "passive.small"},
		{NameKey: "ship.delta", Speed: 5.0, SizeScale: 1.0, Lives: 4, PassiveKey: "passive.life"},
	}

	// 创建菜单状态
	menuState := world.ECS.World.Entry(world.ECS.World.Create(components.MenuState))
	components.MenuState.Set(menuState, &components.MenuStateData{
		SelectedIndex: 0,
		OptionCount:   len(templates),
		Confirmed:     false,
		ShipTemplates: templates,
	})

	return scene
}

// Update 更新战机选择场景
func (s *ShipSelectScene) Update() error {
	s.inputSystem.Update(s.world.ECS.World)
	s.inputSystem.ProcessMenuInput(s.world.ECS.World)
	return nil
}

// Draw 绘制战机选择场景
func (s *ShipSelectScene) Draw(screen *ebiten.Image) {
	s.menuSystem.DrawShipSelect(s.world.ECS.World, screen)
}

// IsConfirmed 检查是否已确认
func (s *ShipSelectScene) IsConfirmed() bool {
	var confirmed bool
	components.MenuState.Each(s.world.ECS.World, func(entry *donburi.Entry) {
		state := components.MenuState.Get(entry)
		confirmed = state.Confirmed
	})
	return confirmed
}

// GetOptions 获取玩家选项
func (s *ShipSelectScene) GetOptions() PlayerOptions {
	var opts PlayerOptions
	components.MenuState.Each(s.world.ECS.World, func(entry *donburi.Entry) {
		state := components.MenuState.Get(entry)
		if state.SelectedIndex < len(state.ShipTemplates) {
			t := state.ShipTemplates[state.SelectedIndex]
			opts = PlayerOptions{
				Speed:      t.Speed,
				SizeScale:  t.SizeScale,
				Lives:      t.Lives,
				PassiveKey: t.PassiveKey,
			}
		}
	})
	return opts
}

