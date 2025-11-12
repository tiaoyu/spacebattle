package battle

import (
	"spacebattle/internal/sound"
	"time"
)

// 碰撞与爆炸逻辑（从 shooter_game.go 拆分）

func (sg *ShooterGame) checkCollisions() {
	// 子弹与小敌机
	for i := len(sg.bullets) - 1; i >= 0; i-- {
		if !sg.bullets[i].Active {
			continue
		}
		// 先与小敌机碰撞
		for j := len(sg.enemies) - 1; j >= 0; j-- {
			if !sg.enemies[j].Active {
				continue
			}
			if sg.bullets[i].X < sg.enemies[j].X+sg.enemies[j].Width &&
				sg.bullets[i].X+sg.bullets[i].Width > sg.enemies[j].X &&
				sg.bullets[i].Y < sg.enemies[j].Y+sg.enemies[j].Height &&
				sg.bullets[i].Y+sg.bullets[i].Height > sg.enemies[j].Y {
				// 处理命中（沿用原逻辑）
				if sg.bullets[i].RemainingPenetration > 0 {
					sg.bullets[i].RemainingPenetration--
				} else {
					sg.bullets[i].Active = false
				}
				sg.enemies[j].Health--
				sound.PlayHit()
				if sg.enemies[j].Health <= 0 {
					sg.enemies[j].Active = false
					sg.score += 10
					sg.killedEnemyCount++
					sg.createExplosion(sg.enemies[j].X+sg.enemies[j].Width/2, sg.enemies[j].Y+sg.enemies[j].Height/2)
				}
				break
			}
		}
		// 再与Boss碰撞
		if sg.bossActive && sg.boss.Active {
			if sg.bullets[i].X < sg.boss.X+sg.boss.Width &&
				sg.bullets[i].X+sg.bullets[i].Width > sg.boss.X &&
				sg.bullets[i].Y < sg.boss.Y+sg.boss.Height &&
				sg.bullets[i].Y+sg.bullets[i].Height > sg.boss.Y {
				if sg.bullets[i].RemainingPenetration > 0 {
					sg.bullets[i].RemainingPenetration--
				} else {
					sg.bullets[i].Active = false
				}
				// Boss 软收束：进入Boss阶段后按时间提高对Boss的伤害
				// 10s后×2，14s后×3（防止超时）
				dmg := 1
				if d := time.Since(sg.startTime) - sg.smallPhaseDuration; d > 0 {
					if d >= 14*time.Second {
						dmg = 3
					} else if d >= 10*time.Second {
						dmg = 2
					}
				}
				sg.boss.Health -= dmg
				sound.PlayHit()
				if sg.boss.Health <= 0 {
					sg.boss.Active = false
					sg.bossActive = false
					sg.victory = true
					sg.bossKilled = true
					sg.createExplosion(sg.boss.X+sg.boss.Width/2, sg.boss.Y+sg.boss.Height/2)
				}
			}
		}
	}

	// 检查玩家与敌机的碰撞
	for _, enemy := range sg.enemies {
		if !enemy.Active {
			continue
		}

		if sg.player.X < enemy.X+enemy.Width &&
			sg.player.X+sg.player.Width > enemy.X &&
			sg.player.Y < enemy.Y+enemy.Height &&
			sg.player.Y+sg.player.Height > enemy.Y {

			// 碰撞发生
			sg.lives--
			if sg.lives <= 0 {
				sg.gameOver = true
			} else {
				// 重置玩家位置
				sg.player.X = 400
				sg.player.Y = 500
			}

			// 创建爆炸效果
			sg.createExplosion(sg.player.X+sg.player.Width/2, sg.player.Y+sg.player.Height/2)
			break
		}
	}
}

func (sg *ShooterGame) createExplosion(x, y float64) {
	explosion := Explosion{
		X:         x,
		Y:         y,
		Radius:    0,
		MaxRadius: 30,
		Active:    true,
		Timer:     0,
	}
	sg.explosions = append(sg.explosions, explosion)
}
