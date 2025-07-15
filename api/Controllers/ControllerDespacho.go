package Controllers

import (
	modelos "backend-inventario/api/Models"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"gorm.io/gorm"
)

// Estructuras para la respuesta de Google Maps Distance Matrix API
type DistanceMatrixResponse struct {
	Status string `json:"status"`
	Rows   []struct {
		Elements []struct {
			Status   string `json:"status"`
			Distance struct {
				Text  string `json:"text"`
				Value int    `json:"value"` // distancia en metros
			} `json:"distance"`
			Duration struct {
				Text  string `json:"text"`
				Value int    `json:"value"` // duración en segundos
			} `json:"duration"`
		} `json:"elements"`
	} `json:"rows"`
}

// DespachoConTotales es una estructura que incluye un despacho y sus totales calculados
type DespachoConTotales struct {
	modelos.Despacho
	CantidadItems     int                         `json:"cantidad_items"`
	TotalKg           float64                     `json:"total_kg"`
	TotalPrecio       float64                     `json:"total_precio"`
	ProductosDespacho []ProductoDespachoDetallado `json:"productos"`
	DistanciaKm       float64                     `json:"distancia_km,omitempty"`      // Distancia en kilómetros
	DuracionMinutos   int                         `json:"duracion_minutos,omitempty"`  // Duración en minutos
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
	return GetDespachosConDistancia(db, "")
}

func GetDespachosConDistancia(db *gorm.DB, googleMapsAPIKey string) ([]DespachoConTotales, error) {
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

		// Calcular distancia si se proporciona la API key
		var distanciaKm float64
		var duracionMinutos int
		if googleMapsAPIKey != "" {
			if dist, dur, err := CalcularDistanciaDespacho(db, googleMapsAPIKey, &despacho); err == nil {
				distanciaKm = dist
				duracionMinutos = dur
			}
			// Si hay error en el cálculo de distancia, simplemente no la incluimos pero continuamos
		}

		resultado = append(resultado, DespachoConTotales{
			Despacho:          despacho,
			CantidadItems:     totalItems,
			TotalKg:           totalKg,
			TotalPrecio:       totalPrecio,
			DistanciaKm:       distanciaKm,
			DuracionMinutos:   duracionMinutos,
			ProductosDespacho: productosDetallados,
		})
	}
	return resultado, nil
}

func GetDespachoByID(db *gorm.DB, id uint) (*DespachoConTotales, error) {
	return GetDespachoByIDConDistancia(db, id, "")
}

