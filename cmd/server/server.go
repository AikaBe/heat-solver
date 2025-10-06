package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"heat-solver/internal/config"
	"heat-solver/internal/solver"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("./web")))

	http.HandleFunc("/simulate", func(w http.ResponseWriter, r *http.Request) {
		method := r.URL.Query().Get("method")
		if method == "" {
			method = "FTCS"
		}
		dx, _ := strconv.ParseFloat(r.URL.Query().Get("dx"), 64)
		if dx == 0 {
			dx = 0.1
		}
		dt, _ := strconv.ParseFloat(r.URL.Query().Get("dt"), 64)
		if dt == 0 {
			dt = 0.001
		}
		tmax, _ := strconv.ParseFloat(r.URL.Query().Get("tmax"), 64)
		if tmax == 0 {
			tmax = 1.0
		}

		nx := int(1.0 / dx)
		nt := int(tmax / dt)

		params := config.Params{
			Method: method,
			Dx:     dx,
			Dt:     dt,
			Tmax:   tmax,
		}

		var u [][]float64
		switch params.Method {
		case "FTCS":
			u = solver.SolveFTCS(nx, nt, params.Dx, params.Dt)
		case "BTCS":
			u = solver.SolveBTCS(nx, nt, params.Dx, params.Dt)
		case "CN":
			u = solver.SolveCrankNicolson(nx, nt, params.Dx, params.Dt)
		default:
			http.Error(w, "Unknown method", http.StatusBadRequest)
			return
		}

		response := map[string]interface{}{
			"dx": params.Dx,
			"dt": params.Dt,
			"u":  u,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	log.Println("🚀 Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
