package Controllers

import (
	modelos "backend-inventario/api/Models"
	"errors"
	"fmt"
	"time"

	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"

	"gorm.io/gorm"
)

//"strconv"

// DespachoConTotales es una estructura que incluye un despacho y sus totales calculados
type DespachoConTotales struct {
	modelos.Despacho
	CantidadItems     int     `json:"cantidad_items"`
	TotalKg           float64 `json:"total_kg"`
	TotalPrecio       float64 `json:"total_precio"`
	IVA               float64
	ProductosDespacho []ProductoDespachoDetallado `json:"items"`
}

// Se define una estructura interna para manejar cada unidad de producto
type Unidad struct {
	SKU        string
	Peso       float64
	Volumen    float64
	SucursalID uint
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
		Preload("Cotizacion.Cliente.Tipo").
		Preload("Cotizacion.Usuario.Rol").
		Preload("Camion.Tipo").
		Preload("OrigenSucursal.Tipo").
		Preload("DestinoDirCliente.Cliente.Tipo").
		Preload("ProductosDespacho.Producto").
		Find(&despachos).Error
	if err != nil {
		return nil, errors.New("error al consultar despachos en la base de datos: " + err.Error())
	}

	var resultado []DespachoConTotales

	for _, despacho := range despachos {
		totalKg := 0.0
		totalItems := 0
		totalPrecio := 0.0
		var productosDetallados []ProductoDespachoDetallado

		for _, producto := range despacho.ProductosDespacho {
			totalItems += producto.Cantidad
			totalKg += float64(producto.Cantidad) * producto.Producto.Peso
			totalPrecio += float64(producto.Cantidad) * producto.Producto.Precio

			// Crear producto detallado
			detallado := ProductoDespachoDetallado{
				DespachoID:  producto.DespachoID,
				ProductoID:  producto.ProductoID,
				SKU:         producto.Producto.SKU,
				Nombre:      producto.Producto.Nombre,
				Descripcion: producto.Producto.Descripcion,
				Cantidad:    producto.Cantidad,
				Peso:        producto.Producto.Peso,
				Alto:        producto.Producto.Alto,
				Ancho:       producto.Producto.Ancho,
				Largo:       producto.Producto.Largo,
				Precio:      producto.Producto.Precio,
				PesoTotal:   producto.Producto.Peso * float64(producto.Cantidad),
				PrecioTotal: producto.Producto.Precio * float64(producto.Cantidad),
			}
			productosDetallados = append(productosDetallados, detallado)
		}

		resultado = append(resultado, DespachoConTotales{
			Despacho:          despacho,
			CantidadItems:     totalItems,
			TotalKg:           totalKg,
			TotalPrecio:       totalPrecio,
			ProductosDespacho: productosDetallados,
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
		Preload("ProductosDespacho.Producto.Proveedor").
		First(&despacho, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	var totalKg float64
	var totalItems int
	var totalPrecio float64
	var productosDetallados []ProductoDespachoDetallado

	for _, p := range despacho.ProductosDespacho {
		totalItems += p.Cantidad
		totalKg += float64(p.Cantidad) * p.Producto.Peso
		totalPrecio += float64(p.Cantidad) * p.Producto.Precio

		// Crear producto detallado
		detallado := ProductoDespachoDetallado{
			DespachoID:  p.DespachoID,
			ProductoID:  p.ProductoID,
			SKU:         p.Producto.SKU,
			Nombre:      p.Producto.Nombre,
			Descripcion: p.Producto.Descripcion,
			Cantidad:    p.Cantidad,
			Peso:        p.Producto.Peso,
			Alto:        p.Producto.Alto,
			Ancho:       p.Producto.Ancho,
			Largo:       p.Producto.Largo,
			Precio:      p.Producto.Precio,
			PesoTotal:   p.Producto.Peso * float64(p.Cantidad),
			PrecioTotal: p.Producto.Precio * float64(p.Cantidad),
		}
		productosDetallados = append(productosDetallados, detallado)
	}

	iva := totalPrecio * 0.19

	resultado := DespachoConTotales{
		Despacho:          despacho,
		CantidadItems:     totalItems,
		TotalKg:           totalKg,
		TotalPrecio:       totalPrecio,
		IVA:               iva,
		ProductosDespacho: productosDetallados,
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

	var tiposDisponibles []modelos.TipoCamion
	if err := db.Order("peso_maximo ASC").Find(&tiposDisponibles).Error; err != nil {
		return nil, err
	}
	if len(tiposDisponibles) == 0 {
		return nil, errors.New("no hay tipos de camión disponibles")
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
		actuales = append(actuales, u)
		pesoActual += u.Peso
		volumenActual += u.Volumen

		cabe := false
		for _, tipo := range tiposDisponibles {
			if tipo.PesoMaximo >= pesoActual && tipo.Volumen >= volumenActual {
				cabe = true
				break
			}
		}

		if !cabe {
			grupoValido := actuales[:len(actuales)-1] // se quita la última unidad que no cabe
			grupos = append(grupos, grupoValido)

			actuales = []Unidad{u} // reiniciar el grupo con la unidad actual
			pesoActual = u.Peso
			volumenActual = u.Volumen
		}
	}
	if len(actuales) > 0 {
		grupos = append(grupos, actuales) // se agrega el último grupo
	}

	// 🔢 Índice de precio por km (puede venir de una config .env o base de datos)
	precioPorKm := 500.0 // Ejemplo: 500 CLP por km

	// Dirección origen y destino (como string)
	origenStr := fmt.Sprintf("%s, %s", items[0].Sucursal.Direccion, items[0].Sucursal.Ciudad)
	destinoStr := fmt.Sprintf("%s, %s", destino.Direccion, destino.Ciudad)

	// Obtener distancia
	distanciaKm, err := obtenerDistanciaEnKm(origenStr, destinoStr)
	if err != nil {
		return nil, fmt.Errorf("error al obtener distancia: %v", err)
	}

	// Calcular costo total de envío
	costoTotalEnvio := float64(len(grupos)) * distanciaKm * precioPorKm

	var despachos []modelos.Despacho

	// 🚛 Por cada grupo de unidades, se crea un despacho nuevo
	for _, grupo := range grupos {
		var tipoCamionID uint = 0
		for _, tipo := range tiposDisponibles {
			if tipo.PesoMaximo >= pesoTotal(grupo) && tipo.Volumen >= volumenTotal(grupo) {
				tipoCamionID = tipo.ID
				break
			}
		}
		if tipoCamionID == 0 {
			return nil, errors.New("no hay tipo de camión disponible para un grupo de productos")
		}

		var camion modelos.Camion
		if err := db.Where("tipo_id = ? AND activo = true", tipoCamionID).First(&camion).Error; err != nil {
			return nil, fmt.Errorf("no hay camiones disponibles del tipo %d", tipoCamionID)
		}

		// Crear el despacho con la información del grupo
		despacho := modelos.Despacho{
			CotizacionID:  cotID,
			CamionID:      camion.ID,
			Origen:        grupo[0].SucursalID,
			Destino:       destino.ID,
			FechaDespacho: time.Now().AddDate(0, 0, 1),            // Fecha de despacho al día siguiente
			ValorDespacho: costoTotalEnvio / float64(len(grupos)), // repartir el total entre camiones
			Estado:        "pendiente",                            // Estado inicial
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

// Funciones auxiliares
// Cambia el estado de los despachos asociados a una cotización
func CambiarEstadoDespachosPorCotizacion(db *gorm.DB, cotizacionID uint, estado string) error {
	if estado != "pendiente" && estado != "entregado" && estado != "aprobado" {
		return errors.New("estado no permitido")
	}
	result := db.Model(&modelos.Despacho{}).
		Where("cotizacion_id = ?", cotizacionID).
		Update("estado", estado)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no se encontraron despachos para la cotización")
	}
	return nil
}
func pesoTotal(grupo []Unidad) float64 {
	var total float64
	for _, u := range grupo {
		total += u.Peso
	}
	return total
}

func volumenTotal(grupo []Unidad) float64 {
	var total float64
	for _, u := range grupo {
		total += u.Volumen
	}
	return total
}

// GetDespachoDistanciaByID obtiene un despacho con información completa para rutas
func GetDespachoDistanciaByID(db *gorm.DB, id uint) (*modelos.DespachoDistanciaResponse, error) {
	var despacho modelos.Despacho

	// Hacer join con todas las tablas relacionadas
	err := db.Preload("Cotizacion.Cliente").
		Preload("Camion").
		Preload("OrigenSucursal").
		Preload("DestinoDirCliente.Cliente").
		Preload("ProductosDespacho.Producto").
		First(&despacho, id).Error

	if err != nil {
		return nil, err
	}

	// Calcular totales
	cantidadItems := 0
	totalKg := 0.0

	for _, producto := range despacho.ProductosDespacho {
		cantidadItems += producto.Cantidad
		totalKg += float64(producto.Cantidad) * producto.Producto.Peso
	}

	// Construir la respuesta
	response := &modelos.DespachoDistanciaResponse{
		ID:                 despacho.ID,
		CotizacionID:       despacho.CotizacionID,
		CamionID:           despacho.CamionID,
		Origen:             despacho.Origen,
		Destino:            despacho.Destino,
		FechaDespacho:      despacho.FechaDespacho,
		ValorDespacho:      despacho.ValorDespacho,
		CantidadItems:      cantidadItems,
		TotalKg:            totalKg,
		DistanciaCalculada: despacho.DistanciaCalculada,
		TiempoEstimado:     despacho.TiempoEstimado,
	}

	// Agregar información de cotización si existe
	if despacho.Cotizacion.ID != 0 {
		response.Cotizacion = &modelos.CotizacionDistanciaResponse{
			Estado: despacho.Cotizacion.Estado,
		}

		if despacho.Cotizacion.Cliente.Rut != "" {
			response.Cotizacion.Cliente = &modelos.ClienteDistanciaResponse{
				Nombre: despacho.Cotizacion.Cliente.Nombre,
				Email:  despacho.Cotizacion.Cliente.Email,
			}
		}
	}

	// Agregar información del camión
	if despacho.Camion.ID != 0 {
		response.Camion = &modelos.CamionDistanciaResponse{
			Patente: despacho.Camion.Patente,
		}
	}

	// Agregar información de la sucursal origen
	if despacho.OrigenSucursal.ID != 0 {
		response.OrigenSucursal = &modelos.SucursalDistanciaResponse{
			Nombre:    despacho.OrigenSucursal.Nombre,
			Direccion: despacho.OrigenSucursal.Direccion,
			Comuna:    despacho.OrigenSucursal.Comuna,
			Ciudad:    despacho.OrigenSucursal.Ciudad,
		}
	}

	// Agregar información del destino
	if despacho.DestinoDirCliente.ID != 0 {
		response.DestinoDirCliente = &modelos.DirClienteDistanciaResponse{
			Direccion: despacho.DestinoDirCliente.Direccion,
			Comuna:    despacho.DestinoDirCliente.Comuna,
			Ciudad:    despacho.DestinoDirCliente.Ciudad,
		}
	}

	return response, nil
}

// GetDespachosDistancia obtiene todos los despachos con información para rutas
func GetDespachosDistancia(db *gorm.DB) ([]modelos.DespachoDistanciaResponse, error) {
	var despachos []modelos.Despacho

	// Obtener todos los despachos con sus relaciones
	err := db.Preload("Cotizacion.Cliente").
		Preload("Camion").
		Preload("OrigenSucursal").
		Preload("DestinoDirCliente.Cliente").
		Preload("ProductosDespacho.Producto").
		Find(&despachos).Error

	if err != nil {
		return nil, err
	}

	var response []modelos.DespachoDistanciaResponse

	for _, despacho := range despachos {
		// Calcular totales para cada despacho
		cantidadItems := 0
		totalKg := 0.0

		for _, producto := range despacho.ProductosDespacho {
			cantidadItems += producto.Cantidad
			totalKg += float64(producto.Cantidad) * producto.Producto.Peso
		}

		// Construir respuesta individual
		despachoResponse := modelos.DespachoDistanciaResponse{
			ID:                 despacho.ID,
			CotizacionID:       despacho.CotizacionID,
			CamionID:           despacho.CamionID,
			Origen:             despacho.Origen,
			Destino:            despacho.Destino,
			FechaDespacho:      despacho.FechaDespacho,
			ValorDespacho:      despacho.ValorDespacho,
			CantidadItems:      cantidadItems,
			TotalKg:            totalKg,
			DistanciaCalculada: despacho.DistanciaCalculada,
			TiempoEstimado:     despacho.TiempoEstimado,
		}

		// Agregar información de cotización si existe
		if despacho.Cotizacion.ID != 0 {
			despachoResponse.Cotizacion = &modelos.CotizacionDistanciaResponse{
				Estado: despacho.Cotizacion.Estado,
			}

			if despacho.Cotizacion.Cliente.Rut != "" {
				despachoResponse.Cotizacion.Cliente = &modelos.ClienteDistanciaResponse{
					Nombre: despacho.Cotizacion.Cliente.Nombre,
					Email:  despacho.Cotizacion.Cliente.Email,
				}
			}
		}

		// Agregar información del camión
		if despacho.Camion.ID != 0 {
			despachoResponse.Camion = &modelos.CamionDistanciaResponse{
				Patente: despacho.Camion.Patente,
			}
		}

		// Agregar información de la sucursal origen
		if despacho.OrigenSucursal.ID != 0 {
			despachoResponse.OrigenSucursal = &modelos.SucursalDistanciaResponse{
				Nombre:    despacho.OrigenSucursal.Nombre,
				Direccion: despacho.OrigenSucursal.Direccion,
				Comuna:    despacho.OrigenSucursal.Comuna,
				Ciudad:    despacho.OrigenSucursal.Ciudad,
			}
		}

		// Agregar información del destino
		if despacho.DestinoDirCliente.ID != 0 {
			despachoResponse.DestinoDirCliente = &modelos.DirClienteDistanciaResponse{
				Direccion: despacho.DestinoDirCliente.Direccion,
				Comuna:    despacho.DestinoDirCliente.Comuna,
				Ciudad:    despacho.DestinoDirCliente.Ciudad,
			}
		}

		response = append(response, despachoResponse)
	}

	return response, nil
}

// ActualizarDistanciaDespacho actualiza la distancia y tiempo de un despacho
func ActualizarDistanciaDespacho(db *gorm.DB, id uint, distancia, tiempo string) error {
	result := db.Model(&modelos.Despacho{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"distancia_calculada": distancia,
			"tiempo_estimado":     tiempo,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("despacho no encontrado")
	}

	return nil
}

// GetFacturaElectronicaByDespachoID retorna una factura electrónica simulada para un despacho
func GetFacturaElectronicaByDespachoID(db *gorm.DB, despachoID uint) (map[string]interface{}, error) {
	ficha, err := GetDespachoByID(db, despachoID)
	if err != nil {
		return nil, err
	}
	factura := map[string]interface{}{
		"folio":       ficha.ID,
		"fecha":       ficha.FechaDespacho,
		"cliente":     ficha.Cotizacion.Cliente.Nombre,
		"rut_cliente": ficha.Cotizacion.Cliente.Rut,
		"total":       ficha.TotalPrecio,
		"estado":      ficha.Cotizacion.Estado,
		"tipo":        "Factura Electrónica",
	}
	return factura, nil
}

func obtenerDistanciaEnKm(origen, destino string) (float64, error) {
	apiKey := os.Getenv("GOOGLE_MAPS_API_KEY")
	urlGoogle := fmt.Sprintf("https://maps.googleapis.com/maps/api/directions/json?origin=%s&destination=%s&key=%s",
		url.QueryEscape(origen), url.QueryEscape(destino), apiKey)

	resp, err := http.Get(urlGoogle)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var data map[string]interface{}
	json.Unmarshal(body, &data)

	routes := data["routes"].([]interface{})
	if len(routes) == 0 {
		return 0, fmt.Errorf("no se encontraron rutas")
	}

	legs := routes[0].(map[string]interface{})["legs"].([]interface{})
	distancia := legs[0].(map[string]interface{})["distance"].(map[string]interface{})["value"].(float64) // metros

	return distancia / 1000, nil // convertimos a km
}
