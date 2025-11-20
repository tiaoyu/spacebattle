package components

import "github.com/yohamta/donburi"

// PositionData 位置组件数据
type PositionData struct {
	X, Y float64
}

// Position 位置组件
var Position = donburi.NewComponentType[PositionData]()

