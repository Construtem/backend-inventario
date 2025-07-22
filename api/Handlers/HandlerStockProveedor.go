package Handlers

import (
	"backend-inventario/api/Controllers"
	modelos "backend-inventario/api/Models"
	"net/http"
	"strconv"
	"time"

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

		// Retornar el stock creado con preloads
		var stockCreado modelos.StockProveedor
		if err := db.Preload("Proveedor").Preload("Producto").
			First(&stockCreado, "proveedor_id = ? AND sku = ?", nuevo.ProveedorID, nuevo.ProductoID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener el stock creado"})
			return
		}

		c.JSON(http.StatusCreated, stockCreado)
	}
}

// StockProveedorCompleto para crear producto, proveedor y stock en una sola operación
type StockProveedorCompleto struct {
	ProveedorID  uint                     `json:"proveedor_id"`
	ProductoID   string                   `json:"sku"`
	Stock        int                      `json:"stock"`
	FechaIngreso time.Time                `json:"fecha_ingreso"`
	Producto     *modelos.Producto        `json:"producto,omitempty"`
	Proveedor    *modelos.Proveedor       `json:"proveedor,omitempty"`
}

func CreateStockProveedorCompletoHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var data StockProveedorCompleto
		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos", "details": err.Error()})
			return
		}

		// Transacción para asegurar consistencia
		tx := db.Begin()
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()

		// Crear o verificar proveedor
		if data.Proveedor != nil && data.Proveedor.ID == 0 {
			if err := tx.Create(data.Proveedor).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear proveedor", "details": err.Error()})
				return
			}
			data.ProveedorID = data.Proveedor.ID
		} else if data.ProveedorID > 0 {
			var proveedor modelos.Proveedor
			if err := tx.First(&proveedor, data.ProveedorID).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusBadRequest, gin.H{"error": "Proveedor no encontrado"})
				return
			}
		}

		// Crear o verificar producto
		if data.Producto != nil && data.Producto.SKU != "" {
			data.ProductoID = data.Producto.SKU
			var existingProducto modelos.Producto
			if err := tx.First(&existingProducto, "sku = ?", data.Producto.SKU).Error; err != nil {
				// El producto no existe, lo creamos
				data.Producto.ProveedorID = data.ProveedorID
				if err := tx.Create(data.Producto).Error; err != nil {
					tx.Rollback()
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear producto", "details": err.Error()})
					return
				}
			}
		} else if data.ProductoID != "" {
			var producto modelos.Producto
			if err := tx.First(&producto, "sku = ?", data.ProductoID).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusBadRequest, gin.H{"error": "Producto no encontrado"})
				return
			}
		}

		// Crear stock
		stockProveedor := modelos.StockProveedor{
			ProveedorID:  data.ProveedorID,
			ProductoID:   data.ProductoID,
			Stock:        data.Stock,
			FechaIngreso: data.FechaIngreso,
		}

		if stockProveedor.FechaIngreso.IsZero() {
			stockProveedor.FechaIngreso = time.Now()
		}

		// Verificar si ya existe
		var existingStock modelos.StockProveedor
		if err := tx.First(&existingStock, "proveedor_id = ? AND sku = ?", data.ProveedorID, data.ProductoID).Error; err == nil {
			tx.Rollback()
			c.JSON(http.StatusConflict, gin.H{"error": "Ya existe stock para este producto y proveedor"})
			return
		}

		if err := tx.Create(&stockProveedor).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear stock", "details": err.Error()})
			return
		}

		tx.Commit()

		// Retornar el stock creado con preloads
		var stockCreado modelos.StockProveedor
		if err := db.Preload("Proveedor").Preload("Producto").
			First(&stockCreado, "proveedor_id = ? AND sku = ?", data.ProveedorID, data.ProductoID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener el stock creado"})
			return
		}

		c.JSON(http.StatusCreated, stockCreado)
	}
}

func UpdateStockProveedorHandler(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        proveedorID, _ := strconv.Atoi(c.Param("proveedor_id"))
        sku := c.Param("sku") // <-- aquí el cambio
        var stock modelos.StockProveedor
        if err := c.ShouldBindJSON(&stock); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos", "details": err.Error()})
            return
        }
        if err := Controllers.UpdateStockProveedor(db, uint(proveedorID), sku, stock); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar stock de proveedor", "details": err.Error()})
            return
        }
        var updated modelos.StockProveedor
        if err := db.Preload("Producto").Preload("Proveedor").
            First(&updated, "proveedor_id = ? AND sku = ?", proveedorID, sku).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo obtener el stock actualizado"})
            return
        }
        c.JSON(http.StatusOK, updated)
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
