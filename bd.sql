CREATE SCHEMA IF NOT EXISTS inventario;

SET search_path TO inventario;

CREATE TABLE roles (
    id SERIAL,
    nombre VARCHAR(50) NOT NULL UNIQUE,
    PRIMARY KEY (id)
);

CREATE TABLE usuarios (
    id SERIAL,
    nombre TEXT NOT NULL,
    correo TEXT NOT NULL UNIQUE,
    contrasena TEXT NOT NULL,
    rol_id ,
    PRIMARY KEY (id)
);

CREATE TABLE ubicaciones (
    id SERIAL,
    nombre TEXT NOT NULL UNIQUE,
    tipo VARCHAR(20) NOT NULL CHECK (tipo IN ('bodega', 'tienda')),
    PRIMARY KEY (id)
);

CREATE TABLE clientes (
    id SERIAL,
    nombre TEXT NOT NULL,
    email TEXT,
    telefono VARCHAR(20),
    direccion TEXT,
    PRIMARY KEY (id)
);

CREATE TABLE proveedores (
    id SERIAL,
    nombre TEXT NOT NULL,
    email TEXT,
    telefono VARCHAR(20),
    direccion TEXT,
    PRIMARY KEY (id)
);

CREATE TABLE categorias (
    id SERIAL,
    nombre TEXT NOT NULL UNIQUE,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS productos (
    id SERIAL,
    codigo VARCHAR(50) NOT NULL UNIQUE,
    nombre TEXT NOT NULL,
    descripcion TEXT,
    categoria_id INT NOT NULL,
    proveedor_id INT NOT NULL,
    precio_costo NUMERIC(10,2) NOT NULL,
    precio_venta NUMERIC(10,2) NOT NULL,
    activo BOOLEAN NOT NULL DEFAULT TRUE,
    PRIMARY KEY (id)
    FOREIGN KEY (categoria_id) REFERENCES categorias(id),
    FOREIGN KEY (proveedor_id) REFERENCES proveedores(id)
);

CREATE TABLE inventario (
    producto_id INT NOT NULL,
    ubicacion_id INT NOT NULL,
    cantidad INT NOT NULL DEFAULT 0,
    PRIMARY KEY (producto_id, ubicacion_id),
    FOREIGN KEY (producto_id) REFERENCES productos(id) ON DELETE CASCADE,
    FOREIGN KEY (ubicacion_id) REFERENCES ubicaciones(id) ON DELETE CASCADE
);

CREATE TABLE cotizaciones (
    id SERIAL,
    fecha TIMESTAMP NOT NULL DEFAULT NOW(),
    cliente_id INT NOT NULL,
    vendedor_id INT NOT NULL,
    ubicacion_id INT NOT NULL,
    estado VARCHAR(20) NOT NULL DEFAULT 'Pendiente' CHECK (estado IN ('Pendiente', 'Aprobada', 'Rechazada', 'Vencida')),
    aprobada_por INT NULL,
    fecha_aprobacion TIMESTAMP NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (cliente_id) REFERENCES clientes(id),
    FOREIGN KEY (vendedor_id) REFERENCES usuarios(id),
    FOREIGN KEY (ubicacion_id) REFERENCES ubicaciones(id),
    FOREIGN KEY (aprobada_por) REFERENCES usuarios(id)
);

CREATE TABLE detalle_cotizacion (
    cotizacion_id INT NOT NULL,
    producto_id INT NOT NULL,
    cantidad INT NOT NULL,
    precio_unitario NUMERIC(10,2) NOT NULL,
    PRIMARY KEY (cotizacion_id, producto_id),
    FOREIGN KEY (cotizacion_id) REFERENCES cotizaciones(id) ON DELETE CASCADE,
    FOREIGN KEY (producto_id) REFERENCES productos(id)
);

CREATE TABLE pedidos (
    id SERIAL,
    fecha TIMESTAMP NOT NULL DEFAULT NOW(),
    cliente_id INT NOT NULL,
    vendedor_id INT NOT NULL,
    cotizacion_id INT NULL,
    ubicacion_id INT NOT NULL,
    estado VARCHAR(20) NOT NULL DEFAULT 'Pendiente' CHECK (estado IN ('Pendiente', 'Despachado', 'Completado', 'Cancelado')),
    fecha_despacho TIMESTAMP NULL,
    despachado_por INT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (cliente_id) REFERENCES clientes(id),
    FOREIGN KEY (vendedor_id) REFERENCES usuarios(id),
    FOREIGN KEY (cotizacion_id) REFERENCES cotizaciones(id),
    FOREIGN KEY (ubicacion_id) REFERENCES ubicaciones(id),
    FOREIGN KEY (despachado_por) REFERENCES usuarios(id)
);

CREATE TABLE detalle_pedido (
    pedido_id INT NOT NULL,
    producto_id INT NOT NULL,
    cantidad INT NOT NULL,
    precio_unitario NUMERIC(10,2) NOT NULL,
    PRIMARY KEY (pedido_id, producto_id),
    FOREIGN KEY (pedido_id) REFERENCES pedidos(id) ON DELETE CASCADE,
    FOREIGN KEY (producto_id) REFERENCES productos(id)
);

CREATE TABLE ordenes_compra (
    id SERIAL,
    fecha TIMESTAMP NOT NULL DEFAULT NOW(),
    proveedor_id INT NOT NULL,
    solicitado_por INT NOT NULL,
    ubicacion_id INT NOT NULL,
    estado VARCHAR(20) NOT NULL DEFAULT 'Pendiente' CHECK (estado IN ('Pendiente', 'Recibido', 'Cancelado')),
    fecha_recepcion TIMESTAMP NULL,
    recibido_por INT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (proveedor_id) REFERENCES proveedores(id),
    FOREIGN KEY (solicitado_por) REFERENCES usuarios(id),
    FOREIGN KEY (ubicacion_id) REFERENCES ubicaciones(id),
    FOREIGN KEY (recibido_por) REFERENCES usuarios(id)
);

CREATE TABLE detalle_orden_compra (
    orden_id INT NOT NULL,
    producto_id INT NOT NULL,
    cantidad INT NOT NULL,
    precio_costo NUMERIC(10,2) NOT NULL,
    PRIMARY KEY (orden_id, producto_id),
    FOREIGN KEY (orden_id) REFERENCES ordenes_compra(id) ON DELETE CASCADE,
    FOREIGN KEY (producto_id) REFERENCES productos(id)
);

CREATE TABLE camiones (
    id SERIAL,
    patente VARCHAR(10) NOT NULL UNIQUE,
    marca TEXT,
    modelo TEXT,
    capacidad_kg NUMERIC(10,2),
    activo BOOLEAN DEFAULT TRUE,
    PRIMARY KEY (id)
);

CREATE TABLE despachos (
    id SERIAL,
    pedido_id INT NOT NULL,
    fecha_salida TIMESTAMP NOT NULL DEFAULT NOW(),
    fecha_entrega TIMESTAMP,
    estado VARCHAR(20) NOT NULL DEFAULT 'En ruta' CHECK (estado IN ('Devuelto', 'En ruta', 'Entregado', 'Cancelado')),
    camion_id INT NOT NULL,
    origen_id INT NOT NULL,
    destino TEXT NOT NULL,
    observaciones TEXT,
    PRIMARY KEY (id),
    FOREIGN KEY (pedido_id) REFERENCES pedidos(id),
    FOREIGN KEY (camion_id) REFERENCES camiones(id),
    FOREIGN KEY (origen_id) REFERENCES ubicaciones(id)
);

CREATE TABLE detalle_despacho (
    despacho_id INT NOT NULL,
    producto_id INT NOT NULL,
    cantidad_despachada INT NOT NULL,
    PRIMARY KEY (despacho_id, producto_id),
    FOREIGN KEY (despacho_id) REFERENCES despachos(id) ON DELETE CASCADE,
    FOREIGN KEY (producto_id) REFERENCES productos(id)
);
