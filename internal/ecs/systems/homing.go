package systems

import (
	"math"
	"time"

	"spacebattle/internal/ecs"
	"spacebattle/internal/ecs/components"
	"spacebattle/internal/ecs/tags"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
)

// HomingSystem 追踪系统
type HomingSystem struct{}

// NewHomingSystem 创建追踪系统
func NewHomingSystem() *HomingSystem {
	return &HomingSystem{}
}

// Update 更新追踪子弹
func (s *HomingSystem) Update(w donburi.World) {
	now := time.Now()
	
	// 先统计每个敌人被锁定的子弹数量
	targetCounts := make(map[donburi.Entity]int)
	
	homingQuery := query.NewQuery(
		filter.Contains(tags.Bullet, components.Homing),
	)
	
	homingQuery.Each(w, func(entry *donburi.Entry) {
		homing := components.Homing.Get(entry)
		if homing.TargetEntity == 0 {
			return
		}
		// 检查目标是否仍然存在
		if w.Valid(homing.TargetEntity) {
			targetCounts[homing.TargetEntity]++
		} else {
			// 目标已被摧毁，清除锁定
			homing.TargetEntity = 0
		}
	})

	// 查找所有追踪子弹并更新
	homingQuery.Each(w, func(entry *donburi.Entry) {
		pos := components.Position.Get(entry)
		vel := components.Velocity.Get(entry)
		homing := components.Homing.Get(entry)

		// 检查是否需要重新分配目标
		needRetarget := false
		if homing.TargetEntity == 0 {
			needRetarget = true
		} else if !w.Valid(homing.TargetEntity) {
			// 目标已被摧毁
			needRetarget = true
		} else if now.Sub(homing.LastRetargetTime) >= homing.RetargetInterval {
			// 定期重新评估目标（避免过度集中）
			needRetarget = true
		}

		if needRetarget {
			// 找到最优目标
			target, found := s.FindBestTarget(w, pos.X, pos.Y, targetCounts)
			if !found {
				return
			}
			
			// 更新旧目标计数
			if homing.TargetEntity != 0 && w.Valid(homing.TargetEntity) {
				if count, ok := targetCounts[homing.TargetEntity]; ok && count > 0 {
					targetCounts[homing.TargetEntity] = count - 1
				}
			}
			
			// 锁定新目标
			homing.TargetEntity = target
			homing.LastRetargetTime = now
			targetCounts[target]++
		}

		// 获取目标位置
		if homing.TargetEntity == 0 || !w.Valid(homing.TargetEntity) {
			return
		}

		targetEntry := w.Entry(homing.TargetEntity)
		if !targetEntry.HasComponent(components.Position) || !targetEntry.HasComponent(components.Size) {
			return
		}

		targetPos := components.Position.Get(targetEntry)
		targetSize := components.Size.Get(targetEntry)

		// 计算目标中心点
		tx := targetPos.X + targetSize.Width/2
		ty := targetPos.Y + targetSize.Height/2

		// 计算当前角度和目标角度
		currentAngle := math.Atan2(vel.VY, vel.VX)
		desiredAngle := math.Atan2(ty-pos.Y, tx-pos.X)

		// 计算角度差
		delta := ecs.WrapAngle(desiredAngle - currentAngle)

		// 限制转向速率
		turn := ecs.ClampFloat64(delta, -homing.TurnRate, homing.TurnRate)

		// 应用转向
		newAngle := currentAngle + turn

		// 更新速度方向
		vel.VX = math.Cos(newAngle) * homing.Speed
		vel.VY = math.Sin(newAngle) * homing.Speed
	})
}

// FindBestTarget 找到最优目标（综合考虑距离、威胁度和已锁定数量）
func (s *HomingSystem) FindBestTarget(w donburi.World, x, y float64, targetCounts map[donburi.Entity]int) (donburi.Entity, bool) {
	bestScore := math.Inf(-1)
	var bestTarget donburi.Entity
	found := false

	// 查找所有敌机和 Boss（使用 OR 组合）
	enemyQuery := query.NewQuery(
		filter.And(
			filter.Or(
				filter.Contains(tags.Enemy),
				filter.Contains(tags.EnemyShooter),
				filter.Contains(tags.EnemyZigzag),
				filter.Contains(tags.EnemyTank),
				filter.Contains(tags.Boss),
			),
			filter.Contains(components.Position, components.Size, components.Health),
		),
	)

	enemyQuery.Each(w, func(entry *donburi.Entry) {
		pos := components.Position.Get(entry)
		size := components.Size.Get(entry)
		health := components.Health.Get(entry)

		// 计算敌人中心点
		ex := pos.X + size.Width/2
		ey := pos.Y + size.Height/2

		// 计算距离
		dx := ex - x
		dy := ey - y
		dist := math.Sqrt(dx*dx + dy*dy)

		// 如果距离太远，跳过（避免追踪屏幕外的敌人）
		if dist > 800 {
			return
		}

		// 已锁定此目标的子弹数量
		lockCount := targetCounts[entry.Entity()]

		// 计算评分（分数越高越优先）
		// 1. 距离因子：距离越近越好（反比）
		distFactor := 1.0 / (1.0 + dist/100.0)

		// 2. 血量因子：血量越低越优先（快速清理残血）
		healthRatio := float64(health.Current) / float64(health.Max)
		healthFactor := 2.0 - healthRatio // 0-1 映射到 1-2

		// 3. 锁定惩罚：已经有很多子弹锁定的目标降低优先级
		lockPenalty := 1.0 / (1.0 + float64(lockCount)*0.5)

		// 4. Boss 加成：Boss 优先级更高
		bossBonus := 1.0
		if entry.HasComponent(tags.Boss) {
			bossBonus = 1.5
		}

		// 综合评分
		score := distFactor * healthFactor * lockPenalty * bossBonus

		if score > bestScore {
			bestScore = score
			bestTarget = entry.Entity()
			found = true
		}
	})

	return bestTarget, found
}

