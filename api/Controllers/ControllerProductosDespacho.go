package Controllers

import (
	modelos "backend-inventario/api/Models"

	"gorm.io/gorm"
)

// GetProductosDespacho obtiene todos los productos de despacho con información relacionada
func GetProductosDespacho(db *gorm.DB) ([]modelos.ProductosDespacho, error) {
	var productosDespacho []modelos.ProductosDespacho

	// Precargar las relaciones con Despacho y Producto
	err := db.Preload("Despacho").
		Preload("Despacho.Cotizacion").
		Preload("Despacho.Camion").
		Preload("Despacho.OrigenSucursal").
		Preload("Despacho.DestinoDirCliente").
		Preload("Producto").
		Preload("Producto.Proveedor").
		Find(&productosDespacho).Error

	return productosDespacho, err
}

// GetProductosDespachoByDespachoID obtiene todos los productos de un despacho específico
func GetProductosDespachoByDespachoID(db *gorm.DB, despachoID uint) ([]modelos.ProductosDespacho, error) {
	var productosDespacho []modelos.ProductosDespacho

	err := db.Preload("Despacho").
		Preload("Despacho.Cotizacion").
		Preload("Despacho.Camion").
		Preload("Despacho.OrigenSucursal").
		Preload("Despacho.DestinoDirCliente").
		Preload("Producto").
		Preload("Producto.Proveedor").
		Where("despacho_id = ?", despachoID).
		Find(&productosDespacho).Error

	return productosDespacho, err
}

// GetProductoDespacho obtiene un producto de despacho específico
func GetProductoDespacho(db *gorm.DB, despachoID uint, productoID string) (modelos.ProductosDespacho, error) {
	var productoDespacho modelos.ProductosDespacho

	err := db.Preload("Despacho").
		Preload("Despacho.Cotizacion").
		Preload("Despacho.Camion").
		Preload("Despacho.OrigenSucursal").
		Preload("Despacho.DestinoDirCliente").
		Preload("Producto").
		Preload("Producto.Proveedor").
		Where("despacho_id = ? AND producto_id = ?", despachoID, productoID).
		First(&productoDespacho).Error

	return productoDespacho, err
}

// CreateProductoDespacho crea un nuevo producto de despacho
func CreateProductoDespacho(db *gorm.DB, productoDespacho modelos.ProductosDespacho) (modelos.ProductosDespacho, error) {
	err := db.Create(&productoDespacho).Error
	if err != nil {
		return productoDespacho, err
	}

	// Recargar con relaciones
	err = db.Preload("Despacho").
		Preload("Producto").
		Where("despacho_id = ? AND producto_id = ?", productoDespacho.DespachoID, productoDespacho.ProductoID).
		First(&productoDespacho).Error

	return productoDespacho, err
}

// UpdateProductoDespacho actualiza un producto de despacho existente
func UpdateProductoDespacho(db *gorm.DB, despachoID uint, productoID string, productoDespacho modelos.ProductosDespacho) (modelos.ProductosDespacho, error) {
	err := db.Model(&modelos.ProductosDespacho{}).
		Where("despacho_id = ? AND producto_id = ?", despachoID, productoID).
		Updates(productoDespacho).Error

	if err != nil {
		return productoDespacho, err
	}

	// Recargar con relaciones
	err = db.Preload("Despacho").
		Preload("Producto").
		Where("despacho_id = ? AND producto_id = ?", despachoID, productoID).
		First(&productoDespacho).Error

	return productoDespacho, err
}

// DeleteProductoDespacho elimina un producto de despacho
func DeleteProductoDespacho(db *gorm.DB, despachoID uint, productoID string) error {
	err := db.Where("despacho_id = ? AND producto_id = ?", despachoID, productoID).
		Delete(&modelos.ProductosDespacho{}).Error

	return err
}
