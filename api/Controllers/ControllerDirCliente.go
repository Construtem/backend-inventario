package Controllers

import (
	modelos "backend-inventario/api/Models"
	"errors"

	"gorm.io/gorm"
)

func GetDirCliente(db *gorm.DB) ([]modelos.DirCliente, error) {
	var direcciones []modelos.DirCliente
	if err := db.Preload("Cliente.Tipo").Find(&direcciones).Error; err != nil {
		return nil, err
	}
	return direcciones, nil
}

func GetDirClienteByID(db *gorm.DB, id uint) (*modelos.DirCliente, error) {
	var direccion modelos.DirCliente
	if err := db.Preload("Cliente.Tipo").First(&direccion, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &direccion, nil
}

func CreateDirCliente(db *gorm.DB, nueva *modelos.DirCliente) error {
	if nueva.RutCliente == "" {
		return errors.New("el cliente es obligatorio")
	}
	if nueva.Direccion == "" {
		return errors.New("la dirección no puede estar vacía")
	}
	return db.Create(nueva).Error
}

func UpdateDirCliente(db *gorm.DB, id uint, actualizada *modelos.DirCliente) (*modelos.DirCliente, error) {
	var existente modelos.DirCliente
	if err := db.First(&existente, id).Error; err != nil {
		return nil, errors.New("dirección no encontrada")
	}

	if actualizada.RutCliente == "" {
		return nil, errors.New("el cliente es obligatorio")
	}
	if actualizada.Direccion == "" {
		return nil, errors.New("la dirección no puede estar vacía")
	}

	err := db.Model(&existente).Updates(actualizada).Error
	if err != nil {
		return nil, err
	}
	return &existente, nil
}

func DeleteDirCliente(db *gorm.DB, id uint) error {
	return db.Delete(&modelos.DirCliente{}, id).Error
}
