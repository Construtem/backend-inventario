package Controllers

import (
	"net/http"
	"strconv"
	"time"

	modelos "backend-inventario/api/Models"
	"backend-inventario/api/db"

	"github.com/gin-gonic/gin"
)

// GetPedidos obtiene todos los pedidos
func GetPedidos(c *gin.Context) {
	var pedidos []modelos.Pedido
	// Preload las relaciones relevantes
	if err := db.DB.Preload("Cliente").Preload("Vendedor").Preload("Cotizacion").Preload("Ubicacion").Preload("Despachador").Find(&pedidos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener pedidos"})
		return
	}
	c.JSON(http.StatusOK, pedidos)
}

// GetPedidoByID obtiene un pedido por su ID
func GetPedidoByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var pedido modelos.Pedido
	if err := db.DB.Preload("Cliente").Preload("Vendedor").Preload("Cotizacion").Preload("Ubicacion").Preload("Despachador").First(&pedido, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pedido no encontrado"})
		return
	}
	c.JSON(http.StatusOK, pedido)
}

// CreatePedido crea un nuevo pedido
func CreatePedido(c *gin.Context) {
	var pedido modelos.Pedido
	if err := c.ShouldBindJSON(&pedido); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Establecer la fecha actual si no se proporciona
	if pedido.Fecha.IsZero() {
		pedido.Fecha = time.Now()
	}

	if err := db.DB.Create(&pedido).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear pedido"})
		return
	}
	// Carga las relaciones para la respuesta
	db.DB.Preload("Cliente").Preload("Vendedor").Preload("Cotizacion").Preload("Ubicacion").Preload("Despachador").First(&pedido, pedido.ID)
	c.JSON(http.StatusCreated, pedido)
}

// UpdatePedido actualiza un pedido existente
func UpdatePedido(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var pedido modelos.Pedido
	if err := db.DB.First(&pedido, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pedido no encontrado"})
		return
	}

	if err := c.ShouldBindJSON(&pedido); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Save(&pedido).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar pedido"})
		return
	}
	// Carga las relaciones para la respuesta
	db.DB.Preload("Cliente").Preload("Vendedor").Preload("Cotizacion").Preload("Ubicacion").Preload("Despachador").First(&pedido, pedido.ID)
	c.JSON(http.StatusOK, pedido)
}

// DeletePedido elimina un pedido
func DeletePedido(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var pedido modelos.Pedido
	if err := db.DB.First(&pedido, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pedido no encontrado"})
		return
	}

	if err := db.DB.Delete(&pedido).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar pedido"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Pedido eliminado exitosamente"})
}
