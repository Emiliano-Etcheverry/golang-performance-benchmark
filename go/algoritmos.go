package main

import (
	"fmt"       // Formating package - sirve para hacer los print
	"math"      // Me brinda las funciones para hacer operaciones como sqrt y demas
	"math/rand" // Es uno de los paquetes de rand tiene el de math y el de crypto, el de math es para uso gral y es mas rapido
	"sort"      // Me sirve para hacer el ordenamiento de estructuras de datos como slice
	"strconv"   // Me sirve para hcaer la conversion de strings y otros tipos de datos
	"strings"   // Me sirve para separar strings
)

// Funcion que reconvierte un string con los tiempos a float64
// Por tanto partes[0] = hora partes[1] = minutos partes[2] = segundos
func convertirTiempo(tiempo string) float64 {
	partes := strings.Split(tiempo, ":")
	if len(partes) < 3 {
		return 0.0 // Si el tiempo es inválido
	}
	horas, _ := strconv.ParseFloat(partes[0], 64)
	minutos, _ := strconv.ParseFloat(partes[1], 64)
	segundos, _ := strconv.ParseFloat(partes[2], 64)
	return (horas * 3600) + (minutos * 60) + segundos
}

// Funcion que hace lo opuesto a la anterior, toma un float64 y devuelve un string de ese tiempo reconvertido con una diferencia de +- 5s
func cambiarTiempo(tiempo float64) string {
	nuevoTiempo := tiempo + (rand.Float64() * 10) - 5
	horas := int(nuevoTiempo) / 3600
	minutos := (int(nuevoTiempo) % 3600) / 60
	segundos := nuevoTiempo - float64(horas*3600+minutos*60)
	return fmt.Sprintf("%d:%02d:%06.3f", horas, minutos, segundos)
}

// Funcion para generar los tiempos de vuelta de manera aleatoria tomando como base la VueltaRapida
func generarTiemposVuelta(vueltas int, tiempoTotal string, vueltaRapida float64) []float64 {
	totalSegundos := convertirTiempo(tiempoTotal)
	tiempos := make([]float64, vueltas)
	sumaTiempos := 0.0
	for i := range vueltas {
		tiempos[i] = vueltaRapida + rand.Float64()*10 // El por 10  es para agregarle a las vueltas generadas hasta 10s
		sumaTiempos = sumaTiempos + tiempos[i]
	}

	factorTiempo := totalSegundos / sumaTiempos // El sentido de esto es que todos los tiempos sumados den exactamente el tiempoTotal porque al generarlos random puede que la suma total de mas que el tiempoTotal entonces los reescalo para que coincida

	for i := range vueltas {
		tiempos[i] = tiempos[i] * factorTiempo
	}

	return tiempos
}

// Me devuelve el tiempo promedio entre las vueltas de un piloto en una carrera
func promedioVueltas(tiempoTotal float64, vueltas int) float64 {
	tiempoPromedio := tiempoTotal / float64(vueltas)
	return tiempoPromedio
}

// Funcion para encontrar el mejor tiempo de entre todos los pilotos por carrera
func mejorTiempoCarrera(carrera Carrera) float64 {
	mejor := convertirTiempo(carrera.Pilotos[0].TiempoTotal)
	for _, piloto := range carrera.Pilotos {
		tiempo := convertirTiempo(piloto.TiempoTotal)
		if tiempo < mejor {
			mejor = tiempo
		}
	}
	return mejor
}

// Lo mismo que la anterior pero al reves
func peorTiempoxVuelta(tiempos []float64) float64 {
	peorTiempo := tiempos[0]
	for i := range tiempos {
		if peorTiempo < tiempos[i] {
			peorTiempo = tiempos[i]
		}
	}
	return peorTiempo
}

// Funcion que calcula la "Consistencia" y para esta voy a utilizar el varemo de la desviacion estandar
// Σ(tiempo[i] - promedio)²
func consistencia(tiempos []float64, tiempoTotal float64, vueltas int) float64 {
	promedio := tiempoTotal / float64(vueltas)
	suma := 0.0
	for i := range tiempos {
		suma += math.Pow(tiempos[i]-promedio, 2) // el math.Pow hace la potencia que seleccione de un valor
	}
	desviacion := math.Sqrt(suma / float64(vueltas))
	return desviacion
}

