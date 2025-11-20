package components

import "github.com/yohamta/donburi"

// HealthData 生命值组件数据
type HealthData struct {
	Current int
	Max     int
}

// Health 生命值组件
var Health = donburi.NewComponentType[HealthData]()

