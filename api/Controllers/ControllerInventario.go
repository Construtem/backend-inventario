package Controllers

import (
	"encoding/csv"
	"io"
	"net/http"
	"strconv"
	"strings"

	modelos "backend-inventario/api/Models"
	"backend-inventario/api/db"

	"github.com/gin-gonic/gin"
)

// GetInventario obtiene todos los registros de inventario
func GetInventario(c *gin.Context) {
	var inventarios []modelos.Inventario
	// Preload las relaciones relevantes
	if err := db.DB.Preload("Producto").Preload("Ubicacion").Find(&inventarios).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener inventario"})
		return
	}
	c.JSON(http.StatusOK, inventarios)
}

// GetInventarioByID obtiene un registro de inventario por su ID
func GetInventarioByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var inventario modelos.Inventario
	if err := db.DB.Preload("Producto").Preload("Ubicacion").First(&inventario, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Registro de inventario no encontrado"})
		return
	}
	c.JSON(http.StatusOK, inventario)
}

// CreateInventario crea un nuevo registro de inventario
func CreateInventario(c *gin.Context) {
	var inventario modelos.Inventario
	if err := c.ShouldBindJSON(&inventario); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Considerar validaciones adicionales, como si el ProductoID y UbicacionID existen

	if err := db.DB.Create(&inventario).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear registro de inventario"})
		return
	}
	// Carga las relaciones para la respuesta
	db.DB.Preload("Producto").Preload("Ubicacion").First(&inventario, inventario.ID)
	c.JSON(http.StatusCreated, inventario)
}

// UpdateInventario actualiza un registro de inventario existente
func UpdateInventario(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var inventario modelos.Inventario
	if err := db.DB.First(&inventario, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Registro de inventario no encontrado"})
		return
	}

	if err := c.ShouldBindJSON(&inventario); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Save(&inventario).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar registro de inventario"})
		return
	}
	// Carga las relaciones para la respuesta
	db.DB.Preload("Producto").Preload("Ubicacion").First(&inventario, inventario.ID)
	c.JSON(http.StatusOK, inventario)
}

// DeleteInventario elimina un registro de inventario
func DeleteInventario(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var inventario modelos.Inventario
	if err := db.DB.First(&inventario, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Registro de inventario no encontrado"})
		return
	}

	if err := db.DB.Delete(&inventario).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar registro de inventario"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Registro de inventario eliminado exitosamente"})
}

func CargaMasivaInventario(c *gin.Context) {
	// Validacion de permisos (opcional)
	// Descomentar y ajustar según sea necesario para la autenticación y autorización
	/*
		user, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
			return
		}
		if user.(modelos.Usuario).RolID != 1 {
			c.JSON(http.StatusForbidden, gin.H{"error": "Acceso restringido"})
			return
		}
	*/

	// recepcion del archivo .csv
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Archivo no proporcionado"})
		return
	}

	openedFile, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al abrir el archivo"})
		return
	}
	defer openedFile.Close()

	// lectura del archivo CSV
	reader := csv.NewReader(openedFile)
	reader.Comma = ','
	reader.LazyQuotes = true

	// encabezados
	headers, err := reader.Read()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al leer el encabezado CSV"})
		return
	}

	// mapa para obtener indices por nombre
	index := map[string]int{}
	for i, h := range headers {
		index[h] = i
	}

	const UbicacionID = 1 // ID de la ubicación por defecto, ajustar según sea necesario

	var productos []modelos.Producto
	var inventarios []modelos.Inventario
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al leer el registro CSV"})
			return
		}

		// validacion de datos
		precioCosto, err := strconv.ParseFloat(record[index["Costo base (c/u)"]], 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Costo base inválido en csv"})
			return
		}

		precioVenta, err := strconv.ParseFloat(record[index["Precio de venta (c/u"]], 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Precio de venta inválido en csv"})
			return
		}

		activo := strings.ToLower(record[index["Activo"]]) == "activo"

		stock, err := strconv.Atoi(record[index["Stock"]])
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Stock inválido en csv"})
			return
		}

		producto := modelos.Producto{
			Codigo:      record[index["SKU"]],
			Descripcion: record[index["Descripción"]],
			PrecioCosto: precioCosto,
			PrecioVenta: precioVenta,
			Activo:      activo,
		}
		productos = append(productos, producto)

		inventario := modelos.Inventario{
			UbicacionID: UbicacionID,
			Cantidad:    stock,
		}
		inventarios = append(inventarios, inventario)

	}
	for i := range productos {
		if err := db.DB.Create(&productos[i]).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error insertando productos"})
			return
		}
		inventarios[i].ProductoID = productos[i].ID
		if err := db.DB.Create(&inventarios[i]).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error insertando inventarios"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Carga masiva de inventario completada exitosamente"})
}
