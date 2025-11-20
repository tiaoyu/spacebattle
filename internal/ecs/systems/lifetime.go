package systems

import (
	"spacebattle/internal/ecs/components"
	"spacebattle/internal/ecs/tags"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
)

// LifetimeSystem 生命周期系统
type LifetimeSystem struct{}

// NewLifetimeSystem 创建生命周期系统
func NewLifetimeSystem() *LifetimeSystem {
	return &LifetimeSystem{}
}

// Update 更新生命周期
func (s *LifetimeSystem) Update(w donburi.World) {
	// 更新爆炸效果
	s.UpdateExplosions(w)
}

// UpdateExplosions 更新爆炸效果
func (s *LifetimeSystem) UpdateExplosions(w donburi.World) {
	explosionQuery := query.NewQuery(
		filter.Contains(tags.Explosion, components.Lifetime),
	)

	var toRemove []*donburi.Entry

	explosionQuery.Each(w, func(entry *donburi.Entry) {
		lifetime := components.Lifetime.Get(entry)

		// 更新计时器
		lifetime.Timer += 0.1
		// 更新半径（从 0 渐变到最大，然后渐变回 0）
		lifetime.Radius = lifetime.MaxRadius * (1 - lifetime.Timer)

		// 如果生命周期结束，标记为移除
		if lifetime.Timer >= 1.0 {
			toRemove = append(toRemove, entry)
		}
	})

	// 移除过期的爆炸效果
	for _, entry := range toRemove {
		w.Remove(entry.Entity())
	}
}

