package main

import (
	"sync" // Esta esta para sincronizar los goroutines, coordinando el acceso a los recursos de cada uno
	"time" // Esta es para trabajar con el tiempo, en este caso sirve de cronometro para ver cuanto tarda cada procesoer
)

func procesarWorkerPool(carreras []Carrera, equipos []string, Workers int) time.Duration {
	var muWP sync.Mutex
	var wgWP sync.WaitGroup

	// En lugar de lanzar una goroutine por carrera, usamos un numero fijo de workers
	// que van agarrando trabajo de una cola - mas controlado y eficiente
	puntosTotalesWP := make(map[int]int)
	puntajeInternoWP := make(map[int]int)
	nombrePilotoWP := make(map[int]string)

	numWorkers := Workers                         // cantidad de workers igual a los nucleos del procesador
	jobsChan := make(chan Carrera, len(carreras)) // channel que actua como cola de trabajo

	inicioWP := time.Now()

	// lanzo los 8 workers - cada uno va agarrando carreras del channel
	for w := 0; w < numWorkers; w++ {
		wgWP.Add(1)
		go func() {
			defer wgWP.Done()
			for c := range jobsChan { // el worker agarra carreras hasta que el channel se cierre
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

				muWP.Lock() // lock solo para escribir en los maps compartidos
				for numero, puntos := range puntosLocales {
					puntosTotalesWP[numero] += puntos
					nombrePilotoWP[numero] = nombresLocales[numero]
				}
				for _, numero := range ganadores {
					puntajeInternoWP[numero]++
				}
				muWP.Unlock()
			}
		}()
	}

	// cargo todas las carreras en el channel
	for _, carrera := range carreras {
		jobsChan <- carrera
	}
	close(jobsChan) // cierro el channel para que los workers sepan que no hay mas trabajo

	wgWP.Wait() // espero que todos los workers terminen
	tiempoWorker := time.Since(inicioWP)

	return tiempoWorker
}
