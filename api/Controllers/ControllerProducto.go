package Controllers

import (
	modelos "backend-inventario/api/Models"
	"errors"

	"gorm.io/gorm"
)

// GetProductos obtiene todos los productos
func GetProductos(db *gorm.DB) ([]modelos.Producto, error) {
	var productos []modelos.Producto
	if err := db.Preload("Proveedor").Find(&productos).Error; err != nil {
		return nil, err
	}
	return productos, nil
}

// GetProductoBySKU obtiene un producto por su SKU
func GetProductoBySKU(db *gorm.DB, sku string) (*modelos.Producto, error) {
	var producto modelos.Producto
	if err := db.First(&producto, "sku = ?", sku).Error; err != nil {
		return nil, err
	}
	return &producto, nil
}

// CreateProducto crea un nuevo producto
func CreateProducto(db *gorm.DB, producto *modelos.Producto) error {
	return db.Create(producto).Error
}

// UpdateProducto actualiza un producto existente
func UpdateProducto(db *gorm.DB, sku string, nuevo *modelos.Producto) (*modelos.Producto, error) {
	var existente modelos.Producto
	if err := db.First(&existente, "sku = ?", sku).Error; err != nil {
		return nil, errors.New("producto no encontrado")
	}

	err := db.Model(&existente).Updates(modelos.Producto{
		Nombre:      nuevo.Nombre,
		Descripcion: nuevo.Descripcion,
		Marca:       nuevo.Marca,
		Peso:        nuevo.Peso,
		Largo:       nuevo.Largo,
		Ancho:       nuevo.Ancho,
		Alto:        nuevo.Alto,
		Precio:      nuevo.Precio,
	}).Error

	if err != nil {
		return nil, err
	}

	return &existente, nil
}

// DeleteProducto elimina un producto
func DeleteProducto(db *gorm.DB, sku string) error {
	return db.Delete(&modelos.Producto{}, "sku = ?", sku).Error
}
