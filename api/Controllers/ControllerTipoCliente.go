package Controllers

import (
	modelos "backend-inventario/api/Models"
	"errors"

	"gorm.io/gorm"
)

func GetTiposCliente(db *gorm.DB) ([]modelos.TipoCliente, error) {
	var tipos []modelos.TipoCliente
	if err := db.Find(&tipos).Error; err != nil {
		return nil, err
	}
	return tipos, nil
}

func GetTipoClienteByID(db *gorm.DB, id uint) (*modelos.TipoCliente, error) {
	var tipo modelos.TipoCliente
	if err := db.First(&tipo, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &tipo, nil
}

func CreateTipoCliente(db *gorm.DB, nuevo *modelos.TipoCliente) error {
	if nuevo.Nombre == "" {
		return errors.New("el nombre del tipo de cliente no puede estar vacío")
	}
	return db.Create(nuevo).Error
}

func UpdateTipoCliente(db *gorm.DB, id uint, nuevo *modelos.TipoCliente) (*modelos.TipoCliente, error) {
	var existente modelos.TipoCliente
	if err := db.First(&existente, id).Error; err != nil {
		return nil, errors.New("tipo de cliente no encontrado")
	}

	if nuevo.Nombre == "" {
		return nil, errors.New("el nombre del tipo de cliente no puede estar vacío")
	}

	err := db.Model(&existente).Updates(modelos.TipoCliente{
		Nombre: nuevo.Nombre,
	}).Error
	if err != nil {
		return nil, err
	}
	return &existente, nil
}

func DeleteTipoCliente(db *gorm.DB, id uint) error {
	return db.Delete(&modelos.TipoCliente{}, id).Error
}
