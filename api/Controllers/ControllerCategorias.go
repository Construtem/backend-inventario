package Controllers

import (
	"net/http"
	"strconv"

	modelos "backend-inventario/api/Models"
	"backend-inventario/api/db"

	"github.com/gin-gonic/gin"
)

// GetCategorias obtiene todas las categorías
func GetCategorias(c *gin.Context) {
	var categorias []modelos.Categoria
	if err := db.DB.Find(&categorias).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener categorías"})
		return
	}
	c.JSON(http.StatusOK, categorias)
}

// GetCategoriaByID obtiene una categoría por su ID
func GetCategoriaByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var categoria modelos.Categoria
	if err := db.DB.First(&categoria, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Categoría no encontrada"})
		return
	}
	c.JSON(http.StatusOK, categoria)
}

// CreateCategoria crea una nueva categoría
func CreateCategoria(c *gin.Context) {
	var categoria modelos.Categoria
	if err := c.ShouldBindJSON(&categoria); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := db.DB.Create(&categoria).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear categoría"})
		return
	}
	c.JSON(http.StatusCreated, categoria)
}

// UpdateCategoria actualiza una categoría existente
func UpdateCategoria(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var categoria modelos.Categoria
	if err := db.DB.First(&categoria, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Categoría no encontrada"})
		return
	}

	if err := c.ShouldBindJSON(&categoria); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Save(&categoria).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar categoría"})
		return
	}
	c.JSON(http.StatusOK, categoria)
}

// DeleteCategoria elimina una categoría
func DeleteCategoria(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var categoria modelos.Categoria
	if err := db.DB.First(&categoria, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Categoría no encontrada"})
		return
	}

	if err := db.DB.Delete(&categoria).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar categoría"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Categoría eliminada exitosamente"})
}
