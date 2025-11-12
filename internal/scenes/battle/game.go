package battle

import (
	"math"
	"math/rand"
	"time"

	"spacebattle/internal/config"
	"spacebattle/internal/sound"
	"spacebattle/internal/utils"
)

// ShooterGame 太空射击游戏场景
type ShooterGame struct {
	// 游戏状态
	player           PlayerShip
	initialOptions   PlayerOptions
	bullets          []Bullet
	enemies          []EnemyShip
	explosions       []Explosion
	stars            []Star
	score            int
	difficultyMul    float64
	lives            int
	settled          bool
	gameOver         bool
	victory          bool
	killedEnemyCount int
	bossKilled       bool
	rewardCached     int
	input            *utils.InputManager
	lastShot         time.Time
	shotDelay        time.Duration
	lastEnemy        time.Time
	enemyDelay       time.Duration
	gmOpen           bool
	gmIndex          int
	gmTab            int // 0: 技能, 1: 声效
	scheduledShots   []time.Time

	// 关卡与Boss

	spawnedEnemyCount int // 已生成数量
	maxSimultaneous   int // 同屏最大敌机
	batchSize         int // 每次生成数量
	baseEnemyHP       int // 敌机基础血量
	bossSpawned       bool
	bossActive        bool
	boss              EnemyShip

	// 波次
	// 时间制波次配置（60s 总时长：小怪45s，Boss15s；5波×9s）
	startTime          time.Time
	totalDuration      time.Duration
	smallPhaseDuration time.Duration
	waveLength         time.Duration
	waveCount          int
	waveIndex          int // 当前时间波次索引 0..(waveCount-1)
	waveMinIntervals   []time.Duration
}

// NewShooterGame 创建太空射击游戏
func NewShooterGame() *ShooterGame {
	cfg := config.DefaultConfig()
	game := &ShooterGame{
		player: PlayerShip{
			X:      400,
			Y:      500,
			Width:  40,
			Height: 30,
			Speed:  5,
			Fire: FireSkillConfig{
				FireRateHz:        5.0,
				BulletsPerShot:    1,
				SpreadDeg:         2,
				BulletSpeed:       8.0,
				BurstChance:       0.0,
				PenetrationCount:  0,
				EnableHoming:      false,
				HomingTurnRateRad: 0.01,
				BurstInterval:     60 * time.Millisecond,
			},
		},
		bullets:    make([]Bullet, 0),
		enemies:    make([]EnemyShip, 0),
		explosions: make([]Explosion, 0),
		stars:      make([]Star, 0),
		score:      0,
		lives:      3,
		gameOver:   false,
		victory:    false,
		input:      utils.NewInputManager(),
		lastShot:   time.Now(),
		shotDelay:  200 * time.Millisecond,
		lastEnemy:  time.Now(),
		enemyDelay: 800 * time.Millisecond,
		gmOpen:     false,
		gmIndex:    0,

		// 关卡参数（同屏与批量）
		spawnedEnemyCount: 0,
		maxSimultaneous:   cfg.MaxSimultaneous,
		batchSize:         cfg.BatchSize,
		baseEnemyHP:       1,
		bossSpawned:       false,
		bossActive:        false,
	}

	// 初始化音频
	sound.Init()

	// 初始化音画与计时
	game.initStars()
	game.recomputeShotDelay()

	// 时间制波次配置（与设计一致）
	game.startTime = time.Now()
	game.totalDuration = cfg.TotalDuration
	game.smallPhaseDuration = cfg.SmallPhaseDuration
	game.waveLength = cfg.WaveLength
	game.waveCount = cfg.WaveCount
	game.waveIndex = 0
	game.waveMinIntervals = cfg.WaveMinIntervals

	return game
}

