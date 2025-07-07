package Controllers

import (
	modelos "backend-inventario/api/Models"
	"errors"

	"gorm.io/gorm"
)

func GetRoles(db *gorm.DB) ([]modelos.Rol, error) {
	var roles []modelos.Rol
	if err := db.Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

func GetRolByID(db *gorm.DB, id uint) (*modelos.Rol, error) {
	var rol modelos.Rol
	if err := db.First(&rol, id).Error; err != nil {
		return nil, errors.New("rol no encontrado")
	}
	return &rol, nil
}

func CreateRol(db *gorm.DB, nuevo *modelos.Rol) error {
	if nuevo.Nombre == "" {
		return errors.New("el nombre del rol no puede estar vacío")
	}
	return db.Create(nuevo).Error
}

func UpdateRol(db *gorm.DB, id uint, actualizado *modelos.Rol) (*modelos.Rol, error) {
	var existente modelos.Rol
	if err := db.First(&existente, id).Error; err != nil {
		return nil, errors.New("rol no encontrado")
	}

	if actualizado.Nombre == "" {
		return nil, errors.New("el nombre del rol no puede estar vacío")
	}

	err := db.Model(&existente).Updates(modelos.Rol{
		Nombre: actualizado.Nombre,
	}).Error
	if err != nil {
		return nil, err
	}
	return &existente, nil
}

func DeleteRol(db *gorm.DB, id uint) error {
	if err := db.First(&modelos.Rol{}, id).Error; err != nil {
		return errors.New("rol no encontrado")
	}
	return db.Delete(&modelos.Rol{}, id).Error
}
