package Controllers

import (
	"net/http"
	"strconv"
	"time"

	modelos "backend-inventario/api/Models"
	"backend-inventario/api/db"

	"github.com/gin-gonic/gin"
)

// GetOrdenesCompra obtiene todas las ordenes de compra
func GetOrdenesCompra(c *gin.Context) {
	var ordenesCompra []modelos.OrdenCompra
	// Preload las relaciones relevantes
	if err := db.DB.Preload("Proveedor").Preload("Solicitante").Preload("Ubicacion").Preload("Receptor").Find(&ordenesCompra).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener ordenes de compra"})
		return
	}
	c.JSON(http.StatusOK, ordenesCompra)
}

// GetOrdenCompraByID obtiene una orden de compra por su ID
func GetOrdenCompraByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var ordenCompra modelos.OrdenCompra
	if err := db.DB.Preload("Proveedor").Preload("Solicitante").Preload("Ubicacion").Preload("Receptor").First(&ordenCompra, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Orden de compra no encontrada"})
		return
	}
	c.JSON(http.StatusOK, ordenCompra)
}

// CreateOrdenCompra crea una nueva orden de compra
func CreateOrdenCompra(c *gin.Context) {
	var ordenCompra modelos.OrdenCompra
	if err := c.ShouldBindJSON(&ordenCompra); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Establecer la fecha actual si no se proporciona
	if ordenCompra.Fecha.IsZero() {
		ordenCompra.Fecha = time.Now()
	}

	if err := db.DB.Create(&ordenCompra).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear orden de compra"})
		return
	}
	// Carga las relaciones para la respuesta
	db.DB.Preload("Proveedor").Preload("Solicitante").Preload("Ubicacion").Preload("Receptor").First(&ordenCompra, ordenCompra.ID)
	c.JSON(http.StatusCreated, ordenCompra)
}

// UpdateOrdenCompra actualiza una orden de compra existente
func UpdateOrdenCompra(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var ordenCompra modelos.OrdenCompra
	if err := db.DB.First(&ordenCompra, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Orden de compra no encontrada"})
		return
	}

	if err := c.ShouldBindJSON(&ordenCompra); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Save(&ordenCompra).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar orden de compra"})
		return
	}
	// Carga las relaciones para la respuesta
	db.DB.Preload("Proveedor").Preload("Solicitante").Preload("Ubicacion").Preload("Receptor").First(&ordenCompra, ordenCompra.ID)
	c.JSON(http.StatusOK, ordenCompra)
}

// DeleteOrdenCompra elimina una orden de compra
func DeleteOrdenCompra(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var ordenCompra modelos.OrdenCompra
	if err := db.DB.First(&ordenCompra, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Orden de compra no encontrada"})
		return
	}

	if err := db.DB.Delete(&ordenCompra).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar orden de compra"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Orden de compra eliminada exitosamente"})
}
