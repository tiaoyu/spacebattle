package systems

import (
	"fmt"
	"image/color"
	"time"

	"spacebattle/internal/config"
	"spacebattle/internal/ecs/components"
	"spacebattle/internal/ecs/tags"
	"spacebattle/internal/fonts"
	"spacebattle/internal/i18n"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
)

// RenderSystem 渲染系统
type RenderSystem struct{}

// NewRenderSystem 创建渲染系统
func NewRenderSystem() *RenderSystem {
	return &RenderSystem{}
}

// Draw 绘制所有实体
func (s *RenderSystem) Draw(w donburi.World, screen *ebiten.Image) {
	// 绘制背景
	screen.Fill(config.DefaultConfig().BackgroundColor)

	// 绘制星星
	s.DrawStars(w, screen)

	// 绘制玩家
	s.DrawEntities(w, screen, tags.Player)

	// 绘制子弹
	s.DrawEntities(w, screen, tags.Bullet)

	// 绘制敌机（所有类型）
	s.DrawEntities(w, screen, tags.Enemy)
	s.DrawEntities(w, screen, tags.EnemyShooter)
	s.DrawEntities(w, screen, tags.EnemyZigzag)
	s.DrawEntities(w, screen, tags.EnemyTank)

	// 绘制敌机子弹
	s.DrawEntities(w, screen, tags.EnemyBullet)

	// 绘制 Boss
	s.DrawBoss(w, screen)

	// 绘制爆炸效果
	s.DrawExplosions(w, screen)

	// 绘制 HUD
	s.DrawHUD(w, screen)

	// 绘制结算界面
	s.DrawGameOver(w, screen)
}

// DrawStars 绘制星星
func (s *RenderSystem) DrawStars(w donburi.World, screen *ebiten.Image) {
	starQuery := query.NewQuery(
		filter.Contains(tags.Star, components.Position, components.Star, components.Sprite),
	)

	starQuery.Each(w, func(entry *donburi.Entry) {
		pos := components.Position.Get(entry)
		star := components.Star.Get(entry)
		sprite := components.Sprite.Get(entry)

		c := sprite.Color.(color.RGBA)
		vector.DrawFilledCircle(screen, float32(pos.X), float32(pos.Y), float32(star.Size), c, true)
	})
}

// DrawEntities 绘制指定标签的实体
func (s *RenderSystem) DrawEntities(w donburi.World, screen *ebiten.Image, tag donburi.IComponentType) {
	entityQuery := query.NewQuery(
		filter.Contains(tag, components.Position, components.Size, components.Sprite),
	)

	entityQuery.Each(w, func(entry *donburi.Entry) {
		pos := components.Position.Get(entry)
		size := components.Size.Get(entry)
		sprite := components.Sprite.Get(entry)

		c := sprite.Color.(color.RGBA)

		if sprite.Shape == "circle" {
			vector.DrawFilledCircle(
				screen,
				float32(pos.X+size.Width/2),
				float32(pos.Y+size.Height/2),
				float32(size.Width/2),
				c,
				true,
			)
		} else {
			// 默认为矩形
			vector.DrawFilledRect(
				screen,
				float32(pos.X),
				float32(pos.Y),
				float32(size.Width),
				float32(size.Height),
				c,
				true,
			)
		}
	})
}

// DrawBoss 绘制 Boss（带血条）
func (s *RenderSystem) DrawBoss(w donburi.World, screen *ebiten.Image) {
	bossQuery := query.NewQuery(
		filter.Contains(tags.Boss, components.Position, components.Size, components.Health, components.Sprite),
	)

	bossQuery.Each(w, func(entry *donburi.Entry) {
		pos := components.Position.Get(entry)
		size := components.Size.Get(entry)
		health := components.Health.Get(entry)
		sprite := components.Sprite.Get(entry)

		// 绘制 Boss 本体
		c := sprite.Color.(color.RGBA)
		vector.DrawFilledRect(
			screen,
			float32(pos.X),
			float32(pos.Y),
			float32(size.Width),
			float32(size.Height),
			c,
			true,
		)

		// 绘制血条背景
		cfg := config.DefaultConfig()
		barWidth := size.Width
		barHeight := 5.0
		vector.DrawFilledRect(
			screen,
			float32(pos.X),
			float32(pos.Y-10),
			float32(barWidth),
			float32(barHeight),
			cfg.BossHPBarBg,
			true,
		)

		// 绘制当前血量
		healthPercent := float64(health.Current) / float64(health.Max)
		vector.DrawFilledRect(
			screen,
			float32(pos.X),
			float32(pos.Y-10),
			float32(barWidth*healthPercent),
			float32(barHeight),
			cfg.BossHPBarFg,
			true,
		)
	})
}

// DrawExplosions 绘制爆炸效果
func (s *RenderSystem) DrawExplosions(w donburi.World, screen *ebiten.Image) {
	explosionQuery := query.NewQuery(
		filter.Contains(tags.Explosion, components.Position, components.Lifetime, components.Sprite),
	)

	explosionQuery.Each(w, func(entry *donburi.Entry) {
		pos := components.Position.Get(entry)
		lifetime := components.Lifetime.Get(entry)
		sprite := components.Sprite.Get(entry)

		c := sprite.Color.(color.RGBA)
		// 透明度随时间衰减
		c.A = uint8(255 * (1 - lifetime.Timer))

		vector.DrawFilledCircle(
			screen,
			float32(pos.X),
			float32(pos.Y),
			float32(lifetime.Radius),
			c,
			true,
		)
	})
}

