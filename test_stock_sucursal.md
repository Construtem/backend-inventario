# Test de API Stock Sucursal

## Cambios realizados

Se ha actualizado la estructura de `stock_sucursal` para usar `sku` en lugar de `producto_id`:

### Modelo actualizado:

```json
{
  "sku": "H001",
  "sucursal_id": 4,
  "cantidad": 50,
  "descuento": 0.0
}
```

## Endpoints actualizados:

### 1. Obtener todo el stock

```
GET /api/stock-sucursal
```

### 2. Obtener stock específico

```
GET /api/stock-sucursal/{sucursal_id}/{sku}
Ejemplo: GET /api/stock-sucursal/4/H001
```

### 3. Crear nuevo registro de stock

```
POST /api/stock-sucursal
Body:
{
  "sku": "H001",
  "sucursal_id": 4,
  "cantidad": 50,
  "descuento": 0.00
}
```

### 4. Actualizar stock existente

```
PUT /api/stock-sucursal/{sucursal_id}/{sku}
Ejemplo: PUT /api/stock-sucursal/4/H001
Body:
{
  "cantidad": 45,
  "descuento": 5.00
}
```

### 5. Eliminar registro de stock

```
DELETE /api/stock-sucursal/{sucursal_id}/{sku}
Ejemplo: DELETE /api/stock-sucursal/4/H001
```

## Ejemplos de uso con JavaScript/TypeScript (Frontend):

### Obtener stock específico:

```javascript
const response = await fetch("http://localhost:8080/api/stock-sucursal/4/H001");
const stock = await response.json();
console.log(stock);
```

### Crear nuevo stock:

```javascript
const nuevoStock = {
  sku: "H001",
  sucursal_id: 4,
  cantidad: 50,
  descuento: 0.0,
};

const response = await fetch("http://localhost:8080/api/stock-sucursal", {
  method: "POST",
  headers: {
    "Content-Type": "application/json",
  },
  body: JSON.stringify(nuevoStock),
});
```

### Actualizar stock:

```javascript
const stockActualizado = {
  cantidad: 45,
  descuento: 5.0,
};

const response = await fetch(
  "http://localhost:8080/api/stock-sucursal/4/H001",
  {
    method: "PUT",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(stockActualizado),
  }
);
```

## Notas importantes:

1. El campo `sku` reemplaza completamente a `producto_id`
2. Las validaciones están activas:
   - `sku` es requerido
   - `sucursal_id` es requerido
   - `cantidad` debe ser >= 0
   - `descuento` debe estar entre 0 y 100
3. Los datos se corresponden con la estructura de tu base de datos
4. Las relaciones con `Producto` y `Sucursal` siguen funcionando correctamente
