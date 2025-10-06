package main

import (
	"flag"
	"log/slog"
	"os"
	"time"

	"heat-solver/internal/config"
	"heat-solver/internal/io"
	"heat-solver/internal/solver"
)

func main() {
	method := flag.String("method", "FTCS", "Numerical method: FTCS, BTCS, or CN")
	dx := flag.Float64("dx", 0.1, "Spatial step size")
	dt := flag.Float64("dt", 0.001, "Time step size")
	tmax := flag.Float64("tmax", 1.0, "Maximum simulation time")
	outfile := flag.String("out", "results.csv", "Output CSV file")

	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	params := config.Params{
		Method:  *method,
		Dx:      *dx,
		Dt:      *dt,
		Tmax:    *tmax,
		Outfile: *outfile,
	}

	nx := int(1.0 / params.Dx)
	nt := int(params.Tmax / params.Dt)

	slog.Info("Simulation parameters",
		"method", params.Method,
		"dx", params.Dx,
		"dt", params.Dt,
		"tmax", params.Tmax,
		"outfile", params.Outfile,
	)
	slog.Info("Grid configuration", "nx", nx, "nt", nt)

	start := time.Now()

	var u [][]float64

	switch params.Method {
	case "FTCS":
		u = solver.SolveFTCS(nx, nt, params.Dx, params.Dt)
	case "BTCS":
		u = solver.SolveBTCS(nx, nt, params.Dx, params.Dt)
	case "CN":
		u = solver.SolveCrankNicolson(nx, nt, params.Dx, params.Dt)
	default:
		slog.Error("Unknown method", "method", params.Method)
		os.Exit(1)
	}

	elapsed := time.Since(start)
	slog.Info("Computation completed", "runtime_sec", elapsed.Seconds())

	if err := io.SaveToCSV(u, params.Dx, params.Dt, params.Outfile); err != nil {
		slog.Error("Error saving results", "error", err)
		os.Exit(1)
	}

	slog.Info("Results successfully saved", "file", params.Outfile)
}
