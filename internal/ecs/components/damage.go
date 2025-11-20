package components

import "github.com/yohamta/donburi"

// DamageData 伤害数据
type DamageData struct {
	Value int // 伤害值
}

// Damage 伤害组件
var Damage = donburi.NewComponentType[DamageData]()
