package menu

import (
	"fmt"
	"image/color"
	"time"

	"spacebattle/internal/balance"
	"spacebattle/internal/config"
	"spacebattle/internal/fonts"
	"spacebattle/internal/i18n"
	"spacebattle/internal/progress"
	"spacebattle/internal/scenes/battle"
	"spacebattle/internal/utils"

	"github.com/hajimehoshi/ebiten/v2"
)

// DeployScene 出征场景（选择难度倍率）
type DeployScene struct {
	input     *utils.InputManager
	opts      battle.PlayerOptions
	confirmed bool
	// 连续可调难度
	difficulty    float64
	minDifficulty float64
	maxDifficulty float64
	// 长按加速
	leftHoldStart  time.Time
	rightHoldStart time.Time
	lastAdjust     time.Time
}

func NewDeployScene(opts battle.PlayerOptions) *DeployScene {
	cfg := config.DefaultConfig()
	return &DeployScene{
		input:         utils.NewInputManager(),
		opts:          opts,
		difficulty:    cfg.DifficultyMin,
		minDifficulty: cfg.DifficultyMin,
		maxDifficulty: cfg.DifficultyMax,
	}
}

func (d *DeployScene) Confirmed() bool                  { return d.confirmed }
func (d *DeployScene) GetOptions() battle.PlayerOptions { return d.opts }

func (d *DeployScene) Update() error {
	d.input.Update()
	now := time.Now()

	// 计算步进与加速
	step := func() float64 {
		v := d.difficulty
		switch {
		case v < 2:
			return 0.1
		case v < 10:
			return 0.5
		case v < 100:
			return 1
		case v < 1000:
			return 10
		case v < 10000:
			return 100
		case v < 100000:
			return 1000
		case v < 1000000:
			return 10000
		default:
			return 100000
		}
	}()
	repeat := 150 * time.Millisecond

	// 右键（增加）
	if d.input.IsKeyJustPressed(ebiten.KeyRight) || d.input.IsKeyJustPressed(ebiten.KeyArrowRight) {
		d.rightHoldStart = now
		d.difficulty = clampFloat(d.difficulty+step, d.minDifficulty, d.maxDifficulty)
		d.lastAdjust = now
	} else if d.input.IsKeyPressed(ebiten.KeyRight) || d.input.IsKeyPressed(ebiten.KeyArrowRight) {
		hold := now.Sub(d.rightHoldStart)
		if hold > 500*time.Millisecond {
			repeat = 30 * time.Millisecond
		}
		if now.Sub(d.lastAdjust) >= repeat {
			d.difficulty = clampFloat(d.difficulty+step, d.minDifficulty, d.maxDifficulty)
			d.lastAdjust = now
		}
	}
	// 左键（减少）
	if d.input.IsKeyJustPressed(ebiten.KeyLeft) || d.input.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		d.leftHoldStart = now
		d.difficulty = clampFloat(d.difficulty-step, d.minDifficulty, d.maxDifficulty)
		d.lastAdjust = now
	} else if d.input.IsKeyPressed(ebiten.KeyLeft) || d.input.IsKeyPressed(ebiten.KeyArrowLeft) {
		hold := now.Sub(d.leftHoldStart)
		if hold > 500*time.Millisecond {
			repeat = 30 * time.Millisecond
		}
		if now.Sub(d.lastAdjust) >= repeat {
			d.difficulty = clampFloat(d.difficulty-step, d.minDifficulty, d.maxDifficulty)
			d.lastAdjust = now
		}
	}

	if d.input.IsKeyJustPressed(ebiten.KeyEnter) || d.input.IsKeyJustPressed(ebiten.KeySpace) {
		cost := balance.DifficultyCost(d.difficulty)
		if progress.SpendMerits(cost) {
			d.opts.DifficultyMultiplier = d.difficulty
			d.confirmed = true
		}
	}
	return nil
}

func (d *DeployScene) Draw(screen *ebiten.Image) {
	fonts.DrawTextCenteredLarge(screen, i18n.T("deploy.title"), 0, 120, 800, color.White)
	fonts.DrawTextCentered(screen, i18n.T("deploy.hint"), 0, 180, 800, color.RGBA{R: 200, G: 200, B: 200, A: 255})

	// 显示当前功勋
	fonts.DrawTextCentered(screen, i18n.T("common.merits")+": "+fmt.Sprintf("%d", progress.GetMerits()), 0, 210, 800, color.White)

	// 显示当前难度与成本
	cost := balance.DifficultyCost(d.difficulty)
	text := fmt.Sprintf("x%.2f  (%s:%d)", d.difficulty, i18n.T("upgrade.cost"), cost)
	fonts.DrawTextCentered(screen, text, 0, 260, 800, color.White)
	fonts.DrawTextCentered(screen, i18n.T("select.hint"), 0, 300, 800, color.RGBA{R: 200, G: 200, B: 200, A: 255})
}

func clampFloat(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
