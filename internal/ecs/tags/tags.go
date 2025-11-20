package tags

import "github.com/yohamta/donburi"

// 标签组件用于标识实体类型

// PlayerTag 玩家标记
type PlayerTag struct{}

var Player = donburi.NewTag()

// EnemyTag 敌机标记
type EnemyTag struct{}

var Enemy = donburi.NewTag()

// BossTag Boss 标记
type BossTag struct{}

var Boss = donburi.NewTag()

// BulletTag 子弹标记
type BulletTag struct{}

var Bullet = donburi.NewTag()

// ExplosionTag 爆炸效果标记
type ExplosionTag struct{}

var Explosion = donburi.NewTag()

// StarTag 星星背景标记
type StarTag struct{}

var Star = donburi.NewTag()

// EnemyShooterTag 射击型敌机标记
type EnemyShooterTag struct{}

var EnemyShooter = donburi.NewTag()

// EnemyZigzagTag 之字型敌机标记
type EnemyZigzagTag struct{}

var EnemyZigzag = donburi.NewTag()

// EnemyTankTag 肉盾型敌机标记
type EnemyTankTag struct{}

var EnemyTank = donburi.NewTag()

// EnemyBulletTag 敌机子弹标记
type EnemyBulletTag struct{}

var EnemyBullet = donburi.NewTag()

// ParticleTagType 粒子标记
type ParticleTagType struct{}

var Particle = donburi.NewTag()

