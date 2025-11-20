package components

import "github.com/yohamta/donburi"

// ScreenShakeData 屏幕震动数据
type ScreenShakeData struct {
	Intensity float64 // 震动强度（像素）
	Duration  float64 // 持续时间（秒）
	Elapsed   float64 // 已过时间（秒）
	OffsetX   float64 // 当前X轴偏移
	OffsetY   float64 // 当前Y轴偏移
}

// ScreenShake 屏幕震动组件
var ScreenShake = donburi.NewComponentType[ScreenShakeData]()

