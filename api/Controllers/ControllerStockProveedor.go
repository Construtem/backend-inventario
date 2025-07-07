package Controllers

import (
	modelos "backend-inventario/api/Models"
	"errors"

	"gorm.io/gorm"
)

func GetStockProveedor(db *gorm.DB) ([]modelos.StockProveedor, error) {
	var stock []modelos.StockProveedor
	if err := db.
		Preload("Proveedor").
		Preload("Producto").
		Find(&stock).
		Error; err != nil {
		return nil, err
	}
	return stock, nil
}

func GetStockProveedorByID(db *gorm.DB, proveedorID uint, productoID string) (*modelos.StockProveedor, error) {
	var stock modelos.StockProveedor
	if err := db.
		Preload("Proveedor").
		Preload("Producto").
		First(&stock, "proveedor_id = ? AND producto_id = ?", proveedorID, productoID).
		Error; err != nil {
		return nil, err
	}
	return &stock, nil
}

func CreateStockProveedor(db *gorm.DB, stock *modelos.StockProveedor) error {
	return db.Create(stock).Error
}

func UpdateStockProveedor(db *gorm.DB, proveedorID uint, productoID string, cantidad int) error {
	var stock modelos.StockProveedor
	if err := db.First(&stock, "proveedor_id = ? AND producto_id = ?", proveedorID, productoID).Error; err != nil {
		return errors.New("stock proveedor no encontrado")
	}

	stock.Stock = cantidad
	return db.Save(&stock).Error
}

func DeleteStockProveedor(db *gorm.DB, proveedorID uint, productoID string) error {
	return db.Delete(&modelos.StockProveedor{}, "proveedor_id = ? AND producto_id = ?", proveedorID, productoID).Error
}
