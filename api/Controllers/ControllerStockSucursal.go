package Controllers

import (
	modelos "backend-inventario/api/Models"

	"gorm.io/gorm"
)

// GetStockSucursal obtiene todos los registros de stock por sucursal
func GetStockSucursal(db *gorm.DB) ([]modelos.StockSucursal, error) {
	var stocks []modelos.StockSucursal
	if err := db.Preload("Producto").Preload("Sucursal").Find(&stocks).Error; err != nil {
		return nil, err
	}
	return stocks, nil
}

// GetStockSucursalByID obtiene un registro de stock por su ID
func GetStockSucursalByID(db *gorm.DB, sucursalID uint, sku string) (*modelos.StockSucursal, error) {
	var stock modelos.StockSucursal
	if err := db.Where("sucursal_id = ? AND sku = ?", sucursalID, sku).Preload("Producto").Preload("Sucursal").First(&stock).Error; err != nil {
		return nil, err
	}
	return &stock, nil
}

// CreateStockSucursal crea un nuevo registro de stock
func CreateStockSucursal(db *gorm.DB, nuevo *modelos.StockSucursal) error {
	return db.Create(nuevo).Error
}

// UpdateStockSucursal actualiza un registro de stock existente
func UpdateStockSucursal(db *gorm.DB, sku string, sucursalID uint, actualizado *modelos.StockSucursal) error {
	return db.Model(&modelos.StockSucursal{}).
		Where("sku = ? AND sucursal_id = ?", sku, sucursalID).
		Updates(actualizado).Error
}

// DeleteStockSucursal elimina un registro de stock
func DeleteStockSucursal(db *gorm.DB, sku string, sucursalID uint) error {
	return db.Where("sku = ? AND sucursal_id = ?", sku, sucursalID).Delete(&modelos.StockSucursal{}).Error
}
