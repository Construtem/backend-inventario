package Controllers

import (
	modelos "backend-inventario/api/Models"
	"errors"
	"time"

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
		First(&stock, "proveedor_id = ? AND sku = ?", proveedorID, productoID).
		Error; err != nil {
		return nil, err
	}
	return &stock, nil
}

func CreateStockProveedor(db *gorm.DB, stock *modelos.StockProveedor) error {
	// Verificar que el proveedor existe
	var proveedor modelos.Proveedor
	if err := db.First(&proveedor, stock.ProveedorID).Error; err != nil {
		return errors.New("proveedor no encontrado")
	}

	// Verificar que el producto existe
	var producto modelos.Producto
	if err := db.First(&producto, "sku = ?", stock.ProductoID).Error; err != nil {
		return errors.New("producto no encontrado")
	}

	// Verificar si ya existe el registro de stock
	var existingStock modelos.StockProveedor
	if err := db.First(&existingStock, "proveedor_id = ? AND sku = ?", stock.ProveedorID, stock.ProductoID).Error; err == nil {
		return errors.New("ya existe stock para este producto y proveedor")
	}

	if stock.FechaIngreso.IsZero() {
		stock.FechaIngreso = time.Now()
	}
	
	return db.Create(stock).Error
}

func UpdateStockProveedor(db *gorm.DB, proveedorID uint, sku string, actualizado modelos.StockProveedor) error {
    var stock modelos.StockProveedor
    // Busca el registro con preload de Producto y Proveedor
    if err := db.Preload("Producto").Preload("Proveedor").
        First(&stock, "proveedor_id = ? AND sku = ?", proveedorID, sku).Error; err != nil {
        return errors.New("stock proveedor no encontrado")
    }

    // Actualiza solo los campos necesarios
    stock.Stock = actualizado.Stock

    // Producto: nombre y precio
    if stock.Producto.SKU != "" {
        stock.Producto.Nombre = actualizado.Producto.Nombre
        stock.Producto.Precio = actualizado.Producto.Precio
        db.Save(&stock.Producto)
    }

    // Proveedor: nombre/marca
    if stock.Proveedor.ID != 0 {
        stock.Proveedor.Marca = actualizado.Proveedor.Marca
        db.Save(&stock.Proveedor)
    }

    return db.Save(&stock).Error
}

func DeleteStockProveedor(db *gorm.DB, proveedorID uint, sku string) error {
    return db.Delete(&modelos.StockProveedor{}, "proveedor_id = ? AND sku = ?", proveedorID, sku).Error
}