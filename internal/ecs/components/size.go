package components

import "github.com/yohamta/donburi"

// SizeData 尺寸组件数据
type SizeData struct {
	Width, Height float64
}

// Size 尺寸组件
var Size = donburi.NewComponentType[SizeData]()

