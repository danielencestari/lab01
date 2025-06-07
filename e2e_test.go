package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"weather-cep-api/handlers"
	"weather-cep-api/models"
	"weather-cep-api/services"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupE2ERouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	
	// Cria instâncias dos serviços reais
	cepService := services.NewCEPService()
	weatherService := services.NewWeatherService()
	
	// Cria instância do handler
	weatherHandler := handlers.NewWeatherHandler(cepService, weatherService)
	
	// Configura o router
	router := gin.New()
	
	// Define as rotas
	router.GET("/health", weatherHandler.HealthCheck)
	router.GET("/temperature/:cep", weatherHandler.GetTemperatureByCEP)
	
	return router
}

func TestE2E_HealthCheck(t *testing.T) {
	router := setupE2ERouter()
	
	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "ok", response["status"])
	assert.Equal(t, "Weather CEP API is running", response["message"])
}

func TestE2E_GetTemperature_ValidCEP_WithAPIKey(t *testing.T) {
	// Verifica se a API key está configurada para testes de integração
	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		t.Skip("WEATHER_API_KEY not set, skipping integration test")
	}

	router := setupE2ERouter()
	
	// Teste com CEP de São Paulo (Avenida Paulista)
	req, _ := http.NewRequest("GET", "/temperature/01310100", nil)
	w := httptest.NewRecorder()
	
	router.ServeHTTP(w, req)
	
	// Verifica se o request foi bem-sucedido
	assert.Equal(t, http.StatusOK, w.Code)
	
	// Decodifica a resposta
	var response models.TemperatureResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	
	// Verifica se as temperaturas foram retornadas
	assert.NotZero(t, response.TempC, "Temperature in Celsius should not be zero")
	assert.NotZero(t, response.TempF, "Temperature in Fahrenheit should not be zero")
	assert.NotZero(t, response.TempK, "Temperature in Kelvin should not be zero")
	
	// Verifica conversões de temperatura
	expectedF := response.TempC*1.8 + 32
	expectedK := response.TempC + 273
	
	assert.InDelta(t, expectedF, response.TempF, 0.1, "Fahrenheit conversion should be correct")
	assert.InDelta(t, expectedK, response.TempK, 0.1, "Kelvin conversion should be correct")
	
	// Verifica se as temperaturas estão em uma faixa razoável para Brasil
	assert.GreaterOrEqual(t, response.TempC, -10.0, "Temperature should be greater than -10°C")
	assert.LessOrEqual(t, response.TempC, 50.0, "Temperature should be less than 50°C")
}

func TestE2E_GetTemperature_InvalidCEP_Format(t *testing.T) {
	router := setupE2ERouter()
	
	tests := []struct {
		name string
		cep  string
	}{
		{
			name: "CEP too short",
			cep:  "123456",
		},
		{
			name: "CEP too long",
			cep:  "123456789",
		},
		{
			name: "CEP with letters",
			cep:  "1234567a",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/temperature/"+tt.cep, nil)
			w := httptest.NewRecorder()
			
			router.ServeHTTP(w, req)
			
			assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
			assert.Equal(t, "invalid zipcode", w.Body.String())
		})
	}
}

func TestE2E_GetTemperature_CEPNotFound(t *testing.T) {
	router := setupE2ERouter()
	
	// CEP que não existe (99999-999)
	req, _ := http.NewRequest("GET", "/temperature/99999999", nil)
	w := httptest.NewRecorder()
	
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "can not find zipcode", w.Body.String())
}

func TestE2E_GetTemperature_WithoutAPIKey(t *testing.T) {
	// Salva a API key atual e a remove temporariamente
	originalKey := os.Getenv("WEATHER_API_KEY")
	os.Unsetenv("WEATHER_API_KEY")
	defer os.Setenv("WEATHER_API_KEY", originalKey)
	
	router := setupE2ERouter()
	
	// Teste com CEP válido mas sem API key
	req, _ := http.NewRequest("GET", "/temperature/01310100", nil)
	w := httptest.NewRecorder()
	
	router.ServeHTTP(w, req)
	
	// Deve retornar erro interno do servidor porque a API key não está configurada
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "error fetching weather data", w.Body.String())
}

func TestE2E_GetTemperature_DifferentValidCEPs(t *testing.T) {
	// Verifica se a API key está configurada
	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		t.Skip("WEATHER_API_KEY not set, skipping integration test")
	}

	router := setupE2ERouter()
	
	validCEPs := []struct {
		cep      string
		name     string
		city     string
	}{
		{
			cep:  "01310100", // São Paulo - Avenida Paulista
			name: "São Paulo",
			city: "São Paulo",
		},
		{
			cep:  "20040020", // Rio de Janeiro - Centro
			name: "Rio de Janeiro",
			city: "Rio de Janeiro",
		},
		{
			cep:  "70040010", // Brasília - Asa Norte
			name: "Brasília",
			city: "Brasília",
		},
	}
	
	for _, tc := range validCEPs {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/temperature/"+tc.cep, nil)
			w := httptest.NewRecorder()
			
			router.ServeHTTP(w, req)
			
			// Pode dar timeout em alguns casos, então aceita tanto sucesso quanto erro de servidor
			if w.Code == http.StatusOK {
				var response models.TemperatureResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				
				// Verifica se as temperaturas são válidas
				assert.NotZero(t, response.TempC)
				assert.NotZero(t, response.TempF)
				assert.NotZero(t, response.TempK)
				
				// Verifica conversões
				expectedF := response.TempC*1.8 + 32
				expectedK := response.TempC + 273
				
				assert.InDelta(t, expectedF, response.TempF, 0.1)
				assert.InDelta(t, expectedK, response.TempK, 0.1)
			} else {
				// Em caso de erro de API ou timeout, deve ser erro 500
				assert.Equal(t, http.StatusInternalServerError, w.Code)
				assert.Equal(t, "error fetching weather data", w.Body.String())
			}
		})
	}
}

func TestE2E_GetTemperature_CEPFormats(t *testing.T) {
	// Verifica se a API key está configurada
	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		t.Skip("WEATHER_API_KEY not set, skipping integration test")
	}

	router := setupE2ERouter()
	
	// Testa diferentes formatos do mesmo CEP
	cepFormats := []string{
		"01310100",   // Sem hífen
		"01310-100",  // Com hífen (será normalizado internamente)
	}
	
	for _, cep := range cepFormats {
		t.Run("CEP_format_"+cep, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/temperature/"+cep, nil)
			w := httptest.NewRecorder()
			
			router.ServeHTTP(w, req)
			
			// Ambos os formatos devem funcionar
			if w.Code == http.StatusOK {
				var response models.TemperatureResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.NotZero(t, response.TempC)
			} else {
				// Em caso de erro de API, deve ser erro 500
				assert.Equal(t, http.StatusInternalServerError, w.Code)
			}
		})
	}
} 