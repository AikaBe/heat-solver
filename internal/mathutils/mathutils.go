package mathutils

import (
	"log/slog"
	"math"
)

// Аналитическое решение u(x,t) = exp(-π²t) * sin(πx)
func AnalyticalSolution(x, t float64) float64 {
	result := math.Exp(-math.Pi*math.Pi*t) * math.Sin(math.Pi*x)
	slog.Debug("AnalyticalSolution computed", "x", x, "t", t, "u_exact", result)
	return result
}

// Начальное условие u(x,0) = sin(πx)
func InitialCondition(x float64) float64 {
	result := math.Sin(math.Pi * x)
	slog.Debug("InitialCondition computed", "x", x, "u0", result)
	return result
}
