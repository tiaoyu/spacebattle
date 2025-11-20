package balance

import (
	"spacebattle/internal/config"
	"math"
)

// DifficultyCost 依据当前难度倍率计算成本
// mul <= 1.0 时成本为 0
func DifficultyCost(mul float64) int {
	cfg := config.DefaultConfig()
	if mul <= 1.0 {
		return 0
	}
	base := max(1, cfg.DifficultyCostBase)
	// 渐进增益：在线性基础上乘以对数系数，使难度越高单次提升成本越大且不至于爆炸
	logFactor := 1.0 + 0.5*math.Log2(math.Max(1.0, mul))
	cost := int((mul - 1.0) * float64(base) * logFactor)
	if cost < 1 {
		cost = 1
	}
	return cost
}

// MaxAffordableDifficulty 根据可用功勋计算能支付的最大难度
// 使用二分搜索找到最大的难度使得 DifficultyCost(difficulty) <= merits
func MaxAffordableDifficulty(merits int) float64 {
	cfg := config.DefaultConfig()
	
	// 如果功勋为 0 或负数，只能选择最低难度
	if merits <= 0 {
		return cfg.DifficultyMin
	}
	
	// 二分搜索范围
	minDiff := cfg.DifficultyMin
	maxDiff := cfg.DifficultyMax
	
	// 如果最大难度都支付得起，直接返回
	if DifficultyCost(maxDiff) <= merits {
		return maxDiff
	}
	
	// 二分搜索，精度 0.01
	epsilon := 0.01
	for maxDiff - minDiff > epsilon {
		mid := (minDiff + maxDiff) / 2.0
		cost := DifficultyCost(mid)
		
		if cost <= merits {
			minDiff = mid
		} else {
			maxDiff = mid
		}
	}
	
	return minDiff
}