// Funcion que devuelve al piloto que menor indice de consistencia tuvo en comparacion con su compañero de equipo
func mejorxEscuderia(carrera Carrera, escuderia string) Piloto {
	var piloto1 Piloto
	var piloto2 Piloto

	cont := 0
	for _, piloto := range carrera.Pilotos {
		if piloto.Equipo == escuderia {
			if cont == 0 {
				piloto1 = piloto
			} else if cont == 1 {
				piloto2 = piloto
			}
			cont++
		}
	}

	tiempoTotal1 := convertirTiempo(piloto1.TiempoTotal)
	tiempoTotal2 := convertirTiempo(piloto2.TiempoTotal)
	if cont == 0 {
		return Piloto{}
	} else if cont < 2 {
		return piloto1
	}
	if consistencia(piloto1.Tiempos, tiempoTotal1, piloto1.Vueltas) < consistencia(piloto2.Tiempos, tiempoTotal2, piloto2.Vueltas) {
		return piloto1
	} else {
		return piloto2
	}
}

// Funcion que genera nuevas temporadas con sus respectivas carreras y pilotos, y ademas genera los tiempos de todas estas basandose en el JSON principal con los datos reales del campeonato 2025
func generarTemporadas(carrerasBase []Carrera, anioInicio int, anioFin int) []Carrera {
	var nuevasCarreras []Carrera

	for anio := anioInicio; anio <= anioFin; anio++ {
		for _, carreraBase := range carrerasBase {
			nuevaCarrera := Carrera{
				Nombre: strings.ReplaceAll(carreraBase.Nombre, "2025", fmt.Sprintf("%d", anio)), // Por cada carrera crea una nueva con el mismo nombre pero cambiando el año
			}

			mejorTiempoBase := mejorTiempoCarrera(carreraBase) // Intente que todos los tiempos totales que se van a generar partan del mejor tiempo de la carrera de la temporada 2025 asi todos tenian tiempos distintos en cada carrera y cualquiera podia ganar

			for _, piloto := range carreraBase.Pilotos {

				totalSegundos := convertirTiempo(piloto.TiempoTotal)
				diferencia := totalSegundos - mejorTiempoBase
				variacion := (rand.Float64()*2 - 1) * diferencia
				nuevoTotal := totalSegundos + variacion
				nuevoTiempoStr := cambiarTiempo(nuevoTotal) // Solo para guardar el string
				factor := nuevoTotal / convertirTiempo(piloto.TiempoTotal)

				nuevoPiloto := Piloto{
					Numero:       piloto.Numero,
					Nombre:       piloto.Nombre,
					Nacionalidad: piloto.Nacionalidad,
					Equipo:       piloto.Equipo,
					Vueltas:      piloto.Vueltas,
					TiempoTotal:  nuevoTiempoStr,
					VueltaRapida: piloto.VueltaRapida * factor,
					Posicion:     0, // la asigno luego
					Puntos:       0, // la asigno luego tambien
				}
				nuevoPiloto.Tiempos = generarTiemposVuelta(nuevoPiloto.Vueltas, nuevoPiloto.TiempoTotal, nuevoPiloto.VueltaRapida)
				nuevaCarrera.Pilotos = append(nuevaCarrera.Pilotos, nuevoPiloto)
			}

			// Ordeno los pilotos en base al tiempo total y les asigo sus nuevas posiciones y puntos
			tiemposEnSegundos := make([]float64, len(nuevaCarrera.Pilotos))
			for i, p := range nuevaCarrera.Pilotos {
				tiemposEnSegundos[i] = convertirTiempo(p.TiempoTotal)
			}
			sort.Slice(nuevaCarrera.Pilotos, func(i, j int) bool {
				return tiemposEnSegundos[i] < tiemposEnSegundos[j]
			})
			puntosPorPosicion := []int{0, 25, 18, 15, 12, 10, 8, 6, 4, 2, 1}
			for i := range nuevaCarrera.Pilotos {
				nuevaCarrera.Pilotos[i].Posicion = i + 1
				if i+1 <= 10 {
					nuevaCarrera.Pilotos[i].Puntos = puntosPorPosicion[i+1]
				} else {
					nuevaCarrera.Pilotos[i].Puntos = 0 // Si no estan en los primeros 10 no tienen puntos
				}
			}

			nuevasCarreras = append(nuevasCarreras, nuevaCarrera) // Agrego la nueva carrera entera al final del slice
		}
	}
	return nuevasCarreras
}
