# Finite Difference Solution of the 1D Heat Equation in Go
*A comparative study of FTCS, BTCS, and Crank–Nicolson with a CLI, plots, and a tiny web demo.*

> Repository: https://github.com/AikaBe/heat-solver

This repository contains clean Go implementations of three classical finite-difference schemes for the 1D heat equation
\[ \(u_t = \alpha\,u_{xx}\) on \(x\in[0,L]\) with Dirichlet boundaries \],
plus utilities to generate CSV outputs and publication‑ready figures. A lightweight Go web server provides an interactive demo that calls the **same solvers** as the CLI.

The code accompanies the article:

> **A. Bekmanova. & A.Beisikeyva** *Finite Difference Solution of the One-Dimensional Heat Equation in Golang: A Comparative Study of FTCS, BTCS, and Crank–Nicolson.* arXiv:XXXX.XXXXX, 2025.

---

## Contents

```
heat-solver/
  cmd/
    head/        # CLI runner: FTCS / BTCS / CN → CSV
    server/      # Local web demo (http://localhost:8080)
  internal/
    config/      # Params (dx, dt, tmax, alpha, L, method, I/O)
    io/          # CSV writer (x, t, u, u_exact, error)
    solver/      # FTCS, BTCS, Crank–Nicolson, Thomas tridiagonal
    mathutils/   # Exact solution: e^{-π² α t} sin(π x), IC, helpers
  plots/
    plot_results.py  # Matplotlib script: figures as vector PDFs
  web/           # Minimal HTML/JS front end for the server
  go.mod         # Go module (>= 1.22 recommended)
  README.md      # You are here
```

---

## Quick start

### Requirements
- **Go** ≥ 1.22 with `go` on your PATH.  
- **Python** ≥ 3.9 (optional, for plots) with `matplotlib` and `pandas`:
  ```bash
  pip install matplotlib pandas numpy
  ```

### Build & test
```bash
go build ./...
go test  ./...    # includes Thomas solver tests and smoke tests for schemes
```

### Command‑line usage (CSV generation)
Run FTCS / BTCS / Crank–Nicolson and save results to CSV. All flags have defaults.

```bash
go run cmd/head/main.go \
  --method=FTCS \          # FTCS | BTCS | CN
  --dx=0.01 \              # spatial step
  --dt=0.0005 \            # time step
  --tmax=1.0 \             # final time
  --alpha=1.0 \            # thermal diffusivity α
  --L=1.0 \                # domain length [0, L]
  --out=ftcs.csv
```

CSV columns (English, stable):
```
x,t,u,u_exact,error
```

### Publication‑ready figures (vector PDFs)
Use the Python script to create the 4‑panel overview figure and a cross‑method comparison. It saves **vector PDFs** with embedded fonts (also PNGs for convenience).

```bash
# Single method overview (e.g., FTCS)
python plots/plot_results.py ftcs.csv FTCS
# → ftcs_plot.pdf  (vector)
# → ftcs_plot.png  (300 dpi)

# Compare several methods at the final time
python plots/plot_results.py compare ftcs.csv FTCS btcs.csv BTCS cn.csv CN
# → methods_comparison.pdf
# → methods_comparison.png
```

> The plotting script automatically selects time snapshots and is agnostic to `tmax`. The heat map uses raster tiles inside a vector PDF for reasonable file sizes. If you need fully vector heat maps, switch to `pcolormesh` in the script (may produce very large PDFs).

---

## Reproducing the paper’s results

### 1) FTCS stability demonstration
```bash
# r = α dt / dx^2
go run cmd/head/main.go --method=FTCS --dx=0.02 --dt=0.00019 --alpha=1 --L=1 --tmax=0.1 --out=ftcs_stable.csv   # r≈0.475
go run cmd/head/main.go --method=FTCS --dx=0.02 --dt=0.00021 --alpha=1 --L=1 --tmax=0.1 --out=ftcs_unstable.csv # r≈0.525

python plots/plot_results.py ftcs_stable.csv FTCS
python plots/plot_results.py ftcs_unstable.csv FTCS
```

### 2) Temporal convergence (dx fixed)
```bash
DX=0.01
for DT in 0.01 0.005 0.0025 0.00125; do
  go run cmd/head/main.go --method=BTCS --dx=$DX --dt=$DT --tmax=0.1 --out=btcs_dt_${DT}.csv
  go run cmd/head/main.go --method=CN   --dx=$DX --dt=$DT --tmax=0.1 --out=cn_dt_${DT}.csv
done
# compute slopes with your own script or a small notebook
```

