import pandas as pd
import matplotlib.pyplot as plt
import matplotlib.ticker as ticker

# Leer el CSV
df = pd.read_csv("tiempos.csv", 
                  header=None,
                  names=["carreras", "secuencial", "paralelo", "workerpool", "fanout"])

# Ordenar por cantidad de carreras
df = df.sort_values("carreras")

fig, ax = plt.subplots(figsize=(12, 7))

ax.plot(df["carreras"], df["secuencial"], label="Secuencial",    color="red",    linewidth=0.5)
ax.plot(df["carreras"], df["paralelo"], label="Paralelo",      color="blue",   linewidth=0.5)
ax.plot(df["carreras"], df["workerpool"], label="Worker Pool",   color="green",  linewidth=0.5)
ax.plot(df["carreras"], df["fanout"], label="Fan-Out/Fan-In",color="purple", linewidth=0.5)

# Etiquetas y titulo
ax.set_xlabel("Cantidad de carreras procesadas", fontsize=13)
ax.set_ylabel("Tiempo de ejecucion (ms)", fontsize=13)
ax.set_title("Comparativa de performance: Secuencial vs Paralelo vs Worker Pool vs Fan-Out/In", fontsize=14)

# Grilla y leyenda
ax.grid(True, linestyle="--", alpha=0.6)
ax.legend(fontsize=12)

# Formato del eje x
ax.xaxis.set_major_formatter(ticker.FuncFormatter(lambda x, _: f"{int(x):,}"))

plt.tight_layout()
plt.savefig("comparativa_performance.png", dpi=150)
plt.show()

print("Grafico guardado en comparativa_performance.png")