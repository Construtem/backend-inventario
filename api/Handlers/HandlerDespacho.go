package Handlers

import (
	"backend-inventario/api/Controllers"
	modelos "backend-inventario/api/Models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DespachoRequest struct {
	Despacho  modelos.Despacho            `json:"valor_despacho"`
	Productos []modelos.ProductosDespacho `json:"productos"`
}

func CreateDespachoHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request DespachoRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Formato inválido", "details": err.Error()})
			return
		}

		if err := Controllers.CreateDespacho(db, &request.Despacho, request.Productos); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear despacho", "details": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Despacho creado exitosamente", "despacho_id": request.Despacho.ID})
	}
}

func GetDespachosHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		despachos, err := Controllers.GetDespachos(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Error al obtener despachos",
				"details": err.Error(),
				"message": "Error interno del servidor al consultar despachos",
			})
			return
		}
		c.JSON(http.StatusOK, despachos)
	}
}

func GetDespachoByIDHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
			return
		}

		despacho, err := Controllers.GetDespachoByID(db, uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Despacho no encontrado", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, despacho)
	}
}

func UpdateDespachoHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
			return
		}

		var actualizado modelos.Despacho
		if err := c.ShouldBindJSON(&actualizado); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos", "details": err.Error()})
			return
		}

		if err := Controllers.UpdateDespacho(db, uint(id), &actualizado); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar despacho", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Despacho actualizado exitosamente"})
	}
}

func DeleteDespachoHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
			return
		}

		if err := Controllers.DeleteDespacho(db, uint(id)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar despacho", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Despacho eliminado exitosamente"})
	}
}
