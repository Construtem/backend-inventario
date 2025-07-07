package Controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/phpdave11/gofpdf"
	"gorm.io/gorm"
)

// Handler para /api/despachos/:id/pdf
func GenerarDespachoPDF(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
			return
		}

		// Obtener el despacho con totales
		despacho, err := GetDespachoByID(db, uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Despacho no encontrado"})
			return
		}

		// Crear PDF
		pdf := gofpdf.New("P", "mm", "A4", "")
		pdf.SetTitle("Guía de Despacho Electrónica", false)
		pdf.AddPage()

		// Encabezado de la empresa
		generarEncabezadoEmpresa(pdf)

		// RUT y Guía de despacho electrónica en recuadro rojo
		generarRecuadroRUT(pdf, strconv.Itoa(int(despacho.ID)))

		// Información del destinatario y despacho
		generarInfoDestinatario(pdf, despacho)

		// Tabla de productos
		generarTablaProductos(pdf, despacho)

		// Totales
		generarTotales(pdf, despacho)

		// Pie de página con código de barras simulado
		generarPiePagina(pdf)

		// PDF como stream
		c.Header("Content-Type", "application/pdf")
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=guia_despacho_%d.pdf", despacho.ID))
		err = pdf.Output(c.Writer)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo generar el PDF"})
		}
	}
}

func generarEncabezadoEmpresa(pdf *gofpdf.Fpdf) {
	// Logo/Nombre de la empresa (lado izquierdo)
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 8, "COMERCIAL FRANCISCO TOSO LTDA")
	pdf.Ln(6)
	
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 5, "Grandes Tiendas - Productos de Ferretería y para el")
	pdf.Ln(4)
	pdf.Cell(0, 5, "hogar")
	pdf.Ln(4)
	pdf.Cell(0, 5, "Matucana 90")
	pdf.Ln(4)
	pdf.Cell(0, 5, "Villa Alemana")
	pdf.Ln(4)
	pdf.Cell(0, 5, "Villa Alemana")
	pdf.Ln(10)
}

func generarRecuadroRUT(pdf *gofpdf.Fpdf, folioGuia string) {
	// Posicionar en la esquina superior derecha
	pdf.SetXY(130, 15)
	
	// Recuadro rojo para RUT y Guía
	pdf.SetDrawColor(255, 0, 0)
	pdf.SetLineWidth(1.5)
	pdf.Rect(130, 15, 65, 35, "D")
	
	// Contenido del recuadro
	pdf.SetXY(135, 20)
	pdf.SetFont("Arial", "B", 12)
	pdf.SetTextColor(255, 0, 0)
	pdf.Cell(0, 6, "R.U.T.: 76008058-6")
	pdf.SetXY(135, 28)
	pdf.Cell(0, 6, "Guía de despacho electrónica")
	pdf.SetXY(135, 36)
	pdf.Cell(0, 6, "Folio N°"+folioGuia)
	
	// Restaurar color de texto
	pdf.SetTextColor(0, 0, 0)
	pdf.SetDrawColor(0, 0, 0)
	pdf.SetLineWidth(0.2)
	
	pdf.SetXY(180, 52)
	pdf.SetFont("Arial", "", 8)
	pdf.Cell(0, 4, "UNIDAD S.I.L.")
	pdf.Ln(15)
}

func generarInfoDestinatario(pdf *gofpdf.Fpdf, despacho *DespachoConTotales) {
	// Recuadro para información del destinatario
	pdf.SetXY(10, 65)
	pdf.Rect(10, 65, 190, 20, "D")
	
	// Contenido del recuadro
	pdf.SetXY(15, 70)
	pdf.SetFont("Arial", "", 9)
	
	// Primera línea
	pdf.Cell(30, 4, "Fecha:")
	pdf.Cell(40, 4, despacho.FechaDespacho.Format("2006-01-02"))
	pdf.Cell(30, 4, "Señor(es):")
	pdf.Cell(40, 4, despacho.Cotizacion.Cliente.Nombre)
	pdf.Ln(5)
	
	// Segunda línea  
	pdf.SetX(15)
	pdf.Cell(30, 4, "Giro:")
	pdf.Cell(40, 4, "COMERCIO AL POR MENOR") // Campo genérico o agregar a modelo Cliente
	pdf.Cell(30, 4, "Ciudad:")
	pdf.Cell(30, 4, despacho.DestinoDirCliente.Ciudad)
	pdf.Cell(20, 4, "R.U.T.:")
	pdf.Cell(25, 4, despacho.Cotizacion.Cliente.Rut)
	pdf.Ln(5)
	
	// Tercera línea
	pdf.SetX(15)
	pdf.Cell(30, 4, "Comuna:")
	pdf.Cell(40, 4, despacho.DestinoDirCliente.Comuna)
	pdf.Cell(30, 4, "Dirección:")
	pdf.Cell(70, 4, despacho.DestinoDirCliente.Direccion)
	pdf.Cell(20, 4, "Forma de pago:")
	pdf.Ln(10)
}

