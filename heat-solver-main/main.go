package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math"
	"os"
	"strconv"
	"time"
)

type Params struct {
	method  string
	dx      float64
	dt      float64
	tmax    float64
	alpha   float64
	outfile string
}

// Analytical solution: u(x,t) = exp(-π² α t) * sin(π x)
func analyticalSolution(x, t, alpha float64) float64 {
	return math.Exp(-math.Pi*math.Pi*alpha*t) * math.Sin(math.Pi*x)
}

// Initial condition: u(x,0) = sin(π x)
func initialCondition(x float64) float64 {
	return math.Sin(math.Pi * x)
}

// Thomas algorithm for tridiagonal systems
// a: sub-diagonal (a[0] is ignored / should be 0)
// b: main diagonal
// c: super-diagonal (c[n-1] ignored / should be 0)
// d: right-hand side; returns x solving Ax = d.
func thomasAlgorithm(a, b, c, d []float64) []float64 {
	n := len(d)
	cp := make([]float64, n)
	dp := make([]float64, n)
	x := make([]float64, n)

	// modified coefficients
	cp[0] = c[0] / b[0]
	dp[0] = d[0] / b[0]

	for i := 1; i < n; i++ {
		denom := b[i] - a[i]*cp[i-1]
		cp[i] = 0.0
		if i < n-1 {
			cp[i] = c[i] / denom
		}
		dp[i] = (d[i] - a[i]*dp[i-1]) / denom
	}

	// back substitution
	x[n-1] = dp[n-1]
	for i := n - 2; i >= 0; i-- {
		x[i] = dp[i] - cp[i]*x[i+1]
	}

	return x
}

// FTCS (explicit, first-order in time, second-order in space, conditionally stable)
func solveFTCS(nx, nt int, dx, dt, alpha float64) [][]float64 {
	r := alpha * dt / (dx * dx)
	fmt.Printf("FTCS: r = %.4f\n", r)
	if r > 0.5 {
		fmt.Printf("Warning: FTCS may be unstable (r = %.4f > 0.5)\n", r)
	}

	u := make([][]float64, nt+1)
	for n := range u {
		u[n] = make([]float64, nx+1)
	}

	// initial condition at t = 0
	for i := 0; i <= nx; i++ {
		x := float64(i) * dx
		u[0][i] = initialCondition(x)
	}

	// time stepping
	for n := 0; n < nt; n++ {
		// enforce boundaries at current time level
		u[n][0] = 0.0
		u[n][nx] = 0.0

		// update interior
		for i := 1; i < nx; i++ {
			u[n+1][i] = u[n][i] + r*(u[n][i+1]-2*u[n][i]+u[n][i-1])
		}
	}

	// final boundaries
	u[nt][0] = 0.0
	u[nt][nx] = 0.0

	return u
}

// BTCS (implicit, unconditionally stable, first-order in time)
func solveBTCS(nx, nt int, dx, dt, alpha float64) [][]float64 {
	r := alpha * dt / (dx * dx)
	fmt.Printf("BTCS: r = %.4f\n", r)

	u := make([][]float64, nt+1)
	for n := range u {
		u[n] = make([]float64, nx+1)
	}

	// initial condition
	for i := 0; i <= nx; i++ {
		x := float64(i) * dx
		u[0][i] = initialCondition(x)
	}

	m := nx - 1 // number of interior points (i = 1..nx-1)

	a := make([]float64, m)
	b := make([]float64, m)
	c := make([]float64, m)
	d := make([]float64, m)

	// tridiagonal matrix for BTCS
	for j := 0; j < m; j++ {
		a[j] = -r
		b[j] = 1 + 2*r
		c[j] = -r
	}
	a[0] = 0.0      // no sub-diagonal in first row
	c[m-1] = 0.0    // no super-diagonal in last row

	for n := 0; n < nt; n++ {
		// boundary values at new time level
		u[n+1][0] = 0.0
		u[n+1][nx] = 0.0

		// build RHS from previous time level n
		for j := 0; j < m; j++ {
			i := j + 1 // physical index
			d[j] = u[n][i] // boundaries are zero, so no extra terms
		}

		sol := thomasAlgorithm(a, b, c, d)
		for j := 0; j < m; j++ {
			u[n+1][j+1] = sol[j]
		}
	}

	return u
}

