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
	method string
	dx     float64
	dt     float64
	tmax   float64
	outfile string
}

// Аналитическое решение u(x,t) = exp(-π²t) * sin(πx)
func analyticalSolution(x, t float64) float64 {
	return math.Exp(-math.Pi*math.Pi*t) * math.Sin(math.Pi*x)
}

// Начальное условие u(x,0) = sin(πx)
func initialCondition(x float64) float64 {
	return math.Sin(math.Pi * x)
}

// FTCS (явная схема)
func solveFTCS(nx int, nt int, dx, dt float64) [][]float64 {
	r := dt / (dx * dx)
	
	if r > 0.5 {
		fmt.Printf("Warning: FTCS may be unstable! r = %.4f > 0.5\n", r)
	}
	
	u := make([][]float64, nt+1)
	for i := range u {
		u[i] = make([]float64, nx+1)
	}
	
	for i := 0; i <= nx; i++ {
		x := float64(i) * dx
		u[0][i] = initialCondition(x)
	}
	
	for n := 0; n <= nt; n++ {
		u[n][0] = 0.0
		u[n][nx] = 0.0
	}

	for n := 0; n < nt; n++ {
		for i := 1; i < nx; i++ {
			u[n+1][i] = u[n][i] + r*(u[n][i+1] - 2*u[n][i] + u[n][i-1])
		}
	}
	
	return u
}

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
	
	return x
}

func solveBTCS(nx int, nt int, dx, dt float64) [][]float64 {
	r := dt / (dx * dx)
	
	u := make([][]float64, nt+1)
	for i := range u {
		u[i] = make([]float64, nx+1)
	}
	
	for i := 0; i <= nx; i++ {
		x := float64(i) * dx
		u[0][i] = initialCondition(x)
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
	
	return u
}

func solveCrankNicolson(nx int, nt int, dx, dt float64) [][]float64 {
	r := dt / (dx * dx)
	
	u := make([][]float64, nt+1)
	for i := range u {
		u[i] = make([]float64, nx+1)
	}
	
	for i := 0; i <= nx; i++ {
		x := float64(i) * dx
		u[0][i] = initialCondition(x)
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
	
	return u
}

func computeErrors(u [][]float64, dx, dt float64) (float64, float64) {
	nt := len(u) - 1
	nx := len(u[0]) - 1
	
	var l2Error, linfError float64
	var sumSq float64
	
	t := float64(nt) * dt
	
	for i := 0; i <= nx; i++ {
		x := float64(i) * dx
		exact := analyticalSolution(x, t)
		err := math.Abs(u[nt][i] - exact)
		
		sumSq += err * err
		if err > linfError {
			linfError = err
		}
	}
	
	l2Error = math.Sqrt(sumSq / float64(nx+1))
	
	return l2Error, linfError
}

func saveToCSV(u [][]float64, dx, dt float64, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	
	writer := csv.NewWriter(file)
	defer writer.Flush()
	
	writer.Write([]string{"x", "t", "u_numeric", "u_exact", "error"})
	
	nt := len(u) - 1
	nx := len(u[0]) - 1
	
	for n := 0; n <= nt; n++ {
		t := float64(n) * dt
		for i := 0; i <= nx; i++ {
			x := float64(i) * dx
			exact := analyticalSolution(x, t)
			err := math.Abs(u[n][i] - exact)
			
			writer.Write([]string{
				strconv.FormatFloat(x, 'f', 6, 64),
				strconv.FormatFloat(t, 'f', 6, 64),
				strconv.FormatFloat(u[n][i], 'f', 6, 64),
				strconv.FormatFloat(exact, 'f', 6, 64),
				strconv.FormatFloat(err, 'f', 6, 64),
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
	outfile := flag.String("out", "results.csv", "Output CSV file")
	
	flag.Parse()
	
	params := Params{
		method:  *method,
		dx:      *dx,
		dt:      *dt,
		tmax:    *tmax,
		outfile: *outfile,
	}
	
	nx := int(1.0 / params.dx)
	nt := int(params.tmax / params.dt)
	
	fmt.Printf("Method: %s\n", params.method)
	fmt.Printf("dx = %.4f, dt = %.6f, steps = %d\n", params.dx, params.dt, nt)
	fmt.Printf("Grid size: nx = %d, nt = %d\n", nx, nt)
	
	var u [][]float64
	start := time.Now()
	
	switch params.method {
	case "FTCS":
		u = solveFTCS(nx, nt, params.dx, params.dt)
	case "BTCS":
		u = solveBTCS(nx, nt, params.dx, params.dt)
	case "CN":
		u = solveCrankNicolson(nx, nt, params.dx, params.dt)
	default:
		fmt.Printf("Unknown method: %s\n", params.method)
		fmt.Println("Available methods: FTCS, BTCS, CN")
		os.Exit(1)
	}
	
	elapsed := time.Since(start)
	
	l2Error, linfError := computeErrors(u, params.dx, params.dt)
	
	fmt.Printf("L2 error = %.6f\n", l2Error)
	fmt.Printf("Linf error = %.6f\n", linfError)
	fmt.Printf("Runtime = %.6fs\n", elapsed.Seconds())
	
	err := saveToCSV(u, params.dx, params.dt, params.outfile)
	if err != nil {
		fmt.Printf("Error saving results: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Printf("Results saved to %s\n", params.outfile)
}