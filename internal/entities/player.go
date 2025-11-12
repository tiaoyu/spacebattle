package entities

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Player 玩家实体
type Player struct {
	X, Y     float64 // 位置
	Width    float64 // 宽度
	Height   float64 // 高度
	Velocity float64 // 速度
}

// NewPlayer 创建新玩家
func NewPlayer(x, y float64) *Player {
	return &Player{
		X:        x,
		Y:        y,
		Width:    32,
		Height:   32,
		Velocity: 5.0,
	}
}

// Update 更新玩家状态
func (p *Player) Update() {
	// TODO: 处理玩家输入和移动逻辑
}

// Draw 绘制玩家
func (p *Player) Draw(screen *ebiten.Image) {
	// TODO: 绘制玩家精灵
	// 临时使用矩形表示
	// 这里需要实际的绘制代码
}
