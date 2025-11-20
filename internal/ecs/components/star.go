package components

import "github.com/yohamta/donburi"

// StarData 星星背景数据
type StarData struct {
	Speed float64 // 滚动速度
	Size  float64 // 大小
}

// Star 星星组件
var Star = donburi.NewComponentType[StarData]()

