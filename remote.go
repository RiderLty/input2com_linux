package main

import (
	"math"
	"math/rand"
	"time"
)

// 辅助函数：计算绝对值
func abs32(x int32) int32 {
	if x < 0 {
		return -x
	}
	return x
}

// GeneratePullTrajectory (保持原逻辑不变，这是平滑的核心)
func GeneratePullTrajectory(totalDisp float64, totalFrames int, noiseStd, easeAmt float64) []float64 {
	if totalFrames <= 1 {
		return []float64{totalDisp}
	}
	sign := 1.0
	if totalDisp < 0 {
		sign = -1.0
		totalDisp = -totalDisp
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	ease := func(t float64) float64 {
		if easeAmt <= 0 {
			return t
		}
		if easeAmt >= 1 {
			return t * t * (3 - 2*t)
		}
		return t*(1.0-easeAmt) + (t*t*(3-2*t))*easeAmt
	}

	pos := make([]float64, totalFrames)
	for i := 0; i < totalFrames; i++ {
		t := float64(i) / float64(totalFrames-1)
		pos[i] = ease(t) * totalDisp
	}

	deltas := make([]float64, totalFrames)
	prev := 0.0
	for i := 0; i < totalFrames; i++ {
		deltas[i] = pos[i] - prev
		prev = pos[i]
	}

	minPositive := totalDisp / float64(totalFrames) * 1e-3
	for i := 0; i < totalFrames; i++ {
		if noiseStd > 0 {
			u1, u2 := r.Float64(), r.Float64()
			if u1 < 1e-12 {
				u1 = 1e-12
			}
			z0 := math.Sqrt(-2.0*math.Log(u1)) * math.Cos(2*math.Pi*u2)
			deltas[i] += z0 * noiseStd
		}
		if deltas[i] < minPositive {
			deltas[i] = minPositive
		}
	}

	sum := 0.0
	for _, v := range deltas {
		sum += v
	}
	scale := totalDisp / sum
	for i := range deltas {
		deltas[i] *= scale
	}

	for i := range deltas {
		deltas[i] *= sign
	}
	return deltas
}
