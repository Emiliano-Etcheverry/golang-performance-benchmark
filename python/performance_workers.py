import pandas as pd
import matplotlib.pyplot as plt
import matplotlib.ticker as ticker

df = pd.read_csv("workers.csv",
                  header=None,
                  names=["carreras", "workers", "tiempo"])

df = df.sort_values(["carreras", "workers"])


fig, (ax1, ax2) = plt.subplots(2, 1, figsize=(12, 14))

colores = {
    1:    "red",
    2:    "orange",
    3:    "green",
    4:   "blue",
    5:   "purple",
    6:   "brown",
    7:  "pink",
    8:  "gray",
    9:  "black",
    10: "violet"
}

ventana = 10
workers_excluidos = [2, 4]

for numWorkers in df["workers"].unique():
    if numWorkers in workers_excluidos:
        continue
    df_worker = df[df["workers"] == numWorkers].copy()
    df_worker["tiempo_smooth"] = df_worker["tiempo"].rolling(window=ventana).mean()
    ax1.plot(df_worker["carreras"], df_worker["tiempo_smooth"],
             label=f"{numWorkers} workers",
             color=colores[numWorkers],
             linewidth=0.5)

ax1.set_xlabel("Cantidad de carreras procesadas", fontsize=13)
ax1.set_ylabel("Tiempo de ejecucion (ms)", fontsize=13)
ax1.set_title("Worker Pool - Comparativa por cantidad de workers", fontsize=14)
ax1.grid(True, linestyle="--", alpha=0.6)
ax1.legend(fontsize=12)
ax1.xaxis.set_major_formatter(ticker.FuncFormatter(lambda x, _: f"{int(x):,}"))

# Grafico 2 - tiempo promedio por cantidad de workers
df_promedio = df.groupby("workers")["tiempo"].mean().reset_index()
ax2.bar(df_promedio["workers"].astype(str), df_promedio["tiempo"],
        color=list(colores.values()))
ax2.set_xlabel("Cantidad de workers", fontsize=13)
ax2.set_ylabel("Tiempo promedio (ms)", fontsize=13)
ax2.set_title("Tiempo promedio por cantidad de workers", fontsize=14)
ax2.grid(True, linestyle="--", alpha=0.6, axis="y")

plt.tight_layout()
plt.savefig("comparativa_workers.png", dpi=150)
plt.show()

print("Grafico guardado en comparativa_workers.png")