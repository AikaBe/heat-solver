package io

import (
	"encoding/csv"
	"log/slog"
	"math"
	"os"
	"strconv"

	"heat-solver/internal/mathutils"
)

func SaveToCSV(u [][]float64, dx, dt float64, filename string) error {
	slog.Info("Saving results to CSV", "file", filename)

	file, err := os.Create(filename)
	if err != nil {
		slog.Error("Failed to create output file", "file", filename, "error", err)
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			slog.Warn("Failed to close file", "file", filename, "error", err)
		}
	}()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write([]string{"x", "t", "u_numeric", "u_exact", "error"}); err != nil {
		slog.Error("Failed to write CSV header", "error", err)
		return err
	}

	nt := len(u) - 1
	nx := len(u[0]) - 1

	slog.Info("Writing simulation results to CSV",
		"rows", (nt+1)*(nx+1),
		"nx", nx,
		"nt", nt,
	)

	for n := 0; n <= nt; n++ {
		t := float64(n) * dt
		for i := 0; i <= nx; i++ {
			x := float64(i) * dx
			exact := mathutils.AnalyticalSolution(x, t)
			errVal := math.Abs(u[n][i] - exact)

			if err := writer.Write([]string{
				strconv.FormatFloat(x, 'f', 6, 64),
				strconv.FormatFloat(t, 'f', 6, 64),
				strconv.FormatFloat(u[n][i], 'f', 6, 64),
				strconv.FormatFloat(exact, 'f', 6, 64),
				strconv.FormatFloat(errVal, 'f', 6, 64),
			}); err != nil {
				slog.Error("Failed to write CSV record", "row", n, "col", i, "error", err)
				return err
			}
		}
	}

	slog.Info("CSV file successfully written", "file", filename)
	return nil
}
