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
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Hubo un problema al obtener los clientes.",
				"details": err.Error(),
			})
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
			c.JSON(http.StatusBadRequest, gin.H{"error": "El ID proporcionado no es válido."})
			return
		}

		cliente, err := Controllers.GetClienteByID(db, uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "No se encontró un cliente con el ID especificado."})
			return
		}
		c.JSON(http.StatusOK, cliente)
	}
}

func CreateClienteHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var nuevo modelos.Cliente
		if err := c.ShouldBindJSON(&nuevo); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Los datos del cliente enviados no son válidos.",
				"details": err.Error(),
			})
			return
		}
		if err := Controllers.CreateCliente(db, &nuevo); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "No se pudo registrar el cliente.",
				"details": err.Error(),
			})
			return
		}
		c.JSON(http.StatusCreated, gin.H{
			"message": "Cliente creado exitosamente.",
			"cliente": nuevo,
		})
	}
}

func UpdateClienteHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "El ID proporcionado no es válido."})
			return
		}

		var actualizado modelos.Cliente
		if err := c.ShouldBindJSON(&actualizado); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Los datos enviados para actualizar el cliente no son válidos.",
				"details": err.Error(),
			})
			return
		}

		cliente, err := Controllers.UpdateCliente(db, uint(id), &actualizado)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "No se pudo actualizar el cliente.",
				"details": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "Cliente actualizado exitosamente.",
			"cliente": cliente,
		})
	}
}

func DeleteClienteHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "El ID proporcionado no es válido."})
			return
		}

		if err := Controllers.DeleteCliente(db, uint(id)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "No se pudo eliminar el cliente.",
				"details": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "Cliente eliminado exitosamente.",
		})
	}
}
