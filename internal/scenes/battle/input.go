package battle

import (
	"spacebattle/internal/sound"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

// 处理输入（从 shooter_game.go 拆分）
func (sg *ShooterGame) handleInput() {
	// GM 调试开关
	if sg.input.IsKeyJustPressed(ebiten.KeyG) {
		sg.gmOpen = !sg.gmOpen
	}

	// GM 面板输入
	if sg.gmOpen {
		// Tab 切换面板
		if sg.input.IsKeyJustPressed(ebiten.KeyTab) {
			sg.gmTab = (sg.gmTab + 1) % 2
		}

		if sg.gmTab == 0 {
			// 技能面板
			if sg.input.IsKeyJustPressed(ebiten.KeyUp) || sg.input.IsKeyJustPressed(ebiten.KeyArrowUp) {
				sg.gmIndex--
				if sg.gmIndex < 0 {
					sg.gmIndex = 7
				}
			}
			if sg.input.IsKeyJustPressed(ebiten.KeyDown) || sg.input.IsKeyJustPressed(ebiten.KeyArrowDown) {
				sg.gmIndex++
				if sg.gmIndex > 7 {
					sg.gmIndex = 0
				}
			}

			// 调整值（左右）
			if sg.input.IsKeyJustPressed(ebiten.KeyLeft) || sg.input.IsKeyJustPressed(ebiten.KeyArrowLeft) {
				sg.adjustGMValue(false)
			}
			if sg.input.IsKeyJustPressed(ebiten.KeyRight) || sg.input.IsKeyJustPressed(ebiten.KeyArrowRight) {
				sg.adjustGMValue(true)
			}
		} else {
			// 声效面板
			if sg.input.IsKeyJustPressed(ebiten.KeyUp) || sg.input.IsKeyJustPressed(ebiten.KeyArrowUp) {
				sg.gmIndex--
				if sg.gmIndex < 0 {
					sg.gmIndex = 6
				}
			}
			if sg.input.IsKeyJustPressed(ebiten.KeyDown) || sg.input.IsKeyJustPressed(ebiten.KeyArrowDown) {
				sg.gmIndex++
				if sg.gmIndex > 6 {
					sg.gmIndex = 0
				}
			}
			if sg.input.IsKeyJustPressed(ebiten.KeyLeft) || sg.input.IsKeyJustPressed(ebiten.KeyArrowLeft) {
				sg.adjustGMSound(false)
			}
			if sg.input.IsKeyJustPressed(ebiten.KeyRight) || sg.input.IsKeyJustPressed(ebiten.KeyArrowRight) {
				sg.adjustGMSound(true)
			}
		}
		return
	}

	// 移动
	if sg.input.IsKeyPressed(ebiten.KeyLeft) && sg.player.X > 0 {
		sg.player.X -= sg.player.Speed
	}
	if sg.input.IsKeyPressed(ebiten.KeyRight) && sg.player.X < 800-sg.player.Width {
		sg.player.X += sg.player.Speed
	}
	if sg.input.IsKeyPressed(ebiten.KeyUp) && sg.player.Y > 0 {
		sg.player.Y -= sg.player.Speed
	}
	if sg.input.IsKeyPressed(ebiten.KeyDown) && sg.player.Y < 600-sg.player.Height {
		sg.player.Y += sg.player.Speed
	}

	// 射击
	if sg.input.IsKeyPressed(ebiten.KeySpace) && time.Since(sg.lastShot) >= sg.shotDelay {
		sg.shoot()
		sound.PlayShoot()
		if rand.Float64() < sg.player.Fire.BurstChance {
			delay := sg.player.Fire.BurstInterval
			if delay <= 0 {
				delay = time.Duration(float64(sg.shotDelay) * 0.3)
			}
			sg.scheduledShots = append(sg.scheduledShots, time.Now().Add(delay))
		}
		sg.lastShot = time.Now()
	}
}
