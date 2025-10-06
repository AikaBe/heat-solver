document.getElementById("run").addEventListener("click", async () => {
  const method = document.getElementById("method").value;
  const dx = document.getElementById("dx").value;
  const dt = document.getElementById("dt").value;
  const tmax = document.getElementById("tmax").value;

  const res = await fetch(`/simulate?method=${method}&dx=${dx}&dt=${dt}&tmax=${tmax}`);
  const data = await res.json();

  const u = data.u;
  const dxVal = data.dx;
  const dtVal = data.dt;
  animateHeat(u, dxVal, dtVal);
});

function animateHeat(u, dx, dt) {
  const canvas = document.getElementById("heatCanvas");
  const ctx = canvas.getContext("2d");
  const nx = u[0].length;
  const width = canvas.width;
  const height = canvas.height;

  let frame = 0;
  const nt = u.length;
  let animationId = null;

  function drawFrame() {
    ctx.clearRect(0, 0, width, height);

    // --- Рисуем линию ---
    ctx.beginPath();
    ctx.moveTo(0, height / 2);
    for (let i = 0; i < nx; i++) {
      const x = (i / nx) * width;
      const y = height / 2 - u[frame][i] * 200;
      ctx.lineTo(x, y);
    }
    ctx.strokeStyle = "orange";
    ctx.lineWidth = 2;
    ctx.stroke();

    // --- Получаем температуру в центре ---
    const centerIndex = Math.floor(nx / 2);
    const tempCenter = u[frame][centerIndex];
    const timeNow = (frame * dt).toFixed(4);

    // --- Подписи ---
    ctx.fillStyle = "rgba(0, 0, 0, 0.6)";
    ctx.fillRect(10, height - 45, 270, 35); // фон для текста
    ctx.fillStyle = "white";
    ctx.font = "16px Arial";
    ctx.fillText(`🕒 Time: ${timeNow} s`, 20, height - 25);
    ctx.fillText(`🌡️ Temp(center): ${tempCenter.toFixed(4)}`, 160, height - 25);

    // --- Проверяем условие остановки ---
    if (Math.abs(tempCenter) < 1e-4) {
      cancelAnimationFrame(animationId);
      console.log(`Simulation stopped: temperature at center ≈ 0 (t = ${timeNow}s)`);
      return;
    }

    // --- Следующий кадр ---
    frame++;
    if (frame < nt) {
      animationId = requestAnimationFrame(drawFrame);
    }
  }

  drawFrame();
}
