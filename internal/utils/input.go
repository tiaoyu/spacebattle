package utils

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// InputManager 输入管理器
type InputManager struct {
	keys map[ebiten.Key]bool
}

// NewInputManager 创建输入管理器
func NewInputManager() *InputManager {
	return &InputManager{
		keys: make(map[ebiten.Key]bool),
	}
}

// Update 更新输入状态
func (im *InputManager) Update() {
	// TODO: 更新按键状态
}

// IsKeyPressed 检查按键是否被按下
func (im *InputManager) IsKeyPressed(key ebiten.Key) bool {
	return ebiten.IsKeyPressed(key)
}

// IsKeyJustPressed 检查按键是否刚被按下
func (im *InputManager) IsKeyJustPressed(key ebiten.Key) bool {
	return inpututil.IsKeyJustPressed(key)
}

// IsKeyJustReleased 检查按键是否刚被释放
func (im *InputManager) IsKeyJustReleased(key ebiten.Key) bool {
	return inpututil.IsKeyJustReleased(key)
}
