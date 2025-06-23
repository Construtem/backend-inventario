package modelos

import "time"

type Cliente struct {
	ID        uint   `gorm:"primaryKey"`
	Nombre    string `gorm:"size:255;not null"`
	Email     string `gorm:"size:255;unique;not null"`
	Telefono  string `gorm:"size:20"`
	Direccion string `gorm:"size:255"`
}
type Producto struct {
	ID          uint   `gorm:"primaryKey"`
	Codigo      string `gorm:"size:50;unique;not null"`
	Descripcion string `gorm:"size:255;not null"`
	CategoriaID uint
	Categoria   Categoria
	PrecioCosto float64 `gorm:"type:numeric(10,2);not null"`
	PrecioVenta float64 `gorm:"type:numeric(10,2);not null"`
	Activo      bool    `gorm:"default:true"`
}
type Categoria struct {
	ID     uint   `gorm:"primaryKey"`
	Nombre string `gorm:"size:100;unique;not null"`
}

type Camion struct {
	ID          uint    `gorm:"primaryKey"`
	Patente     string  `gorm:"size:20;unique;not null"`
	Marca       string  `gorm:"size:100"`
	Modelo      string  `gorm:"size:100"`
	CapacidadKg float64 `gorm:"type:numeric(10,2)"`
	Activo      bool    `gorm:"default:true"`
}
type Pedido struct {
	ID            uint      `gorm:"primaryKey"`
	Fecha         time.Time `gorm:"not null"`
	ClienteID     uint
	Cliente       Cliente `gorm:"foreignKey:ClienteID"` // Relación con Cliente
	VendedorID    string
	Vendedor      Usuario     `gorm:"foreignKey:VendedorID;references:UID"` // Relación con Usuario (Vendedor)
	CotizacionID  *uint       // Puede ser nulo
	Cotizacion    *Cotizacion `gorm:"foreignKey:CotizacionID"` // Relación con Cotizacion (puede ser nula)
	UbicacionID   uint
	Ubicacion     Ubicacion `gorm:"foreignKey:UbicacionID"` // Relación con Ubicacion
	Estado        string    `gorm:"size:50;not null"`
	FechaDespacho *time.Time
	DespachadoPor *string
	Despachador   *Usuario `gorm:"foreignKey:DespachadoPor;references:UID"` // Relación con Usuario (Despachador, puede ser nulo)
}
type Cotizacion struct {
	ID              uint      `gorm:"primaryKey"`
	Fecha           time.Time `gorm:"not null"`
	ClienteID       uint
	Cliente         Cliente `gorm:"foreignKey:ClienteID"` // Relación con Cliente
	VendedorID      string
	Vendedor        Usuario `gorm:"foreignKey:VendedorID;references:UID"` // Relación con Usuario (Vendedor)
	UbicacionID     uint
	Ubicacion       Ubicacion `gorm:"foreignKey:UbicacionID"` // Relación con Ubicacion
	Estado          string    `gorm:"size:50;not null"`
	FechaAprobacion *time.Time
}
type Usuario struct {
	UID    string `gorm:"primaryKey;size:100;not null;unique" json:"uid"`
	Nombre string `gorm:"size:255;not null"`
	RolID  uint
	Rol    Rol `gorm:"foreignKey:RolID"` // Relación con Rol
}
type Rol struct {
	ID     uint   `gorm:"primaryKey"`
	Nombre string `gorm:"size:50;unique;not null"`
}
type Ubicacion struct {
	ID     uint   `gorm:"primaryKey"`
	Nombre string `gorm:"size:255;not null"`
	Tipo   string `gorm:"size:50"` // Ej: 'bodega', 'tienda', 'cliente'
}
type Proveedor struct {
	ID        uint   `gorm:"primaryKey"`
	Nombre    string `gorm:"size:255;not null"`
	Email     string `gorm:"size:255;unique"`
	Telefono  string `gorm:"size:20"`
	Direccion string `gorm:"size:255"`
}
type Inventario struct {
	ID          uint `gorm:"primaryKey"`
	ProductoID  uint
	Producto    Producto `gorm:"foreignKey:ProductoID"` // Relación con Producto
	UbicacionID uint
	Ubicacion   Ubicacion `gorm:"foreignKey:UbicacionID"` // Relación con Ubicacion
	Cantidad    int       `gorm:"not null"`
}

