package Controllers

import (
	"net/http"
	"strconv"

	modelos "backend-inventario/api/Models"
	"backend-inventario/api/db"

	"github.com/gin-gonic/gin"
)

// GetProveedores obtiene todos los proveedores
func GetProveedores(c *gin.Context) {
	var proveedores []modelos.Proveedor
	if err := db.DB.Find(&proveedores).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener proveedores"})
		return
	}
	c.JSON(http.StatusOK, proveedores)
}

// GetProveedorByID obtiene un proveedor por su ID
func GetProveedorByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var proveedor modelos.Proveedor
	if err := db.DB.First(&proveedor, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Proveedor no encontrado"})
		return
	}
	c.JSON(http.StatusOK, proveedor)
}

// CreateProveedor crea un nuevo proveedor
func CreateProveedor(c *gin.Context) {
	var proveedor modelos.Proveedor
	if err := c.ShouldBindJSON(&proveedor); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := db.DB.Create(&proveedor).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear proveedor"})
		return
	}
	c.JSON(http.StatusCreated, proveedor)
}

// UpdateProveedor actualiza un proveedor existente
func UpdateProveedor(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var proveedor modelos.Proveedor
	if err := db.DB.First(&proveedor, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Proveedor no encontrado"})
		return
	}

	if err := c.ShouldBindJSON(&proveedor); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Save(&proveedor).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar proveedor"})
		return
	}
	c.JSON(http.StatusOK, proveedor)
}

// DeleteProveedor elimina un proveedor
func DeleteProveedor(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var proveedor modelos.Proveedor
	if err := db.DB.First(&proveedor, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Proveedor no encontrado"})
		return
	}

	if err := db.DB.Delete(&proveedor).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar proveedor"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Proveedor eliminado exitosamente"})
}
