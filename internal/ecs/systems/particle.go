package systems

import (
	"image/color"
	"math"
	"math/rand"

	"spacebattle/internal/config"
	"spacebattle/internal/ecs"
	"spacebattle/internal/ecs/components"
	"spacebattle/internal/ecs/tags"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
)

// ParticleSystem 粒子系统
type ParticleSystem struct {
	world *ecs.World
	cfg   *config.Config
}

// NewParticleSystem 创建粒子系统
func NewParticleSystem(world *ecs.World) *ParticleSystem {
	return &ParticleSystem{
		world: world,
		cfg:   config.DefaultConfig(),
	}
}

// Update 更新粒子
func (s *ParticleSystem) Update(w donburi.World, dt float64) {
	particleQuery := query.NewQuery(
		filter.Contains(tags.Particle, components.Position, components.Particle),
	)

	var toRemove []*donburi.Entry

	particleQuery.Each(w, func(entry *donburi.Entry) {
		pos := components.Position.Get(entry)
		particle := components.Particle.Get(entry)

		// 更新位置
		pos.X += particle.VX * dt * 60 // 60 FPS 标准化
		pos.Y += particle.VY * dt * 60

		// 更新生命周期
		particle.Life -= dt
		if particle.Life <= 0 {
			toRemove = append(toRemove, entry)
			return
		}

		// 更新透明度（淡出效果）
		ratio := particle.Life / particle.MaxLife
		particle.Alpha = uint8(ratio * 255)
	})

	// 清理死亡粒子
	for _, entry := range toRemove {
		w.Remove(entry.Entity())
	}
}

// Draw 绘制粒子
func (s *ParticleSystem) Draw(w donburi.World, screen *ebiten.Image) {
	particleQuery := query.NewQuery(
		filter.Contains(tags.Particle, components.Position, components.Particle),
	)

	particleQuery.Each(w, func(entry *donburi.Entry) {
		pos := components.Position.Get(entry)
		particle := components.Particle.Get(entry)

		// 绘制粒子（小方块）
		col := color.RGBA{
			R: particle.ColorR,
			G: particle.ColorG,
			B: particle.ColorB,
			A: particle.Alpha,
		}

		vector.DrawFilledRect(
			screen,
			float32(pos.X),
			float32(pos.Y),
			float32(particle.Size),
			float32(particle.Size),
			col,
			false,
		)
	})
}

// CreateExplosionParticles 创建爆炸粒子效果
func (s *ParticleSystem) CreateExplosionParticles(x, y float64) {
	// 检查粒子数量限制
	count := 0
	query.NewQuery(filter.Contains(tags.Particle)).Each(s.world.ECS.World, func(entry *donburi.Entry) {
		count++
	})

	if count >= s.cfg.ParticleMaxCount {
		return // 达到上限，不再生成
	}

	numParticles := s.cfg.ParticleCountOnKill
	for i := 0; i < numParticles; i++ {
		// 随机速度方向
		angle := rand.Float64() * 2 * math.Pi
		speed := rand.Float64()*2 + 1
		vx := speed * math.Cos(angle)
		vy := speed * math.Sin(angle)

		// 随机大小
		size := rand.Float64()*2 + 2

		// 使用配置中的粒子颜色
		s.world.CreateParticle(x, y, vx, vy, size, s.cfg.ParticleColor.R, s.cfg.ParticleColor.G, s.cfg.ParticleColor.B)
	}
}
