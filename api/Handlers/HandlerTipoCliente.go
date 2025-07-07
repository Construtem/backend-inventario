package Handlers

import (
	"backend-inventario/api/Controllers"
	modelos "backend-inventario/api/Models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetTipoClienteHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tipos, err := Controllers.GetTiposCliente(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener tipos de clientes", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, tipos)
	}
}

func GetTipoClienteByIDHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalido"})
			return
		}

		tipoCliente, err := Controllers.GetTipoClienteByID(db, uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Tipo de cliente no encontrado", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, tipoCliente)
	}
}

func CreateTipoClienteHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var nuevo modelos.TipoCliente
		if err := c.ShouldBindJSON(&nuevo); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos", "details": err.Error()})
			return
		}
		if err := Controllers.CreateTipoCliente(db, &nuevo); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear tipo de cliente", "details": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, nuevo)
	}
}

func UpdateTipoClienteHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalido"})
			return
		}

		var actualizado modelos.TipoCliente
		if err := c.ShouldBindJSON(&actualizado); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos", "details": err.Error()})
			return
		}

		tipo, err := Controllers.UpdateTipoCliente(db, uint(id), &actualizado)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar tipo de cliente", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, tipo)
	}
}

func DeleteTipoClienteHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalido"})
			return
		}

		if err := Controllers.DeleteTipoCliente(db, uint(id)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar tipo de cliente", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Tipo de cliente eliminado exitosamente"})
	}
}
