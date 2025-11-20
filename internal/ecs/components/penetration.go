package components

import "github.com/yohamta/donburi"

// PenetrationData 穿透能力数据
type PenetrationData struct {
	Remaining int // 剩余穿透次数
}

// Penetration 穿透组件
var Penetration = donburi.NewComponentType[PenetrationData]()

