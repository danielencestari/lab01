package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"weather-cep-api/models"
	"weather-cep-api/utils"
)

// WeatherServiceInterface define o contrato para serviços de clima
type WeatherServiceInterface interface {
	GetTemperatureByCity(city, state string) (*models.TemperatureResponse, error)
}

// WeatherService implementa o serviço de consulta de clima
type WeatherService struct {
	httpClient HTTPClientInterface
	apiKey     string
}

// NewWeatherService cria uma nova instância do serviço de clima
func NewWeatherService() *WeatherService {
	return &WeatherService{
		httpClient: &http.Client{},
		apiKey:     os.Getenv("WEATHER_API_KEY"),
	}
}

// NewWeatherServiceWithClient cria uma nova instância com HTTP client customizado (para testes)
func NewWeatherServiceWithClient(client HTTPClientInterface, apiKey string) *WeatherService {
	return &WeatherService{
		httpClient: client,
		apiKey:     apiKey,
	}
}

// GetTemperatureByCity consulta a temperatura atual de uma cidade usando WeatherAPI
func (s *WeatherService) GetTemperatureByCity(city, state string) (*models.TemperatureResponse, error) {
	// Verifica se a API key está configurada
	if s.apiKey == "" {
		return nil, fmt.Errorf("weather API key not configured")
	}

	// Constrói a query de localização (cidade, estado, Brasil)
	location := fmt.Sprintf("%s, %s, Brazil", city, state)
	encodedLocation := url.QueryEscape(location)

	// Constrói URL da WeatherAPI
	weatherURL := fmt.Sprintf("https://api.weatherapi.com/v1/current.json?key=%s&q=%s&aqi=no", 
		s.apiKey, encodedLocation)

	// Faz a requisição HTTP
	resp, err := s.httpClient.Get(weatherURL)
	if err != nil {
		return nil, fmt.Errorf("error fetching weather data: %w", err)
	}
	defer resp.Body.Close()

	// Verifica status da resposta
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error fetching weather data: status %d", resp.StatusCode)
	}

	// Decodifica a resposta JSON
	var weatherResp models.WeatherAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&weatherResp); err != nil {
		return nil, fmt.Errorf("error decoding weather response: %w", err)
	}

	// Converte temperaturas para todas as escalas
	tempC, tempF, tempK := utils.ConvertTemperatures(weatherResp.Current.TempC)

	// Cria resposta final
	response := &models.TemperatureResponse{
		TempC: tempC,
		TempF: tempF,
		TempK: tempK,
	}

	return response, nil
} 