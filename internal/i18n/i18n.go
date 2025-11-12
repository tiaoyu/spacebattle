package i18n

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Language 语言类型
type Language string

const (
	Chinese Language = "zh"
	English Language = "en"
	Russian Language = "ru"
)

// TextMap 文本映射
type TextMap map[string]string

// I18n 国际化管理器
type I18n struct {
	currentLanguage Language
	texts           map[Language]TextMap
}

var globalI18n *I18n

// Init 初始化国际化系统
func Init() error {
	globalI18n = &I18n{
		currentLanguage: Chinese, // 默认中文
		texts:           make(map[Language]TextMap),
	}

	// 加载中文语言包
	if err := globalI18n.loadLanguage(Chinese); err != nil {
		return err
	}

	// 加载英文语言包
	if err := globalI18n.loadLanguage(English); err != nil {
		return err
	}

	// 加载俄文语言包
	if err := globalI18n.loadLanguage(Russian); err != nil {
		return err
	}

	return nil
}

// loadLanguage 加载指定语言的语言包
func (i *I18n) loadLanguage(lang Language) error {
	// 尝试从文件加载
	filename := filepath.Join("./assets", "i18n", string(lang)+".json")
	if data, err := os.ReadFile(filename); err == nil {
		var texts TextMap
		if err := json.Unmarshal(data, &texts); err == nil {
			i.texts[lang] = texts
			return nil
		}
	}

	// 如果文件不存在，使用内置的语言包
	i.texts[lang] = getBuiltinTexts(lang)
	return nil
}

// getBuiltinTexts 获取内置语言包
func getBuiltinTexts(lang Language) TextMap {
	switch lang {
	case Chinese:
		return TextMap{
			// 主菜单
			"menu.title":        "太空射击游戏",
			"menu.start":        "开始游戏",
			"menu.exit":         "退出",
			"menu.instructions": "使用方向键选择，回车确认",
			"menu.esc_hint":     "ESC键返回主菜单",

			// 通用
			"common.score":     "分数",
			"common.lives":     "生命",
			"common.game_over": "游戏结束",
			"common.restart":   "按 R 重新开始",
			"common.back_menu": "按 ESC 返回主菜单",
			"common.controls":  "控制说明",
			"common.on":        "开",
			"common.off":       "关",
			"unit.ms":          "毫秒",

			// 太空射击游戏
			"shooter.instructions":     "方向键移动飞船，空格键射击",
			"shooter.destroy_enemies":  "消灭敌机获得分数，避免被撞击",
			"shooter.fire.rate":        "射速",
			"shooter.fire.per_shot":    "每次发射数量",
			"shooter.fire.spread":      "散射角度",
			"shooter.fire.speed":       "子弹速度",
			"shooter.fire.penetration": "穿透次数",
			"shooter.fire.homing":      "追踪效果",
			"shooter.fire.turn_rate":   "追踪转向速率",
			"shooter.fire.burst":       "连发概率",
		}
	case English:
		return TextMap{
			// 主菜单
			"menu.title":        "Space Shooter",
			"menu.start":        "Space Shooter",
			"menu.exit":         "Exit",
			"menu.instructions": "Use arrow keys to select, Enter to confirm",
			"menu.esc_hint":     "ESC to return to main menu",

			// 通用
			"common.score":     "Score",
			"common.lives":     "Lives",
			"common.game_over": "Game Over",
			"common.restart":   "Press R to restart",
			"common.back_menu": "Press ESC to return to menu",
			"common.controls":  "Controls",
			"common.on":        "On",
			"common.off":       "Off",
			"unit.ms":          "ms",

			// 太空射击游戏
			"shooter.instructions":     "Arrow keys to move, Space to shoot",
			"shooter.destroy_enemies":  "Destroy enemies for points, avoid collisions",
			"shooter.fire.rate":        "Fire rate",
			"shooter.fire.per_shot":    "Bullets/shot",
			"shooter.fire.spread":      "Spread",
			"shooter.fire.speed":       "Bullet speed",
			"shooter.fire.penetration": "Penetration",
			"shooter.fire.homing":      "Homing",
			"shooter.fire.turn_rate":   "Turn rate",
			"shooter.fire.burst":       "Burst",
		}
	case Russian:
		return TextMap{
			// 主菜单
			"menu.title":        "Космический шутер",
			"menu.start":        "Космический шутер",
			"menu.exit":         "Выход",
			"menu.instructions": "Стрелки для выбора, Enter для подтверждения",
			"menu.esc_hint":     "ESC для возврата в главное меню",

			// 通用
			"common.score":     "Счёт",
			"common.lives":     "Жизни",
			"common.game_over": "Игра окончена",
			"common.restart":   "Нажмите R для перезапуска",
			"common.back_menu": "Нажмите ESC для возврата в меню",
			"common.controls":  "Управление",

			// 太空射击游戏
			"shooter.instructions":    "Стрелки для движения корабля, Пробел для стрельбы",
			"shooter.destroy_enemies": "Уничтожайте врагов за очки, избегайте столкновений",
		}
	default:
		return TextMap{}
	}
}

// GetText 获取文本
func GetText(key string) string {
	if globalI18n == nil {
		return key
	}

	if texts, exists := globalI18n.texts[globalI18n.currentLanguage]; exists {
		if text, exists := texts[key]; exists {
			return text
		}
	}

	return key
}

// SetLanguage 设置当前语言
func SetLanguage(lang Language) {
	if globalI18n != nil {
		globalI18n.currentLanguage = lang
	}
}

// GetCurrentLanguage 获取当前语言
func GetCurrentLanguage() Language {
	if globalI18n != nil {
		return globalI18n.currentLanguage
	}
	return Chinese
}

// T 获取文本的便捷函数
func T(key string) string {
	return GetText(key)
}
