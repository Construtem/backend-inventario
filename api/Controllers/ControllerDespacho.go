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

	err := db.
		Preload("Cotizacion.Cliente.Tipo").
		Preload("Cotizacion.Usuario.Rol").
		Preload("Camion.Tipo").
		Preload("OrigenSucursal.Tipo").
		Preload("DestinoDirCliente.Cliente.Tipo").
		Preload("ProductosDespacho.Producto").
		Find(&despachos).Error
	if err != nil {
		return nil, err
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
			Despacho:      despacho,
			CantidadItems: totalItems,
			TotalKg:       totalKg,
		})
	}
	return resultado, nil
}

func GetDespachoByID(db *gorm.DB, id uint) (*DespachoConTotales, error) {
	var despacho modelos.Despacho
	err := db.
		Preload("Cotizacion.Cliente.Tipo").
		Preload("Cotizacion.Usuario.Rol").
		Preload("Camion.Tipo").
		Preload("OrigenSucursal.Tipo").
		Preload("DestinoDirCliente.Cliente.Tipo").
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
		Despacho:      despacho,
		CantidadItems: totalItems,
		TotalKg:       totalKg,
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

func CalcularDespacho(db *gorm.DB, cotID uint) ([]modelos.Despacho, error) {
	// Se eliminan previamente los despachos existentes para esta cotización (si los hay)
	if err := db.Where("cotizacion_id = ?", cotID).Delete(&modelos.Despacho{}).Error; err != nil {
		return nil, err
	}

	// Se buscan los ítems de la cotización con sus productos y la cotización en sí
	var items []modelos.CotizacionItem
	err := db.
		Preload("Producto").
		Preload("Cotizacion").
		Where("cotizacion_id = ?", cotID).
		Find(&items).Error
	if err != nil {
		return nil, err
	}
	// 🚨 Si no hay productos asociados a la cotización, se devuelve un error
	if len(items) == 0 {
		return nil, errors.New("no hay productos en la cotización")
	}

	// 🚚 Se obtiene el tipo de camión con ID 1 (esto puede mejorarse a futuro para soportar múltiples tipos)
	var tipo modelos.TipoCamion
	if err := db.First(&tipo, 1).Error; err != nil {
		return nil, err
	}

	// 🧱 Se define una estructura interna para manejar cada unidad de producto
	type Unidad struct {
		SKU        string
		Peso       float64
		Volumen    float64
		SucursalID uint
	}
	var unidades []Unidad

	// 🧮 Se desglosan los ítems en unidades individuales (uno por cantidad), calculando su volumen
	for _, item := range items {
		vol := item.Producto.Largo * item.Producto.Ancho * item.Producto.Alto / 1_000_000
		for i := 0; i < item.Cantidad; i++ {
			unidades = append(unidades, Unidad{
				SKU:        item.Producto.SKU,
				Peso:       item.Producto.Peso,
				Volumen:    vol,
				SucursalID: item.SucursalID,
			})
		}
	}

	// 🏠 Se obtiene la dirección de destino del cliente usando el RUT presente en la cotización
	var destino modelos.DirCliente
	if err := db.Where("rut_cliente = ?", items[0].Cotizacion.RutCliente).First(&destino).Error; err != nil {
		return nil, err
	}

	// 🧩 Agrupamiento de unidades en despachos según capacidad de camión (peso y volumen)
	var grupos [][]Unidad
	var actuales []Unidad
	var pesoActual, volumenActual float64

	for _, u := range unidades {
		if pesoActual+u.Peso > tipo.PesoMaximo || volumenActual+u.Volumen > tipo.Volumen {
			grupos = append(grupos, actuales) // se cierra el grupo actual
			actuales = []Unidad{}             // se empieza un nuevo grupo
			pesoActual = 0
			volumenActual = 0
		}
		actuales = append(actuales, u)
		pesoActual += u.Peso
		volumenActual += u.Volumen
	}
	if len(actuales) > 0 {
		grupos = append(grupos, actuales) // se agrega el último grupo
	}

	var despachos []modelos.Despacho

	// 🚛 Por cada grupo de unidades, se crea un despacho nuevo
	for _, grupo := range grupos {
		despacho := modelos.Despacho{
			CotizacionID:  cotID,
			CamionID:      1,
			Origen:        grupo[0].SucursalID,
			Destino:       destino.ID,
			FechaDespacho: time.Now().AddDate(0, 0, 1), // Fecha de despacho al día siguiente
			ValorDespacho: 0,                           // Valor del despacho aún no calculado
			Estado:        "pendiente",                 // Estado inicial
		}

		if err := db.Create(&despacho).Error; err != nil {
			return nil, err
		}

		// 🧮 Se agrupan las unidades por SKU para registrar la cantidad total por producto en el despacho
		mapSKU := make(map[string]int)
		for _, u := range grupo {
			mapSKU[u.SKU]++
		}

		// 📦 Se crea el detalle del despacho (productos_despacho) por SKU y cantidad
		for sku, cantidad := range mapSKU {
			prod := modelos.ProductosDespacho{
				DespachoID: despacho.ID,
				ProductoID: sku,
				Cantidad:   cantidad,
			}
			if err := db.Create(&prod).Error; err != nil {
				return nil, err
			}
		}

		// 🧾 Se guarda el despacho generado
		despachos = append(despachos, despacho)
	}

	return despachos, nil
}

func GetDespachosPorCotizacion(db *gorm.DB, cotID uint) ([]modelos.Despacho, error) {
	// Se obtienen todos los despachos asociados a la cotización especificada
	var despachos []modelos.Despacho
	err := db.
		Preload("Cotizacion.Cliente.Tipo").
		Preload("Cotizacion.Usuario.Rol").
		Preload("Camion.Tipo").
		Preload("OrigenSucursal.Tipo").
		Preload("DestinoDirCliente.Cliente.Tipo").
		Preload("ProductosDespacho.Producto.proveedor").
		Where("cotizacion_id = ?", cotID).
		Find(&despachos).Error
	if err != nil {
		return nil, err
	}
	return despachos, nil
}

func AprobarDespacho(db *gorm.DB, cotID uint) error {
	// Se actualiza el estado del despacho a "aprobado" para la cotización especificada
	result := db.Model(&modelos.Despacho{}).
		Where("cotizacion_id = ?", cotID).
		Update("estado", "aprobado")

	if result.RowsAffected == 0 {
		return errors.New("no se encontró despacho para la cotización especificada")
	}
	return result.Error
}
