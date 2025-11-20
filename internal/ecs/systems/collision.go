package systems

import (
	"time"

	"spacebattle/internal/ecs"
	"spacebattle/internal/ecs/components"
	"spacebattle/internal/ecs/tags"
	"spacebattle/internal/sound"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
)

// CollisionSystem 碰撞检测系统
type CollisionSystem struct {
	world             *ecs.World
	shipAbilitySystem *ShipAbilitySystem
	particleSystem    *ParticleSystem
	shakeSystem       *ScreenShakeSystem
}

// NewCollisionSystem 创建碰撞检测系统
func NewCollisionSystem(world *ecs.World, shipAbility *ShipAbilitySystem, particle *ParticleSystem, shake *ScreenShakeSystem) *CollisionSystem {
	return &CollisionSystem{
		world:             world,
		shipAbilitySystem: shipAbility,
		particleSystem:    particle,
		shakeSystem:       shake,
	}
}

// Update 更新碰撞检测
func (s *CollisionSystem) Update(w donburi.World) {
	s.CheckBulletEnemyCollisions(w)
	s.CheckBulletBossCollisions(w)
	s.CheckPlayerEnemyCollisions(w)
	s.CheckEnemyBulletPlayerCollisions(w)
}

// CheckBulletEnemyCollisions 检查子弹与敌机碰撞
func (s *CollisionSystem) CheckBulletEnemyCollisions(w donburi.World) {
	bulletQuery := query.NewQuery(
		filter.Contains(tags.Bullet, components.Position, components.Size, components.Damage),
	)

	// 创建查询所有类型的敌机（基础型、射击型、之字型、肉盾型）
	allEnemyFilters := []query.Query{
		*query.NewQuery(filter.Contains(tags.Enemy, components.Position, components.Size, components.Health)),
		*query.NewQuery(filter.Contains(tags.EnemyShooter, components.Position, components.Size, components.Health)),
		*query.NewQuery(filter.Contains(tags.EnemyZigzag, components.Position, components.Size, components.Health)),
		*query.NewQuery(filter.Contains(tags.EnemyTank, components.Position, components.Size, components.Health)),
	}

	var bulletsToRemove []*donburi.Entry
	var enemiesToRemove []*donburi.Entry

	// 获取游戏状态
	var gameState *components.GameStateData
	query.NewQuery(filter.Contains(components.GameState)).Each(w, func(entry *donburi.Entry) {
		gameState = components.GameState.Get(entry)
	})

	// 获取玩家实体用于被动触发
	var playerEntry *donburi.Entry
	query.NewQuery(filter.Contains(tags.Player)).Each(w, func(entry *donburi.Entry) {
		playerEntry = entry
	})

	bulletQuery.Each(w, func(bullet *donburi.Entry) {
		bulletPos := components.Position.Get(bullet)
		bulletSize := components.Size.Get(bullet)
		bulletDamage := components.Damage.Get(bullet)
		bulletRemoved := false

		// 检查所有类型的敌机
		for _, enemyQuery := range allEnemyFilters {
			enemyQuery.Each(w, func(enemy *donburi.Entry) {
				if bulletRemoved {
					return
				}

				enemyPos := components.Position.Get(enemy)
				enemySize := components.Size.Get(enemy)
				enemyHealth := components.Health.Get(enemy)

				// AABB 碰撞检测
				if s.CheckAABB(
					bulletPos.X, bulletPos.Y, bulletSize.Width, bulletSize.Height,
					enemyPos.X, enemyPos.Y, enemySize.Width, enemySize.Height,
				) {
					// 处理穿透
					if bullet.HasComponent(components.Penetration) {
						pen := components.Penetration.Get(bullet)
						if pen.Remaining > 0 {
							pen.Remaining--
						} else {
							bulletsToRemove = append(bulletsToRemove, bullet)
							bulletRemoved = true
						}
					} else {
						bulletsToRemove = append(bulletsToRemove, bullet)
						bulletRemoved = true
					}

					// 减少敌机血量（使用子弹伤害）
					damage := bulletDamage.Value
					if damage <= 0 {
						damage = 1 // 最小伤害为1
					}
					enemyHealth.Current -= damage
					sound.PlayHit()

					if enemyHealth.Current <= 0 {
						// 敌机被击毁
						enemiesToRemove = append(enemiesToRemove, enemy)
						if gameState != nil {
							gameState.Score += 10
							gameState.KilledEnemyCount++
						}

						// 触发被动技能（Alpha回血、Beta叠buff）
						if s.shipAbilitySystem != nil && playerEntry != nil {
							s.shipAbilitySystem.OnEnemyKilled(w, playerEntry)
						}

						// 创建粒子效果
						if s.particleSystem != nil {
							s.particleSystem.CreateExplosionParticles(
								enemyPos.X+enemySize.Width/2,
								enemyPos.Y+enemySize.Height/2,
							)
						}

						// 触发屏幕震动
						if s.shakeSystem != nil {
							s.shakeSystem.TriggerShake(2.0, 0.15)
						}

						// 创建爆炸效果
						s.world.CreateExplosion(
							enemyPos.X+enemySize.Width/2,
							enemyPos.Y+enemySize.Height/2,
							30,
						)
					}
				}
			})
		}
	})

	// 移除被标记的实体
	for _, bullet := range bulletsToRemove {
		w.Remove(bullet.Entity())
	}
	for _, enemy := range enemiesToRemove {
		w.Remove(enemy.Entity())
	}
}

