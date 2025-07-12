package Controllers

import (
	modelos "backend-inventario/api/Models"
	"errors"
	"time"

	"gorm.io/gorm"
)

type DespachoConTotales struct {
	modelos.Despacho
	CantidadItems     int                         `json:"cantidad_items"`
	TotalKg           float64                     `json:"total_kg"`
	ProductosDespacho []modelos.ProductosDespacho `json:"items"`
}

func CreateDespacho(db *gorm.DB, despacho *modelos.Despacho, productos []modelos.ProductosDespacho) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// Validar que la fecha de despacho no sea en el pasado
		if despacho.FechaDespacho.Before(time.Now().Add(-24 * time.Hour)) {
			return errors.New("la fecha de despacho no puede ser en el pasado")
		}

		if err := tx.Create(despacho).Error; err != nil {
			return err
		}

		for _, p := range productos {
			p.DespachoID = despacho.ID

			if err := tx.Create(&p).Error; err != nil {
				return err
			}

			res := tx.Model(&modelos.StockSucursal{}).
				Where("producto_id = ? AND sucursal_id = ?", p.ProductoID, despacho.Origen).
				Update("cantidad", gorm.Expr("cantidad - ?", p.Cantidad))

			if res.Error != nil {
				return res.Error
			}
			if res.RowsAffected == 0 {
				return errors.New("no se encontró stock suficiente para el producto " + p.ProductoID)
			}
		}
		return nil
	})
}

func GetDespachos(db *gorm.DB) ([]DespachoConTotales, error) {
	var despachos []modelos.Despacho

	// Verificar la conexión a la base de datos
	if db == nil {
		return nil, errors.New("conexión a la base de datos no disponible")
	}

	err := db.
		Preload("Cotizacion.Cliente").
		Preload("Camion").
		Preload("OrigenSucursal").
		Preload("DestinoDirCliente").
		Preload("ProductosDespacho.Producto").
		Find(&despachos).Error
	if err != nil {
		return nil, errors.New("error al consultar despachos en la base de datos: " + err.Error())
	}

	var resultado []DespachoConTotales

	for _, despacho := range despachos {
		totalKg := 0.0
		totalItems := 0

		for _, producto := range despacho.ProductosDespacho {
			totalItems += producto.Cantidad
			totalKg += float64(producto.Cantidad) * producto.Producto.Peso
		}

		resultado = append(resultado, DespachoConTotales{
			Despacho:          despacho,
			CantidadItems:     totalItems,
			TotalKg:           totalKg,
			ProductosDespacho: despacho.ProductosDespacho,
		})
	}
	return resultado, nil
}

func GetDespachoByID(db *gorm.DB, id uint) (*DespachoConTotales, error) {
	var despacho modelos.Despacho
	err := db.
		Preload("Cotizacion.Cliente").
		Preload("Cotizacion").
		Preload("Camion").
		Preload("OrigenSucursal").
		Preload("DestinoDirCliente").
		Preload("ProductosDespacho.Producto").
		First(&despacho, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	var totalKg float64
	var totalItems int
	for _, p := range despacho.ProductosDespacho {
		totalItems += p.Cantidad
		totalKg += float64(p.Cantidad) * p.Producto.Peso
	}

	resultado := DespachoConTotales{
		Despacho:          despacho,
		CantidadItems:     totalItems,
		TotalKg:           totalKg,
		ProductosDespacho: despacho.ProductosDespacho,
	}
	return &resultado, nil
}

func UpdateDespacho(db *gorm.DB, id uint, actualizado *modelos.Despacho) error {
	var existente modelos.Despacho
	if err := db.First(&existente, id).Error; err != nil {
		return errors.New("despacho no encontrado")
	}
	return db.Model(&existente).Updates(actualizado).Error
}

func DeleteDespacho(db *gorm.DB, id uint) error {
	return db.Delete(&modelos.Despacho{}, id).Error
}
