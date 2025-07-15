package Controllers

import (
	modelos "backend-inventario/api/Models"

	"gorm.io/gorm"
)

// GetSucursales obtiene todas las sucursales
func GetSucursales(db *gorm.DB) ([]modelos.Sucursal, error) {
	var sucursales []modelos.Sucursal
	if err := db.Preload("Tipo").Find(&sucursales).Error; err != nil {
		return nil, err
	}
	return sucursales, nil
}

// GetBodegas obtiene solo las sucursales que son bodegas
func GetBodegas(db *gorm.DB) ([]modelos.Sucursal, error) {
	var bodegas []modelos.Sucursal

	// Por ahora, devolvemos todas las sucursales para test
	// TODO: Aplicar filtro por tipo bodega una vez que verifiquemos que funciona
	if err := db.Preload("Tipo").Find(&bodegas).Error; err != nil {
		return nil, err
	}

	return bodegas, nil
}

// GetSucursalByID obtiene una ubicación por su ID
func GetSucursalByID(db *gorm.DB, id uint) (*modelos.Sucursal, error) {
	var sucursal modelos.Sucursal
	if err := db.Preload("Tipo").First(&sucursal, id).Error; err != nil {
		return nil, err
	}
	return &sucursal, nil
}

// CreateSucursal crea una nueva ubicación
func CreateSucursal(db *gorm.DB, nueva *modelos.Sucursal) error {
	return db.Create(nueva).Error
}

// UpdateSucursal actualiza una ubicación existente
func UpdateSucursal(db *gorm.DB, id uint, actualizada *modelos.Sucursal) (*modelos.Sucursal, error) {
	var existente modelos.Sucursal
	if err := db.First(&existente, id).Error; err != nil {
		return nil, err
	}
	if err := db.Model(&existente).Updates(actualizada).Error; err != nil {
		return nil, err
	}

	if err := db.Preload("Tipo").First(&existente, id).Error; err != nil {
		return nil, err
	}
	return &existente, nil
}

// DeleteSucursal elimina una ubicación
func DeleteSucursal(db *gorm.DB, id uint) error {
	return db.Delete(&modelos.Sucursal{}, id).Error
}
