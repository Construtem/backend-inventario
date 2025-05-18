CREATE SCHEMA IF NOT EXISTS inventario;

SET search_path TO inventario;

CREATE TABLE IF NOT EXISTS sucursal (
    id_sucursal SERIAL,
    nombre VARCHAR(50) NOT NULL,
    direccion VARCHAR(255) NOT NULL,
    telefono VARCHAR(20) NOT NULL,
    PRIMARY KEY (id_sucursal)
);

CREATE TABLE IF NOT EXISTS producto (
    sku VARCHAR(10) NOT NULL,
    nombre VARCHAR(50) NOT NULL,
    descripcion VARCHAR(255),
    peso_kg DECIMAL(5,2) NOT NULL,
    dim_largo_cm DECIMAL(5,2) NOT NULL,
    dim_ancho_cm DECIMAL(5,2) NOT NULL,
    dim_alto_cm DECIMAL(5,2) NOT NULL,
    precio_unitario DECIMAL(10,2) NOT NULL,
    PRIMARY KEY (sku)
);

CREATE TABLE IF NOT EXISTS bodega (
    id_bodega SERIAL,
    nombre VARCHAR(50) NOT NULL,
    direccion VARCHAR(255) NOT NULL,
    telefono VARCHAR(20) NOT NULL,
    PRIMARY KEY (id_bodega)
);

CREATE TABLE IF NOT EXISTS stock_bodega (
    id_stock_bodega SERIAL,
    id_bodega INT NOT NULL,
    sku VARCHAR(10) NOT NULL,
    cantidad INT NOT NULL,
    PRIMARY KEY (id_stock_bodega),
    FOREIGN KEY (id_bodega) REFERENCES bodega(id_bodega),
    FOREIGN KEY (sku) REFERENCES producto(sku)
);

CREATE TABLE IF NOT EXISTS orden_pedido (
    id_orden_pedido SERIAL,
    fecha_crea TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    id_sucursal INT NOT NULL,
    id_bodega INT NOT NULL,
    estado VARCHAR(20) NOT NULL,
    PRIMARY KEY (id_orden_pedido),
    FOREIGN KEY (id_sucursal) REFERENCES sucursal(id_sucursal),
    FOREIGN KEY (id_bodega) REFERENCES bodega(id_bodega)
);

CREATE TABLE IF NOT EXISTS orden_pedido_item (
    id_pedido_item SERIAL,
    id_orden_pedido INT NOT NULL,
    sku VARCHAR(10) NOT NULL,
    cantidad INT NOT NULL,
    PRIMARY KEY (id_pedido_item),
    FOREIGN KEY (id_orden_pedido) REFERENCES orden_pedido(id_orden_pedido),
    FOREIGN KEY (sku) REFERENCES producto(sku)
);

CREATE TABLE IF NOT EXISTS tipo_camion (
    id_tipo_camion SERIAL,
    nombre VARCHAR(50) NOT NULL,
    descripcion VARCHAR(255),
    PRIMARY KEY (id_tipo_camion)
);

CREATE TABLE IF NOT EXISTS estado_camion(
    id_estado_camion SERIAL,
    nombre VARCHAR(20) NOT NULL,
    PRIMARY KEY (id_estado_camion)
);

CREATE TABLE IF NOT EXISTS camion (
    id_camion SERIAL,
    patente VARCHAR(10) NOT NULL UNIQUE,
    id_tipo_camion INT NOT NULL,
    id_estado_camion INT NOT NULL,
    capacidad_kg DECIMAL(10,2) NOT NULL,
    dim_largo_cm DECIMAL(5,2) NOT NULL,
    dim_ancho_cm DECIMAL(5,2) NOT NULL,
    dim_alto_cm DECIMAL(5,2) NOT NULL,
    PRIMARY KEY (id_camion),
    FOREIGN KEY (id_tipo_camion) REFERENCES tipo_camion(id_tipo_camion),
    FOREIGN KEY (id_estado_camion) REFERENCES estado_camion(id_estado_camion)
);

CREATE TABLE IF NOT EXISTS estado_envio (
    id_estado_envio SERIAL,
    nombre VARCHAR(20) NOT NULL,
    PRIMARY KEY (id_estado_envio)
);

CREATE TABLE IF NOT EXISTS envio (
    id_envio SERIAL,
    id_orden_pedido INT NOT NULL,
    id_camion INT NOT NULL,
    id_estado_envio INT NOT NULL,
    fecha_envio DATE NOT NULL,
    PRIMARY KEY (id_envio),
    FOREIGN KEY (id_orden_pedido) REFERENCES orden_pedido(id_orden_pedido),
    FOREIGN KEY (id_camion) REFERENCES camion(id_camion),
    FOREIGN KEY (id_estado_envio) REFERENCES estado_envio(id_estado_envio)
);