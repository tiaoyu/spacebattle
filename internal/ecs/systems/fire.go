package systems

import (
	"math"
	"math/rand"
	"time"

	"spacebattle/internal/ecs"
	"spacebattle/internal/ecs/components"
	"spacebattle/internal/ecs/tags"
	"spacebattle/internal/sound"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
)

// FireSystem 射击系统
type FireSystem struct {
	world *ecs.World
}

// NewFireSystem 创建射击系统
func NewFireSystem(world *ecs.World) *FireSystem {
	return &FireSystem{
		world: world,
	}
}

// Update 更新射击系统
func (s *FireSystem) Update(w donburi.World, firePressed bool) {
	// 处理计划中的连射
	s.ProcessScheduledShots(w)

	// 如果没有按下射击键，直接返回
	if !firePressed {
		return
	}

	// 查找玩家实体
	playerQuery := query.NewQuery(
		filter.Contains(tags.Player, components.Position, components.Size, components.FireSkill),
	)

	playerQuery.Each(w, func(entry *donburi.Entry) {
		pos := components.Position.Get(entry)
		size := components.Size.Get(entry)
		fireSkill := components.FireSkill.Get(entry)

		now := time.Now()

		// 检查射击冷却
		if now.Sub(fireSkill.LastShot) < fireSkill.ShotDelay {
			return
		}

		// 执行射击
		s.Fire(w, pos.X+size.Width/2-2, pos.Y, fireSkill)
		sound.PlayShoot()

		fireSkill.LastShot = now

		// 处理连发
		if fireSkill.BurstChance > 0 && rand.Float64() < fireSkill.BurstChance {
			burstTime := now.Add(fireSkill.BurstInterval)
			fireSkill.ScheduledShots = append(fireSkill.ScheduledShots, burstTime)
		}
	})
}

// Fire 执行射击
func (s *FireSystem) Fire(w donburi.World, x, y float64, fireSkill *components.FireSkillData) {
	numBullets := max(fireSkill.BulletsPerShot, 1)
	baseAngle := -math.Pi / 2
	spreadRad := fireSkill.SpreadDeg * math.Pi / 180.0

	for i := range numBullets {
		var angle float64
		if numBullets == 1 || spreadRad == 0 {
			angle = baseAngle
		} else {
			t := float64(i) / float64(numBullets-1)
			angle = baseAngle + (t-0.5)*spreadRad
		}

		vx := math.Cos(angle) * fireSkill.BulletSpeed
		vy := math.Sin(angle) * fireSkill.BulletSpeed

		s.world.CreateBullet(
			x, y, vx, vy,
			fireSkill.BulletSpeed,
			fireSkill.BulletDamage,
			fireSkill.PenetrationCount,
			fireSkill.EnableHoming,
			fireSkill.HomingTurnRateRad,
		)
	}
}

// ProcessScheduledShots 处理计划中的连射
func (s *FireSystem) ProcessScheduledShots(w donburi.World) {
	playerQuery := query.NewQuery(
		filter.Contains(tags.Player, components.Position, components.Size, components.FireSkill),
	)

	playerQuery.Each(w, func(entry *donburi.Entry) {
		pos := components.Position.Get(entry)
		size := components.Size.Get(entry)
		fireSkill := components.FireSkill.Get(entry)

		if len(fireSkill.ScheduledShots) == 0 {
			return
		}

		now := time.Now()
		idx := 0

		for idx < len(fireSkill.ScheduledShots) && fireSkill.ScheduledShots[idx].Before(now) {
			s.Fire(w, pos.X+size.Width/2-2, pos.Y, fireSkill)
			sound.PlayShoot()
			idx++
		}

		if idx > 0 {
			fireSkill.ScheduledShots = fireSkill.ScheduledShots[idx:]
		}
	})
}

// InitializeFireSkill 初始化射击技能
func InitializeFireSkill(fireSkill *components.FireSkillData) {
	fireSkill.ShotDelay = ecs.ComputeShotDelay(fireSkill.FireRateHz)
	fireSkill.LastShot = time.Now().Add(-fireSkill.ShotDelay) // 允许立即射击
	fireSkill.ScheduledShots = make([]time.Time, 0)
}

