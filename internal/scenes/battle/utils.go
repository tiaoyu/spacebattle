package battle

import "math"

// 工具函数

// wrapAngle 将角度限制在-π到π之间
func wrapAngle(a float64) float64 {
	for a > math.Pi {
		a -= 2 * math.Pi
	}
	for a < -math.Pi {
		a += 2 * math.Pi
	}
	return a
}

// clamp 限制值在最小值和最大值之间，返回值在min和max之间
func clamp(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

// ternary 三元运算符
func ternary(cond bool, a, b string) string {
	if cond {
		return a
	}
	return b
}