// NewShooterGameWithOptions 创建带玩家选项的太空射击游戏
func NewShooterGameWithOptions(opts PlayerOptions) *ShooterGame {
	g := NewShooterGame()
	cfg := config.DefaultConfig()
	// 保存初始选项以便重开时复用
	g.initialOptions = opts
	// 应用外部传入的玩家属性
	if opts.Speed > 0 {
		g.player.Speed = opts.Speed
	}
	if opts.SizeScale > 0 {
		g.player.Width *= opts.SizeScale
		g.player.Height *= opts.SizeScale
	}
	if opts.Lives > 0 {
		g.lives = opts.Lives
	}
	if opts.DifficultyMultiplier > 0 {
		g.difficultyMul = opts.DifficultyMultiplier
	} else {
		g.difficultyMul = 1.0
	}
	// 应用被动
	switch opts.PassiveKey {
	case "passive.speed":
		g.player.Speed *= 1.5
	case "passive.small":
		g.player.Width *= 0.8
		g.player.Height *= 0.8
	case "passive.life":
		g.lives += 1
	}
	// 应用升级加成
	if opts.ModFireRateHz != 0 {
		g.player.Fire.FireRateHz = clamp(g.player.Fire.FireRateHz+opts.ModFireRateHz, 0.5, cfg.MaxFireRateHz)
		g.recomputeShotDelay()
	}
	if opts.ModBulletsPerShot > 0 {
		g.player.Fire.BulletsPerShot += opts.ModBulletsPerShot
		if g.player.Fire.BulletsPerShot > cfg.MaxBulletsPerShot {
			g.player.Fire.BulletsPerShot = cfg.MaxBulletsPerShot
		}
	}
	if opts.ModPenetration > 0 {
		g.player.Fire.PenetrationCount += opts.ModPenetration
		if g.player.Fire.PenetrationCount > cfg.MaxPenetration {
			g.player.Fire.PenetrationCount = cfg.MaxPenetration
		}
	}
	if opts.ModSpreadDeltaDeg != 0 {
		g.player.Fire.SpreadDeg = clamp(g.player.Fire.SpreadDeg+opts.ModSpreadDeltaDeg, 0, cfg.MaxSpreadDeg)
	}
	if opts.ModBulletSpeed != 0 {
		g.player.Fire.BulletSpeed = clamp(g.player.Fire.BulletSpeed+opts.ModBulletSpeed, 1, cfg.MaxBulletSpeed)
	}
	if opts.ModBurstChance != 0 {
		g.player.Fire.BurstChance = clamp(g.player.Fire.BurstChance+opts.ModBurstChance, 0, cfg.MaxBurstChance)
	}
	if opts.ModEnableHoming {
		g.player.Fire.EnableHoming = true
	}
	if opts.ModTurnRateRad != 0 {
		g.player.Fire.HomingTurnRateRad = clamp(g.player.Fire.HomingTurnRateRad+opts.ModTurnRateRad, 0, cfg.MaxTurnRateRad)
	}
	return g
}

// recomputeShotDelay 根据技能的 FireRateHz 重新计算发射冷却
func (sg *ShooterGame) recomputeShotDelay() {
	if sg.player.Fire.FireRateHz <= 0 {
		sg.shotDelay = 300 * time.Millisecond
		return
	}
	perShotSeconds := 1.0 / sg.player.Fire.FireRateHz
	sg.shotDelay = time.Duration(perShotSeconds * float64(time.Second))
}

// shoot 射击
func (sg *ShooterGame) shoot() {
	numBullets := max(sg.player.Fire.BulletsPerShot, 1)

	baseAngle := -math.Pi / 2
	spreadRad := sg.player.Fire.SpreadDeg * math.Pi / 180.0

	for i := range numBullets {
		var angle float64
		if numBullets == 1 || spreadRad == 0 {
			angle = baseAngle
		} else {
			t := float64(i) / float64(numBullets-1)
			angle = baseAngle + (t-0.5)*spreadRad
		}

		vx := math.Cos(angle) * sg.player.Fire.BulletSpeed
		vy := math.Sin(angle) * sg.player.Fire.BulletSpeed

		bullet := Bullet{
			X:                    sg.player.X + sg.player.Width/2 - 2,
			Y:                    sg.player.Y,
			VX:                   vx,
			VY:                   vy,
			Width:                4,
			Height:               10,
			Active:               true,
			Speed:                sg.player.Fire.BulletSpeed,
			RemainingPenetration: sg.player.Fire.PenetrationCount,
			Homing:               sg.player.Fire.EnableHoming,
			HomingTurnRate:       sg.player.Fire.HomingTurnRateRad,
		}
		sg.bullets = append(sg.bullets, bullet)
	}
}

