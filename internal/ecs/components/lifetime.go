package components

import "github.com/yohamta/donburi"

// LifetimeData 生命周期数据
type LifetimeData struct {
	Timer     float64 // 当前计时器 (0.0 - 1.0)
	MaxRadius float64 // 最大半径（用于爆炸等效果）
	Radius    float64 // 当前半径
}

// Lifetime 生命周期组件
var Lifetime = donburi.NewComponentType[LifetimeData]()

