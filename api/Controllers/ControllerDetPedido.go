package Controllers

import (
	"net/http"
	"strconv"

	modelos "backend-inventario/api/Models"
	"backend-inventario/api/db"

	"github.com/gin-gonic/gin"
)

// GetDetallesPedido obtiene todos los detalles de pedido
func GetDetallesPedido(c *gin.Context) {
	var detalles []modelos.DetallePedido
	// Preload las relaciones relevantes
	if err := db.DB.Preload("Pedido").Preload("Producto").Find(&detalles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener detalles de pedido"})
		return
	}
	c.JSON(http.StatusOK, detalles)
}

// GetDetallePedidoByID obtiene un detalle de pedido por su ID
func GetDetallePedidoByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var detalle modelos.DetallePedido
	if err := db.DB.Preload("Pedido").Preload("Producto").First(&detalle, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Detalle de pedido no encontrado"})
		return
	}
	c.JSON(http.StatusOK, detalle)
}

// CreateDetallePedido crea un nuevo detalle de pedido
func CreateDetallePedido(c *gin.Context) {
	var detalle modelos.DetallePedido
	if err := c.ShouldBindJSON(&detalle); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Create(&detalle).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear detalle de pedido"})
		return
	}
	// Carga las relaciones para la respuesta
	db.DB.Preload("Pedido").Preload("Producto").First(&detalle, detalle.ID)
	c.JSON(http.StatusCreated, detalle)
}

// UpdateDetallePedido actualiza un detalle de pedido existente
func UpdateDetallePedido(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var detalle modelos.DetallePedido
	if err := db.DB.First(&detalle, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Detalle de pedido no encontrado"})
		return
	}

	if err := c.ShouldBindJSON(&detalle); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Save(&detalle).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar detalle de pedido"})
		return
	}
	// Carga las relaciones para la respuesta
	db.DB.Preload("Pedido").Preload("Producto").First(&detalle, detalle.ID)
	c.JSON(http.StatusOK, detalle)
}

// DeleteDetallePedido elimina un detalle de pedido
func DeleteDetallePedido(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var detalle modelos.DetallePedido
	if err := db.DB.First(&detalle, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Detalle de pedido no encontrado"})
		return
	}

	if err := db.DB.Delete(&detalle).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar detalle de pedido"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Detalle de pedido eliminado exitosamente"})
}