func generarTablaProductos(pdf *gofpdf.Fpdf, despacho *DespachoConTotales) {
	// Encabezados de la tabla de referencia
	pdf.SetXY(10, 95)
	pdf.Rect(10, 95, 190, 12, "D") // Recuadro para encabezados
	
	pdf.SetFont("Arial", "B", 9)
	pdf.SetXY(15, 98)
	pdf.Cell(20, 6, "Tipo documento")
	pdf.Cell(30, 6, "Motivo referencia")
	pdf.Cell(20, 6, "Folio")
	pdf.Cell(20, 6, "Fecha")
	pdf.Ln(8)
	
	pdf.SetXY(15, 106)
	pdf.Cell(30, 6, "Cotización")
	pdf.Cell(50, 6, "Venta")
	pdf.Cell(20, 6, fmt.Sprintf("%d", despacho.CotizacionID))
	pdf.Cell(20, 6, despacho.Cotizacion.FechaCrea.Format("2006-01-02"))
	pdf.Ln(10)
	
	// Tabla de productos
	pdf.SetXY(10, 120)
	pdf.Rect(10, 120, 190, 12, "D") // Encabezado productos
	
	pdf.SetFont("Arial", "B", 9)
	pdf.SetXY(15, 123)
	pdf.Cell(20, 6, "CODIGO")
	pdf.Cell(70, 6, "DESCRIPCION")
	pdf.Cell(15, 6, "CANT.")
	pdf.Cell(25, 6, "P.UNITARIO")
	pdf.Cell(20, 6, "RECARGO")
	pdf.Cell(20, 6, "DESCUENTO")
	pdf.Cell(20, 6, "TOTAL")
	pdf.Ln(8)
	
	// Productos reales del despacho
	pdf.SetFont("Arial", "", 8)
	y := 135
	var subtotal float64
	
	for _, productoDespacho := range despacho.ProductosDespacho {
		pdf.SetXY(15, float64(y))
		pdf.Cell(20, 5, productoDespacho.Producto.SKU)
		pdf.Cell(70, 5, productoDespacho.Producto.Nombre)
		pdf.Cell(15, 5, fmt.Sprintf("%d", productoDespacho.Cantidad))
		pdf.Cell(25, 5, fmt.Sprintf("%.0f", productoDespacho.Producto.Precio))
		pdf.Cell(20, 5, "-") // recargo
		pdf.Cell(20, 5, "-") // descuento
		
		totalLinea := float64(productoDespacho.Cantidad) * productoDespacho.Producto.Precio
		subtotal += totalLinea
		pdf.Cell(20, 5, fmt.Sprintf("%.0f", totalLinea))
		y += 8
		
		if y > 200 { // Si se acaba el espacio, agregar nueva página
			break
		}
	}
	
	// Guardar subtotal para usar en totales
	pdf.SetXY(15, float64(y+5))
	pdf.SetFont("Arial", "I", 8)
	pdf.Cell(0, 4, fmt.Sprintf("Total de productos: %d | Peso total: %.2f kg", 
		despacho.CantidadItems, despacho.TotalKg))
}

func generarTotales(pdf *gofpdf.Fpdf, despacho *DespachoConTotales) {
	// Calcular totales basados en productos reales
	var subtotal float64
	for _, producto := range despacho.ProductosDespacho {
		subtotal += float64(producto.Cantidad) * producto.Producto.Precio
	}
	
	iva := subtotal * 0.19
	total := subtotal + iva
	
	// Totales en la parte inferior derecha
	pdf.SetXY(150, 200)
	pdf.SetFont("Arial", "", 10)
	
	pdf.Cell(20, 6, "Neto")
	pdf.Cell(20, 6, fmt.Sprintf("%.0f", subtotal))
	pdf.Ln(5)
	
	pdf.SetX(150)
	pdf.Cell(20, 6, "IVA (19%)")
	pdf.Cell(20, 6, fmt.Sprintf("%.0f", iva))
	pdf.Ln(5)
	
	pdf.SetX(150)
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(20, 6, "Total")
	pdf.Cell(20, 6, fmt.Sprintf("%.0f", total))
}

func generarPiePagina(pdf *gofpdf.Fpdf) {
	// Información adicional en la parte inferior
	pdf.SetXY(10, 230)
	pdf.SetFont("Arial", "", 8)
	
	// Simulación de código de barras
	pdf.Rect(25, 230, 50, 20, "D")
	pdf.SetXY(30, 235)
	for i := 0; i < 45; i++ {
		pdf.Line(30+float64(i), 235, 30+float64(i), 245)
	}
	
	// Texto del pie
	pdf.SetXY(85, 235)
	pdf.Cell(0, 4, "Timbre electrónico SII")
	pdf.Ln(3)
	pdf.SetX(85)
	pdf.Cell(0, 4, "Resolución 0 del 2021 - Verifique su documento en: www.sii.cl")
	
	// Información adicional del lado derecho
	pdf.SetXY(130, 230)
	pdf.Cell(20, 4, "R.U.T.:")
	pdf.Ln(3)
	pdf.SetX(130)
	pdf.Cell(20, 4, "NOMBRE:")
	pdf.Ln(3)
	pdf.SetX(130)
	pdf.Cell(20, 4, "FECHA:")
	pdf.Cell(30, 4, "RECINTO:")
	pdf.Ln(8)
	pdf.SetX(130)
	pdf.Cell(20, 4, "FIRMA:")
	
	// Texto legal en la parte inferior
	pdf.SetXY(10, 255)
	pdf.SetFont("Arial", "", 7)
	pdf.Cell(0, 3, "\"El acuse de recibo que se declara en este acto, de acuerdo a lo dispuesto en la letra")
	pdf.Ln(3)
	pdf.Cell(0, 3, "b) del artículo 4°, y la letra c) del artículo 5° de la ley 19.983, acredita que la entrega de")
	pdf.Ln(3)
	pdf.Cell(0, 3, "mercaderías o servicios) prestado(s) ha(n) sido recibido(s)\"")
}
