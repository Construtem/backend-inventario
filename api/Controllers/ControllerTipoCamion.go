package Controllers

import (
	modelos "backend-inventario/api/Models"
	"errors"

	"gorm.io/gorm"
)

func GetTiposCamion(db *gorm.DB) ([]modelos.TipoCamion, error) {
	var tipos []modelos.TipoCamion
	if err := db.Find(&tipos).Error; err != nil {
		return nil, err
	}
	return tipos, nil
}

func GetTipoCamionByID(db *gorm.DB, id uint) (*modelos.TipoCamion, error) {
	var tipo modelos.TipoCamion
	if err := db.First(&tipo, id).Error; err != nil {
		return nil, err
	}
	return &tipo, nil
}

func CreateTipoCamion(db *gorm.DB, nuevo *modelos.TipoCamion) error {
	if nuevo.Volumen <= 0 || nuevo.PesoMaximo <= 0 {
		return errors.New("volumen y peso máximo deben ser mayores a cero")
	}
	return db.Create(nuevo).Error
}

func UpdateTipoCamion(db *gorm.DB, id uint, nuevo *modelos.TipoCamion) (*modelos.TipoCamion, error) {
	var existente modelos.TipoCamion
	if err := db.First(&existente, id).Error; err != nil {
		return nil, errors.New("tipo de camión no encontrado")
	}
	if nuevo.Volumen <= 0 || nuevo.PesoMaximo <= 0 {
		return nil, errors.New("volumen y peso máximo deben ser mayores a cero")
	}
	err := db.Model(&existente).Updates(modelos.TipoCamion{
		Volumen:    nuevo.Volumen,
		PesoMaximo: nuevo.PesoMaximo,
	}).Error
	if err != nil {
		return nil, err
	}
	return &existente, nil
}

func DeleteTipoCamion(db *gorm.DB, id uint) error {
	return db.Delete(&modelos.TipoCamion{}, id).Error
}
