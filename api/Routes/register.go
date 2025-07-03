package Routes

import (
	"backend-inventario/api/Handlers"

	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(router *gin.Engine, db *gorm.DB) {
	// Configuración de CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Permite tu frontend de Next.js
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour, // 24 horas
	}))

	// Grupo de rutas para la API
	api := router.Group("/api")

	// Rutas para Productos
	api.GET("/productos", Handlers.GetProductosHandler(db))
	api.GET("/productos/:sku", Handlers.GetProductoBySKUHandler(db))
	api.POST("/productos", Handlers.CreateProductoHandler(db))
	api.PUT("/productos/:sku", Handlers.UpdateProductoHandler(db))
	api.DELETE("/productos/:sku", Handlers.DeleteProductoHandler(db))

	// Rutas para Sucursales
	api.GET("/sucursales", Handlers.GetSucursalesHandler(db))
	api.GET("/sucursales/:id", Handlers.GetSucursalByIDHandler(db))
	api.POST("/sucursales", Handlers.CreateSucursalHandler(db))
	api.PUT("/sucursales/:id", Handlers.UpdateSucursalHandler(db))
	api.DELETE("/sucursales/:id", Handlers.DeleteSucursalHandler(db))

	// Rutas para Stock por Sucursal
	api.GET("/stock-sucursal", Handlers.GetStockSucursalHandler(db))
	api.GET("/stock-sucursal/:sucursal_id/:producto_id", Handlers.GetStockSucursalByIDHandler(db))
	api.POST("/stock-sucursal", Handlers.CreateStockSucursalHandler(db))
	api.PUT("/stock-sucursal/:sucursal_id/:producto_id", Handlers.UpdateStockSucursalHandler(db))
	api.DELETE("/stock-sucursal/:sucursal_id/:producto_id", Handlers.DeleteStockSucursalHandler(db))

	// Rutas para Tipo de Sucursal
	api.GET("/tipos-sucursal", Handlers.GetTipoSucursalHandler(db))
	api.GET("/tipos-sucursal/:id", Handlers.GetTipoSucursalByIDHandler(db))
	api.POST("/tipos-sucursal", Handlers.CreateTipoSucursalHandler(db))
	api.PUT("/tipos-sucursal/:id", Handlers.UpdateTipoSucursalHandler(db))
	api.DELETE("/tipos-sucursal/:id", Handlers.DeleteTipoSucursalHandler(db))
}
