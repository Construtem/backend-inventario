package Controllers

import (
	"fmt"
	"log"
	"net/http"
	"os"
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

// Estructura del encabezado del PDF
type CompanyConfig struct {
	RazonSocial string
	Giro        string
	Direccion   string
	Comuna      string
	Ciudad      string
	Telefono    string
	Email       string
	LogoPath    string
	TimbrePath  string
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

		config := CompanyConfig{
			RazonSocial: "Constructem Ltda.",
			Giro:        "Grandes Tiendas - Productos de Ferretería y para el hogar",
			Direccion:   "Matucana 90, Villa Alemana",
			Comuna:      "Providencia",
			Ciudad:      "Santiago",
			Telefono:    "+56 32 2345678",
			Email:       "Email: contacto@franciscotoso.cl",
			LogoPath:    "img/construtem.png",
			TimbrePath:  "img/TimbreElectronico.png",
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
		err = generarPDFEstructurado(pdf, tr, despacho, config)
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

func generarPDFEstructurado(pdf *gofpdf.Fpdf, tr func(string) string, despacho *DespachoConTotales, config CompanyConfig) error {
	// 1. Datos del Emisor y Título de Guía de Despacho
	initialY := 20.0

	// --- Sección del Título y RUT (Recuadro Rojo a la Derecha) ---
	pdf.SetDrawColor(255, 0, 0)
	pdf.SetFillColor(255, 255, 255)
	pdf.SetTextColor(0, 0, 0)
	pdf.SetLineWidth(0.5)

	// Posición y dimensiones del recuadro rojo
	rectRightX := 120.0
	rectRightY := 5.0
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
	// Cargar logo
	logoPath := config.LogoPath
	logoX := 15.0
	logoY := 2.5
	logoMaxWidth := 60.0
	logoMaxHeight := 25.0
	if _, err := os.Stat(logoPath); err == nil {
		info := pdf.RegisterImage(logoPath, "")
		if info != nil {
			// Calcular dimensiones manteniendo proporción
			logoWidth := logoMaxWidth
			logoHeight := logoWidth * info.Height() / info.Width()

			// Ajustar si la altura excede el máximo
			if logoHeight > logoMaxHeight {
				logoHeight = logoMaxHeight
				logoWidth = logoHeight * info.Width() / info.Height()
			}

			pdf.Image(logoPath, logoX, logoY, logoWidth, logoHeight, false, "", 0, "")
			log.Printf("INFO: Logo cargado correctamente desde %s", logoPath)
		} else {
			log.Printf("ADVERTENCIA: No se pudo procesar la imagen del logo en %s", logoPath)
		}
	} else {
		log.Printf("ADVERTENCIA: Logo no encontrado en %s", logoPath)
	}

	// Posicionar debajo del logo
	pdf.SetY(0 + logoY + 25) // Ajusta según la altura del logo

	// Razon Social
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Arial", "B", 14)
	pdf.SetXY(15, initialY+5) // Asegúrate que initialY esté definido (puedes usar 35 o 40 si no)
	pdf.Cell(0, 6, tr(config.RazonSocial))
	pdf.Ln(6)

	// Giro
	pdf.SetFont("Arial", "", 9)
	pdf.SetX(15)
	pdf.Cell(0, 4, tr(config.Giro))
	pdf.Ln(4)

	//Dirección y comuna
	pdf.SetX(15)
	pdf.Cell(0, 4, tr(fmt.Sprintf("Dirección: %s, %s", config.Direccion, config.Comuna)))
	pdf.Ln(4)

	// Email y teléfono
	pdf.SetX(15)
	pdf.Cell(0, 4, tr(fmt.Sprintf("Email: %s | Tel: %s", config.Email, config.Telefono)))
	pdf.Ln(15)

	// 2. Rectángulo naranja con información del despacho
	currentY := pdf.GetY()
	pdf.SetFillColor(255, 102, 0)
	pdf.Rect(10, currentY, 190, 30, "F")

	// Texto en blanco dentro del rectángulo
	pdf.SetTextColor(255, 255, 255)
	pdf.SetFont("Arial", "B", 10)

	// Línea 1 - N° Despacho y Fecha
	pdf.SetXY(15, currentY+5)
	pdf.Cell(0, 5, tr(fmt.Sprintf("N° Despacho: %d", despacho.ID)))

	pdf.SetXY(100, currentY+5)
	pdf.Cell(0, 5, tr(fmt.Sprintf("Fecha de Despacho: %s", formatDate(despacho.FechaDespacho))))

	// Línea 2 - Cliente y RUT
	pdf.SetXY(15, currentY+10)
	pdf.Cell(0, 5, tr(fmt.Sprintf("Cliente: %s", despacho.Cotizacion.Cliente.Nombre)))

	pdf.SetXY(100, currentY+10)
	pdf.Cell(0, 5, tr(fmt.Sprintf("RUT Cliente: %s", despacho.Cotizacion.Cliente.Rut)))

	// Línea 3 - Origen y Destino
	pdf.SetXY(15, currentY+15)
	pdf.Cell(0, 5, tr(fmt.Sprintf("Origen: %s", despacho.OrigenSucursal.Nombre)))

	pdf.SetXY(100, currentY+15)
	pdf.Cell(0, 5, tr(fmt.Sprintf("Destino: %s, %s, %s", despacho.DestinoDirCliente.Direccion, despacho.DestinoDirCliente.Comuna, despacho.DestinoDirCliente.Ciudad)))

	// Línea 4 - Camión y Estado
	pdf.SetXY(15, currentY+20)
	pdf.Cell(0, 5, tr(fmt.Sprintf("Camión: %s", despacho.Camion.Patente)))

	pdf.SetXY(100, currentY+20)
	pdf.Cell(0, 5, tr(fmt.Sprintf("Estado: %s", despacho.Cotizacion.Estado)))

	// Mover el cursor después del bloque
	pdf.SetY(currentY + 30)
	pdf.SetTextColor(0, 0, 0)

	// Posicionar después del rectángulo
	pdf.SetY(currentY + 50) // avanzar para continuar debajo del rectángulo
	pdf.SetTextColor(0, 0, 0)

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
	pdf.SetFillColor(255, 102, 0)   // Fondo naranja
	pdf.SetTextColor(255, 255, 255) // Texto blanco

	pdf.CellFormat(15, 8, tr("Código"), "", 0, "C", true, 0, "")
	pdf.CellFormat(55, 8, tr("Nombre"), "", 0, "C", true, 0, "")
	pdf.CellFormat(15, 8, tr("Cantidad"), "", 0, "C", true, 0, "")
	pdf.CellFormat(25, 8, tr("Peso Unit (kg)"), "", 0, "C", true, 0, "")
	pdf.CellFormat(25, 8, tr("Peso Total (kg)"), "", 0, "C", true, 0, "")
	pdf.CellFormat(25, 8, tr("Precio Unit"), "", 0, "C", true, 0, "")
	pdf.CellFormat(30, 8, tr("Precio Total Neto"), "", 1, "C", true, 0, "")

	// 6. Filas de datos
	pdf.SetFont("Arial", "", 9)
	pdf.SetTextColor(0, 0, 0) // Texto negro

	for i, item := range despacho.ProductosDespacho {
		// Color de fondo alternado
		if i%2 == 0 {
			pdf.SetFillColor(255, 255, 255) // Blanco
		} else {
			pdf.SetFillColor(240, 240, 240) // Gris claro
		}

		// Bordes invisibles
		pdf.SetDrawColor(255, 255, 255)

		pdf.CellFormat(15, 8, tr(item.SKU), "", 0, "C", true, 0, "")
		pdf.CellFormat(55, 8, tr(item.Nombre), "", 0, "L", true, 0, "")
		pdf.CellFormat(15, 8, fmt.Sprintf("%d", item.Cantidad), "", 0, "C", true, 0, "")
		pdf.CellFormat(25, 8, fmt.Sprintf("%.2f", item.Peso), "", 0, "R", true, 0, "")
		pdf.CellFormat(25, 8, fmt.Sprintf("%.2f", item.PesoTotal), "", 0, "R", true, 0, "")
		pdf.CellFormat(25, 8, fmt.Sprintf("$%.2f", item.Precio), "", 0, "R", true, 0, "")
		pdf.CellFormat(30, 8, fmt.Sprintf("$%.2f", item.PrecioTotal), "", 1, "R", true, 0, "")
	}

	// 7. Línea bajo la tabla
	pdf.SetDrawColor(0, 0, 0) // Negro para la línea
	pdf.SetLineWidth(0.6)
	pdf.Line(10, (pdf.GetY() + 5), 200, (pdf.GetY() + 5))
	pdf.SetLineWidth(0.2)
	pdf.Ln(5)

	// 8. Totales en recuadro
	pdf.Ln(5)
	pdf.SetFont("Arial", "", 10)

	// Posición para los totales
	startTotalsY := pdf.GetY()
	rectX := 130.0
	rectWidth := 70.0
	numRows := 6
	rowHeight := 7.0
	rectHeight := float64(numRows) * rowHeight

	// Dibujar recuadro de totales
	pdf.SetDrawColor(0, 0, 0)
	pdf.SetFillColor(255, 255, 255) // Fondo blanco o el color que desees
	pdf.SetLineWidth(0.2)
	pdf.Rect(rectX, startTotalsY, rectWidth, rectHeight, "FD")

	pdf.SetY(startTotalsY)

	// Total Items
	pdf.SetX(rectX)
	pdf.CellFormat(45, rowHeight, "Total Items:", "", 0, "L", false, 0, "")
	pdf.CellFormat(20, rowHeight, fmt.Sprintf("%d", despacho.CantidadItems), "", 1, "R", false, 0, "")

	// Total Peso
	pdf.SetX(rectX)
	pdf.CellFormat(45, rowHeight, "Total Peso (kg):", "", 0, "L", false, 0, "")
	pdf.CellFormat(20, rowHeight, fmt.Sprintf("%.2f", despacho.TotalKg), "", 1, "R", false, 0, "")

	// Total Precio Neto
	pdf.SetX(rectX)
	pdf.CellFormat(45, rowHeight, "Total Precio Neto:", "", 0, "L", false, 0, "")
	pdf.CellFormat(20, rowHeight, fmt.Sprintf("$%.0f", despacho.TotalPrecio), "", 1, "R", false, 0, "")

	// IVA 19%
	pdf.SetX(rectX)
	pdf.CellFormat(45, rowHeight, "IVA (19%):", "", 0, "L", false, 0, "")
	pdf.CellFormat(20, rowHeight, fmt.Sprintf("$%.0f", despacho.IVA), "", 1, "R", false, 0, "")

	// Total con IVA
	totalConIVA := despacho.TotalPrecio + despacho.IVA

	// Valor del Despacho
	ValorDespacho := despacho.ValorDespacho
	pdf.SetX(rectX)
	pdf.CellFormat(45, rowHeight, "Despacho:", "", 0, "L", false, 0, "")
	pdf.CellFormat(20, rowHeight, fmt.Sprintf("$%.0f", ValorDespacho), "", 1, "R", false, 0, "")

	// Total final
	TotalFinal := totalConIVA + ValorDespacho
	pdf.SetX(rectX)
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(45, rowHeight, "Total:", "", 0, "L", false, 0, "")
	pdf.CellFormat(20, rowHeight, fmt.Sprintf("$%.0f", TotalFinal), "", 1, "R", false, 0, "")

	// 9. Timbre electrónico
	timbreY := pdf.GetY() - (pdf.GetY() - startTotalsY)
	timbreWidth := 90.0
	var timbreHeight float64

	if _, err := os.Stat(config.TimbrePath); os.IsNotExist(err) {
		log.Printf("ADVERTENCIA: Timbre electrónico no encontrado en %s.", config.TimbrePath)
	} else {
		info := pdf.RegisterImage(config.TimbrePath, "")
		if info != nil {
			timbreHeight = timbreWidth * info.Height() / info.Width()
			pdf.Image(config.TimbrePath, 15, timbreY, timbreWidth, timbreHeight, false, "", 0, "")

			// Texto bajo el timbre
			pdf.SetFont("Arial", "", 9)
			pdf.SetTextColor(0, 0, 0)
			pdf.SetXY(25, timbreY+timbreHeight+1)
			pdf.Cell(timbreWidth, 4, tr("Res.99 de 2014 Verifique documento: www.sii.cl"))
		} else {
			log.Printf("ADVERTENCIA: No se pudo procesar la imagen del timbre en %s.", config.TimbrePath)
		}
	}

	// 10. Mensaje final (justo antes del pie)
	pdf.SetFont("Arial", "", 10)
	mensaje := "Gracias por confiar en nosotros"
	mensajeWidth := pdf.GetStringWidth(mensaje)
	pageWidth, pageHeight := pdf.GetPageSize()

	pdf.SetY(pageHeight - 40) // 5 mm arriba del pie de página
	pdf.SetX((pageWidth - mensajeWidth) / 2)
	pdf.CellFormat(mensajeWidth, 5, mensaje, "", 1, "C", false, 0, "")

	// 11. Pie de página
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