type Despacho struct {
	ID            uint `gorm:"primaryKey"`
	PedidoID      uint
	Pedido        Pedido    `gorm:"foreignKey:PedidoID"` // Relación con Pedido
	FechaSalida   time.Time `gorm:"not null"`
	FechaEntrega  *time.Time
	Estado        string `gorm:"size:50;not null"` // Ej: 'pendiente', 'en_ruta', 'entregado', 'cancelado'
	CamionID      uint
	Camion        Camion `gorm:"foreignKey:CamionID"` // Relación con Camion
	OrigenID      uint
	Origen        Ubicacion `gorm:"foreignKey:OrigenID"` // Relación con Ubicacion (Origen)
	DestinoID     uint
	Destino       Ubicacion `gorm:"foreignKey:DestinoID"` // Relación con Ubicacion (Destino)
	Observaciones string    `gorm:"size:500"`
}

type OrdenCompra struct {
	ID             uint      `gorm:"primaryKey"`
	Fecha          time.Time `gorm:"not null"`
	ProveedorID    uint
	Proveedor      Proveedor `gorm:"foreignKey:ProveedorID"` // Relación con Proveedor
	SolicitadoPor  string
	Solicitante    Usuario `gorm:"foreignKey:SolicitadoPor;references:UID"` // Relación con Usuario (Solicitante)
	UbicacionID    uint
	Ubicacion      Ubicacion `gorm:"foreignKey:UbicacionID"` // Relación con Ubicacion (donde se recibirá)
	Estado         string    `gorm:"size:50;not null"`       // Ej: 'pendiente', 'aprobada', 'recibida', 'cancelada'
	FechaRecepcion *time.Time
	RecibidoPor    *string
	Receptor       *Usuario `gorm:"foreignKey:RecibidoPor;references:UID"` // Relación con Usuario (Receptor, puede ser nulo)
}

type DetalleCotizacion struct {
	ID             uint `gorm:"primaryKey"`
	CotizacionID   uint
	Cotizacion     Cotizacion `gorm:"foreignKey:CotizacionID"` // Relación con Cotizacion
	ProductoID     uint
	Producto       Producto `gorm:"foreignKey:ProductoID"` // Relación con Producto
	Cantidad       int      `gorm:"not null"`
	PrecioUnitario float64  `gorm:"type:numeric(10,2);not null"`
}
type DetalleOrdenCompra struct {
	ID          uint `gorm:"primaryKey"`
	OrdenID     uint
	OrdenCompra OrdenCompra `gorm:"foreignKey:OrdenID"` // Relación con OrdenCompra
	ProductoID  uint
	Producto    Producto `gorm:"foreignKey:ProductoID"` // Relación con Producto
	Cantidad    int      `gorm:"not null"`
	PrecioCosto float64  `gorm:"type:numeric(10,2);not null"`
}

type DetallePedido struct {
	ID             uint `gorm:"primaryKey"`
	PedidoID       uint
	Pedido         Pedido `gorm:"foreignKey:PedidoID"` // Relación con Pedido
	ProductoID     uint
	Producto       Producto `gorm:"foreignKey:ProductoID"` // Relación con Producto
	Cantidad       int      `gorm:"not null"`
	PrecioUnitario float64  `gorm:"type:numeric(10,2);not null"`
}
type DetalleDespacho struct {
	ID                 uint `gorm:"primaryKey"`
	DespachoID         uint
	Despacho           Despacho `gorm:"foreignKey:DespachoID"` // Relación con Despacho
	ProductoID         uint
	Producto           Producto `gorm:"foreignKey:ProductoID"` // Relación con Producto
	CantidadDespachada int      `gorm:"not null"`
}
