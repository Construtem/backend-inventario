package Handlers

import (
	"backend-inventario/api/Controllers"
	modelos "backend-inventario/api/Models"
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
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Los datos del despacho enviados no son válidos.",
				"details": err.Error(),
			})
			return
		}

		if err := Controllers.CreateDespacho(db, &request.Despacho, request.Productos); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "No se pudo registrar el despacho.",
				"details": err.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message":  "Despacho creado exitosamente.",
			"despacho": request.Despacho,
		})
	}
}

func GetDespachosHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		despachos, err := Controllers.GetDespachos(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Hubo un problema al obtener los despachos.",
				"details": err.Error(),
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
			c.JSON(http.StatusBadRequest, gin.H{"error": "El ID proporcionado no es válido."})
			return
		}

		despacho, err := Controllers.GetDespachoByID(db, uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "No se encontró un despacho con el ID especificado.",
				"details": err.Error(),
			})
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
			c.JSON(http.StatusBadRequest, gin.H{"error": "El ID proporcionado no es válido."})
			return
		}

		var actualizado modelos.Despacho
		if err := c.ShouldBindJSON(&actualizado); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Los datos enviados para actualizar el despacho no son válidos.",
				"details": err.Error(),
			})
			return
		}

		if err := Controllers.UpdateDespacho(db, uint(id), &actualizado); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "No se pudo actualizar el despacho.",
				"details": err.Error(),
			})
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
			c.JSON(http.StatusBadRequest, gin.H{"error": "El ID proporcionado no es válido."})
			return
		}

		if err := Controllers.DeleteDespacho(db, uint(id)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "No se pudo eliminar el despacho.",
				"details": err.Error(),
			})
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
			c.JSON(http.StatusBadRequest, gin.H{
				"error":    "Solicitud inválida",
				"mensaje":  "No se pudo procesar el cálculo. Verifica los datos enviados.",
				"detalles": err.Error(),
			})
			return
		}

		despacho, err := Controllers.CalcularDespacho(db, req.CotizacionID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":    "Error interno",
				"mensaje":  "No fue posible calcular el despacho.",
				"detalles": err.Error(),
			})
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
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "ID inválido",
				"mensaje": "El ID de cotización proporcionado no es válido.",
			})
			return
		}

		despachos, err := Controllers.GetDespachosPorCotizacion(db, uint(id))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":    "Error interno",
				"mensaje":  "No se pudieron obtener los despachos asociados.",
				"detalles": err.Error(),
			})
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
			c.JSON(http.StatusBadRequest, gin.H{
				"error":    "Solicitud inválida",
				"mensaje":  "No se pudo procesar la aprobación del despacho. Verifica los datos enviados.",
				"detalles": err.Error(),
			})
			return
		}

		if err := Controllers.AprobarDespacho(db, req.CotizacionID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":    "Error al aprobar",
				"mensaje":  "No se pudo aprobar el despacho.",
				"detalles": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Despacho aprobado exitosamente"})
	}
}

// Cambia el estado de los despachos asociados a una cotización
func CambiarEstadoDespachosHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			CotizacionID uint   `json:"cotizacion_id"`
			Estado       string `json:"estado"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":    "Solicitud inválida",
				"mensaje":  "No se pudo cambiar el estado. Verifica los datos enviados.",
				"detalles": err.Error(),
			})
			return
		}
		err := Controllers.CambiarEstadoDespachosPorCotizacion(db, req.CotizacionID, req.Estado)
		if err != nil {
			switch err.Error() {
			case "Estado no permitido":
				c.JSON(http.StatusBadRequest, gin.H{"error": "Estado no permitido", "mensaje": "El estado enviado no es válido."})
			case "No se encontraron despachos para la cotización":
				c.JSON(http.StatusNotFound, gin.H{"error": "No encontrado", "mensaje": "No hay despachos registrados para la cotización indicada."})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":    "Error interno",
					"mensaje":  "No se pudo actualizar el estado de los despachos.",
					"detalles": err.Error(),
				})
			}
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Estado de despachos actualizado exitosamente"})
	}
}

// GetFichaDespachoHandler retorna la ficha de despacho y la factura electrónica en un solo JSON
func GetFichaDespachoHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "ID inválido",
				"mensaje": "El ID de despacho proporcionado no es válido.",
			})
			return
		}

		// Obtener la ficha de despacho
		ficha, err := Controllers.GetDespachoByID(db, uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error":    "No encontrado",
				"mensaje":  "No se encontró el despacho solicitado.",
				"detalles": err.Error(),
			})
			return
		}

		// Obtener la factura electrónica usando el controlador
		factura, err := Controllers.GetFacturaElectronicaByDespachoID(db, uint(id))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":    "Error al obtener factura",
				"mensaje":  "No se pudo recuperar la factura electrónica del despacho.",
				"detalles": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"ficha_despacho":      ficha,
			"factura_electronica": factura,
		})
	}
}
