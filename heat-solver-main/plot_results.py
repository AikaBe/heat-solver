import pandas as pd
import matplotlib.pyplot as plt
import numpy as np
import sys
import os

def plot_results(csv_file, method_name):
    """Plotting graphs from a CSV file"""

    if not os.path.exists(csv_file):
        print(f"Error: File {csv_file} not found!")
        return

    print(f"Loading data from {csv_file}...")
    df = pd.read_csv(csv_file)

    # --- Recompute error in Python, ignore any broken CSV error column ---
    for col in ["u_numeric", "u_exact"]:
        df[col] = pd.to_numeric(df[col], errors="coerce")

    df["error"] = np.abs(df["u_numeric"] - df["u_exact"])
    df.replace([np.inf, -np.inf], np.nan, inplace=True)

    times = sorted(df["t"].unique())
    print(f"Found {len(times)} time layers")

    fig = plt.figure(figsize=(16, 10))

    # ==================== Graph 1: Evolution of the solution ====================
    ax1 = plt.subplot(2, 2, 1)

    num_snapshots = min(6, len(times))
    snapshot_indices = np.linspace(0, len(times) - 1, num_snapshots, dtype=int)

    cmap = plt.cm.viridis
    colors = [cmap(i / max(1, num_snapshots - 1)) for i in range(num_snapshots)]

    for idx, t_idx in enumerate(snapshot_indices):
        t = times[t_idx]
        data = df[df["t"] == t].sort_values("x")
        ax1.plot(
            data["x"],
            data["u_numeric"],
            color=colors[idx],
            linewidth=2,
            label=f"t={t:.3f}",
        )

    ax1.set_xlabel("x", fontsize=12)
    ax1.set_ylabel("u(x,t)", fontsize=12)
    ax1.set_title(f"Evolution of the solution ({method_name})",
                  fontsize=14, fontweight="bold")
    ax1.legend(loc="best")
    ax1.grid(True, alpha=0.3)

    # ==================== Graph 2: Comparison with analytic ====================
    ax2 = plt.subplot(2, 2, 2)

    final_t = times[-1]
    final_data = df[df["t"] == final_t].sort_values("x")

    ax2.plot(
        final_data["x"],
        final_data["u_numeric"],
        "b-",
        linewidth=2,
        label="Numerical solution",
    )
    ax2.plot(
        final_data["x"],
        final_data["u_exact"],
        "r--",
        linewidth=2,
        label="Analytical solution",
    )
    ax2.set_xlabel("x", fontsize=12)
    ax2.set_ylabel("u(x,t)", fontsize=12)
    ax2.set_title(f"Comparison with t={final_t:.3f}",
                  fontsize=14, fontweight="bold")
    ax2.legend()
    ax2.grid(True, alpha=0.3)

    # ==================== Graph 3: Error distribution ====================
    ax3 = plt.subplot(2, 2, 3)

    valid = np.isfinite(final_data["error"])
    x_valid = final_data["x"][valid]
    err_valid = final_data["error"][valid]

    # ------- ONLY CHANGE YOU REQUESTED -------
    if len(err_valid) > 0:
        epsilon = 1e-20
        err_plot = np.maximum(err_valid, epsilon)

        ax3.plot(
            x_valid,
            err_plot,
            "g-",
            linewidth=2,
            marker="o",
            markersize=4,
        )
        ax3.set_yscale("log")

        # make tiny errors visible
        ax3.set_ylim(
            max(np.min(err_plot) * 0.5, epsilon),
            np.max(err_plot) * 2
        )

    else:
        ax3.text(
            0.5,
            0.5,
            "No valid error data",
            transform=ax3.transAxes,
            ha="center",
            va="center",
        )
    # ------- END OF CHANGE -------

    ax3.set_xlabel("x", fontsize=12)
    ax3.set_ylabel("Absolute error", fontsize=12)
    ax3.set_title(f"Error distribution at t={final_t:.3f}",
                  fontsize=14, fontweight="bold")
    ax3.grid(True, alpha=0.3)

    # ==================== Graph 4: Heatmap ====================
    ax4 = plt.subplot(2, 2, 4)

    x_vals = sorted(df["x"].unique())
    t_vals = sorted(df["t"].unique())

    step = max(1, len(t_vals) // 50)
    t_vals_subset = t_vals[::step]

    u_matrix = np.zeros((len(t_vals_subset), len(x_vals)))

    for i, t in enumerate(t_vals_subset):
        data = df[df["t"] == t].sort_values("x")
        row = pd.to_numeric(data["u_numeric"], errors="coerce").fillna(0.0)
        u_matrix[i, :] = row.values

    im = ax4.imshow(
        u_matrix,
        aspect="auto",
        origin="lower",
        extent=[min(x_vals), max(x_vals), min(t_vals_subset), max(t_vals_subset)],
        cmap="hot",
        interpolation="bilinear",
    )
    ax4.set_xlabel("x", fontsize=12)
    ax4.set_ylabel("t", fontsize=12)
    ax4.set_title("Temperature evolution (heat map)",
                  fontsize=14, fontweight="bold")
    cbar = plt.colorbar(im, ax=ax4)
    cbar.set_label("u(x,t)", fontsize=10)

    # ==================== Global error statistics ====================
    if len(err_valid) > 0:
        l2_error = float(np.sqrt(np.mean(err_valid**2)))
        linf_error = float(np.max(err_valid))
    else:
        l2_error = float("nan")
        linf_error = float("nan")

    fig.suptitle(
        f"Results of numerical solution of the heat equation\n"
        f"Method: {method_name} | L2 error: {l2_error:.3e} | L∞ error: {linf_error:.3e}",
        fontsize=16,
        fontweight="bold",
        y=0.98,
    )

    plt.tight_layout(rect=[0, 0, 1, 0.96])

    base, _ = os.path.splitext(csv_file)
    pdf_file = base + "_plot.pdf"
    png_file = base + "_plot.png"
    plt.savefig(pdf_file, bbox_inches="tight")
    plt.savefig(png_file, dpi=300, bbox_inches="tight")
    print(f"✓ Figures saved: {pdf_file} and {png_file}")

    plt.show()

    print("\n" + "=" * 60)
    print(f"ERROR STATISTICS (t={final_t:.3f}):")
    print(f"  L2 error:   {l2_error:.8e}")
    print(f"  L∞ error:   {linf_error:.8e}")
    print(f"  Mean error: {np.mean(err_valid) if len(err_valid)>0 else np.nan:.8e}")
    print("=" * 60)


# ---------- Comparison unchanged ----------
def compare_methods(csv_files, method_names):
    fig, axes = plt.subplots(2, 2, figsize=(16, 10))

    colors = ["blue", "red", "green", "orange"]
    markers = ["o", "s", "^", "d"]

    for idx, (csv_file, method) in enumerate(zip(csv_files, method_names)):
        if not os.path.exists(csv_file):
            print(f"WARNING: {csv_file} not found, skipped")
            continue

        df = pd.read_csv(csv_file)

        for col in ["u_numeric", "u_exact"]:
            df[col] = pd.to_numeric(df[col], errors="coerce")
        df["error"] = np.abs(df["u_numeric"] - df["u_exact"])
        df.replace([np.inf, -np.inf], np.nan, inplace=True)

        final_t = df["t"].max()
        final_data = df[df["t"] == final_t].sort_values("x")

        valid = np.isfinite(final_data["error"])
        x_valid = final_data["x"][valid]
        u_valid = final_data["u_numeric"][valid]
        err_valid = final_data["error"][valid]

        color = colors[idx]
        marker = markers[idx]

        axes[0, 0].plot(x_valid, u_valid, color=color, linewidth=2, label=method)
        axes[0, 1].plot(x_valid, err_valid, color=color, linewidth=2, marker=marker, markersize=4, label=method)

        l2_error = float(np.sqrt(np.mean(err_valid**2)))
        linf_error = float(np.max(err_valid))

        axes[1, 0].bar(idx, l2_error, color=color)
        axes[1, 1].bar(idx, linf_error, color=color)

    if csv_files:
        df0 = pd.read_csv(csv_files[0])
        df0["u_exact"] = pd.to_numeric(df0["u_exact"], errors="coerce")
        final_t0 = df0["t"].max()
        final_data0 = df0[df0["t"] == final_t0].sort_values("x")
        axes[0, 0].plot(
            final_data0["x"], final_data0["u_exact"], "k--", linewidth=2, label="Analytical"
        )

    axes[0, 0].set_xlabel("x")
    axes[0, 0].set_ylabel("u(x,t)")
    axes[0, 0].set_title("Solution comparison")
    axes[0, 0].legend()
    axes[0, 0].grid(True, alpha=0.3)

    axes[0, 1].set_xlabel("x")
    axes[0, 1].set_ylabel("Absolute error")
    axes[0, 1].set_title("Error comparison")
    axes[0, 1].legend()
    axes[0, 1].grid(True, alpha=0.3)
    axes[0, 1].set_yscale("log")

    axes[1, 0].set_ylabel("L2 error")
    axes[1, 0].set_title("Norm L2")
    axes[1, 0].set_xticks(range(len(method_names)))
    axes[1, 0].set_xticklabels(method_names)

    axes[1, 1].set_ylabel("L∞ error")
    axes[1, 1].set_title("Norm L∞")
    axes[1, 1].set_xticks(range(len(method_names)))
    axes[1, 1].set_xticklabels(method_names)

    plt.tight_layout()
    plt.savefig("methods_comparison.png", dpi=300, bbox_inches="tight")
    print("✓ Method Comparison saved: methods_comparison.png")
    plt.show()


if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Usage:")
        print("  python plot_results.py <csv_file> [method_name]")
        print("  python plot_results.py compare <file1> <method1> <file2> <method2> ...")
        sys.exit(1)

    if sys.argv[1] == "compare":
        if len(sys.argv) < 4 or len(sys.argv) % 2 != 0:
            print("ERROR: need pairs <file> <method> for comparison")
            sys.exit(1)

        csv_files = sys.argv[2::2]
        method_names = sys.argv[3::2]
        compare_methods(csv_files, method_names)
    else:
        csv_file = sys.argv[1]
        method_name = sys.argv[2] if len(sys.argv) > 2 else "Unknown"
        plot_results(csv_file, method_name)
