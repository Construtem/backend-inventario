package main

import (
	"backend-inventario/api/db"
	"backend-inventario/config"
	"fmt"
	"log"
)

func main() {
	// Cargar variables de entorno
	config.LoadEnv()

	// Conectar a la base de datos
	database, err := db.ConectarDB()
	if err != nil {
		log.Fatalf("Error al conectar a la base de datos: %v", err)
	}

	fmt.Println("Iniciando migración para agregar campos de distancia...")

	// Agregar las nuevas columnas a la tabla despacho
	sqlStatements := []string{
		`ALTER TABLE despacho ADD COLUMN IF NOT EXISTS distancia_calculada VARCHAR(50);`,
		`ALTER TABLE despacho ADD COLUMN IF NOT EXISTS tiempo_estimado VARCHAR(50);`,
	}

	for _, stmt := range sqlStatements {
		if err := database.Exec(stmt).Error; err != nil {
			log.Printf("Error ejecutando: %s - %v", stmt, err)
		} else {
			fmt.Printf("Ejecutado exitosamente: %s\n", stmt)
		}
	}

	fmt.Println("Migración completada exitosamente!")
}