// updateBullets 更新子弹
func (sg *ShooterGame) updateBullets() {
	for i := len(sg.bullets) - 1; i >= 0; i-- {
		if !sg.bullets[i].Active {
			sg.bullets = append(sg.bullets[:i], sg.bullets[i+1:]...)
			continue
		}

		// 追踪逻辑
		if sg.bullets[i].Homing {
			tx, ty, found := sg.findNearestEnemy(sg.bullets[i].X, sg.bullets[i].Y)
			if found {
				currentAngle := math.Atan2(sg.bullets[i].VY, sg.bullets[i].VX)
				desiredAngle := math.Atan2(ty-sg.bullets[i].Y, tx-sg.bullets[i].X)
				delta := wrapAngle(desiredAngle - currentAngle)
				turn := clamp(delta, -sg.bullets[i].HomingTurnRate, sg.bullets[i].HomingTurnRate)
				newAngle := currentAngle + turn
				sg.bullets[i].VX = math.Cos(newAngle) * sg.bullets[i].Speed
				sg.bullets[i].VY = math.Sin(newAngle) * sg.bullets[i].Speed
			}
		}

		sg.bullets[i].X += sg.bullets[i].VX
		sg.bullets[i].Y += sg.bullets[i].VY

		// 移除超出屏幕的子弹
		if sg.bullets[i].Y < 0 || sg.bullets[i].Y > 600 {
			sg.bullets[i].Active = false
		}
	}
}

func (sg *ShooterGame) findNearestEnemy(x, y float64) (float64, float64, bool) {
	minDist2 := math.MaxFloat64
	var tx, ty float64
	found := false
	for _, e := range sg.enemies {
		if !e.Active {
			continue
		}
		ex := e.X + e.Width/2
		ey := e.Y + e.Height/2
		dx := ex - x
		dy := ey - y
		d := dx*dx + dy*dy
		if d < minDist2 {
			minDist2 = d
			tx, ty = ex, ey
			found = true
		}
	}
	return tx, ty, found
}

// updateEnemies 更新敌机
func (sg *ShooterGame) updateEnemies() {
	for i := len(sg.enemies) - 1; i >= 0; i-- {
		if !sg.enemies[i].Active {
			sg.enemies = append(sg.enemies[:i], sg.enemies[i+1:]...)
			continue
		}

		sg.enemies[i].X += sg.enemies[i].VX
		sg.enemies[i].Y += sg.enemies[i].VY

		// 移除超出屏幕的敌机
		if sg.enemies[i].Y > 600 {
			sg.enemies[i].Active = false
		}
	}
}

// updateExplosions 更新爆炸效果
func (sg *ShooterGame) updateExplosions() {
	for i := len(sg.explosions) - 1; i >= 0; i-- {
		if !sg.explosions[i].Active {
			sg.explosions = append(sg.explosions[:i], sg.explosions[i+1:]...)
			continue
		}

		sg.explosions[i].Timer += 0.1
		sg.explosions[i].Radius = sg.explosions[i].MaxRadius * (1 - sg.explosions[i].Timer)

		if sg.explosions[i].Timer >= 1.0 {
			sg.explosions[i].Active = false
		}
	}
}

// initStars 初始化星星
func (sg *ShooterGame) initStars() {
	for range 100 {
		star := Star{
			X:     rand.Float64() * 800,
			Y:     rand.Float64() * 600,
			Speed: 1 + rand.Float64()*2, // 1-3 的随机速度
			Size:  1 + rand.Float64(),   // 1-2 的随机大小
		}
		sg.stars = append(sg.stars, star)
	}
}

// updateStars 更新星星位置
func (sg *ShooterGame) updateStars() {
	for i := range sg.stars {
		// 星星向下移动
		sg.stars[i].Y += sg.stars[i].Speed

		// 如果星星移出屏幕底部，重新从顶部开始
		if sg.stars[i].Y > 600 {
			sg.stars[i].Y = -sg.stars[i].Size
			sg.stars[i].X = rand.Float64() * 800
		}
	}
}

