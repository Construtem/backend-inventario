# Script SQL para agregar columnas de distancia y tiempo

```sql
-- Ejecutar estos comandos directamente en tu base de datos PostgreSQL
-- para agregar las nuevas columnas al modelo Despacho

ALTER TABLE despacho ADD COLUMN IF NOT EXISTS distancia_calculada VARCHAR(50);
ALTER TABLE despacho ADD COLUMN IF NOT EXISTS tiempo_estimado VARCHAR(50);

-- Verificar que las columnas se agregaron correctamente
SELECT column_name, data_type, character_maximum_length
FROM information_schema.columns 
WHERE table_name = 'despacho' 
AND column_name IN ('distancia_calculada', 'tiempo_estimado');
```

## Comandos alternativos usando psql:

```bash
# Conectarse a la base de datos y ejecutar los comandos
psql -h localhost -p 5432 -U tu_usuario -d tu_base_de_datos

# Dentro de psql, ejecutar:
ALTER TABLE despacho ADD COLUMN IF NOT EXISTS distancia_calculada VARCHAR(50);
ALTER TABLE despacho ADD COLUMN IF NOT EXISTS tiempo_estimado VARCHAR(50);
```

## Usando el CLI de tu gestor de base de datos favorito:
- pgAdmin
- DBeaver  
- DataGrip
- etc.

Simplemente ejecuta los dos comandos ALTER TABLE en tu herramienta de administración de base de datos preferida.
