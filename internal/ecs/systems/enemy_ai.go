package systems

import (
	"math"
	"time"

	"spacebattle/internal/config"
	"spacebattle/internal/ecs"
	"spacebattle/internal/ecs/components"
	"spacebattle/internal/ecs/tags"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
)

// EnemyAISystem 敌机AI系统
type EnemyAISystem struct {
	world *ecs.World
	cfg   *config.Config
}

// NewEnemyAISystem 创建敌机AI系统
func NewEnemyAISystem(world *ecs.World) *EnemyAISystem {
	return &EnemyAISystem{
		world: world,
		cfg:   config.DefaultConfig(),
	}
}

// Update 更新敌机AI行为
func (s *EnemyAISystem) Update(w donburi.World, dt float64) {
	// 获取玩家位置（用于射击型敌机瞄准）
	var playerPos *components.PositionData
	query.NewQuery(filter.Contains(tags.Player, components.Position)).Each(w, func(entry *donburi.Entry) {
		playerPos = components.Position.Get(entry)
	})

	// 处理射击型敌机
	s.processShooterEnemies(w, playerPos)

	// 处理之字型敌机
	s.processZigzagEnemies(w, dt)
}

// processShooterEnemies 处理射击型敌机
func (s *EnemyAISystem) processShooterEnemies(w donburi.World, playerPos *components.PositionData) {
	shooterQuery := query.NewQuery(
		filter.Contains(tags.EnemyShooter, components.Position, components.Size, components.EnemyAI),
	)

	now := time.Now()
	shooterQuery.Each(w, func(entry *donburi.Entry) {
		ai := components.EnemyAI.Get(entry)
		pos := components.Position.Get(entry)
		size := components.Size.Get(entry)

		// 检查敌机是否在屏幕内
		if !s.isInScreen(pos.X, pos.Y, size.Width, size.Height) {
			return // 不在屏幕内，不发射子弹
		}

		// 检查是否到达射击时间
		if now.Sub(ai.LastShotTime) >= ai.ShootInterval {
			if playerPos != nil {
				// 向玩家位置射击
				dx := playerPos.X - pos.X
				dy := playerPos.Y - pos.Y
				dist := math.Sqrt(dx*dx + dy*dy)
				if dist > 0 {
					// 创建敌机子弹
					vx := (dx / dist) * 4.0
					vy := (dy / dist) * 4.0
					s.world.CreateEnemyBullet(pos.X, pos.Y, vx, vy)
				}
			}
			ai.LastShotTime = now
		}
	})
}

// isInScreen 检查实体是否在屏幕内
func (s *EnemyAISystem) isInScreen(x, y, width, height float64) bool {
	screenWidth := float64(s.cfg.WindowWidth)
	screenHeight := float64(s.cfg.WindowHeight)

	// 实体完全超出屏幕边界则返回 false
	if x+width < 0 || x > screenWidth || y+height < 0 || y > screenHeight {
		return false
	}
	return true
}

// processZigzagEnemies 处理之字型敌机
func (s *EnemyAISystem) processZigzagEnemies(w donburi.World, dt float64) {
	zigzagQuery := query.NewQuery(
		filter.Contains(tags.EnemyZigzag, components.Velocity, components.EnemyAI),
	)

	zigzagQuery.Each(w, func(entry *donburi.Entry) {
		vel := components.Velocity.Get(entry)
		ai := components.EnemyAI.Get(entry)

		// 更新摆动相位
		ai.ZigzagPhase += ai.ZigzagSpeed * dt

		// 根据正弦函数计算横向速度
		vel.VX = math.Sin(ai.ZigzagPhase) * 2.0
	})
}
