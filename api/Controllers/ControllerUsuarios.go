package Controllers

import (
	"net/http"
	"strconv"

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
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var usuario modelos.Usuario
	// Preload la relación con Rol
	if err := db.DB.Preload("Rol").First(&usuario, id).Error; err != nil {
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

	// // Opcional: Hashear la contraseña antes de guardar
	// hashedPassword, err := bcrypt.GenerateFromPassword([]byte(usuario.Contrasena), bcrypt.DefaultCost)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al hashear contraseña"})
	// 	return
	// }
	// usuario.Contrasena = string(hashedPassword)

	if err := db.DB.Create(&usuario).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear usuario"})
		return
	}
	// Carga el rol para la respuesta
	db.DB.Preload("Rol").First(&usuario, usuario.Correo)
	c.JSON(http.StatusCreated, usuario)
}

// UpdateUsuario actualiza un usuario existente
func UpdateUsuario(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var usuario modelos.Usuario
	if err := db.DB.First(&usuario, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}

	// Guardar la contraseña actual antes de vincular los nuevos datos
	//	currentPassword := usuario.Contrasena

	if err := c.ShouldBindJSON(&usuario); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// // Opcional: Si la contraseña fue actualizada, hashearla
	// if usuario.Contrasena != "" && usuario.Contrasena != currentPassword {
	// 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(usuario.Contrasena), bcrypt.DefaultCost)
	// 	if err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al hashear contraseña"})
	// 		return
	// 	}
	// 	usuario.Contrasena = string(hashedPassword)
	// } else {
	// 	usuario.Contrasena = currentPassword // Mantener la contraseña si no se envió una nueva
	// }

	if err := db.DB.Save(&usuario).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar usuario"})
		return
	}
	// Carga el rol para la respuesta
	db.DB.Preload("Rol").First(&usuario, usuario.Correo)
	c.JSON(http.StatusOK, usuario)
}

// DeleteUsuario elimina un usuario
func DeleteUsuario(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var usuario modelos.Usuario
	if err := db.DB.First(&usuario, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}

	if err := db.DB.Delete(&usuario).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar usuario"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Usuario eliminado exitosamente"})
}
