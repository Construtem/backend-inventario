package Handlers

import (
	"backend-inventario/api/Controllers"
	modelos "backend-inventario/api/Models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetSucursalesHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		sucursales, err := Controllers.GetSucursales(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener sucursales", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, sucursales)
	}
}

func GetSucursalByIDHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalido"})
			return
		}

		sucursal, err := Controllers.GetSucursalByID(db, uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Sucursal no encontrada", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, sucursal)
	}
}

func CreateSucursalHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var nueva modelos.Sucursal
		if err := c.ShouldBindJSON(&nueva); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos", "details": err.Error()})
			return
		}
		if err := Controllers.CreateSucursal(db, &nueva); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear sucursal", "details": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, nueva)
	}
}

func UpdateSucursalHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalido"})
			return
		}

		var actualizada modelos.Sucursal
		if err := c.ShouldBindJSON(&actualizada); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos", "details": err.Error()})
			return
		}

		sucursal, err := Controllers.UpdateSucursal(db, uint(id), &actualizada)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar sucursal", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, sucursal)
	}
}

func DeleteSucursalHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalido"})
			return
		}

		if err := Controllers.DeleteSucursal(db, uint(id)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar sucursal", "details": err.Error()})
			return
		}
		c.JSON(http.StatusNoContent, nil)
	}
}
