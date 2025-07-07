package Handlers

import (
	"backend-inventario/api/Controllers"
	modelos "backend-inventario/api/Models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetClientesHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientes, err := Controllers.GetClientes(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener clientes", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, clientes)
	}
}

func GetClienteByIDHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalido"})
			return
		}

		cliente, err := Controllers.GetClienteByID(db, uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Cliente no encontrado", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, cliente)
	}
}

func CreateClienteHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var nuevo modelos.Cliente
		if err := c.ShouldBindJSON(&nuevo); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos", "details": err.Error()})
			return
		}
		if err := Controllers.CreateCliente(db, &nuevo); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear cliente", "details": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, nuevo)
	}
}

func UpdateClienteHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalido"})
			return
		}

		var actualizado modelos.Cliente
		if err := c.ShouldBindJSON(&actualizado); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos", "details": err.Error()})
			return
		}

		cliente, err := Controllers.UpdateCliente(db, uint(id), &actualizado)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar cliente", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, cliente)
	}
}

func DeleteClienteHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalido"})
			return
		}

		if err := Controllers.DeleteCliente(db, uint(id)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar cliente", "details": err.Error()})
			return
		}
		c.JSON(http.StatusNoContent, gin.H{"message": "Cliente eliminado exitosamente"})
	}
}
