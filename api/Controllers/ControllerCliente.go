package Controllers

import (
	modelos "backend-inventario/api/Models"
	"errors"

	"gorm.io/gorm"
)

func GetClientes(db *gorm.DB) ([]modelos.Cliente, error) {
	var clientes []modelos.Cliente
	if err := db.Preload("Tipo").Find(&clientes).Error; err != nil {
		return nil, err
	}
	return clientes, nil
}

func GetClienteByID(db *gorm.DB, id uint) (*modelos.Cliente, error) {
	var cliente modelos.Cliente
	if err := db.Preload("Tipo").First(&cliente, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &cliente, nil
}

func CreateCliente(db *gorm.DB, nuevo *modelos.Cliente) error {
	if nuevo.Nombre == "" {
		return errors.New("el nombre del cliente no puede estar vacío")
	}
	if nuevo.TipoID == 0 {
		return errors.New("el tipo de cliente es obligatorio")
	}
	return db.Create(nuevo).Error
}

func UpdateCliente(db *gorm.DB, id uint, actualizado *modelos.Cliente) (*modelos.Cliente, error) {
	var existente modelos.Cliente
	if err := db.First(&existente, id).Error; err != nil {
		return nil, errors.New("cliente no encontrado")
	}

	if actualizado.Nombre == "" {
		return nil, errors.New("el nombre del cliente no puede estar vacío")
	}
	if actualizado.TipoID == 0 {
		return nil, errors.New("el tipo de cliente es obligatorio")
	}

	err := db.Model(&existente).Updates(modelos.Cliente{
		Nombre:      actualizado.Nombre,
		Telefono:    actualizado.Telefono,
		Email:       actualizado.Email,
		RazonSocial: actualizado.RazonSocial,
		Rut:         actualizado.Rut,
		TipoID:      actualizado.TipoID,
	}).Error
	if err != nil {
		return nil, err
	}

	return &existente, nil
}

func DeleteCliente(db *gorm.DB, id uint) error {
	var cliente modelos.Cliente
	if err := db.First(&cliente, id).Error; err != nil {
		return errors.New("cliente no encontrado")
	}

	// Verificar si el cliente tiene despachos asociados
	var count int64
	if err := db.Model(&modelos.Despacho{}).Where("destino_dir_cliente_id = ?", id).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("no se puede eliminar el cliente porque tiene despachos asociados")
	}

	return db.Delete(&cliente, id).Error
}
