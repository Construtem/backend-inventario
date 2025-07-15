package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// Estructura para el request de cálculo de distancia
type DistanciaRequest struct {
	Origen  string `json:"origen"`
	Destino string `json:"destino"`
}

// Estructura para la respuesta
type DistanciaResponse struct {
	DistanciaKm     float64 `json:"distancia_km"`
	DuracionMinutos int     `json:"duracion_minutos"`
	Origen          string  `json:"origen"`
	Destino         string  `json:"destino"`
	Status          string  `json:"status"`
}

// Estructura para despachos con distancia
type DespachoConDistancia struct {
	ID              uint    `json:"id"`
	CotizacionID    uint    `json:"cotizacion_id"`
	CamionID        uint    `json:"camion_id"`
	Origen          uint    `json:"origen"`
	Destino         uint    `json:"destino"`
	FechaDespacho   string  `json:"fecha_despacho"`
	ValorDespacho   float64 `json:"valor_despacho"`
	Estado          string  `json:"estado"`
	CantidadItems   int     `json:"cantidad_items"`
	TotalKg         float64 `json:"total_kg"`
	TotalPrecio     float64 `json:"total_precio"`
	DistanciaKm     float64 `json:"distancia_km"`
	DuracionMinutos int     `json:"duracion_minutos"`
}

func main() {
	baseURL := "http://localhost:8080/api"
	
	fmt.Println("🚀 Probando funcionalidades de cálculo de distancias")
	fmt.Println("=" * 60)
	
	// Test 1: Verificar que el servidor esté funcionando
	fmt.Println("\n1. 🔍 Verificando conexión al servidor...")
	if !verificarServidor(baseURL) {
		log.Fatal("❌ El servidor no está disponible. Asegúrate de que esté corriendo en puerto 8080")
	}
	fmt.Println("✅ Servidor disponible")
	
	// Test 2: Probar cálculo de distancia simple
	fmt.Println("\n2. 📍 Probando cálculo de distancia básico...")
	probarCalculoDistancia(baseURL)
	
	// Test 3: Obtener despachos con distancias
	fmt.Println("\n3. 🚛 Obteniendo despachos con distancias calculadas...")
	obtenerDespachosConDistancia(baseURL)
	
	// Test 4: Probar con direcciones chilenas específicas
	fmt.Println("\n4. 🇨🇱 Probando con direcciones chilenas...")
	probarDistanciasChilenas(baseURL)
	
	fmt.Println("\n🎉 Pruebas completadas!")
}

func verificarServidor(baseURL string) bool {
	resp, err := http.Get(baseURL + "/despachos")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == 200
}

func probarCalculoDistancia(baseURL string) {
	// Direcciones de ejemplo en Santiago
	request := DistanciaRequest{
		Origen:  "Santiago Centro, Santiago, Chile",
		Destino: "Las Condes, Santiago, Chile",
	}
	
	jsonData, err := json.Marshal(request)
	if err != nil {
		log.Printf("❌ Error al serializar JSON: %v", err)
		return
	}
	
	resp, err := http.Post(baseURL+"/calcular-distancia", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("❌ Error en petición HTTP: %v", err)
		return
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("❌ Error al leer respuesta: %v", err)
		return
	}
	
	if resp.StatusCode != 200 {
		fmt.Printf("❌ Error %d: %s\n", resp.StatusCode, string(body))
		
		// Mostrar mensaje de ayuda si es error de API key
		if resp.StatusCode == 500 {
			fmt.Println("\n💡 Posibles soluciones:")
			fmt.Println("   1. Verifica que GOOGLE_MAPS_API_KEY esté configurada en tu .env")
			fmt.Println("   2. Asegúrate de que la API key tenga acceso a Distance Matrix API")
			fmt.Println("   3. Verifica que la API key no esté restringida")
		}
		return
	}
	
	var response DistanciaResponse
	if err := json.Unmarshal(body, &response); err != nil {
		log.Printf("❌ Error al parsear respuesta: %v", err)
		return
	}
	
	fmt.Printf("✅ Distancia calculada: %.2f km\n", response.DistanciaKm)
	fmt.Printf("✅ Duración estimada: %d minutos\n", response.DuracionMinutos)
}

func obtenerDespachosConDistancia(baseURL string) {
	resp, err := http.Get(baseURL + "/despachos-distancia")
	if err != nil {
		log.Printf("❌ Error en petición: %v", err)
		return
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("❌ Error al leer respuesta: %v", err)
		return
	}
	
	if resp.StatusCode != 200 {
		fmt.Printf("❌ Error %d: %s\n", resp.StatusCode, string(body))
		return
	}
	
	var despachos []DespachoConDistancia
	if err := json.Unmarshal(body, &despachos); err != nil {
		log.Printf("❌ Error al parsear despachos: %v", err)
		return
	}
	
	fmt.Printf("✅ Se obtuvieron %d despachos\n", len(despachos))
	
	for i, despacho := range despachos {
		if i >= 3 { // Mostrar solo los primeros 3
			fmt.Printf("   ... y %d más\n", len(despachos)-3)
			break
		}
		
		fmt.Printf("   📦 Despacho #%d:\n", despacho.ID)
		fmt.Printf("      - Estado: %s\n", despacho.Estado)
		fmt.Printf("      - Items: %d\n", despacho.CantidadItems)
		fmt.Printf("      - Peso: %.1f kg\n", despacho.TotalKg)
		
		if despacho.DistanciaKm > 0 {
			fmt.Printf("      - Distancia: %.2f km\n", despacho.DistanciaKm)
			fmt.Printf("      - Duración: %d min\n", despacho.DuracionMinutos)
		} else {
			fmt.Printf("      - Distancia: No calculada (sin API key)\n")
		}
		fmt.Println()
	}
}

func probarDistanciasChilenas(baseURL string) {
	// Direcciones comunes en Chile para testing
	rutas := []DistanciaRequest{
		{
			Origen:  "Providencia, Santiago, Chile",
			Destino: "Ñuñoa, Santiago, Chile",
		},
		{
			Origen:  "Valparaíso, Chile",
			Destino: "Santiago, Chile",
		},
		{
			Origen:  "Av. Libertador Bernardo O'Higgins 1449, Santiago, Chile",
			Destino: "Av. Apoquindo 3000, Las Condes, Santiago, Chile",
		},
	}
	
	for i, ruta := range rutas {
		fmt.Printf("   Ruta %d: %s → %s\n", i+1, ruta.Origen, ruta.Destino)
		
		jsonData, _ := json.Marshal(ruta)
		resp, err := http.Post(baseURL+"/calcular-distancia", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Printf("      ❌ Error: %v\n", err)
			continue
		}
		defer resp.Body.Close()
		
		body, _ := io.ReadAll(resp.Body)
		
		if resp.StatusCode == 200 {
			var response DistanciaResponse
			json.Unmarshal(body, &response)
			fmt.Printf("      ✅ %.2f km, %d minutos\n", response.DistanciaKm, response.DuracionMinutos)
		} else {
			fmt.Printf("      ❌ Error %d\n", resp.StatusCode)
		}
		
		// Pequeña pausa para no sobrecargar la API
		time.Sleep(200 * time.Millisecond)
	}
}

func repeat(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}
