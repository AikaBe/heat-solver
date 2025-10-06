import pandas as pd
import matplotlib.pyplot as plt
import numpy as np
import sys
import os

def plot_results(csv_file, method_name):
    """Построение графиков из CSV файла"""
    
    if not os.path.exists(csv_file):
        print(f"Ошибка: файл {csv_file} не найден!")
        return
    
    print(f"Загрузка данных из {csv_file}...")
    df = pd.read_csv(csv_file)

    times = sorted(df['t'].unique())
    print(f"Найдено {len(times)} временных слоёв")
    
    fig = plt.figure(figsize=(16, 10))
    
    # ==================== График 1: Эволюция решения ====================
    ax1 = plt.subplot(2, 2, 1)
    
    num_snapshots = min(6, len(times))
    snapshot_indices = np.linspace(0, len(times)-1, num_snapshots, dtype=int)
    
    cmap = plt.cm.viridis
    colors = [cmap(i/num_snapshots) for i in range(num_snapshots)]
    
    for idx, t_idx in enumerate(snapshot_indices):
        t = times[t_idx]
        data = df[df['t'] == t]
        ax1.plot(data['x'], data['u_numeric'], 
                color=colors[idx], linewidth=2, 
                label=f't={t:.3f}')
    
    ax1.set_xlabel('x', fontsize=12)
    ax1.set_ylabel('u(x,t)', fontsize=12)
    ax1.set_title(f'Эволюция решения ({method_name})', fontsize=14, fontweight='bold')
    ax1.legend(loc='best')
    ax1.grid(True, alpha=0.3)
    
    # ==================== График 2: Сравнение с аналитикой ====================
    ax2 = plt.subplot(2, 2, 2)
    
    final_t = times[-1]
    final_data = df[df['t'] == final_t].sort_values('x')
    
    ax2.plot(final_data['x'], final_data['u_numeric'], 
            'b-', linewidth=2, label='Численное решение')
    ax2.plot(final_data['x'], final_data['u_exact'], 
            'r--', linewidth=2, label='Аналитическое решение')
    ax2.set_xlabel('x', fontsize=12)
    ax2.set_ylabel('u(x,t)', fontsize=12)
    ax2.set_title(f'Сравнение при t={final_t:.3f}', fontsize=14, fontweight='bold')
    ax2.legend()
    ax2.grid(True, alpha=0.3)
    
    # ==================== График 3: Распределение ошибки ====================
    ax3 = plt.subplot(2, 2, 3)
    
    ax3.plot(final_data['x'], final_data['error'], 
            'g-', linewidth=2, marker='o', markersize=4)
    ax3.set_xlabel('x', fontsize=12)
    ax3.set_ylabel('Абсолютная ошибка', fontsize=12)
    ax3.set_title(f'Распределение ошибки при t={final_t:.3f}', fontsize=14, fontweight='bold')
    ax3.grid(True, alpha=0.3)
    ax3.set_yscale('log')
    
    # ==================== График 4: Тепловая карта (heatmap) ====================
    ax4 = plt.subplot(2, 2, 4)
    
    x_vals = sorted(df['x'].unique())
    t_vals = sorted(df['t'].unique())
  
    step = max(1, len(t_vals) // 50)
    t_vals_subset = t_vals[::step]
    
    u_matrix = np.zeros((len(t_vals_subset), len(x_vals)))
    
    for i, t in enumerate(t_vals_subset):
        data = df[df['t'] == t].sort_values('x')
        u_matrix[i, :] = data['u_numeric'].values
    
    im = ax4.imshow(u_matrix, aspect='auto', origin='lower', 
                    extent=[0, 1, 0, final_t], cmap='hot', interpolation='bilinear')
    ax4.set_xlabel('x', fontsize=12)
    ax4.set_ylabel('t', fontsize=12)
    ax4.set_title('Эволюция температуры (тепловая карта)', fontsize=14, fontweight='bold')
    cbar = plt.colorbar(im, ax=ax4)
    cbar.set_label('u(x,t)', fontsize=10)
    
    # ==================== Общая информация ====================
    l2_error = np.sqrt(np.mean(final_data['error']**2))
    linf_error = final_data['error'].max()
    
    fig.suptitle(f'Результаты численного решения уравнения теплопроводности\n'
                 f'Метод: {method_name} | L2 error: {l2_error:.6f} | L∞ error: {linf_error:.6f}',
                 fontsize=16, fontweight='bold', y=0.98)
    
    plt.tight_layout(rect=[0, 0, 1, 0.96])
    
    output_file = csv_file.replace('.csv', '_plot.png')
    plt.savefig(output_file, dpi=300, bbox_inches='tight')
    print(f"✓ График сохранён: {output_file}")
    
    plt.show()
    
    print("\n" + "="*60)
    print(f"СТАТИСТИКА ОШИБОК (t={final_t:.3f}):")
    print(f"  L2 error:   {l2_error:.8f}")
    print(f"  L∞ error:   {linf_error:.8f}")
    print(f"  Mean error: {final_data['error'].mean():.8f}")
    print("="*60)

def compare_methods(csv_files, method_names):
    """Сравнение нескольких методов на одном графике"""
    
    fig, axes = plt.subplots(2, 2, figsize=(16, 10))
    
    colors = ['blue', 'red', 'green', 'orange']
    markers = ['o', 's', '^', 'd']
    
    for idx, (csv_file, method) in enumerate(zip(csv_files, method_names)):
        if not os.path.exists(csv_file):
            print(f"Предупреждение: {csv_file} не найден, пропускаем")
            continue
            
        df = pd.read_csv(csv_file)
        final_t = df['t'].max()
        final_data = df[df['t'] == final_t].sort_values('x')
        
        color = colors[idx % len(colors)]
        marker = markers[idx % len(markers)]
        
        axes[0, 0].plot(final_data['x'], final_data['u_numeric'], 
                       color=color, linewidth=2, label=method)
        
        axes[0, 1].plot(final_data['x'], final_data['error'], 
                       color=color, linewidth=2, marker=marker, 
                       markersize=4, label=method)
        
        l2_error = np.sqrt(np.mean(final_data['error']**2))
        axes[1, 0].bar(idx, l2_error, color=color, label=method)
        
        linf_error = final_data['error'].max()
        axes[1, 1].bar(idx, linf_error, color=color, label=method)
    
    if csv_files:
        df = pd.read_csv(csv_files[0])
        final_t = df['t'].max()
        final_data = df[df['t'] == final_t].sort_values('x')
        axes[0, 0].plot(final_data['x'], final_data['u_exact'], 
                       'k--', linewidth=2, label='Аналитическое')
    
    axes[0, 0].set_xlabel('x')
    axes[0, 0].set_ylabel('u(x,t)')
    axes[0, 0].set_title('Сравнение решений')
    axes[0, 0].legend()
    axes[0, 0].grid(True, alpha=0.3)
    
    axes[0, 1].set_xlabel('x')
    axes[0, 1].set_ylabel('Абсолютная ошибка')
    axes[0, 1].set_title('Сравнение ошибок')
    axes[0, 1].legend()
    axes[0, 1].grid(True, alpha=0.3)
    axes[0, 1].set_yscale('log')
    
    axes[1, 0].set_ylabel('L2 error')
    axes[1, 0].set_title('Норма L2')
    axes[1, 0].set_xticks(range(len(method_names)))
    axes[1, 0].set_xticklabels(method_names)
    axes[1, 0].grid(True, alpha=0.3, axis='y')
    
    axes[1, 1].set_ylabel('L∞ error')
    axes[1, 1].set_title('Норма L∞')
    axes[1, 1].set_xticks(range(len(method_names)))
    axes[1, 1].set_xticklabels(method_names)
    axes[1, 1].grid(True, alpha=0.3, axis='y')
    
    plt.tight_layout()
    plt.savefig('methods_comparison.png', dpi=300, bbox_inches='tight')
    print("✓ Сравнение методов сохранено: methods_comparison.png")
    plt.show()

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Использование:")
        print("  python plot_results.py <csv_file> [method_name]")
        print("  python plot_results.py compare <file1> <method1> <file2> <method2> ...")
        print("\nПримеры:")
        print("  python plot_results.py results.csv FTCS")
        print("  python plot_results.py compare ftcs.csv FTCS btcs.csv BTCS cn.csv CN")
        sys.exit(1)
    
    if sys.argv[1] == "compare":
        if len(sys.argv) < 4 or len(sys.argv) % 2 != 0:
            print("Ошибка: для compare нужны пары <file> <method>")
            sys.exit(1)
        
        csv_files = sys.argv[2::2]
        method_names = sys.argv[3::2]
        compare_methods(csv_files, method_names)
    else:
        csv_file = sys.argv[1]
        method_name = sys.argv[2] if len(sys.argv) > 2 else "Unknown"
        plot_results(csv_file, method_name)