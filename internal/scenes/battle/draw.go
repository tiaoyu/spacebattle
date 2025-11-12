package battle

import (
	"fmt"
	"spacebattle/internal/fonts"
	"spacebattle/internal/i18n"
	"spacebattle/internal/progress"
	"spacebattle/internal/sound"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Draw 渲染逻辑（从 shooter_game.go 拆分）
func (sg *ShooterGame) Draw(screen *ebiten.Image) {
	// 背景
	screen.Fill(color.RGBA{R: 0, G: 0, B: 20, A: 255})

	// 星星
	for _, star := range sg.stars {
		vector.DrawFilledRect(screen, float32(star.X), float32(star.Y), float32(star.Size), float32(star.Size), color.RGBA{R: 255, G: 255, B: 255, A: 255}, false)
	}

	// 玩家
	sg.drawPlayer(screen)

	// 子弹
	for _, bullet := range sg.bullets {
		if bullet.Active {
			sg.drawBullet(screen, bullet)
		}
	}

	// 敌机
	for _, enemy := range sg.enemies {
		if enemy.Active {
			sg.drawEnemy(screen, enemy)
		}
	}

	// Boss
	if sg.bossActive && sg.boss.Active {
		sg.drawBoss(screen)
	}

	// 爆炸
	for _, explosion := range sg.explosions {
		if explosion.Active {
			for r := 0.0; r < explosion.Radius; r += 3 {
				alpha := 1.0 - r/explosion.MaxRadius
				c := color.RGBA{R: uint8(255 * alpha), G: uint8(100 * alpha), B: 0, A: 255}
				vector.DrawFilledRect(screen, float32(explosion.X-r), float32(explosion.Y-r), float32(r*2), float32(r*2), c, false)
			}
		}
	}

	// UI
	scoreText := fmt.Sprintf("%s: %d", i18n.T("common.score"), sg.score)
	livesText := fmt.Sprintf("%s: %d", i18n.T("common.lives"), sg.lives)
	fonts.DrawText(screen, scoreText, 10, 30, color.White)
	fonts.DrawText(screen, livesText, 10, 50, color.White)

	// 右上角：时间与阶段（波次/Boss）
	elapsed := time.Since(sg.startTime)
	if elapsed < 0 {
		elapsed = 0
	}
	remain := sg.totalDuration - elapsed
	if remain < 0 {
		remain = 0
	}
	remainSec := int(remain.Seconds() + 0.5)
	timeText := fmt.Sprintf("%s: %d%s", i18n.T("common.time_left"), remainSec, i18n.T("unit.s"))
	fonts.DrawText(screen, timeText, 650, 10, color.White)

	phaseText := ""
	if (sg.bossActive && sg.boss.Active) || elapsed >= sg.smallPhaseDuration {
		phaseText = i18n.T("common.boss")
	} else {
		phaseText = fmt.Sprintf("%s %d/%d", i18n.T("common.wave"), sg.waveIndex+1, sg.waveCount)
	}
	fonts.DrawText(screen, phaseText, 650, 30, color.White)

	fs := sg.player.Fire
	uiY := 90
	fonts.DrawText(screen, fmt.Sprintf("%s: %.1f", i18n.T("shooter.fire.rate"), fs.FireRateHz), 10, uiY, color.White)
	uiY += 18
	fonts.DrawText(screen, fmt.Sprintf("%s: %d", i18n.T("shooter.fire.per_shot"), fs.BulletsPerShot), 10, uiY, color.White)
	uiY += 18
	fonts.DrawText(screen, fmt.Sprintf("%s: %.0f°", i18n.T("shooter.fire.spread"), fs.SpreadDeg), 10, uiY, color.White)
	uiY += 18
	fonts.DrawText(screen, fmt.Sprintf("%s: %.1f", i18n.T("shooter.fire.speed"), fs.BulletSpeed), 10, uiY, color.White)
	uiY += 18
	fonts.DrawText(screen, fmt.Sprintf("%s: %d", i18n.T("shooter.fire.penetration"), fs.PenetrationCount), 10, uiY, color.White)
	uiY += 18
	fonts.DrawText(screen, fmt.Sprintf("%s: %s", i18n.T("shooter.fire.homing"), ternary(fs.EnableHoming, i18n.T("common.on"), i18n.T("common.off"))), 10, uiY, color.White)
	uiY += 18
	fonts.DrawText(screen, fmt.Sprintf("%s: %.2f", i18n.T("shooter.fire.turn_rate"), fs.HomingTurnRateRad), 10, uiY, color.White)
	uiY += 18
	fonts.DrawText(screen, fmt.Sprintf("%s: %.2f (%d%s)", i18n.T("shooter.fire.burst"), fs.BurstChance, fs.BurstInterval.Milliseconds(), i18n.T("unit.ms")), 10, uiY, color.White)

	// GM 面板
	if sg.gmOpen {
		// 背板
		vector.DrawFilledRect(screen, 520, 10, 270, 260, color.RGBA{R: 20, G: 20, B: 60, A: 255}, false)
		title := "[GM]"
		if sg.gmTab == 0 {
			title += " 技能 (Tab 切换)"
		} else {
			title += " 声效 (Tab 切换)"
		}
		fonts.DrawText(screen, title, 530, 25, color.RGBA{R: 255, G: 255, B: 0, A: 255})

		if sg.gmTab == 0 {
			options := []string{
				fmt.Sprintf("0 FireRateHz: %.2f", sg.player.Fire.FireRateHz),
				fmt.Sprintf("1 BulletsPer: %d", sg.player.Fire.BulletsPerShot),
				fmt.Sprintf("2 SpreadDeg: %.1f", sg.player.Fire.SpreadDeg),
				fmt.Sprintf("3 BulletSpeed: %.2f", sg.player.Fire.BulletSpeed),
				fmt.Sprintf("4 BurstChance: %.2f", sg.player.Fire.BurstChance),
				fmt.Sprintf("5 Penetration: %d", sg.player.Fire.PenetrationCount),
				fmt.Sprintf("6 Homing: %v", sg.player.Fire.EnableHoming),
				fmt.Sprintf("7 TurnRate: %.2f", sg.player.Fire.HomingTurnRateRad),
			}
			y := 50
			for i, s := range options {
				prefix := "  "
				col := color.RGBA{R: 220, G: 220, B: 220, A: 255}
				if i == sg.gmIndex {
					prefix = "> "
					col = color.RGBA{R: 255, G: 255, B: 0, A: 255}
				}
				fonts.DrawText(screen, prefix+s, 530, y, col)
				y += 20
			}
			fonts.DrawText(screen, "上下选择 左右调整", 530, y+5, color.RGBA{R: 180, G: 180, B: 180, A: 255})
		} else {
			cfg := sound.GetShootConfig()
			options := []string{
				fmt.Sprintf("0 Duration: %.2fs", cfg.DurationSec),
				fmt.Sprintf("1 BaseFreq: %.0f Hz", cfg.BaseFreq),
				fmt.Sprintf("2 MinFreq: %.0f Hz", cfg.MinFreq),
				fmt.Sprintf("3 Sweep: %.2f", cfg.SweepFactor),
				fmt.Sprintf("4 Decay: %.1f", cfg.Decay),
				fmt.Sprintf("5 Amp: %.2f", cfg.Amplitude),
				fmt.Sprintf("6 Wave: %s", cfg.Waveform),
			}
			y := 50
			for i, s := range options {
				prefix := "  "
				col := color.RGBA{R: 220, G: 220, B: 220, A: 255}
				if i == sg.gmIndex {
					prefix = "> "
					col = color.RGBA{R: 255, G: 255, B: 0, A: 255}
				}
				fonts.DrawText(screen, prefix+s, 530, y, col)
				y += 20
			}
			fonts.DrawText(screen, "上下选择 左右调整 (Tab切换)", 530, y+5, color.RGBA{R: 180, G: 180, B: 180, A: 255})
		}
	}

	// 说明与状态
	fonts.DrawText(screen, i18n.T("shooter.instructions"), 10, 10, color.RGBA{R: 200, G: 200, B: 200, A: 255})
	fonts.DrawText(screen, i18n.T("shooter.destroy_enemies"), 10, 70, color.RGBA{R: 200, G: 200, B: 200, A: 255})

	if sg.gameOver {
		fonts.DrawTextCentered(screen, i18n.T("common.game_over"), 0, 250, 800, color.RGBA{R: 255, G: 0, B: 0, A: 255})
		fonts.DrawTextCentered(screen, i18n.T("common.restart"), 0, 280, 800, color.White)
		// 失败：显示缓存奖励（为胜利奖励的 1/3）
		meritsText := fmt.Sprintf("%s: +%d  (%s: %d)", i18n.T("common.merits"), sg.rewardCached, i18n.T("common.merits_total"), progress.GetMerits())
		fonts.DrawTextCentered(screen, meritsText, 0, 340, 800, color.White)
		fonts.DrawTextCentered(screen, i18n.T("common.back_menu"), 0, 310, 800, color.White)
	}
	if !sg.gameOver && sg.victory {
		fonts.DrawTextCentered(screen, i18n.T("common.victory"), 0, 250, 800, color.RGBA{R: 0, G: 200, B: 0, A: 255})
		fonts.DrawTextCentered(screen, i18n.T("common.restart"), 0, 280, 800, color.White)
		// 胜利：显示缓存奖励
		meritsText := fmt.Sprintf("%s: +%d  (%s: %d)", i18n.T("common.merits"), sg.rewardCached, i18n.T("common.merits_total"), progress.GetMerits())
		fonts.DrawTextCentered(screen, meritsText, 0, 340, 800, color.White)
		fonts.DrawTextCentered(screen, i18n.T("common.back_menu"), 0, 310, 800, color.White)
	}
}

// drawPlayer 绘制玩家飞船（主体+发光+机翼）
func (sg *ShooterGame) drawPlayer(screen *ebiten.Image) {
	// 外发光层（柔和蓝绿）
	glow := color.RGBA{R: 40, G: 220, B: 180, A: 80}
	vector.DrawFilledRect(screen, float32(sg.player.X-6), float32(sg.player.Y-6), float32(sg.player.Width+12), float32(sg.player.Height+12), glow, false)

	// 主体
	body := color.RGBA{R: 30, G: 230, B: 140, A: 255}
	vector.DrawFilledRect(screen, float32(sg.player.X), float32(sg.player.Y), float32(sg.player.Width), float32(sg.player.Height), body, false)

	// 机翼与机头（矩形细节替代）
	wingH := float32(sg.player.Height * 0.4)
	vector.DrawFilledRect(screen, float32(sg.player.X-6), float32(float32(sg.player.Y+sg.player.Height)-wingH), 6, wingH, color.RGBA{R: 60, G: 230, B: 190, A: 160}, false)
	vector.DrawFilledRect(screen, float32(sg.player.X+sg.player.Width), float32(float32(sg.player.Y+sg.player.Height)-wingH), 6, wingH, color.RGBA{R: 60, G: 230, B: 190, A: 160}, false)
	// 机头亮块
	vector.DrawFilledRect(screen, float32(sg.player.X+sg.player.Width/2-6), float32(sg.player.Y-6), 12, 6, color.RGBA{R: 120, G: 255, B: 230, A: 180}, false)

	// 中心亮条
	cx := float32(sg.player.X + sg.player.Width/2 - 2)
	vector.DrawFilledRect(screen, cx, float32(sg.player.Y+2), 4, float32(sg.player.Height-4), color.RGBA{R: 230, G: 255, B: 255, A: 200}, false)
}

// drawBullet 绘制子弹（核心+辉光）
func (sg *ShooterGame) drawBullet(screen *ebiten.Image, b Bullet) {
	core := color.RGBA{R: 255, G: 240, B: 80, A: 255}
	glow := color.RGBA{R: 255, G: 200, B: 50, A: 100}
	vector.DrawFilledRect(screen, float32(b.X-1), float32(b.Y-1), float32(b.Width+2), float32(b.Height+2), glow, false)
	vector.DrawFilledRect(screen, float32(b.X), float32(b.Y), float32(b.Width), float32(b.Height), core, false)
}

// drawEnemy 绘制敌机（主体+描边+血条）
func (sg *ShooterGame) drawEnemy(screen *ebiten.Image, e EnemyShip) {
	// 背景辉光
	glow := color.RGBA{R: 255, G: 60, B: 60, A: 80}
	vector.DrawFilledRect(screen, float32(e.X-4), float32(e.Y-4), float32(e.Width+8), float32(e.Height+8), glow, false)

	// 主体
	body := color.RGBA{R: 240, G: 50, B: 50, A: 255}
	vector.DrawFilledRect(screen, float32(e.X), float32(e.Y), float32(e.Width), float32(e.Height), body, false)

	// 顶部血条
	if e.MaxHealth > 0 {
		w := float32(e.Width)
		rate := float32(e.Health) / float32(e.MaxHealth)
		rate = float32(clamp(float64(rate), 0, 1))
		// 背板
		vector.DrawFilledRect(screen, float32(e.X), float32(e.Y-6), w, 3, color.RGBA{R: 40, G: 0, B: 0, A: 200}, false)
		// 进度
		col := color.RGBA{R: 255, G: uint8(80 + 120*rate), B: 0, A: 220}
		vector.DrawFilledRect(screen, float32(e.X), float32(e.Y-6), w*rate, 3, col, false)
	}
}

// drawBoss 绘制Boss（更强烈的发光和更粗血条）
func (sg *ShooterGame) drawBoss(screen *ebiten.Image) {
	b := sg.boss
	// 外圈辉光
	glow := color.RGBA{R: 220, G: 50, B: 220, A: 80}
	vector.DrawFilledRect(screen, float32(b.X-10), float32(b.Y-10), float32(b.Width+20), float32(b.Height+20), glow, false)
	// 主体
	vector.DrawFilledRect(screen, float32(b.X), float32(b.Y), float32(b.Width), float32(b.Height), color.RGBA{R: 200, G: 0, B: 200, A: 255}, false)
	// 血条
	if b.MaxHealth > 0 {
		w := float32(b.Width)
		rate := float32(b.Health) / float32(b.MaxHealth)
		rate = float32(clamp(float64(rate), 0, 1))
		vector.DrawFilledRect(screen, float32(b.X), float32(b.Y-8), w, 5, color.RGBA{R: 20, G: 0, B: 20, A: 200}, false)
		vector.DrawFilledRect(screen, float32(b.X), float32(b.Y-8), w*rate, 5, color.RGBA{R: 255, G: 80, B: 255, A: 230}, false)
	}
}
