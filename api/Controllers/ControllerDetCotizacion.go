package Controllers

import (
	"net/http"
	"strconv"

	modelos "backend-inventario/api/Models"
	"backend-inventario/api/db"

	"github.com/gin-gonic/gin"
)

// GetDetallesCotizacion obtiene todos los detalles de cotización
func GetDetallesCotizacion(c *gin.Context) {
	var detalles []modelos.DetalleCotizacion
	// Preload las relaciones relevantes
	if err := db.DB.Preload("Cotizacion").Preload("Producto").Find(&detalles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener detalles de cotización"})
		return
	}
	c.JSON(http.StatusOK, detalles)
}

// GetDetalleCotizacionByID obtiene un detalle de cotización por su ID
func GetDetalleCotizacionByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var detalle modelos.DetalleCotizacion
	if err := db.DB.Preload("Cotizacion").Preload("Producto").First(&detalle, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Detalle de cotización no encontrado"})
		return
	}
	c.JSON(http.StatusOK, detalle)
}

// CreateDetalleCotizacion crea un nuevo detalle de cotización
func CreateDetalleCotizacion(c *gin.Context) {
	var detalle modelos.DetalleCotizacion
	if err := c.ShouldBindJSON(&detalle); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Create(&detalle).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear detalle de cotización"})
		return
	}
	// Carga las relaciones para la respuesta
	db.DB.Preload("Cotizacion").Preload("Producto").First(&detalle, detalle.ID)
	c.JSON(http.StatusCreated, detalle)
}

// UpdateDetalleCotizacion actualiza un detalle de cotización existente
func UpdateDetalleCotizacion(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var detalle modelos.DetalleCotizacion
	if err := db.DB.First(&detalle, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Detalle de cotización no encontrado"})
		return
	}

	if err := c.ShouldBindJSON(&detalle); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Save(&detalle).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar detalle de cotización"})
		return
	}
	// Carga las relaciones para la respuesta
	db.DB.Preload("Cotizacion").Preload("Producto").First(&detalle, detalle.ID)
	c.JSON(http.StatusOK, detalle)
}

// DeleteDetalleCotizacion elimina un detalle de cotización
func DeleteDetalleCotizacion(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var detalle modelos.DetalleCotizacion
	if err := db.DB.First(&detalle, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Detalle de cotización no encontrado"})
		return
	}

	if err := db.DB.Delete(&detalle).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar detalle de cotización"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Detalle de cotización eliminado exitosamente"})
}
