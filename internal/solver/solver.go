package solver

import (
	"log/slog"
	"heat-solver/internal/mathutils"
)

// FTCS (явная схема)
func SolveFTCS(nx, nt int, dx, dt float64) [][]float64 {
	r := dt / (dx * dx)
	if r > 0.5 {
		slog.Warn("FTCS may be unstable", "r", r)
	} else {
		slog.Debug("FTCS stability check passed", "r", r)
	}

	slog.Info("Starting FTCS solver", "nx", nx, "nt", nt, "dx", dx, "dt", dt)

	u := make([][]float64, nt+1)
	for i := range u {
		u[i] = make([]float64, nx+1)
	}

	// Начальное условие
	for i := 0; i <= nx; i++ {
		x := float64(i) * dx
		u[0][i] = mathutils.InitialCondition(x)
	}

	// Граничные условия
	for n := 0; n <= nt; n++ {
		u[n][0] = 0.0
		u[n][nx] = 0.0
	}

	// Основной цикл
	for n := 0; n < nt; n++ {
		for i := 1; i < nx; i++ {
			u[n+1][i] = u[n][i] + r*(u[n][i+1]-2*u[n][i]+u[n][i-1])
		}
	}

	slog.Info("FTCS solver finished successfully")
	return u
}

// BTCS (неявная схема)
func SolveBTCS(nx, nt int, dx, dt float64) [][]float64 {
	r := dt / (dx * dx)
	slog.Info("Starting BTCS solver", "nx", nx, "nt", nt, "dx", dx, "dt", dt, "r", r)

	u := make([][]float64, nt+1)
	for i := range u {
		u[i] = make([]float64, nx+1)
	}

	for i := 0; i <= nx; i++ {
		x := float64(i) * dx
		u[0][i] = mathutils.InitialCondition(x)
	}

	for n := 0; n <= nt; n++ {
		u[n][0] = 0.0
		u[n][nx] = 0.0
	}

	a := make([]float64, nx-1)
	b := make([]float64, nx-1)
	c := make([]float64, nx-1)
	d := make([]float64, nx-1)

	for i := 0; i < nx-1; i++ {
		a[i] = -r
		b[i] = 1 + 2*r
		c[i] = -r
	}

	for n := 0; n < nt; n++ {
		for i := 0; i < nx-1; i++ {
			d[i] = u[n][i+1]
		}

		d[0] += r * u[n+1][0]
		d[nx-2] += r * u[n+1][nx]

		solution := thomasAlgorithm(a, b, c, d)
		for i := 0; i < nx-1; i++ {
			u[n+1][i+1] = solution[i]
		}
	}

	slog.Info("BTCS solver finished successfully")
	return u
}

// Crank–Nicolson (полуявная схема)
func SolveCrankNicolson(nx, nt int, dx, dt float64) [][]float64 {
	r := dt / (dx * dx)
	slog.Info("Starting Crank–Nicolson solver", "nx", nx, "nt", nt, "dx", dx, "dt", dt, "r", r)

	u := make([][]float64, nt+1)
	for i := range u {
		u[i] = make([]float64, nx+1)
	}

	for i := 0; i <= nx; i++ {
		x := float64(i) * dx
		u[0][i] = mathutils.InitialCondition(x)
	}

	for n := 0; n <= nt; n++ {
		u[n][0] = 0.0
		u[n][nx] = 0.0
	}

	a := make([]float64, nx-1)
	b := make([]float64, nx-1)
	c := make([]float64, nx-1)
	d := make([]float64, nx-1)

	for i := 0; i < nx-1; i++ {
		a[i] = -r / 2
		b[i] = 1 + r
		c[i] = -r / 2
	}

	for n := 0; n < nt; n++ {
		for i := 0; i < nx-1; i++ {
			d[i] = (r/2)*u[n][i] + (1-r)*u[n][i+1] + (r/2)*u[n][i+2]
		}
		d[0] += (r / 2) * u[n+1][0]
		d[nx-2] += (r / 2) * u[n+1][nx]

		solution := thomasAlgorithm(a, b, c, d)
		for i := 0; i < nx-1; i++ {
			u[n+1][i+1] = solution[i]
		}
	}

	slog.Info("Crank–Nicolson solver finished successfully")
	return u
}

// Алгоритм Томаса (метод прогонки)
func thomasAlgorithm(a, b, c, d []float64) []float64 {
	n := len(d)
	cp := make([]float64, n)
	dp := make([]float64, n)
	x := make([]float64, n)

	cp[0] = c[0] / b[0]
	dp[0] = d[0] / b[0]

	for i := 1; i < n; i++ {
		denom := b[i] - a[i]*cp[i-1]
		cp[i] = c[i] / denom
		dp[i] = (d[i] - a[i]*dp[i-1]) / denom
	}

	x[n-1] = dp[n-1]
	for i := n - 2; i >= 0; i-- {
		x[i] = dp[i] - cp[i]*x[i+1]
	}

	slog.Debug("Thomas algorithm executed", "n", n)
	return x
}
