package Controllers

import (
	"net/http"
	"strconv" // Necesario para ParseUint

	modelos "backend-inventario/api/Models"
	"backend-inventario/api/db"

	"github.com/gin-gonic/gin"
)

// GetClientes obtiene todos los clientes
func GetClientes(c *gin.Context) {
	var clientes []modelos.Cliente
	// preload para cargar las relaciones si las tuvieran, Cliente no tiene en este momento.
	if err := db.DB.Find(&clientes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener clientes"})
		return
	}
	c.JSON(http.StatusOK, clientes)
}

// GetClienteByID obtiene un cliente por su ID
func GetClienteByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32) // Convertir string a uint
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var cliente modelos.Cliente
	if err := db.DB.First(&cliente, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cliente no encontrado"})
		return
	}
	c.JSON(http.StatusOK, cliente)
}

// CreateCliente crea un nuevo cliente
func CreateCliente(c *gin.Context) {
	var cliente modelos.Cliente
	if err := c.ShouldBindJSON(&cliente); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := db.DB.Create(&cliente).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear cliente"})
		return
	}
	c.JSON(http.StatusCreated, cliente)
}

// UpdateCliente actualiza un cliente existente
func UpdateCliente(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var cliente modelos.Cliente
	if err := db.DB.First(&cliente, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cliente no encontrado"})
		return
	}

	// Bind con los datos JSON actualizados
	if err := c.ShouldBindJSON(&cliente); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Save(&cliente).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar cliente"})
		return
	}
	c.JSON(http.StatusOK, cliente)
}

// DeleteCliente elimina un cliente
func DeleteCliente(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var cliente modelos.Cliente
	if err := db.DB.First(&cliente, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cliente no encontrado"})
		return
	}

	if err := db.DB.Delete(&cliente).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar cliente"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Cliente eliminado exitosamente"})
}