// adjustGMValue 调整 GM 面板的值
func (sg *ShooterGame) adjustGMValue(increase bool) {
	// 属性顺序：0 FireRateHz(float) 1 BulletsPerShot(int) 2 SpreadDeg(float)
	// 3 BulletSpeed(float) 4 BurstChance(float) 5 PenetrationCount(int)
	// 6 EnableHoming(bool) 7 HomingTurnRateRad(float)
	p := sg.player.Fire
	if increase {
		switch sg.gmIndex {
		case 0:
			p.FireRateHz = clamp(p.FireRateHz+0.5, 0.5, 30)
			sg.player.Fire = p
			sg.recomputeShotDelay()
		case 1:
			if p.BulletsPerShot < 20 {
				p.BulletsPerShot++
			}
		case 2:
			p.SpreadDeg = clamp(p.SpreadDeg+5, 0, 180)
		case 3:
			p.BulletSpeed = clamp(p.BulletSpeed+0.5, 1, 30)
		case 4:
			p.BurstChance = clamp(p.BurstChance+0.05, 0, 1)
		case 5:
			if p.PenetrationCount < 10 {
				p.PenetrationCount++
			}
		case 6:
			p.EnableHoming = !p.EnableHoming
		case 7:
			p.HomingTurnRateRad = clamp(p.HomingTurnRateRad+0.02, 0, 1)
		}
	} else {
		switch sg.gmIndex {
		case 0:
			p.FireRateHz = clamp(p.FireRateHz-0.5, 0.5, 30)
			sg.player.Fire = p
			sg.recomputeShotDelay()
		case 1:
			if p.BulletsPerShot > 1 {
				p.BulletsPerShot--
			}
		case 2:
			p.SpreadDeg = clamp(p.SpreadDeg-5, 0, 180)
		case 3:
			p.BulletSpeed = clamp(p.BulletSpeed-0.5, 1, 30)
		case 4:
			p.BurstChance = clamp(p.BurstChance-0.05, 0, 1)
		case 5:
			if p.PenetrationCount > 0 {
				p.PenetrationCount--
			}
		case 6:
			p.EnableHoming = !p.EnableHoming
		case 7:
			p.HomingTurnRateRad = clamp(p.HomingTurnRateRad-0.02, 0, 1)
		}
	}
	sg.player.Fire = p
}

// adjustGMSound 调整声效面板的值
func (sg *ShooterGame) adjustGMSound(increase bool) {
	cfg := sound.GetShootConfig()
	step := func(v, s, min, max float64, up bool) float64 {
		if up {
			return clamp(v+s, min, max)
		}
		return clamp(v-s, min, max)
	}
	if increase {
		switch sg.gmIndex {
		case 0:
			cfg.DurationSec = step(cfg.DurationSec, 0.01, 0.02, 0.5, true)
		case 1:
			cfg.BaseFreq = step(cfg.BaseFreq, 20, 50, 4000, true)
		case 2:
			cfg.MinFreq = step(cfg.MinFreq, 10, 20, cfg.BaseFreq, true)
		case 3:
			cfg.SweepFactor = step(cfg.SweepFactor, 0.02, 0, 1, true)
		case 4:
			cfg.Decay = step(cfg.Decay, 1, 1, 80, true)
		case 5:
			cfg.Amplitude = step(cfg.Amplitude, 0.05, 0, 1, true)
		case 6:
			// Waveform 轮换
			switch cfg.Waveform {
			case "square":
				cfg.Waveform = "triangle"
			case "triangle":
				cfg.Waveform = "noise"
			default:
				cfg.Waveform = "square"
			}
		}
	} else {
		switch sg.gmIndex {
		case 0:
			cfg.DurationSec = step(cfg.DurationSec, 0.01, 0.02, 0.5, false)
		case 1:
			cfg.BaseFreq = step(cfg.BaseFreq, 20, 50, 4000, false)
		case 2:
			cfg.MinFreq = step(cfg.MinFreq, 10, 20, cfg.BaseFreq, false)
		case 3:
			cfg.SweepFactor = step(cfg.SweepFactor, 0.02, 0, 1, false)
		case 4:
			cfg.Decay = step(cfg.Decay, 1, 1, 80, false)
		case 5:
			cfg.Amplitude = step(cfg.Amplitude, 0.05, 0, 1, false)
		case 6:
			// Waveform 轮换（相反方向）
			switch cfg.Waveform {
			case "square":
				cfg.Waveform = "noise"
			case "noise":
				cfg.Waveform = "triangle"
			default:
				cfg.Waveform = "square"
			}
		}
	}
	sound.SetShootConfig(cfg)
	sound.PlayShoot()
}
