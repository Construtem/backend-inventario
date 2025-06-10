package Controllers

import (
	"net/http"
	"strconv"

	modelos "backend-inventario/api/Models"
	"backend-inventario/api/db"

	"github.com/gin-gonic/gin"
)

// GetUbicaciones obtiene todas las ubicaciones
func GetUbicaciones(c *gin.Context) {
	var ubicaciones []modelos.Ubicacion
	if err := db.DB.Find(&ubicaciones).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener ubicaciones"})
		return
	}
	c.JSON(http.StatusOK, ubicaciones)
}

// GetUbicacionByID obtiene una ubicación por su ID
func GetUbicacionByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var ubicacion modelos.Ubicacion
	if err := db.DB.First(&ubicacion, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ubicación no encontrada"})
		return
	}
	c.JSON(http.StatusOK, ubicacion)
}

// CreateUbicacion crea una nueva ubicación
func CreateUbicacion(c *gin.Context) {
	var ubicacion modelos.Ubicacion
	if err := c.ShouldBindJSON(&ubicacion); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := db.DB.Create(&ubicacion).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear ubicación"})
		return
	}
	c.JSON(http.StatusCreated, ubicacion)
}

// UpdateUbicacion actualiza una ubicación existente
func UpdateUbicacion(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var ubicacion modelos.Ubicacion
	if err := db.DB.First(&ubicacion, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ubicación no encontrada"})
		return
	}

	if err := c.ShouldBindJSON(&ubicacion); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Save(&ubicacion).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar ubicación"})
		return
	}
	c.JSON(http.StatusOK, ubicacion)
}

// DeleteUbicacion elimina una ubicación
func DeleteUbicacion(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var ubicacion modelos.Ubicacion
	if err := db.DB.First(&ubicacion, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ubicación no encontrada"})
		return
	}

	if err := db.DB.Delete(&ubicacion).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar ubicación"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Ubicación eliminada exitosamente"})
}