func GetDespachoByIDConDistancia(db *gorm.DB, id uint, googleMapsAPIKey string) (*DespachoConTotales, error) {
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

	// Calcular distancia si se proporciona la API key
	var distanciaKm float64
	var duracionMinutos int
	if googleMapsAPIKey != "" {
		if dist, dur, err := CalcularDistanciaDespacho(db, googleMapsAPIKey, &despacho); err == nil {
			distanciaKm = dist
			duracionMinutos = dur
		}
		// Si hay error en el cálculo de distancia, simplemente no la incluimos pero continuamos
	}

	resultado := DespachoConTotales{
		Despacho:          despacho,
		CantidadItems:     totalItems,
		TotalKg:           totalKg,
		TotalPrecio:       totalPrecio,
		DistanciaKm:       distanciaKm,
		DuracionMinutos:   duracionMinutos,
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
	return CalcularDespachoConDistancia(db, cotID, "")
}

func CalcularDespachoConDistancia(db *gorm.DB, cotID uint, googleMapsAPIKey string) ([]modelos.Despacho, error) {
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
			FechaDespacho: time.Now().AddDate(0, 0, 1), // Fecha de despacho al día siguiente
			ValorDespacho: 0,                           // Valor del despacho aún no calculado
			Estado:        "pendiente",                 // Estado inicial
		}

		// 🗺️ Calcular distancia si se proporciona la API key
		if googleMapsAPIKey != "" {
			if distancia, _, err := CalcularDistanciaDespacho(db, googleMapsAPIKey, &despacho); err == nil {
				// Aquí podrías calcular el valor del despacho basado en la distancia
				// Por ejemplo: despacho.ValorDespacho = distancia * tarifaPorKm
				_ = distancia // Por ahora solo calculamos pero no asignamos valor automáticamente
			}
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
	return GetDespachosPorCotizacionConDistancia(db, cotID, "")
}

func GetDespachosPorCotizacionConDistancia(db *gorm.DB, cotID uint, googleMapsAPIKey string) ([]modelos.Despacho, error) {
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

// GetDespachosPorCotizacionDetallado obtiene despachos con totales y distancias calculadas
func GetDespachosPorCotizacionDetallado(db *gorm.DB, cotID uint, googleMapsAPIKey string) ([]DespachoConTotales, error) {
	var despachos []modelos.Despacho
	err := db.
		Preload("Cotizacion.Cliente.Tipo").
		Preload("Cotizacion.Usuario.Rol").
		Preload("Camion.Tipo").
		Preload("OrigenSucursal.Tipo").
		Preload("DestinoDirCliente.Cliente.Tipo").
		Preload("ProductosDespacho.Producto").
		Where("cotizacion_id = ?", cotID).
		Find(&despachos).Error
	if err != nil {
		return nil, err
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

		// Calcular distancia si se proporciona la API key
		var distanciaKm float64
		var duracionMinutos int
		if googleMapsAPIKey != "" {
			if dist, dur, err := CalcularDistanciaDespacho(db, googleMapsAPIKey, &despacho); err == nil {
				distanciaKm = dist
				duracionMinutos = dur
			}
		}

		resultado = append(resultado, DespachoConTotales{
			Despacho:          despacho,
			CantidadItems:     totalItems,
			TotalKg:           totalKg,
			TotalPrecio:       totalPrecio,
			DistanciaKm:       distanciaKm,
			DuracionMinutos:   duracionMinutos,
			ProductosDespacho: productosDetallados,
		})
	}

	return resultado, nil
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

// CalcularDistancia calcula la distancia entre origen y destino usando Google Maps Distance Matrix API
func CalcularDistancia(apiKey, origen, destino string) (float64, int, error) {
	if apiKey == "" {
		return 0, 0, errors.New("API key de Google Maps no configurada")
	}

	// Construir la URL para la API de Google Maps Distance Matrix
	baseURL := "https://maps.googleapis.com/maps/api/distancematrix/json"
	params := url.Values{}
	params.Add("origins", origen)
	params.Add("destinations", destino)
	params.Add("mode", "driving")
	params.Add("language", "es")
	params.Add("key", apiKey)

	requestURL := baseURL + "?" + params.Encode()
	fmt.Printf("🌐 URL de la API: %s\n", requestURL)

	// Realizar la petición HTTP
	resp, err := http.Get(requestURL)
	if err != nil {
		fmt.Printf("❌ Error en petición HTTP: %v\n", err)
		return 0, 0, fmt.Errorf("error al consultar Google Maps API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("❌ Error HTTP status: %d\n", resp.StatusCode)
		return 0, 0, fmt.Errorf("error HTTP %d al consultar Google Maps API", resp.StatusCode)
	}

	// Leer la respuesta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("❌ Error leyendo respuesta: %v\n", err)
		return 0, 0, fmt.Errorf("error al leer respuesta de Google Maps API: %v", err)
	}

	fmt.Printf("📡 Respuesta de Google Maps API: %s\n", string(body))

	// Parsear la respuesta JSON
	var response DistanceMatrixResponse
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Printf("❌ Error parseando JSON: %v\n", err)
		return 0, 0, fmt.Errorf("error al parsear respuesta de Google Maps API: %v", err)
	}

	// Verificar el estado de la respuesta
	if response.Status != "OK" {
		fmt.Printf("❌ Error en status de API: %s\n", response.Status)
		return 0, 0, fmt.Errorf("error en Google Maps API: %s", response.Status)
	}

	// Verificar que hay resultados
	if len(response.Rows) == 0 || len(response.Rows[0].Elements) == 0 {
		return 0, 0, errors.New("no se encontraron resultados de distancia")
	}

	element := response.Rows[0].Elements[0]
	if element.Status != "OK" {
		return 0, 0, fmt.Errorf("error en cálculo de distancia: %s", element.Status)
	}

	// Convertir de metros a kilómetros y de segundos a minutos
	distanciaKm := float64(element.Distance.Value) / 1000.0
	duracionMinutos := element.Duration.Value / 60

	return distanciaKm, duracionMinutos, nil
}

// CalcularDistanciaDespacho calcula la distancia de un despacho específico
func CalcularDistanciaDespacho(db *gorm.DB, apiKey string, despacho *modelos.Despacho) (float64, int, error) {
	// Obtener información de origen (sucursal)
	var sucursal modelos.Sucursal
	if err := db.First(&sucursal, despacho.Origen).Error; err != nil {
		return 0, 0, fmt.Errorf("no se encontró la sucursal de origen: %v", err)
	}

	// Obtener información de destino (dirección del cliente)
	var direccion modelos.DirCliente
	if err := db.First(&direccion, despacho.Destino).Error; err != nil {
		return 0, 0, fmt.Errorf("no se encontró la dirección de destino: %v", err)
	}

	// Construir las direcciones para la API
	origen := sucursal.Direccion + ", " + sucursal.Comuna + ", " + sucursal.Ciudad
	destino := direccion.Direccion + ", " + direccion.Comuna + ", " + direccion.Ciudad

	return CalcularDistancia(apiKey, origen, destino)
}
