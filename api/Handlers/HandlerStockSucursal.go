package Handlers

import (
	"backend-inventario/api/Controllers"
	modelos "backend-inventario/api/Models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetStockSucursalHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		stocks, err := Controllers.GetStockSucursal(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener stock por sucursal", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, stocks)
	}
}

func GetStockSucursalByIDHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		sku := c.Param("sku")
		sucursalIDStr := c.Param("sucursal_id")
		sucursalID, err := strconv.ParseUint(sucursalIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalido"})
			return
		}

		stock, err := Controllers.GetStockSucursalByID(db, uint(sucursalID), sku)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Registro de stock no encontrado", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, stock)
	}
}

func CreateStockSucursalHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var nuevo modelos.StockSucursal
		if err := c.ShouldBindJSON(&nuevo); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := Controllers.CreateStockSucursal(db, &nuevo); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear registro de stock", "details": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, nuevo)
	}
}

func UpdateStockSucursalHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		sku := c.Param("sku")
		sucursalIDStr := c.Param("sucursal_id")
		sucursalID, err := strconv.ParseUint(sucursalIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalido"})
			return
		}

		var actualizado modelos.StockSucursal
		if err := c.ShouldBindJSON(&actualizado); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos", "details": err.Error()})
			return
		}

		if err := Controllers.UpdateStockSucursal(db, sku, uint(sucursalID), &actualizado); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar registro de stock", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Registro de stock actualizado correctamente"})
	}
}

func DeleteStockSucursalHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		sku := c.Param("sku")
		sucursalIDStr := c.Param("sucursal_id")
		sucursalID, err := strconv.ParseUint(sucursalIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalido"})
			return
		}

		if err := Controllers.DeleteStockSucursal(db, sku, uint(sucursalID)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar registro de stock", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Registro de stock eliminado correctamente"})
	}
}
