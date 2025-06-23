package Controllers

import (
	"net/http"
	"strconv"

	modelos "backend-inventario/api/Models"
	"backend-inventario/api/db"

	"github.com/gin-gonic/gin"
)

// GetRoles obtiene todos los roles
func GetRoles(c *gin.Context) {
	var roles []modelos.Rol
	if err := db.DB.Find(&roles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener roles"})
		return
	}
	c.JSON(http.StatusOK, roles)
}

// GetRolByID obtiene un rol por su ID
func GetRolByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var rol modelos.Rol
	if err := db.DB.First(&rol, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Rol no encontrado"})
		return
	}
	c.JSON(http.StatusOK, rol)
}

// CreateRol crea un nuevo rol
func CreateRol(c *gin.Context) {
	var rol modelos.Rol
	if err := c.ShouldBindJSON(&rol); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := db.DB.Create(&rol).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear rol"})
		return
	}
	c.JSON(http.StatusCreated, rol)
}

// UpdateRol actualiza un rol existente
func UpdateRol(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var rol modelos.Rol
	if err := db.DB.First(&rol, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Rol no encontrado"})
		return
	}

	if err := c.ShouldBindJSON(&rol); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Save(&rol).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar rol"})
		return
	}
	c.JSON(http.StatusOK, rol)
}

// DeleteRol elimina un rol
func DeleteRol(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var rol modelos.Rol
	if err := db.DB.First(&rol, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Rol no encontrado"})
		return
	}

	if err := db.DB.Delete(&rol).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar rol"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Rol eliminado exitosamente"})
}
