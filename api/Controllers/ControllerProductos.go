package Controllers

import (
	"net/http"
	"strconv"

	modelos "backend-inventario/api/Models"
	"backend-inventario/api/db"

	"github.com/gin-gonic/gin"
)

// GetProductos obtiene todos los productos
func GetProductos(c *gin.Context) {
	var productos []modelos.Producto
	// Preload la relación con Categoria para obtener los datos de la categoría
	if err := db.DB.Preload("Categoria").Find(&productos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener productos"})
		return
	}
	c.JSON(http.StatusOK, productos)
}

// GetProductoByID obtiene un producto por su ID
func GetProductoByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var producto modelos.Producto
	// Preload la relación con Categoria
	if err := db.DB.Preload("Categoria").First(&producto, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Producto no encontrado"})
		return
	}
	c.JSON(http.StatusOK, producto)
}

// CreateProducto crea un nuevo producto
func CreateProducto(c *gin.Context) {
	var producto modelos.Producto
	if err := c.ShouldBindJSON(&producto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Asegúrate de que la CategoriaID exista si no se maneja la creación de categorías anidadas
	if producto.CategoriaID == 0 { // O alguna validación para asegurar que se proporciona
		c.JSON(http.StatusBadRequest, gin.H{"error": "CategoriaID es requerida"})
		return
	}

	if err := db.DB.Create(&producto).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear producto"})
		return
	}
	// Carga la categoría para la respuesta, si es necesario
	db.DB.Preload("Categoria").First(&producto, producto.ID)
	c.JSON(http.StatusCreated, producto)
}

// UpdateProducto actualiza un producto existente
func UpdateProducto(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var producto modelos.Producto
	if err := db.DB.First(&producto, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Producto no encontrado"})
		return
	}

	if err := c.ShouldBindJSON(&producto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Save(&producto).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar producto"})
		return
	}
	// Carga la categoría para la respuesta
	db.DB.Preload("Categoria").First(&producto, producto.ID)
	c.JSON(http.StatusOK, producto)
}

// DeleteProducto elimina un producto
func DeleteProducto(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var producto modelos.Producto
	if err := db.DB.First(&producto, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Producto no encontrado"})
		return
	}

	if err := db.DB.Delete(&producto).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar producto"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Producto eliminado exitosamente"})
}
