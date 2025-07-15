package main

import (
	"fmt"
	"log"
	"os"

	//	modelos "backend-inventario/api/Models"
	"backend-inventario/api/Routes"
	"backend-inventario/api/db"
	"backend-inventario/config"
	"backend-inventario/handlers"
	"backend-inventario/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload" // Carga automáticamente el archivo .env
)

func main() {
	// En producción (Kubernetes), solo intenta cargar .env en local, pero nunca detengas la app si no existe
	config.LoadEnv() // Esto ya maneja el log si no existe .env

	// Conectar a la base de datos
	log.Printf("Conectando a la base de datos %s...", os.Getenv("DB_NAME"))
	database, err := db.ConectarDB()
	if err != nil {
		log.Fatalf("Error al conectar a la base de datos: %v", err)
	}
	fmt.Println("Conexión a la base de datos exitosa")

	services.InitFirebase()

	//	modelos.MigrarTablas(database)
	//	fmt.Println("Migración de tablas exitosa")

	router := gin.Default()

	// Configurando CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{ // Lista de URLs permitidas para CORS
			os.Getenv("FRONT_VENTAS_URL"),      // URL del frontend de ventas
			os.Getenv("FRONT_INVENTARIO_URL"),  // URL del frontend de inventario
			os.Getenv("FRONT_FACTURACION_URL"), // URL del frontend de facturación
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	router.POST("/auth/verify", handlers.VerifyToken) //Ruta para autenticacion firebase
	Routes.RegisterRoutes(router, database)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Puerto alternativo
	}

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
	fmt.Printf("Servidor corriendo en http://localhost:%s\n", port)
}
