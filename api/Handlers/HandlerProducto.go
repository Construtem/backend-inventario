package Handlers

import (
	"backend-inventario/api/Controllers"
	modelos "backend-inventario/api/Models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetProductosHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		productos, err := Controllers.GetProductos(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener productos", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, productos)
	}
}

func GetProductoBySKUHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		sku := c.Param("sku")
		producto, err := Controllers.GetProductoBySKU(db, sku)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Producto no encontrado"})
			return
		}
		c.JSON(http.StatusOK, producto)
	}
}

func CreateProductoHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var nuevo modelos.Producto
		if err := c.ShouldBindJSON(&nuevo); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
			return
		}
		if err := Controllers.CreateProducto(db, &nuevo); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear producto", "details": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, nuevo)
	}
}

func UpdateProductoHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		sku := c.Param("sku")
		var actualizado modelos.Producto
		if err := c.ShouldBindJSON(&actualizado); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos", "details": err.Error()})
			return
		}
		producto, err := Controllers.UpdateProducto(db, sku, &actualizado)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, producto)
	}
}

func DeleteProductoHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		sku := c.Param("sku")
		if err := Controllers.DeleteProducto(db, sku); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar producto", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Producto eliminado exitosamente"})
	}
}
