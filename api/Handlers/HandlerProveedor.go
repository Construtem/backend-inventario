package Handlers

import (
	"backend-inventario/api/Controllers"
	modelos "backend-inventario/api/Models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetProveedoresHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		proveedores, err := Controllers.GetProveedores(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener proveedores", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, proveedores)
	}
}

func GetProveedorByIDHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
			return
		}

		proveedor, err := Controllers.GetProveedorByID(db, uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Proveedor no encontrado", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, proveedor)
	}
}

func CreateProveedorHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var nuevo modelos.Proveedor
		if err := c.ShouldBindJSON(&nuevo); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos", "details": err.Error()})
			return
		}
		if err := Controllers.CreateProveedor(db, &nuevo); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear proveedor", "details": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, nuevo)
	}
}

func UpdateProveedorHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
			return
		}

		var actualizado modelos.Proveedor
		if err := c.ShouldBindJSON(&actualizado); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos", "details": err.Error()})
			return
		}

		proveedor, err := Controllers.UpdateProveedor(db, uint(id), &actualizado)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar proveedor", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, proveedor)
	}
}

func DeleteProveedorHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
			return
		}

		if err := Controllers.DeleteProveedor(db, uint(id)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar proveedor", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Proveedor eliminado exitosamente"})
	}
}
