package Controllers

import (
	"net/http"
	"strconv"

	modelos "backend-inventario/api/Models"
	"backend-inventario/api/db"

	"github.com/gin-gonic/gin"
)

// GetDetallesOrdenCompra obtiene todos los detalles de orden de compra
func GetDetallesOrdenCompra(c *gin.Context) {
	var detalles []modelos.DetalleOrdenCompra
	// Preload las relaciones relevantes
	if err := db.DB.Preload("OrdenCompra").Preload("Producto").Find(&detalles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener detalles de orden de compra"})
		return
	}
	c.JSON(http.StatusOK, detalles)
}

// GetDetalleOrdenCompraByID obtiene un detalle de orden de compra por su ID
func GetDetalleOrdenCompraByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var detalle modelos.DetalleOrdenCompra
	if err := db.DB.Preload("OrdenCompra").Preload("Producto").First(&detalle, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Detalle de orden de compra no encontrado"})
		return
	}
	c.JSON(http.StatusOK, detalle)
}

// CreateDetalleOrdenCompra crea un nuevo detalle de orden de compra
func CreateDetalleOrdenCompra(c *gin.Context) {
	var detalle modelos.DetalleOrdenCompra
	if err := c.ShouldBindJSON(&detalle); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Create(&detalle).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear detalle de orden de compra"})
		return
	}
	// Carga las relaciones para la respuesta
	db.DB.Preload("OrdenCompra").Preload("Producto").First(&detalle, detalle.ID)
	c.JSON(http.StatusCreated, detalle)
}

// UpdateDetalleOrdenCompra actualiza un detalle de orden de compra existente
func UpdateDetalleOrdenCompra(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var detalle modelos.DetalleOrdenCompra
	if err := db.DB.First(&detalle, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Detalle de orden de compra no encontrado"})
		return
	}

	if err := c.ShouldBindJSON(&detalle); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Save(&detalle).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar detalle de orden de compra"})
		return
	}
	// Carga las relaciones para la respuesta
	db.DB.Preload("OrdenCompra").Preload("Producto").First(&detalle, detalle.ID)
	c.JSON(http.StatusOK, detalle)
}

// DeleteDetalleOrdenCompra elimina un detalle de orden de compra
func DeleteDetalleOrdenCompra(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var detalle modelos.DetalleOrdenCompra
	if err := db.DB.First(&detalle, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Detalle de orden de compra no encontrado"})
		return
	}

	if err := db.DB.Delete(&detalle).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar detalle de orden de compra"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Detalle de orden de compra eliminado exitosamente"})
}
