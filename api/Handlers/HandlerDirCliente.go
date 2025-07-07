package Handlers

import (
	"backend-inventario/api/Controllers"
	modelos "backend-inventario/api/Models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetDirClientesHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		dirClientes, err := Controllers.GetDirCliente(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener direcciones de clientes", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, dirClientes)
	}
}

func GetDirClienteByIDHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalido"})
			return
		}

		dirCliente, err := Controllers.GetDirClienteByID(db, uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Dirección de cliente no encontrada", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, dirCliente)
	}
}

func CreateDirClienteHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var nueva modelos.DirCliente
		if err := c.ShouldBindJSON(&nueva); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos", "details": err.Error()})
			return
		}
		if err := Controllers.CreateDirCliente(db, &nueva); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear dirección de cliente", "details": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, nueva)
	}
}

func UpdateDirClienteHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalido"})
			return
		}

		var actualizada modelos.DirCliente
		if err := c.ShouldBindJSON(&actualizada); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos", "details": err.Error()})
			return
		}

		dir, err := Controllers.UpdateDirCliente(db, uint(id), &actualizada)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar dirección de cliente", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, dir)
	}
}

func DeleteDirClienteHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalido"})
			return
		}

		if err := Controllers.DeleteDirCliente(db, uint(id)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar dirección de cliente", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Dirección de cliente eliminada exitosamente"})
	}
}
