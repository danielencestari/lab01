package services

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWeatherService_GetTemperatureByCity_Success(t *testing.T) {
	// Mock da resposta da WeatherAPI
	mockResponse := `{
		"location": {
			"name": "São Paulo",
			"region": "Sao Paulo",
			"country": "Brazil"
		},
		"current": {
			"temp_c": 25.0,
			"temp_f": 77.0
		}
	}`

	// Cria resposta HTTP mock
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(mockResponse)),
		Header:     make(http.Header),
	}
	resp.Header.Set("Content-Type", "application/json")

	// Configura mock do HTTP client
	mockClient := new(MockHTTPClient)
	expectedURL := "https://api.weatherapi.com/v1/current.json?key=test-api-key&q=S%C3%A3o+Paulo%2C+SP%2C+Brazil&aqi=no"
	mockClient.On("Get", expectedURL).Return(resp, nil)

	// Cria service com HTTP client mockado
	service := NewWeatherServiceWithClient(mockClient, "test-api-key")

	// Executa o teste
	result, err := service.GetTemperatureByCity("São Paulo", "SP")

	// Assertions
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 25.0, result.TempC)
	assert.Equal(t, 77.0, result.TempF)
	assert.Equal(t, 298.0, result.TempK) // 25 + 273

	// Verifica se o mock foi chamado corretamente
	mockClient.AssertExpectations(t)
}

func TestWeatherService_GetTemperatureByCity_NoAPIKey(t *testing.T) {
	// Cria service sem API key
	service := NewWeatherServiceWithClient(nil, "")

	// Executa o teste
	result, err := service.GetTemperatureByCity("São Paulo", "SP")

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "weather API key not configured")
}

func TestWeatherService_GetTemperatureByCity_APIError(t *testing.T) {
	// Cria resposta HTTP mock de erro 401
	resp := &http.Response{
		StatusCode: http.StatusUnauthorized,
		Body:       io.NopCloser(strings.NewReader(`{"error": {"code": 1002, "message": "API key not provided."}}`)),
		Header:     make(http.Header),
	}

	// Configura mock do HTTP client
	mockClient := new(MockHTTPClient)
	expectedURL := "https://api.weatherapi.com/v1/current.json?key=invalid-key&q=S%C3%A3o+Paulo%2C+SP%2C+Brazil&aqi=no"
	mockClient.On("Get", expectedURL).Return(resp, nil)

	// Cria service com HTTP client mockado
	service := NewWeatherServiceWithClient(mockClient, "invalid-key")

	// Executa o teste
	result, err := service.GetTemperatureByCity("São Paulo", "SP")

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "error fetching weather data: status 401")

	// Verifica se o mock foi chamado corretamente
	mockClient.AssertExpectations(t)
}

func TestWeatherService_GetTemperatureByCity_InvalidJSON(t *testing.T) {
	// Mock da resposta com JSON inválido
	mockResponse := `{invalid json}`

	// Cria resposta HTTP mock
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(mockResponse)),
		Header:     make(http.Header),
	}
	resp.Header.Set("Content-Type", "application/json")

	// Configura mock do HTTP client
	mockClient := new(MockHTTPClient)
	expectedURL := "https://api.weatherapi.com/v1/current.json?key=test-api-key&q=S%C3%A3o+Paulo%2C+SP%2C+Brazil&aqi=no"
	mockClient.On("Get", expectedURL).Return(resp, nil)

	// Cria service com HTTP client mockado
	service := NewWeatherServiceWithClient(mockClient, "test-api-key")

	// Executa o teste
	result, err := service.GetTemperatureByCity("São Paulo", "SP")

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "error decoding weather response")

	// Verifica se o mock foi chamado corretamente
	mockClient.AssertExpectations(t)
}

func TestWeatherService_GetTemperatureByCity_HTTPError(t *testing.T) {
	// Configura mock do HTTP client para retornar erro de conexão
	mockClient := new(MockHTTPClient)
	expectedURL := "https://api.weatherapi.com/v1/current.json?key=test-api-key&q=S%C3%A3o+Paulo%2C+SP%2C+Brazil&aqi=no"
	mockClient.On("Get", expectedURL).Return(nil, errors.New("connection error"))

	// Cria service com HTTP client mockado
	service := NewWeatherServiceWithClient(mockClient, "test-api-key")

	// Executa o teste
	result, err := service.GetTemperatureByCity("São Paulo", "SP")

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "error fetching weather data")

	// Verifica se o mock foi chamado corretamente
	mockClient.AssertExpectations(t)
}

func TestWeatherService_GetTemperatureByCity_DifferentTemperatures(t *testing.T) {
	tests := []struct {
		name      string
		tempC     float64
		expectedF float64
		expectedK float64
	}{
		{
			name:      "Zero Celsius",
			tempC:     0.0,
			expectedF: 32.0,
			expectedK: 273.0,
		},
		{
			name:      "Negative temperature",
			tempC:     -10.0,
			expectedF: 14.0,
			expectedK: 263.0,
		},
		{
			name:      "High temperature",
			tempC:     40.0,
			expectedF: 104.0,
			expectedK: 313.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock da resposta da WeatherAPI
			mockResponse := `{
				"location": {
					"name": "São Paulo",
					"region": "Sao Paulo",
					"country": "Brazil"
				},
				"current": {
					"temp_c": ` + fmt.Sprintf("%.1f", tt.tempC) + `,
					"temp_f": 77.0
				}
			}`

			// Cria resposta HTTP mock
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(mockResponse)),
				Header:     make(http.Header),
			}
			resp.Header.Set("Content-Type", "application/json")

			// Configura mock do HTTP client
			mockClient := new(MockHTTPClient)
			expectedURL := "https://api.weatherapi.com/v1/current.json?key=test-api-key&q=S%C3%A3o+Paulo%2C+SP%2C+Brazil&aqi=no"
			mockClient.On("Get", expectedURL).Return(resp, nil)

			// Cria service com HTTP client mockado
			service := NewWeatherServiceWithClient(mockClient, "test-api-key")

			// Executa o teste
			result, err := service.GetTemperatureByCity("São Paulo", "SP")

			// Assertions
			require.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, tt.tempC, result.TempC)
			assert.Equal(t, tt.expectedF, result.TempF)
			assert.Equal(t, tt.expectedK, result.TempK)

			// Verifica se o mock foi chamado corretamente
			mockClient.AssertExpectations(t)
		})
	}
} 