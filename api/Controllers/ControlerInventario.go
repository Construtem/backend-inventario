package Controllers

import (
	"net/http"
	"strconv"

	modelos "backend-inventario/api/Models"
	"backend-inventario/api/db"

	"github.com/gin-gonic/gin"
)

// GetInventario obtiene todos los registros de inventario
func GetInventario(c *gin.Context) {
	var inventarios []modelos.Inventario
	// Preload las relaciones relevantes
	if err := db.DB.Preload("Producto").Preload("Ubicacion").Find(&inventarios).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener inventario"})
		return
	}
	c.JSON(http.StatusOK, inventarios)
}

// GetInventarioByID obtiene un registro de inventario por su ID
func GetInventarioByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var inventario modelos.Inventario
	if err := db.DB.Preload("Producto").Preload("Ubicacion").First(&inventario, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Registro de inventario no encontrado"})
		return
	}
	c.JSON(http.StatusOK, inventario)
}

// CreateInventario crea un nuevo registro de inventario
func CreateInventario(c *gin.Context) {
	var inventario modelos.Inventario
	if err := c.ShouldBindJSON(&inventario); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Considerar validaciones adicionales, como si el ProductoID y UbicacionID existen

	if err := db.DB.Create(&inventario).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear registro de inventario"})
		return
	}
	// Carga las relaciones para la respuesta
	db.DB.Preload("Producto").Preload("Ubicacion").First(&inventario, inventario.ID)
	c.JSON(http.StatusCreated, inventario)
}

// UpdateInventario actualiza un registro de inventario existente
func UpdateInventario(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var inventario modelos.Inventario
	if err := db.DB.First(&inventario, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Registro de inventario no encontrado"})
		return
	}

	if err := c.ShouldBindJSON(&inventario); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Save(&inventario).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar registro de inventario"})
		return
	}
	// Carga las relaciones para la respuesta
	db.DB.Preload("Producto").Preload("Ubicacion").First(&inventario, inventario.ID)
	c.JSON(http.StatusOK, inventario)
}

// DeleteInventario elimina un registro de inventario
func DeleteInventario(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var inventario modelos.Inventario
	if err := db.DB.First(&inventario, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Registro de inventario no encontrado"})
		return
	}

	if err := db.DB.Delete(&inventario).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar registro de inventario"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Registro de inventario eliminado exitosamente"})
}
