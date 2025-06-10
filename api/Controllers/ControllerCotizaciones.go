package Controllers

import (
	"net/http"
	"strconv"
	"time"

	modelos "backend-inventario/api/Models"
	"backend-inventario/api/db"

	"github.com/gin-gonic/gin"
)

// GetCotizaciones obtiene todas las cotizaciones
func GetCotizaciones(c *gin.Context) {
	var cotizaciones []modelos.Cotizacion
	// Preload las relaciones relevantes
	if err := db.DB.Preload("Cliente").Preload("Vendedor").Preload("Ubicacion").Find(&cotizaciones).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener cotizaciones"})
		return
	}
	c.JSON(http.StatusOK, cotizaciones)
}

// GetCotizacionByID obtiene una cotizacion por su ID
func GetCotizacionByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var cotizacion modelos.Cotizacion
	if err := db.DB.Preload("Cliente").Preload("Vendedor").Preload("Ubicacion").First(&cotizacion, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cotización no encontrada"})
		return
	}
	c.JSON(http.StatusOK, cotizacion)
}

// CreateCotizacion crea una nueva cotizacion
func CreateCotizacion(c *gin.Context) {
	var cotizacion modelos.Cotizacion
	if err := c.ShouldBindJSON(&cotizacion); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Establecer la fecha actual si no se proporciona
	if cotizacion.Fecha.IsZero() {
		cotizacion.Fecha = time.Now()
	}

	if err := db.DB.Create(&cotizacion).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear cotización"})
		return
	}
	// Carga las relaciones para la respuesta
	db.DB.Preload("Cliente").Preload("Vendedor").Preload("Ubicacion").First(&cotizacion, cotizacion.ID)
	c.JSON(http.StatusCreated, cotizacion)
}

// UpdateCotizacion actualiza una cotizacion existente
func UpdateCotizacion(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var cotizacion modelos.Cotizacion
	if err := db.DB.First(&cotizacion, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cotización no encontrada"})
		return
	}

	if err := c.ShouldBindJSON(&cotizacion); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Save(&cotizacion).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar cotización"})
		return
	}
	// Carga las relaciones para la respuesta
	db.DB.Preload("Cliente").Preload("Vendedor").Preload("Ubicacion").First(&cotizacion, cotizacion.ID)
	c.JSON(http.StatusOK, cotizacion)
}

// DeleteCotizacion elimina una cotizacion
func DeleteCotizacion(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var cotizacion modelos.Cotizacion
	if err := db.DB.First(&cotizacion, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cotización no encontrada"})
		return
	}

	if err := db.DB.Delete(&cotizacion).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar cotización"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Cotización eliminada exitosamente"})
}