// DrawHUD 绘制 HUD
func (s *RenderSystem) DrawHUD(w donburi.World, screen *ebiten.Image) {
	// 获取游戏状态
	var gameState *components.GameStateData
	query.NewQuery(filter.Contains(components.GameState)).Each(w, func(entry *donburi.Entry) {
		gameState = components.GameState.Get(entry)
	})

	if gameState == nil {
		return
	}

	// 绘制分数和生命
	scoreText := fmt.Sprintf("%s: %d", i18n.T("common.score"), gameState.Score)
	fonts.DrawText(screen, scoreText, 10, 10, color.White)

	livesText := fmt.Sprintf("%s: %d", i18n.T("common.lives"), gameState.Lives)
	fonts.DrawText(screen, livesText, 10, 30, color.White)

	// 绘制时间/波次信息
	cfg := config.DefaultConfig()
	elapsed := time.Since(gameState.StartTime)
	if elapsed < gameState.SmallPhaseDuration {
		waveText := fmt.Sprintf("Wave: %d/%d", gameState.WaveIndex+1, gameState.WaveCount)
		fonts.DrawText(screen, waveText, 700, 10, color.White)

		timeLeft := gameState.SmallPhaseDuration - elapsed
		timeText := fmt.Sprintf("Time: %ds", int(timeLeft.Seconds()))
		fonts.DrawText(screen, timeText, 700, 30, color.White)
	} else {
		fonts.DrawText(screen, "BOSS!", 700, 10, cfg.UIBossWarningColor)
	}
}

// DrawGameOver 绘制游戏结束界面
func (s *RenderSystem) DrawGameOver(w donburi.World, screen *ebiten.Image) {
	// 获取游戏状态
	var gameState *components.GameStateData
	query.NewQuery(filter.Contains(components.GameState)).Each(w, func(entry *donburi.Entry) {
		gameState = components.GameState.Get(entry)
	})

	if gameState == nil {
		return
	}

	if !gameState.GameOver && !gameState.Victory {
		return
	}

	// 半透明背景
	cfg := config.DefaultConfig()
	vector.DrawFilledRect(screen, 0, 0, 800, 600, cfg.UIOverlayColor, true)

	// 结果标题
	var title string
	var titleColor color.Color
	if gameState.Victory {
		title = i18n.T("common.victory")
		titleColor = cfg.UIVictoryColor
	} else {
		title = i18n.T("common.game_over")
		titleColor = cfg.UIGameOverColor
	}

	fonts.DrawTextCenteredLarge(screen, title, 0, 200, 800, titleColor)

	// 得分
	scoreText := fmt.Sprintf("%s: %d", i18n.T("common.score"), gameState.Score)
	fonts.DrawTextCentered(screen, scoreText, 0, 250, 800, color.White)

	// 功勋奖励详情（如果已结算）
	if gameState.Settled {
		y := 285
		breakdown := gameState.RewardBreakdown

		// 显示表现评分
		if breakdown.PerformanceScore > 0 {
			performanceText := fmt.Sprintf("Performance: %.0f/100", breakdown.PerformanceScore)
			performanceColor := cfg.UITextColor
			if breakdown.PerformanceScore >= 90 {
				performanceColor = cfg.UIVictoryColor // 绿色
			} else if breakdown.PerformanceScore >= 70 {
				performanceColor = cfg.UIMeritColor // 金色
			}
			fonts.DrawTextCentered(screen, performanceText, 0, y, 800, performanceColor)
			y += 30
		}

		// 奖励分解
		if breakdown.BaseReward > 0 {
			fonts.DrawTextCentered(screen, fmt.Sprintf("Base: +%d", breakdown.BaseReward), 0, y, 800, cfg.UITextColor)
			y += 22
		}
		if breakdown.DifficultyBonus > 0 {
			fonts.DrawTextCentered(screen, fmt.Sprintf("Difficulty: +%d", breakdown.DifficultyBonus), 0, y, 800, cfg.UIMeritColor)
			y += 22
		}
		if breakdown.KillBonus > 0 {
			fonts.DrawTextCentered(screen, fmt.Sprintf("Kill: +%d", breakdown.KillBonus), 0, y, 800, cfg.UIMeritColor)
			y += 22
		}
		if breakdown.SpeedBonus > 0 {
			fonts.DrawTextCentered(screen, fmt.Sprintf("Speed: +%d", breakdown.SpeedBonus), 0, y, 800, cfg.UIMeritColor)
			y += 22
		}
		if breakdown.BossBonus > 0 {
			fonts.DrawTextCentered(screen, fmt.Sprintf("Boss: +%d", breakdown.BossBonus), 0, y, 800, cfg.UIVictoryColor)
			y += 22
		}
		if breakdown.PerfectBonus > 0 {
			fonts.DrawTextCentered(screen, fmt.Sprintf("PERFECT: +%d", breakdown.PerfectBonus), 0, y, 800, cfg.UIVictoryColor)
			y += 22
		}

		// 总计
		if breakdown.TotalReward > 0 {
			y += 8
			totalText := fmt.Sprintf("Total Merits: +%d", breakdown.TotalReward)
			fonts.DrawTextCentered(screen, totalText, 0, y, 800, cfg.UIMeritColor)
		}
	}

	// 操作提示
	restartText := i18n.T("common.restart")
	fonts.DrawTextCentered(screen, restartText, 0, 520, 800, cfg.UIHintColor)

	backText := i18n.T("common.back_menu")
	fonts.DrawTextCentered(screen, backText, 0, 550, 800, cfg.UIHintColor)
}
