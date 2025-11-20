package systems

import (
	"math/rand"

	"spacebattle/internal/ecs"
	"spacebattle/internal/ecs/components"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
)

// ScreenShakeSystem 屏幕震动系统
type ScreenShakeSystem struct {
	world *ecs.World
}

// NewScreenShakeSystem 创建屏幕震动系统
func NewScreenShakeSystem(world *ecs.World) *ScreenShakeSystem {
	return &ScreenShakeSystem{
		world: world,
	}
}

// Update 更新震动状态
func (s *ScreenShakeSystem) Update(w donburi.World, dt float64) {
	shakeQuery := query.NewQuery(
		filter.Contains(components.ScreenShake),
	)

	var toRemove []*donburi.Entry

	shakeQuery.Each(w, func(entry *donburi.Entry) {
		shake := components.ScreenShake.Get(entry)

		shake.Elapsed += dt
		if shake.Elapsed >= shake.Duration {
			// 震动结束
			shake.OffsetX = 0
			shake.OffsetY = 0
			toRemove = append(toRemove, entry)
			return
		}

		// 计算当前震动强度（随时间衰减）
		progress := shake.Elapsed / shake.Duration
		currentIntensity := shake.Intensity * (1.0 - progress)

		// 随机偏移
		shake.OffsetX = (rand.Float64()*2 - 1) * currentIntensity
		shake.OffsetY = (rand.Float64()*2 - 1) * currentIntensity
	})

	// 清理结束的震动
	for _, entry := range toRemove {
		w.Remove(entry.Entity())
	}
}

// GetOffset 获取当前震动偏移量
func (s *ScreenShakeSystem) GetOffset(w donburi.World) (float64, float64) {
	var totalX, totalY float64

	query.NewQuery(filter.Contains(components.ScreenShake)).Each(w, func(entry *donburi.Entry) {
		shake := components.ScreenShake.Get(entry)
		totalX += shake.OffsetX
		totalY += shake.OffsetY
	})

	return totalX, totalY
}

// TriggerShake 触发震动效果
func (s *ScreenShakeSystem) TriggerShake(intensity, duration float64) {
	s.world.CreateScreenShake(intensity, duration)
}

