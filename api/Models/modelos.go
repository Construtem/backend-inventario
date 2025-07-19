package modelos

import "time"

type Producto struct {
	SKU         string  `gorm:"primaryKey;size:20" json:"sku"`
	Nombre      string  `gorm:"size:100;not null" json:"nombre"`
	Descripcion string  `gorm:"type:text" json:"descripcion"`
	ProveedorID uint    `gorm:"not null" json:"proveedor_id"`
	Peso        float64 `gorm:"type:numeric(10,2);not null" json:"peso"`
	Largo       float64 `gorm:"type:numeric(10,2);not null" json:"largo"`
	Ancho       float64 `gorm:"type:numeric(10,2);not null" json:"ancho"`
	Alto        float64 `gorm:"type:numeric(10,2);not null" json:"alto"`
	Precio      float64 `gorm:"type:numeric(10,2);not null" json:"precio"`
	CategoriaID *uint   `gorm:"column:categoria_id" json:"categoria_id"`
	Estado      bool    `gorm:"default:true" json:"estado"`

	Proveedor Proveedor `gorm:"foreignKey:ProveedorID;references:ID;constraint:OnDelete:CASCADE" json:"proveedor"`
	Categoria Categoria `gorm:"foreignKey:CategoriaID;references:ID;constraint:OnDelete:SET NULL" json:"categoria,omitempty"`
}

func (Producto) TableName() string {
	return "productos"
}

type Categoria struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	Nombre string `gorm:"size:100;not null;unique" json:"nombre"`
}

type Proveedor struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	Marca     string `gorm:"size:30;not null" json:"marca"`
	Email     string `gorm:"size:100;not null;unique" json:"email"`
	Telefono  string `gorm:"size:20;not null" json:"telefono"`
	Direccion string `gorm:"size:200;not null" json:"direccion"`
}

func (Proveedor) TableName() string {
	return "proveedores"
}

type StockProveedor struct {
	ProveedorID  uint      `gorm:"primaryKey;column:proveedor_id" json:"proveedor_id"`
	ProductoID   string    `gorm:"primaryKey;column:sku" json:"sku"`
	Stock        int       `gorm:"not null" json:"stock"`
	FechaIngreso time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"fecha_ingreso"`

	Proveedor Proveedor `gorm:"foreignKey:ProveedorID;references:ID;constraint:OnDelete:CASCADE" json:"proveedor"`
	Producto  Producto  `gorm:"foreignKey:ProductoID;references:SKU;constraint:OnDelete:CASCADE" json:"producto"`
}

func (StockProveedor) TableName() string {
	return "stock_proveedor"
}

type TipoSucursal struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	Nombre string `gorm:"size:50;not null" json:"nombre"`
}

func (TipoSucursal) TableName() string {
	return "tipo_sucursal"
}

type Sucursal struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	Nombre    string `gorm:"size:100;not null" json:"nombre"`
	Telefono  string `gorm:"size:20;not null" json:"telefono"`
	Direccion string `gorm:"size:200;not null" json:"direccion"`
	Comuna    string `gorm:"size:100;not null" json:"comuna"`
	Ciudad    string `gorm:"size:100;not null" json:"ciudad"`
	TipoID    uint   `gorm:"column:tipo_id;not null" json:"tipo_id"`

	Tipo TipoSucursal `gorm:"foreignKey:TipoID;references:ID;constraint:OnDelete:CASCADE" json:"tipo"`
}

func (Sucursal) TableName() string {
	return "sucursales"
}

type StockSucursal struct {
	SKU        string  `gorm:"primaryKey;size:20;column:sku" json:"sku" binding:"required"`
	SucursalID uint    `gorm:"primaryKey;column:sucursal_id" json:"sucursal_id" binding:"required"`
	Cantidad   int     `gorm:"not null" json:"cantidad" binding:"required,min=0"`
	Descuento  float64 `gorm:"type:numeric(5,2);default:0;check:descuento >= 0 AND descuento <= 100" json:"descuento" binding:"min=0,max=100"`

	Producto Producto `gorm:"foreignKey:SKU;references:SKU;constraint:OnDelete:CASCADE" json:"producto,omitempty"`
	Sucursal Sucursal `gorm:"foreignKey:SucursalID;references:ID;constraint:OnDelete:CASCADE" json:"sucursal,omitempty"`
}

func (StockSucursal) TableName() string {
	return "stock_sucursal"
}

type Rol struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	Nombre string `gorm:"size:50;not null" json:"nombre"`
}

func (Rol) TableName() string {
	return "roles"
}

type Usuario struct {
	Email  string `gorm:"primaryKey;size:100;not null;unique" json:"email"`
	Nombre string `gorm:"size:50;not null" json:"nombre"`
	RolID  uint   `gorm:"column:rol_id;not null" json:"rol_id"`

	Rol Rol `gorm:"foreignKey:RolID;references:ID;constraint:OnDelete:CASCADE" json:"rol"`
}

type TipoCliente struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	Nombre string `gorm:"size:50;not null" json:"nombre"`
}

func (TipoCliente) TableName() string {
	return "tipo_cliente"
}

type Cliente struct {
	Rut         string `gorm:"primaryKey;size:12;not null;unique" json:"rut"`
	Nombre      string `gorm:"size:100;not null" json:"nombre"`
	Telefono    string `gorm:"size:20" json:"telefono"`
	Email       string `gorm:"size:100;not null;unique" json:"email"`
	RazonSocial string `gorm:"size:100" json:"razon_social"`
	TipoID      uint   `gorm:"column:tipo_id;not null" json:"tipo_id"`

	Tipo TipoCliente `gorm:"foreignKey:TipoID;references:ID;constraint:OnDelete:CASCADE" json:"tipo"`
}

type DirCliente struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	RutCliente string `gorm:"column:rut_cliente;not null" json:"rut_cliente"`
	Nombre     string `gorm:"size:100;not null" json:"nombre"`
	Direccion  string `gorm:"size:200;not null" json:"direccion"`
	Comuna     string `gorm:"size:100;not null" json:"comuna"`
	Ciudad     string `gorm:"size:100;not null" json:"ciudad"`

	Cliente Cliente `gorm:"foreignKey:RutCliente;references:Rut;constraint:OnDelete:CASCADE" json:"cliente"`
}

func (DirCliente) TableName() string {
	return "dir_cliente"
}

type Cotizacion struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	FechaCrea    time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"fecha_crea"`
	Estado       string    `gorm:"size:20;not null" json:"estado"`
	CostoEnvio   float64   `gorm:"type:numeric(10,2);not null" json:"costo_envio"`
	RutCliente   string    `gorm:"column:rut_cliente;not null" json:"rut_cliente"`
	UserID       string    `gorm:"size:100;column:user_id;not null" json:"user_id"`
	TipoDespacho string    `gorm:"size:50;not null" json:"tipo_despacho"`

	Cliente Cliente `gorm:"foreignKey:RutCliente;references:Rut;constraint:OnDelete:CASCADE" json:"cliente"`
	Usuario Usuario `gorm:"foreignKey:UserID;references:Email;constraint:OnDelete:CASCADE" json:"usuario"`
}

