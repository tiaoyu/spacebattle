package balance

import (
	"math"
	"time"
)

// ComputeMeritReward 计算通关后的功勋奖励。
// - 基于难度花费 cost = DifficultyCost(difficultyMul)
// - 考虑击杀率与通关速度两个维度，映射为 [0,1] 的绩效分
// - 最终奖励 = cost * (2.0 + performance)，并钳制在 [2x, 3x]
// - 非胜利时返回 0
func ComputeMeritReward(
	difficultyMul float64,
	killedCount int,
	spawnedCount int,
	elapsed time.Duration,
	total time.Duration,
	victory bool,
) int {
	if !victory {
		return 0
	}
	cost := DifficultyCost(difficultyMul)
	if cost <= 0 {
		return 10 // 难度花费为 0 时默认给最低奖励
	}

	// 击杀率：考虑总生成数量（含上限），钳制 0..1
	killRatio := 0.0
	if spawnedCount > 0 {
		killRatio = float64(killedCount) / float64(spawnedCount)
		if killRatio < 0 {
			killRatio = 0
		}
		if killRatio > 1 {
			killRatio = 1
		}
	}

	// 时间效率：越快越好。elapsed 可能略大于 total（硬收束后），钳制 0..1
	timeFactor := 0.0
	if total > 0 {
		t := (total.Seconds() - elapsed.Seconds()) / total.Seconds()
		if t < 0 {
			t = 0
		}
		if t > 1 {
			t = 1
		}
		timeFactor = t
	}

	// 简单平均为绩效分，可按需要调整权重
	performance := 0.5*killRatio + 0.5*timeFactor
	// 将绩效分映射到 [2,3]
	multiplier := 2.0 + performance
	if multiplier < 2 {
		multiplier = 2
	}
	if multiplier > 3 {
		multiplier = 3
	}

	reward := max(int(math.Round(float64(cost)*multiplier)), 0)
	return reward
}
