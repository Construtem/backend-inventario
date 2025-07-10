package main

import (
	"fmt"
	"log"
	"os"

	//	modelos "backend-inventario/api/Models"
	"backend-inventario/config"
	"backend-inventario/api/Routes"
	"backend-inventario/api/db"
	"backend-inventario/services"
	"backend-inventario/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	   // En producción (Kubernetes), solo intenta cargar .env en local, pero nunca detengas la app si no existe
    config.LoadEnv() // Esto ya maneja el log si no existe .env

	database, err := db.ConectarDB()
	if err != nil {
		log.Fatalf("Error al conectar a la base de datos: %v", err)
	}
	fmt.Println("Conexión a la base de datos exitosa")

	services.InitFirebase()

	//	modelos.MigrarTablas(database)
	//	fmt.Println("Migración de tablas exitosa")

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://inventario.tssw.cl"}, // Permite tu frontend de Next.js
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},	// Agregado el OPTIONS para autenticacion
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	router.POST("/auth/verify", handlers.VerifyToken)	//Ruta para autenticacion firebase
	Routes.RegisterRoutes(router, database)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Puerto por defecto
	}

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
	fmt.Printf("Servidor corriendo en http://localhost:%s\n", port)
}
