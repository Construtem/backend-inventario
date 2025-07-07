package Handlers

import (
	"backend-inventario/api/Controllers"
	modelos "backend-inventario/api/Models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetUsuariosHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		usuarios, err := Controllers.GetUsuarios(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener usuarios", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, usuarios)
	}
}

func GetUsuarioByEmailHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		email := c.Param("email")
		usuario, err := Controllers.GetUsuarioByEmail(db, email)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, usuario)
	}
}

func CreateUsuarioHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var nuevo modelos.Usuario
		if err := c.ShouldBindJSON(&nuevo); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos", "details": err.Error()})
			return
		}
		if err := Controllers.CreateUsuario(db, &nuevo); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear usuario", "details": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, nuevo)
	}
}

func UpdateUsuarioHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		email := c.Param("email")
		var actualizado modelos.Usuario
		if err := c.ShouldBindJSON(&actualizado); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos", "details": err.Error()})
			return
		}
		usuario, err := Controllers.UpdateUsuario(db, email, &actualizado)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar usuario", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, usuario)
	}
}

func DeleteUsuarioHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		email := c.Param("email")
		if err := Controllers.DeleteUsuario(db, email); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar usuario", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Usuario eliminado exitosamente"})
	}
}
