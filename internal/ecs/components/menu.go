package components

import (
	"image/color"

	"github.com/yohamta/donburi"
)

// MenuOptionData 菜单选项数据
type MenuOptionData struct {
	Index    int
	TextKey  string
	Selected bool
}

// MenuOption 菜单选项组件
var MenuOption = donburi.NewComponentType[MenuOptionData]()

// UITextData UI 文本数据
type UITextData struct {
	Text  string
	Color color.Color
	Large bool // 是否使用大字体
}

// UIText UI 文本组件
var UIText = donburi.NewComponentType[UITextData]()

// MenuStateData 菜单状态数据
type MenuStateData struct {
	SelectedIndex int
	OptionCount   int
	Confirmed     bool
	// 战机选择相关
	ShipTemplates []ShipTemplate
	// 升级相关
	AvailableMerits int
	Upgrades        map[string]int
	// 出征相关
	DifficultyMul float64
}

// ShipTemplate 战机模板
type ShipTemplate struct {
	NameKey    string
	Speed      float64
	SizeScale  float64
	Lives      int
	PassiveKey string
}

// MenuState 菜单状态组件
var MenuState = donburi.NewComponentType[MenuStateData]()

