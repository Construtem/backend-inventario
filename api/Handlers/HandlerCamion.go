package Handlers

import (
	"backend-inventario/api/Controllers"
	modelos "backend-inventario/api/Models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetCamionesHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		camiones, err := Controllers.GetCamiones(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Ocurrió un error al obtener los camiones",
				"details": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, camiones)
	}
}

func GetCamionByIDHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "El ID del camión no es válido"})
			return
		}

		camion, err := Controllers.GetCamionByID(db, uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Camión no encontrado"})
			return
		}
		c.JSON(http.StatusOK, camion)
	}
}

func CreateCamionHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var nuevo modelos.Camion
		if err := c.ShouldBindJSON(&nuevo); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Los datos enviados son inválidos",
				"details": err.Error(),
			})
			return
		}
		if err := Controllers.CreateCamion(db, &nuevo); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{ // <-- cambiaste de 500 a 400 si es validación
				"error":   "No se pudo registrar el camión",
				"details": err.Error(),
			})
			return
		}
		c.JSON(http.StatusCreated, nuevo)
	}
}

func UpdateCamionHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "El ID del camión no es válido"})
			return
		}

		var actualizado modelos.Camion
		if err := c.ShouldBindJSON(&actualizado); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Los datos enviados son inválidos",
				"details": err.Error(),
			})
			return
		}

		camion, err := Controllers.UpdateCamion(db, uint(id), &actualizado)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "No se pudo actualizar el camión",
				"details": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, camion)
	}
}

func DeleteCamionHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "El ID del camión no es válido"})
			return
		}

		if err := Controllers.DeleteCamion(db, uint(id)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "No se pudo eliminar el camión",
				"details": err.Error(),
			})
			return
		}
		c.JSON(http.StatusNoContent, nil)
	}
}
