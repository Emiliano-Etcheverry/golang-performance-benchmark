package main

// El datos.go lo hice para separar los archivos de como se obtienen los datos, de que son las cosas (eso esta en modelos.go)

import (
	"encoding/json" // Libreria para trabajar con JSONs
	"os"            // Esta me sirve para interactuar con el sistema operativo, sus archivos y las variables de entorno
)

func obtenerCarreras() []Carrera {
	archivo, _ := os.ReadFile("carreras.json") // Aca simplemente pido que lea el archivo y devuelva el contenido en archivo como bytes
	var carreras []Carrera
	json.Unmarshal(archivo, &carreras) // El Unmarshal toma los bytes del JSON y los transforma en un slice de Carrera, el & es porque necesita la direccion de memoria para escribir ahi los datos

	return carreras
}
