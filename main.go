package main

import (
	"fmt"
	"log"
	"os"

	modelos "backend-inventario/api/Models"
	"backend-inventario/api/Routes"
	"backend-inventario/api/db"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error al cargar .env")
	}

	db.Conectar()
	fmt.Println("Conexión a la base de datos exitosa")

	modelos.MigrarTablas(db.DB)
	fmt.Println("Migración de tablas exitosa")

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))

	Routes.RegisterRoutes(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Puerto por defecto
	}

	err = router.Run(":" + port)
	if err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
	fmt.Printf("Servidor corriendo en http://localhost:%s\n", port)
}
