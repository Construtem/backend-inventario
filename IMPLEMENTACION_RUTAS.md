# 🚀 Guía de Implementación - Endpoints de Rutas y Distancias

## 📋 Resumen de Implementación

Se han implementado exitosamente los endpoints solicitados en el README para el sistema de rutas y despachos con integración a Google Maps API.

## 🔧 Archivos Creados/Modificados

### Nuevos Archivos
- `services/google_maps.go` - Servicio para integración con Google Maps API
- `api/Handlers/HandlerDistancia.go` - Handlers para endpoints de distancia
- `migrate_distancia.go` - Script de migración de base de datos

### Archivos Modificados
- `api/Models/modelos.go` - Agregadas estructuras para respuestas de distancia
- `api/Controllers/ControllerDespacho.go` - Nuevas funciones para manejo de distancias
- `api/Routes/register.go` - Registradas las nuevas rutas
- `main.go` - Configuración mejorada de CORS
- `.env.example` - Variables de entorno para Google Maps

## 🌐 Endpoints Implementados

### 1. Obtener Despacho por ID con Información de Ruta
```
GET /api/despachos-distancia/{id}
```

**Ejemplo de respuesta:**
```json
{
  "id": 1,
  "cotizacion_id": 123,
  "camion_id": 456,
  "origen": 1,
  "destino": 2,
  "fecha_despacho": "2025-07-16T10:00:00Z",
  "valor_despacho": 25000,
  "cantidad_items": 5,
  "total_kg": 150.5,
  "distancia_calculada": "12.5 km",
  "tiempo_estimado": "25 min",
  "cotizacion": {
    "cliente": {
      "nombre": "Empresa Cliente S.A.",
      "email": "cliente@empresa.com"
    },
    "estado": "En ruta"
  },
  "camion": {
    "patente": "ABC-123"
  },
  "origen_sucursal": {
    "nombre": "Sucursal Central",
    "direccion": "Av. Principal 123",
    "comuna": "Santiago",
    "ciudad": "Santiago"
  },
  "destino_dir_cliente": {
    "direccion": "Calle Cliente 456",
    "comuna": "Las Condes",
    "ciudad": "Santiago"
  }
}
```

### 2. Listar Todos los Despachos con Información de Ruta
```
GET /api/despachos-distancia
```

Retorna un array con la misma estructura del endpoint individual.

### 3. Calcular Distancia Manualmente
```
POST /api/despachos/{id}/calcular-distancia
```

**Cuerpo de la petición:**
```json
{
  "origen": "Av. Principal 123, Santiago",
  "destino": "Calle Cliente 456, Las Condes, Santiago"
}
```

**Respuesta:**
```json
{
  "distancia": "12.5 km",
  "duracion": "25 min",
  "ruta_optimizada": true
}
```

### 4. Calcular Distancia Automáticamente (BONUS)
```
POST /api/despachos/{id}/calcular-distancia-automatico
```

Este endpoint toma automáticamente las direcciones del despacho (sucursal origen y dirección del cliente) para calcular la distancia sin necesidad de enviar direcciones en el cuerpo de la petición.

## ⚙️ Configuración Requerida

### 1. Variables de Entorno
Agregar estas variables a tu archivo `.env`:

```env
# Google Maps API Configuration
GOOGLE_MAPS_API_KEY=tu_api_key_de_google_maps
GOOGLE_MAPS_DISTANCE_API_URL=https://maps.googleapis.com/maps/api/distancematrix/json

# CORS Configuration para desarrollo
ALLOWED_ORIGINS=http://localhost:3000
```

### 2. Migración de Base de Datos
Ejecutar el script de migración para agregar las nuevas columnas:

```bash
go run migrate_distancia.go
```

Esto agregará las columnas:
- `distancia_calculada` (VARCHAR(50))
- `tiempo_estimado` (VARCHAR(50))

### 3. Obtener API Key de Google Maps

1. Ve a [Google Cloud Console](https://console.cloud.google.com/)
2. Crea un nuevo proyecto o selecciona uno existente
3. Habilita la API "Distance Matrix API"
4. Crea credenciales (API Key)
5. Configura restricciones si es necesario
6. Agrega la API key a tu archivo `.env`

## 🧪 Testing

### 1. Testing Local
Asegúrate de que el servidor esté corriendo en `http://localhost:8080`

### 2. Datos de Prueba
Puedes usar estos endpoints para testing:

```bash
# Obtener un despacho específico
curl -X GET "http://localhost:8080/api/despachos-distancia/1"

# Obtener todos los despachos
curl -X GET "http://localhost:8080/api/despachos-distancia"

# Calcular distancia manualmente
curl -X POST "http://localhost:8080/api/despachos/1/calcular-distancia" \
  -H "Content-Type: application/json" \
  -d '{
    "origen": "Dieciocho 161, Santiago, Región Metropolitana",
    "destino": "Av. José Pedro Alessandri 1242, Ñuñoa, Región Metropolitana"
  }'

# Calcular distancia automáticamente
curl -X POST "http://localhost:8080/api/despachos/1/calcular-distancia-automatico"
```

## 🔍 Verificación de Funcionamiento

### 1. Base de Datos
Verifica que las nuevas columnas existan:
```sql
DESCRIBE despacho;
```

### 2. Endpoints
Verifica que los endpoints respondan correctamente:
- Status 200 para GET requests exitosos
- Status 404 para despachos no encontrados
- Status 400 para IDs inválidos o datos faltantes
- Status 500 para errores de Google Maps API

### 3. Google Maps Integration
El servicio maneja automáticamente:
- ✅ Validación de direcciones
- ✅ Formato de direcciones chilenas
- ✅ Manejo de errores de la API
- ✅ Timeout de 30 segundos
- ✅ Respuestas en español

## 🚨 Manejo de Errores

### Errores Comunes y Soluciones

1. **Error: "API key inválida"**
   - Verifica que `GOOGLE_MAPS_API_KEY` esté configurada correctamente
   - Asegúrate de que la API key tenga permisos para Distance Matrix API

2. **Error: "Dirección no encontrada"**
   - Verifica que las direcciones estén bien formateadas
   - Incluye comuna y ciudad para mejor precisión

3. **Error: "CORS"**
   - Verifica que `ALLOWED_ORIGINS` incluya la URL de tu frontend
   - Para desarrollo local usa: `http://localhost:3000`

4. **Error: "Despacho no encontrado"**
   - Verifica que el ID del despacho exista en la base de datos
   - Asegúrate de que el despacho tenga direcciones de origen y destino

## 📊 Monitoreo y Logs

El sistema incluye logs detallados para:
- Peticiones a Google Maps API
- Errores de conexión
- Validación de datos
- Actualizaciones de base de datos

## 🔄 Próximos Pasos

1. **Cache de Distancias**: Implementar cache para evitar recalcular distancias repetidas
2. **Rate Limiting**: Implementar límites para evitar exceder cuotas de Google Maps
3. **Batch Processing**: Calcular distancias en lote para múltiples despachos
4. **Analytics**: Agregar métricas de uso de la API

## 🎯 Integración con Frontend

El backend ahora está completamente preparado para recibir requests del frontend en:
- `http://localhost:8080/api/despachos-distancia/{id}`
- `http://localhost:8080/api/despachos-distancia`

La estructura de respuesta coincide exactamente con lo especificado en el README original.

---

**Estado**: ✅ Implementación Completa  
**Compatibilidad**: Frontend Next.js 15.3.4 + TypeScript  
**API Externa**: Google Maps Distance Matrix API  
**Base de Datos**: PostgreSQL con GORM