func (Cotizacion) TableName() string {
	return "cotizaciones"
}

type CotizacionItem struct {
	CotizacionID uint   `gorm:"primaryKey;column:cotizacion_id" json:"cotizacion_id"`
	ProductoID   string `gorm:"primaryKey;size:20;column:producto_id" json:"producto_id"`
	SucursalID   uint   `gorm:"primaryKey;column:sucursal_id" json:"sucursal_id"`
	Cantidad     int    `gorm:"not null" json:"cantidad"`

	Cotizacion Cotizacion `gorm:"foreignKey:CotizacionID;references:ID;constraint:OnDelete:CASCADE" json:"cotizacion"`
	Producto   Producto   `gorm:"foreignKey:ProductoID;references:SKU;constraint:OnDelete:CASCADE" json:"producto"`
	Sucursal   Sucursal   `gorm:"foreignKey:SucursalID;references:ID;constraint:OnDelete:CASCADE" json:"sucursal"`
}

func (CotizacionItem) TableName() string {
	return "cotizacion_item"
}

type TipoCamion struct {
	ID         uint    `gorm:"primaryKey" json:"id"`
	Volumen    float64 `gorm:"type:numeric(10,2);not null" json:"volumen"`
	PesoMaximo float64 `gorm:"type:numeric(10,2);not null" json:"peso_maximo"`
}

func (TipoCamion) TableName() string {
	return "tipo_camion"
}

type Camion struct {
	ID      uint   `gorm:"primaryKey" json:"id"`
	Patente string `gorm:"size:20;not null;unique" json:"patente"`
	TipoID  uint   `gorm:"column:tipo_id;not null" json:"tipo_id"`
	Activo  bool   `gorm:"default:true" json:"activo"`

	Tipo TipoCamion `gorm:"foreignKey:TipoID;references:ID;constraint:OnDelete:CASCADE" json:"tipo"`
}

func (Camion) TableName() string {
	return "camiones"
}

