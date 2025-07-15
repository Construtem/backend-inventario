package main

import (
	"backend-inventario/api/Controllers"
	"fmt"
	"log"
)

// Ejemplo de uso de la nueva funcionalidad de cálculo de distancias
func main() {
	// Ejemplo: calcular distancia entre dos direcciones
	apiKey := "TU_API_KEY_AQUI" // Reemplaza con tu API key real
	origen := "Santiago Centro, Santiago, Chile"
	destino := "Las Condes, Santiago, Chile"

	fmt.Println("Calculando distancia entre:")
	fmt.Printf("Origen: %s\n", origen)
	fmt.Printf("Destino: %s\n", destino)

	distancia, duracion, err := Controllers.CalcularDistancia(apiKey, origen, destino)
	if err != nil {
		log.Printf("Error al calcular distancia: %v", err)
		return
	}

	fmt.Printf("\nResultado:\n")
	fmt.Printf("Distancia: %.2f km\n", distancia)
	fmt.Printf("Duración: %d minutos\n", duracion)
}
