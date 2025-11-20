package balance

import (
	"math"
	"time"
)

// RewardBreakdown 功勋奖励详细分解
type RewardBreakdown struct {
	BaseReward       int     // 基础奖励
	DifficultyBonus  int     // 难度加成
	KillBonus        int     // 击杀加成
	SpeedBonus       int     // 速度加成
	PerfectBonus     int     // 完美通关加成
	BossBonus        int     // Boss击杀加成
	TotalReward      int     // 总奖励
	PerformanceScore float64 // 综合表现分数 (0-100)
}

// ComputeMeritReward 计算通关后的功勋奖励（优化版）
// 新的奖励系统更加细化和有成就感：
// 1. 基础奖励：根据难度提供保底奖励
// 2. 难度加成：高难度提供更高的倍率加成
// 3. 击杀加成：根据击杀率提供额外奖励
// 4. 速度加成：快速通关提供额外奖励
// 5. 完美通关：全击杀+Boss击杀+快速通关的额外大幅加成
// 6. Boss加成：击杀Boss提供固定加成
func ComputeMeritReward(
	difficultyMul float64,
	killedCount int,
	spawnedCount int,
	elapsed time.Duration,
	total time.Duration,
	victory bool,
) int {
	breakdown := ComputeDetailedReward(difficultyMul, killedCount, spawnedCount, elapsed, total, victory)
	return breakdown.TotalReward
}

// ComputeDetailedReward 计算详细的功勋奖励分解
func ComputeDetailedReward(
	difficultyMul float64,
	killedCount int,
	spawnedCount int,
	elapsed time.Duration,
	total time.Duration,
	victory bool,
) RewardBreakdown {
	breakdown := RewardBreakdown{}

	if !victory {
		return breakdown // 失败不给奖励
	}

	// === 1. 基础奖励：与难度成本相关 ===
	cost := DifficultyCost(difficultyMul)
	if cost <= 0 {
		// 低难度给予最低基础奖励
		breakdown.BaseReward = 15
	} else {
		// 基础奖励 = 成本 * 1.5（保证回本并有盈余）
		breakdown.BaseReward = int(math.Round(float64(cost) * 1.5))
	}

	// === 2. 难度加成：高难度有更高的回报 ===
	// 难度越高，倍率越高（对数增长，避免爆炸）
	if difficultyMul > 1.0 {
		// 公式：baseReward * (0.5 * log2(difficulty))
		difficultyBonusRate := 0.5 * math.Log2(difficultyMul)
		breakdown.DifficultyBonus = int(math.Round(float64(breakdown.BaseReward) * difficultyBonusRate))
	}

	// === 3. 击杀加成：根据击杀率 ===
	killRatio := 0.0
	if spawnedCount > 0 {
		killRatio = float64(killedCount) / float64(spawnedCount)
		if killRatio > 1.0 {
			killRatio = 1.0
		}
	}

	// 击杀率 > 80% 开始有显著加成
	if killRatio >= 0.8 {
		killBonusRate := (killRatio - 0.7) / 0.3 // 70%-100% 映射到 0.33-1.0
		if killBonusRate > 1.0 {
			killBonusRate = 1.0
		}
		breakdown.KillBonus = int(math.Round(float64(breakdown.BaseReward) * 0.5 * killBonusRate))
	}

	// === 4. 速度加成：快速通关 ===
	speedRatio := 0.0
	if total > 0 && elapsed > 0 {
		// 剩余时间百分比
		remainingTime := total - elapsed
		if remainingTime > 0 {
			speedRatio = float64(remainingTime) / float64(total)
			if speedRatio > 1.0 {
				speedRatio = 1.0
			}
		}
	}

	// 剩余时间 > 20% 开始有加成
	if speedRatio > 0.2 {
		speedBonusRate := (speedRatio - 0.2) / 0.8 // 20%-100% 映射到 0-1.0
		breakdown.SpeedBonus = int(math.Round(float64(breakdown.BaseReward) * 0.3 * speedBonusRate))
	}

	// === 5. Boss加成：击杀Boss额外奖励 ===
	// 判断是否击杀了Boss（killedCount >= spawnedCount 说明全击杀包括Boss）
	bossKilled := killedCount >= spawnedCount && spawnedCount > 0
	if bossKilled {
		// Boss加成为基础奖励的 30%
		breakdown.BossBonus = int(math.Round(float64(breakdown.BaseReward) * 0.3))
	}

	// === 6. 完美通关加成：全击杀 + Boss击杀 + 速度快 ===
	// 条件：击杀率 >= 95%，速度比 > 30%，Boss已击杀
	if killRatio >= 0.95 && speedRatio > 0.3 && bossKilled {
		// 完美通关给予巨额加成（总和的 50%）
		currentTotal := breakdown.BaseReward + breakdown.DifficultyBonus + 
			breakdown.KillBonus + breakdown.SpeedBonus + breakdown.BossBonus
		breakdown.PerfectBonus = int(math.Round(float64(currentTotal) * 0.5))
	}

	// === 计算总奖励 ===
	breakdown.TotalReward = breakdown.BaseReward + 
		breakdown.DifficultyBonus + 
		breakdown.KillBonus + 
		breakdown.SpeedBonus + 
		breakdown.BossBonus + 
		breakdown.PerfectBonus

	// 确保最低奖励
	if breakdown.TotalReward < 10 {
		breakdown.TotalReward = 10
	}

	// === 计算综合表现分数（用于显示） ===
	breakdown.PerformanceScore = calculatePerformanceScore(killRatio, speedRatio, bossKilled)

	return breakdown
}

// calculatePerformanceScore 计算综合表现分数 (0-100)
func calculatePerformanceScore(killRatio, speedRatio float64, bossKilled bool) float64 {
	score := 0.0
	
	// 击杀率占 40 分
	score += killRatio * 40.0
	
	// 速度占 30 分
	score += speedRatio * 30.0
	
	// Boss击杀占 30 分
	if bossKilled {
		score += 30.0
	}
	
	// 钳制到 0-100
	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}
	
	return math.Round(score)
}
