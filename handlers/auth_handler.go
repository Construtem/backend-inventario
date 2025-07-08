// Este archivo es usado para extraer datos con el token de Firebase para la autenticacion
package handlers

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"backend-inventario/services" // cambia según tu nombre real del módulo
)

func VerifyToken(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token no proporcionado"})
		return
	}

	idToken := strings.TrimPrefix(authHeader, "Bearer ")

	// Verificar token con Firebase
	token, err := services.FirebaseAuth.VerifyIDToken(context.Background(), idToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
		return
	}

	// Extraer información del token
	email := token.Claims["email"].(string)
	name := token.Claims["name"].(string)
	uid := token.Claims["user_id"].(string)
	picture := token.Claims["picture"].(string)


	c.JSON(http.StatusOK, gin.H{
		"name":     name,
		"email":    email,
		"photoURL": picture,
		"uid":      uid,
	})
}
