package modelos

import "time"

type Producto struct {
	SKU         string  `gorm:"primaryKey;size:20" json:"sku"`
	Nombre      string  `gorm:"size:100;not null" json:"nombre"`
	Descripcion string  `gorm:"type:text" json:"descripcion"`
	Marca       string  `gorm:"size:30;not null" json:"marca"`
	Peso        float64 `gorm:"type:numeric(10,2);not null" json:"peso"`
	Largo       float64 `gorm:"type:numeric(10,2);not null" json:"largo"`
	Ancho       float64 `gorm:"type:numeric(10,2);not null" json:"ancho"`
	Alto        float64 `gorm:"type:numeric(10,2);not null" json:"alto"`
	Precio      float64 `gorm:"type:numeric(10,2);not null" json:"precio"`
}

type Proveedor struct {
	ID    uint   `gorm:"primaryKey" json:"id"`
	Marca string `gorm:"size:30;not null" json:"marca"`
}

type StockProveedor struct {
	ProveedorID uint   `gorm:"primaryKey;column:proveedor_id" json:"proveedor_id"`
	ProductoID  string `gorm:"primaryKey;size:20;column:producto_id" json:"producto_id"`
	Stock       int    `gorm:"not null" json:"stock"`

	Proveedor Proveedor `gorm:"foreignKey:ProveedorID;references:ID;constraint:OnDelete:CASCADE" json:"proveedor"`
	Producto  Producto  `gorm:"foreignKey:ProductoID;references:SKU;constraint:OnDelete:CASCADE" json:"producto"`
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
	ProductoID string  `gorm:"primaryKey;size:20;column:producto_id" json:"producto_id"`
	SucursalID uint    `gorm:"primaryKey;column:sucursal_id" json:"sucursal_id"`
	Cantidad   int     `gorm:"not null" json:"cantidad"`
	Descuento  float64 `gorm:"type:numeric(5,2);default:0;check:descuento >= 0 AND descuento <= 100" json:"descuento"`

	Producto Producto `gorm:"foreignKey:ProductoID;references:SKU;constraint:OnDelete:CASCADE" json:"producto"`
	Sucursal Sucursal `gorm:"foreignKey:SucursalID;references:ID;constraint:OnDelete:CASCADE" json:"sucursal"`
}

type Rol struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	Nombre string `gorm:"size:50;not null" json:"nombre"`
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

type Cliente struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Nombre      string `gorm:"size:100;not null" json:"nombre"`
	Telefono    string `gorm:"size:20" json:"telefono"`
	Email       string `gorm:"size:100;not null;unique" json:"email"`
	RazonSocial string `gorm:"size:100" json:"razon_social"`
	Rut         string `gorm:"size:12;not null;unique" json:"rut"`
	TipoID      uint   `gorm:"column:tipo_id;not null" json:"tipo_id"`

	Tipo TipoCliente `gorm:"foreignKey:TipoID;references:ID;constraint:OnDelete:CASCADE" json:"tipo"`
}

type DirCliente struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	ClienteID uint   `gorm:"column:cliente_id;not null" json:"cliente_id"`
	Nombre    string `gorm:"size:100;not null" json:"nombre"`
	Direccion string `gorm:"size:200;not null" json:"direccion"`
	Comuna    string `gorm:"size:100;not null" json:"comuna"`
	Ciudad    string `gorm:"size:100;not null" json:"ciudad"`

	Cliente Cliente `gorm:"foreignKey:ClienteID;references:ID;constraint:OnDelete:CASCADE" json:"cliente"`
}

func (DirCliente) TableName() string {
	return "dir_cliente"
}

type Cotizacion struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	FechaCrea    time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"fecha_crea"`
	Estado       string    `gorm:"size:20;not null" json:"estado"`
	CostoEnvio   float64   `gorm:"type:numeric(10,2);not null" json:"costo_envio"`
	ClienteID    uint      `gorm:"column:cliente_id;not null" json:"cliente_id"`
	UserID       string    `gorm:"size:100;column:user_id;not null" json:"user_id"`
	TipoDespacho string    `gorm:"size:50;not null" json:"tipo_despacho"`

	Cliente Cliente `gorm:"foreignKey:ClienteID;references:ID;constraint:OnDelete:CASCADE" json:"cliente"`
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
	ID            uint      `gorm:"primaryKey" json:"id"`
	CotizacionID  uint      `gorm:"column:cotizacion_id;not null" json:"cotizacion_id"`
	CamionID      uint      `gorm:"column:camion_id;not null" json:"camion_id"`
	Origen        uint      `gorm:"not null" json:"origen"`  // FK a sucursales.id
	Destino       uint      `gorm:"not null" json:"destino"` // FK a dir_cliente.id
	FechaDespacho time.Time `gorm:"not null" json:"fecha_despacho"`
	ValorDespacho float64   `gorm:"type:numeric(10,2);not null" json:"valor_despacho"`

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
	ProductoID string `gorm:"primaryKey;size:20;column:producto_id" json:"producto_id"`
	Cantidad   int    `gorm:"not null" json:"cantidad"`

	Despacho Despacho `gorm:"foreignKey:DespachoID;references:ID;constraint:OnDelete:CASCADE" json:"despacho"`
	Producto Producto `gorm:"foreignKey:ProductoID;references:SKU;constraint:OnDelete:CASCADE" json:"producto"`
}

func (ProductosDespacho) TableName() string {
	return "productos_despacho"
}
