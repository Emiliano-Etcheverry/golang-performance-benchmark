package main

import (
	"fmt"
	"sync" // Esta esta para sincronizar los goroutines, coordinando el acceso a los recursos de cada uno
	"time" // Esta es para trabajar con el tiempo, en este caso sirve de cronometro para ver cuanto tarda cada procesoer
)

func procesarFanOut(carreras []Carrera, equipos []string) time.Duration {
	puntosTotalesFan := make(map[int]int)
	puntajeInternoFan := make(map[int]int)
	nombrePilotoFan := make(map[int]string)

	// Resultado que cada goroutine va a enviar al fan-in
	type ResultadoFan struct {
		puntosLocales  map[int]int
		nombresLocales map[int]string
		ganadores      map[string]int
	}

	// Esta funcion distribuye las carreras a workers y me devuelve un channel con los resultados
	fanOut := func(carreras []Carrera, equipos []string) <-chan ResultadoFan { // el <- es para unicamente leer
		resultadosChan := make(chan ResultadoFan, len(carreras))
		var wgFan sync.WaitGroup

		for _, carrera := range carreras {
			wgFan.Add(1)
			go func(c Carrera) {
				defer wgFan.Done()

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

				// cada goroutine manda su resultado al channel
				resultadosChan <- ResultadoFan{puntosLocales, nombresLocales, ganadores}
			}(carrera)
		}

		// cuando todas terminan cerramos el channel
		go func() {
			wgFan.Wait()
			close(resultadosChan)
		}()

		return resultadosChan
	}

	// FAN-IN - funcion que junta todos los resultados del channel en los maps finales
	fanIn := func(resultadosChan <-chan ResultadoFan) {
		for r := range resultadosChan {
			for numero, puntos := range r.puntosLocales {
				puntosTotalesFan[numero] += puntos
				nombrePilotoFan[numero] = r.nombresLocales[numero]
			}
			for _, numero := range r.ganadores {
				puntajeInternoFan[numero]++
			}
		}
	}

	inicioFan := time.Now()
	resultadosChan := fanOut(carreras, equipos) // distribuye trabajo
	fanIn(resultadosChan)                       // junta resultados

	tiempoFan := time.Since(inicioFan)
	fmt.Println("Tiempo Fan-Out/Fan-In:", tiempoFan)

	return tiempoFan
}
