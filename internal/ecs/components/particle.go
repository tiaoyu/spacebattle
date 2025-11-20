package components

import "github.com/yohamta/donburi"

// ParticleData 粒子效果数据
type ParticleData struct {
	VX, VY    float64 // 速度向量
	Life      float64 // 剩余生命时间（秒）
	MaxLife   float64 // 最大生命时间
	Size      float64 // 粒子大小
	DecayRate float64 // 衰减速率
	ColorR    uint8   // 红色分量
	ColorG    uint8   // 绿色分量
	ColorB    uint8   // 蓝色分量
	Alpha     uint8   // 透明度
}

// Particle 粒子组件
var Particle = donburi.NewComponentType[ParticleData]()

