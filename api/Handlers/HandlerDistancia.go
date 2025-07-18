package Handlers

import (
	"backend-inventario/api/Controllers"
	modelos "backend-inventario/api/Models"
	"backend-inventario/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetDespachoDistanciaByIDHandler maneja GET /api/despachos-distancia/{id}
func GetDespachoDistanciaByIDHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "ID inválido",
				"details": "El ID debe ser un número entero válido",
			})
			return
		}

		despacho, err := Controllers.GetDespachoDistanciaByID(db, uint(id))
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{
					"error":   "Despacho no encontrado",
					"details": "No existe un despacho con el ID especificado",
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Error interno del servidor",
				"details": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, despacho)
	}
}

// GetDespachosDistanciaHandler maneja GET /api/despachos-distancia
func GetDespachosDistanciaHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		despachos, err := Controllers.GetDespachosDistancia(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Error al obtener despachos",
				"details": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, despachos)
	}
}

// CalcularDistanciaDespachoHandler maneja POST /api/despachos/{id}/calcular-distancia
func CalcularDistanciaDespachoHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "ID inválido",
				"details": "El ID debe ser un número entero válido",
			})
			return
		}

		var request modelos.DireccionesRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Datos inválidos",
				"details": err.Error(),
			})
			return
		}

		// Crear servicio de Google Maps
		googleMapsService := services.NewGoogleMapsService()

		// Validar direcciones
		if _, err := googleMapsService.ValidarDireccion(request.Origen); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Dirección de origen inválida",
				"details": err.Error(),
			})
			return
		}

		if _, err := googleMapsService.ValidarDireccion(request.Destino); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Dirección de destino inválida",
				"details": err.Error(),
			})
			return
		}

		// Calcular distancia usando Google Maps
		distanciaCalculada, err := googleMapsService.CalcularDistancia(request.Origen, request.Destino)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Error al calcular distancia",
				"details": err.Error(),
			})
			return
		}

		// Actualizar el despacho con la distancia calculada
		err = Controllers.ActualizarDistanciaDespacho(db, uint(id), distanciaCalculada.Distancia, distanciaCalculada.Duracion)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Error al actualizar despacho",
				"details": err.Error(),
			})
			return
		}

		// Preparar respuesta
		response := modelos.DistanciaResponse{
			Distancia:      distanciaCalculada.Distancia,
			Duracion:       distanciaCalculada.Duracion,
			RutaOptimizada: distanciaCalculada.RutaOptimizada,
		}

		c.JSON(http.StatusOK, response)
	}
}

// CalcularDistanciaAutomaticoHandler calcula automáticamente la distancia usando las direcciones del despacho
func CalcularDistanciaAutomaticoHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "ID inválido",
				"details": "El ID debe ser un número entero válido",
			})
			return
		}

		// Obtener el despacho con sus direcciones
		despacho, err := Controllers.GetDespachoDistanciaByID(db, uint(id))
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{
					"error":   "Despacho no encontrado",
					"details": "No existe un despacho con el ID especificado",
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Error interno del servidor",
				"details": err.Error(),
			})
			return
		}

		// Verificar que tengamos las direcciones necesarias
		if despacho.OrigenSucursal == nil || despacho.DestinoDirCliente == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Información incompleta",
				"details": "El despacho debe tener información completa de origen y destino",
			})
			return
		}

		// Formatear direcciones
		origenCompleto := services.FormatearDireccionCompleta(
			despacho.OrigenSucursal.Direccion,
			despacho.OrigenSucursal.Comuna,
			despacho.OrigenSucursal.Ciudad,
		)

		destinoCompleto := services.FormatearDireccionCompleta(
			despacho.DestinoDirCliente.Direccion,
			despacho.DestinoDirCliente.Comuna,
			despacho.DestinoDirCliente.Ciudad,
		)

		// Crear servicio de Google Maps
		googleMapsService := services.NewGoogleMapsService()

		// Calcular distancia
		distanciaCalculada, err := googleMapsService.CalcularDistancia(origenCompleto, destinoCompleto)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Error al calcular distancia",
				"details": err.Error(),
			})
			return
		}

		// Actualizar el despacho con la distancia calculada
		err = Controllers.ActualizarDistanciaDespacho(db, uint(id), distanciaCalculada.Distancia, distanciaCalculada.Duracion)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Error al actualizar despacho",
				"details": err.Error(),
			})
			return
		}

		// Preparar respuesta
		response := modelos.DistanciaResponse{
			Distancia:      distanciaCalculada.Distancia,
			Duracion:       distanciaCalculada.Duracion,
			RutaOptimizada: distanciaCalculada.RutaOptimizada,
		}

		c.JSON(http.StatusOK, response)
	}
}
