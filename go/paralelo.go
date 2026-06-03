package main

import (
	"fmt"
	"sync" // Esta esta para sincronizar los goroutines, coordinando el acceso a los recursos de cada uno
	"time" // Esta es para trabajar con el tiempo, en este caso sirve de cronometro para ver cuanto tarda cada procesoer
)

func procesarParalelo(carreras []Carrera, equipos []string) time.Duration {
	var wg sync.WaitGroup
	var mu sync.Mutex

	puntosTotales := make(map[int]int)
	puntajeInterno := make(map[int]int)
	nombrePiloto := make(map[int]string)
	inicioP := time.Now() // Con esto capturo el momento para arrancar a medir el tiempo de las go routines

	for _, carrera := range carreras {
		wg.Add(1) // Esto le avisa al WaitGroup que hay una goroutine mas esperando

		// Con esto lanzo la funcion en una goroutine separada, por lo que no espera a que termine, si no que pasa a la proxima carrera y a esa tambien la genera en una rutina separada
		go func(c Carrera) { // Se le pasa el parametro de Carrera para que las goroutines no compartan la misma variable, asi que se van a diferenciar por cada carrera

			defer wg.Done() // Esto es pone para que cuando la goroutine termine le avise al WaitGroup que termino

			// Resultados independientes de cada goroutine - sin usar mutex
			puntosLocales := make(map[int]int)
			nombresLocales := make(map[int]string)
			ganadores := make(map[string]int)

			for _, piloto := range c.Pilotos {
				tiempoTotal := convertirTiempo(piloto.TiempoTotal)
				_ = promedioVueltas(tiempoTotal, piloto.Vueltas)
				_ = consistencia(piloto.Tiempos, tiempoTotal, piloto.Vueltas)
				_ = peorTiempoxVuelta(piloto.Tiempos)
				puntosLocales[piloto.Numero] += piloto.Puntos
				nombresLocales[piloto.Numero] = piloto.Nombre
			}

			for _, equipo := range equipos {
				ganador := mejorxEscuderia(c, equipo)
				if ganador.Nombre != "" {
					ganadores[equipo] = ganador.Numero
					nombresLocales[ganador.Numero] = ganador.Nombre
				}
			}

			mu.Lock() // Hago el lock al momento de tener que escribir en los maps compartidos por todas las goroutines en lugar de lockear por cada operacion

			for numero, puntos := range puntosLocales {
				puntosTotales[numero] += puntos
				nombrePiloto[numero] = nombresLocales[numero]
			}

			for _, numero := range ganadores {
				puntajeInterno[numero]++
			}

			mu.Unlock()
		}(carrera)
	}

	wg.Wait() // Esto es para esperar a que todas terminen, si no lo pongo el main podria terminar tranquilamente antes de que todas las rutinas terminen de ejecutarse

	tiempoParalelo := time.Since(inicioP)
	fmt.Println("Tiempo paralelo (goroutine por carrera):", tiempoParalelo)

	return tiempoParalelo
}
