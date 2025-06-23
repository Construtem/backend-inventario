package Controllers

import (
	"net/http"
	"strconv"

	modelos "backend-inventario/api/Models"
	"backend-inventario/api/db"

	"github.com/gin-gonic/gin"
)

// GetCamiones obtiene todos los camiones
func GetCamiones(c *gin.Context) {
	var camiones []modelos.Camion
	if err := db.DB.Find(&camiones).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener camiones"})
		return
	}
	c.JSON(http.StatusOK, camiones)
}

// GetCamionByID obtiene un camión por su ID
func GetCamionByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var camion modelos.Camion
	if err := db.DB.First(&camion, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Camión no encontrado"})
		return
	}
	c.JSON(http.StatusOK, camion)
}

// CreateCamion crea un nuevo camión
func CreateCamion(c *gin.Context) {
	var camion modelos.Camion
	if err := c.ShouldBindJSON(&camion); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := db.DB.Create(&camion).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear camión"})
		return
	}
	c.JSON(http.StatusCreated, camion)
}

// UpdateCamion actualiza un camión existente
func UpdateCamion(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var camion modelos.Camion
	if err := db.DB.First(&camion, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Camión no encontrado"})
		return
	}

	if err := c.ShouldBindJSON(&camion); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Save(&camion).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar camión"})
		return
	}
	c.JSON(http.StatusOK, camion)
}

// DeleteCamion elimina un camión
func DeleteCamion(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var camion modelos.Camion
	if err := db.DB.First(&camion, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Camión no encontrado"})
		return
	}

	if err := db.DB.Delete(&camion).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar camión"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Camión eliminado exitosamente"})
}
