# API de Cálculo de Distancias para Despachos

Esta documentación describe las nuevas funcionalidades agregadas al sistema de despachos para calcular distancias usando la API de Google Maps.

## Configuración

Para usar las funcionalidades de cálculo de distancia, necesitas configurar la API key de Google Maps en las variables de entorno:

```bash
GOOGLE_MAPS_API_KEY=tu_api_key_aqui
```

## Endpoints Disponibles

### 1. Obtener todos los despachos con distancia
```
GET /api/despachos-distancia
```

**Respuesta:**
```json
[
  {
    "id": 1,
    "cotizacion_id": 1,
    "camion_id": 1,
    "origen": 1,
    "destino": 1,
    "fecha_despacho": "2025-07-16T10:00:00Z",
    "valor_despacho": 50000,
    "estado": "pendiente",
    "cantidad_items": 10,
    "total_kg": 250.5,
    "total_precio": 125000,
    "distancia_km": 15.3,
    "duracion_minutos": 25,
    "items": [...]
  }
]
```

### 2. Obtener un despacho específico con distancia
```
GET /api/despachos-distancia/:id
```

### 3. Calcular despacho con distancias
```
POST /api/despachos/calcular-distancia
```

**Body:**
```json
{
  "cotizacion_id": 1
}
```

### 4. Obtener despachos detallados por cotización
```
GET /api/despachos/cotizacion/:id/detallado
```

### 5. Calcular distancia entre dos direcciones
```
POST /api/calcular-distancia
```

**Body:**
```json
{
  "origen": "Av. Providencia 1234, Providencia, Santiago",
  "destino": "Av. Las Condes 567, Las Condes, Santiago"
}
```

**Respuesta:**
```json
{
  "distancia_km": 8.5,
  "duracion_minutos": 18,
  "origen": "Av. Providencia 1234, Providencia, Santiago",
  "destino": "Av. Las Condes 567, Las Condes, Santiago"
}
```

## Funciones de Controller Disponibles

### En el código Go:

```go
// Funciones principales
Controllers.GetDespachosConDistancia(db, googleMapsAPIKey)
Controllers.GetDespachoByIDConDistancia(db, id, googleMapsAPIKey)
Controllers.CalcularDespachoConDistancia(db, cotizacionID, googleMapsAPIKey)
Controllers.GetDespachosPorCotizacionDetallado(db, cotizacionID, googleMapsAPIKey)

// Función para calcular distancia directamente
Controllers.CalcularDistancia(apiKey, origen, destino)
Controllers.CalcularDistanciaDespacho(db, apiKey, despacho)
```

## Estructura de Datos

### DespachoConTotales (actualizada)
```go
type DespachoConTotales struct {
    modelos.Despacho
    CantidadItems     int                         `json:"cantidad_items"`
    TotalKg           float64                     `json:"total_kg"`
    TotalPrecio       float64                     `json:"total_precio"`
    DistanciaKm       float64                     `json:"distancia_km,omitempty"`
    DuracionMinutos   int                         `json:"duracion_minutos,omitempty"`
    ProductosDespacho []ProductoDespachoDetallado `json:"items"`
}
```

## Migración desde el código React

El código React original:
```javascript
const getRoute = (e) => {
  e.preventDefault();
  if (!origin || !destination) {
    alert("Debes ingresar origen y destino");
    return;
  }

  const directionsService = new google.maps.DirectionsService();
  directionsService.route({
    origin,
    destination,
    travelMode: google.maps.TravelMode.DRIVING,
  }, (result, status) => {
    if (status === "OK" && result) {
      let totalDistanceMeters = 0;
      result.routes[0].legs.forEach((leg) => {
        totalDistanceMeters += leg.distance?.value || 0;
      });
      const totalDistanceKm = totalDistanceMeters / 1000;
      onRouteResult(result, totalDistanceKm);
    }
  });
};
```

Ahora en Go equivale a:
```go
func CalcularDistancia(apiKey, origen, destino string) (float64, int, error) {
    // Implementación usando Google Maps Distance Matrix API
    // Retorna distancia en km y duración en minutos
}
```

## Ejemplos de Uso

### Desde el Frontend
```javascript
// Obtener despachos con distancia
fetch('/api/despachos-distancia')
  .then(response => response.json())
  .then(data => {
    console.log('Despachos con distancia:', data);
    data.forEach(despacho => {
      console.log(`Despacho ${despacho.id}: ${despacho.distancia_km} km`);
    });
  });

// Calcular distancia específica
fetch('/api/calcular-distancia', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    origen: 'Santiago Centro',
    destino: 'Las Condes'
  })
})
.then(response => response.json())
.then(data => {
  console.log(`Distancia: ${data.distancia_km} km`);
  console.log(`Duración: ${data.duracion_minutos} minutos`);
});
```

## Notas Importantes

1. **API Key**: Es necesario tener una API key válida de Google Maps con acceso a Distance Matrix API
2. **Límites de uso**: Google Maps tiene límites de consultas por día/mes según tu plan
3. **Manejo de errores**: Si la API key no está configurada o hay error en la consulta, las distancias aparecerán como 0
4. **Compatibilidad**: Las funciones originales siguen funcionando sin cambios para mantener compatibilidad

## Variables de Entorno Requeridas

```bash
# .env
GOOGLE_MAPS_API_KEY=AIzaSyC...  # Tu API key de Google Maps
```

## Dependencias Adicionales

El código utiliza las librerías estándar de Go:
- `net/http` para peticiones HTTP
- `encoding/json` para parsear respuestas JSON
- `net/url` para construir URLs

No se requieren librerías externas adicionales.
