package Controllers

import (
	"net/http"
	"strconv"

	modelos "backend-inventario/api/Models"
	"backend-inventario/api/db"

	"github.com/gin-gonic/gin"
)

// GetDetallesDespacho obtiene todos los detalles de despacho
func GetDetallesDespacho(c *gin.Context) {
	var detalles []modelos.DetalleDespacho
	// Preload las relaciones relevantes
	if err := db.DB.Preload("Despacho").Preload("Producto").Find(&detalles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener detalles de despacho"})
		return
	}
	c.JSON(http.StatusOK, detalles)
}

// GetDetalleDespachoByID obtiene un detalle de despacho por su ID
func GetDetalleDespachoByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var detalle modelos.DetalleDespacho
	if err := db.DB.Preload("Despacho").Preload("Producto").First(&detalle, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Detalle de despacho no encontrado"})
		return
	}
	c.JSON(http.StatusOK, detalle)
}

// CreateDetalleDespacho crea un nuevo detalle de despacho
func CreateDetalleDespacho(c *gin.Context) {
	var detalle modelos.DetalleDespacho
	if err := c.ShouldBindJSON(&detalle); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Create(&detalle).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear detalle de despacho"})
		return
	}
	// Carga las relaciones para la respuesta
	db.DB.Preload("Despacho").Preload("Producto").First(&detalle, detalle.ID)
	c.JSON(http.StatusCreated, detalle)
}

// UpdateDetalleDespacho actualiza un detalle de despacho existente
func UpdateDetalleDespacho(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var detalle modelos.DetalleDespacho
	if err := db.DB.First(&detalle, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Detalle de despacho no encontrado"})
		return
	}

	if err := c.ShouldBindJSON(&detalle); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Save(&detalle).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar detalle de despacho"})
		return
	}
	// Carga las relaciones para la respuesta
	db.DB.Preload("Despacho").Preload("Producto").First(&detalle, detalle.ID)
	c.JSON(http.StatusOK, detalle)
}

// DeleteDetalleDespacho elimina un detalle de despacho
func DeleteDetalleDespacho(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var detalle modelos.DetalleDespacho
	if err := db.DB.First(&detalle, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Detalle de despacho no encontrado"})
		return
	}

	if err := db.DB.Delete(&detalle).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar detalle de despacho"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Detalle de despacho eliminado exitosamente"})
}
