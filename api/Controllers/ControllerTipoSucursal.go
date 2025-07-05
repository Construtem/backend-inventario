package Controllers

import (
	modelos "backend-inventario/api/Models"

	"gorm.io/gorm"
)

func GetTipoSucursal(db *gorm.DB) ([]modelos.TipoSucursal, error) {
	var tipos []modelos.TipoSucursal
	if err := db.Find(&tipos).Error; err != nil {
		return nil, err
	}
	return tipos, nil
}

func GetTipoSucursalByID(db *gorm.DB, id uint) (*modelos.TipoSucursal, error) {
	var tipo modelos.TipoSucursal
	if err := db.First(&tipo, id).Error; err != nil {
		return nil, err
	}
	return &tipo, nil
}

func CreateTipoSucursal(db *gorm.DB, nuevo *modelos.TipoSucursal) error {
	return db.Create(nuevo).Error
}

func UpdateTipoSucursal(db *gorm.DB, id uint, data map[string]interface{}) (*modelos.TipoSucursal, error) {
	var tipo modelos.TipoSucursal
	if err := db.First(&tipo, id).Error; err != nil {
		return nil, err
	}

	if err := db.Model(&tipo).Updates(data).Error; err != nil {
		return nil, err
	}
	return &tipo, nil
}

func DeleteTipoSucursal(db *gorm.DB, id uint) error {
	return db.Delete(&modelos.TipoSucursal{}, id).Error
}
