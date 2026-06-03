package main

type Piloto struct {
	Numero       int       `json:"Numero"`
	Nombre       string    `json:"Nombre"`
	Nacionalidad string    `json:"Nacionalidad"`
	Equipo       string    `json:"Equipo"`
	Vueltas      int       `json:"Vueltas"`
	TiempoTotal  string    `json:"TiempoTotal"`
	VueltaRapida float64   `json:"VueltaRapida"`
	Posicion     int       `json:"Posicion"`
	Puntos       int       `json:"Puntos"`
	Tiempos      []float64 `json:"Tiempos"`
}

type Carrera struct {
	Nombre  string   `json:"Nombre"`
	Pilotos []Piloto `json:"Pilotos"`
}

// Representa al piloto y su puntaje final
// La uso para armar los rankings al final de main asi puedo ordenarlos con sort.Slice
type ResultadoPiloto struct {
	Nombre    string
	Victorias int
}
