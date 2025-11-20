package components

import (
	"time"

	"github.com/yohamta/donburi"
)

// HomingData 追踪能力数据
type HomingData struct {
	TurnRate          float64       // 转向速率（弧度/帧）
	Speed             float64       // 子弹速度（用于重新计算方向）
	TargetEntity      donburi.Entity // 当前锁定的目标实体
	LastRetargetTime  time.Time     // 上次重新锁定时间
	RetargetInterval  time.Duration // 重新锁定间隔（避免频繁切换目标）
}

// Homing 追踪组件
var Homing = donburi.NewComponentType[HomingData]()

