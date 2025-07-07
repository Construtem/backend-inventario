package Controllers

import (
	modelos "backend-inventario/api/Models"
	"errors"

	"gorm.io/gorm"
)

func GetProveedores(db *gorm.DB) ([]modelos.Proveedor, error) {
	var proveedores []modelos.Proveedor
	if err := db.Find(&proveedores).Error; err != nil {
		return nil, err
	}
	return proveedores, nil
}

func GetProveedorByID(db *gorm.DB, id uint) (*modelos.Proveedor, error) {
	var proveedor modelos.Proveedor
	if err := db.First(&proveedor, id).Error; err != nil {
		return nil, errors.New("proveedor no encontrado")
	}
	return &proveedor, nil
}

func CreateProveedor(db *gorm.DB, proveedor *modelos.Proveedor) error {
	return db.Create(proveedor).Error
}

func UpdateProveedor(db *gorm.DB, id uint, actualizado *modelos.Proveedor) (*modelos.Proveedor, error) {
	var existente modelos.Proveedor
	if err := db.First(&existente, id).Error; err != nil {
		return nil, errors.New("proveedor no encontrado")
	}

	existente.Marca = actualizado.Marca
	if err := db.Save(&existente).Error; err != nil {
		return nil, err
	}
	return &existente, nil
}

func DeleteProveedor(db *gorm.DB, id uint) error {
	return db.Delete(&modelos.Proveedor{}, id).Error
}
