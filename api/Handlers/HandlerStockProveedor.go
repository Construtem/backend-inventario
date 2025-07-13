package Handlers

import (
	"backend-inventario/api/Controllers"
	modelos "backend-inventario/api/Models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetStockProveedorHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		stockProveedores, err := Controllers.GetStockProveedor(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener stock de proveedores", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, stockProveedores)
	}
}

func GetStockProveedorByIDHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		proveedorID, _ := strconv.Atoi(c.Param("proveedor_id"))
		sku := c.Param("sku")
		stock, err := Controllers.GetStockProveedorByID(db, uint(proveedorID), sku)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Stock no encontrado"})
			return
		}
		c.JSON(http.StatusOK, stock)
	}
}

func CreateStockProveedorHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var nuevo modelos.StockProveedor
		if err := c.ShouldBindJSON(&nuevo); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos", "details": err.Error()})
			return
		}
		if err := Controllers.CreateStockProveedor(db, &nuevo); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear stock de proveedor", "details": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, nuevo)
	}
}

func UpdateStockProveedorHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		proveedorID, _ := strconv.Atoi(c.Param("proveedor_id"))
		sku := c.Param("sku")
		var stock modelos.StockProveedor
		if err := c.ShouldBindJSON(&stock); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos", "details": err.Error()})
			return
		}
		if err := Controllers.UpdateStockProveedor(db, uint(proveedorID), sku, stock); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar stock de proveedor", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Stock actualizado correctamente"})
	}
}

func DeleteStockProveedorHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		proveedorID, _ := strconv.Atoi(c.Param("proveedor_id"))
		sku := c.Param("sku")
		if err := Controllers.DeleteStockProveedor(db, uint(proveedorID), sku); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar stock de proveedor", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Stock de proveedor eliminado correctamente"})
	}
}
