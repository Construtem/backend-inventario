package Controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/phpdave11/gofpdf"
	"gorm.io/gorm"
)

// Función auxiliar para formatear fechas
func formatDate(t time.Time) string {
	return t.Format("02/01/2006")
}

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

		// Crear PDF con configuración profesional
		pdf := gofpdf.New("P", "mm", "A4", "")
		pdf.SetTitle("Guía de Despacho Electrónica", false)
		pdf.SetAutoPageBreak(false, 0)
		pdf.AddPage()

		// Configurar transformación UTF-8
		tr := pdf.UnicodeTranslatorFromDescriptor("")
		pdf.SetFont("Arial", "", 10)

		// Generar el PDF con estructura profesional
		err = generarPDFEstructurado(pdf, tr, despacho)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al generar PDF: " + err.Error()})
			return
		}

		// PDF como stream
		c.Header("Content-Type", "application/pdf")
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=guia_despacho_%d.pdf", despacho.ID))
		err = pdf.Output(c.Writer)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo generar el PDF"})
		}
	}
}

func generarPDFEstructurado(pdf *gofpdf.Fpdf, tr func(string) string, despacho *DespachoConTotales) error {
	// 1. Datos del Emisor y Título de Guía de Despacho
	initialY := pdf.GetY()

	// --- Sección del Título y RUT (Recuadro Rojo a la Derecha) ---
	pdf.SetDrawColor(255, 0, 0)
	pdf.SetFillColor(255, 255, 255)
	pdf.SetTextColor(0, 0, 0)
	pdf.SetLineWidth(0.5)

	// Posición y dimensiones del recuadro rojo
	rectRightX := 120.0
	rectRightY := 10.0
	rectRightWidth := 80.0
	rectRightHeight := 30.0

	pdf.Rect(rectRightX, rectRightY, rectRightWidth, rectRightHeight, "D")

	// Contenido dentro del recuadro rojo
	pdf.SetTextColor(255, 0, 0)

	// R.U.T. centrado
	pdf.SetFont("Arial", "B", 10)
	pdf.SetXY(rectRightX+5, rectRightY+3)
	pdf.CellFormat(rectRightWidth-10, 5, tr("R.U.T.: 76008058-6"), "", 0, "C", false, 0, "")

	// "Guía de Despacho Electrónica" centrado
	pdf.SetFont("Arial", "B", 12)
	pdf.SetXY(rectRightX+5, rectRightY+10)
	pdf.MultiCell(rectRightWidth-10, 5, tr("GUIA DE DESPACHO ELECTRONICA"), "", "C", false)

	// N° folio centrado
	pdf.SetFont("Arial", "B", 10)
	pdf.SetXY(rectRightX+5, rectRightY+22)
	pdf.CellFormat(rectRightWidth-10, 5, tr(fmt.Sprintf("Folio N° %d", despacho.ID)), "", 0, "C", false, 0, "")

	// Restablecer color de dibujo a negro
	pdf.SetDrawColor(0, 0, 0)

	// --- Sección de Datos de la Empresa (A la Izquierda) ---
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Arial", "B", 14)
	pdf.SetXY(15, initialY+5)
	pdf.Cell(0, 6, tr("COMERCIAL FRANCISCO TOSO LTDA"))
	pdf.Ln(6)

	pdf.SetFont("Arial", "", 9)
	pdf.SetX(15)
	pdf.Cell(0, 4, tr("Grandes Tiendas - Productos de Ferretería y para el hogar"))
	pdf.Ln(4)
	pdf.SetX(15)
	pdf.Cell(0, 4, tr("Dirección: Matucana 90, Villa Alemana"))
	pdf.Ln(4)
	pdf.SetX(15)
	pdf.Cell(0, 4, tr("Email: contacto@franciscotoso.cl | Tel: +56 32 2345678"))
	pdf.Ln(15)

	// 2. Rectángulo naranja con información del despacho
	currentY := pdf.GetY()
	pdf.SetFillColor(255, 102, 0)
	pdf.Rect(10, currentY, 190, 25, "F")

	pdf.SetTextColor(255, 255, 255)
	pdf.SetFont("Arial", "B", 10)

	// Primera línea
	pdf.SetXY(15, currentY+5)
	pdf.Cell(0, 5, tr(fmt.Sprintf("N° Despacho: %d", despacho.ID)))

	pdf.SetXY(110, currentY+5)
	pdf.Cell(0, 5, tr(fmt.Sprintf("Fecha de Despacho: %s", formatDate(despacho.FechaDespacho))))

	// Segunda línea
	pdf.SetXY(15, currentY+15)
	pdf.Cell(0, 5, tr(fmt.Sprintf("Cliente: %s", despacho.Cotizacion.Cliente.Nombre)))

	pdf.SetXY(110, currentY+15)
	pdf.Cell(0, 5, tr(fmt.Sprintf("RUT Cliente: %s", despacho.Cotizacion.Cliente.Rut)))

	// Posicionar después del rectángulo
	pdf.SetY(currentY + 30)
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Arial", "B", 12)
	pdf.Ln(5)

	// 3. Información adicional del despacho
	pdf.SetFont("Arial", "", 9)
	pdf.SetX(15)
	pdf.Cell(0, 4, tr(fmt.Sprintf("Origen: %s", despacho.OrigenSucursal.Nombre)))
	pdf.Ln(4)
	pdf.SetX(15)
	pdf.Cell(0, 4, tr(fmt.Sprintf("Destino: %s, %s, %s", despacho.DestinoDirCliente.Direccion, despacho.DestinoDirCliente.Comuna, despacho.DestinoDirCliente.Ciudad)))
	pdf.Ln(4)
	pdf.SetX(15)
	pdf.Cell(0, 4, tr(fmt.Sprintf("Camión: %s", despacho.Camion.Patente)))
	pdf.Ln(4)
	pdf.SetX(15)
	pdf.Cell(0, 4, tr(fmt.Sprintf("Estado: %s", despacho.Cotizacion.Estado)))
	pdf.Ln(10)

	// 4. Línea separadora superior
	pdf.SetLineWidth(0.6)
	pdf.Line(10, pdf.GetY(), 200, pdf.GetY())
	pdf.SetLineWidth(0.2)
	pdf.Ln(5)

	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(0, 8, tr("DETALLES DEL DESPACHO"), "", 1, "C", false, 0, "")
	pdf.Ln(5)

	// Línea separadora inferior
	pdf.SetLineWidth(0.6)
	pdf.Line(10, pdf.GetY(), 200, pdf.GetY())
	pdf.SetLineWidth(0.2)
	pdf.Ln(5)

	// 5. Encabezados de la tabla
	pdf.SetFont("Arial", "B", 9)
	pdf.SetFillColor(240, 240, 240)
	pdf.SetTextColor(0, 0, 0)

	pdf.CellFormat(30, 8, tr("Código"), "1", 0, "C", true, 0, "")
	pdf.CellFormat(70, 8, tr("Descripción"), "1", 0, "C", true, 0, "")
	pdf.CellFormat(20, 8, tr("Cantidad"), "1", 0, "C", true, 0, "")
	pdf.CellFormat(35, 8, tr("Peso Unit (kg)"), "1", 0, "C", true, 0, "")
	pdf.CellFormat(35, 8, tr("Peso Total (kg)"), "1", 1, "C", true, 0, "")

	// 6. Filas de datos
	pdf.SetFont("Arial", "", 9)
	pdf.SetFillColor(255, 255, 255)

	for _, item := range despacho.ProductosDespacho {
		pesoTotal := float64(item.Cantidad) * item.Producto.Peso
		pdf.CellFormat(30, 8, tr(item.Producto.SKU), "1", 0, "C", false, 0, "")
		pdf.CellFormat(70, 8, tr(item.Producto.Nombre), "1", 0, "L", false, 0, "")
		pdf.CellFormat(20, 8, fmt.Sprintf("%d", item.Cantidad), "1", 0, "C", false, 0, "")
		pdf.CellFormat(35, 8, fmt.Sprintf("%.2f", item.Producto.Peso), "1", 0, "R", false, 0, "")
		pdf.CellFormat(35, 8, fmt.Sprintf("%.2f", pesoTotal), "1", 1, "R", false, 0, "")
	}

	// 7. Línea bajo la tabla
	pdf.Ln(5)
	pdf.SetLineWidth(0.6)
	pdf.Line(10, pdf.GetY(), 200, pdf.GetY())
	pdf.SetLineWidth(0.2)
	pdf.Ln(5)

	// 8. Totales en recuadro
	pdf.Ln(5)
	pdf.SetFont("Arial", "", 10)

	// Posición para los totales
	startTotalsY := pdf.GetY()
	rectX := 130.0
	rectWidth := 70.0
	numRows := 2
	rowHeight := 7.0
	rectHeight := float64(numRows) * rowHeight

	// Dibujar recuadro de totales
	pdf.SetDrawColor(0, 0, 0)
	pdf.SetLineWidth(0.2)
	pdf.Rect(rectX, startTotalsY, rectWidth, rectHeight, "D")

	pdf.SetY(startTotalsY)

	// Total Items
	pdf.SetX(rectX)
	pdf.CellFormat(45, rowHeight, "Total Items:", "", 0, "L", false, 0, "")
	pdf.CellFormat(20, rowHeight, fmt.Sprintf("%d", despacho.CantidadItems), "", 1, "R", false, 0, "")

	// Total Peso
	pdf.SetX(rectX)
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(45, rowHeight, "Total Peso (kg):", "", 0, "L", false, 0, "")
	pdf.CellFormat(20, rowHeight, fmt.Sprintf("%.2f", despacho.TotalKg), "", 1, "R", false, 0, "")

	pdf.Ln(20)

	// 9. Mensaje final
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(0, 5, "Gracias por confiar en nosotros", "", 1, "C", false, 0, "")
	pdf.Ln(10)

	// 10. Pie de página
	generarPieDespacho(pdf, tr, despacho)

	return nil
}

