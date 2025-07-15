package Handlers

import (
	"backend-inventario/api/Controllers"
	modelos "backend-inventario/api/Models"
	"backend-inventario/config"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DespachoRequest struct {
	Despacho  modelos.Despacho            `json:"valor_despacho"`
	Productos []modelos.ProductosDespacho `json:"productos"`
}

func CreateDespachoHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request DespachoRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Formato inválido", "details": err.Error()})
			return
		}

		if err := Controllers.CreateDespacho(db, &request.Despacho, request.Productos); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear despacho", "details": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Despacho creado exitosamente", "despacho_id": request.Despacho.ID})
	}
}

func GetDespachosHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		despachos, err := Controllers.GetDespachos(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Error al obtener despachos",
				"details": err.Error(),
				"message": "Error interno del servidor al consultar despachos",
			})
			return
		}
		c.JSON(http.StatusOK, despachos)
	}
}

func GetDespachoByIDHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
			return
		}

		despacho, err := Controllers.GetDespachoByID(db, uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Despacho no encontrado", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, despacho)
	}
}

func UpdateDespachoHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
			return
		}

		var actualizado modelos.Despacho
		if err := c.ShouldBindJSON(&actualizado); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos", "details": err.Error()})
			return
		}

		if err := Controllers.UpdateDespacho(db, uint(id), &actualizado); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar despacho", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Despacho actualizado exitosamente"})
	}
}

func DeleteDespachoHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
			return
		}

		if err := Controllers.DeleteDespacho(db, uint(id)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar despacho", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Despacho eliminado exitosamente"})
	}
}

func CalcularDespachoHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		type Req struct {
			CotizacionID uint `json:"cotizacion_id"`
		}

		var req Req
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Solicitud inválida", "details": err.Error()})
			return
		}

		despacho, err := Controllers.CalcularDespacho(db, req.CotizacionID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al calcular despacho", "details": err.Error()})
			return
		}

		c.JSON(http.StatusOK, despacho)
	}
}

func GetDespachosPorCotizacionHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
			return
		}

		despachos, err := Controllers.GetDespachosPorCotizacion(db, uint(id))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener despachos", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, despachos)
	}
}

func AprobarDespachoHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			CotizacionID uint `json:"cotizacion_id"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Solicitud inválida", "details": err.Error()})
			return
		}

		if err := Controllers.AprobarDespacho(db, req.CotizacionID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al aprobar despacho", "details": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Despacho aprobado exitosamente"})
	}
}

func GetDespachosConDistanciaHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener la API key desde la configuración
		apiKey := config.GetGoogleMapsAPIKey()

		despachos, err := Controllers.GetDespachosConDistancia(db, apiKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Error al obtener despachos",
				"details": err.Error(),
				"message": "Error interno del servidor al consultar despachos",
			})
			return
		}
		c.JSON(http.StatusOK, despachos)
	}
}

func GetDespachoByIDConDistanciaHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
			return
		}

		// Obtener la API key desde la configuración
		apiKey := config.GetGoogleMapsAPIKey()

		despacho, err := Controllers.GetDespachoByIDConDistancia(db, uint(id), apiKey)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Despacho no encontrado", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, despacho)
	}
}

func CalcularDespachoConDistanciaHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		type Req struct {
			CotizacionID uint `json:"cotizacion_id"`
		}

		var req Req
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Solicitud inválida", "details": err.Error()})
			return
		}

		// Obtener la API key desde la configuración
		apiKey := config.GetGoogleMapsAPIKey()

		despacho, err := Controllers.CalcularDespachoConDistancia(db, req.CotizacionID, apiKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al calcular despacho", "details": err.Error()})
			return
		}

		c.JSON(http.StatusOK, despacho)
	}
}

func GetDespachosPorCotizacionDetalladoHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
			return
		}

		// Obtener la API key desde la configuración
		apiKey := config.GetGoogleMapsAPIKey()

		despachos, err := Controllers.GetDespachosPorCotizacionDetallado(db, uint(id), apiKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener despachos", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, despachos)
	}
}

// Handler para calcular distancia entre dos direcciones específicas
func CalcularDistanciaHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		type Req struct {
			Origen  string `json:"origen" binding:"required"`
			Destino string `json:"destino" binding:"required"`
		}

		var req Req
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Solicitud inválida - campos requeridos: origen, destino", 
				"details": err.Error(),
				"ejemplo": gin.H{
					"origen": "Av. Providencia 1234, Santiago",
					"destino": "Las Condes, Santiago",
				},
			})
			return
		}

		// Obtener la API key desde la configuración
		apiKey := config.GetGoogleMapsAPIKey()
		if apiKey == "" {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "API key de Google Maps no configurada",
				"message": "Asegúrate de configurar GOOGLE_MAPS_API_KEY en las variables de entorno",
			})
			return
		}

		fmt.Printf("🗺️ Calculando distancia de '%s' a '%s' con API key: %s...\n", req.Origen, req.Destino, apiKey[:10]+"...")

		distancia, duracion, err := Controllers.CalcularDistancia(apiKey, req.Origen, req.Destino)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error al calcular distancia", 
				"details": err.Error(),
				"origen": req.Origen,
				"destino": req.Destino,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"distancia_km":      distancia,
			"duracion_minutos":  duracion,
			"origen":           req.Origen,
			"destino":          req.Destino,
			"status":           "success",
		})
	}
}

// Handler para probar distancia con un despacho real específico
func ProbarDistanciaDespachoRealHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener el primer despacho disponible
		var despacho modelos.Despacho
		err := db.
			Preload("OrigenSucursal").
			Preload("DestinoDirCliente").
			First(&despacho).Error
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "No se encontraron despachos",
				"details": err.Error(),
			})
			return
		}

		// Obtener la API key
		apiKey := config.GetGoogleMapsAPIKey()
		if apiKey == "" {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "API key de Google Maps no configurada",
				"message": "Configura GOOGLE_MAPS_API_KEY en las variables de entorno",
			})
			return
		}

		// Construir las direcciones
		origen := despacho.OrigenSucursal.Direccion + ", " + despacho.OrigenSucursal.Comuna + ", " + despacho.OrigenSucursal.Ciudad
		destino := despacho.DestinoDirCliente.Direccion + ", " + despacho.DestinoDirCliente.Comuna + ", " + despacho.DestinoDirCliente.Ciudad

		fmt.Printf("🚛 Calculando distancia para despacho #%d\n", despacho.ID)
		fmt.Printf("📍 Origen: %s\n", origen)
		fmt.Printf("📍 Destino: %s\n", destino)

		// Calcular distancia
		distancia, duracion, err := Controllers.CalcularDistancia(apiKey, origen, destino)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error al calcular distancia",
				"details": err.Error(),
				"despacho_id": despacho.ID,
				"origen": origen,
				"destino": destino,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"despacho_id":       despacho.ID,
			"origen":           origen,
			"destino":          destino,
			"distancia_km":     distancia,
			"duracion_minutos": duracion,
			"status":           "success",
			"sucursal_origen": gin.H{
				"id":        despacho.OrigenSucursal.ID,
				"nombre":    despacho.OrigenSucursal.Nombre,
				"direccion": despacho.OrigenSucursal.Direccion,
				"comuna":    despacho.OrigenSucursal.Comuna,
				"ciudad":    despacho.OrigenSucursal.Ciudad,
			},
			"cliente_destino": gin.H{
				"id":        despacho.DestinoDirCliente.ID,
				"nombre":    despacho.DestinoDirCliente.Nombre,
				"direccion": despacho.DestinoDirCliente.Direccion,
				"comuna":    despacho.DestinoDirCliente.Comuna,
				"ciudad":    despacho.DestinoDirCliente.Ciudad,
			},
		})
	}
}
