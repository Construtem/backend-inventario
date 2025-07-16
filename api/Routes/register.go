package Routes

import (
	"backend-inventario/api/Controllers"
	"backend-inventario/api/Handlers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(router *gin.Engine, db *gorm.DB) {
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

	// Rutas específicas para Bodegas
	api.GET("/bodegas", Handlers.GetBodegasHandler(db))

	// Ruta temporal de debug para ver tipos de sucursal
	api.GET("/debug/tipos-sucursal", Handlers.GetTiposSucursalDebugHandler(db))

	// Rutas para Stock por Sucursal
	api.GET("/stock-sucursal", Handlers.GetStockSucursalHandler(db))
	api.GET("/stock-sucursal/:sucursal_id/:sku", Handlers.GetStockSucursalByIDHandler(db))
	api.POST("/stock-sucursal", Handlers.CreateStockSucursalHandler(db))
	api.PUT("/stock-sucursal/:sucursal_id/:sku", Handlers.UpdateStockSucursalHandler(db))
	api.DELETE("/stock-sucursal/:sucursal_id/:sku", Handlers.DeleteStockSucursalHandler(db))

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
	api.POST("/despachos/calcular", Handlers.CalcularDespachoHandler(db))
	api.GET("/despachos/cotizacion/:id", Handlers.GetDespachosPorCotizacionHandler(db))
	api.POST("/despachos/aprobar", Handlers.AprobarDespachoHandler(db))
	api.GET("/despachos/:id/pdf", Controllers.GenerarDespachoPDF(db))

	// Rutas específicas para el sistema de rutas y distancias
	api.GET("/despachos-distancia", Handlers.GetDespachosDistanciaHandler(db))
	api.GET("/despachos-distancia/:id", Handlers.GetDespachoDistanciaByIDHandler(db))
	api.POST("/despachos/:id/calcular-distancia", Handlers.CalcularDistanciaDespachoHandler(db))
	api.POST("/despachos/:id/calcular-distancia-automatico", Handlers.CalcularDistanciaAutomaticoHandler(db))

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

	// Rutas para Clientes
	api.GET("/clientes", Handlers.GetClientesHandler(db))
	api.GET("/clientes/:id", Handlers.GetClienteByIDHandler(db))
	api.POST("/clientes", Handlers.CreateClienteHandler(db))
	api.PUT("/clientes/:id", Handlers.UpdateClienteHandler(db))
	api.DELETE("/clientes/:id", Handlers.DeleteClienteHandler(db))

	// Rutas para Tipo de Clientes
	api.GET("/tipos-clientes", Handlers.GetTipoClienteHandler(db))
	api.GET("/tipos-clientes/:id", Handlers.GetTipoClienteByIDHandler(db))
	api.POST("/tipos-clientes", Handlers.CreateTipoClienteHandler(db))
	api.PUT("/tipos-clientes/:id", Handlers.UpdateTipoClienteHandler(db))
	api.DELETE("/tipos-clientes/:id", Handlers.DeleteTipoClienteHandler(db))

	// Rutas para Direcciones de Clientes
	api.GET("/direcciones-clientes", Handlers.GetDirClientesHandler(db))
	api.GET("/direcciones-clientes/:id", Handlers.GetDirClienteByIDHandler(db))
	api.POST("/direcciones-clientes", Handlers.CreateDirClienteHandler(db))
	api.PUT("/direcciones-clientes/:id", Handlers.UpdateDirClienteHandler(db))
	api.DELETE("/direcciones-clientes/:id", Handlers.DeleteDirClienteHandler(db))

	// Rutas para Roles
	api.GET("/rol", Handlers.GetRolesHandler(db))
	api.GET("/rol/:id", Handlers.GetRolByIDHandler(db))
	api.POST("/rol", Handlers.CreateRolHandler(db))
	api.PUT("/rol/:id", Handlers.UpdateRolHandler(db))
	api.DELETE("/rol/:id", Handlers.DeleteRolHandler(db))

	// Rutas para Usuarios
	api.GET("/usuarios", Handlers.GetUsuariosHandler(db))
	api.GET("/usuarios/:email", Handlers.GetUsuarioByEmailHandler(db))
	api.POST("/usuarios", Handlers.CreateUsuarioHandler(db))
	api.PUT("/usuarios/:email", Handlers.UpdateUsuarioHandler(db))
	api.DELETE("/usuarios/:email", Handlers.DeleteUsuarioHandler(db))

	// Rutas para Proveedores
	api.GET("/proveedores", Handlers.GetProveedoresHandler(db))
	api.GET("/proveedores/:id", Handlers.GetProveedorByIDHandler(db))
	api.POST("/proveedores", Handlers.CreateProveedorHandler(db))
	api.PUT("/proveedores/:id", Handlers.UpdateProveedorHandler(db))
	api.DELETE("/proveedores/:id", Handlers.DeleteProveedorHandler(db))

	// Rutas para Stock de Proveedores
	api.GET("/stock-proveedor", Handlers.GetStockProveedorHandler(db))
	api.GET("/stock-proveedor/:proveedor_id/:producto_id", Handlers.GetStockProveedorByIDHandler(db))
	api.POST("/stock-proveedor", Handlers.CreateStockProveedorHandler(db))
	api.PUT("/stock-proveedor/:proveedor_id/:producto_id", Handlers.UpdateStockProveedorHandler(db))
	api.DELETE("/stock-proveedor/:proveedor_id/:producto_id", Handlers.DeleteStockProveedorHandler(db))

	// Rutas para Productos de Despacho
	api.GET("/productos_despacho", Handlers.GetProductosDespachoHandler(db))
	api.GET("/productos_despacho/detallado", Handlers.GetProductosDespachoDetalladoHandler(db))
	api.GET("/productos_despacho/despacho/:despacho_id", Handlers.GetProductosDespachoByDespachoIDHandler(db))
	api.GET("/productos_despacho/despacho/:despacho_id/detallado", Handlers.GetProductosDespachoDetalladoByDespachoIDHandler(db))
	api.GET("/productos_despacho/:despacho_id/:producto_id", Handlers.GetProductoDespachoByIDHandler(db))
	api.POST("/productos_despacho", Handlers.CreateProductoDespachoHandler(db))
	api.PUT("/productos_despacho/:despacho_id/:producto_id", Handlers.UpdateProductoDespachoHandler(db))
	api.DELETE("/productos_despacho/:despacho_id/:producto_id", Handlers.DeleteProductoDespachoHandler(db))
}