func generarPieDespacho(pdf *gofpdf.Fpdf, tr func(string) string, despacho *DespachoConTotales) {
	pageWidth, pageHeight := pdf.GetPageSize()

	// Rectángulo naranja para pie de página
	pdf.SetFillColor(255, 102, 0)
	pdf.Rect(0, pageHeight-35, pageWidth, 35, "F")

	// Texto dentro del rectángulo naranja
	pdf.SetTextColor(255, 255, 255)

	// Frase legal
	pdf.SetFont("Arial", "I", 7)
	pdf.SetXY(10, pageHeight-30)
	pdf.MultiCell(pageWidth-20, 3, tr("Esta guía de despacho es representación fiel del documento electrónico firmado digitalmente según ley N° 19.799"), "", "C", false)

	// S.I.I. - Santiago
	pdf.SetXY(10, pageHeight-22)
	pdf.SetFont("Arial", "", 8)
	pdf.CellFormat(0, 4, "S.I.I. - Santiago", "", 1, "C", false, 0, "")

	// Información adicional
	pdf.SetXY(10, pageHeight-15)
	pdf.SetFont("Arial", "", 7)
	pdf.CellFormat(0, 4, "Timbre electrónico - Resolución 0 del 2021", "", 1, "C", false, 0, "")

	// Número de despacho en la esquina derecha
	pdf.SetXY(pageWidth-40, pageHeight-12)
	pdf.SetFont("Arial", "B", 8)
	pdf.Cell(0, 5, fmt.Sprintf("Despacho N° %d", despacho.ID))
}