// CheckBulletBossCollisions 检查子弹与 Boss 碰撞
func (s *CollisionSystem) CheckBulletBossCollisions(w donburi.World) {
	bulletQuery := query.NewQuery(
		filter.Contains(tags.Bullet, components.Position, components.Size, components.Damage),
	)

	bossQuery := query.NewQuery(
		filter.Contains(tags.Boss, components.Position, components.Size, components.Health),
	)

	var bulletsToRemove []*donburi.Entry
	var bossToRemove []*donburi.Entry

	// 获取游戏状态
	var gameState *components.GameStateData
	query.NewQuery(filter.Contains(components.GameState)).Each(w, func(entry *donburi.Entry) {
		gameState = components.GameState.Get(entry)
	})

	bulletQuery.Each(w, func(bullet *donburi.Entry) {
		bulletPos := components.Position.Get(bullet)
		bulletSize := components.Size.Get(bullet)
		bulletDamage := components.Damage.Get(bullet)
		bulletRemoved := false

		bossQuery.Each(w, func(boss *donburi.Entry) {
			if bulletRemoved {
				return
			}

			bossPos := components.Position.Get(boss)
			bossSize := components.Size.Get(boss)
			bossHealth := components.Health.Get(boss)

			// AABB 碰撞检测
			if s.CheckAABB(
				bulletPos.X, bulletPos.Y, bulletSize.Width, bulletSize.Height,
				bossPos.X, bossPos.Y, bossSize.Width, bossSize.Height,
			) {
				// 处理穿透
				if bullet.HasComponent(components.Penetration) {
					pen := components.Penetration.Get(bullet)
					if pen.Remaining > 0 {
						pen.Remaining--
					} else {
						bulletsToRemove = append(bulletsToRemove, bullet)
						bulletRemoved = true
					}
				} else {
					bulletsToRemove = append(bulletsToRemove, bullet)
					bulletRemoved = true
				}

				// 基础伤害
				baseDmg := bulletDamage.Value
				if baseDmg <= 0 {
					baseDmg = 1
				}

				// Boss 软收束：根据时间提高伤害倍率
				dmgMultiplier := 1
				if gameState != nil {
					elapsed := time.Since(gameState.StartTime)
					bossElapsed := elapsed - gameState.SmallPhaseDuration
					if bossElapsed > 0 {
						if bossElapsed >= 14*time.Second {
							dmgMultiplier = 3
						} else if bossElapsed >= 10*time.Second {
							dmgMultiplier = 2
						}
					}
				}

				// 应用伤害
				totalDmg := baseDmg * dmgMultiplier
				bossHealth.Current -= totalDmg
				sound.PlayHit()

				if bossHealth.Current <= 0 {
					// Boss 被击毁
					bossToRemove = append(bossToRemove, boss)
					if gameState != nil {
						gameState.Victory = true
						gameState.BossKilled = true
						gameState.Score += 200
					}
					// 创建爆炸效果
					s.world.CreateExplosion(
						bossPos.X+bossSize.Width/2,
						bossPos.Y+bossSize.Height/2,
						50,
					)
				}
			}
		})
	})

	// 移除被标记的实体
	for _, bullet := range bulletsToRemove {
		w.Remove(bullet.Entity())
	}
	for _, boss := range bossToRemove {
		w.Remove(boss.Entity())
	}
}

