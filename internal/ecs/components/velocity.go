package components

import "github.com/yohamta/donburi"

// VelocityData 速度组件数据
type VelocityData struct {
	VX, VY float64
}

// Velocity 速度组件
var Velocity = donburi.NewComponentType[VelocityData]()

