package components

import "github.com/yohamta/donburi"

// PlayerInputData 玩家输入数据
type PlayerInputData struct {
	Speed float64 // 移动速度
}

// PlayerInput 玩家输入组件
var PlayerInput = donburi.NewComponentType[PlayerInputData]()

