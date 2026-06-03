package main

import (
	"fmt"
	"os"
)

func main() {
	equipos := []string{
		"McLaren Formula 1 Team",
		"Mercedes-AMG PETRONAS F1 Team",
		"Scuderia Ferrari HP",
		"Oracle Red Bull Racing",
		"Aston Martin Aramco F1 Team",
		"Atlassian Williams Racing",
		"MoneyGram Haas F1 Team",
		"Visa Cash App Racing Bulls F1 Team",
		"Kick Sauber F1 Team",
		"BWT Alpine F1 Team",
	}

	cantidadWorkers := []int{1, 2, 3, 4, 5, 6, 7, 8}

	archivo, _ := os.OpenFile("./Python/workers.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	defer archivo.Close()

	// loop por cantidad de carreras
	for anioFin := 2026; anioFin <= 3024; anioFin += 1 {

		carreras := obtenerCarreras()
		carreras = append(carreras, generarTemporadas(carreras, 2025, anioFin)...)

		for i := range carreras {
			for j := range carreras[i].Pilotos {
				p := &carreras[i].Pilotos[j]
				p.Tiempos = generarTiemposVuelta(p.Vueltas, p.TiempoTotal, p.VueltaRapida)
			}
		}

		// loop por cantidad de workers
		for _, numWorkers := range cantidadWorkers {
			tiempo := procesarWorkerPool(carreras, equipos, numWorkers)

			fmt.Printf("Carreras: %d | Workers: %d | Tiempo: %dms\n",
				len(carreras), numWorkers, tiempo.Milliseconds())

			fmt.Fprintf(archivo, "%d,%d,%d\n",
				len(carreras),
				numWorkers,
				tiempo.Milliseconds(),
			)
		}
	}
}
