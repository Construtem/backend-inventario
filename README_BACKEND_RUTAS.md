# 📍 Integración Backend - Sistema de Rutas y Despachos

## 🎯 Objetivo
Implementar endpoints en el backend para soportar el sistema de rutas de despacho que actualmente usa Google Maps API en el frontend.

## 📋 Estado Actual del Frontend

### Página de Rutas: `/admin/despacho/ruta`
- **Archivo principal**: `src/app/admin/despacho/ruta/page.tsx`
- **Componente de mapa**: `src/app/admin/despacho/ruta/MapView/`
- **API Frontend**: Configurada para `http://localhost:8080`

### Funcionalidades Implementadas
- ✅ Visualización de rutas en Google Maps
- ✅ Cálculo automático de distancia y tiempo con Google Maps
- ✅ Interfaz responsive con información del despacho
- ✅ Manejo de errores y datos de respaldo

## 🔗 Endpoints Necesarios

### 1. Obtener Despacho por ID
```
GET /api/despachos-distancia/{id}
```

**Respuesta esperada:**
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

### 2. Listar Todos los Despachos
```
GET /api/despachos-distancia
```

**Respuesta esperada:**
```json
[
  {
    "id": 1,
    // ... mismo formato que el endpoint individual
  },
  {
    "id": 2,
    // ... siguiente despacho
  }
]
```

### 3. Calcular Distancia (Opcional - Integración con Google Maps)
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

**Respuesta esperada:**
```json
{
  "distancia": "12.5 km",
  "duracion": "25 min",
  "ruta_optimizada": true
}
```

## 🗂️ Estructura de Datos Frontend

### Interface DespachoInfo (Frontend)
```typescript
interface DespachoInfo {
  id: string;
  cliente: string;
  origen: string;
  destino: string;
  estado: string;
  distancia?: string;
  duracion?: string;
}
```

### Interface DespachoBackend (Esperada del Backend)
```typescript
interface DespachoBackend {
  id: number;
  cotizacion_id: number;
  camion_id: number;
  origen: number;
  destino: number;
  fecha_despacho: string;
  valor_despacho: number;
  cantidad_items: number;
  total_kg: number;
  distancia_calculada?: string;
  tiempo_estimado?: string;
  cotizacion?: {
    cliente?: {
      nombre: string;
      email: string;
    };
    estado?: string;
  };
  camion?: {
    patente: string;
  };
  origen_sucursal?: {
    nombre: string;
    direccion?: string;
    comuna?: string;
    ciudad?: string;
  };
  destino_dir_cliente?: {
    direccion: string;
    comuna: string;
    ciudad: string;
  };
}
```

## 🔧 Configuración Requerida

### Variables de Entorno Backend
```env
# Google Maps API
GOOGLE_MAPS_API_KEY=your_google_maps_api_key
GOOGLE_MAPS_DISTANCE_API_URL=https://maps.googleapis.com/maps/api/distancematrix/json

# Base de datos
DATABASE_URL=your_database_connection

# CORS (para desarrollo)
ALLOWED_ORIGINS=http://localhost:3000
```

### Dependencias Sugeridas
- **Google Maps API Client** (para calcular distancias)
- **Axios o Fetch** (para llamadas HTTP a Google Maps)
- **CORS middleware** (para desarrollo local)

## 📊 Flujo de Datos

### 1. Frontend → Backend
```
Frontend (página de ruta) 
    ↓ GET /api/despachos-distancia/{id}
Backend (retorna datos del despacho)
    ↓ 
Frontend (muestra información + mapa de Google)
```

### 2. Cálculo de Distancia (Opcional)
```
Backend → Google Maps Distance Matrix API
    ↓ (obtiene distancia y tiempo)
Backend → Guarda en base de datos
    ↓ 
Frontend → Muestra datos calculados
```

## 🛠️ Implementación Sugerida

### 1. Controller de Despachos
```java
@RestController
@RequestMapping("/api")
public class DespachoController {
    
    @GetMapping("/despachos-distancia")
    public ResponseEntity<List<DespachoDTO>> getAllDespachos() {
        // Implementar lógica para obtener todos los despachos
    }
    
    @GetMapping("/despachos-distancia/{id}")
    public ResponseEntity<DespachoDTO> getDespachoById(@PathVariable Long id) {
        // Implementar lógica para obtener despacho específico
    }
    
    @PostMapping("/despachos/{id}/calcular-distancia")
    public ResponseEntity<DistanciaDTO> calcularDistancia(
        @PathVariable Long id, 
        @RequestBody DireccionesDTO direcciones
    ) {
        // Implementar cálculo con Google Maps API
    }
}
```

### 2. Servicio de Google Maps
```java
@Service
public class GoogleMapsService {
    
    @Value("${google.maps.api.key}")
    private String apiKey;
    
    public DistanciaDTO calcularDistancia(String origen, String destino) {
        // Implementar llamada a Google Maps Distance Matrix API
        // Retornar distancia y tiempo estimado
    }
}
```

### 3. DTO Responses
```java
public class DespachoDTO {
    private Long id;
    private Long cotizacionId;
    private Long camionId;
    // ... otros campos según la estructura mostrada arriba
}

public class DistanciaDTO {
    private String distancia;
    private String duracion;
    private Boolean rutaOptimizada;
    // getters y setters
}
```

## 🔍 Testing

### Datos de Prueba
El frontend está configurado con datos de ejemplo para testing:

```json
{
  "id": "1",
  "cliente": "Universidad Tecnologica Metropolitana del Estado de Chile",
  "origen": "Dieciocho 161, 8330383 Santiago, Región Metropolitana",
  "destino": "Av. José Pedro Alessandri 1242, Ñuñoa, Región Metropolitana",
  "estado": "En ruta",
  "distancia": "6.2 km",
  "duracion": "20 min"
}
```

### URLs de Testing
- **Frontend Development**: `http://localhost:3000/admin/despacho/ruta?despachoId=1`
- **Backend Expected**: `http://localhost:8080/api/despachos-distancia/1`

## 🚨 Consideraciones Importantes

### Manejo de Errores
- El frontend ya maneja errores de conexión
- Retorna datos de ejemplo si el backend no está disponible
- Logs detallados para debugging

### CORS
Configurar CORS para permitir requests desde:
```
http://localhost:3000 (desarrollo)
https://tu-dominio.com (producción)
```

### Rate Limiting
Google Maps API tiene límites de requests - implementar cache para distancias calculadas.

### Seguridad
- Validar IDs de despacho
- Sanitizar direcciones antes de enviar a Google Maps
- Proteger API key de Google Maps (solo en backend)

## 📞 Contacto

Si tienes dudas sobre la implementación o necesitas más detalles sobre algún endpoint específico, contacta al equipo de frontend.

---

**Desarrollado para**: Sistema de Inventario - Módulo de Despachos  
**Tecnología Frontend**: Next.js 15.3.4 + TypeScript  
**Integración**: Google Maps JavaScript API  
**Estado**: Listo para integración backend
