package Handlers

import (
	"backend-inventario/api/Controllers"
	modelos "backend-inventario/api/Models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetProductosDespachohHandler obtiene todos los productos de despacho
func GetProductosDespachoHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		productosDespacho, err := Controllers.GetProductosDespacho(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener productos de despacho: " + err.Error()})
			return
		}
		c.JSON(http.StatusOK, productosDespacho)
	}
}

// GetProductosDespachoDetalladoHandler obtiene todos los productos de despacho con información detallada
func GetProductosDespachoDetalladoHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		productosDespacho, err := Controllers.GetProductosDespachoDetallado(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener productos de despacho detallados: " + err.Error()})
			return
		}
		c.JSON(http.StatusOK, productosDespacho)
	}
}

// GetProductosDespachoByDespachoIDHandler obtiene todos los productos de un despacho específico
func GetProductosDespachoByDespachoIDHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		despachoIDStr := c.Param("despacho_id")
		despachoID, err := strconv.ParseUint(despachoIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID de despacho inválido"})
			return
		}

		productosDespacho, err := Controllers.GetProductosDespachoByDespachoID(db, uint(despachoID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener productos del despacho: " + err.Error()})
			return
		}
		c.JSON(http.StatusOK, productosDespacho)
	}
}

// GetProductosDespachoDetalladoByDespachoIDHandler obtiene productos de un despacho con información detallada
func GetProductosDespachoDetalladoByDespachoIDHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		despachoIDStr := c.Param("despacho_id")
		despachoID, err := strconv.ParseUint(despachoIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID de despacho inválido"})
			return
		}

		productosDespacho, err := Controllers.GetProductosDespachoDetalladoByDespachoID(db, uint(despachoID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener productos detallados del despacho: " + err.Error()})
			return
		}
		c.JSON(http.StatusOK, productosDespacho)
	}
}

// GetProductoDespachoByIDHandler obtiene un producto de despacho específico
func GetProductoDespachoByIDHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		despachoIDStr := c.Param("despacho_id")
		productoID := c.Param("producto_id")

		despachoID, err := strconv.ParseUint(despachoIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID de despacho inválido"})
			return
		}

		productoDespacho, err := Controllers.GetProductoDespacho(db, uint(despachoID), productoID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "Producto de despacho no encontrado"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener producto de despacho: " + err.Error()})
			return
		}
		c.JSON(http.StatusOK, productoDespacho)
	}
}

// CreateProductoDespachoHandler crea un nuevo producto de despacho
func CreateProductoDespachoHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var productoDespacho modelos.ProductosDespacho
		if err := c.ShouldBindJSON(&productoDespacho); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: " + err.Error()})
			return
		}

		createdProductoDespacho, err := Controllers.CreateProductoDespacho(db, productoDespacho)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear producto de despacho: " + err.Error()})
			return
		}
		c.JSON(http.StatusCreated, createdProductoDespacho)
	}
}

// UpdateProductoDespachoHandler actualiza un producto de despacho existente
func UpdateProductoDespachoHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		despachoIDStr := c.Param("despacho_id")
		productoID := c.Param("producto_id")

		despachoID, err := strconv.ParseUint(despachoIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID de despacho inválido"})
			return
		}

		var productoDespacho modelos.ProductosDespacho
		if err := c.ShouldBindJSON(&productoDespacho); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: " + err.Error()})
			return
		}

		updatedProductoDespacho, err := Controllers.UpdateProductoDespacho(db, uint(despachoID), productoID, productoDespacho)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar producto de despacho: " + err.Error()})
			return
		}
		c.JSON(http.StatusOK, updatedProductoDespacho)
	}
}

// DeleteProductoDespachoHandler elimina un producto de despacho
func DeleteProductoDespachoHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		despachoIDStr := c.Param("despacho_id")
		productoID := c.Param("producto_id")

		despachoID, err := strconv.ParseUint(despachoIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID de despacho inválido"})
			return
		}

		err = Controllers.DeleteProductoDespacho(db, uint(despachoID), productoID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar producto de despacho: " + err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Producto de despacho eliminado exitosamente"})
	}
}
