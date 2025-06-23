package Controllers

import (
	"net/http"
	"strconv"
	"time"

	modelos "backend-inventario/api/Models"
	"backend-inventario/api/db"

	"github.com/gin-gonic/gin"
)

// GetDespachos obtiene todos los despachos
func GetDespachos(c *gin.Context) {
	var despachos []modelos.Despacho
	// Preload las relaciones relevantes
	if err := db.DB.Preload("Pedido").Preload("Camion").Preload("Origen").Preload("Destino").Find(&despachos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener despachos"})
		return
	}
	c.JSON(http.StatusOK, despachos)
}

// GetDespachoByID obtiene un despacho por su ID
func GetDespachoByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var despacho modelos.Despacho
	if err := db.DB.Preload("Pedido").Preload("Camion").Preload("Origen").Preload("Destino").First(&despacho, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Despacho no encontrado"})
		return
	}
	c.JSON(http.StatusOK, despacho)
}

// CreateDespacho crea un nuevo despacho
func CreateDespacho(c *gin.Context) {
	var despacho modelos.Despacho
	if err := c.ShouldBindJSON(&despacho); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Establecer la fecha actual si no se proporciona
	if despacho.FechaSalida.IsZero() {
		despacho.FechaSalida = time.Now()
	}

	if err := db.DB.Create(&despacho).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear despacho"})
		return
	}
	// Carga las relaciones para la respuesta
	db.DB.Preload("Pedido").Preload("Camion").Preload("Origen").Preload("Destino").First(&despacho, despacho.ID)
	c.JSON(http.StatusCreated, despacho)
}

// UpdateDespacho actualiza un despacho existente
func UpdateDespacho(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var despacho modelos.Despacho
	if err := db.DB.First(&despacho, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Despacho no encontrado"})
		return
	}

	if err := c.ShouldBindJSON(&despacho); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Save(&despacho).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar despacho"})
		return
	}
	// Carga las relaciones para la respuesta
	db.DB.Preload("Pedido").Preload("Camion").Preload("Origen").Preload("Destino").First(&despacho, despacho.ID)
	c.JSON(http.StatusOK, despacho)
}

// DeleteDespacho elimina un despacho
func DeleteDespacho(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var despacho modelos.Despacho
	if err := db.DB.First(&despacho, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Despacho no encontrado"})
		return
	}

	if err := db.DB.Delete(&despacho).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar despacho"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Despacho eliminado exitosamente"})
}
