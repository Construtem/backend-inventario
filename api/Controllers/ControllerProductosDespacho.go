package Controllers

import (
	modelos "backend-inventario/api/Models"

	"gorm.io/gorm"
)

// ProductoDespachoDetallado estructura mejorada para la respuesta del frontend
type ProductoDespachoDetallado struct {
	DespachoID  uint    `json:"despacho_id"`
	ProductoID  string  `json:"producto_id"`
	SKU         string  `json:"sku"`
	Nombre      string  `json:"nombre"`
	Descripcion string  `json:"descripcion"`
	Cantidad    int     `json:"cantidad"`
	Peso        float64 `json:"peso"`
	Alto        float64 `json:"alto"`
	Ancho       float64 `json:"ancho"`
	Largo       float64 `json:"largo"`
	Precio      float64 `json:"precio"`
	PesoTotal   float64 `json:"peso_total"`   // peso * cantidad
	PrecioTotal float64 `json:"precio_total"` // precio * cantidad
}

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

// GetProductosDespachoDetallado obtiene todos los productos de despacho con información detallada para el frontend
func GetProductosDespachoDetallado(db *gorm.DB) ([]ProductoDespachoDetallado, error) {
	var productosDespacho []modelos.ProductosDespacho
	var resultado []ProductoDespachoDetallado

	err := db.Preload("Producto").Find(&productosDespacho).Error
	if err != nil {
		return nil, err
	}

	for _, pd := range productosDespacho {
		detallado := ProductoDespachoDetallado{
			DespachoID:  pd.DespachoID,
			ProductoID:  pd.ProductoID,
			SKU:         pd.Producto.SKU,
			Nombre:      pd.Producto.Nombre,
			Descripcion: pd.Producto.Descripcion,
			Cantidad:    pd.Cantidad,
			Peso:        pd.Producto.Peso,
			Alto:        pd.Producto.Alto,
			Ancho:       pd.Producto.Ancho,
			Largo:       pd.Producto.Largo,
			Precio:      pd.Producto.Precio,
			PesoTotal:   pd.Producto.Peso * float64(pd.Cantidad),
			PrecioTotal: pd.Producto.Precio * float64(pd.Cantidad),
		}
		resultado = append(resultado, detallado)
	}

	return resultado, nil
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

// GetProductosDespachoDetalladoByDespachoID obtiene productos de un despacho específico con información detallada
func GetProductosDespachoDetalladoByDespachoID(db *gorm.DB, despachoID uint) ([]ProductoDespachoDetallado, error) {
	var productosDespacho []modelos.ProductosDespacho
	var resultado []ProductoDespachoDetallado

	err := db.Preload("Producto").
		Where("despacho_id = ?", despachoID).
		Find(&productosDespacho).Error
	if err != nil {
		return nil, err
	}

	for _, pd := range productosDespacho {
		detallado := ProductoDespachoDetallado{
			DespachoID:  pd.DespachoID,
			ProductoID:  pd.ProductoID,
			SKU:         pd.Producto.SKU,
			Nombre:      pd.Producto.Nombre,
			Descripcion: pd.Producto.Descripcion,
			Cantidad:    pd.Cantidad,
			Peso:        pd.Producto.Peso,
			Alto:        pd.Producto.Alto,
			Ancho:       pd.Producto.Ancho,
			Largo:       pd.Producto.Largo,
			Precio:      pd.Producto.Precio,
			PesoTotal:   pd.Producto.Peso * float64(pd.Cantidad),
			PrecioTotal: pd.Producto.Precio * float64(pd.Cantidad),
		}
		resultado = append(resultado, detallado)
	}

	return resultado, nil
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