// Crank–Nicolson (implicit, unconditionally stable, second-order in time)
func solveCrankNicolson(nx, nt int, dx, dt, alpha float64) [][]float64 {
	r := alpha * dt / (dx * dx)
	fmt.Printf("Crank–Nicolson: r = %.4f\n", r)

	u := make([][]float64, nt+1)
	for n := range u {
		u[n] = make([]float64, nx+1)
	}

	// initial condition
	for i := 0; i <= nx; i++ {
		x := float64(i) * dx
		u[0][i] = initialCondition(x)
	}

	m := nx - 1 // interior points

	a := make([]float64, m)
	b := make([]float64, m)
	c := make([]float64, m)
	d := make([]float64, m)

	// matrix A for CN: -r/2, 1+r, -r/2
	for j := 0; j < m; j++ {
		a[j] = -r / 2.0
		b[j] = 1.0 + r
		c[j] = -r / 2.0
	}
	a[0] = 0.0
	c[m-1] = 0.0

	for n := 0; n < nt; n++ {
		// boundaries at new time level
		u[n+1][0] = 0.0
		u[n+1][nx] = 0.0

		// build RHS d = B u^n
		for j := 0; j < m; j++ {
			i := j + 1 // physical index (1..nx-1)

			left := u[n][i-1]
			center := u[n][i]
			right := u[n][i+1]

			d[j] = (r/2.0)*left + (1.0-r)*center + (r/2.0)*right
		}

		sol := thomasAlgorithm(a, b, c, d)
		for j := 0; j < m; j++ {
			u[n+1][j+1] = sol[j]
		}
	}

	return u
}

// Compute L2 and L∞ error norms at final time
func computeErrors(u [][]float64, dx, dt, alpha float64) (float64, float64) {
	nt := len(u) - 1
	nx := len(u[0]) - 1
	t := float64(nt) * dt

	var sumSq, linf float64
	count := 0

	for i := 0; i <= nx; i++ {
		x := float64(i) * dx
		exact := analyticalSolution(x, t, alpha)
		num := u[nt][i]

		if math.IsNaN(num) || math.IsInf(num, 0) {
			continue
		}

		err := math.Abs(num - exact)
		sumSq += err * err
		if err > linf {
			linf = err
		}
		count++
	}

	if count == 0 {
		return math.NaN(), math.NaN()
	}

	l2 := math.Sqrt(sumSq / float64(count))
	return l2, linf
}

// Save full space–time solution to CSV
func saveToCSV(u [][]float64, dx, dt, alpha float64, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// header
	writer.Write([]string{"x", "t", "u_numeric", "u_exact", "error"})

	nt := len(u) - 1
	nx := len(u[0]) - 1

	for n := 0; n <= nt; n++ {
		t := float64(n) * dt
		for i := 0; i <= nx; i++ {
			x := float64(i) * dx

			num := u[n][i]
			exact := analyticalSolution(x, t, alpha)

			errVal := math.NaN()
			if !(math.IsNaN(num) || math.IsInf(num, 0)) {
				errVal = math.Abs(num - exact)
			}

			writer.Write([]string{
				strconv.FormatFloat(x, 'f', 6, 64),
				strconv.FormatFloat(t, 'f', 6, 64),
				strconv.FormatFloat(num, 'e', 8, 64),
				strconv.FormatFloat(exact, 'e', 8, 64),
				strconv.FormatFloat(errVal, 'e', 8, 64),
			})
		}
	}

	return nil
}

func main() {
	method := flag.String("method", "FTCS", "Numerical method: FTCS, BTCS, or CN")
	dx := flag.Float64("dx", 0.1, "Spatial step size")
	dt := flag.Float64("dt", 0.001, "Time step size")
	tmax := flag.Float64("tmax", 1.0, "Maximum simulation time")
	alpha := flag.Float64("alpha", 1.0, "Thermal diffusivity α")
	outfile := flag.String("out", "results.csv", "Output CSV file")

	flag.Parse()

	params := Params{
		method:  *method,
		dx:      *dx,
		dt:      *dt,
		tmax:    *tmax,
		alpha:   *alpha,
		outfile: *outfile,
	}

	// grid sizes (0..1)
	nx := int(1.0 / params.dx)
	nt := int(params.tmax / params.dt)

	fmt.Printf("Method: %s\n", params.method)
	fmt.Printf("dx = %.4f, dt = %.6f, tmax = %.3f, alpha = %.4f\n",
		params.dx, params.dt, params.tmax, params.alpha)
	fmt.Printf("Grid: nx = %d, nt = %d\n", nx, nt)

	var u [][]float64
	start := time.Now()

	switch params.method {
	case "FTCS":
		u = solveFTCS(nx, nt, params.dx, params.dt, params.alpha)
	case "BTCS":
		u = solveBTCS(nx, nt, params.dx, params.dt, params.alpha)
	case "CN":
		u = solveCrankNicolson(nx, nt, params.dx, params.dt, params.alpha)
	default:
		fmt.Printf("Unknown method: %s\n", params.method)
		fmt.Println("Available methods: FTCS, BTCS, CN")
		os.Exit(1)
	}

	elapsed := time.Since(start)

	l2Error, linfError := computeErrors(u, params.dx, params.dt, params.alpha)
	fmt.Printf("L2  error = %.8e\n", l2Error)
	fmt.Printf("L∞ error = %.8e\n", linfError)
	fmt.Printf("Runtime   = %.6fs\n", elapsed.Seconds())

	if err := saveToCSV(u, params.dx, params.dt, params.alpha, params.outfile); err != nil {
		fmt.Printf("Error saving results: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Results saved to %s\n", params.outfile)
}
