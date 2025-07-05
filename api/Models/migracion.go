package modelos

import (
	"log"

	"gorm.io/gorm"
)

func MigrarTablas(db *gorm.DB) {
	err := db.AutoMigrate(
		&Producto{},
		&Proveedor{},
		&StockProveedor{},
		&TipoSucursal{},
		&Sucursal{},
		&StockSucursal{},
		&Rol{},
		&Usuario{},
		&TipoCliente{},
		&Cliente{},
		&DirCliente{},
		&Cotizacion{},
		&CotizacionItem{},
		&TipoCamion{},
		&Camion{},
		&Despacho{},
		&ProductosDespacho{},
	)
	if err != nil {
		log.Fatal("Error al migrar la base de datos:", err)
	}
	log.Println("Migraciones de base de datos completadas")

}