### 3) Spatial convergence (dt very small)
```bash
DT=0.0001
for DX in 0.04 0.02 0.01 0.005; do
  go run cmd/head/main.go --method=BTCS --dx=$DX --dt=$DT --tmax=0.1 --out=btcs_dx_${DX}.csv
  go run cmd/head/main.go --method=CN   --dx=$DX --dt=$DT --tmax=0.1 --out=cn_dx_${DX}.csv
done
```

### 4) Simple timing comparison (cost)
```bash
go run cmd/head/main.go --method=FTCS --dx=0.005 --dt=0.00001 --tmax=0.5 --out=ftcs_timing.csv
go run cmd/head/main.go --method=BTCS --dx=0.005 --dt=0.0005  --tmax=0.5 --out=btcs_timing.csv
go run cmd/head/main.go --method=CN   --dx=0.005 --dt=0.0005  --tmax=0.5 --out=cn_timing.csv
# wall‑clock times are printed by the program; copy them into a small table
```

---

## Interactive web demo (optional)

Run the local server and open the browser. The server calls the same solvers as the CLI via the shared `internal/solver` package.

```bash
go run cmd/server/server.go
# open http://localhost:8080
```

API (informal): `POST /simulate` with JSON
```json
{
  "method": "CN",     // "FTCS" | "BTCS" | "CN"
  "dx": 0.01,
  "dt": 0.0005,
  "tmax": 1.0,
  "alpha": 1.0,
  "L": 1.0
}
```
Response: arrays `x` (space), `t` (selected times), and `u` (matrix [time][space]).

> The web demo is for pedagogy/visualization only; all results in the paper were regenerated from the CLI and plotted from CSVs.

---

## Notes on correctness and stability

- **Stability (FTCS):** requires \( r = \alpha\,\Delta t /\Delta x^2 \le 1/2 \).  
- **BTCS and CN:** unconditionally stable; expected temporal orders are 1 (BTCS) and 2 (CN). Both are second order in space.  
- **Tridiagonal solve:** implemented with a numerically stable Thomas algorithm (`internal/solver`). Unit tests validate residuals \( \|Ax-b\|_\infty \le 10^{-12} \).

---

## Troubleshooting

- **Go toolchain error**: ensure `go version` prints ≥ 1.22 and `go.mod` specifies a real version (e.g., `go 1.22`).  
- **Empty/odd plots**: verify CSV header is `x,t,u,u_exact,error`.  
- **FTCS “blows up”**: decrease `dt` or increase `dx` so that \( r \le 0.5 \).  
- **Non‑English headers**: use the provided I/O writer in `internal/io` which writes English headers.

---

## Reproducibility checklist

- [ ] `go build ./...` succeeds; `go test ./...` passes.  
- [ ] CSVs regenerated with the exact commands above.  
- [ ] Figures exported as **PDF** (vector) and included in the paper.  
- [ ] All parameters reported (method, dx, dt, L, α, \(r\), T).  
- [ ] Server demo tested locally (optional).

---

## How to cite and how to link in the paper

- **In‑paper code availability sentence (LaTeX):**
  ```latex
  The Go source code and plotting scripts are available at
  \url{https://github.com/AikaBe/heat-solver} (commit \texttt{<hash>}).
  ```
  Consider archiving a release on Zenodo and citing the DOI for permanence.

- **BibTeX for the article (example):**
  ```bibtex
  @misc{bekmanova2025heatgo,
    author       = {Aizhan Bekmanova},
    title        = {Finite Difference Solution of the One-Dimensional Heat Equation in Golang:
                    A Comparative Study of FTCS, BTCS, and Crank--Nicolson},
    year         = {2025},
    eprint       = {XXXX.XXXXX},
    archivePrefix= {arXiv},
    primaryClass = {math.NA},
    url          = {https://github.com/AikaBe/heat-solver}
  }
  ```

---

## License

See [`LICENSE`](./LICENSE). If unspecified, we recommend using MIT for teaching and reproducibility.

---

## Acknowledgments

This work was completed as part of an undergraduate project on numerical PDEs; thanks to open educational references on finite differences and to the Go community for a robust toolchain.
