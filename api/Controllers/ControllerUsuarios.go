package Controllers

import (
	"net/http"

	modelos "backend-inventario/api/Models"
	"backend-inventario/api/db"

	"github.com/gin-gonic/gin"
)

// GetUsuarios obtiene todos los usuarios
func GetUsuarios(c *gin.Context) {
	var usuarios []modelos.Usuario
	// Preload la relación con Rol
	if err := db.DB.Preload("Rol").Find(&usuarios).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener usuarios"})
		return
	}
	c.JSON(http.StatusOK, usuarios)
}

// GetUsuarioByID obtiene un usuario por su ID
func GetUsuarioByID(c *gin.Context) {
	uid := c.Param("id")

	var usuario modelos.Usuario
	// Preload la relación con Rol
	if err := db.DB.Preload("Rol").First(&usuario, "uid = ?", uid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}
	c.JSON(http.StatusOK, usuario)
}

// CreateUsuario crea un nuevo usuario
func CreateUsuario(c *gin.Context) {
	var usuario modelos.Usuario
	if err := c.ShouldBindJSON(&usuario); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Create(&usuario).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear usuario"})
		return
	}
	// Carga el rol para la respuesta
	db.DB.Preload("Rol").First(&usuario, "uid = ?", usuario.UID)
	c.JSON(http.StatusCreated, usuario)
}

// UpdateUsuario actualiza un usuario existente
func UpdateUsuario(c *gin.Context) {
	uid := c.Param("id")

	var usuario modelos.Usuario
	if err := db.DB.First(&usuario, "uid = ?", uid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}

	if err := c.ShouldBindJSON(&usuario); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Save(&usuario).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar usuario"})
		return
	}
	// Carga el rol para la respuesta
	db.DB.Preload("Rol").First(&usuario, "uid = ?", uid)
	c.JSON(http.StatusOK, usuario)
}

// DeleteUsuario elimina un usuario
func DeleteUsuario(c *gin.Context) {
	uid := c.Param("id")

	var usuario modelos.Usuario
	if err := db.DB.First(&usuario, "uid = ?", uid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}

	if err := db.DB.Delete(&usuario).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar usuario"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Usuario eliminado exitosamente"})
}
