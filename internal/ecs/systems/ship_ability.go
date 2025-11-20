package systems

import (
	"time"

	"spacebattle/internal/config"
	"spacebattle/internal/ecs/components"
	"spacebattle/internal/ecs/tags"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
)

// ShipAbilitySystem 战机被动技能系统
type ShipAbilitySystem struct {
	cfg *config.Config
}

// NewShipAbilitySystem 创建战机被动系统
func NewShipAbilitySystem() *ShipAbilitySystem {
	return &ShipAbilitySystem{
		cfg: config.DefaultConfig(),
	}
}

// Update 更新被动技能状态
func (s *ShipAbilitySystem) Update(w donburi.World, dt float64) {
	// 查找玩家实体
	playerQuery := query.NewQuery(
		filter.Contains(tags.Player, components.ShipAbility, components.Health, components.FireSkill),
	)

	playerQuery.Each(w, func(entry *donburi.Entry) {
		ability := components.ShipAbility.Get(entry)

		switch ability.AbilityType {
		case "speed_frenzy":
			// Beta - 速度狂热：检查buff是否过期
			if ability.FrenzyStacks > 0 {
				elapsed := time.Since(ability.LastKillTime).Seconds()
				if elapsed > s.cfg.FrenzyStackDuration {
					// buff过期，清空所有层数
					ability.FrenzyStacks = 0
				}
			}

		case "dodge_master":
			// Gamma - 闪避大师：更新无敌时间和冷却
			if ability.IsInvulnerable {
				ability.InvulnTime -= dt
				if ability.InvulnTime <= 0 {
					ability.IsInvulnerable = false
					ability.InvulnTime = 0
				}
			}
			if ability.InvulnCooldown > 0 {
				ability.InvulnCooldown -= dt
				if ability.InvulnCooldown < 0 {
					ability.InvulnCooldown = 0
				}
			}

		case "energy_shield":
			// Delta - 能量护盾：5秒不受伤害后回复护盾
			if ability.ShieldCurrent < ability.ShieldMax {
				elapsed := time.Since(ability.LastDamageTime).Seconds()
				if elapsed >= s.cfg.ShieldRegenDelay {
					// 每秒回复1点护盾
					ability.ShieldCurrent += int(dt)
					if ability.ShieldCurrent > ability.ShieldMax {
						ability.ShieldCurrent = ability.ShieldMax
					}
				}
			}
		}
	})
}

// OnEnemyKilled 敌机被击杀时触发
func (s *ShipAbilitySystem) OnEnemyKilled(w donburi.World, playerEntry *donburi.Entry) {
	if !playerEntry.HasComponent(components.ShipAbility) {
		return
	}

	ability := components.ShipAbility.Get(playerEntry)
	health := components.Health.Get(playerEntry)

	switch ability.AbilityType {
	case "harvest":
		// Alpha - 战场收割：每5次击杀回复1点生命
		ability.KillCounter++
		if ability.KillCounter >= s.cfg.HarvestKillsRequired {
			if health.Current < health.Max {
				health.Current++
			}
			ability.KillCounter = 0
		}

	case "speed_frenzy":
		// Beta - 速度狂热：叠加射速buff
		ability.LastKillTime = time.Now()
		if ability.FrenzyStacks < s.cfg.FrenzyMaxStacks {
			ability.FrenzyStacks++
		}
	}
}

// OnPlayerDamaged 玩家受到伤害时触发
func (s *ShipAbilitySystem) OnPlayerDamaged(w donburi.World, playerEntry *donburi.Entry) int {
	if !playerEntry.HasComponent(components.ShipAbility) {
		return 1 // 默认伤害
	}

	ability := components.ShipAbility.Get(playerEntry)

	switch ability.AbilityType {
	case "dodge_master":
		// Gamma - 闪避大师：检查是否处于无敌状态
		if ability.IsInvulnerable {
			return 0 // 无敌，不受伤害
		}
		// 触发无敌效果（如果冷却完成）
		if ability.InvulnCooldown <= 0 {
			ability.IsInvulnerable = true
			ability.InvulnTime = s.cfg.DodgeInvulnDuration
			ability.InvulnCooldown = s.cfg.DodgeInvulnCooldown
			ability.LastDamageTaken = time.Now()
		}
		return 1

	case "energy_shield":
		// Delta - 能量护盾：先扣除护盾
		ability.LastDamageTime = time.Now()
		if ability.ShieldCurrent > 0 {
			ability.ShieldCurrent--
			return 0 // 护盾吸收了伤害
		}
		return 1 // 护盾耗尽，扣除生命
	}

	return 1 // 默认伤害
}

// GetFireRateMultiplier 获取射速加成倍率
func (s *ShipAbilitySystem) GetFireRateMultiplier(playerEntry *donburi.Entry) float64 {
	if !playerEntry.HasComponent(components.ShipAbility) {
		return 1.0
	}

	ability := components.ShipAbility.Get(playerEntry)
	if ability.AbilityType == "speed_frenzy" && ability.FrenzyStacks > 0 {
		return 1.0 + float64(ability.FrenzyStacks)*s.cfg.FrenzySpeedBonus
	}

	return 1.0
}

// IsInvulnerable 检查玩家是否无敌
func (s *ShipAbilitySystem) IsInvulnerable(playerEntry *donburi.Entry) bool {
	if !playerEntry.HasComponent(components.ShipAbility) {
		return false
	}

	ability := components.ShipAbility.Get(playerEntry)
	return ability.AbilityType == "dodge_master" && ability.IsInvulnerable
}