type Despacho struct {
	ID                 uint      `gorm:"primaryKey" json:"id"`
	CotizacionID       uint      `gorm:"column:cotizacion_id;not null" json:"cotizacion_id"`
	CamionID           uint      `gorm:"column:camion_id;not null" json:"camion_id"`
	Origen             uint      `gorm:"not null" json:"origen"`  // FK a sucursales.id
	Destino            uint      `gorm:"not null" json:"destino"` // FK a dir_cliente.id
	FechaDespacho      time.Time `gorm:"not null" json:"fecha_despacho"`
	ValorDespacho      float64   `gorm:"type:numeric(10,2);not null" json:"valor_despacho"`
	Estado             string    `gorm:"size:20;not null;default:'pendiente'" json:"estado"`
	DistanciaCalculada *string   `gorm:"size:50" json:"distancia_calculada,omitempty"`
	TiempoEstimado     *string   `gorm:"size:50" json:"tiempo_estimado,omitempty"`

	Cotizacion        Cotizacion          `gorm:"foreignKey:CotizacionID;references:ID;constraint:OnDelete:CASCADE" json:"cotizacion"`
	Camion            Camion              `gorm:"foreignKey:CamionID;references:ID;constraint:OnDelete:CASCADE" json:"camion"`
	OrigenSucursal    Sucursal            `gorm:"foreignKey:Origen;references:ID;constraint:OnDelete:CASCADE" json:"origen_sucursal"`
	DestinoDirCliente DirCliente          `gorm:"foreignKey:Destino;references:ID;constraint:OnDelete:CASCADE" json:"destino_dir_cliente"`
	ProductosDespacho []ProductosDespacho `gorm:"foreignKey:DespachoID;references:ID;constraint:OnDelete:CASCADE" json:"productos"`
}

func (Despacho) TableName() string {
	return "despacho"
}

type ProductosDespacho struct {
	DespachoID uint   `gorm:"primaryKey;column:despacho_id" json:"despacho_id"`
	ProductoID string `gorm:"primaryKey;size:20;column:sku" json:"producto_id"`
	Cantidad   int    `gorm:"not null" json:"cantidad"`

	Despacho Despacho `gorm:"foreignKey:DespachoID;references:ID;constraint:OnDelete:CASCADE" json:"despacho"`
	Producto Producto `gorm:"foreignKey:ProductoID;references:SKU;constraint:OnDelete:CASCADE" json:"producto"`
}

func (ProductosDespacho) TableName() string {
	return "productos_despacho"
}

// DespachoDistanciaResponse es la estructura de respuesta para los endpoints de rutas
type DespachoDistanciaResponse struct {
	ID                 uint                         `json:"id"`
	CotizacionID       uint                         `json:"cotizacion_id"`
	CamionID           uint                         `json:"camion_id"`
	Origen             uint                         `json:"origen"`
	Destino            uint                         `json:"destino"`
	FechaDespacho      time.Time                    `json:"fecha_despacho"`
	ValorDespacho      float64                      `json:"valor_despacho"`
	CantidadItems      int                          `json:"cantidad_items"`
	TotalKg            float64                      `json:"total_kg"`
	DistanciaCalculada *string                      `json:"distancia_calculada,omitempty"`
	TiempoEstimado     *string                      `json:"tiempo_estimado,omitempty"`
	Cotizacion         *CotizacionDistanciaResponse `json:"cotizacion,omitempty"`
	Camion             *CamionDistanciaResponse     `json:"camion,omitempty"`
	OrigenSucursal     *SucursalDistanciaResponse   `json:"origen_sucursal,omitempty"`
	DestinoDirCliente  *DirClienteDistanciaResponse `json:"destino_dir_cliente,omitempty"`
}

// CotizacionDistanciaResponse es la estructura simplificada de cotización para rutas
type CotizacionDistanciaResponse struct {
	Cliente *ClienteDistanciaResponse `json:"cliente,omitempty"`
	Estado  string                    `json:"estado"`
}

// ClienteDistanciaResponse es la estructura simplificada de cliente para rutas
type ClienteDistanciaResponse struct {
	Nombre string `json:"nombre"`
	Email  string `json:"email"`
}

// CamionDistanciaResponse es la estructura simplificada de camión para rutas
type CamionDistanciaResponse struct {
	Patente string `json:"patente"`
}

// SucursalDistanciaResponse es la estructura simplificada de sucursal para rutas
type SucursalDistanciaResponse struct {
	Nombre    string `json:"nombre"`
	Direccion string `json:"direccion,omitempty"`
	Comuna    string `json:"comuna,omitempty"`
	Ciudad    string `json:"ciudad,omitempty"`
}

// DirClienteDistanciaResponse es la estructura simplificada de dirección de cliente para rutas
type DirClienteDistanciaResponse struct {
	Direccion string `json:"direccion"`
	Comuna    string `json:"comuna"`
	Ciudad    string `json:"ciudad"`
}

// DireccionesRequest es la estructura para la petición de cálculo de distancia
type DireccionesRequest struct {
	Origen  string `json:"origen" binding:"required"`
	Destino string `json:"destino" binding:"required"`
}

// DistanciaResponse es la estructura de respuesta para el cálculo de distancia
type DistanciaResponse struct {
	Distancia      string `json:"distancia"`
	Duracion       string `json:"duracion"`
	RutaOptimizada bool   `json:"ruta_optimizada"`
}
