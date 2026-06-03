package main

import (
	"fmt"
	"time" // Esta es para trabajar con el tiempo, en este caso sirve de cronometro para ver cuanto tarda cada procesoer
)

func procesarSecuencial(carreras []Carrera, equipos []string) time.Duration {
	// Procesos ejecutados de manera secuencial para comparar con el paralelo
	puntosTotalesSeq := make(map[int]int)
	puntajeInternoSeq := make(map[int]int)
	nombrePilotoSeq := make(map[int]string)
	inicioS := time.Now()

	for _, carrera := range carreras {
		puntosLocales := make(map[int]int)
		nombresLocales := make(map[int]string)
		ganadores := make(map[string]int)

		for _, piloto := range carrera.Pilotos {
			tiempoTotal := convertirTiempo(piloto.TiempoTotal)
			_ = promedioVueltas(tiempoTotal, piloto.Vueltas)
			_ = consistencia(piloto.Tiempos, tiempoTotal, piloto.Vueltas)
			_ = peorTiempoxVuelta(piloto.Tiempos)
			puntosLocales[piloto.Numero] += piloto.Puntos
			nombresLocales[piloto.Numero] = piloto.Nombre
		}

		for _, equipo := range equipos {
			ganador := mejorxEscuderia(carrera, equipo)
			if ganador.Nombre != "" {
				ganadores[equipo] = ganador.Numero
				nombresLocales[ganador.Numero] = ganador.Nombre
			}
		}

		for numero, puntos := range puntosLocales {
			puntosTotalesSeq[numero] += puntos
			nombrePilotoSeq[numero] = nombresLocales[numero]
		}

		for _, numero := range ganadores {
			puntajeInternoSeq[numero]++
		}
	}

	tiempoSecuencial := time.Since(inicioS)
	fmt.Println("Tiempo secuencial:", tiempoSecuencial)

	return tiempoSecuencial
}