// CheckPlayerEnemyCollisions 检查玩家与敌机碰撞
func (s *CollisionSystem) CheckPlayerEnemyCollisions(w donburi.World) {
	playerQuery := query.NewQuery(
		filter.Contains(tags.Player, components.Position, components.Size, components.Health),
	)

	// 创建查询所有类型的敌机
	allEnemyFilters := []query.Query{
		*query.NewQuery(filter.Contains(tags.Enemy, components.Position, components.Size)),
		*query.NewQuery(filter.Contains(tags.EnemyShooter, components.Position, components.Size)),
		*query.NewQuery(filter.Contains(tags.EnemyZigzag, components.Position, components.Size)),
		*query.NewQuery(filter.Contains(tags.EnemyTank, components.Position, components.Size)),
	}

	playerQuery.Each(w, func(player *donburi.Entry) {
		playerPos := components.Position.Get(player)
		playerSize := components.Size.Get(player)
		playerHealth := components.Health.Get(player)

		// 检查所有类型的敌机
		for _, enemyFilter := range allEnemyFilters {
			enemyFilter.Each(w, func(enemy *donburi.Entry) {
				enemyPos := components.Position.Get(enemy)
				enemySize := components.Size.Get(enemy)

				// AABB 碰撞检测
				if s.CheckAABB(
					playerPos.X, playerPos.Y, playerSize.Width, playerSize.Height,
					enemyPos.X, enemyPos.Y, enemySize.Width, enemySize.Height,
				) {
					// 检查被动技能（无敌、护盾）
					damage := 1
					if s.shipAbilitySystem != nil {
						damage = s.shipAbilitySystem.OnPlayerDamaged(w, player)
					}

					// 应用伤害
					if damage > 0 {
						playerHealth.Current -= damage
						sound.PlayHit()

						// 触发屏幕震动
						if s.shakeSystem != nil {
							s.shakeSystem.TriggerShake(4.0, 0.2)
						}
					}

					// 移除敌机
					w.Remove(enemy.Entity())

					// 重置玩家位置
					playerPos.X = 400
					playerPos.Y = 500

					// 创建爆炸效果
					s.world.CreateExplosion(
						playerPos.X+playerSize.Width/2,
						playerPos.Y+playerSize.Height/2,
						30,
					)
				}
			})
		}
	})
}

// CheckEnemyBulletPlayerCollisions 检查敌机子弹与玩家碰撞
func (s *CollisionSystem) CheckEnemyBulletPlayerCollisions(w donburi.World) {
	bulletQuery := query.NewQuery(
		filter.Contains(tags.EnemyBullet, components.Position, components.Size),
	)

	playerQuery := query.NewQuery(
		filter.Contains(tags.Player, components.Position, components.Size, components.Health),
	)

	var bulletsToRemove []*donburi.Entry

	bulletQuery.Each(w, func(bullet *donburi.Entry) {
		bulletPos := components.Position.Get(bullet)
		bulletSize := components.Size.Get(bullet)

		playerQuery.Each(w, func(player *donburi.Entry) {
			playerPos := components.Position.Get(player)
			playerSize := components.Size.Get(player)
			playerHealth := components.Health.Get(player)

			if s.CheckAABB(
				bulletPos.X, bulletPos.Y, bulletSize.Width, bulletSize.Height,
				playerPos.X, playerPos.Y, playerSize.Width, playerSize.Height,
			) {
				// 标记子弹移除
				bulletsToRemove = append(bulletsToRemove, bullet)

				// 检查被动技能（无敌、护盾）
				damage := 1
				if s.shipAbilitySystem != nil {
					damage = s.shipAbilitySystem.OnPlayerDamaged(w, player)
				}

				// 应用伤害
				if damage > 0 {
					playerHealth.Current -= damage
					sound.PlayHit()

					// 触发屏幕震动
					if s.shakeSystem != nil {
						s.shakeSystem.TriggerShake(4.0, 0.2)
					}

					// 重置玩家位置
					playerPos.X = 400
					playerPos.Y = 500

					// 创建爆炸效果
					s.world.CreateExplosion(
						playerPos.X+playerSize.Width/2,
						playerPos.Y+playerSize.Height/2,
						30,
					)
				}
			}
		})
	})

	// 移除被标记的子弹
	for _, bullet := range bulletsToRemove {
		w.Remove(bullet.Entity())
	}
}

// CheckAABB 检查 AABB 碰撞
func (s *CollisionSystem) CheckAABB(x1, y1, w1, h1, x2, y2, w2, h2 float64) bool {
	return x1 < x2+w2 && x1+w1 > x2 && y1 < y2+h2 && y1+h1 > y2
}
