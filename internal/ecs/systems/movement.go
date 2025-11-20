package systems

import (
	"math/rand"

	"spacebattle/internal/ecs/components"
	"spacebattle/internal/ecs/tags"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
)

// MovementSystem 移动系统
type MovementSystem struct{}

// NewMovementSystem 创建移动系统
func NewMovementSystem() *MovementSystem {
	return &MovementSystem{}
}

// Update 更新所有有速度的实体位置
func (s *MovementSystem) Update(w donburi.World) {
	// 更新所有有位置和速度的实体
	movableQuery := query.NewQuery(
		filter.Contains(components.Position, components.Velocity),
	)

	movableQuery.Each(w, func(entry *donburi.Entry) {
		pos := components.Position.Get(entry)
		vel := components.Velocity.Get(entry)

		pos.X += vel.VX
		pos.Y += vel.VY
	})

	// Boss 特殊移动逻辑（边界反弹）
	s.UpdateBoss(w)

	// 星星特殊移动逻辑（循环滚动）
	s.UpdateStars(w)
}

// UpdateBoss 更新 Boss 移动（边界反弹）
func (s *MovementSystem) UpdateBoss(w donburi.World) {
	bossQuery := query.NewQuery(
		filter.Contains(tags.Boss, components.Position, components.Velocity, components.Size),
	)

	bossQuery.Each(w, func(entry *donburi.Entry) {
		pos := components.Position.Get(entry)
		vel := components.Velocity.Get(entry)
		size := components.Size.Get(entry)

		// X 轴边界反弹
		if pos.X < 50 || pos.X+size.Width > 750 {
			vel.VX = -vel.VX
		}

		// Y 轴边界反弹
		if pos.Y < 20 || pos.Y+size.Height > 300 {
			vel.VY = -vel.VY
		}
	})
}

// UpdateStars 更新星星背景（循环滚动）
func (s *MovementSystem) UpdateStars(w donburi.World) {
	starQuery := query.NewQuery(
		filter.Contains(tags.Star, components.Position, components.Star),
	)

	starQuery.Each(w, func(entry *donburi.Entry) {
		pos := components.Position.Get(entry)
		star := components.Star.Get(entry)

		// 向下移动
		pos.Y += star.Speed

		// 超出底部则回到顶部
		if pos.Y > 600 {
			pos.Y = -star.Size
			pos.X = rand.Float64() * 800
		}
	})
}

// CleanOutOfBounds 清理超出屏幕的实体
func (s *MovementSystem) CleanOutOfBounds(w donburi.World) {
	// 清理子弹
	bulletQuery := query.NewQuery(
		filter.Contains(tags.Bullet, components.Position),
	)

	var toRemove []*donburi.Entry

	bulletQuery.Each(w, func(entry *donburi.Entry) {
		pos := components.Position.Get(entry)
		if pos.Y < 0 || pos.Y > 600 {
			toRemove = append(toRemove, entry)
		}
	})

	for _, entry := range toRemove {
		w.Remove(entry.Entity())
	}

	// 清理敌机
	toRemove = nil
	enemyQuery := query.NewQuery(
		filter.Contains(tags.Enemy, components.Position),
	)

	enemyQuery.Each(w, func(entry *donburi.Entry) {
		pos := components.Position.Get(entry)
		if pos.Y > 600 {
			toRemove = append(toRemove, entry)
		}
	})

	for _, entry := range toRemove {
		w.Remove(entry.Entity())
	}
}

