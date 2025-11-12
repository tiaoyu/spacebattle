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
