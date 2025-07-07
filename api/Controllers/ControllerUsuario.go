package Controllers

import (
	modelos "backend-inventario/api/Models"
	"errors"

	"gorm.io/gorm"
)

func GetUsuarios(db *gorm.DB) ([]modelos.Usuario, error) {
	var usuarios []modelos.Usuario
	if err := db.Preload("Rol").Find(&usuarios).Error; err != nil {
		return nil, err
	}
	return usuarios, nil
}

func GetUsuarioByEmail(db *gorm.DB, email string) (*modelos.Usuario, error) {
	var usuario modelos.Usuario
	if err := db.Preload("Rol").First(&usuario, "email = ?", email).Error; err != nil {
		return nil, errors.New("usuario no encontrado")
	}
	return &usuario, nil
}

func CreateUsuario(db *gorm.DB, nuevo *modelos.Usuario) error {
	if nuevo.Email == "" {
		return errors.New("el email del usuario no puede estar vacío")
	}
	if nuevo.Nombre == "" {
		return errors.New("el nombre del usuario no puede estar vacío")
	}
	return db.Create(nuevo).Error
}

func UpdateUsuario(db *gorm.DB, email string, actualizado *modelos.Usuario) (*modelos.Usuario, error) {
	var existente modelos.Usuario
	if err := db.First(&existente, "email = ?", email).Error; err != nil {
		return nil, errors.New("usuario no encontrado")
	}

	err := db.Model(&existente).Updates(modelos.Usuario{
		Nombre: actualizado.Nombre,
		RolID:  actualizado.RolID,
	}).Error
	if err != nil {
		return nil, err
	}
	return &existente, nil
}

func DeleteUsuario(db *gorm.DB, email string) error {
	return db.Delete(&modelos.Usuario{}, "email = ?", email).Error
}
