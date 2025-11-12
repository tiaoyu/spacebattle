package battle

import (
	"spacebattle/internal/balance"
	"spacebattle/internal/progress"
	"spacebattle/internal/sound"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

// Update 更新游戏逻辑（从 shooter_game.go 拆分）
func (sg *ShooterGame) Update() error {
	if sg.gameOver || sg.victory {
		if sg.input.IsKeyJustPressed(ebiten.KeyR) {
			*sg = *NewShooterGameWithOptions(sg.initialOptions)
		}
		// 结算功勋（一次性）：胜利发放 2x~3x；失败发放胜利奖励的三分之一
		if !sg.settled {
			// 统计值：将 Boss 计为 1 个敌机
			spawned := sg.spawnedEnemyCount
			if sg.bossSpawned {
				spawned += 1
			}
			kills := sg.killedEnemyCount
			if sg.bossKilled {
				kills += 1
			}
			elapsed := time.Since(sg.startTime)
			winReward := balance.ComputeMeritReward(sg.difficultyMul, kills, spawned, elapsed, sg.totalDuration, true)
			if sg.victory {
				sg.rewardCached = winReward
			} else {
				sg.rewardCached = winReward / 3
			}
			if sg.rewardCached > 0 {
				progress.AddMerits(sg.rewardCached)
			}
			sg.settled = true
		}
		return nil
	}

	sg.input.Update()

	// 处理玩家输入
	sg.handleInput()

	// 如果 GM 面板打开，则不进行游戏逻辑更新
	if sg.gmOpen {
		return nil
	}

	// 硬收束：总时长达到后若未胜利，直接结算为胜利
	if time.Since(sg.startTime) >= sg.totalDuration && !sg.victory {
		sg.victory = true
		return nil
	}

	// 处理计划中的连射
	sg.processScheduledShots()

	// 更新子弹
	sg.updateBullets()
	// 更新敌机与Boss
	sg.updateEnemies()
	// 更新爆炸
	sg.updateExplosions()
	// 星星
	sg.updateStars()
	// 碰撞
	sg.checkCollisions()
	// 生成（按波次）
	sg.spawnEnemies()

	return nil
}

// processScheduledShots 处理计划中的连射（从 shooter_game.go 拆分）
func (sg *ShooterGame) processScheduledShots() {
	if len(sg.scheduledShots) == 0 {
		return
	}
	now := time.Now()
	idx := 0
	for idx < len(sg.scheduledShots) && sg.scheduledShots[idx].Before(now) {
		sg.shoot()
		sound.PlayShoot()
		idx++
	}
	if idx > 0 {
		sg.scheduledShots = sg.scheduledShots[idx:]
	}
}
