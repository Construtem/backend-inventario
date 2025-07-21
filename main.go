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
	// Cargar variables de entorno
	config.LoadEnv()

	// Conectar a la base de datos
	log.Printf("Conectando a la base de datos %s...", os.Getenv("DB_NAME"))
	database, err := db.ConectarDB()
	if err != nil {
		log.Fatalf("Error al conectar a la base de datos: %v", err)
	}
	fmt.Println("Conexión a la base de datos exitosa")

	// Inicializar Firebase
	services.InitFirebase()

	// Crear router
	router := gin.Default()

	// Configurar CORS para PRODUCCION
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			os.Getenv("FRONT_INVENTARIO_URL"),
			os.Getenv("FRONT_VENTAS_URL"),
			os.Getenv("FRONT_FACTURACION_URL"),
			//os.Getenv("BACK_FACTURACION_URL")
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Ruta para verificación de autenticación Firebase
	router.POST("/auth/verify", handlers.VerifyToken)

	// Registrar rutas del backend
	Routes.RegisterRoutes(router, database)

	// Puerto por defecto
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Iniciar servidor
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
	fmt.Printf("Servidor corriendo en http://localhost:%s\n", port)
}