package Routes

import (
	"backend-inventario/api/Controllers"
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

	// Rutas para Despachos
	api.GET("/despachos", Handlers.GetDespachosHandler(db))
	api.GET("/despachos/:id", Handlers.GetDespachoByIDHandler(db))
	api.POST("/despachos", Handlers.CreateDespachoHandler(db))
	api.PUT("/despachos/:id", Handlers.UpdateDespachoHandler(db))
	api.DELETE("/despachos/:id", Handlers.DeleteDespachoHandler(db))
	api.GET("/despachos/:id/pdf", Controllers.GenerarDespachoPDF(db))

	// Rutas para Tipos de Camion
	api.GET("/tipos-camion", Handlers.GetTipoCamionHandler(db))
	api.GET("/tipos-camion/:id", Handlers.GetTipoCamionByIDHandler(db))
	api.POST("/tipos-camion", Handlers.CreateTipoCamionHandler(db))
	api.PUT("/tipos-camion/:id", Handlers.UpdateTipoCamionHandler(db))
	api.DELETE("/tipos-camion/:id", Handlers.DeleteTipoCamionHandler(db))

	// Rutas para Camiones
	api.GET("/camiones", Handlers.GetCamionesHandler(db))
	api.GET("/camiones/:id", Handlers.GetCamionByIDHandler(db))
	api.POST("/camiones", Handlers.CreateCamionHandler(db))
	api.PUT("/camiones/:id", Handlers.UpdateCamionHandler(db))
	api.DELETE("/camiones/:id", Handlers.DeleteCamionHandler(db))
}
