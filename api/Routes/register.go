package Routes

import (
	"backend-inventario/api/Controllers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	// Configuración de CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Permite tu frontend de Next.js
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true, MaxAge: 86400, // 24 horas
	}))

	// Grupo de rutas para la API
	api := router.Group("/api")

	// Rutas para Clientes
	api.GET("/clientes", Controllers.GetClientes)
	api.GET("/clientes/:id", Controllers.GetClienteByID)
	api.POST("/clientes", Controllers.CreateCliente)
	api.PUT("/clientes/:id", Controllers.UpdateCliente)
	api.DELETE("/clientes/:id", Controllers.DeleteCliente)

	// Rutas para Productos
	api.GET("/productos", Controllers.GetProductos)
	api.GET("/productos/:id", Controllers.GetProductoByID)
	api.POST("/productos", Controllers.CreateProducto)
	api.PUT("/productos/:id", Controllers.UpdateProducto)
	api.DELETE("/productos/:id", Controllers.DeleteProducto)

	// Rutas para Categorías
	api.GET("/categorias", Controllers.GetCategorias)
	api.GET("/categorias/:id", Controllers.GetCategoriaByID)
	api.POST("/categorias", Controllers.CreateCategoria)
	api.PUT("/categorias/:id", Controllers.UpdateCategoria)
	api.DELETE("/categorias/:id", Controllers.DeleteCategoria)

	// Rutas para Usuarios
	api.GET("/usuarios", Controllers.GetUsuarios)
	api.POST("/usuarios", Controllers.CreateUsuario)
	api.GET("/usuarios/:id", Controllers.GetUsuarioByID)
	api.PUT("/usuarios/:id", Controllers.UpdateUsuario)
	api.DELETE("/usuarios/:id", Controllers.DeleteUsuario)

	// Rutas para Roles
	api.GET("/roles", Controllers.GetRoles)
	api.GET("/roles/:id", Controllers.GetRolByID)
	api.POST("/roles", Controllers.CreateRol)
	api.PUT("/roles/:id", Controllers.UpdateRol)
	api.DELETE("/roles/:id", Controllers.DeleteRol)

	// Rutas para Ubicaciones
	api.GET("/ubicaciones", Controllers.GetUbicaciones)
	api.GET("/ubicaciones/:id", Controllers.GetUbicacionByID)
	api.POST("/ubicaciones", Controllers.CreateUbicacion)
	api.PUT("/ubicaciones/:id", Controllers.UpdateUbicacion)
	api.DELETE("/ubicaciones/:id", Controllers.DeleteUbicacion)

	// Rutas para Proveedores
	api.GET("/proveedores", Controllers.GetProveedores)
	api.GET("/proveedores/:id", Controllers.GetProveedorByID)
	api.POST("/proveedores", Controllers.CreateProveedor)
	api.PUT("/proveedores/:id", Controllers.UpdateProveedor)
	api.DELETE("/proveedores/:id", Controllers.DeleteProveedor)

	// Rutas para Camiones
	api.GET("/camiones", Controllers.GetCamiones)
	api.GET("/camiones/:id", Controllers.GetCamionByID)
	api.POST("/camiones", Controllers.CreateCamion)
	api.PUT("/camiones/:id", Controllers.UpdateCamion)
	api.DELETE("/camiones/:id", Controllers.DeleteCamion)

	// Rutas para Pedidos
	api.GET("/pedidos", Controllers.GetPedidos)
	api.GET("/pedidos/:id", Controllers.GetPedidoByID)
	api.POST("/pedidos", Controllers.CreatePedido)
	api.PUT("/pedidos/:id", Controllers.UpdatePedido)
	api.DELETE("/pedidos/:id", Controllers.DeletePedido)

	// Rutas para Cotizaciones
	api.GET("/cotizaciones", Controllers.GetCotizaciones)
	api.GET("/cotizaciones/:id", Controllers.GetCotizacionByID)
	api.POST("/cotizaciones", Controllers.CreateCotizacion)
	api.PUT("/cotizaciones/:id", Controllers.UpdateCotizacion)
	api.DELETE("/cotizaciones/:id", Controllers.DeleteCotizacion)

	// Rutas para Ordenes de Compra
	api.GET("/ordenes_compra", Controllers.GetOrdenesCompra)
	api.GET("/ordenes_compra/:id", Controllers.GetOrdenCompraByID)
	api.POST("/ordenes_compra", Controllers.CreateOrdenCompra)
	api.PUT("/ordenes_compra/:id", Controllers.UpdateOrdenCompra)
	api.DELETE("/ordenes_compra/:id", Controllers.DeleteOrdenCompra)

	// Rutas para Despachos
	api.GET("/despachos", Controllers.GetDespachos)
	api.GET("/despachos/:id", Controllers.GetDespachoByID)
	api.POST("/despachos", Controllers.CreateDespacho)
	api.PUT("/despachos/:id", Controllers.UpdateDespacho)
	api.DELETE("/despachos/:id", Controllers.DeleteDespacho)

	// Rutas para Inventario
	api.GET("/inventario", Controllers.GetInventario)
	api.GET("/inventario/:id", Controllers.GetInventarioByID)
	api.POST("/inventario", Controllers.CreateInventario)
	api.PUT("/inventario/:id", Controllers.UpdateInventario)
	api.DELETE("/inventario/:id", Controllers.DeleteInventario)

	// Rutas para DetalleCotizacion (Considera si necesitas CRUD completo para detalles o solo dentro del contexto de la entidad padre)
	api.GET("/detalle_cotizacion", Controllers.GetDetallesCotizacion)
	api.GET("/detalle_cotizacion/:id", Controllers.GetDetalleCotizacionByID)
	api.POST("/detalle_cotizacion", Controllers.CreateDetalleCotizacion)
	api.PUT("/detalle_cotizacion/:id", Controllers.UpdateDetalleCotizacion)
	api.DELETE("/detalle_cotizacion/:id", Controllers.DeleteDetalleCotizacion)

	// Rutas para DetalleOrdenCompra
	api.GET("/detalle_orden_compra", Controllers.GetDetallesOrdenCompra)
	api.GET("/detalle_orden_compra/:id", Controllers.GetDetalleOrdenCompraByID)
	api.POST("/detalle_orden_compra", Controllers.CreateDetalleOrdenCompra)
	api.PUT("/detalle_orden_compra/:id", Controllers.UpdateDetalleOrdenCompra)
	api.DELETE("/detalle_orden_compra/:id", Controllers.DeleteDetalleOrdenCompra)

	// Rutas para DetallePedido
	api.GET("/detalle_pedido", Controllers.GetDetallesPedido)
	api.GET("/detalle_pedido/:id", Controllers.GetDetallePedidoByID)
	api.POST("/detalle_pedido", Controllers.CreateDetallePedido)
	api.PUT("/detalle_pedido/:id", Controllers.UpdateDetallePedido)
	api.DELETE("/detalle_pedido/:id", Controllers.DeleteDetallePedido)

	// Rutas para DetalleDespacho
	api.GET("/detalle_despacho", Controllers.GetDetallesDespacho)
	api.GET("/detalle_despacho/:id", Controllers.GetDetalleDespachoByID)
	api.POST("/detalle_despacho", Controllers.CreateDetalleDespacho)
	api.PUT("/detalle_despacho/:id", Controllers.UpdateDetalleDespacho)
	api.DELETE("/detalle_despacho/:id", Controllers.DeleteDetalleDespacho)

}
