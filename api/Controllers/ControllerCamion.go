package Controllers

import (
	modelos "backend-inventario/api/Models"
	"errors"

	"gorm.io/gorm"
)

func GetCamiones(db *gorm.DB) ([]modelos.Camion, error) {
	var camiones []modelos.Camion
	if err := db.Preload("Tipo").Find(&camiones).Error; err != nil {
		return nil, err
	}
	return camiones, nil
}

func GetCamionByID(db *gorm.DB, id uint) (*modelos.Camion, error) {
	var camion modelos.Camion
	if err := db.Preload("Tipo").First(&camion, id).Error; err != nil {
		return nil, err
	}
	return &camion, nil
}

func CreateCamion(db *gorm.DB, nuevo *modelos.Camion) error {
	if nuevo.Patente == "" {
		return errors.New("la patente no puede estar vacía")
	}
	if nuevo.TipoID == 0 {
		return errors.New("el tipo de camión es obligatorio")
	}
	return db.Create(nuevo).Error
}

func UpdateCamion(db *gorm.DB, id uint, actualizado *modelos.Camion) (*modelos.Camion, error) {
	var existente modelos.Camion
	if err := db.First(&existente, id).Error; err != nil {
		return nil, errors.New("camión no encontrado")
	}

	if actualizado.Patente == "" {
		return nil, errors.New("la patente no puede estar vacía")
	}
	if actualizado.TipoID == 0 {
		return nil, errors.New("el tipo de camión es obligatorio")
	}

	err := db.Model(&existente).Updates(modelos.Camion{
		Patente: actualizado.Patente,
		TipoID:  actualizado.TipoID,
		Activo:  actualizado.Activo,
	}).Error
	if err != nil {
		return nil, err
	}

	return &existente, nil
}

func DeleteCamion(db *gorm.DB, id uint) error {
	return db.Delete(&modelos.Camion{}, id).Error
}
