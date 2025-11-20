package components

import (
	"image/color"

	"github.com/yohamta/donburi"
)

// SpriteData 精灵渲染数据
type SpriteData struct {
	Color color.Color
	Shape string // "rect", "circle", etc.
}

// Sprite 精灵组件
var Sprite = donburi.NewComponentType[SpriteData]()

