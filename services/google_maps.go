package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

// GoogleMapsService maneja las interacciones con la API de Google Maps
type GoogleMapsService struct {
	APIKey  string
	BaseURL string
	Client  *http.Client
}

// DistanceMatrixResponse representa la respuesta de la API de Distance Matrix
type DistanceMatrixResponse struct {
	Rows []struct {
		Elements []struct {
			Distance struct {
				Text  string `json:"text"`
				Value int    `json:"value"`
			} `json:"distance"`
			Duration struct {
				Text  string `json:"text"`
				Value int    `json:"value"`
			} `json:"duration"`
			Status string `json:"status"`
		} `json:"elements"`
	} `json:"rows"`
	Status string `json:"status"`
}

// DistanciaCalculada representa el resultado del cálculo de distancia
type DistanciaCalculada struct {
	Distancia      string `json:"distancia"`
	Duracion       string `json:"duracion"`
	RutaOptimizada bool   `json:"ruta_optimizada"`
}

// NewGoogleMapsService crea una nueva instancia del servicio
func NewGoogleMapsService() *GoogleMapsService {
	apiKey := os.Getenv("GOOGLE_MAPS_API_KEY")
	if apiKey == "" {
		apiKey = "YOUR_API_KEY_HERE" // Valor por defecto para desarrollo
	}

	baseURL := os.Getenv("GOOGLE_MAPS_DISTANCE_API_URL")
	if baseURL == "" {
		baseURL = "https://maps.googleapis.com/maps/api/distancematrix/json"
	}

	return &GoogleMapsService{
		APIKey:  apiKey,
		BaseURL: baseURL,
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// CalcularDistancia calcula la distancia y duración entre dos direcciones
func (g *GoogleMapsService) CalcularDistancia(origen, destino string) (*DistanciaCalculada, error) {
	// Construir la URL con los parámetros
	params := url.Values{}
	params.Add("origins", origen)
	params.Add("destinations", destino)
	params.Add("units", "metric")
	params.Add("language", "es")
	params.Add("key", g.APIKey)

	requestURL := fmt.Sprintf("%s?%s", g.BaseURL, params.Encode())

	// Realizar la petición HTTP
	resp, err := g.Client.Get(requestURL)
	if err != nil {
		return nil, fmt.Errorf("error al realizar petición a Google Maps: %w", err)
	}
	defer resp.Body.Close()

	// Leer la respuesta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error al leer respuesta de Google Maps: %w", err)
	}

	// Parsear la respuesta JSON
	var distanceResponse DistanceMatrixResponse
	if err := json.Unmarshal(body, &distanceResponse); err != nil {
		return nil, fmt.Errorf("error al parsear respuesta de Google Maps: %w", err)
	}

	// Validar la respuesta
	if distanceResponse.Status != "OK" {
		return nil, fmt.Errorf("error en API de Google Maps: %s", distanceResponse.Status)
	}

	if len(distanceResponse.Rows) == 0 || len(distanceResponse.Rows[0].Elements) == 0 {
		return nil, fmt.Errorf("no se encontraron resultados de distancia")
	}

	element := distanceResponse.Rows[0].Elements[0]
	if element.Status != "OK" {
		return nil, fmt.Errorf("error al calcular distancia: %s", element.Status)
	}

	// Crear el resultado
	resultado := &DistanciaCalculada{
		Distancia:      element.Distance.Text,
		Duracion:       element.Duration.Text,
		RutaOptimizada: true, // Asumimos que Google Maps siempre devuelve la ruta optimizada
	}

	return resultado, nil
}

// ValidarDireccion valida si una dirección es válida usando la API de Geocoding
func (g *GoogleMapsService) ValidarDireccion(direccion string) (bool, error) {
	// Esta función podría implementarse usando la API de Geocoding de Google
	// Por ahora, solo validamos que no esté vacía
	if direccion == "" {
		return false, fmt.Errorf("la dirección no puede estar vacía")
	}
	return true, nil
}

// FormatearDireccionCompleta combina dirección, comuna y ciudad en un formato optimizado para Google Maps
func FormatearDireccionCompleta(direccion, comuna, ciudad string) string {
	if comuna != "" && ciudad != "" {
		return fmt.Sprintf("%s, %s, %s", direccion, comuna, ciudad)
	} else if ciudad != "" {
		return fmt.Sprintf("%s, %s", direccion, ciudad)
	}
	return direccion
}
