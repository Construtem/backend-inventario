package modelos

import (
	"log"

	"gorm.io/gorm"
)

func MigrarTablas(db *gorm.DB) {
	err := db.AutoMigrate(
		&Categoria{},
		&Cliente{},
		&Producto{},
		&Camion{},
		&Pedido{},
		&Cotizacion{},
		&Usuario{},
		&Rol{},
		&Ubicacion{},
		&Proveedor{},
		&Inventario{},
		&Despacho{},
		&OrdenCompra{},
		&DetalleCotizacion{},
		&DetalleOrdenCompra{},
		&DetallePedido{},
		&DetalleDespacho{},
	)
	if err != nil {
		log.Fatal("Error al migrar la base de datos:", err)
	}
	log.Println("Migraciones de base de datos completadas")

}
